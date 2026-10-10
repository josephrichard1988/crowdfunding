package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Dual-channel combined chaincodes
//   governance → governance-validation-channel (Startup, Validator, Platform)
//   investment → investment-execution-channel  (all 4 orgs)
const (
	GovernanceChannel     = "governance-validation-channel"
	InvestmentChannel     = "investment-execution-channel"
	GovernanceChaincode   = "governance"
	InvestmentChaincode   = "investment"
	CurrentChaincodeName  = GovernanceChaincode
)

func getChannelForCollection(collection string) string {
	switch collection {
	case StartupPrivateCollection, ValidatorPrivateCollection,
		StartupValidatorCollection, StartupPlatformCollection, ValidatorPlatformCollection:
		return GovernanceChannel
	case InvestorPrivateCollection, PlatformPrivateCollection,
		StartupInvestorCollection, InvestorValidatorCollection, InvestorPlatformCollection,
		ThreePartyCollection, AllOrgsCollection:
		return InvestmentChannel
	default:
		return InvestmentChannel
	}
}

func getChaincodeForCollection(collection string) string {
	if getChannelForCollection(collection) == GovernanceChannel {
		return GovernanceChaincode
	}
	return InvestmentChaincode
}

func CrossChannelPut(ctx contractapi.TransactionContextInterface, collection string, key string, value []byte) error {
	channelName := getChannelForCollection(collection)
	chaincodeName := getChaincodeForCollection(collection)

	if channelName == ctx.GetStub().GetChannelID() && chaincodeName == CurrentChaincodeName {
		return ctx.GetStub().PutState(collection+"_"+key, value)
	}

	args := [][]byte{[]byte("DataBridgeContract:StoreData"), []byte(collection), []byte(key), value}
	response := ctx.GetStub().InvokeChaincode(chaincodeName, args, channelName)
	if response.Status != 200 {
		return fmt.Errorf("cross-channel put failed on %s/%s: %s", channelName, chaincodeName, response.Message)
	}
	return nil
}

func CrossChannelGet(ctx contractapi.TransactionContextInterface, collection string, key string) ([]byte, error) {
	channelName := getChannelForCollection(collection)
	chaincodeName := getChaincodeForCollection(collection)

	if channelName == ctx.GetStub().GetChannelID() && chaincodeName == CurrentChaincodeName {
		val, err := ctx.GetStub().GetState(collection + "_" + key)
		return val, err
	}

	args := [][]byte{[]byte("DataBridgeContract:ReadData"), []byte(collection), []byte(key)}
	response := ctx.GetStub().InvokeChaincode(chaincodeName, args, channelName)
	if response.Status != 200 {
		return nil, fmt.Errorf("cross-channel get failed on %s/%s: %s", channelName, chaincodeName, response.Message)
	}

	var unquoted string
	if err := json.Unmarshal(response.Payload, &unquoted); err == nil {
		return []byte(unquoted), nil
	}
	return response.Payload, nil
}
