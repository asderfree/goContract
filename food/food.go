package food

import (
	"encoding/json"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/common/flogging"
	"strconv"
)

type ChainCode struct {

}

//  食品数据结构体
type Infos struct {
	FoodID 			string		`json:"FoodID"`			// 食品id
	FoodProInfo		ProInfo 	`json:"FoodProInfo"`	// 生产信息
	FoodIngInfo		[]IngInfo	`json:"FoodIngInfo"`	// 配料信息
	FoodLogInfo		LogInfo 	`json:"FoodLogInfo"`	// 物流信息
}

type AllInfo struct{
	FoodID 			string 		`json:"FoodId"`
	FoodProInfo 	ProInfo 	`json:"FoodProInfo"`
	FoodIngInfo 	[]IngInfo 	`json:"FoodIngInfo"`
	FoodLogInfo 	[]LogInfo 	`json:"FoodLogInfo"`
}

// 生产信息
type ProInfo struct {
	FoodName		string		`json:"FoodName"`		// 食品名称
	FoodSpec		string		`json:"FoodSpec"`		// 食品规格
	FoodMFGDate		string		`json:"FoodMFGData"`	// 食品生产日期
	FoodEXPDate 	string 		`json:"FoodEXPDate"`    // 食品保质期
	FoodLOT 		string 		`json:"FoodLOT"`   		// 食品批次号
	FoodQSID 		string 		`json:"FoodQSID"`  		// 食品生产许可证编号
	FoodMFRSName 	string 		`json:"FoodMFRSName"`  	// 食品生产商名称
	FoodProPrice 	string 		`json:"FoodProPrice"` 	// 食品生产价格
	FoodProPlace 	string 		`json:"FoodProPlace"`   // 食品生产所在地
}

//  配料信息
type IngInfo struct {
	IngID 			string		`json:"IngID"`			// 配料ID
	IngName 		string 		`json:"IngName"`		// 配料名称
}

// 日志信息
type LogInfo struct{
	LogDepartureTm 	string 		`json:"LogDepartureTm"`	// 出发时间
	LogArrivalTm 	string 		`json:"LogArrivalTm"`  	// 到达时间
	LogMission 		string 		`json:"LogMission"`   	// 处理业务(储存or运输)
	LogDeparturePl 	string 		`json:"LogDeparturePl"` // 出发地
	LogDest 		string 		`json:"LogDest"`     	// 目的地
	LogToSeller 	string 		`json:"LogToSeller"`   	// 销售商
	LogStorageTm 	string 		`json:"LogStorageTm"`   // 存储时间
	LogMOT 			string 		`json:"LogMOT"`         // 运送方式
	LogCopName 		string 		`json:"LogCopName"` 	// 物流公司名称
	LogCost 		string 		`json:"LogCost"`   		// 费用
}

func (c *ChainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

var logger = flogging.MustGetLogger("food_cc")

func (c *ChainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	switch function {
	case "addProInfo":
		return c.addProInfo(stub, args)
	default:
		return shim.Error("Recevied unkown function invocation")
	}
}

// 新增食品信息
func (a *ChainCode) addProInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 10{
		msg := "Incorrect number of arguments. Expecting 10, But " + strconv.Itoa(len(args)) + " get!"
		logger.Error(msg)
		return shim.Error(msg)
	}

	var FoodInfos Infos
	FoodInfos.FoodID = args[0]
	if FoodInfos.FoodID == ""{
		msg := "FoodId can not be empty, Please Reset it!"
		logger.Error(msg)
		return shim.Error(msg)
	}

	FoodInfos.FoodProInfo.FoodName = args[1]
	FoodInfos.FoodProInfo.FoodSpec = args[2]
	FoodInfos.FoodProInfo.FoodMFGDate = args[3]
	FoodInfos.FoodProInfo.FoodEXPDate = args[4]
	FoodInfos.FoodProInfo.FoodLOT = args[5]
	FoodInfos.FoodProInfo.FoodQSID = args[6]
	FoodInfos.FoodProInfo.FoodMFGDate = args[7]
	FoodInfos.FoodProInfo.FoodProPrice = args[8]
	FoodInfos.FoodProInfo.FoodProPlace = args[9]

	ProInfosJSONasBytes, err := json.Marshal(FoodInfos)
	if err != nil{
		logger.Error(err.Error())
		return shim.Error(err.Error())
	}

	err = stub.PutState(FoodInfos.FoodID, ProInfosJSONasBytes)
	if err != nil {
		logger.Error(err.Error())
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// 新增配料信息
func (a *ChainCode) addIngInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var FoodInfos Infos
	var IngInfoitem IngInfo

	if (len(args) - 1) % 2 != 0 || len(args) == 1 {
		logger.Error("Incorrect number of arguments! Expecting number greater than 1 and an odd number!")
		return shim.Error("Incorrect number of arguments! Expecting number greater than 1 and an odd number!")
	}

	//  参数处理与赋值
	FID := args[0]
	for i := 1; i < len(args); i += 2 {
		IngInfoitem.IngID = args[i]
		IngInfoitem.IngName = args[i + 1]
		FoodInfos.FoodIngInfo = append(FoodInfos.FoodIngInfo, IngInfoitem)
	}

	FoodInfos.FoodID = FID
	IngInfoJsonAsBytes, err  := json.Marshal(FoodInfos)
	if err != nil {
		logger.Error(err.Error())
		return shim.Error(err.Error())
	}

	err = stub.PutState(FoodInfos.FoodID, IngInfoJsonAsBytes)
	if err != nil {
		logger.Error(err.Error())
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

