#!/bin/sh
ca65 test.asm
ld65 test.o -o test.bin -t none
