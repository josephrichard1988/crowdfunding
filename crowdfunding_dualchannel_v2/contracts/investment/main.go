package main

import (
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Combined INVESTMENT chaincode for investment-execution-channel.
// Orgs: StartupOrg, ValidatorOrg, InvestorOrg, PlatformOrg.
func main() {
	cc, err := contractapi.NewChaincode(
		&DataBridgeContract{},
		&StartupContract{},
		&InvestorContract{},
		&ValidatorContract{},
		&PlatformContract{},
		&TokenContract{},
	)
	if err != nil {
		log.Panicf("Error creating investment chaincode: %v", err)
	}

	cc.Info.Title = "Dual-Channel Investment Chaincode"
	cc.Info.Version = "1.0.0"

	if err := cc.Start(); err != nil {
		log.Panicf("Error starting investment chaincode: %v", err)
	}
	fmt.Println("Investment chaincode started successfully")
}
