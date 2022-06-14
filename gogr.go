package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/kopoli/appkit"
	gogr "github.com/kopoli/gogr/lib"
)

var (
	version     = "Undefined"
	timestamp   = "Undefined"
	buildGOOS   = "Undefined"
	buildGOARCH = "Undefined"
	progVersion = "" + version
)

func fault(err error, message string, arg ...string) {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %s%s: %s\n", message, strings.Join(arg, " "), err)
		os.Exit(1)
	}
}

func main() {
	opts := appkit.NewOptions()

	opts.Set("program-name", os.Args[0])
	opts.Set("program-version", progVersion)
	opts.Set("program-timestamp", timestamp)
	opts.Set("program-buildgoos", buildGOOS)
	opts.Set("program-buildgoarch", buildGOARCH)
	opts.Set("configuration-file", "config.json")

	err := gogr.Main(os.Args, opts)
	if err == gogr.ErrLicenses {
		l, err := GetLicenses()
		fault(err, "Getting licenses failed")
		s, err := appkit.LicenseString(l)
		fault(err, "Interpreting licenses failed")
		fmt.Print(s)
		os.Exit(0)
	}

	if err != nil {
		if err != gogr.ErrHandled {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		}
		os.Exit(1)
	}
	os.Exit(0)
}
