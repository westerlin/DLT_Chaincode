/*
Based on IBMs own Hyperledger Fabric Example02
March 2018
*/

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// LOM account data structure
type LOM struct {
	ObjectType string `json:"docType"` //docType is used to distinguish the various types of objects in state database
	ID         string `json:"id"`      //the fieldtags are needed to keep case from bouncing around
	Owner      string `json:"owner"`
	Balance    int    `json:"balance"`
	Limit      int    `json:"limit"`
}

// Init function runs when chaincode initializes
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke function runs whenever changes are made to the database - primary interface for all calls
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "transfer" {
		// Make tranfer between two accounts
		return t.transfer(stub, args)
	} else if function == "delete" {
		// Deletes an entity from its state
		return t.delete(stub, args)
	} else if function == "add" {
		// creates an account
		return t.createAccount(stub, args)
	} else if function == "history" {
		// list all states for one account
		return t.getHistory(stub, args)
	} else if function == "list" {
		// list all accounts
		return t.list(stub, args)
	} else if function == "whoami" {
		// list all accounts
		return t.whoAmI(stub, args)
	} else if function == "query" {
		// the old "Query" is now implemtned in invoke
		return t.query(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"transfer\" \"add\" \"delete\" \"query\"  \"list\" \"history\"")
}

func (t *SimpleChaincode) whoAmI(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	creator, err := stub.GetCreator()
	if err != nil {
		return shim.Error("Error while getting creator:" + err.Error())
	}
	return shim.Success(creator)
}

func (t *SimpleChaincode) getHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var normaltime time.Time

	if len(args) != 1 {
		return shim.Error("Require one argument")
	}

	output := "\n{\"history\": ["

	history, err := stub.GetHistoryForKey(args[0])
	if err != nil {
		return shim.Error("could not get History" + ":" + err.Error())
	}

	stateHistory, err := history.Next()
	if err != nil {
		return shim.Error("could not get state story" + ":" + err.Error())
	}

	normaltime = time.Unix(stateHistory.Timestamp.Seconds, int64(stateHistory.Timestamp.Nanos)).UTC()
	//normaltime, _ := Timestamp(stateHistory.Timestamp)
	output += "\n\t{\"value\":" + string(stateHistory.Value) + ",\"timestamp\":" + normaltime.Format("2006-01-02 15:04:05") + ","
	if stateHistory.IsDelete {
		output += "\"Deleted\":True"
	} else {
		output += "\"Deleted\":False"
	}
	for history.HasNext() {

		output += "},"
		stateHistory, err = history.Next()
		if err != nil {
			return shim.Error("could not get state story" + ":" + err.Error())
		}
		normaltime = time.Unix(stateHistory.Timestamp.Seconds, int64(stateHistory.Timestamp.Nanos)).UTC()
		//normaltime, _ = Timestamp(stateHistory.Timestamp)
		output += "\n\t{\"value\":" + string(stateHistory.Value) + ",\"timestamp\":" + normaltime.Format("2006-01-02 15:04:05") + ","
		if stateHistory.IsDelete {
			output += "\"Deleted\":True"
		} else {
			output += "\"Deleted\":False"
		}
	}
	output += "}\n\t]\n}"
	return shim.Success([]byte(output))

}

func getAccount(stub shim.ChaincodeStubInterface, accountID string) (LOM, error) {
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

/*
func errcheck(err error, message string) {
	if err != nil {
		shim.Error(message + ":" + err.Error())
	}
}
*/
// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) transfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var accountSender, accountRecipient LOM // Entities
	var transferValue int                   // Transaction value
	var err error

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	accountSender, err = getAccount(stub, args[0])
	if err != nil {
		return shim.Error("Could not get senders '" + args[0] + "' account:" + err.Error())
	}
	accountRecipient, err = getAccount(stub, args[1])
	if err != nil {
		return shim.Error("Could not get receivers '" + args[0] + "' account:" + err.Error())
	}
	transferValue, err = strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Invalid transfer amount. Expecting integer value" + ":" + err.Error())
	}

	if accountSender.Balance-transferValue < accountSender.Limit {
		fmt.Printf("Account holder %s are breaching lower limit at %d width balance of %d", accountSender.ID, accountSender.Limit, accountSender.Balance)
		return shim.Error("Limit breached. Transfer rejected.")
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
		return shim.Error("Error put world state for sender" + ":" + err.Error())
	}
	accountJSONasBytes, _ = json.Marshal(accountRecipient)
	err = stub.PutState(accountRecipient.ID, accountJSONasBytes)
	if err != nil {
		return shim.Error("Error put world state for receiver" + ":" + err.Error())
	}

	fmt.Printf("========================================")
	fmt.Printf("Transaction of %s to %s", accountSender.ID, accountRecipient.ID)
	fmt.Printf("Sender's balance = %d, Recipient's balance = %d\n", accountSender.Balance, accountRecipient.Balance)
	fmt.Printf("========================================")
	return shim.Success([]byte("Transaction was completed OK !"))
}

func (t *SimpleChaincode) createAccount(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("Creating new Account")
	//_, args := stub.GetFunctionAndParameters()
	var accountID string // Entities
	var balance int      // Asset holdings
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
		return shim.Error("Could not convert account balance to a number" + ":" + err.Error())
	}
	lowerLimit, err = strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("Could not convert account limit to a number" + ":" + err.Error())
	}

	err = createAccount(stub, accountID, owner, balance, lowerLimit)
	if err != nil {
		return shim.Error("Error creating account" + ":" + err.Error())
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
		return shim.Error("Failed to delete state" + ":" + err.Error())
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

	outstring = "\nFollowing accounts are registered:\n=======================================\n"
	res, err := stub.GetStateByRange(keyStart, keyEnd)

	if err != nil {
		return shim.Error("Failed to get world states" + ":" + err.Error())
	}
	kv, err := res.Next()
	if err != nil {
		return shim.Error("Failed to get next world state" + ":" + err.Error())
	}
	account, _ := getAccount(stub, kv.GetKey())
	lomJSONasBytes, err := json.Marshal(account)
	if err != nil {
		return shim.Error("Could not get LOM JSON as object" + ":" + err.Error())
	}
	outstring += " - Account: " + string(lomJSONasBytes[:]) + "\n"
	totalliq := account.Balance

	for res.HasNext() {
		kv, err = res.Next()
		if err != nil {
			return shim.Error("Failed to get next world state" + ":" + err.Error())
		}
		account, _ = getAccount(stub, kv.GetKey())
		lomJSONasBytes, err = json.Marshal(account)
		if err != nil {
			return shim.Error("Could not get LOM JSON as object" + ":" + err.Error())
		}
		outstring += " - Account: " + string(lomJSONasBytes[:]) + "\n"
		totalliq += account.Balance
	}

	res.Close()
	outstring += "\nTotal liquidity supply: " + strconv.Itoa(totalliq) + "\n"
	return shim.Success([]byte(outstring))
}

func (t *SimpleChaincode) listCouchDB(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"LOM\"}}")
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error("Problem when querying world state DB" + ":" + err.Error())
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

func createAccount(stub shim.ChaincodeStubInterface, accountID string, owner string, balance int, limit int) error {

	// ==== Check if LOM already exists ====
	accountAsBytes, err := stub.GetState(accountID)
	if err != nil {
		return err
	} else if accountAsBytes != nil {
		fmt.Println("This account already exists: " + accountID)
		return errors.New("Account already exists:" + accountID)
	}

	// ==== Create Account object and marshal to JSON ====
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
		return shim.Error("Incorrect number of arguments. Expecting ID of the Account to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetState(name) //get the LOM from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Account does not exist: " + name + "\"}"
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
