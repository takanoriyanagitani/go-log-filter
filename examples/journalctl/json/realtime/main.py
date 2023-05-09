import json
import sys
import functools
import operator
import struct

partial = functools.partial
itemget = operator.itemgetter

reducer = lambda state, f: f(state)
swap = lambda f: lambda a,b: f(b,a)
has = swap(operator.contains)

s = struct.Struct(">d")

functools.reduce(
	reducer,
	[
		partial(map, json.loads),
		partial(filter, partial(has, "__REALTIME_TIMESTAMP")),
		partial(map, itemget("__REALTIME_TIMESTAMP")),
		partial(map, float),
		partial(map, partial(operator.mul, 1e-6)),
		partial(map, s.pack),
		partial(map, sys.stdout.buffer.write),
		lambda writes: sum(1 for _ in writes),
	],
	sys.stdin,
)
