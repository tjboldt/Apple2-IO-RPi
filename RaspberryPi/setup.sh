#!/bin/sh

sudo apt install git -y
wget https://golang.org/dl/go1.17.3.linux-armv6l.tar.gz
sudo tar -C /usr/local -xzf go1.17.3.linux-armv6l.tar.gz
sudo ln -s /usr/local/go/bin/go /usr/bin/go
git clone https://github.com/tjboldt/ProDOS-Utilities.git
cd ProDOS-Utilities || exit
go build
cd ~ || exit
sudo ln -s /home/pi/ProDOS-Utilities/ProDOS-Utilities /usr/bin/ProDOS-Utilities
git clone https://github.com/tjboldt/Apple2-IO-RPi.git
cd Apple2-IO-RPi/RaspberryPi/apple2driver || exit
go build
sudo apt install cc65 vim -y
cd ~ || exit
sudo bash -c 'cat > /boot/config.txt << EOF
disable_splash=1
dtoverlay=disable-bt
boot_delay=0
EOF'
sudo bash -c 'echo " quiet" >> /boot/cmdline.txt'
sudo bash -c 'cat > /etc/systemd/system/apple2driver.service << EOF
[Unit]
Description=Apple2-IO-RPi Driver

[Service]
ExecStart=/home/pi/Apple2-IO-RPi/RaspberryPi/apple2driver/apple2driver
StandardOutput=syslog
StandardError=syslog
SyslogIdentifier=apple2driver
User=pi
Group=pi
WorkingDirectory=/home/pi/Apple2-IO-RPi/RaspberryPi/apple2driver

[Install]
WantedBy=basic.target
EOF'
sudo systemctl start apple2driver
sudo systemctl enable apple2driver
sudo systemctl disable avahi-daemon.service
sudo systemctl disable triggerhappy.service
sudo systemctl disable raspi-config.service
