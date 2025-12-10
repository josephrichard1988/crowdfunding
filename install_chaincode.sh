#!/bin/bash

# =============================================================================
# Chaincode Installation Script for Crowdfunding Platform
# =============================================================================
# Usage:
#   ./install_chaincode.sh <org>
#
# Examples:
#   ./install_chaincode.sh startup    # Install all chaincodes for StartupOrg
#   ./install_chaincode.sh validator  # Install all chaincodes for ValidatorOrg
#   ./install_chaincode.sh investor   # Install all chaincodes for InvestorOrg
#   ./install_chaincode.sh platform   # Install all chaincodes for PlatformOrg
# =============================================================================

set -e

# Configuration - Update these paths as needed
CONTRACTS_DIR="${CONTRACTS_DIR:-./contracts}"
STARTUP_CC_PATH="${CONTRACTS_DIR}/startuporg"
VALIDATOR_CC_PATH="${CONTRACTS_DIR}/validatororg"
INVESTOR_CC_PATH="${CONTRACTS_DIR}/investororg"
PLATFORM_CC_PATH="${CONTRACTS_DIR}/platformorg"

# Chaincode labels
STARTUP_CC_LABEL="startup_1.0"
VALIDATOR_CC_LABEL="validator_1.0"
INVESTOR_CC_LABEL="investor_1.0"
PLATFORM_CC_LABEL="platform_1.0"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# =============================================================================
# Helper Functions
# =============================================================================

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_section() {
    echo ""
    echo -e "${CYAN}=============================================================================${NC}"
    echo -e "${CYAN} $1${NC}"
    echo -e "${CYAN}=============================================================================${NC}"
}

# =============================================================================
# Package Chaincode
# =============================================================================

package_chaincode() {
    local cc_name=$1
    local cc_path=$2
    local cc_label=$3
    local output_file="${cc_name}.tar.gz"
    
    log_info "Packaging chaincode '${cc_name}' from '${cc_path}'..."
    
    # Check if path exists
    if [ ! -d "$cc_path" ]; then
        log_error "Chaincode path does not exist: $cc_path"
        return 1
    fi
    
    # Package the chaincode
    peer lifecycle chaincode package ${output_file} \
        --path ${cc_path} \
        --lang golang \
        --label ${cc_label}
    
    if [ $? -eq 0 ]; then
        log_success "Packaged '${cc_name}' -> ${output_file}"
        echo "${output_file}"
    else
        log_error "Failed to package '${cc_name}'"
        return 1
    fi
}

# =============================================================================
# Install Chaincode
# =============================================================================

install_chaincode() {
    local cc_name=$1
    local package_file=$2
    
    log_info "Installing chaincode '${cc_name}' from '${package_file}'..."
    
    # Check if package exists
    if [ ! -f "$package_file" ]; then
        log_error "Package file does not exist: $package_file"
        return 1
    fi
    
    # Install the chaincode
    peer lifecycle chaincode install ${package_file}
    
    if [ $? -eq 0 ]; then
        log_success "Installed '${cc_name}'"
    else
        log_error "Failed to install '${cc_name}'"
        return 1
    fi
}

# =============================================================================
# Query Installed Chaincodes
# =============================================================================

query_installed() {
    log_info "Querying installed chaincodes..."
    peer lifecycle chaincode queryinstalled
}

# =============================================================================
# Get Package ID
# =============================================================================

get_package_id() {
    local cc_label=$1
    
    # Query and extract package ID
    local package_id=$(peer lifecycle chaincode queryinstalled --output json | \
        jq -r ".installed_chaincodes[] | select(.label==\"${cc_label}\") | .package_id")
    
    if [ -n "$package_id" ]; then
        echo "$package_id"
    else
        log_error "Could not find package ID for label: ${cc_label}"
        return 1
    fi
}

# =============================================================================
# Package All Chaincodes (run once)
# =============================================================================

package_all_chaincodes() {
    log_section "Packaging All Chaincodes"
    
    package_chaincode "startup" "$STARTUP_CC_PATH" "$STARTUP_CC_LABEL"
    package_chaincode "validator" "$VALIDATOR_CC_PATH" "$VALIDATOR_CC_LABEL"
    package_chaincode "investor" "$INVESTOR_CC_PATH" "$INVESTOR_CC_LABEL"
    package_chaincode "platform" "$PLATFORM_CC_PATH" "$PLATFORM_CC_LABEL"
    
    log_success "All chaincodes packaged!"
    echo ""
    echo "Package files created:"
    ls -la *.tar.gz 2>/dev/null || echo "No .tar.gz files found"
}

# =============================================================================
# Install All Chaincodes for an Org
# =============================================================================

install_all_for_org() {
    local org=$1
    
    log_section "Installing All Chaincodes for ${org}"
    
    # Check if packages exist, if not, package them
    if [ ! -f "startup.tar.gz" ] || [ ! -f "validator.tar.gz" ] || \
       [ ! -f "investor.tar.gz" ] || [ ! -f "platform.tar.gz" ]; then
        log_warning "Package files not found. Packaging chaincodes first..."
        package_all_chaincodes
    fi
    
    # Install all 4 chaincodes
    install_chaincode "startup" "startup.tar.gz"
    install_chaincode "validator" "validator.tar.gz"
    install_chaincode "investor" "investor.tar.gz"
    install_chaincode "platform" "platform.tar.gz"
    
    log_success "All chaincodes installed for ${org}!"
    
    # Query and display installed chaincodes
    echo ""
    query_installed
    
    # Extract and display package IDs
    echo ""
    log_section "Package IDs for ${org}"
    echo ""
    echo "Export these for the deploy script:"
    echo ""
    
    local startup_id=$(get_package_id "$STARTUP_CC_LABEL" 2>/dev/null)
    local validator_id=$(get_package_id "$VALIDATOR_CC_LABEL" 2>/dev/null)
    local investor_id=$(get_package_id "$INVESTOR_CC_LABEL" 2>/dev/null)
    local platform_id=$(get_package_id "$PLATFORM_CC_LABEL" 2>/dev/null)
    
    echo "export STARTUP_CC_PACKAGE_ID=\"${startup_id}\""
    echo "export VALIDATOR_CC_PACKAGE_ID=\"${validator_id}\""
    echo "export INVESTOR_CC_PACKAGE_ID=\"${investor_id}\""
    echo "export PLATFORM_CC_PACKAGE_ID=\"${platform_id}\""
    echo ""
}

# =============================================================================
# Usage Information
# =============================================================================

show_usage() {
    echo ""
    echo "Usage: $0 <command> [org]"
    echo ""
    echo "Commands:"
    echo "  package              - Package all chaincodes (run once)"
    echo "  install <org>        - Install all chaincodes for specified org"
    echo "  query                - Query installed chaincodes"
    echo "  ids                  - Show package IDs for export"
    echo ""
    echo "Organizations:"
    echo "  startup     - Install for StartupOrg"
    echo "  validator   - Install for ValidatorOrg"
    echo "  investor    - Install for InvestorOrg"
    echo "  platform    - Install for PlatformOrg"
    echo ""
    echo "Examples:"
    echo "  $0 package                  # Package all chaincodes first"
    echo "  $0 install startup          # Install all chaincodes for StartupOrg"
    echo "  $0 install validator        # Install all chaincodes for ValidatorOrg"
    echo "  $0 query                    # Show installed chaincodes"
    echo "  $0 ids                      # Show package IDs to export"
    echo ""
    echo "Workflow:"
    echo "  1. Package chaincodes:  $0 package"
    echo "  2. Switch to org context (set CORE_PEER_* env vars)"
    echo "  3. Install for org:     $0 install <org>"
    echo "  4. Repeat steps 2-3 for each org"
    echo "  5. Export package IDs and run deploy_chaincode.sh"
    echo ""
}

# =============================================================================
# Show Package IDs
# =============================================================================

show_package_ids() {
    log_section "Package IDs"
    
    echo ""
    echo "Copy and export these:"
    echo ""
    
    local startup_id=$(get_package_id "$STARTUP_CC_LABEL" 2>/dev/null || echo "<not installed>")
    local validator_id=$(get_package_id "$VALIDATOR_CC_LABEL" 2>/dev/null || echo "<not installed>")
    local investor_id=$(get_package_id "$INVESTOR_CC_LABEL" 2>/dev/null || echo "<not installed>")
    local platform_id=$(get_package_id "$PLATFORM_CC_LABEL" 2>/dev/null || echo "<not installed>")
    
    echo "export STARTUP_CC_PACKAGE_ID=\"${startup_id}\""
    echo "export VALIDATOR_CC_PACKAGE_ID=\"${validator_id}\""
    echo "export INVESTOR_CC_PACKAGE_ID=\"${investor_id}\""
    echo "export PLATFORM_CC_PACKAGE_ID=\"${platform_id}\""
    echo ""
}

# =============================================================================
# Main Entry Point
# =============================================================================

main() {
    local command=$1
    local org=$2
    
    if [ -z "$command" ]; then
        show_usage
        exit 1
    fi
    
    case $command in
        package)
            package_all_chaincodes
            ;;
        install)
            if [ -z "$org" ]; then
                log_error "Please specify an organization"
                show_usage
                exit 1
            fi
            case $org in
                startup|validator|investor|platform)
                    install_all_for_org "${org^}Org"
                    ;;
                *)
                    log_error "Unknown organization: $org"
                    show_usage
                    exit 1
                    ;;
            esac
            ;;
        query)
            query_installed
            ;;
        ids)
            show_package_ids
            ;;
        help|--help|-h)
            show_usage
            ;;
        *)
            log_error "Unknown command: $command"
            show_usage
            exit 1
            ;;
    esac
}

# Run main function
main "$@"
