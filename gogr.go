package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/jawher/mow.cli"
	"github.com/kopoliitti/gogr/lib"
)

var (
	progName     = "gogr"
	majorVersion = "0"
	version      = "Undefined"
	timestamp    = "Undefined"
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

func main() {
	tags := gogr.NewTagManager("test.json")

	app := cli.App(progName, "Run commands in multiple directories")
	app.Spec = "[OPTIONS] [ARG...]"

	app.Version("version v", fmt.Sprintf("%s: %s\ngo: %s\nBuilt with: %s on %s/%s\n", progName, progVersion,
		runtime.Version(), runtime.Compiler, runtime.GOOS, runtime.GOARCH))
	optVerbose := app.BoolOpt("verbose V", false, "Print verbose output")

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
				fmt.Println(*dirArg)
				fmt.Println(*tagArg)
				dirs := parseDir(*dirArg)
				if len(dirs) == 0 {
					wd, err := os.Getwd()
					if err != nil {
						fault(err, "Getting working directory failed")
					}
					dirs = append(dirs, wd)
				}
				tags.Add(*tagArg, dirs...)
				tags.Save()
			}
		})
		cmd.Command("rm", "Remove tag from paths", func(cmd *cli.Cmd) {
			cmd.Spec = "TAG [DIR...]"
			tagArg := cmd.StringArg("TAG", "", "Tag to remove")
			dirArg := cmd.StringsArg("DIR", nil, "Directories to remove the tag from")

			cmd.Action = func() {
				tags.Remove(*tagArg, parseDir(*dirArg)...)
				tags.Save()
			}
		})
	})

	argArg := app.StringsArg("ARG", nil, "Directories and command to be run")

	app.Action = func() {
		fmt.Println(*optVerbose)
		dirs, rest, err := gogr.ParseDirectories(*argArg)
		if err != nil {
			fault(err, "Directory parsing failed")
			return
		}
		fmt.Println(dirs)
		fmt.Println(rest)

		// for _, dir := range dirs {
		// 	err = gogr.RunCommand(dir, rest[0], rest[1:]...)
		// 	if err != nil {
		// 		fault(err, "Directory parsing failed")
		// 	}
		// }
	}

	app.Run(os.Args)
}
