package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Dual-channel topology (MICROFAB.txt):
//   governance-validation-channel  — StartupOrg, ValidatorOrg, PlatformOrg
//   investment-execution-channel   — StartupOrg, ValidatorOrg, InvestorOrg, PlatformOrg
const (
	GovernanceChannel = "governance-validation-channel"
	InvestmentChannel = "investment-execution-channel"
)

// getChannelForCollection maps logical PDC names to the dual-channel topology.
func getChannelForCollection(collection string) string {
	switch collection {
	case StartupPrivateCollection:
		return GovernanceChannel
	case ValidatorPrivateCollection:
		return GovernanceChannel
	case StartupValidatorCollection:
		return GovernanceChannel
	case StartupPlatformCollection:
		return GovernanceChannel
	case ValidatorPlatformCollection:
		return GovernanceChannel

	case InvestorPrivateCollection:
		return InvestmentChannel
	case PlatformPrivateCollection:
		return InvestmentChannel
	case StartupInvestorCollection:
		return InvestmentChannel
	case InvestorValidatorCollection:
		return InvestmentChannel
	case InvestorPlatformCollection:
		return InvestmentChannel
	case ThreePartyCollection:
		return InvestmentChannel
	case AllOrgsCollection:
		return InvestmentChannel

	default:
		return InvestmentChannel
	}
}

func getChaincodeForCollection(collection string) string {
	switch collection {
	case StartupPrivateCollection, StartupInvestorCollection, StartupValidatorCollection, StartupPlatformCollection:
		return "startup"
	case InvestorPrivateCollection, InvestorValidatorCollection, InvestorPlatformCollection:
		return "investor"
	case ValidatorPrivateCollection, ValidatorPlatformCollection:
		return "validator"
	case PlatformPrivateCollection, AllOrgsCollection, ThreePartyCollection:
		return "platform"
	default:
		return "platform"
	}
}

const CurrentChaincodeName = "investor"

func CrossChannelPut(ctx contractapi.TransactionContextInterface, collection string, key string, value []byte) error {
	channelName := getChannelForCollection(collection)
	chaincodeName := getChaincodeForCollection(collection)

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

func CrossChannelGet(ctx contractapi.TransactionContextInterface, collection string, key string) ([]byte, error) {
	channelName := getChannelForCollection(collection)
	chaincodeName := getChaincodeForCollection(collection)

	if channelName == ctx.GetStub().GetChannelID() && chaincodeName == CurrentChaincodeName {
		val, err := ctx.GetStub().GetState(collection + "_" + key)
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
