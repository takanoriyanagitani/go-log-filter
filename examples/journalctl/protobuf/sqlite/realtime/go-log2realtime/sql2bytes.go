package main

import (
	"context"
	"database/sql"

	lf "github.com/takanoriyanagitani/go-log-filter"
	sqlsrc "github.com/takanoriyanagitani/go-log-filter/data/source/sql"
)

type bytesSourceBuilder struct {
	db   *sql.DB
	conv sqlsrc.Rows2Byte
}

func newBytesSourceBuilder(db *sql.DB, conv sqlsrc.Rows2Byte) bytesSourceBuilder {
	return bytesSourceBuilder{db, conv}
}

func (b bytesSourceBuilder) Close() error { return b.db.Close() }

func (b bytesSourceBuilder) ToBytesSource(
	ctx context.Context,
	query string,
	args ...any,
) sqlsrc.BytesSourceCloser {
	rows, e := b.db.QueryContext(ctx, query, args...)
	if nil != e {
		return lf.BytesSourceErrNew(e)
	}
	if nil != rows.Err() {
		return lf.BytesSourceErrNew(rows.Err())
	}
	return sqlsrc.NewSqlBytesSource(rows, b.conv)
}
