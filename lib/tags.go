package gogr

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/kopoli/appkit"
)

// TagManager is a repository for tags
type TagManager struct {
	ConfFile string              `json:"-"`
	Tags     map[string][]string `json:"tags"`
}

// NewTagManager creates a repository for tags, which it writes to the given
// "configuration-file" from opts.
func NewTagManager(opts appkit.Options) (ret TagManager) {
	ret.ConfFile = opts.Get("configuration-file", "config.json")
	ret.Tags = make(map[string][]string)
	_ = ret.Load()

	return
}

// Save saves the tags into a configuration file.
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

// Load loads the tags from a configuration file.
func (t *TagManager) Load() (err error) {
	b, err := ioutil.ReadFile(t.ConfFile)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &t)
	return
}

// deduplicate removes duplicates from a list of strings
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

// cleanup returns absolute directory names from a given list of directories.
// Returns only directories that exist.
func cleanup(dirs []string) (ret []string) {
	for _, dir := range dirs {
		dir, err := filepath.Abs(dir)
		if err == nil && isDirectory(dir) {
			ret = append(ret, filepath.Clean(dir))
		}
	}
	return
}

// ValidateTag validates the tag string. Returns true if valid.
func (t *TagManager) ValidateTag(tag string) bool {
	re := regexp.MustCompile("[a-zA-Z0-9]+")
	return re.MatchString(tag)
}

// Add adds given directories to given tag. The tag is created if necessary.
func (t *TagManager) Add(tag string, dirs ...string) {
	t.Tags[tag] = deduplicate(cleanup(append(t.Tags[tag], dirs...)))
}

// Remove removes either given directories from a tag. Alternatively, if the
// list of directories is empty, removes the whole tag.
func (t *TagManager) Remove(tag string, dirs ...string) {
	if len(dirs) == 0 {
		delete(t.Tags, tag)
		return
	}

	dirs = deduplicate(cleanup(dirs))

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

// Dirs returns a combined list of directories of given tags.
func (t *TagManager) Dirs(tags []string, dirs []string) (ret []string) {
	for _, tag := range tags {
		ret = append(ret, t.Tags[tag]...)
	}
	ret = deduplicate(cleanup(append(ret, dirs...)))

	sort.Strings(ret)

	return
}

// AreProper checks if the given tags exist. Returns the list of non-existing
// tags.
func (t *TagManager) AreProper(tags []string) (invalid []string) {
	for _, tag := range tags {
		_, ok := t.Tags[tag]
		if !ok {
			invalid = append(invalid, tag)
		}
	}
	return
}

// isDirectory checks if given path is a directory.
func isDirectory(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// ParseDirectory parses a list of strings into list of directories and list
// of arguments. The directories are the first part of the given list of
// strings. The rest arguments are the rest.
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

// ItemType denotes different types of strings within command line arguments.
type ItemType int

const (
	Tag ItemType = iota
	Arg
)

// Operation lists the different types of operations for command line
// arguments,
type Operation int

const (
	None Operation = iota
	Add
	Remove
)

// TagItem is a parsed result of a command line parsing with ParseTags.
type TagItem struct {
	Type ItemType
	Op   Operation
	Str  string
}

// ParseTags parses a list of strings into a list of TagItem structures.
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

type VerifyTagsRet struct {
	Command TagItem
	Tags    []string
	Dirs    []string
	Args    []string
}

// VerifyTags verifies the logic for adding and removing tags on the command line.
func VerifyTags(items []TagItem) (*VerifyTagsRet, error) {
	ret := &VerifyTagsRet{}
	var err error

	if len(items) == 0 {
		err = errors.New("arguments required")
		return nil, err
	}

	for i, item := range items {
		switch {
		case item.Op != None:
			if i > 0 {
				err = errors.New("tagging must be the first argument")
				return nil, err
			}
			ret.Command = item
		case item.Type == Tag:
			if len(ret.Dirs) > 0 || len(ret.Args) > 0 {
				err = errors.New("tags must precede directories or commands")
				return nil, err
			}
			ret.Tags = append(ret.Tags, item.Str)
		case item.Type == Arg:
			if len(ret.Args) > 0 || !isDirectory(item.Str) {
				ret.Args = append(ret.Args, item.Str)
			} else {
				ret.Dirs = append(ret.Dirs, item.Str)
			}
		}
	}

	if ret.Command.Str == "" && len(ret.Args) == 0 {
		err = errors.New("no command to run given")
		return nil, err
	}

	if ret.Command.Str != "" && (len(ret.Dirs) == 0 || len(ret.Args) > 0) {
		err = errors.New("tagging requires one or more directories and zero commands")
		return nil, err
	}

	return ret, nil
}
