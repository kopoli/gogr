package gogr

import (
	"bytes"
	"io"
	"sync"
)

type PrefixedWriter struct {
	Prefix []byte
	Eol    []byte
	Out    io.Writer
	buf    *bytes.Buffer // buffer to house incomplete lines

	sync.Mutex
}

func NewPrefixedWriter(prefix string, writer io.Writer) *PrefixedWriter {
	return &PrefixedWriter{
		Prefix: []byte(prefix),
		Eol:    []byte("\n"),
		Out:    writer,
		buf:    &bytes.Buffer{},
	}
}

func (p *PrefixedWriter) Write(buf []byte) (n int, err error) {
	// If no bytes to write
	if len(buf) == 0 {
		return 0, nil
	}

	p.Lock()
	defer p.Unlock()

	n = len(buf)

	// If only one line to write without newline
	lastLineIdx := bytes.LastIndexByte(buf, '\n')
	if lastLineIdx < 0 {
		// If there is nothing in the buffer
		if p.buf.Len() == 0 {
			p.buf.Write(p.Prefix)
		}

		// Write only into the buffer
		_, err = p.buf.Write(buf)
		return n, err
	}

	endsInNewline := (buf[len(buf)-1] == '\n')
	lines := bytes.Split(buf, []byte{'\n'})

	// If given data ends in newline, skip the last line
	if endsInNewline {
		lines = lines[:len(lines)-1]
	}

	for i := range lines {
		// If either not first line or first and nothing in buffer
		if i > 0 || (i == 0 && p.buf.Len() == 0) {
			p.buf.Write(p.Prefix)
		}

		// If either not last line or last line when ends in a newline
		if i < len(lines)-1 || (i == len(lines)-1 && endsInNewline) {
			p.buf.Write(lines[i])
			p.buf.Write(p.Eol)
		}
	}

	// Write to output
	_, err = p.buf.WriteTo(p.Out)
	if err != nil {
		return 0, err
	}

	// Write the last line to buffer for next time if newline isn't present
	if !endsInNewline {
		_, err = p.buf.Write(lines[len(lines)-1])
		if err != nil {
			return 0, err
		}
	}

	return n, nil
}
