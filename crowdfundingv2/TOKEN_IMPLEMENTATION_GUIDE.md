# Token-Based Fee System Implementation Guide

## Overview

This implementation adds a **token-based payment system** for handling fees, investments, and payments in the crowdfunding platform.

## Token Types

1. **FEE_TOKEN** - Platform fees collected from startups
2. **PAYMENT_TOKEN** - Investment payments and transfers
3. **REWARD_TOKEN** - Rewards and incentives

---

## Architecture

### Token Operations Contract
- **TokenContract** - Manages token accounts, transfers, and fee collection
- Integrated with existing 4 org contracts
- Uses Fabric ledger (not external Token SDK initially for simplicity)

### Key Features
✅ Token account management  
✅ Token issuance (minting)  
✅ Token transfers  
✅ Fee collection with tokens  
✅ Freeze/unfreeze (for escrow)  
✅ Balance tracking  
✅ Transfer history

---

## Token Account Structure

```json
{
  "accountId": "STARTUP001",
  "owner": "STARTUP001",
  "ownerType": "STARTUP",
  "balances": {
    "FEE_TOKEN": 10000,
    "PAYMENT_TOKEN": 50000
  },
  "frozenAmount": {
    "PAYMENT_TOKEN": 5000
  },
  "createdAt": "2025-01-15T10:00:00Z",
  "updatedAt": "2025-01-15T11:00:00Z"
}
```

---

## Complete Workflow with Tokens

### 1. Create Token Accounts for All Users

```bash
# Create token account for Startup
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding -c '{"function":"TokenContract:CreateTokenAccount","Args":["STARTUP001","STARTUP001","STARTUP","{\"PAYMENT_TOKEN\":100000,\"FEE_TOKEN\":0}"]}'

# Create token account for Investor
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding -c '{"function":"TokenContract:CreateTokenAccount","Args":["INVESTOR001","INVESTOR001","INVESTOR","{\"PAYMENT_TOKEN\":500000,\"FEE_TOKEN\":0}"]}'

# Create token account for Platform
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding -c '{"function":"TokenContract:CreateTokenAccount","Args":["PLATFORM","PLATFORM","PLATFORM","{\"PAYMENT_TOKEN\":0,\"FEE_TOKEN\":0}"]}'
```

---

### 2. Issue Initial Tokens (Mint)

```bash
# Issue payment tokens to investors
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding -c '{"function":"TokenContract:IssueTokens","Args":["TOKEN001","PAYMENT_TOKEN","INVESTOR001","100000","USD","PLATFORM","Initial token allocation"]}'

# Issue payment tokens to startups
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding -c '{"function":"TokenContract:IssueTokens","Args":["TOKEN002","PAYMENT_TOKEN","STARTUP001","50000","USD","PLATFORM","Operational tokens"]}'
```

---

### 3. Check Token Balance

```bash
# Query startup's payment token balance
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding -c '{"function":"TokenContract:GetBalance","Args":["STARTUP001","PAYMENT_TOKEN"]}'

# Query investor's payment token balance
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding -c '{"function":"TokenContract:GetBalance","Args":["INVESTOR001","PAYMENT_TOKEN"]}'

# Get full account details
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding -c '{"function":"TokenContract:GetTokenAccount","Args":["STARTUP001"]}'
```

---

### 4. Campaign Fee Collection with Tokens

**Scenario:** Startup raises $50,000, platform collects 5% fee ($2,500)

```bash
# Collect fee tokens from startup
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding -c '{"function":"TokenContract:CollectFeeTokens","Args":["FEE001","CAMP001","STARTUP001","50000","5"]}'
```

**What Happens:**
1. Calculates fee: 5% of $50,000 = $2,500
2. Transfers 2,500 FEE_TOKEN from STARTUP001 to PLATFORM
3. Records fee collection transaction
4. Updates both account balances

---

### 5. Investment Payment with Tokens

**Scenario:** Investor invests $25,000 in campaign

```bash
# Transfer payment tokens from investor to startup
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding -c '{"function":"TokenContract:TransferTokens","Args":["TRANSFER001","PAYMENT_TOKEN","INVESTOR001","STARTUP001","25000","USD","INVESTMENT","CAMP001"]}'
```

**What Happens:**
1. Checks investor has sufficient PAYMENT_TOKEN balance
2. Deducts 25,000 from investor's account
3. Credits 25,000 to startup's account
4. Records transfer transaction with purpose "INVESTMENT"

---

### 6. Freeze Tokens (For Disputes/Escrow)

**Scenario:** Dispute raised, freeze investment amount

```bash
# Freeze tokens in startup's account
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding -c '{"function":"TokenContract:FreezeTokens","Args":["STARTUP001","PAYMENT_TOKEN","25000","Dispute DISPUTE001 - freeze until resolution"]}'
```

**What Happens:**
1. Moves 25,000 from `balances` to `frozenAmount`
2. Tokens still owned by startup but cannot be transferred
3. Records freeze reason

---

### 7. Unfreeze Tokens (After Dispute Resolution)

```bash
# Unfreeze tokens after dispute resolved
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding -c '{"function":"TokenContract:UnfreezeTokens","Args":["STARTUP001","PAYMENT_TOKEN","25000"]}'
```

---

### 8. Query Transfer History

```bash
# Get all transfers for an account
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding -c '{"function":"TokenContract:GetTransferHistory","Args":["STARTUP001"]}'
```

**Response:**
```json
[
  {
    "transferId": "TRANSFER001",
    "tokenType": "PAYMENT_TOKEN",
    "from": "INVESTOR001",
    "to": "STARTUP001",
    "amount": 25000,
    "currency": "USD",
    "purpose": "INVESTMENT",
    "campaignId": "CAMP001",
    "status": "COMPLETED",
    "transferredAt": "2025-01-15T10:00:00Z",
    "completedAt": "2025-01-15T10:00:01Z"
  },
  {
    "transferId": "FEE_TRANSFER_CAMP001_2025-01-15T10:05:00Z",
    "tokenType": "FEE_TOKEN",
    "from": "STARTUP001",
    "to": "PLATFORM",
    "amount": 2500,
    "currency": "USD",
    "purpose": "CAMPAIGN_FEE",
    "campaignId": "CAMP001",
    "status": "COMPLETED",
    "transferredAt": "2025-01-15T10:05:00Z"
  }
]
```

---

## Integration with Existing Contracts

### Update PlatformContract to Use Tokens

Replace old wallet functions with token operations:

**Before:**
```bash
peer chaincode invoke -n platformorg -c '{"function":"CollectCampaignFee","Args":[...]}'
```

**After:**
```bash
peer chaincode invoke -n crowdfunding -c '{"function":"TokenContract:CollectFeeTokens","Args":[...]}'
```

---

## Complete E2E Flow with Tokens

### Step 1: Setup Token Accounts
```bash
# Create accounts for all participants
TokenContract:CreateTokenAccount (STARTUP001, INVESTOR001, PLATFORM, VALIDATOR001)
```

### Step 2: Issue Initial Tokens
```bash
# Platform issues tokens to users
TokenContract:IssueTokens (100000 PAYMENT_TOKEN to INVESTOR001)
TokenContract:IssueTokens (50000 PAYMENT_TOKEN to STARTUP001)
```

### Step 3: Campaign Creation & Validation
```bash
# Normal workflow
StartupContract:CreateCampaign
StartupContract:SubmitForValidation
ValidatorContract:ValidateCampaign
ValidatorContract:ApproveOrRejectCampaign
```

### Step 4: Campaign Publishing
```bash
StartupContract:ShareCampaignToPlatform
PlatformContract:PublishCampaignToPortal
```

### Step 5: Investment with Tokens
```bash
# Investor transfers payment tokens to startup
TokenContract:TransferTokens (INVESTOR001 → STARTUP001, 25000 PAYMENT_TOKEN, "INVESTMENT")

# Record investment in InvestorContract
InvestorContract:MakeInvestment
```

### Step 6: Fee Collection with Tokens
```bash
# Platform collects fee from startup
TokenContract:CollectFeeTokens (STARTUP001, campaignAmount=50000, feePercent=5)
# Automatically transfers 2500 FEE_TOKEN to PLATFORM
```

### Step 7: Dispute Handling
```bash
# If dispute occurs
TokenContract:FreezeTokens (STARTUP001, 25000 PAYMENT_TOKEN)

# After resolution
TokenContract:UnfreezeTokens (STARTUP001, 25000 PAYMENT_TOKEN)
```

---

## Benefits of Token-Based System

✅ **Atomic Transactions** - Transfer and record in single transaction  
✅ **Audit Trail** - Complete transfer history on-chain  
✅ **Escrow Support** - Freeze/unfreeze for dispute resolution  
✅ **Multi-Currency** - Support different token types  
✅ **Transparent** - All transactions recorded on ledger  
✅ **Programmable** - Can add logic (e.g., automatic fee calculation)

---

## Deployment

The token contract is already integrated in `main.go`:

```go
crowdfundingChaincode, err := contractapi.NewChaincode(
    &StartupContract{},
    &InvestorContract{},
    &ValidatorContract{},
    &PlatformContract{},
    &TokenContract{}, // ← Token operations
)
```

Deploy the combined chaincode as usual:
```bash
cd $HOME/crowdfunding/crowdfundingv2/contracts
go mod tidy
peer lifecycle chaincode package crowdfunding-combined.tar.gz --path . --lang golang --label crowdfunding_1.0
# ... install, approve, commit
```

---

## Advanced: Hyperledger Fabric Token SDK (Optional)

For production, you can migrate to **Hyperledger Fabric Token SDK** which provides:
- UTXO-based token model
- Privacy-preserving transfers (zero-knowledge proofs)
- Atomic swaps
- Auditor support

**Migration Path:**
1. Use current simple token implementation for MVP
2. Test and validate business logic
3. Migrate to Fabric Token SDK for production

---

## Testing Commands Summary

```bash
# 1. Create accounts
TokenContract:CreateTokenAccount

# 2. Issue tokens
TokenContract:IssueTokens

# 3. Check balance
TokenContract:GetBalance
TokenContract:GetTokenAccount

# 4. Transfer tokens
TokenContract:TransferTokens

# 5. Collect fees
TokenContract:CollectFeeTokens

# 6. Freeze/unfreeze
TokenContract:FreezeTokens
TokenContract:UnfreezeTokens

# 7. Query history
TokenContract:GetTransferHistory
```

---

## Next Steps

1. ✅ Deploy combined chaincode with TokenContract
2. ✅ Create token accounts for all users
3. ✅ Issue initial token allocations
4. ✅ Test complete workflow with token transfers
5. ✅ Integrate token operations into existing E2E tests
6. Consider Fabric Token SDK for production deployment
