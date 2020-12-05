# Apple2-IO-RPi
Apple II expansion card using a Raspberry Pi for I/O

![Image of Board](/Hardware/Apple2IORPi.jpg)

## Project Status
This is a very early stage project. Currently one board has been assembled and tested. There is no firmware, however with some simple Go and Node.js code, it was determined that communications are working and Go is significantly faster.

## Purpose
The purpose of this project is to provide I/O for an Apple II series 8 bit computer via a Raspberry Pi Zero W which is powered by the Apple II expansion bus. Initially this would be storage via virtual ProDOS compatible drive. Next might be adding virtual serial card support over wifi. Future enhancements could use the RPi for more complex processing as per request from the Apple II. For example, the Apple II could request a web page or application and the RPi could calculate this in Apple II hi-res graphics mode and send the image data back to the II for display purposes.

## Setup
1. Have PCBs made from the gerber and drill files in the Hardware folder
2. Solder chips, header and capacitors in place
3. Attach Raspberry Pi Zero W facing outward from the card
4. Install Raspberry Pi OS on microSD card
5. Configure wifi in boot/wpa_supplicant.conf
6. Add empty ssh file boot (for ssh access over wifi)
7. Put microSD card in the RPi
8. Install the expansion card into the Apple II
9. Power on the Apple II
10. Install NodeJS 11 for ARMv6 on the RPi over ssh
11. sudo npm install -g pm2
12. pm2 startup systemd
13. Run the command from the output of the previous step
14. pm2 start server.js

## Similar Project
If you're looking for complete hardware design or prefer having Apple II peripherals control a control a Raspberry Pi rather than simply using the Raspberry Pi to provide storage and network access to the Apple II, have a look at David Schmenk's excellent [Apple2Pi](https://github.com/dschmenk/apple2pi) project.  


