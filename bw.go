package main

import (
"fmt"
"encoding/json"
"github.com/hyperledger/fabric/core/chaincode/shim"
"github.com/hyperledger/fabric/protos/peer"
"errors"
)
// #########################
//    Type Definitions
// #########################

type AuraBlock struct {
}

type Loan struct {
	LoanId string `json:"loan_id"`
	Type string `json:"loan_type"`
}

type Transaction struct {
	TxLoan Loan
}
// ##################
//       INIT
// ##################


func (t *AuraBlock) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

// ########################
//     Invocations
// ########################


func (t *AuraBlock) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
        // Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()
	fmt.Println("DEBUG: invoke is running " + fn)
	fmt.Println("DEBUG: args %+v",args)

	if fn == "onboardLoan" {
		return t.onboardLoan(stub, args)
	} 

	fmt.Println("invoke did not find function: " + fn)

	return shim.Error("Recieved unknown function invocation")
}

func (t *AuraBlock) onboardLoan(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("- starting onboardLoan")
	
	var loanJSON Loan
	var err error

	err = marshallRequest(args, &loanJSON)
	if err != nil { return shim.Error("Failed to marshall request" + err.Error())}


	var tx Transaction
	tx.TxLoan = loanJSON

	key, err := stub.CreateCompositeKey("txKey", []string{loanJSON.LoanId, loanJSON.Type})
	if err != nil { return shim.Error(err.Error())}

	txAsBytes, err := json.Marshal(tx)
	if err != nil { return shim.Error(err.Error())}

	stub.PutState(key, txAsBytes)
	fmt.Println("- end onboardLoan")
	return shim.Success(nil)
}


// ================
// UTILS
// ===============

func marshallRequest(args []string, loan *Loan) error {
	var err error
	if len(args) != 1 {
                return errors.New("Incorrect number of arguments, expecting 1")
        }

	err = json.Unmarshal([]byte(args[2]), &loan)
	if err != nil {	return err }
	fmt.Println("DEBUG: generated loan %+v", loan)

	return nil
}

// ================
// MAIN
// ================

func main() {
	if err := shim.Start(new(AuraBlock)); err != nil {
		fmt.Printf("Error starting AuraBlock chaincode: %s", err)
	}
}
