package gogr

import (
	"encoding/json"
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
