#!/usr/bin/env bash

# Script for automating the creation of gRPC
# gateways/services and Swagger definitions.

set -e

# TODO:
# USE GO TEMPLATES!!!!!!!!!!!!!!!!!!!!!!!!


declare -A CONFIG_MAP
SERVICE_NAME=
PROJECT=
ME=$(basename $0)
USAGE_EXAMPLE="$ME --service alfa --proto"

function usage {
    BOLD=$(tput bold)
    RESET=$(tput sgr0)
    cat <<-END
${BOLD}NAME${RESET}

    $ME --

${BOLD}SYNOPSIS${RESET}

    $ME [-h | --help] [-s | --service NAME] [-p | --proto]

${BOLD}DESCRIPTION${RESET}

    $ME .

${BOLD}OPTIONS${RESET}

    -s | --service NAME

        Microservice name.

    -p | --project

        Project path, e.g. "github.com/user/project".

    -h | --help

        Print this help message.


${BOLD}EXAMPLES${RESET}

    $USAGE_EXAMPLE

END
}

function trap_fn {
    echo $USAGE_EXAMPLE
}

trap trap_fn ERR SIGHUP SIGTERM SIGINT SIGFPE


TEMP=`getopt -o h,p:,s: --long help,project:,service: \
             -n "$ME" -- "$@"`
eval set -- "$TEMP"

while true; do
    case "$1" in

        -p | --project )
            PROJECT="$2"
            shift 2
            ;;

        -s | --service )
            SERVICE_NAME="${2,,}"
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

# script vars
BASE=$GOPATH/src/$PROJECT
DESTINATION_PATH=$BASE/src/$SERVICE_NAME
PROTOS_PATH=$DESTINATION_PATH/pb
CONFIG_FILE=$DESTINATION_PATH/conf/*.json

STUBS_PATH=$DESTINATION_PATH/pb
DEFINITIONS_PATH=$BASE/conf/

# generate configuration
CONFIG_FILE=$BASE/conf/global_conf.toml
$BASE/scripts/config.sh --out $CONFIG_FILE
CONFIG_DATA="$(cat $BASE/conf/global_conf.json)"


function load_config {

    PATTERN='del(.ConfigFile) | keys[] as $k | "\($k) \(.[$k] | .)"'
    DATA=$(jq -rc "$PATTERN" <<<$CONFIG_DATA)

    while read SERVICE
    do
        NAME="${SERVICE%% *}"
        DATA="${SERVICE#* }"
        CONFIG_MAP[${NAME,,}]="$DATA"
    done <<<"$DATA"
}

load_config

TYPE="service"

function is_gateway {
    jq -r '.Server.IsGateway?' <<<"${CONFIG_MAP[$1]}"
}

#is_gateway $SERVICE_NAME
IS_GATEWAY=$(is_gateway $SERVICE_NAME)
[[ $IS_GATEWAY == "true" ]] && {
    TYPE="gateway"
}

TEMPLATE_PATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd)/templates/$TYPE"
SERVER_FILE=${DESTINATION_PATH}/cmd/main.go
[[ -d ${DESTINATION_PATH}/cmd ]] && cp ${TEMPLATE_PATH}/${TYPE}.go $SERVER_FILE


echo -e \
"Generating:\t\t$SERVICE_NAME ($TYPE)\n
Using protos in:\t${PROTOS_PATH}
Using config in:\t${CONFIG_FILE}
Using template in:\t${TEMPLATE_PATH}
Files written into:\t${DESTINATION_PATH}\n"


mkdir -p $DESTINATION_PATH
mkdir -p $STUBS_PATH
mkdir -p $DEFINITIONS_PATH


# get all the proto files
PROTO_FILES=$(find ${PROTOS_PATH} -type f -name "*.proto")


function generate_stubs {
    echo "Generating gRPC stubs"

    for F in "${PROTO_FILES}"
    do
        protoc \
            -I${PROTOS_PATH} \
            -I${GOPATH}/src \
            -I${GOPATH}/src/github.com/golang/protobuf/ptypes/empty \
            -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
            --go_out=plugins=grpc:${STUBS_PATH} \
            $F
    done
            #-I${GOPATH}/src/github.com/google/protobuf/src \
}


function generate_reverse_proxy {

    echo "Generating gRPC reverse proxy"
    for F in "${PROTO_FILES}"
    do
        protoc \
            -I${PROTOS_PATH} \
            -I${GOPATH}/src \
            -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
            --grpc-gateway_out=logtostderr=true:${STUBS_PATH} \
	        $F
    done
}

function generate_swagger_definitions {

    echo "Generating swagger definitions"
    for F in "${PROTO_FILES}"
    do
        protoc \
            -I${PROTOS_PATH} \
            -I${GOPATH}/src \
            -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
            --swagger_out=logtostderr=true:${DEFINITIONS_PATH} \
            $F
    done
}

function add_services {

    declare -a IMPORTS
    declare -a ENDPOINTS
    declare -a SERVICES

    [[ $IS_GATEWAY != "true" ]] && {
        sed -i "s#\(get service configuration\)#\1\nservice = conf.${SERVICE_NAME^}#" $SERVER_FILE

        FN="pb.Register${SERVICE_NAME^}ServiceServer(srv, ${SERVICE_NAME}.NewService(conf))"
        sed -i "s#\(register the gRPC service server\)#\1\n$FN#" $SERVER_FILE

        IMPORT_PATH="${PROJECT}/src/${SERVICE_NAME}"
        IMPORTS+=("pb \"${IMPORT_PATH}/pb"\")
    }

    DATA="$(jq -r '.Server.Services[]?' <<<"${CONFIG_MAP[$SERVICE_NAME]}")"

    [[ "$DATA" != "" ]] && {

        SEP=

        echo -en "Registering services:\t"
        while read SERVICE
        do

            SERVICE=${SERVICE,,}
            echo -n "${SEP}${SERVICE}"
            SEP=", "

            # format the import
            IMPORT_PATH="${PROJECT}/src/${SERVICE}"
            PB_PACKAGE="${SERVICE}pb"
            IMPORTS+=("${PB_PACKAGE,,} \"${IMPORT_PATH}/pb"\")
            IMPORTS+=("${SERVICE} \"${IMPORT_PATH}"\")

            # format the endpoint
            ENDPOINT="Register${SERVICE^}ServiceHandlerFromEndpoint"

            LIS_ADDR="conf.${SERVICE^}"
            PARAMETERS="${LIS_ADDR}, ${PB_PACKAGE}.${ENDPOINT}"
            [[ $(is_gateway $SERVICE) == "true" ]] && {
                ENDPOINTS+=("gw.LoadEndpoint(ctx, ${PARAMETERS})")
            }
            SERVICES+=("${SERVICE}.Serve(ctx, conf)")

        done <<<"$DATA"

        echo
    }

    # import packages
    IMPORTS_DATA="$(printf '%s\\n' "${IMPORTS[@]}")"
    sed -i "s#\(Project packages\)#\1\n${IMPORTS_DATA}#" $SERVER_FILE

    # add the gateway endpoints
    ENDPOINTS_DATA="$(printf '%s\\n' "${ENDPOINTS[@]}")"
    sed -i "s#\(List of gateway endpoints\)#\1\n${ENDPOINTS_DATA}#" $SERVER_FILE

    # add the services
    SERVICES_DATA="$(printf '%s\\n' "${SERVICES[@]}")"
    sed -i "s#\(List of services\)#\1\n${SERVICES_DATA}#" $SERVER_FILE
}


function configure_grpc_server {

    [[ -e $SERVER_FILE ]] || {
        true
        return
    }

    echo "Configuring gRPC server"

    add_services

}


generate_stubs

configure_grpc_server

[[ $IS_GATEWAY == "true" ]] && {
    generate_reverse_proxy
    generate_swagger_definitions
}

goimports -w $SERVER_FILE

echo "Done"
