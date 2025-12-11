package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// PlatformContract provides functions for PlatformOrg operations
type PlatformContract struct {
	contractapi.Contract
}

// ============================================================================
// DATA STRUCTURES
// ============================================================================

// PublishedCampaign represents a campaign published on the platform portal
type PublishedCampaign struct {
	CampaignID          string      `json:"campaignId"`
	StartupID           string      `json:"startupId"`
	ProjectName         string      `json:"projectName"`
	Category            string      `json:"category"`
	Description         string      `json:"description"`
	GoalAmount          float64     `json:"goalAmount"`
	FundsRaisedAmount   float64     `json:"fundsRaisedAmount"`
	FundsRaisedPercent  float64     `json:"fundsRaisedPercent"`
	Currency            string      `json:"currency"`
	OpenDate            string      `json:"openDate"`
	CloseDate           string      `json:"closeDate"`
	DurationDays        int         `json:"durationDays"`
	ValidationScore     float64     `json:"validationScore"`
	ValidationHash      string      `json:"validationHash"` // Hash verified with ValidatorOrg
	ValidationVerified  bool        `json:"validationVerified"`
	Status              string      `json:"status"` // PENDING_VERIFICATION, PUBLISHED, ACTIVE, FUNDED, COMPLETED, CLOSED
	InvestorCount       int         `json:"investorCount"`
	TotalConfirmed      float64     `json:"totalConfirmed"`
	Milestones          []Milestone `json:"milestones"`
	AgreementIDs        []string    `json:"agreementIds"`
	PublishedAt         string      `json:"publishedAt"`
	UpdatedAt           string      `json:"updatedAt"`
}

// Milestone represents a funding milestone
type Milestone struct {
	MilestoneID   string  `json:"milestoneId"`
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	TargetAmount  float64 `json:"targetAmount"`
	TargetDate    string  `json:"targetDate"`
	Status        string  `json:"status"` // PENDING, IN_PROGRESS, COMPLETED, VERIFIED
	FundsReleased bool    `json:"fundsReleased"`
	ReleasedAt    string  `json:"releasedAt"`
}

// Agreement represents investment agreement (Platform as witness)
type Agreement struct {
	AgreementID       string      `json:"agreementId"`
	CampaignID        string      `json:"campaignId"`
	StartupID         string      `json:"startupId"`
	InvestorID        string      `json:"investorId"`
	InvestmentAmount  float64     `json:"investmentAmount"`
	Currency          string      `json:"currency"`
	Milestones        []Milestone `json:"milestones"`
	Terms             string      `json:"terms"`
	Status            string      `json:"status"` // PROPOSED, NEGOTIATING, ACCEPTED, ACTIVE, COMPLETED, CANCELLED
	StartupAccepted   bool        `json:"startupAccepted"`
	InvestorAccepted  bool        `json:"investorAccepted"`
	PlatformWitnessed bool        `json:"platformWitnessed"`
	WitnessedAt       string      `json:"witnessedAt"`
	CreatedAt         string      `json:"createdAt"`
	AcceptedAt        string      `json:"acceptedAt"`
}

// FundEscrow represents funds held in escrow by Platform
type FundEscrow struct {
	EscrowID     string  `json:"escrowId"`
	AgreementID  string  `json:"agreementId"`
	CampaignID   string  `json:"campaignId"`
	InvestorID   string  `json:"investorId"`
	StartupID    string  `json:"startupId"`
	TotalAmount  float64 `json:"totalAmount"`
	ReleasedAmount float64 `json:"releasedAmount"`
	HeldAmount   float64 `json:"heldAmount"`
	Currency     string  `json:"currency"`
	Status       string  `json:"status"` // ACTIVE, PARTIALLY_RELEASED, FULLY_RELEASED, REFUNDED
	CreatedAt    string  `json:"createdAt"`
	UpdatedAt    string  `json:"updatedAt"`
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

// FundRelease represents fund release to startup (milestone-based)
type FundRelease struct {
	ReleaseID     string  `json:"releaseId"`
	EscrowID      string  `json:"escrowId"`
	AgreementID   string  `json:"agreementId"`
	CampaignID    string  `json:"campaignId"`
	MilestoneID   string  `json:"milestoneId"`
	StartupID     string  `json:"startupId"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	Status        string  `json:"status"` // PENDING, RELEASED
	TriggerReason string  `json:"triggerReason"`
	ReleasedAt    string  `json:"releasedAt"`
}

// CampaignClosure represents campaign closure record
type CampaignClosure struct {
	ClosureID          string  `json:"closureId"`
	CampaignID         string  `json:"campaignId"`
	FinalStatus        string  `json:"finalStatus"` // SUCCESSFUL, FAILED, CANCELLED
	FinalAmount        float64 `json:"finalAmount"`
	FinalInvestorCount int     `json:"finalInvestorCount"`
	ClosureReason      string  `json:"closureReason"`
	ClosedAt           string  `json:"closedAt"`
}

// GlobalMetrics for common-channel (privacy-preserving)
type GlobalMetrics struct {
	MetricsID           string `json:"metricsId"`
	TotalCampaigns      int    `json:"totalCampaigns"`
	ActiveCampaigns     int    `json:"activeCampaigns"`
	SuccessfulCampaigns int    `json:"successfulCampaigns"`
	TotalInvestorCount  int    `json:"totalInvestorCount"`
	MetricsHash         string `json:"metricsHash"`
	PublishedAt         string `json:"publishedAt"`
}

// InvestorQuery represents investor's query to Platform
type InvestorQuery struct {
	QueryID     string   `json:"queryId"`
	CampaignID  string   `json:"campaignId"`
	InvestorID  string   `json:"investorId"`
	Questions   []string `json:"questions"`
	Status      string   `json:"status"` // PENDING, ASSIGNED, ANSWERED
	AssignedValidators []string `json:"assignedValidators"`
	Responses   []string `json:"responses"`
	CreatedAt   string   `json:"createdAt"`
	AnsweredAt  string   `json:"answeredAt"`
}

// InitLedger initializes the PlatformOrg ledger
func (p *PlatformContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	fmt.Println("PlatformOrg contract initialized - Merged Version")
	return nil
}

// ============================================================================
// COMMON-CHANNEL FUNCTIONS - PUBLISHING & AGREEMENTS
// Endorsed by: All orgs for multi-party visibility
// ============================================================================

// PublishCampaignToPortal makes validated campaign visible to investors
// Step 5: Platform publishes validated campaign to all orgs
// Channel: common-channel
// Endorsers: PlatformOrg (multi-party visibility)
// REFACTORED: Now only requires campaignID - fetches all data from ledger
func (p *PlatformContract) PublishCampaignToPortal(
	ctx contractapi.TransactionContextInterface,
	campaignID string,
) (string, error) {
	// Check if already published
	existing, err := ctx.GetStub().GetState(campaignID)
	if err != nil {
		return "", fmt.Errorf("failed to read state: %v", err)
	}
	if existing != nil {
		return "", fmt.Errorf("campaign %s already published", campaignID)
	}

	// Step 1: Get the validator decision record from validator-platform-channel
	// This was recorded when Platform received the validation report
	decisionKey := fmt.Sprintf("DECISION_%s", campaignID)
	decisionJSON, err := ctx.GetStub().GetState(decisionKey)
	if err != nil {
		return "", fmt.Errorf("failed to read validator decision: %v", err)
	}
	if decisionJSON == nil {
		return "", fmt.Errorf("no validator decision found for campaign %s. Record validator decision first", campaignID)
	}

	var decision ValidatorDecisionRecord
	err = json.Unmarshal(decisionJSON, &decision)
	if err != nil {
		return "", fmt.Errorf("failed to parse validator decision: %v", err)
	}

	// Verify campaign was approved
	if !decision.Approved {
		return "", fmt.Errorf("campaign %s was not approved by validator. Cannot publish", campaignID)
	}

	// Step 2: Cross-channel query to get campaign details from StartupOrg
	// via startup-platform-channel
	campaignDataJSON, err := p.InvokeStartupOrgGetCampaign(ctx, campaignID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch campaign from StartupOrg: %v", err)
	}

	// Parse the campaign data
	var campaignData map[string]interface{}
	err = json.Unmarshal([]byte(campaignDataJSON), &campaignData)
	if err != nil {
		return "", fmt.Errorf("failed to parse campaign data: %v", err)
	}

	now := time.Now().Format(time.RFC3339)

	// Extract fields from campaign data with safe type assertions
	startupID, _ := campaignData["startupId"].(string)
	projectName, _ := campaignData["projectName"].(string)
	category, _ := campaignData["category"].(string)
	description, _ := campaignData["description"].(string)
	currency, _ := campaignData["currency"].(string)
	openDate, _ := campaignData["openDate"].(string)
	closeDate, _ := campaignData["close_date"].(string)
	validationHash, _ := campaignData["validationHash"].(string)

	// Handle numeric types
	goalAmount := float64(0)
	if ga, ok := campaignData["goalAmount"].(float64); ok {
		goalAmount = ga
	}

	durationDays := 0
	if dd, ok := campaignData["duration_days"].(float64); ok {
		durationDays = int(dd)
	}

	// Parse milestones if present
	var milestones []Milestone
	if ms, ok := campaignData["milestones"].([]interface{}); ok {
		for _, m := range ms {
			if mMap, ok := m.(map[string]interface{}); ok {
				milestone := Milestone{
					MilestoneID:   getString(mMap, "milestoneId"),
					Title:         getString(mMap, "title"),
					Description:   getString(mMap, "description"),
					TargetAmount:  getFloat(mMap, "targetAmount"),
					TargetDate:    getString(mMap, "targetDate"),
					Status:        "PENDING",
					FundsReleased: false,
				}
				milestones = append(milestones, milestone)
			}
		}
	}

	// Create published campaign using fetched data
	campaign := PublishedCampaign{
		CampaignID:         campaignID,
		StartupID:          startupID,
		ProjectName:        projectName,
		Category:           category,
		Description:        description,
		GoalAmount:         goalAmount,
		FundsRaisedAmount:  0,
		FundsRaisedPercent: 0,
		Currency:           currency,
		OpenDate:           openDate,
		CloseDate:          closeDate,
		DurationDays:       durationDays,
		ValidationScore:    decision.OverallScore,
		ValidationHash:     validationHash,
		ValidationVerified: true, // Already verified via RecordValidatorDecision
		Status:             "PUBLISHED",
		InvestorCount:      0,
		TotalConfirmed:     0,
		Milestones:         milestones,
		AgreementIDs:       []string{},
		PublishedAt:        now,
		UpdatedAt:          now,
	}

	campaignJSON, err := json.Marshal(campaign)
	if err != nil {
		return "", err
	}

	// Store on common-channel
	err = ctx.GetStub().PutState(campaignID, campaignJSON)
	if err != nil {
		return "", err
	}

	// Emit event
	eventPayload := map[string]interface{}{
		"campaignId":      campaignID,
		"startupId":       startupID,
		"projectName":     projectName,
		"validationScore": decision.OverallScore,
		"validationHash":  validationHash,
		"status":          "PUBLISHED",
		"channel":         "common-channel",
		"action":          "CAMPAIGN_PUBLISHED",
		"timestamp":       now,
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("CampaignPublished", eventJSON)

	response := map[string]interface{}{
		"message":         "Campaign published successfully to portal",
		"campaignId":      campaignID,
		"startupId":       startupID,
		"projectName":     projectName,
		"goalAmount":      goalAmount,
		"validationScore": decision.OverallScore,
		"status":          "PUBLISHED",
		"publishedAt":     now,
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// Helper functions for safe type extraction
func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func getFloat(m map[string]interface{}, key string) float64 {
	if v, ok := m[key].(float64); ok {
		return v
	}
	return 0
}

// VerifyAndPublish verifies campaign hash with ValidatorOrg and publishes
// Step 5.1: Platform verifies with Validator before publishing
// Channel: common-channel
// Endorsers: PlatformOrg (multi-party visibility)
func (p *PlatformContract) VerifyAndPublish(
	ctx contractapi.TransactionContextInterface,
	campaignID string,
	verifiedHash string, // Hash verified from Validator
	validatorConfirmed bool,
) (string, error) {
	// Retrieve campaign
	campaignJSON, err := ctx.GetStub().GetState(campaignID)
	if err != nil {
		return "", fmt.Errorf("failed to read campaign: %v", err)
	}
	if campaignJSON == nil {
		return "", fmt.Errorf("campaign %s does not exist", campaignID)
	}

	var campaign PublishedCampaign
	err = json.Unmarshal(campaignJSON, &campaign)
	if err != nil {
		return "", err
	}

	// Verify hash matches
	if campaign.ValidationHash != verifiedHash {
		return "", fmt.Errorf("validation hash mismatch. Campaign hash: %s, Verified hash: %s", campaign.ValidationHash, verifiedHash)
	}

	if !validatorConfirmed {
		return "", fmt.Errorf("validator did not confirm the campaign validity")
	}

	now := time.Now().Format(time.RFC3339)

	// Update to PUBLISHED
	campaign.ValidationVerified = true
	campaign.Status = "PUBLISHED"
	campaign.UpdatedAt = now

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
		"campaignId":         campaignID,
		"validationVerified": true,
		"status":             "PUBLISHED",
		"channel":            "common-channel",
		"action":             "CAMPAIGN_PUBLISHED",
		"timestamp":          now,
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("CampaignPublished", eventJSON)

	response := map[string]interface{}{
		"message":    "Campaign verified and published to investors",
		"campaignId": campaignID,
		"status":     "PUBLISHED",
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// WitnessAgreement records Platform as witness to startup-investor agreement
// Step 9: Platform and Validator witness the agreement for multi-party visibility
// Channel: common-channel
// Endorsers: PlatformOrg (multi-party visibility)
func (p *PlatformContract) WitnessAgreement(
	ctx contractapi.TransactionContextInterface,
	agreementID string,
	campaignID string,
	startupID string,
	investorID string,
	investmentAmount float64,
	currency string,
	terms string,
	milestonesJSON string,
) (string, error) {
	// Check if both parties have accepted
	existingJSON, err := ctx.GetStub().GetState(agreementID)
	if err != nil {
		return "", fmt.Errorf("failed to read agreement: %v", err)
	}

	var agreement Agreement
	now := time.Now().Format(time.RFC3339)

	if existingJSON != nil {
		// Update existing agreement
		err = json.Unmarshal(existingJSON, &agreement)
		if err != nil {
			return "", err
		}

		// Check both parties accepted
		if !agreement.StartupAccepted || !agreement.InvestorAccepted {
			return "", fmt.Errorf("both startup and investor must accept before Platform can witness. Startup: %v, Investor: %v", agreement.StartupAccepted, agreement.InvestorAccepted)
		}
	} else {
		// Parse milestones
		var milestones []Milestone
		if milestonesJSON != "" {
			if err := json.Unmarshal([]byte(milestonesJSON), &milestones); err != nil {
				return "", fmt.Errorf("failed to parse milestones: %v", err)
			}
		}

		// Create new agreement
		agreement = Agreement{
			AgreementID:      agreementID,
			CampaignID:       campaignID,
			StartupID:        startupID,
			InvestorID:       investorID,
			InvestmentAmount: investmentAmount,
			Currency:         currency,
			Milestones:       milestones,
			Terms:            terms,
			Status:           "PROPOSED",
			StartupAccepted:  false,
			InvestorAccepted: false,
			CreatedAt:        now,
		}
	}

	// Platform witnesses the agreement
	agreement.PlatformWitnessed = true
	agreement.WitnessedAt = now
	agreement.Status = "ACTIVE"
	agreement.AcceptedAt = now

	agreementJSON, err := json.Marshal(agreement)
	if err != nil {
		return "", err
	}

	// Store agreement
	err = ctx.GetStub().PutState(agreementID, agreementJSON)
	if err != nil {
		return "", err
	}

	// Update campaign with agreement
	campaignJSON, _ := ctx.GetStub().GetState(campaignID)
	if campaignJSON != nil {
		var campaign PublishedCampaign
		json.Unmarshal(campaignJSON, &campaign)
		campaign.AgreementIDs = append(campaign.AgreementIDs, agreementID)
		campaign.UpdatedAt = now
		updatedCampaignJSON, _ := json.Marshal(campaign)
		ctx.GetStub().PutState(campaignID, updatedCampaignJSON)
	}

	// Create escrow for the funds
	escrow := FundEscrow{
		EscrowID:       fmt.Sprintf("ESCROW_%s", agreementID),
		AgreementID:    agreementID,
		CampaignID:     campaignID,
		InvestorID:     investorID,
		StartupID:      startupID,
		TotalAmount:    investmentAmount,
		ReleasedAmount: 0,
		HeldAmount:     investmentAmount,
		Currency:       currency,
		Status:         "ACTIVE",
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	escrowJSON, _ := json.Marshal(escrow)
	ctx.GetStub().PutState(escrow.EscrowID, escrowJSON)

	// Emit event
	eventPayload := map[string]interface{}{
		"agreementId":      agreementID,
		"campaignId":       campaignID,
		"startupId":        startupID,
		"investorId":       investorID,
		"investmentAmount": investmentAmount,
		"escrowId":         escrow.EscrowID,
		"status":           "ACTIVE",
		"channel":          "common-channel",
		"action":           "AGREEMENT_WITNESSED",
		"timestamp":        now,
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("AgreementWitnessed", eventJSON)

	response := map[string]interface{}{
		"message":          "Agreement witnessed by Platform. Funds held in escrow.",
		"agreementId":      agreementID,
		"escrowId":         escrow.EscrowID,
		"investmentAmount": investmentAmount,
		"status":           "ACTIVE",
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// TriggerFundRelease releases funds to startup based on milestone completion
// Step 13: Platform releases funds from escrow when milestone is verified
// Channel: common-channel
// Endorsers: PlatformOrg (multi-party visibility)
func (p *PlatformContract) TriggerFundRelease(
	ctx contractapi.TransactionContextInterface,
	releaseID string,
	escrowID string,
	agreementID string,
	campaignID string,
	milestoneID string,
	startupID string,
	amount float64,
	currency string,
	triggerReason string,
) (string, error) {
	// Retrieve escrow
	escrowJSON, err := ctx.GetStub().GetState(escrowID)
	if err != nil {
		return "", fmt.Errorf("failed to read escrow: %v", err)
	}
	if escrowJSON == nil {
		return "", fmt.Errorf("escrow %s does not exist", escrowID)
	}

	var escrow FundEscrow
	err = json.Unmarshal(escrowJSON, &escrow)
	if err != nil {
		return "", err
	}

	// Check sufficient funds in escrow
	if amount > escrow.HeldAmount {
		return "", fmt.Errorf("insufficient funds in escrow. Available: %f, Requested: %f", escrow.HeldAmount, amount)
	}

	now := time.Now().Format(time.RFC3339)

	// Create fund release record
	release := FundRelease{
		ReleaseID:     releaseID,
		EscrowID:      escrowID,
		AgreementID:   agreementID,
		CampaignID:    campaignID,
		MilestoneID:   milestoneID,
		StartupID:     startupID,
		Amount:        amount,
		Currency:      currency,
		Status:        "RELEASED",
		TriggerReason: triggerReason,
		ReleasedAt:    now,
	}

	releaseJSON, err := json.Marshal(release)
	if err != nil {
		return "", err
	}

	// Store release record
	err = ctx.GetStub().PutState(releaseID, releaseJSON)
	if err != nil {
		return "", err
	}

	// Update escrow
	escrow.ReleasedAmount += amount
	escrow.HeldAmount -= amount
	escrow.UpdatedAt = now
	if escrow.HeldAmount <= 0 {
		escrow.Status = "FULLY_RELEASED"
	} else {
		escrow.Status = "PARTIALLY_RELEASED"
	}

	updatedEscrowJSON, _ := json.Marshal(escrow)
	ctx.GetStub().PutState(escrowID, updatedEscrowJSON)

	// Update campaign
	campaignJSON, err := ctx.GetStub().GetState(campaignID)
	if err == nil && campaignJSON != nil {
		var campaign PublishedCampaign
		json.Unmarshal(campaignJSON, &campaign)
		
		// Update milestone status
		for i, m := range campaign.Milestones {
			if m.MilestoneID == milestoneID {
				campaign.Milestones[i].FundsReleased = true
				campaign.Milestones[i].ReleasedAt = now
				campaign.Milestones[i].Status = "VERIFIED"
				break
			}
		}
		
		campaign.FundsRaisedAmount += amount
		if campaign.GoalAmount > 0 {
			campaign.FundsRaisedPercent = (campaign.FundsRaisedAmount / campaign.GoalAmount) * 100
		}
		campaign.Status = "FUNDED"
		campaign.UpdatedAt = now
		updatedCampaignJSON, _ := json.Marshal(campaign)
		ctx.GetStub().PutState(campaignID, updatedCampaignJSON)
	}

	// Emit event
	eventPayload := map[string]interface{}{
		"releaseId":     releaseID,
		"escrowId":      escrowID,
		"agreementId":   agreementID,
		"campaignId":    campaignID,
		"milestoneId":   milestoneID,
		"startupId":     startupID,
		"amount":        amount,
		"triggerReason": triggerReason,
		"channel":       "common-channel",
		"action":        "FUNDS_RELEASED",
		"timestamp":     now,
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("FundsReleased", eventJSON)

	response := map[string]interface{}{
		"message":       "Funds released to startup from escrow",
		"releaseId":     releaseID,
		"milestoneId":   milestoneID,
		"amount":        amount,
		"escrowBalance": escrow.HeldAmount,
		"status":        "RELEASED",
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// CloseCampaign closes a campaign after deadline or completion
// Step 14: Campaign closure for multi-party visibility
// Channel: common-channel
// Endorsers: PlatformOrg (multi-party visibility)
func (p *PlatformContract) CloseCampaign(
	ctx contractapi.TransactionContextInterface,
	closureID string,
	campaignID string,
	finalStatus string,
	finalAmount float64,
	finalInvestorCount int,
	closureReason string,
) (string, error) {
	// Create closure record
	closure := CampaignClosure{
		ClosureID:          closureID,
		CampaignID:         campaignID,
		FinalStatus:        finalStatus,
		FinalAmount:        finalAmount,
		FinalInvestorCount: finalInvestorCount,
		ClosureReason:      closureReason,
		ClosedAt:           time.Now().Format(time.RFC3339),
	}

	closureJSON, err := json.Marshal(closure)
	if err != nil {
		return "", err
	}

	// Store closure record
	err = ctx.GetStub().PutState(closureID, closureJSON)
	if err != nil {
		return "", err
	}

	// Update campaign status
	campaignJSON, err := ctx.GetStub().GetState(campaignID)
	if err == nil && campaignJSON != nil {
		var campaign PublishedCampaign
		json.Unmarshal(campaignJSON, &campaign)
		campaign.Status = "CLOSED"
		campaign.FundsRaisedAmount = finalAmount
		campaign.InvestorCount = finalInvestorCount
		if campaign.GoalAmount > 0 {
			campaign.FundsRaisedPercent = (finalAmount / campaign.GoalAmount) * 100
		}
		campaign.UpdatedAt = time.Now().Format(time.RFC3339)
		updatedCampaignJSON, _ := json.Marshal(campaign)
		ctx.GetStub().PutState(campaignID, updatedCampaignJSON)
	}

	// Emit event
	eventPayload := map[string]interface{}{
		"closureId":     closureID,
		"campaignId":    campaignID,
		"finalStatus":   finalStatus,
		"closureReason": closureReason,
		"channel":       "common-channel",
		"action":        "CAMPAIGN_CLOSED",
		"timestamp":     closure.ClosedAt,
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("CampaignClosed", eventJSON)

	response := map[string]interface{}{
		"message":     "Campaign closed successfully",
		"closureId":   closureID,
		"campaignId":  campaignID,
		"finalStatus": finalStatus,
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// ============================================================================
// INVESTOR-PLATFORM-CHANNEL FUNCTIONS
// Endorsed by: InvestorOrg, PlatformOrg
// ============================================================================

// RecordInvestorConfirmation records investment confirmation from InvestorOrg
// Channel: investor-platform-channel
// Endorsers: InvestorOrg, PlatformOrg
func (p *PlatformContract) RecordInvestorConfirmation(
	ctx contractapi.TransactionContextInterface,
	recordID string,
	confirmationID string,
	campaignID string,
	investorID string,
	amount float64,
	currency string,
) (string, error) {
	// Create confirmation record
	record := InvestorConfirmationRecord{
		RecordID:       recordID,
		ConfirmationID: confirmationID,
		CampaignID:     campaignID,
		InvestorID:     investorID,
		Amount:         amount,
		Currency:       currency,
		RecordedAt:     time.Now().Format(time.RFC3339),
	}

	recordJSON, err := json.Marshal(record)
	if err != nil {
		return "", err
	}

	// Store on investor-platform-channel
	err = ctx.GetStub().PutState(recordID, recordJSON)
	if err != nil {
		return "", err
	}

	// Update campaign funding totals
	campaignJSON, err := ctx.GetStub().GetState(campaignID)
	if err == nil && campaignJSON != nil {
		var campaign PublishedCampaign
		json.Unmarshal(campaignJSON, &campaign)
		campaign.TotalConfirmed += amount
		campaign.FundsRaisedAmount += amount
		campaign.InvestorCount++
		if campaign.GoalAmount > 0 {
			campaign.FundsRaisedPercent = (campaign.FundsRaisedAmount / campaign.GoalAmount) * 100
		}
		campaign.UpdatedAt = time.Now().Format(time.RFC3339)
		updatedCampaignJSON, _ := json.Marshal(campaign)
		ctx.GetStub().PutState(campaignID, updatedCampaignJSON)
	}

	// Emit event
	eventPayload := map[string]interface{}{
		"recordId":       recordID,
		"confirmationId": confirmationID,
		"campaignId":     campaignID,
		"investorId":     investorID,
		"channel":        "investor-platform-channel",
		"action":         "INVESTOR_CONFIRMATION_RECORDED",
		"timestamp":      record.RecordedAt,
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("InvestorConfirmationRecorded", eventJSON)

	response := map[string]interface{}{
		"message":    "Investor confirmation recorded",
		"recordId":   recordID,
		"campaignId": campaignID,
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// ============================================================================
// VALIDATOR-PLATFORM-CHANNEL FUNCTIONS
// Endorsed by: ValidatorOrg, PlatformOrg
// ============================================================================

// RecordValidatorDecision records final validation decision from ValidatorOrg
// Channel: validator-platform-channel
// Endorsers: ValidatorOrg, PlatformOrg
func (p *PlatformContract) RecordValidatorDecision(
	ctx contractapi.TransactionContextInterface,
	recordID string,
	campaignID string,
	validationID string,
	approved bool,
	overallScore float64,
	reportHash string,
) (string, error) {
	// Create decision record
	record := ValidatorDecisionRecord{
		RecordID:     recordID,
		CampaignID:   campaignID,
		ValidationID: validationID,
		Approved:     approved,
		OverallScore: overallScore,
		ReportHash:   reportHash,
		RecordedAt:   time.Now().Format(time.RFC3339),
	}

	recordJSON, err := json.Marshal(record)
	if err != nil {
		return "", err
	}

	// Store on validator-platform-channel
	err = ctx.GetStub().PutState(recordID, recordJSON)
	if err != nil {
		return "", err
	}

	// Store by campaign for lookup
	decisionKey := fmt.Sprintf("DECISION_%s", campaignID)
	err = ctx.GetStub().PutState(decisionKey, recordJSON)
	if err != nil {
		return "", err
	}

	// Emit event
	eventPayload := map[string]interface{}{
		"recordId":     recordID,
		"campaignId":   campaignID,
		"validationId": validationID,
		"approved":     approved,
		"overallScore": overallScore,
		"channel":      "validator-platform-channel",
		"action":       "VALIDATOR_DECISION_RECORDED",
		"timestamp":    record.RecordedAt,
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("ValidatorDecisionRecorded", eventJSON)

	response := map[string]interface{}{
		"message":      "Validator decision recorded",
		"recordId":     recordID,
		"campaignId":   campaignID,
		"approved":     approved,
		"overallScore": overallScore,
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// ============================================================================
// COMMON-CHANNEL FUNCTIONS
// Read by: All Orgs, Write by: PlatformOrg (aggregated metrics only - privacy preserving)
// ============================================================================

// PublishGlobalMetrics publishes aggregated system-wide metrics to common-channel
// Channel: common-channel
// Purpose: Privacy-preserving global metrics (never raw data)
func (p *PlatformContract) PublishGlobalMetrics(
	ctx contractapi.TransactionContextInterface,
	metricsID string,
	totalCampaigns int,
	activeCampaigns int,
	successfulCampaigns int,
	totalInvestorCount int,
) (string, error) {
	// Generate metrics hash (no sensitive data)
	metricsData := map[string]interface{}{
		"metricsId":           metricsID,
		"totalCampaigns":      totalCampaigns,
		"activeCampaigns":     activeCampaigns,
		"successfulCampaigns": successfulCampaigns,
		"totalInvestorCount":  totalInvestorCount,
		"timestamp":           time.Now().Format(time.RFC3339),
	}
	metricsDataJSON, _ := json.Marshal(metricsData)
	metricsHash := generateHash(string(metricsDataJSON))

	// Create global metrics
	metrics := GlobalMetrics{
		MetricsID:           metricsID,
		TotalCampaigns:      totalCampaigns,
		ActiveCampaigns:     activeCampaigns,
		SuccessfulCampaigns: successfulCampaigns,
		TotalInvestorCount:  totalInvestorCount,
		MetricsHash:         metricsHash,
		PublishedAt:         time.Now().Format(time.RFC3339),
	}

	metricsJSON, err := json.Marshal(metrics)
	if err != nil {
		return "", err
	}

	// Store on common-channel
	err = ctx.GetStub().PutState(metricsID, metricsJSON)
	if err != nil {
		return "", err
	}

	// Also store as latest metrics
	err = ctx.GetStub().PutState("COMMON_GLOBAL_METRICS_LATEST", metricsJSON)
	if err != nil {
		return "", err
	}

	// Emit event
	eventPayload := map[string]interface{}{
		"metricsId":           metricsID,
		"totalCampaigns":      totalCampaigns,
		"activeCampaigns":     activeCampaigns,
		"successfulCampaigns": successfulCampaigns,
		"totalInvestorCount":  totalInvestorCount,
		"metricsHash":         metricsHash,
		"channel":             "common-channel",
		"action":              "GLOBAL_METRICS_PUBLISHED",
		"timestamp":           metrics.PublishedAt,
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("GlobalMetricsPublished", eventJSON)

	response := map[string]interface{}{
		"message":             "Global metrics published to common channel",
		"metricsId":           metricsID,
		"totalCampaigns":      totalCampaigns,
		"activeCampaigns":     activeCampaigns,
		"successfulCampaigns": successfulCampaigns,
		"metricsHash":         metricsHash,
		"channel":             "common-channel",
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// ============================================================================
// QUERY FUNCTIONS
// ============================================================================

// GetPublishedCampaign retrieves published campaign by ID
func (p *PlatformContract) GetPublishedCampaign(ctx contractapi.TransactionContextInterface, campaignID string) (*PublishedCampaign, error) {
	campaignJSON, err := ctx.GetStub().GetState(campaignID)
	if err != nil {
		return nil, fmt.Errorf("failed to read campaign: %v", err)
	}
	if campaignJSON == nil {
		return nil, fmt.Errorf("campaign %s does not exist", campaignID)
	}

	var campaign PublishedCampaign
	err = json.Unmarshal(campaignJSON, &campaign)
	if err != nil {
		return nil, err
	}

	return &campaign, nil
}

// GetActiveCampaigns returns all active published campaigns
func (p *PlatformContract) GetActiveCampaigns(ctx contractapi.TransactionContextInterface) (string, error) {
	queryString := `{"selector":{"status":"PUBLISHED"}}`

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

		var campaign PublishedCampaign
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

// GetValidatorDecision retrieves validator decision for a campaign
func (p *PlatformContract) GetValidatorDecision(ctx contractapi.TransactionContextInterface, campaignID string) (*ValidatorDecisionRecord, error) {
	decisionKey := fmt.Sprintf("DECISION_%s", campaignID)
	decisionJSON, err := ctx.GetStub().GetState(decisionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read decision: %v", err)
	}
	if decisionJSON == nil {
		return nil, fmt.Errorf("validator decision for campaign %s does not exist", campaignID)
	}

	var decision ValidatorDecisionRecord
	err = json.Unmarshal(decisionJSON, &decision)
	if err != nil {
		return nil, err
	}

	return &decision, nil
}

// GetLatestGlobalMetrics retrieves the latest global metrics
func (p *PlatformContract) GetLatestGlobalMetrics(ctx contractapi.TransactionContextInterface) (*GlobalMetrics, error) {
	metricsJSON, err := ctx.GetStub().GetState("COMMON_GLOBAL_METRICS_LATEST")
	if err != nil {
		return nil, fmt.Errorf("failed to read metrics: %v", err)
	}
	if metricsJSON == nil {
		return nil, fmt.Errorf("no global metrics published yet")
	}

	var metrics GlobalMetrics
	err = json.Unmarshal(metricsJSON, &metrics)
	if err != nil {
		return nil, err
	}

	return &metrics, nil
}

// GetAgreement retrieves a witnessed agreement by ID
// Channel: common-channel
func (p *PlatformContract) GetAgreement(ctx contractapi.TransactionContextInterface, agreementID string) (*Agreement, error) {
	agreementJSON, err := ctx.GetStub().GetState(agreementID)
	if err != nil {
		return nil, fmt.Errorf("failed to read agreement: %v", err)
	}
	if agreementJSON == nil {
		return nil, fmt.Errorf("agreement %s does not exist", agreementID)
	}

	var agreement Agreement
	err = json.Unmarshal(agreementJSON, &agreement)
	if err != nil {
		return nil, err
	}

	return &agreement, nil
}

// GetAgreementsByCampaign retrieves all agreements for a campaign
// Channel: common-channel
func (p *PlatformContract) GetAgreementsByCampaign(ctx contractapi.TransactionContextInterface, campaignID string) (string, error) {
	queryString := fmt.Sprintf(`{"selector":{"campaignId":"%s","agreementId":{"$exists":true}}}`, campaignID)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return "", err
	}
	defer resultsIterator.Close()

	var agreements []map[string]interface{}
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return "", err
		}

		var agreement Agreement
		err = json.Unmarshal(queryResponse.Value, &agreement)
		if err != nil {
			continue
		}

		agreementMap := map[string]interface{}{
			"Key":    queryResponse.Key,
			"Record": agreement,
		}
		agreements = append(agreements, agreementMap)
	}

	agreementsJSON, err := json.Marshal(agreements)
	if err != nil {
		return "", err
	}

	return string(agreementsJSON), nil
}

// GetFundRelease retrieves a fund release record by ID
// Channel: common-channel
func (p *PlatformContract) GetFundRelease(ctx contractapi.TransactionContextInterface, releaseID string) (*FundRelease, error) {
	releaseJSON, err := ctx.GetStub().GetState(releaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to read fund release: %v", err)
	}
	if releaseJSON == nil {
		return nil, fmt.Errorf("fund release %s does not exist", releaseID)
	}

	var release FundRelease
	err = json.Unmarshal(releaseJSON, &release)
	if err != nil {
		return nil, err
	}

	return &release, nil
}

// GetEscrow retrieves an escrow record by ID
// Channel: common-channel
func (p *PlatformContract) GetEscrow(ctx contractapi.TransactionContextInterface, escrowID string) (*FundEscrow, error) {
	escrowJSON, err := ctx.GetStub().GetState(escrowID)
	if err != nil {
		return nil, fmt.Errorf("failed to read escrow: %v", err)
	}
	if escrowJSON == nil {
		return nil, fmt.Errorf("escrow %s does not exist", escrowID)
	}

	var escrow FundEscrow
	err = json.Unmarshal(escrowJSON, &escrow)
	if err != nil {
		return nil, err
	}

	return &escrow, nil
}

// ============================================================================
// CROSS-CHANNEL INVOCATION HELPER FUNCTIONS
// ============================================================================

// InvokeStartupOrgGetCampaign reads full campaign data from StartupOrg
// Cross-channel READ from startup-platform-channel
func (p *PlatformContract) InvokeStartupOrgGetCampaign(
	ctx contractapi.TransactionContextInterface,
	campaignID string,
) (string, error) {
	args := [][]byte{
		[]byte("GetCampaign"),
		[]byte(campaignID),
	}

	response := ctx.GetStub().InvokeChaincode(
		"startuporg",
		args,
		"startup-platform-channel",
	)

	if response.Status != 200 {
		return "", fmt.Errorf("cross-channel query to StartupOrg failed: %s", response.Message)
	}

	return string(response.Payload), nil
}

// InvokeValidatorOrgGetValidation reads validation data from ValidatorOrg
// Cross-channel READ from validator-platform-channel
func (p *PlatformContract) InvokeValidatorOrgGetValidation(
	ctx contractapi.TransactionContextInterface,
	validationID string,
) (string, error) {
	args := [][]byte{
		[]byte("GetValidation"),
		[]byte(validationID),
	}

	response := ctx.GetStub().InvokeChaincode(
		"validatororg",
		args,
		"validator-platform-channel",
	)

	if response.Status != 200 {
		return "", fmt.Errorf("cross-channel query to ValidatorOrg failed: %s", response.Message)
	}

	return string(response.Payload), nil
}

// InvokeStartupOrgNotifyFundRelease notifies StartupOrg about fund release
// Cross-channel call to common-channel
func (p *PlatformContract) InvokeStartupOrgNotifyFundRelease(
	ctx contractapi.TransactionContextInterface,
	campaignID string,
	amount string,
	status string,
) (string, error) {
	args := [][]byte{
		[]byte("ReceiveFunding"),
		[]byte(campaignID),
		[]byte(amount),
		[]byte(status),
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
		"campaignId":     campaignID,
		"targetChannel":  "common-channel",
		"targetContract": "startuporg",
		"action":         "CROSS_CHANNEL_FUND_RELEASE_NOTIFY",
		"timestamp":      time.Now().Format(time.RFC3339),
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("CrossChannelInvoke", eventJSON)

	return string(response.Payload), nil
}

// InvokeCommonChannelPublish publishes to common-channel (all orgs can read)
func (p *PlatformContract) InvokeCommonChannelPublish(
	ctx contractapi.TransactionContextInterface,
	metricsID string,
	metricsHash string,
) (string, error) {
	// Common channel is directly accessible, but this shows cross-channel pattern
	args := [][]byte{
		[]byte("PublishGlobalMetrics"),
		[]byte(metricsID),
		[]byte(metricsHash),
	}

	response := ctx.GetStub().InvokeChaincode(
		"platformorg", // same contract but different channel
		args,
		"common-channel",
	)

	if response.Status != 200 {
		return "", fmt.Errorf("cross-channel invoke to common-channel failed: %s", response.Message)
	}

	return string(response.Payload), nil
}

// generateHash generates SHA256 hash
func generateHash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func main() {
	platformChaincode, err := contractapi.NewChaincode(&PlatformContract{})
	if err != nil {
		fmt.Printf("Error creating PlatformOrg chaincode: %v\n", err)
		return
	}

	if err := platformChaincode.Start(); err != nil {
		fmt.Printf("Error starting PlatformOrg chaincode: %v\n", err)
	}
}
