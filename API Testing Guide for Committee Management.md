# REST API Endpoints for Crowdfunding Platform
## Hyperledger Fabric Chaincode Functions

This document provides all REST API endpoints for the crowdfunding platform based on the Hyperledger Fabric chaincode contracts.

Base API structure:
```
/api/startup/       - StartupOrg operations
/api/investor/      - InvestorOrg operations
/api/validator/     - ValidatorOrg operations
/api/platform/      - PlatformOrg operations
```

---

## 1. STARTUP ORG API

Base path: `/api/startup/`

### 1.1 Campaign Management

#### Create Campaign
- **POST** `/api/startup/campaigns`
- **Body:**
```json
{
  "startupId": "string",
  "projectName": "string",
  "description": "string",
  "category": "string",
  "goalAmount": 100000,
  "currency": "USD",
  "openDate": "2025-01-01",
  "closeDate": "2025-06-01",
  "durationDays": 150,
  "productStage": "string",
  "projectType": "string",
  "tags": ["tech", "innovation"]
}
```

#### Get Campaign
- **GET** `/api/startup/campaigns/{campaignId}`

#### Get Campaigns by Startup
- **GET** `/api/startup/campaigns?startupId={startupId}`

#### Get Campaigns by Category
- **GET** `/api/startup/campaigns?category={category}`

#### Update Campaign
- **PUT** `/api/startup/campaigns/{campaignId}`
- **Body:**
```json
{
  "fieldName": "description",
  "newValue": "Updated description",
  "updateReason": "Clarification needed"
}
```

#### Get Campaign Update History
- **GET** `/api/startup/campaigns/{campaignId}/update-history`

#### Get Campaign Document History
- **GET** `/api/startup/campaigns/{campaignId}/document-history`

#### Get Campaign Validation Hash
- **GET** `/api/startup/campaigns/{campaignId}/validation-hash`

### 1.2 Campaign Submission & Validation

#### Submit for Validation
- **POST** `/api/startup/campaigns/{campaignId}/submit-validation`
- **Body:**
```json
{
  "documents": ["ipfsHash1", "ipfsHash2"],
  "submissionNotes": "Initial submission"
}
```

#### Update Campaign Documents
- **PUT** `/api/startup/campaigns/{campaignId}/documents`
- **Body:**
```json
{
  "documents": ["newIpfsHash1", "newIpfsHash2"],
  "submissionNotes": "Revised documents as requested"
}
```

#### Submit for Publishing
- **POST** `/api/startup/campaigns/{campaignId}/submit-publishing`

#### Mark Campaign Completed
- **PUT** `/api/startup/campaigns/{campaignId}/mark-completed`

### 1.3 Investment & Proposals

#### Acknowledge Investment
- **POST** `/api/startup/investments/{investmentId}/acknowledge`
- **Body:**
```json
{
  "campaignId": "string",
  "investorId": "string"
}
```

#### Get Proposal
- **GET** `/api/startup/proposals/{proposalId}`

#### Get Proposals by Campaign
- **GET** `/api/startup/proposals?campaignId={campaignId}`

#### Get Proposals by Startup
- **GET** `/api/startup/proposals?startupId={startupId}`

#### Respond to Investment Proposal
- **POST** `/api/startup/proposals/{proposalId}/respond`
- **Body:**
```json
{
  "action": "ACCEPT|COUNTER|REJECT",
  "counterAmount": 150000,
  "counterTerms": "Revised terms...",
  "modifiedMilestones": []
}
```

### 1.4 Agreement Management

#### Get Agreement
- **GET** `/api/startup/agreements/{agreementId}`

### 1.5 Milestone Reporting

#### Submit Milestone Report
- **POST** `/api/startup/milestones/reports`
- **Body:**
```json
{
  "campaignId": "string",
  "milestoneId": "string",
  "agreementId": "string",
  "title": "Phase 1 Complete",
  "description": "All deliverables met",
  "evidence": ["ipfsHash1", "ipfsHash2"]
}
```

#### Get Milestone Report
- **GET** `/api/startup/milestones/reports/{reportId}`

#### Receive Funding
- **POST** `/api/startup/funding/receive`
- **Body:**
```json
{
  "releaseId": "string",
  "agreementId": "string",
  "milestoneId": "string",
  "amount": 50000
}
```

### 1.6 Disputes & Fees

#### Submit Dispute
- **POST** `/api/startup/disputes`
- **Body:**
```json
{
  "disputeType": "AGAINST_INVESTOR|AGAINST_VALIDATOR|AGAINST_PLATFORM",
  "targetId": "string",
  "targetType": "string",
  "campaignId": "string",
  "agreementId": "string",
  "title": "Dispute title",
  "description": "Dispute details",
  "claimedAmount": 10000,
  "evidenceHashes": ["ipfsHash1"]
}
```

#### Submit Dispute Evidence
- **POST** `/api/startup/disputes/{disputeId}/evidence`
- **Body:**
```json
{
  "evidenceHashes": ["ipfsHash1", "ipfsHash2"],
  "evidenceDescription": "Additional proof"
}
```

#### Respond to Dispute
- **POST** `/api/startup/disputes/{disputeId}/respond`
- **Body:**
```json
{
  "responseText": "Response to dispute",
  "counterEvidenceHashes": ["ipfsHash1"]
}
```

#### Get Startup Disputes
- **GET** `/api/startup/disputes?startupId={startupId}`

#### Record Fee Payment
- **POST** `/api/startup/fees/record`
- **Body:**
```json
{
  "campaignId": "string",
  "feeType": "CAMPAIGN_FEE|DISPUTE_FEE",
  "amount": 500,
  "transactionHash": "string"
}
```

### 1.7 Common Channel Operations

#### Publish Summary Hash
- **POST** `/api/startup/campaigns/{campaignId}/publish-summary`

---

## 2. INVESTOR ORG API

Base path: `/api/investor/`

### 2.1 Campaign Viewing

#### View Campaign
- **GET** `/api/investor/campaigns/{campaignId}/view`
- **Query Params:** `investorId={investorId}`

### 2.2 Investment Management

#### Make Investment
- **POST** `/api/investor/investments`
- **Body:**
```json
{
  "campaignId": "string",
  "investorId": "string",
  "amount": 5000,
  "currency": "USD"
}
```

#### Get Investment
- **GET** `/api/investor/investments/{investmentId}`

#### Get Investments by Investor
- **GET** `/api/investor/investments?investorId={investorId}`

#### Get Investments by Campaign
- **GET** `/api/investor/investments?campaignId={campaignId}`

#### Withdraw Investment
- **PUT** `/api/investor/investments/{investmentId}/withdraw`

### 2.3 Investment Proposals

#### Create Investment Proposal
- **POST** `/api/investor/proposals`
- **Body:**
```json
{
  "campaignId": "string",
  "startupId": "string",
  "investorId": "string",
  "investmentAmount": 100000,
  "currency": "USD",
  "proposedTerms": "Terms and conditions...",
  "milestones": [
    {
      "title": "Phase 1",
      "description": "Initial development",
      "targetDate": "2025-03-01",
      "fundPercentage": 30
    }
  ]
}
```

#### Get Proposal
- **GET** `/api/investor/proposals/{proposalId}`

#### Get Proposals by Campaign
- **GET** `/api/investor/proposals?campaignId={campaignId}`

#### Get Proposals by Investor
- **GET** `/api/investor/proposals?investorId={investorId}`

#### Respond to Counter Offer
- **POST** `/api/investor/proposals/{proposalId}/counter-response`
- **Body:**
```json
{
  "action": "ACCEPT|COUNTER|REJECT",
  "counterAmount": 120000,
  "counterTerms": "Modified terms..."
}
```

### 2.4 Agreement & Funding

#### Accept Agreement
- **POST** `/api/investor/agreements/{agreementId}/accept`

#### Confirm Funding Commitment
- **POST** `/api/investor/funding/confirm`
- **Body:**
```json
{
  "proposalId": "string",
  "agreementId": "string",
  "campaignId": "string",
  "startupId": "string",
  "amount": 100000,
  "milestones": []
}
```

#### Confirm Investment to Platform
- **POST** `/api/investor/platform/confirm-investment`
- **Body:**
```json
{
  "investmentId": "string",
  "campaignId": "string",
  "amount": 100000
}
```

### 2.5 Milestone Verification

#### Verify Milestone
- **POST** `/api/investor/milestones/verify`
- **Body:**
```json
{
  "milestoneId": "string",
  "agreementId": "string",
  "campaignId": "string",
  "approved": true,
  "feedback": "Milestone completed satisfactorily"
}
```

### 2.6 Risk Insights

#### Request Risk Insights
- **POST** `/api/investor/risk/request`
- **Body:**
```json
{
  "campaignId": "string",
  "investorId": "string"
}
```

#### Record Risk Insight Response
- **POST** `/api/investor/risk/record-response`
- **Body:**
```json
{
  "requestId": "string",
  "campaignId": "string",
  "riskScore": 7.5,
  "riskLevel": "MEDIUM",
  "riskFactors": "Market volatility, team experience",
  "recommendation": "Proceed with caution"
}
```

### 2.7 Disputes & Refunds

#### Submit Dispute
- **POST** `/api/investor/disputes`
- **Body:**
```json
{
  "disputeType": "AGAINST_STARTUP|AGAINST_VALIDATOR|AGAINST_PLATFORM",
  "targetId": "string",
  "targetType": "string",
  "campaignId": "string",
  "agreementId": "string",
  "title": "Dispute title",
  "description": "Dispute details",
  "claimedAmount": 10000,
  "evidenceHashes": ["ipfsHash1"]
}
```

#### Submit Dispute Evidence
- **POST** `/api/investor/disputes/{disputeId}/evidence`
- **Body:**
```json
{
  "evidenceHashes": ["ipfsHash1"],
  "evidenceDescription": "Supporting documents"
}
```

#### Respond to Dispute
- **POST** `/api/investor/disputes/{disputeId}/respond`
- **Body:**
```json
{
  "responseText": "Response details",
  "counterEvidenceHashes": []
}
```

#### Get Investor Disputes
- **GET** `/api/investor/disputes?investorId={investorId}`

#### Request Refund
- **POST** `/api/investor/refunds`
- **Body:**
```json
{
  "campaignId": "string",
  "agreementId": "string",
  "startupId": "string",
  "originalAmount": 100000,
  "requestedAmount": 85000,
  "refundReason": "EARLY_WITHDRAWAL|MID_AGREEMENT_WITHDRAWAL|DISPUTE_RESOLUTION|CAMPAIGN_FAILED",
  "deductionPercent": 15
}
```

#### Get Refund Request
- **GET** `/api/investor/refunds/{requestId}`

### 2.8 Common Channel Operations

#### Publish Investment Summary
- **POST** `/api/investor/campaigns/{campaignId}/publish-summary`

#### Receive Campaign Notification
- **POST** `/api/investor/notifications/campaign`

#### Receive Risk Insight
- **POST** `/api/investor/risk/receive`

---

## 3. VALIDATOR ORG API

Base path: `/api/validator/`

### 3.1 Campaign Validation

#### Validate Campaign
- **POST** `/api/validator/campaigns/{campaignId}/validate`
- **Body:**
```json
{
  "validatorId": "string",
  "campaignHash": "string",
  "documentsReviewed": ["doc1", "doc2"]
}
```

#### Approve or Reject Campaign
- **PUT** `/api/validator/campaigns/{campaignId}/decision`
- **Body:**
```json
{
  "validationId": "string",
  "status": "APPROVED|ON_HOLD|REJECTED|BLACKLISTED",
  "dueDiligenceScore": 8.5,
  "riskScore": 6.2,
  "riskLevel": "MEDIUM",
  "comments": "Review complete",
  "issues": [],
  "requiredDocuments": ""
}
```

#### Verify Campaign Hash
- **GET** `/api/validator/campaigns/{campaignId}/verify-hash`
- **Query Params:** `campaignHash={hash}`

#### Is Campaign Blacklisted
- **GET** `/api/validator/campaigns/{campaignId}/blacklisted`

#### Get Validations by Campaign
- **GET** `/api/validator/validations?campaignId={campaignId}`

#### Get Validation
- **GET** `/api/validator/validations/{validationId}`

#### Get Campaign (from StartupOrg)
- **GET** `/api/validator/campaigns/{campaignId}`

### 3.2 Milestone Verification

#### Verify Milestone Completion
- **POST** `/api/validator/milestones/verify`
- **Body:**
```json
{
  "milestoneId": "string",
  "campaignId": "string",
  "startupId": "string",
  "milestoneReportHash": "string",
  "deliverablesVerified": true,
  "qualityScore": 8.5,
  "comments": "Milestone objectives met",
  "approved": true
}
```

#### Get Milestone Verification
- **GET** `/api/validator/milestones/verifications/{verificationId}`

#### Get Milestone Verification by Milestone
- **GET** `/api/validator/milestones/verifications?milestoneId={milestoneId}`

### 3.3 Risk Assessment

#### Assign Risk Score
- **POST** `/api/validator/risk/assign`
- **Body:**
```json
{
  "campaignId": "string",
  "investorId": "string",
  "riskScore": 7.5,
  "riskLevel": "HIGH",
  "riskFactors": ["Market volatility", "Limited track record"],
  "queryResponse": "Response to investor query",
  "recommendation": "High risk investment"
}
```

#### Get Risk Insight
- **GET** `/api/validator/risk/{campaignId}`

### 3.4 Reporting & Witnessing

#### Send Validation Report to Platform
- **POST** `/api/validator/reports/send-to-platform`
- **Body:**
```json
{
  "campaignId": "string",
  "validationId": "string",
  "campaignHash": "string",
  "overallScore": 8.5,
  "documentScore": 9.0,
  "complianceScore": 8.0,
  "riskScore": 7.5,
  "approved": true,
  "reportSummary": "Campaign approved after thorough review",
  "reportHash": "string"
}
```

#### Get Validation Report
- **GET** `/api/validator/reports?campaignId={campaignId}`

#### Witness Agreement
- **POST** `/api/validator/agreements/witness`
- **Body:**
```json
{
  "agreementId": "string",
  "campaignId": "string",
  "startupId": "string",
  "investorId": "string",
  "investmentAmount": 100000,
  "validatorComments": "Agreement terms verified"
}
```

#### Get Agreement Witness
- **GET** `/api/validator/agreements/{agreementId}/witness`

#### Confirm Campaign Completion
- **POST** `/api/validator/campaigns/{campaignId}/confirm-completion`
- **Body:**
```json
{
  "finalStatus": "SUCCESSFUL|FAILED",
  "completionNotes": "All milestones achieved"
}
```

#### Get Campaign Completion
- **GET** `/api/validator/campaigns/{campaignId}/completion`

### 3.5 Dispute Investigation

#### Accept Dispute Investigation
- **POST** `/api/validator/disputes/{disputeId}/accept-investigation`
- **Body:**
```json
{
  "validatorId": "string",
  "initiatorId": "string",
  "initiatorType": "STARTUP|INVESTOR",
  "respondentId": "string",
  "respondentType": "STARTUP|INVESTOR",
  "campaignId": "string"
}
```

#### Record Investigation Finding
- **POST** `/api/validator/disputes/{disputeId}/findings`
- **Body:**
```json
{
  "investigationId": "string",
  "findingType": "EVIDENCE_VERIFIED|EVIDENCE_INVALID|RULE_VIOLATION|POLICY_BREACH|FRAUD_DETECTED",
  "description": "Finding details",
  "severity": "LOW|MEDIUM|HIGH|CRITICAL",
  "relatedEvidence": "ipfsHash"
}
```

#### Complete Investigation
- **PUT** `/api/validator/disputes/{disputeId}/complete-investigation`
- **Body:**
```json
{
  "investigationId": "string",
  "recommendation": "FAVOR_INITIATOR|FAVOR_RESPONDENT|PARTIAL|DISMISS",
  "recommendedPenalty": "Penalty details"
}
```

#### Investigate Milestone Dispute
- **POST** `/api/validator/disputes/milestone/investigate`
- **Body:**
```json
{
  "disputeId": "string",
  "milestoneId": "string",
  "campaignId": "string",
  "validatorId": "string",
  "milestoneReviewed": true,
  "deliverableStatus": "COMPLETED|PARTIAL|NOT_DELIVERED",
  "qualityAssessment": 75.0,
  "timelineAssessment": "ON_TIME|DELAYED|SEVERELY_DELAYED",
  "delayJustified": false,
  "recommendedAction": "RELEASE_FUNDS|PARTIAL_REFUND|FULL_REFUND",
  "comments": "Investigation notes"
}
```

#### Respond to Dispute (when Validator is respondent)
- **POST** `/api/validator/disputes/{disputeId}/respond`
- **Body:**
```json
{
  "validatorId": "string",
  "responseText": "Response to dispute",
  "justification": "Justification details",
  "supportingDocs": ["ipfsHash1"]
}
```

#### Get Investigation
- **GET** `/api/validator/disputes/investigations/{investigationId}`

#### Get Validator Disputes
- **GET** `/api/validator/disputes?validatorId={validatorId}`

### 3.6 Common Channel Operations

#### Publish Validation Proof
- **POST** `/api/validator/campaigns/{campaignId}/publish-proof`
- **Body:**
```json
{
  "validationHash": "string",
  "status": "APPROVED"
}
```

---

## 4. PLATFORM ORG API

Base path: `/api/platform/`

### 4.1 Campaign Management

#### Publish Campaign to Portal
- **POST** `/api/platform/campaigns/publish`
- **Body:**
```json
{
  "campaignId": "string",
  "startupId": "string",
  "projectName": "string",
  "category": "string",
  "description": "string",
  "goalAmount": 100000,
  "currency": "USD",
  "openDate": "2025-01-01",
  "closeDate": "2025-06-01",
  "durationDays": 150,
  "validationScore": 8.5,
  "validationHash": "string",
  "milestones": []
}
```

#### Verify and Publish
- **POST** `/api/platform/campaigns/verify-publish`
- **Body:**
```json
{
  "campaignId": "string",
  "validationHash": "string"
}
```

#### Get Published Campaign
- **GET** `/api/platform/campaigns/{campaignId}`

#### Get Active Campaigns
- **GET** `/api/platform/campaigns?status=ACTIVE`

#### Close Campaign
- **PUT** `/api/platform/campaigns/{campaignId}/close`
- **Body:**
```json
{
  "finalStatus": "SUCCESSFUL|FAILED|CANCELLED",
  "closureReason": "Funding goal achieved"
}
```

### 4.2 Agreement & Escrow Management

#### Witness Agreement
- **POST** `/api/platform/agreements/witness`
- **Body:**
```json
{
  "agreementId": "string",
  "campaignId": "string",
  "startupId": "string",
  "investorId": "string",
  "investmentAmount": 100000,
  "milestones": []
}
```

#### Get Agreement
- **GET** `/api/platform/agreements/{agreementId}`

#### Get Agreements by Campaign
- **GET** `/api/platform/agreements?campaignId={campaignId}`

#### Get Escrow
- **GET** `/api/platform/escrows/{escrowId}`

### 4.3 Fund Management

#### Trigger Fund Release
- **POST** `/api/platform/funds/release`
- **Body:**
```json
{
  "escrowId": "string",
  "agreementId": "string",
  "campaignId": "string",
  "milestoneId": "string",
  "startupId": "string",
  "amount": 30000,
  "triggerReason": "MILESTONE_VERIFIED"
}
```

#### Get Fund Release
- **GET** `/api/platform/funds/releases/{releaseId}`

### 4.4 Recording & Verification

#### Record Investor Confirmation
- **POST** `/api/platform/records/investor-confirmation`
- **Body:**
```json
{
  "confirmationId": "string",
  "campaignId": "string",
  "investorId": "string",
  "amount": 5000
}
```

#### Record Validator Decision
- **POST** `/api/platform/records/validator-decision`
- **Body:**
```json
{
  "campaignId": "string",
  "validationId": "string",
  "campaignHash": "string",
  "approved": true,
  "overallScore": 8.5,
  "reportHash": "string"
}
```

#### Get Validator Decision
- **GET** `/api/platform/records/validator-decision?campaignId={campaignId}`

### 4.5 Wallet & Token Management

#### Create Wallet
- **POST** `/api/platform/wallets`
- **Body:**
```json
{
  "userId": "string",
  "userType": "STARTUP|INVESTOR|VALIDATOR|PLATFORM",
  "initialBalance": 0
}
```

#### Get Wallet
- **GET** `/api/platform/wallets/{walletId}`

#### Get Wallet by User
- **GET** `/api/platform/wallets?userType={userType}&userId={userId}`

#### Deposit Tokens
- **POST** `/api/platform/wallets/{walletId}/deposit`
- **Body:**
```json
{
  "amount": 1000,
  "reference": "Initial deposit"
}
```

#### Transfer Tokens
- **POST** `/api/platform/wallets/transfer`
- **Body:**
```json
{
  "fromWalletId": "string",
  "toWalletId": "string",
  "amount": 500,
  "transactionType": "TRANSFER|FEE|PENALTY|REFUND|ESCROW|RELEASE",
  "reference": "Campaign payment"
}
```

#### Set Exchange Rate
- **POST** `/api/platform/exchange-rates`
- **Body:**
```json
{
  "currency": "USD",
  "tokenRate": 1.0,
  "effectiveAt": "2025-01-01",
  "expiresAt": "2025-12-31"
}
```

### 4.6 Fee Management

#### Initialize Fee Tiers
- **POST** `/api/platform/fees/tiers/initialize`

#### Get Fee Tier
- **GET** `/api/platform/fees/tiers?goalAmount={goalAmount}`

#### Collect Campaign Fee
- **POST** `/api/platform/fees/campaign/collect`
- **Body:**
```json
{
  "campaignId": "string",
  "startupId": "string",
  "startupWalletId": "string",
  "campaignGoalAmount": 100000
}
```

#### Initialize Dispute Fee Tiers
- **POST** `/api/platform/fees/dispute-tiers/initialize`

#### Get Dispute Fee Tier
- **GET** `/api/platform/fees/dispute-tiers?claimAmount={claimAmount}`

#### Collect Dispute Fee
- **POST** `/api/platform/fees/dispute/collect`
- **Body:**
```json
{
  "disputeId": "string",
  "initiatorId": "string",
  "initiatorWalletId": "string",
  "claimAmount": 10000
}
```

#### Process Dispute Fee Outcome
- **PUT** `/api/platform/fees/dispute/{feeRecordId}/outcome`
- **Body:**
```json
{
  "outcome": "INITIATOR_WIN|RESPONDENT_WIN|PARTIAL|DISMISSED",
  "refundPercentage": 100
}
```

#### Get Dispute Fee Record
- **GET** `/api/platform/fees/dispute/{feeRecordId}`

#### Get Dispute Fee by Dispute ID
- **GET** `/api/platform/fees/dispute?disputeId={disputeId}`

#### Get Dispute Fee Tiers
- **GET** `/api/platform/fees/dispute-tiers`

#### Get User Dispute Fees
- **GET** `/api/platform/fees/dispute?userId={userId}`

### 4.7 Rating & Reputation

#### Record Rating
- **POST** `/api/platform/ratings`
- **Body:**
```json
{
  "ratedUserType": "STARTUP|INVESTOR|VALIDATOR",
  "ratedUserId": "string",
  "raterUserType": "STARTUP|INVESTOR|VALIDATOR",
  "raterUserId": "string",
  "context": "AGREEMENT|MILESTONE|DISPUTE",
  "contextId": "string",
  "rating": 4.5,
  "comment": "Excellent collaboration"
}
```

#### Get Rating Aggregate
- **GET** `/api/platform/ratings/aggregate?userType={userType}&userId={userId}`

#### Get Reputation
- **GET** `/api/platform/reputation?userType={userType}&userId={userId}`

#### Check User Status
- **GET** `/api/platform/users/status?userType={userType}&userId={userId}`

### 4.8 Dispute Management

#### Create Dispute
- **POST** `/api/platform/disputes`
- **Body:**
```json
{
  "initiatorType": "STARTUP|INVESTOR|VALIDATOR",
  "initiatorId": "string",
  "respondentType": "STARTUP|INVESTOR|VALIDATOR",
  "respondentId": "string",
  "disputeType": "MILESTONE_DISPUTE|AGREEMENT_BREACH|FRAUD|MISREPRESENTATION|REFUND_DISPUTE",
  "campaignId": "string",
  "agreementId": "string",
  "title": "Dispute title",
  "description": "Dispute details",
  "claimAmount": 10000,
  "evidenceHashes": ["ipfsHash1"]
}
```

#### Get Dispute
- **GET** `/api/platform/disputes/{disputeId}`

#### Submit Evidence
- **POST** `/api/platform/disputes/{disputeId}/evidence`
- **Body:**
```json
{
  "submittedBy": "string",
  "submitterType": "STARTUP|INVESTOR|VALIDATOR",
  "evidenceHashes": ["ipfsHash1"],
  "description": "Evidence description"
}
```

#### Assign Investigator
- **PUT** `/api/platform/disputes/{disputeId}/assign-investigator`
- **Body:**
```json
{
  "investigatorId": "string"
}
```

#### Add Investigation Note
- **POST** `/api/platform/disputes/{disputeId}/investigation-notes`
- **Body:**
```json
{
  "note": "Investigation progress update"
}
```

#### Enable Voting
- **PUT** `/api/platform/disputes/{disputeId}/enable-voting`
- **Body:**
```json
{
  "eligibleVoters": ["voter1", "voter2", "voter3"]
}
```

#### Commit Vote
- **POST** `/api/platform/disputes/{disputeId}/vote/commit`
- **Body:**
```json
{
  "voterId": "string",
  "voteHash": "hashedVote"
}
```

#### Reveal Vote
- **POST** `/api/platform/disputes/{disputeId}/vote/reveal`
- **Body:**
```json
{
  "voterId": "string",
  "vote": "INITIATOR|RESPONDENT",
  "secret": "string"
}
```

#### Tally Votes
- **POST** `/api/platform/disputes/{disputeId}/vote/tally`

#### Resolve Dispute
- **PUT** `/api/platform/disputes/{disputeId}/resolve`
- **Body:**
```json
{
  "resolution": "INITIATOR_WIN|RESPONDENT_WIN|PARTIAL|DISMISSED",
  "resolutionDetails": "Resolution explanation",
  "penaltyForRespondent": {
    "type": "FINANCIAL|REPUTATION|SUSPENSION|PERMANENT_BAN",
    "amount": 5000
  },
  "refundOrder": {
    "amount": 8000,
    "recipient": "string"
  }
}
```

### 4.9 Refund Management

#### Request Refund
- **POST** `/api/platform/refunds`
- **Body:**
```json
{
  "investorId": "string",
  "campaignId": "string",
  "agreementId": "string",
  "originalAmount": 100000,
  "requestedAmount": 85000,
  "refundReason": "EARLY_WITHDRAWAL",
  "deductionPercent": 15
}
```

### 4.10 Common Channel Operations

#### Publish Global Metrics
- **POST** `/api/platform/metrics/publish`
- **Body:**
```json
{
  "totalCampaigns": 150,
  "activeCampaigns": 45,
  "successfulCampaigns": 80,
  "totalInvestorCount": 5000
}
```

#### Get Latest Global Metrics
- **GET** `/api/platform/metrics/latest`

---

## Authentication & Authorization

All API endpoints require JWT authentication:

**Headers:**
```
Authorization: Bearer <ACCESS_TOKEN>
```

Each organization has its own authentication context and channel access permissions.

---

## HTTP Methods Summary

We use only **3 HTTP methods** for simplicity:

### 1. **GET** - Retrieve/Read Data
- Used for: Reading existing data, querying information
- Examples: `GetCampaign`, `GetInvestment`, `GetValidation`
- No request body needed
- Returns data in response

### 2. **POST** - Create/Submit New Data
- Used for: Creating new resources, submitting actions
- Examples: `CreateCampaign`, `MakeInvestment`, `ValidateCampaign`, `SubmitDispute`
- Requires request body with data
- Returns created resource or confirmation

### 3. **PUT** - Update Existing Data
- Used for: Modifying existing resources, changing status
- Examples: `UpdateCampaign`, `ApproveOrRejectCampaign`, `WithdrawInvestment`
- Requires request body with updated data
- Returns updated resource or confirmation

**Note:** We don't use PATCH or DELETE methods for this platform.

---

## Example API Calls

### Example 1: Create a Campaign (Startup)
```bash
curl -X POST http://localhost:3000/api/startup/campaigns \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "startupId": "startup_001",
    "projectName": "EcoTech Solutions",
    "description": "Sustainable tech for cleaner energy",
    "category": "Technology",
    "goalAmount": 250000,
    "currency": "USD",
    "durationDays": 90
  }'
```

### Example 2: Make an Investment (Investor)
```bash
curl -X POST http://localhost:3000/api/investor/investments \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "campaignId": "campaign_001",
    "investorId": "investor_001",
    "amount": 5000,
    "currency": "USD"
  }'
```

### Example 3: Validate Campaign (Validator)
```bash
curl -X POST http://localhost:3000/api/validator/campaigns/campaign_001/validate \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "validatorId": "validator_001",
    "campaignHash": "abc123hash",
    "documentsReviewed": ["doc1.pdf", "doc2.pdf"]
  }'
```

### Example 4: Get Active Campaigns (Platform)
```bash
curl -X GET http://localhost:3000/api/platform/campaigns?status=ACTIVE \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```

---

## Response Formats

All successful responses return JSON format:

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": { /* response data */ },
  "message": "Operation completed successfully"
}
```

**Error Response (4xx/5xx):**
```json
{
  "success": false,
  "error": "Error message description",
  "code": "ERROR_CODE"
}
```

---

## Notes

1. All write operations (POST, PUT) require proper authentication and authorization
2. Cross-organization invocations are handled through Hyperledger Fabric's chaincode-to-chaincode calls
3. Hash verification ensures data integrity across organizations
4. Dispute resolution involves multi-party consensus through voting mechanisms
5. Token-based wallet system manages all financial transactions on the platform
