/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

//WARNING - this chaincode's ID is hard-coded in chaincode_example04 to illustrate one way of
//calling chaincode from a chaincode. If this example is modified, chaincode_example04.go has
//to be modified as well with the new ID of chaincode_example02.
//chaincode_example05 show's how chaincode ID can be passed in as a parameter instead of
//hard-coding.

import (
	"fmt"
	"bytes"
	"strconv"
	"encoding/json"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type LOM struct {
	objectType string `json:"docType"` //docType is used to distinguish the various types of objects in state database
	ID         string `json:"id"`    //the fieldtags are needed to keep case from bouncing around
	Owner      string `json:"owner"`
	Balance    int    `json:"balance"`
	Limit      int    `json:"limit"`
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	/*
	fmt.Println("ex02 Init")
	_, args := stub.GetFunctionAndParameters()
	var A string    // Entities
	var Aval,i int
	var err error // Asset holdings

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	for i=0; i<len(args); i+=2 {

		// ==== Get account id ====
		A = args[i]

		// ==== Get account balance ====
		Aval, err = strconv.Atoi(args[i+1])
		if err != nil {
			return shim.Error("Expecting integer value for asset holding")
		}

		err = createAccount(stub,A,"NA",Aval,-99999)
		if err != nil {
			return shim.Error("Error creating account")
		}

	}
	*/
	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "transfer" {
		// Make payment of X units from A to B
		return t.transfer(stub, args)
	} else if function == "delete" {
		// Deletes an entity from its state
		return t.delete(stub, args)
	} else if function == "account" {
		// Deletes an entity from its state
		return t.account(stub, args)
	} else if function == "list" {
		// Deletes an entity from its state
		return t.list(stub, args)
	} else if function == "query" {
		// the old "Query" is now implemtned in invoke
		return t.query(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"delete\" \"query\"")
}

func getAccount (stub shim.ChaincodeStubInterface, accountID string) (LOM, error) {
	// Get the state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	accountAsBytes, err := stub.GetState(accountID)
	account := LOM{}

	if err != nil {
		return account, err
	}
	if accountAsBytes == nil {
		return account, errors.New("Entity not found")
	}
	
	err = json.Unmarshal(accountAsBytes, &account) //unmarshal it aka JSON.parse()
	if err != nil {
		return account, err
	}
	return account, nil
}

// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) transfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var accountSender, accountRecipient LOM    // Entities
	var transferValue int          // Transaction value
	var err error

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	accountSender, err = getAccount(stub,args[0])
	if err != nil {
		shim.Error("Could not get sender's account: "+args[0])
	}

	accountRecipient, err = getAccount(stub,args[1])
	if err != nil {
		shim.Error("Could not get receiver's account: "+args[1])
	}

	transferValue, err = strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Invalid transaction amount, expecting a integer value")
	}

	// Business logic
	accountSender.Balance = accountSender.Balance - transferValue
	accountRecipient.Balance = accountRecipient.Balance + transferValue
	fmt.Printf("Account Sender Balance :%d\n", accountSender.Balance)
	fmt.Printf("Account Receipient Balance :%d\n", accountRecipient.Balance)

	// Write the state back to the ledger
	accountJSONasBytes, _ := json.Marshal(accountSender)
	err = stub.PutState(accountSender.ID, accountJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	accountJSONasBytes, _ = json.Marshal(accountRecipient)
	err = stub.PutState(accountRecipient.ID, accountJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Printf("========================================")
	fmt.Printf("Transaction of %s to %s", accountSender.ID,accountRecipient.ID)
	fmt.Printf("Sender's balance = %d, Recipient's balance = %d\n", accountSender.Balance, accountRecipient.Balance)
	fmt.Printf("========================================")
	return shim.Success([]byte("Transaction was completed OK !"))
}

func (t *SimpleChaincode) account(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("Creating new Account")
	//_, args := stub.GetFunctionAndParameters()
	var accountID string    // Entities
	var balance int // Asset holdings
	var err error
	var lowerLimit int

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4 \n Args = accountID, Owner, Balance, Lower Limit")
	}

	// Initialize the chaincode
	accountID = args[0]
	owner := args[1]
	balance, err = strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	lowerLimit, err = strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("Expecting integer value for lower limit")
	}

	err = createAccount(stub,accountID,owner,balance,lowerLimit)

	if err != nil {
		return shim.Error("Error creating account"+err.Error())
	}

	fmt.Printf("\n========================================")
	fmt.Printf("\nNew account created for %s ", accountID)
	fmt.Printf("\nBalnce = %d\n", balance)
	fmt.Printf("========================================")
	return shim.Success([]byte("Account was created successfully! OK !"))
}


// Deletes an entity from state
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	A := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) list(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var outstring string

	keyStart := ""
	keyEnd := ""

	if len(args) >= 1 {
		keyStart = args[0]
	}
	if len(args) >= 2 {
		keyEnd = args[1]
	}

	outstring = "\nFollowing accounts are registered:\n"
	res, err := stub.GetStateByRange(keyStart, keyEnd)

	if err != nil {
		return shim.Error(err.Error())
	}
	kv, err1 := res.Next();
	if err1 != nil {
		return shim.Error(err1.Error())
	}
	outstring += " - Account: " + kv.GetKey() + "\n"

	for res.HasNext() {
		kv, err1 = res.Next();
		if err1 != nil {
			return shim.Error(err1.Error())
		}
		outstring += " - Account: " + kv.GetKey() + "\n"
	}

	res.Close();
	return shim.Success([]byte(outstring))
}

func (t *SimpleChaincode) listCouchDB(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"LOM\"}}")

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}


func createAccount(stub shim.ChaincodeStubInterface, accountID string, owner string, balance int, limit int ) error {

		// ==== Check if LOM already exists ====
		accountAsBytes, err := stub.GetState(accountID)
		if err != nil {
			return err
		} else if accountAsBytes != nil {
			fmt.Println("This account already exists: " + accountID)
			return errors.New("Account already exists:" + accountID)
		}

		// ==== Create marble object and marshal to JSON ====
		objectType := "LOM"
		lom := &LOM{objectType, accountID, owner, balance, limit}
		lomJSONasBytes, err := json.Marshal(lom)
		if err != nil {
			return err
		}

		// === Save Account to state ===
		err = stub.PutState(accountID, lomJSONasBytes)
		if err != nil {
			return err
		}

		return nil
}

// query callback representing the query of a chaincode
func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var name string // Entities
	var err error
	var jsonResp string

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the marble to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetState(name) //get the LOM from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Marble does not exist: " + name + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)

}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
