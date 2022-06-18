package gogr

import (
	"path/filepath"

	"github.com/OpenPeeDeeP/xdg"
	"github.com/kopoli/appkit"
)

// DefaultConfigFile gets the default configuration file name based on given
// appkit.Options
func defaultConfigFile(opts appkit.Options) string {
	path := xdg.New("", opts.Get("application-name", "gogr")).ConfigHome()
	return filepath.Join(path, opts.Get("configuration-file", "config.json"))
}

var DefaultConfigFile = defaultConfigFile
