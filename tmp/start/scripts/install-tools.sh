#!/usr/bin/env bash

set -xe


# Check if the user is an admin
if id -Gn $USER | grep -vq -w admin;
then
    echo "User '$USER' has to be and admin";
    exit 1
fi


AUTHOR='Luigi Riefolo'
AUTHORREF='LR'
COMPANY='Daimler Mobility Services'
EMAIL='luigi.riefolo@daimler.com'
SSH_KEY='mfh%499a-1.a/an1FQbxjan1$3zn!bEa'
BKP_DIR=${HOME}/BKP/

cd ${HOME}/Desktop/setup


# Enable encryption
function enable_filevault() {
    sudo fdesetup enable -user $USER -verbose -authrestart
}


# Enable firewall
function enable_firewall() {
    /usr/libexec/ApplicationFirewall/socketfilterfw --setstealthmode on
    /usr/libexec/ApplicationFirewall/socketfilterfw --setloggingmode on
    /usr/libexec/ApplicationFirewall/socketfilterfw --setallowsigned on
    /usr/libexec/ApplicationFirewall/socketfilterfw --setglobalstate on
    /usr/libexec/ApplicationFirewall/socketfilterfw --getglobalstate
}


# Xcode
function install_xcode() {
    echo "Installing Xcode"
    xcode-select --install

    # Check version
    xcode-select -p
    gcc --version
}


# Ruby version manager
function set_ruby() {
    echo "Installing Ruby Version Manager"
    \curl -sSL https://get.rvm.io | bash -s stable --ruby
    source ~/.rvm/scripts/rvm

    # TODO: open xcode and accept the license agreement
}


# Brew
function set_brew() {
#    echo "Installing Brew"
#    sudo gem update --system
#    ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
#
#    brew update
#    brew tap homebrew/dupes
#    brew tap caskroom/cask
#    brew tap phinze/homebrew-cask
#    brew doctor || true
#
#    # bash v4 needs to be adde to approved shells and use chsh
#    echo -e "\tInstalling Brew packages"
#    brew install bash jq gdb python pyenv pyenv-virtualenv gpatch less m4 make \
#        nano binutils diffutils nmap node wget aws-shell dos2unix pv awscli \
#        hh tmux unrar tree links git htop gawk bash-completion ack colordiff \
#        coreutils watch automake autoconf python3 gnutls screen \
#        gzip openssh openssl perl svn rsync unzip brew-cask-completion ipcalc \
#        moreutils ssh-copy-id rbenv ruby-build rbenv-default-gems \
#        postgresql go lynx calc gnu-getopt valgrind \
#        gtk
#
#    # Delve
#    #tar xvpf delve-1.0.0-rc.1.tar.gz
#    #./scripts/gencert.sh
#    # brew install go-delve/delve/delve
#
#    brew install ed --with-default-names
#    brew install findutils --with-default-names
#    brew install gnu-indent --with-default-names
#    brew install gnu-sed --with-default-names
#    brew install gnu-tar --with-default-names
#    brew install gnu-which --with-default-names
#    brew install grep --with-default-names
#    brew install wdiff --with-gettext
#    brew install vim --override-system-vi --with-lua
    #brew install macvim --override-system-vim --custom-system-icons
    brew install macvim  --custom-system-icons

    # Applications
    echo -e "\tInstalling Brew applications"
    #brew cask install cmake firefox caffeine osxfuse google-chrome google-backup-and-sync slack virtualbox java atom burp-suite shiftit libreoffice libreoffice-language-pack chromium chrome-devtools macfusion cpuinfo docker docker-toolbox google-hangouts wireshark netbeans virtualbox-extension-pack
    brew cask install vlc meld

    # Relink
    echo -e "\tRelinking Brew binaries"
    brew list -1 | while read line;do
        brew unlink $line
        brew link --force --overwrite $line
    done

    # NOTE: do we need this???
    brew link --force gnu-getopt

    # Clean
    echo -e "\tCleaning Brew cache"
    rm ${HOME}/Library/Caches/Homebrew/Cask/*.dmg

    echo -e "\tUpdating Brew"
    brew update && brew upgrade

    brew missing &> ${HOME}/Desktop/brew-missing.txt
    brew ource /etc/bashrc
uoctor &> ${HOME}/Desktop/brew-doctor.txt
}


# Python
function set_python() {
    echo "Setting Python"
    sudo -H pip install --upgrade pip
    sudo -H pip install numpy requests flake8 flake8-docstrings pep8 scrapy nose selenium

    # Upgrade packages
    echo -e "\tUpgrading Python packages"
    pip freeze --local | grep -v '^\-e' | cut -d = -f 1  | xargs -n1 sudo -H pip install -U

    # setuptools downgrades pip
    sudo -H pip install --upgrade pip
}


# Atom packages
function set_atom() {
    echo "Installing Atom packages"
    apm install atom-beautify auto-detect-indentation autocomplete-go autocomplete-java autocomplete-modules autocomplete-python autocomplete-ruby builder-go environment go-config go-debug go-get go-plus godoc gofmt gometalinter-linter gorename java-generator java-importer linter linter-flake8 linter-javac linter-pep8 linter-ruby navigator-go pandoc platformio-ide-terminal python-indent python-tools ruby-slim sort-lines terminal-plus tester-go whitespace
}


# Vim
function set_vim() {
    echo "Setting VIM"
    cp ${HOME}/.vimrc $BKP_DIR || true

    # Create the .vimrc file
	echo -e "\tCreating the .vimrc file"
    echo -e "set nocompatible\nset runtimepath=~/.vim,$(find /usr/local/Cellar/vim/ -type d -regex ".*vim[0-9]+")\n" > ~/.vimrc
    cat ${HOME}/Desktop/setup/Conf/vimrc >> ${HOME}/.vimrc

    # Install the plugins
	echo -e "\t Installing the plugins"
    git clone https://github.com/VundleVim/Vundle.vim.git ${HOME}/.vim/bundle/Vundle.vim
    vim +PluginInstall +qall
    mkdir -p ${HOME}/.vim/autoload ~/.vim/bundle && curl -LSso ${HOME}/.vim/autoload/pathogen.vim https://tpo.pe/pathogen.vim

    # Fix the bash-support help command
	echo -e "\tSetting bash-support"
    sed -i 's/help -m/help /g' ${HOME}/.vim/bundle/bash-support.vim/plugin/bash-support.vim

	VIM_TEMPLATE_DIR=${HOME}/.vim/bundle/bash-support.vim/bash-support/templates

    declare -A bash_settings
    bash_settings=(
        ['AUTHOR']=$AUTHOR
        ['AUTHORREF']=$AUTHORREF
        ['COMPANY']=$COMPANY
        ['EMAIL']=$EMAIL)

    for K in "${!bash_settings[@]}";do
        V="${bash_settings[$K]}"
        sed -ie "s/'$K',\(\s*\)''/'$K',\1'$V'/" ${VIM_TEMPLATE_DIR}/Templates
    done

    sed -i "s/bin\/bash -/usr\/bin\/env bash/" ${VIM_TEMPLATE_DIR}/comments.templates

    cp -r ${HOME}/Desktop/setup/syntax ${HOME}/.vim/
}



# Github
function set_github() {

# Setup github account
# open -a firefox 'https://github.com/'
# Github Personal access tokens
# open -a firefox 'https://github.com/settings/tokens'

    echo "Setting Github"
    git config --global user.name "$AUTHOR"
    git config --global user.email "$EMAIL"

    # Generating SSH key
    ssh-keygen -f "${HOME}/.ssh/id_rsa" -t rsa -b 4096 -C "$EMAIL" -N "$SSH_KEY"

#    echo "Checking SSH agent"
#    eval "$(ssh-agent -s)"
#    ssh-add ~/.ssh/id_rsa
#
#    echo "Testing SSH key"
#    ssh -vT git@github.com

    # TODO set key!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
#    git_api_token="$CWD$CWDXXXXXXXXXX"
#
#    # We'll use the HTTPS to push a ssh key to git, SSH for pull/push configuration
#    gitrepo_ssh="git@github.com:$GIT_NAME/repo.git"
#    gitrepo_https="https://github.com/$GIT_NAME/repo.git"
#
#    sslpub="$(cat ${HOME}/.ssh/id_rsa.pub |tail -1)"
#
#    # git API path for posting a new ssh-key:
#    git_api_addkey="https://api.github.com/user/keys"
#
#    git_ssl_keyname="$(hostname)_simplesurance"
#
#    # Post this ssh key:
#    curl -H "Authorization: token ${git_api_token}" -H "Content-Type: application/json" -X POST -d "{\"title\":\"${git_ssl_keyname}\",\"key\":\"${sslpub}\"}" ${git_api_addkey}
#
#    # git remote set-url origin git@github.com-user1:user1/your-repo-name.git
}


# Workspace
function set_workspace() {
#    echo "Setting Workspace"
#    mkdir ${HOME}/Workspace
#    cp -R ${HOME}/Desktop/setup/Desktop/* ${HOME}/Desktop/
#
#    # Go
#    CWD=$(pwd)
#    cp -R ${HOME}/Desktop/setup/Go ${HOME}/Workspace/
#    git clone https://go.googlesource.com/go ${HOME}/Workspace/Go/source/go
#    cd ${HOME}/Workspace/Go/source/go
#    git checkout master
#    cd src
#    ./all.bash
#    cd $CWD

#    WGET_URL="https://cloud.google.com/appengine/docs/go/download"
#    DOWNLOAD_PATTERN="href=\"\K(https://storage.googleapis.com/appengine-sdks/featured/go_appengine_sdk_darwin_amd64.*zip)"
#    wget -P ${HOME}/Downloads/ $(curl -s --insecure $WGET_URL | grep -Po "$DOWNLOAD_PATTERN" )
#    unzip -q -d ${HOME}/Workspace/Go/source/ ${HOME}/Downloads/go_appengine_sdk_darwin_amd64*.zip
#    rm -f ${HOME}/Downloads/go_appengine_sdk_darwin_amd64.*zip

#    mkdir -p ${HOME}/Workspace/Go/examples/go_web_programming
#    git clone https://github.com/sausheong/gwp ${HOME}/Workspace/Go/examples/go_web_programming/gwp

    GOROOT=
#
#    # Go packages
#    GO_PACKAGES=( \
#        golang.org/x/tools/cmd/... \
#        github.com/golang/lint/golint \
#        github.com/constabulary/gb/... \
#        github.com/fatih/motion \
#        github.com/gorilla/mux \
#        github.com/gorilla/context \
#        github.com/gorilla/securecookie \
#        github.com/gorilla/sessions \
#        github.com/gorilla/rpc \
#        github.com/gorilla/websocket \
#        github.com/josharian/impl \
#        github.com/jstemmer/gotags \
#        github.com/kisielk/errcheck \
#        github.com/klauspost/asmfmt/cmd/asmfmt \
#        github.com/lib/pq \
#        github.com/nsf/gocode \
#        github.com/pkg/errors \
#        github.com/rogpeppe/godef \
#        golang.org/x/tools/cmd/guru \
#        github.com/tools/godep \
#        github.com/zmb3/gogetdoc \
#        google.golang.org/grpc \
#        github.com/golang/protobuf/{proto,protoc-gen-go}
#        github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
#        github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger \
#        google.golang.org/grpc/status \
#        github.com/grpc-ecosystem/go-grpc-middleware/... \
#        gopkg.in/alecthomas/kingpin.v2 \
#        github.com/go-sql-driver/mysql \
#        github.com/alecthomas/gometalinter \
#        github.com/BurntSushi/toml \
#        golang.org/x/time/rate \
        github.com/BurntSushi/toml
#        github.com/BurntSushi/toml/cmd/tomlv )
#    for PKG in "${GO_PACKAGES[@]}"
#    do
#        go get $PKG
#    done
#
#    # gRPC and ProtoBuffers
#    cd ${HOME}/Workspace/Go
#    git clone https://github.com/google/protobuf
#    cd protobuf
#    brew unlink libtool
#    brew link --overwrite libtool
#    ./autogen.sh
#    ./configure
#    make
#    make check
#    sudo make install
#    which protoc
#    protoc --version
#    cd $CWD

    #gometalinter --install

    cp -R ${HOME}/Desktop/setup/Python ${HOME}/Workspace
    cp -R ${HOME}/Desktop/setup/Bash ${HOME}/Workspace
    cp -R ${HOME}/Desktop/setup/Notes ${HOME}/Workspace

    GIT_HUB_ACCOUNT="https://github.com/luigi-riefolo"
    git clone $GIT_HUB_ACCOUNT/network_crawler ${HOME}/Workspace/Python/network_crawler
    git clone $GIT_HUB_ACCOUNT/Address_book ${HOME}/Workspace/Python/Address_book
    git clone $GIT_HUB_ACCOUNT/Web_Crawler ${HOME}/Workspace/Python/Web_Crawler
}


# Set Firefox
function set_firefox() {
    echo "Setting Firefox"

    # Addons
    sudo ${HOME}/Desktop/setup/scripts/install-addons.sh
}



# Main
#enable_filevault
#enable_firewall

#install_xcode

#set_ruby
#set_brew

#set_python
#set_atom
#set_vim

#set_github
#set_workspace

set_firefox

#${HOME}/Desktop/setup/scripts/install-vm.sh --deploy ${HOME}/Desktop/setup/Fedora-64.ova --name "Fedora-64"
