#!/bin/bash

# Default values
FROM="github.com/aziontech/azionapi-v4-go-sdk"
TO="github.com/aziontech/azionapi-v4-go-sdk-dev"
DIRECTORIES=(".") # Default: full project

# Show usage help
show_help() {
    cat <<EOF
Usage: ./update_imports.sh [OPTIONS]

Options:
  --undo                    Reverts the change (sdk-dev → sdk)
  --directories DIR [DIR..] Limit updates to specific directories (recursively)
  --help                    Show this help message

Examples:
  ./update_imports.sh
      Replace sdk → sdk-dev in all .go files in the project

  ./update_imports.sh --undo
      Revert sdk-dev → sdk in all .go files

  ./update_imports.sh --directories ./pkg ./cmd
      Update only inside ./pkg and ./cmd

  ./update_imports.sh --undo --directories ./pkg
      Revert only inside ./pkg
EOF
}

# Parse arguments
while [[ "$#" -gt 0 ]]; do
    case "$1" in
        --undo)
            FROM="github.com/aziontech/azionapi-v4-go-sdk-dev"
            TO="github.com/aziontech/azionapi-v4-go-sdk"
            shift
            ;;
        --directories)
            shift
            DIRECTORIES=()
            while [[ "$#" -gt 0 && "$1" != --* ]]; do
                DIRECTORIES+=("$1")
                shift
            done
            ;;
        --help)
            show_help
            exit 0
            ;;
        *)
            echo "❌ Unknown option: $1"
            echo "Use --help for usage instructions."
            exit 1
            ;;
    esac
done

# Info
echo "🔁 Replacing:"
echo "   FROM: $FROM"
echo "     TO: $TO"
echo "📂 In directories: ${DIRECTORIES[*]}"

# Apply replacements
for dir in "${DIRECTORIES[@]}"; do
    find "$dir" -type f -name "*.go" -print0 | while IFS= read -r -d '' file; do
        if grep -q "$FROM" "$file"; then
            echo "➤ Updating $file"
            sed -i '' "s|$FROM|$TO|g" "$file"
        fi
    done
done

echo "✅ Done."
