# Apple2-IO-RPi
Apple II expansion card using a Raspberry Pi for I/O

![Image of Board](/Hardware/Apple2IORPi.jpg)

## Purpose
The purpose of this project is to provide I/O for an Apple II series 8 bit computer via a Raspberry Pi Zero W which is powered by the Apple II expansion bus. Initially this is storage via virtual ProDOS compatible drive. There is also support for executing Linux commands from the Apple II via a rudimentary console. Next might be adding virtual serial card support over wifi. Future enhancements could use the RPi for more complex processing as per request from the Apple II. For example, the Apple II could request a web page or application and the RPi could calculate this in Apple II hi-res graphics mode and send the image data back to the II for display purposes.

## Project Status
This is an early stage project. Currently four prototype revisions have been assembled and tested. It is now possible for the Apple II to boot from and write to a virutal hard drive image stored on the RPi in any slot and execute simple commands on the RPi via an Apple II application. The code has no error handling or tests yet and is incomplete.

## Roadmap
1. DONE - Build initial prototype that reads/writes virtual hard drive
2. DONE - Create firmware to make the card a bootable device
3. DONE - Fix board with updated second prototype PCB
4. DONE - Add ProDOS clock support
5. WORKING - Add simple RPi command execution (does not keep context between commands yet)
6. DONE - Create new schematic/PCB with faster data transfera
7. DONE - Add hardware support for in-system firmware updates
8. Add firmware update software
9. Add proper ProDOS clock driver (currently just directly sets values on block reads)
10. Add directory context to RPi command execution 
11. Add full terminal emulation 
12. Add support for direct file read/write without drive image
13. Add web service call support
14. Proxy VNC connection, rendering as Apple II compatible graphics

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
10. Update firmware with utility (not written yet) or use EPROM programmer
11. Install Git and Go on the RPi over ssh
12. go build Driver.go
13. Copy a ProDOS hard drive image onto the RPi
14. ./Driver HARDDRIVE.hdv
15. Setup the Driver as a service or to autostart via cronjob

## Similar Project
If you prefer having Apple II peripherals control a Raspberry Pi rather than simply using the Raspberry Pi to provide storage, network access and processing to the Apple II, have a look at David Schmenk's excellent [Apple2Pi](https://github.com/dschmenk/apple2pi) project.  

