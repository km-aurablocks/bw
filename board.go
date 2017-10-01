// BoardWalks.go
	

package main
	

import (
"fmt"
"time"
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
	Client string `json:”client”`
	OrderItem string `json:”order_item”`
	Supplier string `json:”supplier”`
	InvoiceNum string `json:”invoice_num”`
	OrderDetail string `json:”order_detail”`
	OrderState string `json:”order_state”`
	Amount float64 `json:”amount,string"`
	CompletedDate time.Time `json:”completed_date”`
	DaysRemaining float64 `json:”daysremaining,string"`
	Comments string `json:”comments”`		
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
	} else if fn == "getInvoice" {
		return t.getInvoice(stub, args)
	} else if fn == "updateInvoice" {
		return t.updateInvoice(stub, args)
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

	key, err := stub.CreateCompositeKey("txKey", []string{invoiceJSON.BomNum, invoiceJSON.Client})
	if err != nil { return shim.Error(err.Error())}
	

	txAsBytes, err := json.Marshal(tx)
	if err != nil { return shim.Error(err.Error())}

	stub.PutState(key, txAsBytes)
	fmt.Println("- end onboardInvoice")
	return shim.Success(nil)
}
	


func (t *BoardWalk) getInvoice(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("- starting getInvoice")
	var invoiceQuery INVOICE
	var err error
	
	err = marshallRequest(args, &invoiceQuery)
        if err != nil { return shim.Error("Failed to marshall request: " + err.Error())}
	
	key, err := stub.CreateCompositeKey("txKey", []string{invoiceQuery.BomNum, invoiceQuery.Client})
	if err != nil { return shim.Error(err.Error())}
	
	txBytes, err  := stub.GetState(key)
	if err != nil {
		return shim.Error("Failed to get tx: " + err.Error())
	} else if txBytes == nil {
		return shim.Error("Tx does not exist. ")
	}
	

	var tx Transaction
	err = json.Unmarshal(txBytes, &tx)
	if err != nil { return shim.Error(err.Error()) }	

	txBytesOut, err := json.Marshal(tx)
	if err != nil { return shim.Error(err.Error())}
	

	fmt.Println("- end getInvoice")
	return shim.Success(txBytesOut)
}
	


func (t *BoardWalk) updateInvoice(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("- starting updateInvoice")


	var invoiceQuery INVOICE
	var err error

	
	err = marshallRequest(args, &invoiceQuery)
        if err != nil { return shim.Error("Failed to marshall request: " + err.Error())}
	

	key, err := stub.CreateCompositeKey("txKey", []string{invoiceQuery.BomNum, invoiceQuery.Client})
        if err != nil { return shim.Error(err.Error())}
	

	txBytes, err  := stub.GetState(key)
	if err != nil { return shim.Error(err.Error())}
	

	tx := Transaction{}
	err = json.Unmarshal(txBytes, &tx)
	if err != nil { return shim.Error(err.Error()) }

	tx.TxINVOICE.OrderItem = invoiceQuery.OrderItem
	tx.TxINVOICE.Supplier = invoiceQuery.Supplier
	tx.TxINVOICE.InvoiceNum = invoiceQuery.InvoiceNum
	tx.TxINVOICE.OrderDetail = invoiceQuery.OrderDetail
	tx.TxINVOICE.OrderState = invoiceQuery.OrderState
	tx.TxINVOICE.Amount = invoiceQuery.Amount
	tx.TxINVOICE.CompletedDate = invoiceQuery.CompletedDate
	tx.TxINVOICE.DaysRemaining = invoiceQuery.DaysRemaining
	tx.TxINVOICE.Comments = invoiceQuery.Comments


	txAsBytes, err := json.Marshal(tx)
	if err != nil { return shim.Error(err.Error())}
	

	stub.PutState(key, txAsBytes)
		
	fmt.Println("- end updateCreditReceipts")
	return shim.Success(nil)
}
	
	

// ================
// UTILS
// ===============

func marshallRequest(args []string, invoice *INVOICE) error {
	var err error
	if len(args) != 3 {
                return errors.New("Incorrect number of arguments, expecting 3")
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

