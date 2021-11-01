#!/bin/bash
#
#

set -ev

echo "============ Init Accounts =============="


rm -r wallet keystore
go run initShard.go mychannel A 10000
go run initShard.go mychannel B 10000

go run initShard.go channel1 A 10000

go run initShard.go channel2 B 10000