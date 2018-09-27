package control

import (
	"contract/log"
	"contract/module"
	"contract/services"
	"encoding/json"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

func (t *Contract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	log.Logger.Info("Invoke")
	funcation, args := stub.GetFunctionAndParameters()
	lowFuncation := strings.ToLower(funcation)

	if lowFuncation == "init" { // 合同上链
		return t.Init(stub, args)
	}
	if lowFuncation == "transfer" { // 交易
		return t.Transfer(stub, args)
	}
	if lowFuncation == "query" { // 查询
		return t.Query(stub, args)
	}
	return shim.Error("Invalid invoke function name. " + funcation)
}

/**  **/
func (t *Contract) Init(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	log.Logger.Info("##############调用合同初始化接口开始###############")
	returnInfo := module.ReturnInfo{}
	if len(args) >= 1 {
		var initParam module.InitParam
		err := json.Unmarshal([]byte(args[0]), &initParam)
		if err != nil {
			log.Logger.Error("Init:err" + err.Error())
			returnInfo.Success = false
			returnInfo.Info = err.Error()
		} else {
			chaninfo := services.ToInit(stub, initParam)
			// return response
			jsonreturn, err := json.Marshal(chaninfo)
			if err != nil {
				return shim.Error("err:" + err.Error())
			}
			return shim.Success(jsonreturn)
		}
	} else {
		log.Logger.Error("Init:参数不对，请核实参数信息。")
		returnInfo.Success = false
		returnInfo.Info = "参数不对，请核实参数信息"
	}
	jsonreturn, err := json.Marshal(returnInfo)
	if err != nil {
		return shim.Error("err:" + err.Error())
	}
	return shim.Success(jsonreturn)
}

/** 交易上链 **/
func (t *Contract) Transfer(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	log.Logger.Info("##############调用交易接口开始###############")
	returnInfo := module.ReturnInfo{}
	if len(args) >= 1 {
		var transferParam module.TransferParam
		err := json.Unmarshal([]byte(args[0]), &transferParam)
		if err != nil {
			log.Logger.Error("transfer :err" + err.Error())
			returnInfo.Success = false
			returnInfo.Info = err.Error()
		} else {
			chaninfo := services.ToTransfer(stub, transferParam)
			// return response
			jsonreturn, err := json.Marshal(chaninfo)
			if err != nil {
				return shim.Error("err:" + err.Error())
			}
			return shim.Success(jsonreturn)
		}
	} else {
		log.Logger.Error("transfer:参数不对，请核实参数信息。")
		returnInfo.Success = false
		returnInfo.Info = "参数不对，请核实参数信息"
	}
	jsonreturn, err := json.Marshal(returnInfo)
	if err != nil {
		return shim.Error("err:" + err.Error())
	}
	return shim.Success(jsonreturn)
}

/** 查询上链 **/
func (t *Contract) Query(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	log.Logger.Info("##############调用查询接口开始###############")
	returnInfo := module.ReturnInfo{}
	if len(args) >= 1 {
		var queryParam module.QueryParam
		err := json.Unmarshal([]byte(args[0]), &queryParam)
		if err != nil {
			log.Logger.Error("query :err" + err.Error())
			returnInfo.Success = false
			returnInfo.Info = err.Error()
		} else {
			chaninfo := services.ToQuery(stub, queryParam)
			// return response
			jsonreturn, err := json.Marshal(chaninfo)
			if err != nil {
				return shim.Error("err:" + err.Error())
			}
			return shim.Success(jsonreturn)
		}
	} else {
		log.Logger.Error("query:参数不对，请核实参数信息。")
		returnInfo.Success = false
		returnInfo.Info = "参数不对，请核实参数信息"
	}
	jsonreturn, err := json.Marshal(returnInfo)
	if err != nil {
		return shim.Error("err:" + err.Error())
	}
	return shim.Success(jsonreturn)
}
