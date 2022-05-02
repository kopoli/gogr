package gogr

import (
	"github.com/casimir/xdg-go"
	"github.com/kopoli/appkit"
)

// DefaultConfigFile gets the default configuration file name based on given
// appkit.Options
func DefaultConfigFile(opts appkit.Options) string {
	app := xdg.App{Name: opts.Get("application-name", "gogr")}
	return app.ConfigPath(opts.Get("configuration-file", "config.json"))
}
