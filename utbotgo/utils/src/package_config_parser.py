#!/usr/bin/python3
import json
import os.path
import sys

def extract_data(config, fields):
    for field in fields:
        config = config[field]
    if type(config) == str:
        return config
    elif type(config) == list:
        return ' '.join(config)

def write_if_necessary(file, new_data):
    if os.path.exists(file):
        with open(file, 'r') as f:
            old_data = f.read()
    else:
        old_data = None
    if old_data != new_data:
        os.makedirs(os.path.dirname(file), exist_ok=True)
        with open(file, 'w') as f:
            f.write(new_data)

config = json.load(sys.stdin)
amt_files = (len(sys.argv) - 1) // 2
for i in range(amt_files):
    fields, file = sys.argv[i*2+1].split('.'), sys.argv[i*2+2]
    data = extract_data(config, fields)
    write_if_necessary(file, data)
