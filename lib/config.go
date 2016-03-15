package gogr


import "github.com/casimir/xdg-go"


func DefaultConfigFile(opts Options) string {
	app := xdg.App{Name: opts.Get("application-name", "gogr")}
	return app.ConfigPath(opts.Get("configuration-file", "config.json"))
}


