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

type BytesSourceErr struct{ err error }

func BytesSourceErrNew(err error) BytesSourceErr { return BytesSourceErr{err} }

func (b BytesSourceErr) Scan() bool    { return false }
func (b BytesSourceErr) Bytes() []byte { return nil }
func (b BytesSourceErr) Err() error    { return b.err }

func (b BytesSourceErr) AsIf() BytesSource { return b }
func (b BytesSourceErr) Close() error      { return nil }
