package Info

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	sc "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/common/flogging"
	"strconv"
)

// SmartContract Define the Smart Contract structure
type SmartContract struct {

}

// Info 上链数据的结构体
type Info struct{
	Token 		string   	`json:"token"`		//被上链数据在状态数据库中的token
	HashCode	string 		`josn:"hashcode"`	//被上链数据的hash指纹
}

// Init 初始化链码的方法
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response{
	return shim.Success(nil)
}

// 全局日至变量
var logger = flogging.MustGetLogger("info_cc")

// Invoke 提交合约以及合约执行的主方法 所有的合约函数的执行都需要通过该方法
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	function, args := APIstub.GetFunctionAndParameters()
	logger.Infof("Function name is : %d", function)
	logger.Infof("Args length is : %d", len(args))
	for i := 0 ;i < len(args); i++{
		logger.Infof("Args %d is : %s", i+1, args[i])
	}
	return shim.Error("Invalid Smart Contract function name!")
}

// queryInfo : 查询上链信息
func (s *SmartContract) queryInfo(APIstub shim.ChaincodeStubInterface, args []string) sc.Response{
	if len(args) != 1{
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	infoAsBytes , _ := APIstub.GetState(args[0])
	return shim.Success(infoAsBytes)
}

// 初始化合约 用于测试使用
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response{
	infos := []Info{
		Info{Token: "token1", HashCode: "Hash1"},
		{Token:"token2",HashCode: "Hash2"},
		{Token:"token3",HashCode: "Hash3"},
		{Token:"token4",HashCode: "Hash4"},
		{Token:"token5",HashCode: "Hash5"},
		{Token:"token6",HashCode: "Hash6"},
		{Token:"token7",HashCode: "Hash7"},
		{Token:"token8",HashCode: "Hash8"},
	}
	i := 0
	for i <len(infos){
		infoAsBytes, _ := json.Marshal(infos[i])
		APIstub.PutState("INFO"+strconv.Itoa(i),infoAsBytes)
		i = i+1
	}
	return shim.Success(nil)
}

func (s *SmartContract) createInfo(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) !=3{
		return shim.Error("Incorrect number of arguments , expecting 3")
	}

	var info = Info{Token: args[1], HashCode: args[2]}
	infoAsByte, _ := json.Marshal(info)
	APIstub.PutState(args[0],infoAsByte)
	return shim.Success(infoAsByte)
}

// 根据token查找hash
func (s *SmartContract) queryInfoByToken(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args)!= 1{
		return shim.Error("Incorrect number of arguments")
	}
	token := args[0]

	tokenAndResultIterator, err := APIstub.GetStateByPartialCompositeKey("token", []string{token})
	if err != nil{
		return shim.Error(err.Error())
	}

	defer tokenAndResultIterator.Close()

	var i int
	var id string

	var infos []byte
	bArrayMemberAlreadyWeitten := false

	for i = 0; tokenAndResultIterator.HasNext(); i++{
		responseRange, err := tokenAndResultIterator.Next()
		if err != nil{
			return shim.Error(err.Error())
		}

		objectType, compositeKeyParts ,err := APIstub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return shim.Error(err.Error())
		}

		id = compositeKeyParts[1]
		assetAsBytes ,err := APIstub.GetState(id)

		if bArrayMemberAlreadyWeitten == true{
			newBytes := append([]byte(","), assetAsBytes...)
			infos = append(infos,newBytes...)
		} else{
			infos = append(infos, assetAsBytes...)
		}

		fmt.Printf("Found a asset for index : %s asset id : ", objectType, compositeKeyParts[0], compositeKeyParts[1])
		bArrayMemberAlreadyWeitten = true
	}

	infos = append(infos, []byte("]")...)

	return shim.Success(infos)
}








