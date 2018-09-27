package main

import (
	"contract/control"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func main() {
	err := shim.Start(new(control.Contract))
	if err != nil {
		fmt.Printf("Error starting Contract: %s", err)
	}
}
