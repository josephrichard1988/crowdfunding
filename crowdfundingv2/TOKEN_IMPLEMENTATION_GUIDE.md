# Token-Based Fee System Implementation Guide

## Overview

This platform uses **ERC-20 style tokens** for all payments:
- **CFT (CrowdToken)** - Utility token for fees, payments, investments
- **CFRT (CrowdRewardToken)** - Reward tokens convertible to CFT

## Exchange Rates

| Currency | CFT Value |
|----------|-----------|
| **1 INR** | **2.5 CFT** |
| **1 USD** | **83 CFT** |
| **1 CFRT** | **10 CFT** |

---

## ⚠️ Prerequisites

**Before testing, you MUST upgrade the chaincode:**
```bash
cd ~/crowdfunding/crowdfundingv2
./deploy_chaincode.sh upgrade
```

---

## Step-by-Step Testing

### Step 1: Initialize Tokens (Platform Only)

```bash
# Switch to Platform
source ./deploy_chaincode.sh switch platform

# Initialize CFT - 1 billion initial supply, unlimited max
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 \
  --channelID crowdfunding-channel -n crowdfunding \
  --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:9090 \
  -c '{"function":"TokenContract:InitializeToken","Args":["CFT","CrowdToken","2","1000000000","0"]}'

# Initialize CFRT - 0 initial supply, unlimited max
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 \
  --channelID crowdfunding-channel -n crowdfunding \
  --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:9090 \
  -c '{"function":"TokenContract:InitializeToken","Args":["CFRT","CrowdRewardToken","2","0","0"]}'
```

### Step 2: Set Exchange Rates

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

### Step 3: Query Exchange Rate

```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 \
  --channelID crowdfunding-channel -n crowdfunding \
  -c '{"function":"TokenContract:GetExchangeRate","Args":["INR"]}'
# Expected: 2.5
```

---

## Purchase Tokens (Fiat → CFT)

### Startup Purchases Tokens

```bash
source ./deploy_chaincode.sh switch startup

# Startup pays ₹600 INR → receives 1,500 CFT (600 × 2.5)
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 \
  --channelID crowdfunding-channel -n crowdfunding \
  --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 \
  -c '{"function":"TokenContract:PurchaseTokens","Args":["STARTUP001","INR","600"]}'
```

### Check Balance

```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 \
  --channelID crowdfunding-channel -n crowdfunding \
  -c '{"function":"TokenContract:GetBalance","Args":["STARTUP001","CFT"]}'
# Expected: 1500
```

---

## Transfer Tokens (Pay Fees)

```bash
# Pay 250 CFT (₹100) registration fee to Platform
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 \
  --channelID crowdfunding-channel -n crowdfunding \
  --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 \
  -c '{"function":"TokenContract:TransferTokens","Args":["TXN001","CFT","STARTUP001","PLATFORM","250","INR","REGISTRATION_FEE",""]}'

# Pay 1,250 CFT (₹500) campaign creation fee
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 \
  --channelID crowdfunding-channel -n crowdfunding \
  --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 \
  -c '{"function":"TokenContract:TransferTokens","Args":["TXN002","CFT","STARTUP001","PLATFORM","1250","INR","CAMPAIGN_CREATION_FEE","CAMP001"]}'
```

---

## Reward Distribution (Platform → User)

```bash
source ./deploy_chaincode.sh switch platform

# Give 100 CFRT reward to startup
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 \
  --channelID crowdfunding-channel -n crowdfunding \
  --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:9090 \
  -c '{"function":"TokenContract:DistributeRewards","Args":["RWD001","STARTUP001","100","FIRST_CAMPAIGN"]}'
```

### Redeem Rewards (CFRT → CFT)

```bash
source ./deploy_chaincode.sh switch startup

# Convert 100 CFRT → 1,000 CFT (100 × 10)
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 \
  --channelID crowdfunding-channel -n crowdfunding \
  --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 \
  -c '{"function":"TokenContract:RedeemRewardTokens","Args":["STARTUP001","100"]}'
```

---

## Withdrawal (CFT → Fiat)

```bash
# Withdraw 1,000 CFT → ₹400 INR (after 1% fee)
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 \
  --channelID crowdfunding-channel -n crowdfunding \
  --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 \
  -c '{"function":"TokenContract:WithdrawToFiat","Args":["STARTUP001","1000","INR"]}'
```

---

## Platform Fees Summary

| Fee Type | CFT | INR |
|----------|-----|-----|
| Registration | 250 CFT | ₹100 |
| Campaign Creation | 1,250 CFT | ₹500 |
| Campaign Publishing | 2,500 CFT | ₹1,000 |
| Validation | 500 CFT | ₹200 |
| Dispute Filing | 750 CFT | ₹300 |
| Investment | 5% | 5% |
| Withdrawal | 1% | 1% |

---

## Quick Reference

| Function | Purpose |
|----------|---------|
| `InitializeToken` | Create CFT/CFRT tokens |
| `SetExchangeRate` | Set INR/USD rates |
| `GetExchangeRate` | Query current rate |
| `PurchaseTokens` | Fiat → CFT |
| `MintTokens` | Platform mints tokens |
| `TransferTokens` | Transfer CFT between accounts |
| `GetBalance` | Check CFT/CFRT balance |
| `DistributeRewards` | Give CFRT rewards |
| `RedeemRewardTokens` | CFRT → CFT |
| `WithdrawToFiat` | CFT → Fiat |
| `BurnTokens` | Destroy tokens |
| `FreezeTokens` | Freeze for escrow/disputes |
| `UnfreezeTokens` | Unfreeze after resolution |

---

## Legacy Functions (Still Available)

These older functions are still available for backwards compatibility:
- `CreateTokenAccount`
- `IssueTokens`
- `CollectFeeTokens`
- `GetTransferHistory`

