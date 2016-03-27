package gogr

import (
	"bytes"
	"io"
)

type PrefixedWriter struct {
	Prefix string
	writer io.Writer
}

func NewPrefixedWriter(prefix string, writer io.Writer) (ret PrefixedWriter) {
	ret.Prefix = prefix
	ret.writer = writer
	return
}

func (p *PrefixedWriter) Write(buf []byte) (n int, err error) {
	var wr = func(buf []byte) bool {
		_, err = p.writer.Write(buf)
		if err != nil {
			return false
		}
		return true
	}

	n = len(buf)

	pos := -1
	for len(buf) > 0 {
		pos = bytes.IndexRune(buf, '\n')
		if pos == -1 {
			pos = len(buf) - 1
		}

		if !wr([]byte(p.Prefix)) || !wr(buf[:pos+1]) {
			return
		}

		buf = buf[pos+1:]
	}
	return
}
