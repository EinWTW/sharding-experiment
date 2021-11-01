/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-samples/sharding-assessment/chaincode-go/chaincode"
)

func main() {
	shardingChaincode, err := contractapi.NewChaincode(&chaincode.SmartContract{})
	if err != nil {
		log.Panicf("Error creating sharding-assessment chaincode: %v", err)
	}

	if err := shardingChaincode.Start(); err != nil {
		log.Panicf("Error starting sharding-assessment chaincode: %v", err)
	}
}
