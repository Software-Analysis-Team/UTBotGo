#!/usr/bin/env python3
import json

with open('utbotgo/functions/addInt/args.json', 'r') as f:
    args = json.load(f)
    assert len(args) == 1
    arg = args[0]
    assert list(arg.keys()) == ['x', 'y']
    assert type(arg['x']) == int and type(arg['y']) == int

with open('utbotgo/functions/addInt/results.json', 'w') as f:
    f.write(json.dumps([{'R': arg['x'] + arg['y']}]))

with open('utbotgo/functions/multInt/args.json', 'r') as f:
    args = json.load(f)
    assert len(args) == 1
    arg = args[0]
    assert list(arg.keys()) == ['x', 'y']
    assert type(arg['x']) == int and type(arg['y']) == int

with open('utbotgo/functions/multInt/results.json', 'w') as f:
    f.write(json.dumps([{'R': arg['x'] * arg['y']}]))
