#!/usr/bin/env bash

namespace="test"
exampleFiles="./examples"

if [[ $1 == "-u" ]];then
    kubectl -n $namespace delete -f $exampleFiles
    kubectl delete namespace $namespace
    exit 0
fi

kubectl create namespace $namespace 2> /dev/null
kubectl -n $namespace apply -f $exampleFiles
