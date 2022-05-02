package gogr

import (
	"path/filepath"
	"testing"

	"github.com/kopoli/appkit"
)

func TestDefaultConfigFile(t *testing.T) {
	opts := appkit.NewOptions()
	path := DefaultConfigFile(opts)

	if !filepath.IsAbs(path) {
		t.Error("Default path should be proper and absolute:", path)
	}

	opts.Set("application-name", "something-else")
	path2 := DefaultConfigFile(opts)

	if !filepath.IsAbs(path2) || path2 == path {
		t.Error("Changing application name should have different path than default:",
			path, path2)
	}

	opts.Set("configuration-file", "non-default.json")

	path3 := DefaultConfigFile(opts)

	if !filepath.IsAbs(path3) || path3 == path2 {
		t.Error("Changing config file name should have different path than default:",
			path, path2, path3)
	}
}
