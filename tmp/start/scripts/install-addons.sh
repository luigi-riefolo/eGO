#!/usr/bin/env bash
# -------------------------------------------------------
#  Command line script to install
#  Firefox add-ons as global
#  extensions available to all users
#
# -------------------------------------------------------

set -xe

me="$(basename $0)"

if [[ "$EUID" -ne 0 ]]
then
    echo "The script '$me' needs sudo priviledges"
    exit 1
fi


# set global extensions installation path
# you may have to adapt it to your environment
PATH_FIREFOX="$(ls -d ${HOME}/Library/Application\ Support/Firefox/Profiles/*/extensions)"

declare -A ADDONS
ADDONS=(
    [Firebug]=1843
    [FireXpath]=11900
    [Adblock Plus]=1865
    [New Tab Homepage]=777
    [Flash Control]=495154 )


ADDONS_URL="https://addons.mozilla.org/firefox/downloads/latest"

function install_addons {
    for K in "${!ADDONS[@]}";do
        echo -e "\tInstalling: $K"
        ID=${ADDONS[$K]}

        # download extension
        wget -O addon.xpi "${ADDONS_URL}/${ID}/addon-${ID}-latest.xpi" || exit $?

        # get extension UID from install.rdf
        UID_ADDON=$(unzip -p addon.xpi install.rdf | grep "<em:id>" | head -n 1 | sed 's/^.*>\(.*\)<.*$/\1/g')

        # move extension to default installation path
        unzip addon.xpi -d "$PATH_FIREFOX"/$UID_ADDON
        rm addon.xpi

        # set root ownership
        chown -R $SUDO_USER:staff "$PATH_FIREFOX/$UID_ADDON"
        chmod -R a+rX "$PATH_FIREFOX/$UID_ADDON"

        # end message
        echo -e "\tAdd-on added under $PATH_FIREFOX/$UID_ADDON"
    done

    # Manually allow the plugins
    open -a Firefox
}

install_addons
