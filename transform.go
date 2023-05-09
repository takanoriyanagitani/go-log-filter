package filter

import (
	"github.com/takanoriyanagitani/go-log-filter/util"
)

type Transform[T any] func([]byte) (T, error)
type SkipTransformed[T any] func(T) (skip bool)
type ConsumeTransformed[T any] func(T) error

func (s SkipTransformed[T]) Not() SkipTransformed[T] { return func(t T) (skip bool) { return !s(t) } }

func (t Transform[T]) BytesConsumerNew(consumer ConsumeTransformed[T]) BytesConsumer {
	return func(b []byte) error {
		converted, e := t(b)
		if nil != e {
			return e
		}
		return consumer(converted)
	}
}

func (t Transform[T]) SkipBytesNew(skip SkipTransformed[T]) SkipBytes {
	return func(b []byte) bool {
		converted, e := t(b)
		if nil != e {
			return true
		}
		return skip(converted)
	}
}

func (t Transform[T]) ProcessorNew(s SkipTransformed[T]) func(ConsumeTransformed[T]) BytesProcessor {
	return func(c ConsumeTransformed[T]) BytesProcessor {
		var sb SkipBytes = t.SkipBytesNew(s)
		var bc BytesConsumer = t.BytesConsumerNew(c)
		return BytesProcessorNew(sb)(bc)
	}
}

func TransformAdd[T, U any](original Transform[T], conv func(T) (U, error)) Transform[U] {
	return util.Compose(
		original,
		conv,
	)
}
