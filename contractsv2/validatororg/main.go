package main

import (
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	chaincode, err := contractapi.NewChaincode(&ValidatorContract{})
	if err != nil {
		log.Panicf("Error creating validator chaincode: %v", err)
	}

	chaincode.Info.Title = "Validator Org Chaincode"
	chaincode.Info.Version = "1.0.0"

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting validator chaincode: %v", err)
	}
	fmt.Println("Validator chaincode started successfully")
}
