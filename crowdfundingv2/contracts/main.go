package main

import (
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	// Create new chaincode with all 5 contracts (4 org contracts + token contract)
	crowdfundingChaincode, err := contractapi.NewChaincode(
		&StartupContract{},
		&InvestorContract{},
		&ValidatorContract{},
		&PlatformContract{},
		&TokenContract{}, // Token operations for fees and payments
	)

	if err != nil {
		log.Panicf("Error creating crowdfunding chaincode: %v", err)
	}

	// Set chaincode info
	crowdfundingChaincode.Info.Title = "Crowdfunding Platform Chaincode with Token SDK"
	crowdfundingChaincode.Info.Version = "1.0.0"

	// Start chaincode
	if err := crowdfundingChaincode.Start(); err != nil {
		log.Panicf("Error starting crowdfunding chaincode: %v", err)
	}

	fmt.Println("Crowdfunding chaincode with Token SDK started successfully")
}
