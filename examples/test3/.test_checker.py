#!/usr/bin/env python3
import json

with open('utbotgo/functions/g/args.json', 'r') as f:
    args = json.load(f)
    assert len(args) == 4
    cases = [False] * len(args)
    for params in args:
        assert list(params.keys()) == ['a', 'b', 'c']
        a = params['a']
        b = params['b']
        c = params['c']
        assert type(a) == type(b) == type(c) == int
        if a < b and b < c:
            cases[0] = True
        elif a < b and b >= c:
            cases[1] = True
        elif a >= b and b < c:
            cases[2] = True
        elif a >= b and b >= c:
            cases[3] = True
    assert cases[0] and cases[1] and cases[2] and cases[3]

with open('utbotgo/functions/g/answers.json', 'w') as f:
    answers = [0] * len(args)
    for i in range(len(args)):
        answers[i] = {'d': max(args[i]['a'], args[i]['b'], args[i]['c'])}
    f.write(json.dumps(answers))
