/*
Copyright 2020 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run crossShards.go <accountA> <acountB> <amount>")
		return
	}

	account1 := os.Args[1]
	account2 := os.Args[2]
	amount := os.Args[3]

	channel1 := "channel1"
	chaincode1 := "sharding1"
	channel2 := "channel2"
	chaincode2 := "sharding2"

	if strings.HasPrefix(account1, "A") {
		channel1 = "channel1"
		chaincode1 = "sharding1"
	} else if strings.HasPrefix(account1, "B") {
		channel1 = "channel2"
		chaincode1 = "sharding2"
	}

	if strings.HasPrefix(account2, "A") {
		channel2 = "channel1"
		chaincode2 = "sharding1"
	} else if strings.HasPrefix(account2, "B") {
		channel2 = "channel2"
		chaincode2 = "sharding2"
	}
	//log.Println("============ cross-shards transfer starts ============")

	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if err != nil {
		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
	}

	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		log.Fatalf("Failed to create wallet: %v", err)
	}

	if !wallet.Exists("appUser") {
		err = populateWallet(wallet)
		if err != nil {
			log.Fatalf("Failed to populate wallet contents: %v", err)
		}
	}

	ccpPath := filepath.Join(
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"connection-org1.yaml",
	)

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "appUser"),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()

	network1, err := gw.GetNetwork(channel1)
	if err != nil {
		log.Fatalf("Failed to get network: %v", err)
	}
	contract1 := network1.GetContract(chaincode1)

	network2, err := gw.GetNetwork(channel2)
	if err != nil {
		log.Fatalf("Failed to get network: %v", err)
	}
	contract2 := network2.GetContract(chaincode2)

	//log.Println("--> Submit Transaction: AcquireLock ")

	// Add lock to accounts
	identifier := account1 + account2
	result, err := contract1.SubmitTransaction("AcquireLock", account1, identifier)
	if err != nil {
		log.Fatalf("Failed to get write lock: %v", err)
	}
	//log.Println(string(result))
	result, err = contract2.SubmitTransaction("AcquireLock", account2, identifier)
	if err != nil {
		log.Fatalf("Failed to get write lock: %v", err)
	}
	//log.Println(string(result))

	// result, err = contract1.SubmitTransaction("WLockExistsCrossShards", account2, identifier)
	// if err != nil {
	// 	log.Fatalf("Failed to get cross shard write lock: %v", err)
	// }
	// log.Println(string(result))

	// Transfer on shards
	_, err = contract1.SubmitTransaction("SendAmountCrossShards", account1, account2, amount)
	if err != nil {
		log.Printf("Failed to Submit transaction: %v", err)
	}
	_, err = contract2.SubmitTransaction("SendAmountCrossShards", account1, account2, amount)
	if err != nil {
		log.Printf("Failed to Submit transaction: %v", err)
	}

	// Remove lock to accounts
	_, err = contract1.SubmitTransaction("DeleteLock", account1, identifier)
	if err != nil {
		log.Printf("Failed to delete write lock: %v", err)
	}
	_, err = contract2.SubmitTransaction("DeleteLock", account2, identifier)
	if err != nil {
		log.Printf("Failed to delete write lock: %v", err)
	}

	// Result
	// result, err = contract1.EvaluateTransaction("GetBalance", account1)
	// if err != nil {
	// 	log.Fatalf("Failed to evaluate transaction: %v\n", err)
	// }
	// log.Println("--> GetBalance: " + account1 + " " + string(result))
	// result, err = contract2.EvaluateTransaction("GetBalance", account2)
	// if err != nil {
	// 	log.Fatalf("Failed to evaluate transaction: %v\n", err)
	// }
	// log.Println("--> GetBalance: " + account2 + " " + string(result))

	// log.Println("============ application-golang ends ============")
}

func populateWallet(wallet *gateway.Wallet) error {
	log.Println("============ Populating wallet ============")
	credPath := filepath.Join(
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"users",
		"User1@org1.example.com",
		"msp",
	)

	certPath := filepath.Join(credPath, "signcerts", "User1@org1.example.com-cert.pem")
	// read the certificate pem
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return fmt.Errorf("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))

	return wallet.Put("appUser", identity)
}
