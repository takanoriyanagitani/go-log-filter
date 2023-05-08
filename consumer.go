package filter

type BytesConsumer func([]byte) error

func BytesCounterNew(countBytes func([]byte)) BytesConsumer {
	return func(b []byte) error {
		countBytes(b)
		return nil
	}
}
