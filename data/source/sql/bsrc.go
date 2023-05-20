package sqldata

import (
	"database/sql"
	"io"

	lf "github.com/takanoriyanagitani/go-log-filter"
)

type Rows2Byte func(*sql.Rows) ([]byte, error)

var Val2Byte Rows2Byte = func(rows *sql.Rows) (val []byte, e error) {
	e = rows.Scan(&val)
	return
}

type BytesSourceCloser interface {
	lf.BytesSource
	io.Closer
}

type Rows2Bytes struct {
	rows *sql.Rows
	err  error
	conv Rows2Byte
}

func NewSqlBytesSource(rows *sql.Rows, conv Rows2Byte) *Rows2Bytes {
	var err error = nil
	return &Rows2Bytes{
		rows,
		err,
		conv,
	}
}

func (r *Rows2Bytes) AsIf() BytesSourceCloser { return r }

func (r *Rows2Bytes) Scan() bool   { return r.rows.Next() }
func (r *Rows2Bytes) Close() error { return r.rows.Close() }

func (r *Rows2Bytes) Bytes() []byte {
	var ret []byte
	ret, r.err = r.conv(r.rows)
	return ret
}

func (r *Rows2Bytes) Err() error { return r.err }
