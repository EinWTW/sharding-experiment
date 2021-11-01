#!/bin/bash
#

npx caliper launch manager --caliper-workspace ./ --caliper-networkconfig networks/noShardingNetworkConfig.yaml --caliper-benchconfig benchmarks/noShardingBenchmark.yaml --caliper-flow-only-test --caliper-fabric-gateway-enabled
