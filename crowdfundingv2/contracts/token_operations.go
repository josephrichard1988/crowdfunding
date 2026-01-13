package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// TokenContract handles token-based transactions for fees and payments
type TokenContract struct {
	contractapi.Contract
}

// Token Types - ERC-20 Style
const (
	TokenTypeCFT     = "CFT"           // CrowdToken - Utility token for fees/payments
	TokenTypeCFRT    = "CFRT"          // CrowdRewardToken - Reward incentives
	TokenTypeFee     = "FEE_TOKEN"     // Legacy - Platform fee tokens
	TokenTypePayment = "PAYMENT_TOKEN" // Legacy - Investment/payment tokens
	TokenTypeReward  = "REWARD_TOKEN"  // Legacy - Reward tokens
)

// Default Exchange Rates (can be updated by Platform)
const (
	DefaultINRRate      = 2.5  // 1 INR = 2.5 CFT
	DefaultUSDRate      = 83.0 // 1 USD = 83 CFT
	DefaultCFRTtoCFT    = 10.0 // 1 CFRT = 10 CFT
	PlatformAccountID   = "PLATFORM"
)

// ============================================================================
// DATA STRUCTURES
// ============================================================================

// TokenMetadata stores token configuration (ERC-20 style)
type TokenMetadata struct {
	TokenType   string  `json:"tokenType"`   // CFT or CFRT
	Name        string  `json:"name"`        // CrowdToken or CrowdRewardToken
	Symbol      string  `json:"symbol"`      // CFT or CFRT
	Decimals    int     `json:"decimals"`    // Usually 2
	TotalSupply float64 `json:"totalSupply"` // Current total supply
	MaxSupply   float64 `json:"maxSupply"`   // 0 = unlimited
	Owner       string  `json:"owner"`       // PlatformOrg MSP ID
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}

// ExchangeRate stores fiat-to-CFT conversion rates
type ExchangeRate struct {
	Currency  string  `json:"currency"`  // INR, USD, EUR, etc.
	Rate      float64 `json:"rate"`      // How many CFT per 1 unit of currency
	UpdatedAt string  `json:"updatedAt"`
	UpdatedBy string  `json:"updatedBy"` // Who updated the rate
}

// Token represents a fungible token
type Token struct {
	TokenID      string  `json:"tokenId"`
	TokenType    string  `json:"tokenType"`
	Owner        string  `json:"owner"`
	Amount       float64 `json:"amount"`
	Currency     string  `json:"currency"`
	Status       string  `json:"status"` // ACTIVE, LOCKED, BURNED
	IssuedBy     string  `json:"issuedBy"`
	IssuedAt     string  `json:"issuedAt"`
	LastModified string  `json:"lastModified"`
	Metadata     string  `json:"metadata"` // Additional info (JSON)
}

// TokenAccount represents a user's token balance
type TokenAccount struct {
	AccountID    string             `json:"accountId"`
	Owner        string             `json:"owner"`
	OwnerType    string             `json:"ownerType"` // STARTUP, INVESTOR, VALIDATOR, PLATFORM
	Balances     map[string]float64 `json:"balances"`  // tokenType -> amount
	FrozenAmount map[string]float64 `json:"frozenAmount"`
	CreatedAt    string             `json:"createdAt"`
	UpdatedAt    string             `json:"updatedAt"`
}

// TokenTransfer represents a token transfer transaction
type TokenTransfer struct {
	TransferID   string  `json:"transferId"`
	TokenType    string  `json:"tokenType"`
	From         string  `json:"from"`
	To           string  `json:"to"`
	Amount       float64 `json:"amount"`
	Currency     string  `json:"currency"`
	Purpose      string  `json:"purpose"` // FEE_PAYMENT, INVESTMENT, REFUND, etc.
	CampaignID   string  `json:"campaignId,omitempty"`
	Status       string  `json:"status"` // PENDING, COMPLETED, FAILED
	TransferredAt string `json:"transferredAt"`
	CompletedAt  string  `json:"completedAt,omitempty"`
}

// ============================================================================
// INIT
// ============================================================================

func (t *TokenContract) InitTokenLedger(ctx contractapi.TransactionContextInterface) error {
	fmt.Println("TokenContract initialized")
	return nil
}

// ============================================================================
// TOKEN ACCOUNT MANAGEMENT
// ============================================================================

// CreateTokenAccount creates a token account for a user
func (t *TokenContract) CreateTokenAccount(
	ctx contractapi.TransactionContextInterface,
	accountID string,
	owner string,
	ownerType string,
	initialBalanceJSON string,
) error {

	// Check if account already exists
	existingJSON, err := ctx.GetStub().GetState("TOKEN_ACCOUNT_" + accountID)
	if err == nil && existingJSON != nil {
		return fmt.Errorf("token account %s already exists", accountID)
	}

	timestamp := time.Now().Format(time.RFC3339)

	// Parse initial balances
	var initialBalance map[string]float64
	if initialBalanceJSON != "" {
		err := json.Unmarshal([]byte(initialBalanceJSON), &initialBalance)
		if err != nil {
			return fmt.Errorf("failed to parse initial balance: %v", err)
		}
	} else {
		initialBalance = make(map[string]float64)
		initialBalance[TokenTypePayment] = 0
		initialBalance[TokenTypeFee] = 0
	}

	account := TokenAccount{
		AccountID:    accountID,
		Owner:        owner,
		OwnerType:    ownerType,
		Balances:     initialBalance,
		FrozenAmount: make(map[string]float64),
		CreatedAt:    timestamp,
		UpdatedAt:    timestamp,
	}

	accountJSON, err := json.Marshal(account)
	if err != nil {
		return fmt.Errorf("failed to marshal token account: %v", err)
	}

	err = ctx.GetStub().PutState("TOKEN_ACCOUNT_"+accountID, accountJSON)
	if err != nil {
		return fmt.Errorf("failed to create token account: %v", err)
	}

	return nil
}

// GetTokenAccount retrieves a token account
func (t *TokenContract) GetTokenAccount(
	ctx contractapi.TransactionContextInterface,
	accountID string,
) (*TokenAccount, error) {

	accountJSON, err := ctx.GetStub().GetState("TOKEN_ACCOUNT_" + accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to read token account: %v", err)
	}
	if accountJSON == nil {
		return nil, fmt.Errorf("token account %s does not exist", accountID)
	}

	var account TokenAccount
	err = json.Unmarshal(accountJSON, &account)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal token account: %v", err)
	}

	return &account, nil
}

// ============================================================================
// TOKEN ISSUANCE (MINTING)
// ============================================================================

// IssueTokens mints new tokens and credits to an account
func (t *TokenContract) IssueTokens(
	ctx contractapi.TransactionContextInterface,
	tokenID string,
	tokenType string,
	recipient string,
	amount float64,
	currency string,
	issuer string,
	metadata string,
) error {

	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	timestamp := time.Now().Format(time.RFC3339)

	// Create token record
	token := Token{
		TokenID:      tokenID,
		TokenType:    tokenType,
		Owner:        recipient,
		Amount:       amount,
		Currency:     currency,
		Status:       "ACTIVE",
		IssuedBy:     issuer,
		IssuedAt:     timestamp,
		LastModified: timestamp,
		Metadata:     metadata,
	}

	tokenJSON, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("failed to marshal token: %v", err)
	}

	err = ctx.GetStub().PutState("TOKEN_"+tokenID, tokenJSON)
	if err != nil {
		return fmt.Errorf("failed to store token: %v", err)
	}

	// Credit recipient's account
	account, err := t.GetTokenAccount(ctx, recipient)
	if err != nil {
		// Create account if it doesn't exist
		err = t.CreateTokenAccount(ctx, recipient, recipient, "USER", "{}")
		if err != nil {
			return fmt.Errorf("failed to create account: %v", err)
		}
		account, _ = t.GetTokenAccount(ctx, recipient)
	}

	if account.Balances == nil {
		account.Balances = make(map[string]float64)
	}
	account.Balances[tokenType] += amount
	account.UpdatedAt = timestamp

	updatedAccountJSON, _ := json.Marshal(account)
	err = ctx.GetStub().PutState("TOKEN_ACCOUNT_"+recipient, updatedAccountJSON)
	if err != nil {
		return fmt.Errorf("failed to update recipient account: %v", err)
	}

	return nil
}

// ============================================================================
// TOKEN TRANSFER
// ============================================================================

// TransferTokens transfers tokens between accounts
func (t *TokenContract) TransferTokens(
	ctx contractapi.TransactionContextInterface,
	transferID string,
	tokenType string,
	from string,
	to string,
	amount float64,
	currency string,
	purpose string,
	campaignID string,
) error {

	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	timestamp := time.Now().Format(time.RFC3339)

	// Get sender account
	fromAccount, err := t.GetTokenAccount(ctx, from)
	if err != nil {
		return fmt.Errorf("sender account not found: %v", err)
	}

	// Check balance
	if fromAccount.Balances[tokenType] < amount {
		return fmt.Errorf("insufficient balance. Available: %.2f, Required: %.2f", 
			fromAccount.Balances[tokenType], amount)
	}

	// Get or create recipient account
	toAccount, err := t.GetTokenAccount(ctx, to)
	if err != nil {
		err = t.CreateTokenAccount(ctx, to, to, "USER", "{}")
		if err != nil {
			return fmt.Errorf("failed to create recipient account: %v", err)
		}
		toAccount, _ = t.GetTokenAccount(ctx, to)
	}

	// Perform transfer
	fromAccount.Balances[tokenType] -= amount
	fromAccount.UpdatedAt = timestamp

	if toAccount.Balances == nil {
		toAccount.Balances = make(map[string]float64)
	}
	toAccount.Balances[tokenType] += amount
	toAccount.UpdatedAt = timestamp

	// Update accounts
	fromJSON, _ := json.Marshal(fromAccount)
	toJSON, _ := json.Marshal(toAccount)

	ctx.GetStub().PutState("TOKEN_ACCOUNT_"+from, fromJSON)
	ctx.GetStub().PutState("TOKEN_ACCOUNT_"+to, toJSON)

	// Record transfer
	transfer := TokenTransfer{
		TransferID:    transferID,
		TokenType:     tokenType,
		From:          from,
		To:            to,
		Amount:        amount,
		Currency:      currency,
		Purpose:       purpose,
		CampaignID:    campaignID,
		Status:        "COMPLETED",
		TransferredAt: timestamp,
		CompletedAt:   timestamp,
	}

	transferJSON, _ := json.Marshal(transfer)
	ctx.GetStub().PutState("TOKEN_TRANSFER_"+transferID, transferJSON)

	return nil
}

// ============================================================================
// FEE COLLECTION WITH TOKENS
// ============================================================================

// CollectFeeTokens collects fee tokens from a campaign
func (t *TokenContract) CollectFeeTokens(
	ctx contractapi.TransactionContextInterface,
	feeID string,
	campaignID string,
	startupID string,
	campaignAmount float64,
	feePercent float64,
) error {

	feeAmount := (campaignAmount * feePercent) / 100
	timestamp := time.Now().Format(time.RFC3339)

	// Transfer fee tokens from startup to platform
	transferID := fmt.Sprintf("FEE_TRANSFER_%s_%s", campaignID, timestamp)
	
	err := t.TransferTokens(
		ctx,
		transferID,
		TokenTypeFee,
		startupID,
		"PLATFORM",
		feeAmount,
		"USD",
		"CAMPAIGN_FEE",
		campaignID,
	)

	if err != nil {
		return fmt.Errorf("failed to collect fee: %v", err)
	}

	// Record fee collection
	feeRecord := map[string]interface{}{
		"feeId":          feeID,
		"campaignId":     campaignID,
		"startupId":      startupID,
		"campaignAmount": campaignAmount,
		"feePercent":     feePercent,
		"feeAmount":      feeAmount,
		"transferId":     transferID,
		"collectedAt":    timestamp,
		"status":         "COLLECTED",
	}

	feeJSON, _ := json.Marshal(feeRecord)
	ctx.GetStub().PutState("FEE_COLLECTION_"+feeID, feeJSON)

	return nil
}

// ============================================================================
// FREEZE/UNFREEZE (FOR ESCROW)
// ============================================================================

// FreezeTokens freezes tokens (for escrow during disputes)
func (t *TokenContract) FreezeTokens(
	ctx contractapi.TransactionContextInterface,
	accountID string,
	tokenType string,
	amount float64,
	reason string,
) error {

	account, err := t.GetTokenAccount(ctx, accountID)
	if err != nil {
		return fmt.Errorf("account not found: %v", err)
	}

	if account.Balances[tokenType] < amount {
		return fmt.Errorf("insufficient balance to freeze")
	}

	timestamp := time.Now().Format(time.RFC3339)

	// Move from balance to frozen
	account.Balances[tokenType] -= amount
	if account.FrozenAmount == nil {
		account.FrozenAmount = make(map[string]float64)
	}
	account.FrozenAmount[tokenType] += amount
	account.UpdatedAt = timestamp

	accountJSON, _ := json.Marshal(account)
	ctx.GetStub().PutState("TOKEN_ACCOUNT_"+accountID, accountJSON)

	// Record freeze action
	freezeRecord := map[string]interface{}{
		"accountId": accountID,
		"tokenType": tokenType,
		"amount":    amount,
		"reason":    reason,
		"frozenAt":  timestamp,
	}
	freezeJSON, _ := json.Marshal(freezeRecord)
	ctx.GetStub().PutState("FREEZE_"+accountID+"_"+timestamp, freezeJSON)

	return nil
}

// UnfreezeTokens unfreezes tokens
func (t *TokenContract) UnfreezeTokens(
	ctx contractapi.TransactionContextInterface,
	accountID string,
	tokenType string,
	amount float64,
) error {

	account, err := t.GetTokenAccount(ctx, accountID)
	if err != nil {
		return fmt.Errorf("account not found: %v", err)
	}

	if account.FrozenAmount[tokenType] < amount {
		return fmt.Errorf("insufficient frozen amount")
	}

	timestamp := time.Now().Format(time.RFC3339)

	// Move from frozen to balance
	account.FrozenAmount[tokenType] -= amount
	account.Balances[tokenType] += amount
	account.UpdatedAt = timestamp

	accountJSON, _ := json.Marshal(account)
	ctx.GetStub().PutState("TOKEN_ACCOUNT_"+accountID, accountJSON)

	return nil
}

// ============================================================================
// QUERY FUNCTIONS
// ============================================================================

// GetBalance returns token balance for an account
func (t *TokenContract) GetBalance(
	ctx contractapi.TransactionContextInterface,
	accountID string,
	tokenType string,
) (float64, error) {

	account, err := t.GetTokenAccount(ctx, accountID)
	if err != nil {
		return 0, err
	}

	return account.Balances[tokenType], nil
}

// GetTransferHistory retrieves transfer history
func (t *TokenContract) GetTransferHistory(
	ctx contractapi.TransactionContextInterface,
	accountID string,
) (string, error) {

	// Query transfers where accountID is sender or recipient
	queryString := fmt.Sprintf(`{
		"selector": {
			"$or": [
				{"from": "%s"},
				{"to": "%s"}
			]
		}
	}`, accountID, accountID)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return "", fmt.Errorf("failed to query transfers: %v", err)
	}
	defer resultsIterator.Close()

	var transfers []map[string]interface{}
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			continue
		}

		var transfer map[string]interface{}
		json.Unmarshal(queryResponse.Value, &transfer)
		transfers = append(transfers, transfer)
	}

	transfersJSON, _ := json.Marshal(transfers)
	return string(transfersJSON), nil
}

// ============================================================================
// TOKEN INITIALIZATION (ERC-20 Style)
// ============================================================================

// InitializeToken creates a new token type (CFT or CFRT) - Platform only
func (t *TokenContract) InitializeToken(
	ctx contractapi.TransactionContextInterface,
	tokenType string,
	name string,
	decimals int,
	initialSupply float64,
	maxSupply float64,
) error {

	// Check if token already exists
	existingJSON, err := ctx.GetStub().GetState("TOKEN_METADATA_" + tokenType)
	if err == nil && existingJSON != nil {
		return fmt.Errorf("token %s already initialized", tokenType)
	}

	timestamp := time.Now().Format(time.RFC3339)

	metadata := TokenMetadata{
		TokenType:   tokenType,
		Name:        name,
		Symbol:      tokenType,
		Decimals:    decimals,
		TotalSupply: initialSupply,
		MaxSupply:   maxSupply, // 0 = unlimited
		Owner:       PlatformAccountID,
		CreatedAt:   timestamp,
		UpdatedAt:   timestamp,
	}

	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %v", err)
	}

	err = ctx.GetStub().PutState("TOKEN_METADATA_"+tokenType, metadataJSON)
	if err != nil {
		return fmt.Errorf("failed to store token metadata: %v", err)
	}

	// Create Platform account with initial supply if > 0
	if initialSupply > 0 {
		err = t.CreateTokenAccount(ctx, PlatformAccountID, PlatformAccountID, "PLATFORM", "{}")
		if err != nil {
			// Account may already exist, continue
		}
		account, _ := t.GetTokenAccount(ctx, PlatformAccountID)
		if account != nil {
			if account.Balances == nil {
				account.Balances = make(map[string]float64)
			}
			account.Balances[tokenType] = initialSupply
			account.UpdatedAt = timestamp
			accountJSON, _ := json.Marshal(account)
			ctx.GetStub().PutState("TOKEN_ACCOUNT_"+PlatformAccountID, accountJSON)
		}
	}

	return nil
}

// GetTokenMetadata retrieves token metadata
func (t *TokenContract) GetTokenMetadata(
	ctx contractapi.TransactionContextInterface,
	tokenType string,
) (*TokenMetadata, error) {

	metadataJSON, err := ctx.GetStub().GetState("TOKEN_METADATA_" + tokenType)
	if err != nil {
		return nil, fmt.Errorf("failed to read token metadata: %v", err)
	}
	if metadataJSON == nil {
		return nil, fmt.Errorf("token %s not initialized", tokenType)
	}

	var metadata TokenMetadata
	err = json.Unmarshal(metadataJSON, &metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %v", err)
	}

	return &metadata, nil
}

// ============================================================================
// EXCHANGE RATE MANAGEMENT
// ============================================================================

// SetExchangeRate sets the exchange rate for a currency (Platform only)
func (t *TokenContract) SetExchangeRate(
	ctx contractapi.TransactionContextInterface,
	currency string,
	rate float64,
) error {

	if rate <= 0 {
		return fmt.Errorf("rate must be positive")
	}

	timestamp := time.Now().Format(time.RFC3339)

	exchangeRate := ExchangeRate{
		Currency:  currency,
		Rate:      rate,
		UpdatedAt: timestamp,
		UpdatedBy: PlatformAccountID,
	}

	rateJSON, err := json.Marshal(exchangeRate)
	if err != nil {
		return fmt.Errorf("failed to marshal exchange rate: %v", err)
	}

	err = ctx.GetStub().PutState("EXCHANGE_RATE_"+currency, rateJSON)
	if err != nil {
		return fmt.Errorf("failed to store exchange rate: %v", err)
	}

	return nil
}

// GetExchangeRate retrieves the exchange rate for a currency
func (t *TokenContract) GetExchangeRate(
	ctx contractapi.TransactionContextInterface,
	currency string,
) (float64, error) {

	rateJSON, err := ctx.GetStub().GetState("EXCHANGE_RATE_" + currency)
	if err != nil {
		return 0, fmt.Errorf("failed to read exchange rate: %v", err)
	}
	if rateJSON == nil {
		// Return default rates
		switch currency {
		case "INR":
			return DefaultINRRate, nil
		case "USD":
			return DefaultUSDRate, nil
		default:
			return 0, fmt.Errorf("exchange rate for %s not found", currency)
		}
	}

	var exchangeRate ExchangeRate
	err = json.Unmarshal(rateJSON, &exchangeRate)
	if err != nil {
		return 0, fmt.Errorf("failed to unmarshal exchange rate: %v", err)
	}

	return exchangeRate.Rate, nil
}

// ============================================================================
// TOKEN PURCHASE (Fiat to CFT)
// ============================================================================

// PurchaseTokens converts fiat payment to CFT tokens
func (t *TokenContract) PurchaseTokens(
	ctx contractapi.TransactionContextInterface,
	accountID string,
	currency string,
	fiatAmount float64,
) error {

	if fiatAmount <= 0 {
		return fmt.Errorf("fiat amount must be positive")
	}

	// Get exchange rate
	rate, err := t.GetExchangeRate(ctx, currency)
	if err != nil {
		return fmt.Errorf("failed to get exchange rate: %v", err)
	}

	// Calculate CFT amount
	cftAmount := fiatAmount * rate

	timestamp := time.Now().Format(time.RFC3339)

	// Get or create account
	account, err := t.GetTokenAccount(ctx, accountID)
	if err != nil {
		err = t.CreateTokenAccount(ctx, accountID, accountID, "USER", "{}")
		if err != nil {
			return fmt.Errorf("failed to create account: %v", err)
		}
		account, _ = t.GetTokenAccount(ctx, accountID)
	}

	if account.Balances == nil {
		account.Balances = make(map[string]float64)
	}
	account.Balances[TokenTypeCFT] += cftAmount
	account.UpdatedAt = timestamp

	accountJSON, _ := json.Marshal(account)
	err = ctx.GetStub().PutState("TOKEN_ACCOUNT_"+accountID, accountJSON)
	if err != nil {
		return fmt.Errorf("failed to update account: %v", err)
	}

	// Update total supply
	metadata, _ := t.GetTokenMetadata(ctx, TokenTypeCFT)
	if metadata != nil {
		metadata.TotalSupply += cftAmount
		metadata.UpdatedAt = timestamp
		metadataJSON, _ := json.Marshal(metadata)
		ctx.GetStub().PutState("TOKEN_METADATA_"+TokenTypeCFT, metadataJSON)
	}

	// Record purchase
	purchaseRecord := map[string]interface{}{
		"accountId":  accountID,
		"currency":   currency,
		"fiatAmount": fiatAmount,
		"cftAmount":  cftAmount,
		"rate":       rate,
		"timestamp":  timestamp,
	}
	purchaseJSON, _ := json.Marshal(purchaseRecord)
	ctx.GetStub().PutState("TOKEN_PURCHASE_"+accountID+"_"+timestamp, purchaseJSON)

	return nil
}

// ============================================================================
// MINT TOKENS (Platform only)
// ============================================================================

// MintTokens mints new tokens to an account (Platform only)
func (t *TokenContract) MintTokens(
	ctx contractapi.TransactionContextInterface,
	tokenType string,
	recipient string,
	amount float64,
) error {

	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	timestamp := time.Now().Format(time.RFC3339)

	// Get or create recipient account
	account, err := t.GetTokenAccount(ctx, recipient)
	if err != nil {
		err = t.CreateTokenAccount(ctx, recipient, recipient, "USER", "{}")
		if err != nil {
			return fmt.Errorf("failed to create account: %v", err)
		}
		account, _ = t.GetTokenAccount(ctx, recipient)
	}

	if account.Balances == nil {
		account.Balances = make(map[string]float64)
	}
	account.Balances[tokenType] += amount
	account.UpdatedAt = timestamp

	accountJSON, _ := json.Marshal(account)
	ctx.GetStub().PutState("TOKEN_ACCOUNT_"+recipient, accountJSON)

	// Update total supply
	metadata, _ := t.GetTokenMetadata(ctx, tokenType)
	if metadata != nil {
		// Check max supply
		if metadata.MaxSupply > 0 && (metadata.TotalSupply+amount) > metadata.MaxSupply {
			return fmt.Errorf("exceeds max supply of %.2f", metadata.MaxSupply)
		}
		metadata.TotalSupply += amount
		metadata.UpdatedAt = timestamp
		metadataJSON, _ := json.Marshal(metadata)
		ctx.GetStub().PutState("TOKEN_METADATA_"+tokenType, metadataJSON)
	}

	return nil
}

// ============================================================================
// BURN TOKENS
// ============================================================================

// BurnTokens destroys tokens from an account
func (t *TokenContract) BurnTokens(
	ctx contractapi.TransactionContextInterface,
	accountID string,
	tokenType string,
	amount float64,
) error {

	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	account, err := t.GetTokenAccount(ctx, accountID)
	if err != nil {
		return fmt.Errorf("account not found: %v", err)
	}

	if account.Balances[tokenType] < amount {
		return fmt.Errorf("insufficient balance to burn")
	}

	timestamp := time.Now().Format(time.RFC3339)

	account.Balances[tokenType] -= amount
	account.UpdatedAt = timestamp

	accountJSON, _ := json.Marshal(account)
	ctx.GetStub().PutState("TOKEN_ACCOUNT_"+accountID, accountJSON)

	// Update total supply
	metadata, _ := t.GetTokenMetadata(ctx, tokenType)
	if metadata != nil {
		metadata.TotalSupply -= amount
		metadata.UpdatedAt = timestamp
		metadataJSON, _ := json.Marshal(metadata)
		ctx.GetStub().PutState("TOKEN_METADATA_"+tokenType, metadataJSON)
	}

	return nil
}

// ============================================================================
// REWARD TOKEN DISTRIBUTION
// ============================================================================

// DistributeRewards distributes CFRT reward tokens (Platform only)
func (t *TokenContract) DistributeRewards(
	ctx contractapi.TransactionContextInterface,
	rewardID string,
	recipient string,
	amount float64,
	reason string,
) error {

	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	timestamp := time.Now().Format(time.RFC3339)

	// Mint CFRT to recipient
	err := t.MintTokens(ctx, TokenTypeCFRT, recipient, amount)
	if err != nil {
		return fmt.Errorf("failed to mint reward tokens: %v", err)
	}

	// Record reward distribution
	rewardRecord := map[string]interface{}{
		"rewardId":  rewardID,
		"recipient": recipient,
		"amount":    amount,
		"reason":    reason,
		"timestamp": timestamp,
	}
	rewardJSON, _ := json.Marshal(rewardRecord)
	ctx.GetStub().PutState("REWARD_DISTRIBUTION_"+rewardID, rewardJSON)

	return nil
}

// RedeemRewardTokens converts CFRT to CFT (1 CFRT = 10 CFT)
func (t *TokenContract) RedeemRewardTokens(
	ctx contractapi.TransactionContextInterface,
	accountID string,
	cfrtAmount float64,
) error {

	if cfrtAmount <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	account, err := t.GetTokenAccount(ctx, accountID)
	if err != nil {
		return fmt.Errorf("account not found: %v", err)
	}

	if account.Balances[TokenTypeCFRT] < cfrtAmount {
		return fmt.Errorf("insufficient CFRT balance. Available: %.2f, Required: %.2f",
			account.Balances[TokenTypeCFRT], cfrtAmount)
	}

	// Calculate CFT to credit
	cftAmount := cfrtAmount * DefaultCFRTtoCFT

	timestamp := time.Now().Format(time.RFC3339)

	// Deduct CFRT, credit CFT
	account.Balances[TokenTypeCFRT] -= cfrtAmount
	account.Balances[TokenTypeCFT] += cftAmount
	account.UpdatedAt = timestamp

	accountJSON, _ := json.Marshal(account)
	ctx.GetStub().PutState("TOKEN_ACCOUNT_"+accountID, accountJSON)

	// Record redemption
	redemptionRecord := map[string]interface{}{
		"accountId":  accountID,
		"cfrtAmount": cfrtAmount,
		"cftAmount":  cftAmount,
		"rate":       DefaultCFRTtoCFT,
		"timestamp":  timestamp,
	}
	redemptionJSON, _ := json.Marshal(redemptionRecord)
	ctx.GetStub().PutState("CFRT_REDEMPTION_"+accountID+"_"+timestamp, redemptionJSON)

	return nil
}

// ============================================================================
// WITHDRAWAL (CFT to Fiat)
// ============================================================================

// WithdrawToFiat converts CFT back to fiat (for withdrawal)
func (t *TokenContract) WithdrawToFiat(
	ctx contractapi.TransactionContextInterface,
	accountID string,
	cftAmount float64,
	currency string,
) (float64, error) {

	if cftAmount <= 0 {
		return 0, fmt.Errorf("amount must be positive")
	}

	account, err := t.GetTokenAccount(ctx, accountID)
	if err != nil {
		return 0, fmt.Errorf("account not found: %v", err)
	}

	if account.Balances[TokenTypeCFT] < cftAmount {
		return 0, fmt.Errorf("insufficient CFT balance")
	}

	// Get exchange rate
	rate, err := t.GetExchangeRate(ctx, currency)
	if err != nil {
		return 0, fmt.Errorf("failed to get exchange rate: %v", err)
	}

	// Calculate fiat amount (CFT / rate = fiat)
	fiatAmount := cftAmount / rate

	// Apply 1% withdrawal fee
	withdrawalFee := fiatAmount * 0.01
	netFiatAmount := fiatAmount - withdrawalFee

	timestamp := time.Now().Format(time.RFC3339)

	// Deduct CFT from account
	account.Balances[TokenTypeCFT] -= cftAmount
	account.UpdatedAt = timestamp

	accountJSON, _ := json.Marshal(account)
	ctx.GetStub().PutState("TOKEN_ACCOUNT_"+accountID, accountJSON)

	// Record withdrawal
	withdrawalRecord := map[string]interface{}{
		"accountId":      accountID,
		"cftAmount":      cftAmount,
		"currency":       currency,
		"grossFiat":      fiatAmount,
		"fee":            withdrawalFee,
		"netFiat":        netFiatAmount,
		"rate":           rate,
		"timestamp":      timestamp,
		"status":         "PENDING", // Actual fiat transfer happens off-chain
	}
	withdrawalJSON, _ := json.Marshal(withdrawalRecord)
	ctx.GetStub().PutState("WITHDRAWAL_"+accountID+"_"+timestamp, withdrawalJSON)

	return netFiatAmount, nil
}

// ============================================================================
// DISPUTE RESOLUTION SYSTEM
// ============================================================================

// DisputeTicket represents a dispute ticket on-chain
type DisputeTicket struct {
	TicketID          string   `json:"ticketId"`
	DisputeType       string   `json:"disputeType"`       // STARTUP_INVESTOR, STARTUP_VALIDATOR, etc.
	InitiatorID       string   `json:"initiatorId"`
	InitiatorType     string   `json:"initiatorType"`     // STARTUP, INVESTOR, VALIDATOR
	RespondentID      string   `json:"respondentId"`
	RespondentType    string   `json:"respondentType"`
	CampaignID        string   `json:"campaignId,omitempty"`
	AgreementID       string   `json:"agreementId,omitempty"`
	ClaimAmount       float64  `json:"claimAmount"`       // Amount in CFT
	PenaltyCategory   string   `json:"penaltyCategory"`   // FRAUD, DELAY, BREACH, etc.
	EvidenceHashes    []string `json:"evidenceHashes"`
	Status            string   `json:"status"`            // OPEN, INVESTIGATING, RESOLVED, CLOSED
	InvestigatorID    string   `json:"investigatorId,omitempty"`
	FrozenAmount      float64  `json:"frozenAmount"`      // CFT frozen as penalty hold
	PenaltyAmount     float64  `json:"penaltyAmount"`     // Final penalty in CFT
	Resolution        string   `json:"resolution,omitempty"`
	RatingImpact      float64  `json:"ratingImpact"`      // Impact on ML rating
	CreatedAt         string   `json:"createdAt"`
	ResolvedAt        string   `json:"resolvedAt,omitempty"`
}

// MLRating represents the ML-based rating for a user
type MLRating struct {
	UserID            string  `json:"userId"`
	UserType          string  `json:"userType"`          // STARTUP, INVESTOR, VALIDATOR
	OverallScore      float64 `json:"overallScore"`      // 0-100
	TrustScore        float64 `json:"trustScore"`        // 0-100
	DisputeScore      float64 `json:"disputeScore"`      // Impact from disputes
	ComplianceScore   float64 `json:"complianceScore"`   // Platform rule compliance
	EngagementScore   float64 `json:"engagementScore"`   // Platform activity
	TotalDisputes     int     `json:"totalDisputes"`
	DisputesWon       int     `json:"disputesWon"`
	DisputesLost      int     `json:"disputesLost"`
	TotalPenalties    float64 `json:"totalPenalties"`    // CFT penalties received
	TotalRewards      float64 `json:"totalRewards"`      // CFRT rewards received
	BlacklistStatus   bool    `json:"blacklistStatus"`
	BlacklistReason   string  `json:"blacklistReason,omitempty"`
	FeeTier           string  `json:"feeTier"`           // STANDARD, PREMIUM, TRUSTED
	CreatedAt         string  `json:"createdAt"`
	UpdatedAt         string  `json:"updatedAt"`
}

// CreateDisputeTicket creates a new dispute ticket
func (t *TokenContract) CreateDisputeTicket(
	ctx contractapi.TransactionContextInterface,
	ticketID string,
	disputeType string,
	initiatorID string,
	initiatorType string,
	respondentID string,
	respondentType string,
	campaignID string,
	agreementID string,
	claimAmount float64,
	penaltyCategory string,
	evidenceHashesJSON string,
) error {

	// Validate dispute type
	validDisputeTypes := map[string]bool{
		"STARTUP_INVESTOR":  true,
		"STARTUP_VALIDATOR": true,
		"INVESTOR_VALIDATOR": true,
		"STARTUP_PLATFORM":  true,
		"INVESTOR_PLATFORM": true,
		"MULTILATERAL":      true,
	}
	if !validDisputeTypes[disputeType] {
		return fmt.Errorf("invalid dispute type: %s", disputeType)
	}

	// Parse evidence hashes
	var evidenceHashes []string
	if evidenceHashesJSON != "" {
		json.Unmarshal([]byte(evidenceHashesJSON), &evidenceHashes)
	}

	timestamp := time.Now().Format(time.RFC3339)

	ticket := DisputeTicket{
		TicketID:        ticketID,
		DisputeType:     disputeType,
		InitiatorID:     initiatorID,
		InitiatorType:   initiatorType,
		RespondentID:    respondentID,
		RespondentType:  respondentType,
		CampaignID:      campaignID,
		AgreementID:     agreementID,
		ClaimAmount:     claimAmount,
		PenaltyCategory: penaltyCategory,
		EvidenceHashes:  evidenceHashes,
		Status:          "OPEN",
		FrozenAmount:    0,
		PenaltyAmount:   0,
		RatingImpact:    0,
		CreatedAt:       timestamp,
	}

	ticketJSON, err := json.Marshal(ticket)
	if err != nil {
		return fmt.Errorf("failed to marshal ticket: %v", err)
	}

	// Store on world state
	err = ctx.GetStub().PutState("DISPUTE_TICKET_"+ticketID, ticketJSON)
	if err != nil {
		return fmt.Errorf("failed to store dispute ticket: %v", err)
	}

	return nil
}

// FreezePenaltyTokens freezes tokens from respondent as penalty hold
func (t *TokenContract) FreezePenaltyTokens(
	ctx contractapi.TransactionContextInterface,
	ticketID string,
	respondentID string,
	freezeAmount float64,
) error {

	if freezeAmount <= 0 {
		return fmt.Errorf("freeze amount must be positive")
	}

	// Get dispute ticket
	ticketJSON, err := ctx.GetStub().GetState("DISPUTE_TICKET_" + ticketID)
	if err != nil || ticketJSON == nil {
		return fmt.Errorf("dispute ticket not found: %v", err)
	}

	var ticket DisputeTicket
	err = json.Unmarshal(ticketJSON, &ticket)
	if err != nil {
		return fmt.Errorf("failed to unmarshal ticket: %v", err)
	}

	if ticket.Status != "OPEN" && ticket.Status != "INVESTIGATING" {
		return fmt.Errorf("cannot freeze tokens for dispute in status: %s", ticket.Status)
	}

	// Get respondent account
	account, err := t.GetTokenAccount(ctx, respondentID)
	if err != nil {
		return fmt.Errorf("respondent account not found: %v", err)
	}

	if account.Balances[TokenTypeCFT] < freezeAmount {
		return fmt.Errorf("insufficient balance to freeze. Available: %.2f, Required: %.2f", 
			account.Balances[TokenTypeCFT], freezeAmount)
	}

	timestamp := time.Now().Format(time.RFC3339)

	// Move from balance to frozen
	account.Balances[TokenTypeCFT] -= freezeAmount
	if account.FrozenAmount == nil {
		account.FrozenAmount = make(map[string]float64)
	}
	account.FrozenAmount[TokenTypeCFT] += freezeAmount
	account.UpdatedAt = timestamp

	// Update account
	accountJSON, _ := json.Marshal(account)
	ctx.GetStub().PutState("TOKEN_ACCOUNT_"+respondentID, accountJSON)

	// Update ticket
	ticket.FrozenAmount += freezeAmount
	ticket.Status = "INVESTIGATING"
	ticketJSON, _ = json.Marshal(ticket)
	ctx.GetStub().PutState("DISPUTE_TICKET_"+ticketID, ticketJSON)

	return nil
}

// ExecutePenalty executes the penalty - transfers frozen tokens to platform/refund
func (t *TokenContract) ExecutePenalty(
	ctx contractapi.TransactionContextInterface,
	ticketID string,
	penaltyAmount float64,
	refundToInitiator float64,
	resolution string,
	ratingImpact float64,
) error {

	// Get dispute ticket
	ticketJSON, err := ctx.GetStub().GetState("DISPUTE_TICKET_" + ticketID)
	if err != nil || ticketJSON == nil {
		return fmt.Errorf("dispute ticket not found: %v", err)
	}

	var ticket DisputeTicket
	err = json.Unmarshal(ticketJSON, &ticket)
	if err != nil {
		return fmt.Errorf("failed to unmarshal ticket: %v", err)
	}

	if ticket.Status != "INVESTIGATING" {
		return fmt.Errorf("dispute must be in INVESTIGATING status to execute penalty")
	}

	if penaltyAmount + refundToInitiator > ticket.FrozenAmount {
		return fmt.Errorf("penalty + refund exceeds frozen amount")
	}

	timestamp := time.Now().Format(time.RFC3339)

	// Get respondent account
	respondentAccount, err := t.GetTokenAccount(ctx, ticket.RespondentID)
	if err != nil {
		return fmt.Errorf("respondent account not found: %v", err)
	}

	// Deduct from frozen amount
	respondentAccount.FrozenAmount[TokenTypeCFT] -= (penaltyAmount + refundToInitiator)
	respondentAccount.UpdatedAt = timestamp

	// Return excess frozen amount to balance
	excessAmount := ticket.FrozenAmount - penaltyAmount - refundToInitiator
	if excessAmount > 0 {
		respondentAccount.Balances[TokenTypeCFT] += excessAmount
	}

	accountJSON, _ := json.Marshal(respondentAccount)
	ctx.GetStub().PutState("TOKEN_ACCOUNT_"+ticket.RespondentID, accountJSON)

	// Transfer penalty to platform
	if penaltyAmount > 0 {
		platformAccount, _ := t.GetTokenAccount(ctx, PlatformAccountID)
		if platformAccount == nil {
			platformAccount = &TokenAccount{
				AccountID: PlatformAccountID,
				Owner:     PlatformAccountID,
				OwnerType: "PLATFORM",
				Balances:  make(map[string]float64),
				CreatedAt: timestamp,
			}
		}
		platformAccount.Balances[TokenTypeCFT] += penaltyAmount
		platformAccount.UpdatedAt = timestamp
		platformJSON, _ := json.Marshal(platformAccount)
		ctx.GetStub().PutState("TOKEN_ACCOUNT_"+PlatformAccountID, platformJSON)
	}

	// Refund to initiator
	if refundToInitiator > 0 {
		initiatorAccount, err := t.GetTokenAccount(ctx, ticket.InitiatorID)
		if err == nil && initiatorAccount != nil {
			initiatorAccount.Balances[TokenTypeCFT] += refundToInitiator
			initiatorAccount.UpdatedAt = timestamp
			initiatorJSON, _ := json.Marshal(initiatorAccount)
			ctx.GetStub().PutState("TOKEN_ACCOUNT_"+ticket.InitiatorID, initiatorJSON)
		}
	}

	// Update ticket
	ticket.Status = "RESOLVED"
	ticket.PenaltyAmount = penaltyAmount
	ticket.Resolution = resolution
	ticket.RatingImpact = ratingImpact
	ticket.ResolvedAt = timestamp

	ticketJSON, _ = json.Marshal(ticket)
	ctx.GetStub().PutState("DISPUTE_TICKET_"+ticketID, ticketJSON)

	// Update ML ratings for both parties
	if ratingImpact != 0 {
		// Reduce rating for loser
		t.updateMLRatingInternal(ctx, ticket.RespondentID, ticket.RespondentType, -ratingImpact, penaltyAmount, false)
		// Increase rating for winner
		t.updateMLRatingInternal(ctx, ticket.InitiatorID, ticket.InitiatorType, ratingImpact/2, 0, true)
	}

	return nil
}

// ApplyBlacklist blacklists a user for severe violations
func (t *TokenContract) ApplyBlacklist(
	ctx contractapi.TransactionContextInterface,
	userID string,
	userType string,
	reason string,
) error {

	timestamp := time.Now().Format(time.RFC3339)

	// Get or create ML rating
	rating, _ := t.GetMLRating(ctx, userID)
	if rating == nil {
		rating = &MLRating{
			UserID:       userID,
			UserType:     userType,
			OverallScore: 50,
			TrustScore:   50,
			FeeTier:      "STANDARD",
			CreatedAt:    timestamp,
		}
	}

	rating.BlacklistStatus = true
	rating.BlacklistReason = reason
	rating.OverallScore = 0
	rating.TrustScore = 0
	rating.FeeTier = "BLACKLISTED"
	rating.UpdatedAt = timestamp

	ratingJSON, _ := json.Marshal(rating)
	ctx.GetStub().PutState("ML_RATING_"+userID, ratingJSON)

	return nil
}

// GetDisputeTicket retrieves a dispute ticket
func (t *TokenContract) GetDisputeTicket(
	ctx contractapi.TransactionContextInterface,
	ticketID string,
) (*DisputeTicket, error) {

	ticketJSON, err := ctx.GetStub().GetState("DISPUTE_TICKET_" + ticketID)
	if err != nil || ticketJSON == nil {
		return nil, fmt.Errorf("dispute ticket not found: %v", err)
	}

	var ticket DisputeTicket
	err = json.Unmarshal(ticketJSON, &ticket)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal ticket: %v", err)
	}

	return &ticket, nil
}

// ============================================================================
// ML RATING SYSTEM
// ============================================================================

// InitializeMLRating creates initial ML rating for a user
func (t *TokenContract) InitializeMLRating(
	ctx contractapi.TransactionContextInterface,
	userID string,
	userType string,
) error {

	// Check if rating already exists
	existing, _ := t.GetMLRating(ctx, userID)
	if existing != nil {
		return fmt.Errorf("ML rating already exists for user: %s", userID)
	}

	timestamp := time.Now().Format(time.RFC3339)

	rating := MLRating{
		UserID:          userID,
		UserType:        userType,
		OverallScore:    70.0, // Start with neutral-good rating
		TrustScore:      70.0,
		DisputeScore:    0,
		ComplianceScore: 80.0,
		EngagementScore: 50.0,
		TotalDisputes:   0,
		DisputesWon:     0,
		DisputesLost:    0,
		TotalPenalties:  0,
		TotalRewards:    0,
		BlacklistStatus: false,
		FeeTier:         "STANDARD",
		CreatedAt:       timestamp,
		UpdatedAt:       timestamp,
	}

	ratingJSON, err := json.Marshal(rating)
	if err != nil {
		return fmt.Errorf("failed to marshal rating: %v", err)
	}

	err = ctx.GetStub().PutState("ML_RATING_"+userID, ratingJSON)
	if err != nil {
		return fmt.Errorf("failed to store rating: %v", err)
	}

	return nil
}

// UpdateMLRating updates user's ML rating based on activity/dispute outcome
func (t *TokenContract) UpdateMLRating(
	ctx contractapi.TransactionContextInterface,
	userID string,
	userType string,
	ratingChange float64,
	penaltyAmount float64,
	isDisputeWinner bool,
	reason string,
) error {

	return t.updateMLRatingInternal(ctx, userID, userType, ratingChange, penaltyAmount, isDisputeWinner)
}

// Internal helper for ML rating updates
func (t *TokenContract) updateMLRatingInternal(
	ctx contractapi.TransactionContextInterface,
	userID string,
	userType string,
	ratingChange float64,
	penaltyAmount float64,
	isDisputeWinner bool,
) error {

	timestamp := time.Now().Format(time.RFC3339)

	// Get or create rating
	rating, _ := t.GetMLRating(ctx, userID)
	if rating == nil {
		rating = &MLRating{
			UserID:          userID,
			UserType:        userType,
			OverallScore:    70.0,
			TrustScore:      70.0,
			ComplianceScore: 80.0,
			EngagementScore: 50.0,
			FeeTier:         "STANDARD",
			CreatedAt:       timestamp,
		}
	}

	// Update rating scores
	rating.OverallScore += ratingChange
	rating.TrustScore += ratingChange * 0.8
	rating.DisputeScore += ratingChange

	// Clamp values between 0 and 100
	rating.OverallScore = clampRating(rating.OverallScore)
	rating.TrustScore = clampRating(rating.TrustScore)

	// Update dispute stats
	rating.TotalDisputes++
	if isDisputeWinner {
		rating.DisputesWon++
	} else {
		rating.DisputesLost++
		rating.TotalPenalties += penaltyAmount
	}

	// Update fee tier based on score
	rating.FeeTier = calculateFeeTier(rating.OverallScore)
	rating.UpdatedAt = timestamp

	ratingJSON, _ := json.Marshal(rating)
	ctx.GetStub().PutState("ML_RATING_"+userID, ratingJSON)

	return nil
}

// GetMLRating retrieves user's ML rating
func (t *TokenContract) GetMLRating(
	ctx contractapi.TransactionContextInterface,
	userID string,
) (*MLRating, error) {

	ratingJSON, err := ctx.GetStub().GetState("ML_RATING_" + userID)
	if err != nil || ratingJSON == nil {
		return nil, fmt.Errorf("ML rating not found for user: %s", userID)
	}

	var rating MLRating
	err = json.Unmarshal(ratingJSON, &rating)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal rating: %v", err)
	}

	return &rating, nil
}

// AddRewardBonus adds CFRT reward and updates engagement score
func (t *TokenContract) AddRewardBonus(
	ctx contractapi.TransactionContextInterface,
	userID string,
	rewardAmount float64,
	reason string,
) error {

	timestamp := time.Now().Format(time.RFC3339)

	// Get rating
	rating, _ := t.GetMLRating(ctx, userID)
	if rating == nil {
		return fmt.Errorf("ML rating not found for user: %s", userID)
	}

	// Update engagement and total rewards
	rating.TotalRewards += rewardAmount
	rating.EngagementScore += rewardAmount / 100 // Small boost per reward
	rating.EngagementScore = clampRating(rating.EngagementScore)
	rating.UpdatedAt = timestamp

	ratingJSON, _ := json.Marshal(rating)
	ctx.GetStub().PutState("ML_RATING_"+userID, ratingJSON)

	return nil
}

// Helper function to clamp rating between 0 and 100
func clampRating(value float64) float64 {
	if value < 0 {
		return 0
	}
	if value > 100 {
		return 100
	}
	return value
}

// Helper function to calculate fee tier based on score
func calculateFeeTier(score float64) string {
	if score >= 85 {
		return "TRUSTED"   // Lower fees
	} else if score >= 60 {
		return "STANDARD"
	} else if score >= 30 {
		return "PROBATION" // Higher fees
	}
	return "RESTRICTED"
}

// GetPenaltySchedule returns the penalty amounts for various dispute scenarios
func (t *TokenContract) GetPenaltySchedule(
	ctx contractapi.TransactionContextInterface,
) (string, error) {

	// Penalty schedule based on 1 INR = 2.5 CFT
	schedule := map[string]interface{}{
		"startupPenalties": map[string]interface{}{
			"fraudulentCampaign": map[string]float64{"cfT": 2500, "inr": 1000},
			"misuseOfFunds":      map[string]float64{"cfT": 3750, "inr": 1500},
			"milestoneDefault":   map[string]float64{"cfT": 1250, "inr": 500},
			"documentFraud":      map[string]float64{"cfT": 2000, "inr": 800},
		},
		"investorPenalties": map[string]interface{}{
			"fraudulentClaim":    map[string]float64{"cfT": 500, "inr": 200},
			"refundAbuse":        map[string]float64{"cfT": 375, "inr": 150},
			"harassingStartup":   map[string]float64{"cfT": 625, "inr": 250},
		},
		"validatorPenalties": map[string]interface{}{
			"fraudApproval":      map[string]float64{"cfT": 1250, "inr": 500},
			"biasedValidation":   map[string]float64{"cfT": 625, "inr": 250},
			"delayedValidation":  map[string]float64{"cfT": 250, "inr": 100},
		},
	}

	scheduleJSON, err := json.Marshal(schedule)
	if err != nil {
		return "", fmt.Errorf("failed to marshal schedule: %v", err)
	}

	return string(scheduleJSON), nil
}


