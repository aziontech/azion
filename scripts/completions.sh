#!/bin/sh
# scripts/completions.sh
set -e
rm -rf completions
mkdir completions
for sh in bash zsh fish; do
  go run cmd/azion/main.go completion "$sh" >"completions/azion.$sh"
done
