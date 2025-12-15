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

// Token Types
const (
	TokenTypeFee      = "FEE_TOKEN"      // Platform fee tokens
	TokenTypePayment  = "PAYMENT_TOKEN"  // Investment/payment tokens
	TokenTypeReward   = "REWARD_TOKEN"   // Reward tokens
)

// ============================================================================
// DATA STRUCTURES
// ============================================================================

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
