package gogr

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/kopoli/appkit"
)

func RunCommand(directory string, program string, args ...string) (err error) {
	dir := filepath.Base(directory)
	pwo := NewPrefixedWriter(fmt.Sprintf("%s: ", dir), os.Stdout)
	pwe := NewPrefixedWriter(fmt.Sprintf("%s(err): ", dir), os.Stderr)
	cmd := exec.Command(program, args...)
	cmd.Dir = directory
	cmd.Stdout = &pwo
	cmd.Stderr = &pwe

	err = cmd.Run()
	if err != nil {
		fmt.Fprintf(&pwe, "Command failed: %s\n", err)
		return
	}

	return
}

func RunCommands(opts appkit.Options, dirs []string, args []string) (err error) {
	concurrent, err := strconv.ParseBool(opts.Get("concurrent", "false"))
	if err != nil {
		return
	}

	if concurrent {
		wg := sync.WaitGroup{}
		for _, dir := range dirs {
			wg.Add(1)
			go func(dir string) {
				defer wg.Done()
				_ = RunCommand(dir, args[0], args[1:]...)
			}(dir)
		}
		wg.Wait()
	} else {
		for _, dir := range dirs {
			_ = RunCommand(dir, args[0], args[1:]...)
		}
	}
	return
}
