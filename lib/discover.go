package gogr

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var filepathWalk = filepath.Walk

// Discover directories beginning from root that contain the given file
func Discover(opts Options, root string, file string) (dirs []string, err error) {
	maxDepth, err := strconv.ParseInt(opts.Get("discover-max-depth", "5"), 10, 0)
	if err != nil {
		err = fmt.Errorf("Parsing maximum discovery depth failed: %s", err)
		return
	}

	dw := func(path string, info os.FileInfo, err error) (ret error) {
		relpath, err := filepath.Rel(root, path)
		pathlist := strings.Split(relpath, string(filepath.Separator))

		if len(pathlist) >= int(maxDepth) {
			ret = filepath.SkipDir
			return
		}
		if pathlist[len(pathlist)-1] == file {
			dirs = append(dirs, filepath.Dir(path))
		}

		return
	}

	err = filepathWalk(root, dw)
	return
}
