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
