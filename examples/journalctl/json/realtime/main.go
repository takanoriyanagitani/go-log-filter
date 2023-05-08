package main

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"strconv"

	jsonit "github.com/json-iterator/go"

	lf "github.com/takanoriyanagitani/go-log-filter"
)

var json = jsonit.ConfigCompatibleWithStandardLibrary

type Raw struct {
	Realtime string `json:"__REALTIME_TIMESTAMP"`
}

type Log struct {
	Realtime float64
}

func fastTransformer() lf.Transform[Log] {
	var buf Raw
	return func(b []byte) (l Log, e error) {
		ejson := json.Unmarshal(b, &buf)
		realtime, eint := strconv.ParseInt(buf.Realtime, 10, 64)
		l.Realtime = float64(realtime) * 1e-6
		return l, errors.Join(ejson, eint)
	}
}

func fastConsumer(w io.Writer) lf.ConsumeTransformed[Log] {
	var bo binary.ByteOrder = binary.BigEndian
	return func(l Log) error {
		return binary.Write(w, bo, &l.Realtime)
	}
}

func fastSkip() lf.SkipTransformed[Log] { return func(l Log) (skip bool) { return false } }

func main() {
	var w io.Writer = os.Stdout
	var bw *bufio.Writer = bufio.NewWriter(w)

	var st lf.SkipTransformed[Log] = fastSkip()
	var ct lf.ConsumeTransformed[Log] = fastConsumer(bw)

	var tf lf.Transform[Log] = fastTransformer()
	var bp lf.BytesProcessor = tf.ProcessorNew(st)(ct)

	var r io.Reader = os.Stdin
	var bs lf.BytesSource = lf.NewBytesSource(r)

	e := bp(bs)
	if nil != e {
		panic(e)
	}

	e = bw.Flush()
	if nil != e {
		panic(e)
	}
}
