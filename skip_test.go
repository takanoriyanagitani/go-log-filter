package filter_test

import (
	"testing"

	lf "github.com/takanoriyanagitani/go-log-filter"
)

func TestSkip(t *testing.T) {

	t.Parallel()

	t.Run("SkipByLenRangeNew", func(t *testing.T) {
		t.Parallel()

		var s lf.SkipBytes = lf.SkipByLenRangeNew(1, 7)

		tf := func(input []byte, expected bool) func(*testing.T) {
			var got bool = s(input)
			return assertEq(got, expected)
		}

		t.Run("zero", tf(nil, true))
		t.Run("single", tf([]byte("1"), false))
		t.Run("max", tf([]byte("123456"), false))
		t.Run("too many bytes", tf([]byte("1234567"), true))
	})

	t.Run("SkipBytes", func(t *testing.T) {
		t.Parallel()
		t.Run("And", func(t *testing.T) {
			t.Parallel()

			var zero lf.SkipBytes = lf.SkipByLenNew(func(l int) (skip bool) { return 0 == l })
			var non0 lf.SkipBytes = zero.Not()
			var chk1 lf.SkipBytes = func(b []byte) (skip bool) { return '{' == b[0] }
			var chk lf.SkipBytes = non0.And(chk1).Not()

			t.Run("empty", assertEq(chk(nil), true))
			t.Run("single", assertEq(chk([]byte("}")), true))
			t.Run("single", assertEq(chk([]byte("[")), true))
			t.Run("single", assertEq(chk([]byte("]")), true))
			t.Run("single", assertEq(chk([]byte("{")), false))
		})
	})

}
