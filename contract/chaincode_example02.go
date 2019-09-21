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
	"encoding/binary"
	"encoding/json"
	"fmt"
	"encoding/base64"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
	urldata map[string]DataStruct
}

type DataStruct struct{
	Urlhash string
	Counter uint32
	Keypool []string
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Init")
	t.urldata = make(map[string]DataStruct)
	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "downloadquery" {
		// Make payment of X units from A to B
		return t.downloadquery(stub, args)
	} else if function == "delete" {
		// Deletes an entity from its state
		return t.delete(stub, args)
	} else if function == "upload" {
		// the old "Query" is now implemtned in invoke
		return t.upload(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"delete\" \"query\"")
}

// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) downloadquery(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var url string   // Entities
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	url = args[0]

	// Get the state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	Avalbytes, err := stub.GetState(url)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Avalbytes == nil {
		return shim.Error("Entity not found")
	}
	counter := binary.LittleEndian.Uint32(Avalbytes)

	retdata := t.urldata[url].Keypool[counter]

	counter += 1

	a := make([]byte, 4)
	binary.LittleEndian.PutUint32(a, counter)
	// Write the state back to the ledger
	err = stub.PutState(url,a)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(retdata))
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

// query callback representing the query of a chaincode
func (t *SimpleChaincode) upload(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	// Initialize the chaincode
	url := args[0]

	data := args[1]

	fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	databyte, _ := base64.StdEncoding.DecodeString(data)
	counter := uint32(0)
	var ss []string
	json.Unmarshal(databyte,&ss)

	t.urldata[url] = DataStruct{url,counter,ss}
	a := make([]byte, 4)
	binary.LittleEndian.PutUint32(a, counter)
	// Write the state to the ledger
	err = stub.PutState(url, a)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)


}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
