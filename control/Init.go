package control

import (
	"contract/log"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type Contract struct {
}

func (t *Contract) Init(stub shim.ChaincodeStubInterface) peer.Response {
	log.Logger.Info("Init")
	return shim.Success(nil)
}
