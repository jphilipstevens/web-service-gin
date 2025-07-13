#!/usr/bin/env bash
set -euo pipefail

# Overwrite version.txt with the provided version number.
# Usage: set_version.sh <version>

version="$1"

if [[ -z "$version" ]]; then
  echo "Version is required" >&2
  exit 1
fi

echo "$version" > version.txt
