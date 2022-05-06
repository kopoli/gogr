package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/kopoli/appkit"
	gogr "github.com/kopoli/gogr/lib"
)

var (
	version     = "Undefined"
	timestamp   = "Undefined"
	buildGOOS   = "Undefined"
	buildGOARCH = "Undefined"
	progVersion = "" + version
)

var tagman gogr.TagManager
var opts appkit.Options

func fault(err error, message string, arg ...string) {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %s%s: %s\n", message, strings.Join(arg, " "), err)
		os.Exit(1)
	}
}

func addTag(tagman gogr.TagManager, tag string, dirs []string) {
	if len(dirs) == 0 {
		wd, err := os.Getwd()
		fault(err, "Getting working directory failed")
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
		fault(fmt.Errorf("parsing failed"), "Improper tag string")
	}
	tagman.Remove(tag, dirs...)
	err := tagman.Save()
	fault(err, "Saving configuration failed.")
}

func main() {
	opts = appkit.NewOptions()

	opts.Set("program-name", os.Args[0])
	opts.Set("program-version", progVersion)
	opts.Set("program-timestamp", timestamp)
	opts.Set("program-buildgoos", buildGOOS)
	opts.Set("program-buildgoarch", buildGOARCH)
	opts.Set("configuration-file", "config.json")

	base := appkit.NewCommand(nil, "", "Run commands in multiple directories")
	base.ArgumentHelp = "| [+-]@<tag> [CMD ...]"
	optVersion := base.Flags.Bool("version", false, "Display version")
	base.Flags.BoolVar(optVersion, "v", false, "Display version")
	optVerbose := base.Flags.Bool("verbose", false, "Print verbose output")
	base.Flags.BoolVar(optVerbose, "V", false, "Print verbose output")

	optConfig := base.Flags.String("config", gogr.DefaultConfigFile(opts), "Configuration file")
	base.Flags.StringVar(optConfig, "c", gogr.DefaultConfigFile(opts), "Configuration file")
	optConcurrent := base.Flags.Bool("concurrent", false, "Run the commands concurrently")
	base.Flags.BoolVar(optConcurrent, "j", false, "Run the commands concurrently")
	optLicenses := base.Flags.Bool("licenses", false, "Display the licenses")

	tag := appkit.NewCommand(base, "tag", "Tag management")
	tag.SubCommandHelp = "<COMMAND>"
	tlist := appkit.NewCommand(tag, "list", "List all tags or directories of given tag. This is the default action.")
	tlist.ArgumentHelp = "[TAG ...]"
	tadd := appkit.NewCommand(tag, "add", "Add tag to path")
	tadd.ArgumentHelp = "TAG [DIR ...]"
	trm := appkit.NewCommand(tag, "rm", "Remove tag from paths")
	trm.ArgumentHelp = "TAG [DIR ...]"

	discover := appkit.NewCommand(base, "discover", "Discover directories containing a certain file")
	discover.ArgumentHelp = "TAG [ROOT ...]"
	optDepth := discover.Flags.Int("max-depth", 5, "Maximum depth of discovery")
	optFile := discover.Flags.String("File", ".git", "File or directory do discover")

	err := base.Parse(os.Args[1:], opts)
	if err == flag.ErrHelp {
		os.Exit(0)
	}

	errorShowHelp := func(message string) {
		fmt.Fprintf(os.Stderr, "Error: %s\n\n", message)
		base.Flags.Usage()
		os.Exit(1)
	}

	if *optVersion {
		fmt.Println(appkit.VersionString(opts))
		os.Exit(0)
	}
	opts.Set("configuration-file", *optConfig)
	if *optVerbose {
		opts.Set("flag-verbose", "t")
	}
	if *optConcurrent {
		opts.Set("concurrent", "t")
	}
	if *optLicenses {
		l, err := GetLicenses()
		fault(err, "Getting licenses failed")
		s, err := appkit.LicenseString(l)
		fault(err, "Interpreting licenses failed")
		fmt.Print(s)
		os.Exit(0)
	}

	opts.Set("discover-max-depth", strconv.Itoa(*optDepth))
	opts.Set("discover-file", *optFile)

	cmd := opts.Get("cmdline-command", "")
	argstr := opts.Get("cmdline-args", "")
	args := appkit.SplitArguments(argstr)

	if cmd == "" && argstr == "" {
		errorShowHelp("Arguments required")
	}

	tagman = gogr.NewTagManager(opts)

	parseDir := func(dirs []string) ([]string) {
		ret := []string{}
		if len(dirs) > 0 {
			_, rest, err := gogr.ParseDirectories(dirs)
			fault(err, "Directory parsing failed")
			if len(rest) != 0 {
				fault(fmt.Errorf("not directories: %v", rest), "Invalid arguments")
			}
			ret = dirs
		}
		return ret
	}

	parseTagDirArg := func(args []string) (tag string, dirs []string) {
		if len(args) < 1 {
			fault(fmt.Errorf("Not enough arguments"), "Command line parsing failed")
		}

		return args[0], args[1:]
	}

	switch cmd {
	case "tag":
		fallthrough
	case "tag list":
		if len(args) == 0 {
			for tag := range tagman.Tags {
				fmt.Println(tag)
			}
		} else {
			dirs := tagman.Dirs(args, nil)
			for _, dir := range dirs {
				fmt.Println(dir)
			}
		}
	case "tag add":
		tag, dirs := parseTagDirArg(args)
		addTag(tagman, tag, parseDir(dirs))
	case "tag rm":
		tag, dirs := parseTagDirArg(args)
		rmTag(tagman, tag, parseDir(dirs))
	case "discover":
		tag, roots := parseTagDirArg(args)
		if len(roots) == 0 {
			wd, err := os.Getwd()
			fault(err, "Could not get the current directory")
			roots = append(roots, wd)
		}
	
		var dirs []string
		for _, root := range roots {
			tmp, err := gogr.Discover(opts, root, *optFile)
			fault(err, "Discovery failed")
			dirs = append(dirs, tmp...)
		}
		for _, dir := range dirs {
			fmt.Println(dir)
		}

		rmTag(tagman, tag, []string{})
		addTag(tagman, tag, dirs)
	default:
		tagitems := gogr.ParseTags(args)
		vt, err := gogr.VerifyTags(tagitems)
		fault(err, "Parsing arguments failed")
		if len(vt.Dirs) == 0 && len(vt.Tags) == 0 {
			errorShowHelp("Directories or tags are required")
		}

		if vt.Command.Str != "" {
			switch vt.Command.Op {
			case gogr.Add:
				addTag(tagman, vt.Command.Str, vt.Dirs)
			case gogr.Remove:
				rmTag(tagman, vt.Command.Str, vt.Dirs)
			case gogr.None:
				fallthrough
			default:
				fault(fmt.Errorf("improper command"), "Internal error")
			}
			return
		}

		invalid := tagman.AreProper(vt.Tags)
		if len(invalid) > 0 {
			fault(fmt.Errorf("improper tags"),
				"Tags: ", strings.Join(invalid, ", "))
		}

		vt.Dirs = tagman.Dirs(vt.Tags, vt.Dirs)

		err = gogr.RunCommands(opts, vt.Dirs, vt.Args)
		fault(err, "Running command failed")
	}
}
