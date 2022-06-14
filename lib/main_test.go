package gogr

import (
	"reflect"
	"testing"
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
		t.Run(tt.name, func(t *testing.T) {
			if got := escapeTagArgs(tt.args, tt.unescape); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("escapeTagArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}
