#!/bin/bash 

## Please ensure run this script after  login docker regitry
set -e 
source get_env.sh

BuildImage="golang:1.15"

function UT(){
  echo "|------UT-------|"
  docker run --rm -v $CODESPACE/span:/root/go  $BuildImage sh -c "cd /root/go && sh -x test_cover.sh"
}

UT
