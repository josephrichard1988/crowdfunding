package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// ValidatorContract provides functions for ValidatorOrg operations
type ValidatorContract struct {
	contractapi.Contract
}

// ValidationRecord represents a campaign validation result
type ValidationRecord struct {
	ValidationID       string   `json:"validationId"`
	CampaignID         string   `json:"campaignId"`
	CampaignHash       string   `json:"campaignHash"` // Hash from StartupOrg for verification
	ValidatorID        string   `json:"validatorId"`
	// Status: PENDING, IN_PROGRESS, APPROVED, ON_HOLD, REJECTED, BLACKLISTED
	Status             string   `json:"status"`
	DocumentsVerified  bool     `json:"documentsVerified"`
	ComplianceCheck    bool     `json:"complianceCheck"`
	DueDiligenceScore  float64  `json:"dueDiligenceScore"`
	RiskScore          float64  `json:"riskScore"`
	RiskLevel          string   `json:"riskLevel"` // LOW, MEDIUM, HIGH
	Comments           []string `json:"comments"`
	Issues             []string `json:"issues"`
	RequiredDocuments  string   `json:"requiredDocuments"` // Docs needed if ON_HOLD
	ValidationAttempts []ValidationAttempt `json:"validationAttempts"` // Track all attempts
	ValidatedAt        string   `json:"validatedAt"`
	CreatedAt          string   `json:"createdAt"`
}

// ValidationAttempt tracks each validation attempt (linked by CampaignID)
type ValidationAttempt struct {
	AttemptID       string   `json:"attemptId"`
	AttemptNumber   int      `json:"attemptNumber"`
	DocumentsReviewed []string `json:"documentsReviewed"`
	Status          string   `json:"status"` // APPROVED, ON_HOLD, REJECTED
	Score           float64  `json:"score"`
	Comments        string   `json:"comments"`
	RequiredDocs    string   `json:"requiredDocs"`
	AttemptedAt     string   `json:"attemptedAt"`
}

// RiskInsight represents risk information shared with investors
type RiskInsight struct {
	InsightID      string   `json:"insightId"`
	CampaignID     string   `json:"campaignId"`
	InvestorID     string   `json:"investorId"` // If requested by specific investor
	RiskScore      float64  `json:"riskScore"`
	RiskLevel      string   `json:"riskLevel"`
	RiskFactors    []string `json:"riskFactors"`
	QueryResponse  string   `json:"queryResponse"` // Response to investor's query
	Recommendation string   `json:"recommendation"`
	CreatedAt      string   `json:"createdAt"`
}

// ValidationReport represents detailed report sent to PlatformOrg
type ValidationReport struct {
	ReportID        string  `json:"reportId"`
	CampaignID      string  `json:"campaignId"`
	ValidationID    string  `json:"validationId"`
	CampaignHash    string  `json:"campaignHash"` // For Platform to verify
	OverallScore    float64 `json:"overallScore"`
	DocumentScore   float64 `json:"documentScore"`
	ComplianceScore float64 `json:"complianceScore"`
	RiskScore       float64 `json:"riskScore"`
	Approved        bool    `json:"approved"`
	ReportSummary   string  `json:"reportSummary"`
	ReportHash      string  `json:"reportHash"`
	CreatedAt       string  `json:"createdAt"`
}

// ValidationProof for common-channel (privacy-preserving)
type ValidationProof struct {
	ProofID        string `json:"proofId"`
	CampaignID     string `json:"campaignId"`
	ValidationHash string `json:"validationHash"`
	Status         string `json:"status"`
	PublishedAt    string `json:"publishedAt"`
}

// BlacklistedCampaign tracks rejected campaigns that cannot be resubmitted
type BlacklistedCampaign struct {
	CampaignID     string `json:"campaignId"`
	Reason         string `json:"reason"`
	BlacklistedAt  string `json:"blacklistedAt"`
	BlacklistedBy  string `json:"blacklistedBy"`
}

// MilestoneValidation represents milestone verification by Validator
// Used in Phase 12: startup-validator-channel
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
// Used in Phase 9: common-channel
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

// ============================================================================
// DISPUTE INVESTIGATION STRUCTURES FOR VALIDATOR
// ============================================================================

// DisputeInvestigation represents validator's investigation of a dispute
type DisputeInvestigation struct {
	InvestigationID   string                `json:"investigationId"`
	DisputeID         string                `json:"disputeId"`
	ValidatorID       string                `json:"validatorId"`
	DisputeType       string                `json:"disputeType"`
	InitiatorID       string                `json:"initiatorId"`
	InitiatorType     string                `json:"initiatorType"`
	RespondentID      string                `json:"respondentId"`
	RespondentType    string                `json:"respondentType"`
	CampaignID        string                `json:"campaignId"`
	Status            string                `json:"status"` // ASSIGNED, IN_PROGRESS, COMPLETED
	Findings          []InvestigationFinding `json:"findings"`
	EvidenceReviewed  []string              `json:"evidenceReviewed"`
	TransactionLogs   []string              `json:"transactionLogs"`
	Recommendation    string                `json:"recommendation"` // FAVOR_INITIATOR, FAVOR_RESPONDENT, PARTIAL, DISMISS
	RecommendedPenalty string               `json:"recommendedPenalty"`
	AssignedAt        string                `json:"assignedAt"`
	CompletedAt       string                `json:"completedAt"`
}

// InvestigationFinding represents a finding during investigation
type InvestigationFinding struct {
	FindingID     string `json:"findingId"`
	FindingType   string `json:"findingType"` // EVIDENCE_VERIFIED, EVIDENCE_INVALID, RULE_VIOLATION, POLICY_BREACH, FRAUD_DETECTED
	Description   string `json:"description"`
	Severity      string `json:"severity"` // LOW, MEDIUM, HIGH, CRITICAL
	RelatedEvidence string `json:"relatedEvidence"` // IPFS hash
	RecordedAt    string `json:"recordedAt"`
}

// ValidatorDisputeResponse represents validator's response when they are respondent
type ValidatorDisputeResponse struct {
	ResponseID      string   `json:"responseId"`
	DisputeID       string   `json:"disputeId"`
	ValidatorID     string   `json:"validatorId"`
	ResponseText    string   `json:"responseText"`
	Justification   string   `json:"justification"`
	SupportingDocs  []string `json:"supportingDocs"` // IPFS hashes
	RespondedAt     string   `json:"respondedAt"`
}

// MilestoneInvestigation for investigating milestone-related disputes
type MilestoneInvestigation struct {
	InvestigationID     string  `json:"investigationId"`
	DisputeID           string  `json:"disputeId"`
	MilestoneID         string  `json:"milestoneId"`
	CampaignID          string  `json:"campaignId"`
	ValidatorID         string  `json:"validatorId"`
	MilestoneReviewed   bool    `json:"milestoneReviewed"`
	DeliverableStatus   string  `json:"deliverableStatus"` // COMPLETED, PARTIAL, NOT_DELIVERED
	QualityAssessment   float64 `json:"qualityAssessment"` // 0-100
	TimelineAssessment  string  `json:"timelineAssessment"` // ON_TIME, DELAYED, SEVERELY_DELAYED
	DelayJustified      bool    `json:"delayJustified"`
	RecommendedAction   string  `json:"recommendedAction"` // RELEASE_FUNDS, PARTIAL_REFUND, FULL_REFUND
	Comments            string  `json:"comments"`
	InvestigatedAt      string  `json:"investigatedAt"`
}

// InitLedger initializes the ValidatorOrg ledger
func (v *ValidatorContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	fmt.Println("ValidatorOrg contract initialized - Merged Version")
	return nil
}

// ============================================================================
// STARTUP-VALIDATOR-CHANNEL FUNCTIONS
// Endorsed by: StartupOrg, ValidatorOrg
// ============================================================================

// ValidateCampaign performs validation checks on submitted campaign (ML Model)
// Step 3: Validator validates campaign - can APPROVE, ON_HOLD, or REJECT
// Channel: startup-validator-channel
// Endorsers: StartupOrg, ValidatorOrg
func (v *ValidatorContract) ValidateCampaign(
	ctx contractapi.TransactionContextInterface,
	validationID string,
	campaignID string,
	campaignHash string, // Hash from StartupOrg for verification
	validatorID string,
	documentsVerified bool,
	complianceCheck bool,
	dueDiligenceScore float64,
	riskScore float64,
	decision string, // APPROVED, ON_HOLD, REJECTED
	commentsJSON string,
	requiredDocuments string, // If ON_HOLD, what docs are needed
) (string, error) {
	// Check if campaign is already blacklisted
	blacklistKey := fmt.Sprintf("BLACKLIST_%s", campaignID)
	blacklisted, _ := ctx.GetStub().GetState(blacklistKey)
	if blacklisted != nil {
		return "", fmt.Errorf("campaign %s is blacklisted and cannot be validated", campaignID)
	}

	// Parse comments
	var comments []string
	if commentsJSON != "" {
		if err := json.Unmarshal([]byte(commentsJSON), &comments); err != nil {
			return "", fmt.Errorf("failed to parse comments: %v", err)
		}
	}

	// Determine risk level based on score
	var riskLevel string
	if riskScore < 3.0 {
		riskLevel = "LOW"
	} else if riskScore < 7.0 {
		riskLevel = "MEDIUM"
	} else {
		riskLevel = "HIGH"
	}

	now := time.Now().Format(time.RFC3339)

	// Check if validation record exists (for revalidation after ON_HOLD)
	existingJSON, _ := ctx.GetStub().GetState(validationID)
	var validation ValidationRecord
	var attemptNumber int

	if existingJSON != nil {
		// Existing validation - this is a revalidation
		err := json.Unmarshal(existingJSON, &validation)
		if err != nil {
			return "", err
		}
		attemptNumber = len(validation.ValidationAttempts) + 1
	} else {
		// New validation
		validation = ValidationRecord{
			ValidationID:       validationID,
			CampaignID:         campaignID,
			CampaignHash:       campaignHash,
			ValidatorID:        validatorID,
			DocumentsVerified:  documentsVerified,
			ComplianceCheck:    complianceCheck,
			DueDiligenceScore:  dueDiligenceScore,
			RiskScore:          riskScore,
			RiskLevel:          riskLevel,
			Comments:           comments,
			ValidationAttempts: []ValidationAttempt{},
			CreatedAt:          now,
		}
		attemptNumber = 1
	}

	// Create validation attempt
	attempt := ValidationAttempt{
		AttemptID:     fmt.Sprintf("ATT_%s_%d", validationID, attemptNumber),
		AttemptNumber: attemptNumber,
		Status:        decision,
		Score:         dueDiligenceScore,
		Comments:      commentsJSON,
		RequiredDocs:  requiredDocuments,
		AttemptedAt:   now,
	}
	validation.ValidationAttempts = append(validation.ValidationAttempts, attempt)

	// Update validation based on decision
	validation.Status = decision
	validation.ValidatedAt = now
	validation.RiskScore = riskScore
	validation.RiskLevel = riskLevel
	validation.DueDiligenceScore = dueDiligenceScore

	if decision == "ON_HOLD" {
		validation.RequiredDocuments = requiredDocuments
	}

	// If REJECTED due to fraud, blacklist the campaign
	if decision == "REJECTED" {
		// Check if fraud detected (high risk + documents not verified)
		if !documentsVerified && riskScore >= 8.0 {
			validation.Status = "BLACKLISTED"
			// Create blacklist entry
			blacklistEntry := BlacklistedCampaign{
				CampaignID:    campaignID,
				Reason:        "Fraudulent documents detected",
				BlacklistedAt: now,
				BlacklistedBy: validatorID,
			}
			blacklistJSON, _ := json.Marshal(blacklistEntry)
			ctx.GetStub().PutState(blacklistKey, blacklistJSON)
		}
	}

	validationJSON, err := json.Marshal(validation)
	if err != nil {
		return "", err
	}

	// Store validation record
	err = ctx.GetStub().PutState(validationID, validationJSON)
	if err != nil {
		return "", err
	}

	// Store by campaign ID for easy lookup
	campaignValKey := fmt.Sprintf("CAMPAIGN_VAL_%s", campaignID)
	err = ctx.GetStub().PutState(campaignValKey, validationJSON)
	if err != nil {
		return "", err
	}

	// Emit event
	eventPayload := map[string]interface{}{
		"validationId":      validationID,
		"campaignId":        campaignID,
		"decision":          decision,
		"attemptNumber":     attemptNumber,
		"riskLevel":         riskLevel,
		"requiredDocuments": requiredDocuments,
		"channel":           "startup-validator-channel",
		"action":            "CAMPAIGN_VALIDATED",
		"timestamp":         now,
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("CampaignValidated", eventJSON)

	response := map[string]interface{}{
		"message":           fmt.Sprintf("Campaign validation completed: %s", decision),
		"validationId":      validationID,
		"campaignId":        campaignID,
		"status":            validation.Status,
		"attemptNumber":     attemptNumber,
		"riskLevel":         riskLevel,
		"requiredDocuments": requiredDocuments,
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// ApproveOrRejectCampaign makes final decision on campaign approval
// This is a separate function for final approval after validation
// Channel: startup-validator-channel
// Endorsers: StartupOrg, ValidatorOrg
func (v *ValidatorContract) ApproveOrRejectCampaign(
	ctx contractapi.TransactionContextInterface,
	validationID string,
	decision string, // APPROVED, REJECTED, ON_HOLD
	finalComments string,
	requiredDocuments string, // If ON_HOLD
) (string, error) {
	validationJSON, err := ctx.GetStub().GetState(validationID)
	if err != nil {
		return "", fmt.Errorf("failed to read validation: %v", err)
	}
	if validationJSON == nil {
		return "", fmt.Errorf("validation %s does not exist", validationID)
	}

	var validation ValidationRecord
	err = json.Unmarshal(validationJSON, &validation)
	if err != nil {
		return "", err
	}

	now := time.Now().Format(time.RFC3339)

	// Update validation status
	validation.Status = decision
	validation.ValidatedAt = now
	validation.Comments = append(validation.Comments, finalComments)

	if decision == "ON_HOLD" {
		validation.RequiredDocuments = requiredDocuments
	}

	// If REJECTED, blacklist the campaign
	if decision == "REJECTED" {
		blacklistKey := fmt.Sprintf("BLACKLIST_%s", validation.CampaignID)
		blacklistEntry := BlacklistedCampaign{
			CampaignID:    validation.CampaignID,
			Reason:        finalComments,
			BlacklistedAt: now,
			BlacklistedBy: validation.ValidatorID,
		}
		blacklistJSON, _ := json.Marshal(blacklistEntry)
		ctx.GetStub().PutState(blacklistKey, blacklistJSON)
		validation.Status = "BLACKLISTED"
	}

	updatedValidationJSON, err := json.Marshal(validation)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(validationID, updatedValidationJSON)
	if err != nil {
		return "", err
	}

	// Update campaign key as well
	campaignValKey := fmt.Sprintf("CAMPAIGN_VAL_%s", validation.CampaignID)
	ctx.GetStub().PutState(campaignValKey, updatedValidationJSON)

	// Emit event
	eventPayload := map[string]interface{}{
		"validationId":      validationID,
		"campaignId":        validation.CampaignID,
		"decision":          decision,
		"status":            validation.Status,
		"requiredDocuments": requiredDocuments,
		"channel":           "startup-validator-channel",
		"action":            "VALIDATION_DECISION",
		"timestamp":         now,
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("CampaignDecision", eventJSON)

	response := map[string]interface{}{
		"message":           fmt.Sprintf("Campaign validation decision: %s", validation.Status),
		"validationId":      validationID,
		"campaignId":        validation.CampaignID,
		"status":            validation.Status,
		"requiredDocuments": requiredDocuments,
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// VerifyCampaignHash allows Platform/Investor to verify campaign hash
// Used by PlatformOrg and InvestorOrg to verify campaign validity
// Channel: validator-platform-channel or investor-validator-channel
func (v *ValidatorContract) VerifyCampaignHash(
	ctx contractapi.TransactionContextInterface,
	campaignID string,
	hashToVerify string,
) (string, error) {
	campaignValKey := fmt.Sprintf("CAMPAIGN_VAL_%s", campaignID)
	validationJSON, err := ctx.GetStub().GetState(campaignValKey)
	if err != nil {
		return "", fmt.Errorf("failed to read validation: %v", err)
	}
	if validationJSON == nil {
		return "", fmt.Errorf("no validation record for campaign %s", campaignID)
	}

	var validation ValidationRecord
	err = json.Unmarshal(validationJSON, &validation)
	if err != nil {
		return "", err
	}

	isValid := validation.CampaignHash == hashToVerify

	response := map[string]interface{}{
		"campaignId":   campaignID,
		"hashValid":    isValid,
		"storedHash":   validation.CampaignHash,
		"providedHash": hashToVerify,
		"status":       validation.Status,
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// IsCampaignBlacklisted checks if a campaign is blacklisted
func (v *ValidatorContract) IsCampaignBlacklisted(
	ctx contractapi.TransactionContextInterface,
	campaignID string,
) (string, error) {
	blacklistKey := fmt.Sprintf("BLACKLIST_%s", campaignID)
	blacklistJSON, err := ctx.GetStub().GetState(blacklistKey)
	if err != nil {
		return "", fmt.Errorf("failed to check blacklist: %v", err)
	}

	if blacklistJSON == nil {
		response := map[string]interface{}{
			"campaignId":   campaignID,
			"blacklisted":  false,
		}
		responseJSON, _ := json.Marshal(response)
		return string(responseJSON), nil
	}

	var blacklist BlacklistedCampaign
	json.Unmarshal(blacklistJSON, &blacklist)

	response := map[string]interface{}{
		"campaignId":    campaignID,
		"blacklisted":   true,
		"reason":        blacklist.Reason,
		"blacklistedAt": blacklist.BlacklistedAt,
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// ============================================================================
// MILESTONE VERIFICATION ON STARTUP-VALIDATOR-CHANNEL
// Phase 12: Validator verifies milestone completion
// ============================================================================

// VerifyMilestoneCompletion verifies milestone completion report from Startup
// Step 12: Validator verifies milestone for fund release
// Channel: startup-validator-channel
// Endorsers: StartupOrg, ValidatorOrg
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
) (string, error) {
	now := time.Now().Format(time.RFC3339)

	// Create milestone verification record
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
		VerifiedAt:           now,
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

	// Store by milestone for lookup
	milestoneVerifyKey := fmt.Sprintf("MILESTONE_VERIFY_%s", milestoneID)
	ctx.GetStub().PutState(milestoneVerifyKey, verificationJSON)

	// Emit event
	eventPayload := map[string]interface{}{
		"verificationId":       verificationID,
		"milestoneId":          milestoneID,
		"campaignId":           campaignID,
		"startupId":            startupID,
		"deliverablesVerified": deliverablesVerified,
		"qualityScore":         qualityScore,
		"approved":             approved,
		"channel":              "startup-validator-channel",
		"action":               "MILESTONE_VERIFIED",
		"timestamp":            now,
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("MilestoneVerified", eventJSON)

	status := "REJECTED"
	if approved {
		status = "APPROVED"
	}

	response := map[string]interface{}{
		"message":        "Milestone verification completed",
		"verificationId": verificationID,
		"milestoneId":    milestoneID,
		"approved":       approved,
		"status":         status,
		"qualityScore":   qualityScore,
		"nextStep":       "Platform to release funds from escrow on common-channel",
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// ============================================================================
// INVESTOR-VALIDATOR-CHANNEL FUNCTIONS
// Endorsed by: InvestorOrg, ValidatorOrg
// ============================================================================

// AssignRiskScore assigns and shares risk score with InvestorOrg
// Step 7: Validator responds to investor's risk query
// Channel: investor-validator-channel
// Endorsers: InvestorOrg, ValidatorOrg
func (v *ValidatorContract) AssignRiskScore(
	ctx contractapi.TransactionContextInterface,
	insightID string,
	campaignID string,
	investorID string, // Requesting investor (optional for general risk)
	riskScore float64,
	riskFactorsJSON string,
	investorQuery string, // Query from investor
	queryResponse string, // Response to query
	recommendation string,
) (string, error) {
	// Parse risk factors
	var riskFactors []string
	if riskFactorsJSON != "" {
		if err := json.Unmarshal([]byte(riskFactorsJSON), &riskFactors); err != nil {
			return "", fmt.Errorf("failed to parse risk factors: %v", err)
		}
	}

	// Determine risk level based on score
	var riskLevel string
	if riskScore < 3.0 {
		riskLevel = "LOW"
	} else if riskScore < 7.0 {
		riskLevel = "MEDIUM"
	} else {
		riskLevel = "HIGH"
	}

	// Create risk insight for investors
	insight := RiskInsight{
		InsightID:      insightID,
		CampaignID:     campaignID,
		InvestorID:     investorID,
		RiskScore:      riskScore,
		RiskLevel:      riskLevel,
		RiskFactors:    riskFactors,
		QueryResponse:  queryResponse,
		Recommendation: recommendation,
		CreatedAt:      time.Now().Format(time.RFC3339),
	}

	insightJSON, err := json.Marshal(insight)
	if err != nil {
		return "", err
	}

	// Store on investor-validator-channel
	err = ctx.GetStub().PutState(insightID, insightJSON)
	if err != nil {
		return "", err
	}

	// Also store by campaign ID for easy lookup
	campaignRiskKey := fmt.Sprintf("RISK_%s", campaignID)
	err = ctx.GetStub().PutState(campaignRiskKey, insightJSON)
	if err != nil {
		return "", err
	}

	// If investor-specific, store by investor too
	if investorID != "" {
		investorRiskKey := fmt.Sprintf("INVESTOR_RISK_%s_%s", investorID, campaignID)
		ctx.GetStub().PutState(investorRiskKey, insightJSON)
	}

	// Emit event
	eventPayload := map[string]interface{}{
		"insightId":  insightID,
		"campaignId": campaignID,
		"investorId": investorID,
		"riskScore":  riskScore,
		"riskLevel":  riskLevel,
		"channel":    "investor-validator-channel",
		"action":     "RISK_SCORE_ASSIGNED",
		"timestamp":  insight.CreatedAt,
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("RiskScoreAssigned", eventJSON)

	response := map[string]interface{}{
		"message":       "Risk analysis provided to investor",
		"insightId":     insightID,
		"campaignId":    campaignID,
		"investorId":    investorID,
		"riskScore":     riskScore,
		"riskLevel":     riskLevel,
		"queryResponse": queryResponse,
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// ============================================================================
// VALIDATOR-PLATFORM-CHANNEL FUNCTIONS
// Endorsed by: ValidatorOrg, PlatformOrg
// ============================================================================

// SendValidationReportToPlatform sends detailed validation report to PlatformOrg
// Step 3: Validator sends report to Platform for verification
// Channel: validator-platform-channel
// Endorsers: ValidatorOrg, PlatformOrg
func (v *ValidatorContract) SendValidationReportToPlatform(
	ctx contractapi.TransactionContextInterface,
	reportID string,
	campaignID string,
	validationID string,
	campaignHash string, // Hash for Platform to verify with StartupOrg
	overallScore float64,
	documentScore float64,
	complianceScore float64,
	riskScore float64,
	approved bool,
	reportSummary string,
) (string, error) {
	// Generate report hash
	reportData := map[string]interface{}{
		"reportId":        reportID,
		"campaignId":      campaignID,
		"campaignHash":    campaignHash,
		"overallScore":    overallScore,
		"documentScore":   documentScore,
		"complianceScore": complianceScore,
		"riskScore":       riskScore,
		"approved":        approved,
		"timestamp":       time.Now().Format(time.RFC3339),
	}
	reportDataJSON, _ := json.Marshal(reportData)
	reportHash := generateHash(string(reportDataJSON))

	// Create validation report
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
		CreatedAt:       time.Now().Format(time.RFC3339),
	}

	reportJSON, err := json.Marshal(report)
	if err != nil {
		return "", err
	}

	// Store on validator-platform-channel
	err = ctx.GetStub().PutState(reportID, reportJSON)
	if err != nil {
		return "", err
	}

	// Also store by campaign ID for platform lookup
	platformReportKey := fmt.Sprintf("PLATFORM_REPORT_%s", campaignID)
	err = ctx.GetStub().PutState(platformReportKey, reportJSON)
	if err != nil {
		return "", err
	}

	// Emit event
	eventPayload := map[string]interface{}{
		"reportId":     reportID,
		"campaignId":   campaignID,
		"campaignHash": campaignHash,
		"overallScore": overallScore,
		"approved":     approved,
		"reportHash":   reportHash,
		"channel":      "validator-platform-channel",
		"action":       "VALIDATION_REPORT_SENT",
		"timestamp":    report.CreatedAt,
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("ValidationReportSent", eventJSON)

	response := map[string]interface{}{
		"message":      "Validation report sent to PlatformOrg",
		"reportId":     reportID,
		"campaignId":   campaignID,
		"campaignHash": campaignHash,
		"overallScore": overallScore,
		"approved":     approved,
		"reportHash":   reportHash,
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// ============================================================================
// COMMON-CHANNEL FUNCTIONS
// Multi-party visibility for agreements, fund releases, and proofs
// ============================================================================

// WitnessAgreement records Validator as witness to startup-investor agreement
// Step 9: Validator witnesses the agreement for multi-party visibility
// Channel: common-channel
// Endorsers: ValidatorOrg (multi-party visibility)
func (v *ValidatorContract) WitnessAgreement(
	ctx contractapi.TransactionContextInterface,
	witnessID string,
	agreementID string,
	campaignID string,
	startupID string,
	investorID string,
	investmentAmount float64,
	validatorComments string,
) (string, error) {
	now := time.Now().Format(time.RFC3339)

	// Create witness record
	witness := AgreementWitness{
		WitnessID:         witnessID,
		AgreementID:       agreementID,
		CampaignID:        campaignID,
		StartupID:         startupID,
		InvestorID:        investorID,
		InvestmentAmount:  investmentAmount,
		ValidatorComments: validatorComments,
		WitnessedAt:       now,
	}

	witnessJSON, err := json.Marshal(witness)
	if err != nil {
		return "", err
	}

	// Store witness record
	err = ctx.GetStub().PutState(witnessID, witnessJSON)
	if err != nil {
		return "", err
	}

	// Store by agreement for lookup
	agreementWitnessKey := fmt.Sprintf("VALIDATOR_WITNESS_%s", agreementID)
	ctx.GetStub().PutState(agreementWitnessKey, witnessJSON)

	// Emit event
	eventPayload := map[string]interface{}{
		"witnessId":        witnessID,
		"agreementId":      agreementID,
		"campaignId":       campaignID,
		"startupId":        startupID,
		"investorId":       investorID,
		"investmentAmount": investmentAmount,
		"channel":          "common-channel",
		"action":           "VALIDATOR_WITNESSED_AGREEMENT",
		"timestamp":        now,
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("ValidatorWitnessedAgreement", eventJSON)

	response := map[string]interface{}{
		"message":          "Agreement witnessed by Validator",
		"witnessId":        witnessID,
		"agreementId":      agreementID,
		"investmentAmount": investmentAmount,
		"witnessedAt":      now,
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// ConfirmCampaignCompletion confirms campaign has successfully completed
// Step 14: Validator confirms campaign completion for multi-party visibility
// Channel: common-channel
// Endorsers: ValidatorOrg (multi-party visibility)
func (v *ValidatorContract) ConfirmCampaignCompletion(
	ctx contractapi.TransactionContextInterface,
	confirmationID string,
	campaignID string,
	validationID string,
	allMilestonesCompleted bool,
	finalReport string,
) (string, error) {
	now := time.Now().Format(time.RFC3339)

	// Create completion confirmation
	confirmation := map[string]interface{}{
		"confirmationId":         confirmationID,
		"campaignId":             campaignID,
		"validationId":           validationID,
		"allMilestonesCompleted": allMilestonesCompleted,
		"finalReport":            finalReport,
		"confirmedAt":            now,
	}

	confirmationJSON, err := json.Marshal(confirmation)
	if err != nil {
		return "", err
	}

	// Store confirmation
	err = ctx.GetStub().PutState(confirmationID, confirmationJSON)
	if err != nil {
		return "", err
	}

	// Store by campaign for lookup
	completionKey := fmt.Sprintf("CAMPAIGN_COMPLETION_%s", campaignID)
	ctx.GetStub().PutState(completionKey, confirmationJSON)

	// Emit event
	eventPayload := map[string]interface{}{
		"confirmationId":         confirmationID,
		"campaignId":             campaignID,
		"allMilestonesCompleted": allMilestonesCompleted,
		"channel":                "common-channel",
		"action":                 "CAMPAIGN_COMPLETION_CONFIRMED",
		"timestamp":              now,
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("CampaignCompletionConfirmed", eventJSON)

	response := map[string]interface{}{
		"message":                "Campaign completion confirmed by Validator",
		"confirmationId":         confirmationID,
		"campaignId":             campaignID,
		"allMilestonesCompleted": allMilestonesCompleted,
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// PublishValidationProof publishes cryptographic validation proof to common-channel
// Channel: common-channel
// Purpose: Privacy-preserving validation proof (no confidential details)
func (v *ValidatorContract) PublishValidationProof(
	ctx contractapi.TransactionContextInterface,
	proofID string,
	campaignID string,
	validationID string,
	status string,
) (string, error) {
	// Generate validation hash (no sensitive data)
	proofData := map[string]interface{}{
		"proofId":      proofID,
		"campaignId":   campaignID,
		"validationId": validationID,
		"status":       status,
		"timestamp":    time.Now().Format(time.RFC3339),
	}
	proofDataJSON, _ := json.Marshal(proofData)
	validationHash := generateHash(string(proofDataJSON))

	// Create validation proof
	proof := ValidationProof{
		ProofID:        proofID,
		CampaignID:     campaignID,
		ValidationHash: validationHash,
		Status:         status,
		PublishedAt:    time.Now().Format(time.RFC3339),
	}

	proofJSON, err := json.Marshal(proof)
	if err != nil {
		return "", err
	}

	// Store on common-channel
	commonKey := fmt.Sprintf("COMMON_VALIDATION_%s", campaignID)
	err = ctx.GetStub().PutState(commonKey, proofJSON)
	if err != nil {
		return "", err
	}

	// Emit event
	eventPayload := map[string]interface{}{
		"proofId":        proofID,
		"campaignId":     campaignID,
		"validationHash": validationHash,
		"status":         status,
		"channel":        "common-channel",
		"action":         "VALIDATION_PROOF_PUBLISHED",
		"timestamp":      proof.PublishedAt,
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("ValidationProofPublished", eventJSON)

	response := map[string]interface{}{
		"message":        "Validation proof published to common channel",
		"proofId":        proofID,
		"campaignId":     campaignID,
		"validationHash": validationHash,
		"channel":        "common-channel",
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// ============================================================================
// QUERY FUNCTIONS
// ============================================================================

// GetCampaign retrieves campaign from StartupOrg's data on startup-validator-channel
// Since ValidatorOrg and StartupOrg share the same ledger on startup-validator-channel,
// ValidatorOrg can directly query the campaign data that StartupOrg wrote
func (v *ValidatorContract) GetCampaign(ctx contractapi.TransactionContextInterface, campaignID string) (string, error) {
	// Query the campaign directly from the shared ledger
	// StartupOrg stores campaigns using just the campaignID as the key
	campaignJSON, err := ctx.GetStub().GetState(campaignID)
	if err != nil {
		return "", fmt.Errorf("failed to read campaign: %v", err)
	}
	if campaignJSON == nil {
		return "", fmt.Errorf("campaign %s does not exist", campaignID)
	}

	return string(campaignJSON), nil
}

// GetValidation retrieves validation record by ID
func (v *ValidatorContract) GetValidation(ctx contractapi.TransactionContextInterface, validationID string) (*ValidationRecord, error) {
	validationJSON, err := ctx.GetStub().GetState(validationID)
	if err != nil {
		return nil, fmt.Errorf("failed to read validation: %v", err)
	}
	if validationJSON == nil {
		return nil, fmt.Errorf("validation %s does not exist", validationID)
	}

	var validation ValidationRecord
	err = json.Unmarshal(validationJSON, &validation)
	if err != nil {
		return nil, err
	}

	return &validation, nil
}

// GetRiskInsight retrieves risk insight by campaign ID
func (v *ValidatorContract) GetRiskInsight(ctx contractapi.TransactionContextInterface, campaignID string) (*RiskInsight, error) {
	riskKey := fmt.Sprintf("RISK_%s", campaignID)
	insightJSON, err := ctx.GetStub().GetState(riskKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read risk insight: %v", err)
	}
	if insightJSON == nil {
		return nil, fmt.Errorf("risk insight for campaign %s does not exist", campaignID)
	}

	var insight RiskInsight
	err = json.Unmarshal(insightJSON, &insight)
	if err != nil {
		return nil, err
	}

	return &insight, nil
}

// GetValidationReport retrieves validation report by campaign ID
func (v *ValidatorContract) GetValidationReport(ctx contractapi.TransactionContextInterface, campaignID string) (*ValidationReport, error) {
	reportKey := fmt.Sprintf("PLATFORM_REPORT_%s", campaignID)
	reportJSON, err := ctx.GetStub().GetState(reportKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read validation report: %v", err)
	}
	if reportJSON == nil {
		return nil, fmt.Errorf("validation report for campaign %s does not exist", campaignID)
	}

	var report ValidationReport
	err = json.Unmarshal(reportJSON, &report)
	if err != nil {
		return nil, err
	}

	return &report, nil
}

// GetMilestoneVerification retrieves milestone verification by ID
// Channel: startup-validator-channel
func (v *ValidatorContract) GetMilestoneVerification(ctx contractapi.TransactionContextInterface, verificationID string) (*MilestoneValidation, error) {
	verificationJSON, err := ctx.GetStub().GetState(verificationID)
	if err != nil {
		return nil, fmt.Errorf("failed to read milestone verification: %v", err)
	}
	if verificationJSON == nil {
		return nil, fmt.Errorf("milestone verification %s does not exist", verificationID)
	}

	var verification MilestoneValidation
	err = json.Unmarshal(verificationJSON, &verification)
	if err != nil {
		return nil, err
	}

	return &verification, nil
}

// GetMilestoneVerificationByMilestone retrieves milestone verification by milestone ID
// Channel: startup-validator-channel
func (v *ValidatorContract) GetMilestoneVerificationByMilestone(ctx contractapi.TransactionContextInterface, milestoneID string) (*MilestoneValidation, error) {
	milestoneVerifyKey := fmt.Sprintf("MILESTONE_VERIFY_%s", milestoneID)
	verificationJSON, err := ctx.GetStub().GetState(milestoneVerifyKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read milestone verification: %v", err)
	}
	if verificationJSON == nil {
		return nil, fmt.Errorf("milestone verification for milestone %s does not exist", milestoneID)
	}

	var verification MilestoneValidation
	err = json.Unmarshal(verificationJSON, &verification)
	if err != nil {
		return nil, err
	}

	return &verification, nil
}

// GetAgreementWitness retrieves validator's agreement witness record
// Channel: common-channel
func (v *ValidatorContract) GetAgreementWitness(ctx contractapi.TransactionContextInterface, agreementID string) (*AgreementWitness, error) {
	agreementWitnessKey := fmt.Sprintf("VALIDATOR_WITNESS_%s", agreementID)
	witnessJSON, err := ctx.GetStub().GetState(agreementWitnessKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read agreement witness: %v", err)
	}
	if witnessJSON == nil {
		return nil, fmt.Errorf("agreement witness for %s does not exist", agreementID)
	}

	var witness AgreementWitness
	err = json.Unmarshal(witnessJSON, &witness)
	if err != nil {
		return nil, err
	}

	return &witness, nil
}

// GetCampaignCompletion retrieves campaign completion confirmation
// Channel: common-channel
func (v *ValidatorContract) GetCampaignCompletion(ctx contractapi.TransactionContextInterface, campaignID string) (string, error) {
	completionKey := fmt.Sprintf("CAMPAIGN_COMPLETION_%s", campaignID)
	completionJSON, err := ctx.GetStub().GetState(completionKey)
	if err != nil {
		return "", fmt.Errorf("failed to read campaign completion: %v", err)
	}
	if completionJSON == nil {
		return "", fmt.Errorf("campaign completion for %s does not exist", campaignID)
	}

	return string(completionJSON), nil
}

// GetValidationsByCampaign retrieves all validations for a campaign
// Channel: startup-validator-channel
func (v *ValidatorContract) GetValidationsByCampaign(ctx contractapi.TransactionContextInterface, campaignID string) (string, error) {
	queryString := fmt.Sprintf(`{"selector":{"campaignId":"%s","validationId":{"$exists":true}}}`, campaignID)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return "", err
	}
	defer resultsIterator.Close()

	var validations []map[string]interface{}
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return "", err
		}

		var validation ValidationRecord
		err = json.Unmarshal(queryResponse.Value, &validation)
		if err != nil {
			continue
		}

		validationMap := map[string]interface{}{
			"Key":    queryResponse.Key,
			"Record": validation,
		}
		validations = append(validations, validationMap)
	}

	validationsJSON, err := json.Marshal(validations)
	if err != nil {
		return "", err
	}

	return string(validationsJSON), nil
}

// ============================================================================
// CROSS-CHANNEL INVOCATION HELPER FUNCTIONS
// ============================================================================

// InvokeStartupOrgGetCampaign reads campaign data from StartupOrg on startup-validator-channel
// This is a cross-channel READ operation
func (v *ValidatorContract) InvokeStartupOrgGetCampaign(
	ctx contractapi.TransactionContextInterface,
	campaignID string,
) (string, error) {
	args := [][]byte{
		[]byte("GetCampaign"),
		[]byte(campaignID),
	}

	// Cross-channel query to startup chaincode
	// Note: chaincode name is "startup" as deployed, not "startuporg"
	response := ctx.GetStub().InvokeChaincode(
		"startup",
		args,
		"startup-validator-channel",
	)

	if response.Status != 200 {
		return "", fmt.Errorf("cross-channel query to StartupOrg failed: %s", response.Message)
	}

	return string(response.Payload), nil
}

// InvokePlatformOrgRecordDecision sends validation decision to PlatformOrg
// Cross-channel call from startup-validator-channel to validator-platform-channel
func (v *ValidatorContract) InvokePlatformOrgRecordDecision(
	ctx contractapi.TransactionContextInterface,
	recordID string,
	campaignID string,
	validationID string,
	approved string,
	overallScore string,
	reportHash string,
) (string, error) {
	args := [][]byte{
		[]byte("RecordValidatorDecision"),
		[]byte(recordID),
		[]byte(campaignID),
		[]byte(validationID),
		[]byte(approved),
		[]byte(overallScore),
		[]byte(reportHash),
	}

	// Note: chaincode name is "platform" as deployed, not "platformorg"
	response := ctx.GetStub().InvokeChaincode(
		"platform",
		args,
		"validator-platform-channel",
	)

	if response.Status != 200 {
		return "", fmt.Errorf("cross-channel invoke to PlatformOrg failed: %s", response.Message)
	}

	// Emit cross-channel event
	eventPayload := map[string]interface{}{
		"campaignId":     campaignID,
		"validationId":   validationID,
		"targetChannel":  "validator-platform-channel",
		"targetContract": "platformorg",
		"action":         "CROSS_CHANNEL_RECORD_DECISION",
		"timestamp":      time.Now().Format(time.RFC3339),
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("CrossChannelInvoke", eventJSON)

	return string(response.Payload), nil
}

// InvokeInvestorOrgShareRisk shares risk insights with InvestorOrg
// Cross-channel call to investor-validator-channel
func (v *ValidatorContract) InvokeInvestorOrgShareRisk(
	ctx contractapi.TransactionContextInterface,
	insightID string,
	campaignID string,
	riskScore string,
	riskLevel string,
	recommendation string,
) (string, error) {
	args := [][]byte{
		[]byte("ReceiveRiskInsight"),
		[]byte(insightID),
		[]byte(campaignID),
		[]byte(riskScore),
		[]byte(riskLevel),
		[]byte(recommendation),
	}

	// Note: chaincode name is "investor" as deployed, not "investororg"
	response := ctx.GetStub().InvokeChaincode(
		"investor",
		args,
		"investor-validator-channel",
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

// ============================================================================
// DISPUTE INVESTIGATION FUNCTIONS FOR VALIDATOR
// ============================================================================

// AcceptDisputeInvestigation accepts assignment to investigate a dispute
func (v *ValidatorContract) AcceptDisputeInvestigation(
	ctx contractapi.TransactionContextInterface,
	investigationID string,
	disputeID string,
	validatorID string,
	disputeType string,
	initiatorID string,
	initiatorType string,
	respondentID string,
	respondentType string,
	campaignID string,
) (string, error) {
	now := time.Now().Format(time.RFC3339)

	investigation := DisputeInvestigation{
		InvestigationID: investigationID,
		DisputeID:       disputeID,
		ValidatorID:     validatorID,
		DisputeType:     disputeType,
		InitiatorID:     initiatorID,
		InitiatorType:   initiatorType,
		RespondentID:    respondentID,
		RespondentType:  respondentType,
		CampaignID:      campaignID,
		Status:          "ASSIGNED",
		Findings:        []InvestigationFinding{},
		EvidenceReviewed: []string{},
		TransactionLogs: []string{},
		AssignedAt:      now,
	}

	investigationKey := fmt.Sprintf("DISPUTE_INVESTIGATION_%s", investigationID)
	investigationJSON, _ := json.Marshal(investigation)
	ctx.GetStub().PutState(investigationKey, investigationJSON)

	// Emit event
	eventPayload := map[string]interface{}{
		"investigationId": investigationID,
		"disputeId":       disputeID,
		"validatorId":     validatorID,
		"status":          "ASSIGNED",
		"action":          "INVESTIGATION_ACCEPTED",
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("InvestigationAccepted", eventJSON)

	response := map[string]interface{}{
		"message":         "Investigation assignment accepted",
		"investigationId": investigationID,
		"disputeId":       disputeID,
		"status":          "ASSIGNED",
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// RecordInvestigationFinding records a finding during investigation
func (v *ValidatorContract) RecordInvestigationFinding(
	ctx contractapi.TransactionContextInterface,
	investigationID string,
	findingID string,
	findingType string,
	description string,
	severity string,
	relatedEvidence string,
) (string, error) {
	investigationKey := fmt.Sprintf("DISPUTE_INVESTIGATION_%s", investigationID)
	investigationJSON, err := ctx.GetStub().GetState(investigationKey)
	if err != nil || investigationJSON == nil {
		return "", fmt.Errorf("investigation %s not found", investigationID)
	}

	var investigation DisputeInvestigation
	json.Unmarshal(investigationJSON, &investigation)

	now := time.Now().Format(time.RFC3339)
	finding := InvestigationFinding{
		FindingID:       findingID,
		FindingType:     findingType,
		Description:     description,
		Severity:        severity,
		RelatedEvidence: relatedEvidence,
		RecordedAt:      now,
	}

	investigation.Findings = append(investigation.Findings, finding)
	investigation.Status = "IN_PROGRESS"

	if relatedEvidence != "" {
		investigation.EvidenceReviewed = append(investigation.EvidenceReviewed, relatedEvidence)
	}

	investigationUpdatedJSON, _ := json.Marshal(investigation)
	ctx.GetStub().PutState(investigationKey, investigationUpdatedJSON)

	response := map[string]interface{}{
		"message":   "Finding recorded",
		"findingId": findingID,
		"severity":  severity,
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// CompleteInvestigation completes investigation with recommendation
func (v *ValidatorContract) CompleteInvestigation(
	ctx contractapi.TransactionContextInterface,
	investigationID string,
	recommendation string,
	recommendedPenalty string,
) (string, error) {
	investigationKey := fmt.Sprintf("DISPUTE_INVESTIGATION_%s", investigationID)
	investigationJSON, err := ctx.GetStub().GetState(investigationKey)
	if err != nil || investigationJSON == nil {
		return "", fmt.Errorf("investigation %s not found", investigationID)
	}

	var investigation DisputeInvestigation
	json.Unmarshal(investigationJSON, &investigation)

	now := time.Now().Format(time.RFC3339)
	investigation.Status = "COMPLETED"
	investigation.Recommendation = recommendation
	investigation.RecommendedPenalty = recommendedPenalty
	investigation.CompletedAt = now

	investigationUpdatedJSON, _ := json.Marshal(investigation)
	ctx.GetStub().PutState(investigationKey, investigationUpdatedJSON)

	// Emit event for Platform to process
	eventPayload := map[string]interface{}{
		"investigationId":    investigationID,
		"disputeId":          investigation.DisputeID,
		"validatorId":        investigation.ValidatorID,
		"recommendation":     recommendation,
		"recommendedPenalty": recommendedPenalty,
		"status":             "COMPLETED",
		"action":             "INVESTIGATION_COMPLETED",
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("InvestigationCompleted", eventJSON)

	response := map[string]interface{}{
		"message":            "Investigation completed",
		"investigationId":    investigationID,
		"recommendation":     recommendation,
		"recommendedPenalty": recommendedPenalty,
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// InvestigateMilestoneDispute investigates milestone-related disputes
func (v *ValidatorContract) InvestigateMilestoneDispute(
	ctx contractapi.TransactionContextInterface,
	investigationID string,
	disputeID string,
	milestoneID string,
	campaignID string,
	validatorID string,
	deliverableStatus string,
	qualityAssessment float64,
	timelineAssessment string,
	delayJustified bool,
	recommendedAction string,
	comments string,
) (string, error) {
	now := time.Now().Format(time.RFC3339)

	milestoneInvestigation := MilestoneInvestigation{
		InvestigationID:    investigationID,
		DisputeID:          disputeID,
		MilestoneID:        milestoneID,
		CampaignID:         campaignID,
		ValidatorID:        validatorID,
		MilestoneReviewed:  true,
		DeliverableStatus:  deliverableStatus,
		QualityAssessment:  qualityAssessment,
		TimelineAssessment: timelineAssessment,
		DelayJustified:     delayJustified,
		RecommendedAction:  recommendedAction,
		Comments:           comments,
		InvestigatedAt:     now,
	}

	investigationKey := fmt.Sprintf("MILESTONE_INVESTIGATION_%s", investigationID)
	investigationJSON, _ := json.Marshal(milestoneInvestigation)
	ctx.GetStub().PutState(investigationKey, investigationJSON)

	// Emit event
	eventPayload := map[string]interface{}{
		"investigationId":   investigationID,
		"disputeId":         disputeID,
		"milestoneId":       milestoneID,
		"recommendedAction": recommendedAction,
		"action":            "MILESTONE_INVESTIGATED",
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("MilestoneInvestigated", eventJSON)

	response := map[string]interface{}{
		"message":           "Milestone investigation completed",
		"investigationId":   investigationID,
		"deliverableStatus": deliverableStatus,
		"recommendedAction": recommendedAction,
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// RespondToDispute allows validator to respond when they are the respondent
func (v *ValidatorContract) RespondToDispute(
	ctx contractapi.TransactionContextInterface,
	responseID string,
	disputeID string,
	validatorID string,
	responseText string,
	justification string,
	supportingDocsJSON string,
) (string, error) {
	var supportingDocs []string
	if supportingDocsJSON != "" {
		json.Unmarshal([]byte(supportingDocsJSON), &supportingDocs)
	}

	now := time.Now().Format(time.RFC3339)

	disputeResponse := ValidatorDisputeResponse{
		ResponseID:     responseID,
		DisputeID:      disputeID,
		ValidatorID:    validatorID,
		ResponseText:   responseText,
		Justification:  justification,
		SupportingDocs: supportingDocs,
		RespondedAt:    now,
	}

	responseKey := fmt.Sprintf("VALIDATOR_DISPUTE_RESPONSE_%s", responseID)
	responseJSON, _ := json.Marshal(disputeResponse)
	ctx.GetStub().PutState(responseKey, responseJSON)

	// Emit event
	eventPayload := map[string]interface{}{
		"responseId":  responseID,
		"disputeId":   disputeID,
		"validatorId": validatorID,
		"action":      "VALIDATOR_DISPUTE_RESPONSE",
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("ValidatorDisputeResponse", eventJSON)

	return fmt.Sprintf(`{"message": "Response submitted", "responseId": "%s"}`, responseID), nil
}

// GetInvestigation retrieves an investigation by ID
func (v *ValidatorContract) GetInvestigation(ctx contractapi.TransactionContextInterface, investigationID string) (string, error) {
	investigationKey := fmt.Sprintf("DISPUTE_INVESTIGATION_%s", investigationID)
	investigationJSON, err := ctx.GetStub().GetState(investigationKey)
	if err != nil || investigationJSON == nil {
		return "", fmt.Errorf("investigation %s not found", investigationID)
	}
	return string(investigationJSON), nil
}

// GetValidatorDisputes retrieves all disputes involving a validator
func (v *ValidatorContract) GetValidatorDisputes(ctx contractapi.TransactionContextInterface, validatorID string) (string, error) {
	queryString := fmt.Sprintf(`{"selector":{"validatorId":"%s"}}`, validatorID)
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return "[]", nil
	}
	defer resultsIterator.Close()

	var investigations []DisputeInvestigation
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			continue
		}

		var inv DisputeInvestigation
		if err := json.Unmarshal(queryResponse.Value, &inv); err == nil {
			investigations = append(investigations, inv)
		}
	}

	investigationsJSON, _ := json.Marshal(investigations)
	return string(investigationsJSON), nil
}

func main() {
	validatorChaincode, err := contractapi.NewChaincode(&ValidatorContract{})
	if err != nil {
		fmt.Printf("Error creating ValidatorOrg chaincode: %v\n", err)
		return
	}

	if err := validatorChaincode.Start(); err != nil {
		fmt.Printf("Error starting ValidatorOrg chaincode: %v\n", err)
	}
}
