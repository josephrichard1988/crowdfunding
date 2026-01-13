package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// InvestorContract provides functions for InvestorOrg operations using PDC
type InvestorContract struct {
	contractapi.Contract
}

// ============================================================================
// DATA STRUCTURES
// ============================================================================

// CampaignView represents campaign details visible to investors (22-parameter format)
type CampaignView struct {
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
	OpenDate           string  `json:"open_date"`
	CloseDate          string  `json:"close_date"`
	FundsRaisedAmount  float64 `json:"funds_raised_amount"`
	FundsRaisedPercent float64 `json:"funds_raised_percent"`
	ValidationScore    float64 `json:"validationScore"`
	RiskLevel          string  `json:"riskLevel"`
	InvestorCount      int     `json:"investorCount"`
	Status             string  `json:"status"`
	ViewedAt           string  `json:"viewedAt"`
}

// InvestmentProposal represents an investment proposal with terms
type InvestmentProposal struct {
	ProposalID       string             `json:"proposalId"`
	CampaignID       string             `json:"campaignId"`
	StartupID        string             `json:"startupId"`
	InvestorID       string             `json:"investorId"`
	InvestmentAmount float64            `json:"investmentAmount"`
	Currency         string             `json:"currency"`
	ProposedTerms    string             `json:"proposedTerms"`
	Milestones       []Milestone        `json:"milestones"`
	Status           string             `json:"status"` // PROPOSED, COUNTERED, ACCEPTED, REJECTED, EXPIRED
	NegotiationRound int                `json:"negotiationRound"`
	History          []NegotiationEntry `json:"history"`
	CreatedAt        string             `json:"createdAt"`
	UpdatedAt        string             `json:"updatedAt"`
}

// FundingCommitment represents confirmed funding commitment
type FundingCommitment struct {
	CommitmentID string      `json:"commitmentId"`
	ProposalID   string      `json:"proposalId"`
	AgreementID  string      `json:"agreementId"`
	CampaignID   string      `json:"campaignId"`
	StartupID    string      `json:"startupId"`
	InvestorID   string      `json:"investorId"`
	Amount       float64     `json:"amount"`
	Currency     string      `json:"currency"`
	Milestones   []Milestone `json:"milestones"`
	Status       string      `json:"status"` // COMMITTED, ESCROWED, PARTIALLY_RELEASED, RELEASED
	CommittedAt  string      `json:"committedAt"`
}

// MilestoneVerification represents investor verification of milestone
type MilestoneVerification struct {
	VerificationID string `json:"verificationId"`
	MilestoneID    string `json:"milestoneId"`
	AgreementID    string `json:"agreementId"`
	CampaignID     string `json:"campaignId"`
	InvestorID     string `json:"investorId"`
	Approved       bool   `json:"approved"`
	Feedback       string `json:"feedback"`
	VerifiedAt     string `json:"verifiedAt"`
}

// RiskInsightRequest represents investor's request for risk info
type RiskInsightRequest struct {
	RequestID   string `json:"requestId"`
	CampaignID  string `json:"campaignId"`
	InvestorID  string `json:"investorId"`
	Status      string `json:"status"` // PENDING, FULFILLED
	RequestedAt string `json:"requestedAt"`
	FulfilledAt string `json:"fulfilledAt"`
}

// RiskInsightResponse represents response from Validator
type RiskInsightResponse struct {
	ResponseID     string  `json:"responseId"`
	RequestID      string  `json:"requestId"`
	CampaignID     string  `json:"campaignId"`
	InvestorID     string  `json:"investorId"`
	RiskScore      float64 `json:"riskScore"`
	RiskLevel      string  `json:"riskLevel"`
	RiskFactors    string  `json:"riskFactors"`
	Recommendation string  `json:"recommendation"`
	ReceivedAt     string  `json:"receivedAt"`
}

// InvestmentConfirmation represents confirmation sent to PlatformOrg
type InvestmentConfirmation struct {
	ConfirmationID string  `json:"confirmationId"`
	InvestmentID   string  `json:"investmentId"`
	CampaignID     string  `json:"campaignId"`
	InvestorID     string  `json:"investorId"`
	Amount         float64 `json:"amount"`
	Currency       string  `json:"currency"`
	ConfirmedAt    string  `json:"confirmedAt"`
}

// InvestorDisputeSubmission represents a dispute submitted by investor
type InvestorDisputeSubmission struct {
	SubmissionID   string   `json:"submissionId"`
	DisputeID      string   `json:"disputeId"`
	InvestorID     string   `json:"investorId"`
	DisputeType    string   `json:"disputeType"`
	TargetID       string   `json:"targetId"`
	TargetType     string   `json:"targetType"`
	CampaignID     string   `json:"campaignId"`
	AgreementID    string   `json:"agreementId"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	ClaimedAmount  float64  `json:"claimedAmount"`
	EvidenceHashes []string `json:"evidenceHashes"`
	Status         string   `json:"status"`
	CreatedAt      string   `json:"createdAt"`
}

// RefundRequest represents a refund request from investor
type RefundRequest struct {
	RequestID        string  `json:"requestId"`
	InvestorID       string  `json:"investorId"`
	CampaignID       string  `json:"campaignId"`
	AgreementID      string  `json:"agreementId"`
	StartupID        string  `json:"startupId"`
	OriginalAmount   float64 `json:"originalAmount"`
	RequestedAmount  float64 `json:"requestedAmount"`
	RefundReason     string  `json:"refundReason"`
	DeductionPercent float64 `json:"deductionPercent"`
	ExpectedRefund   float64 `json:"expectedRefund"`
	Status           string  `json:"status"`
	RequestedAt      string  `json:"requestedAt"`
	ProcessedAt      string  `json:"processedAt"`
}

// ============================================================================
// INIT
// ============================================================================

func (i *InvestorContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	fmt.Println("InvestorOrg contract initialized with PDC support")
	return nil
}

// ============================================================================
// CAMPAIGN VIEWING - Using Public Ledger
// ============================================================================

// ViewCampaign allows investor to view campaign details from public ledger
// Logs the view for tracking
func (i *InvestorContract) ViewCampaign(
	ctx contractapi.TransactionContextInterface,
	viewRecordID string,
	campaignID string,
	investorID string,
) (*CampaignView, error) {

	// Get campaign from public world state
	publicJSON, err := ctx.GetStub().GetState("CAMPAIGN_PUBLIC_" + campaignID)
	if err != nil {
		return nil, fmt.Errorf("failed to read campaign: %v", err)
	}
	if publicJSON == nil {
		return nil, fmt.Errorf("campaign not found")
	}

	var publicInfo map[string]interface{}
	err = json.Unmarshal(publicJSON, &publicInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal campaign: %v", err)
	}

	timestamp := time.Now().Format(time.RFC3339)

	// Create campaign view
	projectName, _ := publicInfo["projectName"].(string)
	category, _ := publicInfo["category"].(string)
	goalAmount, _ := publicInfo["goalAmount"].(float64)
	currency, _ := publicInfo["currency"].(string)
	status, _ := publicInfo["status"].(string)
	description, _ := publicInfo["description"].(string)

	// Handle Arrays safely
	tags := []string{}
	if t, ok := publicInfo["tags"].([]interface{}); ok {
		for _, tag := range t {
			if s, ok := tag.(string); ok {
				tags = append(tags, s)
			}
		}
	}

	documents := []string{}
	if d, ok := publicInfo["documents"].([]interface{}); ok {
		for _, doc := range d {
			if s, ok := doc.(string); ok {
				documents = append(documents, s)
			}
		}
	}

	campaignView := CampaignView{
		CampaignID:  campaignID,
		ProjectName: projectName,
		Category:    category,
		GoalAmount:  goalAmount,
		Currency:    currency,
		Status:      status,
		Description: description,
		Tags:        tags,
		Documents:   documents,
		ViewedAt:    timestamp,
	}

	// Log view in investor's private collection
	viewJSON, _ := json.Marshal(campaignView)
	// Use the provided viewRecordID for the key
	viewKey := "CAMPAIGN_VIEW_" + viewRecordID
	ctx.GetStub().PutPrivateData(InvestorPrivateCollection, viewKey, viewJSON)

	return &campaignView, nil
}

// ViewCampaignDetails retrieves public campaign details without logging a view
// This provides the detailed "convincing" info for an investor to make a decision
func (i *InvestorContract) ViewCampaignDetails(
	ctx contractapi.TransactionContextInterface,
	campaignID string,
) (*CampaignView, error) {

	// Get campaign from public world state
	publicJSON, err := ctx.GetStub().GetState("CAMPAIGN_PUBLIC_" + campaignID)
	if err != nil {
		return nil, fmt.Errorf("failed to read campaign: %v", err)
	}
	if publicJSON == nil {
		return nil, fmt.Errorf("campaign not found")
	}

	var publicInfo map[string]interface{}
	err = json.Unmarshal(publicJSON, &publicInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal campaign: %v", err)
	}

	timestamp := time.Now().Format(time.RFC3339)

	// Create campaign view
	projectName, _ := publicInfo["projectName"].(string)
	startupID, _ := publicInfo["startupId"].(string)
	category, _ := publicInfo["category"].(string)
	goalAmount, _ := publicInfo["goalAmount"].(float64)
	currency, _ := publicInfo["currency"].(string)
	status, _ := publicInfo["status"].(string)
	description, _ := publicInfo["description"].(string)
	deadline, _ := publicInfo["deadline"].(string)
	riskLevel, _ := publicInfo["riskLevel"].(string)
	validationScore, _ := publicInfo["dueDiligenceScore"].(float64) // Use DD score or Risk Score

	// Handle Arrays safely
	tags := []string{}
	if t, ok := publicInfo["tags"].([]interface{}); ok {
		for _, tag := range t {
			if s, ok := tag.(string); ok {
				tags = append(tags, s)
			}
		}
	}

	documents := []string{}
	// Only show documents if explicitly desired, but user said "convince to invest" so docs are important
	if d, ok := publicInfo["documents"].([]interface{}); ok {
		for _, doc := range d {
			if s, ok := doc.(string); ok {
				documents = append(documents, s)
			}
		}
	}

	campaignView := CampaignView{
		CampaignID:      campaignID,
		StartupID:       startupID,
		ProjectName:     projectName,
		Category:        category,
		GoalAmount:      goalAmount,
		Currency:        currency,
		Status:          status,
		Description:     description,
		Deadline:        deadline,
		RiskLevel:       riskLevel,
		ValidationScore: validationScore,
		Tags:            tags,
		Documents:       documents,
		ViewedAt:        timestamp,
	}

	return &campaignView, nil
}

// GetAvailableCampaigns retrieves all published campaigns with minimal details (list view)
func (i *InvestorContract) GetAvailableCampaigns(
	ctx contractapi.TransactionContextInterface,
) ([]*CampaignView, error) {

	// Range query for all public campaigns
	resultsIterator, err := ctx.GetStub().GetStateByRange("CAMPAIGN_PUBLIC_", "CAMPAIGN_PUBLIC_~")
	if err != nil {
		return nil, fmt.Errorf("failed to get campaigns: %v", err)
	}
	defer resultsIterator.Close()

	var campaigns []*CampaignView

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var publicInfo map[string]interface{}
		err = json.Unmarshal(queryResponse.Value, &publicInfo)
		if err != nil {
			continue // Skip malformed records
		}

		// Extract status first - ONLY include PUBLISHED campaigns
		status, _ := publicInfo["status"].(string)
		if status != "PUBLISHED" {
			continue // Skip non-published campaigns
		}

		// Extract ID
		campaignID := ""
		if id, ok := publicInfo["campaignId"].(string); ok {
			campaignID = id
		}

		projectName, _ := publicInfo["projectName"].(string)
		startupID, _ := publicInfo["startupId"].(string)
		category, _ := publicInfo["category"].(string)
		sector, _ := publicInfo["sector"].(string)
		description, _ := publicInfo["description"].(string)
		deadline, _ := publicInfo["deadline"].(string)
		currency, _ := publicInfo["currency"].(string)
		goalAmount, _ := publicInfo["goalAmount"].(float64)
		fundsRaisedPercent, _ := publicInfo["fundsRaisedPercent"].(float64)
		validationScore, _ := publicInfo["validationScore"].(float64)
		riskLevel, _ := publicInfo["riskLevel"].(string)

		// Extract tags
		var tags []string
		if tagsInterface, ok := publicInfo["tags"].([]interface{}); ok {
			for _, t := range tagsInterface {
				if tagStr, ok := t.(string); ok {
					tags = append(tags, tagStr)
				}
			}
		}
		if tags == nil {
			tags = []string{}
		}

		// Published campaign view - rich metadata for investors
		summary := &CampaignView{
			CampaignID:         campaignID,
			StartupID:          startupID,
			ProjectName:        projectName,
			Category:           category,
			Sector:             sector,
			Description:        description,
			Deadline:           deadline,
			Currency:           currency,
			GoalAmount:         goalAmount,
			FundsRaisedPercent: fundsRaisedPercent,
			ValidationScore:    validationScore,
			RiskLevel:          riskLevel,
			Status:             status,
			Tags:               tags,
			Documents:          []string{}, // Empty slice, not nil
		}
		campaigns = append(campaigns, summary)
	}

	if campaigns == nil {
		campaigns = []*CampaignView{}
	}

	return campaigns, nil
}

// ============================================================================
// INVESTMENT MANAGEMENT - Using PDC
// ============================================================================

// MakeInvestment creates an investment commitment
// Stored in StartupInvestorCollection (shared with startup)
func (i *InvestorContract) MakeInvestment(
	ctx contractapi.TransactionContextInterface,
	investmentID string,
	campaignID string,
	investorID string,
	amount float64,
	currency string,
) error {

	timestamp := time.Now().Format(time.RFC3339)

	// Check for duplicate investment ID
	existingInvJSON, _ := ctx.GetStub().GetPrivateData(StartupInvestorCollection, "INVESTMENT_"+investmentID)
	if existingInvJSON != nil {
		return fmt.Errorf("investment %s already exists", investmentID)
	}

	investment := Investment{
		InvestmentID: investmentID,
		CampaignID:   campaignID,
		InvestorID:   investorID,
		Amount:       amount,
		Currency:     currency,
		Status:       "COMMITTED",
		CommittedAt:  timestamp,
	}

	investmentJSON, err := json.Marshal(investment)
	if err != nil {
		return fmt.Errorf("failed to marshal investment: %v", err)
	}

	// Store in shared collection with StartupOrg
	// This enables StartupContract:GetCampaign to see the investment
	err = ctx.GetStub().PutPrivateData(StartupInvestorCollection, "INVESTMENT_"+investmentID, investmentJSON)
	if err != nil {
		return fmt.Errorf("failed to create investment: %v", err)
	}

	// Also track in investor's private collection
	err = ctx.GetStub().PutPrivateData(InvestorPrivateCollection, "MY_INVESTMENT_"+investmentID, investmentJSON)
	if err != nil {
		return fmt.Errorf("failed to track investment: %v", err)
	}

	return nil
}

// WithdrawInvestment withdraws an investment before acknowledgment
func (i *InvestorContract) WithdrawInvestment(
	ctx contractapi.TransactionContextInterface,
	investmentID string,
) error {

	// Get investment from shared collection
	investmentJSON, err := ctx.GetStub().GetPrivateData(StartupInvestorCollection, "INVESTMENT_"+investmentID)
	if err != nil || investmentJSON == nil {
		return fmt.Errorf("investment not found: %v", err)
	}

	var investment Investment
	err = json.Unmarshal(investmentJSON, &investment)
	if err != nil {
		return fmt.Errorf("failed to unmarshal investment: %v", err)
	}

	if investment.Status != "COMMITTED" {
		return fmt.Errorf("investment cannot be withdrawn in current status")
	}

	timestamp := time.Now().Format(time.RFC3339)
	investment.Status = "WITHDRAWN"
	investment.WithdrawnAt = timestamp

	investmentJSON, _ = json.Marshal(investment)
	err = ctx.GetStub().PutPrivateData(StartupInvestorCollection, "INVESTMENT_"+investmentID, investmentJSON)
	if err != nil {
		return fmt.Errorf("failed to withdraw investment: %v", err)
	}

	// Update in private collection
	ctx.GetStub().PutPrivateData(InvestorPrivateCollection, "MY_INVESTMENT_"+investmentID, investmentJSON)

	return nil
}

// ============================================================================
// INVESTMENT PROPOSALS - Using PDC
// ============================================================================

// CreateInvestmentProposal creates an investment proposal with terms
// Stored in StartupInvestorCollection
func (i *InvestorContract) CreateInvestmentProposal(
	ctx contractapi.TransactionContextInterface,
	proposalID string,
	campaignID string,
	investorID string,
	startupID string,
	investmentAmount float64,
	currency string,
	equity string,
	duration string,
	milestonesJSON string,
	proposedTerms string,
) error {

	var milestones []Milestone
	if milestonesJSON != "" {
		err := json.Unmarshal([]byte(milestonesJSON), &milestones)
		if err != nil {
			return fmt.Errorf("failed to parse milestones: %v. Input was: %s", err, milestonesJSON)
		}
	}

	timestamp := time.Now().Format(time.RFC3339)

	// Construct terms from equity and duration if separate fields aren't in struct
	fullTerms := fmt.Sprintf("Equity: %s, Duration: %s. %s", equity, duration, proposedTerms)

	proposal := InvestmentProposal{
		ProposalID:       proposalID,
		CampaignID:       campaignID,
		StartupID:        startupID,
		InvestorID:       investorID,
		InvestmentAmount: investmentAmount,
		Currency:         currency,
		ProposedTerms:    fullTerms,
		Milestones:       milestones,
		Status:           "PROPOSED",
		NegotiationRound: 1,
		History:          []NegotiationEntry{},
		CreatedAt:        timestamp,
		UpdatedAt:        timestamp,
	}

	// Add initial history entry
	historyEntry := NegotiationEntry{
		Round:     1,
		Party:     "INVESTOR",
		Action:    "PROPOSE",
		Amount:    investmentAmount,
		Terms:     proposedTerms,
		Timestamp: timestamp,
	}
	proposal.History = append(proposal.History, historyEntry)

	proposalJSON, err := json.Marshal(proposal)
	if err != nil {
		return fmt.Errorf("failed to marshal proposal: %v", err)
	}

	// Store in shared collection with StartupOrg
	err = ctx.GetStub().PutPrivateData(StartupInvestorCollection, "PROPOSAL_"+proposalID, proposalJSON)
	if err != nil {
		return fmt.Errorf("failed to create proposal: %v", err)
	}

	return nil
}

// RespondToCounterOffer responds to a counter offer from startup
func (i *InvestorContract) RespondToCounterOffer(
	ctx contractapi.TransactionContextInterface,
	proposalID string,
	action string,
	counterAmount float64,
	counterTerms string,
) error {

	// Get proposal from shared collection
	proposalJSON, err := ctx.GetStub().GetPrivateData(StartupInvestorCollection, "PROPOSAL_"+proposalID)
	if err != nil || proposalJSON == nil {
		return fmt.Errorf("proposal not found: %v", err)
	}

	var proposal InvestmentProposal
	err = json.Unmarshal(proposalJSON, &proposal)
	if err != nil {
		return fmt.Errorf("failed to unmarshal proposal: %v", err)
	}

	timestamp := time.Now().Format(time.RFC3339)

	if action == "ACCEPT" {
		proposal.Status = "ACCEPTED"
	} else if action == "REJECT" {
		proposal.Status = "REJECTED"
	} else if action == "COUNTER" {
		proposal.Status = "COUNTERED"
		proposal.InvestmentAmount = counterAmount
		proposal.ProposedTerms = counterTerms
		proposal.NegotiationRound++

		// Add history entry
		historyEntry := NegotiationEntry{
			Round:     proposal.NegotiationRound,
			Party:     "INVESTOR",
			Action:    "COUNTER",
			Amount:    counterAmount,
			Terms:     counterTerms,
			Timestamp: timestamp,
		}
		proposal.History = append(proposal.History, historyEntry)
	}

	proposal.UpdatedAt = timestamp

	proposalJSON, _ = json.Marshal(proposal)
	err = ctx.GetStub().PutPrivateData(StartupInvestorCollection, "PROPOSAL_"+proposalID, proposalJSON)
	if err != nil {
		return fmt.Errorf("failed to update proposal: %v", err)
	}

	return nil
}

// ============================================================================
// AGREEMENT & FUNDING - Using PDC
// ============================================================================

// AcceptAgreement accepts a finalized agreement
// Agreement moved to ThreePartyCollection (Startup, Investor, Platform)
func (i *InvestorContract) AcceptAgreement(
	ctx contractapi.TransactionContextInterface,
	agreementID string,
	investorID string,
) error {

	// Get agreement from StartupInvestorCollection
	agreementJSON, err := ctx.GetStub().GetPrivateData(StartupInvestorCollection, "AGREEMENT_"+agreementID)
	if err != nil || agreementJSON == nil {
		return fmt.Errorf("agreement not found: %v", err)
	}

	var agreement map[string]interface{}
	err = json.Unmarshal(agreementJSON, &agreement)
	if err != nil {
		return fmt.Errorf("failed to unmarshal agreement: %v", err)
	}

	timestamp := time.Now().Format(time.RFC3339)
	agreement["investorAccepted"] = true
	agreement["acceptedAt"] = timestamp
	agreement["status"] = "ACCEPTED"

	// Update agreement
	agreementJSON, _ = json.Marshal(agreement)
	err = ctx.GetStub().PutPrivateData(StartupInvestorCollection, "AGREEMENT_"+agreementID, agreementJSON)
	if err != nil {
		return fmt.Errorf("failed to accept agreement: %v", err)
	}

	// Copy to ThreePartyCollection for Platform visibility
	err = ctx.GetStub().PutPrivateData(ThreePartyCollection, "AGREEMENT_"+agreementID, agreementJSON)
	if err != nil {
		return fmt.Errorf("failed to share agreement: %v", err)
	}

	return nil
}

// ConfirmFundingCommitment confirms funding commitment
func (i *InvestorContract) ConfirmFundingCommitment(
	ctx contractapi.TransactionContextInterface,
	commitmentID string,
	proposalID string,
	agreementID string,
	campaignID string,
	startupID string,
	investorID string,
	amount float64,
	currency string,
	milestonesJSON string,
) error {

	var milestones []Milestone
	if milestonesJSON != "" {
		json.Unmarshal([]byte(milestonesJSON), &milestones)
	}

	timestamp := time.Now().Format(time.RFC3339)

	commitment := FundingCommitment{
		CommitmentID: commitmentID,
		ProposalID:   proposalID,
		AgreementID:  agreementID,
		CampaignID:   campaignID,
		StartupID:    startupID,
		InvestorID:   investorID,
		Amount:       amount,
		Currency:     currency,
		Milestones:   milestones,
		Status:       "COMMITTED",
		CommittedAt:  timestamp,
	}

	commitmentJSON, err := json.Marshal(commitment)
	if err != nil {
		return fmt.Errorf("failed to marshal commitment: %v", err)
	}

	// Store in ThreePartyCollection
	err = ctx.GetStub().PutPrivateData(ThreePartyCollection, "FUNDING_COMMITMENT_"+commitmentID, commitmentJSON)
	if err != nil {
		return fmt.Errorf("failed to confirm funding: %v", err)
	}

	// Track in private collection
	ctx.GetStub().PutPrivateData(InvestorPrivateCollection, "MY_COMMITMENT_"+commitmentID, commitmentJSON)

	return nil
}

// ConfirmInvestmentToPlatform sends investment confirmation to Platform
func (i *InvestorContract) ConfirmInvestmentToPlatform(
	ctx contractapi.TransactionContextInterface,
	confirmationID string,
	investmentID string,
	campaignID string,
	investorID string,
	amount float64,
	currency string,
) error {

	timestamp := time.Now().Format(time.RFC3339)

	confirmation := InvestmentConfirmation{
		ConfirmationID: confirmationID,
		InvestmentID:   investmentID,
		CampaignID:     campaignID,
		InvestorID:     investorID,
		Amount:         amount,
		Currency:       currency,
		ConfirmedAt:    timestamp,
	}

	confirmationJSON, err := json.Marshal(confirmation)
	if err != nil {
		return fmt.Errorf("failed to marshal confirmation: %v", err)
	}

	// Store in InvestorPlatformCollection
	err = ctx.GetStub().PutPrivateData(InvestorPlatformCollection, "INVESTMENT_CONFIRMATION_"+confirmationID, confirmationJSON)
	if err != nil {
		return fmt.Errorf("failed to send confirmation: %v", err)
	}

	return nil
}

// ============================================================================
// MILESTONE VERIFICATION - Using PDC
// ============================================================================

// VerifyMilestone verifies a milestone completion
// Verification stored in ThreePartyCollection
func (i *InvestorContract) VerifyMilestone(
	ctx contractapi.TransactionContextInterface,
	verificationID string,
	milestoneID string,
	agreementID string,
	campaignID string,
	investorID string,
	approved bool,
	feedback string,
) error {

	timestamp := time.Now().Format(time.RFC3339)

	verification := MilestoneVerification{
		VerificationID: verificationID,
		MilestoneID:    milestoneID,
		AgreementID:    agreementID,
		CampaignID:     campaignID,
		InvestorID:     investorID,
		Approved:       approved,
		Feedback:       feedback,
		VerifiedAt:     timestamp,
	}

	verificationJSON, err := json.Marshal(verification)
	if err != nil {
		return fmt.Errorf("failed to marshal verification: %v", err)
	}

	// Store in ThreePartyCollection (visible to Startup, Investor, Platform)
	err = ctx.GetStub().PutPrivateData(ThreePartyCollection, "MILESTONE_VERIFICATION_"+verificationID, verificationJSON)
	if err != nil {
		return fmt.Errorf("failed to verify milestone: %v", err)
	}

	return nil
}

// ============================================================================
// RISK INSIGHTS - Using PDC
// ============================================================================

// RequestRiskInsights requests risk information from Validator
func (i *InvestorContract) RequestRiskInsights(
	ctx contractapi.TransactionContextInterface,
	requestID string,
	campaignID string,
	investorID string,
) error {

	timestamp := time.Now().Format(time.RFC3339)

	request := RiskInsightRequest{
		RequestID:   requestID,
		CampaignID:  campaignID,
		InvestorID:  investorID,
		Status:      "PENDING",
		RequestedAt: timestamp,
	}

	requestJSON, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}

	// Store in InvestorValidatorCollection
	err = ctx.GetStub().PutPrivateData(InvestorValidatorCollection, "RISK_REQUEST_"+requestID, requestJSON)
	if err != nil {
		return fmt.Errorf("failed to request risk insights: %v", err)
	}

	return nil
}

// RecordRiskInsightResponse records risk insight response from Validator
func (i *InvestorContract) RecordRiskInsightResponse(
	ctx contractapi.TransactionContextInterface,
	responseID string,
	requestID string,
	campaignID string,
	investorID string,
	riskScore float64,
	riskLevel string,
	riskFactors string,
	recommendation string,
) error {

	timestamp := time.Now().Format(time.RFC3339)

	response := RiskInsightResponse{
		ResponseID:     responseID,
		RequestID:      requestID,
		CampaignID:     campaignID,
		InvestorID:     investorID,
		RiskScore:      riskScore,
		RiskLevel:      riskLevel,
		RiskFactors:    riskFactors,
		Recommendation: recommendation,
		ReceivedAt:     timestamp,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %v", err)
	}

	// Store in InvestorValidatorCollection
	err = ctx.GetStub().PutPrivateData(InvestorValidatorCollection, "RISK_RESPONSE_"+responseID, responseJSON)
	if err != nil {
		return fmt.Errorf("failed to record response: %v", err)
	}

	// Update request status
	requestJSON, err := ctx.GetStub().GetPrivateData(InvestorValidatorCollection, "RISK_REQUEST_"+requestID)
	if err == nil && requestJSON != nil {
		var request RiskInsightRequest
		json.Unmarshal(requestJSON, &request)
		request.Status = "FULFILLED"
		request.FulfilledAt = timestamp
		requestJSON, _ = json.Marshal(request)
		ctx.GetStub().PutPrivateData(InvestorValidatorCollection, "RISK_REQUEST_"+requestID, requestJSON)
	}

	return nil
}

// RequestValidationDetails requests validation score and risk insights from Validator
// Investor sends campaignID, Validator responds with validation details
func (i *InvestorContract) RequestValidationDetails(
	ctx contractapi.TransactionContextInterface,
	requestID string,
	campaignID string,
	investorID string,
) error {

	timestamp := time.Now().Format(time.RFC3339)

	// Create validation request
	request := map[string]interface{}{
		"requestId":   requestID,
		"campaignId":  campaignID,
		"investorId":  investorID,
		"requestedAt": timestamp,
		"status":      "PENDING",
	}

	requestJSON, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}

	// Store request in InvestorValidatorCollection so Validator can see it
	err = ctx.GetStub().PutPrivateData(InvestorValidatorCollection, "VALIDATION_REQUEST_"+requestID, requestJSON)
	if err != nil {
		return fmt.Errorf("failed to store validation request: %v", err)
	}

	return nil
}

// GetValidationResponse retrieves validation details response from Validator
func (i *InvestorContract) GetValidationResponse(
	ctx contractapi.TransactionContextInterface,
	requestID string,
) (string, error) {

	// Read response from InvestorValidatorCollection (written by Validator)
	responseJSON, err := ctx.GetStub().GetPrivateData(InvestorValidatorCollection, "VALIDATION_RESPONSE_"+requestID)
	if err != nil || responseJSON == nil {
		return "", fmt.Errorf("validation response not found for request %s", requestID)
	}

	return string(responseJSON), nil
}

// ReceiveRiskInsight is called by ValidatorOrg to provide risk insights
func (i *InvestorContract) ReceiveRiskInsight(
	ctx contractapi.TransactionContextInterface,
	responseID string,
	requestID string,
	campaignID string,
	riskScore float64,
	riskLevel string,
	riskFactors string,
	recommendation string,
) error {
	// This would be called via chaincode-to-chaincode invocation in multi-channel
	// In PDC, Validator writes directly to InvestorValidatorCollection
	return nil
}

// ReceiveCampaignNotification receives campaign notification (from public ledger)
func (i *InvestorContract) ReceiveCampaignNotification(
	ctx contractapi.TransactionContextInterface,
	notificationID string,
	campaignID string,
	projectName string,
	category string,
	goalAmount float64,
) error {

	timestamp := time.Now().Format(time.RFC3339)

	notification := map[string]interface{}{
		"notificationId": notificationID,
		"campaignId":     campaignID,
		"projectName":    projectName,
		"category":       category,
		"goalAmount":     goalAmount,
		"receivedAt":     timestamp,
	}

	notificationJSON, _ := json.Marshal(notification)
	ctx.GetStub().PutPrivateData(InvestorPrivateCollection, "NOTIFICATION_"+notificationID, notificationJSON)

	return nil
}

// ============================================================================
// DISPUTE MANAGEMENT - Using PDC
// ============================================================================

// SubmitDispute submits a dispute
func (i *InvestorContract) SubmitDispute(
	ctx contractapi.TransactionContextInterface,
	disputeID string,
	investorID string,
	disputeType string,
	targetID string,
	targetType string,
	campaignID string,
	agreementID string,
	title string,
	description string,
	claimedAmount float64,
	evidenceHashesJSON string,
) error {

	var evidenceHashes []string
	json.Unmarshal([]byte(evidenceHashesJSON), &evidenceHashes)

	timestamp := time.Now().Format(time.RFC3339)
	submissionID := fmt.Sprintf("DISPUTE_SUB_%s", disputeID)

	dispute := InvestorDisputeSubmission{
		SubmissionID:   submissionID,
		DisputeID:      disputeID,
		InvestorID:     investorID,
		DisputeType:    disputeType,
		TargetID:       targetID,
		TargetType:     targetType,
		CampaignID:     campaignID,
		AgreementID:    agreementID,
		Title:          title,
		Description:    description,
		ClaimedAmount:  claimedAmount,
		EvidenceHashes: evidenceHashes,
		Status:         "SUBMITTED",
		CreatedAt:      timestamp,
	}

	disputeJSON, err := json.Marshal(dispute)
	if err != nil {
		return fmt.Errorf("failed to marshal dispute: %v", err)
	}

	// Store in AllOrgsCollection
	err = ctx.GetStub().PutPrivateData(AllOrgsCollection, "DISPUTE_"+disputeID, disputeJSON)
	if err != nil {
		return fmt.Errorf("failed to submit dispute: %v", err)
	}

	return nil
}

// SubmitDisputeEvidence submits additional evidence
func (i *InvestorContract) SubmitDisputeEvidence(
	ctx contractapi.TransactionContextInterface,
	disputeID string,
	evidenceHashesJSON string,
	evidenceDescription string,
) error {

	var evidenceHashes []string
	json.Unmarshal([]byte(evidenceHashesJSON), &evidenceHashes)

	timestamp := time.Now().Format(time.RFC3339)

	evidence := map[string]interface{}{
		"disputeId":           disputeID,
		"evidenceHashes":      evidenceHashes,
		"evidenceDescription": evidenceDescription,
		"submittedBy":         "investor",
		"submittedAt":         timestamp,
	}

	evidenceJSON, _ := json.Marshal(evidence)
	evidenceKey := fmt.Sprintf("DISPUTE_EVIDENCE_%s_%s", disputeID, timestamp)

	err := ctx.GetStub().PutPrivateData(AllOrgsCollection, evidenceKey, evidenceJSON)
	if err != nil {
		return fmt.Errorf("failed to submit evidence: %v", err)
	}

	return nil
}

// RespondToDispute responds to a dispute filed against investor
func (i *InvestorContract) RespondToDispute(
	ctx contractapi.TransactionContextInterface,
	disputeID string,
	responseText string,
	counterEvidenceHashesJSON string,
) error {

	var counterEvidenceHashes []string
	json.Unmarshal([]byte(counterEvidenceHashesJSON), &counterEvidenceHashes)

	timestamp := time.Now().Format(time.RFC3339)

	response := map[string]interface{}{
		"disputeId":             disputeID,
		"responseText":          responseText,
		"counterEvidenceHashes": counterEvidenceHashes,
		"respondedBy":           "investor",
		"respondedAt":           timestamp,
	}

	responseJSON, _ := json.Marshal(response)

	err := ctx.GetStub().PutPrivateData(AllOrgsCollection, "DISPUTE_RESPONSE_"+disputeID, responseJSON)
	if err != nil {
		return fmt.Errorf("failed to respond to dispute: %v", err)
	}

	return nil
}

// ============================================================================
// REFUND MANAGEMENT - Using PDC
// ============================================================================

// RequestRefund requests a refund
func (i *InvestorContract) RequestRefund(
	ctx contractapi.TransactionContextInterface,
	requestID string,
	investorID string,
	campaignID string,
	agreementID string,
	startupID string,
	originalAmount float64,
	requestedAmount float64,
	refundReason string,
	deductionPercent float64,
) error {

	timestamp := time.Now().Format(time.RFC3339)
	expectedRefund := originalAmount - (originalAmount * deductionPercent / 100)

	refund := RefundRequest{
		RequestID:        requestID,
		InvestorID:       investorID,
		CampaignID:       campaignID,
		AgreementID:      agreementID,
		StartupID:        startupID,
		OriginalAmount:   originalAmount,
		RequestedAmount:  requestedAmount,
		RefundReason:     refundReason,
		DeductionPercent: deductionPercent,
		ExpectedRefund:   expectedRefund,
		Status:           "PENDING",
		RequestedAt:      timestamp,
	}

	refundJSON, err := json.Marshal(refund)
	if err != nil {
		return fmt.Errorf("failed to marshal refund: %v", err)
	}

	// Store in InvestorPlatformCollection
	err = ctx.GetStub().PutPrivateData(InvestorPlatformCollection, "REFUND_REQUEST_"+requestID, refundJSON)
	if err != nil {
		return fmt.Errorf("failed to request refund: %v", err)
	}

	return nil
}

// ============================================================================
// QUERY FUNCTIONS
// ============================================================================

// GetInvestment retrieves an investment
func (i *InvestorContract) GetInvestment(ctx contractapi.TransactionContextInterface, investmentID string) (*Investment, error) {
	investmentJSON, err := ctx.GetStub().GetPrivateData(InvestorPrivateCollection, "MY_INVESTMENT_"+investmentID)
	if err != nil || investmentJSON == nil {
		return nil, fmt.Errorf("investment not found: %v", err)
	}

	var investment Investment
	err = json.Unmarshal(investmentJSON, &investment)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal investment: %v", err)
	}

	return &investment, nil
}

// GetInvestmentsByInvestor retrieves all investments by an investor
func (i *InvestorContract) GetInvestmentsByInvestor(ctx contractapi.TransactionContextInterface, investorID string) (string, error) {
	iterator, err := ctx.GetStub().GetPrivateDataByRange(InvestorPrivateCollection, "MY_INVESTMENT_", "MY_INVESTMENT_~")
	if err != nil {
		return "[]", fmt.Errorf("failed to query investments: %v", err)
	}
	defer iterator.Close()

	var investments []Investment
	for iterator.HasNext() {
		queryResponse, err := iterator.Next()
		if err != nil {
			continue
		}

		var investment Investment
		err = json.Unmarshal(queryResponse.Value, &investment)
		if err != nil {
			continue
		}

		// Filter by investorID
		if investment.InvestorID == investorID {
			investments = append(investments, investment)
		}
	}

	if investments == nil {
		investments = []Investment{}
	}

	investmentsJSON, _ := json.Marshal(investments)
	return string(investmentsJSON), nil
}

// GetMyInvestments retrieves all investments for the current user
// This is an alias for GetInvestmentsByInvestor when called without filter
func (i *InvestorContract) GetMyInvestments(ctx contractapi.TransactionContextInterface) (string, error) {
	iterator, err := ctx.GetStub().GetPrivateDataByRange(InvestorPrivateCollection, "MY_INVESTMENT_", "MY_INVESTMENT_~")
	if err != nil {
		return "[]", fmt.Errorf("failed to query investments: %v", err)
	}
	defer iterator.Close()

	var investments []Investment
	for iterator.HasNext() {
		queryResponse, err := iterator.Next()
		if err != nil {
			continue
		}

		var investment Investment
		err = json.Unmarshal(queryResponse.Value, &investment)
		if err != nil {
			continue
		}

		investments = append(investments, investment)
	}

	if investments == nil {
		investments = []Investment{}
	}

	investmentsJSON, _ := json.Marshal(investments)
	return string(investmentsJSON), nil
}

// GetInvestmentsByCampaign retrieves all investments for a campaign
func (i *InvestorContract) GetInvestmentsByCampaign(ctx contractapi.TransactionContextInterface, campaignID string) (string, error) {
	iterator, err := ctx.GetStub().GetPrivateDataByRange(InvestorPrivateCollection, "MY_INVESTMENT_", "MY_INVESTMENT_~")
	if err != nil {
		return "[]", fmt.Errorf("failed to query investments: %v", err)
	}
	defer iterator.Close()

	var investments []Investment
	for iterator.HasNext() {
		queryResponse, err := iterator.Next()
		if err != nil {
			continue
		}

		var investment Investment
		err = json.Unmarshal(queryResponse.Value, &investment)
		if err != nil {
			continue
		}

		// Filter by campaignID
		if investment.CampaignID == campaignID {
			investments = append(investments, investment)
		}
	}

	if investments == nil {
		investments = []Investment{}
	}

	investmentsJSON, _ := json.Marshal(investments)
	return string(investmentsJSON), nil
}

// GetViewedCampaigns retrieves all viewed campaigns
func (i *InvestorContract) GetViewedCampaigns(ctx contractapi.TransactionContextInterface) (string, error) {
	iterator, err := ctx.GetStub().GetPrivateDataByRange(InvestorPrivateCollection, "VIEWED_", "VIEWED_~")
	if err != nil {
		return "[]", fmt.Errorf("failed to query viewed campaigns: %v", err)
	}
	defer iterator.Close()

	var viewed []map[string]interface{}
	for iterator.HasNext() {
		queryResponse, err := iterator.Next()
		if err != nil {
			continue
		}

		var record map[string]interface{}
		err = json.Unmarshal(queryResponse.Value, &record)
		if err != nil {
			continue
		}

		viewed = append(viewed, record)
	}

	if viewed == nil {
		viewed = []map[string]interface{}{}
	}

	viewedJSON, _ := json.Marshal(viewed)
	return string(viewedJSON), nil
}

// GetViewedCampaignsByInvestor retrieves viewed campaigns for a specific investor
func (i *InvestorContract) GetViewedCampaignsByInvestor(ctx contractapi.TransactionContextInterface, investorID string) (string, error) {
	iterator, err := ctx.GetStub().GetPrivateDataByRange(InvestorPrivateCollection, "VIEWED_", "VIEWED_~")
	if err != nil {
		return "[]", fmt.Errorf("failed to query viewed campaigns: %v", err)
	}
	defer iterator.Close()

	var viewed []map[string]interface{}
	for iterator.HasNext() {
		queryResponse, err := iterator.Next()
		if err != nil {
			continue
		}

		var record map[string]interface{}
		err = json.Unmarshal(queryResponse.Value, &record)
		if err != nil {
			continue
		}

		// Filter by investorID
		if id, ok := record["investorId"].(string); ok && id == investorID {
			viewed = append(viewed, record)
		}
	}

	if viewed == nil {
		viewed = []map[string]interface{}{}
	}

	viewedJSON, _ := json.Marshal(viewed)
	return string(viewedJSON), nil
}

// GetProposal retrieves a proposal
func (i *InvestorContract) GetProposal(ctx contractapi.TransactionContextInterface, proposalID string) (string, error) {
	proposalJSON, err := ctx.GetStub().GetPrivateData(StartupInvestorCollection, "PROPOSAL_"+proposalID)
	if err != nil || proposalJSON == nil {
		return "", fmt.Errorf("proposal not found: %v", err)
	}

	return string(proposalJSON), nil
}

// GetProposalsByCampaign retrieves all proposals for a campaign
func (i *InvestorContract) GetProposalsByCampaign(ctx contractapi.TransactionContextInterface, campaignID string) (string, error) {
	iterator, err := ctx.GetStub().GetPrivateDataByRange(StartupInvestorCollection, "PROPOSAL_", "PROPOSAL_~")
	if err != nil {
		return "[]", fmt.Errorf("failed to query proposals: %v", err)
	}
	defer iterator.Close()

	var proposals []InvestmentProposal
	for iterator.HasNext() {
		queryResponse, err := iterator.Next()
		if err != nil {
			continue
		}

		var proposal InvestmentProposal
		err = json.Unmarshal(queryResponse.Value, &proposal)
		if err != nil {
			continue
		}

		if proposal.CampaignID == campaignID {
			proposals = append(proposals, proposal)
		}
	}

	if proposals == nil {
		proposals = []InvestmentProposal{}
	}

	proposalsJSON, _ := json.Marshal(proposals)
	return string(proposalsJSON), nil
}

// GetProposalsByInvestor retrieves all proposals by an investor
func (i *InvestorContract) GetProposalsByInvestor(ctx contractapi.TransactionContextInterface, investorID string) (string, error) {
	iterator, err := ctx.GetStub().GetPrivateDataByRange(StartupInvestorCollection, "PROPOSAL_", "PROPOSAL_~")
	if err != nil {
		return "[]", fmt.Errorf("failed to query proposals: %v", err)
	}
	defer iterator.Close()

	var proposals []InvestmentProposal
	for iterator.HasNext() {
		queryResponse, err := iterator.Next()
		if err != nil {
			continue
		}

		var proposal InvestmentProposal
		err = json.Unmarshal(queryResponse.Value, &proposal)
		if err != nil {
			continue
		}

		if proposal.InvestorID == investorID {
			proposals = append(proposals, proposal)
		}
	}

	if proposals == nil {
		proposals = []InvestmentProposal{}
	}

	proposalsJSON, _ := json.Marshal(proposals)
	return string(proposalsJSON), nil
}

// GetRefundRequest retrieves a refund request
func (i *InvestorContract) GetRefundRequest(ctx contractapi.TransactionContextInterface, requestID string) (string, error) {
	refundJSON, err := ctx.GetStub().GetPrivateData(InvestorPlatformCollection, "REFUND_REQUEST_"+requestID)
	if err != nil || refundJSON == nil {
		return "", fmt.Errorf("refund request not found: %v", err)
	}

	return string(refundJSON), nil
}

// GetInvestorDisputes retrieves all disputes for an investor
func (i *InvestorContract) GetInvestorDisputes(ctx contractapi.TransactionContextInterface, investorID string) (string, error) {
	iterator, err := ctx.GetStub().GetPrivateDataByRange(AllOrgsCollection, "DISPUTE_", "DISPUTE_~")
	if err != nil {
		return "[]", fmt.Errorf("failed to query disputes: %v", err)
	}
	defer iterator.Close()

	var disputes []InvestorDisputeSubmission
	for iterator.HasNext() {
		queryResponse, err := iterator.Next()
		if err != nil {
			continue
		}

		var dispute InvestorDisputeSubmission
		err = json.Unmarshal(queryResponse.Value, &dispute)
		if err != nil {
			continue
		}

		if dispute.InvestorID == investorID {
			disputes = append(disputes, dispute)
		}
	}

	if disputes == nil {
		disputes = []InvestorDisputeSubmission{}
	}

	disputesJSON, _ := json.Marshal(disputes)
	return string(disputesJSON), nil
}

// PublishInvestmentSummary publishes investment summary to public ledger
func (i *InvestorContract) PublishInvestmentSummary(
	ctx contractapi.TransactionContextInterface,
	summaryID string,
	campaignID string,
	investorCount int,
) error {

	timestamp := time.Now().Format(time.RFC3339)
	summaryData := fmt.Sprintf("%s|%s|%d", summaryID, campaignID, investorCount)
	hash := sha256.Sum256([]byte(summaryData))
	summaryHash := hex.EncodeToString(hash[:])

	summary := map[string]interface{}{
		"summaryId":     summaryID,
		"campaignId":    campaignID,
		"investorCount": investorCount,
		"summaryHash":   summaryHash,
		"publishedAt":   timestamp,
	}

	summaryJSON, _ := json.Marshal(summary)

	// Store on public world state
	err := ctx.GetStub().PutState("INVESTMENT_SUMMARY_"+summaryID, summaryJSON)
	if err != nil {
		return fmt.Errorf("failed to publish summary: %v", err)
	}

	return nil
}

// ============================================================================
// TOKEN INTEGRATION FUNCTIONS (CFT/CFRT)
// ============================================================================

// InvestorFeeRecord represents fee record for investors
type InvestorFeeRecord struct {
	RecordID        string  `json:"recordId"`
	InvestorID      string  `json:"investorId"`
	CampaignID      string  `json:"campaignId"`
	FeeType         string  `json:"feeType"`
	AmountCFT       float64 `json:"amountCFT"`
	TransactionHash string  `json:"transactionHash"`
	PaidAt          string  `json:"paidAt"`
}

// GetInvestorFees returns the fees applicable to investors
func (i *InvestorContract) GetInvestorFees(
	ctx contractapi.TransactionContextInterface,
) (string, error) {

	// Fee schedule based on 1 INR = 2.5 CFT
	fees := map[string]interface{}{
		"investmentFee": map[string]interface{}{
			"percentage":  5.0,
			"description": "5% fee on investment amount (deducted at investment)",
		},
		"disputeFilingFee": map[string]interface{}{
			"amountCFT":   750,
			"amountINR":   300,
			"description": "Fee for filing a dispute",
		},
		"withdrawalFee": map[string]interface{}{
			"percentage":  1.0,
			"description": "1% fee for CFT to Fiat withdrawal",
		},
	}

	feesJSON, err := json.Marshal(fees)
	if err != nil {
		return "", fmt.Errorf("failed to marshal fees: %v", err)
	}

	return string(feesJSON), nil
}

// RecordInvestmentFee records investment fee paid
func (i *InvestorContract) RecordInvestmentFee(
	ctx contractapi.TransactionContextInterface,
	recordID string,
	investorID string,
	campaignID string,
	investmentAmountCFT float64,
	feeAmountCFT float64,
	transactionHash string,
) error {

	timestamp := time.Now().Format(time.RFC3339)

	record := InvestorFeeRecord{
		RecordID:        recordID,
		InvestorID:      investorID,
		CampaignID:      campaignID,
		FeeType:         "INVESTMENT_FEE",
		AmountCFT:       feeAmountCFT,
		TransactionHash: transactionHash,
		PaidAt:          timestamp,
	}

	recordJSON, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal fee record: %v", err)
	}

	err = ctx.GetStub().PutPrivateData(InvestorPrivateCollection, "FEE_PAYMENT_"+recordID, recordJSON)
	if err != nil {
		return fmt.Errorf("failed to record fee: %v", err)
	}

	return nil
}

// GetInvestorFeePayments retrieves all fee payments for an investor
func (i *InvestorContract) GetInvestorFeePayments(
	ctx contractapi.TransactionContextInterface,
	investorID string,
) (string, error) {

	iterator, err := ctx.GetStub().GetPrivateDataByRange(InvestorPrivateCollection, "FEE_PAYMENT_", "FEE_PAYMENT_~")
	if err != nil {
		return "[]", fmt.Errorf("failed to query fee payments: %v", err)
	}
	defer iterator.Close()

	var payments []InvestorFeeRecord
	for iterator.HasNext() {
		queryResponse, err := iterator.Next()
		if err != nil {
			continue
		}

		var record InvestorFeeRecord
		err = json.Unmarshal(queryResponse.Value, &record)
		if err != nil {
			continue
		}

		if record.InvestorID == investorID {
			payments = append(payments, record)
		}
	}

	if payments == nil {
		payments = []InvestorFeeRecord{}
	}

	paymentsJSON, err := json.Marshal(payments)
	if err != nil {
		return "[]", fmt.Errorf("failed to marshal payments: %v", err)
	}

	return string(paymentsJSON), nil
}

// GetDisputePenaltyScheduleForInvestor returns penalty amounts for investor disputes
func (i *InvestorContract) GetDisputePenaltyScheduleForInvestor(
	ctx contractapi.TransactionContextInterface,
) (string, error) {

	// Penalty schedule based on 1 INR = 2.5 CFT
	penalties := map[string]interface{}{
		"fraudulentClaim": map[string]interface{}{
			"penaltyCFT":  500,
			"penaltyINR":  200,
			"description": "Filing fraudulent claims against startups",
		},
		"refundAbuse": map[string]interface{}{
			"penaltyCFT":  375,
			"penaltyINR":  150,
			"description": "Abusing refund mechanisms",
		},
	}

	penaltiesJSON, err := json.Marshal(penalties)
	if err != nil {
		return "", fmt.Errorf("failed to marshal penalties: %v", err)
	}

	return string(penaltiesJSON), nil
}
