sudo apt install cmake gcc-arm-none-eabi libnewlib-arm-none-eabi libstdc++-arm-none-eabi-newlib -y
cd ~
git clone https://github.com/raspberrypi/pico-sdk.git
cd ~/pico-sdk/
git submodule update --init
