#!/bin/bash
VERSION=$(cat version.txt)
IFS='.' read -ra VERSIONS <<< "$VERSION"

case "$1" in
  major)
    VERSIONS[0]=$((VERSIONS[0] + 1))
    VERSIONS[1]=0
    VERSIONS[2]=0
    ;;
  minor)
    VERSIONS[1]=$((VERSIONS[1] + 1))
    VERSIONS[2]=0
    ;;
  patch)
    VERSIONS[2]=$((VERSIONS[2] + 1))
    ;;
  *)
    echo "Invalid version bump type"
    exit 1
    ;;
esac

NEW_VERSION="${VERSIONS[0]}.${VERSIONS[1]}.${VERSIONS[2]}"
echo "Bumped version to $NEW_VERSION"
echo "$NEW_VERSION" > version.txt
