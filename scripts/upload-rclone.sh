#!/usr/bin/env bash
set -euo pipefail

# Upload static files to Azion S3-compatible storage using rclone.
# Mirrors the concurrency profile of `azion deploy-remote` (20 parallel transfers).

SETTINGS="${AZION_SETTINGS:-$HOME/.azion/ProdV4/settings.toml}"
SRC="${SRC:-bin/adventurous-simonpearson/.edge/storage/adventurous-simonpearson1}"
REQUIRED_RCLONE_VERSION="v1.73.1"
TRANSFERS="${TRANSFERS:-20}"
ENDPOINT="${ENDPOINT:-https://s3.us-east-005.azionstorage.net}"
REGION="${REGION:-us-east-005}"
BUCKET=""

usage() {
    cat <<EOF
Usage: $(basename "$0") --bucket <name> [--src <dir>] [--transfers <n>]

Required:
  --bucket <name>      Destination S3 bucket name

Optional:
  --src <dir>          Source directory (default: $SRC)
  --transfers <n>      Parallel transfers (default: $TRANSFERS)
  -h, --help           Show this help
EOF
}

while [[ $# -gt 0 ]]; do
    case "$1" in
        --bucket)
            BUCKET="${2:-}"
            shift 2
            ;;
        --bucket=*)
            BUCKET="${1#*=}"
            shift
            ;;
        --src)
            SRC="${2:-}"
            shift 2
            ;;
        --src=*)
            SRC="${1#*=}"
            shift
            ;;
        --transfers)
            TRANSFERS="${2:-}"
            shift 2
            ;;
        --transfers=*)
            TRANSFERS="${1#*=}"
            shift
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            echo "error: unknown argument: $1" >&2
            usage >&2
            exit 1
            ;;
    esac
done

if [[ -z "$BUCKET" ]]; then
    echo "error: --bucket is required" >&2
    usage >&2
    exit 1
fi

if ! command -v rclone >/dev/null 2>&1; then
    echo "error: rclone is not installed (need $REQUIRED_RCLONE_VERSION)" >&2
    exit 1
fi

installed_version="$(rclone version | head -n1 | awk '{print $2}')"
if [[ "$installed_version" != "$REQUIRED_RCLONE_VERSION" ]]; then
    echo "error: rclone $REQUIRED_RCLONE_VERSION required, found $installed_version" >&2
    exit 1
fi

if [[ ! -f "$SETTINGS" ]]; then
    echo "error: settings file not found at $SETTINGS" >&2
    exit 1
fi

if [[ ! -d "$SRC" ]]; then
    echo "error: source directory not found at $SRC" >&2
    exit 1
fi

# Pull S3 creds from settings.toml. Values are single-quoted in the file.
toml_get() {
    grep -E "^[[:space:]]*$1[[:space:]]*=" "$SETTINGS" \
        | head -n1 \
        | sed -E "s/^[[:space:]]*$1[[:space:]]*=[[:space:]]*['\"]?([^'\"]*)['\"]?[[:space:]]*\$/\1/"
}

ACCESS_KEY="$(toml_get S3AccessKey)"
SECRET_KEY="$(toml_get S3SecretKey)"

if [[ -z "$ACCESS_KEY" || -z "$SECRET_KEY" ]]; then
    echo "error: missing S3 credentials in $SETTINGS" >&2
    exit 1
fi

export RCLONE_CONFIG_AZION_TYPE=s3
export RCLONE_CONFIG_AZION_PROVIDER=Other
export RCLONE_CONFIG_AZION_ENDPOINT="$ENDPOINT"
export RCLONE_CONFIG_AZION_REGION="$REGION"
export RCLONE_CONFIG_AZION_ACCESS_KEY_ID="$ACCESS_KEY"
export RCLONE_CONFIG_AZION_SECRET_ACCESS_KEY="$SECRET_KEY"
export RCLONE_CONFIG_AZION_ACL=private

echo "rclone $installed_version"
echo "src:      $SRC"
echo "dest:     azion:$BUCKET"
echo "endpoint: $ENDPOINT"
echo "transfers: $TRANSFERS"
echo

# --copy-links: source contains a symlink (20260529103001 -> ../public); dereference it.
# --transfers / --checkers: parallelism, matching deploy-remote's 20-worker upload pool.
# --s3-no-check-bucket: skip the HeadBucket round-trip (bucket is known to exist).
rclone copy "$SRC" "azion:$BUCKET" \
    --transfers "$TRANSFERS" \
    --checkers "$TRANSFERS" \
    --copy-links \
    --s3-no-check-bucket \
    --s3-no-head \
    --progress \
    --stats 1s
