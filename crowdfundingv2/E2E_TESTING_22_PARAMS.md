# E2E Testing Guide - 22 Parameter Campaign Format

## Complete Campaign Lifecycle Testing with PDC

---

## Environment Setup

```bash
# Switch to StartupOrg using deployment script
source ./deploy_chaincode.sh switch startup
```

---

## Test Flow 1: Complete Campaign Lifecycle

### 1.1 INVOKE: Create Campaign (22 Parameters)

```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses investororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses validatororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:9090 -c '{"function":"StartupContract:CreateCampaign","Args":["CAMP004","STARTUP001","Technology","2025-03-31","USD","false","false","2025-01-01","Prototype","Hardware","[\"IoT\",\"SmartHome\",\"AI\"]","false","false","90","1","1","2025","50000","50K-100K","Smart Home IoT Platform","An innovative IoT platform for smart home automation with AI-powered features","[\"business_plan.pdf\",\"pitch_deck.pdf\",\"financials.xlsx\"]"]}'
```

**Parameters Breakdown:**
1. `campaignID`: "CAMP004"
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

### 1.2 QUERY: Verify Campaign Created

```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding -c '{"function":"StartupContract:GetCampaign","Args":["CAMP004"]}'
```

**Expected Response:**
```json
{
  "campaignId": "CAMP004",
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

---

### 1.3 INVOKE: Submit for Validation (StartupValidatorShared PDC)

```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 -c '{"function":"StartupContract:SubmitForValidation","Args":["CAMP004","[\"business_plan_v2.pdf\"]","Please validate our IoT platform campaign"]}'
```

### 1.4 QUERY: Check Validation Status

---

### 1.5 Switch to ValidatorOrg: Validate Campaign

```bash
# Switch to ValidatorOrg
source ./deploy_chaincode.sh switch validator

# View pending campaign (from StartupValidatorShared)
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding -c '{"function":"ValidatorContract:GetCampaign","Args":["CAMP003"]}'

# Validate campaign
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses investororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses validatororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:9090 -c '{"function":"ValidatorContract:Valid
ateCampaign","Args":["VAL001","CAMP001","VALIDATOR001"]}'
```

### 1.6 INVOKE: Approve Campaign (Generates Digital Signature)

```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses investororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses validatororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:9090 -c '{"function":"ValidatorContract:ApproveOrRejectCampaign","Args":["VAL001","CAMP004","APPROVED","8.5","3.2","LOW","[\"Strong technical team\",\"Viable market\"]","[]",""]}'
```

**What Happens:**
- Validator generates digital signature (validationHash)
- Stores approval in StartupValidatorShared (Startup can see)
- Stores approval in ValidatorPlatformShared (Platform can verify)

---

### 1.7 Switch to StartupOrg: Share Campaign with Platform

### 1.7 Switch to StartupOrg: Share Campaign with Platform

```bash
# Switch back to StartupOrg
source ./deploy_chaincode.sh switch startup

# Share campaign with Platform (includes validator hash)
# Note: Use the actual validationHash returned from validator approval
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses investororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses validatororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:9090 -c '{"function":"StartupContract:ShareCampaignToPlatform","Args":["CAMP004","VALIDATOR_HASH_HERE"]}'
```

**What Happens:**
- Startup verifies campaign is APPROVED by validator
- Verifies hash matches what validator provided
- Copies all 22 campaign parameters to StartupPlatformShared
- Campaign status updated to PENDING_PLATFORM_APPROVAL

---

### 1.8 Switch to PlatformOrg: Verify Hash & Publish Campaign

```bash
# Switch to PlatformOrg
source ./deploy_chaincode.sh switch platform

# Verify shared campaign details
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding -c '{"function":"PlatformContract:GetSharedCampaign","Args":["CAMP001"]}'

# Publish campaign (only 2 parameters: campaignID + validationHash)
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses investororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses validatororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:9090 -c '{"function":"PlatformContract:PublishCampaignToPortal","Args":["CAMP004","VALIDATOR_HASH_HERE"]}'
```

**What Happens:**
- Platform reads campaign from StartupPlatformShared
- Platform reads validation approval from ValidatorPlatformShared
- Verifies hash matches (3-way verification)
- If verified: publishes campaign + sends success notification to Startup
- If mismatch: rejects (tampering detected)

---

### 1.9 Switch to StartupOrg: Check Publish Notification

```bash
# Switch back to StartupOrg
source ./deploy_chaincode.sh switch startup

# Check if Platform published the campaign
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding -c '{"function":"StartupContract:CheckPublishNotification","Args":["CAMP004"]}'
```

**Expected Response:**
```json
{
  "campaignId": "CAMP004",
  "status": "PUBLISHED",
  "message": "Campaign 'Smart Home IoT Platform' has been successfully published on the platform",
  "publishedAt": "2025-01-15T10:00:00Z",
  "validationScore": 8.5,
  "riskLevel": "LOW"
}
```

---

### 1.10 Switch to InvestorOrg: View Published Campaign

```bash
# Switch to InvestorOrg peer
export CORE_PEER_LOCALMSPID="InvestorOrgMSP"
export CORE_PEER_ADDRESS=investororgpeer-api.127-0-0-1.nip.io:9090
export CORE_PEER_MSPCONFIGPATH=$HOME/crowdfunding/crowdfundingv2/_msp/InvestorOrg/investororgadmin/msp

# View public campaign details (prior to importing)
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding -c '{"function":"PlatformContract:GetPublishedCampaign","Args":["CAMP001"]}'

# View campaign (stores in InvestorPrivateData)
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses investororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses validatororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:9090 -c '{"function":"InvestorContract:ViewCampaign","Args":["VIEW001","CAMP004","INVESTOR001"]}'
```

---

### 1.11 InvestorOrg: Request Validation Details from Validator

```bash
# Request validation score and risk insights from Validator
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses investororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses validatororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:9090 -c '{"function":"InvestorContract:RequestValidationDetails","Args":["REQ001","CAMP004","INVESTOR001"]}'
```

**What Happens:**
- Investor creates request in InvestorValidatorShared
- Validator can see the request

---

### 1.12 Switch to ValidatorOrg: Provide Validation Details

```bash
# Switch to ValidatorOrg
source ./deploy_chaincode.sh switch validator

# Provide validation details to investor
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses investororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses validatororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:9090 -c '{"function":"ValidatorContract:ProvideValidationDetailsToInvestor","Args":["REQ001","CAMP004"]}'
```

**What Happens:**
- Validator reads validation approval
- Creates response with scores and risk level
- Stores in InvestorValidatorShared

---

### 1.13 Switch to InvestorOrg: Read Validation Response

```bash
# Switch back to InvestorOrg
source ./deploy_chaincode.sh switch investor

# Read validation response
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding -c '{"function":"InvestorContract:GetValidationResponse","Args":["REQ001"]}'
```

**Expected Response:**
```json
{
  "requestId": "REQ001",
  "campaignId": "CAMP004",
  "validatorId": "VALIDATOR001",
  "dueDiligenceScore": 8.5,
  "riskScore": 3.2,
  "riskLevel": "LOW",
  "validationHash": "abc123def456...",
  "approvedAt": "2025-01-15T09:00:00Z",
  "respondedAt": "2025-01-15T11:00:00Z"
}
```

---

### 1.14 InvestorOrg: Make Investment or Negotiate

```bash
# Option 1: Direct investment
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses investororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses validatororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:9090 -c '{"function":"InvestorContract:MakeInvestment","Args":["INV001","CAMP004","INVESTOR001","25000","USD"]}'

# Option 2: Create investment proposal (negotiate with Startup)
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses investororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses validatororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:9090 -c '{"function":"InvestorContract:CreateInvestmentProposal","Args":["PROPOSAL001","CAMP004","INVESTOR001","STARTUP001","25000","USD","15","3 years","[{\"milestoneId\":\"M1\",\"title\":\"Beta Launch\",\"amount\":10000}]","Equity terms with milestone-based release"]}'
```

---

## Test Flow 2: Multiple Campaigns with Different Sectors

### 2.1 Create SaaS Campaign

```bash
export CORE_PEER_LOCALMSPID="StartupOrgMSP"
export CORE_PEER_ADDRESS=startuporgpeer-api.127-0-0-1.nip.io:9090
export CORE_PEER_MSPCONFIGPATH=$HOME/crowdfunding/crowdfundingv2/_msp/StartupOrg/startuporgadmin/msp

peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses investororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses validatororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:9090 -c '{"function":"StartupContract:CreateCampaign","Args":["CAMP004","STARTUP002","SaaS","2025-06-30","USD","true","false","2024-05-15","MVP","Software","[\"SaaS\",\"B2B\",\"Analytics\"]","true","true","120","15","2","2025","150000","100K-500K","Enterprise Analytics Platform","Next-generation analytics platform for enterprise customers with AI-driven insights","[\"pitch_deck.pdf\",\"customer_testimonials.pdf\",\"financial_projections.xlsx\",\"mvp_demo.mp4\"]"]}'
```

### 2.2 Create HealthTech Campaign

```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses investororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses validatororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:9090 -c '{"function":"StartupContract:CreateCampaign","Args":["CAMP004","STARTUP003","HealthTech","2025-04-15","EUR","false","true","2024-11-01","Prototype","Healthcare","[\"HealthTech\",\"Medical\",\"Diagnostics\"]","true","false","60","1","3","2025","75000","50K-100K","AI Medical Diagnostics Tool","AI-powered diagnostic tool for early disease detection with FDA approval in progress","[\"business_plan.pdf\",\"clinical_trial_data.pdf\",\"regulatory_roadmap.pdf\"]"]}'
```

---

## Test Flow 3: Privacy Verification (PDC Testing)

### 3.1 StartupOrg: Query Own Campaign (Should Work)

```bash
source ./deploy_chaincode.sh switch startup

peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding -c '{"function":"StartupContract:GetCampaign","Args":["CAMP004"]}'
```

**Expected:** ✅ Full campaign data returned

### 3.2 InvestorOrg: Try to Query StartupPrivateData (Should Fail)

```bash
export CORE_PEER_LOCALMSPID="InvestorOrgMSP"
export CORE_PEER_ADDRESS=investororgpeer-api.127-0-0-1.nip.io:9090
export CORE_PEER_MSPCONFIGPATH=$HOME/crowdfunding/crowdfundingv2/_msp/InvestorOrg/investororgadmin/msp

# This will fail because InvestorOrg cannot access StartupPrivateData collection
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding -c '{"function":"StartupContract:GetCampaign","Args":["CAMP004"]}'
```

**Expected:** ❌ Error or empty response (cannot access StartupPrivateData)

### 3.3 InvestorOrg: Query Public Campaign Info (Should Work)

```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding -c '{"function":"PlatformContract:GetPublishedCampaign","Args":["CAMP004"]}'
```

**Expected:** ✅ Basic public campaign info (name, category, goal, status)

---

## Test Flow 4: Investment Proposal with 22-Param Campaign

### 4.1 InvestorOrg: Create Investment Proposal

```bash
export CORE_PEER_LOCALMSPID="InvestorOrgMSP"
export CORE_PEER_ADDRESS=investororgpeer-api.127-0-0-1.nip.io:9090
export CORE_PEER_MSPCONFIGPATH=$HOME/crowdfunding/crowdfundingv2/_msp/InvestorOrg/investororgadmin/msp

peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses investororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses validatororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:9090 -c '{"function":"InvestorContract:CreateInvestmentProposal","Args":["PROPOSAL001","CAMP004","INVESTOR001","STARTUP001","25000","USD","15","3 years","[{\"milestoneId\":\"M1\",\"title\":\"Beta Launch\",\"amount\":10000},{\"milestoneId\":\"M2\",\"title\":\"100 Users\",\"amount\":15000}]","Standard equity terms with milestone-based fund release"]}'
```

### 4.2 StartupOrg: Acknowledge Investment

```bash
export CORE_PEER_LOCALMSPID="StartupOrgMSP"
export CORE_PEER_ADDRESS=startuporgpeer-api.127-0-0-1.nip.io:9090
export CORE_PEER_MSPCONFIGPATH=$HOME/crowdfunding/crowdfundingv2/_msp/StartupOrg/startuporgadmin/msp

peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses investororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses validatororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:9090 -c '{"function":"StartupContract:AcknowledgeInvestment","Args":["ACK001","INV001","CAMP004","STARTUP001","INVESTOR001","Thank you for your investment. We are excited to work together!"]}'
```

---

## Test Flow 5: Milestone Submission & Verification

### 5.1 StartupOrg: Submit Milestone Report

```bash
export CORE_PEER_LOCALMSPID="StartupOrgMSP"
export CORE_PEER_ADDRESS=startuporgpeer-api.127-0-0-1.nip.io:9090
export CORE_PEER_MSPCONFIGPATH=$HOME/crowdfunding/crowdfundingv2/_msp/StartupOrg/startuporgadmin/msp

peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses investororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses validatororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:9090 -c '{"function":"StartupContract:SubmitMilestoneReport","Args":["MILESTONE_RPT001","CAMP004","STARTUP001","M1","Beta Launch Completed","Successfully launched beta version with 50 test users. All core features operational.","milestone_evidence_hash_123"]}'
```

### 5.2 ValidatorOrg: Verify Milestone

```bash
source ./deploy_chaincode.sh switch validator

peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses investororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses validatororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:9090 -c '{"function":"ValidatorContract:VerifyMilestoneCompletion","Args":["MILESTONE_VER001","MILESTONE_001","CAMP004","STARTUP001","report_hash_123","true","9.5","Milestone completed as described. Beta platform is functional.","true"]}'
```

### 5.3 PlatformOrg: Release Funds

```bash
export CORE_PEER_LOCALMSPID="PlatformOrgMSP"
export CORE_PEER_ADDRESS=platformorgpeer-api.127-0-0-1.nip.io:9090
export CORE_PEER_MSPCONFIGPATH=$HOME/crowdfunding/crowdfundingv2/_msp/PlatformOrg/platformorgadmin/msp

peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses investororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses validatororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:9090 -c '{"function":"PlatformContract:TriggerFundRelease","Args":["RELEASE001","ESCROW_AGREEMENT_001","AGREEMENT001","CAMP004","M1","STARTUP001","10000","Milestone M1 verified and approved by validator"]}'
```

---

## Test Flow 6: Dispute Resolution

### 6.1 InvestorOrg: Create Dispute

```bash
export CORE_PEER_LOCALMSPID="InvestorOrgMSP"
export CORE_PEER_ADDRESS=investororgpeer-api.127-0-0-1.nip.io:9090
export CORE_PEER_MSPCONFIGPATH=$HOME/crowdfunding/crowdfundingv2/_msp/InvestorOrg/investororgadmin/msp

peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses investororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses validatororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:9090 -c '{"function":"PlatformContract:CreateDispute","Args":["DISPUTE001","INVESTOR","INVESTOR001","STARTUP","STARTUP001","MILESTONE_DISPUTE","CAMP004","AGREEMENT001","Milestone M2 Not Completed","Startup claims 100 users achieved but evidence shows only 45 active users","15000"]}'
```

### 6.2 PlatformOrg: Assign Investigator

```bash
export CORE_PEER_LOCALMSPID="PlatformOrgMSP"
export CORE_PEER_ADDRESS=platformorgpeer-api.127-0-0-1.nip.io:9090
export CORE_PEER_MSPCONFIGPATH=$HOME/crowdfunding/crowdfundingv2/_msp/PlatformOrg/platformorgadmin/msp

peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses investororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses validatororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:9090 -c '{"function":"PlatformContract:AssignInvestigator","Args":["DISPUTE001","VALIDATOR002"]}'
```

---

## Test Flow 7: Wallet & Fee Management

### 7.1 PlatformOrg: Create Wallets

```bash
export CORE_PEER_LOCALMSPID="PlatformOrgMSP"
export CORE_PEER_ADDRESS=platformorgpeer-api.127-0-0-1.nip.io:9090
export CORE_PEER_MSPCONFIGPATH=$HOME/crowdfunding/crowdfundingv2/_msp/PlatformOrg/platformorgadmin/msp

# Create Startup Wallet
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses investororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses validatororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:9090 -c '{"function":"PlatformContract:CreateWallet","Args":["WALLET_STARTUP001","STARTUP001","STARTUP","0"]}'

# Create Investor Wallet
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses investororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses validatororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:9090 -c '{"function":"PlatformContract:CreateWallet","Args":["WALLET_INVESTOR001","INVESTOR001","INVESTOR","100000"]}'
```

### 7.2 PlatformOrg: Set Fee Tiers

```bash
# Campaign Fee Tier
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses investororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses validatororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:9090 -c '{"function":"PlatformContract:SetCampaignFeeTier","Args":["TIER_001","0","100000","5","5% fee for campaigns under $100K"]}'

# Dispute Fee Tier
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses investororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses validatororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:9090 -c '{"function":"PlatformContract:SetDisputeFeeTier","Args":["DISPUTE_TIER_001","0","50000","500","$500 fee for disputes under $50K"]}'
```

### 7.3 PlatformOrg: Collect Campaign Fee

```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses investororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses validatororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:9090 -c '{"function":"PlatformContract:CollectCampaignFee","Args":["FEE_COLLECTION_001","CAMP004","STARTUP001","50000","5"]}'
```

---

## Complete Test Script (All in One)

```bash
#!/bin/bash

# Complete E2E test with 22-parameter campaigns

echo "=== TEST 1: Create Campaign with 22 Parameters ==="
source ./deploy_chaincode.sh switch startup

peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses investororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses validatororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:9090 -c '{"function":"StartupContract:CreateCampaign","Args":["CAMP004","STARTUP001","Technology","2025-03-31","USD","false","false","2025-01-01","Prototype","Hardware","[\"IoT\",\"SmartHome\",\"AI\"]","false","false","90","1","1","2025","50000","50K-100K","Smart Home IoT Platform","An innovative IoT platform for smart home automation with AI-powered features","[\"business_plan.pdf\",\"pitch_deck.pdf\",\"financials.xlsx\"]"]}'

echo "=== TEST 2: Verify Campaign ==="
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding -c '{"function":"StartupContract:GetCampaign","Args":["CAMP004"]}'

echo "=== TEST 3: Submit for Validation ==="
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses investororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses validatororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:9090 -c '{"function":"StartupContract:SubmitForValidation","Args":["VAL001","CAMP004","STARTUP001","Please validate"]}'

echo "=== TEST 4: Validate Campaign (as Validator) ==="
source ./deploy_chaincode.sh switch validator

peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses investororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses validatororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:9090 -c '{"function":"ValidatorContract:ValidateCampaign","Args":["VAL001","CAMP004","VALIDATOR001"]}'

peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses investororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses validatororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:9090 -c '{"function":"ValidatorContract:ApproveOrRejectCampaign","Args":["VAL001","CAMP004","APPROVED","8.5","3.2","LOW","[\"Good campaign\"]","[]",""]}'

echo "=== TEST 5: Startup Shares Campaign with Platform ==="
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses investororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses validatororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:9090 -c '{"function":"StartupContract:ShareCampaignToPlatform","Args":["CAMP004","VALIDATOR_HASH"]}'

echo "=== TEST 6: Platform Verifies and Publishes ==="
source ./deploy_chaincode.sh switch platform

peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID crowdfunding-channel -n crowdfunding --peerAddresses startuporgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses investororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses validatororgpeer-api.127-0-0-1.nip.io:9090 --peerAddresses platformorgpeer-api.127-0-0-1.nip.io:9090 -c '{"function":"PlatformContract:PublishCampaignToPortal","Args":["CAMP004","VALIDATOR_HASH"]}'

echo "=== All Tests Completed ==="
```

---

## Updated Workflow Summary

### Key Changes from Previous Version:

1. **Validator Approval** (Step 1.6):
   - Now generates **digital signature (validationHash)**
   - Parameters: `validationID, campaignID, status, dueDiligenceScore, riskScore, riskLevel, commentsJSON, issuesJSON, requiredDocuments`
   - Stores approval in both `StartupValidatorShared` and `ValidatorPlatformShared`

2. **Startup Shares Campaign** (Step 1.7 - NEW):
   - Function: `ShareCampaignToPlatform(campaignID, validationHash)`
   - Verifies campaign is APPROVED
   - Verifies hash matches validator's hash
   - Copies all 22 parameters to `StartupPlatformShared`

3. **Platform Publishes** (Step 1.8 - UPDATED):
   - Function: `PublishCampaignToPortal(campaignID, validationHash)` - **Only 2 parameters!**
   - Reads campaign from `StartupPlatformShared`
   - Reads validation from `ValidatorPlatformShared`
   - Verifies hash (3-way verification)
   - Publishes if verified, sends success notification to Startup

4. **Startup Checks Notification** (Step 1.9 - NEW):
   - Function: `CheckPublishNotification(campaignID)`
   - Reads success message from Platform

5. **Investor Requests Validation** (Steps 1.11-1.13 - NEW):
   - Investor: `RequestValidationDetails(requestID, campaignID, investorID)`
   - Validator: `ProvideValidationDetailsToInvestor(requestID, campaignID)`
   - Investor: `GetValidationResponse(requestID)` → receives scores and risk level

### Benefits:

✅ **Efficient**: Platform function now takes 2 params (not 25!)  
✅ **Secure**: 3-way hash verification prevents tampering  
✅ **Proper PDC**: Uses collections for data sharing  
✅ **Clear workflow**: Validator → Startup → Platform → Investor

---

## Why collections_config.json is Required

### Without PDC Config:
```bash
# ❌ All data would be visible to ALL organizations
# ❌ StartupOrg private campaigns visible to InvestorOrg
# ❌ InvestorOrg portfolios visible to StartupOrg  
# ❌ ValidatorOrg assessments visible to everyone
# ❌ PlatformOrg wallets/fees visible to all
```

### With PDC Config:
```bash
# ✅ StartupOrg sees only their campaigns (StartupPrivateData)
# ✅ InvestorOrg sees only their investments (InvestorPrivateData)
# ✅ Shared data uses paired collections (StartupInvestorShared)
# ✅ Disputes visible to all for transparency (AllOrgsShared)
# ✅ Platform manages wallets privately (PlatformPrivateData)
```

### Deployment Command:
```bash
peer lifecycle chaincode approveformyorg \
  --channelID crowdfunding-channel \
  --name crowdfunding \
  --version 1.0 \
  --package-id $PACKAGE_ID \
  --sequence 1 \
  --collections-config $HOME/crowdfunding/crowdfundingv2/collections_config.json
  # ☝️ THIS IS MANDATORY for PDC to work!
```

---

## Summary

✅ **22-Parameter Campaign Format Implemented**  
✅ **PDC Collections for Privacy**  
✅ **Complete E2E Testing Flow**  
✅ **Privacy Verification Tests**  
✅ **Dispute Resolution Flow**  
✅ **Wallet & Fee Management**  

**Key Advantage:** All campaign details match your exact API format while maintaining privacy through PDC!
