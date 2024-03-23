package main

import (
	"errors"
	"fmt"
	"strconv"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// ABstore Chaincode implementation
type ABstore struct {
	contractapi.Contract
}

func (t *ABstore) Init(ctx contractapi.TransactionContextInterface, sell string, sellVal int, buy string, buyVal int, company string, companyVal int) error {
	fmt.Println("ABstore Init")
	var err error
	// Initialize the chaincode
	fmt.Printf("sellVal = %d, buyVal = %d, companyVal = %d\n", sellVal, buyVal, companyVal)
	// Write the state to the ledger
	err = ctx.GetStub().PutState(sell, []byte(strconv.Itoa(sellVal)))
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(buy, []byte(strconv.Itoa(buyVal)))
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(company, []byte(strconv.Itoa(companyVal)))
	if err != nil {
		return err
	}

	return nil
}

func (t *ABstore) Invoke(ctx contractapi.TransactionContextInterface, sell, buy, company, itemName string, X int) error {
	var err error
	var sellVal int
	var buyVal int
	var companyVal int
	// Get the state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	sellValBytes, err := ctx.GetStub().GetState(sell)
	if err != nil {
		return fmt.Errorf("Failed to get state")
	}
	if sellValBytes == nil {
		return fmt.Errorf("Entity not found")
	}
	sellVal, _ = strconv.Atoi(string(sellValBytes))

	buyValBytes, err := ctx.GetStub().GetState(buy)
	if err != nil {
		return fmt.Errorf("Failed to get state")
	}
	if buyValBytes == nil {
		return fmt.Errorf("Entity not found")
	}
	buyVal, _ = strconv.Atoi(string(buyValBytes))

	companyValBytes, err := ctx.GetStub().GetState(company)
	if err != nil {
		return fmt.Errorf("Failed to get state")
	}
	if companyValBytes == nil {
		return fmt.Errorf("Entity not found")
	}
	companyVal, _ = strconv.Atoi(string(companyValBytes))
	
	// Calculate fee to be paid to the company
	fee := int(float64(X) * 0.03)
	
	// Perform the execution
	sellVal += X
	buyVal = buyVal-X-fee
	companyVal += fee
	
	fmt.Printf("sellVal = %d, buyVal = %d, companyVal = %d\n", sellVal, buyVal, companyVal)

	// Write the state back to the ledger
	err = ctx.GetStub().PutState(sell+"-"+itemName, []byte(strconv.Itoa(sellVal)))
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(buy+"-"+itemName, []byte(strconv.Itoa(buyVal)))
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(company+"-"+itemName, []byte(strconv.Itoa(companyVal)))
	if err != nil {
		return err
	}

	return nil
}


// Delete an entity from state
func (t *ABstore) Delete(ctx contractapi.TransactionContextInterface, sell string) error {
	// Delete the key from the state in ledger
	err := ctx.GetStub().DelState(sell)
	if err != nil {
		return fmt.Errorf("Failed to delete state")
	}

	return nil
}

// Query callback representing the query of a chaincode
func (t *ABstore) Query(ctx contractapi.TransactionContextInterface, sell string) (string, error) {
	var err error
	// Get the state from the ledger
	sellValBytes, err := ctx.GetStub().GetState(sell)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + sell + "\"}"
		return "", errors.New(jsonResp)
	}

	if sellValBytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + sell + "\"}"
		return "", errors.New(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + sell + "\",\"Amount\":\"" + string(sellValBytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return string(sellValBytes), nil
}

func (t *ABstore) GetAllQuery(ctx contractapi.TransactionContextInterface) ([]string, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	var wallet []string
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		jsonResp := "{\"Name\":\"" + string(queryResponse.Key) + "\",\"Amount\":\"" + string(queryResponse.Value) + "\"}"
		wallet = append(wallet, jsonResp)
	}
	return wallet, nil
}

func main() {
	cc, err := contractapi.NewChaincode(new(ABstore))
	if err != nil {
		panic(err.Error())
	}
	if err := cc.Start(); err != nil {
		fmt.Printf("Error starting ABstore chaincode: %s", err)
	}
}
