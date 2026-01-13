# Complete E2E Testing Guide with Token Integration

## Overview

This guide combines campaign lifecycle testing with the token-based fee system. **All fees are paid by the Startup** using CFT tokens.

## Fee Schedule

| Action | Fee (CFT) | Recipient |
|--------|-----------|-----------|
| Campaign Creation | **10 CFT** | Platform |
| Submit to Validation | **50 CFT** | Validator |
| Publish to Portal | **50 CFT** | Platform |
| **Total Campaign Journey** | **110 CFT** | |

## Exchange Rates

| Currency | CFT Value |
|----------|-----------|
| 1 INR | 2.5 CFT |
| 1 USD | 83 CFT |
| 1 CFRT | 10 CFT (reward redemption) |

---

## Phase 1: Token System Setup (One-Time)

### 1.1 Initialize Tokens (Platform Only)

```bash
source ./deploy_chaincode.sh switch platform

# Initialize CFT - 1 billion initial supply
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 \
  --channelID crowdfunding-channel -n crowdfunding \
  --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:9090 \
  -c '{"function":"TokenContract:InitializeToken","Args":["CFT","CrowdToken","2","1000000000","0"]}'

# Initialize CFRT (Rewards)
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 \
  --channelID crowdfunding-channel -n crowdfunding \
  --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:9090 \
  -c '{"function":"TokenContract:InitializeToken","Args":["CFRT","CrowdRewardToken","2","0","0"]}'
```

### 1.2 Set Exchange Rates

```bash
# 1 INR = 2.5 CFT
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 \
  --channelID crowdfunding-channel -n crowdfunding \
  --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:9090 \
  -c '{"function":"TokenContract:SetExchangeRate","Args":["INR","2.5"]}'

# 1 USD = 83 CFT  
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 \
  --channelID crowdfunding-channel -n crowdfunding \
  --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:9090 \
  -c '{"function":"TokenContract:SetExchangeRate","Args":["USD","83"]}'
```

---

## Phase 2: Startup Purchases Tokens

```bash
source ./deploy_chaincode.sh switch startup

# Startup buys ₹100 INR worth of tokens (receives 250 CFT)
# This covers: 10 + 50 + 50 = 110 CFT for full campaign journey
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 \
  --channelID crowdfunding-channel -n crowdfunding \
  --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 \
  -c '{"function":"TokenContract:PurchaseTokens","Args":["STARTUP001","INR","100"]}'

# Verify balance (should show 250 CFT)
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 \
  --channelID crowdfunding-channel -n crowdfunding \
  -c '{"function":"TokenContract:GetBalance","Args":["STARTUP001","CFT"]}'
```

---

## Phase 3: Campaign Lifecycle with Token Fees

### 3.1 Pay Creation Fee + Create Campaign (23 Parameters)

```bash
# Step 1: Pay campaign creation fee (10 CFT to Platform)
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 \
  --channelID crowdfunding-channel -n crowdfunding \
  --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 \
  -c '{"function":"TokenContract:TransferTokens","Args":["STARTUP001","PLATFORM_TREASURY","10","CFT","Campaign creation fee for CAMP001"]}'

# Step 2: Create campaign with 22 parameters + feePaymentTxId
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 \
  --channelID crowdfunding-channel -n crowdfunding \
  --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 \
  -c '{"function":"StartupContract:CreateCampaign","Args":["CAMP001","STARTUP001","Technology","2025-03-31","USD","false","false","2025-01-01","Prototype","Hardware","[\"IoT\",\"SmartHome\",\"AI\"]","false","false","90","1","1","2025","50000","50K-100K","Smart Home IoT Platform","An innovative IoT platform for smart home automation","[\"business_plan.pdf\",\"pitch_deck.pdf\"]"]}'
```

### 3.2 Pay Validation Fee + Submit for Validation

```bash
# Step 1: Pay validation fee (50 CFT to Validator)
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 \
  --channelID crowdfunding-channel -n crowdfunding \
  --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 \
  -c '{"function":"TokenContract:TransferTokens","Args":["STARTUP001","VALIDATOR_TREASURY","50","CFT","Validation fee for CAMP001"]}'

# Step 2: Submit for validation
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 \
  --channelID crowdfunding-channel -n crowdfunding \
  --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 \
  -c '{"function":"StartupContract:SubmitForValidation","Args":["CAMP001","[\"business_plan_v2.pdf\"]","Please validate our IoT platform campaign"]}'
```

### 3.3 Validator Approves Campaign

```bash
source ./deploy_chaincode.sh switch validator

# Validate and approve
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 \
  --channelID crowdfunding-channel -n crowdfunding \
  --peerAddresses validatororgpeer-api.127-0-0-1.nip.io:9090 \
  -c '{"function":"ValidatorContract:ApproveOrRejectCampaign","Args":["VAL001","CAMP001","APPROVED","8.5","3.2","LOW","[\"Strong team\",\"Viable market\"]","[]",""]}'
```

### 3.4 Pay Publishing Fee + Share to Platform

```bash
source ./deploy_chaincode.sh switch startup

# Step 1: Pay publishing fee (50 CFT to Platform)
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 \
  --channelID crowdfunding-channel -n crowdfunding \
  --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 \
  -c '{"function":"TokenContract:TransferTokens","Args":["STARTUP001","PLATFORM_TREASURY","50","CFT","Publishing fee for CAMP001"]}'

# Step 2: Share to platform (use validationHash from step 3.3)
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 \
  --channelID crowdfunding-channel -n crowdfunding \
  --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 \
  -c '{"function":"StartupContract:ShareCampaignToPlatform","Args":["CAMP001","<VALIDATION_HASH>"]}'
```

### 3.5 Platform Publishes Campaign

```bash
source ./deploy_chaincode.sh switch platform

peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 \
  --channelID crowdfunding-channel -n crowdfunding \
  --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:9090 \
  -c '{"function":"PlatformContract:PublishCampaign","Args":["CAMP001","<VALIDATION_HASH>"]}'
```

---

## Phase 4: Investor Flow

```bash
source ./deploy_chaincode.sh switch investor

# Investor purchases tokens
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 \
  --channelID crowdfunding-channel -n crowdfunding \
  --peerAddresses investororgpeer-api.127-0-0-1.nip.io:9090 \
  -c '{"function":"TokenContract:PurchaseTokens","Args":["INVESTOR001","USD","1000"]}'

# Invest in campaign (5% fee auto-deducted)
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 \
  --channelID crowdfunding-channel -n crowdfunding \
  --peerAddresses investororgpeer-api.127-0-0-1.nip.io:9090 \
  -c '{"function":"InvestorContract:InvestInCampaign","Args":["INVESTOR001","CAMP001","10000","USD"]}'
```

---

## Token Balance Summary (After Full Campaign)

| Entity | Starting CFT | After Campaign | Change |
|--------|--------------|----------------|--------|
| STARTUP001 | 250 | 140 | -110 (fees) |
| PLATFORM_TREASURY | 0 | 60 | +60 (10+50) |
| VALIDATOR_TREASURY | 0 | 50 | +50 |

---

## Quick Reference: All Token Functions

| Function | Description |
|----------|-------------|
| `TokenContract:InitializeToken` | Create token type (Platform only) |
| `TokenContract:SetExchangeRate` | Set fiat→CFT rate |
| `TokenContract:PurchaseTokens` | Buy CFT with fiat |
| `TokenContract:TransferTokens` | Send CFT to another account |
| `TokenContract:GetBalance` | Check account balance |
| `TokenContract:WithdrawToFiat` | Convert CFT back to fiat |
| `TokenContract:RedeemRewardTokens` | Convert CFRT to CFT |
