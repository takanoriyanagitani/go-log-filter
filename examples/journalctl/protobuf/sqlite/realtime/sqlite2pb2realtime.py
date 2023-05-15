import functools
import operator
import sqlite3
import struct
import sys

import jtyp_pb2

reducer = lambda state, f: f(state)
curry = lambda f: lambda x: lambda y: f(x,y)
compose = lambda f, g: lambda x: g(f(x))

def float2bytes(s, f):
	return s.pack(f)

def bytes2log(l, b):
	l.Clear()
	l.ParseFromString(b)
	return l

bytes2logNew = curry(bytes2log)
b2l = bytes2logNew(jtyp_pb2.Log())
log2items = operator.attrgetter("items")
items2rtime = operator.itemgetter("__REALTIME_TIMESTAMP")
rtime2s = lambda s: float(s) * 1e-6

float2bytesNew = curry(float2bytes)
f2b = float2bytesNew(struct.Struct(">d"))

serialized2items = compose(b2l, log2items)
items2realtime = compose(items2rtime, rtime2s)
items2bytes = compose(items2realtime, f2b)
serialized2bytes = compose(serialized2items, items2bytes)

with sqlite3.connect("file:journalctl.sqlite3.db?mode=ro", uri=True) as con:
	functools.reduce(
		reducer,
		[
			functools.partial(map, operator.itemgetter(0)),
			functools.partial(map, serialized2bytes),
			functools.partial(map, sys.stdout.buffer.write),
			lambda writes: sum(1 for _ in writes),
		],
		con.execute('''
			SELECT val FROM logs
			ORDER BY key
			LIMIT 1048576
		'''),
	)
	pass
