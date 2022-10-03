#!/bin/sh

sudo apt install git -y
wget https://golang.org/dl/go1.19.linux-armv6l.tar.gz
sudo tar -C /usr/local -xzf go1.19.linux-armv6l.tar.gz
sudo ln -s /usr/local/go/bin/go /usr/bin/go
git clone https://github.com/oliverschmidt/ProDOS-Utilities.git
cd ProDOS-Utilities || exit
go mod tidy
go build
cd ~ || exit
sudo ln -s $HOME/ProDOS-Utilities/ProDOS-Utilities /usr/bin/ProDOS-Utilities
git clone https://github.com/oliverschmidt/Apple2-IO-RPi.git
cd Apple2-IO-RPi/RaspberryPi/apple2driver || exit
go mod tidy
go build
sudo apt install cc65 vim -y
cd ~ || exit
sudo bash -c 'cat > /boot/config.txt << EOF
disable_splash=1
dtoverlay=disable-bt
boot_delay=0
EOF'
sudo bash -c 'echo " quiet" >> /boot/cmdline.txt'
wget hhttps://github.com/2a-4am/4cade/releases/latest/download/Total.Replay.hdv
sudo --preserve-env=HOME --preserve-env=USER bash -c 'cat > /etc/systemd/system/apple2driver.service << EOF
[Unit]
Description=Apple2-IO-RPi Driver

[Service]
ExecStart=$HOME/Apple2-IO-RPi/RaspberryPi/apple2driver/apple2driver -cdc -d1 $HOME/Apple2-IO-RPi/RaspberryPi/Apple2-IO-RPi.hdv -d2 $HOME/Total.Replay.hdv
SyslogIdentifier=apple2driver
User=$USER
Group=$USER
WorkingDirectory=$HOME/Apple2-IO-RPi/RaspberryPi/apple2driver

[Install]
WantedBy=basic.target
EOF'
sudo systemctl start apple2driver
sudo systemctl enable apple2driver
sudo systemctl disable avahi-daemon.service
sudo systemctl disable triggerhappy.service
sudo systemctl disable raspi-config.service
