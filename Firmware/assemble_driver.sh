#!/bin/sh
ca65 Driver.asm --listing Driver.lst
ld65 Driver.o -o Driver.bin -t none
