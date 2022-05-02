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
	n = len(buf)

	pos := -1
	line := make([]byte, 0, 1024)
	for len(buf) > 0 {
		pos = bytes.IndexRune(buf, '\n')
		if pos == -1 {
			pos = len(buf) - 1
		}

		line = line[:0]
		line = append(line, []byte(p.Prefix)...)
		line = append(line, buf[:pos+1]...)

		_, err = p.writer.Write(line)
		if err != nil {
			return
		}

		buf = buf[pos+1:]
	}
	return
}
