#!/usr/bin/env bash

BINARY=git-auto

set -e
rm -rf completions
mkdir completions

for sh in bash zsh fish; do
  go run main.go completion "$sh" >"completions/$BINARY.$sh"
done
