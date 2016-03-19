package gogr

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type endOfProcessing bool

func (endOfProcessing) Error() string {
	return "End of processing"
}

func Discover(opts Options, root string, file string) (dirs []string, err error) {
	maxDepth, err := strconv.ParseInt(opts.Get("discover-max-depth", "5"), 10, 0)
	if err != nil {
		err = errors.New(fmt.Sprintf("Parsing maximum discovery depth failed: %s", err))
		return
	}

	dw := func(path string, info os.FileInfo, err error) (ret error) {
		relpath, err := filepath.Rel(root, path)
		pathlist := strings.Split(relpath, string(filepath.Separator))
		// fmt.Println(pathlist, len(pathlist))

		if len(pathlist) >= int(maxDepth) {
			ret = filepath.SkipDir
			return
		}
		if pathlist[len(pathlist)-1] == file {
			// fmt.Println("Found",path,"as directory", filepath.Dir(path))
			dirs = append(dirs, filepath.Dir(path))
		}

		return
	}

	err = filepath.Walk(root, dw)

	// fmt.Println("Hakemistot", dirs)
	return
}
