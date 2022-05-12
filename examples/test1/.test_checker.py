#!/usr/bin/env python3
import json

with open('utbotgo/functions/getSign/args.json', 'r') as f:
    args = json.load(f)
    assert len(args) == 3
    values = [0] * len(args)
    for i, params in enumerate(args):
        assert list(params.keys()) == ['x']
        values[i] = (params['x'], i)
    values.sort()
    assert values[0][0] < 0 and values[1][0] == 0 and values[2][0] > 0

with open('utbotgo/functions/getSign/results.json', 'w') as f:
    answers = [0] * len(args)
    for value, i in values:
        if value == 0:
            answer = 0
        elif value < 0:
            answer = -1
        else:
            answer = 1
        answers[i] = {'R': answer}
    f.write(json.dumps(answers))
