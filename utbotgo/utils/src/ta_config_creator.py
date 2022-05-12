#!/usr/bin/python3
import json
import sys

with open(sys.argv[1], 'r') as f:
    package_name = f.read().strip()
with open(sys.argv[2], 'r') as f:
    file_templates = json.load(f)
with open(sys.argv[3], 'r') as f:
    ta_types = json.load(f)
with open(sys.argv[4], 'r') as f:
    raw_args = json.load(f)
if len(sys.argv) == 6:
    with open(sys.argv[5], 'r') as f:
        raw_answers = json.load(f)
else:
    raw_answers = []
files = dict()
for name, template in file_templates.items():
    files[name] = {"name": template.replace('%', ta_types["name"])}
args = [dict() for _ in range(len(raw_args))]
for i, raw_arg in enumerate(raw_args):
    for key, value in raw_arg.items():
        args[i][key] = str(value)
answers = [dict() for _ in range(len(raw_answers))]
for i, raw_answer in enumerate(raw_answers):
    for key, value in raw_answer.items():
        answers[i][key] = str(value) if type(value) != str else f'"{value}"'
print(json.dumps({"package": {"name": package_name}, "files": files, "func": ta_types, "args": args, "answers": answers}))
