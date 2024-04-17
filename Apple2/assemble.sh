#!/bin/sh
HW_TYPE=${1-0}

ca65 DriveFirmware.asm -D HW_TYPE=$HW_TYPE -D SLOT=0 -o DriveSlot0.o || exit 1
ca65 DriveFirmware.asm -D HW_TYPE=$HW_TYPE -D SLOT=1 -o DriveSlot1.o || exit 1
ca65 DriveFirmware.asm -D HW_TYPE=$HW_TYPE -D SLOT=2 -o DriveSlot2.o || exit 1
ca65 DriveFirmware.asm -D HW_TYPE=$HW_TYPE -D SLOT=3 -o DriveSlot3.o || exit 1
ca65 DriveFirmware.asm -D HW_TYPE=$HW_TYPE -D SLOT=4 -o DriveSlot4.o || exit 1
ca65 DriveFirmware.asm -D HW_TYPE=$HW_TYPE -D SLOT=5 -o DriveSlot5.o || exit 1
ca65 DriveFirmware.asm -D HW_TYPE=$HW_TYPE -D SLOT=6 -o DriveSlot6.o || exit 1
ca65 DriveFirmware.asm -D HW_TYPE=$HW_TYPE -D SLOT=7 -o DriveSlot7.o --listing DriveFirmware.lst --list-bytes 255 || exit 1
ld65 DriveSlot0.o DriveSlot1.o DriveSlot2.o DriveSlot3.o DriveSlot4.o DriveSlot5.o DriveSlot6.o DriveSlot7.o -o DriveFirmware.bin -C ../.cicd/none.cfg || exit 1

ca65 MenuFirmware.asm -D HW_TYPE=$HW_TYPE -D SLOT=0 -o MenuSlot0.o || exit 1
ca65 MenuFirmware.asm -D HW_TYPE=$HW_TYPE -D SLOT=1 -o MenuSlot1.o || exit 1
ca65 MenuFirmware.asm -D HW_TYPE=$HW_TYPE -D SLOT=2 -o MenuSlot2.o || exit 1
ca65 MenuFirmware.asm -D HW_TYPE=$HW_TYPE -D SLOT=3 -o MenuSlot3.o || exit 1
ca65 MenuFirmware.asm -D HW_TYPE=$HW_TYPE -D SLOT=4 -o MenuSlot4.o || exit 1
ca65 MenuFirmware.asm -D HW_TYPE=$HW_TYPE -D SLOT=5 -o MenuSlot5.o || exit 1
ca65 MenuFirmware.asm -D HW_TYPE=$HW_TYPE -D SLOT=6 -o MenuSlot6.o || exit 1
ca65 MenuFirmware.asm -D HW_TYPE=$HW_TYPE -D SLOT=7 -o MenuSlot7.o --listing MenuFirmware.lst --list-bytes 255 || exit 1
ld65 MenuSlot0.o MenuSlot1.o MenuSlot2.o MenuSlot3.o MenuSlot4.o MenuSlot5.o MenuSlot6.o MenuSlot7.o -o MenuFirmware.bin -C ../.cicd/none.cfg || exit 1

ca65 CommandFirmware.asm -D HW_TYPE=$HW_TYPE -D SLOT=0 -o CommandSlot0.o || exit 1
ca65 CommandFirmware.asm -D HW_TYPE=$HW_TYPE -D SLOT=1 -o CommandSlot1.o || exit 1
ca65 CommandFirmware.asm -D HW_TYPE=$HW_TYPE -D SLOT=2 -o CommandSlot2.o || exit 1
ca65 CommandFirmware.asm -D HW_TYPE=$HW_TYPE -D SLOT=3 -o CommandSlot3.o || exit 1
ca65 CommandFirmware.asm -D HW_TYPE=$HW_TYPE -D SLOT=4 -o CommandSlot4.o || exit 1
ca65 CommandFirmware.asm -D HW_TYPE=$HW_TYPE -D SLOT=5 -o CommandSlot5.o || exit 1
ca65 CommandFirmware.asm -D HW_TYPE=$HW_TYPE -D SLOT=6 -o CommandSlot6.o || exit 1
ca65 CommandFirmware.asm -D HW_TYPE=$HW_TYPE -D SLOT=7 -o CommandSlot7.o --listing CommandFirmware.lst --list-bytes 255 || exit 1
ld65 CommandSlot0.o CommandSlot1.o CommandSlot2.o CommandSlot3.o CommandSlot4.o CommandSlot5.o CommandSlot6.o CommandSlot7.o -o CommandFirmware.bin -C ../.cicd/none.cfg || exit 1

ca65 FileAccessFirmware.asm -D HW_TYPE=$HW_TYPE -D SLOT=0 -o FileAccessSlot0.o || exit 1
ca65 FileAccessFirmware.asm -D HW_TYPE=$HW_TYPE -D SLOT=1 -o FileAccessSlot1.o || exit 1
ca65 FileAccessFirmware.asm -D HW_TYPE=$HW_TYPE -D SLOT=2 -o FileAccessSlot2.o || exit 1
ca65 FileAccessFirmware.asm -D HW_TYPE=$HW_TYPE -D SLOT=3 -o FileAccessSlot3.o || exit 1
ca65 FileAccessFirmware.asm -D HW_TYPE=$HW_TYPE -D SLOT=4 -o FileAccessSlot4.o || exit 1
ca65 FileAccessFirmware.asm -D HW_TYPE=$HW_TYPE -D SLOT=5 -o FileAccessSlot5.o || exit 1
ca65 FileAccessFirmware.asm -D HW_TYPE=$HW_TYPE -D SLOT=6 -o FileAccessSlot6.o || exit 1
ca65 FileAccessFirmware.asm -D HW_TYPE=$HW_TYPE -D SLOT=7 -o FileAccessSlot7.o --listing FileAccessFirmware.lst --list-bytes 255 || exit 1
ld65 FileAccessSlot0.o FileAccessSlot1.o FileAccessSlot2.o FileAccessSlot3.o FileAccessSlot4.o FileAccessSlot5.o FileAccessSlot6.o FileAccessSlot7.o -o FileAccessFirmware.bin -C ../.cicd/none.cfg  || exit 1

cat \
DriveFirmware.bin CommandFirmware.bin FileAccessFirmware.bin MenuFirmware.bin \
> ../RaspberryPi/driveimage/AT28C64B.bin

ca65 Shell.asm -D HW_TYPE=$HW_TYPE -o Shell.o --listing Shell.lst || exit 1
ld65 Shell.o -o ../RaspberryPi/driveimage/Shell.bin -C ../.cicd/none.cfg  || exit 1

ca65 RPi.Command.asm -D HW_TYPE=$HW_TYPE -o RPi.Command.o --listing RPi.Command.lst || exit 1
ld65 RPi.Command.o -o ../RaspberryPi/driveimage/RPi.Command.bin -C ../.cicd/none.cfg  || exit 1

ca65 Clock.Driver.asm -D HW_TYPE=$HW_TYPE -o Clock.Driver.o --listing Clock.Driver.lst || exit 1
ld65 Clock.Driver.o -o ../RaspberryPi/driveimage/Clock.Driver.bin -C ../.cicd/none.cfg  || exit 1

rm ./*.o
rm DriveFirmware.bin
rm MenuFirmware.bin
rm CommandFirmware.bin
rm FileAccessFirmware.bin
