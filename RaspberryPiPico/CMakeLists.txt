set(PROJECT_NAME Apple2-IO-RPi)

cmake_minimum_required(VERSION 3.12)

include(pico_sdk_import.cmake)
pico_sdk_init()

project(${PROJECT_NAME} C CXX ASM)
add_executable(${PROJECT_NAME})
pico_add_extra_outputs(${PROJECT_NAME})

pico_enable_stdio_uart(${PROJECT_NAME} 0)
pico_enable_stdio_usb(${PROJECT_NAME} 1)

pico_generate_pio_header(${PROJECT_NAME} ${CMAKE_CURRENT_LIST_DIR}/bus.pio)

target_sources(${PROJECT_NAME} PRIVATE
        main.c
        board.c
        incbin.S
        )

target_link_libraries(${PROJECT_NAME} PRIVATE
        pico_stdlib
        pico_multicore
        hardware_pio
        )
