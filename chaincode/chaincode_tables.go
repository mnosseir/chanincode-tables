/*
Copyright IBM Corp 2016 All Rights Reserved.

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

import (
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	fmt.Println("Init", args)

	err := stub.CreateTable(args[0], []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "EMP_ID", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "EMP_LNAME", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "EMP_FNAME", Type: shim.ColumnDefinition_STRING, Key: false},
	})

	// err := stub.PutState("hello_world", []byte(args[0]))
    if err != nil {
        return nil, err
    }

	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    fmt.Println("invoke is running " + function)

    // Handle different functions
    if function == "init" {
        return t.Init(stub, "init", args)
    } else if function == "write" {
        return t.write(stub, args)
    } else if function == "new_emp" {
			err := t.create_new_emp(stub, args)
			return nil, err
		}
    fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

func (t *SimpleChaincode) get_emp_by_id(stub shim.ChaincodeStubInterface,  emp_id string) ([]byte, error) {
		var columns []shim.Column
		col1 := shim.Column{Value: &shim.Column_String_{String_: emp_id}}
		columns = append(columns, col1)
		fmt.Println("Querying for row with with:", columns)
		row, err := stub.GetRow("EMP", columns)
		if err != nil {
			return nil, err
		}
		fmt.Println("row found:", row )
		var s string = ""
		for i :=0 ; i<len(row.Columns) ; i++ {
			s += row.Columns[i].GetString_()
		}
		emp_data := []byte(s)

		if err != nil {
			return nil, err
		}
		return emp_data, nil
}

func (t *SimpleChaincode) create_new_emp(stub shim.ChaincodeStubInterface, args []string) (error) {
		fmt.Println("create_new_emp", args)

		if len(args) != 3 {
				return errors.New("Incorrect number of arguments. Expecting 3")
		}
		row :=  shim.Row{
				Columns: []*shim.Column{
					&shim.Column{Value: &shim.Column_String_{String_: args[0]}},
					&shim.Column{Value: &shim.Column_String_{String_: args[1]}},
					&shim.Column{Value: &shim.Column_String_{String_: args[2]}},
				},
		}
		fmt.Println("inserting row:", row)

		ok, err := stub.InsertRow("EMP", row)
		if !ok || err != nil {
			if err != nil {
				return err
			} else {
					return errors.New("Error inserting a new row")
			}
		}
		return nil
}

func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    var key, value string
    var err error
    fmt.Println("running write()")

    if len(args) != 2 {
        return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
    }

    key = args[0]                            //rename for fun
    value = args[1]
    err = stub.PutState(key, []byte(value))  //write the variable into the chaincode state
    if err != nil {
        return nil, err
    }
    return nil, nil
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    fmt.Println("query is running " + function)

    // Handle different functions
		if function == "get_emp_by_id" {
			return t.get_emp_by_id(stub, args[0])
		}
		fmt.Println("query did not find func: " + function)

    return nil, errors.New("Received unknown function query: " + function)
}
