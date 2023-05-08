package filter

import (
	"bufio"
	"io"
)

type BytesSource interface {
	Scan() bool
	Bytes() []byte
	Err() error
}

func NewBytesSource(r io.Reader) BytesSource { return bufio.NewScanner(r) }
