package main

// ============================================================================
// SHARED CONSTANTS - Private Data Collection Names
// ============================================================================

const (
	// Org-specific private collections
	StartupPrivateCollection   = "StartupPrivateData"
	InvestorPrivateCollection  = "InvestorPrivateData"
	ValidatorPrivateCollection = "ValidatorPrivateData"
	PlatformPrivateCollection  = "PlatformPrivateData"

	// Two-party shared collections
	StartupInvestorCollection   = "StartupInvestorShared"
	StartupValidatorCollection  = "StartupValidatorShared"
	StartupPlatformCollection   = "StartupPlatformShared"
	InvestorValidatorCollection = "InvestorValidatorShared"
	InvestorPlatformCollection  = "InvestorPlatformShared"
	ValidatorPlatformCollection = "ValidatorPlatformShared"

	// Multi-party shared collections
	ThreePartyCollection = "ThreePartyShared"
	AllOrgsCollection    = "AllOrgsShared"
)

// ============================================================================
// SHARED DATA STRUCTURES
// ============================================================================

// Investment represents an investment commitment (shared across contracts)
type Investment struct {
	InvestmentID string  `json:"investmentId"`
	CampaignID   string  `json:"campaignId"`
	InvestorID   string  `json:"investorId"`
	Amount       float64 `json:"amount"`
	Currency     string  `json:"currency"`
	Status       string  `json:"status"` // COMMITTED, ACKNOWLEDGED, CONFIRMED, WITHDRAWN
	CommittedAt  string  `json:"committedAt"`
	ConfirmedAt  string  `json:"confirmedAt"`
	WithdrawnAt  string  `json:"withdrawnAt"`
}

// Milestone for milestone-based fund release (shared structure)
type Milestone struct {
	MilestoneID    string  `json:"milestoneId"`
	Title          string  `json:"title"`
	Description    string  `json:"description"`
	TargetDate     string  `json:"targetDate"`
	FundPercentage float64 `json:"fundPercentage"`
	TargetAmount   float64 `json:"targetAmount"`
	Status         string  `json:"status"`
	FundsReleased  bool    `json:"fundsReleased"`
	ReleasedAt     string  `json:"releasedAt"`
}

// Agreement represents investment agreement (shared across contracts)
type Agreement struct {
	AgreementID       string      `json:"agreementId"`
	CampaignID        string      `json:"campaignId"`
	StartupID         string      `json:"startupId"`
	InvestorID        string      `json:"investorId"`
	InvestmentAmount  float64     `json:"investmentAmount"`
	Currency          string      `json:"currency"`
	Milestones        []Milestone `json:"milestones"`
	Terms             string      `json:"terms"`
	Status            string      `json:"status"`
	StartupAccepted   bool        `json:"startupAccepted"`
	InvestorAccepted  bool        `json:"investorAccepted"`
	PlatformWitnessed bool        `json:"platformWitnessed"`
	WitnessedAt       string      `json:"witnessedAt"`
	CreatedAt         string      `json:"createdAt"`
	AcceptedAt        string      `json:"acceptedAt"`
}

// NegotiationEntry tracks negotiation history (shared structure)
type NegotiationEntry struct {
	Round     int     `json:"round"`
	Party     string  `json:"party"`
	Action    string  `json:"action"`
	Amount    float64 `json:"amount"`
	Terms     string  `json:"terms"`
	Timestamp string  `json:"timestamp"`
}

// Campaign represents a startup crowdfunding campaign with all required fields (22 parameters)
type Campaign struct {
	CampaignID string `json:"campaignId"`
	StartupID  string `json:"startupId"`

	// Core Campaign Fields (22 Parameters)
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
	OpenDate            string  `json:"open_date"`
	CloseDate           string  `json:"close_date"`
	FundsRaisedAmount   float64 `json:"funds_raised_amount"`
	FundsRaisedPercent  float64 `json:"funds_raised_percent"`
	Status              string  `json:"status"`
	ValidationStatus    string  `json:"validationStatus"`
	ValidationScore     float64 `json:"validationScore"`
	RiskLevel           string  `json:"riskLevel"`
	SubmissionHash      string  `json:"submissionHash"`
	ValidationProofHash string  `json:"validationProofHash"` /* Was ValidationHash */
	InvestorCount       int     `json:"investorCount"`
	PlatformStatus      string  `json:"platformStatus"`

	// Timestamps
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	ApprovedAt  string `json:"approvedAt"`
	PublishedAt string `json:"publishedAt"`
}
