#!/bin/sh
ca65 DriveFirmware.asm -D SLOT=0 -o DriveSlot0.o || exit 1
ca65 DriveFirmware.asm -D SLOT=1 -o DriveSlot1.o || exit 1
ca65 DriveFirmware.asm -D SLOT=2 -o DriveSlot2.o || exit 1
ca65 DriveFirmware.asm -D SLOT=3 -o DriveSlot3.o || exit 1
ca65 DriveFirmware.asm -D SLOT=4 -o DriveSlot4.o || exit 1
ca65 DriveFirmware.asm -D SLOT=5 -o DriveSlot5.o || exit 1
ca65 DriveFirmware.asm -D SLOT=6 -o DriveSlot6.o || exit 1
ca65 DriveFirmware.asm -D SLOT=7 -o DriveSlot7.o --listing DriveFirmware.lst.new --list-bytes 255 || exit 1
ld65 DriveSlot0.o DriveSlot1.o DriveSlot2.o DriveSlot3.o DriveSlot4.o DriveSlot5.o DriveSlot6.o DriveSlot7.o -o DriveFirmware.bin.new -C ../.cicd/none.cfg || exit 1

ca65 MenuFirmware.asm -D SLOT=0 -o MenuSlot0.o || exit 1
ca65 MenuFirmware.asm -D SLOT=1 -o MenuSlot1.o || exit 1
ca65 MenuFirmware.asm -D SLOT=2 -o MenuSlot2.o || exit 1
ca65 MenuFirmware.asm -D SLOT=3 -o MenuSlot3.o || exit 1
ca65 MenuFirmware.asm -D SLOT=4 -o MenuSlot4.o || exit 1
ca65 MenuFirmware.asm -D SLOT=5 -o MenuSlot5.o || exit 1
ca65 MenuFirmware.asm -D SLOT=6 -o MenuSlot6.o || exit 1
ca65 MenuFirmware.asm -D SLOT=7 -o MenuSlot7.o --listing MenuFirmware.lst.new --list-bytes 255 || exit 1
ld65 MenuSlot0.o MenuSlot1.o MenuSlot2.o MenuSlot3.o MenuSlot4.o MenuSlot5.o MenuSlot6.o MenuSlot7.o -o MenuFirmware.bin.new -C ../.cicd/none.cfg || exit 1

ca65 CommandFirmware.asm -D SLOT=0 -o CommandSlot0.o || exit 1
ca65 CommandFirmware.asm -D SLOT=1 -o CommandSlot1.o || exit 1
ca65 CommandFirmware.asm -D SLOT=2 -o CommandSlot2.o || exit 1
ca65 CommandFirmware.asm -D SLOT=3 -o CommandSlot3.o || exit 1
ca65 CommandFirmware.asm -D SLOT=4 -o CommandSlot4.o || exit 1
ca65 CommandFirmware.asm -D SLOT=5 -o CommandSlot5.o || exit 1
ca65 CommandFirmware.asm -D SLOT=6 -o CommandSlot6.o || exit 1
ca65 CommandFirmware.asm -D SLOT=7 -o CommandSlot7.o --listing CommandFirmware.lst.new --list-bytes 255 || exit 1
ld65 CommandSlot0.o CommandSlot1.o CommandSlot2.o CommandSlot3.o CommandSlot4.o CommandSlot5.o CommandSlot6.o CommandSlot7.o -o CommandFirmware.bin.new -C ../.cicd/none.cfg || exit 1

ca65 FileAccessFirmware.asm -D SLOT=0 -o FileAccessSlot0.o || exit 1
ca65 FileAccessFirmware.asm -D SLOT=1 -o FileAccessSlot1.o || exit 1
ca65 FileAccessFirmware.asm -D SLOT=2 -o FileAccessSlot2.o || exit 1
ca65 FileAccessFirmware.asm -D SLOT=3 -o FileAccessSlot3.o || exit 1
ca65 FileAccessFirmware.asm -D SLOT=4 -o FileAccessSlot4.o || exit 1
ca65 FileAccessFirmware.asm -D SLOT=5 -o FileAccessSlot5.o || exit 1
ca65 FileAccessFirmware.asm -D SLOT=6 -o FileAccessSlot6.o || exit 1
ca65 FileAccessFirmware.asm -D SLOT=7 -o FileAccessSlot7.o --listing FileAccessFirmware.lst.new --list-bytes 255 || exit 1
ld65 FileAccessSlot0.o FileAccessSlot1.o FileAccessSlot2.o FileAccessSlot3.o FileAccessSlot4.o FileAccessSlot5.o FileAccessSlot6.o FileAccessSlot7.o -o FileAccessFirmware.bin.new -C ../.cicd/none.cfg  || exit 1

cat \
DriveFirmware.bin.new CommandFirmware.bin.new FileAccessFirmware.bin.new MenuFirmware.bin.new \
> AT28C64B.bin.new

ca65 Shell.asm -o Shell.o --listing Shell.lst.new || exit 1
ld65 Shell.o -o Shell.bin.new -C ../.cicd/none.cfg  || exit 1

ca65 RPi.Command.asm -o RPi.Command.o --listing RPi.Command.lst.new || exit 1
ld65 RPi.Command.o -o RPi.Command.bin.new -C ../.cicd/none.cfg  || exit 1

rm ./*.o
diff RPi.Command.bin RPi.Command.bin.new || exit 1
diff AT28C64B.bin AT28C64B.bin.new || exit 1
diff Shell.bin Shell.bin.new || exit 1

rm ./*.new
