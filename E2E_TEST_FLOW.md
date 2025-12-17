# End-to-End Test Flow - Crowdfunding Platform

## 🎯 Complete Flow: "Smart Home IoT Platform" Campaign

This document provides all invoke and query commands in sequential order for complete end-to-end testing.

---

## ⚙️ PREREQUISITES

Before running commands, set the peer context for each organization:

```bash
# For StartupOrg
export CORE_PEER_LOCALMSPID=StartupOrgMSP
export CORE_PEER_MSPCONFIGPATH=/home/kajal/crowdfunding/_msp/StartupOrg/startuporgadmin/msp
export CORE_PEER_ADDRESS=startuporgpeer-api.127-0-0-1.nip.io:9090

# For ValidatorOrg
export CORE_PEER_LOCALMSPID=ValidatorOrgMSP
export CORE_PEER_MSPCONFIGPATH=/home/kajal/crowdfunding/_msp/ValidatorOrg/validatororgadmin/msp
export CORE_PEER_ADDRESS=validatororgpeer-api.127-0-0-1.nip.io:9090

# For InvestorOrg
export CORE_PEER_LOCALMSPID=InvestorOrgMSP
export CORE_PEER_MSPCONFIGPATH=/home/kajal/crowdfunding/_msp/InvestorOrg/investororgadmin/msp
export CORE_PEER_ADDRESS=investororgpeer-api.127-0-0-1.nip.io:9090

# For PlatformOrg
export CORE_PEER_LOCALMSPID=PlatformOrgMSP
export CORE_PEER_MSPCONFIGPATH=/home/kajal/crowdfunding/_msp/PlatformOrg/platformorgadmin/msp
export CORE_PEER_ADDRESS=platformorgpeer-api.127-0-0-1.nip.io:9090
```

---

# 🚀 PHASE 1: CAMPAIGN CREATION
**Organization: StartupOrg | Channel: startup-validator-channel**

### 1.1 INVOKE: Create Campaign
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n startup -c '{"function":"CreateCampaign","Args":["CAMP001","STARTUP001","Technology","2025-03-31","USD","false","false","2025-01-01","Prototype","Hardware","[\"IoT\",\"SmartHome\",\"AI\"]","false","false","90","1","1","2025","50000","50K-100K","Smart Home IoT Platform","An innovative IoT platform for smart home automation with AI-powered features","[\"business_plan.pdf\",\"pitch_deck.pdf\",\"financials.xlsx\"]"]}'
```

### 1.2 QUERY: Verify Campaign Created
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n startup -c '{"function":"GetCampaign","Args":["CAMP001"]}'
```

### 1.3 OPTIONAL: Update Campaign (if errors found before validation)
**Purpose:** Allow startup to update campaign details BEFORE submitting for validation
**All changes are recorded in ledger history**

#### Example 1: Update Campaign Title
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n startup -c '{"function":"UpdateCampaign","Args":["CAMP001","title","Smart Home IoT Platform v2","Fixed typo in original title"]}'
```

#### Example 2: Update Goal Amount
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n startup -c '{"function":"UpdateCampaign","Args":["CAMP001","goalAmount","60000","Increased based on revised budget estimate"]}'
```

#### Example 3: Update Description
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n startup -c '{"function":"UpdateCampaign","Args":["CAMP001","description","An innovative IoT platform for smart home automation with AI-powered features and voice control","Added voice control feature to description"]}'
```

#### Example 4: Update Tags
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n startup -c '{"function":"UpdateCampaign","Args":["CAMP001","tags","[\"IoT\",\"SmartHome\",\"AI\",\"VoiceControl\"]","Added VoiceControl tag"]}'
```

### 1.3.1 QUERY: Get Campaign Update History
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n startup -c '{"function":"GetCampaignUpdateHistory","Args":["CAMP001"]}'
```

### 1.4 INVOKE: Submit for Validation
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n startup -c '{"function":"SubmitForValidation","Args":["CAMP001","Please review our IoT platform proposal"]}'
```

### 1.5 QUERY: Verify Status Changed to PENDING_VALIDATION
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n startup -c '{"function":"GetCampaign","Args":["CAMP001"]}'
```

### 1.6 QUERY: Get Campaign Validation Hash (for Validator)
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n startup -c '{"function":"GetCampaignValidationHash","Args":["CAMP001"]}'
```

---

# 🔍 PHASE 2: CAMPAIGN VALIDATION
**Organization: ValidatorOrg | Channel: startup-validator-channel**

> ⚠️ **KNOWN ISSUE:** Querying the campaign from ValidatorOrg chaincode before validation may return a 500 error. 
> This is because the validator chaincode doesn't have direct access to StartupOrg's private data.
> **Workaround:** Use the startup chaincode to query the campaign (both orgs share startup-validator-channel):

### 2.1 QUERY: Validator Views Campaign Before Validating (Using Startup Chaincode)
**NOTE:** Query through startup chaincode since both orgs share startup-validator-channel
```bash
# Set ValidatorOrg context first, then query startup chaincode
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n startup -c '{"function":"GetCampaign","Args":["CAMP001"]}'
```

### 2.2 INVOKE: Validate Campaign (APPROVED)
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n validator -c '{"function":"ValidateCampaign","Args":["VAL001","CAMP001","<HASH_FROM_STEP_1.6>","VALIDATOR001","true","true","8.5","2.5","APPROVED","[\"Documents verified\",\"Team credentials confirmed\"]",""]}'
```

### 2.3 QUERY: Verify Validation Record
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n validator -c '{"function":"GetValidation","Args":["VAL001"]}'
```

### 2.4 QUERY: Verify Campaign Hash
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n validator -c '{"function":"VerifyCampaignHash","Args":["CAMP001","<HASH_FROM_STEP_1.6>"]}'
```

### 2.5 QUERY: Check Campaign Not Blacklisted
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n validator -c '{"function":"IsCampaignBlacklisted","Args":["CAMP001"]}'
```

### 2.6 QUERY: Get All Validations for Campaign
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n validator -c '{"function":"GetValidationsByCampaign","Args":["CAMP001"]}'
```

---

# 📤 PHASE 3: VALIDATOR SENDS APPROVED CAMPAIGN TO PLATFORM
**Organization: ValidatorOrg | Channel: validator-platform-channel**
**NOTE:** ValidatorOrg sends approved campaign to PlatformOrg (NOT StartupOrg)

### 3.1 INVOKE: Send Validation Report to Platform
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID validator-platform-channel -n validator -c '{"function":"SendValidationReportToPlatform","Args":["REPORT001","CAMP001","VAL001","<HASH_FROM_STEP_1.6>","8.5","9.0","8.0","2.5","true","Campaign fully verified. Low risk. Recommended for publication."]}'
```

### 3.2 QUERY: Verify Validation Report Sent
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID validator-platform-channel -n validator -c '{"function":"GetValidationReport","Args":["CAMP001"]}'
```

---

# 🌐 PHASE 4: PLATFORM RECEIVES & PUBLISHES CAMPAIGN
**Organization: PlatformOrg | Channel: validator-platform-channel → common-channel**
**NOTE:** Platform queries from validator-platform-channel, then publishes to common-channel

### 4.1 QUERY: Platform Views Validation Report (from ValidatorOrg)
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID validator-platform-channel -n platform -c '{"function":"GetValidationReport","Args":["CAMP001"]}'
```

### 4.2 INVOKE: Record Validator Decision
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID validator-platform-channel -n platform -c '{"function":"RecordValidatorDecision","Args":["REC001","CAMP001","VAL001","<HASH_FROM_STEP_1.6>","true","8.5","report_hash_here"]}'
```

### 4.3 QUERY: Verify Validator Decision Recorded
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID validator-platform-channel -n platform -c '{"function":"GetValidatorDecision","Args":["CAMP001"]}'
```

### 4.4 INVOKE: Publish Campaign to Common Channel
**NOTE:** Now only requires campaignID - fetches all data from validator decision and StartupOrg
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"PublishCampaignToPortal","Args":["CAMP001"]}'
```

### 4.5 QUERY: Verify Published Campaign on Common Channel
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"GetPublishedCampaign","Args":["CAMP001"]}'
```

### 4.6 QUERY: Get Active Campaigns
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"GetActiveCampaigns","Args":[]}'
```

---

# 📢 PHASE 5: PLATFORM NOTIFIES STARTUP OF PUBLICATION
**Organization: PlatformOrg | Channel: startup-platform-channel**
**NOTE:** Platform updates StartupOrg's ledger about publication status

### 5.1 INVOKE: Notify Startup of Campaign Publication
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-platform-channel -n platform -c '{"function":"NotifyStartupOfPublication","Args":["CAMP001","PUBLISHED","Campaign successfully published to portal"]}'
```

### 5.2 QUERY: Verify Notification Sent
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-platform-channel -n platform -c '{"function":"GetPublicationNotification","Args":["CAMP001"]}'
```

**Switch to StartupOrg**

### 5.3 QUERY: Startup Verifies Publication Status
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-platform-channel -n startup -c '{"function":"GetPublicationStatus","Args":["CAMP001"]}'
```

---

# 👀 PHASE 6: INVESTOR VIEWS & INVESTS
**Organization: InvestorOrg | Channel: common-channel → platform-investor-channel**
**NOTE:** Investor first views campaign from common-channel, then invests on platform-investor-channel

### 6.1 QUERY: Investor Views Published Campaign (from common-channel)
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"GetPublishedCampaign","Args":["CAMP001"]}'
```

### 6.2 QUERY: Get All Active Campaigns
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"GetActiveCampaigns","Args":[]}'
```

### 6.3 INVOKE: Record Campaign View (InvestorOrg records interest)
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID platform-investor-channel -n investor -c '{"function":"ViewCampaign","Args":["CAMP001","INV001","Smart Home IoT Platform","Technology","An innovative IoT platform","50000","0","USD","2025-01-01","2025-03-31","Prototype","Hardware","[\"IoT\",\"SmartHome\",\"AI\"]","90","8.5","LOW","0","PUBLISHED"]}'
```

### 6.4 INVOKE: Make Investment Commitment
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID platform-investor-channel -n investor -c '{"function":"MakeInvestment","Args":["INV_001","CAMP001","INV001","10000","USD"]}'
```

### 6.5 QUERY: Verify Investment Created
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID platform-investor-channel -n investor -c '{"function":"GetInvestment","Args":["INV_001"]}'
```

### 6.6 QUERY: Get All Investments by Investor
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID platform-investor-channel -n investor -c '{"function":"GetInvestmentsByInvestor","Args":["INV001"]}'
```

### 6.7 QUERY: Get All Investments for Campaign
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID platform-investor-channel -n investor -c '{"function":"GetInvestmentsByCampaign","Args":["CAMP001"]}'
```

---

# ⚠️ PHASE 7: RISK INSIGHTS REQUEST
**Organization: InvestorOrg | Channel: investor-validator-channel**

### 7.1 INVOKE: Request Risk Insights
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID investor-validator-channel -n investor -c '{"function":"RequestRiskInsights","Args":["RISK_REQ001","CAMP001","INV001"]}'
```

**Switch to ValidatorOrg**

### 7.2 INVOKE: Validator Assigns Risk Score
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID investor-validator-channel -n validator -c '{"function":"AssignRiskScore","Args":["INSIGHT001","CAMP001","INV001","2.5","[\"Strong team\",\"Good market potential\",\"Early stage product\"]","What are the main risks?","The main risks are market competition and execution timeline. Team has strong track record.","RECOMMENDED - Low risk investment with good potential"]}'
```

### 7.3 QUERY: Verify Risk Insight
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID investor-validator-channel -n validator -c '{"function":"GetRiskInsight","Args":["CAMP001"]}'
```

---

# 🤝 PHASE 8: INVESTMENT PROPOSAL & NEGOTIATION
**Organization: InvestorOrg | Channel: startup-investor-channel**

### 8.1 INVOKE: Create Investment Proposal
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-investor-channel -n investor -c '{"function":"CreateInvestmentProposal","Args":["PROP001","CAMP001","STARTUP001","INV001","25000","USD","10% equity stake with board observer rights","[{\"milestoneId\":\"MS001\",\"title\":\"Prototype Development\",\"description\":\"Complete working prototype\",\"targetDate\":\"2025-02-01\",\"fundPercentage\":30,\"status\":\"PENDING\",\"fundsReleased\":false,\"releasedAt\":\"\"},{\"milestoneId\":\"MS002\",\"title\":\"Beta Testing\",\"description\":\"Complete beta testing\",\"targetDate\":\"2025-02-28\",\"fundPercentage\":40,\"status\":\"PENDING\",\"fundsReleased\":false,\"releasedAt\":\"\"},{\"milestoneId\":\"MS003\",\"title\":\"Production Launch\",\"description\":\"Launch production\",\"targetDate\":\"2025-03-31\",\"fundPercentage\":30,\"status\":\"PENDING\",\"fundsReleased\":false,\"releasedAt\":\"\"}]"]}'
```

### 8.2 QUERY: Verify Proposal Created (InvestorOrg)
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-investor-channel -n investor -c '{"function":"GetProposal","Args":["PROP001"]}'
```

### 8.3 QUERY: Get All Proposals by Investor
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-investor-channel -n investor -c '{"function":"GetProposalsByInvestor","Args":["INV001"]}'
```

**Switch to StartupOrg**

### 8.4 QUERY: Startup Views Proposal
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-investor-channel -n startup -c '{"function":"GetProposal","Args":["PROP001"]}'
```

### 8.5 QUERY: Get All Proposals for Startup
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-investor-channel -n startup -c '{"function":"GetProposalsByStartup","Args":["STARTUP001"]}'
```

### 8.6 QUERY: Get All Proposals for Campaign
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-investor-channel -n startup -c '{"function":"GetProposalsByCampaign","Args":["CAMP001"]}'
```

### 8.7 INVOKE: Startup Counter Offer
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-investor-channel -n startup -c '{"function":"RespondToInvestmentProposal","Args":["PROP001","COUNTER","8% equity stake with quarterly updates","30000"]}'
```

### 8.8 QUERY: Verify Proposal Status Updated
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-investor-channel -n startup -c '{"function":"GetProposal","Args":["PROP001"]}'
```

**Switch to InvestorOrg**

### 8.9 QUERY: Investor Views Counter Offer from Startup
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-investor-channel -n investor -c '{"function":"GetProposal","Args":["PROP001"]}'
```

### 8.10 QUERY: Get All Proposals Pending Response
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-investor-channel -n investor -c '{"function":"GetProposalsByInvestor","Args":["INV001"]}'
```

### 8.11 INVOKE: Investor Accepts Counter Offer
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-investor-channel -n investor -c '{"function":"RespondToCounterOffer","Args":["PROP001","INV001","ACCEPT","9% equity stake - final offer","27500"]}'
```

### 8.12 INVOKE: Investor Formally Accepts Agreement
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-investor-channel -n investor -c '{"function":"AcceptAgreement","Args":["PROP001","AGR001","INV001"]}'
```

### 8.13 QUERY: Verify Proposal Final Status
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-investor-channel -n investor -c '{"function":"GetProposal","Args":["PROP001"]}'
```

---

# ✅ PHASE 9: PLATFORM & VALIDATOR WITNESS AGREEMENT
**Organization: PlatformOrg | Channel: common-channel**

### 9.1 INVOKE: Platform Witnesses Agreement
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"WitnessAgreement","Args":["AGR001","CAMP001","STARTUP001","INV001","27500","USD","9% equity stake - final offer","[{\"milestoneId\":\"MS001\",\"title\":\"Prototype Development\",\"description\":\"Complete working prototype\",\"targetAmount\":8250,\"targetDate\":\"2025-02-01\",\"status\":\"PENDING\",\"fundsReleased\":false,\"releasedAt\":\"\"},{\"milestoneId\":\"MS002\",\"title\":\"Beta Testing\",\"description\":\"Complete beta testing\",\"targetAmount\":11000,\"targetDate\":\"2025-02-28\",\"status\":\"PENDING\",\"fundsReleased\":false,\"releasedAt\":\"\"},{\"milestoneId\":\"MS003\",\"title\":\"Production Launch\",\"description\":\"Launch production\",\"targetAmount\":8250,\"targetDate\":\"2025-03-31\",\"status\":\"PENDING\",\"fundsReleased\":false,\"releasedAt\":\"\"}]"]}'
```

### 9.2 QUERY: Verify Agreement Created (Platform)
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"GetAgreement","Args":["AGR001"]}'
```

### 9.3 QUERY: Get All Agreements for Campaign
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"GetAgreementsByCampaign","Args":["CAMP001"]}'
```

**Switch to ValidatorOrg**

### 9.4 INVOKE: Validator Witnesses Agreement
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n validator -c '{"function":"WitnessAgreement","Args":["WITNESS001","AGR001","CAMP001","STARTUP001","INV001","27500","Agreement terms verified and compliant"]}'
```

### 9.5 QUERY: Verify Validator Witness Record
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n validator -c '{"function":"GetAgreementWitness","Args":["AGR001"]}'
```

---

# 💰 PHASE 10: FUNDING CONFIRMATION
**Organization: InvestorOrg | Channel: investor-platform-channel**

### 10.1 INVOKE: Confirm Funding Commitment
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID investor-platform-channel -n investor -c '{"function":"ConfirmFundingCommitment","Args":["COMMIT001","PROP001","AGR001","CAMP001","STARTUP001","INV001","27500","USD","[{\"milestoneId\":\"MS001\",\"title\":\"Prototype Development\",\"description\":\"Complete working prototype\",\"targetDate\":\"2025-02-01\",\"fundPercentage\":30,\"status\":\"PENDING\",\"fundsReleased\":false,\"releasedAt\":\"\"},{\"milestoneId\":\"MS002\",\"title\":\"Beta Testing\",\"description\":\"Complete beta testing\",\"targetDate\":\"2025-02-28\",\"fundPercentage\":40,\"status\":\"PENDING\",\"fundsReleased\":false,\"releasedAt\":\"\"},{\"milestoneId\":\"MS003\",\"title\":\"Production Launch\",\"description\":\"Launch production\",\"targetDate\":\"2025-03-31\",\"fundPercentage\":30,\"status\":\"PENDING\",\"fundsReleased\":false,\"releasedAt\":\"\"}]"]}'
```

### 10.2 INVOKE: Confirm Investment to Platform
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID investor-platform-channel -n investor -c '{"function":"ConfirmInvestmentToPlatform","Args":["CONFIRM001","INV_001","CAMP001","INV001","27500","USD"]}'
```

### 10.3 QUERY: Verify Investment Status Updated
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID platform-investor-channel -n investor -c '{"function":"GetInvestment","Args":["INV_001"]}'
```

**Switch to PlatformOrg**

### 10.4 INVOKE: Platform Records Investor Confirmation
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID investor-platform-channel -n platform -c '{"function":"RecordInvestorConfirmation","Args":["PLAT_REC001","CONFIRM001","CAMP001","INV001","27500","USD"]}'
```

---

# 📨 PHASE 11: STARTUP ACKNOWLEDGES INVESTMENT
**Organization: StartupOrg | Channel: common-channel**

### 11.1 INVOKE: Acknowledge Investment Receipt
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n startup -c '{"function":"AcknowledgeInvestment","Args":["INV_001","CAMP001","INV001","27500"]}'
```

### 11.2 QUERY: Verify Campaign from StartupOrg
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n startup -c '{"function":"GetCampaign","Args":["CAMP001"]}'
```

---

# 🏆 PHASE 12: MILESTONE COMPLETION & VERIFICATION
**Organization: StartupOrg | Channel: startup-validator-channel**

### 12.1 INVOKE: Submit Milestone Report
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n startup -c '{"function":"SubmitMilestoneReport","Args":["REPORT_MS001","CAMP001","MS001","AGR001","Prototype Development Complete","Successfully completed working prototype with all core features","[\"prototype_demo.mp4\",\"test_results.pdf\",\"code_review.pdf\"]"]}'
```

### 12.2 QUERY: Verify Milestone Report Created
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n startup -c '{"function":"GetMilestoneReport","Args":["REPORT_MS001"]}'
```

**Switch to ValidatorOrg**

### 12.3 INVOKE: Verify Milestone Completion
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n validator -c '{"function":"VerifyMilestoneCompletion","Args":["VERIFY_MS001","MS001","CAMP001","STARTUP001","milestone_report_hash_123","true","8.5","Milestone verified - prototype meets all requirements","true"]}'
```

### 12.4 QUERY: Verify Milestone Verification Record
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n validator -c '{"function":"GetMilestoneVerification","Args":["VERIFY_MS001"]}'
```

### 12.5 QUERY: Get Verification by Milestone ID
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n validator -c '{"function":"GetMilestoneVerificationByMilestone","Args":["MS001"]}'
```

---

# 💸 PHASE 13: FUND RELEASE
**Organization: InvestorOrg | Channel: common-channel**

### 13.1 INVOKE: Investor Approves Milestone
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n investor -c '{"function":"VerifyMilestone","Args":["INV_VERIFY_MS001","MS001","AGR001","CAMP001","INV001","true","Excellent work on the prototype. All features working as expected."]}'
```

**Switch to PlatformOrg**

### 13.2 INVOKE: Platform Triggers Fund Release
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"TriggerFundRelease","Args":["RELEASE001","ESCROW001","AGR001","CAMP001","MS001","STARTUP001","8250","USD","Milestone 1 verified by validator and investor"]}'
```

### 13.3 QUERY: Verify Fund Release Record
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"GetFundRelease","Args":["RELEASE001"]}'
```

### 13.4 QUERY: Verify Escrow Status
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"GetEscrow","Args":["ESCROW001"]}'
```

**Switch to StartupOrg**

### 13.5 INVOKE: Startup Receives Funding
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n startup -c '{"function":"ReceiveFunding","Args":["RELEASE001","CAMP001","MS001","8250","USD"]}'
```

### 13.6 QUERY: Verify Campaign Status Updated
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"GetPublishedCampaign","Args":["CAMP001"]}'
```

---

# 🎉 PHASE 14: CAMPAIGN COMPLETION
**Organization: StartupOrg | Channel: startup-validator-channel**

### 14.1 INVOKE: Mark Campaign Completed
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n startup -c '{"function":"MarkCampaignCompleted","Args":["CAMP001","50000","5"]}'
```

### 14.2 QUERY: Verify Campaign Status
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n startup -c '{"function":"GetCampaign","Args":["CAMP001"]}'
```

**Switch to ValidatorOrg**

### 14.3 INVOKE: Validator Confirms Completion
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n validator -c '{"function":"ConfirmCampaignCompletion","Args":["COMPLETION001","CAMP001","VAL001","true","All milestones verified and completed successfully"]}'
```

### 14.4 QUERY: Verify Campaign Completion Record
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n validator -c '{"function":"GetCampaignCompletion","Args":["CAMP001"]}'
```

**Switch to PlatformOrg**

### 14.5 INVOKE: Platform Closes Campaign
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"CloseCampaign","Args":["CLOSE001","CAMP001","SUCCESSFUL","50000","5","All milestones completed successfully"]}'
```

### 14.6 QUERY: Verify Final Campaign Status
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"GetPublishedCampaign","Args":["CAMP001"]}'
```

---

# 📢 PHASE 15: COMMON CHANNEL PUBLICATIONS
**Privacy-preserving summaries visible to all organizations**

**StartupOrg:**
### 15.1 INVOKE: Publish Campaign Summary Hash
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n startup -c '{"function":"PublishSummaryHash","Args":["SUMMARY001","CAMP001","Technology"]}'
```

**ValidatorOrg:**
### 15.2 INVOKE: Publish Validation Proof
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n validator -c '{"function":"PublishValidationProof","Args":["PROOF001","CAMP001","VAL001","APPROVED"]}'
```

**InvestorOrg:**
### 15.3 INVOKE: Publish Investment Summary
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n investor -c '{"function":"PublishInvestmentSummary","Args":["INV_SUMMARY001","CAMP001","5"]}'
```

**PlatformOrg:**
### 15.4 INVOKE: Publish Global Metrics
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"PublishGlobalMetrics","Args":["METRICS001","10","5","8","25"]}'
```

### 15.5 QUERY: Get Latest Global Metrics
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"GetLatestGlobalMetrics","Args":[]}'
```

---

# 📊 SUMMARY QUERIES (RUN ANYTIME)

## StartupOrg Queries
```bash
# Get campaign
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n startup -c '{"function":"GetCampaign","Args":["CAMP001"]}'

# Get campaigns by category
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n startup -c '{"function":"GetCampaignsByCategory","Args":["Technology"]}'

# Get campaigns by startup
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n startup -c '{"function":"GetCampaignsByStartup","Args":["STARTUP001"]}'

# Get document history
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n startup -c '{"function":"GetCampaignDocumentHistory","Args":["CAMP001"]}'

# Get agreement
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-investor-channel -n startup -c '{"function":"GetAgreement","Args":["AGR001"]}'

# Get milestone report
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n startup -c '{"function":"GetMilestoneReport","Args":["REPORT_MS001"]}'

# Get proposals
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-investor-channel -n startup -c '{"function":"GetProposal","Args":["PROP001"]}'
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-investor-channel -n startup -c '{"function":"GetProposalsByStartup","Args":["STARTUP001"]}'
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-investor-channel -n startup -c '{"function":"GetProposalsByCampaign","Args":["CAMP001"]}'
```

## ValidatorOrg Queries
```bash
# Get validation
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n validator -c '{"function":"GetValidation","Args":["VAL001"]}'
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n validator -c '{"function":"GetValidationsByCampaign","Args":["CAMP001"]}'

# Get validation report
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID validator-platform-channel -n validator -c '{"function":"GetValidationReport","Args":["CAMP001"]}'

# Get risk insight
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID investor-validator-channel -n validator -c '{"function":"GetRiskInsight","Args":["CAMP001"]}'

# Get milestone verification
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n validator -c '{"function":"GetMilestoneVerification","Args":["VERIFY_MS001"]}'
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n validator -c '{"function":"GetMilestoneVerificationByMilestone","Args":["MS001"]}'

# Get agreement witness
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n validator -c '{"function":"GetAgreementWitness","Args":["AGR001"]}'

# Get campaign completion
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n validator -c '{"function":"GetCampaignCompletion","Args":["CAMP001"]}'

# Verify hash / blacklist
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n validator -c '{"function":"VerifyCampaignHash","Args":["CAMP001","<HASH>"]}'
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n validator -c '{"function":"IsCampaignBlacklisted","Args":["CAMP001"]}'
```

## InvestorOrg Queries
```bash
# Get investment
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID platform-investor-channel -n investor -c '{"function":"GetInvestment","Args":["INV_001"]}'
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID platform-investor-channel -n investor -c '{"function":"GetInvestmentsByInvestor","Args":["INV001"]}'
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID platform-investor-channel -n investor -c '{"function":"GetInvestmentsByCampaign","Args":["CAMP001"]}'

# Get proposals
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-investor-channel -n investor -c '{"function":"GetProposal","Args":["PROP001"]}'
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-investor-channel -n investor -c '{"function":"GetProposalsByInvestor","Args":["INV001"]}'
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-investor-channel -n investor -c '{"function":"GetProposalsByCampaign","Args":["CAMP001"]}'
```

## PlatformOrg Queries
```bash
# Get published campaign
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"GetPublishedCampaign","Args":["CAMP001"]}'
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"GetActiveCampaigns","Args":[]}'

# Get validator decision
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID validator-platform-channel -n platform -c '{"function":"GetValidatorDecision","Args":["CAMP001"]}'

# Get agreement
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"GetAgreement","Args":["AGR001"]}'
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"GetAgreementsByCampaign","Args":["CAMP001"]}'

# Get fund release / escrow
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"GetFundRelease","Args":["RELEASE001"]}'
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"GetEscrow","Args":["ESCROW001"]}'

# Get global metrics
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"GetLatestGlobalMetrics","Args":[]}'
```

---

# 🔄 QUICK REFERENCE: ORG CONTEXT SWITCHING

```bash
# StartupOrg
export CORE_PEER_LOCALMSPID=StartupOrgMSP && export CORE_PEER_MSPCONFIGPATH=/home/kajal/crowdfunding/_msp/StartupOrg/startuporgadmin/msp && export CORE_PEER_ADDRESS=startuporgpeer-api.127-0-0-1.nip.io:9090

# ValidatorOrg
export CORE_PEER_LOCALMSPID=ValidatorOrgMSP && export CORE_PEER_MSPCONFIGPATH=/home/kajal/crowdfunding/_msp/ValidatorOrg/validatororgadmin/msp && export CORE_PEER_ADDRESS=validatororgpeer-api.127-0-0-1.nip.io:9090

# InvestorOrg
export CORE_PEER_LOCALMSPID=InvestorOrgMSP && export CORE_PEER_MSPCONFIGPATH=/home/kajal/crowdfunding/_msp/InvestorOrg/investororgadmin/msp && export CORE_PEER_ADDRESS=investororgpeer-api.127-0-0-1.nip.io:9090

# PlatformOrg
export CORE_PEER_LOCALMSPID=PlatformOrgMSP && export CORE_PEER_MSPCONFIGPATH=/home/kajal/crowdfunding/_msp/PlatformOrg/platformorgadmin/msp && export CORE_PEER_ADDRESS=platformorgpeer-api.127-0-0-1.nip.io:9090
```

---

# ⚖️ PHASE 16: DISPUTE SYSTEM WITH FILING FEES
**Anti-spam mechanism: Filing fees prevent frivolous disputes**

## 📋 Dispute Fee Tiers

| Tier | Name | Claim Amount Range | Filing Fee | Refund on Win |
|------|------|-------------------|------------|---------------|
| MINOR | Minor Dispute | < $500 | 10 tokens | 100% |
| STANDARD | Standard Dispute | $500 - $5,000 | 25 tokens | 100% |
| MAJOR | Major Dispute | $5,000 - $50,000 | 50 tokens | 100% |
| CRITICAL | Critical Dispute | > $50,000 | 100 tokens | 100% |

## 16.1 Initialize Dispute Fee System (One-time Setup)
**Organization: PlatformOrg | Channel: common-channel**

### INVOKE: Initialize Dispute Fee Tiers
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"InitializeDisputeFeeTiers","Args":[]}'
```

### QUERY: Get All Dispute Fee Tiers
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"GetDisputeFeeTiers","Args":[]}'
```

## 16.2 Wallet Setup (Required Before Filing Disputes)
**Organization: PlatformOrg | Channel: common-channel**

### INVOKE: Create Wallet for Investor
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"CreateWallet","Args":["WALLET_INV001","INV001","INVESTOR"]}'
```

### INVOKE: Create Wallet for Startup
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"CreateWallet","Args":["WALLET_STARTUP001","STARTUP001","STARTUP"]}'
```

### INVOKE: Deposit Tokens to Wallet
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"DepositTokens","Args":["WALLET_INV001","500","Initial deposit for platform participation"]}'
```

### QUERY: Check Wallet Balance
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"GetWallet","Args":["WALLET_INV001"]}'
```

## 16.3 Filing a Dispute (Fee Auto-Collected)
**Organization: InvestorOrg or StartupOrg | Channel: common-channel**

> 💡 **How it works:** When you create a dispute, the filing fee is automatically deducted from your wallet.
> - If you WIN the dispute → Fee is refunded 100%
> - If you LOSE or dismiss → Fee is forfeited

### INVOKE: Create Dispute (Investor vs Startup - Missed Milestone)
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"CreateDispute","Args":["DISP001","MILESTONE_FAILURE","MISSED_DEADLINE","INV001","INVESTOR","STARTUP001","STARTUP","CAMP001","AGR001","Missed Milestone Deadline","Startup failed to deliver Milestone MS002 by the agreed deadline","5000","[]"]}'
```

**Expected Response:**
```json
{
  "message": "Dispute created successfully",
  "disputeId": "DISP001",
  "ticketNumber": "DISP-123456",
  "status": "OPEN",
  "filingFee": 25,
  "feeNote": "Filing fee is locked. Will be refunded if dispute resolved in your favor, forfeited otherwise."
}
```

### QUERY: Get Dispute Details
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"GetDispute","Args":["DISP001"]}'
```

### QUERY: Get Dispute Fee Record
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"GetDisputeFeeByDisputeID","Args":["DISP001"]}'
```

## 16.4 Dispute Investigation (ValidatorOrg)
**Organization: ValidatorOrg | Channel: common-channel**

### INVOKE: Accept Dispute Investigation
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n validator -c '{"function":"AcceptDisputeInvestigation","Args":["INVEST001","DISP001","VALIDATOR001"]}'
```

### INVOKE: Record Investigation Finding
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n validator -c '{"function":"RecordInvestigationFinding","Args":["FINDING001","INVEST001","DISP001","Startup provided evidence of technical issues causing delay. Documentation shows 2-week delay but good faith effort.","[\"delay_evidence.pdf\",\"communication_logs.pdf\"]","PARTIAL_VIOLATION"]}'
```

### INVOKE: Complete Investigation
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n validator -c '{"function":"CompleteInvestigation","Args":["INVEST001","DISP001","PARTIAL","Both parties share responsibility. Recommend partial refund.","[\"Both parties should improve communication\",\"Milestone deadline was ambitious\"]"]}'
```

## 16.5 Anonymous Voting on Dispute (If Enabled)
**Organization: PlatformOrg | Channel: common-channel**

### INVOKE: Enable Voting for Dispute
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"EnableDisputeVoting","Args":["DISP001","2025-12-15T00:00:00Z","2025-12-20T00:00:00Z","5"]}'
```

### INVOKE: Commit Vote (Phase 1 - Anonymous)
```bash
# Voter commits a hash of their vote (vote + secret)
# Example: SHA256("FAVOR_INITIATOR" + "mysecret123") = <VOTE_HASH>
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"CommitVote","Args":["VOTE001","DISP001","VALIDATOR002","<VOTE_HASH>"]}'
```

### INVOKE: Reveal Vote (Phase 2 - After Commit Period)
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"RevealVote","Args":["VOTE001","DISP001","VALIDATOR002","FAVOR_INITIATOR","mysecret123"]}'
```

### INVOKE: Tally Votes
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"TallyVotes","Args":["DISP001"]}'
```

## 16.6 Resolve Dispute (Fee Refund/Forfeit Auto-Processed)
**Organization: PlatformOrg | Channel: common-channel**

### INVOKE: Resolve Dispute - Favor Initiator (Fee Refunded)
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"ResolveDispute","Args":["DISP001","FAVOR_INITIATOR","Investigation confirmed startup missed deadline without valid excuse.","[{\"penaltyId\":\"PEN001\",\"userId\":\"STARTUP001\",\"userType\":\"STARTUP\",\"tokenAmount\":500,\"reputationDeduct\":10,\"description\":\"Missed milestone deadline\"}]","[{\"refundOrderId\":\"REF001\",\"fromUserId\":\"STARTUP001\",\"fromUserType\":\"STARTUP\",\"toUserId\":\"INV001\",\"toUserType\":\"INVESTOR\",\"amount\":2500,\"deductionPercent\":0,\"reason\":\"Partial refund for missed milestone\"}]"]}'
```

### INVOKE: Resolve Dispute - Favor Respondent (Fee Forfeited)
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"ResolveDispute","Args":["DISP001","FAVOR_RESPONDENT","Investigation found dispute claim was invalid.","[]","[]"]}'
```

### QUERY: Verify Dispute Resolution
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"GetDispute","Args":["DISP001"]}'
```

### QUERY: Check Fee Outcome
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"GetDisputeFeeByDisputeID","Args":["DISP001"]}'
```

## 16.7 Query User's Dispute Fee History
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"GetUserDisputeFees","Args":["INV001"]}'
```

**Expected Response:**
```json
{
  "fees": [...],
  "totalCount": 3,
  "totalPaid": 75,
  "totalRefunded": 25,
  "totalForfeited": 50
}
```

---

# 🌟 PHASE 17: REPUTATION & RATING SYSTEM
**Organization: PlatformOrg | Channel: common-channel**

## 17.1 ML Rating System (0-100 Scale)

### INVOKE: Record Rating (ML Model Output)
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"RecordRating","Args":["RATING001","STARTUP001","STARTUP","CREDIBILITY","85","ML_MODEL_V1","Based on document verification and historical data"]}'
```

### QUERY: Get Rating Aggregate
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"GetRatingAggregate","Args":["STARTUP001","STARTUP"]}'
```

## 17.2 Reputation Management

### QUERY: Get User Reputation
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"GetReputation","Args":["STARTUP001","STARTUP"]}'
```

### QUERY: Get Reputation History
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"GetReputationHistory","Args":["STARTUP001"]}'
```

## 17.3 Auto-Suspension Rules

| Condition | Action |
|-----------|--------|
| Reputation Score ≤ 30 | Auto-Suspended for 30 days |
| 3 Consecutive Penalties | Auto-Suspended for 30 days |
| Reputation Score ≤ 15 | Permanently Blacklisted |

---

# 💳 PHASE 18: CAMPAIGN FEE SYSTEM
**Organization: PlatformOrg | Channel: common-channel**

## Campaign Fee Tiers

| Tier | Goal Amount Range | Platform Fee |
|------|------------------|--------------|
| SMALL | < $10,000 | 50 tokens |
| MEDIUM | $10,000 - $50,000 | 150 tokens |
| LARGE | $50,000 - $200,000 | 400 tokens |
| ENTERPRISE | > $200,000 | 1000 tokens |

### INVOKE: Initialize Campaign Fee Tiers
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"InitializeFeeTiers","Args":[]}'
```

### INVOKE: Collect Campaign Fee (Before Publication)
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"CollectCampaignFee","Args":["CFEE001","CAMP001","STARTUP001","WALLET_STARTUP001","50000"]}'
```

### QUERY: Get Campaign Fee Tiers
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"GetFeeTiers","Args":[]}'
```

---

## 📝 NOTES

1. **Replace `<HASH_FROM_STEP_1.6>`** with the actual hash returned from `GetCampaignValidationHash`
2. **Execute commands in order** - each phase depends on the previous
3. **Switch org context** before running commands for that organization
4. **Query after each invoke** to verify the operation succeeded
5. **IDs must be unique** - change CAMP001, INV001, etc. for new test runs
6. **Dispute Filing Fees** - Ensure wallet has sufficient tokens before creating disputes
7. **Anonymous Voting** - Uses commit-reveal scheme: commit hash first, reveal vote later
8. **Auto-Suspension** - System automatically suspends users with low reputation scores

---

## 🔧 TROUBLESHOOTING

### Common Issues

| Error | Cause | Solution |
|-------|-------|----------|
| `500 Internal Server Error` on validator query | Chaincode trying to access data from different org's namespace | Use the source org's chaincode for queries (e.g., query `-n startup` even from ValidatorOrg context) |
| `insufficient balance for dispute filing fee` | Wallet doesn't have enough tokens | Deposit tokens using `DepositTokens` function |
| `user is suspended` | Reputation score fell below threshold | Wait for suspension period or appeal |
| `user is blacklisted` | Permanent ban due to severe violations | No automated recovery |
| `dispute fee already processed` | Trying to process fee twice | Fee is auto-processed on resolution |
| `cross-channel query failed` | Wrong chaincode name in InvokeChaincode | Ensure using deployed names: `startup`, `validator`, `platform`, `investor` |

### Important: Chaincode Names vs Package Labels

When upgrading chaincodes, remember:
- **Package Label**: `startup_1`, `startup_2`, etc. (changes with each version)
- **Chaincode Name**: `startup` (stays constant, used in `-n` flag and `InvokeChaincode`)

Cross-chaincode invocations always use the **chaincode name**, not the label:
```go
// Correct - uses deployed chaincode name
ctx.GetStub().InvokeChaincode("startup", args, "startup-validator-channel")

// Wrong - uses file name or label
ctx.GetStub().InvokeChaincode("startuporg", args, "startup-validator-channel")  // ❌
ctx.GetStub().InvokeChaincode("startup_2", args, "startup-validator-channel")   // ❌
```
