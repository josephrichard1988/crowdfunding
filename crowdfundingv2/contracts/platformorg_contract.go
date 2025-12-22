package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// PlatformContract provides functions for PlatformOrg operations using PDC
type PlatformContract struct {
	contractapi.Contract
}

// ============================================================================
// DATA STRUCTURES
// ============================================================================

// PublishedCampaign represents a campaign published on the platform portal (22-parameter format)
type PublishedCampaign struct {
	CampaignID string `json:"campaignId"`
	StartupID  string `json:"startupId"`

	// 22 Core Parameters
	Category          string   `json:"category"`
	Deadline          string   `json:"deadline"`
	Currency          string   `json:"currency"`
	HasRaised         bool     `json:"has_raised"`
	HasGovGrants      bool     `json:"has_gov_grants"`
	IncorpDate        string   `json:"incorp_date"`
	ProjectStage      string   `json:"project_stage"`
	Sector            string   `json:"sector"`
	Tags              []string `json:"tags"`
	TeamAvailable     bool     `json:"team_available"`
	InvestorCommitted bool     `json:"investor_committed"`
	Duration          int      `json:"duration"`
	FundingDay        int      `json:"funding_day"`
	FundingMonth      int      `json:"funding_month"`
	FundingYear       int      `json:"funding_year"`
	GoalAmount        float64  `json:"goal_amount"`
	InvestmentRange   string   `json:"investment_range"`
	ProjectName       string   `json:"project_name"`
	Description       string   `json:"description"`
	Documents         []string `json:"documents"`

	// Calculated/Status Fields
	OpenDate           string      `json:"open_date"`
	CloseDate          string      `json:"close_date"`
	FundsRaisedAmount  float64     `json:"funds_raised_amount"`
	FundsRaisedPercent float64     `json:"funds_raised_percent"`
	ValidationScore    float64     `json:"validationScore"`
	ValidationHash     string      `json:"validationHash"`
	ValidationVerified bool        `json:"validationVerified"`
	RiskScore          float64     `json:"riskScore"`
	RiskLevel          string      `json:"riskLevel"`
	Status             string      `json:"status"`
	InvestorCount      int         `json:"investorCount"`
	TotalConfirmed     float64     `json:"totalConfirmed"`
	Milestones         []Milestone `json:"milestones"`
	AgreementIDs       []string    `json:"agreementIds"`
	PublishedAt        string      `json:"publishedAt"`
	UpdatedAt          string      `json:"updatedAt"`
}

// Agreement represents investment agreement
// FundEscrow represents funds held in escrow by Platform
type FundEscrow struct {
	EscrowID       string  `json:"escrowId"`
	AgreementID    string  `json:"agreementId"`
	CampaignID     string  `json:"campaignId"`
	InvestorID     string  `json:"investorId"`
	StartupID      string  `json:"startupId"`
	TotalAmount    float64 `json:"totalAmount"`
	ReleasedAmount float64 `json:"releasedAmount"`
	HeldAmount     float64 `json:"heldAmount"`
	Currency       string  `json:"currency"`
	Status         string  `json:"status"`
	CreatedAt      string  `json:"createdAt"`
	UpdatedAt      string  `json:"updatedAt"`
}

// InvestorConfirmationRecord represents recorded investor confirmation
type InvestorConfirmationRecord struct {
	RecordID       string  `json:"recordId"`
	ConfirmationID string  `json:"confirmationId"`
	CampaignID     string  `json:"campaignId"`
	InvestorID     string  `json:"investorId"`
	Amount         float64 `json:"amount"`
	Currency       string  `json:"currency"`
	RecordedAt     string  `json:"recordedAt"`
}

// ValidatorDecisionRecord represents recorded validator decision
type ValidatorDecisionRecord struct {
	RecordID     string  `json:"recordId"`
	CampaignID   string  `json:"campaignId"`
	ValidationID string  `json:"validationId"`
	CampaignHash string  `json:"campaignHash"`
	Approved     bool    `json:"approved"`
	OverallScore float64 `json:"overallScore"`
	ReportHash   string  `json:"reportHash"`
	RecordedAt   string  `json:"recordedAt"`
}

// FundRelease represents fund release to startup
type FundRelease struct {
	ReleaseID     string  `json:"releaseId"`
	EscrowID      string  `json:"escrowId"`
	AgreementID   string  `json:"agreementId"`
	CampaignID    string  `json:"campaignId"`
	MilestoneID   string  `json:"milestoneId"`
	StartupID     string  `json:"startupId"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	Status        string  `json:"status"`
	TriggerReason string  `json:"triggerReason"`
	ReleasedAt    string  `json:"releasedAt"`
}

// CampaignClosure represents campaign closure record
type CampaignClosure struct {
	ClosureID          string  `json:"closureId"`
	CampaignID         string  `json:"campaignId"`
	FinalStatus        string  `json:"finalStatus"`
	FinalAmount        float64 `json:"finalAmount"`
	FinalInvestorCount int     `json:"finalInvestorCount"`
	ClosureReason      string  `json:"closureReason"`
	ClosedAt           string  `json:"closedAt"`
}

// GlobalMetrics for public ledger
type GlobalMetrics struct {
	MetricsID           string `json:"metricsId"`
	TotalCampaigns      int    `json:"totalCampaigns"`
	ActiveCampaigns     int    `json:"activeCampaigns"`
	SuccessfulCampaigns int    `json:"successfulCampaigns"`
	TotalInvestorCount  int    `json:"totalInvestorCount"`
	MetricsHash         string `json:"metricsHash"`
	PublishedAt         string `json:"publishedAt"`
}

// Wallet represents a user's cryptographic wallet
type Wallet struct {
	WalletID       string  `json:"walletId"`
	UserID         string  `json:"userId"`
	UserType       string  `json:"userType"`
	Balance        float64 `json:"balance"`
	LockedBalance  float64 `json:"lockedBalance"`
	TotalDeposited float64 `json:"totalDeposited"`
	TotalWithdrawn float64 `json:"totalWithdrawn"`
	TotalPenalties float64 `json:"totalPenalties"`
	CreatedAt      string  `json:"createdAt"`
	UpdatedAt      string  `json:"updatedAt"`
}

// TokenTransaction represents a token transfer/exchange transaction
type TokenTransaction struct {
	TransactionID   string  `json:"transactionId"`
	FromWalletID    string  `json:"fromWalletId"`
	ToWalletID      string  `json:"toWalletId"`
	Amount          float64 `json:"amount"`
	TransactionType string  `json:"transactionType"`
	Reference       string  `json:"reference"`
	Status          string  `json:"status"`
	CreatedAt       string  `json:"createdAt"`
}

// TokenExchangeRate represents exchange rate
type TokenExchangeRate struct {
	RateID      string  `json:"rateId"`
	Currency    string  `json:"currency"`
	TokenRate   float64 `json:"tokenRate"`
	EffectiveAt string  `json:"effectiveAt"`
	ExpiresAt   string  `json:"expiresAt"`
	SetBy       string  `json:"setBy"`
}

// FeeTier represents platform fee tier
type FeeTier struct {
	TierID        string  `json:"tierId"`
	MinGoalAmount float64 `json:"minGoalAmount"`
	MaxGoalAmount float64 `json:"maxGoalAmount"`
	FeePercentage float64 `json:"feePercentage"`
	Description   string  `json:"description"`
}

// FeeCollection represents collected fee
type FeeCollection struct {
	CollectionID  string  `json:"collectionId"`
	CampaignID    string  `json:"campaignId"`
	StartupID     string  `json:"startupId"`
	FeeType       string  `json:"feeType"`
	Amount        float64 `json:"amount"`
	GoalAmount    float64 `json:"goalAmount"`
	FeePercentage float64 `json:"feePercentage"`
	Status        string  `json:"status"`
	CollectedAt   string  `json:"collectedAt"`
}

// DisputeFeeTier represents dispute fee tier
type DisputeFeeTier struct {
	TierID         string  `json:"tierId"`
	MinClaimAmount float64 `json:"minClaimAmount"`
	MaxClaimAmount float64 `json:"maxClaimAmount"`
	FeeAmount      float64 `json:"feeAmount"`
	Description    string  `json:"description"`
}

// DisputeFeeRecord represents dispute fee payment
type DisputeFeeRecord struct {
	FeeRecordID  string  `json:"feeRecordId"`
	DisputeID    string  `json:"disputeId"`
	InitiatorID  string  `json:"initiatorId"`
	ClaimAmount  float64 `json:"claimAmount"`
	FeeAmount    float64 `json:"feeAmount"`
	Status       string  `json:"status"`
	RefundAmount float64 `json:"refundAmount"`
	CollectedAt  string  `json:"collectedAt"`
	ProcessedAt  string  `json:"processedAt"`
}

// RatingRecord represents a rating
type RatingRecord struct {
	RatingID      string  `json:"ratingId"`
	RatedUserType string  `json:"ratedUserType"`
	RatedUserId   string  `json:"ratedUserId"`
	RaterUserType string  `json:"raterUserType"`
	RaterUserId   string  `json:"raterUserId"`
	Context       string  `json:"context"`
	ContextId     string  `json:"contextId"`
	Rating        float64 `json:"rating"`
	Comment       string  `json:"comment"`
	CreatedAt     string  `json:"createdAt"`
}

// RatingAggregate represents aggregated ratings
type RatingAggregate struct {
	UserType      string  `json:"userType"`
	UserId        string  `json:"userId"`
	TotalRatings  int     `json:"totalRatings"`
	AverageRating float64 `json:"averageRating"`
	UpdatedAt     string  `json:"updatedAt"`
}

// ReputationScore represents user reputation
type ReputationScore struct {
	UserType        string  `json:"userType"`
	UserId          string  `json:"userId"`
	ReputationScore float64 `json:"reputationScore"`
	DisputesWon     int     `json:"disputesWon"`
	DisputesLost    int     `json:"disputesLost"`
	TotalDisputes   int     `json:"totalDisputes"`
	SuccessfulDeals int     `json:"successfulDeals"`
	Status          string  `json:"status"`
	UpdatedAt       string  `json:"updatedAt"`
}

// Dispute represents a dispute
type Dispute struct {
	DisputeID          string           `json:"disputeId"`
	InitiatorType      string           `json:"initiatorType"`
	InitiatorID        string           `json:"initiatorId"`
	RespondentType     string           `json:"respondentType"`
	RespondentID       string           `json:"respondentId"`
	DisputeType        string           `json:"disputeType"`
	CampaignID         string           `json:"campaignId"`
	AgreementID        string           `json:"agreementId"`
	Title              string           `json:"title"`
	Description        string           `json:"description"`
	ClaimAmount        float64          `json:"claimAmount"`
	EvidenceHashes     []string         `json:"evidenceHashes"`
	Status             string           `json:"status"`
	InvestigatorID     string           `json:"investigatorId"`
	InvestigationNotes []string         `json:"investigationNotes"`
	VotingEnabled      bool             `json:"votingEnabled"`
	EligibleVoters     []string         `json:"eligibleVoters"`
	Votes              []VoteCommitment `json:"votes"`
	Resolution         string           `json:"resolution"`
	ResolutionDetails  string           `json:"resolutionDetails"`
	CreatedAt          string           `json:"createdAt"`
	ResolvedAt         string           `json:"resolvedAt"`
}

// VoteCommitment for commit-reveal voting
type VoteCommitment struct {
	VoterID     string `json:"voterId"`
	VoteHash    string `json:"voteHash"`
	Revealed    bool   `json:"revealed"`
	Vote        string `json:"vote"`
	CommittedAt string `json:"committedAt"`
	RevealedAt  string `json:"revealedAt"`
}

// Penalty represents a penalty
type Penalty struct {
	PenaltyID   string  `json:"penaltyId"`
	UserType    string  `json:"userType"`
	UserID      string  `json:"userId"`
	DisputeID   string  `json:"disputeId"`
	PenaltyType string  `json:"penaltyType"`
	Amount      float64 `json:"amount"`
	AppliedAt   string  `json:"appliedAt"`
}

// RefundOrder represents a refund order
type RefundOrder struct {
	RefundID    string  `json:"refundId"`
	DisputeID   string  `json:"disputeId"`
	Recipient   string  `json:"recipient"`
	Amount      float64 `json:"amount"`
	ProcessedAt string  `json:"processedAt"`
}

// ============================================================================
// INIT
// ============================================================================

func (p *PlatformContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	fmt.Println("PlatformOrg contract initialized with PDC support")
	return nil
}

// ============================================================================
// CAMPAIGN MANAGEMENT - Using PDC
// ============================================================================

// PublishCampaignToPortal verifies validation and publishes campaign
// Platform reads campaign from StartupPlatformCollection using campaignID
// Verifies hash with ValidatorPlatformCollection before publishing
func (p *PlatformContract) PublishCampaignToPortal(
	ctx contractapi.TransactionContextInterface,
	campaignID string,
	validationHash string,
) error {

	// Step 1: Read campaign data from StartupPlatformCollection
	campaignDataJSON, err := ctx.GetStub().GetPrivateData(StartupPlatformCollection, "CAMPAIGN_SHARE_"+campaignID)
	if err != nil || campaignDataJSON == nil {
		return fmt.Errorf("campaign not found in StartupPlatformCollection. Startup must share campaign first: %v", err)
	}

	var campaignData map[string]interface{}
	err = json.Unmarshal(campaignDataJSON, &campaignData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal campaign data: %v", err)
	}

	// Step 2: Verify validation hash with Validator
	validationApprovalJSON, err := ctx.GetStub().GetPrivateData(ValidatorPlatformCollection, "VALIDATION_APPROVAL_"+campaignID)
	if err != nil || validationApprovalJSON == nil {
		return fmt.Errorf("validation approval not found in ValidatorPlatformCollection: %v", err)
	}

	var validationApproval map[string]interface{}
	err = json.Unmarshal(validationApprovalJSON, &validationApproval)
	if err != nil {
		return fmt.Errorf("failed to unmarshal validation approval: %v", err)
	}

	// Step 3: Verify hash matches
	validatorHash, ok := validationApproval["validationHash"].(string)
	if !ok || validatorHash != validationHash {
		return fmt.Errorf("validation hash mismatch. Platform verification failed. Expected: %s, Received: %s", validatorHash, validationHash)
	}

	campaignHash, ok := campaignData["validationHash"].(string)
	if !ok || campaignHash != validationHash {
		return fmt.Errorf("campaign validation hash mismatch. Cannot publish")
	}

	// Check if already published
	publicCheckJSON, _ := ctx.GetStub().GetState("CAMPAIGN_PUBLIC_" + campaignID)
	if publicCheckJSON != nil {
		var publicCheck map[string]interface{}
		err := json.Unmarshal(publicCheckJSON, &publicCheck)
		if err == nil {
			if status, ok := publicCheck["status"].(string); ok && status == "PUBLISHED" {
				return fmt.Errorf("campaign %s is already published", campaignID)
			}
		}
	}

	// Step 4: Hash verification successful - proceed with publishing
	timestamp := time.Now().Format(time.RFC3339)

	// Extract campaign fields
	campaignID = campaignData["campaignId"].(string)
	startupID := campaignData["startupId"].(string)
	category := campaignData["category"].(string)
	deadline := campaignData["deadline"].(string)
	currency := campaignData["currency"].(string)
	projectName := campaignData["projectName"].(string)
	description := campaignData["description"].(string)

	// Handle type conversions for boolean and numeric fields
	hasRaised := campaignData["hasRaised"].(bool)
	hasGovGrants := campaignData["hasGovGrants"].(bool)
	teamAvailable := campaignData["teamAvailable"].(bool)
	investorCommitted := campaignData["investorCommitted"].(bool)

	incorpDate := campaignData["incorpDate"].(string)
	projectStage := campaignData["projectStage"].(string)
	sector := campaignData["sector"].(string)
	investmentRange := campaignData["investmentRange"].(string)

	duration := int(campaignData["duration"].(float64))
	fundingDay := int(campaignData["fundingDay"].(float64))
	fundingMonth := int(campaignData["fundingMonth"].(float64))
	fundingYear := int(campaignData["fundingYear"].(float64))
	goalAmount := campaignData["goalAmount"].(float64)

	// Extract validation scores
	validationScore := validationApproval["dueDiligenceScore"].(float64)
	riskScore := validationApproval["riskScore"].(float64)
	riskLevel := validationApproval["riskLevel"].(string)

	// Extract arrays
	tagsInterface := campaignData["tags"].([]interface{})
	tags := make([]string, len(tagsInterface))
	for i, v := range tagsInterface {
		tags[i] = v.(string)
	}

	docsInterface := campaignData["documents"].([]interface{})
	documents := make([]string, len(docsInterface))
	for i, v := range docsInterface {
		documents[i] = v.(string)
	}

	// Calculate dates
	openDate := fmt.Sprintf("%04d-%02d-%02d", fundingYear, fundingMonth, fundingDay)

	// Create published campaign
	campaign := PublishedCampaign{
		CampaignID:         campaignID,
		StartupID:          startupID,
		Category:           category,
		Deadline:           deadline,
		Currency:           currency,
		HasRaised:          hasRaised,
		HasGovGrants:       hasGovGrants,
		IncorpDate:         incorpDate,
		ProjectStage:       projectStage,
		Sector:             sector,
		Tags:               tags,
		TeamAvailable:      teamAvailable,
		InvestorCommitted:  investorCommitted,
		Duration:           duration,
		FundingDay:         fundingDay,
		FundingMonth:       fundingMonth,
		FundingYear:        fundingYear,
		GoalAmount:         goalAmount,
		InvestmentRange:    investmentRange,
		ProjectName:        projectName,
		Description:        description,
		Documents:          documents,
		OpenDate:           openDate,
		CloseDate:          deadline,
		ValidationScore:    validationScore,
		ValidationHash:     validationHash,
		ValidationVerified: true,
		RiskScore:          riskScore,
		RiskLevel:          riskLevel,
		FundsRaisedAmount:  0,
		FundsRaisedPercent: 0,
		InvestorCount:      0,
		Status:             "PUBLISHED",
		PublishedAt:        timestamp,
		Milestones:         []Milestone{},
	}

	campaignJSON, err := json.Marshal(campaign)
	if err != nil {
		return fmt.Errorf("failed to marshal campaign: %v", err)
	}

	// Store in PlatformPrivateCollection
	err = ctx.GetStub().PutPrivateData(PlatformPrivateCollection, "PUBLISHED_CAMPAIGN_"+campaignID, campaignJSON)
	if err != nil {
		return fmt.Errorf("failed to store published campaign: %v", err)
	}

	// Store FULL campaign info on world state for all investors to see
	// This makes the 22-parameters visible to everyone
	err = ctx.GetStub().PutState("CAMPAIGN_PUBLIC_"+campaignID, campaignJSON)
	if err != nil {
		return fmt.Errorf("failed to store public campaign info: %v", err)
	}

	// Step 5: Send success notification to Startup via StartupPlatformCollection
	notification := map[string]interface{}{
		"campaignId":      campaignID,
		"status":          "PUBLISHED",
		"message":         fmt.Sprintf("Campaign '%s' has been successfully published on the platform", projectName),
		"publishedAt":     timestamp,
		"validationScore": validationScore,
		"riskScore":       riskScore,
		"riskLevel":       riskLevel,
	}

	notificationJSON, _ := json.Marshal(notification)
	err = ctx.GetStub().PutPrivateData(StartupPlatformCollection, "PUBLISH_NOTIFICATION_"+campaignID, notificationJSON)
	if err != nil {
		return fmt.Errorf("failed to send notification to startup: %v", err)
	}

	// Step 6: Update the shared campaign status so Platform and Startup can see it's published
	campaignData["status"] = "PUBLISHED"
	campaignData["publishedAt"] = timestamp
	updatedCampaignJSON, _ := json.Marshal(campaignData)
	err = ctx.GetStub().PutPrivateData(StartupPlatformCollection, "CAMPAIGN_SHARE_"+campaignID, updatedCampaignJSON)
	if err != nil {
		return fmt.Errorf("failed to update shared campaign status: %v", err)
	}

	return nil
}

// VerifyAndPublish verifies validation hash and publishes campaign
func (p *PlatformContract) VerifyAndPublish(
	ctx contractapi.TransactionContextInterface,
	campaignID string,
	validationHash string,
) error {

	// Get validation report from ValidatorPlatformCollection
	reportKey := fmt.Sprintf("VALIDATION_REPORT_%s", campaignID)
	reportJSON, err := ctx.GetStub().GetPrivateData(ValidatorPlatformCollection, reportKey)
	if err != nil || reportJSON == nil {
		return fmt.Errorf("validation report not found: %v", err)
	}

	var report map[string]interface{}
	json.Unmarshal(reportJSON, &report)

	// Verify hash
	if report["campaignHash"].(string) != validationHash {
		return fmt.Errorf("validation hash mismatch")
	}

	// Get campaign
	campaignJSON, err := ctx.GetStub().GetPrivateData(PlatformPrivateCollection, "PUBLISHED_CAMPAIGN_"+campaignID)
	if err != nil || campaignJSON == nil {
		return fmt.Errorf("campaign not found: %v", err)
	}

	var campaign PublishedCampaign
	json.Unmarshal(campaignJSON, &campaign)

	timestamp := time.Now().Format(time.RFC3339)
	campaign.ValidationVerified = true
	campaign.Status = "PUBLISHED"
	campaign.UpdatedAt = timestamp

	campaignJSON, _ = json.Marshal(campaign)
	err = ctx.GetStub().PutPrivateData(PlatformPrivateCollection, "PUBLISHED_CAMPAIGN_"+campaignID, campaignJSON)
	if err != nil {
		return fmt.Errorf("failed to update campaign: %v", err)
	}

	// Update public info
	publicInfo := map[string]interface{}{
		"campaignId":  campaignID,
		"projectName": campaign.ProjectName,
		"category":    campaign.Category,
		"goalAmount":  campaign.GoalAmount,
		"status":      "ACTIVE",
		"publishedAt": campaign.PublishedAt,
	}
	publicJSON, _ := json.Marshal(publicInfo)
	ctx.GetStub().PutState("CAMPAIGN_PUBLIC_"+campaignID, publicJSON)

	return nil
}

// WitnessAgreement witnesses an agreement from ThreePartyCollection
func (p *PlatformContract) WitnessAgreement(
	ctx contractapi.TransactionContextInterface,
	agreementID string,
	campaignID string,
	startupID string,
	investorID string,
	investmentAmount float64,
	currency string,
	milestonesJSON string,
) error {

	var milestones []Milestone
	if milestonesJSON != "" {
		json.Unmarshal([]byte(milestonesJSON), &milestones)
	}

	timestamp := time.Now().Format(time.RFC3339)

	agreement := Agreement{
		AgreementID:       agreementID,
		CampaignID:        campaignID,
		StartupID:         startupID,
		InvestorID:        investorID,
		InvestmentAmount:  investmentAmount,
		Currency:          currency,
		Milestones:        milestones,
		Status:            "ACTIVE",
		PlatformWitnessed: true,
		WitnessedAt:       timestamp,
		CreatedAt:         timestamp,
	}

	agreementJSON, err := json.Marshal(agreement)
	if err != nil {
		return fmt.Errorf("failed to marshal agreement: %v", err)
	}

	// Store in ThreePartyCollection
	err = ctx.GetStub().PutPrivateData(ThreePartyCollection, "AGREEMENT_"+agreementID, agreementJSON)
	if err != nil {
		return fmt.Errorf("failed to witness agreement: %v", err)
	}

	// Create escrow
	escrowID := fmt.Sprintf("ESCROW_%s", agreementID)
	escrow := FundEscrow{
		EscrowID:       escrowID,
		AgreementID:    agreementID,
		CampaignID:     campaignID,
		InvestorID:     investorID,
		StartupID:      startupID,
		TotalAmount:    investmentAmount,
		ReleasedAmount: 0,
		HeldAmount:     investmentAmount,
		Currency:       currency,
		Status:         "ACTIVE",
		CreatedAt:      timestamp,
		UpdatedAt:      timestamp,
	}

	escrowJSON, _ := json.Marshal(escrow)
	ctx.GetStub().PutPrivateData(PlatformPrivateCollection, "ESCROW_"+escrowID, escrowJSON)

	return nil
}

// TriggerFundRelease releases funds based on milestone verification
func (p *PlatformContract) TriggerFundRelease(
	ctx contractapi.TransactionContextInterface,
	releaseID string,
	escrowID string,
	agreementID string,
	campaignID string,
	milestoneID string,
	startupID string,
	amount float64,
	triggerReason string,
) error {

	timestamp := time.Now().Format(time.RFC3339)

	release := FundRelease{
		ReleaseID:     releaseID,
		EscrowID:      escrowID,
		AgreementID:   agreementID,
		CampaignID:    campaignID,
		MilestoneID:   milestoneID,
		StartupID:     startupID,
		Amount:        amount,
		Currency:      "USD",
		Status:        "RELEASED",
		TriggerReason: triggerReason,
		ReleasedAt:    timestamp,
	}

	releaseJSON, err := json.Marshal(release)
	if err != nil {
		return fmt.Errorf("failed to marshal release: %v", err)
	}

	// Store in StartupPlatformCollection so startup can see it
	err = ctx.GetStub().PutPrivateData(StartupPlatformCollection, "FUND_RELEASE_"+releaseID, releaseJSON)
	if err != nil {
		return fmt.Errorf("failed to release funds: %v", err)
	}

	// Update escrow
	escrowJSON, err := ctx.GetStub().GetPrivateData(PlatformPrivateCollection, "ESCROW_"+escrowID)
	if err == nil && escrowJSON != nil {
		var escrow FundEscrow
		json.Unmarshal(escrowJSON, &escrow)
		escrow.ReleasedAmount += amount
		escrow.HeldAmount -= amount
		escrow.UpdatedAt = timestamp
		if escrow.HeldAmount <= 0 {
			escrow.Status = "FULLY_RELEASED"
		} else {
			escrow.Status = "PARTIALLY_RELEASED"
		}
		escrowJSON, _ = json.Marshal(escrow)
		ctx.GetStub().PutPrivateData(PlatformPrivateCollection, "ESCROW_"+escrowID, escrowJSON)
	}

	return nil
}

// CloseCampaign closes a campaign
func (p *PlatformContract) CloseCampaign(
	ctx contractapi.TransactionContextInterface,
	closureID string,
	campaignID string,
	finalStatus string,
	closureReason string,
) error {

	// Get campaign
	campaignJSON, err := ctx.GetStub().GetPrivateData(PlatformPrivateCollection, "PUBLISHED_CAMPAIGN_"+campaignID)
	if err != nil || campaignJSON == nil {
		return fmt.Errorf("campaign not found: %v", err)
	}

	var campaign PublishedCampaign
	json.Unmarshal(campaignJSON, &campaign)

	timestamp := time.Now().Format(time.RFC3339)

	closure := CampaignClosure{
		ClosureID:          closureID,
		CampaignID:         campaignID,
		FinalStatus:        finalStatus,
		FinalAmount:        campaign.FundsRaisedAmount,
		FinalInvestorCount: campaign.InvestorCount,
		ClosureReason:      closureReason,
		ClosedAt:           timestamp,
	}

	closureJSON, _ := json.Marshal(closure)
	ctx.GetStub().PutState("CAMPAIGN_CLOSURE_"+closureID, closureJSON)

	// Update campaign status
	campaign.Status = "CLOSED"
	campaign.UpdatedAt = timestamp
	campaignJSON, _ = json.Marshal(campaign)
	ctx.GetStub().PutPrivateData(PlatformPrivateCollection, "PUBLISHED_CAMPAIGN_"+campaignID, campaignJSON)

	return nil
}

// RecordInvestorConfirmation records investor confirmation from InvestorPlatformCollection
func (p *PlatformContract) RecordInvestorConfirmation(
	ctx contractapi.TransactionContextInterface,
	recordID string,
	confirmationID string,
	campaignID string,
	investorID string,
	amount float64,
	currency string,
) error {

	timestamp := time.Now().Format(time.RFC3339)

	record := InvestorConfirmationRecord{
		RecordID:       recordID,
		ConfirmationID: confirmationID,
		CampaignID:     campaignID,
		InvestorID:     investorID,
		Amount:         amount,
		Currency:       currency,
		RecordedAt:     timestamp,
	}

	recordJSON, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal record: %v", err)
	}

	// Store in PlatformPrivateCollection
	err = ctx.GetStub().PutPrivateData(PlatformPrivateCollection, "INVESTOR_CONFIRMATION_"+recordID, recordJSON)
	if err != nil {
		return fmt.Errorf("failed to record confirmation: %v", err)
	}

	// Update campaign stats
	campaignJSON, err := ctx.GetStub().GetPrivateData(PlatformPrivateCollection, "PUBLISHED_CAMPAIGN_"+campaignID)
	if err == nil && campaignJSON != nil {
		var campaign PublishedCampaign
		json.Unmarshal(campaignJSON, &campaign)
		campaign.TotalConfirmed += amount
		campaign.InvestorCount++
		campaignJSON, _ = json.Marshal(campaign)
		ctx.GetStub().PutPrivateData(PlatformPrivateCollection, "PUBLISHED_CAMPAIGN_"+campaignID, campaignJSON)
	}

	return nil
}

// RecordValidatorDecision records validator decision from ValidatorPlatformCollection
func (p *PlatformContract) RecordValidatorDecision(
	ctx contractapi.TransactionContextInterface,
	recordID string,
	campaignID string,
	validationID string,
	campaignHash string,
	approved bool,
	overallScore float64,
	reportHash string,
) error {

	timestamp := time.Now().Format(time.RFC3339)

	record := ValidatorDecisionRecord{
		RecordID:     recordID,
		CampaignID:   campaignID,
		ValidationID: validationID,
		CampaignHash: campaignHash,
		Approved:     approved,
		OverallScore: overallScore,
		ReportHash:   reportHash,
		RecordedAt:   timestamp,
	}

	recordJSON, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal record: %v", err)
	}

	err = ctx.GetStub().PutPrivateData(PlatformPrivateCollection, "VALIDATOR_DECISION_"+recordID, recordJSON)
	if err != nil {
		return fmt.Errorf("failed to record decision: %v", err)
	}

	return nil
}

// PublishGlobalMetrics publishes platform metrics to public ledger
func (p *PlatformContract) PublishGlobalMetrics(
	ctx contractapi.TransactionContextInterface,
	metricsID string,
	totalCampaigns int,
	activeCampaigns int,
	successfulCampaigns int,
	totalInvestorCount int,
) error {

	timestamp := time.Now().Format(time.RFC3339)

	metricsData := fmt.Sprintf("%s|%d|%d|%d|%d",
		metricsID, totalCampaigns, activeCampaigns, successfulCampaigns, totalInvestorCount)
	hash := sha256.Sum256([]byte(metricsData))
	metricsHash := hex.EncodeToString(hash[:])

	metrics := GlobalMetrics{
		MetricsID:           metricsID,
		TotalCampaigns:      totalCampaigns,
		ActiveCampaigns:     activeCampaigns,
		SuccessfulCampaigns: successfulCampaigns,
		TotalInvestorCount:  totalInvestorCount,
		MetricsHash:         metricsHash,
		PublishedAt:         timestamp,
	}

	metricsJSON, _ := json.Marshal(metrics)
	err := ctx.GetStub().PutState("GLOBAL_METRICS_"+metricsID, metricsJSON)
	if err != nil {
		return fmt.Errorf("failed to publish metrics: %v", err)
	}

	return nil
}

// ============================================================================
// WALLET & TOKEN MANAGEMENT - Using PDC
// ============================================================================

// CreateWallet creates a new wallet
func (p *PlatformContract) CreateWallet(
	ctx contractapi.TransactionContextInterface,
	walletID string,
	userID string,
	userType string,
	initialBalance float64,
) error {

	timestamp := time.Now().Format(time.RFC3339)

	wallet := Wallet{
		WalletID:       walletID,
		UserID:         userID,
		UserType:       userType,
		Balance:        initialBalance,
		LockedBalance:  0,
		TotalDeposited: initialBalance,
		TotalWithdrawn: 0,
		TotalPenalties: 0,
		CreatedAt:      timestamp,
		UpdatedAt:      timestamp,
	}

	walletJSON, err := json.Marshal(wallet)
	if err != nil {
		return fmt.Errorf("failed to marshal wallet: %v", err)
	}

	// Store in PlatformPrivateCollection
	err = ctx.GetStub().PutPrivateData(PlatformPrivateCollection, "WALLET_"+walletID, walletJSON)
	if err != nil {
		return fmt.Errorf("failed to create wallet: %v", err)
	}

	return nil
}

// DepositTokens deposits tokens to wallet
func (p *PlatformContract) DepositTokens(
	ctx contractapi.TransactionContextInterface,
	walletID string,
	amount float64,
	reference string,
) error {

	walletJSON, err := ctx.GetStub().GetPrivateData(PlatformPrivateCollection, "WALLET_"+walletID)
	if err != nil || walletJSON == nil {
		return fmt.Errorf("wallet not found: %v", err)
	}

	var wallet Wallet
	json.Unmarshal(walletJSON, &wallet)

	timestamp := time.Now().Format(time.RFC3339)
	wallet.Balance += amount
	wallet.TotalDeposited += amount
	wallet.UpdatedAt = timestamp

	walletJSON, _ = json.Marshal(wallet)
	err = ctx.GetStub().PutPrivateData(PlatformPrivateCollection, "WALLET_"+walletID, walletJSON)
	if err != nil {
		return fmt.Errorf("failed to deposit: %v", err)
	}

	return nil
}

// TransferTokens transfers tokens between wallets
func (p *PlatformContract) TransferTokens(
	ctx contractapi.TransactionContextInterface,
	transactionID string,
	fromWalletID string,
	toWalletID string,
	amount float64,
	transactionType string,
	reference string,
) error {

	// Get from wallet
	fromJSON, err := ctx.GetStub().GetPrivateData(PlatformPrivateCollection, "WALLET_"+fromWalletID)
	if err != nil || fromJSON == nil {
		return fmt.Errorf("from wallet not found: %v", err)
	}

	var fromWallet Wallet
	json.Unmarshal(fromJSON, &fromWallet)

	if fromWallet.Balance < amount {
		return fmt.Errorf("insufficient balance")
	}

	// Get to wallet
	toJSON, err := ctx.GetStub().GetPrivateData(PlatformPrivateCollection, "WALLET_"+toWalletID)
	if err != nil || toJSON == nil {
		return fmt.Errorf("to wallet not found: %v", err)
	}

	var toWallet Wallet
	json.Unmarshal(toJSON, &toWallet)

	timestamp := time.Now().Format(time.RFC3339)

	// Update balances
	fromWallet.Balance -= amount
	fromWallet.TotalWithdrawn += amount
	fromWallet.UpdatedAt = timestamp

	toWallet.Balance += amount
	toWallet.TotalDeposited += amount
	toWallet.UpdatedAt = timestamp

	// Save wallets
	fromJSON, _ = json.Marshal(fromWallet)
	ctx.GetStub().PutPrivateData(PlatformPrivateCollection, "WALLET_"+fromWalletID, fromJSON)

	toJSON, _ = json.Marshal(toWallet)
	ctx.GetStub().PutPrivateData(PlatformPrivateCollection, "WALLET_"+toWalletID, toJSON)

	// Record transaction
	transaction := TokenTransaction{
		TransactionID:   transactionID,
		FromWalletID:    fromWalletID,
		ToWalletID:      toWalletID,
		Amount:          amount,
		TransactionType: transactionType,
		Reference:       reference,
		Status:          "COMPLETED",
		CreatedAt:       timestamp,
	}

	txJSON, _ := json.Marshal(transaction)
	ctx.GetStub().PutPrivateData(PlatformPrivateCollection, "TRANSACTION_"+transactionID, txJSON)

	return nil
}

// SetExchangeRate sets token exchange rate
func (p *PlatformContract) SetExchangeRate(
	ctx contractapi.TransactionContextInterface,
	rateID string,
	currency string,
	tokenRate float64,
	effectiveAt string,
	expiresAt string,
) error {

	rate := TokenExchangeRate{
		RateID:      rateID,
		Currency:    currency,
		TokenRate:   tokenRate,
		EffectiveAt: effectiveAt,
		ExpiresAt:   expiresAt,
		SetBy:       "PLATFORM",
	}

	rateJSON, _ := json.Marshal(rate)
	err := ctx.GetStub().PutState("EXCHANGE_RATE_"+rateID, rateJSON)
	if err != nil {
		return fmt.Errorf("failed to set exchange rate: %v", err)
	}

	return nil
}

// GetWallet retrieves a wallet
func (p *PlatformContract) GetWallet(ctx contractapi.TransactionContextInterface, walletID string) (string, error) {
	walletJSON, err := ctx.GetStub().GetPrivateData(PlatformPrivateCollection, "WALLET_"+walletID)
	if err != nil || walletJSON == nil {
		return "", fmt.Errorf("wallet not found: %v", err)
	}

	return string(walletJSON), nil
}

// GetWalletByUser retrieves wallet by user ID and type
func (p *PlatformContract) GetWalletByUser(ctx contractapi.TransactionContextInterface, userType string, userID string) (string, error) {
	// Would use rich query in CouchDB
	return `{}`, nil
}

// ============================================================================
// FEE MANAGEMENT - Campaign Fees Using PDC
// ============================================================================

// SetCampaignFeeTier sets a campaign fee tier
func (p *PlatformContract) SetCampaignFeeTier(
	ctx contractapi.TransactionContextInterface,
	tierID string,
	minGoalAmount float64,
	maxGoalAmount float64,
	feePercentage float64,
	description string,
) error {

	tier := FeeTier{
		TierID:        tierID,
		MinGoalAmount: minGoalAmount,
		MaxGoalAmount: maxGoalAmount,
		FeePercentage: feePercentage,
		Description:   description,
	}

	tierJSON, _ := json.Marshal(tier)
	err := ctx.GetStub().PutState("CAMPAIGN_FEE_TIER_"+tierID, tierJSON)
	if err != nil {
		return fmt.Errorf("failed to set fee tier: %v", err)
	}

	return nil
}

// CollectCampaignFee collects fee from startup when campaign reaches goal
func (p *PlatformContract) CollectCampaignFee(
	ctx contractapi.TransactionContextInterface,
	collectionID string,
	campaignID string,
	startupID string,
	goalAmount float64,
	feePercentage float64,
) error {

	amount := goalAmount * (feePercentage / 100.0)
	timestamp := time.Now().Format(time.RFC3339)

	collection := FeeCollection{
		CollectionID:  collectionID,
		CampaignID:    campaignID,
		StartupID:     startupID,
		FeeType:       "CAMPAIGN",
		Amount:        amount,
		GoalAmount:    goalAmount,
		FeePercentage: feePercentage,
		Status:        "COLLECTED",
		CollectedAt:   timestamp,
	}

	collectionJSON, _ := json.Marshal(collection)
	err := ctx.GetStub().PutPrivateData(PlatformPrivateCollection, "FEE_COLLECTION_"+collectionID, collectionJSON)
	if err != nil {
		return fmt.Errorf("failed to collect fee: %v", err)
	}

	return nil
}

// GetCampaignFeeTiers retrieves all campaign fee tiers
func (p *PlatformContract) GetCampaignFeeTiers(ctx contractapi.TransactionContextInterface) (string, error) {
	// Would use rich query in CouchDB
	return `[]`, nil
}

// ============================================================================
// FEE MANAGEMENT - Dispute Fees Using PDC
// ============================================================================

// SetDisputeFeeTier sets a dispute fee tier
func (p *PlatformContract) SetDisputeFeeTier(
	ctx contractapi.TransactionContextInterface,
	tierID string,
	minClaimAmount float64,
	maxClaimAmount float64,
	feeAmount float64,
	description string,
) error {

	tier := DisputeFeeTier{
		TierID:         tierID,
		MinClaimAmount: minClaimAmount,
		MaxClaimAmount: maxClaimAmount,
		FeeAmount:      feeAmount,
		Description:    description,
	}

	tierJSON, _ := json.Marshal(tier)
	err := ctx.GetStub().PutState("DISPUTE_FEE_TIER_"+tierID, tierJSON)
	if err != nil {
		return fmt.Errorf("failed to set dispute fee tier: %v", err)
	}

	return nil
}

// CollectDisputeFee collects fee from dispute initiator
func (p *PlatformContract) CollectDisputeFee(
	ctx contractapi.TransactionContextInterface,
	feeRecordID string,
	disputeID string,
	initiatorID string,
	claimAmount float64,
	feeAmount float64,
) error {

	timestamp := time.Now().Format(time.RFC3339)

	record := DisputeFeeRecord{
		FeeRecordID: feeRecordID,
		DisputeID:   disputeID,
		InitiatorID: initiatorID,
		ClaimAmount: claimAmount,
		FeeAmount:   feeAmount,
		Status:      "COLLECTED",
		CollectedAt: timestamp,
	}

	recordJSON, _ := json.Marshal(record)
	err := ctx.GetStub().PutPrivateData(PlatformPrivateCollection, "DISPUTE_FEE_"+feeRecordID, recordJSON)
	if err != nil {
		return fmt.Errorf("failed to collect dispute fee: %v", err)
	}

	return nil
}

// RefundDisputeFee refunds dispute fee if initiator wins
func (p *PlatformContract) RefundDisputeFee(
	ctx contractapi.TransactionContextInterface,
	feeRecordID string,
	disputeID string,
) error {

	recordJSON, err := ctx.GetStub().GetPrivateData(PlatformPrivateCollection, "DISPUTE_FEE_"+feeRecordID)
	if err != nil || recordJSON == nil {
		return fmt.Errorf("fee record not found: %v", err)
	}

	var record DisputeFeeRecord
	json.Unmarshal(recordJSON, &record)

	timestamp := time.Now().Format(time.RFC3339)
	record.Status = "REFUNDED"
	record.RefundAmount = record.FeeAmount
	record.ProcessedAt = timestamp

	recordJSON, _ = json.Marshal(record)
	ctx.GetStub().PutPrivateData(PlatformPrivateCollection, "DISPUTE_FEE_"+feeRecordID, recordJSON)

	return nil
}

// ForfeitDisputeFee forfeits dispute fee if initiator loses
func (p *PlatformContract) ForfeitDisputeFee(
	ctx contractapi.TransactionContextInterface,
	feeRecordID string,
	disputeID string,
) error {

	recordJSON, err := ctx.GetStub().GetPrivateData(PlatformPrivateCollection, "DISPUTE_FEE_"+feeRecordID)
	if err != nil || recordJSON == nil {
		return fmt.Errorf("fee record not found: %v", err)
	}

	var record DisputeFeeRecord
	json.Unmarshal(recordJSON, &record)

	timestamp := time.Now().Format(time.RFC3339)
	record.Status = "FORFEITED"
	record.RefundAmount = 0
	record.ProcessedAt = timestamp

	recordJSON, _ = json.Marshal(record)
	ctx.GetStub().PutPrivateData(PlatformPrivateCollection, "DISPUTE_FEE_"+feeRecordID, recordJSON)

	return nil
}

// GetDisputeFeeTiers retrieves all dispute fee tiers
func (p *PlatformContract) GetDisputeFeeTiers(ctx contractapi.TransactionContextInterface) (string, error) {
	// Would use rich query in CouchDB
	return `[]`, nil
}

// ProcessDisputeFeeOutcome processes fee based on dispute outcome
func (p *PlatformContract) ProcessDisputeFeeOutcome(
	ctx contractapi.TransactionContextInterface,
	feeRecordID string,
	disputeID string,
	initiatorWon bool,
) error {

	if initiatorWon {
		return p.RefundDisputeFee(ctx, feeRecordID, disputeID)
	}
	return p.ForfeitDisputeFee(ctx, feeRecordID, disputeID)
}

// GetDisputeFeeRecord retrieves dispute fee record
func (p *PlatformContract) GetDisputeFeeRecord(
	ctx contractapi.TransactionContextInterface,
	feeRecordID string,
) (string, error) {
	recordJSON, err := ctx.GetStub().GetPrivateData(PlatformPrivateCollection, "DISPUTE_FEE_"+feeRecordID)
	if err != nil || recordJSON == nil {
		return "", fmt.Errorf("fee record not found: %v", err)
	}

	return string(recordJSON), nil
}

// ============================================================================
// RATING & REPUTATION MANAGEMENT - Using PDC
// ============================================================================

// RecordRating records a rating
func (p *PlatformContract) RecordRating(
	ctx contractapi.TransactionContextInterface,
	ratingID string,
	ratedUserType string,
	ratedUserId string,
	raterUserType string,
	raterUserId string,
	context string,
	contextId string,
	rating float64,
	comment string,
) error {

	timestamp := time.Now().Format(time.RFC3339)

	record := RatingRecord{
		RatingID:      ratingID,
		RatedUserType: ratedUserType,
		RatedUserId:   ratedUserId,
		RaterUserType: raterUserType,
		RaterUserId:   raterUserId,
		Context:       context,
		ContextId:     contextId,
		Rating:        rating,
		Comment:       comment,
		CreatedAt:     timestamp,
	}

	recordJSON, _ := json.Marshal(record)
	err := ctx.GetStub().PutState("RATING_"+ratingID, recordJSON)
	if err != nil {
		return fmt.Errorf("failed to record rating: %v", err)
	}

	// Update aggregate
	p.updateRatingAggregate(ctx, ratedUserType, ratedUserId, rating)
	p.updateReputationScore(ctx, ratedUserType, ratedUserId)

	return nil
}

// updateRatingAggregate updates aggregated ratings
func (p *PlatformContract) updateRatingAggregate(
	ctx contractapi.TransactionContextInterface,
	userType string,
	userId string,
	newRating float64,
) error {

	key := fmt.Sprintf("RATING_AGGREGATE_%s_%s", userType, userId)
	aggregateJSON, err := ctx.GetStub().GetState(key)

	var aggregate RatingAggregate
	timestamp := time.Now().Format(time.RFC3339)

	if err == nil && aggregateJSON != nil {
		json.Unmarshal(aggregateJSON, &aggregate)
		totalRatingsFloat := float64(aggregate.TotalRatings)
		aggregate.AverageRating = ((aggregate.AverageRating * totalRatingsFloat) + newRating) / (totalRatingsFloat + 1)
		aggregate.TotalRatings++
		aggregate.UpdatedAt = timestamp
	} else {
		aggregate = RatingAggregate{
			UserType:      userType,
			UserId:        userId,
			TotalRatings:  1,
			AverageRating: newRating,
			UpdatedAt:     timestamp,
		}
	}

	aggregateJSON, _ = json.Marshal(aggregate)
	ctx.GetStub().PutState(key, aggregateJSON)

	return nil
}

// updateReputationScore updates user reputation
func (p *PlatformContract) updateReputationScore(
	ctx contractapi.TransactionContextInterface,
	userType string,
	userId string,
) error {

	key := fmt.Sprintf("REPUTATION_%s_%s", userType, userId)
	reputationJSON, err := ctx.GetStub().GetState(key)

	var reputation ReputationScore
	timestamp := time.Now().Format(time.RFC3339)

	if err == nil && reputationJSON != nil {
		json.Unmarshal(reputationJSON, &reputation)
	} else {
		reputation = ReputationScore{
			UserType:        userType,
			UserId:          userId,
			ReputationScore: 0,
			DisputesWon:     0,
			DisputesLost:    0,
			TotalDisputes:   0,
			SuccessfulDeals: 0,
			Status:          "ACTIVE",
			UpdatedAt:       timestamp,
		}
	}

	// Get rating aggregate
	aggKey := fmt.Sprintf("RATING_AGGREGATE_%s_%s", userType, userId)
	aggregateJSON, err := ctx.GetStub().GetState(aggKey)
	if err == nil && aggregateJSON != nil {
		var aggregate RatingAggregate
		json.Unmarshal(aggregateJSON, &aggregate)
		reputation.ReputationScore = aggregate.AverageRating * 20 // Convert to 0-100 scale
	}

	reputation.UpdatedAt = timestamp

	reputationJSON, _ = json.Marshal(reputation)
	ctx.GetStub().PutState(key, reputationJSON)

	return nil
}

// GetRatingAggregate retrieves rating aggregate for a user
func (p *PlatformContract) GetRatingAggregate(
	ctx contractapi.TransactionContextInterface,
	userType string,
	userId string,
) (string, error) {

	key := fmt.Sprintf("RATING_AGGREGATE_%s_%s", userType, userId)
	aggregateJSON, err := ctx.GetStub().GetState(key)
	if err != nil || aggregateJSON == nil {
		return "", fmt.Errorf("rating aggregate not found: %v", err)
	}

	return string(aggregateJSON), nil
}

// GetReputationScore retrieves reputation score for a user
func (p *PlatformContract) GetReputationScore(
	ctx contractapi.TransactionContextInterface,
	userType string,
	userId string,
) (string, error) {

	key := fmt.Sprintf("REPUTATION_%s_%s", userType, userId)
	reputationJSON, err := ctx.GetStub().GetState(key)
	if err != nil || reputationJSON == nil {
		return "", fmt.Errorf("reputation not found: %v", err)
	}

	return string(reputationJSON), nil
}

// ============================================================================
// DISPUTE MANAGEMENT - Using AllOrgsCollection for transparency
// ============================================================================

// CreateDispute creates a new dispute (stored in AllOrgsCollection for transparency)
func (p *PlatformContract) CreateDispute(
	ctx contractapi.TransactionContextInterface,
	disputeID string,
	initiatorType string,
	initiatorID string,
	respondentType string,
	respondentID string,
	disputeType string,
	campaignID string,
	agreementID string,
	title string,
	description string,
	claimAmount float64,
) error {

	timestamp := time.Now().Format(time.RFC3339)

	dispute := Dispute{
		DisputeID:          disputeID,
		InitiatorType:      initiatorType,
		InitiatorID:        initiatorID,
		RespondentType:     respondentType,
		RespondentID:       respondentID,
		DisputeType:        disputeType,
		CampaignID:         campaignID,
		AgreementID:        agreementID,
		Title:              title,
		Description:        description,
		ClaimAmount:        claimAmount,
		EvidenceHashes:     []string{},
		Status:             "CREATED",
		InvestigatorID:     "",
		InvestigationNotes: []string{},
		VotingEnabled:      false,
		EligibleVoters:     []string{},
		Votes:              []VoteCommitment{},
		Resolution:         "",
		ResolutionDetails:  "",
		CreatedAt:          timestamp,
	}

	disputeJSON, _ := json.Marshal(dispute)
	err := ctx.GetStub().PutPrivateData(AllOrgsCollection, "DISPUTE_"+disputeID, disputeJSON)
	if err != nil {
		return fmt.Errorf("failed to create dispute: %v", err)
	}

	return nil
}

// SubmitEvidence submits evidence for a dispute
func (p *PlatformContract) SubmitEvidence(
	ctx contractapi.TransactionContextInterface,
	disputeID string,
	evidenceHash string,
) error {

	disputeJSON, err := ctx.GetStub().GetPrivateData(AllOrgsCollection, "DISPUTE_"+disputeID)
	if err != nil || disputeJSON == nil {
		return fmt.Errorf("dispute not found: %v", err)
	}

	var dispute Dispute
	json.Unmarshal(disputeJSON, &dispute)

	dispute.EvidenceHashes = append(dispute.EvidenceHashes, evidenceHash)

	disputeJSON, _ = json.Marshal(dispute)
	ctx.GetStub().PutPrivateData(AllOrgsCollection, "DISPUTE_"+disputeID, disputeJSON)

	return nil
}

// AssignInvestigator assigns validator to investigate dispute
func (p *PlatformContract) AssignInvestigator(
	ctx contractapi.TransactionContextInterface,
	disputeID string,
	investigatorID string,
) error {

	disputeJSON, err := ctx.GetStub().GetPrivateData(AllOrgsCollection, "DISPUTE_"+disputeID)
	if err != nil || disputeJSON == nil {
		return fmt.Errorf("dispute not found: %v", err)
	}

	var dispute Dispute
	json.Unmarshal(disputeJSON, &dispute)

	dispute.InvestigatorID = investigatorID
	dispute.Status = "UNDER_INVESTIGATION"

	disputeJSON, _ = json.Marshal(dispute)
	ctx.GetStub().PutPrivateData(AllOrgsCollection, "DISPUTE_"+disputeID, disputeJSON)

	return nil
}

// AddInvestigationNote adds investigation note
func (p *PlatformContract) AddInvestigationNote(
	ctx contractapi.TransactionContextInterface,
	disputeID string,
	note string,
) error {

	disputeJSON, err := ctx.GetStub().GetPrivateData(AllOrgsCollection, "DISPUTE_"+disputeID)
	if err != nil || disputeJSON == nil {
		return fmt.Errorf("dispute not found: %v", err)
	}

	var dispute Dispute
	json.Unmarshal(disputeJSON, &dispute)

	timestamp := time.Now().Format(time.RFC3339)
	noteEntry := fmt.Sprintf("[%s] %s", timestamp, note)
	dispute.InvestigationNotes = append(dispute.InvestigationNotes, noteEntry)

	disputeJSON, _ = json.Marshal(dispute)
	ctx.GetStub().PutPrivateData(AllOrgsCollection, "DISPUTE_"+disputeID, disputeJSON)

	return nil
}

// EnableVoting enables voting for dispute
func (p *PlatformContract) EnableVoting(
	ctx contractapi.TransactionContextInterface,
	disputeID string,
	eligibleVotersJSON string,
) error {

	var eligibleVoters []string
	json.Unmarshal([]byte(eligibleVotersJSON), &eligibleVoters)

	disputeJSON, err := ctx.GetStub().GetPrivateData(AllOrgsCollection, "DISPUTE_"+disputeID)
	if err != nil || disputeJSON == nil {
		return fmt.Errorf("dispute not found: %v", err)
	}

	var dispute Dispute
	json.Unmarshal(disputeJSON, &dispute)

	dispute.VotingEnabled = true
	dispute.EligibleVoters = eligibleVoters
	dispute.Status = "VOTING_OPEN"

	disputeJSON, _ = json.Marshal(dispute)
	ctx.GetStub().PutPrivateData(AllOrgsCollection, "DISPUTE_"+disputeID, disputeJSON)

	return nil
}

// CommitVote commits a vote (commit phase of commit-reveal voting)
func (p *PlatformContract) CommitVote(
	ctx contractapi.TransactionContextInterface,
	disputeID string,
	voterID string,
	voteHash string,
) error {

	disputeJSON, err := ctx.GetStub().GetPrivateData(AllOrgsCollection, "DISPUTE_"+disputeID)
	if err != nil || disputeJSON == nil {
		return fmt.Errorf("dispute not found: %v", err)
	}

	var dispute Dispute
	json.Unmarshal(disputeJSON, &dispute)

	if !dispute.VotingEnabled {
		return fmt.Errorf("voting not enabled")
	}

	// Check if voter is eligible
	eligible := false
	for _, voter := range dispute.EligibleVoters {
		if voter == voterID {
			eligible = true
			break
		}
	}
	if !eligible {
		return fmt.Errorf("voter not eligible")
	}

	timestamp := time.Now().Format(time.RFC3339)

	commitment := VoteCommitment{
		VoterID:     voterID,
		VoteHash:    voteHash,
		Revealed:    false,
		Vote:        "",
		CommittedAt: timestamp,
	}

	dispute.Votes = append(dispute.Votes, commitment)

	disputeJSON, _ = json.Marshal(dispute)
	ctx.GetStub().PutPrivateData(AllOrgsCollection, "DISPUTE_"+disputeID, disputeJSON)

	return nil
}

// RevealVote reveals a vote (reveal phase of commit-reveal voting)
func (p *PlatformContract) RevealVote(
	ctx contractapi.TransactionContextInterface,
	disputeID string,
	voterID string,
	vote string,
	salt string,
) error {

	disputeJSON, err := ctx.GetStub().GetPrivateData(AllOrgsCollection, "DISPUTE_"+disputeID)
	if err != nil || disputeJSON == nil {
		return fmt.Errorf("dispute not found: %v", err)
	}

	var dispute Dispute
	json.Unmarshal(disputeJSON, &dispute)

	// Find vote commitment
	voteIndex := -1
	for i, v := range dispute.Votes {
		if v.VoterID == voterID {
			voteIndex = i
			break
		}
	}

	if voteIndex == -1 {
		return fmt.Errorf("vote commitment not found")
	}

	// Verify hash
	voteData := vote + salt
	hash := sha256.Sum256([]byte(voteData))
	voteHash := hex.EncodeToString(hash[:])

	if voteHash != dispute.Votes[voteIndex].VoteHash {
		return fmt.Errorf("vote hash mismatch")
	}

	timestamp := time.Now().Format(time.RFC3339)
	dispute.Votes[voteIndex].Revealed = true
	dispute.Votes[voteIndex].Vote = vote
	dispute.Votes[voteIndex].RevealedAt = timestamp

	disputeJSON, _ = json.Marshal(dispute)
	ctx.GetStub().PutPrivateData(AllOrgsCollection, "DISPUTE_"+disputeID, disputeJSON)

	return nil
}

// TallyVotes tallies votes and determines outcome
func (p *PlatformContract) TallyVotes(
	ctx contractapi.TransactionContextInterface,
	disputeID string,
) (string, error) {

	disputeJSON, err := ctx.GetStub().GetPrivateData(AllOrgsCollection, "DISPUTE_"+disputeID)
	if err != nil || disputeJSON == nil {
		return "", fmt.Errorf("dispute not found: %v", err)
	}

	var dispute Dispute
	json.Unmarshal(disputeJSON, &dispute)

	initiatorSupport := 0
	respondentSupport := 0

	for _, vote := range dispute.Votes {
		if !vote.Revealed {
			continue
		}
		if vote.Vote == "SUPPORT_INITIATOR" {
			initiatorSupport++
		} else if vote.Vote == "SUPPORT_RESPONDENT" {
			respondentSupport++
		}
	}

	outcome := "TIED"
	if initiatorSupport > respondentSupport {
		outcome = "INITIATOR_WINS"
	} else if respondentSupport > initiatorSupport {
		outcome = "RESPONDENT_WINS"
	}

	dispute.Status = "VOTING_CLOSED"
	dispute.Resolution = outcome

	disputeJSON, _ = json.Marshal(dispute)
	ctx.GetStub().PutPrivateData(AllOrgsCollection, "DISPUTE_"+disputeID, disputeJSON)

	return outcome, nil
}

// ResolveDispute resolves dispute with final decision
func (p *PlatformContract) ResolveDispute(
	ctx contractapi.TransactionContextInterface,
	disputeID string,
	resolution string,
	resolutionDetails string,
) error {

	disputeJSON, err := ctx.GetStub().GetPrivateData(AllOrgsCollection, "DISPUTE_"+disputeID)
	if err != nil || disputeJSON == nil {
		return fmt.Errorf("dispute not found: %v", err)
	}

	var dispute Dispute
	json.Unmarshal(disputeJSON, &dispute)

	timestamp := time.Now().Format(time.RFC3339)
	dispute.Status = "RESOLVED"
	dispute.Resolution = resolution
	dispute.ResolutionDetails = resolutionDetails
	dispute.ResolvedAt = timestamp

	disputeJSON, _ = json.Marshal(dispute)
	ctx.GetStub().PutPrivateData(AllOrgsCollection, "DISPUTE_"+disputeID, disputeJSON)

	// Update reputations
	if resolution == "INITIATOR_WINS" {
		p.incrementDisputeWin(ctx, dispute.InitiatorType, dispute.InitiatorID)
		p.incrementDisputeLoss(ctx, dispute.RespondentType, dispute.RespondentID)
	} else if resolution == "RESPONDENT_WINS" {
		p.incrementDisputeWin(ctx, dispute.RespondentType, dispute.RespondentID)
		p.incrementDisputeLoss(ctx, dispute.InitiatorType, dispute.InitiatorID)
	}

	return nil
}

// incrementDisputeWin increments dispute win count
func (p *PlatformContract) incrementDisputeWin(
	ctx contractapi.TransactionContextInterface,
	userType string,
	userId string,
) error {

	key := fmt.Sprintf("REPUTATION_%s_%s", userType, userId)
	reputationJSON, err := ctx.GetStub().GetState(key)

	var reputation ReputationScore
	timestamp := time.Now().Format(time.RFC3339)

	if err == nil && reputationJSON != nil {
		json.Unmarshal(reputationJSON, &reputation)
	} else {
		reputation = ReputationScore{
			UserType:        userType,
			UserId:          userId,
			ReputationScore: 0,
			DisputesWon:     0,
			DisputesLost:    0,
			TotalDisputes:   0,
			SuccessfulDeals: 0,
			Status:          "ACTIVE",
		}
	}

	reputation.DisputesWon++
	reputation.TotalDisputes++
	reputation.UpdatedAt = timestamp

	reputationJSON, _ = json.Marshal(reputation)
	ctx.GetStub().PutState(key, reputationJSON)

	return nil
}

// incrementDisputeLoss increments dispute loss count
func (p *PlatformContract) incrementDisputeLoss(
	ctx contractapi.TransactionContextInterface,
	userType string,
	userId string,
) error {

	key := fmt.Sprintf("REPUTATION_%s_%s", userType, userId)
	reputationJSON, err := ctx.GetStub().GetState(key)

	var reputation ReputationScore
	timestamp := time.Now().Format(time.RFC3339)

	if err == nil && reputationJSON != nil {
		json.Unmarshal(reputationJSON, &reputation)
	} else {
		reputation = ReputationScore{
			UserType:        userType,
			UserId:          userId,
			ReputationScore: 0,
			DisputesWon:     0,
			DisputesLost:    0,
			TotalDisputes:   0,
			SuccessfulDeals: 0,
			Status:          "ACTIVE",
		}
	}

	reputation.DisputesLost++
	reputation.TotalDisputes++
	reputation.UpdatedAt = timestamp

	reputationJSON, _ = json.Marshal(reputation)
	ctx.GetStub().PutState(key, reputationJSON)

	return nil
}

// ApplyPenalty applies penalty to user
func (p *PlatformContract) ApplyPenalty(
	ctx contractapi.TransactionContextInterface,
	penaltyID string,
	userType string,
	userID string,
	disputeID string,
	penaltyType string,
	amount float64,
) error {

	timestamp := time.Now().Format(time.RFC3339)

	penalty := Penalty{
		PenaltyID:   penaltyID,
		UserType:    userType,
		UserID:      userID,
		DisputeID:   disputeID,
		PenaltyType: penaltyType,
		Amount:      amount,
		AppliedAt:   timestamp,
	}

	penaltyJSON, _ := json.Marshal(penalty)
	err := ctx.GetStub().PutState("PENALTY_"+penaltyID, penaltyJSON)
	if err != nil {
		return fmt.Errorf("failed to apply penalty: %v", err)
	}

	// Update reputation
	key := fmt.Sprintf("REPUTATION_%s_%s", userType, userID)
	reputationJSON, err := ctx.GetStub().GetState(key)
	if err == nil && reputationJSON != nil {
		var reputation ReputationScore
		json.Unmarshal(reputationJSON, &reputation)
		reputation.ReputationScore -= 10 // Penalty reduces reputation
		if reputation.ReputationScore < 0 {
			reputation.ReputationScore = 0
		}
		reputation.UpdatedAt = timestamp
		reputationJSON, _ = json.Marshal(reputation)
		ctx.GetStub().PutState(key, reputationJSON)
	}

	return nil
}

// ProcessRefund processes refund order
func (p *PlatformContract) ProcessRefund(
	ctx contractapi.TransactionContextInterface,
	refundID string,
	disputeID string,
	recipient string,
	amount float64,
) error {

	timestamp := time.Now().Format(time.RFC3339)

	refund := RefundOrder{
		RefundID:    refundID,
		DisputeID:   disputeID,
		Recipient:   recipient,
		Amount:      amount,
		ProcessedAt: timestamp,
	}

	refundJSON, _ := json.Marshal(refund)
	err := ctx.GetStub().PutPrivateData(InvestorPlatformCollection, "REFUND_"+refundID, refundJSON)
	if err != nil {
		return fmt.Errorf("failed to process refund: %v", err)
	}

	return nil
}

// GetDispute retrieves dispute details
func (p *PlatformContract) GetDispute(
	ctx contractapi.TransactionContextInterface,
	disputeID string,
) (string, error) {

	disputeJSON, err := ctx.GetStub().GetPrivateData(AllOrgsCollection, "DISPUTE_"+disputeID)
	if err != nil || disputeJSON == nil {
		return "", fmt.Errorf("dispute not found: %v", err)
	}

	return string(disputeJSON), nil
}

// GetAllDisputes retrieves all disputes
func (p *PlatformContract) GetAllDisputes(ctx contractapi.TransactionContextInterface) (string, error) {
	// Would use rich query in CouchDB to get all disputes from AllOrgsCollection
	return `[]`, nil
}

// ============================================================================
// QUERY FUNCTIONS
// ============================================================================

// GetPublishedCampaign retrieves published campaign
func (p *PlatformContract) GetPublishedCampaign(
	ctx contractapi.TransactionContextInterface,
	campaignID string,
) (string, error) {

	campaignJSON, err := ctx.GetStub().GetPrivateData(PlatformPrivateCollection, "PUBLISHED_CAMPAIGN_"+campaignID)
	if err != nil || campaignJSON == nil {
		return "", fmt.Errorf("campaign not found: %v", err)
	}

	return string(campaignJSON), nil
}

// GetSharedCampaign retrieves campaign shared by Startup (before publishing)
func (p *PlatformContract) GetSharedCampaign(
	ctx contractapi.TransactionContextInterface,
	campaignID string,
) (string, error) {

	// Check if already published first - if so, return the published version (status: PUBLISHED)
	publishedJSON, err := ctx.GetStub().GetPrivateData(PlatformPrivateCollection, "PUBLISHED_CAMPAIGN_"+campaignID)
	if err == nil && publishedJSON != nil {
		return string(publishedJSON), nil
	}

	// If not published, get from shared collection (status: PENDING_PLATFORM_APPROVAL)
	campaignJSON, err := ctx.GetStub().GetPrivateData(StartupPlatformCollection, "CAMPAIGN_SHARE_"+campaignID)
	if err != nil || campaignJSON == nil {
		return "", fmt.Errorf("shared campaign not found: %v", err)
	}

	return string(campaignJSON), nil
}

// GetAllSharedCampaigns retrieves all shared campaigns (pending + published)
func (p *PlatformContract) GetAllSharedCampaigns(ctx contractapi.TransactionContextInterface) ([]map[string]interface{}, error) {
	var allCampaigns []map[string]interface{}

	// 1. Get pending shared campaigns from StartupPlatformCollection
	pendingIterator, err := ctx.GetStub().GetPrivateDataByRange(StartupPlatformCollection, "CAMPAIGN_SHARE_", "CAMPAIGN_SHARE_~")
	if err == nil {
		defer pendingIterator.Close()
		for pendingIterator.HasNext() {
			queryResponse, err := pendingIterator.Next()
			if err != nil {
				continue
			}

			var campaignMap map[string]interface{}
			err = json.Unmarshal(queryResponse.Value, &campaignMap)
			if err != nil {
				continue
			}
			campaignMap["status"] = "PENDING_PLATFORM_APPROVAL"
			allCampaigns = append(allCampaigns, campaignMap)
		}
	}

	// 2. Get published campaigns from PlatformPrivateCollection
	publishedIterator, err := ctx.GetStub().GetPrivateDataByRange(PlatformPrivateCollection, "PUBLISHED_CAMPAIGN_", "PUBLISHED_CAMPAIGN_~")
	if err == nil {
		defer publishedIterator.Close()
		for publishedIterator.HasNext() {
			queryResponse, err := publishedIterator.Next()
			if err != nil {
				continue
			}

			var campaignMap map[string]interface{}
			err = json.Unmarshal(queryResponse.Value, &campaignMap)
			if err != nil {
				continue
			}
			campaignMap["status"] = "PUBLISHED"
			allCampaigns = append(allCampaigns, campaignMap)
		}
	}

	if allCampaigns == nil {
		allCampaigns = []map[string]interface{}{}
	}

	return allCampaigns, nil
}

// GetAgreement retrieves agreement from ThreePartyCollection
func (p *PlatformContract) GetAgreement(
	ctx contractapi.TransactionContextInterface,
	agreementID string,
) (string, error) {

	agreementJSON, err := ctx.GetStub().GetPrivateData(ThreePartyCollection, "AGREEMENT_"+agreementID)
	if err != nil || agreementJSON == nil {
		return "", fmt.Errorf("agreement not found: %v", err)
	}

	return string(agreementJSON), nil
}

// GetEscrow retrieves escrow details
func (p *PlatformContract) GetEscrow(
	ctx contractapi.TransactionContextInterface,
	escrowID string,
) (string, error) {

	escrowJSON, err := ctx.GetStub().GetPrivateData(PlatformPrivateCollection, "ESCROW_"+escrowID)
	if err != nil || escrowJSON == nil {
		return "", fmt.Errorf("escrow not found: %v", err)
	}

	return string(escrowJSON), nil
}

// ============================================================================
// TOKEN INTEGRATION FUNCTIONS (CFT/CFRT)
// ============================================================================

// CollectPublishingFee collects the publishing fee from startup in CFT
// Publishing fee: 2,500 CFT (₹1,000 at 1 INR = 2.5 CFT)
func (p *PlatformContract) CollectPublishingFee(
	ctx contractapi.TransactionContextInterface,
	feeID string,
	campaignID string,
	startupID string,
) error {

	timestamp := time.Now().Format(time.RFC3339)
	publishingFeeCFT := 2500.0 // ₹1,000 × 2.5

	// Record fee collection
	feeRecord := FeeCollection{
		CollectionID:  feeID,
		CampaignID:    campaignID,
		StartupID:     startupID,
		FeeType:       "PUBLISHING_FEE",
		Amount:        publishingFeeCFT,
		GoalAmount:    0,
		FeePercentage: 0,
		Status:        "PENDING_TOKEN_TRANSFER",
		CollectedAt:   timestamp,
	}

	feeJSON, err := json.Marshal(feeRecord)
	if err != nil {
		return fmt.Errorf("failed to marshal fee record: %v", err)
	}

	err = ctx.GetStub().PutPrivateData(PlatformPrivateCollection, "FEE_COLLECTION_"+feeID, feeJSON)
	if err != nil {
		return fmt.Errorf("failed to store fee record: %v", err)
	}

	// Note: Actual CFT transfer is done via TokenContract:TransferTokens
	// This function records the fee obligation
	return nil
}

// CollectRegistrationFee records startup registration fee obligation
// Registration fee: 250 CFT (₹100 at 1 INR = 2.5 CFT)
func (p *PlatformContract) CollectRegistrationFee(
	ctx contractapi.TransactionContextInterface,
	feeID string,
	startupID string,
) error {

	timestamp := time.Now().Format(time.RFC3339)
	registrationFeeCFT := 250.0 // ₹100 × 2.5

	feeRecord := FeeCollection{
		CollectionID:  feeID,
		CampaignID:    "",
		StartupID:     startupID,
		FeeType:       "REGISTRATION_FEE",
		Amount:        registrationFeeCFT,
		GoalAmount:    0,
		FeePercentage: 0,
		Status:        "PENDING_TOKEN_TRANSFER",
		CollectedAt:   timestamp,
	}

	feeJSON, err := json.Marshal(feeRecord)
	if err != nil {
		return fmt.Errorf("failed to marshal fee record: %v", err)
	}

	err = ctx.GetStub().PutPrivateData(PlatformPrivateCollection, "FEE_COLLECTION_"+feeID, feeJSON)
	if err != nil {
		return fmt.Errorf("failed to store fee record: %v", err)
	}

	return nil
}

// CollectCampaignCreationFee records campaign creation fee obligation
// Creation fee: 1,250 CFT (₹500 at 1 INR = 2.5 CFT)
func (p *PlatformContract) CollectCampaignCreationFee(
	ctx contractapi.TransactionContextInterface,
	feeID string,
	campaignID string,
	startupID string,
) error {

	timestamp := time.Now().Format(time.RFC3339)
	creationFeeCFT := 1250.0 // ₹500 × 2.5

	feeRecord := FeeCollection{
		CollectionID:  feeID,
		CampaignID:    campaignID,
		StartupID:     startupID,
		FeeType:       "CAMPAIGN_CREATION_FEE",
		Amount:        creationFeeCFT,
		GoalAmount:    0,
		FeePercentage: 0,
		Status:        "PENDING_TOKEN_TRANSFER",
		CollectedAt:   timestamp,
	}

	feeJSON, err := json.Marshal(feeRecord)
	if err != nil {
		return fmt.Errorf("failed to marshal fee record: %v", err)
	}

	err = ctx.GetStub().PutPrivateData(PlatformPrivateCollection, "FEE_COLLECTION_"+feeID, feeJSON)
	if err != nil {
		return fmt.Errorf("failed to store fee record: %v", err)
	}

	return nil
}

// ConfirmFeePayment confirms that a fee has been paid via token transfer
func (p *PlatformContract) ConfirmFeePayment(
	ctx contractapi.TransactionContextInterface,
	feeID string,
	transactionID string,
) error {

	feeJSON, err := ctx.GetStub().GetPrivateData(PlatformPrivateCollection, "FEE_COLLECTION_"+feeID)
	if err != nil || feeJSON == nil {
		return fmt.Errorf("fee record not found: %v", err)
	}

	var feeRecord FeeCollection
	err = json.Unmarshal(feeJSON, &feeRecord)
	if err != nil {
		return fmt.Errorf("failed to unmarshal fee record: %v", err)
	}

	timestamp := time.Now().Format(time.RFC3339)
	feeRecord.Status = "PAID"
	feeRecord.CollectedAt = timestamp

	feeJSON, _ = json.Marshal(feeRecord)
	err = ctx.GetStub().PutPrivateData(PlatformPrivateCollection, "FEE_COLLECTION_"+feeID, feeJSON)
	if err != nil {
		return fmt.Errorf("failed to update fee record: %v", err)
	}

	// Record the transaction link
	txLink := map[string]interface{}{
		"feeId":         feeID,
		"transactionId": transactionID,
		"confirmedAt":   timestamp,
	}
	txLinkJSON, _ := json.Marshal(txLink)
	ctx.GetStub().PutPrivateData(PlatformPrivateCollection, "FEE_TX_"+feeID, txLinkJSON)

	return nil
}

// CollectInvestmentFee calculates and records 5% investment fee
func (p *PlatformContract) CollectInvestmentFee(
	ctx contractapi.TransactionContextInterface,
	feeID string,
	campaignID string,
	investorID string,
	investmentAmountCFT float64,
) error {

	timestamp := time.Now().Format(time.RFC3339)
	feePercentage := 5.0
	feeAmountCFT := investmentAmountCFT * (feePercentage / 100)

	feeRecord := FeeCollection{
		CollectionID:  feeID,
		CampaignID:    campaignID,
		StartupID:     investorID, // Using StartupID field for investor
		FeeType:       "INVESTMENT_FEE",
		Amount:        feeAmountCFT,
		GoalAmount:    investmentAmountCFT,
		FeePercentage: feePercentage,
		Status:        "PENDING_TOKEN_TRANSFER",
		CollectedAt:   timestamp,
	}

	feeJSON, err := json.Marshal(feeRecord)
	if err != nil {
		return fmt.Errorf("failed to marshal fee record: %v", err)
	}

	err = ctx.GetStub().PutPrivateData(PlatformPrivateCollection, "FEE_COLLECTION_"+feeID, feeJSON)
	if err != nil {
		return fmt.Errorf("failed to store fee record: %v", err)
	}

	return nil
}

// GetFeeSchedule returns the current fee schedule in CFT
func (p *PlatformContract) GetFeeSchedule(
	ctx contractapi.TransactionContextInterface,
) (string, error) {

	// Fee schedule based on 1 INR = 2.5 CFT
	feeSchedule := map[string]interface{}{
		"exchangeRate": map[string]float64{
			"INR": 2.5,
			"USD": 83.0,
		},
		"fees": map[string]interface{}{
			"registrationFee": map[string]interface{}{
				"amountCFT": 250,
				"amountINR": 100,
				"type":      "FIXED",
			},
			"campaignCreationFee": map[string]interface{}{
				"amountCFT": 1250,
				"amountINR": 500,
				"type":      "FIXED",
			},
			"campaignPublishingFee": map[string]interface{}{
				"amountCFT": 2500,
				"amountINR": 1000,
				"type":      "FIXED",
			},
			"validationFee": map[string]interface{}{
				"amountCFT": 500,
				"amountINR": 200,
				"type":      "FIXED",
			},
			"disputeFilingFee": map[string]interface{}{
				"amountCFT": 750,
				"amountINR": 300,
				"type":      "FIXED",
			},
			"investmentFee": map[string]interface{}{
				"percentage": 5.0,
				"type":       "PERCENTAGE",
			},
			"withdrawalFee": map[string]interface{}{
				"percentage": 1.0,
				"type":       "PERCENTAGE",
			},
		},
	}

	scheduleJSON, err := json.Marshal(feeSchedule)
	if err != nil {
		return "", fmt.Errorf("failed to marshal fee schedule: %v", err)
	}

	return string(scheduleJSON), nil
}

