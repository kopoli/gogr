package gogr

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

type TagManager struct {
	ConfFile string              `json:"-"`
	Tags     map[string][]string `json:"tags"`
}

func NewTagManager(opts Options) (ret TagManager) {
	ret.ConfFile = opts.Get("configuration-file", "config.json")
	ret.Tags = make(map[string][]string)
	_ = ret.Load()

	return
}

func (t *TagManager) Save() (err error) {
	b, err := json.MarshalIndent(t, " ", "    ")
	if err != nil {
		return
	}
	err = os.MkdirAll(filepath.Dir(t.ConfFile), 0755)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(t.ConfFile, b, 0666)
	return
}

func (t *TagManager) Load() (err error) {
	b, err := ioutil.ReadFile(t.ConfFile)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &t)
	return
}

func deduplicate(strings []string) (ret []string) {
	m := make(map[string]bool)
	for _, str := range strings {
		m[str] = true
	}
	for k := range m {
		ret = append(ret, k)
	}
	return
}

func (t *TagManager) ValidateTag(tag string) bool {
	re := regexp.MustCompile("[a-zA-Z0-9]+")
	return re.MatchString(tag)
}

func (t *TagManager) Add(tag string, dirs ...string) {
	t.Tags[tag] = deduplicate(cleanup(append(t.Tags[tag], dirs...)))
}

// If dirs is empty, remove the whole tag, otherwise remove the given
// directories from the tag
func (t *TagManager) Remove(tag string, dirs ...string) {
	if len(dirs) == 0 {
		delete(t.Tags, tag)
		return
	}

	var ret []string

	for _, dir := range t.Tags[tag] {
		remove := false
		for _, rmdir := range dirs {
			if dir == rmdir {
				remove = true
				break
			}
		}
		if !remove {
			ret = append(ret, dir)
		}
	}
	t.Tags[tag] = ret
}

func cleanup(dirs []string) (ret []string) {
	for _, dir := range dirs {
		dir, err := filepath.Abs(dir)
		if err == nil && isDirectory(dir) {
			ret = append(ret, filepath.Clean(dir))
		}
	}
	return
}

func (t *TagManager) Dirs(tags []string, dirs []string) (ret []string) {
	for _, tag := range tags {
		ret = append(ret, t.Tags[tag]...)
	}
	ret = deduplicate(cleanup(append(ret, dirs...)))

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

	re := regexp.MustCompile("([+-]?)@([a-zA-Z0-9]+)")

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
