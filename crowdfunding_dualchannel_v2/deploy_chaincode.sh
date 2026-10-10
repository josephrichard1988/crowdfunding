#!/bin/bash
#
# export MICROFAB_CONFIG=$(cat MICROFAB.txt)
# docker run -d --name microfab_dual -e MICROFAB_CONFIG -p 7070:7070 ibmcom/ibp-microfab
#
# curl -s http://console.127-0-0-1.nip.io:7070/ak/api/v1/components | weft microfab -w ./_wallets -p ./_gateways -m ./_msp -f
# Installing Binaries: curl -sSL https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh | bash -s -- binary
# =============================================================================
# Dual-Channel Combined Chaincodes
# Two Channels + 2 Combined Packages + Private Data Collections
# =============================================================================
#
# TOPOLOGY (from MICROFAB.txt, port 7070):
# ---------------------------------------
#   governance-validation-channel
#     Orgs:       StartupOrg, ValidatorOrg, PlatformOrg   ❌ InvestorOrg NOT a member
#     Chaincode:  governance  (Startup + Validator + Platform + DataBridge)
#     Policy:     AND('ValidatorOrgMSP.peer','PlatformOrgMSP.peer')
#     Collections: collections_config_governance.json
#
#   investment-execution-channel
#     Orgs:       StartupOrg, ValidatorOrg, InvestorOrg, PlatformOrg
#     Chaincode:  investment  (all 4 orgs + Token + DataBridge)
#     Policy:     OR('StartupOrgMSP.peer','ValidatorOrgMSP.peer','InvestorOrgMSP.peer','PlatformOrgMSP.peer')
#     Collections: collections_config_investment.json
#
# INVOKE STYLE (multi-contract packages — same as crowdfundingv2):
#   -n governance  -c '{"function":"StartupContract:CreateCampaign","Args":[...]}'
#   -n investment  -c '{"function":"InvestorContract:MakeInvestment","Args":[...]}'
#
# 4-org-CC variant remains at ../crowdfunding_dualchannel/
#
# =============================================================================
# AVAILABLE COMMANDS:
# =============================================================================
#
#   source ./deploy_chaincode.sh package [governance|investment|all]
#   source ./deploy_chaincode.sh install <org>
#   source ./deploy_chaincode.sh install-all
#   source ./deploy_chaincode.sh deploy <governance|investment|all>
#   source ./deploy_chaincode.sh upgrade <governance|investment>
#   source ./deploy_chaincode.sh upgrade-all
#   source ./deploy_chaincode.sh query-committed <governance|investment>
#   source ./deploy_chaincode.sh check-readiness <governance|investment>
#   source ./deploy_chaincode.sh query-installed <org>
#   source ./deploy_chaincode.sh switch <org>
#   source ./deploy_chaincode.sh help
#
# INITIAL FLOW:
#   source ./deploy_chaincode.sh package
#   source ./deploy_chaincode.sh install-all
#   source ./deploy_chaincode.sh deploy all
#
# =============================================================================

# Change PROJECT_DIR according to your device
PROJECT_DIR="/root/crowdfunding/crowdfunding_dualchannel_v2"

ORDERER_URL="orderer-api.127-0-0-1.nip.io:7070"
MSP_BASE_PATH="/root/crowdfunding/crowdfunding_dualchannel_v2/_msp"
CONTRACTS_BASE_PATH="./contracts"

GOVERNANCE_CHANNEL="governance-validation-channel"
INVESTMENT_CHANNEL="investment-execution-channel"
GOVERNANCE_CC="governance"
INVESTMENT_CC="investment"

COLLECTIONS_GOVERNANCE="./collections_config_governance.json"
COLLECTIONS_INVESTMENT="./collections_config_investment.json"

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

# Setup Fabric environment paths
setup_fabric_env() {
    export PATH=$PATH:${PROJECT_DIR}/bin
    export FABRIC_CFG_PATH=${PROJECT_DIR}/config
    # Microfab Go builds on install can take several minutes
    export CORE_PEER_CLIENT_CONNTIMEOUT=300s
    print_info "Fabric environment paths configured"
}

# Switch to StartupOrg context
switch_to_startup() {
    export CORE_PEER_LOCALMSPID=StartupOrgMSP
    export CORE_PEER_MSPCONFIGPATH=${MSP_BASE_PATH}/StartupOrg/startuporgadmin/msp
    export CORE_PEER_ADDRESS=startuporgpeer-api.127-0-0-1.nip.io:7070
    setup_fabric_env
    print_success "Switched to StartupOrg context"
}

# Switch to ValidatorOrg context
switch_to_validator() {
    export CORE_PEER_LOCALMSPID=ValidatorOrgMSP
    export CORE_PEER_MSPCONFIGPATH=${MSP_BASE_PATH}/ValidatorOrg/validatororgadmin/msp
    export CORE_PEER_ADDRESS=validatororgpeer-api.127-0-0-1.nip.io:7070
    setup_fabric_env
    print_success "Switched to ValidatorOrg context"
}

# Switch to PlatformOrg context
switch_to_platform() {
    export CORE_PEER_LOCALMSPID=PlatformOrgMSP
    export CORE_PEER_MSPCONFIGPATH=${MSP_BASE_PATH}/PlatformOrg/platformorgadmin/msp
    export CORE_PEER_ADDRESS=platformorgpeer-api.127-0-0-1.nip.io:7070
    setup_fabric_env
    print_success "Switched to PlatformOrg context"
}

# Switch to InvestorOrg context
switch_to_investor() {
    export CORE_PEER_LOCALMSPID=InvestorOrgMSP
    export CORE_PEER_MSPCONFIGPATH=${MSP_BASE_PATH}/InvestorOrg/investororgadmin/msp
    export CORE_PEER_ADDRESS=investororgpeer-api.127-0-0-1.nip.io:7070
    setup_fabric_env
    print_success "Switched to InvestorOrg context"
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
# Channel / CC helpers
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

resolve_cc() {
    case "$1" in
        governance|gov) echo "${GOVERNANCE_CC}" ;;
        investment|inv) echo "${INVESTMENT_CC}" ;;
        *) echo "" ;;
    esac
}

channel_for_cc() {
    case "$1" in
        "${GOVERNANCE_CC}") echo "${GOVERNANCE_CHANNEL}" ;;
        "${INVESTMENT_CC}") echo "${INVESTMENT_CHANNEL}" ;;
        *) echo "" ;;
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

orgs_on_channel() {
    case "$1" in
        "${GOVERNANCE_CHANNEL}") echo "startup validator platform" ;;
        "${INVESTMENT_CHANNEL}") echo "startup validator investor platform" ;;
        *) echo "" ;;
    esac
}

# Which CCs an org should install
ccs_for_org() {
    case "$1" in
        startup|validator|platform) echo "governance investment" ;;
        investor) echo "investment" ;;  # Investor never on governance
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
    for file in ${PROJECT_DIR}/${chaincode_name}_*.tar.gz ${PROJECT_DIR}/${chaincode_name}_*.tgz; do
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
    local contract_path="${CONTRACTS_BASE_PATH}/${chaincode_name}"
    local version
    version=$(get_next_package_version "$chaincode_name")
    local package_file="${PROJECT_DIR}/${chaincode_name}_${version}.tar.gz"
    local label="${chaincode_name}_${version}"

    print_header "📦 Packaging chaincode: ${chaincode_name}"
    setup_fabric_env
    cd "${PROJECT_DIR}"

    if [ ! -d "$contract_path" ]; then
        print_error "Contract path does not exist: ${contract_path}"
        return 1
    fi

    # Same as crowdfundingv2: tidy + vendor so peer external builder can compile offline
    print_info "Running go mod tidy in ${contract_path}..."
    (
        cd "${contract_path}" || exit 1
        go mod tidy || exit 1
        print_info "Running go mod vendor..."
        go mod vendor || exit 1
    )
    if [ $? -ne 0 ]; then
        print_error "go mod tidy/vendor failed for ${chaincode_name}"
        return 1
    fi

    print_info "Package label: ${label}"
    print_info "Package file:  ${package_file}"
    print_info "Path:          ${contract_path}"

    peer lifecycle chaincode package "${package_file}" \
        --path "${contract_path}" \
        --lang golang \
        --label "${label}"

    if [ $? -eq 0 ]; then
        print_success "Successfully packaged ${chaincode_name} as ${package_file}"
        print_info "Package includes go.mod, go.sum, vendor/, and all contract .go files"
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
            print_header "📦 Packaging Dual-Channel Combined Chaincodes"
            package_one governance || return 1
            package_one investment || return 1
            print_success "Both chaincodes packaged successfully!"
            print_warning "Next: source ./deploy_chaincode.sh install-all"
            ;;
        governance|investment)
            package_one "$target" || return 1
            ;;
        *)
            print_error "Usage: package [governance|investment|all]"
            return 1
            ;;
    esac
}

latest_package() {
    local chaincode_name=$1
    ls -t ${PROJECT_DIR}/${chaincode_name}_*.tar.gz ${PROJECT_DIR}/${chaincode_name}_*.tgz 2>/dev/null | head -1
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

    print_info "Installing ${pkg} on current peer (Go build may take several minutes)..."
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
    elif echo "$install_output" | grep -qiE "timeout|Unavailable|keepalive"; then
        # Peer often finishes the build after the CLI times out — check queryinstalled
        print_warning "CLI timed out; checking whether peer finished installing..."
        sleep 5
        package_id=$(peer lifecycle chaincode queryinstalled 2>&1 | grep "Label: ${chaincode_name}_" | head -1 | grep -oP "Package ID: \K[^,]+")
        if [ -n "$package_id" ]; then
            print_success "${chaincode_name} is installed on peer (recovered after timeout)"
        else
            print_error "Install failed for ${chaincode_name} (timeout and not found on peer)"
            print_warning "Wait a bit, then retry: source ./deploy_chaincode.sh install <org>"
            return 1
        fi
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

    print_header "📥 Installing Dual-Channel chaincodes on ${org}"
    switch_org "$org" || return 1

    local ccs
    ccs=$(ccs_for_org "$org")
    for cc in $ccs; do
        print_header "Installing ${cc} on ${org}"
        print_warning "Press Enter to install ${cc}..."
        read -r
        install_one_cc "$cc" || return 1
        echo ""
    done

    print_success "Installation complete on ${org}"
}

install_all() {
    print_header "📥 Installing Dual-Channel Chaincodes on All Organizations"
    for org in startup validator platform investor; do
        install_on_org "$org" || return 1
    done
    echo ""
    print_success "Installation complete!"
    print_info "GOVERNANCE_CC_PACKAGE_ID=${GOVERNANCE_CC_PACKAGE_ID}"
    print_info "INVESTMENT_CC_PACKAGE_ID=${INVESTMENT_CC_PACKAGE_ID}"
    echo ""
    print_warning "Next: source ./deploy_chaincode.sh deploy all"
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
    local cc_alias=$1
    local cc channel
    cc=$(resolve_cc "$cc_alias")
    if [ -z "$cc" ]; then
        print_error "Usage: check-readiness <governance|investment>"
        return 1
    fi
    channel=$(channel_for_cc "$cc")

    switch_to_startup
    local sequence version policy collections
    sequence=$(next_sequence "$channel" "$cc")
    version=$(next_version "$channel" "$cc")
    policy=$(get_policy "$channel")
    collections=$(get_collections "$channel")

    print_info "Checking commit readiness for '${cc}' on '${channel}' (v${version}/seq${sequence})..."
    peer lifecycle chaincode checkcommitreadiness \
        --channelID "${channel}" \
        --name "${cc}" \
        --version "${version}" \
        --sequence "${sequence}" \
        --collections-config "${collections}" \
        --signature-policy "${policy}"
}

# =============================================================================
# Channel-centric deploy (one combined CC per channel)
# =============================================================================

deploy_governance() {
    print_header "🚀 Deploying 'governance' on ${GOVERNANCE_CHANNEL}"
    print_info "Members: StartupOrg, ValidatorOrg, PlatformOrg (NO InvestorOrg)"
    require_package_id governance || return 1

    local pid
    pid=$(get_package_id governance)

    for org in startup validator platform; do
        print_header "Approve governance as ${org}"
        switch_org "$org" || return 1
        print_warning "Press Enter to approve as ${org}..."
        read -r
        approve_chaincode "${GOVERNANCE_CHANNEL}" "${GOVERNANCE_CC}" "${pid}" || return 1
    done

    print_header "Commit governance"
    switch_to_startup
    print_warning "Press Enter to commit governance..."
    read -r
    commit_chaincode "${GOVERNANCE_CHANNEL}" "${GOVERNANCE_CC}" || return 1
    print_success "governance deployed on ${GOVERNANCE_CHANNEL}"
}

deploy_investment() {
    print_header "🚀 Deploying 'investment' on ${INVESTMENT_CHANNEL}"
    print_info "Members: StartupOrg, ValidatorOrg, InvestorOrg, PlatformOrg"
    require_package_id investment || return 1

    local pid
    pid=$(get_package_id investment)

    for org in startup validator investor platform; do
        print_header "Approve investment as ${org}"
        switch_org "$org" || return 1
        print_warning "Press Enter to approve as ${org}..."
        read -r
        approve_chaincode "${INVESTMENT_CHANNEL}" "${INVESTMENT_CC}" "${pid}" || return 1
    done

    print_header "Commit investment"
    switch_to_startup
    print_warning "Press Enter to commit investment..."
    read -r
    commit_chaincode "${INVESTMENT_CHANNEL}" "${INVESTMENT_CC}" || return 1
    print_success "investment deployed on ${INVESTMENT_CHANNEL}"
}

deploy_cmd() {
    case "$1" in
        governance|gov)
            deploy_governance
            ;;
        investment|inv)
            deploy_investment
            ;;
        all)
            deploy_governance || return 1
            deploy_investment || return 1
            print_success "Both dual-channel combined chaincodes deployed!"
            ;;
        *)
            print_error "Usage: deploy <governance|investment|all>"
            return 1
            ;;
    esac
}

# =============================================================================
# Upgrade
# =============================================================================

upgrade_one() {
    local cc_alias=$1
    local cc channel pid
    cc=$(resolve_cc "$cc_alias")
    if [ -z "$cc" ]; then
        print_error "Usage: upgrade <governance|investment>"
        return 1
    fi
    channel=$(channel_for_cc "$cc")
    require_package_id "$cc" || return 1
    pid=$(get_package_id "$cc")

    print_header "⬆️  Upgrading '${cc}' on '${channel}'"

    for org in $(orgs_on_channel "$channel"); do
        print_header "Approve upgrade as ${org}"
        switch_org "$org" || return 1
        print_warning "Press Enter to approve upgrade as ${org}..."
        read -r
        approve_chaincode "$channel" "$cc" "$pid" || return 1
    done

    switch_to_startup
    print_warning "Press Enter to commit upgrade..."
    read -r
    commit_chaincode "$channel" "$cc" || return 1
    print_success "Upgraded '${cc}' on '${channel}'"
}

upgrade_all() {
    upgrade_one governance || return 1
    upgrade_one investment || return 1
}

query_committed() {
    local cc_alias=$1
    local cc channel
    cc=$(resolve_cc "$cc_alias")
    if [ -z "$cc" ]; then
        print_error "Usage: query-committed <governance|investment>"
        return 1
    fi
    channel=$(channel_for_cc "$cc")
    switch_to_startup
    print_info "Committed '${cc}' on '${channel}':"
    peer lifecycle chaincode querycommitted --channelID "${channel}" --name "${cc}"
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
    echo "Crowdfunding Dual-Channel v2 — Combined Chaincode Deploy Tool"
    echo "Port 7070 | 2 Channels | 2 Combined CCs + PDC"
    echo "============================================================================="
    echo ""
    echo "CHANNELS / CHAINCODES:"
    echo "  governance  → ${GOVERNANCE_CHANNEL}"
    echo "                CC:   governance (Startup+Validator+Platform+DataBridge)"
    echo "                Orgs: Startup, Validator, Platform   (NO Investor)"
    echo "                Policy: ${POLICY_GOVERNANCE}"
    echo ""
    echo "  investment  → ${INVESTMENT_CHANNEL}"
    echo "                CC:   investment (all 4 orgs + Token + DataBridge)"
    echo "                Orgs: Startup, Validator, Investor, Platform"
    echo "                Policy: ${POLICY_INVESTMENT}"
    echo ""
    echo "PACKAGING:"
    echo "  package [governance|investment|all]"
    echo ""
    echo "INSTALLATION:"
    echo "  install <org>     Install relevant CCs (investor skips governance)"
    echo "  install-all"
    echo ""
    echo "DEPLOYMENT (channel-centric):"
    echo "  deploy governance | deploy investment | deploy all"
    echo ""
    echo "UPGRADE:"
    echo "  upgrade <governance|investment>"
    echo "  upgrade-all"
    echo ""
    echo "QUERY:"
    echo "  query-committed <governance|investment>"
    echo "  check-readiness <governance|investment>"
    echo "  query-installed <org>"
    echo ""
    echo "UTILITY:"
    echo "  switch <startup|validator|investor|platform>"
    echo "  help"
    echo ""
    echo "INVOKE EXAMPLES:"
    echo "  -n governance -c '{\"function\":\"StartupContract:CreateCampaign\",...}'"
    echo "  -n investment -c '{\"function\":\"InvestorContract:MakeInvestment\",...}'"
    echo ""
    echo "INITIAL FLOW:"
    echo "  source ./deploy_chaincode.sh package"
    echo "  source ./deploy_chaincode.sh install-all"
    echo "  source ./deploy_chaincode.sh deploy all"
    echo ""
    echo "4-org-CC variant: ../crowdfunding_dualchannel/"
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
            print_error "Usage: deploy <governance|investment|all>"
            return 1
        fi
        deploy_cmd "$2"
        ;;
    upgrade)
        upgrade_one "$2"
        ;;
    upgrade-all)
        upgrade_all
        ;;
    query-committed)
        query_committed "$2"
        ;;
    check-readiness)
        check_readiness "$2"
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
