package gogr

import (
	"bytes"
	"fmt"
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

type strOp func(string) bool

type checkOp struct {
	what string
	op   strOp
}

type checker struct {
	checks []checkOp
}

func (c *checker) add(what string, op strOp) *checker {
	c.checks = append(c.checks, checkOp{what, op})
	return c
}

func (c *checker) Err(op strOp) *checker {
	return c.add("error", op)
}

func (c *checker) Out(op strOp) *checker {
	return c.add("output", op)
}

func (c *checker) Conf(op strOp) *checker {
	return c.add("config", op)
}

func (c *checker) Check(err error, output, confdata string) error {
	for i := range c.checks {
		ret := true
		switch c.checks[i].what {
		case "error":
			if err == nil {
				err = fmt.Errorf("")
			}
			ret = c.checks[i].op(err.Error())
		case "output":
			ret = c.checks[i].op(output)
		case "config":
			ret = c.checks[i].op(confdata)
		}
		if !ret {
			return fmt.Errorf("op num %d type %s failed", i, c.checks[i].what)
		}
	}
	return nil
}

func Test_Main(t *testing.T) {
	confFile := "test.conf"

	chk := func() *checker {
		return &checker{}
	}

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

	not := func(op strOp) strOp {
		return func(s string) bool {
			return !op(s)
		}
	}

	oneTag := `{"tags": {"one": ["/tmp"]}}`
	twoTags := `{"tags": {"one": ["/tmp"], "two": []}}`

	tests := []struct {
		name     string
		tagsJSON string
		args     []string
		check    *checker
	}{
		{"Empty args", "", []string{},
			chk().Out(isFound("^Error.*Arguments required"))},
		{"Invalid args with no conf", "", []string{"blerg"},
			chk().Out(is("")).Err(isFound("loading tags failed"))},
		{"Invalid args with conf", `{"tags":{}}`, []string{"blerg"},
			chk().Out(isFound("^Error.*Directories or tags are required")).Err(isFound("handled"))},
		{"Help arg, no error", "", []string{"-h"},
			chk().Out(isFound("^Usage")).Err(is(""))},
		{"Help arg, no error 2", "", []string{"--help"},
			chk().Out(isFound("^Usage")).Err(is(""))},

		{"Empty tag list", "{}", []string{"tag"},
			chk().Out(is("")).Err(is("")).Conf(is("{}"))},
		{"Empty tag list 2", "{}", []string{"tag", "list"},
			chk().Out(is("")).Err(is("")).Conf(is("{}"))},
		{"One tag", oneTag, []string{"tag", "list"},
			chk().Out(is("one\n")).Err(is("")).Conf(is(oneTag))},
		{"Two tags", twoTags, []string{"tag", "list"},
			chk().Out(is("one\ntwo\n")).Err(is("")).Conf(is(twoTags))},
		{"One tag, list dir", oneTag, []string{"tag", "list", "one"},
			chk().Out(is("/tmp\n")).Err(is("")).Conf(is(oneTag))},
		{"Two tags, list dir", twoTags, []string{"tag", "list", "one"},
			chk().Out(is("/tmp\n")).Err(is("")).Conf(is(twoTags))},

		{"Add tag", "{}", []string{"tag", "add", "one", "/tmp"},
			chk().Out(is("")).Err(is("")).Conf(isFound("tags")).Conf(isFound("one")).Conf(isFound("/tmp"))},
		{"Add tag @syntax", "{}", []string{"+@one", "/tmp"},
			chk().Out(is("")).Err(is("")).Conf(isFound("tags")).Conf(isFound("one")).Conf(isFound("/tmp"))},
		{"Add tag to existing", oneTag, []string{"tag", "add", "one", "/root"},
			chk().Out(is("")).Err(is("")).Conf(isFound("tags")).Conf(isFound("one")).
				Conf(isFound("/tmp")).Conf(isFound("/root"))},

		{"Remove tag", oneTag, []string{"tag", "delete", "one"},
			chk().Out(is("")).Err(is("")).Conf(isFound("tags")).Conf(not(isFound("one")))},

		{"Remove tag from two", twoTags, []string{"tag", "delete", "two"},
			chk().Out(is("")).Err(is("")).Conf(isFound("tags")).
				Conf(isFound("one")).Conf(not(isFound("two")))},
		{"Remove one tag from two", twoTags, []string{"tag", "delete", "one"},
			chk().Out(is("")).Err(is("")).Conf(isFound("tags")).
				Conf(not(isFound("one"))).Conf(isFound("two"))},

		{"Remove dir from tag", oneTag, []string{"tag", "delete", "one", "/tmp"},
			chk().Out(is("")).Err(is("")).Conf(isFound("tags")).Conf(isFound("one")).Conf(not(isFound("/tmp")))},
		{"Remove dir from tag @syntax", oneTag, []string{"-@one", "/tmp"},
			chk().Out(is("")).Err(is("")).Conf(isFound("tags")).Conf(isFound("one")).Conf(not(isFound("/tmp")))},

		{"Option -version", "", []string{"-version"},
			chk().Out(isFound("Built.*with"))},
		{"Option -v", "", []string{"-v"},
			chk().Out(isFound("Built.*with"))},
		{"Option -licenses", "", []string{"-licenses"},
			chk().Out(is("")).Err(isFound("license display requested"))},

		{"Discover current dir", "{}", []string{"discover", "-file", "main_test.go", "this", "."},
			chk().Out(is(".\n")).Err(is("")).Conf(isFound("/lib"))},
		{"Discover current dir w/o arg", "{}", []string{"discover", "-file", "main_test.go", "this"},
			chk().Out(isFound("lib\n")).Err(is("")).Conf(isFound("/lib"))},

		{"Run command in tag", oneTag, []string{"@one", "pwd"},
			chk().Out(is("tmp: /tmp\n")).Err(is(""))},
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

			tt.args = append([]string{"progname"}, tt.args...)

			DefaultConfigFile = func(o appkit.Options) string {
				return confFile
			}
			defer func() { DefaultConfigFile = defaultConfigFile }()

			err = Main(tt.args, opts)

			// Ignore error, as the file might not be created
			b, _ := ioutil.ReadFile(confFile)

			out := buf.String()
			confData := string(b)
			err2 := tt.check.Check(err, out, confData)
			if err2 != nil {
				t.Errorf("Checking failed with %v\nerror %v\noutput: %s\nconffile: %s\n", err2, err, out, confData)
			}
		})
	}
}
