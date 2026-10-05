package main

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Map collections to channels
func getChannelForCollection(collection string) string {
	switch collection {
	case StartupPrivateCollection:
		return "common-channel"
	case InvestorPrivateCollection:
		return "common-channel"
	case ValidatorPrivateCollection:
		return "common-channel"
	case PlatformPrivateCollection:
		return "common-channel"
	case StartupInvestorCollection:
		return "startup-investor-channel"
	case StartupValidatorCollection:
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
		return "startuporg"
	case InvestorPrivateCollection, InvestorValidatorCollection, InvestorPlatformCollection:
		return "investororg"
	case ValidatorPrivateCollection, ValidatorPlatformCollection:
		return "validatororg"
	case PlatformPrivateCollection, AllOrgsCollection:
		return "platformorg"
	default:
		return "platformorg"
	}
}

// CrossChannelPut replaces PutPrivateData
func CrossChannelPut(ctx contractapi.TransactionContextInterface, collection string, key string, value []byte) error {
	channelName := getChannelForCollection(collection)
	chaincodeName := getChaincodeForCollection(collection)

	args := [][]byte{[]byte("StoreData"), []byte(key), value}
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

	args := [][]byte{[]byte("ReadData"), []byte(key)}
	response := ctx.GetStub().InvokeChaincode(chaincodeName, args, channelName)
	if response.Status != 200 {
		return nil, fmt.Errorf("cross-channel get failed on %s/%s: %s", channelName, chaincodeName, response.Message)
	}
	return response.Payload, nil
}
