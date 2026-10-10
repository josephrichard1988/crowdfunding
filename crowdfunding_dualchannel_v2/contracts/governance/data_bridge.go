package main

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// DataBridgeContract provides StoreData/ReadData for cross-channel state bridging.
// Org contracts keep their own StoreData for compatibility; CrossChannelPut targets this bridge
// via "DataBridgeContract:StoreData" / "DataBridgeContract:ReadData".
type DataBridgeContract struct {
	contractapi.Contract
}

func (d *DataBridgeContract) StoreData(ctx contractapi.TransactionContextInterface, collection string, key string, value string) error {
	err := ctx.GetStub().PutState(collection+"_"+key, []byte(value))
	if err != nil {
		return fmt.Errorf("DataBridge StoreData failed: %v", err)
	}
	return nil
}

func (d *DataBridgeContract) ReadData(ctx contractapi.TransactionContextInterface, collection string, key string) (string, error) {
	val, err := ctx.GetStub().GetState(collection + "_" + key)
	if err != nil {
		return "", fmt.Errorf("DataBridge ReadData failed: %v", err)
	}
	if val == nil {
		return "", fmt.Errorf("key not found: %s_%s", collection, key)
	}
	return string(val), nil
}
