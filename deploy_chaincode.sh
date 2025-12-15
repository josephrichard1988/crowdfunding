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
#    - Smart approval: only approves if chaincode is already committed
#
# 5. Sync-Upgrade Mechanism
#    - sync_upgrade_chaincode: Detects if chaincode is committed, then approves
#    - upgrade_all: Upgrades all 4 chaincodes across all organizations and channels
#    - Smart upgrade with automatic sequence detection
#
# 6. Complete Channel Deployment
#    - Deployment functions for all 7 channels
#    - Proper peer address handling for multi-org channels
#    - Commit readiness checking
#
# GLOBAL VARIABLES:
# -----------------
#   MSP_BASE_PATH       - Base path for MSP configs (/home/kajal/crowdfunding/_msp)
#   CONTRACTS_BASE_PATH - Base path for contracts (./contracts)
#   ORDERER_URL         - Orderer API URL
#
# =============================================================================
# AVAILABLE COMMANDS:
# =============================================================================
#
# PACKAGING COMMANDS:
#   ./deploy_chaincode.sh package                    - Package all 4 chaincodes (auto-version)
#   ./deploy_chaincode.sh package <chaincode>        - Package specific chaincode only
#                                                      (chaincode: startup|validator|investor|platform)
#
# INSTALLATION COMMANDS:
#   ./deploy_chaincode.sh install <org>              - Install all chaincodes on specific org
#                                                      (org: startup|validator|investor|platform)
#   ./deploy_chaincode.sh install-all                - Install all chaincodes on all 4 orgs (interactive)
#
# DEPLOYMENT COMMANDS (Initial):
#   ./deploy_chaincode.sh deploy <org>               - Deploy all channels for specific org
#   ./deploy_chaincode.sh deploy <org> <channel>     - Deploy specific channel only for org
#                                                      (channel: common|startup-validator|startup-investor|
#                                                       startup-platform|investor-validator|investor-platform|
#                                                       validator-platform)
#   ./deploy_chaincode.sh deploy all                 - Deploy for all orgs (interactive)
#
# UPGRADE COMMANDS:
#   ./deploy_chaincode.sh upgrade <cc> <channel>     - Upgrade specific chaincode on specific channel
#                                                      (cc: startup|validator|investor|platform)
#                                                      Auto-approves from all orgs on that channel
#   ./deploy_chaincode.sh upgrade-all                - Upgrade all chaincodes on all 7 channels
#   ./deploy_chaincode.sh sync-upgrade <cc> <ch> <org> - Smart upgrade for specific org
#                                                        (only if chaincode already committed)
#   ./deploy_chaincode.sh approve-chaincode <org> <cc> <ch> - Manually approve chaincode upgrade
#                                                              Use when another org upgraded and you need to approve
#                                                              Example: After StartupOrg upgrades, ValidatorOrg approves
#
# QUERY & CHECK COMMANDS:
#   ./deploy_chaincode.sh query-committed <cc> <ch>  - Query committed chaincode details
#   ./deploy_chaincode.sh check-readiness <cc> <ch>  - Check if chaincode ready to commit
#
# UTILITY COMMANDS:
#   ./deploy_chaincode.sh switch <org>               - Switch peer context to specific org
#   ./deploy_chaincode.sh help                       - Show detailed usage information
#
# =============================================================================
# DEPLOYMENT FLOW - INITIAL DEPLOYMENT:
# =============================================================================
#
# STEP 1: Package All Chaincodes
#   Command: ./deploy_chaincode.sh package
#   Output:  - Creates startup_1.tgz, validator_1.tgz, investor_1.tgz, platform_1.tgz
#            - Shows version numbers and package labels
#            - Each package contains the Go chaincode from ./contracts/<name>org/
#
# STEP 2: Install Chaincodes on All Organizations
#   Command: ./deploy_chaincode.sh install-all
#           (OR individually: ./deploy_chaincode.sh install startup)
#   Output:  For each org (StartupOrg, ValidatorOrg, InvestorOrg, PlatformOrg):
#            - Switches to org context
#            - Installs all 4 chaincode packages on the org's peer
#            - Auto-extracts and exports package IDs as environment variables:
#              * STARTUP_CC_PACKAGE_ID=startup_1:abc123...
#              * VALIDATOR_CC_PACKAGE_ID=validator_1:def456...
#              * INVESTOR_CC_PACKAGE_ID=investor_1:ghi789...
#              * PLATFORM_CC_PACKAGE_ID=platform_1:jkl012...
#            - Prompts to press Enter after each chaincode installation
#            - Package IDs are automatically available for next steps
#
# STEP 3: Deploy on Each Organization (Approve & Commit)
#   Command: ./deploy_chaincode.sh deploy startup
#   Output:  For StartupOrg:
#            - Approves startup chaincode on: startup-validator-channel, 
#              startup-platform-channel, startup-investor-channel, common-channel
#            - Commits startup chaincode on all 4 channels (as owner)
#            - Approves other orgs' chaincodes: validator, investor, platform
#            Shows: Version 1, Sequence 1, Channel names, Approval status, Commit status
#
#   Command: ./deploy_chaincode.sh deploy validator
#   Output:  For ValidatorOrg:
#            - Approves validator chaincode on: startup-validator-channel,
#              investor-validator-channel, validator-platform-channel, common-channel
#            - Commits validator chaincode on all 4 channels (as owner)
#            - Approves other orgs' chaincodes
#
#   Command: ./deploy_chaincode.sh deploy investor
#   Output:  For InvestorOrg:
#            - Approves investor chaincode on: startup-investor-channel,
#              investor-validator-channel, investor-platform-channel, common-channel
#            - Commits investor chaincode on all 4 channels (as owner)
#            - Approves other orgs' chaincodes
#
#   Command: ./deploy_chaincode.sh deploy platform
#   Output:  For PlatformOrg:
#            - Approves platform chaincode on: startup-platform-channel,
#              investor-platform-channel, validator-platform-channel, common-channel
#            - Commits platform chaincode on all 4 channels (as owner)
#            - Approves other orgs' chaincodes
#
#   Result: All 4 chaincodes deployed on all 7 channels with proper endorsement policies
#
# STEP 4: Verify Deployment
#   Command: ./deploy_chaincode.sh query-committed startup common-channel
#   Output:  - Shows: Version 1, Sequence 1, Endorsement Policy, Approvals from all orgs
#            - Confirms chaincode is active and ready for invocation
#
# =============================================================================
# UPGRADE FLOW - UPGRADE EXISTING CHAINCODES:
# =============================================================================
#
# STEP 1: Package New Versions
#   Command: ./deploy_chaincode.sh package
#   Output:  - Auto-detects existing versions (startup_1.tgz -> startup_2.tgz)
#            - Creates startup_2.tgz, validator_2.tgz, investor_2.tgz, platform_2.tgz
#            - Shows new version numbers
#
# STEP 2: Install New Versions on All Organizations
#   Command: ./deploy_chaincode.sh install-all
#   Output:  - Installs new package versions on all org peers
#            - Auto-exports new package IDs (startup_2:xyz789...)
#            - Overwrites previous package ID environment variables
#
# STEP 3: Upgrade All Chaincodes (Automated)
#   Command: ./deploy_chaincode.sh upgrade-all
#   Output:  For each channel and chaincode:
#            - Auto-detects current sequence (e.g., Sequence 1)
#            - Calculates next sequence (Sequence 2)
#            - Approves upgrade from all orgs on each channel
#            - Shows approval status per org
#            - Commits upgrade with new sequence and version
#            - Processes all 7 channels automatically:
#              * common-channel (4 chaincodes, 4 orgs each)
#              * startup-investor-channel (2 chaincodes, 2 orgs each)
#              * startup-validator-channel (2 chaincodes, 2 orgs each)
#              * startup-platform-channel (2 chaincodes, 2 orgs each)
#              * investor-validator-channel (2 chaincodes, 2 orgs each)
#              * investor-platform-channel (2 chaincodes, 2 orgs each)
#              * validator-platform-channel (2 chaincodes, 2 orgs each)
#            - Displays success message for each upgrade
#
# ALTERNATIVE - Upgrade Specific Chaincode on Specific Channel:
#   Command: ./deploy_chaincode.sh upgrade startup common-channel
#   Output:  - Approves upgrade from all 4 orgs (startup, validator, investor, platform)
#            - Auto-increments sequence (1 -> 2)
#            - Commits new version from StartupOrg
#            - Shows sequence progression and commit status
#
# =============================================================================
# EXPECTED OUTPUT EXAMPLES:
# =============================================================================
#
# PACKAGE OUTPUT:
#   📦 Packaging All Chaincodes
#   📦 Packaging chaincode: startup
#      Current version: 1
#      New version: 2
#      Package label: startup_2
#      Package file: startup_2.tgz
#   ✅ Successfully packaged startup as startup_2.tgz
#   [Repeat for validator, investor, platform]
#   ✅ All chaincodes packaged successfully!
#
# INSTALL OUTPUT:
#   📥 Installing Chaincodes on StartupOrg
#   🔄 Switching to StartupOrg...
#   ✅ Now operating as StartupOrg
#   📥 Installing startup_2.tgz on StartupOrg...
#   2024.12.12 10:30:15.123 UTC 0001 INFO [cli.lifecycle.chaincode] submitInstallProposal
#   ✅ Successfully installed startup on StartupOrg
#   🔍 Please query the package ID and export it:
#      peer lifecycle chaincode queryinstalled
#      Then export (example):
#      export STARTUP_CC_PACKAGE_ID="startup_2:abc123def456..."
#   ✅ Auto-exported: STARTUP_CC_PACKAGE_ID=startup_2:abc123def456...
#   Press Enter to continue to next chaincode...
#   [Repeat for all 4 chaincodes]
#
# DEPLOY OUTPUT:
#   🚀 Deploying chaincodes for StartupOrg
#   ============================================================================
#   StartupOrg: startup-validator-channel
#   ============================================================================
#   [INFO] Approving chaincode 'startup' on channel 'startup-validator-channel' 
#          (version: 1, sequence: 1)...
#   2024.12.12 10:35:20.456 UTC 0001 INFO [chaincodeCmd] ClientWait
#   ✅ Successfully approved 'startup' on 'startup-validator-channel' 
#      (version: 1, sequence: 1)
#   [INFO] Committing chaincode 'startup' on channel 'startup-validator-channel' 
#          (version: 1, sequence: 1)...
#   2024.12.12 10:35:25.789 UTC 0001 INFO [chaincodeCmd] ClientWait
#   ✅ Successfully committed 'startup' on 'startup-validator-channel' 
#      (version: 1, sequence: 1)
#   [Repeat for other channels and chaincodes]
#
# QUERY-COMMITTED OUTPUT:
#   [INFO] Querying committed chaincode 'startup' on 'common-channel'...
#   Committed chaincode definition for chaincode 'startup' on channel 'common-channel':
#   Version: 2, Sequence: 2, Endorsement Plugin: escc, Validation Plugin: vscc
#   Approvals: [StartupOrgMSP: true, ValidatorOrgMSP: true, InvestorOrgMSP: true, 
#               PlatformOrgMSP: true]
#
# UPGRADE-ALL OUTPUT:
#   ============================================================================
#   Comprehensive Upgrade of All Chaincodes on All Channels
#   ============================================================================
#   This will upgrade all chaincodes on all 7 channels.
#   Make sure all package IDs are exported:
#     - STARTUP_CC_PACKAGE_ID
#     - VALIDATOR_CC_PACKAGE_ID
#     - INVESTOR_CC_PACKAGE_ID
#     - PLATFORM_CC_PACKAGE_ID
#   Continue with upgrade? (y/n): y
#   ============================================================================
#   Upgrading common-channel
#   ============================================================================
#   [INFO] Smart approve for 'startup' on 'common-channel' by startup...
#   [INFO] startup is committed (sequence: 1). Proceeding with approval...
#   ✅ Approved upgrade for 'startup' on 'common-channel' (sequence: 2)
#   [Repeat for all orgs, then commit]
#   ✅ Committed 'startup' on 'common-channel' (version: 2, sequence: 2)
#   [Process continues for all channels and chaincodes]
#   ✅ All chaincodes upgraded on all channels!
#
# =============================================================================
# TYPICAL WORKFLOW EXECUTION TIME:
# =============================================================================
#   Package:       ~10 seconds (all 4 chaincodes)
#   Install-all:   ~2-3 minutes (4 orgs × 4 chaincodes with prompts)
#   Deploy all:    ~5-7 minutes (4 orgs × multiple channels)
#   Upgrade-all:   ~10-15 minutes (7 channels × 4 chaincodes × multiple orgs)
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
    
    # Build peer addresses based on channel
    local peer_addresses=""
    case "$channel" in
        "common-channel")
            peer_addresses="--peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses validatororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses investororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:9090"
            ;;
        "startup-investor-channel")
            peer_addresses="--peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses investororgpeer-api.127-0-0-1.nip.io:9090"
            ;;
        "startup-validator-channel")
            peer_addresses="--peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses validatororgpeer-api.127-0-0-1.nip.io:9090"
            ;;
        "startup-platform-channel")
            peer_addresses="--peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:9090"
            ;;
        "investor-validator-channel")
            peer_addresses="--peerAddresses investororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses validatororgpeer-api.127-0-0-1.nip.io:9090"
            ;;
        "investor-platform-channel")
            peer_addresses="--peerAddresses investororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:9090"
            ;;
        "validator-platform-channel")
            peer_addresses="--peerAddresses validatororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:9090"
            ;;
    esac
    
    peer lifecycle chaincode commit \
        -o ${ORDERER_URL} \
        --channelID ${channel} \
        --name ${chaincode_name} \
        --version ${version} \
        --sequence ${sequence} \
        ${peer_addresses} \
        --waitForEvent
    
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

# Smart approve - only approves if chaincode is already committed (for upgrades)
smart_approve() {
    local channel=$1
    local chaincode_name=$2
    local package_id=$3
    local org_name=$4
    
    log_info "Smart approve for '${chaincode_name}' on '${channel}' by ${org_name}..."
    
    # Check if chaincode is committed
    local current_sequence=$(get_committed_sequence "$channel" "$chaincode_name")
    
    if [ "$current_sequence" == "0" ]; then
        log_warning "${chaincode_name} not yet committed on ${channel}. Skipping approval for ${org_name}."
        return 0
    fi
    
    log_info "${chaincode_name} is committed (sequence: ${current_sequence}). Proceeding with approval for ${org_name}..."
    
    # Get next sequence and version
    local next_sequence=$((current_sequence + 1))
    local next_version=$(get_next_version "$channel" "$chaincode_name")
    
    # Approve the upgrade
    peer lifecycle chaincode approveformyorg \
        -o ${ORDERER_URL} \
        --channelID ${channel} \
        --name ${chaincode_name} \
        --version ${next_version} \
        --sequence ${next_sequence} \
        --waitForEvent \
        --package-id ${package_id}
    
    if [ $? -eq 0 ]; then
        log_success "Approved upgrade for '${chaincode_name}' on '${channel}' (sequence: ${next_sequence})"
    else
        log_error "Failed to approve upgrade for '${chaincode_name}' on '${channel}'"
        return 1
    fi
}

# Query committed chaincode details
query_committed() {
    local channel=$1
    local chaincode_name=$2
    
    log_info "Querying committed chaincode '${chaincode_name}' on '${channel}'..."
    
    peer lifecycle chaincode querycommitted \
        --channelID ${channel} \
        --name ${chaincode_name}
}

# =============================================================================
# Sync-Upgrade Functions
# =============================================================================

# Sync-upgrade: Smart upgrade for a specific chaincode on a channel for an org
sync_upgrade_chaincode() {
    local chaincode_name=$1
    local channel=$2
    local org_name=$3
    
    log_section "Sync-Upgrade: ${chaincode_name} on ${channel} for ${org_name}"
    
    # Switch to the org
    switch_org "$org_name"
    
    # Get package ID variable name
    local package_id_var="${chaincode_name^^}_CC_PACKAGE_ID"
    local package_id=${!package_id_var}
    
    if [ -z "$package_id" ]; then
        log_error "Package ID not set: ${package_id_var}"
        return 1
    fi
    
    # Use smart approve (only approves if already committed)
    smart_approve "$channel" "$chaincode_name" "$package_id" "$org_name"
}

# Upgrade all chaincodes on all channels
upgrade_all() {
    log_section "Comprehensive Upgrade of All Chaincodes on All Channels"
    
    log_info "This will upgrade all chaincodes on all 7 channels."
    log_warning "Make sure all package IDs are exported:"
    log_warning "  - STARTUP_CC_PACKAGE_ID"
    log_warning "  - VALIDATOR_CC_PACKAGE_ID"
    log_warning "  - INVESTOR_CC_PACKAGE_ID"
    log_warning "  - PLATFORM_CC_PACKAGE_ID"
    echo ""
    read -p "Continue with upgrade? (y/n): " confirm
    
    if [ "$confirm" != "y" ]; then
        log_info "Upgrade cancelled."
        return 0
    fi
    
    # ========== COMMON CHANNEL (All 4 orgs, all 4 chaincodes) ==========
    log_section "Upgrading common-channel"
    
    # Startup chaincode on common-channel
    sync_upgrade_chaincode "startup" "common-channel" "startup"
    sync_upgrade_chaincode "startup" "common-channel" "validator"
    sync_upgrade_chaincode "startup" "common-channel" "investor"
    sync_upgrade_chaincode "startup" "common-channel" "platform"
    switch_to_startup
    commit_chaincode "common-channel" "startup"
    
    # Validator chaincode on common-channel
    sync_upgrade_chaincode "validator" "common-channel" "startup"
    sync_upgrade_chaincode "validator" "common-channel" "validator"
    sync_upgrade_chaincode "validator" "common-channel" "investor"
    sync_upgrade_chaincode "validator" "common-channel" "platform"
    switch_to_validator
    commit_chaincode "common-channel" "validator"
    
    # Investor chaincode on common-channel
    sync_upgrade_chaincode "investor" "common-channel" "startup"
    sync_upgrade_chaincode "investor" "common-channel" "validator"
    sync_upgrade_chaincode "investor" "common-channel" "investor"
    sync_upgrade_chaincode "investor" "common-channel" "platform"
    switch_to_investor
    commit_chaincode "common-channel" "investor"
    
    # Platform chaincode on common-channel
    sync_upgrade_chaincode "platform" "common-channel" "startup"
    sync_upgrade_chaincode "platform" "common-channel" "validator"
    sync_upgrade_chaincode "platform" "common-channel" "investor"
    sync_upgrade_chaincode "platform" "common-channel" "platform"
    switch_to_platform
    commit_chaincode "common-channel" "platform"
    
    # ========== STARTUP-INVESTOR CHANNEL ==========
    log_section "Upgrading startup-investor-channel"
    
    sync_upgrade_chaincode "startup" "startup-investor-channel" "startup"
    sync_upgrade_chaincode "startup" "startup-investor-channel" "investor"
    switch_to_startup
    commit_chaincode "startup-investor-channel" "startup"
    
    sync_upgrade_chaincode "investor" "startup-investor-channel" "startup"
    sync_upgrade_chaincode "investor" "startup-investor-channel" "investor"
    switch_to_investor
    commit_chaincode "startup-investor-channel" "investor"
    
    # ========== STARTUP-VALIDATOR CHANNEL ==========
    log_section "Upgrading startup-validator-channel"
    
    sync_upgrade_chaincode "startup" "startup-validator-channel" "startup"
    sync_upgrade_chaincode "startup" "startup-validator-channel" "validator"
    switch_to_startup
    commit_chaincode "startup-validator-channel" "startup"
    
    sync_upgrade_chaincode "validator" "startup-validator-channel" "startup"
    sync_upgrade_chaincode "validator" "startup-validator-channel" "validator"
    switch_to_validator
    commit_chaincode "startup-validator-channel" "validator"
    
    # ========== STARTUP-PLATFORM CHANNEL ==========
    log_section "Upgrading startup-platform-channel"
    
    sync_upgrade_chaincode "startup" "startup-platform-channel" "startup"
    sync_upgrade_chaincode "startup" "startup-platform-channel" "platform"
    switch_to_startup
    commit_chaincode "startup-platform-channel" "startup"
    
    sync_upgrade_chaincode "platform" "startup-platform-channel" "startup"
    sync_upgrade_chaincode "platform" "startup-platform-channel" "platform"
    switch_to_platform
    commit_chaincode "startup-platform-channel" "platform"
    
    # ========== INVESTOR-VALIDATOR CHANNEL ==========
    log_section "Upgrading investor-validator-channel"
    
    sync_upgrade_chaincode "investor" "investor-validator-channel" "investor"
    sync_upgrade_chaincode "investor" "investor-validator-channel" "validator"
    switch_to_investor
    commit_chaincode "investor-validator-channel" "investor"
    
    sync_upgrade_chaincode "validator" "investor-validator-channel" "investor"
    sync_upgrade_chaincode "validator" "investor-validator-channel" "validator"
    switch_to_validator
    commit_chaincode "investor-validator-channel" "validator"
    
    # ========== INVESTOR-PLATFORM CHANNEL ==========
    log_section "Upgrading investor-platform-channel"
    
    sync_upgrade_chaincode "investor" "investor-platform-channel" "investor"
    sync_upgrade_chaincode "investor" "investor-platform-channel" "platform"
    switch_to_investor
    commit_chaincode "investor-platform-channel" "investor"
    
    sync_upgrade_chaincode "platform" "investor-platform-channel" "investor"
    sync_upgrade_chaincode "platform" "investor-platform-channel" "platform"
    switch_to_platform
    commit_chaincode "investor-platform-channel" "platform"
    
    # ========== VALIDATOR-PLATFORM CHANNEL ==========
    log_section "Upgrading validator-platform-channel"
    
    sync_upgrade_chaincode "validator" "validator-platform-channel" "validator"
    sync_upgrade_chaincode "validator" "validator-platform-channel" "platform"
    switch_to_validator
    commit_chaincode "validator-platform-channel" "validator"
    
    sync_upgrade_chaincode "platform" "validator-platform-channel" "validator"
    sync_upgrade_chaincode "platform" "validator-platform-channel" "platform"
    switch_to_platform
    commit_chaincode "validator-platform-channel" "platform"
    
    log_success "✅ All chaincodes upgraded on all channels!"
}

# Approve chaincode for specific org on specific channel (for manual approval)
# Use this when another org has upgraded a chaincode and you need to approve it
approve_chaincode_for_org() {
    local org=$1
    local chaincode_name=$2
    local channel=$3
    
    if [ -z "$org" ] || [ -z "$chaincode_name" ] || [ -z "$channel" ]; then
        log_error "Usage: approve-chaincode <org> <chaincode> <channel>"
        echo "Example: ./deploy_chaincode.sh approve-chaincode investor startup common-channel"
        return 1
    fi
    
    log_section "Approving ${chaincode_name} chaincode on ${channel} for ${org}"
    
    # Switch to the specified org
    switch_org "$org"
    
    # Get package ID
    local package_id_var="${chaincode_name^^}_CC_PACKAGE_ID"
    local package_id=${!package_id_var}
    
    if [ -z "$package_id" ]; then
        log_error "Package ID not set: ${package_id_var}"
        log_info "Please install chaincode first: ./deploy_chaincode.sh install $org"
        return 1
    fi
    
    # Use smart approve to match committed version
    smart_approve "$channel" "$chaincode_name" "$package_id" "$org"
    
    log_success "Approved ${chaincode_name} on ${channel} for ${org}"
}

# Upgrade specific chaincode on specific channel
upgrade_chaincode() {
    local chaincode_name=$1
    local channel=$2
    
    log_section "Upgrading ${chaincode_name} on ${channel}"
    
    # Get package ID
    local package_id_var="${chaincode_name^^}_CC_PACKAGE_ID"
    local package_id=${!package_id_var}
    
    if [ -z "$package_id" ]; then
        log_error "Package ID not set: ${package_id_var}"
        return 1
    fi
    
    # Determine which orgs are on this channel and approve from each
    case "$channel" in
        "common-channel")
            sync_upgrade_chaincode "$chaincode_name" "$channel" "startup"
            sync_upgrade_chaincode "$chaincode_name" "$channel" "validator"
            sync_upgrade_chaincode "$chaincode_name" "$channel" "investor"
            sync_upgrade_chaincode "$chaincode_name" "$channel" "platform"
            ;;
        "startup-investor-channel")
            sync_upgrade_chaincode "$chaincode_name" "$channel" "startup"
            sync_upgrade_chaincode "$chaincode_name" "$channel" "investor"
            ;;
        "startup-validator-channel")
            sync_upgrade_chaincode "$chaincode_name" "$channel" "startup"
            sync_upgrade_chaincode "$chaincode_name" "$channel" "validator"
            ;;
        "startup-platform-channel")
            sync_upgrade_chaincode "$chaincode_name" "$channel" "startup"
            sync_upgrade_chaincode "$chaincode_name" "$channel" "platform"
            ;;
        "investor-validator-channel")
            sync_upgrade_chaincode "$chaincode_name" "$channel" "investor"
            sync_upgrade_chaincode "$chaincode_name" "$channel" "validator"
            ;;
        "investor-platform-channel")
            sync_upgrade_chaincode "$chaincode_name" "$channel" "investor"
            sync_upgrade_chaincode "$chaincode_name" "$channel" "platform"
            ;;
        "validator-platform-channel")
            sync_upgrade_chaincode "$chaincode_name" "$channel" "validator"
            sync_upgrade_chaincode "$chaincode_name" "$channel" "platform"
            ;;
        *)
            log_error "Unknown channel: $channel"
            return 1
            ;;
    esac
    
    # Commit from the chaincode owner org
    case "$chaincode_name" in
        "startup") switch_to_startup ;;
        "validator") switch_to_validator ;;
        "investor") switch_to_investor ;;
        "platform") switch_to_platform ;;
    esac
    
    commit_chaincode "$channel" "$chaincode_name"
    log_success "Upgrade complete for ${chaincode_name} on ${channel}"
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
    echo "${CYAN}═══════════════════════════════════════════════════════════════════${NC}"
    echo "${YELLOW}  Crowdfunding Platform - Chaincode Deployment Script${NC}"
    echo "${CYAN}═══════════════════════════════════════════════════════════════════${NC}"
    echo ""
    echo "${GREEN}USAGE:${NC} $0 <command> [options]"
    echo ""
    echo "${GREEN}COMMANDS:${NC}"
    echo "  ${CYAN}upgrade${NC} <chaincode> <channel> - Upgrade specific chaincode on specific channel"
    echo "  ${CYAN}upgrade-all${NC}                   - Upgrade all chaincodes on all channels (comprehensive)"
    echo "  ${CYAN}sync-upgrade${NC} <cc> <ch> <org>  - Smart sync-upgrade for specific org (only if committed)"
    echo "  ${CYAN}approve-chaincode${NC} <org> <cc> <ch> - Manually approve chaincode upgrade for specific org"
    echo "                                   ${BLUE}Use when another org upgraded and you need to approve${NC}"
    echo "  ${CYAN}check-readiness${NC} <cc> <ch>     - Check commit readiness for chaincode"
    echo "  ${CYAN}query-committed${NC} <cc> <ch>     - Query committed chaincode details"ific channel"
    echo "  ${CYAN}upgrade-all${NC}                   - Upgrade all chaincodes on all channels (comprehensive)"
    echo "  ${CYAN}sync-upgrade${NC} <cc> <ch> <org>  - Smart sync-upgrade for specific org (only if committed)"
    echo "  ${CYAN}check-readiness${NC} <cc> <ch>     - Check commit readiness for chaincode"
    echo "  ${CYAN}query-committed${NC} <cc> <ch>     - Query committed chaincode details"
    echo "  ${CYAN}switch${NC} <org>                  - Switch to organization context"
    echo "  ${CYAN}help${NC}                          - Show this help message"
    echo ""
    echo "${GREEN}ORGANIZATIONS:${NC}"
    echo "  startup, validator, investor, platform, all"
    echo ""
    echo "${GREEN}CHAINCODES:${NC}"
    echo "  startup, validator, investor, platform"
    echo ""
    echo "${GREEN}CHANNELS:${NC}"
    echo "  common-channel, startup-investor-channel, startup-validator-channel,"
    echo "  startup-platform-channel, investor-validator-channel,"
    echo "  investor-platform-channel, validator-platform-channel"
    echo ""
    echo "${GREEN}INITIAL DEPLOYMENT WORKFLOW:${NC}"
    echo "  ${YELLOW}1.${NC} Package chaincodes:       $0 package"
    echo "  ${YELLOW}2.${NC} Install on all orgs:      $0 install-all"
    echo "     ${BLUE}(Or individually):${NC}       $0 install startup"
    echo "                               $0 install validator"
    echo "                               $0 install investor"
    echo "                               $0 install platform"
    echo "  ${YELLOW}3.${NC} Package IDs are auto-exported during installation"
    echo "  ${YELLOW}4.${NC} Deploy on each org:       $0 deploy startup"
    echo "                               $0 deploy validator"
    echo "                               $0 deploy investor"
    echo "                               $0 deploy platform"
    echo ""
    echo "${GREEN}UPGRADE WORKFLOW:${NC}"
    echo "  ${YELLOW}1.${NC} Package new versions:     $0 package"
    echo "  ${YELLOW}2.${NC} Install on all orgs:      $0 install-all"
    echo "  ${YELLOW}3.${NC} Upgrade all at once:      $0 upgrade-all"
    echo "     ${BLUE}(Or upgrade specific):${NC}    $0 upgrade startup common-channel"
    echo ""
    echo "${GREEN}EXAMPLES:${NC}"
    echo "  ${CYAN}# Initial deployment${NC}"
    echo "  $0 package"
    echo "  $0 install-all"
    echo "  $0 deploy startup"
    echo "  $0 deploy validator"
    echo "  ${CYAN}# Upgrade specific chaincode${NC}"
    echo "  $0 package startup"
    echo "  $0 install startup"
    echo "  $0 upgrade startup common-channel"
    echo ""
    echo "  ${CYAN}# Manual approval (when another org upgraded)${NC}"
    echo "  $0 approve-chaincode investor startup common-channel"
    echo "  $0 approve-chaincode validator platform validator-platform-channel"
    echo ""
    echo "  ${CYAN}# Check status${NC}"
    echo "  ${CYAN}# Upgrade specific chaincode${NC}"
    echo "  $0 package startup"
    echo "  $0 install startup"
    echo "  $0 upgrade startup common-channel"
    echo ""
    echo "  ${CYAN}# Check status${NC}"
    echo "  $0 query-committed startup common-channel"
    echo "  $0 check-readiness validator startup-validator-channel"
    echo ""
    echo "  ${CYAN}# Switch context${NC}"
    echo "  $0 switch validator"
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
        upgrade)
            if [ -z "$arg1" ] || [ -z "$arg2" ]; then
                log_error "Usage: $0 upgrade <chaincode> <channel>"
                exit 1
            fi
            set_package_ids
            upgrade_chaincode "$arg1" "$arg2"
            ;;
        upgrade-all)
            set_package_ids
            upgrade_all
            ;;
        sync-upgrade)
            if [ -z "$arg1" ] || [ -z "$arg2" ] || [ -z "$3" ]; then
                log_error "Usage: $0 sync-upgrade <chaincode> <channel> <org>"
                exit 1
            fi
            set_package_ids
            sync_upgrade_chaincode "$arg1" "$arg2" "$3"
            ;;
        approve-chaincode)
            if [ -z "$arg1" ] || [ -z "$arg2" ] || [ -z "$3" ]; then
                log_error "Usage: $0 approve-chaincode <org> <chaincode> <channel>"
                log_info "Example: $0 approve-chaincode investor startup common-channel"
                exit 1
            fi
            set_package_ids
            approve_chaincode_for_org "$arg1" "$arg2" "$3"
            ;;
        check-readiness)
            if [ -z "$arg1" ] || [ -z "$arg2" ]; then
                log_error "Usage: $0 check-readiness <chaincode> <channel>"
                exit 1
            fi
            check_readiness "$arg2" "$arg1"
            ;;
        query-committed)
            if [ -z "$arg1" ] || [ -z "$arg2" ]; then
                log_error "Usage: $0 query-committed <chaincode> <channel>"
                exit 1
            fi
            query_committed "$arg2" "$arg1"
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
