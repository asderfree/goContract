package store

import (
	"bytes"
	"encoding/json"
	"fmt"
	sc "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric/common/flogging"
	"strconv"
)

type SmartContract struct {

}

type Car struct {
	Make 	string 	`json:"make"`
	Model 	string	`json:"model"`
	Color	string 	`json:"color"`
	Owner	string	`json:"owner"`
}

type carPrivateDetails struct{
	Owner	string 	`json:"owner"`
	Price	string	`json:"price"`
}

// Init 初始化合约
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

var logger = flogging.MustGetLogger("fabcar_cc")

// Invoke 提交合约的方法
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response  {
	function, args := APIstub.GetFunctionAndParameters()
	logger.Infof("Function name is : %s", function)
	logger.Infof("Args length is : %d", len(args))

	switch function{
	case "queryCar":
		return s.queryCar(APIstub, args)
	case "initLedger":
		return s.initLedger(APIstub)
	case "createCar":
		return s.createCar(APIstub, args)
	case "queryAllCars":
		return s.queryAllCars(APIstub)
	default:
		return shim.Error("Invalid Smart Contract function name.")
	}
}


// 查询车辆信息
func (s *SmartContract) queryCar(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1{
		return shim.Error("Incorrect number of arguments. Excepting 1")
	}

	carAsBytes, _ := APIstub.GetState(args[0])
	logger.Infof(string(carAsBytes))
	return shim.Success(carAsBytes)
}


func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	cars := []Car{
		Car{Make: "Toyota", Model: "Prius", Color: "blue", Owner: "Tomoko"},
		Car{Make: "Ford", Model: "Mustang", Color: "red", Owner: "Brad"},
		Car{Make: "Hyundai", Model: "Tucson", Color: "green", Owner: "Jin Soo"},
		Car{Make: "Volkswagen", Model: "Passat", Color: "yellow", Owner: "Max"},
		Car{Make: "Tesla", Model: "S", Color: "black", Owner: "Adriana"},
		Car{Make: "Peugeot", Model: "205", Color: "purple", Owner: "Michel"},
		Car{Make: "Chery", Model: "S22L", Color: "white", Owner: "Aarav"},
		Car{Make: "Fiat", Model: "Punto", Color: "violet", Owner: "Pari"},
		Car{Make: "Tata", Model: "Nano", Color: "indigo", Owner: "Valeria"},
		Car{Make: "Holden", Model: "Barina", Color: "brown", Owner: "Shotaro"},
	}

	i := 0
	for i < len(cars) {
		carAsBytes, _ := json.Marshal(cars[i])
		APIstub.PutState("CAR"+strconv.Itoa(i), carAsBytes)
		i = i + 1
	}

	return shim.Success(nil)
}


func (s *SmartContract) createCar(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 5 {
		logger.Error("Incorrect number of arguments. Expecting 5, But get " , len(args))
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var car = Car{Make: args[1], Model: args[2], Color: args[3], Owner: args[4]}

	carAsBytes, _ := json.Marshal(car)
	APIstub.PutState(args[0], carAsBytes)

	indexName := "owner~key"
	colorNameIndexKey, err := APIstub.CreateCompositeKey(indexName, []string{car.Owner, args[0]})
	if err != nil {
		logger.Error(err.Error())
		return shim.Error(err.Error())
	}
	value := []byte{0x00}
	APIstub.PutState(colorNameIndexKey, value)
	logger.Infof("get Here")
	return shim.Success(carAsBytes)
}

func (s *SmartContract) queryAllCars(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "CAR0"
	endKey := "CAR999"
	logger.Infof(startKey)
	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
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
			buffer.WriteString(",")
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
	buffer.WriteString("]")

	fmt.Printf("- queryAllCars:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}