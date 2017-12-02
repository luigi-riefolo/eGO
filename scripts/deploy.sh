#!/usr/bin/env bash

# Script for deploying a service to Kubernetes.

set -xe

SERVICE_LIST="$(ls -x $BASE/deployments)"

if [[ $# -eq 0 ]]
then
    echo -e "Please supply a valid service name\nAvailable services: $SERVICE_LIST"
    exit
fi

SERVICE="$1"

if [[ ! -d $BASE/deployments/$SERVICE ]]
then
    echo "Service $SERVICE is not supported"
    exit
fi

# Namespace
NAMESPACE="$(grep -Po 'namespace:\s*\K(.*)' $BASE/deployments/$SERVICE/deployment.yaml)"
kubectl get namespace $NAMESPACE || kubectl create namespace $NAMESPACE


# Config map
kubectl delete configmap $SERVICE-config --ignore-not-found=true -n $NAMESPACE
FILES="$(ls $BASE/deployments/$SERVICE/conf)"
CONF_FILES=
for F in $FILES
do
    CONF_FILES+="--from-file=$BASE/deployments/$SERVICE/conf/$F "
done
kubectl create configmap $SERVICE-config $CONF_FILES -o json -n $NAMESPACE


# Deployment
kubectl create -f $BASE/deployments/$SERVICE/deployment.yaml || \
    kubectl replace -f $BASE/deployments/$SERVICE/deployment.yaml --record

kubectl describe deployment $SERVICE-deployment -n $NAMESPACE


# Service
#kubectl delete svc $SERVICE-service --ignore-not-found=true -n $NAMESPACE
kubectl create -f $BASE/deployments/$SERVICE/service.yaml -n $NAMESPACE || \
    kubectl apply -f $BASE/deployments/$SERVICE/service.yaml --record

kubectl get svc -n $NAMESPACE
minikube service $SERVICE -n $NAMESPACE
