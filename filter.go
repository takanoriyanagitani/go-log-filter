package filter

type BytesProcessor func(BytesSource) error

func BytesProcessorNew(s SkipBytes) func(BytesConsumer) BytesProcessor {
	return func(c BytesConsumer) BytesProcessor {
		return func(src BytesSource) error {
			for src.Scan() {
				var b []byte = src.Bytes()
				var skip bool = s(b)
				if skip {
					continue
				}
				var e error = c(b)
				if nil != e {
					return e
				}
			}
			return nil
		}
	}
}
