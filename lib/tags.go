package gogr

import (
	"encoding/json"
	"io/ioutil"
)

type TagManager struct {
	ConfFile string              `json:"-"`
	Tags     map[string][]string `json:"tags"`
}

type TagData struct {
}

func NewTagManager(file string) (ret TagManager) {
	ret.ConfFile = file
	ret.Tags = make(map[string][]string)
	ret.Load()

	return
}

func (t *TagManager) Save() (err error) {
	b, err := json.MarshalIndent(t.Tags, " ", "    ")
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

	err = json.Unmarshal(b, &t.Tags)
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

func (t *TagManager) Add(tag string, dirs ...string) {
	t.Tags[tag] = deduplicate(append(t.Tags[tag], dirs...))
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

func (t *TagManager) Dirs(tag string) (dirs []string) {
	return t.Tags[tag]
}
