package main

import (
	"fmt"
	"os"

	"github.com/kopoliitti/gogr/lib"
	"github.com/jawher/mow.cli"
)

var (
	progName     = "gogr"
	majorVersion = "0"
	version      = "Undefined"
	timestamp    = "Undefined"
	progVersion  = majorVersion + "-" + version
)

func main() {
	app := cli.App(progName, "Run commands in multiple directories")

	app.Spec = "[OPTIONS] [ARG...]"
	app.Version(progName, progVersion)

	app.Command("tag", "Tag management", func(cmd *cli.Cmd) {
		cmd.Command("add", "Add tag to path", nil)
		cmd.Command("rm", "Remove tag from paths", nil)
	})

	optVerbose := app.BoolOpt("verbose V", false, "Print verbose output")
	argArg := app.StringsArg("ARG", nil, "Directories and command to be run")

	app.Action = func() {
		fmt.Println(*optVerbose)
		// fmt.Println(*argArg)
		dirs, rest, err := gogr.ParseDirectories(*argArg)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(dirs)
		fmt.Println(rest)

		for _, dir := range dirs {
			err = gogr.RunCommand(dir, rest[0], rest[1:]...)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	app.Run(os.Args)
}
