#!/usr/bin/env bash

set -xe

VIRTUAL_BOX=${HOME}/VirtualBox\ VMs

DEPLOY_VM=
DOWNLOAD_VM=
HD_SIZE=
ISO=
ISO_FILE=false
OS_TYPE=
RAM_SIZE=
REPLACE=
VM_NAME=

URL="http://fedora.uib.no/fedora/linux/releases/25/Workstation/x86_64/iso/"
VM_URL_PATTERN="Fedora-Workstation-Live-"



function download_vm() {
    echo -e "Downloading $OS_TYPE VM"
    wget -e robots=off -nd -r -l1 --no-parent -P ~/Downloads/ -A "${VM_URL_PATTERN}*.iso" "$URL"
    ISO=$(ls ~/Downloads/*.iso)
}

function invalid_types() {
    echo "Invalid VM type"
    echo "List of valid VM types:"
    VBoxManage list ostypes | ggrep -Po "^ID:\s+\K(.*)" | sed ':a;N;$!ba;s/\n/, /g'
}

function check_os_type() {
    if ! VBoxManage list ostypes | ggrep -Po "^ID:\s+\K(.*)" | grep -wq "$1"
    then
        echo "Invalid OS type"
        invalid_types
        exit 1
    fi
}


function install_vm() {

    # Create VM
    echo -e "Creating VM"

    VBoxManage createvm --name "$VM_NAME" --register

    # VM settings
    echo -e "Setting VM"
    VBoxManage modifyvm "$VM_NAME" --ostype "$OS_TYPE" --memory $RAM_SIZE \
        --vram 128 --ioapic on --boot1 dvd --boot2 disk --nic1 bridged --bridgeadapter1 en1

    VBoxManage createhd --filename "${VIRTUAL_BOX}/${VM_NAME}/${VM_NAME}.vdi" --size ${HD_SIZE}

    VBoxManage storagectl "$VM_NAME" --name "SATA Controller" --add sata --controller IntelAHCI
    VBoxManage storageattach "$VM_NAME" --storagectl "SATA Controller" \
        --port 0 --device 0 --type hdd --medium "${VIRTUAL_BOX}/${VM_NAME}/${VM_NAME}".vdi

    VBoxManage storagectl "$VM_NAME" --name "IDE Controller" --add ide
    VBoxManage storageattach "$VM_NAME" --storagectl "IDE Controller" \
            --port 0 --device 0 --type dvddrive --medium "${ISO}"

    echo -e "Booting VM"
    #VBoxHeadless --startvm "$VM_NAME" &
    VBoxManage startvm "$VM_NAME"

    #ssh -L 3389:127.0.0.1:3389 <host>
}

#function install_guest_additions() {
#    echo -e "Taking snapshot"
#    VBoxManage snapshot "$VM_NAME" take "$(date +%d-%m-%y_%H:%M)"
#
#    VBoxManage modifyvm "$VM_NAME" --dvd none
#
#    VBoxManage storageattach "$VM_NAME" --storagectl "IDE Controller" \
#        --port 0 --device 0 --type hdd \
#        --medium /Applications/VirtualBox.app/Contents/MacOS/VBoxGuestAdditions.iso
#
##    mkdir /mnt/dvd
##    mount -t iso9660 -o ro /dev/dvd /mnt/dvd
##    cd /mnt/dvd
##    ./VBoxLinuxAdditions.run
#
#    VBoxManage controlvm clipboard bidirectional
#    VBoxManage controlvm draganddrop hosttoguest
#
#    vboxmanage sharedfolder add "io" --name share-name \
#        --hostpath /path/to/folder/ --automount
#    vboxmanage sharedfolder remove "io" --name share-name
#
#     # To mount it on the guest,
#     #sudo mount -t vboxsf -o uid=$UID share-name /path/to/folder/share/
#}

function remove_vm() {
    if VBoxManage list vms | grep -wq "$VM_NAME"
    then
        echo "Removing existing VM $VM_NAME"
        VBoxManage unregistervm "$VM_NAME" --delete
    fi
}


function deploy_vm() {
    echo "Deploying VM '$VM_NAME'"
    VBoxManage import "$DEPLOY_VM"
    VBoxManage startvm "$VM_NAME"

    VBoxManage sharedfolder add "$VM_NAME" --name share-name --hostpath ${HOME}/Workspace/ --automount
}


function usage {
	me=$(basename $0)
    bold=$(tput bold)
    underline=$(tput smul)
    reset=$(tput sgr0)
    cat <<-END
${bold}NAME${reset}

    $me -- create a VM

${bold}SYNOPSIS${reset}

    $me [-d | --dowload-iso] [-i | --iso ISO] [-o | --os OS_TYPE] [-s | --hd BYTES] [-r | --ram BYTES]

${bold}DESCRIPTION${reset}

    $me creates a VM. The VM can be downloaded or supplied with a file.

END
}


# Arguments
if [[ $# == 0 ]]
then
    echo "Please supply a valid of arguments"
	usage
    exit 1
fi

if [[ $? != 0 ]]
then
    echo "Invalid arguments" >&2
    exit 1
fi

TEMP=`getopt -o di:n:o:s:r: --long deploy:,download-iso,iso:,name:,os:,hd:,ram:,replace \
             -n 'install-vm.sh' -- "$@"`

# Note the quotes around `$TEMP': they are essential!
eval set -- "$TEMP"

while true; do
    case "$1" in
        -d | --download-iso )
            DOWNLOAD_ISO=true
            shift
            ;;
        --deploy )
            DEPLOY_VM="$2"
            shift 2
            ;;
        -i | --iso )
            [[ -f $2 ]] || {
                echo "ISO file does not exist: $2"
                exit 1
            }
            file $2 | grep -q ISO || {
                echo "Please supply a valid ISO"
                exit 1
            }
            ISO=$2
            ISO_FILE=true
            shift 2
            ;;
        -n | --name )
            VM_NAME="$2"
            shift 2
            ;;
        -o | --os )
            OS_TYPE="$2"
            check_os_type $OS_TYPE
            shift 2
            ;;
        -s | --hd )
            HD_SIZE=$2
            shift 2
            ;;
        -r | --ram )
            RAM_SIZE=$2
            shift 2
            ;;
        --replace )
            REPLACE=true
            shift
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


# Set default values
# 32 GB
[[ -n $HD_SIZE ]] || HD_SIZE=32768
# 2 GB
[[ -n $RAM_SIZE ]] || RAM_SIZE=4096
[[ -n $OS_TYPE ]] || OS_TYPE="Fedora_64"
[[ -n $VM_NAME ]] || VM_NAME="$OS_TYPE"

if [[ "${DOWNLOAD_ISO}" = true &&"${ISO_FILE}" = true ]]
then
    echo "Requested ISO download, but only ISO file will be used"
    DOWNLOAD_ISO=false
fi

#if [[ "${DOWNLOAD_ISO}" = true ]]
#then
#    download_vm
#fi

if [[ "${REPLACE}" = true ]]
then
    remove_vm
else
    if [[ $(VBoxManage list vms | grep -qw "$VM_NAME") ]]
    then
        echo "VM '$VM_NAME' alreardy exist."
        echo "If you want to replace it use --replace."
        exit 1
    fi
fi



#########################################
#                                       #
#               Main                    #
#                                       #
#########################################

# Setup the VM installing GuestAdditions
[[ "$SETUP_VM" = true ]] && {
    install_guest_additions || exit $?
}

# Deploy VM
[[ -n "$DEPLOY_VM" ]] && {
    deploy_VM || exit $?
    exit 0
}

# Install VM (default)
install_vm

