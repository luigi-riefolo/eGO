#!/usr/bin/env bash

# Script for automating the setup of minikube.

minikube version >/dev/null

if [[ $? -ne 0 ]]
then
    echo "Installing minikube"
    [[ -e /usr/local/bin/minikube ]] && {
        minikube delete || true
        sudo rm /usr/local/bin/minikube
        sudo rm -rf ~/.minikube
    }

    curl -sLo minikube https://storage.googleapis.com/minikube/releases/v0.23.0/minikube-darwin-amd64
    chmod +x minikube
    sudo mv minikube /usr/local/bin/
fi


STATUS=$(minikube status --format {{.MinikubeStatus}})
if [[ "$STATUS" != "Running" ]]
then
    minikube start --vm-driver=xhyve
    eval $(minikube docker-env)
fi
