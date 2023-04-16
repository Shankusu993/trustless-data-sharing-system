package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Dispute struct for representing a dispute
type Dispute struct {
	DisputeID     string `json:"disputeID"`
	Raiser        string `json:"raiser"`
	Defendant     string `json:"defendant"`
	Description   string `json:"description"`
	Response      string `json:"response"`
	Status        string `json:"status"`
	ConfirmRaiser bool   `json:"confirmRaiser"`
}

// DisputeResolutionContract contract for managing disputes
type DisputeResolutionContract struct {
	contractapi.Contract
}

// RaiseDispute raises a dispute by storing it on the blockchain
func (dr *DisputeResolutionContract) RaiseDispute(ctx contractapi.TransactionContextInterface, disputeID, raiser, defendant, description string) error {
	dispute := Dispute{
		DisputeID:   disputeID,
		Raiser:      raiser,
		Defendant:   defendant,
		Description: description,
		Status:      "RAISED",
	}

	disputeAsBytes, _ := json.Marshal(dispute)
	return ctx.GetStub().PutState(disputeID, disputeAsBytes)
}

// RespondToDispute responds to a dispute by updating the dispute on the blockchain
func (dr *DisputeResolutionContract) RespondToDispute(ctx contractapi.TransactionContextInterface, disputeID, response string) error {
	disputeAsBytes, err := ctx.GetStub().GetState(disputeID)
	if err != nil {
		return err
	}

	dispute := Dispute{}
	err = json.Unmarshal(disputeAsBytes, &dispute)
	if err != nil {
		return err
	}

	dispute.Response = response
	dispute.Status = "RESPONDED"

	disputeAsBytes, _ = json.Marshal(dispute)
	return ctx.GetStub().PutState(disputeID, disputeAsBytes)
}

// GetDispute retrieves a dispute from the blockchain
func (dr *DisputeResolutionContract) GetDispute(ctx contractapi.TransactionContextInterface, disputeID string) (*Dispute, error) {
	disputeAsBytes, err := ctx.GetStub().GetState(disputeID)
	if err != nil {
		return nil, err
	}

	dispute := Dispute{}
	err = json.Unmarshal(disputeAsBytes, &dispute)
	if err != nil {
		return nil, err
	}

	return &dispute, nil
}

// ConfirmResolution confirms the resolution of a dispute by updating the dispute on the blockchain
func (dr *DisputeResolutionContract) ConfirmResolution(ctx contractapi.TransactionContextInterface, disputeID string) error {
	disputeAsBytes, err := ctx.GetStub().GetState(disputeID)
	if err != nil {
		return err
	}

	dispute := Dispute{}
	err = json.Unmarshal(disputeAsBytes, &dispute)
	if err != nil {
		return err
	}

	dispute.Status = "RESOLVED"
	dispute.ConfirmRaiser = true

	disputeAsBytes, _ = json.Marshal(dispute)
	return ctx.GetStub().PutState(disputeID, disputeAsBytes)
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(DisputeResolutionContract))
	if err != nil {
		fmt.Printf("Error create dispute resolution chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting dispute resolution chaincode: %s", err.Error())
	}
}
