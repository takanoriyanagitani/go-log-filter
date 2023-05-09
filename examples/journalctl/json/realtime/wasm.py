import functools
import json
import operator

from wasmtime import Store, Module, Instance, Func, FuncType

curry = lambda f: lambda x: lambda y: f(x,y)

reducer = lambda state, f: f(state)

replace = lambda neo: lambda _: neo

s = Store()
wasmname2module = curry(Module.from_file)(s.engine)
m = wasmname2module("./j2realtime.wasm")
module2instance = curry(lambda store, module: Instance(store, module, []))(s)
i = module2instance(m)
e = i.exports(s)
to_real = e["to_real"]
addr = e["addr"](s)
esize = e["esize"]
eaddr = e["eaddr"](s)
mem = e["memory"]

functools.reduce(
	reducer,
	[
		operator.methodcaller("encode", "UTF8"),
		lambda b: mem.write(s, b, addr),
		functools.partial(to_real, s),
		functools.partial(operator.mul, 1e6),
		print,
		replace(s),
		esize,
		lambda elen: mem.read(s, eaddr, eaddr+elen),
		print,
	],
	json.dumps(dict(
		__REALTIME_TIMESTAMP = "634",
	)),
)
