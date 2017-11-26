#!/usr/bin/env bash

# Script for setting up local networking

sudo -v

HOST_FILE=/etc/hosts
HOST_IP='127.0.0.1'

BANNER='# Project host files'
FOOTER='####################'

TMP_FILE=$(mktemp)
BKP_FILE=$(mktemp)

function add_server_hostnames {

    cp -p $HOST_FILE "$BKP_FILE"
    echo -e "Hosts file backup:\t$BKP_FILE"

    cp $HOST_FILE "$TMP_FILE"

    TMP_DATA="$BANNER\n"

    for S in ${SERVERS}
    do
        LINE="$HOST_IP\t$S\n"
        TMP_DATA+="$LINE"
    done

    TMP_DATA+="$FOOTER"

    grep -qF "$BANNER" $HOST_FILE ||
        echo -e "\n${BANNER}\n${FOOTER}" >> "$TMP_FILE"

    sed -ir "/$BANNER/,/$FOOTER/c\\$TMP_DATA" "$TMP_FILE"

    sudo cp "$TMP_FILE" $HOST_FILE

    rm "$TMP_FILE"
}


add_server_hostnames
