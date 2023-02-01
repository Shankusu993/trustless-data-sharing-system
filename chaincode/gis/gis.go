package main

import (
    "encoding/json"
    "fmt"

    "github.com/hyperledger/fabric/core/chaincode/shim"
    pb "github.com/hyperledger/fabric/protos/peer"
)

// GeoJSONData struct definition
type GeoJSONData struct {
    Geometry   json.RawMessage `json:"geometry"`
    Properties json.RawMessage `json:"properties"`
}

// SimpleChaincode struct definition
type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
    return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
    function, args := stub.GetFunctionAndParameters()
    if function == "addGeoJSONData" {
        return t.addGeoJSONData(stub, args)
    } else if function == "getGeoJSONData" {
        return t.getGeoJSONData(stub, args)
    }
    return shim.Error("Invalid function name")
}

// addGeoJSONData function definition
func (t *SimpleChaincode) addGeoJSONData(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    if len(args) != 3 {
        return shim.Error("Incorrect number of arguments. Expecting 3")
    }
    
    // extract arguments
    collection := args[0]
    key := args[1]
    value := args[2]
    
    // create new GeoJSONData object
    data := &GeoJSONData{Geometry: json.RawMessage(value)}
    
    // convert GeoJSONData object to bytes
    dataBytes, err := json.Marshal(data)
    if err != nil {
        return shim.Error(fmt.Sprintf("Failed to marshal GeoJSON data: %s", err))
    }
    
    // save data to private data collection
    err = stub.PutPrivateData(collection, key, dataBytes)
    if err != nil {
        return shim.Error(fmt.Sprintf("Failed to save data to collection: %s", err))
    }
    
    return shim.Success(nil)
}

// getGeoJSONData function definition
func (t *SimpleChaincode) getGeoJSONData(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	// extract arguments
	collection := args[0]
	key := args[1]

	// get data from private data collection
	dataBytes, err := stub.GetPrivateData(collection, key)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to get data from collection: %s", err))
	} else if dataBytes == nil {
		return shim.Error("Data not found")
	}

	// create new GeoJSONData object
	data := &GeoJSONData{}

	// unmarshal data bytes
	err = json.Unmarshal(dataBytes, data)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to unmarshal data: %s", err))
	}

	// return data
	return shim.Success(data.Geometry)
}
