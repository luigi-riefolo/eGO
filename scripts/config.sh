#!/usr/bin/env bash

# Script for the dynamic generation of service config files.


set -e

BASE=${GOPATH}/src/github.com/luigi-riefolo/alfa
ME=$(basename $0)
OUT_FILE=/dev/stdout

function usage {
    BOLD=$(tput bold)
    RESET=$(tput sgr0)
    cat <<-END
${BOLD}NAME${RESET}

    $ME -- genates a TOML global config file.

${BOLD}SYNOPSIS${RESET}

    $ME [-h | --help] [-o | --out FILE]

${BOLD}DESCRIPTION${RESET}

    $ME dynamically generates a TOML global config file.

${BOLD}OPTIONS${RESET}

    -o | --out

        Output file. If not specified STDOUT is used.

    -h | --help

        Print this help message.


${BOLD}EXAMPLES${RESET}

    ./$ME --out conf/global_conf.toml

END
}

TEMP=`getopt -o h,o: --long help,out: \
             -n "$ME" -- "$@"`
eval set -- "$TEMP"

while true; do
    case "$1" in

        -o | --out )
            OUT_FILE=$2
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


HEADER="# Dynamically generated file, please do not manually edit it.\n"


function global_conf {

    echo -e $HEADER > $OUT_FILE

    for F in $(ls $BASE/src)
    do
        cat $BASE/src/$F/conf/config.toml >> $OUT_FILE
        echo -e "\n" >> $OUT_FILE
    done

    go run $BASE/pkg/config/cmd/main.go config \
        -config $BASE/conf/global_conf.toml  | json_pp > $BASE/conf/global_conf.json
}

global_conf
