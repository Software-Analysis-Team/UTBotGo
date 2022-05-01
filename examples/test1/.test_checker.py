#!/usr/bin/env python3
import json

with open('utbotgo/functions/getSign/args.json', 'r') as f:
    args = json.load(f)
    values = [0] * len(args)
    for i, params in enumerate(args):
        assert list(params.keys()) == ['x']
        values[i] = params['x']
    values.sort()
    assert values[0] < 0 and values[1] == 0 and values[2] > 0
