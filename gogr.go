package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/jawher/mow.cli"
	"github.com/kopoli/gogr/lib"
)

var (
	progName     = "gogr"
	majorVersion = "1"
	version      = "development"
	timestamp    = "undefined"
	progVersion  = majorVersion + "-" + version
)

var tagman gogr.TagManager
var opts gogr.Options

func printErr(err error, message string, arg ...string) {
	msg := ""
	if err != nil {
		msg = fmt.Sprintf(" (error: %s)", err)
	}
	fmt.Fprintf(os.Stderr, "Error: %s%s.%s\n", message, strings.Join(arg, " "), msg)
}

func fault(err error, message string, arg ...string) {
	printErr(err, message, arg...)
	cli.Exit(1)
}

func faultShowHelp(app *cli.Cli, message string, arg ...string) {
	printErr(nil, message, arg...)
	app.PrintLongHelp()
	cli.Exit(1)
}

func addTag(tagman gogr.TagManager, tag string, dirs []string) {
	if len(dirs) == 0 {
		wd, err := os.Getwd()
		if err != nil {
			fault(err, "Getting working directory failed")
		}
		dirs = append(dirs, wd)
	}
	if !tagman.ValidateTag(tag) {
		fault(nil, "Improper tag string")
	}
	tagman.Add(tag, dirs...)
	err := tagman.Save()
	if err != nil {
		fault(err, "Saving configuration failed.")
	}
}

func rmTag(tagman gogr.TagManager, tag string, dirs []string) {
	if !tagman.ValidateTag(tag) {
		fault(nil, "Improper tag string")
	}
	tagman.Remove(tag, dirs...)
	err := tagman.Save()
	if err != nil {
		fault(err, "Saving configuration failed.")
	}
}

func cmdTag(cmd *cli.Cmd) {
	parseDir := func(dirs []string) (ret []string) {
		if len(dirs) > 0 {
			_, rest, err := gogr.ParseDirectories(dirs)
			if err != nil {
				fault(err, "Directory parsing failed")
			}
			if len(rest) != 0 {
				fault(nil, "The following are not directories: ", rest...)
			}
			ret = dirs
		}
		return
	}

	cmd.Command("add", "Add tag to path", func(cmd *cli.Cmd) {
		cmd.Spec = "TAG [DIR...]"
		tagArg := cmd.StringArg("TAG", "", "Tag to add")
		dirArg := cmd.StringsArg("DIR", nil, "Directories to add the tag to")

		cmd.Action = func() {
			dirs := parseDir(*dirArg)
			addTag(tagman, *tagArg, dirs)
		}
	})
	cmd.Command("rm", "Remove tag from paths", func(cmd *cli.Cmd) {
		cmd.Spec = "TAG [DIR...]"
		tagArg := cmd.StringArg("TAG", "", "Tag to remove")
		dirArg := cmd.StringsArg("DIR", nil, "Directories to remove the tag from")

		cmd.Action = func() {
			rmTag(tagman, *tagArg, parseDir(*dirArg))
		}
	})

	listTags := func() {
		for tag := range tagman.Tags {
			fmt.Println(tag)
		}
	}

	cmd.Command("list", "List all tags or directories of given tag. This is the default action.", func(cmd *cli.Cmd) {
		cmd.Spec = "[TAG]"

		tagArg := cmd.StringArg("TAG", "", "List directories of this tag")

		cmd.Action = func() {
			if *tagArg == "" {
				listTags()
			} else {
				dirs := tagman.Dirs([]string{*tagArg}, nil)
				for _, dir := range dirs {
					fmt.Println(dir)
				}
			}
		}
	})

	cmd.Action = listTags
}

func cmdDiscover(cmd *cli.Cmd) {
	optDepth := cmd.IntOpt("max-depth d", 5, "Maximum depth of discovery")
	wd, err := os.Getwd()
	if err != nil {
		fault(err, "Could not get the current directory")
	}
	tagArg := cmd.StringArg("TAG", "", "Tag to add or update")
	optRoots := cmd.StringsArg("ROOT", []string{wd},
		"Root directory for the discovery")
	optFile := cmd.StringOpt("file f", ".git",
		"File or directory to discover")
	cmd.Spec = "[OPTIONS] TAG [ROOT...]"

	cmd.Action = func() {
		var dirs []string
		opts.Set("discover-max-depth", strconv.Itoa(*optDepth))
		for _, root := range *optRoots {
			tmp, err := gogr.Discover(opts, root, *optFile)
			if err != nil {
				fault(err, "Discovery failed")
			}

			dirs = append(dirs, tmp...)
		}
		for _, dir := range dirs {
			fmt.Println(dir)
		}

		rmTag(tagman, *tagArg, []string{})
		addTag(tagman, *tagArg, dirs)
	}
}

func main() {
	opts = gogr.GetOptions()

	opts.Set("application-name", progName)
	opts.Set("configuration-file", "config.json")

	app := cli.App(progName, "Run commands in multiple directories")
	app.Spec = "[OPTIONS] [-- ARG...]"

	app.Version("version v", fmt.Sprintf(
		"%s: %s\nBuilt with: %s/%s on %s/%s\nBuilt: %s",
		progName, progVersion, runtime.Compiler, runtime.Version(),
		runtime.GOOS, runtime.GOARCH, timestamp))

	optVerbose := app.BoolOpt("verbose V", false, "Print verbose output")
	optConfig := app.StringOpt("config c", gogr.DefaultConfigFile(opts), "Configuration file")
	optConcurrent := app.BoolOpt("concurrent j", false, "Run the commands concurrently")

	app.Before = func() {
		opts.Set("configuration-file", *optConfig)
		tagman = gogr.NewTagManager(opts)
		opts.Set("flag-verbose", strconv.FormatBool(*optVerbose))
		opts.Set("concurrent", strconv.FormatBool(*optConcurrent))
	}

	app.Command("tag", "Tag management", cmdTag)
	app.Command("discover", "Discover directories containing a certain file", cmdDiscover)

	argArg := app.StringsArg("ARG", nil, "Directories and command to be run")

	app.Action = func() {
		if len(*argArg) == 0 {
			faultShowHelp(app, "Arguments required")
		}
		tagitems := gogr.ParseTags(*argArg)
		cmd, tags, dirs, args, err := gogr.VerifyTags(tagitems)
		if err != nil {
			fault(err, "Parsing arguments failed")
		}
		if len(dirs) == 0 && len(tags) == 0 {
			faultShowHelp(app, "Directories or tags are required")
		}

		if cmd.Str != "" {
			switch cmd.Op {
			case gogr.Add:
				addTag(tagman, cmd.Str, dirs)
			case gogr.Remove:
				rmTag(tagman, cmd.Str, dirs)
			default:
				fault(nil, "Internal error: improper command")
			}
			return
		}

		invalid := tagman.AreProper(tags)
		if len(invalid) > 0 {
			fault(nil, "Improper tags: ",strings.Join(invalid,", "))
		}

		dirs = tagman.Dirs(tags, dirs)

		err = gogr.RunCommands(opts, dirs, args)
		if err != nil {
			fault(err, "Running command failed")
		}
	}

	err := app.Run(os.Args)
	if err != nil {
		fault(err, "Argument parsing failed")
	}
}
