// BoardWalk.go
	

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
	

type BoardWalk struct {
}
	
type INVOICE struct {
	BomNum string `json:”bom_num”`
	OrderState string `json:”order_state”`
}
	

type Transaction struct {
	TxINVOICE INVOICE
}
// ##################
//       INIT
// ##################
	
func (t *BoardWalk) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}
	

// ########################
//     Invocations
// ########################
	

func (t *BoardWalk) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
        // Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()
	fmt.Println("DEBUG: invoke is running " + fn)
	fmt.Println("DEBUG: args %+v",args)
	

	if fn == "onboardInvoice" {
		return t.onboardInvoice(stub, args)
	} 

	fmt.Println("invoke did not find function: " + fn)

	return shim.Error("Received unknown function invocation")
}
	

func (t *BoardWalk) onboardInvoice(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("- starting onboardInvoice")
	var invoiceJSON INVOICE
	var err error
	
	err = marshallRequest(args, &invoiceJSON)
	if err != nil { return shim.Error("Failed to marshall request" + err.Error())}
	
	var tx Transaction
	tx.TxINVOICE = invoiceJSON

	key, err := stub.CreateCompositeKey("txKey", []string{invoiceJSON.BomNum})
	if err != nil { return shim.Error(err.Error())}
	

	txAsBytes, err := json.Marshal(tx)
	if err != nil { return shim.Error(err.Error())}

	stub.PutState(key, txAsBytes)
	fmt.Println("- end onboardInvoice")
	return shim.Success(nil)
}
	

// ================
// UTILS
// ===============

func marshallRequest(args []string, invoice *INVOICE) error {
	var err error
	if len(args) != 2 {
                return errors.New("Incorrect number of arguments, expecting 2")
        }
	
	err = json.Unmarshal([]byte(args[0]), &invoice)
	if err != nil {	return err }
	fmt.Println("DEBUG: generated INVOICE %+v", invoice)
	
	return nil
}
	

// ================
// MAIN
// ================
	

func main() {
	if err := shim.Start(new(BoardWalk)); err != nil {
		fmt.Printf("Error starting BoardWalk chaincode: %s", err)
	}
}

