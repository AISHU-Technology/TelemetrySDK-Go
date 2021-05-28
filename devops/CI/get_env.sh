#!/bin/bash

SCRIPTPATH=$(dirname $(readlink -f "$0"))
CODESPACE=$(dirname $(dirname ${SCRIPTPATH}))
BRANCHNAME=${BUILD_SOURCEBRANCHNAME}
SERVICENAME="cloud-control-plane"
IMGREPO="acr.aishu.cn/proton"
BUILDID=${BUILD_BUILDNUMBER}


function defaultValue(){
  if [ "$BRANCHNAME"x == "x" ]
  then
    BRANCHNAME=$(git symbolic-ref --short -q HEAD)
  fi

  if [ "$BUILDID"x == "x" ]
  then
    BUILDID=$RANDOM
  fi
}

defaultValue
