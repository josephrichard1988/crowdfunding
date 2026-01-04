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

// CampaignPrivateDetails - stored in StartupPrivateCollection
type CampaignPrivateDetails struct {
	CampaignID       string               `json:"campaignId"`
	DocumentHistory  []DocumentSubmission `json:"documentHistory"`
	CurrentDocuments []string             `json:"currentDocuments"`
	UpdateHistory    []CampaignUpdate     `json:"updateHistory"`
	InternalNotes    string               `json:"internalNotes"`
	FinancialDetails string               `json:"financialDetails"`
}

// CampaignPublicInfo - stored on public ledger (world state)
type CampaignPublicInfo struct {
	CampaignID  string  `json:"campaignId"`
	StartupID   string  `json:"startupId"`
	ProjectName string  `json:"projectName"`
	Category    string  `json:"category"`
	GoalAmount  float64 `json:"goalAmount"`
	Currency    string  `json:"currency"`
	Status      string  `json:"status"`
	PublishedAt string  `json:"publishedAt"`
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
	SubmissionID   string   `json:"submissionId"`
	DisputeID      string   `json:"disputeId"`
	StartupID      string   `json:"startupId"`
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

// Startup represents a startup entity owned by a user
type Startup struct {
	StartupID   string   `json:"startupId"`
	OwnerID     string   `json:"ownerId"` // orgUserId of the startup owner
	Name        string   `json:"name"`
	Description string   `json:"description"`
	DisplayID   string   `json:"displayId"`   // Human-readable ID like S-001
	CampaignIDs []string `json:"campaignIds"` // List of campaign IDs under this startup
	CreatedAt   string   `json:"createdAt"`
	UpdatedAt   string   `json:"updatedAt"`
}

// DeletionRecord tracks deletion of campaigns/startups with fee information
type DeletionRecord struct {
	DeletionID    string  `json:"deletionId"`
	EntityType    string  `json:"entityType"`    // "CAMPAIGN" or "STARTUP"
	EntityID      string  `json:"entityId"`      // campaignId or startupId
	EntityName    string  `json:"entityName"`    // Name for reference
	OwnerID       string  `json:"ownerId"`       // Owner of deleted entity
	FeeCharged    float64 `json:"feeCharged"`    // Fee in CFT
	FundsRaised   float64 `json:"fundsRaised"`   // Original funds raised
	FeePercentage float64 `json:"feePercentage"` // 60% or 0 for fixed fee
	Reason        string  `json:"reason"`        // User-provided reason
	DeletedAt     string  `json:"deletedAt"`
	TxID          string  `json:"txId"` // Transaction ID
}

// DeletionFeePreview for showing fee before deletion
type DeletionFeePreview struct {
	EntityID      string  `json:"entityId"`
	EntityType    string  `json:"entityType"`
	FundsRaised   float64 `json:"fundsRaised"`
	FeeAmount     float64 `json:"feeAmount"`
	FeePercentage float64 `json:"feePercentage"` // 60 or 0 (for fixed)
	IsFixedFee    bool    `json:"isFixedFee"`
}

// PublicDeletionRecord - minimal info for public visibility (no sensitive data)
type PublicDeletionRecord struct {
	DeletionID string  `json:"deletionId"`
	EntityType string  `json:"entityType"` // "CAMPAIGN" or "STARTUP"
	EntityID   string  `json:"entityId"`   // campaignId or startupId (no name)
	FeeCharged float64 `json:"feeCharged"` // Fee in CFT
	Reason     string  `json:"reason"`     // User-provided reason
	DeletedAt  string  `json:"deletedAt"`
	TxID       string  `json:"txId"` // Transaction ID for verification
}

// StartupDeletionFeeResult wraps multiple return values for CalculateStartupDeletionFee
// ContractAPI only allows single return value + error
type StartupDeletionFeeResult struct {
	TotalFee     float64              `json:"totalFee"`
	CampaignFees []DeletionFeePreview `json:"campaignFees"`
}

// StartupDeletionResult wraps multiple return values for DeleteStartup
// ContractAPI only allows single return value + error
type StartupDeletionResult struct {
	StartupDeletion   DeletionRecord   `json:"startupDeletion"`
	CampaignDeletions []DeletionRecord `json:"campaignDeletions"`
}

// ============================================================================
// INIT
// ============================================================================

func (s *StartupContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	fmt.Println("StartupOrg contract initialized with PDC support")
	return nil
}

// ============================================================================
// STARTUP MANAGEMENT - Startup must be created before campaigns
// ============================================================================

// CreateStartup creates a new startup entity for a user
func (s *StartupContract) CreateStartup(
	ctx contractapi.TransactionContextInterface,
	startupID string,
	ownerID string,
	name string,
	description string,
	displayID string,
) error {
	// Validate inputs
	if startupID == "" || ownerID == "" || name == "" {
		return fmt.Errorf("startupID, ownerID, and name are required")
	}

	// Check if startup already exists
	existingStartup, err := ctx.GetStub().GetPrivateData(StartupPrivateCollection, "STARTUP_"+startupID)
	if err != nil {
		return fmt.Errorf("failed to check existing startup: %v", err)
	}
	if existingStartup != nil {
		return fmt.Errorf("startup %s already exists", startupID)
	}

	// Create startup object
	startup := Startup{
		StartupID:   startupID,
		OwnerID:     ownerID,
		Name:        name,
		Description: description,
		DisplayID:   displayID,
		CampaignIDs: []string{},
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}

	// Store in private data collection
	startupJSON, err := json.Marshal(startup)
	if err != nil {
		return fmt.Errorf("failed to marshal startup: %v", err)
	}

	err = ctx.GetStub().PutPrivateData(StartupPrivateCollection, "STARTUP_"+startupID, startupJSON)
	if err != nil {
		return fmt.Errorf("failed to store startup: %v", err)
	}

	// Also store owner mapping for querying startups by owner
	// Key: OWNER_{ownerID}_{startupID}
	ownerKey := fmt.Sprintf("OWNER_%s_%s", ownerID, startupID)
	err = ctx.GetStub().PutPrivateData(StartupPrivateCollection, ownerKey, startupJSON)
	if err != nil {
		return fmt.Errorf("failed to store owner mapping: %v", err)
	}

	return nil
}

// GetStartup retrieves a startup by its ID
func (s *StartupContract) GetStartup(
	ctx contractapi.TransactionContextInterface,
	startupID string,
) (*Startup, error) {
	startupJSON, err := ctx.GetStub().GetPrivateData(StartupPrivateCollection, "STARTUP_"+startupID)
	if err != nil {
		return nil, fmt.Errorf("failed to read startup: %v", err)
	}
	if startupJSON == nil {
		return nil, fmt.Errorf("startup %s not found", startupID)
	}

	var startup Startup
	err = json.Unmarshal(startupJSON, &startup)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal startup: %v", err)
	}

	return &startup, nil
}

// GetStartupsByOwner retrieves all startups owned by a specific user
func (s *StartupContract) GetStartupsByOwner(
	ctx contractapi.TransactionContextInterface,
	ownerID string,
) ([]Startup, error) {
	// Query using range with owner prefix
	startKey := fmt.Sprintf("OWNER_%s_", ownerID)
	endKey := fmt.Sprintf("OWNER_%s_~", ownerID)

	resultsIterator, err := ctx.GetStub().GetPrivateDataByRange(StartupPrivateCollection, startKey, endKey)
	if err != nil {
		return nil, fmt.Errorf("failed to query startups by owner: %v", err)
	}
	defer resultsIterator.Close()

	var startups []Startup
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			continue
		}

		var startup Startup
		err = json.Unmarshal(queryResponse.Value, &startup)
		if err != nil {
			continue
		}
		startups = append(startups, startup)
	}

	if startups == nil {
		startups = []Startup{}
	}

	return startups, nil
}

// AddCampaignToStartup adds a campaign ID to a startup's campaign list
func (s *StartupContract) AddCampaignToStartup(
	ctx contractapi.TransactionContextInterface,
	startupID string,
	campaignID string,
) error {
	startup, err := s.GetStartup(ctx, startupID)
	if err != nil {
		return err
	}

	// Add campaign to list
	startup.CampaignIDs = append(startup.CampaignIDs, campaignID)
	startup.UpdatedAt = time.Now().Format(time.RFC3339)

	// Update in private data
	startupJSON, err := json.Marshal(startup)
	if err != nil {
		return fmt.Errorf("failed to marshal updated startup: %v", err)
	}

	err = ctx.GetStub().PutPrivateData(StartupPrivateCollection, "STARTUP_"+startupID, startupJSON)
	if err != nil {
		return fmt.Errorf("failed to update startup: %v", err)
	}

	// Update owner mapping too
	ownerKey := fmt.Sprintf("OWNER_%s_%s", startup.OwnerID, startupID)
	err = ctx.GetStub().PutPrivateData(StartupPrivateCollection, ownerKey, startupJSON)
	if err != nil {
		return fmt.Errorf("failed to update owner mapping: %v", err)
	}

	return nil
}

// ============================================================================
// DELETION WITH FEES
// ============================================================================

// CalculateCampaignDeletionFee calculates the fee for deleting a campaign
// Returns 60% of funds raised, or 100 CFT if no funds raised
func (s *StartupContract) CalculateCampaignDeletionFee(
	ctx contractapi.TransactionContextInterface,
	campaignID string,
) (*DeletionFeePreview, error) {
	campaign, err := s.GetCampaign(ctx, campaignID)
	if err != nil {
		return nil, err
	}

	var feeAmount float64
	var feePercentage float64
	isFixedFee := false

	fundsRaised := campaign.FundsRaisedAmount
	if fundsRaised > 0 {
		// 60% of funds raised
		feePercentage = 60.0
		feeAmount = fundsRaised * 0.60
	} else {
		// Fixed 100 CFT for campaigns with no funds
		feeAmount = 100.0
		isFixedFee = true
	}

	return &DeletionFeePreview{
		EntityID:      campaignID,
		EntityType:    "CAMPAIGN",
		FundsRaised:   fundsRaised,
		FeeAmount:     feeAmount,
		FeePercentage: feePercentage,
		IsFixedFee:    isFixedFee,
	}, nil
}

// DeleteCampaign deletes a campaign and records the deletion with fee
func (s *StartupContract) DeleteCampaign(
	ctx contractapi.TransactionContextInterface,
	campaignID string,
	reason string,
) (*DeletionRecord, error) {
	// Get campaign
	campaign, err := s.GetCampaign(ctx, campaignID)
	if err != nil {
		return nil, fmt.Errorf("campaign not found: %v", err)
	}

	// Calculate fee
	feePreview, err := s.CalculateCampaignDeletionFee(ctx, campaignID)
	if err != nil {
		return nil, err
	}

	// Create deletion record
	timestamp := time.Now()
	deletionID := fmt.Sprintf("DEL_CAMP_%s_%d", campaignID, timestamp.Unix())

	deletionRecord := DeletionRecord{
		DeletionID:    deletionID,
		EntityType:    "CAMPAIGN",
		EntityID:      campaignID,
		EntityName:    campaign.ProjectName,
		OwnerID:       campaign.StartupID,
		FeeCharged:    feePreview.FeeAmount,
		FundsRaised:   campaign.FundsRaisedAmount,
		FeePercentage: feePreview.FeePercentage,
		Reason:        reason,
		DeletedAt:     timestamp.Format(time.RFC3339),
		TxID:          ctx.GetStub().GetTxID(),
	}

	// Create public deletion record (no sensitive data - no owner, no name, no funds details)
	publicRecord := PublicDeletionRecord{
		DeletionID: deletionID,
		EntityType: "CAMPAIGN",
		EntityID:   campaignID,
		FeeCharged: feePreview.FeeAmount,
		Reason:     reason,
		DeletedAt:  timestamp.Format(time.RFC3339),
		TxID:       ctx.GetStub().GetTxID(),
	}

	// Store PUBLIC record in world state (visible in CouchDB - no sensitive data)
	publicJSON, err := json.Marshal(publicRecord)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal public deletion record: %v", err)
	}
	err = ctx.GetStub().PutState("DELETION_"+deletionID, publicJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to store public deletion record: %v", err)
	}

	// Store FULL record in private collection (with all details)
	deletionJSON, err := json.Marshal(deletionRecord)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal deletion record: %v", err)
	}
	err = ctx.GetStub().PutPrivateData(StartupPrivateCollection, "DELETION_"+deletionID, deletionJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to store private deletion record: %v", err)
	}

	// Delete campaign from private data
	err = ctx.GetStub().DelPrivateData(StartupPrivateCollection, "CAMPAIGN_"+campaignID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete campaign: %v", err)
	}

	// Delete campaign private details
	ctx.GetStub().DelPrivateData(StartupPrivateCollection, "CAMPAIGN_PRIVATE_"+campaignID)

	// Remove from startup's campaign list
	startup, _ := s.GetStartup(ctx, campaign.StartupID)
	if startup != nil {
		newCampaignIDs := []string{}
		for _, cid := range startup.CampaignIDs {
			if cid != campaignID {
				newCampaignIDs = append(newCampaignIDs, cid)
			}
		}
		startup.CampaignIDs = newCampaignIDs
		startup.UpdatedAt = timestamp.Format(time.RFC3339)
		startupJSON, _ := json.Marshal(startup)
		ctx.GetStub().PutPrivateData(StartupPrivateCollection, "STARTUP_"+startup.StartupID, startupJSON)
		ownerKey := fmt.Sprintf("OWNER_%s_%s", startup.OwnerID, startup.StartupID)
		ctx.GetStub().PutPrivateData(StartupPrivateCollection, ownerKey, startupJSON)
	}

	return &deletionRecord, nil
}

// CalculateStartupDeletionFee calculates total fee for deleting a startup
// Sum of individual campaign fees (60% each or 100 CFT if no funds)
func (s *StartupContract) CalculateStartupDeletionFee(
	ctx contractapi.TransactionContextInterface,
	startupID string,
) (*StartupDeletionFeeResult, error) {
	startup, err := s.GetStartup(ctx, startupID)
	if err != nil {
		return nil, err
	}

	if len(startup.CampaignIDs) == 0 {
		// No campaigns - fixed 100 CFT
		return &StartupDeletionFeeResult{
			TotalFee: 100.0,
			CampaignFees: []DeletionFeePreview{{
				EntityID:    startupID,
				EntityType:  "STARTUP",
				FundsRaised: 0,
				FeeAmount:   100.0,
				IsFixedFee:  true,
			}},
		}, nil
	}

	var totalFee float64
	campaignFees := []DeletionFeePreview{}

	for _, campaignID := range startup.CampaignIDs {
		feePreview, err := s.CalculateCampaignDeletionFee(ctx, campaignID)
		if err != nil {
			// If campaign can't be found, add fixed fee
			feePreview = &DeletionFeePreview{
				EntityID:   campaignID,
				EntityType: "CAMPAIGN",
				FeeAmount:  100.0,
				IsFixedFee: true,
			}
		}
		totalFee += feePreview.FeeAmount
		campaignFees = append(campaignFees, *feePreview)
	}

	return &StartupDeletionFeeResult{
		TotalFee:     totalFee,
		CampaignFees: campaignFees,
	}, nil
}

// DeleteStartup deletes a startup and all its campaigns, recording fees
func (s *StartupContract) DeleteStartup(
	ctx contractapi.TransactionContextInterface,
	startupID string,
	reason string,
) (*StartupDeletionResult, error) {
	startup, err := s.GetStartup(ctx, startupID)
	if err != nil {
		return nil, fmt.Errorf("startup not found: %v", err)
	}

	timestamp := time.Now()
	var campaignDeletions []DeletionRecord
	var totalFee float64

	// Delete all campaigns first
	for _, campaignID := range startup.CampaignIDs {
		campDeletion, err := s.DeleteCampaign(ctx, campaignID, reason+" (parent startup deleted)")
		if err == nil {
			campaignDeletions = append(campaignDeletions, *campDeletion)
			totalFee += campDeletion.FeeCharged
		}
	}

	// If no campaigns, fixed 100 CFT
	if len(startup.CampaignIDs) == 0 {
		totalFee = 100.0
	}

	// Create startup deletion record
	deletionID := fmt.Sprintf("DEL_STU_%s_%d", startupID, timestamp.Unix())

	startupDeletionRecord := DeletionRecord{
		DeletionID:    deletionID,
		EntityType:    "STARTUP",
		EntityID:      startupID,
		EntityName:    startup.Name,
		OwnerID:       startup.OwnerID,
		FeeCharged:    totalFee,
		FundsRaised:   0, // Aggregated from campaigns
		FeePercentage: 0, // Mixed
		Reason:        reason,
		DeletedAt:     timestamp.Format(time.RFC3339),
		TxID:          ctx.GetStub().GetTxID(),
	}

	// Create public deletion record (no sensitive data)
	publicRecord := PublicDeletionRecord{
		DeletionID: deletionID,
		EntityType: "STARTUP",
		EntityID:   startupID,
		FeeCharged: totalFee,
		Reason:     reason,
		DeletedAt:  timestamp.Format(time.RFC3339),
		TxID:       ctx.GetStub().GetTxID(),
	}

	// Store PUBLIC record in world state (visible in CouchDB - no sensitive data)
	publicJSON, _ := json.Marshal(publicRecord)
	ctx.GetStub().PutState("DELETION_"+deletionID, publicJSON)

	// Store FULL record in private collection (with all details)
	deletionJSON, _ := json.Marshal(startupDeletionRecord)
	ctx.GetStub().PutPrivateData(StartupPrivateCollection, "DELETION_"+deletionID, deletionJSON)

	// Delete startup from private data
	ctx.GetStub().DelPrivateData(StartupPrivateCollection, "STARTUP_"+startupID)

	// Delete owner mapping
	ownerKey := fmt.Sprintf("OWNER_%s_%s", startup.OwnerID, startupID)
	ctx.GetStub().DelPrivateData(StartupPrivateCollection, ownerKey)

	return &StartupDeletionResult{
		StartupDeletion:   startupDeletionRecord,
		CampaignDeletions: campaignDeletions,
	}, nil
}

// GetDeletionRecord retrieves PUBLIC deletion record by ID (visible in CouchDB)
// Contains: entityId, entityType, feeCharged, reason, deletedAt, txId
// Does NOT contain: ownerID, entityName, fundsRaised (these are private)
func (s *StartupContract) GetDeletionRecord(
	ctx contractapi.TransactionContextInterface,
	deletionID string,
) (*PublicDeletionRecord, error) {
	deletionJSON, err := ctx.GetStub().GetState("DELETION_" + deletionID)
	if err != nil || deletionJSON == nil {
		return nil, fmt.Errorf("deletion record not found: %s", deletionID)
	}

	var record PublicDeletionRecord
	err = json.Unmarshal(deletionJSON, &record)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal deletion record: %v", err)
	}

	return &record, nil
}

// GetPrivateDeletionRecord retrieves FULL deletion record from private collection
// Contains all details including ownerID, entityName, fundsRaised
func (s *StartupContract) GetPrivateDeletionRecord(
	ctx contractapi.TransactionContextInterface,
	deletionID string,
) (*DeletionRecord, error) {
	deletionJSON, err := ctx.GetStub().GetPrivateData(StartupPrivateCollection, "DELETION_"+deletionID)
	if err != nil || deletionJSON == nil {
		return nil, fmt.Errorf("private deletion record not found: %s", deletionID)
	}

	var record DeletionRecord
	err = json.Unmarshal(deletionJSON, &record)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal deletion record: %v", err)
	}

	return &record, nil
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

	// Check if campaign already exists
	existingCampaignJSON, err := ctx.GetStub().GetPrivateData(StartupPrivateCollection, "CAMPAIGN_"+campaignID)
	if err != nil {
		return fmt.Errorf("failed to read from private collection: %v", err)
	}
	if existingCampaignJSON != nil {
		return fmt.Errorf("campaign %s already exists", campaignID)
	}

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
		CampaignID:          campaignID,
		StartupID:           startupID,
		Category:            category,
		Deadline:            deadline,
		Currency:            currency,
		HasRaised:           hasRaisedBool,
		HasGovGrants:        hasGovGrantsBool,
		IncorpDate:          incorpDate,
		ProjectStage:        projectStage,
		Sector:              sector,
		Tags:                tags,
		TeamAvailable:       teamAvailableBool,
		InvestorCommitted:   investorCommittedBool,
		Duration:            durationInt,
		FundingDay:          fundingDayInt,
		FundingMonth:        fundingMonthInt,
		FundingYear:         fundingYearInt,
		GoalAmount:          goalAmountFloat,
		InvestmentRange:     investmentRange,
		ProjectName:         projectName,
		Description:         description,
		Documents:           documents,
		OpenDate:            openDate,
		CloseDate:           closeDate,
		FundsRaisedAmount:   0,
		FundsRaisedPercent:  0,
		Status:              "DRAFT",
		ValidationStatus:    "NOT_SUBMITTED",
		ValidationScore:     0,
		SubmissionHash:      "",
		ValidationProofHash: "",
		InvestorCount:       0,
		PlatformStatus:      "NOT_PUBLISHED",
		CreatedAt:           timestamp,
		UpdatedAt:           timestamp,
		ApprovedAt:          "",
		PublishedAt:         "",
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

	// Add campaign to startup's campaign list
	err = s.AddCampaignToStartup(ctx, startupID, campaignID)
	if err != nil {
		return fmt.Errorf("failed to add campaign to startup: %v", err)
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

	// Generate campaign hash for validation (Submission Integrity)
	campaignHash := s.generateSubmissionHash(campaign)
	campaign.SubmissionHash = campaignHash

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
		"submissionHash":   campaignHash, /* Was validationHash */
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

// generateSubmissionHash generates a hash for campaign verification (Submission)
func (s *StartupContract) generateSubmissionHash(campaign Campaign) string {
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
	validationProofHash string, /* Was validationHash */
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

	// Sync validation status from StartupValidatorCollection (same as GetCampaign)
	statusJSON, _ := ctx.GetStub().GetPrivateData(StartupValidatorCollection, "VALIDATION_STATUS_"+campaignID)
	if statusJSON != nil {
		var statusUpdate map[string]interface{}
		json.Unmarshal(statusJSON, &statusUpdate)

		if val, ok := statusUpdate["status"].(string); ok {
			campaign.ValidationStatus = val
		}
		if val, ok := statusUpdate["validationProofHash"].(string); ok {
			campaign.ValidationProofHash = val
		}
		if val, ok := statusUpdate["riskLevel"].(string); ok {
			campaign.RiskLevel = val
		}
		if val, ok := statusUpdate["riskScore"].(float64); ok {
			campaign.ValidationScore = val
		}
		if val, ok := statusUpdate["dueDiligenceScore"].(float64); ok {
			if campaign.ValidationScore == 0 {
				campaign.ValidationScore = val
			}
		}
	}

	// Verify campaign is validated
	if campaign.ValidationStatus != "APPROVED" {
		return fmt.Errorf("campaign must be APPROVED before sharing with platform. Current status: %s", campaign.ValidationStatus)
	}

	// Verify campaign is not already shared or published
	if campaign.Status == "PENDING_PLATFORM_APPROVAL" || campaign.Status == "PUBLISHED" {
		return fmt.Errorf("campaign is already shared or published. Status: %s", campaign.Status)
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

	validatorHash, ok := validationStatus["validationProofHash"].(string)
	if !ok || validatorHash != validationProofHash {
		return fmt.Errorf("validation proof hash mismatch. Cannot share with platform")
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
		"validationProofHash": validationProofHash,
		"validationScore":     validationStatus["dueDiligenceScore"],
		"riskScore":           validationStatus["riskScore"],
		"riskLevel":           validationStatus["riskLevel"],
		"validationStatus":    campaign.ValidationStatus,
		"sharedAt":            timestamp,
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
		"campaignId":          campaign.CampaignID,
		"startupId":           campaign.StartupID,
		"projectName":         campaign.ProjectName,
		"description":         campaign.Description,
		"category":            campaign.Category,
		"goalAmount":          campaign.GoalAmount,
		"currency":            campaign.Currency,
		"openDate":            campaign.OpenDate,
		"closeDate":           campaign.CloseDate,
		"durationDays":        campaign.Duration,
		"validationScore":     campaign.ValidationScore,
		"validationProofHash": campaign.ValidationProofHash,
		"submittedAt":         timestamp,
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
		InvestmentID: investmentID,
		CampaignID:   campaignID,
		InvestorID:   investorID,
		Amount:       amount,
		Currency:     "USD",
		Status:       "ACKNOWLEDGED",
		CommittedAt:  timestamp,
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
		"disputeId":           disputeID,
		"evidenceHashes":      evidenceHashes,
		"evidenceDescription": evidenceDescription,
		"submittedBy":         "startup",
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
	// 1. Get fundamental campaign data
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

	// 2. Check for validation updates in StartupValidatorCollection
	statusJSON, _ := ctx.GetStub().GetPrivateData(StartupValidatorCollection, "VALIDATION_STATUS_"+campaignID)
	if statusJSON != nil {
		var statusUpdate map[string]interface{}
		json.Unmarshal(statusJSON, &statusUpdate)

		if val, ok := statusUpdate["status"].(string); ok {
			campaign.ValidationStatus = val
		}
		if val, ok := statusUpdate["validationProofHash"].(string); ok {
			campaign.ValidationProofHash = val
		}
		if val, ok := statusUpdate["riskLevel"].(string); ok {
			campaign.RiskLevel = val // Assuming RiskLevel exists in Campaign struct, otherwise ignore or map
		}
		if val, ok := statusUpdate["riskScore"].(float64); ok {
			campaign.ValidationScore = val // Mapping risk/dd score to ValidationScore
		}
		if val, ok := statusUpdate["dueDiligenceScore"].(float64); ok {
			if campaign.ValidationScore == 0 {
				campaign.ValidationScore = val
			}
		}

		// If status is APPROVED, update main status if it was SUBMITTED
		if campaign.ValidationStatus == "APPROVED" && campaign.Status == "SUBMITTED" {
			// We effectively consider it pre-approved locally, though official status change might require invoke
			// But for display, showing APPROVED is correct
		}
	}

	// 3. Check for publication updates in StartupPlatformCollection
	publishNoteJSON, _ := ctx.GetStub().GetPrivateData(StartupPlatformCollection, "PUBLISH_NOTIFICATION_"+campaignID)
	if publishNoteJSON != nil {
		var publishNote map[string]interface{}
		json.Unmarshal(publishNoteJSON, &publishNote)

		if status, ok := publishNote["status"].(string); ok && status == "PUBLISHED" {
			campaign.Status = "PUBLISHED"
			campaign.PlatformStatus = "PUBLISHED"

			if val, ok := publishNote["publishedAt"].(string); ok {
				campaign.PublishedAt = val
			}
		}
	}

	// 4. Calculate Funds Raised from StartupInvestorCollection
	// We need to iterate over investments for this campaign.
	// Since we can't easily query by partial key in private data without an index or knowing IDs,
	// checking if we have a "INVESTMENT_INDEX" or similar would be best.
	// However, for now, we'll try a range query if supported in private data (GetPrivateDataByRange).

	investmentIterator, err := ctx.GetStub().GetPrivateDataByRange(StartupInvestorCollection, "INVESTMENT_", "INVESTMENT_~")
	if err == nil {
		defer investmentIterator.Close()

		var totalRaised float64
		// simple set to count unique investors
		uniqueInvestors := make(map[string]bool)

		for investmentIterator.HasNext() {
			response, err := investmentIterator.Next()
			if err != nil {
				continue
			}

			var investment map[string]interface{}
			if err := json.Unmarshal(response.Value, &investment); err != nil {
				continue
			}

			// Check if this investment belongs to the current campaign
			if invCampaignID, ok := investment["campaignId"].(string); ok && invCampaignID == campaignID {
				// Check status
				if status, ok := investment["status"].(string); ok && (status == "COMMITTED" || status == "RELEASED") {
					if amount, ok := investment["amount"].(float64); ok {
						totalRaised += amount
					}
					if invID, ok := investment["investorId"].(string); ok {
						uniqueInvestors[invID] = true
					}
				}
			}
		}

		campaign.FundsRaisedAmount = totalRaised
		campaign.InvestorCount = len(uniqueInvestors)
		if campaign.GoalAmount > 0 {
			campaign.FundsRaisedPercent = (totalRaised / campaign.GoalAmount) * 100
		}

		// Map boolean flags if raised > 0
		if totalRaised > 0 {
			campaign.HasRaised = true
		}
	}

	return &campaign, nil
}

// GetAllCampaigns retrieves all campaigns from private collection (only accessible by StartupOrg)
func (s *StartupContract) GetAllCampaigns(ctx contractapi.TransactionContextInterface) ([]Campaign, error) {
	// Use GetPrivateDataByRange to iterate through all campaigns
	resultsIterator, err := ctx.GetStub().GetPrivateDataByRange(StartupPrivateCollection, "CAMPAIGN_", "CAMPAIGN_~")
	if err != nil {
		return nil, fmt.Errorf("failed to get campaigns: %v", err)
	}
	defer resultsIterator.Close()

	var campaigns []Campaign
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			continue
		}

		// Skip CAMPAIGN_PRIVATE_ keys - only process CAMPAIGN_{ID} keys
		key := queryResponse.Key
		if len(key) > 17 && key[:17] == "CAMPAIGN_PRIVATE_" {
			continue
		}

		var campaign Campaign
		err = json.Unmarshal(queryResponse.Value, &campaign)
		if err != nil {
			continue
		}

		// Get validation status updates
		statusJSON, _ := ctx.GetStub().GetPrivateData(StartupValidatorCollection, "VALIDATION_STATUS_"+campaign.CampaignID)
		if statusJSON != nil {
			var statusUpdate map[string]interface{}
			json.Unmarshal(statusJSON, &statusUpdate)
			if val, ok := statusUpdate["status"].(string); ok {
				campaign.ValidationStatus = val
			}
			if val, ok := statusUpdate["validationProofHash"].(string); ok {
				campaign.ValidationProofHash = val
			}
		}

		// Ensure arrays are never null (fix schema validation)
		if campaign.Tags == nil {
			campaign.Tags = []string{}
		}
		if campaign.Documents == nil {
			campaign.Documents = []string{}
		}

		campaigns = append(campaigns, campaign)
	}

	if campaigns == nil {
		campaigns = []Campaign{}
	}

	return campaigns, nil
}

// GetCampaignsByStartupId retrieves campaigns created by a specific startup user
// This enables strict user isolation - each user only sees their own campaigns
func (s *StartupContract) GetCampaignsByStartupId(ctx contractapi.TransactionContextInterface, startupId string) ([]Campaign, error) {
	// Validate input
	if startupId == "" {
		return []Campaign{}, nil
	}

	// Use GetPrivateDataByRange to iterate through all campaigns
	resultsIterator, err := ctx.GetStub().GetPrivateDataByRange(StartupPrivateCollection, "CAMPAIGN_", "CAMPAIGN_~")
	if err != nil {
		return nil, fmt.Errorf("failed to get campaigns: %v", err)
	}
	defer resultsIterator.Close()

	var campaigns []Campaign
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			continue
		}

		// Skip CAMPAIGN_PRIVATE_ keys - only process CAMPAIGN_{ID} keys
		key := queryResponse.Key
		if len(key) > 17 && key[:17] == "CAMPAIGN_PRIVATE_" {
			continue
		}

		var campaign Campaign
		err = json.Unmarshal(queryResponse.Value, &campaign)
		if err != nil {
			continue
		}

		// Filter by startupId - strict user isolation
		if campaign.StartupID != startupId {
			continue
		}

		// Get validation status updates
		statusJSON, _ := ctx.GetStub().GetPrivateData(StartupValidatorCollection, "VALIDATION_STATUS_"+campaign.CampaignID)
		if statusJSON != nil {
			var statusUpdate map[string]interface{}
			json.Unmarshal(statusJSON, &statusUpdate)
			if val, ok := statusUpdate["status"].(string); ok {
				campaign.ValidationStatus = val
			}
			if val, ok := statusUpdate["validationProofHash"].(string); ok {
				campaign.ValidationProofHash = val
			}
		}

		// Ensure arrays are never null (fix schema validation)
		if campaign.Tags == nil {
			campaign.Tags = []string{}
		}
		if campaign.Documents == nil {
			campaign.Documents = []string{}
		}

		campaigns = append(campaigns, campaign)
	}

	if campaigns == nil {
		campaigns = []Campaign{}
	}

	return campaigns, nil
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

// ============================================================================
// TOKEN INTEGRATION FUNCTIONS (CFT/CFRT)
// ============================================================================

// GetFeePayments retrieves all fee payments for a startup
func (s *StartupContract) GetFeePayments(
	ctx contractapi.TransactionContextInterface,
	startupID string,
) (string, error) {

	// Get fee payments from StartupPrivateCollection
	iterator, err := ctx.GetStub().GetPrivateDataByRange(StartupPrivateCollection, "FEE_PAYMENT_", "FEE_PAYMENT_~")
	if err != nil {
		return "[]", fmt.Errorf("failed to query fee payments: %v", err)
	}
	defer iterator.Close()

	var payments []FeePaymentRecord
	for iterator.HasNext() {
		queryResponse, err := iterator.Next()
		if err != nil {
			continue
		}

		var record FeePaymentRecord
		err = json.Unmarshal(queryResponse.Value, &record)
		if err != nil {
			continue
		}

		if record.StartupID == startupID {
			payments = append(payments, record)
		}
	}

	if payments == nil {
		payments = []FeePaymentRecord{}
	}

	paymentsJSON, err := json.Marshal(payments)
	if err != nil {
		return "[]", fmt.Errorf("failed to marshal payments: %v", err)
	}

	return string(paymentsJSON), nil
}

// GetRequiredFees returns the fees required for campaign operations
func (s *StartupContract) GetRequiredFees(
	ctx contractapi.TransactionContextInterface,
) (string, error) {

	// Fee schedule based on 1 INR = 2.5 CFT
	fees := map[string]interface{}{
		"registrationFee": map[string]interface{}{
			"amountCFT":   250,
			"amountINR":   100,
			"description": "One-time startup registration fee",
		},
		"campaignCreationFee": map[string]interface{}{
			"amountCFT":   1250,
			"amountINR":   500,
			"description": "Fee for creating a new campaign",
		},
		"campaignPublishingFee": map[string]interface{}{
			"amountCFT":   2500,
			"amountINR":   1000,
			"description": "Fee for publishing campaign to portal",
		},
		"validationFee": map[string]interface{}{
			"amountCFT":   500,
			"amountINR":   200,
			"description": "Fee for campaign validation",
		},
	}

	feesJSON, err := json.Marshal(fees)
	if err != nil {
		return "", fmt.Errorf("failed to marshal fees: %v", err)
	}

	return string(feesJSON), nil
}
