#!/bin/sh
ca65 DriveFirmware.asm -D STARTSLOT=\$c000 -o Slot0.o 
ca65 DriveFirmware.asm -D STARTSLOT=\$c100 -o Slot1.o 
ca65 DriveFirmware.asm -D STARTSLOT=\$c200 -o Slot2.o
ca65 DriveFirmware.asm -D STARTSLOT=\$c300 -o Slot3.o
ca65 DriveFirmware.asm -D STARTSLOT=\$c400 -o Slot4.o
ca65 DriveFirmware.asm -D STARTSLOT=\$c500 -o Slot5.o
ca65 DriveFirmware.asm -D STARTSLOT=\$c600 -o Slot6.o
ca65 DriveFirmware.asm -D STARTSLOT=\$c700 -o Slot7.o
ld65 Slot0.o Slot1.o Slot2.o Slot3.o Slot4.o Slot5.o Slot6.o Slot7.o -o DriveFirmware.bin -t none
cat \
DriveFirmware.bin DriveFirmware.bin DriveFirmware.bin DriveFirmware.bin \
DriveFirmware.bin DriveFirmware.bin DriveFirmware.bin DriveFirmware.bin \
DriveFirmware.bin DriveFirmware.bin DriveFirmware.bin DriveFirmware.bin \
DriveFirmware.bin DriveFirmware.bin DriveFirmware.bin DriveFirmware.bin \
> Firmware_27256_EPROM.bin
rm Slot*.o
rm DriveFirmware.bin

ca65 Rpi.Command.asm -o Rpi.Command.o
ld65 Rpi.Command.o -o Rpi.Command.bin -t none
rm Rpi.Command.o
