# Apple2-IO-RPi
Apple II expansion card using a Raspberry Pi for I/O

![Image of Board](/Hardware/Apple2IORPi.jpg)

## Purpose
The purpose of this project is to provide I/O for an Apple II series 8 bit computer via a Raspberry Pi Zero W which is powered by the Apple II expansion bus. This includes using the attached RPi Zero W for it's storage, network and processor to provide new functionality for the Apple II.

## Features
1. Boot message which waits for RPi to be ready
2. ProDOS bootable drive from image stored on RPi
3. Linux bash shell to the RPi from the Apple II including some VT100 support via `-SHELL` from ProDOS
4. Load binary files directly from the RPi to the II (temporarily removed from firmware but will come back as a utility or dynamic virtual drive
5. Update Apple II firmware in place from image on RPi (note, this is done per slot)
6. Supports two drive images at the same time
7. Supports "RPI" command from BASIC to execute Linux commands from the command prompt or inside BASIC programs: `10 PRINT CHR$(4);"RPI ls -al /"`

## Project Status
So far, this is still a project and not a finished product. The current prototype is on the sixth revision and about 30 of the previous prototypes have been assembled and tested. The sixth prototype has not been tested or assembled yet, just designed although it is functionally equivalent other than a new jumper to select internal/external power.

The card enables the Apple II to boot from and write to virtual hard drive images stored on the RPi in any slot (except slot 3), execute Linux commands from Applesoft BASIC and run a bash shell with VT100 emulation. The code has very few tests and is incomplete. Note that currently the firmware assumes an 80 column card is in slot 3 and than you have lowercase support. Most development has been done with an enhanced Apple //e with the card in slot 7. If you have other drive controllers earlier in the boot cycle, you can still boot from the Apple2-IO-RPi. For example, if the card was in slot 4, you could type `PR#4` from the BASIC prompt to boot the card. Note that the Raspberry Pi Zero W (and W 2) consume 170 - 250 mA and there is only 500 mA available to all expansion slots according to Apple. It is not recommended to have a lot of other cards in the system at the same time. With the sixth revision of the prototype, it is possible to remove the power jumper and run the RPi on an external USB power source. If configured for external power, note that the card's firmware will hang on boot without USB power on as the latch chips are powered by the 3.3V output of the RPi. 

If you have a problem or idea for enhancement, log an issue [here](https://github.com/tjboldt/Apple2-IO-RPi/issues) or start a [discussion](https://github.com/tjboldt/Apple2-IO-RPi/discussions/categories/general). I recommend starring/watching the project for updates on GitHub. You are welcome to fork the project and submit pull requests which I will review.

## Roadmap
1. Allow more than two virtual hard drives at a time
2. In-memory virtual drive representing current working directory in Linux for ease of copying files between Linux and ProDOS
3. Update GPIO library to reduce power consumption (remove busy-wait loops and use edge triggering)
4. Proper ProDOS clock driver (currently just directly sets values on block reads)
5. Bi-directional image conversion between common formats and HIRES
6. Remote code execution 
7. Proxy VNC connection, rendering as Apple II compatible graphics

## Setup
[Setup card from scratch](https://github.com/tjboldt/Apple2-IO-RPi/discussions/63)

[Setup if you received a complete board from me](https://github.com/tjboldt/Apple2-IO-RPi/discussions/64)

[Update to latest](https://github.com/tjboldt/Apple2-IO-RPi/discussions/65)

## Contributions/Thanks
- Hans Hübner (@hanshuebner) for help cleaning up schematics
- Scott (@figital) for assembling early prototypes and recommending the AT28C64B chip
- Philip A. for help with 3.3V conversion
- David Schmenk (@dschmenk) for creating his Apple2Pi project, proving that Raspberry Pi can be powered by the Apple II expansion bus
- Wyatt Wong (@wyatt-wong) for testing with multiple cards
- (@Abysmal) for shell and rpi.command testing
- (@bfranske) for suggesting adding the 5V jumper
- Tim Boldt (@timboldt) for recommending removing sysfs based GPIO code

## Similar Project
If you prefer having Apple II peripherals control a Raspberry Pi rather than simply using the Raspberry Pi to provide storage, network access and processing to the Apple II, have a look at David Schmenk's excellent [Apple2Pi](https://github.com/dschmenk/apple2pi) project. I am often asked about differences between these two projects. They are similar in some ways but essentially opposite. The Apple2Pi is meant for the primary machine to be the RPi, using the Apple II for it's peripherals. The Apple2-IO-RPi is meant to have the Apple II as the primary machine and just use the RPi for its processing, storage and network.
