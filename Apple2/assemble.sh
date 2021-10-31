#!/bin/sh
ca65 DriveFirmware.asm -D SLOT=0 -o DriveSlot0.o 
ca65 DriveFirmware.asm -D SLOT=1 -o DriveSlot1.o 
ca65 DriveFirmware.asm -D SLOT=2 -o DriveSlot2.o
ca65 DriveFirmware.asm -D SLOT=3 -o DriveSlot3.o
ca65 DriveFirmware.asm -D SLOT=4 -o DriveSlot4.o
ca65 DriveFirmware.asm -D SLOT=5 -o DriveSlot5.o 
ca65 DriveFirmware.asm -D SLOT=6 -o DriveSlot6.o
ca65 DriveFirmware.asm -D SLOT=7 -o DriveSlot7.o --listing DriveFirmware.lst --list-bytes 255
ld65 DriveSlot0.o DriveSlot1.o DriveSlot2.o DriveSlot3.o DriveSlot4.o DriveSlot5.o DriveSlot6.o DriveSlot7.o -o DriveFirmware.bin -t none

ca65 MenuFirmware.asm -D SLOT=0 -o MenuSlot0.o 
ca65 MenuFirmware.asm -D SLOT=1 -o MenuSlot1.o 
ca65 MenuFirmware.asm -D SLOT=2 -o MenuSlot2.o
ca65 MenuFirmware.asm -D SLOT=3 -o MenuSlot3.o
ca65 MenuFirmware.asm -D SLOT=4 -o MenuSlot4.o
ca65 MenuFirmware.asm -D SLOT=5 -o MenuSlot5.o 
ca65 MenuFirmware.asm -D SLOT=6 -o MenuSlot6.o
ca65 MenuFirmware.asm -D SLOT=7 -o MenuSlot7.o --listing MenuFirmware.lst --list-bytes 255
ld65 MenuSlot0.o MenuSlot1.o MenuSlot2.o MenuSlot3.o MenuSlot4.o MenuSlot5.o MenuSlot6.o MenuSlot7.o -o MenuFirmware.bin -t none

ca65 CommandFirmware.asm -D SLOT=0 -o CommandSlot0.o 
ca65 CommandFirmware.asm -D SLOT=1 -o CommandSlot1.o 
ca65 CommandFirmware.asm -D SLOT=2 -o CommandSlot2.o
ca65 CommandFirmware.asm -D SLOT=3 -o CommandSlot3.o
ca65 CommandFirmware.asm -D SLOT=4 -o CommandSlot4.o
ca65 CommandFirmware.asm -D SLOT=5 -o CommandSlot5.o
ca65 CommandFirmware.asm -D SLOT=6 -o CommandSlot6.o
ca65 CommandFirmware.asm -D SLOT=7 -o CommandSlot7.o --listing CommandFirmware.lst --list-bytes 255
ld65 CommandSlot0.o CommandSlot1.o CommandSlot2.o CommandSlot3.o CommandSlot4.o CommandSlot5.o CommandSlot6.o CommandSlot7.o -o CommandFirmware.bin -t none

ca65 FileAccessFirmware.asm -D SLOT=0 -o FileAccessSlot0.o 
ca65 FileAccessFirmware.asm -D SLOT=1 -o FileAccessSlot1.o 
ca65 FileAccessFirmware.asm -D SLOT=2 -o FileAccessSlot2.o
ca65 FileAccessFirmware.asm -D SLOT=3 -o FileAccessSlot3.o
ca65 FileAccessFirmware.asm -D SLOT=4 -o FileAccessSlot4.o
ca65 FileAccessFirmware.asm -D SLOT=5 -o FileAccessSlot5.o
ca65 FileAccessFirmware.asm -D SLOT=6 -o FileAccessSlot6.o
ca65 FileAccessFirmware.asm -D SLOT=7 -o FileAccessSlot7.o --listing FileAccessFirmware.lst --list-bytes 255
ld65 FileAccessSlot0.o FileAccessSlot1.o FileAccessSlot2.o FileAccessSlot3.o FileAccessSlot4.o FileAccessSlot5.o FileAccessSlot6.o FileAccessSlot7.o -o FileAccessFirmware.bin -t none

cat \
DriveFirmware.bin CommandFirmware.bin FileAccessFirmware.bin MenuFirmware.bin \
> AT28C64B.bin

rm *.o
rm DriveFirmware.bin
rm MenuFirmware.bin
rm CommandFirmware.bin
rm FileAccessFirmware.bin

ProDOS-Utilities -d ../RaspberryPi/Apple2-IO-RPi.hdv -c put -i AT28C64B.bin -p /APPLE2.IO.RPI/AT28C64B.BIN
ProDOS-Utilities -d ../RaspberryPi/Apple2-IO-RPi.hdv -c ls
