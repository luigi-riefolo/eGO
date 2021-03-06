#!/usr/bin/env bash

# Script for automating the setup of kubectl.

set -xe

# Get sudo priviledges
echo "Sudoer password"
sudo -v

# Kubernetes
echo "Installing kubectl"
[[ -e /usr/local/bin/kubectl ]] && rm /usr/local/bin/kubectl
STABLE="$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)"
curl -sLO https://storage.googleapis.com/kubernetes-release/release/$STABLE/bin/darwin/amd64/kubectl

chmod +x ./kubectl

sudo mv ./kubectl /usr/local/bin/kubectl

source <(kubectl completion bash)
echo "source <(kubectl completion bash)" >> ~/.bash_profile

brew uninstall docker-machine-driver-xhyve || true
brew install docker-machine-driver-xhyve
sudo chown root:wheel $(brew --prefix)/opt/docker-machine-driver-xhyve/bin/docker-machine-driver-xhyve
sudo chmod u+s $(brew --prefix)/opt/docker-machine-driver-xhyve/bin/docker-machine-driver-xhyve

rm -f /usr/local/bin/docker-machine

brew link --overwrite docker-machine

brew services restart docker-machine
