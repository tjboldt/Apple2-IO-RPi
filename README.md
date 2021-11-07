# Apple2-IO-RPi
Apple II expansion card using a Raspberry Pi for I/O

![Image of Board](/Hardware/Apple2IORPi.jpg)

## Purpose
The purpose of this project is to provide I/O for an Apple II series 8 bit computer via a Raspberry Pi Zero W which is powered by the Apple II expansion bus. This includes using the attached RPi Zero W for it's storage, network and processor to provide new functionality for the Apple II.

## Project Status
So far, this is a project and not a finished product. The current prototype is on the fifth revision and a few have been assembled and tested. It is now possible for the Apple II to boot from and write to a virutal hard drive image stored on the RPi in any slot and execute simple commands on the RPi via the Apple II. The code has very few tests and is incomplete. Note that currently the firmware assumes an 80 column card is in slot 3 and than you have lowercase support. If you have a problem or idea for enhancement, log an issue [here](https://github.com/tjboldt/Apple2-IO-RPi/issues). I recommend starring/watching the project for updates on GitHub. You are welcome to fork the project and submit pull requests which I will review.

## Features
1. Boot menu which waits for RPi to be ready
2. ProDOS bootable drive from image stored on RPi
3. Execute Linux commands on the RPi from the Apple II
4. Load binary files directly from the RPi to the II
5. Update Apple II firmware in place from image on RPi
6. Supports two drive images at the same time (Note: backward compatible with previous firmware but requires firmware update in order to work with two drives)

## Roadmap
1. Extend BASIC.SYSTEM commands
    1. RPI - Execute a single Linux command
    2. SH - Open a Linux shell
    3. WGET - Download and save to file 
2. Proper ProDOS clock driver (currently just directly sets values on block reads)
3. Bi-directional image conversion between common formats and HIRES
4. Full terminal emulation
5. Remote code execution 
6. Proxy VNC connection, rendering as Apple II compatible graphics

## Setup starting from scratch
1. Have PCBs made from the gerber and drill files in the Hardware folder or email me for a blank or fully assembled board
2. Solder chips, header and capacitors in place
3. Attach Raspberry Pi Zero W facing outward from the card
4. Install [Raspberry Pi OS](https://www.raspberrypi.org/software/) on microSD card
5. Configure wifi in [boot/wpa_supplicant.conf](https://howtoraspberrypi.com/how-to-raspberry-pi-headless-setup/)
6. Add empty ssh file boot (for ssh access over wifi)
7. Put microSD card in the RPi
8. Install the expansion card into the Apple II
9. Power on the Apple II
10. Update firmware with utility (can be found on Apple2-IO-RPi.hdv drive image) or use EPROM programmer
11. Install Golang (note that `sudo apt install golang` will install an older version that doesn't work with this code)
    1. `sudo apt install git`
    2. `wget https://golang.org/dl/go1.17.3.linux-armv6l.tar.gz`
    3. `sudo tar -C /usr/local -xzf go1.17.3.linux-armv6l.tar.gz`
    4. Edit `~/.profile` with your favourite editor and add the following:
        1. `PATH=$PATH:/usr/local/go/bin`
        2. `GOPATH=$HOME/golang`
12. `git clone https://github.com/tjboldt/ProDOS-Utilities.git`
13. `cd ProDOS-Utilities`
14. `go build`
15. `cd ~`
16. `git clone https://github.com/tjboldt/Apple2-IO-RPi.git`
17. `cd Apple2-IO-RPi/RaspberryPi/apple2driver`
18. `go build`
19. `./apple2driver`
20. Optional to step above, `./apple2driver -d1 YOUR_DRIVE.hdv` (Apple2-IO-RPi.hdv is automatically selected as drive 2), or `./apple2driver -d1 YOUR_DRIVE.hdv -d2 YOUR_SECOND_DRIVE.hdv`
21. Setup the Driver as a service or to autostart via cronjob (`crontab -e` then add the line `@reboot /home/pi/Apple2-IO-RPi/RaspberryPi/apple2driver/apple2driver -d1 YOUR_DRIVE.hdv -d2 /home/pi/Apple2-IO-RPi/RaspberryPi/Apple2-IO-RPi.hdv > /home/pi/Apple2-IO-RPi/RaspberryPi/Apple2-IO-RPi.log`)
22. If you want to assemble 6502 code on the RPi:
    1. Add to your ~/.profile `PATH=$PATH:/home/pi/ProDOS-Utilities`
    2. Install cc65 with `sudo apt install cc65`

## Setup if you received a complete board from me
1. Put in any slot (slot 7 preferred as it is the first to boot)
2. Turn on your Apple II
3. Wait for the RPi to start up (will show ... until it connects)
4. Select `2` from the menu to run a command
5. Type `a2wifi list` to check that you have a wireless network in range
6. Type `a2wifi select YOUR_SSID YOUR_PASSWORD` to connect the RPi to your network
7. Once connected, you should be able to ssh to the RPi from any computer if you want. From a Mac or Linux (or Windows with Bash) shell, type `ssh pi@raspberrypi`. The default password is `raspberry`. 
8. Rebooting the Apple II, you can then select `1` from the menu and it will boot to ProDOS 2.4.2 and you will have a mostly blank 32 MB drive with a couple of utilities for updating firmware and such.
9. Star and Watch this repo on GitHub for the latest updates

## Update
1. Restart Apple II
2. Select run command
3. `cd /home/pi/Apple2-IO-RPi`
4. `git pull`
5. `cd RaspberryPi/apple2driver`
6. `go build`
7. `sudo reboot`
8. Restart Apple II
9. Select boot
10. `-UPDATE.FIRMWARE`
11. Enter slot number the card is in
12. Wait for firmware update to complete all four pages

### Additional steps to upgrade Golang and use new command line parameters for the service (if you had set up before October 11, 2021)
This must be done via ssh directly into the RPi:
1. `sudo apt remove golang`
2. `wget https://golang.org/dl/go1.17.3.linux-armv6l.tar.gz`
3. `sudo tar -C /usr/local -xzf go1.17.3.linux-armv6l.tar.gz`
4. Edit `~/.profile` with your favourite editor and add the following:
    1. `PATH=$PATH:/usr/local/go/bin`
    2. `GOPATH=$HOME/golang`
5. `cd ~/`
6. `git clone https://github.com/tjboldt/ProDOS-Utilities.git`
7. `cd ProDOS-Utilities`
8. `cd ~/Apple2-IO-RPi/RaspberryPi/apple2driver`
9. `go build`
10. Edit the Driver autostart via cronjob (`crontab -e` then edit the line to be `@reboot /home/pi/Apple2-IO-RPi/RaspberryPi/apple2driver/apple2driver -d1 YOUR_DRIVE.hdv -d2 /home/pi/Apple2-IO-RPi/RaspberryPi/Apple2-IO-RPi.hdv > /home/pi/Apple2-IO-RPi/RaspberryPi/Apple2-IO-RPi.log`) or simply `@reboot /home/pi/Apple2-IO-RPi/RaspberryPi/apple2driver/apple2driver > /home/pi/Apple2-IO-RPi/RaspberryPi/Apple2-IO-RPi.log` if you don't need your own drive image.

## Similar Project
If you prefer having Apple II peripherals control a Raspberry Pi rather than simply using the Raspberry Pi to provide storage, network access and processing to the Apple II, have a look at David Schmenk's excellent [Apple2Pi](https://github.com/dschmenk/apple2pi) project. 
