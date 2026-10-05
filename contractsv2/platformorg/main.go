package main

import (
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	chaincode, err := contractapi.NewChaincode(&PlatformContract{}, &TokenContract{})
	if err != nil {
		log.Panicf("Error creating platform chaincode: %v", err)
	}

	chaincode.Info.Title = "Platform Org Chaincode"
	chaincode.Info.Version = "1.0.0"

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting platform chaincode: %v", err)
	}
	fmt.Println("Platform chaincode started successfully")
}
