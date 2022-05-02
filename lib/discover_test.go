package gogr

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

var mockFileList []string
var mockWalkError error

func mockWalk(root string, WalkFn filepath.WalkFunc) (err error) {

	info, err := os.Stat(".")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Internal error: os.Stat pwd failed:", err)
		panic(nil)
	}

	for _, file := range mockFileList {
		WalkFn(file, info, nil)
	}
	err = mockWalkError
	return
}

func containsItem(list []string, item string) bool {
	for _, elem := range list {
		if elem == item {
			return true
		}
	}
	return false
}

func containsOnly(list []string, items []string) bool {
	for _, elem := range list {
		if !containsItem(items, elem) {
			return false
		}
	}
	return true
}

func allFromSlash(dirs []string) (ret []string) {
	for _, dir := range dirs {
		ret = append(ret, filepath.FromSlash(dir))
	}
	return
}

func TestDiscover(t *testing.T) {
	filepathWalk = mockWalk
	opts := optionMap{}

	opts.Set("discover-max-depth", "abc")
	_, err := Discover(&opts, ".", ".git")
	if err == nil {
		t.Error("Having non-number in discover-max-depth should fail Discover")
	}

	mockFileList = []string{"a/", "b/.git", "c/something/.git", "d"}

	opts.Set("discover-max-depth", "5")
	dirs, err := Discover(&opts, ".", ".git")
	checkDirs := []string{"b", "c/something"}
	checkDirs = allFromSlash(checkDirs)
	if !containsOnly(dirs, checkDirs) {
		t.Error("Discover should have discovered only dirs", checkDirs, "but it found", dirs)
	}

	opts.Set("discover-max-depth", "0")
	dirs, err = Discover(&opts, ".", ".git")
	if len(dirs) != 0 {
		t.Error("When max-depth == 0 Discover should not find anything. It found", dirs)
	}

}
