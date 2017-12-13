#!/usr/bin/env bash

# Script for automating the setup of minikube.
# NOTE: this script will always update the
# current version of minikube if it's out of date.


INSTALLED_VERSION=$(minikube version | grep -Po 'v\d+.*')

MINIKUBE_VERSION=$(curl -s -L -o /dev/null -w %{url_effective} \
    https://github.com/kubernetes/minikube/releases/latest | grep -Po '/\K(v.*)')

ARCH="minikube-linux-amd64"
[[ $OSTYPE == darwin* ]] && ARCH="minikube-darwin-amd64"

if [[ "$INSTALLED_VERSION" != "$MINIKUBE_VERSION" ]]
then
    echo "Installing minikube"
    [[ -e /usr/local/bin/minikube ]] && {
        minikube delete || true
        sudo rm /usr/local/bin/minikube
        sudo rm -rf ~/.minikube
    }

    curl -Lo minikube \
        https://storage.googleapis.com/minikube/releases/$MINIKUBE_VERSION/$ARCH
    chmod +x minikube
    sudo mv minikube /usr/local/bin/
fi

STATUS=$(minikube status --format {{.MinikubeStatus}})
if [[ "$STATUS" != "Running" ]]
then
    minikube start --vm-driver=xhyve
    eval $(minikube docker-env)
fi

minikube status
