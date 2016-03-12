package gogr

import (
	"fmt"
	"os"
	"os/exec"
)

func RunCommand(directory string, program string, args ...string) (err error) {
	// fmt.Println(directory)
	// fmt.Println(program)
	cmd := exec.Command(program, args...)
	cmd.Dir = directory
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return
	}

	return
}

func ParseDirectories(args []string) (dirs []string, rest []string, err error) {
	var info os.FileInfo
	pos := -1
	for i, arg := range args {
		info, err = os.Stat(arg)
		if err != nil || !info.IsDir() {
			pos = i
			err = nil
			break
		}
	}

	if pos == -1 {
		dirs = args
	} else {
		dirs = args[:pos]
		rest = args[pos:]
	}
	return
}
