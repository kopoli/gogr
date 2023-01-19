package gogr

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/kopoli/appkit"
)

func RunCommand(hidePrefix bool, directory string, program string, args ...string) (err error) {
	dir := filepath.Base(directory)
	var pfx, errPfx string
	if !hidePrefix {
		pfx = fmt.Sprintf("%s: ", dir)
		errPfx = fmt.Sprintf("%s(err): ", dir)
	}
	pwo := NewPrefixedWriter(pfx, stdout)
	pwe := NewPrefixedWriter(errPfx, stderr)

	cmd := exec.Command(program, args...)
	cmd.Dir = directory
	cmd.Stdout = pwo
	cmd.Stderr = pwe

	err = cmd.Run()
	if err != nil {
		fmt.Fprintf(pwe, "Command failed: %s\n", err)
		return
	}

	return
}

func RunCommands(opts appkit.Options, dirs []string, args []string) (err error) {
	concurrent := opts.IsSet("concurrent")
	hidePrefix := opts.IsSet("hide-prefix")

	if concurrent {
		wg := sync.WaitGroup{}
		for _, dir := range dirs {
			wg.Add(1)
			go func(dir string) {
				defer wg.Done()
				_ = RunCommand(hidePrefix, dir, args[0], args[1:]...)
			}(dir)
		}
		wg.Wait()
	} else {
		for _, dir := range dirs {
			_ = RunCommand(hidePrefix, dir, args[0], args[1:]...)
		}
	}
	return
}
