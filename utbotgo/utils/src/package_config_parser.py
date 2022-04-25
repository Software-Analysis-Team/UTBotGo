#!/usr/bin/python3
import json
import sys

field = sys.argv[1]
value = json.load(sys.stdin)[field]
if type(value) == str:
    print(value)
elif type(value) == list:
    print(' '.join(value))
else:
    print(value)
