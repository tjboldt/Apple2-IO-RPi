# Apple2-IO-RPi

1. See the [upstream repo](https://github.com/tjboldt/Apple2-IO-RPi) for the original, full info.

2. This fork is based on [A2Pico](https://github.com/oliverschmidt/a2pico).

## Installing

1. Flash [Apple2-IO-RPi.uf2](https://github.com/oliverschmidt/apple2-io-rpi/releases/latest/download/Apple2-IO-RPi.uf2) to the Raspberry Pi Pico.

2. Execute steps __5.)__ - __8.)__ and __10.)__ on [Setup new card from scratch](https://github.com/tjboldt/Apple2-IO-RPi/discussions/63).

3. Execute `wget --no-cache -O - https://raw.githubusercontent.com/oliverschmidt/apple2-io-rpi/main/RaspberryPi/setup.sh | bash`

## Using

1. Access the two ProDOS 8 drives either via cold boot or `PR#n`.

2. Make sure to check out `SHELL` and `RPI.COMMAND` in drive 1.

## Building

1. Build the project in `/RaspberryPiPico` (requires the [Raspberry Pico C/C++ SDK](https://www.raspberrypi.com/documentation/microcontrollers/c_sdk.html)).

2. Execute `go build` in `/RaspberryPi/apple2driver` (requires [Go](https://go.dev/)).
