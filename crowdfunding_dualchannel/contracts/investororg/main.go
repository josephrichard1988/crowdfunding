package main

import (
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	chaincode, err := contractapi.NewChaincode(&InvestorContract{})
	if err != nil {
		log.Panicf("Error creating investor chaincode: %v", err)
	}

	chaincode.Info.Title = "Investor Org Chaincode"
	chaincode.Info.Version = "1.0.0"

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting investor chaincode: %v", err)
	}
	fmt.Println("Investor chaincode started successfully")
}
