#!/usr/bin/env bash

# Script for setting up grafana in Kubernetes.
# See the config files in deployments/grafana for additional info.

set -ex

PD=$(pwd)
cd $(dirname "$(readlink -f "$0")")/../deployments/grafana

kubectl delete configmap grafana-server-conf --ignore-not-found=true -n monitoring
kubectl delete svc grafana-service --ignore-not-found=true -n monitoring
kubectl delete deployment grafana-deployment --ignore-not-found=true -n monitoring

kubectl get namespace monitoring || kubectl create namespace monitoring

#kubectl create -f config-map.yaml

kubectl create -f deployment.yaml

kubectl describe deployment grafana-deployment -n monitoring

kubectl create -f service.yaml

kubectl get svc -n monitoring

minikube service grafana -n monitoring

cd $PD
