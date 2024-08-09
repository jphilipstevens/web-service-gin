#!/bin/bash

# Function to determine the new version based on commit types
function determine_new_version() {
  local current_version=$(cat version.txt)
  local major=$(echo "$current_version" | cut -d'.' -f1)
  local minor=$(echo "$current_version" | cut -d'.' -f2)
  local patch=$(echo "$current_version" | cut -d'.' -f3)

  local has_breaking_change=$(git log --pretty=format:'%s' --since="v$current_version" | grep -E 'BREAKING CHANGE')

  if [ ! -z "$has_breaking_change" ]; then
    ((major++))
    minor=0
    patch=0
  elif git log --pretty=format:'%s' --since="v$current_version" | grep -E '^feat' | wc -l > 0; then
    ((minor++))
    patch=0
  elif git log --pretty=format:'%s' --since="v$current_version" | grep -E '^fix' | wc -l > 0; then
    ((patch++))
  fi

  echo "${major}.${minor}.${patch}"
}

# Update the version file
new_version=$(determine_new_version)
echo "$new_version" > version.txt
