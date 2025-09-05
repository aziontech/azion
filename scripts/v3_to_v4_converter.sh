#!/bin/bash

# v3_to_v4_converter.sh
# Script to convert azion.json from V3 to V4 format
# If domain exists in V3, creates a workload using Azion CLI

set -e

# Default values
INPUT_FILE="azion.json"
OUTPUT_FILE="azion.json.v4"
VERBOSE=false
BACKUP=false
DRY_RUN=false
CREATE_WORKLOAD=false
EXISTING_WORKLOAD_ID=""

# Function to display usage information
usage() {
    echo "Usage: $0 [OPTIONS]"
    echo "Convert azion.json from V3 to V4 format"
    echo ""
    echo "Options:"
    echo "  -i, --input FILE    Input azion.json file (default: azion.json)"
    echo "  -o, --output FILE   Output file (default: azion.json.v4)"
    echo "  -v, --verbose       Enable verbose output"
    echo "  -b, --backup        Create a backup of the original file"
    echo "  -w, --create-workload  Create a workload if domain exists (requires Azion CLI)"
    echo "  -d, --dry-run       Show what would be done without making changes"
    echo "  -W, --workload-id ID  Use an existing workload ID instead of creating a new one"
    echo "  -h, --help          Display this help message"
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

# Create backup if requested
if [ "$BACKUP" = true ] && [ -f "$INPUT_FILE" ]; then
    BACKUP_FILE="${INPUT_FILE}.backup-$(date +%Y%m%d%H%M%S)"
    cp "$INPUT_FILE" "$BACKUP_FILE"
    echo "Backup created: $BACKUP_FILE"
fi

# Function to log verbose messages
log() {
    if [ "$VERBOSE" = true ]; then
        echo "$1"
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

log "Converting $INPUT_FILE from V3 to V4 format..."

# Read the V3 JSON file
V3_JSON=$(cat "$INPUT_FILE")

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

# Extract domain and application info if domain exists
if [ "$HAS_DOMAIN" = "true" ]; then
    DOMAIN_NAME=$(echo "$V3_JSON" | jq -r '.domain.name')
    # Get application ID and ensure it's a number
    APP_ID=$(echo "$V3_JSON" | jq -r '.application.id' | tr -d '\n')
    APP_NAME=$(echo "$V3_JSON" | jq -r '.application.name')
    
    log "Domain found: $DOMAIN_NAME"
    
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
            WORKLOAD_DETAILS=$(azion describe workload --workload-id "$WORKLOAD_ID" --format json 2>&1) || {
                warn "Failed to get workload details: $WORKLOAD_DETAILS"
                log "Continuing with conversion using only the workload ID"
                WORKLOAD_NAME="unknown_workload"
            }
            
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
        log "Creating workload for domain $DOMAIN_NAME..."
        
        # Create a unique workload name based on the domain name
        WORKLOAD_NAME="${DOMAIN_NAME}"
        
        if [ "$DRY_RUN" = true ]; then
            log "[DRY RUN] Would create workload: $WORKLOAD_NAME"
            # Use a dummy ID for dry run
            WORKLOAD_ID="999999"
        else
            # Create the workload using Azion CLI
            log "Creating workload '$WORKLOAD_NAME' using Azion CLI..."
            WORKLOAD_RESULT=$(azion create workload --name "$WORKLOAD_NAME" --active true 2>&1) || {
                warn "Failed to create workload: $WORKLOAD_RESULT"
                log "Continuing with conversion without creating workload"
            }
            
            # Extract workload ID from the result
            WORKLOAD_ID=$(echo "$WORKLOAD_RESULT" | grep -o 'with ID [0-9]*' | awk '{print $3}')
            
            if [ -z "$WORKLOAD_ID" ]; then
                warn "Failed to extract workload ID from CLI output"
                log "Continuing with conversion without workload ID"
            else
                log "Workload created with ID: $WORKLOAD_ID"
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

# Add empty connectors array
V4_JSON=$(echo "$V4_JSON" | jq '.connectors = []')

# Check if origins exist in the V3 JSON
HAS_ORIGINS=$(echo "$V3_JSON" | jq -r 'if .origin and (.origin | length > 0) then "true" else "false" end')

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
    
    # If origins are found, provide guidance on creating connectors
    if [ "$HAS_ORIGINS" = "true" ]; then
        echo ""
        echo "=== IMPORTANT: Origins Found - Connector Creation Required ==="
        echo "Your V3 configuration contains origins that need to be converted to connectors in V4."
        echo ""
        echo "To create a connector using Azion CLI, run:"
        echo "azion create connector --name 'your-connector-name' --type 'connector_type' [flags]"
        echo ""
        echo "Available connector types: http, storage, live_ingest"
        echo ""
        echo "After creating the connector, update your azion.json file with the connector information:"
        echo ""
        echo "Example connector entry in azion.json:"
        echo ""
        echo '"connectors": ['  
        echo '  {'  
        echo '    "id": 123456,'  
        echo '    "name": "your-connector-name",'  
        echo '    "address": ['  
        echo '      {'  
        echo '        "address": "origin.example.com",'  
        echo '        "weight": 1'  
        echo '      }'  
        echo '    ]'  
        echo '  }'  
        echo ']'  
        echo ""
        echo "Please, read the instructions above carefully and create the connector using the Azion CLI."
        echo ""
        echo "For more information, run: azion create connector --help"
        echo "=================================================="
    fi
fi

exit 0
