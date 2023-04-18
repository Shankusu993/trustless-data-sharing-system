package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// ACL struct definition
type ACL struct {
	Id       	   string    `json:"id"`
	Identity       string    `json:"identity"`
	Identifier     string    `json:"identifier"`
	Qualifier      string    `json:"qualifier"`
	Validity       string `json:"validity"`
	Importance     int       `json:"importance"`
	MinBehavior    int       `json:"minBehavior"`
	DisputeSafeguard int      `json:"disputeSafeguard"`
}

// AccessControl contract for managing ACLs
type AccessControlContract struct {
	contractapi.Contract
}

// AddACL adds an ACL entry by storing it on the blockchain
func (ac *AccessControlContract) AddACL(ctx contractapi.TransactionContextInterface, Id, Identity, Identifier, Qualifier, Validity string) error {
	acl := ACL{
		Id:			Id,
		Identity:	Identity,
		Identifier:	Identifier,
		Qualifier:	Qualifier,
		Validity:	Validity, // UTC time
		Importance:	0,
		MinBehavior:	0,
		DisputeSafeguard:	0,
	}

	aclAsBytes, _ := json.Marshal(acl)
	return ctx.GetStub().PutState(Id, aclAsBytes)
}

// Updating the ACL on the blockchain
func (ac *AccessControlContract) UpdateACL(ctx contractapi.TransactionContextInterface, Id, Qualifier string) error {
	aclAsBytes, err := ctx.GetStub().GetState(Id)
	if err != nil {
		return err
	}

	acl := ACL{}
	err = json.Unmarshal(aclAsBytes, &acl)
	if err != nil {
		return err
	}

	acl.Qualifier = Qualifier

	aclAsBytes, _ = json.Marshal(acl)
	return ctx.GetStub().PutState(Id, aclAsBytes)
}


// GetACL retrieves an ACL entry from the blockchain
func (ac *AccessControlContract) GetACL(ctx contractapi.TransactionContextInterface, Id string) (*ACL, error) {
	aclAsBytes, err := ctx.GetStub().GetState(Id)
	if err != nil {
		return nil, err
	}

	acl := ACL{}
	err = json.Unmarshal(aclAsBytes, &acl)
	if err != nil {
		return nil, err
	}

	return &acl, nil
}


func main() {
	chaincode, err := contractapi.NewChaincode(new(AccessControlContract))
	if err != nil {
		fmt.Printf("Error create access control chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting access control chaincode: %s", err.Error())
	}
}