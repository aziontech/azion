#!/bin/bash

# v3_to_v4_converter.sh
# Script to convert azion.json from V3 to V4 format
# If domain exists in V3, creates a workload using Azion CLI
# Optionally reads a manifest.json to configure workloads and connectors
#
# V3 → V4 mapping:
#   - V3 domain  → V4 workload  (manifest.workloads config applied)
#   - V3 origin  → V4 connector (manifest.connectors config applied)
#
# Manifest connector format (ConnectorPolymorphicRequest):
#   Storage: { "name": "...", "active": true, "type": "storage", "attributes": { "bucket": "...", "prefix": "..." } }
#   HTTP:    { "name": "...", "active": true, "type": "http", "attributes": { "addresses": [{ "address": "..." }], ... } }
#
# Manifest workload format (WorkloadManifest):
#   { "name": "...", "active": true, "domains": [], "tls": {...}, "protocols": {...}, ... }

set -e

# Default values
INPUT_FILE="azion.json"
OUTPUT_FILE="azion.json.v4"
MANIFEST_FILE=""
VERBOSE=false
BACKUP=false
DRY_RUN=false
CREATE_WORKLOAD=false
CREATE_CONNECTORS=false
EXISTING_WORKLOAD_ID=""

# Function to display usage information
usage() {
    echo "Usage: $0 [OPTIONS]"
    echo "Convert azion.json from V3 to V4 format"
    echo ""
    echo "Options:"
    echo "  -i, --input FILE           Input azion.json file (default: azion.json)"
    echo "  -o, --output FILE          Output file (default: azion.json.v4)"
    echo "  -m, --manifest FILE        Manifest file (manifest.json) for workload/connector configuration"
    echo "  -v, --verbose              Enable verbose output"
    echo "  -b, --backup               Create a backup of the original file"
    echo "  -w, --create-workload      Create a workload if domain exists (requires Azion CLI)"
    echo "  -c, --create-connectors    Create connectors from V3 origins using manifest config (requires --manifest)"
    echo "  -d, --dry-run              Show what would be done without making changes"
    echo "  -W, --workload-id ID       Use an existing workload ID instead of creating a new one"
    echo "  -h, --help                 Display this help message"
    echo ""
    echo "V3 to V4 Mapping:"
    echo "  domain  → workload    (uses manifest workloads[] config if available)"
    echo "  origin  → connector   (uses manifest connectors[] config if available)"
    echo ""
    echo "Examples:"
    echo "  $0 -i azion.json -o azion.json.v4"
    echo "  $0 -i azion.json -w -m manifest.json"
    echo "  $0 -i azion.json -w -c -m manifest.json"
    echo "  $0 -i azion.json -w -c -m manifest.json -d"
    echo ""
    exit 1
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    key="$1"
    case $key in
        -i|--input)
            INPUT_FILE="$2"
            shift 2
            ;;
        -o|--output)
            OUTPUT_FILE="$2"
            shift 2
            ;;
        -m|--manifest)
            MANIFEST_FILE="$2"
            shift 2
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -b|--backup)
            BACKUP=true
            shift
            ;;
        -w|--create-workload)
            CREATE_WORKLOAD=true
            shift
            ;;
        -c|--create-connectors)
            CREATE_CONNECTORS=true
            shift
            ;;
        -W|--workload-id)
            EXISTING_WORKLOAD_ID="$2"
            shift 2
            ;;
        -d|--dry-run)
            DRY_RUN=true
            shift
            ;;
        -h|--help)
            usage
            ;;
        *)
            echo "Unknown option: $1"
            usage
            ;;
    esac
done

# Check if input file exists
if [ ! -f "$INPUT_FILE" ]; then
    echo "Error: Input file '$INPUT_FILE' not found"
    exit 1
fi

# Check if jq is installed
if ! command -v jq &> /dev/null; then
    echo "Error: jq is required but not installed. Please install jq first."
    echo "You can install it using: brew install jq (on macOS) or apt-get install jq (on Ubuntu)"
    exit 1
fi

# Check if azion CLI is installed
if ! command -v azion &> /dev/null; then
    echo "Error: azion CLI is required but not installed. Please install azion CLI first."
    exit 1
fi

# Validate manifest file if provided
if [ -n "$MANIFEST_FILE" ]; then
    if [ ! -f "$MANIFEST_FILE" ]; then
        echo "Error: Manifest file '$MANIFEST_FILE' not found"
        exit 1
    fi
    # Validate it's valid JSON
    if ! jq empty "$MANIFEST_FILE" 2>/dev/null; then
        echo "Error: Manifest file '$MANIFEST_FILE' is not valid JSON"
        exit 1
    fi
fi

# Validate --create-connectors requires --manifest
if [ "$CREATE_CONNECTORS" = true ] && [ -z "$MANIFEST_FILE" ]; then
    echo "Error: --create-connectors requires --manifest to be specified"
    exit 1
fi

# Create backup if requested
if [ "$BACKUP" = true ] && [ -f "$INPUT_FILE" ]; then
    BACKUP_FILE="${INPUT_FILE}.backup-$(date +%Y%m%d%H%M%S)"
    cp "$INPUT_FILE" "$BACKUP_FILE"
    echo "Backup created: $BACKUP_FILE"
fi

# Function to log verbose messages (outputs to stderr to avoid polluting stdout in subshells)
log() {
    if [ "$VERBOSE" = true ]; then
        echo "$1" >&2
    fi
}

# Function to log errors
error() {
    echo "ERROR: $1" >&2
}

# Function to log warnings
warn() {
    echo "WARNING: $1" >&2
}

# Function to clean up temporary files
cleanup() {
    if [ -n "${TEMP_DIR:-}" ] && [ -d "$TEMP_DIR" ]; then
        rm -rf "$TEMP_DIR"
    fi
}
trap cleanup EXIT

log "Converting $INPUT_FILE from V3 to V4 format..."

# Read the V3 JSON file
V3_JSON=$(cat "$INPUT_FILE")

# Read manifest if provided
MANIFEST_JSON=""
if [ -n "$MANIFEST_FILE" ]; then
    MANIFEST_JSON=$(cat "$MANIFEST_FILE")
    log "Manifest loaded from $MANIFEST_FILE"
fi

# Check if it's already in V4 format (has workloads field)
if echo "$V3_JSON" | jq -e '.workloads' &> /dev/null; then
    log "File appears to already be in V4 format (has workloads field)"
    echo "File appears to already be in V4 format. No conversion needed."
    exit 0
fi

# Convert Function from object to array
V4_JSON=$(echo "$V3_JSON" | jq '.function = [.function]')
if [ $? -ne 0 ]; then
    error "Failed to convert function field to array. Check if the input file is valid JSON."
    exit 1
fi

# Check if domain exists in V3
HAS_DOMAIN=$(echo "$V3_JSON" | jq -r 'if .domain and .domain.name and .domain.name != "" then "true" else "false" end')

# Initialize variables
WORKLOAD_ID=""
DOMAIN_NAME=""
APP_ID=""
APP_NAME=""

# Create temp directory for intermediate JSON files
TEMP_DIR=$(mktemp -d)

# ============================================================================
# Function: find_manifest_workload
# Finds a workload configuration in the manifest that best matches the V3 domain.
#
# V3 Domain structure:  { "id", "name", "domain_name", "url" }
# V4 Workload (manifest): { "name", "active", "infrastructure", "workload_domain_allow_access",
#                            "domains", "tls", "protocols", "mtls" }
#
# Matching strategies (in order):
#   1. Match by domain name in the workload's domains array
#   2. Match by application name (workload name == app name)
#   3. Match by domain name (workload name == domain name)
#   4. If only one workload exists, use it as fallback
# ============================================================================
find_manifest_workload() {
    local domain_name="$1"
    local app_name="$2"

    if [ -z "$MANIFEST_JSON" ]; then
        return 1
    fi

    local workload_count
    workload_count=$(echo "$MANIFEST_JSON" | jq '.workloads | length')

    if [ "$workload_count" = "0" ] || [ "$workload_count" = "null" ]; then
        log "No workloads found in manifest"
        return 1
    fi

    # Strategy 1: Look for a workload that has a matching domain in its domains array
    local idx
    idx=$(echo "$MANIFEST_JSON" | jq --arg domain "$domain_name" \
        '[.workloads | to_entries[] | select(.value.domains[]? == $domain) | .key] | first // empty')

    if [ -n "$idx" ] && [ "$idx" != "null" ]; then
        log "Found manifest workload matching domain '$domain_name' at index $idx"
        echo "$idx"
        return 0
    fi

    # Strategy 2: Match by application name (manifest workload name == app name)
    if [ -n "$app_name" ]; then
        idx=$(echo "$MANIFEST_JSON" | jq --arg name "$app_name" \
            '[.workloads | to_entries[] | select(.value.name == $name) | .key] | first // empty')

        if [ -n "$idx" ] && [ "$idx" != "null" ]; then
            log "Found manifest workload matching app name '$app_name' at index $idx"
            echo "$idx"
            return 0
        fi
    fi

    # Strategy 3: Match by domain name (manifest workload name == domain name)
    idx=$(echo "$MANIFEST_JSON" | jq --arg name "$domain_name" \
        '[.workloads | to_entries[] | select(.value.name == $name) | .key] | first // empty')

    if [ -n "$idx" ] && [ "$idx" != "null" ]; then
        log "Found manifest workload matching domain name '$domain_name' at index $idx"
        echo "$idx"
        return 0
    fi

    # Strategy 4: If there's only one workload in the manifest, use it
    if [ "$workload_count" = "1" ]; then
        log "Using the only workload found in manifest (index 0)"
        echo "0"
        return 0
    fi

    log "No matching workload found in manifest for domain '$domain_name'"
    return 1
}

# ============================================================================
# Function: build_workload_create_json
# Builds a workload creation JSON from the manifest's WorkloadManifest config.
#
# Manifest WorkloadManifest fields → WorkloadRequest:
#   name                         → name
#   active                       → active
#   domains[]                    → domains[]
#   tls.certificate              → tls.certificate
#   tls.ciphers                  → tls.ciphers
#   tls.minimum_version          → tls.minimum_version
#   protocols.http.versions[]    → protocols.http.versions[]
#   protocols.http.http_ports[]  → protocols.http.http_ports[]
#   protocols.http.https_ports[] → protocols.http.https_ports[]
#   protocols.http.quic_ports[]  → protocols.http.quic_ports[]
#   workload_domain_allow_access → workload_domain_allow_access
#   mtls                         → mtls
# ============================================================================
build_workload_create_json() {
    local workload_idx="$1"
    local override_name="$2"

    local workload_json
    workload_json=$(echo "$MANIFEST_JSON" | jq --argjson idx "$workload_idx" '.workloads[$idx]')

    # Start with required fields
    local create_json
    create_json=$(echo "$workload_json" | jq --arg name "$override_name" '{
        name: $name,
        active: (.active // true)
    }')

    # Add TLS configuration if present
    if echo "$workload_json" | jq -e '.tls' &> /dev/null; then
        create_json=$(echo "$create_json" | jq --argjson tls "$(echo "$workload_json" | jq '.tls')" \
            '. + {tls: $tls}')
        log "  Added TLS config (minimum_version: $(echo "$workload_json" | jq -r '.tls.minimum_version // "not set"'))"
    fi

    # Add protocols configuration if present
    if echo "$workload_json" | jq -e '.protocols' &> /dev/null; then
        create_json=$(echo "$create_json" | jq --argjson protocols "$(echo "$workload_json" | jq '.protocols')" \
            '. + {protocols: $protocols}')
        log "  Added protocols config (versions: $(echo "$workload_json" | jq -c '.protocols.http.versions // []'))"
    fi

    # Add workload_domain_allow_access if present
    if echo "$workload_json" | jq -e '.workload_domain_allow_access' &> /dev/null; then
        create_json=$(echo "$create_json" | jq --argjson wda "$(echo "$workload_json" | jq '.workload_domain_allow_access')" \
            '. + {workload_domain_allow_access: $wda}')
        log "  Added workload_domain_allow_access: $(echo "$workload_json" | jq '.workload_domain_allow_access')"
    fi

    # Add domains if present and non-empty
    if echo "$workload_json" | jq -e '.domains | length > 0' &> /dev/null; then
        create_json=$(echo "$create_json" | jq --argjson domains "$(echo "$workload_json" | jq '.domains')" \
            '. + {domains: $domains}')
        log "  Added $(echo "$workload_json" | jq '.domains | length') domain(s)"
    fi

    # Add mTLS configuration if present
    if echo "$workload_json" | jq -e '.mtls' &> /dev/null; then
        create_json=$(echo "$create_json" | jq --argjson mtls "$(echo "$workload_json" | jq '.mtls')" \
            '. + {mtls: $mtls}')
        log "  Added mTLS configuration"
    fi

    echo "$create_json"
}

# ============================================================================
# Function: fetch_v3_domain_details
# Fetches the full V3 domain details from the API using azion describe domain.
# V3 domain fields: name, cnames, cname_access_only, is_active, edge_application_id,
#   digital_certificate_id, is_mtls_enabled, mtls_verification,
#   mtls_trusted_ca_certificate_id, edge_firewall_id, crl_list
# ============================================================================
fetch_v3_domain_details() {
    local domain_id="$1"

    if [ -z "$domain_id" ] || [ "$domain_id" = "null" ]; then
        log "No domain ID available, cannot fetch V3 domain details"
        return 1
    fi

    local result
    result=$(azion describe domain --domain-id "$domain_id" --format json 2>&1) || true

    if echo "$result" | grep -qi "error\|fail"; then
        log "Could not fetch V3 domain details for ID $domain_id: $result"
        return 1
    fi

    echo "$result"
}

# ============================================================================
# Function: build_workload_from_v3_domain
# Builds a workload creation JSON by mapping V3 domain fields → V4 workload.
#
# V3 Domain → V4 Workload field mapping:
#   name                           → name
#   is_active                      → active
#   is_mtls_enabled                → mtls.enabled
#   mtls_verification              → mtls.config.verification
#   mtls_trusted_ca_certificate_id → mtls.certificate
#   crl_list                       → mtls.config.crl
#
# V3 fields NOT mapped (warnings emitted):
#   edge_application_id  → maps to workload_deployment (handled separately)
#   edge_firewall_id     → maps to workload_deployment (handled separately)
#   cnames               → no V4 equivalent
#   cname_access_only    → no V4 equivalent
#   digital_certificate_id → too complex, user must configure manually
# ============================================================================
build_workload_from_v3_domain() {
    local domain_json="$1"
    local override_name="$2"

    # Start with name and active status
    local name
    name=$(echo "$domain_json" | jq -r '.name // ""')
    if [ -n "$override_name" ]; then
        name="$override_name"
    fi

    local is_active
    is_active=$(echo "$domain_json" | jq '.is_active // true')

    local create_json
    create_json=$(jq -n --arg name "$name" --argjson active "$is_active" \
        '{name: $name, active: $active}')

    # Map mTLS fields
    local is_mtls_enabled
    is_mtls_enabled=$(echo "$domain_json" | jq '.is_mtls_enabled // false')

    if [ "$is_mtls_enabled" = "true" ]; then
        local mtls_json='{"enabled": true}'

        # Map mtls_trusted_ca_certificate_id → mtls.certificate
        local mtls_cert_id
        mtls_cert_id=$(echo "$domain_json" | jq '.mtls_trusted_ca_certificate_id // 0')
        if [ "$mtls_cert_id" != "0" ] && [ "$mtls_cert_id" != "null" ]; then
            mtls_json=$(echo "$mtls_json" | jq --argjson cert "$mtls_cert_id" '. + {certificate: $cert}')
        fi

        # Map mtls_verification → mtls.config.verification
        local mtls_verification
        mtls_verification=$(echo "$domain_json" | jq -r '.mtls_verification // ""')
        if [ -n "$mtls_verification" ]; then
            mtls_json=$(echo "$mtls_json" | jq --arg ver "$mtls_verification" \
                '.config = (.config // {}) | .config.verification = $ver')
        fi

        # Map crl_list → mtls.config.crl
        local crl_list
        crl_list=$(echo "$domain_json" | jq '.crl_list // []')
        if [ "$crl_list" != "[]" ] && [ "$crl_list" != "null" ]; then
            mtls_json=$(echo "$mtls_json" | jq --argjson crl "$crl_list" \
                '.config = (.config // {}) | .config.crl = $crl')
        fi

        create_json=$(echo "$create_json" | jq --argjson mtls "$mtls_json" '. + {mtls: $mtls}')
        log "  Added mTLS configuration from V3 domain"
    fi

    # Warn about non-convertible fields
    local digital_cert_id
    digital_cert_id=$(echo "$domain_json" | jq -r '.digital_certificate_id // "null"')
    if [ "$digital_cert_id" != "null" ] && [ "$digital_cert_id" != "" ]; then
        warn "V3 domain has digital_certificate_id ($digital_cert_id). This is too complex to convert automatically."
        warn "Please configure the TLS certificate manually on the new workload after migration."
    fi

    local cnames
    cnames=$(echo "$domain_json" | jq '.cnames // []')
    if [ "$cnames" != "[]" ] && [ "$cnames" != "null" ]; then
        warn "V3 domain has CNAMEs: $(echo "$cnames" | jq -c '.'). CNAMEs have no direct V4 workload equivalent."
        warn "You may need to configure custom domains separately after migration."
    fi

    local cname_access_only
    cname_access_only=$(echo "$domain_json" | jq '.cname_access_only // false')
    if [ "$cname_access_only" = "true" ]; then
        warn "V3 domain has cname_access_only=true. This has no direct V4 workload equivalent."
    fi

    echo "$create_json"
}

# ============================================================================
# Function: create_workload_deployment
# Creates a workload_deployment for the newly created workload.
#
# V3 Domain → V4 WorkloadDeployment mapping:
#   edge_application_id → strategy.attributes.application (as application-id)
#   edge_firewall_id    → strategy.attributes.firewall (as firewall-id)
#
# CLI: azion create workload-deployment --name <name> --workload-id <wid>
#      --application-id <app_id> [--firewall-id <fw_id>] --active true --current true
# ============================================================================
create_workload_deployment() {
    local workload_id="$1"
    local workload_name="$2"
    local edge_app_id="$3"
    local edge_firewall_id="$4"

    if [ -z "$workload_id" ] || [ "$workload_id" = "0" ]; then
        log "No workload ID, skipping deployment creation"
        return 1
    fi

    if [ -z "$edge_app_id" ] || [ "$edge_app_id" = "0" ] || [ "$edge_app_id" = "null" ]; then
        log "No edge_application_id found, skipping deployment creation"
        return 1
    fi

    local deployment_name="${workload_name}"

    if [ "$DRY_RUN" = true ]; then
        echo "[DRY RUN] Would create workload_deployment:"
        echo "  name: $deployment_name"
        echo "  workload_id: $workload_id"
        echo "  application_id: $edge_app_id"
        if [ -n "$edge_firewall_id" ] && [ "$edge_firewall_id" != "0" ] && [ "$edge_firewall_id" != "null" ]; then
            echo "  firewall_id: $edge_firewall_id"
        fi
        return 0
    fi

    log "Creating workload_deployment '$deployment_name' for workload $workload_id..."

    local deploy_cmd="azion create workload-deployment --name \"$deployment_name\" --workload-id $workload_id --application-id $edge_app_id --active true --current true"

    if [ -n "$edge_firewall_id" ] && [ "$edge_firewall_id" != "0" ] && [ "$edge_firewall_id" != "null" ]; then
        deploy_cmd="$deploy_cmd --firewall-id $edge_firewall_id"
    fi

    local deploy_result
    deploy_result=$(eval "$deploy_cmd" 2>&1) || true

    local deploy_id
    deploy_id=$(echo "$deploy_result" | grep -o 'with ID [0-9]*' | awk '{print $3}')

    if [ -n "$deploy_id" ]; then
        echo "Workload deployment '$deployment_name' created with ID: $deploy_id"
        echo "  → V3 edge_application_id ($edge_app_id) mapped to deployment strategy"
        if [ -n "$edge_firewall_id" ] && [ "$edge_firewall_id" != "0" ] && [ "$edge_firewall_id" != "null" ]; then
            echo "  → V3 edge_firewall_id ($edge_firewall_id) mapped to deployment strategy"
        fi
        return 0
    else
        warn "Failed to create workload_deployment: $deploy_result"
        return 1
    fi
}

# ============================================================================
# Function: find_manifest_connector_for_origin
# Finds a connector in the manifest that corresponds to a V3 origin.
#
# V3 Origin structure:
#   { "name", "origin_type", "bucket", "prefix", "addresses": [{"address","weight"}], "host_header", ... }
#
# V4 Manifest Connector (ConnectorPolymorphicRequest) structure:
#   Storage: { "name", "active", "type": "storage", "attributes": { "bucket", "prefix" } }
#   HTTP:    { "name", "active", "type": "http", "attributes": { "addresses": [{"address","active","http_port","https_port"}], "connection_options", "modules" } }
#
# Matching strategies:
#   1. Match by name (origin name == connector name)
#   2. For storage origins: match by bucket name
#   3. For HTTP origins: match by first address
# ============================================================================
find_manifest_connector_for_origin() {
    local origin_name="$1"
    local origin_type="$2"
    local origin_bucket="$3"
    local origin_first_address="$4"

    if [ -z "$MANIFEST_JSON" ]; then
        return 1
    fi

    local connector_count
    connector_count=$(echo "$MANIFEST_JSON" | jq '.connectors | length')

    if [ "$connector_count" = "0" ] || [ "$connector_count" = "null" ]; then
        log "No connectors found in manifest"
        return 1
    fi

    # Strategy 1: Match by name
    local idx
    idx=$(echo "$MANIFEST_JSON" | jq --arg name "$origin_name" \
        '[.connectors | to_entries[] | select(
            (.value.name // "") == $name or
            (.value.ConnectorHTTPRequest.name // "") == $name or
            (.value.ConnectorRequest.name // "") == $name
        ) | .key] | first // empty')

    if [ -n "$idx" ] && [ "$idx" != "null" ]; then
        log "Found manifest connector matching origin name '$origin_name' at index $idx"
        echo "$idx"
        return 0
    fi

    # Strategy 2: For storage-like origins, match by bucket
    if [ -n "$origin_bucket" ]; then
        idx=$(echo "$MANIFEST_JSON" | jq --arg bucket "$origin_bucket" \
            '[.connectors | to_entries[] | select(
                (.value.attributes.bucket // "") == $bucket
            ) | .key] | first // empty')

        if [ -n "$idx" ] && [ "$idx" != "null" ]; then
            log "Found manifest connector matching bucket '$origin_bucket' at index $idx"
            echo "$idx"
            return 0
        fi
    fi

    # Strategy 3: For HTTP origins, match by first address
    if [ -n "$origin_first_address" ]; then
        idx=$(echo "$MANIFEST_JSON" | jq --arg addr "$origin_first_address" \
            '[.connectors | to_entries[] | select(
                (.value.attributes.addresses[]?.address // "") == $addr
            ) | .key] | first // empty')

        if [ -n "$idx" ] && [ "$idx" != "null" ]; then
            log "Found manifest connector matching address '$origin_first_address' at index $idx"
            echo "$idx"
            return 0
        fi
    fi

    log "No matching connector found in manifest for origin '$origin_name'"
    return 1
}

# ============================================================================
# Function: get_connector_type
# Determines the connector type from a manifest connector entry.
# The manifest stores connectors as ConnectorPolymorphicRequest (oneOf: ConnectorHTTPRequest | ConnectorRequest)
# which always has a "type" field: "http" or "storage"
# ============================================================================
get_connector_type() {
    local connector_json="$1"
    echo "$connector_json" | jq -r '.type // "unknown"'
}

# ============================================================================
# Function: create_connectors_from_manifest
# Creates connectors by mapping V3 origins to manifest connectors.
#
# For each V3 origin:
#   1. Find a matching connector in the manifest
#   2. If found, use the manifest connector JSON (which is already in
#      ConnectorPolymorphicRequest format) to create via CLI
#   3. If not found, build a basic connector from the V3 origin data
#
# V3 origin_type mapping:
#   "object_storage"/"single_origin"/"load_balancer" → determines connector type
#   but the manifest connector "type" field takes precedence
# ============================================================================
create_connectors_from_manifest() {
    if [ -z "$MANIFEST_JSON" ]; then
        warn "No manifest provided, cannot create connectors"
        return 1
    fi

    local origin_count
    origin_count=$(echo "$V3_JSON" | jq '.origin | length')

    local connector_count
    connector_count=$(echo "$MANIFEST_JSON" | jq '.connectors | length')

    if [ "$origin_count" = "0" ] || [ "$origin_count" = "null" ]; then
        # No V3 origins - still try to create connectors from manifest directly
        if [ "$connector_count" = "0" ] || [ "$connector_count" = "null" ]; then
            log "No origins in V3 and no connectors in manifest"
            return 0
        fi
        log "No V3 origins found, creating connectors directly from manifest"
        create_manifest_connectors_directly
        return $?
    fi

    echo ""
    echo "=== Creating Connectors (V3 Origins → V4 Connectors) ==="
    log "Found $origin_count V3 origin(s) and $connector_count manifest connector(s)"

    local created_connectors=()
    local i=0

    while [ $i -lt "$origin_count" ]; do
        local origin_json
        origin_json=$(echo "$V3_JSON" | jq --argjson idx "$i" '.origin[$idx]')

        local origin_name
        origin_name=$(echo "$origin_json" | jq -r '.name // ""')

        local origin_type
        origin_type=$(echo "$origin_json" | jq -r '.origin_type // ""')

        local origin_bucket
        origin_bucket=$(echo "$origin_json" | jq -r '.bucket // ""')

        local origin_first_address
        origin_first_address=$(echo "$origin_json" | jq -r '.addresses[0].address // ""')

        log "Processing V3 origin [$i]: name='$origin_name', type='$origin_type'"

        # Try to find a matching manifest connector for this origin
        local manifest_connector_idx=""
        manifest_connector_idx=$(find_manifest_connector_for_origin "$origin_name" "$origin_type" "$origin_bucket" "$origin_first_address") || true

        local create_json=""
        local connector_type=""

        if [ -n "$manifest_connector_idx" ]; then
            # Use the manifest connector config directly (it's already in ConnectorPolymorphicRequest format)
            create_json=$(echo "$MANIFEST_JSON" | jq --argjson idx "$manifest_connector_idx" '.connectors[$idx]')
            connector_type=$(get_connector_type "$create_json")
            log "  Using manifest connector config (index $manifest_connector_idx, type: $connector_type)"
        else
            # Build a basic connector from V3 origin data
            log "  No manifest match, building connector from V3 origin data"

            # Determine if this is a storage or HTTP origin
            # Check: explicit origin_type, bucket field, or name containing "storage"
            local is_storage=false
            if [ "$origin_type" = "object_storage" ] || [ -n "$origin_bucket" ]; then
                is_storage=true
            elif echo "$origin_name" | grep -qi "storage"; then
                is_storage=true
                log "  Detected storage origin by name heuristic ('$origin_name' contains 'storage')"
            fi

            if [ "$is_storage" = true ]; then
                # V3 object_storage origin → V4 storage connector (ConnectorRequest)
                connector_type="storage"

                # Use origin's bucket/prefix, fall back to project-level bucket/prefix
                local bucket prefix
                bucket="$origin_bucket"
                prefix=$(echo "$origin_json" | jq -r '.prefix // ""')

                if [ -z "$bucket" ]; then
                    bucket=$(echo "$V3_JSON" | jq -r '.bucket // ""')
                    log "  Using project-level bucket: $bucket"
                fi
                if [ -z "$prefix" ]; then
                    prefix=$(echo "$V3_JSON" | jq -r '.prefix // ""')
                    log "  Using project-level prefix: $prefix"
                fi

                create_json=$(jq -n \
                    --arg name "$origin_name" \
                    --arg bucket "$bucket" \
                    --arg prefix "$prefix" \
                    '{
                        name: $name,
                        active: true,
                        type: "storage",
                        attributes: {
                            bucket: $bucket,
                            prefix: $prefix
                        }
                    }')

                # Remove empty prefix if not set
                if [ -z "$prefix" ]; then
                    create_json=$(echo "$create_json" | jq 'del(.attributes.prefix)')
                fi
            else
                # V3 single_origin/load_balancer → V4 HTTP connector (ConnectorHTTPRequest)
                connector_type="http"

                # Map V3 addresses [{address, weight}] → V4 addresses [{address, active}]
                local v4_addresses
                v4_addresses=$(echo "$origin_json" | jq '[.addresses[]? | {address: .address, active: true}]')
                if [ "$v4_addresses" = "[]" ] || [ "$v4_addresses" = "null" ]; then
                    v4_addresses="[]"
                fi

                local host_header
                host_header=$(echo "$origin_json" | jq -r '.host_header // ""')

                create_json=$(jq -n \
                    --arg name "$origin_name" \
                    --argjson addresses "$v4_addresses" \
                    '{
                        name: $name,
                        active: true,
                        type: "http",
                        attributes: {
                            addresses: $addresses
                        }
                    }')

                # Add connection_options with host_header if available
                if [ -n "$host_header" ]; then
                    create_json=$(echo "$create_json" | jq --arg hh "$host_header" \
                        '.attributes.connection_options = {host_header: $hh}')
                fi
            fi
        fi

        if [ -z "$create_json" ]; then
            warn "Could not build connector config for origin '$origin_name'. Skipping."
            i=$((i + 1))
            continue
        fi

        # Write the connector JSON to a temp file
        local temp_connector_file="$TEMP_DIR/connector_${i}.json"
        echo "$create_json" | jq '.' > "$temp_connector_file"

        if [ "$DRY_RUN" = true ]; then
            echo "[DRY RUN] Would create $connector_type connector for V3 origin '$origin_name':"
            cat "$temp_connector_file"
            echo ""
        else
            log "Creating $connector_type connector for V3 origin '$origin_name'..."

            CONNECTOR_RESULT=$(azion create connector --type "$connector_type" --file "$temp_connector_file" 2>&1) || true

            # Check if connector was created successfully
            CONNECTOR_ID=$(echo "$CONNECTOR_RESULT" | grep -o 'with ID [0-9]*' | awk '{print $3}')

            # Handle "Invalid bucket name" error for storage connectors
            # This happens when the user changed accounts and the old bucket doesn't exist
            if [ -z "$CONNECTOR_ID" ] && echo "$CONNECTOR_RESULT" | grep -qi "Invalid bucket name"; then
                local old_bucket
                old_bucket=$(echo "$create_json" | jq -r '.attributes.bucket // ""')
                local timestamp
                timestamp=$(date +%Y%m%d%H%M%S)
                local new_bucket="${old_bucket}-${timestamp}"

                log "  Bucket '$old_bucket' not found (likely account change). Creating new bucket '$new_bucket'..."
                echo "Bucket '$old_bucket' is invalid. Creating new bucket '$new_bucket'..."

                # Create the new bucket using the CLI
                local bucket_result
                bucket_result=$(azion create storage bucket --name "$new_bucket" --workloads-access read_only 2>&1) || true

                if echo "$bucket_result" | grep -qi "error\|fail"; then
                    warn "Failed to create bucket '$new_bucket': $bucket_result"
                    warn "Failed to create connector for origin '$origin_name'"
                    i=$((i + 1))
                    continue
                fi

                echo "Bucket '$new_bucket' created successfully"

                # Update the connector JSON with the new bucket name
                create_json=$(echo "$create_json" | jq --arg bucket "$new_bucket" '.attributes.bucket = $bucket')
                echo "$create_json" | jq '.' > "$temp_connector_file"

                log "  Retrying connector creation with new bucket '$new_bucket'..."

                # Retry creating the connector with the new bucket
                CONNECTOR_RESULT=$(azion create connector --type "$connector_type" --file "$temp_connector_file" 2>&1) || true
                CONNECTOR_ID=$(echo "$CONNECTOR_RESULT" | grep -o 'with ID [0-9]*' | awk '{print $3}')
            fi

            if [ -z "$CONNECTOR_ID" ] && echo "$CONNECTOR_RESULT" | grep -qi "error\|fail"; then
                warn "Failed to create connector for origin '$origin_name': $CONNECTOR_RESULT"
                i=$((i + 1))
                continue
            fi

            # Fallback ID extraction
            if [ -z "$CONNECTOR_ID" ]; then
                CONNECTOR_ID=$(echo "$CONNECTOR_RESULT" | grep -oE '[0-9]+' | head -1)
            fi

            local connector_name
            connector_name=$(echo "$create_json" | jq -r '.name')

            if [ -n "$CONNECTOR_ID" ]; then
                echo "Connector '$connector_name' (type: $connector_type) created with ID: $CONNECTOR_ID (from V3 origin: '$origin_name')"
                created_connectors+=("$CONNECTOR_ID:$connector_name:$connector_type")
            else
                warn "Connector for origin '$origin_name' may have been created but could not extract ID from: $CONNECTOR_RESULT"
            fi
        fi

        i=$((i + 1))
    done

    # Update V4 JSON with created connectors info
    if [ ${#created_connectors[@]} -gt 0 ]; then
        local connectors_array="[]"
        for entry in "${created_connectors[@]}"; do
            local cid cname ctype
            cid=$(echo "$entry" | cut -d: -f1)
            cname=$(echo "$entry" | cut -d: -f2)
            ctype=$(echo "$entry" | cut -d: -f3)
            connectors_array=$(echo "$connectors_array" | jq --arg id "$cid" --arg name "$cname" --arg type "$ctype" \
                '. + [{id: ($id | tonumber), name: $name, type: $type}]')
        done
        V4_JSON=$(echo "$V4_JSON" | jq --argjson conns "$connectors_array" '.connectors = $conns')
        log "Updated V4 JSON with ${#created_connectors[@]} connector(s)"
    fi

    echo "=== Connector Creation Complete ==="
    echo ""
}

create_manifest_connectors_directly() {
    local connector_count
    connector_count=$(echo "$MANIFEST_JSON" | jq '.connectors | length')

    echo ""
    echo "=== Creating Connectors from Manifest ==="
    log "Found $connector_count connector(s) in manifest (no V3 origins to map)"

    local created_connectors=()
    local i=0

    while [ $i -lt "$connector_count" ]; do
        local connector_json
        connector_json=$(echo "$MANIFEST_JSON" | jq --argjson idx "$i" '.connectors[$idx]')

        local connector_name
        connector_name=$(echo "$connector_json" | jq -r '.name // "unnamed"')

        local connector_type
        connector_type=$(get_connector_type "$connector_json")

        log "Processing manifest connector [$i]: name='$connector_name', type='$connector_type'"

        # Write the connector JSON to a temp file (manifest format is already the API format)
        local temp_connector_file="$TEMP_DIR/connector_${i}.json"
        echo "$connector_json" | jq '.' > "$temp_connector_file"

        if [ "$DRY_RUN" = true ]; then
            echo "[DRY RUN] Would create $connector_type connector '$connector_name':"
            cat "$temp_connector_file"
            echo ""
        else
            log "Creating $connector_type connector '$connector_name'..."

            CONNECTOR_RESULT=$(azion create connector --type "$connector_type" --file "$temp_connector_file" 2>&1) || true

            CONNECTOR_ID=$(echo "$CONNECTOR_RESULT" | grep -o 'with ID [0-9]*' | awk '{print $3}')

            # Handle "Invalid bucket name" for storage connectors
            if [ -z "$CONNECTOR_ID" ] && echo "$CONNECTOR_RESULT" | grep -qi "Invalid bucket name"; then
                local old_bucket
                old_bucket=$(echo "$connector_json" | jq -r '.attributes.bucket // ""')
                local timestamp
                timestamp=$(date +%Y%m%d%H%M%S)
                local new_bucket="${old_bucket}-${timestamp}"

                log "  Bucket '$old_bucket' not found (likely account change). Creating new bucket '$new_bucket'..."
                echo "Bucket '$old_bucket' is invalid. Creating new bucket '$new_bucket'..."

                local bucket_result
                bucket_result=$(azion create storage bucket --name "$new_bucket" --workloads-access read_only 2>&1) || true

                if echo "$bucket_result" | grep -qi "error\|fail"; then
                    warn "Failed to create bucket '$new_bucket': $bucket_result"
                    warn "Failed to create connector '$connector_name'"
                    i=$((i + 1))
                    continue
                fi

                echo "Bucket '$new_bucket' created successfully"

                connector_json=$(echo "$connector_json" | jq --arg bucket "$new_bucket" '.attributes.bucket = $bucket')
                echo "$connector_json" | jq '.' > "$temp_connector_file"

                CONNECTOR_RESULT=$(azion create connector --type "$connector_type" --file "$temp_connector_file" 2>&1) || true
                CONNECTOR_ID=$(echo "$CONNECTOR_RESULT" | grep -o 'with ID [0-9]*' | awk '{print $3}')
            fi

            if [ -z "$CONNECTOR_ID" ] && echo "$CONNECTOR_RESULT" | grep -qi "error\|fail"; then
                warn "Failed to create connector '$connector_name': $CONNECTOR_RESULT"
                i=$((i + 1))
                continue
            fi

            if [ -z "$CONNECTOR_ID" ]; then
                CONNECTOR_ID=$(echo "$CONNECTOR_RESULT" | grep -oE '[0-9]+' | head -1)
            fi

            if [ -n "$CONNECTOR_ID" ]; then
                echo "Connector '$connector_name' (type: $connector_type) created with ID: $CONNECTOR_ID"
                created_connectors+=("$CONNECTOR_ID:$connector_name:$connector_type")
            else
                warn "Connector '$connector_name' may have been created but could not extract ID from: $CONNECTOR_RESULT"
            fi
        fi

        i=$((i + 1))
    done

    # Update V4 JSON with created connectors info
    if [ ${#created_connectors[@]} -gt 0 ]; then
        local connectors_array="[]"
        for entry in "${created_connectors[@]}"; do
            local cid cname ctype
            cid=$(echo "$entry" | cut -d: -f1)
            cname=$(echo "$entry" | cut -d: -f2)
            ctype=$(echo "$entry" | cut -d: -f3)
            connectors_array=$(echo "$connectors_array" | jq --arg id "$cid" --arg name "$cname" --arg type "$ctype" \
                '. + [{id: ($id | tonumber), name: $name, type: $type}]')
        done
        V4_JSON=$(echo "$V4_JSON" | jq --argjson conns "$connectors_array" '.connectors = $conns')
        log "Updated V4 JSON with ${#created_connectors[@]} connector(s)"
    fi

    echo "=== Connector Creation Complete ==="
    echo ""
}

# Extract domain and application info if domain exists
if [ "$HAS_DOMAIN" = "true" ]; then
    DOMAIN_NAME=$(echo "$V3_JSON" | jq -r '.domain.name')
    DOMAIN_ID=$(echo "$V3_JSON" | jq -r '.domain.id // ""')
    # Get application ID and ensure it's a number
    APP_ID=$(echo "$V3_JSON" | jq -r '.application.id' | tr -d '\n')
    APP_NAME=$(echo "$V3_JSON" | jq -r '.application.name')
    
    log "Domain found: $DOMAIN_NAME (mapping to workload)"
    
    # Fetch full V3 domain details from the API for proper field mapping
    V3_DOMAIN_DETAILS=""
    V3_EDGE_APP_ID=""
    V3_EDGE_FIREWALL_ID=""
    if [ "$DRY_RUN" != true ] && [ -n "$DOMAIN_ID" ] && [ "$DOMAIN_ID" != "null" ]; then
        log "Fetching V3 domain details for ID: $DOMAIN_ID..."
        V3_DOMAIN_DETAILS=$(fetch_v3_domain_details "$DOMAIN_ID") || true
        if [ -n "$V3_DOMAIN_DETAILS" ]; then
            V3_EDGE_APP_ID=$(echo "$V3_DOMAIN_DETAILS" | jq -r '.edge_application_id // 0')
            V3_EDGE_FIREWALL_ID=$(echo "$V3_DOMAIN_DETAILS" | jq -r '.edge_firewall_id // 0')
            log "  V3 domain details fetched: edge_application_id=$V3_EDGE_APP_ID, edge_firewall_id=$V3_EDGE_FIREWALL_ID"
        fi
    fi
    
    # Check if an existing workload ID was provided
    if [ -n "$EXISTING_WORKLOAD_ID" ]; then
        log "Using existing workload ID: $EXISTING_WORKLOAD_ID"
        WORKLOAD_ID="$EXISTING_WORKLOAD_ID"
        
        # Get workload details using the describe command
        if [ "$DRY_RUN" = true ]; then
            log "[DRY RUN] Would get details for workload ID: $WORKLOAD_ID"
            WORKLOAD_NAME="existing_workload"
        else
            log "Getting details for workload ID: $WORKLOAD_ID"
            WORKLOAD_DETAILS=$(azion describe workload --workload-id "$WORKLOAD_ID" --format json 2>&1) || true
            if echo "$WORKLOAD_DETAILS" | grep -qi "error\|fail"; then
                warn "Failed to get workload details: $WORKLOAD_DETAILS"
                log "Continuing with conversion using only the workload ID"
                WORKLOAD_NAME="unknown_workload"
                WORKLOAD_DETAILS=""
            fi
            
            # Extract workload name from the details
            if [ -n "$WORKLOAD_DETAILS" ]; then
                WORKLOAD_NAME=$(echo "$WORKLOAD_DETAILS" | jq -r '.name' 2>/dev/null)
                if [ -z "$WORKLOAD_NAME" ] || [ "$WORKLOAD_NAME" = "null" ]; then
                    WORKLOAD_NAME="unknown_workload"
                    warn "Failed to extract workload name from details"
                else
                    log "Found workload name: $WORKLOAD_NAME"
                fi
            fi
        fi
    # Create workload only if the flag is set and no existing ID was provided
    elif [ "$CREATE_WORKLOAD" = true ]; then
        log "Creating workload for V3 domain '$DOMAIN_NAME' (domain → workload)..."
        
        # Create a unique workload name based on the domain name
        WORKLOAD_NAME="${DOMAIN_NAME}"
        
        # Check if manifest has a matching workload configuration
        MANIFEST_WORKLOAD_IDX=""
        if [ -n "$MANIFEST_JSON" ]; then
            MANIFEST_WORKLOAD_IDX=$(find_manifest_workload "$DOMAIN_NAME" "$APP_NAME") || true
        fi
        
        if [ "$DRY_RUN" = true ]; then
            if [ -n "$MANIFEST_WORKLOAD_IDX" ]; then
                WORKLOAD_CREATE_JSON=$(build_workload_create_json "$MANIFEST_WORKLOAD_IDX" "$WORKLOAD_NAME")
                echo "[DRY RUN] Would create workload '$WORKLOAD_NAME' using manifest config (index $MANIFEST_WORKLOAD_IDX):"
                echo "$WORKLOAD_CREATE_JSON" | jq '.'
            elif [ -n "$V3_DOMAIN_DETAILS" ]; then
                WORKLOAD_CREATE_JSON=$(build_workload_from_v3_domain "$V3_DOMAIN_DETAILS" "$WORKLOAD_NAME")
                echo "[DRY RUN] Would create workload '$WORKLOAD_NAME' from V3 domain details:"
                echo "$WORKLOAD_CREATE_JSON" | jq '.'
            else
                echo "[DRY RUN] Would create workload '$WORKLOAD_NAME' (basic config, no manifest match)"
            fi
            # Use a dummy ID for dry run
            WORKLOAD_ID="999999"

            # Dry run workload_deployment
            if [ -n "$V3_EDGE_APP_ID" ] && [ "$V3_EDGE_APP_ID" != "0" ] && [ "$V3_EDGE_APP_ID" != "null" ]; then
                create_workload_deployment "$WORKLOAD_ID" "$WORKLOAD_NAME" "$V3_EDGE_APP_ID" "$V3_EDGE_FIREWALL_ID"
            elif [ -n "$APP_ID" ] && [ "$APP_ID" != "0" ] && [ "$APP_ID" != "null" ]; then
                echo "[DRY RUN] Would create workload_deployment with application_id from V3 azion.json: $APP_ID"
            fi
        else
            # Build workload creation JSON
            # Priority: 1) manifest config, 2) V3 domain details, 3) basic config
            WORKLOAD_CREATE_JSON=""
            WORKLOAD_CREATE_SOURCE=""

            if [ -n "$MANIFEST_WORKLOAD_IDX" ]; then
                WORKLOAD_CREATE_JSON=$(build_workload_create_json "$MANIFEST_WORKLOAD_IDX" "$WORKLOAD_NAME")
                WORKLOAD_CREATE_SOURCE="manifest"
                log "  Using manifest config for workload creation"
            elif [ -n "$V3_DOMAIN_DETAILS" ]; then
                WORKLOAD_CREATE_JSON=$(build_workload_from_v3_domain "$V3_DOMAIN_DETAILS" "$WORKLOAD_NAME")
                WORKLOAD_CREATE_SOURCE="v3_domain"
                log "  Using V3 domain details for workload creation (with mTLS mapping)"
            fi

            # Attempt to create a workload and handle name collisions
            WORKLOAD_ATTEMPTS=0
            MAX_WORKLOAD_ATTEMPTS=2

            while [ $WORKLOAD_ATTEMPTS -lt $MAX_WORKLOAD_ATTEMPTS ]; do
                WORKLOAD_RESULT=""

                if [ -n "$WORKLOAD_CREATE_JSON" ]; then
                    # Update name in JSON (might have been changed by retry)
                    WORKLOAD_CREATE_JSON=$(echo "$WORKLOAD_CREATE_JSON" | jq --arg name "$WORKLOAD_NAME" '.name = $name')

                    log "Creating workload '$WORKLOAD_NAME' using $WORKLOAD_CREATE_SOURCE config..."
                    log "  Creation payload: $(echo "$WORKLOAD_CREATE_JSON" | jq -c '.')"

                    TEMP_WORKLOAD_FILE="$TEMP_DIR/workload_create.json"
                    echo "$WORKLOAD_CREATE_JSON" | jq '.' > "$TEMP_WORKLOAD_FILE"

                    WORKLOAD_RESULT=$(azion create workload --file "$TEMP_WORKLOAD_FILE" 2>&1) || true
                else
                    # Fallback: basic config with just name and active
                    log "Creating workload '$WORKLOAD_NAME' using Azion CLI (basic config)..."
                    WORKLOAD_RESULT=$(azion create workload --name "$WORKLOAD_NAME" --active true 2>&1) || true
                fi

                # Check if workload was created successfully
                WORKLOAD_ID=$(echo "$WORKLOAD_RESULT" | grep -o 'with ID [0-9]*' | awk '{print $3}')

                if [ -n "$WORKLOAD_ID" ]; then
                    # Success
                    break
                fi

                # Check if the error is "name already in use"
                if echo "$WORKLOAD_RESULT" | grep -qi "already in use"; then
                    RETRY_TIMESTAMP=$(date +%Y%m%d%H%M%S)
                    WORKLOAD_NAME="${DOMAIN_NAME}-${RETRY_TIMESTAMP}"
                    warn "Workload name already in use. Retrying with: $WORKLOAD_NAME"
                    WORKLOAD_ATTEMPTS=$((WORKLOAD_ATTEMPTS + 1))
                    continue
                fi

                # Other error — don't retry
                warn "Failed to create workload: $WORKLOAD_RESULT"
                log "Continuing with conversion without creating workload"
                break
            done
            
            if [ -n "$WORKLOAD_ID" ]; then
                log "Workload created with ID: $WORKLOAD_ID"
                
                if [ "$WORKLOAD_CREATE_SOURCE" = "manifest" ]; then
                    echo "Workload '$WORKLOAD_NAME' created with ID: $WORKLOAD_ID (configured from manifest, V3 domain → V4 workload)"
                elif [ "$WORKLOAD_CREATE_SOURCE" = "v3_domain" ]; then
                    echo "Workload '$WORKLOAD_NAME' created with ID: $WORKLOAD_ID (mapped from V3 domain details)"
                else
                    echo "Workload '$WORKLOAD_NAME' created with ID: $WORKLOAD_ID (V3 domain → V4 workload)"
                fi

                # Create workload_deployment: V3 edge_application_id/edge_firewall_id → V4 deployment strategy
                DEPLOY_APP_ID=""
                DEPLOY_FW_ID=""

                # Prefer V3 domain API details, fall back to azion.json application.id
                if [ -n "$V3_EDGE_APP_ID" ] && [ "$V3_EDGE_APP_ID" != "0" ] && [ "$V3_EDGE_APP_ID" != "null" ]; then
                    DEPLOY_APP_ID="$V3_EDGE_APP_ID"
                    DEPLOY_FW_ID="$V3_EDGE_FIREWALL_ID"
                elif [ -n "$APP_ID" ] && [ "$APP_ID" != "0" ] && [ "$APP_ID" != "null" ]; then
                    DEPLOY_APP_ID="$APP_ID"
                fi

                if [ -n "$DEPLOY_APP_ID" ]; then
                    create_workload_deployment "$WORKLOAD_ID" "$WORKLOAD_NAME" "$DEPLOY_APP_ID" "$DEPLOY_FW_ID"
                fi
            else
                warn "Failed to extract workload ID from CLI output"
                log "Continuing with conversion without workload ID"
            fi
        fi
    fi
fi

# Add workloads field to V4 JSON
if [ -n "$WORKLOAD_ID" ]; then
    # Get domain info
    DOMAIN_NAME=$(echo "$V3_JSON" | jq -r '.domain.name')
    DOMAIN_ID=$(echo "$V3_JSON" | jq -r '.domain.id')
    
    # Create workload object with the domain
    V4_JSON=$(echo "$V4_JSON" | jq --arg name "$WORKLOAD_NAME" \
                                  --arg id "$WORKLOAD_ID" \
                                  --arg domain "$DOMAIN_NAME" \
                                  '.workloads = {
                                      "id": ($id | tonumber),
                                      "name": $name,
                                      "domains": [$domain],
                                      "url": "",
                                      "deployments": []
                                  }')
else
    # Add empty workloads object
    V4_JSON=$(echo "$V4_JSON" | jq '.workloads = {"id": 0, "name": "", "domains": [], "url": "", "deployments": []}')
fi

# Add empty connectors array (will be populated if --create-connectors is used)
V4_JSON=$(echo "$V4_JSON" | jq '.connectors = []')

# Check if origins exist in the V3 JSON
HAS_ORIGINS=$(echo "$V3_JSON" | jq -r 'if .origin and (.origin | length > 0) then "true" else "false" end')


if [ "$CREATE_CONNECTORS" = true ]; then
    create_connectors_from_manifest
fi

# Write the V4 JSON to the output file
if [ "$DRY_RUN" = true ]; then
    echo "[DRY RUN] V4 JSON content would be:"
    echo "$V4_JSON" | jq '.'
else
    echo "$V4_JSON" | jq '.' > "$OUTPUT_FILE"
    log "Conversion complete. V4 format saved to $OUTPUT_FILE"
    
    if [ "$OUTPUT_FILE" != "$INPUT_FILE" ]; then
        echo "To replace the original file with the V4 version, run:"
        echo "mv $OUTPUT_FILE $INPUT_FILE"
    else
        echo "Original file has been updated to V4 format"
    fi
    
    # If origins are found and connectors were NOT auto-created, provide guidance
    if [ "$HAS_ORIGINS" = "true" ] && [ "$CREATE_CONNECTORS" = false ]; then
        echo ""
        echo "=== IMPORTANT: Origins Found - Connector Creation Required ==="
        echo "Your V3 configuration contains origins that need to be converted to connectors in V4."
        echo ""
        echo "V3 origins map to V4 connectors as follows:"
        echo "  object_storage origin → storage connector"
        echo "  single_origin/load_balancer origin → http connector"
        echo ""

        if [ -n "$MANIFEST_FILE" ]; then
            echo "A manifest file was provided. You can automatically create connectors by re-running with:"
            echo "  $0 -i $INPUT_FILE -c -m $MANIFEST_FILE"
            echo ""
        else
            echo "Provide a manifest.json file with connector definitions to automate this:"
            echo "  $0 -i $INPUT_FILE -c -m manifest.json"
            echo ""
        fi

        echo "Connector JSON format for the manifest:"
        echo ""
        echo "Storage connector:"
        echo '  { "name": "my-storage", "active": true, "type": "storage", "attributes": { "bucket": "my-bucket", "prefix": "20260101" } }'
        echo ""
        echo "HTTP connector:"
        echo '  { "name": "my-http", "active": true, "type": "http", "attributes": { "addresses": [{ "address": "origin.example.com" }] } }'
        echo ""
        echo "Or create manually: azion create connector --type <type> --file connector.json"
        echo "For more information, run: azion create connector --help"
        echo "=================================================="
    fi
fi

exit 0
