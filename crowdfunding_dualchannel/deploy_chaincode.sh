#!/bin/bash
#
# export MICROFAB_CONFIG=$(cat MICROFAB.txt)
# docker run -d --name microfab_dual -e MICROFAB_CONFIG -p 7070:7070 ibmcom/ibp-microfab
#
# curl -s http://console.127-0-0-1.nip.io:7070/ak/api/v1/components | weft microfab -w ./_wallets -p ./_gateways -m ./_msp -f
# Installing Binaries: curl -sSL https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh | bash -s -- binary
# =============================================================================
# Chaincode Deployment Script for Crowdfunding Dual-Channel Platform
# Two Channels + 4 Org Chaincodes + Private Data Collections
# =============================================================================
#
# TOPOLOGY (from MICROFAB.txt, port 7070):
# ---------------------------------------
#   governance-validation-channel
#     Orgs:       StartupOrg, ValidatorOrg, PlatformOrg   ❌ InvestorOrg NOT a member
#     Chaincodes: startup, validator, platform            ❌ investor NOT deployed here
#     Policy:     AND('ValidatorOrgMSP.peer','PlatformOrgMSP.peer')
#     Collections: collections_config_governance.json
#
#   investment-execution-channel
#     Orgs:       StartupOrg, ValidatorOrg, InvestorOrg, PlatformOrg
#     Chaincodes: startup, validator, investor, platform
#     Policy:     OR('StartupOrgMSP.peer','ValidatorOrgMSP.peer','InvestorOrgMSP.peer','PlatformOrgMSP.peer')
#     Collections: collections_config_investment.json
#
# FEATURES:
# ---------
# 1. Organization Context Switching
#    - switch_to_startup/validator/platform/investor
#    - Exports CORE_PEER_LOCALMSPID, CORE_PEER_MSPCONFIGPATH, CORE_PEER_ADDRESS (port 7070)
#
# 2. Chaincode Packaging with Auto-Versioning
#    - Packages from contracts/<name>org/ (startup, validator, investor, platform)
#    - Auto-increments version from existing .tar.gz / .tgz files
#
# 3. Interactive Installation
#    - install <org> / install-all
#    - Auto-exports STARTUP_CC_PACKAGE_ID, VALIDATOR_CC_PACKAGE_ID,
#      INVESTOR_CC_PACKAGE_ID, PLATFORM_CC_PACKAGE_ID
#
# 4. Per-Org Deployment (same pattern as 7-channel script)
#    - deploy <org>              → approve/commit on all channels that org belongs to
#    - deploy <org> <channel>    → single channel only
#    - deploy all                → interactive deploy for each org
#    - InvestorOrg never touches governance-validation-channel
#
# 5. Upgrade / Sync-Upgrade / Approve / Check-Readiness / Query-Committed
#
# GLOBAL VARIABLES:
# -----------------
#   MSP_BASE_PATH, CONTRACTS_BASE_PATH, ORDERER_URL
#   GOVERNANCE_CHANNEL, INVESTMENT_CHANNEL
#   COLLECTIONS_GOVERNANCE, COLLECTIONS_INVESTMENT
#
# =============================================================================
# AVAILABLE COMMANDS:
# =============================================================================
#
# PACKAGING:
#   source ./deploy_chaincode.sh package                    - Package all 4 chaincodes
#   source ./deploy_chaincode.sh package <cc>               - Package one (startup|validator|investor|platform)
#
# INSTALLATION:
#   source ./deploy_chaincode.sh install <org>              - Install all 4 CCs on one org
#   source ./deploy_chaincode.sh install-all                - Install on all 4 orgs (interactive)
#
# DEPLOYMENT (Initial — org-centric, like 7-channel script):
#   source ./deploy_chaincode.sh deploy <org>               - Deploy all channels for that org
#   source ./deploy_chaincode.sh deploy <org> <channel>     - Deploy one channel for that org
#                                                      channel: governance | investment
#                                                      (aliases: gov, inv, full channel names)
#   source ./deploy_chaincode.sh deploy all                 - Interactive deploy for all orgs
#
#   Org → channel membership:
#     startup   → governance + investment
#     validator → governance + investment
#     platform  → governance + investment
#     investor  → investment ONLY  (governance is rejected / skipped)
#
# UPGRADE:
#   source ./deploy_chaincode.sh upgrade <cc> <channel>     - Approve from all member orgs + commit
#   source ./deploy_chaincode.sh upgrade-all                - Upgrade all CCs on both channels
#   source ./deploy_chaincode.sh sync-upgrade <cc> <ch> <org>
#   source ./deploy_chaincode.sh approve-chaincode <org> <cc> <ch>
#
# QUERY / CHECK:
#   source ./deploy_chaincode.sh query-committed <cc> <ch>
#   source ./deploy_chaincode.sh check-readiness <cc> <ch>
#   source ./deploy_chaincode.sh query-installed <org>
#
# UTILITY:
#   source ./deploy_chaincode.sh switch <org>
#   source ./deploy_chaincode.sh help
#
# =============================================================================
# DEPLOYMENT FLOW - INITIAL DEPLOYMENT:
# =============================================================================
#
# STEP 1: Package All Chaincodes
#   Command: source ./deploy_chaincode.sh package
#   Output:  startup_1.tar.gz, validator_1.tar.gz, investor_1.tar.gz, platform_1.tar.gz
#            from contracts/startuporg/, validatororg/, investororg/, platformorg/
#
# STEP 2: Install on All Organizations
#   Command: source ./deploy_chaincode.sh install-all
#           (OR: source ./deploy_chaincode.sh install startup  … etc.)
#   Output:  Each org peer gets all 4 packages installed
#            Auto-exports:
#              STARTUP_CC_PACKAGE_ID=startup_1:…
#              VALIDATOR_CC_PACKAGE_ID=validator_1:…
#              INVESTOR_CC_PACKAGE_ID=investor_1:…
#              PLATFORM_CC_PACKAGE_ID=platform_1:…
#
# STEP 3: Deploy Per Organization (Approve & Commit)
#   Command: source ./deploy_chaincode.sh deploy startup
#   Output:  StartupOrg on governance-validation-channel:
#              - Approves+commits own 'startup'
#              - Approves 'validator', 'platform'
#            StartupOrg on investment-execution-channel:
#              - Approves+commits own 'startup'
#              - Approves 'validator', 'investor', 'platform'
#
#   Command: source ./deploy_chaincode.sh deploy validator
#   Output:  Same pattern for validator on BOTH channels
#
#   Command: source ./deploy_chaincode.sh deploy platform
#   Output:  Same pattern for platform on BOTH channels
#
#   Command: source ./deploy_chaincode.sh deploy investor
#   Output:  InvestorOrg on investment-execution-channel ONLY:
#              - Approves+commits own 'investor'
#              - Approves 'startup', 'validator', 'platform'
#            ❌ Skips governance-validation-channel (Investor not a member)
#
#   Command: source ./deploy_chaincode.sh deploy all
#   Output:  Interactive prompts for each org in order
#
# STEP 4: Verify
#   source ./deploy_chaincode.sh query-committed startup governance
#   source ./deploy_chaincode.sh query-committed startup investment
#   source ./deploy_chaincode.sh query-committed investor investment
#   # investor on governance should NOT exist
#
# =============================================================================
# UPGRADE FLOW:
# =============================================================================
#
# STEP 1: source ./deploy_chaincode.sh package
# STEP 2: source ./deploy_chaincode.sh install-all
# STEP 3: source ./deploy_chaincode.sh upgrade-all
#    OR:  source ./deploy_chaincode.sh upgrade startup governance
#         source ./deploy_chaincode.sh upgrade startup investment
#         …
#
# =============================================================================
# EXPECTED CHANNEL / CHAINCODE MATRIX:
# =============================================================================
#
#                    | governance-validation | investment-execution
# -------------------+-----------------------+----------------------
# startup CC         | Startup+Val+Plat      | all 4 orgs
# validator CC       | Startup+Val+Plat      | all 4 orgs
# platform CC        | Startup+Val+Plat      | all 4 orgs
# investor CC        | ❌ NOT DEPLOYED       | all 4 orgs
#
# =============================================================================

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

ORDERER_URL="orderer-api.127-0-0-1.nip.io:7070"
MSP_BASE_PATH="${SCRIPT_DIR}/_msp"
CONTRACTS_BASE_PATH="${SCRIPT_DIR}/contracts"

GOVERNANCE_CHANNEL="governance-validation-channel"
INVESTMENT_CHANNEL="investment-execution-channel"

COLLECTIONS_GOVERNANCE="${SCRIPT_DIR}/collections_config_governance.json"
COLLECTIONS_INVESTMENT="${SCRIPT_DIR}/collections_config_investment.json"

PEER_STARTUP="startuporgpeer-api.127-0-0-1.nip.io:7070"
PEER_VALIDATOR="validatororgpeer-api.127-0-0-1.nip.io:7070"
PEER_INVESTOR="investororgpeer-api.127-0-0-1.nip.io:7070"
PEER_PLATFORM="platformorgpeer-api.127-0-0-1.nip.io:7070"

POLICY_GOVERNANCE="AND('ValidatorOrgMSP.peer','PlatformOrgMSP.peer')"
POLICY_INVESTMENT="OR('StartupOrgMSP.peer','ValidatorOrgMSP.peer','InvestorOrgMSP.peer','PlatformOrgMSP.peer')"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

print_header() {
    echo ""
    echo -e "${CYAN}============================================================================${NC}"
    echo -e "${CYAN}$1${NC}"
    echo -e "${CYAN}============================================================================${NC}"
}
print_success() { echo -e "${GREEN}✅ $1${NC}"; }
print_error() { echo -e "${RED}❌ $1${NC}"; }
print_info() { echo -e "${BLUE}[INFO] $1${NC}"; }
print_warning() { echo -e "${YELLOW}⚠️  $1${NC}"; }

setup_fabric_env() {
    export PATH=$PATH:${SCRIPT_DIR}/bin
    export FABRIC_CFG_PATH=${SCRIPT_DIR}/config
}

# =============================================================================
# Organization Context Switching
# =============================================================================

switch_to_startup() {
    print_info "Switching to StartupOrg..."
    export CORE_PEER_LOCALMSPID=StartupOrgMSP
    export CORE_PEER_MSPCONFIGPATH=${MSP_BASE_PATH}/StartupOrg/startuporgadmin/msp
    export CORE_PEER_ADDRESS=${PEER_STARTUP}
    setup_fabric_env
    print_success "Now operating as StartupOrg (${PEER_STARTUP})"
}

switch_to_validator() {
    print_info "Switching to ValidatorOrg..."
    export CORE_PEER_LOCALMSPID=ValidatorOrgMSP
    export CORE_PEER_MSPCONFIGPATH=${MSP_BASE_PATH}/ValidatorOrg/validatororgadmin/msp
    export CORE_PEER_ADDRESS=${PEER_VALIDATOR}
    setup_fabric_env
    print_success "Now operating as ValidatorOrg (${PEER_VALIDATOR})"
}

switch_to_investor() {
    print_info "Switching to InvestorOrg..."
    export CORE_PEER_LOCALMSPID=InvestorOrgMSP
    export CORE_PEER_MSPCONFIGPATH=${MSP_BASE_PATH}/InvestorOrg/investororgadmin/msp
    export CORE_PEER_ADDRESS=${PEER_INVESTOR}
    setup_fabric_env
    print_success "Now operating as InvestorOrg (${PEER_INVESTOR})"
}

switch_to_platform() {
    print_info "Switching to PlatformOrg..."
    export CORE_PEER_LOCALMSPID=PlatformOrgMSP
    export CORE_PEER_MSPCONFIGPATH=${MSP_BASE_PATH}/PlatformOrg/platformorgadmin/msp
    export CORE_PEER_ADDRESS=${PEER_PLATFORM}
    setup_fabric_env
    print_success "Now operating as PlatformOrg (${PEER_PLATFORM})"
}

switch_org() {
    case "$1" in
        startup) switch_to_startup ;;
        validator) switch_to_validator ;;
        investor) switch_to_investor ;;
        platform) switch_to_platform ;;
        *) print_error "Unknown organization: $1 (startup|validator|investor|platform)"; return 1 ;;
    esac
}

# =============================================================================
# Channel / policy helpers
# =============================================================================

resolve_channel() {
    case "$1" in
        governance|gov|governance-validation|governance-validation-channel)
            echo "${GOVERNANCE_CHANNEL}" ;;
        investment|inv|investment-execution|investment-execution-channel)
            echo "${INVESTMENT_CHANNEL}" ;;
        *) echo "" ;;
    esac
}

# Normalize short aliases used in deploy <org> <channel>
normalize_channel_filter() {
    case "$1" in
        ""|all) echo "" ;;
        governance|gov|governance-validation|governance-validation-channel) echo "governance" ;;
        investment|inv|investment-execution|investment-execution-channel) echo "investment" ;;
        *) echo "INVALID" ;;
    esac
}

get_policy() {
    if [ "$1" = "${GOVERNANCE_CHANNEL}" ]; then
        echo "${POLICY_GOVERNANCE}"
    else
        echo "${POLICY_INVESTMENT}"
    fi
}

get_collections() {
    if [ "$1" = "${GOVERNANCE_CHANNEL}" ]; then
        echo "${COLLECTIONS_GOVERNANCE}"
    else
        echo "${COLLECTIONS_INVESTMENT}"
    fi
}

# Orgs that must approve a chaincode on a given channel
orgs_on_channel() {
    case "$1" in
        "${GOVERNANCE_CHANNEL}") echo "startup validator platform" ;;
        "${INVESTMENT_CHANNEL}") echo "startup validator investor platform" ;;
        *) echo "" ;;
    esac
}

# Chaincodes deployed on a given channel
ccs_on_channel() {
    case "$1" in
        "${GOVERNANCE_CHANNEL}") echo "startup validator platform" ;;
        "${INVESTMENT_CHANNEL}") echo "startup validator investor platform" ;;
        *) echo "" ;;
    esac
}

get_package_id_var() { echo "${1^^}_CC_PACKAGE_ID"; }

get_package_id() {
    local var
    var=$(get_package_id_var "$1")
    echo "${!var}"
}

require_package_id() {
    local cc=$1
    local pid
    pid=$(get_package_id "$cc")
    if [ -z "$pid" ]; then
        print_error "Package ID not set: $(get_package_id_var "$cc")"
        print_warning "Run: source ./deploy_chaincode.sh install-all"
        return 1
    fi
}

# =============================================================================
# Packaging
# =============================================================================

get_next_package_version() {
    local chaincode_name=$1
    local max_version=0
    shopt -s nullglob
    for file in ${SCRIPT_DIR}/${chaincode_name}_*.tar.gz ${SCRIPT_DIR}/${chaincode_name}_*.tgz; do
        [ -f "$file" ] || continue
        local version
        version=$(basename "$file" | sed -n "s/${chaincode_name}_\([0-9]\+\)\.\(tar\.gz\|tgz\)/\1/p")
        if [ -n "$version" ] && [ "$version" -gt "$max_version" ]; then
            max_version=$version
        fi
    done
    shopt -u nullglob
    echo $((max_version + 1))
}

package_one() {
    local chaincode_name=$1
    local contract_path="${CONTRACTS_BASE_PATH}/${chaincode_name}org"
    local version
    version=$(get_next_package_version "$chaincode_name")
    local package_file="${SCRIPT_DIR}/${chaincode_name}_${version}.tar.gz"
    local label="${chaincode_name}_${version}"

    print_header "📦 Packaging chaincode: ${chaincode_name}"
    setup_fabric_env
    cd "${SCRIPT_DIR}"

    if [ ! -d "$contract_path" ]; then
        print_error "Contract path does not exist: ${contract_path}"
        return 1
    fi

    print_info "Current → new version auto-detected"
    print_info "Package label: ${label}"
    print_info "Package file:  ${package_file}"
    print_info "Path:          ${contract_path}"

    peer lifecycle chaincode package "${package_file}" \
        --path "${contract_path}" \
        --lang golang \
        --label "${label}"

    if [ $? -eq 0 ]; then
        print_success "Successfully packaged ${chaincode_name} as ${package_file}"
    else
        print_error "Failed to package ${chaincode_name}"
        return 1
    fi
}

package_all() {
    local target=${1:-all}
    setup_fabric_env
    case "$target" in
        all|"")
            print_header "📦 Packaging All Chaincodes"
            package_one startup || return 1
            package_one validator || return 1
            package_one investor || return 1
            package_one platform || return 1
            print_success "All chaincodes packaged successfully!"
            print_warning "Next: source ./deploy_chaincode.sh install-all"
            ;;
        startup|validator|investor|platform)
            package_one "$target" || return 1
            ;;
        *)
            print_error "Usage: package [startup|validator|investor|platform|all]"
            return 1
            ;;
    esac
}

latest_package() {
    local chaincode_name=$1
    ls -t ${SCRIPT_DIR}/${chaincode_name}_*.tar.gz ${SCRIPT_DIR}/${chaincode_name}_*.tgz 2>/dev/null | head -1
}

# =============================================================================
# Installation
# =============================================================================

install_one_cc() {
    local chaincode_name=$1
    local pkg
    pkg=$(latest_package "$chaincode_name")
    if [ -z "$pkg" ]; then
        print_error "No package for ${chaincode_name}. Run 'source ./deploy_chaincode.sh package' first."
        return 1
    fi

    print_info "Installing ${pkg} on current peer..."
    local install_output install_status package_id
    install_output=$(peer lifecycle chaincode install "${pkg}" 2>&1)
    install_status=$?
    echo "$install_output"

    if echo "$install_output" | grep -q "already successfully installed"; then
        print_warning "${chaincode_name} already installed"
        package_id=$(echo "$install_output" | grep -oP "package ID '\K[^']+")
        if [ -z "$package_id" ]; then
            package_id=$(peer lifecycle chaincode queryinstalled 2>&1 | grep "Label: ${chaincode_name}_" | head -1 | grep -oP "Package ID: \K[^,]+")
        fi
    elif [ $install_status -eq 0 ]; then
        print_success "${chaincode_name} installed successfully"
        package_id=$(echo "$install_output" | grep -oP 'Chaincode code package identifier: \K.*')
    else
        print_error "Install failed for ${chaincode_name}"
        return 1
    fi

    if [ -z "$package_id" ]; then
        print_error "Could not extract package ID for ${chaincode_name}"
        return 1
    fi

    local var_name
    var_name=$(get_package_id_var "$chaincode_name")
    export "${var_name}=${package_id}"
    echo "export ${var_name}=\"${package_id}\"" >> ~/.bashrc
    print_success "Auto-exported: ${var_name}=${package_id}"
}

install_on_org() {
    local org=$1
    case "$org" in
        startup|validator|investor|platform) ;;
        *) print_error "Invalid org: $org (startup|validator|investor|platform)"; return 1 ;;
    esac

    print_header "📥 Installing Chaincodes on ${org}"
    switch_org "$org" || return 1

    for cc in startup validator investor platform; do
        print_header "Installing ${cc} on ${org}"
        print_warning "Press Enter to install ${cc}..."
        read -r
        install_one_cc "$cc" || return 1
        echo ""
    done

    print_success "Installation complete on ${org}"
}

install_all() {
    print_header "📥 Installing Chaincodes on All Organizations"
    for org in startup validator investor platform; do
        install_on_org "$org" || return 1
    done
    echo ""
    print_success "Installation complete on all 4 organizations!"
    print_info "STARTUP_CC_PACKAGE_ID=${STARTUP_CC_PACKAGE_ID}"
    print_info "VALIDATOR_CC_PACKAGE_ID=${VALIDATOR_CC_PACKAGE_ID}"
    print_info "INVESTOR_CC_PACKAGE_ID=${INVESTOR_CC_PACKAGE_ID}"
    print_info "PLATFORM_CC_PACKAGE_ID=${PLATFORM_CC_PACKAGE_ID}"
    echo ""
    print_warning "Next: source ./deploy_chaincode.sh deploy startup"
    print_warning "  then deploy validator, platform, investor (or: deploy all)"
}

# =============================================================================
# Lifecycle: approve / commit / readiness
# =============================================================================

get_committed_sequence() {
    local channel=$1
    local chaincode_name=$2
    local result
    result=$(peer lifecycle chaincode querycommitted --channelID "${channel}" --name "${chaincode_name}" --output json 2>/dev/null || echo "")
    if [ -z "$result" ]; then echo "0"; return; fi
    local seq
    seq=$(echo "$result" | jq -r '.sequence // 0' 2>/dev/null || echo "0")
    echo "${seq:-0}"
}

get_committed_version() {
    local channel=$1
    local chaincode_name=$2
    local result
    result=$(peer lifecycle chaincode querycommitted --channelID "${channel}" --name "${chaincode_name}" --output json 2>/dev/null || echo "")
    if [ -z "$result" ]; then echo "0"; return; fi
    local ver
    ver=$(echo "$result" | jq -r '.version // "0"' 2>/dev/null || echo "0")
    [ -z "$ver" ] || [ "$ver" = "null" ] && echo "0" || echo "$ver"
}

next_sequence() {
    local cur
    cur=$(get_committed_sequence "$1" "$2")
    echo $((cur + 1))
}

next_version() {
    local cur
    cur=$(get_committed_version "$1" "$2")
    if [ "$cur" = "0" ] || [ -z "$cur" ]; then echo "1"
    elif [[ "$cur" =~ ^[0-9]+$ ]]; then echo $((cur + 1))
    else echo "$cur"
    fi
}

approve_chaincode() {
    local channel=$1
    local chaincode_name=$2
    local package_id=$3
    local sequence version policy collections

    sequence=$(next_sequence "$channel" "$chaincode_name")
    version=$(next_version "$channel" "$chaincode_name")
    policy=$(get_policy "$channel")
    collections=$(get_collections "$channel")

    print_info "Approving '${chaincode_name}' on '${channel}' (version: ${version}, sequence: ${sequence})..."

    peer lifecycle chaincode approveformyorg \
        -o "${ORDERER_URL}" \
        --channelID "${channel}" \
        --name "${chaincode_name}" \
        --version "${version}" \
        --sequence "${sequence}" \
        --package-id "${package_id}" \
        --collections-config "${collections}" \
        --signature-policy "${policy}" \
        --waitForEvent

    if [ $? -eq 0 ]; then
        print_success "Approved '${chaincode_name}' on '${channel}' (version: ${version}, sequence: ${sequence})"
    else
        print_error "Failed to approve '${chaincode_name}' on '${channel}'"
        return 1
    fi
}

commit_chaincode() {
    local channel=$1
    local chaincode_name=$2
    local sequence version policy collections

    sequence=$(next_sequence "$channel" "$chaincode_name")
    version=$(next_version "$channel" "$chaincode_name")
    policy=$(get_policy "$channel")
    collections=$(get_collections "$channel")

    print_info "Committing '${chaincode_name}' on '${channel}' (version: ${version}, sequence: ${sequence})..."

    if [ "$channel" = "${GOVERNANCE_CHANNEL}" ]; then
        # 3 orgs only — Investor peer must NOT be included
        peer lifecycle chaincode commit \
            -o "${ORDERER_URL}" \
            --channelID "${channel}" \
            --name "${chaincode_name}" \
            --version "${version}" \
            --sequence "${sequence}" \
            --collections-config "${collections}" \
            --signature-policy "${policy}" \
            --peerAddresses "${PEER_STARTUP}" \
            --peerAddresses "${PEER_VALIDATOR}" \
            --peerAddresses "${PEER_PLATFORM}" \
            --waitForEvent
    else
        # All 4 orgs
        peer lifecycle chaincode commit \
            -o "${ORDERER_URL}" \
            --channelID "${channel}" \
            --name "${chaincode_name}" \
            --version "${version}" \
            --sequence "${sequence}" \
            --collections-config "${collections}" \
            --signature-policy "${policy}" \
            --peerAddresses "${PEER_STARTUP}" \
            --peerAddresses "${PEER_VALIDATOR}" \
            --peerAddresses "${PEER_INVESTOR}" \
            --peerAddresses "${PEER_PLATFORM}" \
            --waitForEvent
    fi

    if [ $? -eq 0 ]; then
        print_success "Committed '${chaincode_name}' on '${channel}' (version: ${version}, sequence: ${sequence})"
    else
        print_error "Failed to commit '${chaincode_name}' on '${channel}'"
        return 1
    fi
}

check_readiness() {
    local chaincode_name=$1
    local channel_alias=$2
    local channel
    channel=$(resolve_channel "$channel_alias")
    if [ -z "$channel" ]; then
        print_error "Usage: check-readiness <cc> <governance|investment>"
        return 1
    fi

    # Investor chaincode / Investor org never on governance
    if [ "$channel" = "${GOVERNANCE_CHANNEL}" ] && [ "$chaincode_name" = "investor" ]; then
        print_error "investor chaincode is NOT deployed on governance-validation-channel"
        return 1
    fi

    switch_to_startup
    local sequence version policy collections
    sequence=$(next_sequence "$channel" "$chaincode_name")
    version=$(next_version "$channel" "$chaincode_name")
    policy=$(get_policy "$channel")
    collections=$(get_collections "$channel")

    print_info "Checking commit readiness for '${chaincode_name}' on '${channel}' (v${version}/seq${sequence})..."
    peer lifecycle chaincode checkcommitreadiness \
        --channelID "${channel}" \
        --name "${chaincode_name}" \
        --version "${version}" \
        --sequence "${sequence}" \
        --collections-config "${collections}" \
        --signature-policy "${policy}"
}

# =============================================================================
# Per-Org Deployment (mirrors 7-channel deploy_<org>_org pattern)
# =============================================================================

deploy_startup_org() {
    local channel_filter
    channel_filter=$(normalize_channel_filter "$1")
    if [ "$channel_filter" = "INVALID" ]; then
        print_error "Invalid channel filter: $1 (governance|investment)"
        return 1
    fi

    print_header "🚀 Deploying chaincodes for StartupOrg"
    switch_to_startup
    require_package_id startup || return 1
    require_package_id validator || return 1
    require_package_id platform || return 1

    # ---- governance-validation-channel (3 orgs) ----
    if [ -z "$channel_filter" ] || [ "$channel_filter" = "governance" ]; then
        print_header "StartupOrg: ${GOVERNANCE_CHANNEL}"
        print_info "Members: StartupOrg, ValidatorOrg, PlatformOrg (NO InvestorOrg)"

        # Own chaincode → approve + commit
        approve_chaincode "${GOVERNANCE_CHANNEL}" "startup" "$(get_package_id startup)" || return 1
        commit_chaincode "${GOVERNANCE_CHANNEL}" "startup" || return 1

        # Other orgs' chaincodes → approve only
        approve_chaincode "${GOVERNANCE_CHANNEL}" "validator" "$(get_package_id validator)" || return 1
        approve_chaincode "${GOVERNANCE_CHANNEL}" "platform" "$(get_package_id platform)" || return 1
    fi

    # ---- investment-execution-channel (4 orgs) ----
    if [ -z "$channel_filter" ] || [ "$channel_filter" = "investment" ]; then
        print_header "StartupOrg: ${INVESTMENT_CHANNEL}"
        print_info "Members: StartupOrg, ValidatorOrg, InvestorOrg, PlatformOrg"
        require_package_id investor || return 1

        approve_chaincode "${INVESTMENT_CHANNEL}" "startup" "$(get_package_id startup)" || return 1
        commit_chaincode "${INVESTMENT_CHANNEL}" "startup" || return 1

        approve_chaincode "${INVESTMENT_CHANNEL}" "validator" "$(get_package_id validator)" || return 1
        approve_chaincode "${INVESTMENT_CHANNEL}" "investor" "$(get_package_id investor)" || return 1
        approve_chaincode "${INVESTMENT_CHANNEL}" "platform" "$(get_package_id platform)" || return 1
    fi

    print_success "StartupOrg deployment complete!"
}

deploy_validator_org() {
    local channel_filter
    channel_filter=$(normalize_channel_filter "$1")
    if [ "$channel_filter" = "INVALID" ]; then
        print_error "Invalid channel filter: $1 (governance|investment)"
        return 1
    fi

    print_header "🚀 Deploying chaincodes for ValidatorOrg"
    switch_to_validator
    require_package_id validator || return 1
    require_package_id startup || return 1
    require_package_id platform || return 1

    if [ -z "$channel_filter" ] || [ "$channel_filter" = "governance" ]; then
        print_header "ValidatorOrg: ${GOVERNANCE_CHANNEL}"
        print_info "Members: StartupOrg, ValidatorOrg, PlatformOrg (NO InvestorOrg)"

        approve_chaincode "${GOVERNANCE_CHANNEL}" "validator" "$(get_package_id validator)" || return 1
        commit_chaincode "${GOVERNANCE_CHANNEL}" "validator" || return 1

        approve_chaincode "${GOVERNANCE_CHANNEL}" "startup" "$(get_package_id startup)" || return 1
        approve_chaincode "${GOVERNANCE_CHANNEL}" "platform" "$(get_package_id platform)" || return 1
    fi

    if [ -z "$channel_filter" ] || [ "$channel_filter" = "investment" ]; then
        print_header "ValidatorOrg: ${INVESTMENT_CHANNEL}"
        require_package_id investor || return 1

        approve_chaincode "${INVESTMENT_CHANNEL}" "validator" "$(get_package_id validator)" || return 1
        commit_chaincode "${INVESTMENT_CHANNEL}" "validator" || return 1

        approve_chaincode "${INVESTMENT_CHANNEL}" "startup" "$(get_package_id startup)" || return 1
        approve_chaincode "${INVESTMENT_CHANNEL}" "investor" "$(get_package_id investor)" || return 1
        approve_chaincode "${INVESTMENT_CHANNEL}" "platform" "$(get_package_id platform)" || return 1
    fi

    print_success "ValidatorOrg deployment complete!"
}

deploy_platform_org() {
    local channel_filter
    channel_filter=$(normalize_channel_filter "$1")
    if [ "$channel_filter" = "INVALID" ]; then
        print_error "Invalid channel filter: $1 (governance|investment)"
        return 1
    fi

    print_header "🚀 Deploying chaincodes for PlatformOrg"
    switch_to_platform
    require_package_id platform || return 1
    require_package_id startup || return 1
    require_package_id validator || return 1

    if [ -z "$channel_filter" ] || [ "$channel_filter" = "governance" ]; then
        print_header "PlatformOrg: ${GOVERNANCE_CHANNEL}"
        print_info "Members: StartupOrg, ValidatorOrg, PlatformOrg (NO InvestorOrg)"

        approve_chaincode "${GOVERNANCE_CHANNEL}" "platform" "$(get_package_id platform)" || return 1
        commit_chaincode "${GOVERNANCE_CHANNEL}" "platform" || return 1

        approve_chaincode "${GOVERNANCE_CHANNEL}" "startup" "$(get_package_id startup)" || return 1
        approve_chaincode "${GOVERNANCE_CHANNEL}" "validator" "$(get_package_id validator)" || return 1
    fi

    if [ -z "$channel_filter" ] || [ "$channel_filter" = "investment" ]; then
        print_header "PlatformOrg: ${INVESTMENT_CHANNEL}"
        require_package_id investor || return 1

        approve_chaincode "${INVESTMENT_CHANNEL}" "platform" "$(get_package_id platform)" || return 1
        commit_chaincode "${INVESTMENT_CHANNEL}" "platform" || return 1

        approve_chaincode "${INVESTMENT_CHANNEL}" "startup" "$(get_package_id startup)" || return 1
        approve_chaincode "${INVESTMENT_CHANNEL}" "validator" "$(get_package_id validator)" || return 1
        approve_chaincode "${INVESTMENT_CHANNEL}" "investor" "$(get_package_id investor)" || return 1
    fi

    print_success "PlatformOrg deployment complete!"
}

deploy_investor_org() {
    local channel_filter
    channel_filter=$(normalize_channel_filter "$1")
    if [ "$channel_filter" = "INVALID" ]; then
        print_error "Invalid channel filter: $1 (governance|investment)"
        return 1
    fi

    print_header "🚀 Deploying chaincodes for InvestorOrg"
    switch_to_investor
    require_package_id investor || return 1
    require_package_id startup || return 1
    require_package_id validator || return 1
    require_package_id platform || return 1

    # Investor is NOT on governance — refuse / skip
    if [ "$channel_filter" = "governance" ]; then
        print_error "InvestorOrg is NOT a member of ${GOVERNANCE_CHANNEL}"
        print_warning "Investor only deploys on ${INVESTMENT_CHANNEL}"
        return 1
    fi

    if [ -z "$channel_filter" ]; then
        print_info "Skipping ${GOVERNANCE_CHANNEL} — InvestorOrg is not a member"
    fi

    # investment-execution-channel only
    if [ -z "$channel_filter" ] || [ "$channel_filter" = "investment" ]; then
        print_header "InvestorOrg: ${INVESTMENT_CHANNEL}"
        print_info "Members: StartupOrg, ValidatorOrg, InvestorOrg, PlatformOrg"

        approve_chaincode "${INVESTMENT_CHANNEL}" "investor" "$(get_package_id investor)" || return 1
        commit_chaincode "${INVESTMENT_CHANNEL}" "investor" || return 1

        approve_chaincode "${INVESTMENT_CHANNEL}" "startup" "$(get_package_id startup)" || return 1
        approve_chaincode "${INVESTMENT_CHANNEL}" "validator" "$(get_package_id validator)" || return 1
        approve_chaincode "${INVESTMENT_CHANNEL}" "platform" "$(get_package_id platform)" || return 1
    fi

    print_success "InvestorOrg deployment complete!"
}

deploy_all_orgs() {
    print_header "🚀 Deploying chaincodes for ALL organizations (dual-channel)"
    print_warning "Order: startup → validator → platform → investor"
    print_warning "Investor only acts on investment-execution-channel"

    echo ""
    read -p "Deploy for StartupOrg? (y/n): " confirm
    if [ "$confirm" = "y" ] || [ "$confirm" = "Y" ]; then
        deploy_startup_org || return 1
    fi

    echo ""
    read -p "Deploy for ValidatorOrg? (y/n): " confirm
    if [ "$confirm" = "y" ] || [ "$confirm" = "Y" ]; then
        deploy_validator_org || return 1
    fi

    echo ""
    read -p "Deploy for PlatformOrg? (y/n): " confirm
    if [ "$confirm" = "y" ] || [ "$confirm" = "Y" ]; then
        deploy_platform_org || return 1
    fi

    echo ""
    read -p "Deploy for InvestorOrg (investment channel only)? (y/n): " confirm
    if [ "$confirm" = "y" ] || [ "$confirm" = "Y" ]; then
        deploy_investor_org || return 1
    fi

    print_success "All dual-channel deployments complete!"
}

# deploy <org|all> [channel]
deploy_cmd() {
    local org=$1
    local channel_filter=$2

    case "$org" in
        startup) deploy_startup_org "$channel_filter" ;;
        validator) deploy_validator_org "$channel_filter" ;;
        platform) deploy_platform_org "$channel_filter" ;;
        investor) deploy_investor_org "$channel_filter" ;;
        all) deploy_all_orgs ;;
        *)
            print_error "Usage: deploy <startup|validator|investor|platform|all> [governance|investment]"
            return 1
            ;;
    esac
}

# =============================================================================
# Upgrade helpers
# =============================================================================

smart_approve() {
    local channel=$1
    local chaincode_name=$2
    local package_id=$3
    local org_name=$4

    print_info "Smart approve for '${chaincode_name}' on '${channel}' by ${org_name}..."

    local current_sequence
    current_sequence=$(get_committed_sequence "$channel" "$chaincode_name")
    if [ "$current_sequence" = "0" ]; then
        print_warning "${chaincode_name} not yet committed on ${channel}. Skipping approval for ${org_name}."
        return 0
    fi

    print_info "${chaincode_name} is committed (sequence: ${current_sequence}). Proceeding..."
    approve_chaincode "$channel" "$chaincode_name" "$package_id"
}

sync_upgrade_chaincode() {
    local chaincode_name=$1
    local channel_alias=$2
    local org=$3
    local channel package_id

    channel=$(resolve_channel "$channel_alias")
    if [ -z "$channel" ] || [ -z "$org" ] || [ -z "$chaincode_name" ]; then
        print_error "Usage: sync-upgrade <cc> <governance|investment> <org>"
        return 1
    fi

    if [ "$channel" = "${GOVERNANCE_CHANNEL}" ]; then
        if [ "$chaincode_name" = "investor" ] || [ "$org" = "investor" ]; then
            print_error "InvestorOrg / investor chaincode cannot sync-upgrade on governance-validation-channel"
            return 1
        fi
    fi

    print_header "Sync-Upgrade: ${chaincode_name} on ${channel} for ${org}"
    switch_org "$org" || return 1
    require_package_id "$chaincode_name" || return 1
    package_id=$(get_package_id "$chaincode_name")
    smart_approve "$channel" "$chaincode_name" "$package_id" "$org"
}

approve_chaincode_cmd() {
    local org=$1
    local chaincode_name=$2
    local channel_alias=$3
    local channel package_id

    channel=$(resolve_channel "$channel_alias")
    if [ -z "$org" ] || [ -z "$chaincode_name" ] || [ -z "$channel" ]; then
        print_error "Usage: approve-chaincode <org> <cc> <governance|investment>"
        return 1
    fi

    if [ "$channel" = "${GOVERNANCE_CHANNEL}" ]; then
        if [ "$org" = "investor" ] || [ "$chaincode_name" = "investor" ]; then
            print_error "InvestorOrg / investor CC cannot approve on governance-validation-channel"
            return 1
        fi
    fi

    print_header "Approving ${chaincode_name} on ${channel} for ${org}"
    switch_org "$org" || return 1
    require_package_id "$chaincode_name" || return 1
    package_id=$(get_package_id "$chaincode_name")
    smart_approve "$channel" "$chaincode_name" "$package_id" "$org"
}

upgrade_chaincode() {
    local chaincode_name=$1
    local channel_alias=$2
    local channel package_id org

    channel=$(resolve_channel "$channel_alias")
    if [ -z "$chaincode_name" ] || [ -z "$channel" ]; then
        print_error "Usage: upgrade <startup|validator|investor|platform> <governance|investment>"
        return 1
    fi

    if [ "$channel" = "${GOVERNANCE_CHANNEL}" ] && [ "$chaincode_name" = "investor" ]; then
        print_error "investor chaincode is NOT deployed on governance-validation-channel"
        return 1
    fi

    print_header "⬆️  Upgrading ${chaincode_name} on ${channel}"
    require_package_id "$chaincode_name" || return 1
    package_id=$(get_package_id "$chaincode_name")

    for org in $(orgs_on_channel "$channel"); do
        switch_org "$org" || return 1
        approve_chaincode "$channel" "$chaincode_name" "$package_id" || return 1
    done

    case "$chaincode_name" in
        startup) switch_to_startup ;;
        validator) switch_to_validator ;;
        investor) switch_to_investor ;;
        platform) switch_to_platform ;;
    esac
    commit_chaincode "$channel" "$chaincode_name" || return 1
    print_success "Upgrade complete for ${chaincode_name} on ${channel}"
}

upgrade_all() {
    print_header "⬆️  Comprehensive Upgrade — Both Dual Channels"
    print_warning "Requires exported package IDs from install-all"
    echo ""
    read -p "Continue with upgrade-all? (y/n): " confirm
    if [ "$confirm" != "y" ] && [ "$confirm" != "Y" ]; then
        print_info "Upgrade aborted"
        return 0
    fi

    print_header "Upgrading ${GOVERNANCE_CHANNEL} (startup, validator, platform)"
    for cc in startup validator platform; do
        upgrade_chaincode "$cc" governance || return 1
    done

    print_header "Upgrading ${INVESTMENT_CHANNEL} (startup, validator, investor, platform)"
    for cc in startup validator investor platform; do
        upgrade_chaincode "$cc" investment || return 1
    done

    print_success "All chaincodes upgraded on both dual channels!"
}

# =============================================================================
# Query helpers
# =============================================================================

query_committed() {
    local chaincode_name=$1
    local channel_alias=$2
    local channel

    if [ -z "$chaincode_name" ] || [ -z "$channel_alias" ]; then
        print_error "Usage: query-committed <cc> <governance|investment>"
        return 1
    fi

    channel=$(resolve_channel "$channel_alias")
    if [ -z "$channel" ]; then
        print_error "Invalid channel: $channel_alias"
        return 1
    fi

    if [ "$channel" = "${GOVERNANCE_CHANNEL}" ] && [ "$chaincode_name" = "investor" ]; then
        print_warning "investor chaincode is not expected on governance-validation-channel"
    fi

    switch_to_startup
    print_info "Querying committed '${chaincode_name}' on '${channel}'..."
    peer lifecycle chaincode querycommitted --channelID "${channel}" --name "${chaincode_name}"
}

query_installed() {
    local org=$1
    if [ -z "$org" ]; then
        print_error "Usage: query-installed <startup|validator|investor|platform>"
        return 1
    fi
    switch_org "$org" || return 1
    print_info "Querying installed chaincodes on ${org}..."
    peer lifecycle chaincode queryinstalled
}

# =============================================================================
# Help
# =============================================================================

show_help() {
    echo ""
    echo "============================================================================="
    echo "Crowdfunding Dual-Channel — Chaincode Deployment Tool"
    echo "Port 7070 | 2 Channels | 4 Org Chaincodes + PDC"
    echo "============================================================================="
    echo ""
    echo "CHANNELS:"
    echo "  governance  → ${GOVERNANCE_CHANNEL}"
    echo "                Orgs: Startup, Validator, Platform   (NO Investor)"
    echo "                CCs:  startup, validator, platform   (NO investor CC)"
    echo "                Policy: ${POLICY_GOVERNANCE}"
    echo ""
    echo "  investment  → ${INVESTMENT_CHANNEL}"
    echo "                Orgs: Startup, Validator, Investor, Platform"
    echo "                CCs:  startup, validator, investor, platform"
    echo "                Policy: ${POLICY_INVESTMENT}"
    echo ""
    echo "PACKAGING:"
    echo "  package [cc|all]                 Package from contracts/<cc>org"
    echo ""
    echo "INSTALLATION:"
    echo "  install <org>                    Install all 4 CCs on one org"
    echo "  install-all                      Install on all 4 orgs"
    echo ""
    echo "DEPLOYMENT (org-centric, like 7-channel script):"
    echo "  deploy <org>                     Deploy all channels for org"
    echo "  deploy <org> <governance|investment>"
    echo "  deploy all                       Interactive deploy for all orgs"
    echo ""
    echo "  Per-org channel cases:"
    echo "    startup   → governance + investment"
    echo "    validator → governance + investment"
    echo "    platform  → governance + investment"
    echo "    investor  → investment ONLY (governance rejected)"
    echo ""
    echo "UPGRADE:"
    echo "  upgrade <cc> <governance|investment>"
    echo "  upgrade-all"
    echo "  sync-upgrade <cc> <ch> <org>"
    echo "  approve-chaincode <org> <cc> <ch>"
    echo ""
    echo "QUERY:"
    echo "  query-committed <cc> <ch>"
    echo "  check-readiness <cc> <ch>"
    echo "  query-installed <org>"
    echo ""
    echo "UTILITY:"
    echo "  switch <startup|validator|investor|platform>"
    echo "  help"
    echo ""
    echo "INITIAL FLOW:"
    echo "  source ./deploy_chaincode.sh package"
    echo "  source ./deploy_chaincode.sh install-all"
    echo "  source ./deploy_chaincode.sh deploy startup"
    echo "  source ./deploy_chaincode.sh deploy validator"
    echo "  source ./deploy_chaincode.sh deploy platform"
    echo "  source ./deploy_chaincode.sh deploy investor"
    echo "  # OR: source ./deploy_chaincode.sh deploy all"
    echo ""
    echo "VERIFY:"
    echo "  source ./deploy_chaincode.sh query-committed startup governance"
    echo "  source ./deploy_chaincode.sh query-committed startup investment"
    echo "  source ./deploy_chaincode.sh query-committed investor investment"
    echo ""
    echo "PATHS:"
    echo "  Contracts: ${CONTRACTS_BASE_PATH}"
    echo "  MSP:       ${MSP_BASE_PATH}"
    echo "  Orderer:   ${ORDERER_URL}"
    echo ""
}

# =============================================================================
# Main
# =============================================================================

case "$1" in
    package)
        package_all "${2:-all}"
        ;;
    install)
        if [ -z "$2" ]; then
            print_error "Please specify organization: startup|validator|investor|platform"
            return 1
        fi
        install_on_org "$2"
        ;;
    install-all)
        install_all
        ;;
    deploy)
        if [ -z "$2" ]; then
            print_error "Usage: deploy <startup|validator|investor|platform|all> [governance|investment]"
            return 1
        fi
        deploy_cmd "$2" "$3"
        ;;
    upgrade)
        upgrade_chaincode "$2" "$3"
        ;;
    upgrade-all)
        upgrade_all
        ;;
    sync-upgrade)
        sync_upgrade_chaincode "$2" "$3" "$4"
        ;;
    approve-chaincode)
        approve_chaincode_cmd "$2" "$3" "$4"
        ;;
    query-committed)
        query_committed "$2" "$3"
        ;;
    check-readiness)
        check_readiness "$2" "$3"
        ;;
    query-installed)
        query_installed "$2"
        ;;
    switch)
        if [ -z "$2" ]; then
            print_error "Please specify organization: startup|validator|investor|platform"
            return 1
        fi
        switch_org "$2"
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        print_error "Invalid command: $1"
        echo ""
        print_info "Run 'source ./deploy_chaincode.sh help' for usage"
        return 1
        ;;
esac
