package Shequ

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/common/flogging"
)

type ShequContract struct {
}

// 全局日至变量
var logger = flogging.MustGetLogger("shequ_cc")

type Shequ struct {
	Token    string `json:"token"`    //对应的token存储值
	HashCode string `json:"hashcode"` //对应的hash值
}

func (a *ShequContract) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (s *ShequContract) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fun, args := stub.GetFunctionAndParameters()
	logger.Infof("Function name is : %s", fun)
	logger.Infof("Args length is : %d", len(args))
	for i := 0; i < len(args); i++ {
		logger.Infof("Args %d is : %s", i+1, args[i])
	}
	switch fun {
	case "InitLedger":
		return s.InitLedger(stub)
	case "AddNewSheQu":
		return s.AddNewSheQu(stub, args)
	case "GetSQByToken":
		return s.GetSQByToken(stub, args)
	case "UpdateSQByToken":
		return s.UpdateSQByToken(stub, args)
	case "QueryAllSQ":
		return s.QueryAllSQ(stub, args)
	}
	return shim.Error("Invalid Smart Contract function name!")
}

func (s *ShequContract) InitLedger(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Infof("Begin to Add 5 infos into the ledger")
	shequs := []Shequ{
		{Token: "token1", HashCode: "hashcode1"},
		{Token: "token2", HashCode: "hashcode2"},
		{Token: "token3", HashCode: "hashcode3"},
		{Token: "token4", HashCode: "hashcode4"},
		{Token: "token5", HashCode: "hashcode5"},
	}
	i := 0
	for i < len(shequs) {
		shequAPIBytes, _ := json.Marshal(shequs[i])
		err := stub.PutState(shequs[i].Token, shequAPIBytes)
		if err != nil {
			logger.Error(err.Error())
			return shim.Error(err.Error())
		}
		logger.Infof("%s putstate success.", shequs[i].Token)
		i = i + 1
	}
	logger.Infof("5 infos has been added into the ledger!")
	return shim.Success(nil)
}

// 添加新的社区信息
func (s *ShequContract) AddNewSheQu(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var newShequ Shequ
	// 检查参数个数是否正确
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments.")
	}
	newShequ.Token = args[0]
	newShequ.HashCode = args[1]

	newSheQuJSONAsBytes, err := json.Marshal(newShequ)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(newShequ.Token, newSheQuJSONAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(newSheQuJSONAsBytes)
}

// 根据token获取信息
func (*ShequContract) GetSQByToken(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		logger.Infof("Incorrect number of arguments.")
		return shim.Error("Incorrect number of arguments.")
	}

	token := args[0]
	//logger.Infow(token)
	data, err := stub.GetState(token)
	//logger.Infow(strconv.Itoa(len(data)))
	if err != nil {
		logger.Errorf("Data get failed : %s", err.Error())
		return shim.Error("data get failed: " + err.Error())
	}
	logger.Infof(string(data))
	return shim.Success(data)
}

// 根据token更新社区信息
func (*ShequContract) UpdateSQByToken(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 检查参数数量是不是正确
	if len(args) != 2 {
		logger.Error("参数数量错误.")
		return shim.Error("Incorrect number of arguments.")
	}
	// 获取区块链上的社区信息
	token := args[0]
	logger.Infof("将要查询的sq Token为:" + token)
	SheQuAsBytes, _ := stub.GetState(token)
	newShequ := Shequ{}

	// 将获取到的信息写入到临时变量shequToUpdate中
	_ = json.Unmarshal(SheQuAsBytes, &newShequ)
	newShequ.HashCode = args[1]

	// 修改完毕的信息重新上链
	SheQuAsBytes, _ = json.Marshal(newShequ)
	stub.PutState(newShequ.Token, SheQuAsBytes)

	// 最后要将json格式的改后信息发送回去
	return shim.Success(SheQuAsBytes)
}

// 查询所有的社区信息
// 由于goapi 只能通过GetStateByRange进行读取信息， 所以这里也只能通过该方法进行读取判断
func (*ShequContract) QueryAllSQ(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 检查参数数量是否正确
	if len(args) != 2 {
		logger.Error("参数数量错误.")
	}
	// 设置起始与终止变量 这里需要注意的是 终止的数字应该设为9以上才有效
	//startkey := args[0]
	//endKey := args[1]
	startkey := "token0"
	endKey := "token91"
	logger.Infof(startkey)
	logger.Infof(endKey)
	resultIterator, error := stub.GetStateByRange(startkey, endKey)
	if error != nil {
		logger.Error(error.Error())
	}
	defer resultIterator.Close()

	// 设置一个json数组包含查询结果
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultIterator.HasNext() {
		queryResponse, err := resultIterator.Next()
		if err != nil {
			logger.Error(err.Error())
			return shim.Error(err.Error())
		}
		// 在数组成员之前添加一个逗号，对第一个数组成员取消显示逗号
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"key\":")
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
	fmt.Printf("- query All Shqequ Date: \n %s\n", buffer.String())
	return shim.Success(buffer.Bytes())
}

func main() {
	err := shim.Start(new(ShequContract))
	if err != nil {
		fmt.Printf("Error starting the contract: %s", err)
	}
}
