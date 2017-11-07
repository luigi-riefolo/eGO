#!/usr/bin/env bash

set -ex

# NOTE: the port is hardcoded, for kubectl
# doesn't currently support templating.

ME=$(basename $0)
DIR=$(dirname "$(readlink -f $0)")

echo $DIR

USAGE_EXAMPLE="$ME --docker-user mario --user-email mario.mario@bros.com"
REGISTRY="localhost:5000"

function usage {
    BOLD=$(tput bold)
    RESET=$(tput sgr0)
    cat <<-END
${BOLD}NAME${RESET}

    $ME -- create a local registry.


${BOLD}SYNOPSIS${RESET}

    $ME [-h | --help] [-u | --docker-user USER] [-e | --docker-email] [-r | --registry host:port]


${BOLD}OPTIONS${RESET}

    -e | --docker-email

        User's Docker email.

    -u | --docker-user

        The Docker user.

    -r | --registry

        Registry address (default localhost:5000).

    -h | --help

        Print this help message.


${BOLD}EXAMPLES${RESET}

    $USAGE_EXAMPLE

END
}


TEMP=`getopt -o h,u:,e:,r: --long help,docker-user:,user-email:,registry: \
             -n "$ME" -- "$@"`
eval set -- "$TEMP"

while true; do
    case "$1" in
        -u | --docker-user )
            DOCKER_USER="$2"
            shift 2
            ;;

        -e | --user-email )
            USER_EMAIL="$2"
            shift 2
            ;;

        -r | --registry )
            REGISTRY="$2"
            shift 2
            ;;

        -h | --help )
            shift
            usage
            exit 0
            ;;

        -- )
            shift
            break
            ;;

        * )
            break
            ;;
    esac
done

[[ -z $DOCKER_USER ]] && {
    echo "Please supply a Docker user"
    exit 1
}

[[ -z $USER_EMAIL ]] && {
    echo "Please supply user's email"
    exit 1
}


read -s -p "Please type the Docker password:    " DOCKER_PASS

echo "Docker login"
docker login --username "${DOCKER_USER}" --password "$DOCKER_PASS"

MINIKUBE_STATUS=$(minikube status --format {{.MinikubeStatus}})
if [[ $MINIKUBE_STATUS != "Running" ]]
then
    echo "Starting minikube"
    minikube start --vm-driver=xhyve
fi

# delete previous pods
echo "Removing previous registry"
# NOTE:
# currently kubectl returns before completely removing a resource
kubectl delete replicationcontroller --force --ignore-not-found=true --namespace=kube-system kube-registry-v0
kubectl delete pod --force --ignore-not-found=true --namespace=kube-system kube-registry-proxy
kubectl delete service --force --ignore-not-found=true --namespace=kube-system kube-registry
kubectl delete secret --force --ignore-not-found=true regsecret


echo "Registry Docker configuration"
cat ~/.docker/config.json

echo "Starting minikube"
minikube start \
    --vm-driver xhyve \
    --insecure-registry "$REGISTRY"

# start talking to the docker daemon inside the minikube VM
eval $(minikube docker-env)

# create a secret
echo "Creating secret"
kubectl create secret docker-registry regsecret \
    --docker-server "${REGISTRY}"/registry \
    --docker-username "$DOCKER_USER" \
    --docker-password "$DOCKER_PASS" \
    --docker-email "$USER_EMAIL"

echo $(kubectl get secret regsecret --output=yaml | grep -Po "cfg: \K(.+)") > $DIR/secret64

kubectl get pods --namespace kube-system


echo "Creating registry"
kubectl create -f $DIR/local-registry.yaml

sleep 15

kubectl get pods --namespace kube-system

POD=$(kubectl get pod -n kube-system | grep kube-registry-v0 | awk '{print $1;}')
while true
do
    POD_STATUS=$(kubectl get pod --namespace kube-system "$POD" --template={{.status.phase}})
    [[ $POD_STATUS == "Running" ]] && break
    sleep 5
done

echo "Starting registry"
PORT=${REGISTRY##*:}
#nohup  kubectl port-forward --namespace kube-system "$POD" $PORT:$PORT &>/dev/null &
kubectl port-forward --namespace kube-system "$POD" $PORT:$PORT

# NOTE: kubectl fails to tunnel the first request
curl -s http://"${REGISTRY}"/v2/_catalog 2>&1 >/dev/null

echo "Done"
