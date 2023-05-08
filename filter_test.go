package filter_test

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	lf "github.com/takanoriyanagitani/go-log-filter"
)

func assertNil(e error) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()

		if nil == e {
			return
		}

		t.Fatalf("Unexpected error: %v\n", e)
	}
}

func mustNil(e error) {
	if nil == e {
		return
	}

	panic(e)
}

func assertEqNew[T any](comp func(a, b T) (same bool)) func(a, b T) func(*testing.T) {
	return func(a, b T) func(*testing.T) {
		return func(t *testing.T) {
			t.Helper()
			var same bool = comp(a, b)
			if same {
				return
			}

			t.Errorf("Unexpected value got.\n")
			t.Errorf("Expected: %v\n", b)
			t.Fatalf("Got:      %v\n", a)
		}
	}
}

func assertEq[T comparable](a, b T) func(*testing.T) {
	comp := func(a, b T) (same bool) { return a == b }
	return assertEqNew(comp)(a, b)
}

func TestFilter(t *testing.T) {
	t.Parallel()

	t.Run("BytesProcessorNew", func(t *testing.T) {
		t.Parallel()

		var s lf.SkipBytes = lf.NopBytesSkipper

		var cnt uint64 = 0
		var c lf.BytesConsumer = lf.BytesCounterNew(func(b []byte) { cnt += uint64(len(b)) })

		var p lf.BytesProcessor = lf.BytesProcessorNew(s)(c)

		var buf bytes.Buffer
		var wtr io.Writer = &buf

		var e error

		_, e = fmt.Fprintln(wtr, "12")
		mustNil(e)
		_, e = fmt.Fprintln(wtr, "34")
		mustNil(e)

		var rdr io.Reader = &buf
		var src lf.BytesSource = lf.NewBytesSource(rdr)

		e = p(src)

		t.Run("no error", assertNil(e))
		t.Run("count check", assertEq(4, cnt))
	})
}
