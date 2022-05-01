#!/usr/bin/python3
import json
import sys

fields = sys.argv[1].split('.')
config = json.load(sys.stdin)
for field in fields:
    config = config[field]
if type(config) == str:
    print(config)
elif type(config) == list:
    print(' '.join(config))
