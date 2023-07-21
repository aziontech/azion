#!/bin/sh
# scripts/completions.sh

rm -rf completions
mkdir completions

generate_completions() {
  local shell=$1
  local output_file="completions/azioncli.$shell"

  if go run cmd/azioncli/main.go completion "$shell" --no-update > "$output_file"; then
    echo "Completions generated for $shell"
  else
    echo "Failed to generate completions for $shell"
  fi
}

shells=("bash" "zsh" "fish")

for shell in "${shells[@]}"; do
  generate_completions "$shell"
done
