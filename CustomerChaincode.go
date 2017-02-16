Skip to content
This repository
Search
Pull requests
Issues
Gist
 @vinodkushwaha
 Sign out
 Watch 0
  Star 0
  Fork 0 vinodkushwaha/CustomerChainCodeExample
 Code  Issues 0  Pull requests 0  Projects 0  Wiki  Pulse  Graphs  Settings
Branch: master Find file Copy pathCustomerChainCodeExample/CustomerChaincode.go
189eb5b  just now
@vinodkushwaha vinodkushwaha Update CustomerChaincode.go
2 contributors @vinodkushwaha @gulalkumar
RawBlameHistory     
196 lines (155 sloc)  5.33 KB
package main

import (
	"errors"
	"fmt"
	//"strconv"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	//"github.com/golang/protobuf/ptypes/timestamp"
)

// Customer Chaincode implementation
type CustomerChaincode struct {
}

var customerIndexTxStr = "_customerIndexTxStr"

type CustomerDoc struct {
    DOCUMENT_NAME string `json:"DOCUMENT_NAME"`
	DOCUMENT_STRING string `json:"DOCUMENT_STRING"`
}

type CustomerData struct{
	CUSTOMER_ID string `json:"CUSTOMER_ID"`
	CUSTOMER_NAME string `json:"CUSTOMER_NAME"`
	CUSTOMER_DOB string `json:"CUSTOMER_DOB"`
	CUSTOMER_KYC_FLAG string `json:"CUSTOMER_KYC_FLAG"`
	CUSTOMER_DOC []CustomerDoc
	}


func (t *CustomerChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var err error
	// Initialize the chaincode

	fmt.Printf("Deployment of Customer ChainCode is completed\n")

	var emptyCustomerTxs []CustomerData
	jsonAsBytes, _ := json.Marshal(emptyCustomerTxs)
	err = stub.PutState(customerIndexTxStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}


	return nil, nil
}

// Add customer data for the policy
func (t *CustomerChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == customerIndexTxStr {
		return t.RegisterCustomer(stub, args)
	}
	return nil, nil
}

func (t *CustomerChaincode)  RegisterCustomer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var CustomerDataObj CustomerData
	var CustomerDataList []CustomerData
	var err error
   	fmt.Printf("********pankaj CUSTOMER_DOC:%d\n", len(args))
	
	if len(args) < 4 {
		return nil, errors.New("Incorrect number of arguments. Need 4 arguments")
	}

	// Initialize the chaincode
	CustomerDataObj.CUSTOMER_ID = args[0]
	CustomerDataObj.CUSTOMER_NAME = args[1]
	CustomerDataObj.CUSTOMER_DOB = args[2]
	CustomerDataObj.CUSTOMER_KYC_FLAG = args[3]
	fmt.Printf("********pankaj CUSTOMER_DOC:%s\n", args[4])
	
	var number_of_docs int
	number_of_docs = (len(args)-4)/2
	
	var CustomerDocObj []CustomerDoc
	
	for i := 0; i < number_of_docs; i++ {
		
		CustomerDocObj[i].DOCUMENT_NAME = args[4]
		fmt.Printf("********pankaj CustomerDocObj[i].DOCUMENT_NAME:%s\n", i)
		CustomerDocObj[i].DOCUMENT_STRING = args[5]
	}
	
	CustomerDataObj.CUSTOMER_DOC = CustomerDocObj
	
	customerTxsAsBytes, err := stub.GetState(customerIndexTxStr)
	if err != nil {
		return nil, errors.New("Failed to get customer transactions")
	}
	json.Unmarshal(customerTxsAsBytes, &CustomerDataList)

	CustomerDataList = append(CustomerDataList, CustomerDataObj)
	jsonAsBytes, _ := json.Marshal(CustomerDataList)

	err = stub.PutState(customerIndexTxStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// Query callback representing the query of a chaincode
func (t *CustomerChaincode) Query(stub shim.ChaincodeStubInterface,function string, args []string) ([]byte, error) {

	var customer_name string // Entities
	var customer_id string
	var customer_dob string
	var err error
	var resAsBytes []byte

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3 parameters to query")
	}

	customer_id = args[0]
	customer_name = args[1]
	customer_dob = args[2]

	resAsBytes, err = t.GetCustomerDetails(stub, customer_name, customer_id, customer_dob)

	fmt.Printf("Query Response:%s\n", resAsBytes)

	if err != nil {
		return nil, err
	}

	return resAsBytes, nil
}

func (t *CustomerChaincode)  GetCustomerDetails(stub shim.ChaincodeStubInterface, customer_name string, customer_id string, customer_dob string) ([]byte, error) {

	//var requiredObj CustomerData
	var objFound bool
	CustomerTxsAsBytes, err := stub.GetState(customerIndexTxStr)
	if err != nil {
		return nil, errors.New("Failed to get Customer Records")
	}
	var CustomerTxObjects []CustomerData
	var CustomerTxObjects1 []CustomerData
	json.Unmarshal(CustomerTxsAsBytes, &CustomerTxObjects)
	length := len(CustomerTxObjects)
	fmt.Printf("Output from chaincode: %s\n", CustomerTxsAsBytes)

	if customer_id == "" {
		res, err := json.Marshal(CustomerTxObjects)
		if err != nil {
		return nil, errors.New("Failed to Marshal the required Obj")
		}
		return res, nil
	}

	objFound = false
	// iterate
	for i := 0; i < length; i++ {
		obj := CustomerTxObjects[i]
		//if ((customer_id == obj.CUSTOMER_ID) && (customer_name == obj.CUSTOMER_NAME) && (customer_dob == obj.CUSTOMER_DOB)) 
		
		fmt.Printf("Output from customer_id: %s\n", customer_id)
		fmt.Printf("Output from obj.CUSTOMER_ID: %s\n", obj.CUSTOMER_ID)
		fmt.Printf("Output from customer_name: %s\n", customer_name)
		fmt.Printf("Output from obj.CUSTOMER_NAME: %s\n", obj.CUSTOMER_NAME)
		fmt.Printf("Output from customer_dob: %s\n", customer_dob)
		fmt.Printf("Output from obj.CUSTOMER_DOB: %s\n", obj.CUSTOMER_DOB)
		if ((obj.CUSTOMER_ID) == customer_id){
			CustomerTxObjects1 = append(CustomerTxObjects1,obj)
			//requiredObj = obj
			objFound = true
		}
	}

	if objFound {
		res, err := json.Marshal(CustomerTxObjects1)
		if err != nil {
		return nil, errors.New("Failed to Marshal the required Obj")
		}
		return res, nil
	} else {
		res, err := json.Marshal("No Data found")
		if err != nil {
		return nil, errors.New("Failed to Marshal the required Obj")
		}
		return res, nil
	}
}

func main() {
	err := shim.Start(new(CustomerChaincode))
	if err != nil {
		fmt.Printf("Error starting Customer Simple chaincode: %s", err)
	}
}
Contact GitHub API Training Shop Blog About
© 2017 GitHub, Inc. Terms Privacy Security Status Help
