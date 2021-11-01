####################### Resources #####################
caliper-workspae/:  caliper project [README, benchmarks, networks, workload]
sharding-assessment/: [chaincode-go, application-go, test-go, bin/]
useful-script-if-needed/: [setupEnv.sh, restart.sh]
README : That's me!

######################   Steps   ######################
# 1.Setup Fabric Environmenton Linux server(skip if your fabric environment is ready)
# The following script will help you install all and create workspace '~/go/src/github.com/hyperledger'
./setupEnv.sh

# 2. Start test-network
# Move restart.sh to folder 'fabric-samples/test-network'
# Move 'sharding-assessment' to folder 'fabric-samples/'
./restart.sh
# The above script will start a blockchain test nework, create channels and deploy chaincode for testing.

# 3. Init accounts
cd fabric-samples/sharding-assessment/
./initAllCounts.sh

# 4. Caliper tests
# Move caliper-workspace to path '~/go/src/github.com/hyperledger/'
# Please refer the 'caliper-workspace/README' to begin your testing.