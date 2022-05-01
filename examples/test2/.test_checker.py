#!/usr/bin/env python3
import json

with open('utbotgo/functions/addInt/args.json', 'r') as f:
    args = json.load(f)
    assert len(args) == 1
    assert list(args[0].keys()) == ['x', 'y']
    assert type(args[0]['x']) == int and type(args[0]['y']) == int

with open('utbotgo/functions/multInt/args.json', 'r') as f:
    args = json.load(f)
    assert len(args) == 1
    assert list(args[0].keys()) == ['x', 'y']
    assert type(args[0]['x']) == int and type(args[0]['y']) == int
