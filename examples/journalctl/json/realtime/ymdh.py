import json
import sys
import functools
import operator
import time
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
		partial(filter, lambda f: 1683493200.0 <= f <= 1683496800.0),
		partial(map, time.localtime),
		partial(map, lambda t: dict(ymdhm=[t.tm_year, t.tm_mon, t.tm_mday, t.tm_hour, t.tm_min])),
		partial(map, print),
		lambda prints: sum(1 for _ in prints),
	],
	sys.stdin,
)
