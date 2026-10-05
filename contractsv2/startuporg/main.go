package main

import (
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	chaincode, err := contractapi.NewChaincode(&StartupContract{})
	if err != nil {
		log.Panicf("Error creating startup chaincode: %v", err)
	}

	chaincode.Info.Title = "Startup Org Chaincode"
	chaincode.Info.Version = "1.0.0"

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting startup chaincode: %v", err)
	}
	fmt.Println("Startup chaincode started successfully")
}
