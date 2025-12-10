# Complete Crowdfunding Platform Workflow Commands

## 🎯 CASE SCENARIO: "Smart Home IoT Platform" Campaign

This demonstrates the complete lifecycle of a crowdfunding campaign from creation to fund release.

---

## 📋 WORKFLOW PHASES OVERVIEW

| Phase | Description | Channel |
|-------|-------------|---------|
| 1 | Campaign Creation & Query | startup-validator-channel |
| 2 | Validation (APPROVED/ON_HOLD/REJECTED scenarios) | startup-validator-channel |
| 3 | Validator Report to Platform | validator-platform-channel |
| 4 | Submit for Publishing | startup-platform-channel |
| 5 | Platform Publishes Campaign | common-channel |
| 6 | Investor Views & Invests | platform-investor-channel |
| 7 | Risk Insights Request | investor-validator-channel |
| 8 | Investment Proposal & Negotiation | startup-investor-channel |
| 9 | Platform and Validator Witnesses Agreement | common-channel |
| 10 | Funding Confirmation | investor-platform-channel |
| 11 | Startup Acknowledges Investment (Platform forwards after receiving funds) | common-channel |
| 12 | Milestone Completion Verification by Validator | startup-validator-channel |
| 13 | Milestone Completion & Fund Release | common-channel |
| 14 | Campaign Completion | startup-validator-channel & common-channel |
| 15 | Common Channel Publications | common-channel |

---

## 📋 PHASE 1: Campaign Creation & Query (startup-validator-channel)

### Step 1.1: Startup Creates Campaign (DRAFT)
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n startup -c '{"function":"CreateCampaign","Args":["CAMP001","STARTUP001","Technology","2025-03-31","USD","false","false","2025-01-01","Prototype","Hardware","[\"IoT\",\"SmartHome\",\"AI\"]","false","false","90","1","1","2025","50000","50K-100K","Smart Home IoT Platform","An innovative IoT platform for smart home automation with AI-powered features","[\"business_plan.pdf\",\"pitch_deck.pdf\",\"financials.xlsx\"]"]}'
```

### Step 1.2: Query Campaign (Startup can view at any time)
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n startup -c '{"function":"GetCampaign","Args":["CAMP001"]}'
```

### Step 1.3: Startup Submits for Validation
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n startup -c '{"function":"SubmitForValidation","Args":["CAMP001","Please review our IoT platform proposal"]}'
```

### Step 1.4: Query Campaign After Submission
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n startup -c '{"function":"GetCampaign","Args":["CAMP001"]}'
```

---

## 📋 PHASE 2: Validator Validates Campaign (startup-validator-channel)

### SCENARIO A: Validator APPROVES Campaign

### Step 2.1A: Validator Validates Campaign (APPROVED)
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n validator -c '{"function":"ValidateCampaign","Args":["VAL001","CAMP001","a4b9cf29a14cda330a06f67bdb4abfe4aa1ecf2e4d1512d5ee466d66cad41e9d","VALIDATOR001","true","true","8.5","2.5","APPROVED","[\"Documents verified\",\"Team credentials confirmed\"]",""]}'
```

### Step 2.2A: Query Validation Record
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n validator -c '{"function":"GetValidation","Args":["VAL001"]}'
```

---

### SCENARIO B: Validator puts Campaign ON_HOLD (needs more docs)

### Step 2.1B: Validator Validates Campaign (ON_HOLD)
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n validator -c '{"function":"ValidateCampaign","Args":["VAL001","CAMP001","a4b9cf29a14cda330a06f67bdb4abfe4aa1ecf2e4d1512d5ee466d66cad41e9d","VALIDATOR001","false","true","5.0","4.5","ON_HOLD","[\"Missing financial projections\",\"Need team credentials\"]","team_credentials.pdf,financial_projections.xlsx"]}'
```

### Step 2.2B: Startup Updates Documents After ON_HOLD
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n startup -c '{"function":"UpdateCampaignDocs","Args":["CAMP001","[\"business_plan.pdf\",\"pitch_deck.pdf\",\"financials.xlsx\",\"team_credentials.pdf\",\"financial_projections.xlsx\"]","Added requested documents: team credentials and financial projections"]}'
```

### Step 2.3B: Validator Re-validates After Document Update (APPROVED)
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n validator -c '{"function":"ValidateCampaign","Args":["VAL001","CAMP001","a4b9cf29a14cda330a06f67bdb4abfe4aa1ecf2e4d1512d5ee466d66cad41e9d","VALIDATOR001","true","true","8.5","2.0","APPROVED","[\"All documents verified\",\"Financial projections look solid\"]",""]}'
```

---

### SCENARIO C: Validator REJECTS Campaign (Fraudulent - Blacklisted)

### Step 2.1C: Validator Validates Campaign (REJECTED - Fraud Detected)
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n validator -c '{"function":"ValidateCampaign","Args":["VAL002","CAMP_FRAUD","fake_hash","VALIDATOR001","false","false","1.0","9.5","REJECTED","[\"Fraudulent documents detected\",\"Identity verification failed\"]",""]}'
```

### Step 2.2C: Check if Campaign is Blacklisted
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n validator -c '{"function":"IsCampaignBlacklisted","Args":["CAMP_FRAUD"]}'
```

---

## 📋 PHASE 3: Validator Sends Report to Platform (validator-platform-channel)

### Step 3.1: Validator Sends Validation Report to Platform
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID validator-platform-channel -n validator -c '{"function":"SendValidationReportToPlatform","Args":["REPORT001","CAMP001","VAL001","a4b9cf29a14cda330a06f67bdb4abfe4aa1ecf2e4d1512d5ee466d66cad41e9d","8.5","9.0","8.0","2.5","true","Campaign fully verified. Low risk. Recommended for publication."]}'
```

### Step 3.2: Query Validation Report
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID validator-platform-channel -n validator -c '{"function":"GetValidationReport","Args":["CAMP001"]}'
```

---

## 📋 PHASE 4: Startup Submits to Platform for Publishing (startup-platform-channel)

### Step 4.1: Startup Submits Campaign for Publishing
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-platform-channel -n startup -c '{"function":"SubmitForPublishing","Args":["CAMP001"]}'
```

---

## 📋 PHASE 5: Platform Publishes Campaign (common-channel)

### Step 5.1: Platform Publishes Campaign to Portal (visible to all orgs)
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"PublishCampaignToPortal","Args":["CAMP001","STARTUP001","Smart Home IoT Platform","Technology","An innovative IoT platform for smart home automation","50000","USD","2025-01-01","2025-03-31","90","8.5","a4b9cf29a14cda330a06f67bdb4abfe4aa1ecf2e4d1512d5ee466d66cad41e9d","[{\"milestoneId\":\"MS001\",\"title\":\"Prototype Development\",\"description\":\"Complete working prototype\",\"targetAmount\":15000,\"targetDate\":\"2025-02-01\",\"status\":\"PENDING\",\"fundsReleased\":false,\"releasedAt\":\"\"},{\"milestoneId\":\"MS002\",\"title\":\"Beta Testing\",\"description\":\"Complete beta testing phase\",\"targetAmount\":20000,\"targetDate\":\"2025-02-28\",\"status\":\"PENDING\",\"fundsReleased\":false,\"releasedAt\":\"\"},{\"milestoneId\":\"MS003\",\"title\":\"Production Launch\",\"description\":\"Launch production version\",\"targetAmount\":15000,\"targetDate\":\"2025-03-31\",\"status\":\"PENDING\",\"fundsReleased\":false,\"releasedAt\":\"\"}]"]}'
```

### Step 5.2: Query Published Campaign (All orgs can query)
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"GetPublishedCampaign","Args":["CAMP001"]}'
```

### Step 5.3: Platform Records Validator Decision
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID validator-platform-channel -n platform -c '{"function":"RecordValidatorDecision","Args":["REC001","CAMP001","VAL001","a4b9cf29a14cda330a06f67bdb4abfe4aa1ecf2e4d1512d5ee466d66cad41e9d","true","8.5","report_hash_here"]}'
```

---

## 📋 PHASE 6: Investor Views & Invests (platform-investor-channel)

### Step 6.1: Investor Views Campaign Details
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID platform-investor-channel -n investor -c '{"function":"ViewCampaign","Args":["CAMP001","INV001","Smart Home IoT Platform","Technology","An innovative IoT platform","50000","0","USD","2025-01-01","2025-03-31","Prototype","Hardware","[\"IoT\",\"SmartHome\",\"AI\"]","90","8.5","LOW","0","PUBLISHED"]}'
```

### Step 6.2: Investor Makes Investment Commitment
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID platform-investor-channel -n investor -c '{"function":"MakeInvestment","Args":["INV_001","CAMP001","INV001","10000","USD"]}'
```

### Step 6.3: Query Investment
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID platform-investor-channel -n investor -c '{"function":"GetInvestment","Args":["INV_001"]}'
```

### Step 6.4: Query All Investments by Investor
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID platform-investor-channel -n investor -c '{"function":"GetInvestmentsByInvestor","Args":["INV001"]}'
```

### Step 6.5: Query All Investments for Campaign
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID platform-investor-channel -n investor -c '{"function":"GetInvestmentsByCampaign","Args":["CAMP001"]}'
```

---

## 📋 PHASE 7: Investor Requests Risk Insights (investor-validator-channel)

### Step 7.1: Investor Requests Risk Insights
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID investor-validator-channel -n investor -c '{"function":"RequestRiskInsights","Args":["RISK_REQ001","CAMP001","INV001"]}'
```

### Step 7.2: Validator Assigns Risk Score to Investor
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID investor-validator-channel -n validator -c '{"function":"AssignRiskScore","Args":["INSIGHT001","CAMP001","INV001","2.5","[\"Strong team\",\"Good market potential\",\"Early stage product\"]","What are the main risks?","The main risks are market competition and execution timeline. Team has strong track record.","RECOMMENDED - Low risk investment with good potential"]}'
```

### Step 7.3: Query Risk Insight
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID investor-validator-channel -n validator -c '{"function":"GetRiskInsight","Args":["CAMP001"]}'
```

---

## 📋 PHASE 8: Investment Proposal & Negotiation (startup-investor-channel)

### Step 8.1: Investor Creates Investment Proposal
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-investor-channel -n investor -c '{"function":"CreateInvestmentProposal","Args":["PROP001","CAMP001","STARTUP001","INV001","25000","USD","10% equity stake with board observer rights","[{\"milestoneId\":\"MS001\",\"title\":\"Prototype Development\",\"description\":\"Complete working prototype\",\"targetDate\":\"2025-02-01\",\"fundPercentage\":30,\"status\":\"PENDING\",\"fundsReleased\":false,\"releasedAt\":\"\"},{\"milestoneId\":\"MS002\",\"title\":\"Beta Testing\",\"description\":\"Complete beta testing\",\"targetDate\":\"2025-02-28\",\"fundPercentage\":40,\"status\":\"PENDING\",\"fundsReleased\":false,\"releasedAt\":\"\"},{\"milestoneId\":\"MS003\",\"title\":\"Production Launch\",\"description\":\"Launch production\",\"targetDate\":\"2025-03-31\",\"fundPercentage\":30,\"status\":\"PENDING\",\"fundsReleased\":false,\"releasedAt\":\"\"}]"]}'
```

### Step 8.2: Startup Responds with Counter Offer
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-investor-channel -n startup -c '{"function":"RespondToInvestmentProposal","Args":["PROP001","COUNTER","8% equity stake with quarterly updates","30000"]}'
```

### Step 8.3: Investor Responds to Counter Offer (ACCEPT)
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-investor-channel -n investor -c '{"function":"RespondToCounterOffer","Args":["PROP001","ACCEPT","9% equity stake - final offer","27500"]}'
```

### Step 8.4: Investor Formally Accepts Agreement
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-investor-channel -n investor -c '{"function":"AcceptAgreement","Args":["PROP001","AGR001"]}'
```

---

## 📋 PHASE 9: Platform and Validator Witnesses Agreement (common-channel)

### Step 9.1: Platform Witnesses Agreement (visible to all parties)
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"WitnessAgreement","Args":["AGR001","CAMP001","STARTUP001","INV001","27500","USD","9% equity stake - final offer","[{\"milestoneId\":\"MS001\",\"title\":\"Prototype Development\",\"description\":\"Complete working prototype\",\"targetAmount\":8250,\"targetDate\":\"2025-02-01\",\"status\":\"PENDING\",\"fundsReleased\":false,\"releasedAt\":\"\"},{\"milestoneId\":\"MS002\",\"title\":\"Beta Testing\",\"description\":\"Complete beta testing\",\"targetAmount\":11000,\"targetDate\":\"2025-02-28\",\"status\":\"PENDING\",\"fundsReleased\":false,\"releasedAt\":\"\"},{\"milestoneId\":\"MS003\",\"title\":\"Production Launch\",\"description\":\"Launch production\",\"targetAmount\":8250,\"targetDate\":\"2025-03-31\",\"status\":\"PENDING\",\"fundsReleased\":false,\"releasedAt\":\"\"}]"]}'
```

### Step 9.2: Validator Witnesses Agreement (adds validation attestation)
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n validator -c '{"function":"WitnessAgreement","Args":["WITNESS001","AGR001","CAMP001","STARTUP001","INV001","27500","Agreement terms verified and compliant"]}'
```

---

## 📋 PHASE 10: Funding Confirmation (investor-platform-channel)

### Step 10.1: Investor Confirms Funding Commitment
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID investor-platform-channel -n investor -c '{"function":"ConfirmFundingCommitment","Args":["COMMIT001","PROP001","AGR001","CAMP001","STARTUP001","INV001","27500","USD","[{\"milestoneId\":\"MS001\",\"title\":\"Prototype Development\",\"description\":\"Complete working prototype\",\"targetDate\":\"2025-02-01\",\"fundPercentage\":30,\"status\":\"PENDING\",\"fundsReleased\":false,\"releasedAt\":\"\"},{\"milestoneId\":\"MS002\",\"title\":\"Beta Testing\",\"description\":\"Complete beta testing\",\"targetDate\":\"2025-02-28\",\"fundPercentage\":40,\"status\":\"PENDING\",\"fundsReleased\":false,\"releasedAt\":\"\"},{\"milestoneId\":\"MS003\",\"title\":\"Production Launch\",\"description\":\"Launch production\",\"targetDate\":\"2025-03-31\",\"fundPercentage\":30,\"status\":\"PENDING\",\"fundsReleased\":false,\"releasedAt\":\"\"}]"]}'
```

### Step 10.2: Investor Confirms Investment to Platform
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID investor-platform-channel -n investor -c '{"function":"ConfirmInvestmentToPlatform","Args":["CONFIRM001","INV_001","CAMP001","INV001","27500","USD"]}'
```

### Step 10.3: Platform Records Investor Confirmation
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID investor-platform-channel -n platform -c '{"function":"RecordInvestorConfirmation","Args":["PLAT_REC001","CONFIRM001","CAMP001","INV001","27500","USD"]}'
```

---

## 📋 PHASE 11: Startup Acknowledges Investment via Platform (common-channel)

Platform receives funds from Investor and Startup acknowledges receipt on common-channel.

### Step 11.1: Startup Acknowledges Investment Receipt
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n startup -c '{"function":"AcknowledgeInvestment","Args":["INV_001","CAMP001","INV001","27500"]}'
```

---

## 📋 PHASE 12: Milestone Completion Verification by Validator (startup-validator-channel)

### Step 12.1: Startup Submits Milestone Report
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n startup -c '{"function":"SubmitMilestoneReport","Args":["REPORT_MS001","CAMP001","MS001","AGR001","Prototype Development Complete","Successfully completed working prototype with all core features","[\"prototype_demo.mp4\",\"test_results.pdf\",\"code_review.pdf\"]"]}'
```

### Step 12.2: Validator Verifies Milestone Completion
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n validator -c '{"function":"VerifyMilestoneCompletion","Args":["VERIFY_MS001","MS001","CAMP001","STARTUP001","milestone_report_hash_123","true","8.5","Milestone verified - prototype meets all requirements","true"]}'
```

### Step 12.3: Query Milestone Verification Status
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n validator -c '{"function":"GetValidation","Args":["VERIFY_MS001"]}'
```

---

## 📋 PHASE 13: Milestone Completion & Fund Release (common-channel)

### Step 13.1: Investor Approves Milestone (on common-channel for transparency)
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n investor -c '{"function":"VerifyMilestone","Args":["INV_VERIFY_MS001","MS001","AGR001","CAMP001","INV001","true","Excellent work on the prototype. All features working as expected."]}'
```

### Step 13.2: Platform Triggers Fund Release for Milestone 1
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"TriggerFundRelease","Args":["RELEASE001","ESCROW001","AGR001","CAMP001","MS001","STARTUP001","8250","USD","Milestone 1 verified by validator and investor"]}'
```

### Step 13.3: Startup Receives Funding
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n startup -c '{"function":"ReceiveFunding","Args":["RELEASE001","CAMP001","MS001","8250","USD"]}'
```

### Step 13.4: Query Updated Campaign Status
```bash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"GetPublishedCampaign","Args":["CAMP001"]}'
```

---

## 📋 PHASE 14: Campaign Completion (startup-validator-channel & common-channel)

### Step 14.1: Startup Marks Campaign Completed (startup-validator-channel)
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n startup -c '{"function":"MarkCampaignCompleted","Args":["CAMP001","50000","5"]}'
```

### Step 14.2: Validator Confirms Campaign Completion
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n validator -c '{"function":"ConfirmCampaignCompletion","Args":["CONFIRM001","CAMP001","VAL001","true","All milestones verified and completed successfully"]}'
```

### Step 14.3: Platform Closes Campaign (common-channel)
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"CloseCampaign","Args":["CLOSE001","CAMP001","SUCCESSFUL","50000","5","All milestones completed successfully"]}'
```

---

## 📋 PHASE 15: Common Channel Publications (common-channel)

Privacy-preserving summaries and proofs visible to all organizations.

### Step 15.1: Startup Publishes Campaign Summary Hash
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n startup -c '{"function":"PublishSummaryHash","Args":["SUMMARY001","CAMP001","Technology"]}'
```

### Step 15.2: Validator Publishes Validation Proof
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n validator -c '{"function":"PublishValidationProof","Args":["PROOF001","CAMP001","VAL001","APPROVED"]}'
```

### Step 15.3: Investor Publishes Investment Summary
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n investor -c '{"function":"PublishInvestmentSummary","Args":["INV_SUMMARY001","CAMP001","5"]}'
```

### Step 15.4: Platform Publishes Global Metrics
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"PublishGlobalMetrics","Args":["METRICS001","10","5","8","25"]}'
```

---

## 📋 ADDITIONAL QUERY COMMANDS

### Startup Queries
```bash
# Get campaigns by category
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n startup -c '{"function":"GetCampaignsByCategory","Args":["Technology"]}'

# Get campaigns by startup
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n startup -c '{"function":"GetCampaignsByStartup","Args":["STARTUP001"]}'

# Get campaign validation hash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n startup -c '{"function":"GetCampaignValidationHash","Args":["CAMP001"]}'

# Get campaign document history
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n startup -c '{"function":"GetCampaignDocumentHistory","Args":["CAMP001"]}'

# Get agreement
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-investor-channel -n startup -c '{"function":"GetAgreement","Args":["AGR001"]}'

# Get milestone report
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-investor-channel -n startup -c '{"function":"GetMilestoneReport","Args":["REPORT_MS001"]}'
```

### Platform Queries
```bash
# Get active campaigns (published campaigns are on common-channel)
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"GetActiveCampaigns","Args":[]}'

# Get validator decision
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID validator-platform-channel -n platform -c '{"function":"GetValidatorDecision","Args":["CAMP001"]}'

# Get latest global metrics
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID common-channel -n platform -c '{"function":"GetLatestGlobalMetrics","Args":[]}'
```

### Validator Queries
```bash
# Verify campaign hash
peer chaincode query -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-validator-channel -n validator -c '{"function":"VerifyCampaignHash","Args":["CAMP001","a4b9cf29a14cda330a06f67bdb4abfe4aa1ecf2e4d1512d5ee466d66cad41e9d"]}'
```

---

## 📋 INVESTOR WITHDRAWAL SCENARIO

### Investor Withdraws Investment (Before Campaign Closes)
```bash
peer chaincode invoke -o orderer-api.127-0-0-1.nip.io:9090 --channelID startup-investor-channel -n investor -c '{"function":"WithdrawInvestment","Args":["INV_001","Changed investment strategy"]}'
```

---

## 🔄 CHANNEL SUMMARY

| Channel | Organizations | Purpose |
|---------|---------------|---------|
| startup-validator-channel | StartupOrg, ValidatorOrg | Campaign creation, validation, milestone verification |
| startup-platform-channel | StartupOrg, PlatformOrg | Submit for publishing |
| startup-investor-channel | StartupOrg, InvestorOrg | Proposal negotiation |
| platform-investor-channel | PlatformOrg, InvestorOrg | Investor views campaigns, funding confirmation |
| investor-validator-channel | InvestorOrg, ValidatorOrg | Risk insights |
| validator-platform-channel | ValidatorOrg, PlatformOrg | Validation reports |
| **common-channel** | **All Organizations** | **Publishing, agreements, fund release, completion, public proofs** |

---

## 📝 NOTES

1. **Hash Values**: Replace `a4b9cf29a14cda330a06f67bdb4abfe4aa1ecf2e4d1512d5ee466d66cad41e9d` with actual hash returned from CreateCampaign
2. **IDs**: All IDs (CAMP001, INV001, etc.) should be unique per invocation
3. **Channel Context**: Make sure to switch peer context to the correct organization before invoking on their behalf
4. **Order of Operations**: Follow the phase sequence for proper workflow
5. **Common Channel**: Used for multi-party visibility operations (publishing, agreements, fund release, acknowledgements)
