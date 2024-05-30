## Setup starting from scratch
1. DIY instructions and PCB coming but it is recommended to buy a board from Ralle Palaveev
2. Update the Pico firmware using Apple2-IO-RPi.uf2 from https://github.com/tjboldt/Apple2-IO-RPi/releases/latest
3. Attach Raspberry Pi Zero 2 W with a micro USB OTG cable
4. Install [Raspberry Pi OS](https://www.raspberrypi.org/software/) on microSD card using a modern computer
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
9. Use `ssh` to connect to the RPi using the password you configured
10. `wget --no-cache -O - https://raw.githubusercontent.com/tjboldt/Apple2-IO-RPi/main/RaspberryPi/setup.sh 1 | bash`

## Options
You can support two drives and change their drive images.
1. Modify the `ExecStart` line in `/etc/systemd/system/apple2driver.service` and make it something like the following: `ExecStart=/home/pi/Apple2-IO-RPi/RaspberryPi/apple2driver/apple2driver -cdc=true -d1 /home/pi/Apple2-IO-RPi/RaspberryPi/Apple2-IO-RPi.hdv -d2 /home/pi/Apple2-IO-RPi/RaspberryPi/TotalReplay401.hdv`
2. `sudo systemctl daemon-reload`
3. `sudo systemctl restart apple2driver.service`
