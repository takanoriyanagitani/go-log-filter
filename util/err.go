package util

func ifOk[T any](valid bool, ok func() (T, error), ng func() (T, error)) (T, error) {
	switch valid {
	case true:
		return ok()
	default:
		return ng()
	}
}
