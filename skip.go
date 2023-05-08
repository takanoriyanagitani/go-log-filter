package filter

type SkipBytes func([]byte) (skip bool)

func (s SkipBytes) And(other SkipBytes) SkipBytes {
	return func(b []byte) (skip bool) {
		return s(b) && other(b)
	}
}

func (s SkipBytes) Not() SkipBytes {
	return func(b []byte) (skip bool) {
		return !s(b)
	}
}

func SkipBytesNewStatic(staticSkip bool) SkipBytes {
	return func(_ []byte) (skip bool) { return staticSkip }
}

var NopBytesSkipper SkipBytes = SkipBytesNewStatic(false)

func SkipByLenNew(check func(int) (skip bool)) SkipBytes {
	return func(b []byte) bool {
		var l int = len(b)
		return check(l)
	}
}

func SkipByLenRangeNew(lbi int, ube int) SkipBytes {
	var s1 SkipBytes = SkipByLenNew(func(l int) (skip bool) { return lbi <= l })
	var s2 SkipBytes = SkipByLenNew(func(l int) (skip bool) { return l < ube })
	return s1.And(s2).Not()
}
