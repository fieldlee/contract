package services

import (
	"contract/common"
	"contract/log"
	"contract/module"
	"encoding/json"
	"fmt"
	"math"

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
	erc.TotalSupply = param.TotalSupply * uint64(math.Pow(10, param.Decimals))
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
	if erc.BalanceOf[param.From] == nil {
		tChan.ContractName = param.Name
		tChan.Success = false
		tChan.Info = fmt.Sprint(param.From, "-- value: 0")
		return
	} else {
		if erc.BalanceOf[param.From] >= param.Value {
			erc.BalanceOf[param.From] = erc.BalanceOf[param.From] - param.Value
			if erc.BalanceOf[param.To] == nil {
				erc.BalanceOf[param.To] = param.Value
			} else {
				erc.BalanceOf[param.To] = erc.BalanceOf[param.To] + param.Value
			}
		} else {
			tChan.ContractName = param.Name
			tChan.Success = false
			tChan.Info = fmt.Sprint(param.From, "-- value:", erc.BalanceOf[param.From])
			return
		}
	}

	jsonByte, err := json.Marshal(erc)
	if err != nil {
		tChan.ContractName = param.Name
		tChan.Success = false
		tChan.Info = err.Error()
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
	if erc.BalanceOf[param.Address] == nil {
		tChan.Address = param.Address
		tChan.Success = false
		tChan.Value = 0
		tChan.Info = "用户没有购买该资产"
		return
	} else {
		tChan.Address = param.Address
		tChan.Success = true
		tChan.Value = erc.BalanceOf[param.Address]
		tChan.Info = "用户没有购买该资产"
	}
	return
}
