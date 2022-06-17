package gogr

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kopoli/appkit"
)

var filepathWalk = filepath.Walk

// Discover directories beginning from root that contain the given file
func Discover(opts appkit.Options, root string, file string) ([]string, error) {
	var dirs []string
	maxDepth, err := strconv.ParseInt(opts.Get("discover-max-depth", "5"), 10, 0)
	if err != nil {
		err = fmt.Errorf("parsing maximum discovery depth failed: %v", err)
		return nil, err
	}

	dw := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		var relpath string
		relpath, _ = filepath.Rel(root, path)
		pathlist := strings.Split(relpath, string(filepath.Separator))

		if len(pathlist) >= int(maxDepth) {
			return filepath.SkipDir
		}
		if pathlist[len(pathlist)-1] == file {
			dirs = append(dirs, filepath.Dir(path))
		}

		return nil
	}

	err = filepathWalk(root, dw)
	return dirs, err
}
