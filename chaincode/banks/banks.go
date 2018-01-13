/*
 * The sample smart contract for documentation topic:
 * cross border funds transfer
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the bank structure, with 3 properties.  Structure tags are used by encoding/json library
type Bank struct {
	Name     string  `json:"name"`
	BankID   string  `json:"bankID"`
	Country  string  `json:"country"`
	Currency string  `json:"currency"`
	Reserves float64 `json:"reserves"`
}

// Define the customer structure, with 3 properties.  Structure tags are used by encoding/json library
type Customer struct {
	Name           string  `json:"name"`
	CustID         string  `json:"custID"`
	Country        string  `json:"country"`
	Currency       string  `json:"currency"`
	Balance        float64 `json:"balance"`
	CustomerBankID string  `json:"customerBankID"`
}

// Define the map for currency pairs, with 3 properties.  Structure tags are used by encoding/json library
type Forex struct {
	Pair string  `json:"pair"`
	Rate float64 `json:"rate"`
}

/*
 * The Init method is called when the Smart Contract "banks" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract ""
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "queryAll" { //return all the assets on the ledger
		return s.queryAll(APIstub, args)
	} else if function == "query" { //single bank or customer or forexPair
		return s.query(APIstub, args)
	} else if function == "pay" { //execute a payment between two currencies
		return s.pay(APIstub, args)
	} else if function == "createBank" {
		return s.createBank(APIstub, args)
	} else if function == "createCustomer" {
		return s.createCustomer(APIstub, args)
	} else if function == "createForex" {
		return s.createForex(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

/** ----------------------------------------------------------------------**/
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	banks := []Bank{
		{Name: "US_Bank", BankID: "US_Bank", Country: "USA", Currency: "USD", Reserves: 1000000.0},
		{Name: "UK_Bank", BankID: "UK_Bank", Country: "UK ", Currency: "GBP", Reserves: 1000000.0},
		{Name: "Japan_Bank", BankID: "Japan_Bank", Country: "JAPAN", Currency: "JPY", Reserves: 10000000.0},
	}

	customers := []Customer{
		{Name: "US_John_Doe", CustID: "123", Country: "US", Currency: "USD", Balance: 10000.0, CustomerBankID: "US_Bank"},
		{Name: "US_Alice", CustID: "456", Country: "US", Currency: "USD", Balance: 10000.0, CustomerBankID: "US_Bank"},
		{Name: "UK_John_Doe", CustID: "123", Country: "UK", Currency: "GBP", Balance: 10000.0, CustomerBankID: "UK_Bank"},
		{Name: "UK_Alice", CustID: "456", Country: "UK", Currency: "GBP", Balance: 10000.0, CustomerBankID: "UK_Bank"},
		{Name: "JPY_John_Doe", CustID: "123", Country: "Japan", Currency: "JPY", Balance: 1000000.0, CustomerBankID: "Japan_Bank"},
		{Name: "JPY_Alice", CustID: "456", Country: "Japan", Currency: "JPY", Balance: 1000000.0, CustomerBankID: "Japan_Bank"},
	}
	/** currency rates 1:currency format */
	forex := []Forex{
		{Pair: "USD:GBP", Rate: 0.75},
		{Pair: "USD:JPY", Rate: 115.0},
		{Pair: "GBP:USD", Rate: 1.35},
		{Pair: "GBP:JPY", Rate: 155.0},
		{Pair: "JPY:USD", Rate: 0.0088},
		{Pair: "JPY:GBP", Rate: 0.0065},
		/* the following are needed in case payments are in the same currency */
		{Pair: "USD:USD", Rate: 1.0},
		{Pair: "GBP:GBP", Rate: 1.0},
		{Pair: "JPY:JPY", Rate: 1.0},
	}

	writeForexToLedger(APIstub, forex)
	writeCustomerToLedger(APIstub, customers)
	writeBankToLedger(APIstub, banks)

	return shim.Success(nil)
}

/** --------------------------------------------------------------------------------------------------------*/
func writeForexToLedger(APIStub shim.ChaincodeStubInterface, forex []Forex) sc.Response {
	for i := 0; i < len(forex); i++ {
		key := forex[i].Pair
		chkBytes, _ := APIStub.GetState(key)
		if chkBytes == nil { //only add if it is not already present
			asBytes, _ := json.Marshal(forex[i])
			err := APIStub.PutState(forex[i].Pair, asBytes)
			if err != nil {
				return shim.Error(err.Error())
			}
		} else {
			msg := " Forex Pair with key:" + key + " already exists.. skipping ......."
			return shim.Error(msg)
		}
	}
	return shim.Success(nil)
}

/** --------------------------------------------------------------------------------------------------------*/
func writeBankToLedger(APIStub shim.ChaincodeStubInterface, banks []Bank) sc.Response {

	for i := 0; i < len(banks); i++ {
		key := banks[i].BankID
		chkBytes, _ := APIStub.GetState(key)
		if chkBytes == nil { //only add if it is not already present
			asBytes, _ := json.Marshal(banks[i])
			err := APIStub.PutState(key, asBytes)
			if err != nil {
				return shim.Error(err.Error())
			}
		} else {
			msg := " Forex Bank with key:" + key + " already exists.. skipping ......."
			return shim.Error(msg)
		}

	}
	return shim.Success(nil)
}

/** --------------------------------------------------------------------------------------------------------*/
func writeCustomerToLedger(APIStub shim.ChaincodeStubInterface, customers []Customer) sc.Response {
	for i := 0; i < len(customers); i++ {
		key := customers[i].Name + "_" + customers[i].CustID
		chkBytes, _ := APIStub.GetState(key)
		if chkBytes == nil { //only add if it is not already present
			asBytes, _ := json.Marshal(customers[i])
			err := APIStub.PutState(customers[i].Name+"_"+customers[i].CustID, asBytes)
			if err != nil {
				return shim.Error(err.Error())
			}
		} else {
			msg := " Forex Customer with key:" + key + " already exists.. skipping ......."
			return shim.Error(msg)
		}
	}
	return shim.Success(nil)
}

/** --------------------------------------------------------------------------------------------------------*/
func (s *SmartContract) queryAll(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments for querying all assets. Expecting 1")
	}
	//collection := args[0]
	//startKey := collection + "0"
	//endKey := collection + "99"

	resultsIterator, err := APIstub.GetStateByRange("", "")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString("\n,")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}\n")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("\n]")

	fmt.Printf("- queryAll:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

/** ----------------------------------------------------------------------**/

func (s *SmartContract) query(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	asBytes, _ := APIstub.GetState(args[0])
	return shim.Success(asBytes)
}

/** ---------------------------------------------------------------------- */
func (s *SmartContract) pay(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	payAmount, _ := strconv.ParseFloat(args[2], 64)

	//get FROM customer from ledger
	fromCustAsBytes, _ := APIstub.GetState(args[0])
	fromCustomer := Customer{}

	json.Unmarshal(fromCustAsBytes, &fromCustomer)
	fromCustomerName := fromCustomer.Name
	fromCustomerCustID := fromCustomer.CustID
	fromCurrency := fromCustomer.Currency
	fromBalance := float64(fromCustomer.Balance)
	fromBank := fromCustomer.CustomerBankID

	//check if customer has enough balance to cover the payment
	if fromBalance < payAmount {
		errMsg := "Insufficent funds in customer: " + fromCustomerName + " Customer ID: " + fromCustomerCustID
		return shim.Error(errMsg)
	}

	//get TO customer from ledger
	toCustAsBytes, _ := APIstub.GetState(args[1])
	toCustomer := Customer{}

	json.Unmarshal(toCustAsBytes, &toCustomer)
	toCustomerName := toCustomer.Name
	toCurrency := toCustomer.Currency
	toBalance := toCustomer.Balance
	toBank := toCustomer.CustomerBankID

	//get exchange rate from the ledger
	toForexPairAsBytes, _ := APIstub.GetState(fromCurrency + ":" + toCurrency)
	forexPair := Forex{}
	json.Unmarshal(toForexPairAsBytes, &forexPair)
	exchangeRate := forexPair.Rate

	//get bank  FROM ledger
	fromCustBankAsBytes, _ := APIstub.GetState(fromBank)

	fromCustomerBank := Bank{}
	json.Unmarshal(fromCustBankAsBytes, &fromCustomerBank)
	fromBankReserves := fromCustomerBank.Reserves

	//check if bank has reserves to cover the transfer
	if fromBankReserves < payAmount {
		errMsg := "Insufficent funds in bank reserves: " + fromCustomerBank.Name + " Bank ID: " + fromCustomerBank.BankID
		return shim.Error(errMsg)
	}

	//reduce FROM customer balance by payment amount
	fromCustomer.Balance = fromBalance - payAmount
	//reduce FROM bank reservers by payment amount
	fromCustomerBank.Reserves = fromBankReserves - payAmount

	//increase TO customer balance by payment amount
	toCustomer.Balance = toBalance + (payAmount * exchangeRate)
	//get bank  TO ledger
	toCustBankAsBytes, _ := APIstub.GetState(toBank)
	toCustomerBank := Bank{}
	json.Unmarshal(toCustBankAsBytes, &toCustomerBank)
	//increase TO bank reservers by payment amount
	toCustomerBank.Reserves = toCustomerBank.Reserves + (payAmount * exchangeRate)

	//write all changed assets to the ledger
	fromCustAsBytes, _ = json.Marshal(fromCustomer)
	err := APIstub.PutState(args[0], fromCustAsBytes)
	if err != nil {
		return shim.Error("Error writing updates to FROM customer account " + fromCustomer.Name)
	}

	toCustAsBytes, _ = json.Marshal(toCustomer)
	err = APIstub.PutState(args[1], toCustAsBytes)
	if err != nil {
		return shim.Error("Error writing updates to TO customer account " + toCustomer.Name)
	}
	fromCustBankAsBytes, _ = json.Marshal(fromCustomerBank)
	err = APIstub.PutState(fromBank, fromCustBankAsBytes)
	if err != nil {
		return shim.Error("Error writing updates to FROM Bank account " + fromCustomerBank.Name)
	}
	toCustBankAsBytes, _ = json.Marshal(toCustomerBank)
	err = APIstub.PutState(toBank, toCustBankAsBytes)
	if err != nil {
		return shim.Error("Error writing updates to TO Bank account " + toCustomerBank.Name)
	}

	if err == nil {
		fmt.Println("~~~~~~~~~~~~~~~~~ Success fully transferred: ", payAmount, " From: ", fromCustomerName, " TO: ", toCustomerName, "~~~~~~~~~~~~~~~~~")
	}

	return shim.Success(nil)
}

/** ----------------------------------------------------------------------------------------------
ceate bank needs 5 args
Name     string  `json:"name"`
BankID	 string  `json:"bankID"`
Country  string  `json:"country"`
Currency string  `json:"currency"`
Reserves float64 `json:"reserves"`
args: ['EU_Bank', 'EU_Bank', 'Europe', 'EURO', '1000000.0'],
*/
func (s *SmartContract) createBank(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments for creating a bank. Expecting 5")
	}
	reserves, _ := strconv.ParseFloat(args[4], 64)
	banks := []Bank{Bank{Name: args[0], BankID: args[1], Country: args[2], Currency: args[3], Reserves: reserves}}

	writeBankToLedger(APIstub, banks)
	return shim.Success(nil)
}

/**----------------------------------------------------------------------------------------------
createCustomer needs 6 args
	Name     string  `json:"name"`
	CustID   string  `json:"custID"`
	Country  string  `json:"country"`
	Currency string  `json:"currency"`
	Balance  float64 `json:"balance"`
	CustomerBankID string `json:"customerBankID"`
	["US_Mary_Jane", "789",  "US", "USD", 100000.0, "US_Bank"],
*/
func (s *SmartContract) createCustomer(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments for creating a customer. Expecting 6")
	}
	balance, _ := strconv.ParseFloat(args[4], 64)
	customers := []Customer{{Name: args[0], CustID: args[1], Country: args[2], Currency: args[3], Balance: balance, CustomerBankID: args[5]}}

	writeCustomerToLedger(APIstub, customers)

	return shim.Success(nil)
}

/**----------------------------------------------------------------------------------------------
createForex needs 2 args
	Pair string  `json:"pair"`
	Rate float64 `json:"rate"`
*/
func (s *SmartContract) createForex(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments for creating a forex pair. Expecting 2")
	}

	rate, _ := strconv.ParseFloat(args[1], 64)
	forex := []Forex{{Pair: args[0], Rate: rate}}

	writeForexToLedger(APIstub, forex)

	return shim.Success(nil)
}

/**----------------------------------------------------------------------------------------------
The main function is only relevant in unit test mode. Only included here for completeness.
*/

func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}

	fmt.Println("successfully initialized smart contract")

}
