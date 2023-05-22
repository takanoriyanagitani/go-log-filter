package util

func Curry[T, U, V any](f func(T, U) (V, error)) func(T) func(U) (V, error) {
	return func(t T) func(U) (V, error) {
		return func(u U) (V, error) {
			return f(t, u)
		}
	}
}
