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
	InvestmentID     string  `json:"investmentId"`
	CampaignID       string  `json:"campaignId"`
	InvestorID       string  `json:"investorId"`
	Amount           float64 `json:"amount"`
	Currency         string  `json:"currency"`
	Status           string  `json:"status"` // COMMITTED, ACKNOWLEDGED, CONFIRMED, WITHDRAWN
	CommittedAt      string  `json:"committedAt"`
	ConfirmedAt      string  `json:"confirmedAt"`
	WithdrawnAt      string  `json:"withdrawnAt"`
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
