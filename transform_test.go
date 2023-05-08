package filter_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"testing"

	lf "github.com/takanoriyanagitani/go-log-filter"
)

func TestTransform(t *testing.T) {

	t.Parallel()

	t.Run("ProcessorNew", func(t *testing.T) {
		t.Parallel()

		var st lf.SkipTransformed[map[string]string] = func(m map[string]string) (skip bool) {
			_, found := m["__REALTIME_TIMESTAMP"]
			return !found
		}

		var cnt int = 0
		var ct lf.ConsumeTransformed[map[string]string] = func(m map[string]string) error {
			stime := m["__REALTIME_TIMESTAMP"]
			_, e := strconv.ParseInt(stime, 10, 64)
			if nil == e {
				cnt += 1
			}
			return e
		}

		var tf lf.Transform[map[string]string] = func(b []byte) (map[string]string, error) {
			var m map[string]string
			var e error = json.Unmarshal(b, &m)
			return m, e
		}

		var bp lf.BytesProcessor = tf.ProcessorNew(st)(ct)

		var buf bytes.Buffer
		var wtr io.Writer = &buf

		var e error
		_, e = fmt.Fprintln(wtr, `{"__REALTIME_TIMESTAMP": "1683489600120173"}`)
		mustNil(e)

		_, e = fmt.Fprintln(wtr, `{"__REALTIME_TIMESTAMP": "9999999999999999"}`)
		mustNil(e)

		_, e = fmt.Fprintln(wtr, `{"__REALTIME_TIMESTAMP": "zzzzzzzzzzzzzzzz"}`)
		mustNil(e)

		_, e = fmt.Fprintln(wtr, `{"__MONOTIME_TIMESTAMP": "1683489600120173"}`)
		mustNil(e)

		_, e = fmt.Fprintln(wtr, `{[`)
		mustNil(e)

		var rdr io.Reader = &buf
		var bs lf.BytesSource = lf.NewBytesSource(rdr)

		e = bp(bs)
		t.Run("error", assertEq(false, nil == e))
		t.Run("2 valid timestamps", assertEq(cnt, 2))
	})

}
