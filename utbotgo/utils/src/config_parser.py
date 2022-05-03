#!/usr/bin/python3
import yaml
import sys
import subprocess

config_file_path = "utbotgo/config.yml"
try:
    with open(config_file_path, 'r') as f:
        config = yaml.safe_load(f)
        functions = set(filter(lambda x: type(x) == str, config["tested-functions"]))
except BaseException as err:
    print(err)
    sys.exit(1)

try:
    subprocess.run(["make", "-f", "/utbotgo/workspace_files/Makefile", f"FNS={' '.join(functions)}", sys.argv[1]], check=True)
except subprocess.CalledProcessError as exc:
    sys.exit(exc.returncode)
