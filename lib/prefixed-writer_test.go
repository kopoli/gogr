package gogr

import (
	"bytes"
	"fmt"
	"testing"
)

func TestPrefixedWriter(t *testing.T) {
	buf := bytes.Buffer{}
	prefix := "str: "
	pw := NewPrefixedWriter(prefix, &buf)

	if buf.String() != "" {
		t.Error("PrefixedWriter should not write anything when it's created.")
	}

	data := "first row"
	fmt.Fprintln(&pw, data)

	result := prefix + data + "\n"
	testContain := func(result string) {
		if buf.String() != result {
			t.Error("The prefixed writer should contain", result, "but it contains", buf.String())
		}
	}
	testContain(result)

	data = "second row"

	fmt.Fprintln(&pw, data)
	result = result + prefix + data + "\n"
	testContain(result)

	fmt.Fprintln(&pw, "")

	testContain(result + prefix + "\n")
}

func TestPrefixWriterNoNewline(t *testing.T) {
	buf := bytes.Buffer{}
	prefix := "str: "
	pw := NewPrefixedWriter(prefix, &buf)

	data := "something"
	fmt.Fprintf(&pw, data)

	result := prefix + data
	if buf.String() != result {
		t.Error("The prefixed writer should contain", result, "But it contains", buf.String())
	}
}

func TestPrefixWriterEmpty(t *testing.T) {
	buf := bytes.Buffer{}
	prefix := "str: "
	pw := NewPrefixedWriter(prefix, &buf)
	fmt.Fprintf(&pw, "")

	if buf.String() != "" {
		t.Error("The prefixed writer should return empty string on empty input. Got", buf.String())
	}
}
