## Setup starting from scratch
1. Have PCBs made from [PCBWay](https://www.pcbway.com/project/shareproject/Apple2_IO_RPi_v7_65457a66.html) or from the gerber and drill files in the Hardware folder or email me for a blank or fully assembled board
2. Solder chips, header and capacitors in place
3. If you have an EPROM programmer, it is preferred to pre-populate the EEPROM with the contents of AT28C68B.bin
4. Attach Raspberry Pi Zero 2 W facing outward from the card
5. Install [Raspberry Pi OS](https://www.raspberrypi.org/software/) on microSD card using a modern computer
    1. Use Raspberry Pi Imager
    2. Click "Choose OS"
    3. Select Other => "Raspberry Pi OS Lite (64 bit)" for the Zero 2 W
    4. Click on the gear icon in the bottom right
    5. Enable ssh
    6. Set password
    7. Configure wifi
    8. Set locale settings  
6. Put microSD card in the RPi
7. Install the expansion card into the Apple II
8. Power on the Apple II
9. Only if you didn't have an EPROM programmer in step 3, install firmware with utility (can be found on Apple2-IO-RPi.hdv drive image in the release on GitHub) 
10. Use `ssh` to connect to the RPi using the password you configured
11. `wget --no-cache -O - https://raw.githubusercontent.com/tjboldt/Apple2-IO-RPi/main/RaspberryPi/setup.sh | bash`

## Options
You can support two drives and change their drive images.
1. Modify the `ExecStart` line in `/etc/systemd/system/apple2driver.service` and make it something like the following: `ExecStart=/home/pi/Apple2-IO-RPi/RaspberryPi/apple2driver/apple2driver -d1 /home/pi/Apple2-IO-RPi/RaspberryPi/Apple2-IO-RPi.hdv -d2 /home/pi/Apple2-IO-RPi/RaspberryPi/TotalReplay401.hdv`
2. `sudo systemctl daemon-reload`
3. `sudo systemctl restart apple2driver.service`

