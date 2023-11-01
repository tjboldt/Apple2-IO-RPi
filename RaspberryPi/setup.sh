#!/bin/sh

sudo apt update
sudo apt install git -y
if [ "$(uname -m)" = 'aarch64' ]; then
  if [ ! -f go1.21.3.linux-arm64.tar.gz ]; then  
    wget https://go.dev/dl/go1.21.3.linux-arm64.tar.gz
    sudo tar -C /usr/local -xzf go1.21.3.linux-arm64.tar.gz
  fi
else
  if [ ! -f go1.21.3.linux-arm6l.tar.gz ]; then  
    wget https://golang.org/dl/go1.21.3.linux-armv6l.tar.gz
    sudo tar -C /usr/local -xzf go1.21.3.linux-armv6l.tar.gz
  fi
fi
if [ -f "/usr/bin/go" ]; then
    sudo rm /usr/bin/go
fi
sudo ln -s /usr/local/go/bin/go /usr/bin/go
if [ ! -d "ProDOS-Utilities" ]; then
  git clone https://github.com/tjboldt/ProDOS-Utilities.git
fi
cd ProDOS-Utilities || exit
go build
cd ~ || exit
if [ -L "/usr/bin/ProDOS-Utilities" ]; then
    sudo rm /usr/bin/ProDOS-Utilities
fi
sudo ln -s /home/pi/ProDOS-Utilities/ProDOS-Utilities /usr/bin/ProDOS-Utilities
if [ ! -d "Apple2-IO-RPi" ]; then
  git clone https://github.com/tjboldt/Apple2-IO-RPi.git
fi
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
bash -c 'cat > apple2driver.service << EOF
[Unit]
Description=Apple2-IO-RPi Driver

[Service]
ExecStart=/home/$USER/Apple2-IO-RPi/RaspberryPi/apple2driver/apple2driver
StandardOutput=syslog
StandardError=syslog
SyslogIdentifier=apple2driver
User=$USER
Group=$USER
WorkingDirectory=/home/$USER/Apple2-IO-RPi/RaspberryPi/apple2driver

[Install]
WantedBy=basic.target
EOF'
sudo mv apple2driver.service /etc/systemd/system/apple2driver.service  
sudo chown root:root /etc/systemd/system/apple2driver.service  
sudo systemctl start apple2driver
sudo systemctl enable apple2driver
sudo systemctl disable avahi-daemon.service
sudo systemctl disable triggerhappy.service
sudo systemctl disable raspi-config.service
sudo systemctl daemon-reload
