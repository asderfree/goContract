package Shequ

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/common/flogging"
)

type ShequContract struct{
	
}

// 全局日至变量
var logger = flogging.MustGetLogger("shequ_cc")

type Shequ struct {
	Token		string 		`json:"token"`		//对应的token存储值
	HashCode	string		`json:"hashcode"`	//对应的hash值
}

func (a *ShequContract) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (s *ShequContract) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fun, args := stub.GetFunctionAndParameters()
	logger.Infof("Function name is : %s", fun)
	logger.Infof("Args length is : %d", len(args))
	for i := 0 ;i < len(args); i++{
		logger.Infof("Args %d is : %s", i+1, args[i])
	}
	switch  fun {
	case "InitLedger":
		return s.InitLedger(stub)
	case "AddNewSheQu":
		return s.AddNewSheQu(stub,args)
	case "GetSQByToken":
		return s.GetSQByToken(stub,args)
	}
	return shim.Error("Invalid Smart Contract function name!")
}

func (s *ShequContract) InitLedger(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Infof("Begin to Add 5 infos into the ledger")
	shequs := []Shequ{
		{Token: "token1",HashCode: "hashcode1"},
		{Token: "token2",HashCode: "hashcode2"},
		{Token: "token3",HashCode: "hashcode3"},
		{Token: "token4",HashCode: "hashcode4"},
		{Token: "token5",HashCode: "hashcode5"},
	}
	i := 0
	for i < len(shequs){
		shequAPIBytes, _ := json.Marshal(shequs[i])
		err:=stub.PutState(shequs[i].Token,shequAPIBytes)
		if err != nil {
			logger.Error(err.Error())
			return shim.Error(err.Error())
		}
		logger.Infof("%s putstate success.",shequs[i].Token)
		i = i + 1
	}
	logger.Infof("5 infos has been added into the ledger!")
	return shim.Success(nil)
}

// 添加新的社区信息
func (s *ShequContract) AddNewSheQu(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	var err error
	var newShequ Shequ
	// 检查参数个数是否正确
	if len(args) != 2{
		return shim.Error("Incorrect number of arguments.")
	}
	newShequ.Token = args[0]
	newShequ.HashCode = args[1]

	newSheQuJSONAsBytes,err := json.Marshal(newShequ)
	if err!= nil{
		return shim.Error(err.Error())
	}

	err = stub.PutState(newShequ.Token,newSheQuJSONAsBytes)
	if err != nil{
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

// 根据token获取信息
func (*ShequContract) GetSQByToken(stub shim.ChaincodeStubInterface,args []string)pb.Response {
	if len(args) != 1{
		logger.Infof("Incorrect number of arguments.")
		return shim.Error("Incorrect number of arguments.")
	}

	token := args[0]
	//logger.Infow(token)
	data,err := stub.GetState(token)
	//logger.Infow(strconv.Itoa(len(data)))
	if err != nil {
		logger.Errorf("Data get failed : %s",err.Error())
		return shim.Error("data get failed: " + err.Error())
	}
	logger.Infof(string(data))
	return shim.Success(data)
}

func main(){
	err := shim.Start(new(ShequContract))
	if err != nil {
		fmt.Printf("Error starting the contract: %s",err)
	}
}


