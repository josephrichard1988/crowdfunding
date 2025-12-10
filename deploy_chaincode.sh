#!/bin/bash

# =============================================================================
# Chaincode Deployment Script for Crowdfunding Platform
# =============================================================================
# Usage:
#   ./deploy_chaincode.sh <org> [channel]
#
# Examples:
#   ./deploy_chaincode.sh startup              # Deploy all chaincodes for StartupOrg
#   ./deploy_chaincode.sh validator            # Deploy all chaincodes for ValidatorOrg
#   ./deploy_chaincode.sh investor             # Deploy all chaincodes for InvestorOrg
#   ./deploy_chaincode.sh platform             # Deploy all chaincodes for PlatformOrg
#   ./deploy_chaincode.sh startup common       # Deploy only on common-channel for StartupOrg //don't use this
#   ./deploy_chaincode.sh all                  # Deploy for all orgs (sequential) //don't use this
# =============================================================================

set -e

# Configuration
ORDERER_URL="orderer-api.127-0-0-1.nip.io:9090"

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
    echo -e "${YELLOW}=============================================================================${NC}"
    echo -e "${YELLOW} $1${NC}"
    echo -e "${YELLOW}=============================================================================${NC}"
}

# =============================================================================
# Sequence and Version Detection Functions
# =============================================================================

# Get the current committed sequence for a chaincode on a channel
# Returns 0 if not committed yet, otherwise returns the current sequence
get_committed_sequence() {
    local channel=$1
    local chaincode_name=$2
    
    # Query committed chaincode and extract sequence
    local result=$(peer lifecycle chaincode querycommitted \
        --channelID ${channel} \
        --name ${chaincode_name} \
        --output json 2>/dev/null || echo "")
    
    if [ -z "$result" ] || [ "$result" == "" ]; then
        echo "0"
        return
    fi
    
    # Extract sequence from JSON
    local sequence=$(echo "$result" | jq -r '.sequence // 0' 2>/dev/null || echo "0")
    
    if [ -z "$sequence" ] || [ "$sequence" == "null" ]; then
        echo "0"
    else
        echo "$sequence"
    fi
}

# Get the current committed version for a chaincode on a channel
# Returns "0" if not committed yet
get_committed_version() {
    local channel=$1
    local chaincode_name=$2
    
    # Query committed chaincode and extract version
    local result=$(peer lifecycle chaincode querycommitted \
        --channelID ${channel} \
        --name ${chaincode_name} \
        --output json 2>/dev/null || echo "")
    
    if [ -z "$result" ] || [ "$result" == "" ]; then
        echo "0"
        return
    fi
    
    # Extract version from JSON
    local version=$(echo "$result" | jq -r '.version // "0"' 2>/dev/null || echo "0")
    
    if [ -z "$version" ] || [ "$version" == "null" ]; then
        echo "0"
    else
        echo "$version"
    fi
}

# Get next sequence number (current + 1, minimum 1)
get_next_sequence() {
    local channel=$1
    local chaincode_name=$2
    
    local current=$(get_committed_sequence "$channel" "$chaincode_name")
    local next=$((current + 1))
    
    echo "$next"
}

# Get next version (current + 1, minimum 1)
get_next_version() {
    local channel=$1
    local chaincode_name=$2
    
    local current=$(get_committed_version "$channel" "$chaincode_name")
    
    # Handle version as integer
    if [ "$current" == "0" ] || [ -z "$current" ]; then
        echo "1"
    else
        # Try to increment if numeric, otherwise append .1
        if [[ "$current" =~ ^[0-9]+$ ]]; then
            local next=$((current + 1))
            echo "$next"
        else
            # For non-numeric versions like "1.0", just increment sequence part
            echo "$current"
        fi
    fi
}

# Show current sequence and version info
show_sequence() {
    local channel=$1
    local chaincode_name=$2
    
    local current_seq=$(get_committed_sequence "$channel" "$chaincode_name")
    local current_ver=$(get_committed_version "$channel" "$chaincode_name")
    local next_seq=$((current_seq + 1))
    local next_ver=$(get_next_version "$channel" "$chaincode_name")
    
    if [ "$current_seq" == "0" ]; then
        echo -e "${CYAN}[INFO]${NC} '${chaincode_name}' on '${channel}': Not committed yet. Will use version: ${next_ver}, sequence: ${next_seq}"
    else
        echo -e "${CYAN}[INFO]${NC} '${chaincode_name}' on '${channel}': Current version: ${current_ver}, sequence: ${current_seq} -> Next version: ${next_ver}, sequence: ${next_seq}"
    fi
}

# =============================================================================
# Approve and Commit Functions
# =============================================================================

# Approve chaincode for org with auto-sequence and auto-version detection
approve_chaincode() {
    local channel=$1
    local chaincode_name=$2
    local package_id=$3
    
    # Get the next sequence and version
    local sequence=$(get_next_sequence "$channel" "$chaincode_name")
    local version=$(get_next_version "$channel" "$chaincode_name")
    
    log_info "Approving chaincode '${chaincode_name}' on channel '${channel}' (version: ${version}, sequence: ${sequence})..."
    
    peer lifecycle chaincode approveformyorg \
        -o ${ORDERER_URL} \
        --channelID ${channel} \
        --name ${chaincode_name} \
        --version ${version} \
        --sequence ${sequence} \
        --waitForEvent \
        --package-id ${package_id}
    
    if [ $? -eq 0 ]; then
        log_success "Approved '${chaincode_name}' on '${channel}' (version: ${version}, sequence: ${sequence})"
    else
        log_error "Failed to approve '${chaincode_name}' on '${channel}'"
        return 1
    fi
}

# Commit chaincode with auto-sequence and auto-version detection
commit_chaincode() {
    local channel=$1
    local chaincode_name=$2
    
    # Get the next sequence and version
    local sequence=$(get_next_sequence "$channel" "$chaincode_name")
    local version=$(get_next_version "$channel" "$chaincode_name")
    
    log_info "Committing chaincode '${chaincode_name}' on channel '${channel}' (version: ${version}, sequence: ${sequence})..."
    
    peer lifecycle chaincode commit \
        -o ${ORDERER_URL} \
        --channelID ${channel} \
        --name ${chaincode_name} \
        --version ${version} \
        --sequence ${sequence}
    
    if [ $? -eq 0 ]; then
        log_success "Committed '${chaincode_name}' on '${channel}' (version: ${version}, sequence: ${sequence})"
    else
        log_error "Failed to commit '${chaincode_name}' on '${channel}'"
        return 1
    fi
}

# Check commit readiness with auto-sequence and auto-version detection
check_readiness() {
    local channel=$1
    local chaincode_name=$2
    
    local sequence=$(get_next_sequence "$channel" "$chaincode_name")
    local version=$(get_next_version "$channel" "$chaincode_name")
    
    log_info "Checking commit readiness for '${chaincode_name}' on '${channel}' (version: ${version}, sequence: ${sequence})..."
    
    peer lifecycle chaincode checkcommitreadiness \
        --channelID ${channel} \
        --name ${chaincode_name} \
        --version ${version} \
        --sequence ${sequence}
}

# =============================================================================
# Package ID Setup Functions
# =============================================================================

set_package_ids() {
    # Set these after installing chaincodes
    # You can either export them before running or modify here
    
    if [ -z "$STARTUP_CC_PACKAGE_ID" ]; then
        log_warning "STARTUP_CC_PACKAGE_ID not set. Please export it or set below."
        # STARTUP_CC_PACKAGE_ID="startup_1.0:xxxx"
    fi
    
    if [ -z "$VALIDATOR_CC_PACKAGE_ID" ]; then
        log_warning "VALIDATOR_CC_PACKAGE_ID not set. Please export it or set below."
        # VALIDATOR_CC_PACKAGE_ID="validator_1.0:xxxx"
    fi
    
    if [ -z "$INVESTOR_CC_PACKAGE_ID" ]; then
        log_warning "INVESTOR_CC_PACKAGE_ID not set. Please export it or set below."
        # INVESTOR_CC_PACKAGE_ID="investor_1.0:xxxx"
    fi
    
    if [ -z "$PLATFORM_CC_PACKAGE_ID" ]; then
        log_warning "PLATFORM_CC_PACKAGE_ID not set. Please export it or set below."
        # PLATFORM_CC_PACKAGE_ID="platform_1.0:xxxx"
    fi
}

# =============================================================================
# StartupOrg Deployment
# =============================================================================

deploy_startup_org() {
    local channel_filter=$1
    
    log_section "Deploying chaincodes for StartupOrg"
    
    # startup-validator-channel
    if [ -z "$channel_filter" ] || [ "$channel_filter" == "startup-validator" ]; then
        log_section "StartupOrg: startup-validator-channel"
        
        # Own chaincode
        approve_chaincode "startup-validator-channel" "startup" "$STARTUP_CC_PACKAGE_ID"
        commit_chaincode "startup-validator-channel" "startup"
        
        # Other org's chaincode
        approve_chaincode "startup-validator-channel" "validator" "$VALIDATOR_CC_PACKAGE_ID"
    fi
    
    # startup-platform-channel
    if [ -z "$channel_filter" ] || [ "$channel_filter" == "startup-platform" ]; then
        log_section "StartupOrg: startup-platform-channel"
        
        # Own chaincode
        approve_chaincode "startup-platform-channel" "startup" "$STARTUP_CC_PACKAGE_ID"
        commit_chaincode "startup-platform-channel" "startup"
        
        # Other org's chaincode
        approve_chaincode "startup-platform-channel" "platform" "$PLATFORM_CC_PACKAGE_ID"
    fi
    
    # startup-investor-channel
    if [ -z "$channel_filter" ] || [ "$channel_filter" == "startup-investor" ]; then
        log_section "StartupOrg: startup-investor-channel"
        
        # Own chaincode
        approve_chaincode "startup-investor-channel" "startup" "$STARTUP_CC_PACKAGE_ID"
        commit_chaincode "startup-investor-channel" "startup"
        
        # Other org's chaincode
        approve_chaincode "startup-investor-channel" "investor" "$INVESTOR_CC_PACKAGE_ID"
    fi
    
    # common-channel
    if [ -z "$channel_filter" ] || [ "$channel_filter" == "common" ]; then
        log_section "StartupOrg: common-channel"
        
        # Own chaincode
        approve_chaincode "common-channel" "startup" "$STARTUP_CC_PACKAGE_ID"
        commit_chaincode "common-channel" "startup"
        
        # Other orgs' chaincodes
        approve_chaincode "common-channel" "validator" "$VALIDATOR_CC_PACKAGE_ID"
        approve_chaincode "common-channel" "investor" "$INVESTOR_CC_PACKAGE_ID"
        approve_chaincode "common-channel" "platform" "$PLATFORM_CC_PACKAGE_ID"
    fi
    
    log_success "StartupOrg deployment complete!"
}

# =============================================================================
# ValidatorOrg Deployment
# =============================================================================

deploy_validator_org() {
    local channel_filter=$1
    
    log_section "Deploying chaincodes for ValidatorOrg"
    
    # startup-validator-channel
    if [ -z "$channel_filter" ] || [ "$channel_filter" == "startup-validator" ]; then
        log_section "ValidatorOrg: startup-validator-channel"
        
        # Own chaincode
        approve_chaincode "startup-validator-channel" "validator" "$VALIDATOR_CC_PACKAGE_ID"
        commit_chaincode "startup-validator-channel" "validator"
        
        # Other org's chaincode
        approve_chaincode "startup-validator-channel" "startup" "$STARTUP_CC_PACKAGE_ID"
    fi
    
    # investor-validator-channel
    if [ -z "$channel_filter" ] || [ "$channel_filter" == "investor-validator" ]; then
        log_section "ValidatorOrg: investor-validator-channel"
        
        # Own chaincode
        approve_chaincode "investor-validator-channel" "validator" "$VALIDATOR_CC_PACKAGE_ID"
        commit_chaincode "investor-validator-channel" "validator"
        
        # Other org's chaincode
        approve_chaincode "investor-validator-channel" "investor" "$INVESTOR_CC_PACKAGE_ID"
    fi
    
    # validator-platform-channel
    if [ -z "$channel_filter" ] || [ "$channel_filter" == "validator-platform" ]; then
        log_section "ValidatorOrg: validator-platform-channel"
        
        # Own chaincode
        approve_chaincode "validator-platform-channel" "validator" "$VALIDATOR_CC_PACKAGE_ID"
        commit_chaincode "validator-platform-channel" "validator"
        
        # Other org's chaincode
        approve_chaincode "validator-platform-channel" "platform" "$PLATFORM_CC_PACKAGE_ID"
    fi
    
    # common-channel
    if [ -z "$channel_filter" ] || [ "$channel_filter" == "common" ]; then
        log_section "ValidatorOrg: common-channel"
        
        # Own chaincode
        approve_chaincode "common-channel" "validator" "$VALIDATOR_CC_PACKAGE_ID"
        commit_chaincode "common-channel" "validator"
        
        # Other orgs' chaincodes
        approve_chaincode "common-channel" "startup" "$STARTUP_CC_PACKAGE_ID"
        approve_chaincode "common-channel" "investor" "$INVESTOR_CC_PACKAGE_ID"
        approve_chaincode "common-channel" "platform" "$PLATFORM_CC_PACKAGE_ID"
    fi
    
    log_success "ValidatorOrg deployment complete!"
}

# =============================================================================
# InvestorOrg Deployment
# =============================================================================

deploy_investor_org() {
    local channel_filter=$1
    
    log_section "Deploying chaincodes for InvestorOrg"
    
    # startup-investor-channel
    if [ -z "$channel_filter" ] || [ "$channel_filter" == "startup-investor" ]; then
        log_section "InvestorOrg: startup-investor-channel"
        
        # Own chaincode
        approve_chaincode "startup-investor-channel" "investor" "$INVESTOR_CC_PACKAGE_ID"
        commit_chaincode "startup-investor-channel" "investor"
        
        # Other org's chaincode
        approve_chaincode "startup-investor-channel" "startup" "$STARTUP_CC_PACKAGE_ID"
    fi
    
    # investor-platform-channel
    if [ -z "$channel_filter" ] || [ "$channel_filter" == "investor-platform" ]; then
        log_section "InvestorOrg: investor-platform-channel"
        
        # Own chaincode
        approve_chaincode "investor-platform-channel" "investor" "$INVESTOR_CC_PACKAGE_ID"
        commit_chaincode "investor-platform-channel" "investor"
        
        # Other org's chaincode
        approve_chaincode "investor-platform-channel" "platform" "$PLATFORM_CC_PACKAGE_ID"
    fi
    
    # investor-validator-channel
    if [ -z "$channel_filter" ] || [ "$channel_filter" == "investor-validator" ]; then
        log_section "InvestorOrg: investor-validator-channel"
        
        # Own chaincode
        approve_chaincode "investor-validator-channel" "investor" "$INVESTOR_CC_PACKAGE_ID"
        commit_chaincode "investor-validator-channel" "investor"
        
        # Other org's chaincode
        approve_chaincode "investor-validator-channel" "validator" "$VALIDATOR_CC_PACKAGE_ID"
    fi
    
    # common-channel
    if [ -z "$channel_filter" ] || [ "$channel_filter" == "common" ]; then
        log_section "InvestorOrg: common-channel"
        
        # Own chaincode
        approve_chaincode "common-channel" "investor" "$INVESTOR_CC_PACKAGE_ID"
        commit_chaincode "common-channel" "investor"
        
        # Other orgs' chaincodes
        approve_chaincode "common-channel" "startup" "$STARTUP_CC_PACKAGE_ID"
        approve_chaincode "common-channel" "validator" "$VALIDATOR_CC_PACKAGE_ID"
        approve_chaincode "common-channel" "platform" "$PLATFORM_CC_PACKAGE_ID"
    fi
    
    log_success "InvestorOrg deployment complete!"
}

# =============================================================================
# PlatformOrg Deployment
# =============================================================================

deploy_platform_org() {
    local channel_filter=$1
    
    log_section "Deploying chaincodes for PlatformOrg"
    
    # startup-platform-channel
    if [ -z "$channel_filter" ] || [ "$channel_filter" == "startup-platform" ]; then
        log_section "PlatformOrg: startup-platform-channel"
        
        # Own chaincode
        approve_chaincode "startup-platform-channel" "platform" "$PLATFORM_CC_PACKAGE_ID"
        commit_chaincode "startup-platform-channel" "platform"
        
        # Other org's chaincode
        approve_chaincode "startup-platform-channel" "startup" "$STARTUP_CC_PACKAGE_ID"
    fi
    
    # investor-platform-channel
    if [ -z "$channel_filter" ] || [ "$channel_filter" == "investor-platform" ]; then
        log_section "PlatformOrg: investor-platform-channel"
        
        # Own chaincode
        approve_chaincode "investor-platform-channel" "platform" "$PLATFORM_CC_PACKAGE_ID"
        commit_chaincode "investor-platform-channel" "platform"
        
        # Other org's chaincode
        approve_chaincode "investor-platform-channel" "investor" "$INVESTOR_CC_PACKAGE_ID"
    fi
    
    # validator-platform-channel
    if [ -z "$channel_filter" ] || [ "$channel_filter" == "validator-platform" ]; then
        log_section "PlatformOrg: validator-platform-channel"
        
        # Own chaincode
        approve_chaincode "validator-platform-channel" "platform" "$PLATFORM_CC_PACKAGE_ID"
        commit_chaincode "validator-platform-channel" "platform"
        
        # Other org's chaincode
        approve_chaincode "validator-platform-channel" "validator" "$VALIDATOR_CC_PACKAGE_ID"
    fi
    
    # common-channel
    if [ -z "$channel_filter" ] || [ "$channel_filter" == "common" ]; then
        log_section "PlatformOrg: common-channel"
        
        # Own chaincode
        approve_chaincode "common-channel" "platform" "$PLATFORM_CC_PACKAGE_ID"
        commit_chaincode "common-channel" "platform"
        
        # Other orgs' chaincodes
        approve_chaincode "common-channel" "startup" "$STARTUP_CC_PACKAGE_ID"
        approve_chaincode "common-channel" "validator" "$VALIDATOR_CC_PACKAGE_ID"
        approve_chaincode "common-channel" "investor" "$INVESTOR_CC_PACKAGE_ID"
    fi
    
    log_success "PlatformOrg deployment complete!"
}

# =============================================================================
# Deploy All Orgs
# =============================================================================

deploy_all_orgs() {
    log_section "Deploying chaincodes for ALL organizations"
    
    log_warning "Make sure you switch peer context between orgs!"
    log_warning "This script assumes you handle org context switching manually."
    
    echo ""
    read -p "Deploy for StartupOrg? (y/n): " confirm
    if [ "$confirm" == "y" ]; then
        deploy_startup_org
    fi
    
    echo ""
    read -p "Deploy for ValidatorOrg? (y/n): " confirm
    if [ "$confirm" == "y" ]; then
        deploy_validator_org
    fi
    
    echo ""
    read -p "Deploy for InvestorOrg? (y/n): " confirm
    if [ "$confirm" == "y" ]; then
        deploy_investor_org
    fi
    
    echo ""
    read -p "Deploy for PlatformOrg? (y/n): " confirm
    if [ "$confirm" == "y" ]; then
        deploy_platform_org
    fi
    
    log_success "All deployments complete!"
}

# =============================================================================
# Usage Information
# =============================================================================

show_usage() {
    echo ""
    echo "Usage: $0 <org> [channel]"
    echo ""
    echo "Organizations:"
    echo "  startup     - Deploy chaincodes for StartupOrg"
    echo "  validator   - Deploy chaincodes for ValidatorOrg"
    echo "  investor    - Deploy chaincodes for InvestorOrg"
    echo "  platform    - Deploy chaincodes for PlatformOrg"
    echo "  all         - Deploy for all organizations (interactive)"
    echo ""
    echo "Optional Channel Filters:"
    echo "  startup-validator   - Only startup-validator-channel"
    echo "  startup-platform    - Only startup-platform-channel"
    echo "  startup-investor    - Only startup-investor-channel"
    echo "  investor-platform   - Only investor-platform-channel"
    echo "  investor-validator  - Only investor-validator-channel"
    echo "  validator-platform  - Only validator-platform-channel"
    echo "  common              - Only common-channel"
    echo ""
    echo "Prerequisites:"
    echo "  Export package IDs before running:"
    echo "    export STARTUP_CC_PACKAGE_ID=<package_id>"
    echo "    export VALIDATOR_CC_PACKAGE_ID=<package_id>"
    echo "    export INVESTOR_CC_PACKAGE_ID=<package_id>"
    echo "    export PLATFORM_CC_PACKAGE_ID=<package_id>"
    echo ""
    echo "Examples:"
    echo "  $0 startup                    # All channels for StartupOrg"
    echo "  $0 startup common             # Only common-channel for StartupOrg"
    echo "  $0 validator startup-validator # Only startup-validator-channel for ValidatorOrg"
    echo "  $0 all                        # Interactive deployment for all orgs"
    echo ""
}

# =============================================================================
# Main Entry Point
# =============================================================================

main() {
    local org=$1
    local channel=$2
    
    if [ -z "$org" ]; then
        show_usage
        exit 1
    fi
    
    # Set package IDs
    set_package_ids
    
    case $org in
        startup)
            deploy_startup_org "$channel"
            ;;
        validator)
            deploy_validator_org "$channel"
            ;;
        investor)
            deploy_investor_org "$channel"
            ;;
        platform)
            deploy_platform_org "$channel"
            ;;
        all)
            deploy_all_orgs
            ;;
        help|--help|-h)
            show_usage
            ;;
        *)
            log_error "Unknown organization: $org"
            show_usage
            exit 1
            ;;
    esac
}

# Run main function
main "$@"
