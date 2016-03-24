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

	last := false
	pos := bytes.IndexRune(buf, '\n')
	for !last {
		if pos == -1 {
			pos = len(buf) - 1
			last = true
		}

		if len(buf) > 0 && (!wr([]byte(p.Prefix)) || !wr(buf[:pos+1])) {
			return
		}

		if !last {
			buf = buf[pos+1:]
			pos = bytes.IndexRune(buf, '\n')
		}
	}
	return
}
