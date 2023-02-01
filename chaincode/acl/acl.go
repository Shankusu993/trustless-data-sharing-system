package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// ACL struct definition
type ACL struct {
	Identity       string    `json:"identity"`
	Identifier     string    `json:"identifier"`
	Qualifier      string    `json:"qualifier"`
	Validity       time.Time `json:"validity"`
	Importance     int       `json:"importance"`
	MinBehavior    int       `json:"minBehavior"`
	DisputeSafeguard int      `json:"disputeSafeguard"`
}

// SimpleChaincode struct definition
type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "addACL" {
		return t.addACL(stub, args)
	} else if function == "getACL" {
		return t.getACL(stub, args)
	}
	return shim.Error("Invalid function name")
}

// addACL function definition
func (t *SimpleChaincode) addACL(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 8 {
		return shim.Error("Incorrect number of arguments. Expecting 8")
	}
	
	// extract arguments
	identity := args[0]
	identifier := args[1]
	qualifier := args[2]
	validity, err := strconv.ParseInt(args[3], 10, 64)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to parse validity timestamp: %s", err))
	}
	validityTime := time.Unix(validity, 0)
	importance, err := strconv.Atoi(args[4])
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to parse importance: %s", err))
	}
	minBehavior, err := strconv.Atoi(args[5])
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to parse minimum behavior score: %s", err))
	}
	disputeSafeguard, err := strconv.Atoi(args[6])
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to parse dispute safeguard: %s", err))
	}

	// create new ACL
	acl := &ACL{Identity: identity, Identifier: identifier, Qualifier: qualifier, Validity: validityTime,
			Importance: importance, MinBehavior: minBehavior, DisputeSafeguard: disputeSafeguard}
	// convert ACL to bytes
	aclBytes, err := json.Marshal(acl)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to marshal ACL: %s", err))
	}

	// save ACL to blockchain
	err = stub.PutState(identity, aclBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to save ACL: %s", err))
	}

	return shim.Success(nil)
}

// getACL function definition
func (t *SimpleChaincode) getACL(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
	return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	// get identity from arguments
	identity := args[0]

	// get ACL from blockchain
	aclBytes, err := stub.GetState(identity)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to get ACL: %s", err))
	} else if aclBytes == nil {
		return shim.Error("ACL not found")
	}

	return shim.Success(aclBytes)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
	fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

