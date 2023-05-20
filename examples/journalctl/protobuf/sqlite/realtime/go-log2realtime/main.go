package main

import (
	"bufio"
	"context"
	"encoding/binary"
	"os"
	"strconv"

	"google.golang.org/protobuf/proto"

	sqlsrc "github.com/takanoriyanagitani/go-log-filter/data/source/sql"
	lg "github.com/takanoriyanagitani/go-log-filter/examples/journalctl/protobuf/sqlite/realtime/go-log2realtime/protobuf/journalctl/easymap"
)

func bytes2log(b []byte, l *lg.Log) error {
	l.Reset()
	return proto.Unmarshal(b, l)
}

func main() {
	db, e := filename2sqlite("./journalctl.sqlite3.db")
	if nil != e {
		panic(e)
	}
	var b bytesSourceBuilder = newBytesSourceBuilder(db, sqlsrc.Val2Byte)
	defer b.Close()

	var bsc sqlsrc.BytesSourceCloser = b.ToBytesSource(
		context.Background(),
		`
			SELECT val FROM logs
			ORDER BY key
			LIMIT 1048576
		`,
	)
	defer bsc.Close()

	var w *bufio.Writer = bufio.NewWriter(os.Stdout)
	var l lg.Log = lg.Log{
		Items: map[string]string{},
	}
	for bsc.Scan() {
		var b []byte = bsc.Bytes()
		e := bytes2log(b, &l)
		if nil != e {
			panic(e)
		}
		var items map[string]string = l.Items
		realtime, found := items["__REALTIME_TIMESTAMP"]
		if !found {
			continue
		}
		parsed, e := strconv.ParseFloat(realtime, 64)
		if nil != e {
			panic(e)
		}
		var rtus float64 = 1e-6 * parsed
		e = binary.Write(w, binary.BigEndian, rtus)
		if nil != e {
			panic(e)
		}
	}
	e = bsc.Err()
	if nil != e {
		panic(e)
	}
	e = w.Flush()
	if nil != e {
		panic(e)
	}
}
