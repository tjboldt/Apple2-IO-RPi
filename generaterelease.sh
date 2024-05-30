#!/bin/sh
mkdir release
cd Apple2
. assemble.sh 1
cd ../RaspberryPiPico
. build.sh
cp Apple2-IO-RPi.elf ../release
cd ../RaspberryPi/driveimage
../../.cicd/ProDOS-Utilities -d ../../release/Apple2-IO-RPi.pico.hardware.hdv -c create -v APPLE2.IO.RPI
../../.cicd/ProDOS-Utilities -d ../../release/Apple2-IO-RPi.pico.hardware.hdv -c putall
../../.cicd/ProDOS-Utilities -d ../../release/Apple2-IO-RPi.pico.hardware.hdv -c ls
cd ../../Apple2
. assemble.sh 0
cd ../RaspberryPi/driveimage
cp AT28C64B.bin ../../release
../../.cicd/ProDOS-Utilities -d ../../release/Apple2-IO-RPi.classic.hardware.hdv -c create -v APPLE2.IO.RPI
../../.cicd/ProDOS-Utilities -d ../../release/Apple2-IO-RPi.classic.hardware.hdv -c putall
../../.cicd/ProDOS-Utilities -d ../../release/Apple2-IO-RPi.classic.hardware.hdv -c ls
