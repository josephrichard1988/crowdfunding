package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
	
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// InvestorContract provides functions for InvestorOrg operations
type InvestorContract struct {
	contractapi.Contract
}

// Investment represents an investment commitment
type Investment struct {
	InvestmentID   string  `json:"investmentId"`
	CampaignID     string  `json:"campaignId"`
	InvestorID     string  `json:"investorId"`
	Amount         float64 `json:"amount"`
	Currency       string  `json:"currency"`
	Status         string  `json:"status"` // COMMITTED, ACKNOWLEDGED, CONFIRMED, WITHDRAWN
	CommittedAt    string  `json:"committedAt"`
	ConfirmedAt    string  `json:"confirmedAt"`
	WithdrawnAt    string  `json:"withdrawnAt"`
}

// CampaignView represents campaign details visible to investors
type CampaignView struct {
	CampaignID         string   `json:"campaignId"`
	ProjectName        string   `json:"projectName"`
	Category           string   `json:"category"`
	Description        string   `json:"description"`
	GoalAmount         float64  `json:"goalAmount"`
	FundsRaisedAmount  float64  `json:"fundsRaisedAmount"`
	FundsRaisedPercent float64  `json:"fundsRaisedPercent"`
	Currency           string   `json:"currency"`
	OpenDate           string   `json:"openDate"`
	CloseDate          string   `json:"closeDate"`
	ProductStage       string   `json:"productStage"`
	ProjectType        string   `json:"projectType"`
	Tags               []string `json:"tags"`
	DurationDays       int      `json:"durationDays"`
	ValidationScore    float64  `json:"validationScore"`
	RiskLevel          string   `json:"riskLevel"`
	InvestorCount      int      `json:"investorCount"`
	Status             string   `json:"status"`
	ViewedAt           string   `json:"viewedAt"`
}

// InvestmentProposal represents an investment proposal with terms
// Step 7: Investor sends investment proposal to startup
type InvestmentProposal struct {
	ProposalID       string      `json:"proposalId"`
	CampaignID       string      `json:"campaignId"`
	StartupID        string      `json:"startupId"`
	InvestorID       string      `json:"investorId"`
	InvestmentAmount float64     `json:"investmentAmount"`
	Currency         string      `json:"currency"`
	ProposedTerms    string      `json:"proposedTerms"`
	Milestones       []Milestone `json:"milestones"`
	Status           string      `json:"status"` // PROPOSED, COUNTERED, ACCEPTED, REJECTED, EXPIRED
	NegotiationRound int         `json:"negotiationRound"`
	History          []NegotiationEntry `json:"history"`
	CreatedAt        string      `json:"createdAt"`
	UpdatedAt        string      `json:"updatedAt"`
}

// Milestone for milestone-based fund release
type Milestone struct {
	MilestoneID     string  `json:"milestoneId"`
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	TargetDate      string  `json:"targetDate"`
	FundPercentage  float64 `json:"fundPercentage"` // Percentage of funds released on completion
	Status          string  `json:"status"`         // PENDING, SUBMITTED, VERIFIED, REJECTED
	FundsReleased   bool    `json:"fundsReleased"`
	ReleasedAt      string  `json:"releasedAt"`
}

// NegotiationEntry tracks negotiation history
type NegotiationEntry struct {
	Round       int     `json:"round"`
	Party       string  `json:"party"` // STARTUP or INVESTOR
	Action      string  `json:"action"` // PROPOSE, COUNTER, ACCEPT, REJECT
	Amount      float64 `json:"amount"`
	Terms       string  `json:"terms"`
	Timestamp   string  `json:"timestamp"`
}

// FundingCommitment represents confirmed funding commitment
// Step 10: Investor confirms funding to Platform
type FundingCommitment struct {
	CommitmentID     string      `json:"commitmentId"`
	ProposalID       string      `json:"proposalId"`
	AgreementID      string      `json:"agreementId"`
	CampaignID       string      `json:"campaignId"`
	StartupID        string      `json:"startupId"`
	InvestorID       string      `json:"investorId"`
	Amount           float64     `json:"amount"`
	Currency         string      `json:"currency"`
	Milestones       []Milestone `json:"milestones"`
	Status           string      `json:"status"` // COMMITTED, ESCROWED, PARTIALLY_RELEASED, RELEASED
	CommittedAt      string      `json:"committedAt"`
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
	ResponseID      string  `json:"responseId"`
	RequestID       string  `json:"requestId"`
	CampaignID      string  `json:"campaignId"`
	InvestorID      string  `json:"investorId"`
	RiskScore       float64 `json:"riskScore"`
	RiskLevel       string  `json:"riskLevel"`
	RiskFactors     string  `json:"riskFactors"`
	Recommendation  string  `json:"recommendation"`
	ReceivedAt      string  `json:"receivedAt"`
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

// InvestmentSummaryHash for common-channel (privacy-preserving)
type InvestmentSummaryHash struct {
	SummaryID     string `json:"summaryId"`
	CampaignID    string `json:"campaignId"`
	InvestorCount int    `json:"investorCount"`
	SummaryHash   string `json:"summaryHash"`
	PublishedAt   string `json:"publishedAt"`
}

// InitLedger initializes the InvestorOrg ledger
func (i *InvestorContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	fmt.Println("InvestorOrg contract initialized")
	return nil
}

// ============================================================================
// PLATFORM-INVESTOR-CHANNEL FUNCTIONS - VIEWING & INVESTING
// Endorsed by: InvestorOrg, PlatformOrg
// ============================================================================

// ViewCampaign allows investors to view approved campaign details
// Step 6: Investor views published campaigns from Platform
// Channel: platform-investor-channel
// Endorsers: InvestorOrg, PlatformOrg
func (i *InvestorContract) ViewCampaign(
	ctx contractapi.TransactionContextInterface,
	campaignID string,
	investorID string,
	projectName string,
	category string,
	description string,
	goalAmount float64,
	fundsRaisedAmount float64,
	currency string,
	openDate string,
	closeDate string,
	productStage string,
	projectType string,
	tagsJSON string,
	durationDays int,
	validationScore float64,
	riskLevel string,
	investorCount int,
	status string,
) (string, error) {
	// Parse tags
	var tags []string
	if tagsJSON != "" {
		if err := json.Unmarshal([]byte(tagsJSON), &tags); err != nil {
			return "", fmt.Errorf("failed to parse tags: %v", err)
		}
	}

	// Calculate funds raised percent
	fundsRaisedPercent := float64(0)
	if goalAmount > 0 {
		fundsRaisedPercent = (fundsRaisedAmount / goalAmount) * 100
	}

	// Create campaign view record
	campaignView := CampaignView{
		CampaignID:         campaignID,
		ProjectName:        projectName,
		Category:           category,
		Description:        description,
		GoalAmount:         goalAmount,
		FundsRaisedAmount:  fundsRaisedAmount,
		FundsRaisedPercent: fundsRaisedPercent,
		Currency:           currency,
		OpenDate:           openDate,
		CloseDate:          closeDate,
		ProductStage:       productStage,
		ProjectType:        projectType,
		Tags:               tags,
		DurationDays:       durationDays,
		ValidationScore:    validationScore,
		RiskLevel:          riskLevel,
		InvestorCount:      investorCount,
		Status:             status,
		ViewedAt:           time.Now().Format(time.RFC3339),
	}

	viewJSON, err := json.Marshal(campaignView)
	if err != nil {
		return "", err
	}

	// Store campaign view on startup-investor-channel
	viewKey := fmt.Sprintf("VIEW_%s_%s", campaignID, investorID)
	err = ctx.GetStub().PutState(viewKey, viewJSON)
	if err != nil {
		return "", err
	}

	// Emit event
	eventPayload := map[string]interface{}{
		"campaignId": campaignID,
		"investorId": investorID,
		"channel":    "platform-investor-channel",
		"action":     "CAMPAIGN_VIEWED",
		"timestamp":  campaignView.ViewedAt,
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("CampaignViewed", eventJSON)

	response := map[string]interface{}{
		"message":    "Campaign details retrieved successfully",
		"campaignId": campaignID,
		"campaign":   campaignView,
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// MakeInvestment commits funds to a campaign
// Step 6: Investor invests in campaign through Platform
// Channel: platform-investor-channel
// Endorsers: InvestorOrg, PlatformOrg
func (i *InvestorContract) MakeInvestment(
	ctx contractapi.TransactionContextInterface,
	investmentID string,
	campaignID string,
	investorID string,
	amount float64,
	currency string,
) (string, error) {
	// Check if investment already exists
	existing, err := ctx.GetStub().GetState(investmentID)
	if err != nil {
		return "", fmt.Errorf("failed to read state: %v", err)
	}
	if existing != nil {
		return "", fmt.Errorf("investment %s already exists", investmentID)
	}

	// Create investment record
	investment := Investment{
		InvestmentID: investmentID,
		CampaignID:   campaignID,
		InvestorID:   investorID,
		Amount:       amount,
		Currency:     currency,
		Status:       "COMMITTED",
		CommittedAt:  time.Now().Format(time.RFC3339),
	}

	investmentJSON, err := json.Marshal(investment)
	if err != nil {
		return "", err
	}

	// Store on startup-investor-channel
	err = ctx.GetStub().PutState(investmentID, investmentJSON)
	if err != nil {
		return "", err
	}

	// Store by campaign for aggregation
	campaignInvestmentKey := fmt.Sprintf("CAMPAIGN_INV_%s_%s", campaignID, investmentID)
	err = ctx.GetStub().PutState(campaignInvestmentKey, investmentJSON)
	if err != nil {
		return "", err
	}

	// Emit event
	eventPayload := map[string]interface{}{
		"investmentId": investmentID,
		"campaignId":   campaignID,
		"investorId":   investorID,
		"channel":      "platform-investor-channel",
		"action":       "INVESTMENT_COMMITTED",
		"timestamp":    investment.CommittedAt,
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("InvestmentCommitted", eventJSON)

	response := map[string]interface{}{
		"message":      "Investment committed successfully",
		"investmentId": investmentID,
		"campaignId":   campaignID,
		"amount":       amount,
		"currency":     currency,
		"status":       "COMMITTED",
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// WithdrawInvestment cancels an investment before deadline
// Channel: startup-investor-channel
// Endorsers: StartupOrg, InvestorOrg
func (i *InvestorContract) WithdrawInvestment(
	ctx contractapi.TransactionContextInterface,
	investmentID string,
	reason string,
) (string, error) {
	investmentJSON, err := ctx.GetStub().GetState(investmentID)
	if err != nil {
		return "", fmt.Errorf("failed to read investment: %v", err)
	}
	if investmentJSON == nil {
		return "", fmt.Errorf("investment %s does not exist", investmentID)
	}

	var investment Investment
	err = json.Unmarshal(investmentJSON, &investment)
	if err != nil {
		return "", err
	}

	// Check if already withdrawn or confirmed
	if investment.Status == "WITHDRAWN" {
		return "", fmt.Errorf("investment %s is already withdrawn", investmentID)
	}
	if investment.Status == "CONFIRMED" {
		return "", fmt.Errorf("investment %s is already confirmed and cannot be withdrawn", investmentID)
	}

	// Update investment status
	investment.Status = "WITHDRAWN"
	investment.WithdrawnAt = time.Now().Format(time.RFC3339)

	updatedInvestmentJSON, err := json.Marshal(investment)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(investmentID, updatedInvestmentJSON)
	if err != nil {
		return "", err
	}

	// Emit event
	eventPayload := map[string]interface{}{
		"investmentId": investmentID,
		"campaignId":   investment.CampaignID,
		"investorId":   investment.InvestorID,
		"reason":       reason,
		"channel":      "startup-investor-channel",
		"action":       "INVESTMENT_WITHDRAWN",
		"timestamp":    investment.WithdrawnAt,
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("InvestmentWithdrawn", eventJSON)

	response := map[string]interface{}{
		"message":      "Investment withdrawn successfully",
		"investmentId": investmentID,
		"status":       "WITHDRAWN",
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// CreateInvestmentProposal sends an investment proposal to startup
// Step 7: Investor sends investment proposal with terms and milestones
// Channel: startup-investor-channel
// Endorsers: StartupOrg, InvestorOrg
func (i *InvestorContract) CreateInvestmentProposal(
	ctx contractapi.TransactionContextInterface,
	proposalID string,
	campaignID string,
	startupID string,
	investorID string,
	investmentAmount float64,
	currency string,
	proposedTerms string,
	milestonesJSON string,
) (string, error) {
	// Check if proposal already exists
	existing, err := ctx.GetStub().GetState(proposalID)
	if err != nil {
		return "", fmt.Errorf("failed to read state: %v", err)
	}
	if existing != nil {
		return "", fmt.Errorf("proposal %s already exists", proposalID)
	}

	// Parse milestones
	var milestones []Milestone
	if milestonesJSON != "" {
		if err := json.Unmarshal([]byte(milestonesJSON), &milestones); err != nil {
			return "", fmt.Errorf("failed to parse milestones: %v", err)
		}
	}

	now := time.Now().Format(time.RFC3339)

	// Create negotiation history entry
	historyEntry := NegotiationEntry{
		Round:     1,
		Party:     "INVESTOR",
		Action:    "PROPOSE",
		Amount:    investmentAmount,
		Terms:     proposedTerms,
		Timestamp: now,
	}

	// Create proposal
	proposal := InvestmentProposal{
		ProposalID:       proposalID,
		CampaignID:       campaignID,
		StartupID:        startupID,
		InvestorID:       investorID,
		InvestmentAmount: investmentAmount,
		Currency:         currency,
		ProposedTerms:    proposedTerms,
		Milestones:       milestones,
		Status:           "PROPOSED",
		NegotiationRound: 1,
		History:          []NegotiationEntry{historyEntry},
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	proposalJSON, err := json.Marshal(proposal)
	if err != nil {
		return "", err
	}

	// Store proposal
	err = ctx.GetStub().PutState(proposalID, proposalJSON)
	if err != nil {
		return "", err
	}

	// Store by campaign for lookup
	campaignProposalKey := fmt.Sprintf("PROPOSAL_%s_%s", campaignID, proposalID)
	ctx.GetStub().PutState(campaignProposalKey, proposalJSON)

	// Emit event
	eventPayload := map[string]interface{}{
		"proposalId":       proposalID,
		"campaignId":       campaignID,
		"startupId":        startupID,
		"investorId":       investorID,
		"investmentAmount": investmentAmount,
		"channel":          "startup-investor-channel",
		"action":           "INVESTMENT_PROPOSED",
		"timestamp":        now,
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("InvestmentProposed", eventJSON)

	response := map[string]interface{}{
		"message":     "Investment proposal sent to startup",
		"proposalId":  proposalID,
		"amount":      investmentAmount,
		"status":      "PROPOSED",
		"nextStep":    "Wait for startup response (ACCEPT/REJECT/COUNTER)",
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// RespondToCounterOffer responds to startup's counter-offer
// Step 8-9: Investor responds to counter-proposal
// Channel: startup-investor-channel
// Endorsers: StartupOrg, InvestorOrg
func (i *InvestorContract) RespondToCounterOffer(
	ctx contractapi.TransactionContextInterface,
	proposalID string,
	investorID string,
	response string, // ACCEPT, REJECT, COUNTER
	counterAmount float64,
	counterTerms string,
) (string, error) {
	// Retrieve proposal
	proposalJSON, err := ctx.GetStub().GetState(proposalID)
	if err != nil {
		return "", fmt.Errorf("failed to read proposal: %v", err)
	}
	if proposalJSON == nil {
		return "", fmt.Errorf("proposal %s does not exist", proposalID)
	}

	var proposal InvestmentProposal
	err = json.Unmarshal(proposalJSON, &proposal)
	if err != nil {
		return "", err
	}

	// Verify investor owns this proposal
	if proposal.InvestorID != investorID {
		return "", fmt.Errorf("investor %s is not the owner of proposal %s", investorID, proposalID)
	}

	// Check proposal status
	if proposal.Status != "COUNTERED" {
		return "", fmt.Errorf("proposal is not in COUNTERED status, current: %s", proposal.Status)
	}

	now := time.Now().Format(time.RFC3339)

	// Create history entry
	historyEntry := NegotiationEntry{
		Round:     proposal.NegotiationRound + 1,
		Party:     "INVESTOR",
		Action:    response,
		Amount:    counterAmount,
		Terms:     counterTerms,
		Timestamp: now,
	}
	proposal.History = append(proposal.History, historyEntry)
	proposal.NegotiationRound++
	proposal.UpdatedAt = now

	switch response {
	case "ACCEPT":
		proposal.Status = "ACCEPTED"
		proposal.InvestmentAmount = proposal.InvestmentAmount // Keep last offered amount
	case "REJECT":
		proposal.Status = "REJECTED"
	case "COUNTER":
		proposal.Status = "PROPOSED" // Back to proposed for startup to respond
		proposal.InvestmentAmount = counterAmount
		proposal.ProposedTerms = counterTerms
	default:
		return "", fmt.Errorf("invalid response: %s. Must be ACCEPT, REJECT, or COUNTER", response)
	}

	updatedProposalJSON, err := json.Marshal(proposal)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(proposalID, updatedProposalJSON)
	if err != nil {
		return "", err
	}

	// Emit event
	eventPayload := map[string]interface{}{
		"proposalId":  proposalID,
		"campaignId":  proposal.CampaignID,
		"investorId":  investorID,
		"response":    response,
		"round":       proposal.NegotiationRound,
		"channel":     "startup-investor-channel",
		"action":      "INVESTOR_RESPONDED",
		"timestamp":   now,
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("InvestorResponded", eventJSON)

	responseData := map[string]interface{}{
		"message":    "Response recorded",
		"proposalId": proposalID,
		"response":   response,
		"status":     proposal.Status,
		"round":      proposal.NegotiationRound,
	}
	responseJSON, _ := json.Marshal(responseData)
	return string(responseJSON), nil
}

// AcceptAgreement accepts the final negotiated agreement
// Step 10: Investor accepts agreement and commits funds
// Channel: startup-investor-channel
// Endorsers: StartupOrg, InvestorOrg
func (i *InvestorContract) AcceptAgreement(
	ctx contractapi.TransactionContextInterface,
	proposalID string,
	agreementID string,
	investorID string,
) (string, error) {
	// Retrieve proposal
	proposalJSON, err := ctx.GetStub().GetState(proposalID)
	if err != nil {
		return "", fmt.Errorf("failed to read proposal: %v", err)
	}
	if proposalJSON == nil {
		return "", fmt.Errorf("proposal %s does not exist", proposalID)
	}

	var proposal InvestmentProposal
	err = json.Unmarshal(proposalJSON, &proposal)
	if err != nil {
		return "", err
	}

	// Check proposal status
	if proposal.Status != "ACCEPTED" {
		return "", fmt.Errorf("proposal must be ACCEPTED before creating agreement, current: %s", proposal.Status)
	}

	now := time.Now().Format(time.RFC3339)

	// Store agreement marker (actual agreement is on Platform)
	agreementMarker := map[string]interface{}{
		"agreementId":      agreementID,
		"proposalId":       proposalID,
		"campaignId":       proposal.CampaignID,
		"startupId":        proposal.StartupID,
		"investorId":       investorID,
		"investmentAmount": proposal.InvestmentAmount,
		"currency":         proposal.Currency,
		"investorAccepted": true,
		"acceptedAt":       now,
	}
	markerJSON, _ := json.Marshal(agreementMarker)
	ctx.GetStub().PutState(fmt.Sprintf("AGREEMENT_INV_%s", agreementID), markerJSON)

	// Emit event for Platform to witness
	eventPayload := map[string]interface{}{
		"agreementId":      agreementID,
		"proposalId":       proposalID,
		"campaignId":       proposal.CampaignID,
		"investorId":       investorID,
		"investmentAmount": proposal.InvestmentAmount,
		"channel":          "startup-investor-channel",
		"action":           "INVESTOR_ACCEPTED_AGREEMENT",
		"timestamp":        now,
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("InvestorAcceptedAgreement", eventJSON)

	responseData := map[string]interface{}{
		"message":      "Agreement accepted. Platform will witness and create escrow.",
		"agreementId":  agreementID,
		"proposalId":   proposalID,
		"amount":       proposal.InvestmentAmount,
		"nextStep":     "Platform to witness agreement and hold funds in escrow",
	}
	responseJSON, _ := json.Marshal(responseData)
	return string(responseJSON), nil
}

// ConfirmFundingCommitment confirms funding commitment to Platform for escrow
// Step 10: Investor commits funds to escrow via Platform
// Channel: investor-platform-channel
// Endorsers: InvestorOrg, PlatformOrg
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
) (string, error) {
	// Parse milestones
	var milestones []Milestone
	if milestonesJSON != "" {
		if err := json.Unmarshal([]byte(milestonesJSON), &milestones); err != nil {
			return "", fmt.Errorf("failed to parse milestones: %v", err)
		}
	}

	now := time.Now().Format(time.RFC3339)

	// Create funding commitment
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
		CommittedAt:  now,
	}

	commitmentJSON, err := json.Marshal(commitment)
	if err != nil {
		return "", err
	}

	// Store commitment
	err = ctx.GetStub().PutState(commitmentID, commitmentJSON)
	if err != nil {
		return "", err
	}

	// Emit event for Platform
	eventPayload := map[string]interface{}{
		"commitmentId": commitmentID,
		"agreementId":  agreementID,
		"campaignId":   campaignID,
		"investorId":   investorID,
		"amount":       amount,
		"channel":      "investor-platform-channel",
		"action":       "FUNDING_COMMITTED",
		"timestamp":    now,
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("FundingCommitted", eventJSON)

	responseData := map[string]interface{}{
		"message":      "Funding committed to escrow",
		"commitmentId": commitmentID,
		"agreementId":  agreementID,
		"amount":       amount,
		"status":       "COMMITTED",
		"nextStep":     "Funds held in escrow. Will be released on milestone completion.",
	}
	responseJSON, _ := json.Marshal(responseData)
	return string(responseJSON), nil
}

// VerifyMilestone verifies and approves a milestone completion
// Step 13: Investor verifies milestone for fund release (multi-party visibility)
// Channel: common-channel
// Endorsers: InvestorOrg (multi-party visibility)
func (i *InvestorContract) VerifyMilestone(
	ctx contractapi.TransactionContextInterface,
	verificationID string,
	milestoneID string,
	agreementID string,
	campaignID string,
	investorID string,
	approved bool,
	feedback string,
) (string, error) {
	now := time.Now().Format(time.RFC3339)

	// Create verification record
	verification := MilestoneVerification{
		VerificationID: verificationID,
		MilestoneID:    milestoneID,
		AgreementID:    agreementID,
		CampaignID:     campaignID,
		InvestorID:     investorID,
		Approved:       approved,
		Feedback:       feedback,
		VerifiedAt:     now,
	}

	verificationJSON, err := json.Marshal(verification)
	if err != nil {
		return "", err
	}

	// Store verification
	err = ctx.GetStub().PutState(verificationID, verificationJSON)
	if err != nil {
		return "", err
	}

	// Emit event for Platform to release funds if approved
	eventPayload := map[string]interface{}{
		"verificationId": verificationID,
		"milestoneId":    milestoneID,
		"agreementId":    agreementID,
		"campaignId":     campaignID,
		"investorId":     investorID,
		"approved":       approved,
		"channel":        "common-channel",
		"action":         "MILESTONE_VERIFIED",
		"timestamp":      now,
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("MilestoneVerified", eventJSON)

	status := "REJECTED"
	if approved {
		status = "APPROVED"
	}

	responseData := map[string]interface{}{
		"message":        "Milestone verification recorded",
		"verificationId": verificationID,
		"milestoneId":    milestoneID,
		"approved":       approved,
		"status":         status,
		"nextStep":       "Platform to release milestone funds from escrow",
	}
	responseJSON, _ := json.Marshal(responseData)
	return string(responseJSON), nil
}

// ============================================================================
// INVESTOR-VALIDATOR-CHANNEL FUNCTIONS
// Endorsed by: InvestorOrg, ValidatorOrg
// ============================================================================

// RequestRiskInsights requests risk information from ValidatorOrg
// Channel: investor-validator-channel
// Endorsers: InvestorOrg, ValidatorOrg
func (i *InvestorContract) RequestRiskInsights(
	ctx contractapi.TransactionContextInterface,
	requestID string,
	campaignID string,
	investorID string,
) (string, error) {
	// Create risk insight request
	request := RiskInsightRequest{
		RequestID:   requestID,
		CampaignID:  campaignID,
		InvestorID:  investorID,
		Status:      "PENDING",
		RequestedAt: time.Now().Format(time.RFC3339),
	}

	requestJSON, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	// Store on investor-validator-channel
	err = ctx.GetStub().PutState(requestID, requestJSON)
	if err != nil {
		return "", err
	}

	// Emit event for ValidatorOrg
	eventPayload := map[string]interface{}{
		"requestId":  requestID,
		"campaignId": campaignID,
		"investorId": investorID,
		"channel":    "investor-validator-channel",
		"action":     "RISK_INSIGHTS_REQUESTED",
		"timestamp":  request.RequestedAt,
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("RiskInsightsRequested", eventJSON)

	response := map[string]interface{}{
		"message":    "Risk insights request sent to ValidatorOrg",
		"requestId":  requestID,
		"campaignId": campaignID,
		"status":     "PENDING",
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// RecordRiskInsightResponse records the risk insight response from Validator
// Channel: investor-validator-channel
// Endorsers: InvestorOrg, ValidatorOrg
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
) (string, error) {
	now := time.Now().Format(time.RFC3339)

	// Create response record
	riskResponse := RiskInsightResponse{
		ResponseID:     responseID,
		RequestID:      requestID,
		CampaignID:     campaignID,
		InvestorID:     investorID,
		RiskScore:      riskScore,
		RiskLevel:      riskLevel,
		RiskFactors:    riskFactors,
		Recommendation: recommendation,
		ReceivedAt:     now,
	}

	responseJSON, err := json.Marshal(riskResponse)
	if err != nil {
		return "", err
	}

	// Store response
	err = ctx.GetStub().PutState(responseID, responseJSON)
	if err != nil {
		return "", err
	}

	// Update original request status
	requestJSON, err := ctx.GetStub().GetState(requestID)
	if err == nil && requestJSON != nil {
		var request RiskInsightRequest
		json.Unmarshal(requestJSON, &request)
		request.Status = "FULFILLED"
		request.FulfilledAt = now
		updatedRequestJSON, _ := json.Marshal(request)
		ctx.GetStub().PutState(requestID, updatedRequestJSON)
	}

	// Emit event
	eventPayload := map[string]interface{}{
		"responseId":  responseID,
		"requestId":   requestID,
		"campaignId":  campaignID,
		"investorId":  investorID,
		"riskLevel":   riskLevel,
		"channel":     "investor-validator-channel",
		"action":      "RISK_INSIGHTS_RECEIVED",
		"timestamp":   now,
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("RiskInsightsReceived", eventJSON)

	responseData := map[string]interface{}{
		"message":        "Risk insights received",
		"responseId":     responseID,
		"riskScore":      riskScore,
		"riskLevel":      riskLevel,
		"recommendation": recommendation,
	}
	responseDataJSON, _ := json.Marshal(responseData)
	return string(responseDataJSON), nil
}

// ============================================================================
// INVESTOR-PLATFORM-CHANNEL FUNCTIONS
// Endorsed by: InvestorOrg, PlatformOrg
// ============================================================================

// ConfirmInvestmentToPlatform sends investment confirmation to PlatformOrg
// Channel: investor-platform-channel
// Endorsers: InvestorOrg, PlatformOrg
func (i *InvestorContract) ConfirmInvestmentToPlatform(
	ctx contractapi.TransactionContextInterface,
	confirmationID string,
	investmentID string,
	campaignID string,
	investorID string,
	amount float64,
	currency string,
) (string, error) {
	// Create confirmation record
	confirmation := InvestmentConfirmation{
		ConfirmationID: confirmationID,
		InvestmentID:   investmentID,
		CampaignID:     campaignID,
		InvestorID:     investorID,
		Amount:         amount,
		Currency:       currency,
		ConfirmedAt:    time.Now().Format(time.RFC3339),
	}

	confirmationJSON, err := json.Marshal(confirmation)
	if err != nil {
		return "", err
	}

	// Store on investor-platform-channel
	err = ctx.GetStub().PutState(confirmationID, confirmationJSON)
	if err != nil {
		return "", err
	}

	// Store by campaign for platform lookup
	platformConfirmKey := fmt.Sprintf("PLATFORM_CONFIRM_%s_%s", campaignID, investmentID)
	err = ctx.GetStub().PutState(platformConfirmKey, confirmationJSON)
	if err != nil {
		return "", err
	}

	// Update original investment status
	investmentJSON, err := ctx.GetStub().GetState(investmentID)
	if err == nil && investmentJSON != nil {
		var investment Investment
		json.Unmarshal(investmentJSON, &investment)
		investment.Status = "CONFIRMED"
		investment.ConfirmedAt = confirmation.ConfirmedAt
		updatedInvestmentJSON, _ := json.Marshal(investment)
		ctx.GetStub().PutState(investmentID, updatedInvestmentJSON)
	}

	// Emit event
	eventPayload := map[string]interface{}{
		"confirmationId": confirmationID,
		"investmentId":   investmentID,
		"campaignId":     campaignID,
		"investorId":     investorID,
		"channel":        "investor-platform-channel",
		"action":         "INVESTMENT_CONFIRMED_TO_PLATFORM",
		"timestamp":      confirmation.ConfirmedAt,
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("InvestmentConfirmedToPlatform", eventJSON)

	response := map[string]interface{}{
		"message":        "Investment confirmation sent to PlatformOrg",
		"confirmationId": confirmationID,
		"investmentId":   investmentID,
		"campaignId":     campaignID,
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// ============================================================================
// COMMON-CHANNEL FUNCTIONS
// Read by: All Orgs, Write by: InvestorOrg (aggregated/hash only - privacy preserving)
// ============================================================================

// PublishInvestmentSummary publishes aggregated investment summary to common-channel
// Channel: common-channel
// Purpose: Privacy-preserving summary (only count and hash, never amounts or identities)
func (i *InvestorContract) PublishInvestmentSummary(
	ctx contractapi.TransactionContextInterface,
	summaryID string,
	campaignID string,
	investorCount int,
) (string, error) {
	// Generate summary hash (no sensitive data - no amounts or investor identities)
	summaryData := map[string]interface{}{
		"summaryId":     summaryID,
		"campaignId":    campaignID,
		"investorCount": investorCount,
		"timestamp":     time.Now().Format(time.RFC3339),
	}
	summaryDataJSON, _ := json.Marshal(summaryData)
	summaryHash := generateHash(string(summaryDataJSON))

	// Create investment summary for common channel
	summary := InvestmentSummaryHash{
		SummaryID:     summaryID,
		CampaignID:    campaignID,
		InvestorCount: investorCount,
		SummaryHash:   summaryHash,
		PublishedAt:   time.Now().Format(time.RFC3339),
	}

	summaryJSON, err := json.Marshal(summary)
	if err != nil {
		return "", err
	}

	// Store on common-channel
	commonKey := fmt.Sprintf("COMMON_INVESTMENT_%s", campaignID)
	err = ctx.GetStub().PutState(commonKey, summaryJSON)
	if err != nil {
		return "", err
	}

	// Emit event
	eventPayload := map[string]interface{}{
		"summaryId":     summaryID,
		"campaignId":    campaignID,
		"investorCount": investorCount,
		"summaryHash":   summaryHash,
		"channel":       "common-channel",
		"action":        "INVESTMENT_SUMMARY_PUBLISHED",
		"timestamp":     summary.PublishedAt,
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("InvestmentSummaryPublished", eventJSON)

	response := map[string]interface{}{
		"message":       "Investment summary published to common channel",
		"summaryId":     summaryID,
		"campaignId":    campaignID,
		"investorCount": investorCount,
		"summaryHash":   summaryHash,
		"channel":       "common-channel",
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// ============================================================================
// QUERY FUNCTIONS
// ============================================================================

// GetInvestment retrieves investment by ID
func (i *InvestorContract) GetInvestment(ctx contractapi.TransactionContextInterface, investmentID string) (*Investment, error) {
	investmentJSON, err := ctx.GetStub().GetState(investmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to read investment: %v", err)
	}
	if investmentJSON == nil {
		return nil, fmt.Errorf("investment %s does not exist", investmentID)
	}

	var investment Investment
	err = json.Unmarshal(investmentJSON, &investment)
	if err != nil {
		return nil, err
	}

	return &investment, nil
}

// GetInvestmentsByInvestor returns all investments by investor
func (i *InvestorContract) GetInvestmentsByInvestor(ctx contractapi.TransactionContextInterface, investorID string) (string, error) {
	queryString := fmt.Sprintf(`{"selector":{"investorId":"%s"}}`, investorID)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return "", err
	}
	defer resultsIterator.Close()

	var investments []map[string]interface{}
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return "", err
		}

		var investment Investment
		err = json.Unmarshal(queryResponse.Value, &investment)
		if err != nil {
			continue
		}

		investmentMap := map[string]interface{}{
			"Key":    queryResponse.Key,
			"Record": investment,
		}
		investments = append(investments, investmentMap)
	}

	investmentsJSON, err := json.Marshal(investments)
	if err != nil {
		return "", err
	}

	return string(investmentsJSON), nil
}

// GetInvestmentsByCampaign returns all investments for a campaign
func (i *InvestorContract) GetInvestmentsByCampaign(ctx contractapi.TransactionContextInterface, campaignID string) (string, error) {
	queryString := fmt.Sprintf(`{"selector":{"campaignId":"%s"}}`, campaignID)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return "", err
	}
	defer resultsIterator.Close()

	var investments []map[string]interface{}
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return "", err
		}

		var investment Investment
		err = json.Unmarshal(queryResponse.Value, &investment)
		if err != nil {
			continue
		}

		investmentMap := map[string]interface{}{
			"Key":    queryResponse.Key,
			"Record": investment,
		}
		investments = append(investments, investmentMap)
	}

	investmentsJSON, err := json.Marshal(investments)
	if err != nil {
		return "", err
	}

	return string(investmentsJSON), nil
}

// GetProposal retrieves an investment proposal by ID
// Channel: startup-investor-channel
func (i *InvestorContract) GetProposal(ctx contractapi.TransactionContextInterface, proposalID string) (*InvestmentProposal, error) {
	proposalJSON, err := ctx.GetStub().GetState(proposalID)
	if err != nil {
		return nil, fmt.Errorf("failed to read proposal: %v", err)
	}
	if proposalJSON == nil {
		return nil, fmt.Errorf("proposal %s does not exist", proposalID)
	}

	var proposal InvestmentProposal
	err = json.Unmarshal(proposalJSON, &proposal)
	if err != nil {
		return nil, err
	}

	return &proposal, nil
}

// GetProposalsByCampaign retrieves all proposals for a campaign
// Channel: startup-investor-channel
func (i *InvestorContract) GetProposalsByCampaign(ctx contractapi.TransactionContextInterface, campaignID string) (string, error) {
	queryString := fmt.Sprintf(`{"selector":{"campaignId":"%s","proposalId":{"$exists":true}}}`, campaignID)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return "", err
	}
	defer resultsIterator.Close()

	var proposals []map[string]interface{}
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return "", err
		}

		var proposal InvestmentProposal
		err = json.Unmarshal(queryResponse.Value, &proposal)
		if err != nil {
			continue
		}

		proposalMap := map[string]interface{}{
			"Key":    queryResponse.Key,
			"Record": proposal,
		}
		proposals = append(proposals, proposalMap)
	}

	proposalsJSON, err := json.Marshal(proposals)
	if err != nil {
		return "", err
	}

	return string(proposalsJSON), nil
}

// GetProposalsByInvestor retrieves all proposals by an investor
// Channel: startup-investor-channel
func (i *InvestorContract) GetProposalsByInvestor(ctx contractapi.TransactionContextInterface, investorID string) (string, error) {
	queryString := fmt.Sprintf(`{"selector":{"investorId":"%s","proposalId":{"$exists":true}}}`, investorID)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return "", err
	}
	defer resultsIterator.Close()

	var proposals []map[string]interface{}
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return "", err
		}

		var proposal InvestmentProposal
		err = json.Unmarshal(queryResponse.Value, &proposal)
		if err != nil {
			continue
		}

		proposalMap := map[string]interface{}{
			"Key":    queryResponse.Key,
			"Record": proposal,
		}
		proposals = append(proposals, proposalMap)
	}

	proposalsJSON, err := json.Marshal(proposals)
	if err != nil {
		return "", err
	}

	return string(proposalsJSON), nil
}

// ============================================================================
// CROSS-CHANNEL INVOCATION HELPER FUNCTIONS
// ============================================================================

// InvokeStartupOrgAcknowledge notifies StartupOrg about investment
// Cross-channel call to common-channel (Step 11)
func (i *InvestorContract) InvokeStartupOrgAcknowledge(
	ctx contractapi.TransactionContextInterface,
	investmentID string,
	campaignID string,
	investorID string,
	amount string,
	currency string,
) (string, error) {
	args := [][]byte{
		[]byte("AcknowledgeInvestment"),
		[]byte(investmentID),
		[]byte(campaignID),
		[]byte(investorID),
		[]byte(amount),
		[]byte(currency),
	}

	response := ctx.GetStub().InvokeChaincode(
		"startuporg",
		args,
		"common-channel",
	)

	if response.Status != 200 {
		return "", fmt.Errorf("cross-channel invoke to StartupOrg failed: %s", response.Message)
	}

	// Emit cross-channel event
	eventPayload := map[string]interface{}{
		"investmentId":   investmentID,
		"campaignId":     campaignID,
		"targetChannel":  "common-channel",
		"targetContract": "startuporg",
		"action":         "CROSS_CHANNEL_ACKNOWLEDGE",
		"timestamp":      time.Now().Format(time.RFC3339),
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("CrossChannelInvoke", eventJSON)

	return string(response.Payload), nil
}

// InvokeValidatorOrgRequestRisk requests risk info from ValidatorOrg
// Cross-channel call to investor-validator-channel
func (i *InvestorContract) InvokeValidatorOrgRequestRisk(
	ctx contractapi.TransactionContextInterface,
	requestID string,
	campaignID string,
	investorID string,
) (string, error) {
	args := [][]byte{
		[]byte("GetRiskInsight"),
		[]byte(requestID),
		[]byte(campaignID),
		[]byte(investorID),
	}

	response := ctx.GetStub().InvokeChaincode(
		"validatororg",
		args,
		"investor-validator-channel",
	)

	if response.Status != 200 {
		return "", fmt.Errorf("cross-channel query to ValidatorOrg failed: %s", response.Message)
	}

	return string(response.Payload), nil
}

// InvokePlatformOrgConfirm sends investment confirmation to PlatformOrg
// Cross-channel call to investor-platform-channel
func (i *InvestorContract) InvokePlatformOrgConfirm(
	ctx contractapi.TransactionContextInterface,
	recordID string,
	confirmationID string,
	campaignID string,
	investorID string,
	amount string,
	currency string,
) (string, error) {
	args := [][]byte{
		[]byte("RecordInvestorConfirmation"),
		[]byte(recordID),
		[]byte(confirmationID),
		[]byte(campaignID),
		[]byte(investorID),
		[]byte(amount),
		[]byte(currency),
	}

	response := ctx.GetStub().InvokeChaincode(
		"platformorg",
		args,
		"investor-platform-channel",
	)

	if response.Status != 200 {
		return "", fmt.Errorf("cross-channel invoke to PlatformOrg failed: %s", response.Message)
	}

	// Emit cross-channel event
	eventPayload := map[string]interface{}{
		"confirmationId": confirmationID,
		"campaignId":     campaignID,
		"targetChannel":  "investor-platform-channel",
		"targetContract": "platformorg",
		"action":         "CROSS_CHANNEL_CONFIRM",
		"timestamp":      time.Now().Format(time.RFC3339),
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("CrossChannelInvoke", eventJSON)

	return string(response.Payload), nil
}

// ReceiveRiskInsight receives risk insight from ValidatorOrg (called via InvokeChaincode)
func (i *InvestorContract) ReceiveRiskInsight(
	ctx contractapi.TransactionContextInterface,
	insightID string,
	campaignID string,
	riskScore string,
	riskLevel string,
	recommendation string,
) (string, error) {
	// Store received risk insight
	insight := map[string]interface{}{
		"insightId":      insightID,
		"campaignId":     campaignID,
		"riskScore":      riskScore,
		"riskLevel":      riskLevel,
		"recommendation": recommendation,
		"receivedAt":     time.Now().Format(time.RFC3339),
	}

	insightJSON, err := json.Marshal(insight)
	if err != nil {
		return "", err
	}

	key := fmt.Sprintf("RISK_INSIGHT_%s", campaignID)
	err = ctx.GetStub().PutState(key, insightJSON)
	if err != nil {
		return "", err
	}

	return string(insightJSON), nil
}

// ReceiveCampaignNotification receives campaign status updates from StartupOrg
func (i *InvestorContract) ReceiveCampaignNotification(
	ctx contractapi.TransactionContextInterface,
	campaignID string,
	status string,
	message string,
) (string, error) {
	notification := map[string]interface{}{
		"campaignId": campaignID,
		"status":     status,
		"message":    message,
		"receivedAt": time.Now().Format(time.RFC3339),
	}

	notificationJSON, err := json.Marshal(notification)
	if err != nil {
		return "", err
	}

	key := fmt.Sprintf("NOTIFICATION_%s_%s", campaignID, time.Now().Format("20060102150405"))
	err = ctx.GetStub().PutState(key, notificationJSON)
	if err != nil {
		return "", err
	}

	return string(notificationJSON), nil
}

// generateHash generates SHA256 hash
func generateHash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func main() {
	investorChaincode, err := contractapi.NewChaincode(&InvestorContract{})
	if err != nil {
		fmt.Printf("Error creating InvestorOrg chaincode: %v\n", err)
		return
	}

	if err := investorChaincode.Start(); err != nil {
		fmt.Printf("Error starting InvestorOrg chaincode: %v\n", err)
	}
}
