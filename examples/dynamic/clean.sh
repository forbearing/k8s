#!/usr/bin/env bash

kubectl delete pod --force pod-unstruct
kubectl delete ns --force ns-unstruct
kubectl delete deploy --force mydep-unstruct
