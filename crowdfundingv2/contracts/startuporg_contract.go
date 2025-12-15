package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// StartupContract provides functions for StartupOrg operations using PDC
type StartupContract struct {
	contractapi.Contract
}

// ============================================================================
// DATA STRUCTURES
// ============================================================================

// Campaign represents a startup crowdfunding campaign with all required fields (22 parameters)
type Campaign struct {
	CampaignID          string   `json:"campaignId"`
	StartupID           string   `json:"startupId"`

	// Core Campaign Fields (22 Parameters)
	Category            string   `json:"category"`
	Deadline            string   `json:"deadline"`
	Currency            string   `json:"currency"`
	HasRaised           bool     `json:"has_raised"`
	HasGovGrants        bool     `json:"has_gov_grants"`
	IncorpDate          string   `json:"incorp_date"`
	ProjectStage        string   `json:"project_stage"`
	Sector              string   `json:"sector"`
	Tags                []string `json:"tags"`
	TeamAvailable       bool     `json:"team_available"`
	InvestorCommitted   bool     `json:"investor_committed"`
	Duration            int      `json:"duration"`
	FundingDay          int      `json:"funding_day"`
	FundingMonth        int      `json:"funding_month"`
	FundingYear         int      `json:"funding_year"`
	GoalAmount          float64  `json:"goal_amount"`
	InvestmentRange     string   `json:"investment_range"`
	ProjectName         string   `json:"project_name"`
	Description         string   `json:"description"`
	Documents           []string `json:"documents"`

	// Calculated/Status Fields
	OpenDate            string   `json:"open_date"`
	CloseDate           string   `json:"close_date"`
	FundsRaisedAmount   float64  `json:"funds_raised_amount"`
	FundsRaisedPercent  float64  `json:"funds_raised_percent"`
	Status              string   `json:"status"`
	ValidationStatus    string   `json:"validationStatus"`
	ValidationScore     float64  `json:"validationScore"`
	ValidationHash      string   `json:"validationHash"`
	InvestorCount       int      `json:"investorCount"`
	PlatformStatus      string   `json:"platformStatus"`

	// Timestamps
	CreatedAt           string   `json:"createdAt"`
	UpdatedAt           string   `json:"updatedAt"`
	ApprovedAt          string   `json:"approvedAt"`
	PublishedAt         string   `json:"publishedAt"`
}

// CampaignPrivateDetails - stored in StartupPrivateCollection
type CampaignPrivateDetails struct {
	CampaignID        string               `json:"campaignId"`
	DocumentHistory   []DocumentSubmission `json:"documentHistory"`
	CurrentDocuments  []string             `json:"currentDocuments"`
	UpdateHistory     []CampaignUpdate     `json:"updateHistory"`
	InternalNotes     string               `json:"internalNotes"`
	FinancialDetails  string               `json:"financialDetails"`
}

// CampaignPublicInfo - stored on public ledger (world state)
type CampaignPublicInfo struct {
	CampaignID     string  `json:"campaignId"`
	StartupID      string  `json:"startupId"`
	ProjectName    string  `json:"projectName"`
	Category       string  `json:"category"`
	GoalAmount     float64 `json:"goalAmount"`
	Currency       string  `json:"currency"`
	Status         string  `json:"status"`
	PublishedAt    string  `json:"publishedAt"`
}

// DocumentSubmission tracks each document submission attempt
type DocumentSubmission struct {
	SubmissionID    string   `json:"submissionId"`
	Documents       []string `json:"documents"`
	SubmittedAt     string   `json:"submittedAt"`
	SubmissionNotes string   `json:"submissionNotes"`
	ResponseStatus  string   `json:"responseStatus"`
	ResponseNotes   string   `json:"responseNotes"`
	ResponseAt      string   `json:"responseAt"`
}

// CampaignUpdate tracks each update made to campaign fields
type CampaignUpdate struct {
	UpdateID     string `json:"updateId"`
	FieldName    string `json:"fieldName"`
	OldValue     string `json:"oldValue"`
	NewValue     string `json:"newValue"`
	UpdateReason string `json:"updateReason"`
	UpdatedAt    string `json:"updatedAt"`
	UpdatedBy    string `json:"updatedBy"`
}

// Milestone represents a funding milestone
// MilestoneReport for progress reporting
// Stored in StartupValidatorCollection when submitted for verification
type MilestoneReport struct {
	ReportID       string   `json:"reportId"`
	CampaignID     string   `json:"campaignId"`
	MilestoneID    string   `json:"milestoneId"`
	AgreementID    string   `json:"agreementId"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	Evidence       []string `json:"evidence"`
	Status         string   `json:"status"`
	SubmittedAt    string   `json:"submittedAt"`
	ReviewedAt     string   `json:"reviewedAt"`
	ReviewComments string   `json:"reviewComments"`
}

// DisputeSubmission represents a dispute submitted by startup
type DisputeSubmission struct {
	SubmissionID    string   `json:"submissionId"`
	DisputeID       string   `json:"disputeId"`
	StartupID       string   `json:"startupId"`
	DisputeType     string   `json:"disputeType"`
	TargetID        string   `json:"targetId"`
	TargetType      string   `json:"targetType"`
	CampaignID      string   `json:"campaignId"`
	AgreementID     string   `json:"agreementId"`
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	ClaimedAmount   float64  `json:"claimedAmount"`
	EvidenceHashes  []string `json:"evidenceHashes"`
	Status          string   `json:"status"`
	CreatedAt       string   `json:"createdAt"`
}

// FeePaymentRecord represents fee payment record
type FeePaymentRecord struct {
	RecordID        string  `json:"recordId"`
	StartupID       string  `json:"startupId"`
	CampaignID      string  `json:"campaignId"`
	FeeType         string  `json:"feeType"`
	Amount          float64 `json:"amount"`
	TransactionHash string  `json:"transactionHash"`
	PaidAt          string  `json:"paidAt"`
}

// ============================================================================
// INIT
// ============================================================================

func (s *StartupContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	fmt.Println("StartupOrg contract initialized with PDC support")
	return nil
}

// ============================================================================
// CAMPAIGN MANAGEMENT - Using PDC
// ============================================================================

// CreateCampaign creates a new campaign with 22 parameters matching API format
func (s *StartupContract) CreateCampaign(
	ctx contractapi.TransactionContextInterface,
	campaignID string,
	startupID string,
	category string,
	deadline string,
	currency string,
	hasRaised string,
	hasGovGrants string,
	incorpDate string,
	projectStage string,
	sector string,
	tagsJSON string,
	teamAvailable string,
	investorCommitted string,
	duration string,
	fundingDay string,
	fundingMonth string,
	fundingYear string,
	goalAmount string,
	investmentRange string,
	projectName string,
	description string,
	documentsJSON string,
) error {

	// Convert string parameters to appropriate types
	hasRaisedBool, _ := strconv.ParseBool(hasRaised)
	hasGovGrantsBool, _ := strconv.ParseBool(hasGovGrants)
	teamAvailableBool, _ := strconv.ParseBool(teamAvailable)
	investorCommittedBool, _ := strconv.ParseBool(investorCommitted)
	durationInt, _ := strconv.Atoi(duration)
	fundingDayInt, _ := strconv.Atoi(fundingDay)
	fundingMonthInt, _ := strconv.Atoi(fundingMonth)
	fundingYearInt, _ := strconv.Atoi(fundingYear)
	goalAmountFloat, _ := strconv.ParseFloat(goalAmount, 64)
	
	// Parse JSON arrays
	var tags []string
	if tagsJSON != "" {
		err := json.Unmarshal([]byte(tagsJSON), &tags)
		if err != nil {
			return fmt.Errorf("failed to parse tags: %v", err)
		}
	}

	var documents []string
	if documentsJSON != "" {
		err := json.Unmarshal([]byte(documentsJSON), &documents)
		if err != nil {
			return fmt.Errorf("failed to parse documents: %v", err)
		}
	}

	timestamp := time.Now().Format(time.RFC3339)

	// Calculate open and close dates based on funding date and duration
	openDate := fmt.Sprintf("%04d-%02d-%02d", fundingYearInt, fundingMonthInt, fundingDayInt)
	// Calculate close date by adding duration days
	closeDate := deadline

	// Create full campaign object
	campaign := Campaign{
		CampaignID:        campaignID,
		StartupID:         startupID,
		Category:          category,
		Deadline:          deadline,
		Currency:          currency,
		HasRaised:         hasRaisedBool,
		HasGovGrants:      hasGovGrantsBool,
		IncorpDate:        incorpDate,
		ProjectStage:      projectStage,
		Sector:            sector,
		Tags:              tags,
		TeamAvailable:     teamAvailableBool,
		InvestorCommitted: investorCommittedBool,
		Duration:          durationInt,
		FundingDay:        fundingDayInt,
		FundingMonth:      fundingMonthInt,
		FundingYear:       fundingYearInt,
		GoalAmount:        goalAmountFloat,
		InvestmentRange:   investmentRange,
		ProjectName:       projectName,
		Description:       description,
		Documents:         documents,
		OpenDate:          openDate,
		CloseDate:         closeDate,
		FundsRaisedAmount: 0,
		FundsRaisedPercent: 0,
		Status:            "DRAFT",
		ValidationStatus:  "NOT_SUBMITTED",
		ValidationScore:   0,
		ValidationHash:    "",
		InvestorCount:     0,
		PlatformStatus:    "NOT_PUBLISHED",
		CreatedAt:         timestamp,
		UpdatedAt:         timestamp,
		ApprovedAt:        "",
		PublishedAt:       "",
	}

	// Store full campaign in StartupPrivateCollection
	campaignJSON, err := json.Marshal(campaign)
	if err != nil {
		return fmt.Errorf("failed to marshal campaign: %v", err)
	}

	err = ctx.GetStub().PutPrivateData(StartupPrivateCollection, "CAMPAIGN_"+campaignID, campaignJSON)
	if err != nil {
		return fmt.Errorf("failed to store campaign in private collection: %v", err)
	}

	// Store public campaign info on world state for visibility
	publicInfo := CampaignPublicInfo{
		CampaignID:  campaignID,
		StartupID:   startupID,
		ProjectName: projectName,
		Category:    category,
		GoalAmount:  goalAmountFloat,
		Currency:    currency,
		Status:      "DRAFT",
		PublishedAt: "",
	}

	publicJSON, err := json.Marshal(publicInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal public info: %v", err)
	}

	err = ctx.GetStub().PutState("CAMPAIGN_PUBLIC_"+campaignID, publicJSON)
	if err != nil {
		return fmt.Errorf("failed to store public campaign info: %v", err)
	}

	// Initialize private details
	privateDetails := CampaignPrivateDetails{
		CampaignID:       campaignID,
		DocumentHistory:  []DocumentSubmission{},
		CurrentDocuments: documents,
		UpdateHistory:    []CampaignUpdate{},
		InternalNotes:    "",
		FinancialDetails: "",
	}

	privateJSON, err := json.Marshal(privateDetails)
	if err != nil {
		return fmt.Errorf("failed to marshal private details: %v", err)
	}

	err = ctx.GetStub().PutPrivateData(StartupPrivateCollection, "CAMPAIGN_PRIVATE_"+campaignID, privateJSON)
	if err != nil {
		return fmt.Errorf("failed to store private details: %v", err)
	}

	return nil
}

// UpdateCampaign updates campaign fields and records the update history
func (s *StartupContract) UpdateCampaign(
	ctx contractapi.TransactionContextInterface,
	campaignID string,
	fieldName string,
	newValue string,
	updateReason string,
	updatedBy string,
) error {

	// Get campaign from private collection
	campaignJSON, err := ctx.GetStub().GetPrivateData(StartupPrivateCollection, "CAMPAIGN_"+campaignID)
	if err != nil {
		return fmt.Errorf("failed to read campaign: %v", err)
	}
	if campaignJSON == nil {
		return fmt.Errorf("campaign %s does not exist", campaignID)
	}

	var campaign Campaign
	err = json.Unmarshal(campaignJSON, &campaign)
	if err != nil {
		return fmt.Errorf("failed to unmarshal campaign: %v", err)
	}

	// Check if campaign can be updated
	if campaign.Status != "DRAFT" && campaign.ValidationStatus != "ON_HOLD" {
		return fmt.Errorf("campaign cannot be updated in current status")
	}

	// Get private details to store update history
	detailsJSON, err := ctx.GetStub().GetPrivateData(StartupPrivateCollection, "CAMPAIGN_PRIVATE_"+campaignID)
	if err != nil {
		return fmt.Errorf("failed to read campaign private details: %v", err)
	}

	var privateDetails CampaignPrivateDetails
	err = json.Unmarshal(detailsJSON, &privateDetails)
	if err != nil {
		return fmt.Errorf("failed to unmarshal private details: %v", err)
	}

	timestamp := time.Now().Format(time.RFC3339)
	updateID := fmt.Sprintf("UPDATE_%s_%d", campaignID, len(privateDetails.UpdateHistory)+1)

	// Record the update (simplified - in real implementation you'd reflect the actual field change)
	updateEntry := CampaignUpdate{
		UpdateID:     updateID,
		FieldName:    fieldName,
		OldValue:     "", // Would retrieve actual old value
		NewValue:     newValue,
		UpdateReason: updateReason,
		UpdatedAt:    timestamp,
		UpdatedBy:    updatedBy,
	}

	privateDetails.UpdateHistory = append(privateDetails.UpdateHistory, updateEntry)
	campaign.UpdatedAt = timestamp

	// Save updated campaign
	campaignJSON, err = json.Marshal(campaign)
	if err != nil {
		return fmt.Errorf("failed to marshal campaign: %v", err)
	}

	err = ctx.GetStub().PutPrivateData(StartupPrivateCollection, "CAMPAIGN_"+campaignID, campaignJSON)
	if err != nil {
		return fmt.Errorf("failed to update campaign: %v", err)
	}

	// Save updated private details
	detailsJSON, err = json.Marshal(privateDetails)
	if err != nil {
		return fmt.Errorf("failed to marshal private details: %v", err)
	}

	err = ctx.GetStub().PutPrivateData(StartupPrivateCollection, "CAMPAIGN_PRIVATE_"+campaignID, detailsJSON)
	if err != nil {
		return fmt.Errorf("failed to update private details: %v", err)
	}

	return nil
}

// SubmitForValidation submits campaign for validation
// Campaign hash is stored in StartupValidatorCollection for validator access
func (s *StartupContract) SubmitForValidation(
	ctx contractapi.TransactionContextInterface,
	campaignID string,
	documentsJSON string,
	submissionNotes string,
) error {

	// Get campaign from private collection
	campaignJSON, err := ctx.GetStub().GetPrivateData(StartupPrivateCollection, "CAMPAIGN_"+campaignID)
	if err != nil || campaignJSON == nil {
		return fmt.Errorf("campaign not found: %v", err)
	}

	var campaign Campaign
	err = json.Unmarshal(campaignJSON, &campaign)
	if err != nil {
		return fmt.Errorf("failed to unmarshal campaign: %v", err)
	}

	if campaign.Status != "DRAFT" && campaign.ValidationStatus != "ON_HOLD" {
		return fmt.Errorf("campaign cannot be submitted in current status")
	}

	// Parse documents
	var documents []string
	err = json.Unmarshal([]byte(documentsJSON), &documents)
	if err != nil {
		return fmt.Errorf("failed to parse documents: %v", err)
	}

	timestamp := time.Now().Format(time.RFC3339)

	// Get private details
	detailsJSON, err := ctx.GetStub().GetPrivateData(StartupPrivateCollection, "CAMPAIGN_PRIVATE_"+campaignID)
	if err != nil {
		return fmt.Errorf("failed to read private details: %v", err)
	}

	var privateDetails CampaignPrivateDetails
	err = json.Unmarshal(detailsJSON, &privateDetails)
	if err != nil {
		return fmt.Errorf("failed to unmarshal private details: %v", err)
	}

	// Record document submission
	submissionID := fmt.Sprintf("SUBMISSION_%s_%d", campaignID, len(privateDetails.DocumentHistory)+1)
	submission := DocumentSubmission{
		SubmissionID:    submissionID,
		Documents:       documents,
		SubmittedAt:     timestamp,
		SubmissionNotes: submissionNotes,
		ResponseStatus:  "PENDING",
	}

	privateDetails.DocumentHistory = append(privateDetails.DocumentHistory, submission)
	privateDetails.CurrentDocuments = documents

	// Update campaign status
	campaign.Status = "SUBMITTED"
	campaign.ValidationStatus = "PENDING_VALIDATION"
	campaign.UpdatedAt = timestamp

	// Generate campaign hash for validation
	campaignHash := s.generateCampaignHash(campaign)
	campaign.ValidationHash = campaignHash

	// Save updated campaign to private collection
	campaignJSON, _ = json.Marshal(campaign)
	err = ctx.GetStub().PutPrivateData(StartupPrivateCollection, "CAMPAIGN_"+campaignID, campaignJSON)
	if err != nil {
		return fmt.Errorf("failed to update campaign: %v", err)
	}

	// Save private details
	detailsJSON, _ = json.Marshal(privateDetails)
	err = ctx.GetStub().PutPrivateData(StartupPrivateCollection, "CAMPAIGN_PRIVATE_"+campaignID, detailsJSON)
	if err != nil {
		return fmt.Errorf("failed to update private details: %v", err)
	}

	// Share campaign data with ValidatorOrg via StartupValidatorCollection
	validatorData := map[string]interface{}{
		"campaignId":       campaign.CampaignID,
		"startupId":        campaign.StartupID,
		"projectName":      campaign.ProjectName,
		"category":         campaign.Category,
		"goalAmount":       campaign.GoalAmount,
		"description":      campaign.Description,
		"documents":        documents,
		"validationHash":   campaignHash,
		"validationStatus": campaign.ValidationStatus,
		"submittedAt":      timestamp,
	}

	validatorDataJSON, _ := json.Marshal(validatorData)
	err = ctx.GetStub().PutPrivateData(StartupValidatorCollection, "VALIDATION_REQUEST_"+campaignID, validatorDataJSON)
	if err != nil {
		return fmt.Errorf("failed to share with validator: %v", err)
	}

	return nil
}

// generateCampaignHash generates a hash for campaign verification
func (s *StartupContract) generateCampaignHash(campaign Campaign) string {
	data := fmt.Sprintf("%s|%s|%s|%s|%.2f|%s|%s",
		campaign.CampaignID,
		campaign.StartupID,
		campaign.ProjectName,
		campaign.Category,
		campaign.GoalAmount,
		campaign.Currency,
		campaign.CreatedAt,
	)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// ShareCampaignToPlatform shares validated campaign with Platform for publishing
// Only callable after Validator has APPROVED the campaign
func (s *StartupContract) ShareCampaignToPlatform(
	ctx contractapi.TransactionContextInterface,
	campaignID string,
	validationHash string,
) error {

	// Get campaign from private collection
	campaignJSON, err := ctx.GetStub().GetPrivateData(StartupPrivateCollection, "CAMPAIGN_"+campaignID)
	if err != nil || campaignJSON == nil {
		return fmt.Errorf("campaign not found: %v", err)
	}

	var campaign Campaign
	err = json.Unmarshal(campaignJSON, &campaign)
	if err != nil {
		return fmt.Errorf("failed to unmarshal campaign: %v", err)
	}

	// Verify campaign is validated
	if campaign.ValidationStatus != "APPROVED" {
		return fmt.Errorf("campaign must be APPROVED before sharing with platform. Current status: %s", campaign.ValidationStatus)
	}

	// Verify validation hash matches what validator provided
	validationStatusJSON, err := ctx.GetStub().GetPrivateData(StartupValidatorCollection, "VALIDATION_STATUS_"+campaignID)
	if err != nil || validationStatusJSON == nil {
		return fmt.Errorf("validation status not found")
	}

	var validationStatus map[string]interface{}
	err = json.Unmarshal(validationStatusJSON, &validationStatus)
	if err != nil {
		return fmt.Errorf("failed to unmarshal validation status: %v", err)
	}

	validatorHash, ok := validationStatus["validationHash"].(string)
	if !ok || validatorHash != validationHash {
		return fmt.Errorf("validation hash mismatch. Cannot share with platform")
	}

	timestamp := time.Now().Format(time.RFC3339)

	// Share complete campaign data with Platform via StartupPlatformCollection
	platformData := map[string]interface{}{
		// All 22 campaign parameters
		"campaignId":        campaign.CampaignID,
		"startupId":         campaign.StartupID,
		"category":          campaign.Category,
		"deadline":          campaign.Deadline,
		"currency":          campaign.Currency,
		"hasRaised":         campaign.HasRaised,
		"hasGovGrants":      campaign.HasGovGrants,
		"incorpDate":        campaign.IncorpDate,
		"projectStage":      campaign.ProjectStage,
		"sector":            campaign.Sector,
		"tags":              campaign.Tags,
		"teamAvailable":     campaign.TeamAvailable,
		"investorCommitted": campaign.InvestorCommitted,
		"duration":          campaign.Duration,
		"fundingDay":        campaign.FundingDay,
		"fundingMonth":      campaign.FundingMonth,
		"fundingYear":       campaign.FundingYear,
		"goalAmount":        campaign.GoalAmount,
		"investmentRange":   campaign.InvestmentRange,
		"projectName":       campaign.ProjectName,
		"description":       campaign.Description,
		"documents":         campaign.Documents,
		
		// Validation info
		"validationHash":     validationHash,
		"validationScore":    validationStatus["dueDiligenceScore"],
		"riskScore":          validationStatus["riskScore"],
		"riskLevel":          validationStatus["riskLevel"],
		"validationStatus":   campaign.ValidationStatus,
		"sharedAt":           timestamp,
	}

	platformJSON, err := json.Marshal(platformData)
	if err != nil {
		return fmt.Errorf("failed to marshal platform data: %v", err)
	}

	// Store in StartupPlatformCollection
	err = ctx.GetStub().PutPrivateData(StartupPlatformCollection, "CAMPAIGN_SHARE_"+campaignID, platformJSON)
	if err != nil {
		return fmt.Errorf("failed to share with platform: %v", err)
	}

	// Update campaign status
	campaign.Status = "PENDING_PLATFORM_APPROVAL"
	campaign.UpdatedAt = timestamp

	campaignJSON, _ = json.Marshal(campaign)
	err = ctx.GetStub().PutPrivateData(StartupPrivateCollection, "CAMPAIGN_"+campaignID, campaignJSON)
	if err != nil {
		return fmt.Errorf("failed to update campaign: %v", err)
	}

	return nil
}

// CheckPublishNotification checks if Platform has published the campaign
func (s *StartupContract) CheckPublishNotification(
	ctx contractapi.TransactionContextInterface,
	campaignID string,
) (string, error) {

	// Read notification from StartupPlatformCollection
	notificationJSON, err := ctx.GetStub().GetPrivateData(StartupPlatformCollection, "PUBLISH_NOTIFICATION_"+campaignID)
	if err != nil || notificationJSON == nil {
		return "", fmt.Errorf("no publish notification found for campaign %s", campaignID)
	}

	return string(notificationJSON), nil
}

// UpdateCampaignDocs updates campaign documents (when validator requests revisions)
func (s *StartupContract) UpdateCampaignDocs(
	ctx contractapi.TransactionContextInterface,
	campaignID string,
	documentsJSON string,
	submissionNotes string,
) error {

	// Similar to SubmitForValidation but for document updates
	return s.SubmitForValidation(ctx, campaignID, documentsJSON, submissionNotes)
}

// SubmitForPublishing submits approved campaign to Platform for publishing
func (s *StartupContract) SubmitForPublishing(
	ctx contractapi.TransactionContextInterface,
	campaignID string,
) error {

	// Get campaign
	campaignJSON, err := ctx.GetStub().GetPrivateData(StartupPrivateCollection, "CAMPAIGN_"+campaignID)
	if err != nil || campaignJSON == nil {
		return fmt.Errorf("campaign not found: %v", err)
	}

	var campaign Campaign
	err = json.Unmarshal(campaignJSON, &campaign)
	if err != nil {
		return fmt.Errorf("failed to unmarshal campaign: %v", err)
	}

	if campaign.ValidationStatus != "APPROVED" {
		return fmt.Errorf("campaign must be approved before publishing")
	}

	timestamp := time.Now().Format(time.RFC3339)
	campaign.PlatformStatus = "PENDING_PLATFORM"
	campaign.UpdatedAt = timestamp

	// Update campaign
	campaignJSON, _ = json.Marshal(campaign)
	err = ctx.GetStub().PutPrivateData(StartupPrivateCollection, "CAMPAIGN_"+campaignID, campaignJSON)
	if err != nil {
		return fmt.Errorf("failed to update campaign: %v", err)
	}

	// Share with Platform via StartupPlatformCollection
	platformData := map[string]interface{}{
		"campaignId":       campaign.CampaignID,
		"startupId":        campaign.StartupID,
		"projectName":      campaign.ProjectName,
		"description":      campaign.Description,
		"category":         campaign.Category,
		"goalAmount":       campaign.GoalAmount,
		"currency":         campaign.Currency,
		"openDate":         campaign.OpenDate,
		"closeDate":        campaign.CloseDate,
		"durationDays":     campaign.Duration,
		"validationScore":  campaign.ValidationScore,
		"validationHash":   campaign.ValidationHash,
		"submittedAt":      timestamp,
	}

	platformDataJSON, _ := json.Marshal(platformData)
	err = ctx.GetStub().PutPrivateData(StartupPlatformCollection, "PUBLISH_REQUEST_"+campaignID, platformDataJSON)
	if err != nil {
		return fmt.Errorf("failed to share with platform: %v", err)
	}

	// Also update public info
	publicInfo := CampaignPublicInfo{
		CampaignID:  campaignID,
		StartupID:   campaign.StartupID,
		ProjectName: campaign.ProjectName,
		Category:    campaign.Category,
		GoalAmount:  campaign.GoalAmount,
		Currency:    campaign.Currency,
		Status:      "PENDING_PLATFORM",
		PublishedAt: "",
	}

	publicJSON, _ := json.Marshal(publicInfo)
	ctx.GetStub().PutState("CAMPAIGN_PUBLIC_"+campaignID, publicJSON)

	return nil
}

// MarkCampaignCompleted marks a campaign as completed
func (s *StartupContract) MarkCampaignCompleted(
	ctx contractapi.TransactionContextInterface,
	campaignID string,
) error {

	campaignJSON, err := ctx.GetStub().GetPrivateData(StartupPrivateCollection, "CAMPAIGN_"+campaignID)
	if err != nil || campaignJSON == nil {
		return fmt.Errorf("campaign not found: %v", err)
	}

	var campaign Campaign
	err = json.Unmarshal(campaignJSON, &campaign)
	if err != nil {
		return fmt.Errorf("failed to unmarshal campaign: %v", err)
	}

	timestamp := time.Now().Format(time.RFC3339)
	campaign.Status = "COMPLETED"
	campaign.UpdatedAt = timestamp

	campaignJSON, _ = json.Marshal(campaign)
	err = ctx.GetStub().PutPrivateData(StartupPrivateCollection, "CAMPAIGN_"+campaignID, campaignJSON)
	if err != nil {
		return fmt.Errorf("failed to update campaign: %v", err)
	}

	// Update public info
	publicInfo := CampaignPublicInfo{
		CampaignID:  campaignID,
		StartupID:   campaign.StartupID,
		ProjectName: campaign.ProjectName,
		Category:    campaign.Category,
		GoalAmount:  campaign.GoalAmount,
		Currency:    campaign.Currency,
		Status:      "COMPLETED",
		PublishedAt: campaign.PublishedAt,
	}

	publicJSON, _ := json.Marshal(publicInfo)
	ctx.GetStub().PutState("CAMPAIGN_PUBLIC_"+campaignID, publicJSON)

	return nil
}

// ============================================================================
// INVESTMENT & PROPOSAL MANAGEMENT - Using PDC
// ============================================================================

// AcknowledgeInvestment acknowledges an investment from investor
// Investment is stored in StartupInvestorCollection
func (s *StartupContract) AcknowledgeInvestment(
	ctx contractapi.TransactionContextInterface,
	investmentID string,
	campaignID string,
	investorID string,
	amount float64,
) error {

	timestamp := time.Now().Format(time.RFC3339)

	investment := Investment{
		InvestmentID:   investmentID,
		CampaignID:     campaignID,
		InvestorID:     investorID,
		Amount:         amount,
		Currency:       "USD",
		Status:         "ACKNOWLEDGED",
		CommittedAt:    timestamp,
	}

	investmentJSON, err := json.Marshal(investment)
	if err != nil {
		return fmt.Errorf("failed to marshal investment: %v", err)
	}

	// Store in shared collection with InvestorOrg
	err = ctx.GetStub().PutPrivateData(StartupInvestorCollection, "INVESTMENT_"+investmentID, investmentJSON)
	if err != nil {
		return fmt.Errorf("failed to acknowledge investment: %v", err)
	}

	// Update campaign investor count and funds raised
	campaignJSON, err := ctx.GetStub().GetPrivateData(StartupPrivateCollection, "CAMPAIGN_"+campaignID)
	if err == nil && campaignJSON != nil {
		var campaign Campaign
		json.Unmarshal(campaignJSON, &campaign)
		campaign.InvestorCount++
		campaign.FundsRaisedAmount += amount
		campaign.FundsRaisedPercent = (campaign.FundsRaisedAmount / campaign.GoalAmount) * 100
		campaign.UpdatedAt = timestamp

		campaignJSON, _ = json.Marshal(campaign)
		ctx.GetStub().PutPrivateData(StartupPrivateCollection, "CAMPAIGN_"+campaignID, campaignJSON)
	}

	return nil
}

// RespondToInvestmentProposal responds to an investment proposal from investor
func (s *StartupContract) RespondToInvestmentProposal(
	ctx contractapi.TransactionContextInterface,
	proposalID string,
	action string,
	counterAmount float64,
	counterTerms string,
	modifiedMilestonesJSON string,
) error {

	// Get proposal from shared collection
	proposalJSON, err := ctx.GetStub().GetPrivateData(StartupInvestorCollection, "PROPOSAL_"+proposalID)
	if err != nil || proposalJSON == nil {
		return fmt.Errorf("proposal not found: %v", err)
	}

	var proposal map[string]interface{}
	err = json.Unmarshal(proposalJSON, &proposal)
	if err != nil {
		return fmt.Errorf("failed to unmarshal proposal: %v", err)
	}

	timestamp := time.Now().Format(time.RFC3339)

	// Update proposal based on action
	if action == "ACCEPT" {
		proposal["status"] = "ACCEPTED"
		proposal["startupAccepted"] = true
		proposal["acceptedAt"] = timestamp
	} else if action == "REJECT" {
		proposal["status"] = "REJECTED"
		proposal["rejectedAt"] = timestamp
	} else if action == "COUNTER" {
		proposal["status"] = "COUNTERED"
		proposal["counterAmount"] = counterAmount
		proposal["counterTerms"] = counterTerms
		proposal["negotiationRound"] = proposal["negotiationRound"].(float64) + 1
	}

	proposal["updatedAt"] = timestamp

	// Save updated proposal
	proposalJSON, _ = json.Marshal(proposal)
	err = ctx.GetStub().PutPrivateData(StartupInvestorCollection, "PROPOSAL_"+proposalID, proposalJSON)
	if err != nil {
		return fmt.Errorf("failed to update proposal: %v", err)
	}

	return nil
}

// ============================================================================
// MILESTONE MANAGEMENT - Using PDC
// ============================================================================

// SubmitMilestoneReport submits a milestone completion report
// Report is stored in StartupValidatorCollection for validator verification
func (s *StartupContract) SubmitMilestoneReport(
	ctx contractapi.TransactionContextInterface,
	reportID string,
	campaignID string,
	milestoneID string,
	agreementID string,
	title string,
	description string,
	evidenceJSON string,
) error {

	var evidence []string
	err := json.Unmarshal([]byte(evidenceJSON), &evidence)
	if err != nil {
		return fmt.Errorf("failed to parse evidence: %v", err)
	}

	timestamp := time.Now().Format(time.RFC3339)

	report := MilestoneReport{
		ReportID:    reportID,
		CampaignID:  campaignID,
		MilestoneID: milestoneID,
		AgreementID: agreementID,
		Title:       title,
		Description: description,
		Evidence:    evidence,
		Status:      "SUBMITTED",
		SubmittedAt: timestamp,
	}

	reportJSON, err := json.Marshal(report)
	if err != nil {
		return fmt.Errorf("failed to marshal report: %v", err)
	}

	// Store in shared collection with ValidatorOrg
	err = ctx.GetStub().PutPrivateData(StartupValidatorCollection, "MILESTONE_REPORT_"+reportID, reportJSON)
	if err != nil {
		return fmt.Errorf("failed to submit report: %v", err)
	}

	return nil
}

// ReceiveFunding records funding received from platform
func (s *StartupContract) ReceiveFunding(
	ctx contractapi.TransactionContextInterface,
	releaseID string,
	agreementID string,
	milestoneID string,
	amount float64,
) error {

	timestamp := time.Now().Format(time.RFC3339)

	fundingRecord := map[string]interface{}{
		"releaseId":   releaseID,
		"agreementId": agreementID,
		"milestoneId": milestoneID,
		"amount":      amount,
		"receivedAt":  timestamp,
	}

	recordJSON, err := json.Marshal(fundingRecord)
	if err != nil {
		return fmt.Errorf("failed to marshal funding record: %v", err)
	}

	// Store in StartupPlatformCollection
	err = ctx.GetStub().PutPrivateData(StartupPlatformCollection, "FUNDING_RECEIVED_"+releaseID, recordJSON)
	if err != nil {
		return fmt.Errorf("failed to record funding: %v", err)
	}

	return nil
}

// ============================================================================
// DISPUTE MANAGEMENT - Using PDC
// ============================================================================

// SubmitDispute submits a dispute (stored in ThreePartyCollection or AllOrgsCollection)
func (s *StartupContract) SubmitDispute(
	ctx contractapi.TransactionContextInterface,
	disputeID string,
	startupID string,
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
	err := json.Unmarshal([]byte(evidenceHashesJSON), &evidenceHashes)
	if err != nil {
		return fmt.Errorf("failed to parse evidence hashes: %v", err)
	}

	timestamp := time.Now().Format(time.RFC3339)
	submissionID := fmt.Sprintf("DISPUTE_SUB_%s", disputeID)

	dispute := DisputeSubmission{
		SubmissionID:   submissionID,
		DisputeID:      disputeID,
		StartupID:      startupID,
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

	// Store in AllOrgsCollection so Platform can manage it
	err = ctx.GetStub().PutPrivateData(AllOrgsCollection, "DISPUTE_"+disputeID, disputeJSON)
	if err != nil {
		return fmt.Errorf("failed to submit dispute: %v", err)
	}

	return nil
}

// SubmitDisputeEvidence submits additional evidence for a dispute
func (s *StartupContract) SubmitDisputeEvidence(
	ctx contractapi.TransactionContextInterface,
	disputeID string,
	evidenceHashesJSON string,
	evidenceDescription string,
) error {

	var evidenceHashes []string
	json.Unmarshal([]byte(evidenceHashesJSON), &evidenceHashes)

	timestamp := time.Now().Format(time.RFC3339)

	evidence := map[string]interface{}{
		"disputeId":            disputeID,
		"evidenceHashes":       evidenceHashes,
		"evidenceDescription":  evidenceDescription,
		"submittedBy":          "startup",
		"submittedAt":          timestamp,
	}

	evidenceJSON, _ := json.Marshal(evidence)
	evidenceKey := fmt.Sprintf("DISPUTE_EVIDENCE_%s_%s", disputeID, timestamp)
	
	err := ctx.GetStub().PutPrivateData(AllOrgsCollection, evidenceKey, evidenceJSON)
	if err != nil {
		return fmt.Errorf("failed to submit evidence: %v", err)
	}

	return nil
}

// RespondToDispute responds to a dispute filed against the startup
func (s *StartupContract) RespondToDispute(
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
		"respondedBy":           "startup",
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
// FEE MANAGEMENT
// ============================================================================

// RecordFeePayment records a fee payment
func (s *StartupContract) RecordFeePayment(
	ctx contractapi.TransactionContextInterface,
	recordID string,
	startupID string,
	campaignID string,
	feeType string,
	amount float64,
	transactionHash string,
) error {

	timestamp := time.Now().Format(time.RFC3339)

	feeRecord := FeePaymentRecord{
		RecordID:        recordID,
		StartupID:       startupID,
		CampaignID:      campaignID,
		FeeType:         feeType,
		Amount:          amount,
		TransactionHash: transactionHash,
		PaidAt:          timestamp,
	}

	feeJSON, err := json.Marshal(feeRecord)
	if err != nil {
		return fmt.Errorf("failed to marshal fee record: %v", err)
	}

	// Store in StartupPlatformCollection
	err = ctx.GetStub().PutPrivateData(StartupPlatformCollection, "FEE_PAYMENT_"+recordID, feeJSON)
	if err != nil {
		return fmt.Errorf("failed to record fee payment: %v", err)
	}

	return nil
}

// ============================================================================
// QUERY FUNCTIONS
// ============================================================================

// GetCampaign retrieves campaign from private collection (only accessible by StartupOrg)
func (s *StartupContract) GetCampaign(ctx contractapi.TransactionContextInterface, campaignID string) (*Campaign, error) {
	campaignJSON, err := ctx.GetStub().GetPrivateData(StartupPrivateCollection, "CAMPAIGN_"+campaignID)
	if err != nil {
		return nil, fmt.Errorf("failed to read campaign: %v", err)
	}
	if campaignJSON == nil {
		return nil, fmt.Errorf("campaign %s does not exist", campaignID)
	}

	var campaign Campaign
	err = json.Unmarshal(campaignJSON, &campaign)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal campaign: %v", err)
	}

	return &campaign, nil
}

// GetCampaignPublic retrieves public campaign info (accessible by all)
func (s *StartupContract) GetCampaignPublic(ctx contractapi.TransactionContextInterface, campaignID string) (*CampaignPublicInfo, error) {
	publicJSON, err := ctx.GetStub().GetState("CAMPAIGN_PUBLIC_" + campaignID)
	if err != nil {
		return nil, fmt.Errorf("failed to read public campaign: %v", err)
	}
	if publicJSON == nil {
		return nil, fmt.Errorf("public campaign info not found")
	}

	var publicInfo CampaignPublicInfo
	err = json.Unmarshal(publicJSON, &publicInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal public info: %v", err)
	}

	return &publicInfo, nil
}

// GetCampaignUpdateHistory retrieves campaign update history
func (s *StartupContract) GetCampaignUpdateHistory(ctx contractapi.TransactionContextInterface, campaignID string) (string, error) {
	detailsJSON, err := ctx.GetStub().GetPrivateData(StartupPrivateCollection, "CAMPAIGN_PRIVATE_"+campaignID)
	if err != nil || detailsJSON == nil {
		return "", fmt.Errorf("private details not found: %v", err)
	}

	var privateDetails CampaignPrivateDetails
	err = json.Unmarshal(detailsJSON, &privateDetails)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal private details: %v", err)
	}

	historyJSON, _ := json.Marshal(privateDetails.UpdateHistory)
	return string(historyJSON), nil
}

// GetCampaignDocumentHistory retrieves document submission history
func (s *StartupContract) GetCampaignDocumentHistory(ctx contractapi.TransactionContextInterface, campaignID string) (string, error) {
	detailsJSON, err := ctx.GetStub().GetPrivateData(StartupPrivateCollection, "CAMPAIGN_PRIVATE_"+campaignID)
	if err != nil || detailsJSON == nil {
		return "", fmt.Errorf("private details not found: %v", err)
	}

	var privateDetails CampaignPrivateDetails
	err = json.Unmarshal(detailsJSON, &privateDetails)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal private details: %v", err)
	}

	historyJSON, _ := json.Marshal(privateDetails.DocumentHistory)
	return string(historyJSON), nil
}

// GetStartupDisputes retrieves all disputes for a startup
func (s *StartupContract) GetStartupDisputes(ctx contractapi.TransactionContextInterface, startupID string) (string, error) {
	// Query disputes from AllOrgsCollection
	// This would require a rich query in CouchDB
	// For now, returning placeholder
	return `[]`, nil
}


