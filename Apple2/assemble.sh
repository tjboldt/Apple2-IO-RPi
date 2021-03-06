#!/bin/sh
ca65 Firmware.asm -D STARTSLOT=\$c000 -o Slot0.o 
ca65 Firmware.asm -D STARTSLOT=\$c100 -o Slot1.o --listing Firmware1.lst
ca65 Firmware.asm -D STARTSLOT=\$c200 -o Slot2.o --listing Firmware2.lst
ca65 Firmware.asm -D STARTSLOT=\$c300 -o Slot3.o --listing Firmware3.lst
ca65 Firmware.asm -D STARTSLOT=\$c400 -o Slot4.o --listing Firmware4.lst
ca65 Firmware.asm -D STARTSLOT=\$c500 -o Slot5.o --listing Firmware5.lst
ca65 Firmware.asm -D STARTSLOT=\$c600 -o Slot6.o --listing Firmware6.lst
ca65 Firmware.asm -D STARTSLOT=\$c700 -o Slot7.o --listing Firmware7.lst
ld65 Slot0.o Slot1.o Slot2.o Slot3.o Slot4.o Slot5.o Slot6.o Slot7.o -o Firmware.bin -t none
cat \
Firmware.bin Firmware.bin Firmware.bin Firmware.bin \
Firmware.bin Firmware.bin Firmware.bin Firmware.bin \
Firmware.bin Firmware.bin Firmware.bin Firmware.bin \
Firmware.bin Firmware.bin Firmware.bin Firmware.bin \
> Firmware_27256_EPROM.bin
