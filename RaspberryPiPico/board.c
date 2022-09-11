/*

MIT License

Copyright (c) 2022 Oliver Schmidt (https://a2retro.de/)

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

*/

#include "pico/multicore.h"

#include "bus.pio.h"

#include "board.h"

extern const __attribute__((aligned(4))) uint8_t firmware[];

static uint32_t page;

void __time_critical_func(board)(void) {
    for (uint gpio = gpio_addr; gpio < gpio_addr + size_addr; gpio++) {
        gpio_init(gpio);
        gpio_set_pulls(gpio, false, false);  // floating
    }

    for (uint gpio = gpio_data; gpio < gpio_data + size_data; gpio++) {
        pio_gpio_init(pio0, gpio);
        gpio_set_pulls(gpio, false, false);  // floating
    }

    gpio_init(gpio_enbl);
    gpio_set_pulls(gpio_enbl, false, false);  // floating

    uint offset;

    offset = pio_add_program(pio0, &enbl_program);
    enbl_program_init(offset);

    offset = pio_add_program(pio0, &write_program);
    write_program_init(offset);

    offset = pio_add_program(pio0, &read_program);
    read_program_init(offset);

    page = 0;

    while (true) {
        uint32_t enbl = pio_sm_get_blocking(pio0, sm_enbl);
        uint32_t addr = enbl & 0x0FFF;
        uint32_t io   = enbl & 0x0F00;  // IOSTRB or IOSEL
        uint32_t strb = enbl & 0x0800;  // IOSTRB
        uint32_t read = enbl & 0x1000;  // R/W

        if (read) {
            if (!io) {  // DEVSEL
                switch (addr & 0x7) {
                    case 0x3:
                        pio_sm_put(pio0, sm_read, !multicore_fifo_rvalid() << 7 |
                                                  !multicore_fifo_wready() << 6);
                        break;
                    case 0x6:
                        pio_sm_put(pio0, sm_read, sio_hw->fifo_rd);
                        break;
                }
            } else {
                if (!strb) {
                    pio_sm_put(pio0, sm_read, firmware[page | addr]);
                }
            }
        } else {
            uint32_t data = pio_sm_get_blocking(pio0, sm_write);
            if (!io) {  // DEVSEL
                switch (addr & 0x7) {
                    case 0x5:
                        sio_hw->fifo_wr = data;
                        break;
                    case 0x7:
                        page = (data & 0x30) << 7;
                        break;
                }
            }
        }
    }
}
