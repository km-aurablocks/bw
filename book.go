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
        

type library struct {
}
        
type book struct {
        BookNum string `json:"book_num"`
        Author string `json:"author"`
        OrderCount float64 `json:"order_count,string"`
        Printer string `json:"printer"`
        Amount float64 `json:"amount,string"`
        OrderDate time.Time `json:"completed_date"`           
}
        

type Transaction struct {
        Txbook book
}
// ##################
//       INIT
// ##################
        
func (t *library) Init(stub shim.ChaincodeStubInterface) peer.Response {
        return shim.Success(nil)
}
        

// ########################
//     Invocations
// ########################
        

func (t *library) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
        // Extract the function and args from the transaction proposal
        fn, args := stub.GetFunctionAndParameters()
        fmt.Println("DEBUG: invoke is running " + fn)
        fmt.Println("DEBUG: args %+v",args)
        

	if fn == "addbook" {
                return t.addbook(stub, args)
        } 

        fmt.Println("invoke did not find function: " + fn)

        return shim.Error("Received unknown function invocation")
}
        

func (t *library) addbook(stub shim.ChaincodeStubInterface, args []string) peer.Response {
        fmt.Println("- starting addbook")
        var bookJSON book
        var err error
        

        err = marshallRequest(args, &bookJSON)
        if err != nil { return shim.Error("Failed to marshall request" + err.Error())}
        
        var tx Transaction
        tx.Txbook = bookJSON

        key, err := stub.CreateCompositeKey("txKey", []string{bookJSON.BookNum, bookJSON.Author})
        if err != nil { return shim.Error(err.Error())}
        

        txAsBytes, err := json.Marshal(tx)
        if err != nil { return shim.Error(err.Error())}

        stub.PutState(key, txAsBytes)
        fmt.Println("- end addbook")
        return shim.Success(nil)
}
        

// ================
// UTILS
// ===============

func marshallRequest(args []string, book *book) error {
        var err error
        if len(args) != 3 {
                return errors.New("Incorrect number of arguments, expecting 3")
        }
        

        err = json.Unmarshal([]byte(args[0]), &book)
        if err != nil { return err }
        fmt.Println("DEBUG: generated book %+v", book)
        

        return nil
}
        

// ================
// MAIN
// ================
        

func main() {
        if err := shim.Start(new(library)); err != nil {
                fmt.Printf("Error starting library chaincode: %s", err)
        }
}
