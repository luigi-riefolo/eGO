#!/usr/bin/env bash

# Treat unset variables as an error
set -o nounset
set -xe


SERVICE_NAME=


PROJECT="github.com/luigi-riefolo/eGO"
BASE=${GOPATH}/src/$PROJECT
SERVICE_DIR="${PROJECT}/src"
CONF_PATH=conf/global_conf.toml

CONFIG_DATA=
FULL_PATH=${GOPATH}/src/${SERVICE_DIR}
EXE_PATH=
REGISTRY=
PORT=
ME=$(basename $0)
ME_DIR=$(dirname $(readlink -f $0))

# TODO
# rollback mechanism, add a list of steps, for each one a function that gets
# the value true in the list of steps when it successfully complete within
# timeout, or fails and calls the rollback function
# detect changes
# detached mode
# reload the service with new vars
# dev/prod env
# use either env var or args and add the var name to help
# run tests, code coverage
# add go gets to dockerfile
# registry argument
# use the go config parser for the service vars
# go-bindata -o data/bindata.go -pkg data data/*.json



function usage {
    BOLD=$(tput bold)
    RESET=$(tput sgr0)
    cat <<-END
${BOLD}NAME${RESET}

    $ME -- microservice toolbox

${BOLD}SYNOPSIS${RESET}

    $ME [-h | --help] [-s | --service NAME] [-r | --registry ADDR]

${BOLD}DESCRIPTION${RESET}

    $ME builds and deploys a microservice.

${BOLD}OPTIONS${RESET}

    -s | --service NAME

        Microservice name.

    -r | --registry ADDR

        Image registry address [default localhost:5000].

    -h | --help

        Print this help message.


${BOLD}EXAMPLES${RESET}

    $ME --service alfa --registry localhost:5000

END
}

TEMP=`getopt -o h,r:,s: --long help,registry:,service: \
             -n "$ME" -- "$@"`
eval set -- "$TEMP"

while true; do
    case "$1" in
        -s | --service )
            SERVICE_NAME="$2"
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


SERVICE_NAME=${SERVICE_NAME,,}
SERVICE_DIR=$FULL_PATH/$SERVICE_NAME
EXE_PATH=$SERVICE_DIR/cmd/main.go
CMD_PATH=$SERVICE_DIR/cmd
REGISTRY="${REGISTRY:-localhost:5000}"


IMAGE_TAG="v1"
REPO_IMAGE="$REGISTRY/$SERVICE_NAME:$IMAGE_TAG"

# remove previous image
docker rmi -f $REPO_IMAGE || true

function rollback {
    echo "errexit on line $(caller)" >&2
    docker rmi -f $REPO_IMAGE || true
    rm $SERVICE_NAME
    clean_up
}

trap rollback ERR SIGHUP SIGTERM SIGINT SIGFPE


# Remove existing service
function clean_up {
    kubectl delete deployment --ignore-not-found=true $SERVICE_NAME
    kubectl delete service --ignore-not-found=true $SERVICE_NAME
}

function load_config {
    CONFIG_FILE=$BASE/$CONF_PATH
    $BASE/scripts/config.sh --out $CONFIG_FILE
    CONFIG_DATA="$(cat $BASE/$CONF_PATH)"
}

function build_image {
    echo -e "Starting build for microservice:\t$SERVICE_NAME"

    # stop running container
    if docker inspect -f '{{.State.Running}}' $SERVICE_NAME > /dev/null 2>&1;
    then
        echo "Stopping container"
        docker stop $SERVICE_NAME
    fi

    echo "Building new image"

    # compile
    #GOOS=linux GOARCH=amd64 go build -x -v -race -a -installsuffix cgo -o ${SERVICE_NAME} ${EXE_PATH}
    #GOOS=darwin
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -x -v \
        -a -installsuffix cgo -o $CMD_PATH/$SERVICE_NAME $EXE_PATH

    # build Docker image
    docker build --pull -t ${REPO_IMAGE} \
        -f ${SERVICE_DIR}/conf/Dockerfile \
        --build-arg GOPATH="$GOPATH" \
        --build-arg SERVICE_NAME="$SERVICE_NAME" \
        --build-arg PROJECT="$PROJECT" \
        --build-arg CONFIG_FILE="$CONF_PATH" \
        $BASE

        #$SERVICE_DIR

    # print info
    docker images $REPO_IMAGE --format 'table {{.ID}}\t{{.Repository}}\t\t{{.Size}}'
}

# publish image
function publish {
    echo "Publishing image"

    docker push $REPO_IMAGE

    docker images
}

# deploy
function deploy {
    echo "Deploying image"

    kubectl create -f $SERVICE_DIR/conf/k8s-deployment.yaml
    kubectl describe deployment $SERVICE_NAME
    kubectl create -f $SERVICE_DIR/conf/k8s-service.yaml

#    PATTERN=".${SERVICE_NAME^}.Server.Port"
#    IS_GATEWAY=$(jq -r ".${SERVICE_NAME^}.Server.IsGateway" <<<"$CONFIG_DATA")
#    if [[ $IS_GATEWAY == "true" ]]
#    then
#        PATTERN=".${SERVICE_NAME^}.Server.GatewayPort"
#    fi
#
#    PORT=$(jq -r $PATTERN <<<"$CONFIG_DATA")
#    #kubectl get pods --all-namespaces -l service-name --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}'
#
    minikube service $SERVICE_NAME --url
}

function go_code {

	${ME_DIR}/go_checker.sh ${SERVICE_DIR}

    go test -coverprofile cover.outs ${SERVICE_DIR}

    # show in browser
    go tool cover -html=cover.out -o cover.html
}

function tidy_up {
    echo "Tidying up"

    #docker rm "/${SERVICE_NAME}"
    rm -f $CMD_PATH/$SERVICE_NAME
}

eval $(minikube docker-env)

# TODO: use a deployment update instead of removing it
clean_up

load_config

#go_code

build_image

publish

deploy

tidy_up
