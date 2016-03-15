package gogr

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"
)

type PrefixedWriter struct {
	Prefix string
	writer io.Writer
}

func NewPrefixedWriter(prefix string, writer io.Writer) (ret PrefixedWriter) {
	ret.Prefix = prefix
	ret.writer = writer
	return
}

func (p *PrefixedWriter) Write(buf []byte) (n int, err error) {
	var wr = func(buf []byte) bool {
		_, err = p.writer.Write(buf)
		if err != nil {
			return false
		}
		return true
	}

	n = len(buf)

	last := false
	pos := bytes.IndexRune(buf, '\n')
	for !last {
		if pos == -1 {
			pos = len(buf) - 1
			last = true
		}

		if len(buf) > 0 && (!wr([]byte(p.Prefix)) || !wr(buf[:pos+1])) {
			return
		}

		if !last {
			buf = buf[pos+1:]
			pos = bytes.IndexRune(buf, '\n')
		}
	}
	return
}

func RunCommand(directory string, program string, args ...string) (err error) {
	dir := filepath.Base(directory)
	pwo := NewPrefixedWriter(fmt.Sprintf("%s: ", dir), os.Stderr)
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

func RunCommands(opts Options, dirs []string, args []string) (err error) {
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

func isDirectory(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func ParseDirectories(args []string) (dirs []string, rest []string, err error) {
	pos := -1
	for i, arg := range args {
		if !isDirectory(arg) {
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

type ItemType int

const (
	Tag ItemType = iota
	Arg
)

type Operation int

const (
	None Operation = iota
	Add
	Remove
)

type TagItem struct {
	Type ItemType
	Op   Operation
	Str  string
}

func ParseTags(args []string) (ret []TagItem) {
	if len(args) == 0 {
		return
	}

	re := regexp.MustCompile("([+-]?)@([[:alpha:]]+)")

	for _, arg := range args {
		var ta TagItem
		tag := re.FindStringSubmatch(arg)
		ta.Op = None
		if len(tag) == 0 {
			ta.Type = Arg
			ta.Str = arg
		} else {
			ta.Str = tag[2]
			if tag[1] == "+" {
				ta.Op = Add
			} else if tag[1] == "-" {
				ta.Op = Remove
			}
		}
		ret = append(ret, ta)
	}

	return
}

func VerifyTags(items []TagItem) (command TagItem, tags []string, dirs []string, args []string, err error) {
	if len(items) == 0 {
		err = errors.New("Arguments required")
		return
	}

	for i, item := range items {
		if item.Op != None {
			if i > 0 {
				err = errors.New("Tagging must be the first argument")
				return
			}
			command = item
			continue
		} else if item.Type == Tag {
			if len(dirs) > 0 || len(args) > 0 {
				err = errors.New("Tags must precede directories or commands")
				return
			}
			tags = append(tags, item.Str)
			continue
		} else if item.Type == Arg {
			if len(args) > 0 || !isDirectory(item.Str) {
				args = append(args, item.Str)
			} else {
				dirs = append(dirs, item.Str)
			}
		}
	}

	if command.Str == "" && len(args) == 0 {
		err = errors.New("No command to run given")
		return
	}

	if command.Str != "" && (len(dirs) == 0 || len(args) > 0) {
		err = errors.New("Tagging requires one or more directories and zero commands.")
	}

	return
}
