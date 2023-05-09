package main

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"math"
	"os"
	"runtime"
	"strconv"

	"github.com/bytecodealliance/wasmtime-go/v8"

	jsonit "github.com/json-iterator/go"

	lf "github.com/takanoriyanagitani/go-log-filter"
	"github.com/takanoriyanagitani/go-log-filter/util"
)

var json = jsonit.ConfigCompatibleWithStandardLibrary

var ErrInvalidRealtime = errors.New("invalid realtime")
var ErrInvalidMemory = errors.New("invalid memory")
var ErrInvalidModule = errors.New("invalid module")
var ErrInvalidNumber = errors.New("invalid number")

type wtInstance struct {
	i *wasmtime.Instance
	f *wasmtime.Func
	s wasmtime.Storelike
	a int32
}

type wtModule struct {
	m *wasmtime.Module
	s wasmtime.Storelike
}

type wtEngine struct{ e *wasmtime.Engine }

func (e wtEngine) newModule(wasm []byte) (*wasmtime.Module, error) {
	return wasmtime.NewModule(e.e, wasm)
}
func (e wtEngine) newStore() *wasmtime.Store { return wasmtime.NewStore(e.e) }

func (e wtEngine) toModule(wasm []byte) (wtModule, error) {
	var s *wasmtime.Store = e.newStore()
	m, err := e.newModule(wasm)
	return wtModule{m, s}, err
}

func wasm2module(wasm []byte) (wtModule, error) {
	var e wtEngine = wtEngine{e: wasmtime.NewEngine()}
	return e.toModule(wasm)
}

var filename2module func(filename string) (wtModule, error) = util.Compose(
	os.ReadFile,
	wasm2module,
)

func (m wtModule) newInstance() (*wasmtime.Instance, error) {
	return wasmtime.NewInstance(m.s, m.m, nil)
}

func (m wtModule) getAddr(f *wasmtime.Func) (int32, error) {
	return util.Compose(
		func(s wasmtime.Storelike) (any, error) { return f.Call(s) },
		func(i any) (int32, error) {
			switch addr := i.(type) {
			case int32:
				return addr, nil
			default:
				return -1, ErrInvalidModule
			}
		},
	)(m.s)
}

func (m wtModule) toInstance(name string) (wtInstance, error) {
	i, e := m.newInstance()
	if nil != e {
		return wtInstance{}, e
	}
	var f *wasmtime.Func = i.GetFunc(m.s, name)
	if nil == f {
		return wtInstance{}, ErrInvalidModule
	}
	var addrFn *wasmtime.Func = i.GetFunc(m.s, "addr")
	if nil == addrFn {
		return wtInstance{}, ErrInvalidModule
	}
	a, e := m.getAddr(addrFn)
	if nil != e {
		return wtInstance{}, e
	}
	var s wasmtime.Storelike = m.s
	return wtInstance{
		i,
		f,
		s,
		a,
	}, nil
}

func (m wtModule) toTransformer(name string) (lf.Transform[float64], error) {
	i, e := m.toInstance(name)
	return i.asTransform(), e
}

func (wi wtInstance) any2f64(i any) (float64, error) {
	switch f := i.(type) {
	case float64:
		return f, nil
	default:
		return 0.0, ErrInvalidRealtime
	}
}

func (wi wtInstance) realtime(size int) (float64, error) {
	var sz int32 = int32(size) & 0xffff
	return util.Compose(
		func(s wasmtime.Storelike) (any, error) { return wi.f.Call(s, sz) },
		util.Compose(
			wi.any2f64,
			wi.f2f,
		),
	)(wi.s)
}

func (wi wtInstance) f2f(f float64) (float64, error) {
	switch math.IsNaN(f) {
	case true:
		return math.NaN(), ErrInvalidNumber
	default:
		return f, nil
	}
}

func (wi wtInstance) toRealtime(raw []byte) (float64, error) {
	var ext *wasmtime.Extern = wi.i.GetExport(wi.s, "memory")
	if nil == ext {
		return 0.0, ErrInvalidMemory
	}
	var mem *wasmtime.Memory = ext.Memory()
	if nil == mem {
		return 0.0, ErrInvalidMemory
	}
	var max int = len(raw) & 0xffff
	var dst []byte = mem.UnsafeData(wi.s)
	var offset int = int(wi.a)
	copy(dst[offset:], raw[:max])
	f, e := wi.realtime(len(raw))
	runtime.KeepAlive(mem)
	return f, e
}

func (wi wtInstance) asTransform() lf.Transform[float64] { return wi.toRealtime }

type Raw struct {
	Realtime string `json:"__REALTIME_TIMESTAMP"`
}

type Log struct {
	Realtime float64
}

func fastTransformer() lf.Transform[Log] {
	var buf Raw
	return func(b []byte) (l Log, e error) {
		ejson := json.Unmarshal(b, &buf)
		realtime, eint := strconv.ParseInt(buf.Realtime, 10, 64)
		l.Realtime = float64(realtime) * 1e-6
		return l, errors.Join(ejson, eint)
	}
}

func fastConsumer(w io.Writer) lf.ConsumeTransformed[Log] {
	var bo binary.ByteOrder = binary.BigEndian
	return func(l Log) error {
		return binary.Write(w, bo, &l.Realtime)
	}
}

func fastSkip() lf.SkipTransformed[Log] { return func(l Log) (skip bool) { return false } }

func native() {
	var w io.Writer = os.Stdout
	var bw *bufio.Writer = bufio.NewWriter(w)

	var st lf.SkipTransformed[Log] = fastSkip()
	var ct lf.ConsumeTransformed[Log] = fastConsumer(bw)

	var tf lf.Transform[Log] = fastTransformer()
	var bp lf.BytesProcessor = tf.ProcessorNew(st)(ct)

	var r io.Reader = os.Stdin
	var bs lf.BytesSource = lf.NewBytesSource(r)

	e := bp(bs)
	if nil != e {
		panic(e)
	}

	e = bw.Flush()
	if nil != e {
		panic(e)
	}
}

func wasm() {
	var w io.Writer = os.Stdout
	var bw *bufio.Writer = bufio.NewWriter(w)

	var st lf.SkipTransformed[Log] = fastSkip()
	var ct lf.ConsumeTransformed[Log] = fastConsumer(bw)

	var wasmname string = "./j2realtime.wasm"
	m, e := filename2module(wasmname)
	if nil != e {
		panic(e)
	}
	t, e := m.toTransformer("to_real")
	if nil != e {
		panic(e)
	}
	var tf lf.Transform[Log] = lf.TransformAdd(t, func(rt float64) (Log, error) {
		return Log{Realtime: rt}, nil
	})

	var bp lf.BytesProcessor = tf.ProcessorNew(st)(ct)

	var r io.Reader = os.Stdin
	var bs lf.BytesSource = lf.NewBytesSource(r)

	e = bp(bs)
	if nil != e {
		panic(e)
	}

	e = bw.Flush()
	if nil != e {
		panic(e)
	}
}

func main() {
	var useWasm string = os.Getenv("USE_WASM")
	switch useWasm {
	case "use":
		wasm()
	default:
		native()
	}
}
