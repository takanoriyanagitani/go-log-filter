package util

func Swap[T, U, V any](f func(t T, u U) (V, error)) func(U, T) (V, error) {
	return func(u U, t T) (V, error) { return f(t, u) }
}
