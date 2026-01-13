# Combined Chaincode Deployment Guide

## Overview

All 4 organization contracts are now combined into a single chaincode package:
- **StartupContract** - Campaign management
- **InvestorContract** - Investment operations  
- **ValidatorContract** - Validation and risk assessment
- **PlatformContract** - Portal, wallet, fees, disputes

## File Structure

```
crowdfundingv2/
├── contracts/
│   ├── main.go                      # Main entry point (registers all contracts)
│   ├── go.mod                       # Go module file
│   ├── startuporg_contract.go       # StartupOrg contract
│   ├── investororg_contract.go      # InvestorOrg contract
│   ├── validatororg_contract.go     # ValidatorOrg contract
│   ├── platformorg_contract.go      # PlatformOrg contract
│   └── startuporg/                  # Original directories (kept for reference)
│       ├── go.mod
│       └── startuporg.go
└── collections_config.json          # PDC configuration
```

---

## Step 1: Prepare Combined Chaincode

```bash
cd $HOME/crowdfunding/crowdfundingv2/contracts

# Run go mod vendor to download dependencies
go mod tidy
go mod vendor
```

---

## Step 2: Package Combined Chaincode

```bash
# Package the combined chaincode (single package for all orgs)
peer lifecycle chaincode package crowdfunding-combined.tar.gz \
  --path . \
  --lang golang \
  --label crowdfunding_1.0

# Verify package created
ls -lh crowdfunding-combined.tar.gz
```

---

## Step 3: Install on All Peers

### Install on StartupOrg Peer

```bash
export CORE_PEER_LOCALMSPID="StartupOrgMSP"
export CORE_PEER_ADDRESS=localhost:7051
export CORE_PEER_MSPCONFIGPATH=$HOME/crowdfunding/_msp/StartupOrg/startuporgadmin/msp

peer lifecycle chaincode install crowdfunding-combined.tar.gz
```

### Install on InvestorOrg Peer

```bash
export CORE_PEER_LOCALMSPID="InvestorOrgMSP"
export CORE_PEER_ADDRESS=localhost:8051
export CORE_PEER_MSPCONFIGPATH=$HOME/crowdfunding/_msp/InvestorOrg/investororgadmin/msp

peer lifecycle chaincode install crowdfunding-combined.tar.gz
```

### Install on ValidatorOrg Peer

```bash
export CORE_PEER_LOCALMSPID="ValidatorOrgMSP"
export CORE_PEER_ADDRESS=localhost:9051
export CORE_PEER_MSPCONFIGPATH=$HOME/crowdfunding/_msp/ValidatorOrg/validatororgadmin/msp

peer lifecycle chaincode install crowdfunding-combined.tar.gz
```

### Install on PlatformOrg Peer

```bash
export CORE_PEER_LOCALMSPID="PlatformOrgMSP"
export CORE_PEER_ADDRESS=localhost:10051
export CORE_PEER_MSPCONFIGPATH=$HOME/crowdfunding/_msp/PlatformOrg/platformorgadmin/msp

peer lifecycle chaincode install crowdfunding-combined.tar.gz
```

---

## Step 4: Get Package ID

```bash
# Query installed chaincodes on any peer
peer lifecycle chaincode queryinstalled

# Output will show:
# Installed chaincodes on peer:
# Package ID: crowdfunding_1.0:abc123def456..., Label: crowdfunding_1.0

# Set the package ID as environment variable
export PACKAGE_ID=crowdfunding_1.0:abc123def456...
```

---

## Step 5: Approve Chaincode for Each Org (with collections_config.json)

### Approve for StartupOrg

```bash
export CORE_PEER_LOCALMSPID="StartupOrgMSP"
export CORE_PEER_ADDRESS=localhost:7051
export CORE_PEER_MSPCONFIGPATH=$HOME/crowdfunding/_msp/StartupOrg/startuporgadmin/msp

peer lifecycle chaincode approveformyorg \
  -o orderer-api.127-0-0-1.nip.io:9090 \
  --channelID crowdfunding-channel \
  --name crowdfunding \
  --version 1.0 \
  --package-id $PACKAGE_ID \
  --sequence 1 \
  --collections-config $HOME/crowdfunding/crowdfundingv2/collections_config.json
```

### Approve for InvestorOrg

```bash
export CORE_PEER_LOCALMSPID="InvestorOrgMSP"
export CORE_PEER_ADDRESS=localhost:8051
export CORE_PEER_MSPCONFIGPATH=$HOME/crowdfunding/_msp/InvestorOrg/investororgadmin/msp

peer lifecycle chaincode approveformyorg \
  -o orderer-api.127-0-0-1.nip.io:9090 \
  --channelID crowdfunding-channel \
  --name crowdfunding \
  --version 1.0 \
  --package-id $PACKAGE_ID \
  --sequence 1 \
  --collections-config $HOME/crowdfunding/crowdfundingv2/collections_config.json
```

### Approve for ValidatorOrg

```bash
export CORE_PEER_LOCALMSPID="ValidatorOrgMSP"
export CORE_PEER_ADDRESS=localhost:9051
export CORE_PEER_MSPCONFIGPATH=$HOME/crowdfunding/_msp/ValidatorOrg/validatororgadmin/msp

peer lifecycle chaincode approveformyorg \
  -o orderer-api.127-0-0-1.nip.io:9090 \
  --channelID crowdfunding-channel \
  --name crowdfunding \
  --version 1.0 \
  --package-id $PACKAGE_ID \
  --sequence 1 \
  --collections-config $HOME/crowdfunding/crowdfundingv2/collections_config.json
```

### Approve for PlatformOrg

```bash
export CORE_PEER_LOCALMSPID="PlatformOrgMSP"
export CORE_PEER_ADDRESS=localhost:10051
export CORE_PEER_MSPCONFIGPATH=$HOME/crowdfunding/_msp/PlatformOrg/platformorgadmin/msp

peer lifecycle chaincode approveformyorg \
  -o orderer-api.127-0-0-1.nip.io:9090 \
  --channelID crowdfunding-channel \
  --name crowdfunding \
  --version 1.0 \
  --package-id $PACKAGE_ID \
  --sequence 1 \
  --collections-config $HOME/crowdfunding/crowdfundingv2/collections_config.json
```

---

## Step 6: Check Commit Readiness

```bash
peer lifecycle chaincode checkcommitreadiness \
  --channelID crowdfunding-channel \
  --name crowdfunding \
  --version 1.0 \
  --sequence 1 \
  --collections-config $HOME/crowdfunding/crowdfundingv2/collections_config.json
```

**Expected Output:**
```
Chaincode definition for chaincode 'crowdfunding', version '1.0', sequence '1' on channel 'crowdfunding-channel' approval status by org:
StartupOrgMSP: true
InvestorOrgMSP: true
ValidatorOrgMSP: true
PlatformOrgMSP: true
```

---

## Step 7: Commit Chaincode to Channel

```bash
peer lifecycle chaincode commit \
  -o orderer-api.127-0-0-1.nip.io:9090 \
  --channelID crowdfunding-channel \
  --name crowdfunding \
  --version 1.0 \
  --sequence 1 \
  --collections-config $HOME/crowdfunding/crowdfundingv2/collections_config.json \
  --peerAddresses localhost:7051 --tlsRootCertFiles $HOME/crowdfunding/_msp/StartupOrg/startuporgadmin/msp/cacerts/ca.crt \
  --peerAddresses localhost:8051 --tlsRootCertFiles $HOME/crowdfunding/_msp/InvestorOrg/investororgadmin/msp/cacerts/ca.crt \
  --peerAddresses localhost:9051 --tlsRootCertFiles $HOME/crowdfunding/_msp/ValidatorOrg/validatororgadmin/msp/cacerts/ca.crt \
  --peerAddresses localhost:10051 --tlsRootCertFiles $HOME/crowdfunding/_msp/PlatformOrg/platformorgadmin/msp/cacerts/ca.crt
```

**Note:** If TLS is not enabled in Microfab, remove the `--tlsRootCertFiles` flags.

---

## Step 8: Query Committed Chaincode

```bash
peer lifecycle chaincode querycommitted \
  --channelID crowdfunding-channel \
  --name crowdfunding
```

---

## Step 9: Test Each Contract

### Test StartupContract

```bash
export CORE_PEER_LOCALMSPID="StartupOrgMSP"
export CORE_PEER_ADDRESS=localhost:7051
export CORE_PEER_MSPCONFIGPATH=$HOME/crowdfunding/_msp/StartupOrg/startuporgadmin/msp

# Note: Use -n crowdfunding (combined chaincode name)
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding -c '{"function":"StartupContract:CreateCampaign","Args":["CAMP001","STARTUP001","Technology","2025-03-31","USD","false","false","2025-01-01","Prototype","Hardware","[\"IoT\",\"SmartHome\"]","false","false","90","1","1","2025","50000","50K-100K","Test Campaign","Test description","[\"doc1.pdf\"]"]}'
```

### Test InvestorContract

```bash
export CORE_PEER_LOCALMSPID="InvestorOrgMSP"
export CORE_PEER_ADDRESS=localhost:8051
export CORE_PEER_MSPCONFIGPATH=$HOME/crowdfunding/_msp/InvestorOrg/investororgadmin/msp

peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding -c '{"function":"InvestorContract:ViewCampaign","Args":["VIEW001","CAMP001","INVESTOR001"]}'
```

### Test ValidatorContract

```bash
export CORE_PEER_LOCALMSPID="ValidatorOrgMSP"
export CORE_PEER_ADDRESS=localhost:9051
export CORE_PEER_MSPCONFIGPATH=$HOME/crowdfunding/_msp/ValidatorOrg/validatororgadmin/msp

peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding -c '{"function":"ValidatorContract:ValidateCampaign","Args":["VAL001","CAMP001","VALIDATOR001","hash123","[]"]}'
```

### Test PlatformContract

```bash
export CORE_PEER_LOCALMSPID="PlatformOrgMSP"
export CORE_PEER_ADDRESS=localhost:10051
export CORE_PEER_MSPCONFIGPATH=$HOME/crowdfunding/_msp/PlatformOrg/platformorgadmin/msp

peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding -c '{"function":"PlatformContract:CreateWallet","Args":["WALLET001","STARTUP001","STARTUP","0"]}'
```

---

## Important Notes

### Function Invocation Format

When using combined chaincode, you must prefix function names with the contract name:

**Format:** `ContractName:FunctionName`

**Examples:**
- `StartupContract:CreateCampaign`
- `InvestorContract:MakeInvestment`
- `ValidatorContract:ValidateCampaign`
- `PlatformContract:PublishCampaignToPortal`

### Single Chaincode Name

All contracts are accessed through the same chaincode name: `crowdfunding`

**Before (separate chaincodes):**
```bash
-n startuporg -c '{"function":"CreateCampaign",...}'
-n investororg -c '{"function":"MakeInvestment",...}'
```

**After (combined chaincode):**
```bash
-n crowdfunding -c '{"function":"StartupContract:CreateCampaign",...}'
-n crowdfunding -c '{"function":"InvestorContract:MakeInvestment",...}'
```

---

## Benefits of Combined Chaincode

✅ **Single Deployment** - One chaincode package for all organizations  
✅ **Easier Management** - Single upgrade process for all contracts  
✅ **Shared State** - All contracts access the same ledger state  
✅ **PDC Support** - All collections defined in one collections_config.json  
✅ **Atomic Transactions** - Cross-contract transactions are easier

---

## Troubleshooting

### Issue: Contract not found

**Error:** `Error: contract with name 'StartupContract' not found`

**Solution:** Ensure you're using the correct contract name prefix:
```bash
# Correct
-c '{"function":"StartupContract:CreateCampaign",...}'

# Wrong
-c '{"function":"CreateCampaign",...}'
```

### Issue: Package ID mismatch

**Solution:** Query installed chaincode and get the exact package ID:
```bash
peer lifecycle chaincode queryinstalled
```

### Issue: Collections config not loaded

**Solution:** Ensure collections_config.json path is correct in all approve and commit commands:
```bash
--collections-config $HOME/crowdfunding/crowdfundingv2/collections_config.json
```

---

## Upgrade Process

To upgrade the chaincode:

1. Update contract code
2. Package with new version: `crowdfunding_2.0`
3. Install on all peers
4. Approve with `--sequence 2 --version 2.0`
5. Commit with `--sequence 2 --version 2.0`

---

## Summary

✅ All 4 contracts combined in single chaincode package  
✅ Single deployment process  
✅ Function calls use `ContractName:FunctionName` format  
✅ All contracts accessed through `-n crowdfunding`  
✅ PDC collections work across all contracts  
✅ Ready for E2E testing with UPDATED_WORKFLOW.md
