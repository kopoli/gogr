package gogr

import (
	"bytes"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"testing"

	"github.com/kopoli/appkit"
)

func Test_escapeTagArgs(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		unescape bool
		want     []string
	}{
		{"Empty", []string{}, true, []string{}},
		{"One arg", []string{"a"}, true, []string{"a"}},
		{"Two args", []string{"a", "b"}, true, []string{"a", "b"}},
		{"Escape prefixed", []string{"-@"}, false, []string{"\\-@"}},
		{"Don't unescape prefixed", []string{"-@"}, true, []string{"-@"}},
		{"Unescape backslashed", []string{"\\-@"}, true, []string{"-@"}},
		{"Escape backslashes in front", []string{"\\"}, false, []string{"\\\\"}},
		{"Escape two", []string{"\\", "-@"}, false, []string{"\\\\", "\\-@"}},
		{"Unescape two", []string{"\\\\", "\\-@"}, true, []string{"\\", "-@"}},
		{"Don't escape in the middle", []string{"ab\\\\", "ab\\-@"}, false, []string{"ab\\\\", "ab\\-@"}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := escapeTagArgs(tt.args, tt.unescape); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("escapeTagArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Main(t *testing.T) {
	confFile := "test.conf"

	type strOp func(string) bool

	is := func(str string) strOp {
		return func(s string) bool {
			return s == str
		}
	}

	isFound := func(re string) strOp {
		rex := regexp.MustCompile(re)
		return func(s string) bool {
			return rex.MatchString(s)
		}
	}

	onetag := `{"tags": {"one": ["/tmp"]}}`

	tests := []struct {
		name     string
		tagsJSON string
		args     []string
		wantErr  bool
		output   strOp
	}{
		{"Empty args", "", []string{}, true, isFound("^Error.*Arguments required")},
		{"Invalid args", "", []string{"blerg"}, true, isFound("^Error.*Directories or tags are required")},
		{"Help arg, no error", "", []string{"-h"}, false, isFound("^Usage")},
		{"Help arg, no error 2", "", []string{"--help"}, false, isFound("^Usage")},
		{"Empty tag list", "", []string{"tag"}, false, is("")},
		{"Empty tag list 2", "", []string{"tag", "list"}, false, is("")},
		{"One tag", onetag, []string{"tag", "list"}, false, is("jeje")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			stdout = buf
			stderr = buf
			var err error
			defer func() { _ = os.Remove(confFile) }()
			if tt.tagsJSON != "" {
				err = ioutil.WriteFile(confFile, []byte(tt.tagsJSON), 0666)
				if err != nil {
					t.Errorf("Could not write conffile: %v", err)
				}
			}
			opts := appkit.NewOptions()
			opts.Set("configuration-file", confFile)

			tt.args = append([]string{"p"}, tt.args...)

			err = Main(tt.args, opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Expected fail %v, error %v", tt.wantErr, err)
			}

			out := buf.String()
			if tt.output != nil && !tt.output(out) {
				t.Errorf("The output didn't match expected:\n %s", out)
			}
		})
	}
}
