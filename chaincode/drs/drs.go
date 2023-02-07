package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// DisputeResolutionChaincode implementation of Chaincode
type DisputeResolutionChaincode struct {
}

// Dispute structure definition
type Dispute struct {
	DisputeID     string `json:"disputeID"`
	Raiser        string `json:"raiser"`
	Defendant     string `json:"defendant"`
	Description   string `json:"description"`
	DefendantResp string `json:"defendantResp"`
	Status        string `json:"status"`
}

// Init function
func (t *DisputeResolutionChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke function
func (t *DisputeResolutionChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "raiseDispute" {
		return t.raiseDispute(stub, args)
	} else if function == "respondToDispute" {
		return t.respondToDispute(stub, args)
	} else if function == "getDispute" {
		return t.getDispute(stub, args)
	} else if function == "confirmResolution" {
		return t.confirmResolution(stub, args)
	}

	return shim.Error("Invalid function name.")
}

// function to get Dispute
func (t *SimpleChaincode) getDispute(stub shim.ChaincodeStubInterface, args []string) peer.Response {
    // Check if the number of arguments passed is correct
    if len(args) != 1 {
        return shim.Error("Incorrect number of arguments. Expecting dispute ID")
    }
    disputeID := args[0]
    
    // Get the dispute from the ledger
    disputeBytes, err := stub.GetState(disputeID)
    if err != nil {
        return shim.Error(fmt.Sprintf("Failed to get dispute with ID %s: %s", disputeID, err))
    }
    if disputeBytes == nil {
        return shim.Error(fmt.Sprintf("Dispute with ID %s not found", disputeID))
    }
    
    // Unmarshal the dispute into a struct
    var dispute Dispute
    err = json.Unmarshal(disputeBytes, &dispute)
    if err != nil {
        return shim.Error(fmt.Sprintf("Failed to unmarshal dispute with ID %s: %s", disputeID, err))
    }
    
    // Return the dispute
    return shim.Success(disputeBytes)
}

// function to raise a Dispute
func (t *DisputeResolution) raiseDispute(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4 arguments: dispute ID, dispute raiser, dispute defendant, and dispute reason")
	}

	disputeID := args[0]
	disputeRaiser := args[1]
	disputeDefendant := args[2]
	disputeReason := args[3]
	disputeResponse := ""

	dispute := Dispute{disputeID, disputeRaiser, disputeDefendant, disputeReason, disputeResponse, "OPEN"}

	disputeAsBytes, err := json.Marshal(dispute)
	if err != nil {
		return shim.Error("Error marshalling dispute")
	}

	err = stub.PutState(disputeID, disputeAsBytes)
	if err != nil {
		return shim.Error("Error putting state for dispute")
	}

	return shim.Success(nil)
}

// function to respond to a Dispute
func (t *SimpleChaincode) respondToDispute(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3: dispute ID, response, and defendant ID")
	}

	disputeID := args[0]
	response := args[1]
	defendant := args[2]

	disputeAsBytes, err := stub.GetState(disputeID)
	if err != nil {
		return shim.Error("Failed to get dispute: " + err.Error())
	} else if disputeAsBytes == nil {
		return shim.Error("Dispute does not exist")
	}

	dispute := Dispute{}
	err = json.Unmarshal(disputeAsBytes, &dispute)
	if err != nil {
		return shim.Error("Failed to decode dispute JSON: " + err.Error())
	}

	if dispute.Defendant != defendant {
		return shim.Error("The defendant does not match the dispute record")
	}

	dispute.Response = response
	disputeAsBytes, err = json.Marshal(dispute)
	if err != nil {
		return shim.Error("Failed to encode dispute JSON: " + err.Error())
	}

	err = stub.PutState(disputeID, disputeAsBytes)
	if err != nil {
		return shim.Error("Failed to update dispute: " + err.Error())
	}

	return shim.Success(nil)
}

//function to confirm resolution
func (t *DisputeResolution) confirmResolution(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1 (disputeID)")
	}
	disputeID := args[0]
	dispute, err := t.getDispute(stub, disputeID)
	if err != nil {
		return shim.Error(err.Error())
	}
	// Check that the caller is the Raiser of the dispute
	caller, err := stub.GetCreator()
	if err != nil {
		return shim.Error(err.Error())
	}
	if !bytes.Equal(caller, dispute.Raiser) {
		return shim.Error("Caller is not authorized to confirm the resolution of this dispute")
	}
	// Update the status of the dispute to "Resolved"
	dispute.Status = "Resolved"
	err = t.putDispute(stub, dispute)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}


func main() {
	err := shim.Start(new(DisputeResolutionChaincode))
	if err != nil {
		fmt.Printf("Error starting DisputeResolutionChaincode chaincode: %s", err)
	}
}
