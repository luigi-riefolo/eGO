#!/usr/bin/env bash

set -ex


# Maximise window
osascript <<END
tell application "System Events"
    set frontmostProcess to name of first process where it is frontmost
end tell

tell application "Finder"
    set screenSize to bounds of window of desktop
end tell

tell application frontmostProcess
    set bounds of the first window to screenSize
end tell
END


# Get sudo priviledges
sudo -v

# Create backup dir
mkdir ${HOME}/BKP || true

# Copy data
DATA_DST=${HOME}/Desktop/setup
echo "Copying data to $DATA_DST"
mkdir $DATA_DST || true
#cp -r * $DATA_DST;
cp -r Bash Bookmarks Conf Desktop Go Notes b.rtf github-recovery-codes.txt init.sh passwd.txt scan.pdf scripts syntax /Users/luigi/Desktop/setup


./scripts/setup-bashrc.sh

# Source bashrc
osascript <<END
tell application "System Events"
    tell application process "Terminal"
        set frontmost to true
        delay 5
        keystroke "source /etc/bashrc"
        keystroke return
        delay 5
		keystroke "./scripts/install-tools.sh"
        keystroke return
    end tell
end tell
END


# Drop sudo priviledges
#sudo -K
