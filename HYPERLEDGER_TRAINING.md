<!-- # Hyperledger Fabric Advanced Training
## Multi-Channel, ERC-20 Tokens, and Private Data Collections (PDC)

This training guide covers three advanced Hyperledger Fabric concepts essential for building enterprise-grade, privacy-preserving decentralized applications.

---

## 1. Multi-Channel Architecture and Creation

### 1.1 Why Multi-Channel?
Channels in Hyperledger Fabric provide **strong data isolation** and **confidentiality** at the infrastructure level. A channel is a private "subnet" of communication between two or more specific network members.
*   **Use Case:** You have a `common-channel` for global state (like token balances) visible to all organizations, but you need a `startup-investor-channel` where only Startups and Investors can negotiate investment terms privately. The Validator and Platform organizations cannot see the transactions on this private channel.

### 1.2 Creating and Joining Channels (Practical Flow)

To create a channel, the network administrator generally follows these steps using the Fabric CLI tools:

**Step 1: Create the Genesis Block (or fetch from Orderer)**
Depending on whether you use the newer `osnadmin` service or the older `configtxgen`, the process generates a genesis block for the new channel.

**Step 2: Orderer Joins the Channel**
```bash
# Using osnadmin to join the orderer to the new channel
osnadmin channel join --channelID startup-investor-channel \
    --config-block /path/to/startup-investor-channel.block \
    -o orderer-api.example.com:7053 \
    --ca-file /path/to/orderer/ca.crt \
    --client-cert /path/to/admin/client.crt \
    --client-key /path/to/admin/client.key
```

**Step 3: Peers Join the Channel**
Each organization that belongs to the channel must join their peers to it.
```bash
# Set context for InvestorOrg
export CORE_PEER_LOCALMSPID="InvestorOrgMSP"
export CORE_PEER_ADDRESS=peer0.investor.example.com:7051

# Fetch the channel block from the orderer
peer channel fetch 0 startup-investor-channel.block -o orderer.example.com:7050 -c startup-investor-channel

# Join the peer to the channel
peer channel join -b startup-investor-channel.block
```

---

## 2. ERC-20 Token Chaincode Sample (Go)

The ERC-20 standard defines a common list of rules for Ethereum tokens, but the pattern is heavily used in Hyperledger Fabric for utility tokens and stablecoins.

### 2.1 ERC-20 State Structure
In Fabric, we use a composite key pattern (e.g., `account~token`) to store balances efficiently.

### 2.2 Go Implementation Sample
Below is a simplified implementation of an ERC-20 style chaincode using the `contractapi-go`.

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// TokenContract provides functions for managing an ERC-20 style token
type TokenContract struct {
	contractapi.Contract
}

// Mint creates new tokens and adds them to the minter's account
func (t *TokenContract) Mint(ctx contractapi.TransactionContextInterface, amount int) error {
	clientMSPID, _ := ctx.GetClientIdentity().GetMSPID()
	
	// Access Control: Only the PlatformOrg can mint new tokens
	if clientMSPID != "PlatformOrgMSP" {
		return fmt.Errorf("unauthorized: only PlatformOrg can mint tokens")
	}

	clientID, _ := ctx.GetClientIdentity().GetID()
	
	// Update total supply
	totalSupplyBytes, _ := ctx.GetStub().GetState("TotalSupply")
	var totalSupply int
	if totalSupplyBytes != nil {
		json.Unmarshal(totalSupplyBytes, &totalSupply)
	}
	totalSupply += amount
	newSupplyBytes, _ := json.Marshal(totalSupply)
	ctx.GetStub().PutState("TotalSupply", newSupplyBytes)

	// Update minter's balance
	balanceBytes, _ := ctx.GetStub().GetState("Balance_" + clientID)
	var balance int
	if balanceBytes != nil {
		json.Unmarshal(balanceBytes, &balance)
	}
	balance += amount
	newBalanceBytes, _ := json.Marshal(balance)
	return ctx.GetStub().PutState("Balance_" + clientID, newBalanceBytes)
}

// Transfer transfers tokens from the caller's account to a recipient
func (t *TokenContract) Transfer(ctx contractapi.TransactionContextInterface, recipientID string, amount int) error {
	senderID, _ := ctx.GetClientIdentity().GetID()

	if amount <= 0 {
		return fmt.Errorf("transfer amount must be greater than zero")
	}

	// Get sender balance
	senderBalanceBytes, _ := ctx.GetStub().GetState("Balance_" + senderID)
	var senderBalance int
	if senderBalanceBytes != nil {
		json.Unmarshal(senderBalanceBytes, &senderBalance)
	}

	if senderBalance < amount {
		return fmt.Errorf("insufficient funds")
	}

	// Get recipient balance
	recipientBalanceBytes, _ := ctx.GetStub().GetState("Balance_" + recipientID)
	var recipientBalance int
	if recipientBalanceBytes != nil {
		json.Unmarshal(recipientBalanceBytes, &recipientBalance)
	}

	// Update balances
	senderBalance -= amount
	recipientBalance += amount

	senderBalJSON, _ := json.Marshal(senderBalance)
	recipientBalJSON, _ := json.Marshal(recipientBalance)

	ctx.GetStub().PutState("Balance_" + senderID, senderBalJSON)
	ctx.GetStub().PutState("Balance_" + recipientID, recipientBalJSON)

	return nil
}

// BalanceOf returns the balance of a specific account
func (t *TokenContract) BalanceOf(ctx contractapi.TransactionContextInterface, accountID string) (int, error) {
	balanceBytes, err := ctx.GetStub().GetState("Balance_" + accountID)
	if err != nil {
		return 0, fmt.Errorf("failed to read from world state: %v", err)
	}
	if balanceBytes == nil {
		return 0, nil // Account doesn't exist yet, balance is 0
	}

	var balance int
	json.Unmarshal(balanceBytes, &balance)
	return balance, nil
}
```

---

## 3. Private Data Collections (PDC)

### 3.1 Why PDC?
While Channels isolate data entirely (the data never leaves the organizations on the channel), Private Data Collections allow organizations on the *same channel* to share confidential data with a subset of members. 
*   **The benefit:** You don't have the overhead of managing dozens of channels. The hash of the private data goes on the public ledger (for auditing and immutability), but the actual payload is distributed only to authorized peers via gossip protocol.

### 3.2 Defining Collections (`collections_config.json`)
When you package and deploy your chaincode, you must provide a collections definition file.

```json
[
  {
    "name": "PlatformPrivateCollection",
    "policy": "OR('PlatformOrgMSP.peer')",
    "requiredPeerCount": 0,
    "maxPeerCount": 1,
    "blockToLive": 0,
    "memberOnlyRead": true,
    "memberOnlyWrite": true
  },
  {
    "name": "InvestorStartupSharedCollection",
    "policy": "OR('InvestorOrgMSP.peer', 'StartupOrgMSP.peer')",
    "requiredPeerCount": 1,
    "maxPeerCount": 2,
    "blockToLive": 1000000,
    "memberOnlyRead": true,
    "memberOnlyWrite": false
  }
]
```

### 3.3 Chaincode Implementation for PDC
To read and write to Private Data Collections, you use `PutPrivateData` and `GetPrivateData` instead of `PutState` and `GetState`.

```go
// StoreConfidentialDocument stores an agreement securely using PDC
func (t *TokenContract) StoreConfidentialDocument(ctx contractapi.TransactionContextInterface, documentID string) error {
	
	// Private data is passed via Transient Map to avoid exposing it in the transaction payload
	transientMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return fmt.Errorf("error getting transient data: %v", err)
	}

	documentJSON, exists := transientMap["document_payload"]
	if !exists {
		return fmt.Errorf("document_payload must be passed as transient data")
	}

	// Write to the private data collection
	err = ctx.GetStub().PutPrivateData("InvestorStartupSharedCollection", "DOC_" + documentID, documentJSON)
	if err != nil {
		return fmt.Errorf("failed to put private data: %v", err)
	}

	return nil
}

// ReadConfidentialDocument retrieves the agreement from the PDC
func (t *TokenContract) ReadConfidentialDocument(ctx contractapi.TransactionContextInterface, documentID string) (string, error) {
	
	// Only organizations explicitly named in the collections_config.json policy can successfully read this
	documentJSON, err := ctx.GetStub().GetPrivateData("InvestorStartupSharedCollection", "DOC_" + documentID)
	if err != nil {
		return "", fmt.Errorf("failed to read private data: %v", err)
	}
	if documentJSON == nil {
		return "", fmt.Errorf("document not found or you are not authorized to view it")
	}

	return string(documentJSON), nil
}
```

### Summary of Differences:
| Feature | Channels | Private Data Collections (PDC) |
| :--- | :--- | :--- |
| **Data Separation** | Entirely separate ledgers and blockchains | Same blockchain, separate private state DBs |
| **Overhead** | High (Ordering service configuration, blocks) | Low (Peer-to-peer gossip dissemination) |
| **Proof on Main Ledger** | No | Yes (Data hash is written to the main channel ledger) | -->
