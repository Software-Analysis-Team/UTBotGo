#!/usr/bin/python3
import json
import sys

with open(sys.argv[1], 'r') as f:
    go_config = json.load(f)
with open(sys.argv[2], 'r') as f:
    file_templates = json.load(f)
with open(sys.argv[3], 'r') as f:
    ta_types = json.load(f)
files = dict()
for name, template in file_templates.items():
    files[name] = {"name": template.replace('%', ta_types["name"])}
print(json.dumps({"package": {"name": go_config["Name"]}, "files": files, "func": ta_types}))
