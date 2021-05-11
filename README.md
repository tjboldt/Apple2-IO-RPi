# Apple2-IO-RPi
Apple II expansion card using a Raspberry Pi for I/O

![Image of Board](/Hardware/Apple2IORPi.jpg)

## Purpose
The purpose of this project is to provide I/O for an Apple II series 8 bit computer via a Raspberry Pi Zero W which is powered by the Apple II expansion bus. This includes using the attached RPi Zero W for it's storage, network and processor to provide new functionality for the Apple II.

## Project Status
So far, this is a project and not a finished product. The current prototype is on the fourth revision and a few have been assembled and tested. It is now possible for the Apple II to boot from and write to a virutal hard drive image stored on the RPi in any slot and execute simple commands on the RPi via the Apple II. The code has no error handling or tests yet and is incomplete.

## Roadmap
1. Detect when RPi is in a ready state
2. Wifi setup tool
3. Proper ProDOS clock driver (currently just directly sets values on block reads)
4. Support for direct file read/write without drive image
5. Image conversion on download
6. Web service call support
7. Full terminal emulation
8. Remote code execution 
9. Proxy VNC connection, rendering as Apple II compatible graphics

## Setup
1. Have PCBs made from the gerber and drill files in the Hardware folder or email me for a blank or fully assembled board
2. Solder chips, header and capacitors in place
3. Attach Raspberry Pi Zero W facing outward from the card
4. Install Raspberry Pi OS on microSD card https://www.raspberrypi.org/software/
5. Configure wifi in boot/wpa_supplicant.conf https://howtoraspberrypi.com/how-to-raspberry-pi-headless-setup/
6. Add empty ssh file boot (for ssh access over wifi)
7. Put microSD card in the RPi
8. Install the expansion card into the Apple II
9. Power on the Apple II
10. Update firmware with utility (not written yet) or use EPROM programmer
11. sudo apt install git golang
12. git clone https://github.com/tjboldt/Apple2-IO-RPi.git
13. cd Apple2-IO-RPi
14. go get
15. go build
16. ./Apple2-IO-RPi Apple2-IO-RPi.hdv
18. Setup the Driver as a service or to autostart via cronjob (crontab -e then add the line @reboot /home/pi/Apple2-IO-RPi/RaspberryPi/Apple2-IO-RPi /home/pi/Apple2-IO-RPi/RaspberryPi/Apple2-IO-RPi.hdv > /home/pi/Apple2-IO-RPi/RaspberryPi/Apple2-IO-RPi.log)

## Similar Project
If you prefer having Apple II peripherals control a Raspberry Pi rather than simply using the Raspberry Pi to provide storage, network access and processing to the Apple II, have a look at David Schmenk's excellent [Apple2Pi](https://github.com/dschmenk/apple2pi) project. 
