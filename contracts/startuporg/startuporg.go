package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// StartupContract provides functions for StartupOrg operations
type StartupContract struct {
	contractapi.Contract
}

// ============================================================================
// DATA STRUCTURES
// ============================================================================

// Campaign represents a startup crowdfunding campaign with all required fields
type Campaign struct {
	CampaignID          string   `json:"campaignId"`
	StartupID           string   `json:"startupId"`

	// Core Campaign Fields (as specified)
	Category            string   `json:"category"`
	CloseDate           string   `json:"close_date"`
	Currency            string   `json:"currency"`
	FundsRaisedAmount   float64  `json:"funds_raised_amount"`
	FundsRaisedPercent  float64  `json:"funds_raised_percent"`
	IsIndemand          bool     `json:"is_indemand"`
	IsPreLaunch         bool     `json:"is_pre_launch"`
	OpenDate            string   `json:"open_date"`
	ProductStage        string   `json:"product_stage"`
	ProjectType         string   `json:"project_type"`
	Tags                []string `json:"tags"`
	IsProven            bool     `json:"is_proven"`
	IsPromoted          bool     `json:"is_promoted"`
	DurationDays        int      `json:"duration_days"`
	LaunchMonth         int      `json:"launch_month"`
	LaunchQuarter       int      `json:"launch_quarter"`
	LaunchYear          int      `json:"launch_year"`
	IsSuccessful        bool     `json:"is_successful"`
	AmountUSD           float64  `json:"amount_usd"`
	GoalAmount          float64  `json:"goal_amount"`
	FundingGoalCategory string   `json:"funding_goal_category"`

	// Additional Campaign Metadata
	ProjectName         string   `json:"projectName"`
	Description         string   `json:"description"`

	// Document History - tracks ALL document submissions (linked by CampaignID)
	DocumentHistory     []DocumentSubmission `json:"documentHistory"`
	CurrentDocuments    []string             `json:"currentDocuments"`

	// Status and Tracking
	// ValidationStatus: DRAFT, PENDING_VALIDATION, ON_HOLD, APPROVED, REJECTED, BLACKLISTED
	// Status: DRAFT, SUBMITTED, VALIDATED, APPROVED, REJECTED, PUBLISHED, FUNDED, COMPLETED, CLOSED
	Status              string   `json:"status"`
	ValidationStatus    string   `json:"validationStatus"`
	ValidationScore     float64  `json:"validationScore"`
	ValidationHash      string   `json:"validationHash"` // Hash for cross-org verification
	ValidationHistory   []ValidationEntry `json:"validationHistory"`
	InvestorCount       int      `json:"investorCount"`

	// Platform Status (after validation approval)
	// NOT_SUBMITTED, PENDING_PLATFORM, PLATFORM_VERIFIED, PUBLISHED
	PlatformStatus      string   `json:"platformStatus"`

	// Milestones and Agreements
	Milestones          []Milestone `json:"milestones"`
	AgreementIDs        []string    `json:"agreementIds"`

	// Timestamps
	CreatedAt           string   `json:"createdAt"`
	UpdatedAt           string   `json:"updatedAt"`
	ApprovedAt          string   `json:"approvedAt"`
	PublishedAt         string   `json:"publishedAt"`
}

// DocumentSubmission tracks each document submission attempt (linked by CampaignID)
type DocumentSubmission struct {
	SubmissionID    string   `json:"submissionId"`
	Documents       []string `json:"documents"`
	SubmittedAt     string   `json:"submittedAt"`
	SubmissionNotes string   `json:"submissionNotes"`
	ResponseStatus  string   `json:"responseStatus"` // PENDING, APPROVED, ON_HOLD, REJECTED
	ResponseNotes   string   `json:"responseNotes"`
	ResponseAt      string   `json:"responseAt"`
}

// ValidationEntry tracks validation history
type ValidationEntry struct {
	ValidationID    string  `json:"validationId"`
	ValidatorID     string  `json:"validatorId"`
	Status          string  `json:"status"` // APPROVED, ON_HOLD, REJECTED
	Score           float64 `json:"score"`
	Comments        string  `json:"comments"`
	RequiredDocs    string  `json:"requiredDocs"` // Documents requested if ON_HOLD
	ValidatedAt     string  `json:"validatedAt"`
}

// Milestone represents a funding milestone
type Milestone struct {
	MilestoneID     string  `json:"milestoneId"`
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	TargetAmount    float64 `json:"targetAmount"`
	TargetDate      string  `json:"targetDate"`
	Status          string  `json:"status"` // PENDING, IN_PROGRESS, COMPLETED, VERIFIED
	CompletionProof string  `json:"completionProof"`
	CompletedAt     string  `json:"completedAt"`
	VerifiedAt      string  `json:"verifiedAt"`
	FundsReleased   bool    `json:"fundsReleased"`
}

// Agreement represents investment agreement between startup and investor
type Agreement struct {
	AgreementID        string       `json:"agreementId"`
	CampaignID         string       `json:"campaignId"`
	StartupID          string       `json:"startupId"`
	InvestorID         string       `json:"investorId"`
	InvestmentAmount   float64      `json:"investmentAmount"`
	Currency           string       `json:"currency"`
	Milestones         []Milestone  `json:"milestones"`
	Terms              string       `json:"terms"`
	Status             string       `json:"status"` // PROPOSED, NEGOTIATING, ACCEPTED, ACTIVE, COMPLETED, CANCELLED
	StartupAccepted    bool         `json:"startupAccepted"`
	InvestorAccepted   bool         `json:"investorAccepted"`
	PlatformWitnessed  bool         `json:"platformWitnessed"`
	NegotiationHistory []NegotiationEntry `json:"negotiationHistory"`
	CreatedAt          string       `json:"createdAt"`
	AcceptedAt         string       `json:"acceptedAt"`
}

// NegotiationEntry tracks agreement negotiations
type NegotiationEntry struct {
	EntryID   string `json:"entryId"`
	FromOrg   string `json:"fromOrg"` // STARTUP or INVESTOR
	Action    string `json:"action"`  // PROPOSE, COUNTER, ACCEPT, REJECT, MODIFY
	Changes   string `json:"changes"`
	Timestamp string `json:"timestamp"`
}

// MilestoneReport for progress reporting
type MilestoneReport struct {
	ReportID       string   `json:"reportId"`
	CampaignID     string   `json:"campaignId"`
	MilestoneID    string   `json:"milestoneId"`
	AgreementID    string   `json:"agreementId"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	Evidence       []string `json:"evidence"`
	Status         string   `json:"status"` // SUBMITTED, UNDER_REVIEW, APPROVED, REJECTED
	SubmittedAt    string   `json:"submittedAt"`
	ReviewedAt     string   `json:"reviewedAt"`
	ReviewComments string   `json:"reviewComments"`
}

// CampaignSummaryHash for common-channel (privacy-preserving)
type CampaignSummaryHash struct {
	CampaignID  string `json:"campaignId"`
	SummaryHash string `json:"summaryHash"`
	Status      string `json:"status"`
	Category    string `json:"category"`
	PublishedAt string `json:"publishedAt"`
}

// Investment represents an investment record
type Investment struct {
	InvestmentID   string  `json:"investmentId"`
	CampaignID     string  `json:"campaignId"`
	InvestorID     string  `json:"investorId"`
	Amount         float64 `json:"amount"`
	Currency       string  `json:"currency"`
	Status         string  `json:"status"` // COMMITTED, ACKNOWLEDGED, WITHDRAWN
	CommittedAt    string  `json:"committedAt"`
	AcknowledgedAt string  `json:"acknowledgedAt"`
}

// InitLedger initializes the StartupOrg ledger
func (s *StartupContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	fmt.Println("StartupOrg contract initialized - Merged Version")
	return nil
}

// ============================================================================
// STARTUP-VALIDATOR-CHANNEL FUNCTIONS
// Endorsed by: StartupOrg, ValidatorOrg
// ============================================================================

// CreateCampaign creates a new campaign with initial documents
// Step 1: Startup creates proposal - Status starts as DRAFT
// Channel: startup-validator-channel
// Endorsers: StartupOrg, ValidatorOrg
func (s *StartupContract) CreateCampaign(
	ctx contractapi.TransactionContextInterface,
	campaignID string,
	startupID string,
	category string,
	closeDate string,
	currency string,
	isIndemand bool,
	isPreLaunch bool,
	openDate string,
	productStage string,
	projectType string,
	tagsJSON string,
	isProven bool,
	isPromoted bool,
	durationDays int,
	launchMonth int,
	launchQuarter int,
	launchYear int,
	goalAmount float64,
	fundingGoalCategory string,
	projectName string,
	description string,
	documentsJSON string,
) (string, error) {
	// Check if campaign already exists
	existing, err := ctx.GetStub().GetState(campaignID)
	if err != nil {
		return "", fmt.Errorf("failed to read state: %v", err)
	}
	if existing != nil {
		return "", fmt.Errorf("campaign %s already exists", campaignID)
	}

	// Check if this campaignID was previously blacklisted (rejected for fraud)
	blacklistKey := fmt.Sprintf("BLACKLIST_%s", campaignID)
	blacklisted, _ := ctx.GetStub().GetState(blacklistKey)
	if blacklisted != nil {
		return "", fmt.Errorf("campaign ID %s has been blacklisted and cannot be reused", campaignID)
	}

	// Parse tags
	var tags []string
	if tagsJSON != "" {
		if err := json.Unmarshal([]byte(tagsJSON), &tags); err != nil {
			return "", fmt.Errorf("failed to parse tags: %v", err)
		}
	}

	// Parse documents
	var documents []string
	if documentsJSON != "" {
		if err := json.Unmarshal([]byte(documentsJSON), &documents); err != nil {
			return "", fmt.Errorf("failed to parse documents: %v", err)
		}
	}

	now := time.Now().Format(time.RFC3339)

	// Create initial document submission (linked by campaignID)
	initialSubmission := DocumentSubmission{
		SubmissionID:    fmt.Sprintf("SUB_%s_1", campaignID),
		Documents:       documents,
		SubmittedAt:     now,
		SubmissionNotes: "Initial submission",
		ResponseStatus:  "PENDING",
		ResponseNotes:   "",
		ResponseAt:      "",
	}

	// Create campaign with DRAFT status
	campaign := Campaign{
		CampaignID:          campaignID,
		StartupID:           startupID,
		Category:            category,
		CloseDate:           closeDate,
		Currency:            currency,
		FundsRaisedAmount:   0,
		FundsRaisedPercent:  0,
		IsIndemand:          isIndemand,
		IsPreLaunch:         isPreLaunch,
		OpenDate:            openDate,
		ProductStage:        productStage,
		ProjectType:         projectType,
		Tags:                tags,
		IsProven:            isProven,
		IsPromoted:          isPromoted,
		DurationDays:        durationDays,
		LaunchMonth:         launchMonth,
		LaunchQuarter:       launchQuarter,
		LaunchYear:          launchYear,
		IsSuccessful:        false,
		AmountUSD:           0,
		GoalAmount:          goalAmount,
		FundingGoalCategory: fundingGoalCategory,
		ProjectName:         projectName,
		Description:         description,
		DocumentHistory:     []DocumentSubmission{initialSubmission},
		CurrentDocuments:    documents,
		Status:              "DRAFT",
		ValidationStatus:    "DRAFT",
		ValidationHistory:   []ValidationEntry{},
		InvestorCount:       0,
		PlatformStatus:      "NOT_SUBMITTED",
		Milestones:          []Milestone{},
		AgreementIDs:        []string{},
		CreatedAt:           now,
		UpdatedAt:           now,
		ValidationScore:     0,
		ApprovedAt:          "",
		PublishedAt:         "",
	}

	// Generate validation hash for cross-org verification
	campaign.ValidationHash = generateCampaignHash(campaign)

	campaignJSON, err := json.Marshal(campaign)
	if err != nil {
		return "", err
	}

	// Store on startup-validator-channel
	err = ctx.GetStub().PutState(campaignID, campaignJSON)
	if err != nil {
		return "", err
	}

	// Emit event for ValidatorOrg
	eventPayload := map[string]interface{}{
		"campaignId":  campaignID,
		"startupId":   startupID,
		"category":    category,
		"projectName": projectName,
		"goalAmount":  goalAmount,
		"status":      "DRAFT",
		"channel":     "startup-validator-channel",
		"action":      "CAMPAIGN_CREATED",
		"timestamp":   campaign.CreatedAt,
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("CampaignCreated", eventJSON)

	response := map[string]interface{}{
		"message":         "Campaign created successfully. Submit for validation to proceed.",
		"campaignId":      campaignID,
		"status":          "DRAFT",
		"validationHash":  campaign.ValidationHash,
		"nextStep":        "Call SubmitForValidation to send to ValidatorOrg (ML Model)",
		"channel":         "startup-validator-channel",
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// SubmitForValidation sends campaign to ValidatorOrg (ML Model) for certification
// Step 2: Startup submits proposal for validation
// Channel: startup-validator-channel
// Endorsers: StartupOrg, ValidatorOrg
func (s *StartupContract) SubmitForValidation(
	ctx contractapi.TransactionContextInterface,
	campaignID string,
	submissionNotes string,
) (string, error) {
	// Retrieve campaign
	campaignJSON, err := ctx.GetStub().GetState(campaignID)
	if err != nil {
		return "", fmt.Errorf("failed to read campaign: %v", err)
	}
	if campaignJSON == nil {
		return "", fmt.Errorf("campaign %s does not exist", campaignID)
	}

	var campaign Campaign
	err = json.Unmarshal(campaignJSON, &campaign)
	if err != nil {
		return "", err
	}

	// Check if already rejected/blacklisted - cannot resubmit
	if campaign.ValidationStatus == "REJECTED" || campaign.ValidationStatus == "BLACKLISTED" {
		return "", fmt.Errorf("campaign %s has been rejected and cannot be resubmitted. Create a new campaign with different ID", campaignID)
	}

	// Can only submit if DRAFT or ON_HOLD (for resubmission after adding docs)
	validStates := map[string]bool{"DRAFT": true, "ON_HOLD": true}
	if !validStates[campaign.ValidationStatus] {
		return "", fmt.Errorf("campaign cannot be submitted for validation in current status: %s", campaign.ValidationStatus)
	}

	now := time.Now().Format(time.RFC3339)

	// Update validation status to PENDING_VALIDATION
	campaign.ValidationStatus = "PENDING_VALIDATION"
	campaign.Status = "SUBMITTED"
	campaign.UpdatedAt = now

	// Update latest document submission with notes
	if len(campaign.DocumentHistory) > 0 {
		lastIdx := len(campaign.DocumentHistory) - 1
		if submissionNotes != "" {
			campaign.DocumentHistory[lastIdx].SubmissionNotes = submissionNotes
		}
		campaign.DocumentHistory[lastIdx].ResponseStatus = "PENDING"
	}

	// Generate new hash for verification
	campaign.ValidationHash = generateCampaignHash(campaign)

	updatedCampaignJSON, err := json.Marshal(campaign)
	if err != nil {
		return "", err
	}

	// Store updated campaign
	err = ctx.GetStub().PutState(campaignID, updatedCampaignJSON)
	if err != nil {
		return "", err
	}

	// Emit event for ValidatorOrg (ML Model)
	eventPayload := map[string]interface{}{
		"campaignId":      campaignID,
		"startupId":       campaign.StartupID,
		"validationHash":  campaign.ValidationHash,
		"documentCount":   len(campaign.CurrentDocuments),
		"submissionCount": len(campaign.DocumentHistory),
		"action":          "SUBMITTED_FOR_VALIDATION",
		"channel":         "startup-validator-channel",
		"timestamp":       now,
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("CampaignSubmittedForValidation", eventJSON)

	response := map[string]interface{}{
		"message":        "Campaign submitted to ValidatorOrg (ML Model) for certification",
		"campaignId":     campaignID,
		"validationHash": campaign.ValidationHash,
		"status":         "PENDING_VALIDATION",
		"nextStep":       "Await ValidatorOrg validation (APPROVED, ON_HOLD, or REJECTED)",
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// UpdateCampaignDocs updates campaign documents after validator puts ON_HOLD
// Step 2.1: If validator requests more documents (ON_HOLD status)
// All documents are linked by campaignID to maintain history
// Channel: startup-validator-channel
// Endorsers: StartupOrg, ValidatorOrg
func (s *StartupContract) UpdateCampaignDocs(
	ctx contractapi.TransactionContextInterface,
	campaignID string,
	updatedDocumentsJSON string,
	submissionNotes string,
) (string, error) {
	campaignJSON, err := ctx.GetStub().GetState(campaignID)
	if err != nil {
		return "", fmt.Errorf("failed to read campaign: %v", err)
	}
	if campaignJSON == nil {
		return "", fmt.Errorf("campaign %s does not exist", campaignID)
	}

	var campaign Campaign
	err = json.Unmarshal(campaignJSON, &campaign)
	if err != nil {
		return "", err
	}

	// Can only update documents if ON_HOLD (validator requested more docs)
	if campaign.ValidationStatus != "ON_HOLD" {
		return "", fmt.Errorf("can only update documents when status is ON_HOLD. Current status: %s", campaign.ValidationStatus)
	}

	// Parse updated documents
	var newDocuments []string
	if updatedDocumentsJSON != "" {
		if err := json.Unmarshal([]byte(updatedDocumentsJSON), &newDocuments); err != nil {
			return "", fmt.Errorf("failed to parse updated documents: %v", err)
		}
	}

	if len(newDocuments) == 0 {
		return "", fmt.Errorf("at least one document is required for resubmission")
	}

	now := time.Now().Format(time.RFC3339)

	// Create new document submission entry (linked by campaignID)
	submissionNum := len(campaign.DocumentHistory) + 1
	newSubmission := DocumentSubmission{
		SubmissionID:    fmt.Sprintf("SUB_%s_%d", campaignID, submissionNum),
		Documents:       newDocuments,
		SubmittedAt:     now,
		SubmissionNotes: submissionNotes,
		ResponseStatus:  "PENDING",
	}

	// Add to document history (maintains full history linked by campaignID)
	campaign.DocumentHistory = append(campaign.DocumentHistory, newSubmission)

	// Update current documents (append new docs to existing)
	campaign.CurrentDocuments = append(campaign.CurrentDocuments, newDocuments...)

	campaign.UpdatedAt = now
	campaign.ValidationHash = generateCampaignHash(campaign)

	updatedCampaignJSON, err := json.Marshal(campaign)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(campaignID, updatedCampaignJSON)
	if err != nil {
		return "", err
	}

	// Emit event
	eventPayload := map[string]interface{}{
		"campaignId":       campaignID,
		"submissionId":     newSubmission.SubmissionID,
		"newDocumentCount": len(newDocuments),
		"totalDocuments":   len(campaign.CurrentDocuments),
		"submissionNumber": submissionNum,
		"action":           "DOCUMENTS_UPDATED",
		"channel":          "startup-validator-channel",
		"timestamp":        now,
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("CampaignDocsUpdated", eventJSON)

	response := map[string]interface{}{
		"message":          "Documents added to campaign. Call SubmitForValidation to resubmit.",
		"campaignId":       campaignID,
		"submissionId":     newSubmission.SubmissionID,
		"totalSubmissions": submissionNum,
		"validationHash":   campaign.ValidationHash,
		"nextStep":         "Call SubmitForValidation to resubmit with new documents",
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// ============================================================================
// STARTUP-PLATFORM-CHANNEL FUNCTIONS
// Endorsed by: StartupOrg, PlatformOrg
// ============================================================================

// SubmitForPublishing submits APPROVED campaign to PlatformOrg for publishing
// Step 3: Only allowed if ValidationStatus == APPROVED
// Channel: startup-platform-channel
// Endorsers: StartupOrg, PlatformOrg
func (s *StartupContract) SubmitForPublishing(
	ctx contractapi.TransactionContextInterface,
	campaignID string,
) (string, error) {
	// Retrieve campaign
	campaignJSON, err := ctx.GetStub().GetState(campaignID)
	if err != nil {
		return "", fmt.Errorf("failed to read campaign: %v", err)
	}
	if campaignJSON == nil {
		return "", fmt.Errorf("campaign %s does not exist", campaignID)
	}

	var campaign Campaign
	err = json.Unmarshal(campaignJSON, &campaign)
	if err != nil {
		return "", err
	}

	// *** CRITICAL CHECK: Only APPROVED campaigns can be submitted for publishing ***
	if campaign.ValidationStatus != "APPROVED" {
		return "", fmt.Errorf("campaign must be APPROVED by ValidatorOrg before submitting for publishing. Current validation status: %s", campaign.ValidationStatus)
	}

	// Check not already submitted to platform
	if campaign.PlatformStatus != "NOT_SUBMITTED" {
		return "", fmt.Errorf("campaign already submitted to platform. Status: %s", campaign.PlatformStatus)
	}

	now := time.Now().Format(time.RFC3339)

	// Update platform status
	campaign.PlatformStatus = "PENDING_PLATFORM"
	campaign.Status = "PENDING_PUBLISHING"
	campaign.UpdatedAt = now

	updatedCampaignJSON, err := json.Marshal(campaign)
	if err != nil {
		return "", err
	}

	// Store on startup-platform-channel
	platformKey := fmt.Sprintf("PLATFORM_%s", campaignID)
	err = ctx.GetStub().PutState(platformKey, updatedCampaignJSON)
	if err != nil {
		return "", err
	}

	// Also update main record
	err = ctx.GetStub().PutState(campaignID, updatedCampaignJSON)
	if err != nil {
		return "", err
	}

	// Emit event for PlatformOrg
	eventPayload := map[string]interface{}{
		"campaignId":       campaignID,
		"startupId":        campaign.StartupID,
		"projectName":      campaign.ProjectName,
		"validationHash":   campaign.ValidationHash,
		"validationScore":  campaign.ValidationScore,
		"validationStatus": campaign.ValidationStatus,
		"action":           "SUBMITTED_FOR_PUBLISHING",
		"channel":          "startup-platform-channel",
		"timestamp":        now,
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("CampaignSubmittedForPublishing", eventJSON)

	response := map[string]interface{}{
		"message":         "Validated campaign submitted to PlatformOrg for publishing",
		"campaignId":      campaignID,
		"validationHash":  campaign.ValidationHash,
		"validationScore": campaign.ValidationScore,
		"platformStatus":  "PENDING_PLATFORM",
		"nextStep":        "PlatformOrg will verify with ValidatorOrg and publish to investors",
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// MarkCampaignCompleted marks campaign as completed after reaching target
// Channel: startup-platform-channel
// Endorsers: StartupOrg, PlatformOrg
func (s *StartupContract) MarkCampaignCompleted(
	ctx contractapi.TransactionContextInterface,
	campaignID string,
	fundsRaisedAmount float64,
	amountUSD float64,
) (string, error) {
	platformKey := fmt.Sprintf("PLATFORM_%s", campaignID)
	campaignJSON, err := ctx.GetStub().GetState(platformKey)
	if err != nil {
		return "", fmt.Errorf("failed to read campaign: %v", err)
	}
	if campaignJSON == nil {
		return "", fmt.Errorf("campaign %s does not exist on platform channel", campaignID)
	}

	var campaign Campaign
	err = json.Unmarshal(campaignJSON, &campaign)
	if err != nil {
		return "", err
	}

	campaign.FundsRaisedAmount = fundsRaisedAmount
	campaign.AmountUSD = amountUSD
	campaign.FundsRaisedPercent = (fundsRaisedAmount / campaign.GoalAmount) * 100
	campaign.IsSuccessful = campaign.FundsRaisedPercent >= 100
	campaign.Status = "COMPLETED"
	campaign.UpdatedAt = time.Now().Format(time.RFC3339)

	updatedCampaignJSON, err := json.Marshal(campaign)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(platformKey, updatedCampaignJSON)
	if err != nil {
		return "", err
	}

	// Emit event
	eventPayload := map[string]interface{}{
		"campaignId":         campaignID,
		"fundsRaisedPercent": campaign.FundsRaisedPercent,
		"isSuccessful":       campaign.IsSuccessful,
		"action":             "CAMPAIGN_COMPLETED",
		"channel":            "startup-platform-channel",
		"timestamp":          campaign.UpdatedAt,
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("CampaignCompleted", eventJSON)

	response := map[string]interface{}{
		"message":            "Campaign marked as completed",
		"campaignId":         campaignID,
		"isSuccessful":       campaign.IsSuccessful,
		"fundsRaisedPercent": campaign.FundsRaisedPercent,
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// RespondToInvestmentProposal allows startup to respond to investor's investment proposal
// Step 10: Startup reviews and responds (accept/reject/counter)
// Channel: startup-platform-channel
// Endorsers: StartupOrg, PlatformOrg
func (s *StartupContract) RespondToInvestmentProposal(
	ctx contractapi.TransactionContextInterface,
	agreementID string,
	action string, // ACCEPT, REJECT, COUNTER
	counterTerms string,
	milestonesJSON string,
) (string, error) {
	// Retrieve agreement
	agreementJSON, err := ctx.GetStub().GetState(agreementID)
	if err != nil {
		return "", fmt.Errorf("failed to read agreement: %v", err)
	}
	if agreementJSON == nil {
		return "", fmt.Errorf("agreement %s does not exist", agreementID)
	}

	var agreement Agreement
	err = json.Unmarshal(agreementJSON, &agreement)
	if err != nil {
		return "", err
	}

	now := time.Now().Format(time.RFC3339)

	// Create negotiation entry
	negotiationEntry := NegotiationEntry{
		EntryID:   fmt.Sprintf("NEG_%s_%d", agreementID, len(agreement.NegotiationHistory)+1),
		FromOrg:   "STARTUP",
		Action:    action,
		Changes:   counterTerms,
		Timestamp: now,
	}
	agreement.NegotiationHistory = append(agreement.NegotiationHistory, negotiationEntry)

	switch action {
	case "ACCEPT":
		agreement.StartupAccepted = true
		if agreement.InvestorAccepted {
			agreement.Status = "ACCEPTED"
			agreement.AcceptedAt = now
		} else {
			agreement.Status = "NEGOTIATING"
		}
	case "REJECT":
		agreement.Status = "CANCELLED"
		agreement.StartupAccepted = false
	case "COUNTER":
		agreement.Status = "NEGOTIATING"
		agreement.StartupAccepted = false
		agreement.InvestorAccepted = false
		if counterTerms != "" {
			agreement.Terms = counterTerms
		}
		if milestonesJSON != "" {
			var milestones []Milestone
			if err := json.Unmarshal([]byte(milestonesJSON), &milestones); err == nil {
				agreement.Milestones = milestones
			}
		}
	default:
		return "", fmt.Errorf("invalid action: %s. Must be ACCEPT, REJECT, or COUNTER", action)
	}

	updatedAgreementJSON, err := json.Marshal(agreement)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(agreementID, updatedAgreementJSON)
	if err != nil {
		return "", err
	}

	// Emit event
	eventPayload := map[string]interface{}{
		"agreementId":     agreementID,
		"campaignId":      agreement.CampaignID,
		"action":          action,
		"startupAccepted": agreement.StartupAccepted,
		"status":          agreement.Status,
		"channel":         "startup-platform-channel",
		"timestamp":       now,
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("StartupRespondedToProposal", eventJSON)

	response := map[string]interface{}{
		"message":     fmt.Sprintf("Startup %s the investment proposal", action),
		"agreementId": agreementID,
		"status":      agreement.Status,
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// SubmitMilestoneReport submits proof of milestone completion for validator verification
// Step 12: Startup submits milestone report for validator verification
// Channel: startup-validator-channel
// Endorsers: StartupOrg, ValidatorOrg
func (s *StartupContract) SubmitMilestoneReport(
	ctx contractapi.TransactionContextInterface,
	reportID string,
	campaignID string,
	milestoneID string,
	agreementID string,
	title string,
	description string,
	evidenceJSON string,
) (string, error) {
	// Parse evidence documents
	var evidence []string
	if evidenceJSON != "" {
		if err := json.Unmarshal([]byte(evidenceJSON), &evidence); err != nil {
			return "", fmt.Errorf("failed to parse evidence: %v", err)
		}
	}

	now := time.Now().Format(time.RFC3339)

	// Create milestone report
	report := MilestoneReport{
		ReportID:    reportID,
		CampaignID:  campaignID,
		MilestoneID: milestoneID,
		AgreementID: agreementID,
		Title:       title,
		Description: description,
		Evidence:    evidence,
		Status:      "SUBMITTED",
		SubmittedAt: now,
	}

	reportJSON, err := json.Marshal(report)
	if err != nil {
		return "", err
	}

	// Store report
	err = ctx.GetStub().PutState(reportID, reportJSON)
	if err != nil {
		return "", err
	}

	// Update milestone status in campaign
	platformKey := fmt.Sprintf("PLATFORM_%s", campaignID)
	campaignJSON, err := ctx.GetStub().GetState(platformKey)
	if err == nil && campaignJSON != nil {
		var campaign Campaign
		if json.Unmarshal(campaignJSON, &campaign) == nil {
			for i, m := range campaign.Milestones {
				if m.MilestoneID == milestoneID {
					campaign.Milestones[i].Status = "COMPLETED"
					campaign.Milestones[i].CompletionProof = reportID
					campaign.Milestones[i].CompletedAt = now
					break
				}
			}
			campaign.UpdatedAt = now
			updatedCampaignJSON, _ := json.Marshal(campaign)
			ctx.GetStub().PutState(platformKey, updatedCampaignJSON)
		}
	}

	// Emit event for ValidatorOrg
	eventPayload := map[string]interface{}{
		"reportId":    reportID,
		"campaignId":  campaignID,
		"milestoneId": milestoneID,
		"agreementId": agreementID,
		"action":      "MILESTONE_REPORT_SUBMITTED",
		"channel":     "startup-validator-channel",
		"timestamp":   now,
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("MilestoneReportSubmitted", eventJSON)

	response := map[string]interface{}{
		"message":     "Milestone report submitted for validator verification",
		"reportId":    reportID,
		"milestoneId": milestoneID,
		"status":      "SUBMITTED",
		"nextStep":    "Awaiting Validator verification, then fund release on common-channel",
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// ReceiveFunding records funding received for a milestone
// Step 13: Startup receives funding from Platform escrow on common-channel
// Channel: common-channel
// Endorsers: All Organizations
func (s *StartupContract) ReceiveFunding(
	ctx contractapi.TransactionContextInterface,
	campaignID string,
	milestoneID string,
	amount float64,
	releaseID string,
) (string, error) {
	// Retrieve campaign
	platformKey := fmt.Sprintf("PLATFORM_%s", campaignID)
	campaignJSON, err := ctx.GetStub().GetState(platformKey)
	if err != nil {
		return "", fmt.Errorf("failed to read campaign: %v", err)
	}
	if campaignJSON == nil {
		return "", fmt.Errorf("campaign %s does not exist", campaignID)
	}

	var campaign Campaign
	err = json.Unmarshal(campaignJSON, &campaign)
	if err != nil {
		return "", err
	}

	now := time.Now().Format(time.RFC3339)

	// Update milestone
	for i, m := range campaign.Milestones {
		if m.MilestoneID == milestoneID {
			campaign.Milestones[i].FundsReleased = true
			campaign.Milestones[i].VerifiedAt = now
			campaign.Milestones[i].Status = "VERIFIED"
			break
		}
	}

	// Update total funds
	campaign.FundsRaisedAmount += amount
	if campaign.GoalAmount > 0 {
		campaign.FundsRaisedPercent = (campaign.FundsRaisedAmount / campaign.GoalAmount) * 100
	}
	if campaign.FundsRaisedPercent >= 100 {
		campaign.IsSuccessful = true
	}
	campaign.UpdatedAt = now

	updatedCampaignJSON, err := json.Marshal(campaign)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(platformKey, updatedCampaignJSON)
	if err != nil {
		return "", err
	}

	// Emit event on common-channel
	eventPayload := map[string]interface{}{
		"campaignId":       campaignID,
		"milestoneId":      milestoneID,
		"amount":           amount,
		"releaseId":        releaseID,
		"totalFundsRaised": campaign.FundsRaisedAmount,
		"percentComplete":  campaign.FundsRaisedPercent,
		"action":           "FUNDING_RECEIVED",
		"channel":          "common-channel",
		"timestamp":        now,
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("FundingReceived", eventJSON)

	response := map[string]interface{}{
		"message":          "Funding received successfully",
		"campaignId":       campaignID,
		"milestoneId":      milestoneID,
		"amountReceived":   amount,
		"totalFundsRaised": campaign.FundsRaisedAmount,
		"percentComplete":  campaign.FundsRaisedPercent,
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// ============================================================================
// COMMON-CHANNEL FUNCTIONS (Multi-party visibility)
// Endorsed by: All Organizations
// ============================================================================

// AcknowledgeInvestment acknowledges investor commitment (Platform forwards after receiving funds)
// Step 11: Startup acknowledges investment via Platform notification on common-channel
// Channel: common-channel
// Endorsers: All Organizations
func (s *StartupContract) AcknowledgeInvestment(
	ctx contractapi.TransactionContextInterface,
	investmentID string,
	campaignID string,
	investorID string,
	amount float64,
	currency string,
) (string, error) {
	investment := Investment{
		InvestmentID:   investmentID,
		CampaignID:     campaignID,
		InvestorID:     investorID,
		Amount:         amount,
		Currency:       currency,
		Status:         "ACKNOWLEDGED",
		AcknowledgedAt: time.Now().Format(time.RFC3339),
	}

	investmentJSON, err := json.Marshal(investment)
	if err != nil {
		return "", err
	}

	// Store investment acknowledgment
	ackKey := fmt.Sprintf("ACK_%s", investmentID)
	err = ctx.GetStub().PutState(ackKey, investmentJSON)
	if err != nil {
		return "", err
	}

	// Update campaign investor count on investor channel
	investorCampaignKey := fmt.Sprintf("INVESTOR_CAMPAIGN_%s", campaignID)
	campaignJSON, err := ctx.GetStub().GetState(investorCampaignKey)

	var campaign Campaign
	if campaignJSON != nil {
		json.Unmarshal(campaignJSON, &campaign)
		campaign.FundsRaisedAmount += amount
		campaign.InvestorCount++
		campaign.FundsRaisedPercent = (campaign.FundsRaisedAmount / campaign.GoalAmount) * 100
	} else {
		campaign = Campaign{
			CampaignID:         campaignID,
			FundsRaisedAmount:  amount,
			InvestorCount:      1,
			FundsRaisedPercent: 0,
		}
	}
	campaign.UpdatedAt = time.Now().Format(time.RFC3339)

	updatedCampaignJSON, _ := json.Marshal(campaign)
	ctx.GetStub().PutState(investorCampaignKey, updatedCampaignJSON)

	// Emit event on common-channel (visible to all orgs)
	eventPayload := map[string]interface{}{
		"investmentId": investmentID,
		"campaignId":   campaignID,
		"investorId":   investorID,
		"amount":       amount,
		"action":       "INVESTMENT_ACKNOWLEDGED",
		"channel":      "common-channel",
		"timestamp":    investment.AcknowledgedAt,
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("InvestmentAcknowledged", eventJSON)

	response := map[string]interface{}{
		"message":      "Investment acknowledged by startup",
		"investmentId": investmentID,
		"campaignId":   campaignID,
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// ============================================================================
// COMMON-CHANNEL FUNCTIONS
// Read by: All Orgs, Write by: StartupOrg (hash only - privacy preserving)
// ============================================================================

// PublishSummaryHash publishes campaign summary hash to common-channel
// Channel: common-channel
// Purpose: Privacy-preserving summary for all organizations (no sensitive data)
func (s *StartupContract) PublishSummaryHash(
	ctx contractapi.TransactionContextInterface,
	campaignID string,
	status string,
	category string,
) (string, error) {
	// Generate summary hash (no sensitive/financial data included)
	summary := map[string]interface{}{
		"campaignId": campaignID,
		"status":     status,
		"category":   category,
		"timestamp":  time.Now().Format(time.RFC3339),
	}

	summaryJSON, _ := json.Marshal(summary)
	summaryHash := generateHash(string(summaryJSON))

	campaignSummary := CampaignSummaryHash{
		CampaignID:  campaignID,
		SummaryHash: summaryHash,
		Status:      status,
		Category:    category,
		PublishedAt: time.Now().Format(time.RFC3339),
	}

	summaryHashJSON, err := json.Marshal(campaignSummary)
	if err != nil {
		return "", err
	}

	// Store on common-channel (read by all orgs)
	commonKey := fmt.Sprintf("COMMON_SUMMARY_%s", campaignID)
	err = ctx.GetStub().PutState(commonKey, summaryHashJSON)
	if err != nil {
		return "", err
	}

	// Emit event
	eventPayload := map[string]interface{}{
		"campaignId":  campaignID,
		"summaryHash": summaryHash,
		"status":      status,
		"channel":     "common-channel",
		"timestamp":   campaignSummary.PublishedAt,
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("SummaryHashPublished", eventJSON)

	response := map[string]interface{}{
		"message":     "Campaign summary hash published to common channel",
		"campaignId":  campaignID,
		"summaryHash": summaryHash,
		"channel":     "common-channel",
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// ============================================================================
// QUERY FUNCTIONS
// ============================================================================

// GetCampaign retrieves campaign by ID
func (s *StartupContract) GetCampaign(ctx contractapi.TransactionContextInterface, campaignID string) (*Campaign, error) {
	campaignJSON, err := ctx.GetStub().GetState(campaignID)
	if err != nil {
		return nil, fmt.Errorf("failed to read campaign: %v", err)
	}
	if campaignJSON == nil {
		return nil, fmt.Errorf("campaign %s does not exist", campaignID)
	}

	var campaign Campaign
	err = json.Unmarshal(campaignJSON, &campaign)
	if err != nil {
		return nil, err
	}

	return &campaign, nil
}

// GetCampaignValidationHash returns the validation hash for cross-org verification
func (s *StartupContract) GetCampaignValidationHash(ctx contractapi.TransactionContextInterface, campaignID string) (string, error) {
	campaign, err := s.GetCampaign(ctx, campaignID)
	if err != nil {
		return "", err
	}

	response := map[string]interface{}{
		"campaignId":       campaignID,
		"validationHash":   campaign.ValidationHash,
		"validationStatus": campaign.ValidationStatus,
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// GetCampaignDocumentHistory returns full document submission history
func (s *StartupContract) GetCampaignDocumentHistory(ctx contractapi.TransactionContextInterface, campaignID string) (string, error) {
	campaign, err := s.GetCampaign(ctx, campaignID)
	if err != nil {
		return "", err
	}

	response := map[string]interface{}{
		"campaignId":       campaignID,
		"documentHistory":  campaign.DocumentHistory,
		"totalSubmissions": len(campaign.DocumentHistory),
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// GetCampaignsByCategory returns campaigns filtered by category
func (s *StartupContract) GetCampaignsByCategory(ctx contractapi.TransactionContextInterface, category string) (string, error) {
	queryString := fmt.Sprintf(`{"selector":{"category":"%s"}}`, category)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return "", err
	}
	defer resultsIterator.Close()

	var campaigns []map[string]interface{}
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return "", err
		}

		var campaign Campaign
		err = json.Unmarshal(queryResponse.Value, &campaign)
		if err != nil {
			continue
		}

		campaignMap := map[string]interface{}{
			"Key":    queryResponse.Key,
			"Record": campaign,
		}
		campaigns = append(campaigns, campaignMap)
	}

	campaignsJSON, err := json.Marshal(campaigns)
	if err != nil {
		return "", err
	}

	return string(campaignsJSON), nil
}

// GetCampaignsByStartup returns all campaigns by startup ID
func (s *StartupContract) GetCampaignsByStartup(ctx contractapi.TransactionContextInterface, startupID string) (string, error) {
	queryString := fmt.Sprintf(`{"selector":{"startupId":"%s"}}`, startupID)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return "", err
	}
	defer resultsIterator.Close()

	var campaigns []map[string]interface{}
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return "", err
		}

		var campaign Campaign
		err = json.Unmarshal(queryResponse.Value, &campaign)
		if err != nil {
			continue
		}

		campaignMap := map[string]interface{}{
			"Key":    queryResponse.Key,
			"Record": campaign,
		}
		campaigns = append(campaigns, campaignMap)
	}

	campaignsJSON, err := json.Marshal(campaigns)
	if err != nil {
		return "", err
	}

	return string(campaignsJSON), nil
}

// GetAgreement retrieves agreement by ID
func (s *StartupContract) GetAgreement(ctx contractapi.TransactionContextInterface, agreementID string) (string, error) {
	agreementJSON, err := ctx.GetStub().GetState(agreementID)
	if err != nil {
		return "", fmt.Errorf("failed to read agreement: %v", err)
	}
	if agreementJSON == nil {
		return "", fmt.Errorf("agreement %s does not exist", agreementID)
	}
	return string(agreementJSON), nil
}

// GetMilestoneReport retrieves milestone report by ID
func (s *StartupContract) GetMilestoneReport(ctx contractapi.TransactionContextInterface, reportID string) (string, error) {
	reportJSON, err := ctx.GetStub().GetState(reportID)
	if err != nil {
		return "", fmt.Errorf("failed to read report: %v", err)
	}
	if reportJSON == nil {
		return "", fmt.Errorf("report %s does not exist", reportID)
	}
	return string(reportJSON), nil
}

// ============================================================================
// CROSS-CHANNEL INVOCATION HELPER FUNCTIONS
// ============================================================================

// InvokePlatformOrgPublish invokes PlatformOrg to publish campaign after validation
// This performs a cross-channel call from startup-validator-channel to startup-platform-channel
func (s *StartupContract) InvokePlatformOrgPublish(
	ctx contractapi.TransactionContextInterface,
	campaignID string,
	startupID string,
	projectName string,
	category string,
	description string,
	goalAmount string,
	currency string,
	openDate string,
	closeDate string,
	durationDays string,
	validationScore string,
) (string, error) {
	// Prepare arguments for cross-channel invocation
	args := [][]byte{
		[]byte("PublishCampaignToPortal"),
		[]byte(campaignID),
		[]byte(startupID),
		[]byte(projectName),
		[]byte(category),
		[]byte(description),
		[]byte(goalAmount),
		[]byte(currency),
		[]byte(openDate),
		[]byte(closeDate),
		[]byte(durationDays),
		[]byte(validationScore),
	}

	// Cross-channel invocation to platformorg on startup-platform-channel
	// NOTE: The peer executing this must be a member of BOTH channels
	response := ctx.GetStub().InvokeChaincode(
		"platformorg",              // chaincode name
		args,                        // function + arguments
		"startup-platform-channel", // target channel
	)

	if response.Status != 200 {
		return "", fmt.Errorf("cross-channel invoke to PlatformOrg failed: %s", response.Message)
	}

	// Emit event for audit trail
	eventPayload := map[string]interface{}{
		"campaignId":     campaignID,
		"targetChannel":  "startup-platform-channel",
		"targetContract": "platformorg",
		"action":         "CROSS_CHANNEL_PUBLISH",
		"timestamp":      time.Now().Format(time.RFC3339),
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("CrossChannelInvoke", eventJSON)

	return string(response.Payload), nil
}

// InvokeInvestorOrgNotify notifies InvestorOrg about campaign status changes
// Cross-channel call from startup-platform-channel to startup-investor-channel
func (s *StartupContract) InvokeInvestorOrgNotify(
	ctx contractapi.TransactionContextInterface,
	campaignID string,
	status string,
	message string,
) (string, error) {
	args := [][]byte{
		[]byte("ReceiveCampaignNotification"),
		[]byte(campaignID),
		[]byte(status),
		[]byte(message),
	}

	response := ctx.GetStub().InvokeChaincode(
		"investororg",
		args,
		"startup-investor-channel",
	)

	if response.Status != 200 {
		return "", fmt.Errorf("cross-channel invoke to InvestorOrg failed: %s", response.Message)
	}

	return string(response.Payload), nil
}

// generateHash generates SHA256 hash
func generateHash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// generateCampaignHash generates SHA256 hash of campaign for cross-org verification
func generateCampaignHash(campaign Campaign) string {
	hashData := map[string]interface{}{
		"campaignId":       campaign.CampaignID,
		"startupId":        campaign.StartupID,
		"projectName":      campaign.ProjectName,
		"goalAmount":       campaign.GoalAmount,
		"currentDocuments": campaign.CurrentDocuments,
		"documentHistory":  len(campaign.DocumentHistory),
	}
	hashJSON, _ := json.Marshal(hashData)
	hash := sha256.Sum256(hashJSON)
	return hex.EncodeToString(hash[:])
}

func main() {
	startupChaincode, err := contractapi.NewChaincode(&StartupContract{})
	if err != nil {
		fmt.Printf("Error creating StartupOrg chaincode: %v\n", err)
		return
	}

	if err := startupChaincode.Start(); err != nil {
		fmt.Printf("Error starting StartupOrg chaincode: %v\n", err)
	}
}
