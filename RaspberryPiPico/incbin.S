.section .time_critical.firmware
.global firmware
.type firmware, %object
.balign 4
firmware:
.incbin "../RaspberryPi/driveimage/AT28C64B.bin"
