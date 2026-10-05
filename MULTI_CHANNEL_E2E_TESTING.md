# Multi-Channel E2E Testing Guide - 22 Parameter Format (Port 8080)

This comprehensive guide maps the complete 22-parameter campaign lifecycle to the fully isolated **Multi-Channel Architecture**. 

> **Why is this required?** Unlike the single-channel deployment (which uses 1 chaincode and 1 channel), the multi-channel architecture requires you to target the exact chaincode (`-n startup`, `-n validator`, etc.) and the exact channel (`--channelID startup-validator-channel`, etc.) where the isolated transaction belongs.

---

## Test Flow 0.0: Platform Token & Wallet Setup (REQUIRED FIRST)

To operate in the network (e.g., pay fees or invest), all parties need tokens and wallets. The platform must mint the initial supply with a strict maximum cap.

### 0.0.1 Switch to PlatformOrg
```bash
source ./deploy_chaincode.sh switch platform
```

### 0.0.2 INVOKE: Initialize Network Tokens (CFT) with a Max Cap
```bash
# Initialize CFT (CrowdFund Token). 
# Note: 1 INR = 2.5 CFT. SEBI requires some investors to have 2Cr INR (~50M CFT) net worth.
# We initialize 100 Million initial supply and 1 Billion max cap to support enterprise funding.
# The 'TokenContract:' prefix is required because these functions are in the secondary contract.
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID common-channel -n platform --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"TokenContract:InitializeToken","Args":["CFT","CrowdToken","18","100000000","1000000000"]}'
```

### 0.0.3 INVOKE: Transfer Tokens to Startup, Investor & Validator
```bash
# Transfer 5,000 CFT to Startup (for listing fees, etc.)
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID common-channel -n platform --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"TokenContract:TransferTokens","Args":["TRX_STARTUP","CFT","PLATFORM","STARTUP001","5000","INR","INITIAL_FUNDING",""]}'

# Transfer 50,000,000 CFT to Investor (Equivalent to 2 Crore INR SEBI Net Worth requirement)
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID common-channel -n platform --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"TokenContract:TransferTokens","Args":["TRX_INVESTOR","CFT","PLATFORM","INV_user123_001","50000000","INR","INITIAL_FUNDING",""]}'

# Transfer 1,000 CFT to Validator (to stake against disputes)
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID common-channel -n platform --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"TokenContract:TransferTokens","Args":["TRX_VALIDATOR","CFT","PLATFORM","VALIDATOR001","1000","INR","INITIAL_FUNDING",""]}'
```
### 0.0.4 QUERY: Verify Token Balances (Optional)
```bash
# Check Startup Token Account Balance
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID common-channel -n platform -c '{"function":"TokenContract:GetTokenAccount","Args":["STARTUP001"]}'

# Check Investor Token Account Balance
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID common-channel -n platform -c '{"function":"TokenContract:GetTokenAccount","Args":["INV_user123_001"]}'

# Check Validator Token Account Balance
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID common-channel -n platform -c '{"function":"TokenContract:GetTokenAccount","Args":["VALIDATOR001"]}'

# Check Platform Token Account Balance (automatically got initial supply during InitializeToken!)
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID common-channel -n platform -c '{"function":"TokenContract:GetTokenAccount","Args":["PLATFORM"]}'
```


<!-- 
### 0.0.4 INVOKE: Create Legacy Wallets (Backward Compatibility)
```bash
# Create Startup Wallet
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID common-channel -n platform --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"CreateWallet","Args":["WALLET_STARTUP001","STARTUP001","STARTUP","5000"]}'

# Create Investor Wallet
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID common-channel -n platform --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"CreateWallet","Args":["WALLET_INVESTOR001","INV_user123_001","INVESTOR","50000000"]}'
```
-->

### 0.0.5 INVOKE: Set Fee Tiers
```bash
# Set 5% fee for campaigns under $100K
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID common-channel -n platform --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"SetCampaignFeeTier","Args":["TIER_001","0","100000","5","5% fee for campaigns under $100K"]}'
```

---

## Environment Setup for Startup

```bash
# Switch to StartupOrg using deployment script
source ./deploy_chaincode.sh switch startup
```

---

## Test Flow 0: Startup Management (REQUIRED BEFORE CAMPAIGNS)

### 0.1 INVOKE: Create Startup
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-validator-channel -n startup --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"CreateStartup","Args":["STARTUP001","STU_user123_001","TechVentures Inc","An innovative technology startup focused on IoT solutions","S-001"]}'
```

### 0.2 QUERY: Get Startup by ID
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-validator-channel -n startup -c '{"function":"GetStartup","Args":["STARTUP001"]}'
```

### 0.3 QUERY: Get Startups By Owner
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-validator-channel -n startup -c '{"function":"GetStartupsByOwner","Args":["STU_user123_001"]}'
```

---

## Test Flow 0.4: Deletion with Fees (Startup & Campaigns)

> **Fee Structure:**
> - Campaign with funds raised: **60% of funds**
> - Campaign with no funds: **100 CFT fixed**

### 0.4.1 QUERY: Get Campaign Deletion Fee Preview
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-validator-channel -n startup -c '{"function":"CalculateCampaignDeletionFee","Args":["CAMP001"]}'
```

### 0.4.2 INVOKE: Delete Campaign
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-validator-channel -n startup --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"DeleteCampaign","Args":["CAMP001","User requested deletion"]}'
```

### 0.4.3 INVOKE: Delete Startup (and all campaigns)
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-validator-channel -n startup --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"DeleteStartup","Args":["STARTUP002","Closing business"]}'
```

---

## Test Flow 1: Complete Campaign Lifecycle

### 1.1 INVOKE: Create Campaign (22 Parameters)
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-validator-channel -n startup --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"CreateCampaign","Args":["CAMP001","STARTUP001","Technology","2025-03-31","USD","false","false","2025-01-01","Prototype","Hardware","[\"IoT\",\"SmartHome\",\"AI\"]","false","false","90","1","1","2025","50000","50K-100K","Smart Home IoT Platform","An innovative IoT platform for smart home automation with AI-powered features","[\"business_plan.pdf\",\"pitch_deck.pdf\",\"financials.xlsx\"]"]}'
```

### 1.2 QUERY: Verify Campaign Created
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-validator-channel -n startup -c '{"function":"GetCampaign","Args":["CAMP001"]}'
```

### 1.2.1 INVOKE: Update Campaign (Optional Edits before validation)
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-validator-channel -n startup --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"UpdateCampaign","Args":["CAMP001","goalAmount","60000","Increased based on revised budget estimate","STU_user123_001"]}'
```

### 1.3 INVOKE: Submit for Validation
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-validator-channel -n startup --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"SubmitForValidation","Args":["CAMP001","[\"business_plan_v2.pdf\"]","Please validate our IoT platform campaign"]}'
```

### 1.4 Switch to ValidatorOrg: Validate & Approve Campaign
```bash
# Switch Context
source ./deploy_chaincode.sh switch validator
```

### 1.4.1 QUERY: Get Pending Validations (Validator sees what needs review)
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-validator-channel -n validator -c '{"function":"GetPendingValidations","Args":[]}'
```

### 1.4.2 QUERY: Read the Campaign details (Cross-channel read from startup-validator-channel)
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-validator-channel -n validator -c '{"function":"GetCampaign","Args":["CAMP001"]}'
```

### 1.4.3 INVOKE: Validate Campaign (Record validation assessment)
```bash
# Use the submissionHash from the GetCampaign query above
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-validator-channel -n validator --peerAddresses validatororgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"ValidateCampaign","Args":["VAL001","CAMP001","VALIDATOR001","eef2c67e8b0c56858ff5a33137d2a604d08c9b3801c0ace7a7cbfa419670d5c5","[\"business_plan_v2.pdf\"]"]}'
```

### 1.4.4 INVOKE: Approve Campaign (Generates digital validationProofHash)
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-validator-channel -n validator --peerAddresses validatororgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"ApproveOrRejectCampaign","Args":["VAL001","CAMP001","APPROVED","8.5","3.2","LOW","[\"Strong technical team\"]","[]",""]}'
```

### 1.4.5 QUERY: Verify Campaign Status after Approval
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-validator-channel -n validator -c '{"function":"GetCampaign","Args":["CAMP001"]}'
```

### 1.4.1 Switch to ValidatorOrg: Share Validation with Platform
```bash
# Share validation approval using validator-platform-channel
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID validator-platform-channel -n validator --peerAddresses validatororgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"ShareValidationToPlatform","Args":["VAL001", "CAMP001"]}'
```

### 1.5 Switch to StartupOrg: Share Campaign with Platform
```bash
# Switch Context
source ./deploy_chaincode.sh switch startup

# Share campaign using startup-platform-channel
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-platform-channel -n startup --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:8080  -c '{"function":"ShareCampaignToPlatform","Args":["CAMP001", "<VALIDATION_PROOF_HASH_FROM_DO_GET_CAMPAIGN_FROM_STRTUP_AND_YOU_WILL_GET_THE_VALIDATION_PROOF_HASH>"]}'
```

### 1.6 Switch to PlatformOrg: Verify & Publish Campaign
```bash
# Switch Context
source ./deploy_chaincode.sh switch platform
```

### 1.6.1 QUERY: See what the Validator shared (read VALIDATION_APPROVAL from validator-platform-channel)
```bash
# NOTE: Cannot call validator's GetCampaign here — that cross-chains to startup-validator-channel
# where PlatformOrg is NOT a member. Use platform chaincode which reads ValidatorPlatformCollection.
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID validator-platform-channel -n platform -c '{"function":"GetValidationApproval","Args":["CAMP001"]}'
```

### 1.6.2 QUERY: See all campaigns received by platform (from startup via startup-platform-channel)
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-platform-channel -n platform -c '{"function":"GetAllSharedCampaigns","Args":[]}'
```

### 1.6.3 QUERY: Get specific shared campaign details (from startup)
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-platform-channel -n platform -c '{"function":"GetSharedCampaign","Args":["CAMP001"]}'
```

### 1.6.4 INVOKE: Pay Publishing Fee (Startup -> Platform)
```bash
# 1. Startup transfers 2,500 CFT (₹1,000) to Platform for the publishing fee
source ./deploy_chaincode.sh switch startup
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID common-channel -n platform --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"TokenContract:TransferTokens","Args":["CFT","STARTUP001","PLATFORM_01","2500"]}'

# 2. Platform records the fee collection
source ./deploy_chaincode.sh switch platform
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID common-channel -n platform --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"CollectPublishingFee","Args":["PUB_FEE_001","CAMP001","STARTUP001"]}'
```

### 1.6.5 INVOKE: Publish Campaign to Portal (verifies proof hash from BOTH channels)
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID common-channel -n platform --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"PublishCampaignToPortal","Args":["CAMP001","8dcebc524d0e55580f8e0a9e7f8cde79531ff704b4af0732f33ebd30a139950c"]}'
```

### 1.6.5 QUERY: Verify Campaign is now Published
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID common-channel -n platform -c '{"function":"GetPublishedCampaign","Args":["CAMP001"]}'
```

### 1.6.6 INVOKE: Notify Startup (via startup-platform-channel)
```bash
# Because Fabric prevents cross-channel state writes, PlatformOrg must explicitly notify StartupOrg
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-platform-channel -n platform --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"NotifyStartupPublished","Args":["CAMP001"]}'
```

### 1.7 Switch to StartupOrg: Check Publish Notification
```bash
# Switch Context
source ./deploy_chaincode.sh switch startup

peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-platform-channel -n startup -c '{"function":"CheckPublishNotification","Args":["CAMP001"]}'
```

## Test Flow 2: Investor Actions (Browse & Due Diligence)

### 2.1 Switch to InvestorOrg & Browse
```bash
source ./deploy_chaincode.sh switch investor

# Glance View (Public data on common-channel from Platform)
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID common-channel -n platform -c '{"function":"GetActiveCampaigns","Args":[]}'

# View Campaign (Stores view log in Investor's local chaincode)
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID investor-platform-channel -n investor --peerAddresses investororgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"ViewCampaign","Args":["VIEW001","CAMP001","INVESTOR001"]}'
```

### 2.2 Request Risk Insights from Validator
```bash

# Investor requests details from Validator on investor-validator-channel
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID investor-validator-channel -n investor --peerAddresses investororgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"RequestValidationDetails","Args":["REQ001","CAMP001","INVESTOR001"]}'
```

### 2.3 Switch to ValidatorOrg: View Request & Respond
```bash
source ./deploy_chaincode.sh switch validator
```

### 2.3.1 QUERY: Validator Views Request (Cross-channel read from investor-validator-channel)
```bash

# Get all pending validation requests (from Investors)
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID investor-validator-channel -n validator -c '{"function":"GetPendingInvestorRequests","Args":[]}'

# Validator views the specific incoming request from the Investor
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID investor-validator-channel -n investor -c '{"function":"ReadData","Args":["InvestorValidatorShared","VALIDATION_REQUEST_REQ001"]}'
```

### 2.3.2 INVOKE: Validator Responds
```bash
# Validator provides response on investor-validator-channel
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID investor-validator-channel -n validator --peerAddresses validatororgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"ProvideValidationDetailsToInvestor","Args":["REQ001","CAMP001"]}'
```

### 2.4 Investor Reads Response
```bash
source ./deploy_chaincode.sh switch investor

peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID investor-validator-channel -n investor -c '{"function":"GetValidationResponse","Args":["REQ001"]}'
```

---

## Test Flow 3: KYC, Proposals & Direct Investment

### 3.1 Register Investor (SEBI Compliant)
> **Note:** Annual income MUST be >= 20000000 (2 Crore INR).
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID investor-platform-channel -n investor --peerAddresses investororgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"RegisterInvestor","Args":["INV_user123_001","Alice Wealthy","25000000","123412341234","ABCDE1234F"]}'
```

### 3.2 Create Investment Proposal (Negotiation)
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-investor-channel -n investor --peerAddresses investororgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"CreateInvestmentProposal","Args":["PROPOSAL001","CAMP001","INV_user123_001","STARTUP001","25000","USD","15","3 years","[{\"milestoneId\":\"M1\",\"title\":\"Beta Launch\",\"amount\":10000}]","Standard equity terms"]}'
```

### 3.2.1 Startup Views and Counters the Proposal
```bash
source ./deploy_chaincode.sh switch startup

# Startup views all proposals for their campaign to see what came in
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-investor-channel -n investor -c '{"function":"GetProposalsByCampaign","Args":["CAMP001"]}'

# Startup responds with a counter offer
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-investor-channel -n startup --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"RespondToInvestmentProposal","Args":["PROPOSAL001","COUNTER","30000","Need more funds for beta launch","[]"]}'
```

### 3.2.2 Investor Views and Counters the Startup's Counter
```bash
source ./deploy_chaincode.sh switch investor

# Investor views all their own proposals across different campaigns
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-investor-channel -n investor -c '{"function":"GetProposalsByInvestor","Args":["INV_user123_001"]}'

# Investor views the specific proposal to see the Startup's counter offer
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-investor-channel -n investor -c '{"function":"GetProposal","Args":["PROPOSAL001"]}'

# Investor counters back
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-investor-channel -n investor --peerAddresses investororgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"RespondToCounterOffer","Args":["PROPOSAL001","COUNTER","28000","Final offer for 15% equity"]}'
```

### 3.2.3 Startup Views and Accepts the Final Offer
```bash
source ./deploy_chaincode.sh switch startup

# Startup views the final counter offer from Investor
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-investor-channel -n investor -c '{"function":"GetProposal","Args":["PROPOSAL001"]}'

# Startup accepts the final offer
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-investor-channel -n startup --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"RespondToInvestmentProposal","Args":["PROPOSAL001","ACCEPT","0","","[]"]}'
```

### 3.2.4 View the Final Negotiated Proposal and Audit History
```bash
# We query the Investor chaincode to see the proposal and the negotiation history
source ./deploy_chaincode.sh switch investor

peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-investor-channel -n investor -c '{"function":"GetProposal","Args":["PROPOSAL001"]}'
```


### 3.3 Make Direct Investment
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-investor-channel -n investor --peerAddresses investororgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"MakeInvestment","Args":["INVEST_001","CAMP001","INV_user123_001","25000","USD"]}'

# Investor transfers the funds to the Platform's Escrow account on the common-channel
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID common-channel -n platform --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"TokenContract:TransferTokens","Args":["TRX_ESCROW_001","CFT","INV_user123_001","PLATFORM","25000","USD","INVESTMENT_ESCROW","CAMP001"]}'
```

### 3.3.1 Investor Views their Investments
```bash
# Because of Hyperledger Fabric's cross-chaincode iterator limitations, we must query the startup chaincode directly to retrieve lists.

# Investor views all investments using GetInvestmentsByInvestor directly on startup chaincode
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-investor-channel -n startup -c '{"function":"GetInvestmentsByInvestor","Args":["INV_user123_001"]}'

# Investor views all investments made into a SPECIFIC campaign (CAMP001) directly on startup chaincode
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-investor-channel -n startup -c '{"function":"GetInvestmentsByCampaign","Args":["CAMP001"]}'
```

### 3.4 Startup Views and Acknowledges Investment
```bash
source ./deploy_chaincode.sh switch startup

# Startup views all investments by querying startup chaincode directly
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-investor-channel -n startup -c '{"function":"GetInvestmentsByInvestor","Args":[""]}'

# Startup views the specific investment details directly from the shared state
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-investor-channel -n startup -c '{"function":"ReadData","Args":["StartupInvestorShared","INVESTMENT_INVEST_001"]}'

# Startup acknowledges the investment
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-investor-channel -n startup --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"AcknowledgeInvestment","Args":["ACK001","INVEST_001","CAMP001","STARTUP001","INV_user123_001","Thank you for your investment!"]}'
```

### 3.4.1 Investor Views the Acknowledgment
```bash
source ./deploy_chaincode.sh switch investor

# Investor can now read the newly created Acknowledgment record using the shared collection
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-investor-channel -n startup -c '{"function":"ReadData","Args":["StartupInvestorShared","ACK_ACK001"]}'
```

---

## Test Flow 4: Milestone Management

### 4.1 Startup Submits Milestone
```bash
source ./deploy_chaincode.sh switch startup

peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-validator-channel -n startup --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"SubmitMilestoneReport","Args":["MILESTONE_RPT001","CAMP001","M1","AGREEMENT001","Beta Launch Completed","Successfully launched beta version.","[\"milestone_evidence_hash_123\"]"]}'
```

### 4.1.1 Validator Views Pending Milestone Reports
```bash
source ./deploy_chaincode.sh switch validator

# View all milestone reports (querying startup chaincode directly because iterators don't work cross-chaincode)
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-validator-channel -n startup -c '{"function":"GetAllMilestoneReports","Args":[]}'

# View a specific milestone report
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-validator-channel -n startup -c '{"function":"GetMilestoneReport","Args":["MILESTONE_RPT001"]}'
```

### 4.2 Validator Verifies Milestone
```bash
source ./deploy_chaincode.sh switch validator

peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-validator-channel -n validator --peerAddresses validatororgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"VerifyMilestoneCompletion","Args":["MILESTONE_VER001","MILESTONE_001","CAMP001","STARTUP001","MILESTONE_RPT001","true","9.5","Milestone completed as described.","true"]}'
```

### 4.2.1 Validator Views Verification Record
```bash
source ./deploy_chaincode.sh switch validator

# View the newly updated Milestone Report to see the VERIFIED status
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-validator-channel -n startup -c '{"function":"GetMilestoneReport","Args":["MILESTONE_RPT001"]}'

# View the detailed MilestoneValidation audit record (which contains the qualityScore)
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-validator-channel -n startup -c '{"function":"ReadData","Args":["StartupValidatorShared","MILESTONE_VERIFICATION_MILESTONE_VER001"]}'
```

### 4.2.2 Validator Shares Milestone with Platform
```bash
source ./deploy_chaincode.sh switch validator

# Validator shares the milestone report to the validator-platform-channel
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID validator-platform-channel -n validator --peerAddresses validatororgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"ShareMilestoneToPlatform","Args":["MILESTONE_RPT001"]}'
```

### 4.2.3 Platform Requests Milestone Validation
```bash
source ./deploy_chaincode.sh switch platform

# Platform formally requests the validator to validate the shared milestone
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID validator-platform-channel -n platform --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"RequestMilestoneValidation","Args":["REQ001","MILESTONE_RPT001","CAMP001"]}'
```

### 4.2.4 Validator Shares Verification with Platform
```bash
source ./deploy_chaincode.sh switch validator

# After verifying the milestone, the Validator shares the verification record to the validator-platform-channel
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID validator-platform-channel -n validator --peerAddresses validatororgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"ShareMilestoneValidationToPlatform","Args":["MILESTONE_VER001"]}'
```

### 4.3 Platform Releases Funds
```bash
source ./deploy_chaincode.sh switch platform

# Platform views the validation report shared by the Validator before releasing funds
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID validator-platform-channel -n validator -c '{"function":"ReadData","Args":["ValidatorPlatformShared","MILESTONE_VERIFICATION_MILESTONE_VER001"]}'

# Platform triggers fund release, providing the verification ID it received from the Validator
# This must be invoked on startup-platform-channel so the release record can be written to the StartupPlatformShared collection
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-platform-channel -n platform --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"TriggerFundRelease","Args":["RELEASE001","ESCROW_001","AGREEMENT001","CAMP001","M1","STARTUP001","10000","MILESTONE_VER001"]}'

# Now the Platform actually executes the token transfer on the common-channel
# This automated function reads the release amount from the ledger, ensures it's only executed once, and restricts access to the platform org
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID common-channel -n platform --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"ExecuteFundReleaseTransfer","Args":["RELEASE001"]}'
```

### 4.4 Startup Views Validation Report and Released Funds
```bash
source ./deploy_chaincode.sh switch startup

# Startup can view the validation report directly from the shared collection
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-validator-channel -n startup -c '{"function":"ReadData","Args":["StartupValidatorShared","MILESTONE_VERIFICATION_MILESTONE_VER001"]}'

# Startup can view the fund release confirmation
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-platform-channel -n startup -c '{"function":"ReadData","Args":["StartupPlatformShared","FUND_RELEASE_RELEASE001"]}'
```

---

## Test Flow 5: Disputes & Privacy (Channel Isolation)

Disputes are recorded on the `common-channel` because they typically require visibility and arbitration by the platform. Below are several comprehensive dispute scenarios representing real-world crowdfunding conflicts.

### 5.1 Scenario A: Investor vs Startup (Resolved in favor of Startup)
```bash
source ./deploy_chaincode.sh switch investor

# 1. Raise dispute visible to all orgs via common-channel routing
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID common-channel -n platform --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"CreateDispute","Args":["DISPUTE001","INVESTOR","INV_user123_001","STARTUP","STARTUP001","MILESTONE_DISPUTE","CAMP001","AGREEMENT001","Milestone M2 Fake","Startup claims 100 users but has 0","15000"]}'

# 2. Investor queries the ledger to verify the dispute was recorded successfully
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID common-channel -n platform -c '{"function":"GetDispute","Args":["DISPUTE001"]}'

# 3. Switch to Platform to resolve the dispute
source ./deploy_chaincode.sh switch platform

# 4. Platform reviews and decides Startup actually met the milestone (Respondent Wins)
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID common-channel -n platform --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"ResolveDispute","Args":["DISPUTE001","RESPONDENT_WINS","Startup provided valid logs proving 100 active users. Dispute dismissed."]}'

# 5. Switch to Startup context to verify they won the dispute
source ./deploy_chaincode.sh switch startup
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID common-channel -n platform -c '{"function":"GetDispute","Args":["DISPUTE001"]}'
```

### 5.2 Scenario B: Investor vs Startup (Resolved in favor of Investor)
```bash
source ./deploy_chaincode.sh switch investor

# 1. Raise another dispute
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID common-channel -n platform --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"CreateDispute","Args":["DISPUTE002","INVESTOR","INV_user123_001","STARTUP","STARTUP001","MISUSE_OF_FUNDS","CAMP001","AGREEMENT001","Funds used for personal expenses","Audit shows funds not used for hardware","25000"]}'

# 2. Investor queries to verify dispute is recorded as 'PENDING'
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID common-channel -n platform -c '{"function":"GetDispute","Args":["DISPUTE002"]}'

# 3. Switch to Platform
source ./deploy_chaincode.sh switch platform

# 4. Platform agrees with Investor (Initiator Wins)
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID common-channel -n platform --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"ResolveDispute","Args":["DISPUTE002","INITIATOR_WINS","Audit confirmed misuse of funds. Escrow frozen and refund initiated."]}'

# 5. Investor verifies the resolution before requesting refund
source ./deploy_chaincode.sh switch investor
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID common-channel -n platform -c '{"function":"GetDispute","Args":["DISPUTE002"]}'
```

### 5.2.1 Post-Dispute Fund Transfer: Processing the Refund (10% Retain)
```bash
source ./deploy_chaincode.sh switch investor

# 1. Investor requests a refund for their original investment of 20,000,000 INR
# Passes '10' for the 10% Platform Refund Retain Percent as per the network rules
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID investor-platform-channel -n investor --peerAddresses investororgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"RequestRefund","Args":["REF_REQ_001","INV_user123_001","CAMP001","AGREEMENT001","STARTUP001","20000000","18000000","Dispute won due to misused funds","10"]}'

# 2. Switch to Platform to process and execute the refund
source ./deploy_chaincode.sh switch platform

# 3. Platform processes the refund back to the Investor's wallet
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID common-channel -n platform --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"ProcessRefund","Args":["REF_001","DISPUTE002","INV_user123_001","18000000"]}'

# 4. Investor checks their token balance to ensure the refund arrived
source ./deploy_chaincode.sh switch investor
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID common-channel -n platform -c '{"function":"TokenContract:GetTokenAccount","Args":["INV_user123_001"]}'
```

### 5.3 Scenario C: Startup vs Validator (Negligence Dispute)
```bash
source ./deploy_chaincode.sh switch startup

# 1. Startup raises dispute against Validator for unfair rejection/negligence
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID common-channel -n platform --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"CreateDispute","Args":["DISPUTE003","STARTUP","STARTUP001","VALIDATOR","VALIDATOR001","VALIDATION_DISPUTE","CAMP001","VAL001","Unfair Risk Assessment","Validator scored us HIGH risk due to a misread document","0"]}'

# 2. Startup queries to verify dispute
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID common-channel -n platform -c '{"function":"GetDispute","Args":["DISPUTE003"]}'

# 3. Switch to Platform
source ./deploy_chaincode.sh switch platform

# 4. Platform resolves the dispute
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID common-channel -n platform --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"ResolveDispute","Args":["DISPUTE003","INITIATOR_WINS","Platform review confirms validator error. Validator score penalized."]}'

# 5. Validator views the outcome of the dispute raised against them
source ./deploy_chaincode.sh switch validator
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID common-channel -n platform -c '{"function":"GetDispute","Args":["DISPUTE003"]}'
```

### 5.4 Scenario D: Investor vs Platform (Platform Fee Dispute)
```bash
source ./deploy_chaincode.sh switch investor

# 1. Investor disputes a platform fee calculation
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID common-channel -n platform --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"CreateDispute","Args":["DISPUTE004","INVESTOR","INV_user123_001","PLATFORM","PLATFORM_01","FEE_DISPUTE","CAMP001","FEE_001","Incorrect Fee Deduction","Platform charged 5% instead of 2% tier","3000"]}'

# 2. Investor queries to verify
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID common-channel -n platform -c '{"function":"GetDispute","Args":["DISPUTE004"]}'

# 3. Switch to Platform
source ./deploy_chaincode.sh switch platform

# 4. Platform acknowledges mistake and resolves (Initiator Wins)
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID common-channel -n platform --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"ResolveDispute","Args":["DISPUTE004","INITIATOR_WINS","System glitch acknowledged. Overcharged fee refunded to wallet."]}'

# 5. Switch back to Investor to verify resolution
source ./deploy_chaincode.sh switch investor
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID common-channel -n platform -c '{"function":"GetDispute","Args":["DISPUTE004"]}'
```

### 5.5 Querying Disputes (All Organizations)
```bash
# Any organization (Investor, Startup, Validator, or Platform) can query the status and outcome of a dispute.
# Since disputes are handled globally, they are queried on the common-channel via the platform chaincode.

# 1. Switch to any context (e.g., Startup)
source ./deploy_chaincode.sh switch startup

# 2. Query the details of a specific dispute (e.g., DISPUTE001)
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID common-channel -n platform -c '{"function":"GetDispute","Args":["DISPUTE001"]}'

# 3. Query another dispute (e.g., DISPUTE002)
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID common-channel -n platform -c '{"function":"GetDispute","Args":["DISPUTE002"]}'
```

### 5.6 Test Channel Privacy Constraint
```bash
source ./deploy_chaincode.sh switch investor

# This MUST fail: Investor cannot read Startup private data directly on startup channels
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-validator-channel -n startup -c '{"function":"GetCampaign","Args":["CAMP001"]}'
```
**Expected Outcome**: Error (channel isolation enforcement). Investor can only read campaigns through Platform's public data or their own `startup-investor-channel` interactions.

---

---

## Complete Multi-Channel Test Script (All in One)

```bash
#!/bin/bash

# Complete E2E test with 22-parameter campaigns on MULTI-CHANNEL

echo "=== TEST 1: Create Campaign (StartupOrg) ==="
source ./deploy_chaincode.sh switch startup

peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-validator-channel -n startup --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"CreateCampaign","Args":["CAMP001","STARTUP001","Technology","2025-03-31","USD","false","false","2025-01-01","Prototype","Hardware","[\"IoT\",\"SmartHome\",\"AI\"]","false","false","90","1","1","2025","50000","50K-100K","Smart Home IoT Platform","An innovative IoT platform for smart home automation with AI-powered features","[\"business_plan.pdf\",\"pitch_deck.pdf\",\"financials.xlsx\"]"]}'

echo "=== TEST 2: Verify Campaign ==="
peer chaincode query -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-validator-channel -n startup -c '{"function":"GetCampaign","Args":["CAMP001"]}'

echo "=== TEST 3: Submit for Validation ==="
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-validator-channel -n startup --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"SubmitForValidation","Args":["CAMP001","[\"business_plan_v2.pdf\"]","Please validate"]}'

echo "=== TEST 4: Validate Campaign (ValidatorOrg) ==="
source ./deploy_chaincode.sh switch validator

peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-validator-channel -n validator --peerAddresses validatororgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"ValidateCampaign","Args":["VAL001","CAMP001","VALIDATOR001","<SUBMISSION_HASH>","[\"business_plan.pdf\"]"]}'

peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-validator-channel -n validator --peerAddresses validatororgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"ApproveOrRejectCampaign","Args":["VAL001","CAMP001","APPROVED","8.5","3.2","LOW","[\"Good campaign\"]","[]",""]}'

echo "=== TEST 5: Startup Shares Campaign with Platform ==="
source ./deploy_chaincode.sh switch startup
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID startup-platform-channel -n startup --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"ShareCampaignToPlatform","Args":["CAMP001","<VALIDATOR_HASH>"]}'

echo "=== TEST 6: Platform Verifies and Publishes ==="
source ./deploy_chaincode.sh switch platform

peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:8080 --channelID common-channel -n platform --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:8080 -c '{"function":"PublishCampaignToPortal","Args":["CAMP001","<VALIDATOR_HASH>"]}'

echo "=== All Tests Completed ==="
```

---

## Expected JSON Responses & Parameter Breakdowns

*Note: For brevity in the main flow, the detailed JSON responses and parameter breakdowns are documented here.*

### Parameters Breakdown for CreateCampaign
1. `campaignID`: "CAMP001"
2. `startupID`: "STARTUP001"
3. `category`: "Technology"
4. `deadline`: "2025-03-31"
5. `currency`: "USD"
6. `hasRaised`: "false"
7. `hasGovGrants`: "false"
8. `incorpDate`: "2025-01-01"
9. `projectStage`: "Prototype"
10. `sector`: "Hardware"
11. `tags`: ["IoT","SmartHome","AI"]
12. `teamAvailable`: "false"
13. `investorCommitted`: "false"
14. `duration`: "90"
15. `fundingDay`: "1"
16. `fundingMonth`: "1"
17. `fundingYear`: "2025"
18. `goalAmount`: "50000"
19. `investmentRange`: "50K-100K"
20. `projectName`: "Smart Home IoT Platform"
21. `description`: "An innovative IoT platform for smart home automation with AI-powered features"
22. `documents`: ["business_plan.pdf","pitch_deck.pdf","financials.xlsx"]

### Expected Response: GetCampaign
```json
{
  "campaignId": "CAMP001",
  "startupId": "STARTUP001",
  "category": "Technology",
  "deadline": "2025-03-31",
  "currency": "USD",
  "has_raised": false,
  "has_gov_grants": false,
  "incorp_date": "2025-01-01",
  "project_stage": "Prototype",
  "sector": "Hardware",
  "tags": ["IoT","SmartHome","AI"],
  "team_available": false,
  "investor_committed": false,
  "duration": 90,
  "funding_day": 1,
  "funding_month": 1,
  "funding_year": 2025,
  "goal_amount": 50000,
  "investment_range": "50K-100K",
  "project_name": "Smart Home IoT Platform",
  "description": "An innovative IoT platform for smart home automation with AI-powered features",
  "documents": ["business_plan.pdf","pitch_deck.pdf","financials.xlsx"],
  "open_date": "2025-01-01",
  "close_date": "2025-03-31",
  "funds_raised_amount": 0,
  "funds_raised_percent": 0,
  "status": "DRAFT",
  "validationStatus": "NOT_SUBMITTED",
  "validationScore": 0,
  "investorCount": 0
}
```

### Expected Response: CheckPublishNotification
```json
{
  "campaignId": "CAMP001",
  "status": "PUBLISHED",
  "message": "Campaign 'Smart Home IoT Platform' has been successfully published on the platform",
  "publishedAt": "2025-01-15T10:00:00Z",
  "validationScore": 8.5,
  "riskLevel": "LOW"
}
```

### Expected Response: GetValidationResponse
```json
{
  "requestId": "REQ001",
  "campaignId": "CAMP001",
  "validatorId": "VALIDATOR001",
  "dueDiligenceScore": 8.5,
  "riskScore": 3.2,
  "riskLevel": "LOW",
  "validationHash": "abc123def456...",
  "approvedAt": "2025-01-15T09:00:00Z",
  "respondedAt": "2025-01-15T11:00:00Z"
}
```
