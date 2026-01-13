package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
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

// ============================================================================
// WALLET & TOKEN SYSTEM
// ============================================================================

// Wallet represents a user's cryptographic wallet for platform tokens
type Wallet struct {
	WalletID        string  `json:"walletId"`
	UserID          string  `json:"userId"`
	UserType        string  `json:"userType"` // STARTUP, INVESTOR, VALIDATOR, PLATFORM
	Balance         float64 `json:"balance"`
	LockedBalance   float64 `json:"lockedBalance"` // Funds locked in escrow/disputes
	TotalDeposited  float64 `json:"totalDeposited"`
	TotalWithdrawn  float64 `json:"totalWithdrawn"`
	TotalPenalties  float64 `json:"totalPenalties"`
	CreatedAt       string  `json:"createdAt"`
	UpdatedAt       string  `json:"updatedAt"`
}

// TokenTransaction represents a token transfer/exchange transaction
type TokenTransaction struct {
	TransactionID   string  `json:"transactionId"`
	FromWalletID    string  `json:"fromWalletId"`
	ToWalletID      string  `json:"toWalletId"`
	Amount          float64 `json:"amount"`
	TransactionType string  `json:"transactionType"` // DEPOSIT, WITHDRAW, TRANSFER, FEE, PENALTY, REFUND, ESCROW, RELEASE
	Reference       string  `json:"reference"`       // CampaignID, DisputeID, etc.
	Status          string  `json:"status"`          // PENDING, COMPLETED, FAILED, REVERSED
	CreatedAt       string  `json:"createdAt"`
}

// TokenExchangeRate represents exchange rate between real currency and platform tokens
type TokenExchangeRate struct {
	RateID       string  `json:"rateId"`
	Currency     string  `json:"currency"`     // USD, EUR, INR, etc.
	TokenRate    float64 `json:"tokenRate"`    // 1 Currency = X Tokens
	EffectiveAt  string  `json:"effectiveAt"`
	ExpiresAt    string  `json:"expiresAt"`
	SetBy        string  `json:"setBy"`
}

// ============================================================================
// PLATFORM FEE SYSTEM
// ============================================================================

// FeeTier represents platform fee tier based on campaign goal amount
type FeeTier struct {
	TierID          string  `json:"tierId"`
	TierName        string  `json:"tierName"`        // SMALL, MEDIUM, LARGE, ENTERPRISE
	MinGoalAmount   float64 `json:"minGoalAmount"`
	MaxGoalAmount   float64 `json:"maxGoalAmount"`
	FixedFee        float64 `json:"fixedFee"`        // Fixed fee in tokens
	PercentageFee   float64 `json:"percentageFee"`   // Percentage of goal (optional for future)
	Description     string  `json:"description"`
	IsActive        bool    `json:"isActive"`
}

// CampaignFee represents fee collected from a startup for campaign
type CampaignFee struct {
	FeeID           string  `json:"feeId"`
	CampaignID      string  `json:"campaignId"`
	StartupID       string  `json:"startupId"`
	TierID          string  `json:"tierId"`
	TierName        string  `json:"tierName"`
	GoalAmount      float64 `json:"goalAmount"`
	FeeAmount       float64 `json:"feeAmount"`
	Status          string  `json:"status"` // PENDING, COLLECTED, REFUNDED
	TransactionID   string  `json:"transactionId"`
	CollectedAt     string  `json:"collectedAt"`
}

// ============================================================================
// ML RATING SYSTEM
// ============================================================================

// Rating scale: 0-100, append-only historical
type Rating struct {
	RatingID        string  `json:"ratingId"`
	TargetUserID    string  `json:"targetUserId"`
	TargetUserType  string  `json:"targetUserType"` // STARTUP, INVESTOR, VALIDATOR
	RaterType       string  `json:"raterType"`      // ML_MODEL, USER, SYSTEM
	RaterID         string  `json:"raterId"`        // ML model ID or user ID
	Score           float64 `json:"score"`          // 0-100 scale
	Category        string  `json:"category"`       // CREDIBILITY, RISK, PERFORMANCE, COMPLIANCE, OVERALL
	Factors         string  `json:"factors"`        // JSON string of rating factors
	EvidenceHash    string  `json:"evidenceHash"`   // IPFS hash of supporting evidence
	CreatedAt       string  `json:"createdAt"`
}

// RatingAggregate stores aggregated rating for a user
type RatingAggregate struct {
	UserID              string  `json:"userId"`
	UserType            string  `json:"userType"`
	OverallScore        float64 `json:"overallScore"`
	CredibilityScore    float64 `json:"credibilityScore"`
	RiskScore           float64 `json:"riskScore"`
	PerformanceScore    float64 `json:"performanceScore"`
	ComplianceScore     float64 `json:"complianceScore"`
	TotalRatings        int     `json:"totalRatings"`
	LastUpdated         string  `json:"lastUpdated"`
}

// ============================================================================
// DISPUTE SYSTEM
// ============================================================================

// DisputeType enum values
const (
	DisputeStartupInvestor   = "STARTUP_INVESTOR"
	DisputeStartupValidator  = "STARTUP_VALIDATOR"
	DisputeInvestorValidator = "INVESTOR_VALIDATOR"
	DisputeStartupPlatform   = "STARTUP_PLATFORM"
	DisputeInvestorPlatform  = "INVESTOR_PLATFORM"
	DisputeMultilateral      = "MULTILATERAL"
)

// DisputeSubType enum values
const (
	SubTypeMisuseOfFunds          = "MISUSE_OF_FUNDS"
	SubTypeDelayedDeliverables    = "DELAYED_DELIVERABLES"
	SubTypeFraudulentClaim        = "FRAUDULENT_CLAIM"
	SubTypeFraudulentApproval     = "FRAUDULENT_APPROVAL"
	SubTypeDelayedValidation      = "DELAYED_VALIDATION"
	SubTypeUnnecessaryDocRequests = "UNNECESSARY_DOC_REQUESTS"
	SubTypeNegligentValidation    = "NEGLIGENT_VALIDATION"
	SubTypeBiasedVerification     = "BIASED_VERIFICATION"
	SubTypeIncorrectSuspension    = "INCORRECT_SUSPENSION"
	SubTypeIncorrectFees          = "INCORRECT_FEES"
	SubTypeIncorrectRefund        = "INCORRECT_REFUND"
	SubTypeWrongfulBan            = "WRONGFUL_BAN"
	SubTypeSmartContractError     = "SMART_CONTRACT_ERROR"
)

// ============================================================================
// DISPUTE FEE SYSTEM - Prevents spam/frivolous disputes
// ============================================================================

// DisputeFeeTier represents fee tiers based on dispute severity
type DisputeFeeTier struct {
	TierID          string  `json:"tierId"`
	TierName        string  `json:"tierName"`        // MINOR, STANDARD, MAJOR, CRITICAL
	FilingFee       float64 `json:"filingFee"`       // Fee to file dispute (in tokens)
	MinClaimAmount  float64 `json:"minClaimAmount"`  // Minimum claim amount for this tier
	MaxClaimAmount  float64 `json:"maxClaimAmount"`  // Maximum claim amount for this tier
	RefundOnWin     bool    `json:"refundOnWin"`     // Whether fee is refunded if initiator wins
	RefundPercent   float64 `json:"refundPercent"`   // Percentage refunded on win (100 = full refund)
	Description     string  `json:"description"`
	IsActive        bool    `json:"isActive"`
}

// DisputeFeeRecord tracks fee collected for a dispute
type DisputeFeeRecord struct {
	FeeRecordID     string  `json:"feeRecordId"`
	DisputeID       string  `json:"disputeId"`
	InitiatorID     string  `json:"initiatorId"`
	InitiatorType   string  `json:"initiatorType"`
	WalletID        string  `json:"walletId"`
	TierID          string  `json:"tierId"`
	TierName        string  `json:"tierName"`
	FeeAmount       float64 `json:"feeAmount"`
	ClaimAmount     float64 `json:"claimAmount"`
	Status          string  `json:"status"`          // COLLECTED, REFUNDED, FORFEITED
	RefundAmount    float64 `json:"refundAmount"`    // Amount refunded (if won)
	TransactionID   string  `json:"transactionId"`
	RefundTxnID     string  `json:"refundTxnId"`
	CollectedAt     string  `json:"collectedAt"`
	ProcessedAt     string  `json:"processedAt"`     // When refund/forfeit happened
}

// Dispute fee tier constants
const (
	DisputeFeeMinor    = 10.0   // Minor disputes (< 500 claimed)
	DisputeFeeStandard = 25.0   // Standard disputes (500 - 5000 claimed)
	DisputeFeeMajor    = 50.0   // Major disputes (5000 - 50000 claimed)
	DisputeFeeCritical = 100.0  // Critical disputes (> 50000 claimed)
)

// Dispute represents a dispute ticket on common-channel
type Dispute struct {
	DisputeID           string              `json:"disputeId"`
	TicketNumber        string              `json:"ticketNumber"`
	DisputeType         string              `json:"disputeType"`    // STARTUP_INVESTOR, STARTUP_VALIDATOR, etc.
	DisputeSubType      string              `json:"disputeSubType"` // MISUSE_OF_FUNDS, DELAYED_DELIVERABLES, etc.
	
	// Parties
	InitiatorID         string              `json:"initiatorId"`
	InitiatorType       string              `json:"initiatorType"` // STARTUP, INVESTOR, VALIDATOR, PLATFORM
	RespondentID        string              `json:"respondentId"`
	RespondentType      string              `json:"respondentType"`
	AdditionalParties   []DisputeParty      `json:"additionalParties"` // For multilateral disputes
	
	// Related entities
	CampaignID          string              `json:"campaignId"`
	AgreementID         string              `json:"agreementId"`
	MilestoneID         string              `json:"milestoneId"`
	TransactionID       string              `json:"transactionId"`
	
	// Dispute details
	Title               string              `json:"title"`
	Description         string              `json:"description"`
	ClaimedAmount       float64             `json:"claimedAmount"` // Amount in dispute
	
	// Dispute Fee (prevents spam)
	FilingFeeID         string              `json:"filingFeeId"`     // Reference to DisputeFeeRecord
	FilingFeePaid       float64             `json:"filingFeePaid"`   // Amount paid to file
	FilingFeeStatus     string              `json:"filingFeeStatus"` // PAID, REFUNDED, FORFEITED
	
	// Evidence (IPFS hashes)
	EvidenceHashes      []Evidence          `json:"evidenceHashes"`
	
	// Investigation
	InvestigatorID      string              `json:"investigatorId"` // Validator or Platform admin
	InvestigationNotes  []InvestigationNote `json:"investigationNotes"`
	SmartContractRules  []string            `json:"smartContractRules"` // Rules invoked
	
	// Voting
	VotingEnabled       bool                `json:"votingEnabled"`
	VotingDeadline      string              `json:"votingDeadline"`
	VotingResult        *VotingResult       `json:"votingResult"`
	
	// Resolution
	Status              string              `json:"status"` // OPEN, UNDER_INVESTIGATION, VOTING, RESOLVED, APPEALED, CLOSED
	Resolution          string              `json:"resolution"` // FAVOR_INITIATOR, FAVOR_RESPONDENT, PARTIAL, DISMISSED
	ResolutionNotes     string              `json:"resolutionNotes"`
	
	// Penalties
	PenaltiesApplied    []Penalty           `json:"penaltiesApplied"`
	RefundsOrdered      []RefundOrder       `json:"refundsOrdered"`
	
	// Timestamps
	CreatedAt           string              `json:"createdAt"`
	UpdatedAt           string              `json:"updatedAt"`
	ResolvedAt          string              `json:"resolvedAt"`
}

// DisputeParty represents a party in a dispute
type DisputeParty struct {
	PartyID    string `json:"partyId"`
	PartyType  string `json:"partyType"`
	Role       string `json:"role"` // INITIATOR, RESPONDENT, WITNESS, INVESTIGATOR
}

// Evidence represents evidence submitted in a dispute
type Evidence struct {
	EvidenceID    string `json:"evidenceId"`
	SubmittedBy   string `json:"submittedBy"`
	SubmitterType string `json:"submitterType"`
	IPFSHash      string `json:"ipfsHash"`
	Description   string `json:"description"`
	EvidenceType  string `json:"evidenceType"` // DOCUMENT, TRANSACTION_LOG, SCREENSHOT, VIDEO, OTHER
	SubmittedAt   string `json:"submittedAt"`
}

// InvestigationNote represents a note during investigation
type InvestigationNote struct {
	NoteID        string `json:"noteId"`
	InvestigatorID string `json:"investigatorId"`
	Note          string `json:"note"`
	FindingType   string `json:"findingType"` // OBSERVATION, EVIDENCE_REVIEW, CONCLUSION
	CreatedAt     string `json:"createdAt"`
}

// ============================================================================
// ANONYMOUS VOTING SYSTEM (Commit-Reveal Scheme)
// ============================================================================

// VoteCommitment represents a committed vote (hash only - anonymous)
type VoteCommitment struct {
	CommitmentID  string `json:"commitmentId"`
	DisputeID     string `json:"disputeId"`
	VoterHash     string `json:"voterHash"`     // Hash of voterID + salt (anonymous)
	VoteHash      string `json:"voteHash"`      // Hash of vote + salt
	CommittedAt   string `json:"committedAt"`
}

// VoteReveal represents a revealed vote (after voting deadline)
type VoteReveal struct {
	RevealID      string `json:"revealId"`
	DisputeID     string `json:"disputeId"`
	CommitmentID  string `json:"commitmentId"`
	VoterOrgType  string `json:"voterOrgType"` // STARTUP, INVESTOR, VALIDATOR, PLATFORM (no specific user ID for anonymity)
	Vote          string `json:"vote"`         // FAVOR_INITIATOR, FAVOR_RESPONDENT, ABSTAIN
	Salt          string `json:"salt"`
	RevealedAt    string `json:"revealedAt"`
	IsValid       bool   `json:"isValid"`      // Whether reveal matches commitment
}

// VotingResult represents the final voting outcome
type VotingResult struct {
	DisputeID         string `json:"disputeId"`
	TotalVotes        int    `json:"totalVotes"`
	FavorInitiator    int    `json:"favorInitiator"`
	FavorRespondent   int    `json:"favorRespondent"`
	Abstained         int    `json:"abstained"`
	InvalidVotes      int    `json:"invalidVotes"`
	Outcome           string `json:"outcome"` // FAVOR_INITIATOR, FAVOR_RESPONDENT, TIE
	TalliedAt         string `json:"talliedAt"`
}

// ============================================================================
// PENALTY & REPUTATION SYSTEM
// ============================================================================

// PenaltySeverity levels
const (
	SeverityLow      = "LOW"
	SeverityMedium   = "MEDIUM"
	SeverityHigh     = "HIGH"
	SeverityCritical = "CRITICAL"
)

// Penalty represents a penalty applied to a user
type Penalty struct {
	PenaltyID       string  `json:"penaltyId"`
	UserID          string  `json:"userId"`
	UserType        string  `json:"userType"`
	DisputeID       string  `json:"disputeId"`
	PenaltyType     string  `json:"penaltyType"` // FINE, REPUTATION_DEDUCTION, SUSPENSION, BLACKLIST
	Severity        string  `json:"severity"`    // LOW, MEDIUM, HIGH, CRITICAL
	TokenAmount     float64 `json:"tokenAmount"` // Fine amount in tokens
	ReputationDeduct float64 `json:"reputationDeduct"` // Points to deduct
	Description     string  `json:"description"`
	Status          string  `json:"status"` // PENDING, APPLIED, REVERSED
	AppliedAt       string  `json:"appliedAt"`
}

// RefundOrder represents a refund ordered as part of dispute resolution
type RefundOrder struct {
	RefundOrderID   string  `json:"refundOrderId"`
	DisputeID       string  `json:"disputeId"`
	FromUserID      string  `json:"fromUserId"`
	FromUserType    string  `json:"fromUserType"`
	ToUserID        string  `json:"toUserId"`
	ToUserType      string  `json:"toUserType"`
	Amount          float64 `json:"amount"`
	DeductionPercent float64 `json:"deductionPercent"` // 15-30% for mid-agreement withdrawal
	NetAmount       float64 `json:"netAmount"`
	Reason          string  `json:"reason"`
	Status          string  `json:"status"` // PENDING, PROCESSED, FAILED
	ProcessedAt     string  `json:"processedAt"`
}

// Reputation tracks user reputation over time
type Reputation struct {
	UserID              string  `json:"userId"`
	UserType            string  `json:"userType"`
	CurrentScore        float64 `json:"currentScore"`        // 0-100
	BaselineScore       float64 `json:"baselineScore"`       // Starting score (default 50)
	TotalDisputes       int     `json:"totalDisputes"`
	DisputesLost        int     `json:"disputesLost"`
	DisputesWon         int     `json:"disputesWon"`
	TotalPenalties      int     `json:"totalPenalties"`
	ConsecutivePenalties int    `json:"consecutivePenalties"` // For auto-suspension
	Status              string  `json:"status"`              // ACTIVE, SUSPENDED, BLACKLISTED
	SuspendedUntil      string  `json:"suspendedUntil"`
	BlacklistedAt       string  `json:"blacklistedAt"`
	BlacklistReason     string  `json:"blacklistReason"`
	UpdatedAt           string  `json:"updatedAt"`
}

// ReputationHistory tracks reputation changes
type ReputationHistory struct {
	HistoryID       string  `json:"historyId"`
	UserID          string  `json:"userId"`
	UserType        string  `json:"userType"`
	PreviousScore   float64 `json:"previousScore"`
	NewScore        float64 `json:"newScore"`
	ChangeAmount    float64 `json:"changeAmount"`
	ChangeReason    string  `json:"changeReason"`
	DisputeID       string  `json:"disputeId"`
	CreatedAt       string  `json:"createdAt"`
}

// Suspension/Blacklist thresholds (configurable)
const (
	ReputationSuspensionThreshold = 30.0  // Below this = temporary suspension
	ReputationBlacklistThreshold  = 15.0  // Below this = permanent blacklist
	ConsecutivePenaltyThreshold   = 3     // 3 consecutive penalties = auto-suspension
	DefaultReputationScore        = 50.0  // Starting reputation
	MaxReputationScore            = 100.0
	MinReputationScore            = 0.0
)

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
		"startup",  // chaincode name (deployed as "startup")
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
		"validator",  // chaincode name (deployed as "validator")
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
		"startup",  // chaincode name (deployed as "startup")
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
		"platform", // same contract but different channel (deployed as "platform")
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

// ============================================================================
// WALLET & TOKEN SYSTEM FUNCTIONS
// ============================================================================

// CreateWallet creates a new wallet for a user
// Channel: common-channel
func (p *PlatformContract) CreateWallet(
	ctx contractapi.TransactionContextInterface,
	walletID string,
	userID string,
	userType string,
	initialDeposit float64,
) (string, error) {
	// Validate user type
	validTypes := map[string]bool{"STARTUP": true, "INVESTOR": true, "VALIDATOR": true, "PLATFORM": true}
	if !validTypes[userType] {
		return "", fmt.Errorf("invalid user type: %s", userType)
	}

	// Check if wallet already exists
	walletKey := fmt.Sprintf("WALLET_%s", walletID)
	existing, _ := ctx.GetStub().GetState(walletKey)
	if existing != nil {
		return "", fmt.Errorf("wallet %s already exists", walletID)
	}

	// Check if user already has a wallet
	userWalletKey := fmt.Sprintf("USER_WALLET_%s_%s", userType, userID)
	existingUserWallet, _ := ctx.GetStub().GetState(userWalletKey)
	if existingUserWallet != nil {
		return "", fmt.Errorf("user %s already has a wallet", userID)
	}

	now := time.Now().Format(time.RFC3339)

	wallet := Wallet{
		WalletID:       walletID,
		UserID:         userID,
		UserType:       userType,
		Balance:        initialDeposit,
		LockedBalance:  0,
		TotalDeposited: initialDeposit,
		TotalWithdrawn: 0,
		TotalPenalties: 0,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	walletJSON, _ := json.Marshal(wallet)
	err := ctx.GetStub().PutState(walletKey, walletJSON)
	if err != nil {
		return "", err
	}

	// Store user-to-wallet mapping
	err = ctx.GetStub().PutState(userWalletKey, []byte(walletID))
	if err != nil {
		return "", err
	}

	// Create initial reputation for the user
	reputationKey := fmt.Sprintf("REPUTATION_%s_%s", userType, userID)
	reputation := Reputation{
		UserID:              userID,
		UserType:            userType,
		CurrentScore:        DefaultReputationScore,
		BaselineScore:       DefaultReputationScore,
		TotalDisputes:       0,
		DisputesLost:        0,
		DisputesWon:         0,
		TotalPenalties:      0,
		ConsecutivePenalties: 0,
		Status:              "ACTIVE",
		UpdatedAt:           now,
	}
	reputationJSON, _ := json.Marshal(reputation)
	ctx.GetStub().PutState(reputationKey, reputationJSON)

	// Record initial deposit transaction if > 0
	if initialDeposit > 0 {
		txnID := fmt.Sprintf("TXN_%s_%d", walletID, time.Now().UnixNano())
		txn := TokenTransaction{
			TransactionID:   txnID,
			FromWalletID:    "SYSTEM",
			ToWalletID:      walletID,
			Amount:          initialDeposit,
			TransactionType: "DEPOSIT",
			Reference:       "Initial deposit",
			Status:          "COMPLETED",
			CreatedAt:       now,
		}
		txnJSON, _ := json.Marshal(txn)
		ctx.GetStub().PutState(fmt.Sprintf("TXN_%s", txnID), txnJSON)
	}

	response := map[string]interface{}{
		"message":   "Wallet created successfully",
		"walletId":  walletID,
		"userId":    userID,
		"userType":  userType,
		"balance":   initialDeposit,
		"createdAt": now,
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// DepositTokens adds tokens to a wallet
func (p *PlatformContract) DepositTokens(
	ctx contractapi.TransactionContextInterface,
	walletID string,
	amount float64,
	reference string,
) (string, error) {
	if amount <= 0 {
		return "", fmt.Errorf("deposit amount must be positive")
	}

	walletKey := fmt.Sprintf("WALLET_%s", walletID)
	walletJSON, err := ctx.GetStub().GetState(walletKey)
	if err != nil || walletJSON == nil {
		return "", fmt.Errorf("wallet %s not found", walletID)
	}

	var wallet Wallet
	json.Unmarshal(walletJSON, &wallet)

	now := time.Now().Format(time.RFC3339)
	wallet.Balance += amount
	wallet.TotalDeposited += amount
	wallet.UpdatedAt = now

	updatedJSON, _ := json.Marshal(wallet)
	ctx.GetStub().PutState(walletKey, updatedJSON)

	// Record transaction
	txnID := fmt.Sprintf("TXN_%s_%d", walletID, time.Now().UnixNano())
	txn := TokenTransaction{
		TransactionID:   txnID,
		FromWalletID:    "SYSTEM",
		ToWalletID:      walletID,
		Amount:          amount,
		TransactionType: "DEPOSIT",
		Reference:       reference,
		Status:          "COMPLETED",
		CreatedAt:       now,
	}
	txnJSON, _ := json.Marshal(txn)
	ctx.GetStub().PutState(fmt.Sprintf("TXN_%s", txnID), txnJSON)

	response := map[string]interface{}{
		"message":    "Deposit successful",
		"walletId":   walletID,
		"amount":     amount,
		"newBalance": wallet.Balance,
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// TransferTokens transfers tokens between wallets
func (p *PlatformContract) TransferTokens(
	ctx contractapi.TransactionContextInterface,
	fromWalletID string,
	toWalletID string,
	amount float64,
	transactionType string,
	reference string,
) (string, error) {
	if amount <= 0 {
		return "", fmt.Errorf("transfer amount must be positive")
	}

	// Get source wallet
	fromKey := fmt.Sprintf("WALLET_%s", fromWalletID)
	fromJSON, err := ctx.GetStub().GetState(fromKey)
	if err != nil || fromJSON == nil {
		return "", fmt.Errorf("source wallet %s not found", fromWalletID)
	}

	var fromWallet Wallet
	json.Unmarshal(fromJSON, &fromWallet)

	if fromWallet.Balance < amount {
		return "", fmt.Errorf("insufficient balance: have %.2f, need %.2f", fromWallet.Balance, amount)
	}

	// Get destination wallet
	toKey := fmt.Sprintf("WALLET_%s", toWalletID)
	toJSON, err := ctx.GetStub().GetState(toKey)
	if err != nil || toJSON == nil {
		return "", fmt.Errorf("destination wallet %s not found", toWalletID)
	}

	var toWallet Wallet
	json.Unmarshal(toJSON, &toWallet)

	now := time.Now().Format(time.RFC3339)

	// Update balances
	fromWallet.Balance -= amount
	fromWallet.UpdatedAt = now
	toWallet.Balance += amount
	toWallet.UpdatedAt = now

	fromUpdatedJSON, _ := json.Marshal(fromWallet)
	ctx.GetStub().PutState(fromKey, fromUpdatedJSON)

	toUpdatedJSON, _ := json.Marshal(toWallet)
	ctx.GetStub().PutState(toKey, toUpdatedJSON)

	// Record transaction
	txnID := fmt.Sprintf("TXN_%d", time.Now().UnixNano())
	txn := TokenTransaction{
		TransactionID:   txnID,
		FromWalletID:    fromWalletID,
		ToWalletID:      toWalletID,
		Amount:          amount,
		TransactionType: transactionType,
		Reference:       reference,
		Status:          "COMPLETED",
		CreatedAt:       now,
	}
	txnJSON, _ := json.Marshal(txn)
	ctx.GetStub().PutState(fmt.Sprintf("TXN_%s", txnID), txnJSON)

	response := map[string]interface{}{
		"message":       "Transfer successful",
		"transactionId": txnID,
		"fromWallet":    fromWalletID,
		"toWallet":      toWalletID,
		"amount":        amount,
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// GetWallet retrieves wallet information
func (p *PlatformContract) GetWallet(ctx contractapi.TransactionContextInterface, walletID string) (string, error) {
	walletKey := fmt.Sprintf("WALLET_%s", walletID)
	walletJSON, err := ctx.GetStub().GetState(walletKey)
	if err != nil || walletJSON == nil {
		return "", fmt.Errorf("wallet %s not found", walletID)
	}
	return string(walletJSON), nil
}

// GetWalletByUser retrieves wallet by user ID and type
func (p *PlatformContract) GetWalletByUser(ctx contractapi.TransactionContextInterface, userType string, userID string) (string, error) {
	userWalletKey := fmt.Sprintf("USER_WALLET_%s_%s", userType, userID)
	walletIDBytes, err := ctx.GetStub().GetState(userWalletKey)
	if err != nil || walletIDBytes == nil {
		return "", fmt.Errorf("no wallet found for user %s", userID)
	}

	walletID := string(walletIDBytes)
	return p.GetWallet(ctx, walletID)
}

// SetExchangeRate sets the exchange rate between currency and tokens
func (p *PlatformContract) SetExchangeRate(
	ctx contractapi.TransactionContextInterface,
	rateID string,
	currency string,
	tokenRate float64,
	validDays int,
) (string, error) {
	now := time.Now()
	expiresAt := now.AddDate(0, 0, validDays)

	rate := TokenExchangeRate{
		RateID:      rateID,
		Currency:    currency,
		TokenRate:   tokenRate,
		EffectiveAt: now.Format(time.RFC3339),
		ExpiresAt:   expiresAt.Format(time.RFC3339),
		SetBy:       "PLATFORM",
	}

	rateKey := fmt.Sprintf("EXCHANGE_RATE_%s", currency)
	rateJSON, _ := json.Marshal(rate)
	ctx.GetStub().PutState(rateKey, rateJSON)

	response := map[string]interface{}{
		"message":   "Exchange rate set",
		"currency":  currency,
		"tokenRate": tokenRate,
		"expiresAt": rate.ExpiresAt,
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// ============================================================================
// PLATFORM FEE SYSTEM FUNCTIONS
// ============================================================================

// InitializeFeeTiers sets up the fee tier structure
func (p *PlatformContract) InitializeFeeTiers(ctx contractapi.TransactionContextInterface) (string, error) {
	tiers := []FeeTier{
		{
			TierID:        "TIER_SMALL",
			TierName:      "SMALL",
			MinGoalAmount: 0,
			MaxGoalAmount: 10000,
			FixedFee:      50,
			PercentageFee: 0,
			Description:   "Small campaigns: $0 - $10,000",
			IsActive:      true,
		},
		{
			TierID:        "TIER_MEDIUM",
			TierName:      "MEDIUM",
			MinGoalAmount: 10001,
			MaxGoalAmount: 50000,
			FixedFee:      150,
			PercentageFee: 0,
			Description:   "Medium campaigns: $10,001 - $50,000",
			IsActive:      true,
		},
		{
			TierID:        "TIER_LARGE",
			TierName:      "LARGE",
			MinGoalAmount: 50001,
			MaxGoalAmount: 200000,
			FixedFee:      400,
			PercentageFee: 0,
			Description:   "Large campaigns: $50,001 - $200,000",
			IsActive:      true,
		},
		{
			TierID:        "TIER_ENTERPRISE",
			TierName:      "ENTERPRISE",
			MinGoalAmount: 200001,
			MaxGoalAmount: 999999999,
			FixedFee:      1000,
			PercentageFee: 0,
			Description:   "Enterprise campaigns: $200,001+",
			IsActive:      true,
		},
	}

	for _, tier := range tiers {
		tierKey := fmt.Sprintf("FEE_TIER_%s", tier.TierID)
		tierJSON, _ := json.Marshal(tier)
		ctx.GetStub().PutState(tierKey, tierJSON)
	}

	// Store tier list
	tierListJSON, _ := json.Marshal([]string{"TIER_SMALL", "TIER_MEDIUM", "TIER_LARGE", "TIER_ENTERPRISE"})
	ctx.GetStub().PutState("FEE_TIER_LIST", tierListJSON)

	return `{"message": "Fee tiers initialized successfully", "tiers": 4}`, nil
}

// GetFeeTier determines the fee tier for a given goal amount
func (p *PlatformContract) GetFeeTier(ctx contractapi.TransactionContextInterface, goalAmount float64) (string, error) {
	tierListJSON, _ := ctx.GetStub().GetState("FEE_TIER_LIST")
	if tierListJSON == nil {
		return "", fmt.Errorf("fee tiers not initialized")
	}

	var tierList []string
	json.Unmarshal(tierListJSON, &tierList)

	for _, tierID := range tierList {
		tierKey := fmt.Sprintf("FEE_TIER_%s", tierID)
		tierJSON, _ := ctx.GetStub().GetState(tierKey)
		if tierJSON == nil {
			continue
		}

		var tier FeeTier
		json.Unmarshal(tierJSON, &tier)

		if tier.IsActive && goalAmount >= tier.MinGoalAmount && goalAmount <= tier.MaxGoalAmount {
			return string(tierJSON), nil
		}
	}

	return "", fmt.Errorf("no matching fee tier for amount %.2f", goalAmount)
}

// CollectCampaignFee collects fee from startup for campaign creation
func (p *PlatformContract) CollectCampaignFee(
	ctx contractapi.TransactionContextInterface,
	feeID string,
	campaignID string,
	startupID string,
	startupWalletID string,
	goalAmount float64,
) (string, error) {
	// Get applicable fee tier
	tierJSON, err := p.GetFeeTier(ctx, goalAmount)
	if err != nil {
		return "", err
	}

	var tier FeeTier
	json.Unmarshal([]byte(tierJSON), &tier)

	// Get startup wallet
	walletKey := fmt.Sprintf("WALLET_%s", startupWalletID)
	walletJSON, err := ctx.GetStub().GetState(walletKey)
	if err != nil || walletJSON == nil {
		return "", fmt.Errorf("startup wallet %s not found", startupWalletID)
	}

	var wallet Wallet
	json.Unmarshal(walletJSON, &wallet)

	if wallet.Balance < tier.FixedFee {
		return "", fmt.Errorf("insufficient balance for fee: have %.2f tokens, need %.2f", wallet.Balance, tier.FixedFee)
	}

	now := time.Now().Format(time.RFC3339)

	// Deduct fee from startup wallet
	wallet.Balance -= tier.FixedFee
	wallet.UpdatedAt = now
	walletUpdatedJSON, _ := json.Marshal(wallet)
	ctx.GetStub().PutState(walletKey, walletUpdatedJSON)

	// Record fee collection
	fee := CampaignFee{
		FeeID:         feeID,
		CampaignID:    campaignID,
		StartupID:     startupID,
		TierID:        tier.TierID,
		TierName:      tier.TierName,
		GoalAmount:    goalAmount,
		FeeAmount:     tier.FixedFee,
		Status:        "COLLECTED",
		TransactionID: fmt.Sprintf("FEE_TXN_%s", feeID),
		CollectedAt:   now,
	}

	feeKey := fmt.Sprintf("CAMPAIGN_FEE_%s", feeID)
	feeJSON, _ := json.Marshal(fee)
	ctx.GetStub().PutState(feeKey, feeJSON)

	// Record transaction
	txn := TokenTransaction{
		TransactionID:   fee.TransactionID,
		FromWalletID:    startupWalletID,
		ToWalletID:      "PLATFORM_TREASURY",
		Amount:          tier.FixedFee,
		TransactionType: "FEE",
		Reference:       fmt.Sprintf("Campaign fee for %s", campaignID),
		Status:          "COMPLETED",
		CreatedAt:       now,
	}
	txnJSON, _ := json.Marshal(txn)
	ctx.GetStub().PutState(fmt.Sprintf("TXN_%s", fee.TransactionID), txnJSON)

	response := map[string]interface{}{
		"message":      "Campaign fee collected",
		"feeId":        feeID,
		"campaignId":   campaignID,
		"tier":         tier.TierName,
		"feeAmount":    tier.FixedFee,
		"newBalance":   wallet.Balance,
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// ============================================================================
// DISPUTE FEE SYSTEM FUNCTIONS - Prevents spam/frivolous disputes
// ============================================================================

// InitializeDisputeFeeTiers sets up the dispute fee tier structure
func (p *PlatformContract) InitializeDisputeFeeTiers(ctx contractapi.TransactionContextInterface) (string, error) {
	tiers := []DisputeFeeTier{
		{
			TierID:         "DISPUTE_MINOR",
			TierName:       "MINOR",
			FilingFee:      10.0,
			MinClaimAmount: 0,
			MaxClaimAmount: 500,
			RefundOnWin:    true,
			RefundPercent:  100.0, // Full refund if won
			Description:    "Minor disputes: Claims $0 - $500",
			IsActive:       true,
		},
		{
			TierID:         "DISPUTE_STANDARD",
			TierName:       "STANDARD",
			FilingFee:      25.0,
			MinClaimAmount: 501,
			MaxClaimAmount: 5000,
			RefundOnWin:    true,
			RefundPercent:  100.0,
			Description:    "Standard disputes: Claims $501 - $5,000",
			IsActive:       true,
		},
		{
			TierID:         "DISPUTE_MAJOR",
			TierName:       "MAJOR",
			FilingFee:      50.0,
			MinClaimAmount: 5001,
			MaxClaimAmount: 50000,
			RefundOnWin:    true,
			RefundPercent:  100.0,
			Description:    "Major disputes: Claims $5,001 - $50,000",
			IsActive:       true,
		},
		{
			TierID:         "DISPUTE_CRITICAL",
			TierName:       "CRITICAL",
			FilingFee:      100.0,
			MinClaimAmount: 50001,
			MaxClaimAmount: 999999999,
			RefundOnWin:    true,
			RefundPercent:  100.0,
			Description:    "Critical disputes: Claims $50,001+",
			IsActive:       true,
		},
	}

	for _, tier := range tiers {
		tierKey := fmt.Sprintf("DISPUTE_FEE_TIER_%s", tier.TierID)
		tierJSON, _ := json.Marshal(tier)
		ctx.GetStub().PutState(tierKey, tierJSON)
	}

	// Store tier list
	tierListJSON, _ := json.Marshal([]string{"DISPUTE_MINOR", "DISPUTE_STANDARD", "DISPUTE_MAJOR", "DISPUTE_CRITICAL"})
	ctx.GetStub().PutState("DISPUTE_FEE_TIER_LIST", tierListJSON)

	return `{"message": "Dispute fee tiers initialized successfully", "tiers": 4}`, nil
}

// GetDisputeFeeTier determines the dispute fee tier for a given claim amount
func (p *PlatformContract) GetDisputeFeeTier(ctx contractapi.TransactionContextInterface, claimAmount float64) (string, error) {
	tierListJSON, _ := ctx.GetStub().GetState("DISPUTE_FEE_TIER_LIST")
	if tierListJSON == nil {
		return "", fmt.Errorf("dispute fee tiers not initialized. Call InitializeDisputeFeeTiers first")
	}

	var tierList []string
	json.Unmarshal(tierListJSON, &tierList)

	for _, tierID := range tierList {
		tierKey := fmt.Sprintf("DISPUTE_FEE_TIER_%s", tierID)
		tierJSON, _ := ctx.GetStub().GetState(tierKey)
		if tierJSON == nil {
			continue
		}

		var tier DisputeFeeTier
		json.Unmarshal(tierJSON, &tier)

		if tier.IsActive && claimAmount >= tier.MinClaimAmount && claimAmount <= tier.MaxClaimAmount {
			return string(tierJSON), nil
		}
	}

	// Default to minor tier if no match
	tierKey := "DISPUTE_FEE_TIER_DISPUTE_MINOR"
	tierJSON, _ := ctx.GetStub().GetState(tierKey)
	if tierJSON != nil {
		return string(tierJSON), nil
	}

	return "", fmt.Errorf("no matching dispute fee tier for claim amount %.2f", claimAmount)
}

// CollectDisputeFee collects fee for filing a dispute (prevents spam)
func (p *PlatformContract) CollectDisputeFee(
	ctx contractapi.TransactionContextInterface,
	feeRecordID string,
	disputeID string,
	initiatorID string,
	initiatorType string,
	claimAmount float64,
) (string, error) {
	// Get user's wallet
	userWalletKey := fmt.Sprintf("USER_WALLET_%s_%s", initiatorType, initiatorID)
	walletIDBytes, err := ctx.GetStub().GetState(userWalletKey)
	if err != nil || walletIDBytes == nil {
		return "", fmt.Errorf("no wallet found for user %s. Create wallet first", initiatorID)
	}

	walletID := string(walletIDBytes)
	walletKey := fmt.Sprintf("WALLET_%s", walletID)
	walletJSON, err := ctx.GetStub().GetState(walletKey)
	if err != nil || walletJSON == nil {
		return "", fmt.Errorf("wallet %s not found", walletID)
	}

	var wallet Wallet
	json.Unmarshal(walletJSON, &wallet)

	// Get applicable dispute fee tier
	tierJSON, err := p.GetDisputeFeeTier(ctx, claimAmount)
	if err != nil {
		return "", err
	}

	var tier DisputeFeeTier
	json.Unmarshal([]byte(tierJSON), &tier)

	// Check balance
	if wallet.Balance < tier.FilingFee {
		return "", fmt.Errorf("insufficient balance for dispute filing fee: have %.2f tokens, need %.2f", wallet.Balance, tier.FilingFee)
	}

	now := time.Now().Format(time.RFC3339)

	// Deduct fee from wallet
	wallet.Balance -= tier.FilingFee
	wallet.UpdatedAt = now
	walletUpdatedJSON, _ := json.Marshal(wallet)
	ctx.GetStub().PutState(walletKey, walletUpdatedJSON)

	// Create fee record
	feeRecord := DisputeFeeRecord{
		FeeRecordID:   feeRecordID,
		DisputeID:     disputeID,
		InitiatorID:   initiatorID,
		InitiatorType: initiatorType,
		WalletID:      walletID,
		TierID:        tier.TierID,
		TierName:      tier.TierName,
		FeeAmount:     tier.FilingFee,
		ClaimAmount:   claimAmount,
		Status:        "COLLECTED",
		TransactionID: fmt.Sprintf("DISPUTE_FEE_TXN_%s", feeRecordID),
		CollectedAt:   now,
	}

	feeKey := fmt.Sprintf("DISPUTE_FEE_%s", feeRecordID)
	feeRecordJSON, _ := json.Marshal(feeRecord)
	ctx.GetStub().PutState(feeKey, feeRecordJSON)

	// Record transaction
	txn := TokenTransaction{
		TransactionID:   feeRecord.TransactionID,
		FromWalletID:    walletID,
		ToWalletID:      "PLATFORM_DISPUTE_POOL",
		Amount:          tier.FilingFee,
		TransactionType: "DISPUTE_FEE",
		Reference:       fmt.Sprintf("Dispute filing fee for %s", disputeID),
		Status:          "COMPLETED",
		CreatedAt:       now,
	}
	txnJSON, _ := json.Marshal(txn)
	ctx.GetStub().PutState(fmt.Sprintf("TXN_%s", feeRecord.TransactionID), txnJSON)

	// Emit event
	eventPayload := map[string]interface{}{
		"feeRecordId": feeRecordID,
		"disputeId":   disputeID,
		"initiatorId": initiatorID,
		"feeAmount":   tier.FilingFee,
		"tier":        tier.TierName,
		"action":      "DISPUTE_FEE_COLLECTED",
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("DisputeFeeCollected", eventJSON)

	response := map[string]interface{}{
		"message":     "Dispute filing fee collected",
		"feeRecordId": feeRecordID,
		"disputeId":   disputeID,
		"tier":        tier.TierName,
		"feeAmount":   tier.FilingFee,
		"newBalance":  wallet.Balance,
		"refundable":  tier.RefundOnWin,
		"note":        "Fee will be refunded if dispute is resolved in your favor",
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// ProcessDisputeFeeOutcome processes dispute fee based on resolution (refund or forfeit)
func (p *PlatformContract) ProcessDisputeFeeOutcome(
	ctx contractapi.TransactionContextInterface,
	feeRecordID string,
	disputeID string,
	resolution string, // FAVOR_INITIATOR, FAVOR_RESPONDENT, PARTIAL, DISMISSED
) (string, error) {
	// Get fee record
	feeKey := fmt.Sprintf("DISPUTE_FEE_%s", feeRecordID)
	feeRecordJSON, err := ctx.GetStub().GetState(feeKey)
	if err != nil || feeRecordJSON == nil {
		return "", fmt.Errorf("dispute fee record %s not found", feeRecordID)
	}

	var feeRecord DisputeFeeRecord
	json.Unmarshal(feeRecordJSON, &feeRecord)

	if feeRecord.Status != "COLLECTED" {
		return "", fmt.Errorf("fee already processed: %s", feeRecord.Status)
	}

	// Get tier info
	tierKey := fmt.Sprintf("DISPUTE_FEE_TIER_%s", feeRecord.TierID)
	tierJSON, _ := ctx.GetStub().GetState(tierKey)
	var tier DisputeFeeTier
	if tierJSON != nil {
		json.Unmarshal(tierJSON, &tier)
	}

	now := time.Now().Format(time.RFC3339)
	var refundAmount float64 = 0

	switch resolution {
	case "FAVOR_INITIATOR":
		// Full refund to initiator (they won)
		if tier.RefundOnWin {
			refundAmount = feeRecord.FeeAmount * (tier.RefundPercent / 100)
			feeRecord.Status = "REFUNDED"
		}
	case "PARTIAL":
		// 50% refund for partial resolution
		refundAmount = feeRecord.FeeAmount * 0.5
		feeRecord.Status = "PARTIALLY_REFUNDED"
	case "FAVOR_RESPONDENT", "DISMISSED":
		// Fee forfeited (initiator lost or dispute was frivolous)
		feeRecord.Status = "FORFEITED"
		refundAmount = 0
	default:
		return "", fmt.Errorf("invalid resolution: %s", resolution)
	}

	// Process refund if applicable
	if refundAmount > 0 {
		// Get initiator's wallet
		walletKey := fmt.Sprintf("WALLET_%s", feeRecord.WalletID)
		walletJSON, _ := ctx.GetStub().GetState(walletKey)
		if walletJSON != nil {
			var wallet Wallet
			json.Unmarshal(walletJSON, &wallet)
			wallet.Balance += refundAmount
			wallet.UpdatedAt = now
			walletUpdatedJSON, _ := json.Marshal(wallet)
			ctx.GetStub().PutState(walletKey, walletUpdatedJSON)

			// Record refund transaction
			refundTxnID := fmt.Sprintf("DISPUTE_REFUND_TXN_%s", feeRecordID)
			refundTxn := TokenTransaction{
				TransactionID:   refundTxnID,
				FromWalletID:    "PLATFORM_DISPUTE_POOL",
				ToWalletID:      feeRecord.WalletID,
				Amount:          refundAmount,
				TransactionType: "DISPUTE_FEE_REFUND",
				Reference:       fmt.Sprintf("Dispute fee refund for %s - %s", disputeID, resolution),
				Status:          "COMPLETED",
				CreatedAt:       now,
			}
			refundTxnJSON, _ := json.Marshal(refundTxn)
			ctx.GetStub().PutState(fmt.Sprintf("TXN_%s", refundTxnID), refundTxnJSON)

			feeRecord.RefundTxnID = refundTxnID
		}
	}

	feeRecord.RefundAmount = refundAmount
	feeRecord.ProcessedAt = now

	feeRecordUpdatedJSON, _ := json.Marshal(feeRecord)
	ctx.GetStub().PutState(feeKey, feeRecordUpdatedJSON)

	// Emit event
	eventPayload := map[string]interface{}{
		"feeRecordId":  feeRecordID,
		"disputeId":    disputeID,
		"resolution":   resolution,
		"status":       feeRecord.Status,
		"refundAmount": refundAmount,
		"action":       "DISPUTE_FEE_PROCESSED",
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("DisputeFeeProcessed", eventJSON)

	response := map[string]interface{}{
		"message":      "Dispute fee processed",
		"feeRecordId":  feeRecordID,
		"disputeId":    disputeID,
		"resolution":   resolution,
		"status":       feeRecord.Status,
		"originalFee":  feeRecord.FeeAmount,
		"refundAmount": refundAmount,
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// GetDisputeFeeRecord retrieves a dispute fee record
func (p *PlatformContract) GetDisputeFeeRecord(ctx contractapi.TransactionContextInterface, feeRecordID string) (string, error) {
	feeKey := fmt.Sprintf("DISPUTE_FEE_%s", feeRecordID)
	feeRecordJSON, err := ctx.GetStub().GetState(feeKey)
	if err != nil || feeRecordJSON == nil {
		return "", fmt.Errorf("dispute fee record %s not found", feeRecordID)
	}
	return string(feeRecordJSON), nil
}

// ============================================================================
// ML RATING SYSTEM FUNCTIONS
// ============================================================================

// RecordRating records a rating from ML model or user (append-only)
func (p *PlatformContract) RecordRating(
	ctx contractapi.TransactionContextInterface,
	ratingID string,
	targetUserID string,
	targetUserType string,
	raterType string,
	raterID string,
	score float64,
	category string,
	factors string,
	evidenceHash string,
) (string, error) {
	// Validate score
	if score < 0 || score > 100 {
		return "", fmt.Errorf("score must be between 0 and 100")
	}

	// Validate category
	validCategories := map[string]bool{
		"CREDIBILITY": true, "RISK": true, "PERFORMANCE": true, "COMPLIANCE": true, "OVERALL": true,
	}
	if !validCategories[category] {
		return "", fmt.Errorf("invalid category: %s", category)
	}

	now := time.Now().Format(time.RFC3339)

	rating := Rating{
		RatingID:       ratingID,
		TargetUserID:   targetUserID,
		TargetUserType: targetUserType,
		RaterType:      raterType,
		RaterID:        raterID,
		Score:          score,
		Category:       category,
		Factors:        factors,
		EvidenceHash:   evidenceHash,
		CreatedAt:      now,
	}

	// Store rating (append-only)
	ratingKey := fmt.Sprintf("RATING_%s", ratingID)
	ratingJSON, _ := json.Marshal(rating)
	ctx.GetStub().PutState(ratingKey, ratingJSON)

	// Update aggregate rating
	p.updateRatingAggregate(ctx, targetUserID, targetUserType, category, score)

	response := map[string]interface{}{
		"message":    "Rating recorded",
		"ratingId":   ratingID,
		"targetUser": targetUserID,
		"category":   category,
		"score":      score,
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// updateRatingAggregate updates the aggregate rating for a user
func (p *PlatformContract) updateRatingAggregate(
	ctx contractapi.TransactionContextInterface,
	userID string,
	userType string,
	category string,
	newScore float64,
) {
	aggKey := fmt.Sprintf("RATING_AGG_%s_%s", userType, userID)
	aggJSON, _ := ctx.GetStub().GetState(aggKey)

	var agg RatingAggregate
	if aggJSON != nil {
		json.Unmarshal(aggJSON, &agg)
	} else {
		agg = RatingAggregate{
			UserID:   userID,
			UserType: userType,
		}
	}

	agg.TotalRatings++
	agg.LastUpdated = time.Now().Format(time.RFC3339)

	// Update specific category score (simple moving average)
	switch category {
	case "CREDIBILITY":
		agg.CredibilityScore = ((agg.CredibilityScore * float64(agg.TotalRatings-1)) + newScore) / float64(agg.TotalRatings)
	case "RISK":
		agg.RiskScore = ((agg.RiskScore * float64(agg.TotalRatings-1)) + newScore) / float64(agg.TotalRatings)
	case "PERFORMANCE":
		agg.PerformanceScore = ((agg.PerformanceScore * float64(agg.TotalRatings-1)) + newScore) / float64(agg.TotalRatings)
	case "COMPLIANCE":
		agg.ComplianceScore = ((agg.ComplianceScore * float64(agg.TotalRatings-1)) + newScore) / float64(agg.TotalRatings)
	case "OVERALL":
		agg.OverallScore = ((agg.OverallScore * float64(agg.TotalRatings-1)) + newScore) / float64(agg.TotalRatings)
	}

	aggUpdatedJSON, _ := json.Marshal(agg)
	ctx.GetStub().PutState(aggKey, aggUpdatedJSON)
}

// GetRatingAggregate retrieves aggregate rating for a user
func (p *PlatformContract) GetRatingAggregate(ctx contractapi.TransactionContextInterface, userType string, userID string) (string, error) {
	aggKey := fmt.Sprintf("RATING_AGG_%s_%s", userType, userID)
	aggJSON, err := ctx.GetStub().GetState(aggKey)
	if err != nil || aggJSON == nil {
		return "", fmt.Errorf("no ratings found for user %s", userID)
	}
	return string(aggJSON), nil
}

// ============================================================================
// DISPUTE SYSTEM FUNCTIONS
// ============================================================================

// CreateDispute creates a new dispute ticket
// Channel: common-channel
func (p *PlatformContract) CreateDispute(
	ctx contractapi.TransactionContextInterface,
	disputeID string,
	disputeType string,
	disputeSubType string,
	initiatorID string,
	initiatorType string,
	respondentID string,
	respondentType string,
	campaignID string,
	agreementID string,
	title string,
	description string,
	claimedAmount float64,
	evidenceHashesJSON string,
) (string, error) {
	// Check if initiator is suspended/blacklisted
	reputationKey := fmt.Sprintf("REPUTATION_%s_%s", initiatorType, initiatorID)
	repJSON, _ := ctx.GetStub().GetState(reputationKey)
	if repJSON != nil {
		var rep Reputation
		json.Unmarshal(repJSON, &rep)
		if rep.Status == "BLACKLISTED" {
			return "", fmt.Errorf("user %s is blacklisted and cannot create disputes", initiatorID)
		}
		if rep.Status == "SUSPENDED" {
			return "", fmt.Errorf("user %s is suspended and cannot create disputes", initiatorID)
		}
	}

	// ========== DISPUTE FILING FEE COLLECTION ==========
	// Collect dispute filing fee to prevent spam/frivolous disputes
	feeID := fmt.Sprintf("DFEE_%s", disputeID)
	feeResult, err := p.CollectDisputeFee(ctx, feeID, disputeID, initiatorID, initiatorType, claimedAmount)
	if err != nil {
		return "", fmt.Errorf("failed to collect dispute filing fee: %v. Ensure wallet has sufficient balance", err)
	}

	// Parse fee result to get fee amount for response
	var feeInfo map[string]interface{}
	json.Unmarshal([]byte(feeResult), &feeInfo)
	feeAmount := feeInfo["feeAmount"]

	now := time.Now().Format(time.RFC3339)
	ticketNumber := fmt.Sprintf("DISP-%d", time.Now().UnixNano()%1000000)

	// Parse evidence hashes
	var evidences []Evidence
	if evidenceHashesJSON != "" {
		json.Unmarshal([]byte(evidenceHashesJSON), &evidences)
	}

	dispute := Dispute{
		DisputeID:          disputeID,
		TicketNumber:       ticketNumber,
		DisputeType:        disputeType,
		DisputeSubType:     disputeSubType,
		InitiatorID:        initiatorID,
		InitiatorType:      initiatorType,
		RespondentID:       respondentID,
		RespondentType:     respondentType,
		AdditionalParties:  []DisputeParty{},
		CampaignID:         campaignID,
		AgreementID:        agreementID,
		Title:              title,
		Description:        description,
		ClaimedAmount:      claimedAmount,
		EvidenceHashes:     evidences,
		InvestigationNotes: []InvestigationNote{},
		SmartContractRules: []string{},
		VotingEnabled:      false,
		Status:             "OPEN",
		PenaltiesApplied:   []Penalty{},
		RefundsOrdered:     []RefundOrder{},
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	disputeKey := fmt.Sprintf("DISPUTE_%s", disputeID)
	disputeJSON, _ := json.Marshal(dispute)
	ctx.GetStub().PutState(disputeKey, disputeJSON)

	// Update dispute count for both parties
	p.incrementDisputeCount(ctx, initiatorType, initiatorID)
	p.incrementDisputeCount(ctx, respondentType, respondentID)

	// Emit event
	eventPayload := map[string]interface{}{
		"disputeId":    disputeID,
		"ticketNumber": ticketNumber,
		"disputeType":  disputeType,
		"initiator":    initiatorID,
		"respondent":   respondentID,
		"status":       "OPEN",
		"filingFee":    feeAmount,
		"feeStatus":    "LOCKED",
		"channel":      "common-channel",
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("DisputeCreated", eventJSON)

	response := map[string]interface{}{
		"message":      "Dispute created successfully",
		"disputeId":    disputeID,
		"ticketNumber": ticketNumber,
		"status":       "OPEN",
		"filingFee":    feeAmount,
		"feeNote":      "Filing fee is locked. Will be refunded if dispute resolved in your favor, forfeited otherwise.",
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// incrementDisputeCount increments dispute count in reputation
func (p *PlatformContract) incrementDisputeCount(ctx contractapi.TransactionContextInterface, userType string, userID string) {
	reputationKey := fmt.Sprintf("REPUTATION_%s_%s", userType, userID)
	repJSON, _ := ctx.GetStub().GetState(reputationKey)
	if repJSON != nil {
		var rep Reputation
		json.Unmarshal(repJSON, &rep)
		rep.TotalDisputes++
		rep.UpdatedAt = time.Now().Format(time.RFC3339)
		repUpdatedJSON, _ := json.Marshal(rep)
		ctx.GetStub().PutState(reputationKey, repUpdatedJSON)
	}
}

// SubmitEvidence allows parties to submit evidence for a dispute
func (p *PlatformContract) SubmitEvidence(
	ctx contractapi.TransactionContextInterface,
	disputeID string,
	evidenceID string,
	submittedBy string,
	submitterType string,
	ipfsHash string,
	description string,
	evidenceType string,
) (string, error) {
	disputeKey := fmt.Sprintf("DISPUTE_%s", disputeID)
	disputeJSON, err := ctx.GetStub().GetState(disputeKey)
	if err != nil || disputeJSON == nil {
		return "", fmt.Errorf("dispute %s not found", disputeID)
	}

	var dispute Dispute
	json.Unmarshal(disputeJSON, &dispute)

	if dispute.Status == "RESOLVED" || dispute.Status == "CLOSED" {
		return "", fmt.Errorf("cannot submit evidence to resolved/closed dispute")
	}

	now := time.Now().Format(time.RFC3339)
	evidence := Evidence{
		EvidenceID:    evidenceID,
		SubmittedBy:   submittedBy,
		SubmitterType: submitterType,
		IPFSHash:      ipfsHash,
		Description:   description,
		EvidenceType:  evidenceType,
		SubmittedAt:   now,
	}

	dispute.EvidenceHashes = append(dispute.EvidenceHashes, evidence)
	dispute.UpdatedAt = now

	disputeUpdatedJSON, _ := json.Marshal(dispute)
	ctx.GetStub().PutState(disputeKey, disputeUpdatedJSON)

	response := map[string]interface{}{
		"message":    "Evidence submitted",
		"disputeId":  disputeID,
		"evidenceId": evidenceID,
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// AssignInvestigator assigns a validator to investigate the dispute
func (p *PlatformContract) AssignInvestigator(
	ctx contractapi.TransactionContextInterface,
	disputeID string,
	investigatorID string,
) (string, error) {
	disputeKey := fmt.Sprintf("DISPUTE_%s", disputeID)
	disputeJSON, err := ctx.GetStub().GetState(disputeKey)
	if err != nil || disputeJSON == nil {
		return "", fmt.Errorf("dispute %s not found", disputeID)
	}

	var dispute Dispute
	json.Unmarshal(disputeJSON, &dispute)

	now := time.Now().Format(time.RFC3339)
	dispute.InvestigatorID = investigatorID
	dispute.Status = "UNDER_INVESTIGATION"
	dispute.UpdatedAt = now

	disputeUpdatedJSON, _ := json.Marshal(dispute)
	ctx.GetStub().PutState(disputeKey, disputeUpdatedJSON)

	response := map[string]interface{}{
		"message":        "Investigator assigned",
		"disputeId":      disputeID,
		"investigatorId": investigatorID,
		"status":         "UNDER_INVESTIGATION",
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// AddInvestigationNote adds a note during investigation
func (p *PlatformContract) AddInvestigationNote(
	ctx contractapi.TransactionContextInterface,
	disputeID string,
	noteID string,
	investigatorID string,
	note string,
	findingType string,
) (string, error) {
	disputeKey := fmt.Sprintf("DISPUTE_%s", disputeID)
	disputeJSON, err := ctx.GetStub().GetState(disputeKey)
	if err != nil || disputeJSON == nil {
		return "", fmt.Errorf("dispute %s not found", disputeID)
	}

	var dispute Dispute
	json.Unmarshal(disputeJSON, &dispute)

	if dispute.InvestigatorID != investigatorID {
		return "", fmt.Errorf("only assigned investigator can add notes")
	}

	now := time.Now().Format(time.RFC3339)
	invNote := InvestigationNote{
		NoteID:         noteID,
		InvestigatorID: investigatorID,
		Note:           note,
		FindingType:    findingType,
		CreatedAt:      now,
	}

	dispute.InvestigationNotes = append(dispute.InvestigationNotes, invNote)
	dispute.UpdatedAt = now

	disputeUpdatedJSON, _ := json.Marshal(dispute)
	ctx.GetStub().PutState(disputeKey, disputeUpdatedJSON)

	return `{"message": "Investigation note added"}`, nil
}

// EnableVoting enables anonymous voting for dispute resolution
func (p *PlatformContract) EnableVoting(
	ctx contractapi.TransactionContextInterface,
	disputeID string,
	votingDeadlineDays int,
) (string, error) {
	disputeKey := fmt.Sprintf("DISPUTE_%s", disputeID)
	disputeJSON, err := ctx.GetStub().GetState(disputeKey)
	if err != nil || disputeJSON == nil {
		return "", fmt.Errorf("dispute %s not found", disputeID)
	}

	var dispute Dispute
	json.Unmarshal(disputeJSON, &dispute)

	now := time.Now()
	deadline := now.AddDate(0, 0, votingDeadlineDays)

	dispute.VotingEnabled = true
	dispute.VotingDeadline = deadline.Format(time.RFC3339)
	dispute.Status = "VOTING"
	dispute.UpdatedAt = now.Format(time.RFC3339)

	disputeUpdatedJSON, _ := json.Marshal(dispute)
	ctx.GetStub().PutState(disputeKey, disputeUpdatedJSON)

	response := map[string]interface{}{
		"message":        "Voting enabled",
		"disputeId":      disputeID,
		"votingDeadline": dispute.VotingDeadline,
		"status":         "VOTING",
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// ============================================================================
// ANONYMOUS VOTING FUNCTIONS (Commit-Reveal Scheme)
// ============================================================================

// CommitVote commits a vote (hash only - anonymous)
func (p *PlatformContract) CommitVote(
	ctx contractapi.TransactionContextInterface,
	commitmentID string,
	disputeID string,
	voterHash string,
	voteHash string,
) (string, error) {
	// Verify dispute is in voting phase
	disputeKey := fmt.Sprintf("DISPUTE_%s", disputeID)
	disputeJSON, err := ctx.GetStub().GetState(disputeKey)
	if err != nil || disputeJSON == nil {
		return "", fmt.Errorf("dispute %s not found", disputeID)
	}

	var dispute Dispute
	json.Unmarshal(disputeJSON, &dispute)

	if dispute.Status != "VOTING" {
		return "", fmt.Errorf("dispute is not in voting phase")
	}

	// Check voting deadline
	deadline, _ := time.Parse(time.RFC3339, dispute.VotingDeadline)
	if time.Now().After(deadline) {
		return "", fmt.Errorf("voting deadline has passed")
	}

	// Check if voter already committed
	existingKey := fmt.Sprintf("VOTE_COMMIT_%s_%s", disputeID, voterHash)
	existing, _ := ctx.GetStub().GetState(existingKey)
	if existing != nil {
		return "", fmt.Errorf("voter has already committed a vote")
	}

	now := time.Now().Format(time.RFC3339)
	commitment := VoteCommitment{
		CommitmentID: commitmentID,
		DisputeID:    disputeID,
		VoterHash:    voterHash,
		VoteHash:     voteHash,
		CommittedAt:  now,
	}

	commitKey := fmt.Sprintf("VOTE_COMMIT_%s", commitmentID)
	commitJSON, _ := json.Marshal(commitment)
	ctx.GetStub().PutState(commitKey, commitJSON)
	ctx.GetStub().PutState(existingKey, []byte(commitmentID))

	response := map[string]interface{}{
		"message":      "Vote committed anonymously",
		"commitmentId": commitmentID,
		"disputeId":    disputeID,
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// RevealVote reveals a previously committed vote
func (p *PlatformContract) RevealVote(
	ctx contractapi.TransactionContextInterface,
	revealID string,
	disputeID string,
	commitmentID string,
	voterOrgType string,
	vote string,
	salt string,
) (string, error) {
	// Get commitment
	commitKey := fmt.Sprintf("VOTE_COMMIT_%s", commitmentID)
	commitJSON, err := ctx.GetStub().GetState(commitKey)
	if err != nil || commitJSON == nil {
		return "", fmt.Errorf("commitment %s not found", commitmentID)
	}

	var commitment VoteCommitment
	json.Unmarshal(commitJSON, &commitment)

	// Verify the reveal matches the commitment
	expectedHash := generateHash(vote + salt)
	isValid := expectedHash == commitment.VoteHash

	now := time.Now().Format(time.RFC3339)
	reveal := VoteReveal{
		RevealID:     revealID,
		DisputeID:    disputeID,
		CommitmentID: commitmentID,
		VoterOrgType: voterOrgType,
		Vote:         vote,
		Salt:         salt,
		RevealedAt:   now,
		IsValid:      isValid,
	}

	revealKey := fmt.Sprintf("VOTE_REVEAL_%s", revealID)
	revealJSON, _ := json.Marshal(reveal)
	ctx.GetStub().PutState(revealKey, revealJSON)

	response := map[string]interface{}{
		"message":  "Vote revealed",
		"revealId": revealID,
		"isValid":  isValid,
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// TallyVotes tallies all votes for a dispute
func (p *PlatformContract) TallyVotes(
	ctx contractapi.TransactionContextInterface,
	disputeID string,
) (string, error) {
	disputeKey := fmt.Sprintf("DISPUTE_%s", disputeID)
	disputeJSON, err := ctx.GetStub().GetState(disputeKey)
	if err != nil || disputeJSON == nil {
		return "", fmt.Errorf("dispute %s not found", disputeID)
	}

	var dispute Dispute
	json.Unmarshal(disputeJSON, &dispute)

	// Query all reveals for this dispute
	queryString := fmt.Sprintf(`{"selector":{"disputeId":"%s"}}`, disputeID)
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		// Fallback: count manually using key prefix iteration
		// For simplicity, we'll use a counter approach
	}
	if resultsIterator != nil {
		defer resultsIterator.Close()
	}

	result := VotingResult{
		DisputeID:       disputeID,
		TotalVotes:      0,
		FavorInitiator:  0,
		FavorRespondent: 0,
		Abstained:       0,
		InvalidVotes:    0,
		TalliedAt:       time.Now().Format(time.RFC3339),
	}

	// Iterate through reveals (simplified)
	if resultsIterator != nil {
		for resultsIterator.HasNext() {
			queryResponse, err := resultsIterator.Next()
			if err != nil {
				continue
			}

			var reveal VoteReveal
			json.Unmarshal(queryResponse.Value, &reveal)

			result.TotalVotes++
			if !reveal.IsValid {
				result.InvalidVotes++
				continue
			}

			switch reveal.Vote {
			case "FAVOR_INITIATOR":
				result.FavorInitiator++
			case "FAVOR_RESPONDENT":
				result.FavorRespondent++
			case "ABSTAIN":
				result.Abstained++
			}
		}
	}

	// Determine outcome
	if result.FavorInitiator > result.FavorRespondent {
		result.Outcome = "FAVOR_INITIATOR"
	} else if result.FavorRespondent > result.FavorInitiator {
		result.Outcome = "FAVOR_RESPONDENT"
	} else {
		result.Outcome = "TIE"
	}

	dispute.VotingResult = &result
	dispute.UpdatedAt = time.Now().Format(time.RFC3339)

	disputeUpdatedJSON, _ := json.Marshal(dispute)
	ctx.GetStub().PutState(disputeKey, disputeUpdatedJSON)

	resultJSON, _ := json.Marshal(result)
	return string(resultJSON), nil
}

// ============================================================================
// DISPUTE RESOLUTION & PENALTY FUNCTIONS
// ============================================================================

// ResolveDispute resolves a dispute and applies penalties/refunds
func (p *PlatformContract) ResolveDispute(
	ctx contractapi.TransactionContextInterface,
	disputeID string,
	resolution string,
	resolutionNotes string,
	penaltiesJSON string,
	refundsJSON string,
) (string, error) {
	disputeKey := fmt.Sprintf("DISPUTE_%s", disputeID)
	disputeJSON, err := ctx.GetStub().GetState(disputeKey)
	if err != nil || disputeJSON == nil {
		return "", fmt.Errorf("dispute %s not found", disputeID)
	}

	var dispute Dispute
	json.Unmarshal(disputeJSON, &dispute)

	now := time.Now().Format(time.RFC3339)

	// Parse and apply penalties
	var penalties []Penalty
	if penaltiesJSON != "" {
		json.Unmarshal([]byte(penaltiesJSON), &penalties)
		for i := range penalties {
			penalties[i].DisputeID = disputeID
			penalties[i].Status = "APPLIED"
			penalties[i].AppliedAt = now

			// Apply penalty to user
			p.applyPenalty(ctx, &penalties[i])
		}
	}

	// Parse and process refunds
	var refunds []RefundOrder
	if refundsJSON != "" {
		json.Unmarshal([]byte(refundsJSON), &refunds)
		for i := range refunds {
			refunds[i].DisputeID = disputeID
			refunds[i].Status = "PROCESSED"
			refunds[i].ProcessedAt = now

			// Process refund
			p.processRefund(ctx, &refunds[i])
		}
	}

	// Update dispute
	dispute.Resolution = resolution
	dispute.ResolutionNotes = resolutionNotes
	dispute.PenaltiesApplied = penalties
	dispute.RefundsOrdered = refunds
	dispute.Status = "RESOLVED"
	dispute.ResolvedAt = now
	dispute.UpdatedAt = now

	// Update reputation for parties based on resolution
	if resolution == "FAVOR_INITIATOR" {
		p.updateReputationAfterDispute(ctx, dispute.InitiatorType, dispute.InitiatorID, true)
		p.updateReputationAfterDispute(ctx, dispute.RespondentType, dispute.RespondentID, false)
	} else if resolution == "FAVOR_RESPONDENT" {
		p.updateReputationAfterDispute(ctx, dispute.InitiatorType, dispute.InitiatorID, false)
		p.updateReputationAfterDispute(ctx, dispute.RespondentType, dispute.RespondentID, true)
	}

	// Process dispute filing fee resolution (refund or forfeit)
	feeID := fmt.Sprintf("DFEE_%s", disputeID)
	feeResult, feeErr := p.ProcessDisputeFeeOutcome(ctx, feeID, disputeID, resolution)
	feeOutcome := ""
	if feeErr == nil {
		var feeResolution map[string]interface{}
		json.Unmarshal([]byte(feeResult), &feeResolution)
		feeOutcome = fmt.Sprintf("%v", feeResolution["action"])
	}

	disputeUpdatedJSON, _ := json.Marshal(dispute)
	ctx.GetStub().PutState(disputeKey, disputeUpdatedJSON)

	// Emit event
	eventPayload := map[string]interface{}{
		"disputeId":   disputeID,
		"resolution":  resolution,
		"status":      "RESOLVED",
		"penalties":   len(penalties),
		"refunds":     len(refunds),
		"filingFee":   feeOutcome,
		"channel":     "common-channel",
	}
	eventJSON, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("DisputeResolved", eventJSON)

	response := map[string]interface{}{
		"message":       "Dispute resolved",
		"disputeId":     disputeID,
		"resolution":    resolution,
		"penalties":     len(penalties),
		"refunds":       len(refunds),
		"filingFeeNote": fmt.Sprintf("Filing fee %s", feeOutcome),
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// applyPenalty applies a penalty to a user
func (p *PlatformContract) applyPenalty(ctx contractapi.TransactionContextInterface, penalty *Penalty) {
	now := time.Now().Format(time.RFC3339)

	// If token fine, deduct from wallet
	if penalty.TokenAmount > 0 {
		userWalletKey := fmt.Sprintf("USER_WALLET_%s_%s", penalty.UserType, penalty.UserID)
		walletIDBytes, _ := ctx.GetStub().GetState(userWalletKey)
		if walletIDBytes != nil {
			walletKey := fmt.Sprintf("WALLET_%s", string(walletIDBytes))
			walletJSON, _ := ctx.GetStub().GetState(walletKey)
			if walletJSON != nil {
				var wallet Wallet
				json.Unmarshal(walletJSON, &wallet)
				wallet.Balance -= penalty.TokenAmount
				if wallet.Balance < 0 {
					wallet.Balance = 0
				}
				wallet.TotalPenalties += penalty.TokenAmount
				wallet.UpdatedAt = now
				walletUpdatedJSON, _ := json.Marshal(wallet)
				ctx.GetStub().PutState(walletKey, walletUpdatedJSON)
			}
		}
	}

	// Update reputation
	reputationKey := fmt.Sprintf("REPUTATION_%s_%s", penalty.UserType, penalty.UserID)
	repJSON, _ := ctx.GetStub().GetState(reputationKey)
	if repJSON != nil {
		var rep Reputation
		json.Unmarshal(repJSON, &rep)

		rep.CurrentScore -= penalty.ReputationDeduct
		if rep.CurrentScore < MinReputationScore {
			rep.CurrentScore = MinReputationScore
		}
		rep.TotalPenalties++
		rep.ConsecutivePenalties++
		rep.UpdatedAt = now

		// Check for auto-suspension/blacklist
		if rep.CurrentScore <= ReputationBlacklistThreshold {
			rep.Status = "BLACKLISTED"
			rep.BlacklistedAt = now
			rep.BlacklistReason = penalty.Description
		} else if rep.CurrentScore <= ReputationSuspensionThreshold || rep.ConsecutivePenalties >= ConsecutivePenaltyThreshold {
			rep.Status = "SUSPENDED"
			suspendedUntil := time.Now().AddDate(0, 0, 30) // 30-day suspension
			rep.SuspendedUntil = suspendedUntil.Format(time.RFC3339)
		}

		repUpdatedJSON, _ := json.Marshal(rep)
		ctx.GetStub().PutState(reputationKey, repUpdatedJSON)

		// Record history
		historyKey := fmt.Sprintf("REP_HISTORY_%s_%d", rep.UserID, time.Now().UnixNano())
		history := ReputationHistory{
			HistoryID:     historyKey,
			UserID:        rep.UserID,
			UserType:      rep.UserType,
			PreviousScore: rep.CurrentScore + penalty.ReputationDeduct,
			NewScore:      rep.CurrentScore,
			ChangeAmount:  -penalty.ReputationDeduct,
			ChangeReason:  penalty.Description,
			DisputeID:     penalty.DisputeID,
			CreatedAt:     now,
		}
		historyJSON, _ := json.Marshal(history)
		ctx.GetStub().PutState(historyKey, historyJSON)
	}

	// Store penalty record
	penaltyKey := fmt.Sprintf("PENALTY_%s", penalty.PenaltyID)
	penaltyJSON, _ := json.Marshal(penalty)
	ctx.GetStub().PutState(penaltyKey, penaltyJSON)
}

// processRefund processes a refund order
func (p *PlatformContract) processRefund(ctx contractapi.TransactionContextInterface, refund *RefundOrder) {
	now := time.Now().Format(time.RFC3339)

	// Calculate net amount after deduction
	deduction := refund.Amount * (refund.DeductionPercent / 100)
	refund.NetAmount = refund.Amount - deduction

	// Get source wallet
	fromWalletKey := fmt.Sprintf("USER_WALLET_%s_%s", refund.FromUserType, refund.FromUserID)
	fromWalletIDBytes, _ := ctx.GetStub().GetState(fromWalletKey)
	if fromWalletIDBytes == nil {
		refund.Status = "FAILED"
		return
	}

	fromWalletID := string(fromWalletIDBytes)
	fromKey := fmt.Sprintf("WALLET_%s", fromWalletID)
	fromJSON, _ := ctx.GetStub().GetState(fromKey)
	if fromJSON == nil {
		refund.Status = "FAILED"
		return
	}

	var fromWallet Wallet
	json.Unmarshal(fromJSON, &fromWallet)

	// Get destination wallet
	toWalletKey := fmt.Sprintf("USER_WALLET_%s_%s", refund.ToUserType, refund.ToUserID)
	toWalletIDBytes, _ := ctx.GetStub().GetState(toWalletKey)
	if toWalletIDBytes == nil {
		refund.Status = "FAILED"
		return
	}

	toWalletID := string(toWalletIDBytes)
	toKey := fmt.Sprintf("WALLET_%s", toWalletID)
	toJSON, _ := ctx.GetStub().GetState(toKey)
	if toJSON == nil {
		refund.Status = "FAILED"
		return
	}

	var toWallet Wallet
	json.Unmarshal(toJSON, &toWallet)

	// Transfer funds
	if fromWallet.Balance >= refund.NetAmount {
		fromWallet.Balance -= refund.NetAmount
		toWallet.Balance += refund.NetAmount
		fromWallet.UpdatedAt = now
		toWallet.UpdatedAt = now

		fromUpdatedJSON, _ := json.Marshal(fromWallet)
		ctx.GetStub().PutState(fromKey, fromUpdatedJSON)

		toUpdatedJSON, _ := json.Marshal(toWallet)
		ctx.GetStub().PutState(toKey, toUpdatedJSON)

		refund.Status = "PROCESSED"
	} else {
		refund.Status = "FAILED"
	}

	// Store refund record
	refundKey := fmt.Sprintf("REFUND_%s", refund.RefundOrderID)
	refundJSON, _ := json.Marshal(refund)
	ctx.GetStub().PutState(refundKey, refundJSON)
}

// updateReputationAfterDispute updates reputation based on dispute outcome
func (p *PlatformContract) updateReputationAfterDispute(ctx contractapi.TransactionContextInterface, userType string, userID string, won bool) {
	reputationKey := fmt.Sprintf("REPUTATION_%s_%s", userType, userID)
	repJSON, _ := ctx.GetStub().GetState(reputationKey)
	if repJSON == nil {
		return
	}

	var rep Reputation
	json.Unmarshal(repJSON, &rep)

	now := time.Now().Format(time.RFC3339)
	previousScore := rep.CurrentScore

	if won {
		rep.DisputesWon++
		rep.ConsecutivePenalties = 0 // Reset consecutive penalties on win
		rep.CurrentScore += 2       // Small boost for winning
		if rep.CurrentScore > MaxReputationScore {
			rep.CurrentScore = MaxReputationScore
		}
	} else {
		rep.DisputesLost++
	}
	rep.UpdatedAt = now

	repUpdatedJSON, _ := json.Marshal(rep)
	ctx.GetStub().PutState(reputationKey, repUpdatedJSON)

	// Record history
	historyKey := fmt.Sprintf("REP_HISTORY_%s_%d", userID, time.Now().UnixNano())
	changeReason := "Dispute lost"
	if won {
		changeReason = "Dispute won"
	}
	history := ReputationHistory{
		HistoryID:     historyKey,
		UserID:        userID,
		UserType:      userType,
		PreviousScore: previousScore,
		NewScore:      rep.CurrentScore,
		ChangeAmount:  rep.CurrentScore - previousScore,
		ChangeReason:  changeReason,
		CreatedAt:     now,
	}
	historyJSON, _ := json.Marshal(history)
	ctx.GetStub().PutState(historyKey, historyJSON)
}

// GetDispute retrieves a dispute by ID
func (p *PlatformContract) GetDispute(ctx contractapi.TransactionContextInterface, disputeID string) (string, error) {
	disputeKey := fmt.Sprintf("DISPUTE_%s", disputeID)
	disputeJSON, err := ctx.GetStub().GetState(disputeKey)
	if err != nil || disputeJSON == nil {
		return "", fmt.Errorf("dispute %s not found", disputeID)
	}
	return string(disputeJSON), nil
}

// GetReputation retrieves reputation for a user
func (p *PlatformContract) GetReputation(ctx contractapi.TransactionContextInterface, userType string, userID string) (string, error) {
	reputationKey := fmt.Sprintf("REPUTATION_%s_%s", userType, userID)
	repJSON, err := ctx.GetStub().GetState(reputationKey)
	if err != nil || repJSON == nil {
		return "", fmt.Errorf("reputation not found for user %s", userID)
	}
	return string(repJSON), nil
}

// CheckUserStatus checks if user is active, suspended, or blacklisted
func (p *PlatformContract) CheckUserStatus(ctx contractapi.TransactionContextInterface, userType string, userID string) (string, error) {
	reputationKey := fmt.Sprintf("REPUTATION_%s_%s", userType, userID)
	repJSON, err := ctx.GetStub().GetState(reputationKey)
	if err != nil || repJSON == nil {
		return `{"status": "ACTIVE", "message": "No reputation record found, assuming active"}`, nil
	}

	var rep Reputation
	json.Unmarshal(repJSON, &rep)

	// Check if suspension has expired
	if rep.Status == "SUSPENDED" && rep.SuspendedUntil != "" {
		suspendedUntil, _ := time.Parse(time.RFC3339, rep.SuspendedUntil)
		if time.Now().After(suspendedUntil) {
			rep.Status = "ACTIVE"
			rep.SuspendedUntil = ""
			rep.ConsecutivePenalties = 0
			rep.UpdatedAt = time.Now().Format(time.RFC3339)
			repUpdatedJSON, _ := json.Marshal(rep)
			ctx.GetStub().PutState(reputationKey, repUpdatedJSON)
		}
	}

	response := map[string]interface{}{
		"userId":           userID,
		"userType":         userType,
		"status":           rep.Status,
		"currentScore":     rep.CurrentScore,
		"suspendedUntil":   rep.SuspendedUntil,
		"blacklistedAt":    rep.BlacklistedAt,
		"blacklistReason":  rep.BlacklistReason,
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// RequestRefund allows investor to request refund (with deduction)
func (p *PlatformContract) RequestRefund(
	ctx contractapi.TransactionContextInterface,
	refundOrderID string,
	investorID string,
	startupID string,
	campaignID string,
	amount float64,
	reason string,
) (string, error) {
	// Determine deduction percentage based on reason (15-30%)
	var deductionPercent float64
	switch reason {
	case "EARLY_WITHDRAWAL":
		deductionPercent = 15.0
	case "MID_AGREEMENT_WITHDRAWAL":
		deductionPercent = 25.0
	case "LATE_WITHDRAWAL":
		deductionPercent = 30.0
	default:
		deductionPercent = 20.0 // Default
	}

	refund := RefundOrder{
		RefundOrderID:    refundOrderID,
		FromUserID:       startupID,
		FromUserType:     "STARTUP",
		ToUserID:         investorID,
		ToUserType:       "INVESTOR",
		Amount:           amount,
		DeductionPercent: deductionPercent,
		Reason:           reason,
		Status:           "PENDING",
	}

	refundKey := fmt.Sprintf("REFUND_REQUEST_%s", refundOrderID)
	refundJSON, _ := json.Marshal(refund)
	ctx.GetStub().PutState(refundKey, refundJSON)

	response := map[string]interface{}{
		"message":          "Refund request created",
		"refundOrderId":    refundOrderID,
		"amount":           amount,
		"deductionPercent": deductionPercent,
		"netAmount":        amount * (1 - deductionPercent/100),
		"status":           "PENDING",
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// ========== DISPUTE FEE QUERIES ==========

// GetDisputeFee retrieves dispute fee details by fee ID
func (p *PlatformContract) GetDisputeFee(ctx contractapi.TransactionContextInterface, feeID string) (string, error) {
	feeKey := fmt.Sprintf("DISPUTE_FEE_%s", feeID)
	feeJSON, err := ctx.GetStub().GetState(feeKey)
	if err != nil || feeJSON == nil {
		return "", fmt.Errorf("dispute fee %s not found", feeID)
	}
	return string(feeJSON), nil
}

// GetDisputeFeeByDisputeID retrieves dispute fee by dispute ID
func (p *PlatformContract) GetDisputeFeeByDisputeID(ctx contractapi.TransactionContextInterface, disputeID string) (string, error) {
	feeID := fmt.Sprintf("DFEE_%s", disputeID)
	return p.GetDisputeFee(ctx, feeID)
}

// GetDisputeFeeTiers retrieves all dispute fee tiers
func (p *PlatformContract) GetDisputeFeeTiers(ctx contractapi.TransactionContextInterface) (string, error) {
	tierIDs := []string{"BASIC", "STANDARD", "COMPLEX", "HIGH_VALUE"}
	var tiers []DisputeFeeTier

	for _, tierID := range tierIDs {
		tierKey := fmt.Sprintf("DISPUTE_FEE_TIER_%s", tierID)
		tierJSON, err := ctx.GetStub().GetState(tierKey)
		if err != nil || tierJSON == nil {
			continue
		}
		var tier DisputeFeeTier
		json.Unmarshal(tierJSON, &tier)
		tiers = append(tiers, tier)
	}

	result := map[string]interface{}{
		"tiers": tiers,
		"note":  "Fee is locked when dispute is created. Refunded if won, forfeited if lost/dismissed.",
	}
	resultJSON, _ := json.Marshal(result)
	return string(resultJSON), nil
}

// GetUserDisputeFees retrieves all dispute fees for a user
func (p *PlatformContract) GetUserDisputeFees(ctx contractapi.TransactionContextInterface, userID string) (string, error) {
	queryString := fmt.Sprintf(`{"selector":{"initiatorId":"%s"}}`, userID)
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return "", fmt.Errorf("failed to query dispute fees: %v", err)
	}
	defer resultsIterator.Close()

	var fees []DisputeFeeRecord
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			continue
		}
		
		// Only include dispute fee records
		if strings.HasPrefix(queryResult.Key, "DISPUTE_FEE_") {
			var fee DisputeFeeRecord
			json.Unmarshal(queryResult.Value, &fee)
			fees = append(fees, fee)
		}
	}

	totalPaid := 0.0
	totalRefunded := 0.0
	totalForfeited := 0.0
	for _, f := range fees {
		totalPaid += f.FeeAmount
		if f.Status == "REFUNDED" {
			totalRefunded += f.FeeAmount
		} else if f.Status == "FORFEITED" {
			totalForfeited += f.FeeAmount
		}
	}

	result := map[string]interface{}{
		"fees":           fees,
		"totalCount":     len(fees),
		"totalPaid":      totalPaid,
		"totalRefunded":  totalRefunded,
		"totalForfeited": totalForfeited,
	}
	resultJSON, _ := json.Marshal(result)
	return string(resultJSON), nil
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
