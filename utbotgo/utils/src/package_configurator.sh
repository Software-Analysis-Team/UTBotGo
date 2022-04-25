#!/bin/bash
set -e
package_config=$(go list -json)
if ! diff <(echo "$package_config") <(cat "$1" 2>&1 || echo "{}") > /dev/null; then
  echo "$package_config" > "$1"
fi
