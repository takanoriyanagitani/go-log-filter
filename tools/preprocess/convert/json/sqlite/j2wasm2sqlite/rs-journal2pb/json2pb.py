import json
import functools
import operator
import sys
import sqlite3

import jtyp_pb2

reducer = lambda state, f: f(state)

curry = lambda f: lambda x: lambda y: f(x, y)

def items2sqlite(conn, items):
    ret = conn.executemany(
        '''
            INSERT INTO logs(val)
            VALUES(:val)
        ''',
        items,
    )
    return ret

def json2log(l, j):
    l.Clear()
    l.items.update(j)
    return l

json2logNew = curry(json2log)
j2l = json2logNew(jtyp_pb2.Log())
l2s = lambda l: l.SerializeToString()

items2sqliteNew = curry(items2sqlite)

bytes2dict = lambda b: dict(val=b)

with sqlite3.connect("./journalctl.sqlite3.db") as con:
    cur = con.execute('''
        DROP TABLE IF EXISTS logs
    ''')
    cur = con.execute('''
        CREATE TABLE IF NOT EXISTS logs(
            key INTEGER PRIMARY KEY,
            val BLOB NOT NULL
        )
    ''')
    i2s = items2sqliteNew(cur)
    functools.reduce(
        reducer,
        [
            functools.partial(map, json.loads),
            functools.partial(map, j2l),
            functools.partial(map, l2s),
            functools.partial(map, bytes2dict),
            lambda items: i2s(items),
            print,
        ],
        sys.stdin,
    )
    pass
