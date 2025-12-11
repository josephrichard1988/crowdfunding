#!/bin/bash
#
# export MICROFAB_CONFIG=$(cat MICROFAB.txt)
# docker run --name microfab -e MICROFAB_CONFIG -p 9090:9090 ibmcom/ibp-microfab
#
# curl -s http://console.127-0-0-1.nip.io:9090/ak/api/v1/components | weft microfab -w ./_wallets -p ./_gateways -m ./_msp -f
# Installing Binaries: curl -sSL https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh | bash -s -- binary
# =============================================================================
# Chaincode Deployment Script for Crowdfunding Platform
# =============================================================================
#
# FEATURES:
# ---------
# 1. Organization Context Switching
#    - switch_to_startup/validator/platform/investor functions
#    - Automatically exports CORE_PEER_LOCALMSPID, CORE_PEER_MSPCONFIGPATH, CORE_PEER_ADDRESS
#    - Sets up PATH and FABRIC_CFG_PATH after each org switch
#    - Uses global MSP_BASE_PATH variable for MSP config paths
#
# 2. Chaincode Packaging with Auto-Versioning
#    - package_chaincode function with auto-incrementing version
#    - Detects existing .tgz files and increments version (startup_1.tgz -> startup_2.tgz)
#    - Uses global CONTRACTS_BASE_PATH for contract paths
#    - package_all_chaincodes packages all 4 chaincodes at once
#
# 3. Interactive Chaincode Installation
#    - install_chaincode_interactive prompts for package ID after each install
#    - install_on_startup/validator/platform/investor installs all chaincodes per org
#    - Automatically finds latest package files
#    - Displays export command for package ID
#
# 4. Approve and Commit with Auto-Sequence Detection
#    - Auto-detects current committed sequence/version
#    - Increments sequence and version automatically
#
# GLOBAL VARIABLES:
# -----------------
#   MSP_BASE_PATH       - Base path for MSP configs (/home/kajal/crowdfunding/_msp)
#   CONTRACTS_BASE_PATH - Base path for contracts (./contracts)
#   ORDERER_URL         - Orderer API URL
#
# =============================================================================
# Usage:
#   source ./deploy_chaincode.sh <command> [options]
#
# Commands:
#   package              - Package all chaincodes with auto-versioning
#   install <org>        - Install chaincodes on specific org (interactive)
#   install-all          - Install chaincodes on all orgs (interactive)
#   deploy <org>         - Deploy (approve/commit) for specific org
#   switch <org>         - Switch to organization context
#
# Switch <org> options:
#  startup: source ./deploy_chaincode.sh switch startup
#  validator: source ./deploy_chaincode.sh switch validator
#  platform: source ./deploy_chaincode.sh switch platform
#  investor: source ./deploy_chaincode.sh switch investor
#
# Examples:
#   source ./deploy_chaincode.sh package              # Package all chaincodes
#   source ./deploy_chaincode.sh install startup      # Install on StartupOrg
#   source ./deploy_chaincode.sh install-all          # Install on all orgs
#   source ./deploy_chaincode.sh deploy startup       # Deploy for StartupOrg
#   source ./deploy_chaincode.sh deploy startup common # Deploy only common-channel
#   source ./deploy_chaincode.sh switch validator     # Switch to ValidatorOrg context
#
# Workflow:
#   1. Package chaincodes:  source ./deploy_chaincode.sh package
#                           source ./deploy_chaincode.sh package startup
#                           source ./deploy_chaincode.sh package validator
#                           source ./deploy_chaincode.sh package platform
#                           source ./deploy_chaincode.sh package investor
#   2. Install per org:     source ./deploy_chaincode.sh install startup
#                           source ./deploy_chaincode.sh install validator
#                           source ./deploy_chaincode.sh install platform
#                           source ./deploy_chaincode.sh install investor
#   3. Export package IDs after each installation (prompted automatically)
#   4. Deploy per org:      source ./deploy_chaincode.sh deploy startup
#                           source ./deploy_chaincode.sh deploy validator
#                           source ./deploy_chaincode.sh deploy platform
#                           source ./deploy_chaincode.sh deploy investor
# =============================================================================

# Disable exit on error for better error handling
# set -e

# =============================================================================
# Global Configuration Variables
# =============================================================================
ORDERER_URL="orderer-api.127-0-0-1.nip.io:9090"
MSP_BASE_PATH="/home/kajal/crowdfunding/_msp"
CONTRACTS_BASE_PATH="./contracts"

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
# Organization Context Switching Functions
# =============================================================================

# Setup Fabric environment paths
setup_fabric_env() {
    export PATH=$PATH:${PWD}/bin
    export FABRIC_CFG_PATH=${PWD}/config
    log_info "Fabric environment paths configured"
}

# Switch to StartupOrg context
switch_to_startup() {
    export CORE_PEER_LOCALMSPID=StartupOrgMSP
    export CORE_PEER_MSPCONFIGPATH=${MSP_BASE_PATH}/StartupOrg/startuporgadmin/msp
    export CORE_PEER_ADDRESS=startuporgpeer-api.127-0-0-1.nip.io:9090
    setup_fabric_env
    log_success "Switched to StartupOrg context"
}

# Switch to ValidatorOrg context
switch_to_validator() {
    export CORE_PEER_LOCALMSPID=ValidatorOrgMSP
    export CORE_PEER_MSPCONFIGPATH=${MSP_BASE_PATH}/ValidatorOrg/validatororgadmin/msp
    export CORE_PEER_ADDRESS=validatororgpeer-api.127-0-0-1.nip.io:9090
    setup_fabric_env
    log_success "Switched to ValidatorOrg context"
}

# Switch to PlatformOrg context
switch_to_platform() {
    export CORE_PEER_LOCALMSPID=PlatformOrgMSP
    export CORE_PEER_MSPCONFIGPATH=${MSP_BASE_PATH}/PlatformOrg/platformorgadmin/msp
    export CORE_PEER_ADDRESS=platformorgpeer-api.127-0-0-1.nip.io:9090
    setup_fabric_env
    log_success "Switched to PlatformOrg context"
}

# Switch to InvestorOrg context
switch_to_investor() {
    export CORE_PEER_LOCALMSPID=InvestorOrgMSP
    export CORE_PEER_MSPCONFIGPATH=${MSP_BASE_PATH}/InvestorOrg/investororgadmin/msp
    export CORE_PEER_ADDRESS=investororgpeer-api.127-0-0-1.nip.io:9090
    setup_fabric_env
    log_success "Switched to InvestorOrg context"
}

# Generic org switcher
switch_org() {
    local org=$1
    case $org in
        startup)
            switch_to_startup
            ;;
        validator)
            switch_to_validator
            ;;
        platform)
            switch_to_platform
            ;;
        investor)
            switch_to_investor
            ;;
        *)
            log_error "Unknown organization: $org"
            return 1
            ;;
    esac
}

# =============================================================================
# Chaincode Packaging Functions
# =============================================================================

# Get next package version by checking existing .tgz files
get_next_package_version() {
    local chaincode_name=$1
    local max_version=0
    
    # Find all existing packages matching pattern
    shopt -s nullglob
    for file in ${chaincode_name}_*.tgz; do
        if [ -f "$file" ]; then
            # Extract version number from filename (e.g., startup_3.tgz -> 3)
            version=$(echo "$file" | sed -n "s/${chaincode_name}_\([0-9]\+\)\.tgz/\1/p")
            if [ -n "$version" ] && [ "$version" -gt "$max_version" ]; then
                max_version=$version
            fi
        fi
    done
    shopt -u nullglob
    
    echo $((max_version + 1))
}

# Package a single chaincode with auto-incrementing version
package_chaincode() {
    local chaincode_name=$1
    local contract_path="${CONTRACTS_BASE_PATH}/${chaincode_name}org"
    
    # Get next version
    local version=$(get_next_package_version "$chaincode_name")
    local package_file="${chaincode_name}_${version}.tgz"
    local label="${chaincode_name}_${version}"
    
    log_section "Packaging ${chaincode_name} chaincode"
    log_info "Version: ${version}"
    log_info "Package: ${package_file}"
    log_info "Label: ${label}"
    log_info "Path: ${contract_path}"
    
    # Check if contract path exists
    if [ ! -d "$contract_path" ]; then
        log_error "Contract path does not exist: $contract_path"
        return 1
    fi
    
    # Package the chaincode
    peer lifecycle chaincode package "$package_file" \
        --path "$contract_path" \
        --lang golang \
        --label "$label"
    
    if [ $? -eq 0 ]; then
        log_success "Packaged ${chaincode_name} as ${package_file}"
        echo "$package_file"  # Return package filename
    else
        log_error "Failed to package ${chaincode_name}"
        return 1
    fi
}

# Package all chaincodes
package_all_chaincodes() {
    log_section "Packaging All Chaincodes"
    
    local startup_pkg=$(package_chaincode "startup")
    local validator_pkg=$(package_chaincode "validator")
    local platform_pkg=$(package_chaincode "platform")
    local investor_pkg=$(package_chaincode "investor")
    
    log_success "All chaincodes packaged successfully!"
    echo ""
    log_info "Package files created:"
    log_info "  - $startup_pkg"
    log_info "  - $validator_pkg"
    log_info "  - $platform_pkg"
    log_info "  - $investor_pkg"
}

# =============================================================================
# Chaincode Installation Functions
# =============================================================================

# Install chaincode on current peer and prompt for package ID
install_chaincode_interactive() {
    local chaincode_name=$1
    local package_file=$2
    
    log_section "Installing ${chaincode_name} chaincode"
    
    if [ ! -f "$package_file" ]; then
        log_error "Package file not found: $package_file"
        return 1
    fi
    
    log_info "Installing ${package_file} on current peer..."
    
    # Install the chaincode
    local install_output=$(peer lifecycle chaincode install "$package_file" 2>&1)
    local install_status=$?
    
    echo "$install_output"
    
    # Check if already installed (extract from error message)
    if echo "$install_output" | grep -q "already successfully installed"; then
        log_warning "${chaincode_name} already installed"
        
        # Extract package ID from error message: "package ID 'startup_4:hash'"
        local package_id=$(echo "$install_output" | grep -oP "package ID '\K[^']+")
        
        if [ -z "$package_id" ]; then
            # Fallback: query installed chaincodes
            log_info "Querying installed chaincodes to get package ID..."
            local query_output=$(peer lifecycle chaincode queryinstalled 2>&1)
            package_id=$(echo "$query_output" | grep "Label: ${chaincode_name}_" | head -1 | grep -oP "Package ID: \K[^,]+")
        fi
        
        if [ ! -z "$package_id" ]; then
            log_info "Package ID: ${package_id}"
            echo ""
            
            # Auto-export the package ID
            local var_name="${chaincode_name^^}_CC_PACKAGE_ID"
            export ${var_name}="${package_id}"
            
            log_success "Auto-exported: ${var_name}=${package_id}"
            echo ""
            log_info "You can also manually export with:"
            echo -e "${CYAN}export ${var_name}=${package_id}${NC}"
            echo ""
            
            # Prompt to continue
            read -p "Press Enter to continue to next chaincode..."
        else
            log_error "Could not extract package ID for ${chaincode_name}"
            return 1
        fi
        
    elif [ $install_status -eq 0 ]; then
        log_success "${chaincode_name} installed successfully!"
        
        # Extract package ID from output
        local package_id=$(echo "$install_output" | grep -o 'Chaincode code package identifier: .*' | cut -d ' ' -f 5)
        
        if [ ! -z "$package_id" ]; then
            log_info "Package ID: ${package_id}"
            echo ""
            
            # Auto-export the package ID
            local var_name="${chaincode_name^^}_CC_PACKAGE_ID"
            export ${var_name}="${package_id}"
            
            log_success "Auto-exported: ${var_name}=${package_id}"
            echo ""
            log_info "You can also manually export with:"
            echo -e "${CYAN}export ${var_name}=${package_id}${NC}"
            echo ""
            
            # Prompt to continue
            read -p "Press Enter to continue to next chaincode..."
        fi
    else
        log_error "Failed to install ${chaincode_name}"
        return 1
    fi
}

# Install chaincode on StartupOrg peer
install_on_startup() {
    log_section "Installing Chaincodes on StartupOrg"
    switch_to_startup
    
    # Find latest package files
    local startup_pkg=$(ls -t startup_*.tgz 2>/dev/null | head -1)
    local validator_pkg=$(ls -t validator_*.tgz 2>/dev/null | head -1)
    local platform_pkg=$(ls -t platform_*.tgz 2>/dev/null | head -1)
    local investor_pkg=$(ls -t investor_*.tgz 2>/dev/null | head -1)
    
    install_chaincode_interactive "startup" "$startup_pkg"
    install_chaincode_interactive "validator" "$validator_pkg"
    install_chaincode_interactive "platform" "$platform_pkg"
    install_chaincode_interactive "investor" "$investor_pkg"
    
    log_success "StartupOrg installation complete!"
}

# Install chaincode on ValidatorOrg peer
install_on_validator() {
    log_section "Installing Chaincodes on ValidatorOrg"
    switch_to_validator
    
    local startup_pkg=$(ls -t startup_*.tgz 2>/dev/null | head -1)
    local validator_pkg=$(ls -t validator_*.tgz 2>/dev/null | head -1)
    local platform_pkg=$(ls -t platform_*.tgz 2>/dev/null | head -1)
    local investor_pkg=$(ls -t investor_*.tgz 2>/dev/null | head -1)
    
    install_chaincode_interactive "startup" "$startup_pkg"
    install_chaincode_interactive "validator" "$validator_pkg"
    install_chaincode_interactive "platform" "$platform_pkg"
    install_chaincode_interactive "investor" "$investor_pkg"
    
    log_success "ValidatorOrg installation complete!"
}

# Install chaincode on PlatformOrg peer
install_on_platform() {
    log_section "Installing Chaincodes on PlatformOrg"
    switch_to_platform
    
    local startup_pkg=$(ls -t startup_*.tgz 2>/dev/null | head -1)
    local validator_pkg=$(ls -t validator_*.tgz 2>/dev/null | head -1)
    local platform_pkg=$(ls -t platform_*.tgz 2>/dev/null | head -1)
    local investor_pkg=$(ls -t investor_*.tgz 2>/dev/null | head -1)
    
    install_chaincode_interactive "startup" "$startup_pkg"
    install_chaincode_interactive "validator" "$validator_pkg"
    install_chaincode_interactive "platform" "$platform_pkg"
    install_chaincode_interactive "investor" "$investor_pkg"
    
    log_success "PlatformOrg installation complete!"
}

# Install chaincode on InvestorOrg peer
install_on_investor() {
    log_section "Installing Chaincodes on InvestorOrg"
    switch_to_investor
    
    local startup_pkg=$(ls -t startup_*.tgz 2>/dev/null | head -1)
    local validator_pkg=$(ls -t validator_*.tgz 2>/dev/null | head -1)
    local platform_pkg=$(ls -t platform_*.tgz 2>/dev/null | head -1)
    local investor_pkg=$(ls -t investor_*.tgz 2>/dev/null | head -1)
    
    install_chaincode_interactive "startup" "$startup_pkg"
    install_chaincode_interactive "validator" "$validator_pkg"
    install_chaincode_interactive "platform" "$platform_pkg"
    install_chaincode_interactive "investor" "$investor_pkg"
    
    log_success "InvestorOrg installation complete!"
}

# Install on all orgs sequentially with prompts
install_all_orgs() {
    log_section "Installing Chaincodes on All Organizations"
    
    echo ""
    read -p "Install on StartupOrg? (y/n): " confirm
    if [ "$confirm" == "y" ]; then
        install_on_startup
    fi
    
    echo ""
    read -p "Install on ValidatorOrg? (y/n): " confirm
    if [ "$confirm" == "y" ]; then
        install_on_validator
    fi
    
    echo ""
    read -p "Install on PlatformOrg? (y/n): " confirm
    if [ "$confirm" == "y" ]; then
        install_on_platform
    fi
    
    echo ""
    read -p "Install on InvestorOrg? (y/n): " confirm
    if [ "$confirm" == "y" ]; then
        install_on_investor
    fi
    
    log_success "Installation on all orgs complete!"
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
    switch_to_startup
    
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
    switch_to_validator
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
    switch_to_investor
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
    switch_to_platform
    
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
show_usage() {
    echo ""
    echo "Usage: $0 <command> [options]"
    echo ""
    echo "Commands:"
    echo "  package [chaincode]  - Package all chaincodes or specific chaincode with auto-versioning"
    echo "  install <org>        - Install chaincodes on specific org"
    echo "  install-all          - Install chaincodes on all orgs (interactive)"
    echo "  deploy <org>         - Deploy (approve/commit) for specific org"
    echo "  switch <org>         - Switch to organization context"
    echo ""
    echo "Chaincodes (for package command):"
    echo "  startup     - StartupOrg chaincode"
    echo "  validator   - ValidatorOrg chaincode"
    echo "  investor    - InvestorOrg chaincode"
    echo "  platform    - PlatformOrg chaincode"
    echo "  (no arg)    - Package all chaincodes"
    echo ""
    echo "Organizations:"
    echo "  startup     - StartupOrg"
    echo "  validator   - ValidatorOrg"
    echo "  investor    - InvestorOrg"
    echo "  platform    - PlatformOrg"
    echo "  all         - All organizations (interactive)"
    echo ""
    echo "  common              - Only common-channel"
    echo ""
    echo "Workflow:"
    echo "  1. Package chaincodes:    $0 package [chaincode]"
    echo "  2. Install on each org:   $0 install startup"
    echo "                            $0 install validator"
    echo "                            $0 install platform"
    echo "                            $0 install investor"
    echo "     (Or use:               $0 install-all)"
    echo "  3. Export package IDs after each installation"
    echo "  4. Deploy on each org:    $0 deploy startup"
    echo "                            $0 deploy validator"
    echo "                            $0 deploy platform"
    echo "                            $0 deploy investor"
    echo ""
    echo "Examples:"
    echo "  $0 package                       # Package all chaincodes"
    echo "  $0 package validator             # Package only validator chaincode"
    echo "  $0 install startup               # Install on StartupOrg"
    echo "  $0 install-all                   # Install on all orgs"
    echo "  $0 deploy startup                # Deploy all channels for StartupOrg"
    echo "  $0 deploy startup common         # Deploy only common-channel"
    echo "  $0 switch validator              # Switch to ValidatorOrg context"
    echo ""
}

# =============================================================================
# Main Entry Point
# =============================================================================

main() {
    local command=$1
    local arg1=$2
    local arg2=$3
    
    if [ -z "$command" ]; then
        show_usage
        exit 1
    fi
    
    case $command in
        package)
            if [ -z "$arg1" ]; then
                # No argument - package all chaincodes
                package_all_chaincodes
            else
                # Specific chaincode specified
                case $arg1 in
                    startup|validator|platform|investor)
                        package_chaincode "$arg1"
                        ;;
                    *)
                        log_error "Unknown chaincode: $arg1"
                        log_info "Valid chaincodes: startup, validator, platform, investor"
                        exit 1
                        ;;
                esac
            fi
            ;;
        install)
            if [ -z "$arg1" ]; then
                log_error "Please specify organization: startup, validator, platform, investor"
                exit 1
            fi
            case $arg1 in
                startup)
                    install_on_startup
                    ;;
                validator)
                    install_on_validator
                    ;;
                platform)
                    install_on_platform
                    ;;
                investor)
                    install_on_investor
                    ;;
                *)
                    log_error "Unknown organization: $arg1"
                    exit 1
                    ;;
            esac
            ;;
        install-all)
            install_all_orgs
            ;;
        deploy)
            if [ -z "$arg1" ]; then
                log_error "Please specify organization for deployment"
                exit 1
            fi
            # Set package IDs
            set_package_ids
            case $arg1 in
                startup)
                    deploy_startup_org "$arg2"
                    ;;
                validator)
                    deploy_validator_org "$arg2"
                    ;;
                investor)
                    deploy_investor_org "$arg2"
                    ;;
                platform)
                    deploy_platform_org "$arg2"
                    ;;
                all)
                    deploy_all_orgs
                    ;;
                *)
                    log_error "Unknown organization: $arg1"
                    exit 1
                    ;;
            esac
            ;;
        switch)
            if [ -z "$arg1" ]; then
                log_error "Please specify organization to switch to"
                exit 1
            fi
            switch_org "$arg1"
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
