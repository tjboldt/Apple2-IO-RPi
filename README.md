# Apple2-IO-RPi

See the [upstream repo](https://github.com/tjboldt/Apple2-IO-RPi) for the original, full info.

## Installing

1. Flash [Apple2-IO-RPi.uf2](https://github.com/oliverschmidt/apple2-io-rpi/releases/latest/download/Apple2-IO-RPi.uf2) to the Raspberry Pi Pico.

2. Execute steps __4.)__ - __10.)__ and __12.)__ on [Setup new card from scratch](https://github.com/tjboldt/Apple2-IO-RPi/discussions/63).

3. Execute `wget --no-cache -O - https://raw.githubusercontent.com/oliverschmidt/Apple2-IO-RPi/main/RaspberryPi/setup.sh | bash`

## Using

1. Access the two ProDOS 8 drives either via cold boot or `PR#n`.

2. Make sure to check out `SHELL` and `RPI.COMMAND` in drive 1.

## Building

1. Execute `assemble 1` in `/Apple2` (requires [cc65](https://cc65.github.io/) and the [ProDOS-Utilities](https://github.com/tjboldt/ProDOS-Utilities)).

2. Build the project in `/RaspberryPiPico` (requires the [Raspberry Pico C/C++ SDK](https://www.raspberrypi.com/documentation/microcontrollers/c_sdk.html)).

3. Execute `go build` in `/RaspberryPi/apple2driver` (requires [Go](https://go.dev/)).
