---
layout: page
title: "Apple2-IO-RPi"
permalink: /
---

# Apple2-IO-RPi
Apple II expansion card using a Raspberry Pi for I/O

![Image of Board](/images/Apple2IORPi.jpg)

## Purpose
The purpose of this project is to provide I/O for an Apple II series 8 bit computer via a Raspberry Pi Zero W which is powered by the Apple II expansion bus. This includes using the attached RPi Zero W for it's storage, network and processor to provide new functionality for the Apple II.

## Project Status
So far, this is a project and not a finished product. The current prototype is on the fifth revision and a few have been assembled and tested. It is now possible for the Apple II to boot from and write to a virutal hard drive image stored on the RPi in any slot and run a bash shell on the RPi via the Apple II. The code has very few tests and is incomplete. Note that currently the firmware assumes an 80 column card is in slot 3 and than you have lowercase support. If you have a problem or idea for enhancement, log an issue [here](https://github.com/tjboldt/Apple2-IO-RPi/issues). I recommend starring/watching the project for updates on GitHub. You are welcome to fork the project and submit pull requests which I will review.

## Features
1. Boot message which waits for RPi to be ready
2. ProDOS bootable drive from image stored on RPi
3. Execute Linux bash shell on the RPi from the Apple II
4. Load binary files directly from the RPi to the II
5. Update Apple II firmware in place from image on RPi
6. Supports two drive images at the same time (Note: backward compatible with previous firmware but requires firmware update in order to work with two drives)
7. Supports "RPI" command from BASIC to execute Linux commands from the command prompt or inside BASIC programs: `10 PRINT CHR$(4);"RPI ls -al /"`

## Roadmap
1. Extend BASIC.SYSTEM commands
    1. RPI - Execute a single Linux command (DONE)
    2. SH - Open a Linux shell
    3. WGET - Download and save to file 
2. Shell improvements
3. Proper ProDOS clock driver (currently just directly sets values on block reads)
4. Bi-directional image conversion between common formats and HIRES
5. Remote code execution 
6. Proxy VNC connection, rendering as Apple II compatible graphics

