#!/bin/bash
#

#./network.sh -h
./network.sh down
./network.sh up createChannel
docker ps -a
./network.sh deployCC -ccn basic -ccp ../asset-transfer-basic/chaincode-go -ccl go
./network.sh deployCC -ccn sharding -ccp ../sharding-assessment/chaincode-go -ccl go
./network.sh createChannel -c channel1
./network.sh deployCC -c channel1 -ccn sharding1 -ccp ../sharding-assessment/chaincode-go -ccl go
./network.sh createChannel -c channel2
./network.sh deployCC -c channel2 -ccn sharding2 -ccp ../sharding-assessment/chaincode-go -ccl go

echo "=================== Success ==================="
