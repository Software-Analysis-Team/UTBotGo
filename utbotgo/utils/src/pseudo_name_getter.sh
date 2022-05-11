#!/bin/bash
set -e
pseudo_name=$(cat "$1")
if [[ "$pseudo_name" != 'main' ]]; then
  pseudo_name=$(cat "$2")
fi
if ! diff <(echo "$pseudo_name") <(cat "$3" 2>&1 || echo "") > /dev/null; then
  echo "$pseudo_name" > "$3"
fi
