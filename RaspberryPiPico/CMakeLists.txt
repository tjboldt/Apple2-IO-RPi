set(PROJECT_NAME Apple2-IO-RPi)

cmake_minimum_required(VERSION 3.13...3.27)
include(pico_sdk_import.cmake)
project(${PROJECT_NAME} C CXX ASM)
pico_sdk_init()

add_executable(${PROJECT_NAME})
pico_add_extra_outputs(${PROJECT_NAME})

pico_enable_stdio_uart(${PROJECT_NAME} 0)
pico_enable_stdio_usb(${PROJECT_NAME} 1)

include(FetchContent)
FetchContent_Declare(a2pico
        GIT_REPOSITORY https://github.com/oliverschmidt/a2pico.git
        GIT_TAG main
        )
FetchContent_MakeAvailable(a2pico)

set_source_files_properties(incbin.S OBJECT_DEPENDS ../RaspberryPi/driveimage/AT28C64B.bin)
add_custom_command(
        WORKING_DIRECTORY ../Apple2
        COMMAND ./assemble.sh 1
        OUTPUT ../RaspberryPi/driveimage/AT28C64B.bin
        VERBATIM
        )

target_sources(${PROJECT_NAME} PRIVATE
        main.c
        board.c
        incbin.S
        )

target_link_libraries(${PROJECT_NAME} PRIVATE
        pico_stdlib
        pico_multicore
        a2pico
        )
