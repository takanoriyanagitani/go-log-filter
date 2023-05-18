package util

func ifOk[T any](valid bool, ok func() (T, error), ng func() (T, error)) (T, error) {
	var f func() (T, error) = Select(ng, ok, valid)
	return f()
}

func Select[T any](f T, t T, cond bool) T {
	switch cond {
	case true:
		return t
	default:
		return f
	}
}
