#!/usr/bin/env bash

set -ex

cd ~/Workspace/TMP

rm -f main

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -x -a -installsuffix cgo -o main server.go

docker build -t localhost:5000/alfa:v10 .

docker push localhost:5000/alfa:v10

#docker run -i -t -p 8080:8080 localhost:5000/alfa:v10


kubectl create -f k8s-deployment.yaml
kubectl describe deployment alfa
kubectl create -f k8s-service.yaml

minikube service alfa --url


#kubectl run $SERVICE_ARTIFACT --image=$REPO_IMAGE --port=$SERVICE_PORT
#
#kubectl get pods --all-namespaces
#
## NOTE: the external IP won't be shown, for
## minikube does not offer  a LoadBalancer
#kubectl expose deployment $SERVICE_ARTIFACT --type=LoadBalancer
#
#kubectl get services $SERVICE_ARTIFACT
#
#minikube service $SERVICE_ARTIFACT
#
#curl $(minikube service $SERVICE_ARTIFACT --url)
