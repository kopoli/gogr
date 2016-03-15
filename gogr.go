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
	version      = "Undefined"
	progVersion  = majorVersion + "-" + version
)

func fault(err error, message string, arg ...string) {
	msg := ""
	if err != nil {
		msg = fmt.Sprintf(" (error: %s)", err)
	}
	fmt.Fprintf(os.Stderr, "Error: %s%s.%s\n", message, strings.Join(arg, " "), msg)
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
	tagman.Add(tag, dirs...)
	err := tagman.Save()
	if err != nil {
		fault(err, "Saving configuration failed.")
	}
}

func rmTag(tagman gogr.TagManager, tag string, dirs []string) {
	tagman.Remove(tag, dirs...)
	err := tagman.Save()
	if err != nil {
		fault(err, "Saving configuration failed.")
	}
}

func main() {
	var tagman gogr.TagManager
	opts := gogr.GetOptions()

	opts.Set("application-name", progName)
	opts.Set("configuration-file", "config.json")

	app := cli.App(progName, "Run commands in multiple directories")
	app.Spec = "[OPTIONS] [-- ARG...]"

	app.Version("version v", fmt.Sprintf("%s: %s\ngo: %s\nBuilt with: %s on %s/%s\n", progName, progVersion,
		runtime.Version(), runtime.Compiler, runtime.GOOS, runtime.GOARCH))
	optVerbose := app.BoolOpt("verbose V", false, "Print verbose output")
	optConfig := app.StringOpt("config c", gogr.DefaultConfigFile(opts), "Configuration file")
	optConcurrent := app.BoolOpt("concurrent j", false, "Run the commands concurrently")

	app.Before = func() {
		opts.Set("configuration-file", *optConfig)
		tagman = gogr.NewTagManager(opts)
		opts.Set("flag-verbose", strconv.FormatBool(*optVerbose))
		opts.Set("concurrent", strconv.FormatBool(*optConcurrent))
	}

	app.Command("tag", "Tag management", func(cmd *cli.Cmd) {

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

		cmd.Command("list", "List all tags or directories of given tag. This is the default action.", func(cmd *cli.Cmd) {
			cmd.Spec = "[TAG]"

			tagArg := cmd.StringArg("TAG", "", "List directories of this tag")

			cmd.Action = func() {
				if *tagArg == "" {
					for tag := range tagman.Tags {
						fmt.Println(tag)
					}
				} else {
					dirs := tagman.Dirs([]string{*tagArg}, nil)
					for _, dir := range dirs {
						fmt.Println(dir)
					}
				}
			}
		})

		cmd.Action = func() {
			for tag := range tagman.Tags {
				fmt.Println(tag)
			}
		}
	})

	argArg := app.StringsArg("ARG", nil, "Directories and command to be run")

	app.Action = func() {
		tagitems := gogr.ParseTags(*argArg)
		cmd, tags, dirs, args, err := gogr.VerifyTags(tagitems)
		if err != nil {
			fault(err, "Parsing arguments failed")
			return
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

		dirs = tagman.Dirs(tags, dirs)

		err = gogr.RunCommands(opts, dirs, args)
		if err != nil {
			fault(err, "Running command failed")
		}
	}

	app.Run(os.Args)
}
