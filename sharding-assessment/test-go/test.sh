#!/bin/bash
#
#

set -ev

###
function runNoShardTest() {
  for (( i=1; i<=$count; ++i)); do
    go run $DIR/noShard.go A"$i" B"$i" &
  done
  wait $!
}

function runInterShardTest() {
  for (( i=1; i<=$count; ++i)); do
    number1=$((1 + $RANDOM % $RANGE))
    number2=$((1 + $RANDOM % $RANGE))
    go run $DIR/interShard.go A"$number1" A"$number2" 1 &
    go run $DIR/interShard.go B"$number1" B"$number2" 1 &
  done
  wait $!
}

function runCrossShardTest() {
  for (( i=1; i<=$count; ++i)); do
    number1=$((1 + $RANDOM % $RANGE))
    time go run $DIR/crossShard.go A"$number1" B"$number1" 1 &
  done
  wait $!
}
###time 

echo
echo "===================== Start Tests ===================== "
echo

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

testcase="$1"
count="$2"
RANGE=10000
starttime=`date +%s.%N`

if [ "$testcase" == "noshard" ]; then
  echo "Test Case: noshard "
  runNoShardTest
elif [ "$testcase" == "intershard" ]; then
  echo "Test case: intershard"
  runInterShardTest
elif [ "$testcase" == "crossshard" ]; then
  echo "Test case: crossshard"
  runCrossShardTest
else
  echo "Warning: Only support test cases [noshard|intershard|crossshard]"
  exit 1
fi

end=`date +%s.%N`
# 
runtime=$( echo "$end - $starttime" | bc -l )
echo "TPS: "
echo "scale=2 ; $count / $runtime" | bc

echo "===================== End Tests ===================== "
echo

exit 0