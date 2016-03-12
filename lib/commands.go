package gogr

import (
	"errors"
	"os"
	"os/exec"
)

func RunCommand(directory string, program string, args ...string) (err error) {
	cmd := exec.Command(program, args...)
	cmd.Dir = directory

	err = cmd.Run()
	if err != nil {
		return
	}

	return
}

func ParseDirectories(args []string) (dirs []string, rest []string, err error) {
	var info os.FileInfo
	for i, arg := range args {
		info, err = os.Stat(arg)
		if err != nil || !info.IsDir() {
			if i == 0 {
				err = errors.New("No directories")
				return
			}
			dirs = args[:i]
			rest = args[i:]
			err = nil
			return
		}
	}

	err = errors.New("No command")
	return
}
