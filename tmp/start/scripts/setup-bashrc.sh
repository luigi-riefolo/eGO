#!/usr/bin/env bash

set -ex

echo "Setting bashrc"

#sudo -- sh -c 'echo "alias zzzzzpotatoxxx=\"echo YES POTATOOOO\"" >> /etc/bashrc'

# Backups
cp /etc/bashrc ${HOME}/BKP/
cp ~/.bashrc ${HOME}/BKP/ || true

sudo cp ${HOME}/Desktop/setup/Conf/bashrc /etc/

echo "Bashrc setting done"
