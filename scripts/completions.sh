#!/bin/sh
# scripts/completions.sh
set -e
rm -rf completions
mkdir completions
for sh in bash zsh fish; do
  go run cmd/azioncli/main.go completion "$sh" --no-update >"completions/azioncli.$sh"
done
