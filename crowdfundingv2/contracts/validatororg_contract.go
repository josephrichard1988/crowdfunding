package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// ValidatorContract provides functions for ValidatorOrg operations using PDC
type ValidatorContract struct {
	contractapi.Contract
}

// ============================================================================
// DATA STRUCTURES
// ============================================================================

// ValidationRecord represents a campaign validation result
type ValidationRecord struct {
	ValidationID       string              `json:"validationId"`
	CampaignID         string              `json:"campaignId"`
	CampaignHash       string              `json:"campaignHash"`
	ValidatorID        string              `json:"validatorId"`
	Status             string              `json:"status"`
	DocumentsVerified  bool                `json:"documentsVerified"`
	ComplianceCheck    bool                `json:"complianceCheck"`
	DueDiligenceScore  float64             `json:"dueDiligenceScore"`
	RiskScore          float64             `json:"riskScore"`
	RiskLevel          string              `json:"riskLevel"`
	Comments           []string            `json:"comments"`
	Issues             []string            `json:"issues"`
	RequiredDocuments  string              `json:"requiredDocuments"`
	ValidationAttempts []ValidationAttempt `json:"validationAttempts"`
	ValidatedAt        string              `json:"validatedAt"`
	CreatedAt          string              `json:"createdAt"`
}

// ValidationAttempt tracks each validation attempt
type ValidationAttempt struct {
	AttemptID         string   `json:"attemptId"`
	AttemptNumber     int      `json:"attemptNumber"`
	DocumentsReviewed []string `json:"documentsReviewed"`
	Status            string   `json:"status"`
	Score             float64  `json:"score"`
	Comments          string   `json:"comments"`
	RequiredDocs      string   `json:"requiredDocs"`
	AttemptedAt       string   `json:"attemptedAt"`
}

// RiskInsight represents risk information shared with investors
type RiskInsight struct {
	InsightID      string   `json:"insightId"`
	CampaignID     string   `json:"campaignId"`
	InvestorID     string   `json:"investorId"`
	RiskScore      float64  `json:"riskScore"`
	RiskLevel      string   `json:"riskLevel"`
	RiskFactors    []string `json:"riskFactors"`
	QueryResponse  string   `json:"queryResponse"`
	Recommendation string   `json:"recommendation"`
	CreatedAt      string   `json:"createdAt"`
}

// ValidationReport represents detailed report sent to PlatformOrg
type ValidationReport struct {
	ReportID        string  `json:"reportId"`
	CampaignID      string  `json:"campaignId"`
	ValidationID    string  `json:"validationId"`
	CampaignHash    string  `json:"campaignHash"`
	OverallScore    float64 `json:"overallScore"`
	DocumentScore   float64 `json:"documentScore"`
	ComplianceScore float64 `json:"complianceScore"`
	RiskScore       float64 `json:"riskScore"`
	Approved        bool    `json:"approved"`
	ReportSummary   string  `json:"reportSummary"`
	ReportHash      string  `json:"reportHash"`
	CreatedAt       string  `json:"createdAt"`
}

// ValidationProof for public ledger
type ValidationProof struct {
	ProofID        string `json:"proofId"`
	CampaignID     string `json:"campaignId"`
	ValidationHash string `json:"validationHash"`
	Status         string `json:"status"`
	PublishedAt    string `json:"publishedAt"`
}

// BlacklistedCampaign tracks rejected campaigns
type BlacklistedCampaign struct {
	CampaignID    string `json:"campaignId"`
	Reason        string `json:"reason"`
	BlacklistedAt string `json:"blacklistedAt"`
	BlacklistedBy string `json:"blacklistedBy"`
}

// MilestoneValidation represents milestone verification by Validator
type MilestoneValidation struct {
	VerificationID       string  `json:"verificationId"`
	MilestoneID          string  `json:"milestoneId"`
	CampaignID           string  `json:"campaignId"`
	StartupID            string  `json:"startupId"`
	MilestoneReportHash  string  `json:"milestoneReportHash"`
	DeliverablesVerified bool    `json:"deliverablesVerified"`
	QualityScore         float64 `json:"qualityScore"`
	Comments             string  `json:"comments"`
	Approved             bool    `json:"approved"`
	VerifiedAt           string  `json:"verifiedAt"`
}

// AgreementWitness represents Validator witnessing an agreement
type AgreementWitness struct {
	WitnessID         string  `json:"witnessId"`
	AgreementID       string  `json:"agreementId"`
	CampaignID        string  `json:"campaignId"`
	StartupID         string  `json:"startupId"`
	InvestorID        string  `json:"investorId"`
	InvestmentAmount  float64 `json:"investmentAmount"`
	ValidatorComments string  `json:"validatorComments"`
	WitnessedAt       string  `json:"witnessedAt"`
}

// DisputeInvestigation represents validator's investigation of a dispute
type DisputeInvestigation struct {
	InvestigationID    string                 `json:"investigationId"`
	DisputeID          string                 `json:"disputeId"`
	ValidatorID        string                 `json:"validatorId"`
	DisputeType        string                 `json:"disputeType"`
	InitiatorID        string                 `json:"initiatorId"`
	InitiatorType      string                 `json:"initiatorType"`
	RespondentID       string                 `json:"respondentId"`
	RespondentType     string                 `json:"respondentType"`
	CampaignID         string                 `json:"campaignId"`
	Status             string                 `json:"status"`
	Findings           []InvestigationFinding `json:"findings"`
	EvidenceReviewed   []string               `json:"evidenceReviewed"`
	TransactionLogs    []string               `json:"transactionLogs"`
	Recommendation     string                 `json:"recommendation"`
	RecommendedPenalty string                 `json:"recommendedPenalty"`
	AssignedAt         string                 `json:"assignedAt"`
	CompletedAt        string                 `json:"completedAt"`
}

// InvestigationFinding represents a finding during investigation
type InvestigationFinding struct {
	FindingID       string `json:"findingId"`
	FindingType     string `json:"findingType"`
	Description     string `json:"description"`
	Severity        string `json:"severity"`
	RelatedEvidence string `json:"relatedEvidence"`
	RecordedAt      string `json:"recordedAt"`
}

// ValidatorDisputeResponse represents validator's response when they are respondent
type ValidatorDisputeResponse struct {
	ResponseID     string   `json:"responseId"`
	DisputeID      string   `json:"disputeId"`
	ValidatorID    string   `json:"validatorId"`
	ResponseText   string   `json:"responseText"`
	Justification  string   `json:"justification"`
	SupportingDocs []string `json:"supportingDocs"`
	RespondedAt    string   `json:"respondedAt"`
}

// MilestoneInvestigation for investigating milestone-related disputes
type MilestoneInvestigation struct {
	InvestigationID    string  `json:"investigationId"`
	DisputeID          string  `json:"disputeId"`
	MilestoneID        string  `json:"milestoneId"`
	CampaignID         string  `json:"campaignId"`
	ValidatorID        string  `json:"validatorId"`
	MilestoneReviewed  bool    `json:"milestoneReviewed"`
	DeliverableStatus  string  `json:"deliverableStatus"`
	QualityAssessment  float64 `json:"qualityAssessment"`
	TimelineAssessment string  `json:"timelineAssessment"`
	DelayJustified     bool    `json:"delayJustified"`
	RecommendedAction  string  `json:"recommendedAction"`
	Comments           string  `json:"comments"`
	InvestigatedAt     string  `json:"investigatedAt"`
}

// ============================================================================
// INIT
// ============================================================================

func (v *ValidatorContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	fmt.Println("ValidatorOrg contract initialized with PDC support")
	return nil
}

// ============================================================================
// CAMPAIGN VALIDATION - Using PDC
// ============================================================================

// ValidateCampaign performs validation checks on submitted campaign
// Reads from StartupValidatorCollection, stores validation in ValidatorPrivateCollection
func (v *ValidatorContract) ValidateCampaign(
	ctx contractapi.TransactionContextInterface,
	validationID string,
	campaignID string,
	validatorID string,
	campaignHash string,
	documentsReviewedJSON string,
) error {

	// Get campaign submission from StartupValidatorCollection
	submissionJSON, err := ctx.GetStub().GetPrivateData(StartupValidatorCollection, "VALIDATION_REQUEST_"+campaignID)
	if err != nil || submissionJSON == nil {
		return fmt.Errorf("validation request not found: %v", err)
	}

	// Check for existing validation to prevent re-validation of final states
	existingValJSON, _ := ctx.GetStub().GetPrivateData(ValidatorPrivateCollection, "VALIDATION_"+validationID)
	if existingValJSON != nil {
		var existingVal ValidationRecord
		json.Unmarshal(existingValJSON, &existingVal)
		if existingVal.Status == "APPROVED" || existingVal.Status == "REJECTED" {
			return fmt.Errorf("validation %s is already final: %s", validationID, existingVal.Status)
		}
	}

	var documentsReviewed []string
	if documentsReviewedJSON != "" {
		json.Unmarshal([]byte(documentsReviewedJSON), &documentsReviewed)
	}

	timestamp := time.Now().Format(time.RFC3339)

	// Create validation record
	validation := ValidationRecord{
		ValidationID:       validationID,
		CampaignID:         campaignID,
		CampaignHash:       campaignHash,
		ValidatorID:        validatorID,
		Status:             "IN_PROGRESS",
		DocumentsVerified:  false,
		ComplianceCheck:    false,
		DueDiligenceScore:  0,
		RiskScore:          0,
		RiskLevel:          "PENDING",
		Comments:           []string{},
		Issues:             []string{},
		RequiredDocuments:  "",
		ValidationAttempts: []ValidationAttempt{},
		CreatedAt:          timestamp,
	}

	// Add initial validation attempt
	attempt := ValidationAttempt{
		AttemptID:         fmt.Sprintf("ATTEMPT_%s_1", validationID),
		AttemptNumber:     1,
		DocumentsReviewed: documentsReviewed,
		Status:            "IN_PROGRESS",
		Score:             0,
		Comments:          "Initial validation started",
		AttemptedAt:       timestamp,
	}
	validation.ValidationAttempts = append(validation.ValidationAttempts, attempt)

	validationJSON, err := json.Marshal(validation)
	if err != nil {
		return fmt.Errorf("failed to marshal validation: %v", err)
	}

	// Store in ValidatorPrivateCollection
	err = ctx.GetStub().PutPrivateData(ValidatorPrivateCollection, "VALIDATION_"+validationID, validationJSON)
	if err != nil {
		return fmt.Errorf("failed to store validation: %v", err)
	}

	return nil
}

// ApproveOrRejectCampaign approves, rejects, or puts campaign on hold
func (v *ValidatorContract) ApproveOrRejectCampaign(
	ctx contractapi.TransactionContextInterface,
	validationID string,
	campaignID string,
	status string,
	dueDiligenceScore float64,
	riskScore float64,
	riskLevel string,
	commentsJSON string,
	issuesJSON string,
	requiredDocuments string,
) error {

	// Get validation from private collection
	validationJSON, err := ctx.GetStub().GetPrivateData(ValidatorPrivateCollection, "VALIDATION_"+validationID)
	if err != nil || validationJSON == nil {
		return fmt.Errorf("validation not found: %v", err)
	}

	var validation ValidationRecord
	err = json.Unmarshal(validationJSON, &validation)
	if err != nil {
		return fmt.Errorf("failed to unmarshal validation: %v", err)
	}

	// Check if already finalized
	if validation.Status == "APPROVED" || validation.Status == "REJECTED" {
		return fmt.Errorf("campaign validation is already final: %s", validation.Status)
	}

	var comments []string
	var issues []string
	if commentsJSON != "" {
		json.Unmarshal([]byte(commentsJSON), &comments)
	}
	if issuesJSON != "" {
		json.Unmarshal([]byte(issuesJSON), &issues)
	}

	timestamp := time.Now().Format(time.RFC3339)

	// Update validation
	validation.Status = status
	validation.DueDiligenceScore = dueDiligenceScore
	validation.RiskScore = riskScore
	validation.RiskLevel = riskLevel
	validation.Comments = comments
	validation.Issues = issues
	validation.RequiredDocuments = requiredDocuments
	validation.ValidatedAt = timestamp

	if status == "APPROVED" {
		validation.DocumentsVerified = true
		validation.ComplianceCheck = true
	}

	validationJSON, _ = json.Marshal(validation)
	err = ctx.GetStub().PutPrivateData(ValidatorPrivateCollection, "VALIDATION_"+validationID, validationJSON)
	if err != nil {
		return fmt.Errorf("failed to update validation: %v", err)
	}

	// Generate digital signature/hash for validation approval
	validationHash := ""
	if status == "APPROVED" {
		hashData := fmt.Sprintf("%s:%s:%s:%f:%f:%s", validationID, campaignID, validation.ValidatorID, dueDiligenceScore, riskScore, timestamp)
		hash := sha256.Sum256([]byte(hashData))
		validationHash = hex.EncodeToString(hash[:])
	}

	// Update validation status in StartupValidatorCollection so startup can see it
	statusUpdate := map[string]interface{}{
		"validationId":      validationID,
		"campaignId":        campaignID,
		"status":            status,
		"riskLevel":         riskLevel,
		"dueDiligenceScore": dueDiligenceScore,
		"riskScore":         riskScore,
		"validationHash":    validationHash,
		"requiredDocuments": requiredDocuments,
		"updatedAt":         timestamp,
	}

	statusJSON, _ := json.Marshal(statusUpdate)
	err = ctx.GetStub().PutPrivateData(StartupValidatorCollection, "VALIDATION_STATUS_"+campaignID, statusJSON)
	if err != nil {
		return fmt.Errorf("failed to update status: %v", err)
	}

	// Store validation approval in ValidatorPlatformCollection for Platform to verify
	if status == "APPROVED" {
		platformVerification := map[string]interface{}{
			"campaignId":        campaignID,
			"validationId":      validationID,
			"validatorId":       validation.ValidatorID,
			"validationHash":    validationHash,
			"dueDiligenceScore": dueDiligenceScore,
			"riskScore":         riskScore,
			"riskLevel":         riskLevel,
			"approvedAt":        timestamp,
		}
		platformJSON, _ := json.Marshal(platformVerification)
		err = ctx.GetStub().PutPrivateData(ValidatorPlatformCollection, "VALIDATION_APPROVAL_"+campaignID, platformJSON)
		if err != nil {
			return fmt.Errorf("failed to store platform verification: %v", err)
		}
	}

	// If blacklisted, record it
	if status == "BLACKLISTED" {
		blacklist := BlacklistedCampaign{
			CampaignID:    campaignID,
			Reason:        "Failed validation - " + requiredDocuments,
			BlacklistedAt: timestamp,
			BlacklistedBy: validation.ValidatorID,
		}
		blacklistJSON, _ := json.Marshal(blacklist)
		ctx.GetStub().PutPrivateData(ValidatorPrivateCollection, "BLACKLISTED_"+campaignID, blacklistJSON)
	}

	return nil
}

// VerifyCampaignHash verifies campaign hash for integrity
func (v *ValidatorContract) VerifyCampaignHash(
	ctx contractapi.TransactionContextInterface,
	campaignID string,
	providedHash string,
) (bool, error) {

	// Get campaign from StartupValidatorCollection
	campaignJSON, err := ctx.GetStub().GetPrivateData(StartupValidatorCollection, "VALIDATION_REQUEST_"+campaignID)
	if err != nil || campaignJSON == nil {
		return false, fmt.Errorf("campaign data not found: %v", err)
	}

	var campaignData map[string]interface{}
	json.Unmarshal(campaignJSON, &campaignData)

	storedHash, ok := campaignData["validationHash"].(string)
	if !ok {
		return false, fmt.Errorf("no validation hash found")
	}

	return storedHash == providedHash, nil
}

// IsCampaignBlacklisted checks if a campaign is blacklisted
func (v *ValidatorContract) IsCampaignBlacklisted(
	ctx contractapi.TransactionContextInterface,
	campaignID string,
) (bool, error) {

	blacklistJSON, err := ctx.GetStub().GetPrivateData(ValidatorPrivateCollection, "BLACKLISTED_"+campaignID)
	if err != nil {
		return false, fmt.Errorf("failed to check blacklist: %v", err)
	}

	return blacklistJSON != nil, nil
}

// ============================================================================
// MILESTONE VERIFICATION - Using PDC
// ============================================================================

// VerifyMilestoneCompletion verifies milestone completion
// Reads from StartupValidatorCollection, stores verification result there
func (v *ValidatorContract) VerifyMilestoneCompletion(
	ctx contractapi.TransactionContextInterface,
	verificationID string,
	milestoneID string,
	campaignID string,
	startupID string,
	milestoneReportHash string,
	deliverablesVerified bool,
	qualityScore float64,
	comments string,
	approved bool,
) error {

	timestamp := time.Now().Format(time.RFC3339)

	verification := MilestoneValidation{
		VerificationID:       verificationID,
		MilestoneID:          milestoneID,
		CampaignID:           campaignID,
		StartupID:            startupID,
		MilestoneReportHash:  milestoneReportHash,
		DeliverablesVerified: deliverablesVerified,
		QualityScore:         qualityScore,
		Comments:             comments,
		Approved:             approved,
		VerifiedAt:           timestamp,
	}

	verificationJSON, err := json.Marshal(verification)
	if err != nil {
		return fmt.Errorf("failed to marshal verification: %v", err)
	}

	// Store in StartupValidatorCollection (shared with startup)
	err = ctx.GetStub().PutPrivateData(StartupValidatorCollection, "MILESTONE_VERIFICATION_"+verificationID, verificationJSON)
	if err != nil {
		return fmt.Errorf("failed to store verification: %v", err)
	}

	return nil
}

// ============================================================================
// RISK ASSESSMENT - Using PDC
// ============================================================================

// AssignRiskScore assigns risk score for a campaign
// Stores in InvestorValidatorCollection for investor access
func (v *ValidatorContract) AssignRiskScore(
	ctx contractapi.TransactionContextInterface,
	insightID string,
	campaignID string,
	investorID string,
	riskScore float64,
	riskLevel string,
	riskFactorsJSON string,
	queryResponse string,
	recommendation string,
) error {

	var riskFactors []string
	if riskFactorsJSON != "" {
		json.Unmarshal([]byte(riskFactorsJSON), &riskFactors)
	}

	timestamp := time.Now().Format(time.RFC3339)

	insight := RiskInsight{
		InsightID:      insightID,
		CampaignID:     campaignID,
		InvestorID:     investorID,
		RiskScore:      riskScore,
		RiskLevel:      riskLevel,
		RiskFactors:    riskFactors,
		QueryResponse:  queryResponse,
		Recommendation: recommendation,
		CreatedAt:      timestamp,
	}

	insightJSON, err := json.Marshal(insight)
	if err != nil {
		return fmt.Errorf("failed to marshal insight: %v", err)
	}

	// Store in InvestorValidatorCollection (shared with investor)
	err = ctx.GetStub().PutPrivateData(InvestorValidatorCollection, "RISK_INSIGHT_"+insightID, insightJSON)
	if err != nil {
		return fmt.Errorf("failed to store risk insight: %v", err)
	}

	// Also respond to risk request if exists
	requestKey := fmt.Sprintf("RISK_REQUEST_%s_%s", campaignID, investorID)
	requestJSON, err := ctx.GetStub().GetPrivateData(InvestorValidatorCollection, requestKey)
	if err == nil && requestJSON != nil {
		var request map[string]interface{}
		json.Unmarshal(requestJSON, &request)
		request["status"] = "FULFILLED"
		request["fulfilledAt"] = timestamp
		requestJSON, _ = json.Marshal(request)
		ctx.GetStub().PutPrivateData(InvestorValidatorCollection, requestKey, requestJSON)
	}

	return nil
}

// ============================================================================
// REPORTING & WITNESSING - Using PDC
// ============================================================================

// SendValidationReportToPlatform sends validation report to Platform
// Stores in ValidatorPlatformCollection
func (v *ValidatorContract) SendValidationReportToPlatform(
	ctx contractapi.TransactionContextInterface,
	reportID string,
	campaignID string,
	validationID string,
	campaignHash string,
	overallScore float64,
	documentScore float64,
	complianceScore float64,
	riskScore float64,
	approved bool,
	reportSummary string,
) error {

	timestamp := time.Now().Format(time.RFC3339)

	// Generate report hash
	reportData := fmt.Sprintf("%s|%s|%.2f|%v", reportID, campaignID, overallScore, approved)
	hash := sha256.Sum256([]byte(reportData))
	reportHash := hex.EncodeToString(hash[:])

	report := ValidationReport{
		ReportID:        reportID,
		CampaignID:      campaignID,
		ValidationID:    validationID,
		CampaignHash:    campaignHash,
		OverallScore:    overallScore,
		DocumentScore:   documentScore,
		ComplianceScore: complianceScore,
		RiskScore:       riskScore,
		Approved:        approved,
		ReportSummary:   reportSummary,
		ReportHash:      reportHash,
		CreatedAt:       timestamp,
	}

	reportJSON, err := json.Marshal(report)
	if err != nil {
		return fmt.Errorf("failed to marshal report: %v", err)
	}

	// Store in ValidatorPlatformCollection
	err = ctx.GetStub().PutPrivateData(ValidatorPlatformCollection, "VALIDATION_REPORT_"+reportID, reportJSON)
	if err != nil {
		return fmt.Errorf("failed to send report: %v", err)
	}

	return nil
}

// WitnessAgreement witnesses an agreement between startup and investor
// Stored on public ledger
func (v *ValidatorContract) WitnessAgreement(
	ctx contractapi.TransactionContextInterface,
	witnessID string,
	agreementID string,
	campaignID string,
	startupID string,
	investorID string,
	investmentAmount float64,
	validatorComments string,
) error {

	timestamp := time.Now().Format(time.RFC3339)

	witness := AgreementWitness{
		WitnessID:         witnessID,
		AgreementID:       agreementID,
		CampaignID:        campaignID,
		StartupID:         startupID,
		InvestorID:        investorID,
		InvestmentAmount:  investmentAmount,
		ValidatorComments: validatorComments,
		WitnessedAt:       timestamp,
	}

	witnessJSON, err := json.Marshal(witness)
	if err != nil {
		return fmt.Errorf("failed to marshal witness: %v", err)
	}

	// Store on public world state
	err = ctx.GetStub().PutState("AGREEMENT_WITNESS_"+witnessID, witnessJSON)
	if err != nil {
		return fmt.Errorf("failed to witness agreement: %v", err)
	}

	return nil
}

// ProvideValidationDetailsToInvestor responds to investor's validation request
// Validator reads request from InvestorValidatorCollection and responds
func (v *ValidatorContract) ProvideValidationDetailsToInvestor(
	ctx contractapi.TransactionContextInterface,
	requestID string,
	campaignID string,
) error {

	// Read investor's request from InvestorValidatorCollection
	requestJSON, err := ctx.GetStub().GetPrivateData(InvestorValidatorCollection, "VALIDATION_REQUEST_"+requestID)
	if err != nil || requestJSON == nil {
		return fmt.Errorf("validation request not found: %v", err)
	}

	// Get validation data for this campaign from ValidatorPrivateCollection
	validationJSON, err := ctx.GetStub().GetPrivateData(ValidatorPlatformCollection, "VALIDATION_APPROVAL_"+campaignID)
	if err != nil || validationJSON == nil {
		return fmt.Errorf("validation approval not found for campaign %s", campaignID)
	}

	var validationApproval map[string]interface{}
	err = json.Unmarshal(validationJSON, &validationApproval)
	if err != nil {
		return fmt.Errorf("failed to unmarshal validation approval: %v", err)
	}

	timestamp := time.Now().Format(time.RFC3339)

	// Create response with validation details
	response := map[string]interface{}{
		"requestId":         requestID,
		"campaignId":        campaignID,
		"validatorId":       validationApproval["validatorId"],
		"dueDiligenceScore": validationApproval["dueDiligenceScore"],
		"riskScore":         validationApproval["riskScore"],
		"riskLevel":         validationApproval["riskLevel"],
		"validationHash":    validationApproval["validationHash"],
		"approvedAt":        validationApproval["approvedAt"],
		"respondedAt":       timestamp,
		"status":            "COMPLETED",
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %v", err)
	}

	// Store response in InvestorValidatorCollection so Investor can read it
	err = ctx.GetStub().PutPrivateData(InvestorValidatorCollection, "VALIDATION_RESPONSE_"+requestID, responseJSON)
	if err != nil {
		return fmt.Errorf("failed to store validation response: %v", err)
	}

	return nil
}

// ConfirmCampaignCompletion confirms campaign completion
func (v *ValidatorContract) ConfirmCampaignCompletion(
	ctx contractapi.TransactionContextInterface,
	campaignID string,
	finalStatus string,
	completionNotes string,
) error {

	timestamp := time.Now().Format(time.RFC3339)

	completion := map[string]interface{}{
		"campaignId":      campaignID,
		"finalStatus":     finalStatus,
		"completionNotes": completionNotes,
		"confirmedAt":     timestamp,
	}

	completionJSON, _ := json.Marshal(completion)

	// Store on public world state
	err := ctx.GetStub().PutState("CAMPAIGN_COMPLETION_"+campaignID, completionJSON)
	if err != nil {
		return fmt.Errorf("failed to confirm completion: %v", err)
	}

	return nil
}

// PublishValidationProof publishes validation proof to public ledger
func (v *ValidatorContract) PublishValidationProof(
	ctx contractapi.TransactionContextInterface,
	proofID string,
	campaignID string,
	validationHash string,
	status string,
) error {

	timestamp := time.Now().Format(time.RFC3339)

	proof := ValidationProof{
		ProofID:        proofID,
		CampaignID:     campaignID,
		ValidationHash: validationHash,
		Status:         status,
		PublishedAt:    timestamp,
	}

	proofJSON, _ := json.Marshal(proof)

	// Store on public world state
	err := ctx.GetStub().PutState("VALIDATION_PROOF_"+proofID, proofJSON)
	if err != nil {
		return fmt.Errorf("failed to publish proof: %v", err)
	}

	return nil
}

// ============================================================================
// DISPUTE INVESTIGATION - Using PDC
// ============================================================================

// AcceptDisputeInvestigation accepts a dispute investigation assignment
func (v *ValidatorContract) AcceptDisputeInvestigation(
	ctx contractapi.TransactionContextInterface,
	investigationID string,
	disputeID string,
	validatorID string,
	initiatorID string,
	initiatorType string,
	respondentID string,
	respondentType string,
	campaignID string,
	disputeType string,
) error {

	timestamp := time.Now().Format(time.RFC3339)

	investigation := DisputeInvestigation{
		InvestigationID:    investigationID,
		DisputeID:          disputeID,
		ValidatorID:        validatorID,
		DisputeType:        disputeType,
		InitiatorID:        initiatorID,
		InitiatorType:      initiatorType,
		RespondentID:       respondentID,
		RespondentType:     respondentType,
		CampaignID:         campaignID,
		Status:             "ASSIGNED",
		Findings:           []InvestigationFinding{},
		EvidenceReviewed:   []string{},
		TransactionLogs:    []string{},
		Recommendation:     "",
		RecommendedPenalty: "",
		AssignedAt:         timestamp,
	}

	investigationJSON, err := json.Marshal(investigation)
	if err != nil {
		return fmt.Errorf("failed to marshal investigation: %v", err)
	}

	// Store in AllOrgsCollection (visible to all for transparency)
	err = ctx.GetStub().PutPrivateData(AllOrgsCollection, "INVESTIGATION_"+investigationID, investigationJSON)
	if err != nil {
		return fmt.Errorf("failed to store investigation: %v", err)
	}

	return nil
}

// RecordInvestigationFinding records an investigation finding
func (v *ValidatorContract) RecordInvestigationFinding(
	ctx contractapi.TransactionContextInterface,
	investigationID string,
	findingID string,
	findingType string,
	description string,
	severity string,
	relatedEvidence string,
) error {

	// Get investigation
	investigationJSON, err := ctx.GetStub().GetPrivateData(AllOrgsCollection, "INVESTIGATION_"+investigationID)
	if err != nil || investigationJSON == nil {
		return fmt.Errorf("investigation not found: %v", err)
	}

	var investigation DisputeInvestigation
	json.Unmarshal(investigationJSON, &investigation)

	timestamp := time.Now().Format(time.RFC3339)

	finding := InvestigationFinding{
		FindingID:       findingID,
		FindingType:     findingType,
		Description:     description,
		Severity:        severity,
		RelatedEvidence: relatedEvidence,
		RecordedAt:      timestamp,
	}

	investigation.Findings = append(investigation.Findings, finding)
	investigation.Status = "IN_PROGRESS"

	investigationJSON, _ = json.Marshal(investigation)
	err = ctx.GetStub().PutPrivateData(AllOrgsCollection, "INVESTIGATION_"+investigationID, investigationJSON)
	if err != nil {
		return fmt.Errorf("failed to record finding: %v", err)
	}

	return nil
}

// CompleteInvestigation completes the investigation with recommendation
func (v *ValidatorContract) CompleteInvestigation(
	ctx contractapi.TransactionContextInterface,
	investigationID string,
	recommendation string,
	recommendedPenalty string,
) error {

	// Get investigation
	investigationJSON, err := ctx.GetStub().GetPrivateData(AllOrgsCollection, "INVESTIGATION_"+investigationID)
	if err != nil || investigationJSON == nil {
		return fmt.Errorf("investigation not found: %v", err)
	}

	var investigation DisputeInvestigation
	json.Unmarshal(investigationJSON, &investigation)

	timestamp := time.Now().Format(time.RFC3339)

	investigation.Status = "COMPLETED"
	investigation.Recommendation = recommendation
	investigation.RecommendedPenalty = recommendedPenalty
	investigation.CompletedAt = timestamp

	investigationJSON, _ = json.Marshal(investigation)
	err = ctx.GetStub().PutPrivateData(AllOrgsCollection, "INVESTIGATION_"+investigationID, investigationJSON)
	if err != nil {
		return fmt.Errorf("failed to complete investigation: %v", err)
	}

	return nil
}

// InvestigateMilestoneDispute investigates milestone-related dispute
func (v *ValidatorContract) InvestigateMilestoneDispute(
	ctx contractapi.TransactionContextInterface,
	investigationID string,
	disputeID string,
	milestoneID string,
	campaignID string,
	validatorID string,
	milestoneReviewed bool,
	deliverableStatus string,
	qualityAssessment float64,
	timelineAssessment string,
	delayJustified bool,
	recommendedAction string,
	comments string,
) error {

	timestamp := time.Now().Format(time.RFC3339)

	milestoneInvestigation := MilestoneInvestigation{
		InvestigationID:    investigationID,
		DisputeID:          disputeID,
		MilestoneID:        milestoneID,
		CampaignID:         campaignID,
		ValidatorID:        validatorID,
		MilestoneReviewed:  milestoneReviewed,
		DeliverableStatus:  deliverableStatus,
		QualityAssessment:  qualityAssessment,
		TimelineAssessment: timelineAssessment,
		DelayJustified:     delayJustified,
		RecommendedAction:  recommendedAction,
		Comments:           comments,
		InvestigatedAt:     timestamp,
	}

	investigationJSON, err := json.Marshal(milestoneInvestigation)
	if err != nil {
		return fmt.Errorf("failed to marshal milestone investigation: %v", err)
	}

	// Store in AllOrgsCollection
	err = ctx.GetStub().PutPrivateData(AllOrgsCollection, "MILESTONE_INVESTIGATION_"+investigationID, investigationJSON)
	if err != nil {
		return fmt.Errorf("failed to store milestone investigation: %v", err)
	}

	return nil
}

// RespondToDispute responds to a dispute when validator is respondent
func (v *ValidatorContract) RespondToDispute(
	ctx contractapi.TransactionContextInterface,
	responseID string,
	disputeID string,
	validatorID string,
	responseText string,
	justification string,
	supportingDocsJSON string,
) error {

	var supportingDocs []string
	if supportingDocsJSON != "" {
		json.Unmarshal([]byte(supportingDocsJSON), &supportingDocs)
	}

	timestamp := time.Now().Format(time.RFC3339)

	response := ValidatorDisputeResponse{
		ResponseID:     responseID,
		DisputeID:      disputeID,
		ValidatorID:    validatorID,
		ResponseText:   responseText,
		Justification:  justification,
		SupportingDocs: supportingDocs,
		RespondedAt:    timestamp,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %v", err)
	}

	// Store in AllOrgsCollection
	err = ctx.GetStub().PutPrivateData(AllOrgsCollection, "VALIDATOR_DISPUTE_RESPONSE_"+responseID, responseJSON)
	if err != nil {
		return fmt.Errorf("failed to respond to dispute: %v", err)
	}

	return nil
}

// ============================================================================
// QUERY FUNCTIONS
// ============================================================================

// GetCampaign retrieves campaign data from StartupValidatorCollection
func (v *ValidatorContract) GetCampaign(ctx contractapi.TransactionContextInterface, campaignID string) (string, error) {
	// 1. Get initial request
	campaignJSON, err := ctx.GetStub().GetPrivateData(StartupValidatorCollection, "VALIDATION_REQUEST_"+campaignID)
	if err != nil || campaignJSON == nil {
		return "", fmt.Errorf("campaign not found: %v", err)
	}

	// 2. See if there is a status update
	statusJSON, err := ctx.GetStub().GetPrivateData(StartupValidatorCollection, "VALIDATION_STATUS_"+campaignID)
	if err != nil {
		// Just return original if error checking status (though typically shouldn't fail)
		return string(campaignJSON), nil
	}

	// If no status update exists, return original
	if statusJSON == nil {
		return string(campaignJSON), nil
	}

	// 3. Merge status into campaign data
	var campaignMap map[string]interface{}
	err = json.Unmarshal(campaignJSON, &campaignMap)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal campaign data: %v", err)
	}

	var statusMap map[string]interface{}
	err = json.Unmarshal(statusJSON, &statusMap)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal status data: %v", err)
	}

	// Overwrite/Add updated fields
	if val, ok := statusMap["status"]; ok {
		campaignMap["validationStatus"] = val
	}
	if val, ok := statusMap["validationHash"]; ok {
		campaignMap["validationHash"] = val
	}
	if val, ok := statusMap["riskLevel"]; ok {
		campaignMap["riskLevel"] = val
	}
	if val, ok := statusMap["dueDiligenceScore"]; ok {
		campaignMap["dueDiligenceScore"] = val
	}
	if val, ok := statusMap["riskScore"]; ok {
		campaignMap["riskScore"] = val
	}
	if val, ok := statusMap["requiredDocuments"]; ok {
		campaignMap["requiredDocuments"] = val
	}

	mergedJSON, err := json.Marshal(campaignMap)
	if err != nil {
		return "", fmt.Errorf("failed to marshal merged data: %v", err)
	}

	return string(mergedJSON), nil
}

// GetPendingValidations retrieves all campaigns pending validation from StartupValidatorCollection
func (v *ValidatorContract) GetPendingValidations(ctx contractapi.TransactionContextInterface) ([]map[string]interface{}, error) {
	// Use GetPrivateDataByRange to iterate through validation requests
	resultsIterator, err := ctx.GetStub().GetPrivateDataByRange(StartupValidatorCollection, "VALIDATION_REQUEST_", "VALIDATION_REQUEST_~")
	if err != nil {
		return nil, fmt.Errorf("failed to get pending validations: %v", err)
	}
	defer resultsIterator.Close()

	var campaigns []map[string]interface{}
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			continue
		}

		var campaignMap map[string]interface{}
		err = json.Unmarshal(queryResponse.Value, &campaignMap)
		if err != nil {
			continue
		}

		// Check status - only include PENDING_VALIDATION
		campaignID, _ := campaignMap["campaignId"].(string)
		statusJSON, _ := ctx.GetStub().GetPrivateData(StartupValidatorCollection, "VALIDATION_STATUS_"+campaignID)

		if statusJSON != nil {
			var statusMap map[string]interface{}
			json.Unmarshal(statusJSON, &statusMap)
			if status, ok := statusMap["status"].(string); ok {
				campaignMap["validationStatus"] = status
				// Skip if already approved/rejected
				if status == "APPROVED" || status == "REJECTED" {
					continue
				}
			}
		} else {
			campaignMap["validationStatus"] = "PENDING_VALIDATION"
		}

		campaigns = append(campaigns, campaignMap)
	}

	if campaigns == nil {
		campaigns = []map[string]interface{}{}
	}

	return campaigns, nil
}

// GetValidation retrieves a validation record
func (v *ValidatorContract) GetValidation(ctx contractapi.TransactionContextInterface, validationID string) (*ValidationRecord, error) {
	validationJSON, err := ctx.GetStub().GetPrivateData(ValidatorPrivateCollection, "VALIDATION_"+validationID)
	if err != nil || validationJSON == nil {
		return nil, fmt.Errorf("validation not found: %v", err)
	}

	var validation ValidationRecord
	err = json.Unmarshal(validationJSON, &validation)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal validation: %v", err)
	}

	return &validation, nil
}

// GetRiskInsight retrieves risk insight for a campaign
func (v *ValidatorContract) GetRiskInsight(ctx contractapi.TransactionContextInterface, campaignID string) (string, error) {
	// Would query from InvestorValidatorCollection
	return `{}`, nil
}

// GetValidationReport retrieves validation report
func (v *ValidatorContract) GetValidationReport(ctx contractapi.TransactionContextInterface, campaignID string) (string, error) {
	// Would query from ValidatorPlatformCollection
	return `{}`, nil
}

// GetMilestoneVerification retrieves milestone verification
func (v *ValidatorContract) GetMilestoneVerification(ctx contractapi.TransactionContextInterface, verificationID string) (string, error) {
	verificationJSON, err := ctx.GetStub().GetPrivateData(StartupValidatorCollection, "MILESTONE_VERIFICATION_"+verificationID)
	if err != nil || verificationJSON == nil {
		return "", fmt.Errorf("verification not found: %v", err)
	}

	return string(verificationJSON), nil
}

// GetMilestoneVerificationByMilestone retrieves verification by milestone ID
func (v *ValidatorContract) GetMilestoneVerificationByMilestone(ctx contractapi.TransactionContextInterface, milestoneID string) (string, error) {
	// Would use rich query
	return `{}`, nil
}

// GetAgreementWitness retrieves agreement witness record
func (v *ValidatorContract) GetAgreementWitness(ctx contractapi.TransactionContextInterface, agreementID string) (string, error) {
	witnessJSON, err := ctx.GetStub().GetState("AGREEMENT_WITNESS_" + agreementID)
	if err != nil || witnessJSON == nil {
		return "", fmt.Errorf("witness not found: %v", err)
	}

	return string(witnessJSON), nil
}

// GetCampaignCompletion retrieves campaign completion record
func (v *ValidatorContract) GetCampaignCompletion(ctx contractapi.TransactionContextInterface, campaignID string) (string, error) {
	completionJSON, err := ctx.GetStub().GetState("CAMPAIGN_COMPLETION_" + campaignID)
	if err != nil || completionJSON == nil {
		return "", fmt.Errorf("completion record not found: %v", err)
	}

	return string(completionJSON), nil
}

// GetValidationsByCampaign retrieves all validations for a campaign
func (v *ValidatorContract) GetValidationsByCampaign(ctx contractapi.TransactionContextInterface, campaignID string) (string, error) {
	// Would use rich query
	return `[]`, nil
}

// GetInvestigation retrieves an investigation
func (v *ValidatorContract) GetInvestigation(ctx contractapi.TransactionContextInterface, investigationID string) (string, error) {
	investigationJSON, err := ctx.GetStub().GetPrivateData(AllOrgsCollection, "INVESTIGATION_"+investigationID)
	if err != nil || investigationJSON == nil {
		return "", fmt.Errorf("investigation not found: %v", err)
	}

	return string(investigationJSON), nil
}

// GetValidatorDisputes retrieves all disputes for a validator
func (v *ValidatorContract) GetValidatorDisputes(ctx contractapi.TransactionContextInterface, validatorID string) (string, error) {
	// Would use rich query
	return `[]`, nil
}

// ============================================================================

// ============================================================================
// TOKEN INTEGRATION FUNCTIONS (CFT/CFRT)
// ============================================================================

// ValidatorFeeRecord represents fee record for validators
type ValidatorFeeRecord struct {
	RecordID        string  `json:"recordId"`
	ValidatorID     string  `json:"validatorId"`
	CampaignID      string  `json:"campaignId"`
	FeeType         string  `json:"feeType"`
	AmountCFT       float64 `json:"amountCFT"`
	TransactionHash string  `json:"transactionHash"`
	PaidAt          string  `json:"paidAt"`
}

// GetValidationFees returns the fees for validation services (paid by startups)
func (v *ValidatorContract) GetValidationFees(
	ctx contractapi.TransactionContextInterface,
) (string, error) {

	// Fee schedule based on 1 INR = 2.5 CFT
	fees := map[string]interface{}{
		"campaignValidationFee": map[string]interface{}{
			"amountCFT":   500,
			"amountINR":   200,
			"description": "Standard campaign validation fee (paid by startup)",
		},
		"milestoneVerificationFee": map[string]interface{}{
			"amountCFT":   125,
			"amountINR":   50,
			"description": "Milestone verification fee",
		},
		"expeditedValidationFee": map[string]interface{}{
			"amountCFT":   1000,
			"amountINR":   400,
			"description": "Expedited validation (24-48 hours)",
		},
	}

	feesJSON, err := json.Marshal(fees)
	if err != nil {
		return "", fmt.Errorf("failed to marshal fees: %v", err)
	}

	return string(feesJSON), nil
}

// RecordValidationFeeReceived records fee received for validation work
func (v *ValidatorContract) RecordValidationFeeReceived(
	ctx contractapi.TransactionContextInterface,
	recordID string,
	validatorID string,
	campaignID string,
	feeType string,
	amountCFT float64,
	transactionHash string,
) error {

	timestamp := time.Now().Format(time.RFC3339)

	record := ValidatorFeeRecord{
		RecordID:        recordID,
		ValidatorID:     validatorID,
		CampaignID:      campaignID,
		FeeType:         feeType,
		AmountCFT:       amountCFT,
		TransactionHash: transactionHash,
		PaidAt:          timestamp,
	}

	recordJSON, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal fee record: %v", err)
	}

	err = ctx.GetStub().PutPrivateData(ValidatorPrivateCollection, "FEE_RECEIVED_"+recordID, recordJSON)
	if err != nil {
		return fmt.Errorf("failed to record fee: %v", err)
	}

	return nil
}

// GetDisputePenaltySchedule returns penalty amounts for dispute scenarios
func (v *ValidatorContract) GetDisputePenaltySchedule(
	ctx contractapi.TransactionContextInterface,
) (string, error) {

	// Penalty schedule based on 1 INR = 2.5 CFT
	penalties := map[string]interface{}{
		"validatorFraudApproval": map[string]interface{}{
			"penaltyCFT":  1250,
			"penaltyINR":  500,
			"description": "Approving fraudulent campaign",
		},
		"biasedValidation": map[string]interface{}{
			"penaltyCFT":  625,
			"penaltyINR":  250,
			"description": "Biased validation (evidence of favoritism)",
		},
		"delayedValidation": map[string]interface{}{
			"penaltyCFT":  250,
			"penaltyINR":  100,
			"description": "Validation delay beyond SLA",
		},
	}

	penaltiesJSON, err := json.Marshal(penalties)
	if err != nil {
		return "", fmt.Errorf("failed to marshal penalties: %v", err)
	}

	return string(penaltiesJSON), nil
}

