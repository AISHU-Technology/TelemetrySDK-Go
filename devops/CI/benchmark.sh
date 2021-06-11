#!/bin/bash 

## Please ensure run this script after  login docker regitry
set -e 
SCRIPTPATH=$(dirname $(readlink -f "$0"))
source $SCRIPTPATH/get_env.sh

BuildImage="golang:1.15"

function Benchmark(){
  echo "|------Benchmark-------|"
  docker run --rm -v $CODESPACE/span:/root/go  $BuildImage sh -c "cd /root/go && sh -x benchtest.sh"
}

Benchmark
