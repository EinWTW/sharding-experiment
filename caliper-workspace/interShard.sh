#!/bin/bash
# interShardingB
npx caliper launch manager --caliper-workspace ../caliper-workspace/ --caliper-networkconfig networks/interShardingA.yaml --caliper-benchconfig benchmarks/interShardingA.yaml --caliper-flow-only-test --caliper-fabric-gateway-enabled & npx caliper launch manager --caliper-workspace ../caliper-workspaceB/ --caliper-networkconfig networks/interShardingB.yaml --caliper-benchconfig benchmarks/interShardingB.yaml --caliper-flow-only-test --caliper-fabric-gateway-enabled && fg


