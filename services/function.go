package services

import (
	"contract/common"
	"contract/log"
	"contract/module"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// 合约初始化
func ToInit(stub shim.ChaincodeStubInterface, param module.InitParam) (tChan module.ChanInfo) {
	// 	verify product if exist or not
	jsonParam, err := stub.GetState(common.CONTRACT_INFO + common.ULINE + param.Name)
	log.Logger.Info("------------------------------------------------------------------")
	log.Logger.Info(string(jsonParam[:]))
	if jsonParam != nil {
		log.Logger.Error("Init -- get contract by contract name -- err: 已经发布" + "	name:" + param.Name)
		tChan.ContractName = param.Name
		tChan.Success = false
		tChan.Info = "已经发布"
		return
	}

	if common.GetUserFromCertification(stub) != param.Issuer {
		log.Logger.Error("Init -- get contract by contract name -- err: 发布人和登录人不对" + "	name:" + param.Name)
		tChan.ContractName = param.Name
		tChan.Success = false
		tChan.Info = "发布人和登录人不符"
		return
	}

	erc := module.ERC{}
	erc.Name = param.Name
	erc.Decimals = param.Decimals
	erc.Symbol = param.Symbol
	erc.TotalSupply = param.TotalSupply * uint64(math.Pow(float64(10), float64(param.Decimals)))
	balance := make(map[string]uint64)
	balance[param.Issuer] = erc.TotalSupply
	erc.BalanceOf = balance

	jsonByte, err := json.Marshal(erc)
	if err != nil {
		log.Logger.Error("Init -- marshal product err:" + err.Error() + "	name:" + param.Name)
		tChan.ContractName = param.Name
		tChan.Success = false
		tChan.Info = err.Error()
		return
	}

	loged := TransferLog(param.Issuer, param.Name, fmt.Sprint("Init ", param.Name), param.Issuer, param.Issuer, erc.TotalSupply)

	if loged == false {
		tChan.ContractName = param.Name
		tChan.Success = false
		tChan.Info = "操作日志记录错误，请重试"
		return
	}

	err = stub.PutState(common.CONTRACT_INFO+common.ULINE+param.Name, jsonByte)
	if err != nil {
		log.Logger.Error("Init -- putState:" + err.Error() + "	name:" + param.Name)
		tChan.ContractName = param.Name
		tChan.Success = false
		tChan.Info = err.Error()
		return
	}
	tChan.ContractName = param.Name
	tChan.Success = true
	tChan.Info = "发布完成"
	return
}

func ToTransfer(stub shim.ChaincodeStubInterface, param module.TransferParam) (tChan module.ChanInfo) {
	// 	verify product if exist or not
	jsonParam, err := stub.GetState(common.CONTRACT_INFO + common.ULINE + param.Name)
	log.Logger.Info("------------------------------------------------------------------")
	log.Logger.Info(string(jsonParam[:]))
	if jsonParam == nil {
		log.Logger.Error("Transfer -- get asset by assetid -- err: 该资产未发布" + "	name:" + param.Name)
		tChan.ContractName = param.Name
		tChan.Success = false
		tChan.Info = "该资产未发布"
		return
	}

	if common.GetUserFromCertification(stub) != param.From {
		tChan.ContractName = param.Name
		tChan.Success = false
		tChan.Info = "登录人和转出人不符"
		return
	}

	erc := module.ERC{}
	err = json.Unmarshal(jsonParam, &erc)
	if err != nil {
		tChan.ContractName = param.Name
		tChan.Success = false
		tChan.Info = err.Error()
		return
	}
	// balance for from more than value
	if val, ok := erc.BalanceOf[param.From]; ok {
		if val >= param.Value {
			erc.BalanceOf[param.From] = val - param.Value
			if toVal, toOk := erc.BalanceOf[param.To]; toOk {
				erc.BalanceOf[param.To] = toVal + param.Value
			} else {
				erc.BalanceOf[param.To] = param.Value
			}
		} else {
			tChan.ContractName = param.Name
			tChan.Success = false
			tChan.Info = fmt.Sprint(param.From, "-- value:", erc.BalanceOf[param.From])
			return
		}
	} else {
		tChan.ContractName = param.Name
		tChan.Success = false
		tChan.Info = fmt.Sprint(param.From, "-- value: 0")
		return
	}

	jsonByte, err := json.Marshal(erc)
	if err != nil {
		tChan.ContractName = param.Name
		tChan.Success = false
		tChan.Info = err.Error()
		return
	}

	loged := TransferLog(param.From, param.Name, fmt.Sprint("Transfer ", param.Name), param.From, param.To, param.Value)

	if loged == false {
		tChan.ContractName = param.Name
		tChan.Success = false
		tChan.Info = "操作日志记录错误，请重试"
		return
	}

	err = stub.PutState(common.CONTRACT_INFO+common.ULINE+param.Name, jsonByte)
	if err != nil {
		log.Logger.Error("Transfer -- putState:" + err.Error() + "	name:" + param.Name)
		tChan.ContractName = param.Name
		tChan.Success = false
		tChan.Info = err.Error()
		return
	}
	tChan.ContractName = param.Name
	tChan.Success = true
	tChan.Info = "交易完成"
	return
}

/** 查询接口 **/
func ToQuery(stub shim.ChaincodeStubInterface, param module.QueryParam) (tChan module.QueryInfo) {
	// 	verify product if exist or not
	jsonParam, err := stub.GetState(common.CONTRACT_INFO + common.ULINE + param.Name)
	log.Logger.Info("------------------------------------------------------------------")
	log.Logger.Info(string(jsonParam[:]))
	if jsonParam == nil {
		log.Logger.Error("Transfer -- get asset by assetid -- err: 该资产未发布" + "	name:" + param.Name)
		tChan.Address = param.Address
		tChan.Success = false
		tChan.Value = 0
		tChan.Info = "该资产未发布"
		return
	}

	erc := module.ERC{}
	err = json.Unmarshal(jsonParam, &erc)
	if err != nil {
		tChan.Address = param.Address
		tChan.Success = false
		tChan.Value = 0
		tChan.Info = err.Error()
		return
	}
	// balance for from more than value
	if val, ok := erc.BalanceOf[param.Address]; ok {
		tChan.Address = param.Address
		tChan.Success = true
		tChan.Value = val
		tChan.Info = ""
	} else {
		tChan.Address = param.Address
		tChan.Success = false
		tChan.Value = 0
		tChan.Info = "用户没有购买该资产"
		return
	}
	return
}

/** 记录日志 **/
func TransferLog(curuser string, name string, operation string, from string, to string, value uint64) bool {
	jsonParam, err := stub.GetState(common.CONTRACT_TRANSFER + common.ULINE + name)
	logTran := module.TransferLog{}
	if jsonParam != nil {
		err = json.Unmarshal(jsonParam, &logTran)
		log.Logger.Error("TransferLog --err:" + err.Error())
		return false
	}
	tran := module.Transfer{}
	tran.TxHash = stub.GetTxID()
	tran.From = from
	tran.To = to
	tran.Value = value
	tran.OperateTime = time.Now().Unix()
	tran.Operation = operation
	tran.Operator = curuser
	if trans, ok := logTran[curuser]; ok {
		logTran[curuser] = append(trans, tran)
	} else {
		tmp := make([]module.Transfer, 0)
		logTran[curuser] = append(tmp, tran)
	}

	jsonByte, err := json.Marshal(logTran)
	if err != nil {
		log.Logger.Error("TransferLog --err:" + err.Error())
		return false
	}

	err = stub.PutState(common.CONTRACT_TRANSFER+common.ULINE+name, jsonByte)
	if err != nil {
		log.Logger.Error("TransferLog -- putState:" + err.Error())
		return false
	}
	return true
}
