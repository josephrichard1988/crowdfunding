package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Map collections to channels
func getChannelForCollection(collection string) string {
	switch collection {
	case StartupPrivateCollection:
		return "startup-validator-channel"
	case InvestorPrivateCollection:
		return "common-channel"
	case PlatformPrivateCollection:
		return "common-channel"
	case StartupInvestorCollection:
		return "startup-investor-channel"
	case StartupValidatorCollection, ValidatorPrivateCollection:
		return "startup-validator-channel"
	case StartupPlatformCollection:
		return "startup-platform-channel"
	case InvestorValidatorCollection:
		return "investor-validator-channel"
	case InvestorPlatformCollection:
		return "investor-platform-channel"
	case ValidatorPlatformCollection:
		return "validator-platform-channel"
	case AllOrgsCollection:
		return "common-channel"
	default:
		return "common-channel"
	}
}

// Map collections to the target chaincode name
func getChaincodeForCollection(collection string) string {
	switch collection {
	case StartupPrivateCollection, StartupInvestorCollection, StartupValidatorCollection, StartupPlatformCollection:
		return "startup"
	case InvestorPrivateCollection, InvestorValidatorCollection, InvestorPlatformCollection:
		return "investor"
	case ValidatorPrivateCollection, ValidatorPlatformCollection:
		return "validator"
	case PlatformPrivateCollection, AllOrgsCollection:
		return "platform"
	default:
		return "platform"
	}
}

const CurrentChaincodeName = "validator"

// CrossChannelPut replaces PutPrivateData
func CrossChannelPut(ctx contractapi.TransactionContextInterface, collection string, key string, value []byte) error {
	channelName := getChannelForCollection(collection)
	chaincodeName := getChaincodeForCollection(collection)

	// If we are invoking our own chaincode on the current channel, use direct state access to avoid self-invocation deadlock
	if channelName == ctx.GetStub().GetChannelID() && chaincodeName == CurrentChaincodeName {
		return ctx.GetStub().PutState(collection+"_"+key, value)
	}

	args := [][]byte{[]byte("StoreData"), []byte(collection), []byte(key), value}
	response := ctx.GetStub().InvokeChaincode(chaincodeName, args, channelName)
	if response.Status != 200 {
		return fmt.Errorf("cross-channel put failed on %s/%s: %s", channelName, chaincodeName, response.Message)
	}
	return nil
}

// CrossChannelGet replaces GetPrivateData
func CrossChannelGet(ctx contractapi.TransactionContextInterface, collection string, key string) ([]byte, error) {
	channelName := getChannelForCollection(collection)
	chaincodeName := getChaincodeForCollection(collection)

	// If we are invoking our own chaincode on the current channel, use direct state access to avoid self-invocation deadlock
	if channelName == ctx.GetStub().GetChannelID() && chaincodeName == CurrentChaincodeName {
		val, err := ctx.GetStub().GetState(collection+"_"+key)
		return val, err
	}

	args := [][]byte{[]byte("ReadData"), []byte(collection), []byte(key)}
	response := ctx.GetStub().InvokeChaincode(chaincodeName, args, channelName)
	if response.Status != 200 {
		return nil, fmt.Errorf("cross-channel get failed on %s/%s: %s", channelName, chaincodeName, response.Message)
	}

	var unquoted string
	err := json.Unmarshal(response.Payload, &unquoted)
	if err == nil {
		return []byte(unquoted), nil
	}

	return response.Payload, nil
}
