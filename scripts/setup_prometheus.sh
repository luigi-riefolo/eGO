#!/usr/bin/env bash

# Script for setting up prometheus in Kubernetes.
# See the config files in deployments/prometheus for additional info.

set -ex

PD=$(pwd)
cd $(dirname "$(readlink -f "$0")")/../deployments/prometheus

kubectl delete configmap prometheus-server-conf --ignore-not-found=true -n monitoring
kubectl delete svc prometheus-service --ignore-not-found=true -n monitoring
kubectl delete deployment prometheus-deployment --ignore-not-found=true -n monitoring
kubectl delete namespace monitoring --ignore-not-found=true

kubectl get namespace monitoring || kubectl create namespace monitoring

kubectl create -f config-map.yaml -n monitoring

kubectl create -f deployment.yaml -n monitoring

kubectl describe deployment prometheus-deployment -n monitoring

kubectl create -f service.yaml -n monitoring

kubectl get svc -n monitoring

minikube service prometheus-service -n monitoring

cd $PD
