package gogr

import (
	"os"
	"reflect"
	"testing"
)

func TestTagManager(t *testing.T) {
	confFile := "test_config.json"
	defer func() { _ = os.Remove(confFile) }()

	tests := []struct {
		name string
		Tags map[string][]string
	}{
		{"Empty", map[string][]string{}},
		{"One tag", map[string][]string{"One": []string{}}},
		{"One tag with dir", map[string][]string{"One": []string{"abc"}}},
		{"Two tags with dir", map[string][]string{
			"One": []string{"abc"},
			"Two": []string{"a", "b"},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := &TagManager{
				ConfFile: confFile,
				Tags:     tt.Tags,
			}
			err := tm.Save()
			if err != nil {
				t.Errorf("TagManager save error = %v", err)
			}

			err = tm.Load()
			if err != nil {
				t.Errorf("TagManager load error = %v", err)
			}

			if !reflect.DeepEqual(tm.Tags, tt.Tags) {
				t.Errorf("Saved and loaded tags differ\nExpected: %s\n---Got: %s\n",
					tt.Tags, tm.Tags)
			}
		})
	}
}

func TestTagManager_ValidateTag(t *testing.T) {
	tests := []struct {
		name  string
		tag   string
		valid bool
	}{
		{"Empty", "", false},
		{"Single char", "a", true},
		{"Dash", "-", false},
		{"String", "longername", true},
		{"Dash string", "longer-name", false},
		{"Underscore", "longer_name", false},
		{"Numbers", "1215", true},
		{"Plus", "1215+abc", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := &TagManager{}
			if got := tm.ValidateTag(tt.tag); got != tt.valid {
				t.Errorf("TagManager.ValidateTag() = %v, want %v", got, tt.valid)
			}
		})
	}
}
