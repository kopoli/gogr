package gogr

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/kopoli/appkit"
)

var stdout io.Writer = os.Stdout
var stderr io.Writer = os.Stderr

func wrapErr(err error, message string) error {
	if err != nil {
		return fmt.Errorf("%s: %v", message, err)
	}
	return nil
}

func addTag(tagman *TagManager, tag string, dirs []string) error {
	if len(dirs) == 0 {
		wd, err := os.Getwd()
		err = wrapErr(err, "getting working directory failed")
		if err != nil {
			return err
		}
		dirs = append(dirs, wd)
	}
	if !tagman.ValidateTag(tag) {
		return fmt.Errorf("improper tag found")
	}
	tagman.Add(tag, dirs...)
	err := tagman.Save()
	return wrapErr(err, "saving configuration failed")
}

func rmTag(tagman *TagManager, tag string, dirs []string) error {
	if !tagman.ValidateTag(tag) {
		return fmt.Errorf("parsing tag string failed")
	}
	tagman.Remove(tag, dirs...)
	err := tagman.Save()
	return wrapErr(err, "saving configuration failed")
}

func escapeTagArgs(args []string, unescape bool) []string {
	for i := range args {
		if !unescape {
			if strings.HasPrefix(args[i], "\\") ||
				strings.HasPrefix(args[i], "-@") {
				args[i] = "\\" + args[i]
			}
		} else {
			args[i] = strings.TrimPrefix(args[i], "\\")
		}
	}
	return args
}

var ErrHandled = fmt.Errorf("error already handled")
var ErrLicenses = fmt.Errorf("license display requested")

func Main(cmdLineArgs []string, opts appkit.Options) error {
	base := appkit.NewCommand(nil, "", "Run commands in multiple directories")
	base.Flags.SetOutput(stderr)
	base.ArgumentHelp = "| [+-]@<tag> [CMD ...]"
	optVersion := base.Flags.Bool("version", false, "Display version")
	base.Flags.BoolVar(optVersion, "v", false, "Display version")
	optVerbose := base.Flags.Bool("verbose", false, "Print verbose output")
	base.Flags.BoolVar(optVerbose, "V", false, "Print verbose output")

	optConfig := base.Flags.String("config", DefaultConfigFile(opts), "Configuration file")
	base.Flags.StringVar(optConfig, "c", DefaultConfigFile(opts), "Configuration file")
	optConcurrent := base.Flags.Bool("concurrent", false, "Run the commands concurrently")
	base.Flags.BoolVar(optConcurrent, "j", false, "Run the commands concurrently")
	optLicenses := base.Flags.Bool("licenses", false, "Display the licenses")

	tag := appkit.NewCommand(base, "tag", "Tag management")
	tag.Flags.SetOutput(stderr)
	tag.SubCommandHelp = "<COMMAND>"

	tlist := appkit.NewCommand(tag, "list l", "List all tags or directories of given tag. This is the default action.")
	tlist.Flags.SetOutput(stderr)
	tlist.ArgumentHelp = "[TAG ...]"
	optRelativeHelp := "Print out directories relative to current directory"
	optRelative := tlist.Flags.Bool("relative", false, optRelativeHelp)
	tlist.Flags.BoolVar(optRelative, "r", false, optRelativeHelp)

	tadd := appkit.NewCommand(tag, "add a", "Add tag to path")
	tadd.Flags.SetOutput(stderr)
	tadd.ArgumentHelp = "TAG [DIR ...]"
	tdel := appkit.NewCommand(tag, "delete del d", "Delete tag or paths from tag")
	tdel.Flags.SetOutput(stderr)
	tdel.ArgumentHelp = "TAG [DIR ...]"

	discover := appkit.NewCommand(base, "discover", "Discover directories containing a certain file")
	discover.Flags.SetOutput(stderr)
	discover.ArgumentHelp = "TAG [ROOT ...]"
	optDepthHelp := "Maximum depth of discovery"
	optDepth := discover.Flags.Int("max-depth", 5, optDepthHelp)
	discover.Flags.IntVar(optDepth, "d", 5, optDepthHelp)
	optFile := discover.Flags.String("file", ".git", "File or directory to discover")
	discover.Flags.StringVar(optFile, "f", ".git", "File or directory to discover")

	args := escapeTagArgs(cmdLineArgs[1:], false)

	err := base.Parse(args, opts)
	if err == flag.ErrHelp {
		return nil
	} else if err != nil {
		return err
	}

	errorShowHelp := func(message string) {
		fmt.Fprintf(stderr, "Error: %s\n\n", message)
		base.Flags.Usage()
	}

	if *optVersion {
		fmt.Fprintln(stdout, appkit.VersionString(opts))
		return nil
	}
	opts.Set("configuration-file", *optConfig)
	if *optVerbose {
		opts.Set("flag-verbose", "t")
	}
	if *optConcurrent {
		opts.Set("concurrent", "t")
	}
	if *optLicenses {
		return ErrLicenses
	}

	opts.Set("discover-max-depth", strconv.Itoa(*optDepth))
	opts.Set("discover-file", *optFile)

	if *optRelative {
		opts.Set("relative-paths", "t")
	}

	cmd := opts.Get("cmdline-command", "")
	argstr := opts.Get("cmdline-args", "")
	args = appkit.SplitArguments(argstr)
	args = escapeTagArgs(args, true)

	if cmd == "" && argstr == "" {
		errorShowHelp("Arguments required")
		return ErrHandled
	}

	tagman, err := NewTagManager(opts)
	if err != nil {
		return fmt.Errorf("loading tags failed: %v", err)
	}

	parseDir := func(dirs []string) ([]string, error) {
		ret := []string{}
		if len(dirs) > 0 {
			_, rest, err := ParseDirectories(dirs)
			err = wrapErr(err, "Directory parsing failed")
			if err != nil {
				return nil, err
			}
			if len(rest) != 0 {
				return nil, wrapErr(fmt.Errorf("not directories: %v", rest), "invalid arguments")
			}
			ret = dirs
		}
		return ret, nil
	}

	parseTagDirArg := func(args []string) (string, []string, error) {
		if len(args) < 1 {
			return "", nil, wrapErr(fmt.Errorf("not enough arguments"), "command line parsing failed")
		}

		return args[0], args[1:], nil
	}

	checkTags := func(tags []string) error {
		invalid := tagman.AreProper(tags)
		if len(invalid) > 0 {
			return fmt.Errorf("improper tags: %s", strings.Join(invalid, ", "))
		}

		return nil
	}

	switch cmd {
	case "tag":
		fallthrough
	case "tag list":
		if len(args) == 0 {
			for tag := range tagman.Tags {
				fmt.Fprintln(stdout, tag)
			}
		} else {
			err := checkTags(args)
			if err != nil {
				return err
			}
			dirs := tagman.Dirs(args, nil)
			if opts.IsSet("relative-paths") {
				dirs = ChangeToRelativePaths(dirs)
			}
			for _, dir := range dirs {
				fmt.Fprintln(stdout, dir)
			}
		}
	case "tag add":
		tag, dirs, err := parseTagDirArg(args)
		if err != nil {
			return err
		}
		dirs, err = parseDir(dirs)
		if err != nil {
			return err
		}
		err = addTag(tagman, tag, dirs)
		if err != nil {
			return err
		}
	case "tag delete":
		tag, dirs, err := parseTagDirArg(args)
		if err != nil {
			return err
		}
		dirs, err = parseDir(dirs)
		if err != nil {
			return err
		}
		err = rmTag(tagman, tag, dirs)
		if err != nil {
			return err
		}
	case "discover":
		tag, roots, err := parseTagDirArg(args)
		if err != nil {
			return err
		}
		if len(roots) == 0 {
			wd, err := os.Getwd()
			err = wrapErr(err, "could not get the current directory")
			if err != nil {
				return err
			}
			roots = append(roots, wd)
		}

		var dirs []string
		for _, root := range roots {
			tmp, err := Discover(opts, root, *optFile)
			err = wrapErr(err, "discovery failed")
			if err != nil {
				return err
			}
			dirs = append(dirs, tmp...)
		}
		for _, dir := range dirs {
			fmt.Fprintln(stdout, dir)
		}

		err = rmTag(tagman, tag, []string{})
		if err != nil {
			return err
		}
		err = addTag(tagman, tag, dirs)
		if err != nil {
			return err
		}
	default:
		tagitems := ParseTags(args)
		vt, err := VerifyTags(tagitems)
		err = wrapErr(err, "parsing arguments failed")
		if err != nil {
			return err
		}
		if len(vt.Dirs) == 0 && len(vt.Tags) == 0 {
			errorShowHelp("Directories or tags are required")
			return ErrHandled
		}

		if vt.Command.Str != "" {
			switch vt.Command.Op {
			case Add:
				err = addTag(tagman, vt.Command.Str, vt.Dirs)
				if err != nil {
					return err
				}
			case Remove:
				err = rmTag(tagman, vt.Command.Str, vt.Dirs)
				if err != nil {
					return err
				}
			case None:
				fallthrough
			default:
				return wrapErr(fmt.Errorf("improper command"), "internal error")
			}
			return nil
		}

		err = checkTags(vt.Tags)
		if err != nil {
			return err
		}

		vt.Dirs = tagman.Dirs(vt.Tags, vt.Dirs)

		err = RunCommands(opts, vt.Dirs, vt.Args)
		return wrapErr(err, "running command failed")
	}
	return nil
}
