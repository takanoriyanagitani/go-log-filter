package main

import (
	"bufio"
	"context"
	"encoding/binary"
	"io"
	"os"
	"strconv"

	"google.golang.org/protobuf/proto"

	lf "github.com/takanoriyanagitani/go-log-filter"
	sqlsrc "github.com/takanoriyanagitani/go-log-filter/data/source/sql"
	lg "github.com/takanoriyanagitani/go-log-filter/examples/journalctl/protobuf/sqlite/realtime/go-log2realtime/protobuf/journalctl/easymap"
	util "github.com/takanoriyanagitani/go-log-filter/util"
)

func bytes2log(b []byte, l *lg.Log) error {
	l.Reset()
	return proto.Unmarshal(b, l)
}

func bytes2logNew() lf.Transform[*lg.Log] {
	buf := lg.Log{}
	return func(serialized []byte) (*lg.Log, error) {
		e := bytes2log(serialized, &buf)
		return &buf, e
	}
}

var skipLog lf.SkipTransformed[*lg.Log] = func(l *lg.Log) (skip bool) {
	var items map[string]string = l.Items
	_, found := items["__REALTIME_TIMESTAMP"]
	return !found
}

var string2float func(bits int, s string) (float64, error) = util.Swap(strconv.ParseFloat)
var str2f6 func(string) (float64, error) = util.Curry(string2float)(64)

var str2realtime func(string) (float64, error) = util.Compose(
	str2f6,
	func(micros float64) (float64, error) { return 1e-6 * micros, nil },
)

func logConsumerNew(w io.Writer) lf.ConsumeTransformed[*lg.Log] {
	float2writer := func(f float64) (int, error) { return 8, binary.Write(w, binary.BigEndian, f) }
	realtime2writer := util.Compose(
		str2realtime,
		float2writer,
	)
	return func(l *lg.Log) error {
		var items map[string]string = l.Items
		var realtime string = items["__REALTIME_TIMESTAMP"]
		_, e := realtime2writer(realtime)
		return e
	}
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
	var log2consumer lf.ConsumeTransformed[*lg.Log] = logConsumerNew(w)
	var b2l lf.Transform[*lg.Log] = bytes2logNew()
	var proc lf.BytesProcessor = b2l.ProcessorNew(skipLog)(log2consumer)
	e = proc(bsc)
	if nil != e {
		panic(e)
	}
	e = w.Flush()
	if nil != e {
		panic(e)
	}
}
