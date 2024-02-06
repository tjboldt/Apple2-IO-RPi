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

#include <stdio.h>
#include <pico/printf.h>
#include <pico/stdlib.h>
#include <pico/multicore.h>

#include "board.h"

#ifdef TRACE
void uart_printf(uart_inst_t *uart, const char *format, ...) {
    static char buffer[0x100];

    va_list va;
    va_start(va, format);
    vsnprintf(buffer, sizeof(buffer), format, va);
    va_end(va);

    buffer[0xFF] = '\0';
    uart_puts(uart, buffer);
}
#endif

void main(void) {
    multicore_launch_core1(board);

    stdio_init_all();
    stdio_set_translate_crlf(&stdio_usb, false);

#ifdef PICO_DEFAULT_LED_PIN
    gpio_init(PICO_DEFAULT_LED_PIN);
    gpio_set_dir(PICO_DEFAULT_LED_PIN, GPIO_OUT);
#endif

#ifdef TRACE
    uart_init(uart0, 115200);
    uart_set_translate_crlf(uart0, true);
    gpio_set_function(PICO_DEFAULT_UART_TX_PIN, GPIO_FUNC_UART);
    gpio_set_function(PICO_DEFAULT_UART_RX_PIN, GPIO_FUNC_UART);
#endif

    while (true) {
        bool conn = stdio_usb_connected();
        if (conn) {
            if (multicore_fifo_rvalid()) {
                uint32_t data = multicore_fifo_pop_blocking();
                putchar(data);
#ifdef TRACE
                uart_printf(uart0, "> %02X\n", data);
#endif
            }
        }

        if (multicore_fifo_wready()) {
            int data = getchar_timeout_us(0);
            if (data != PICO_ERROR_TIMEOUT) {
                multicore_fifo_push_blocking(data);
#ifdef TRACE
                uart_printf(uart0, "< %02X\n", data);
#endif
            }
        }

#ifdef PICO_DEFAULT_LED_PIN
        gpio_put(PICO_DEFAULT_LED_PIN, conn);
#endif
    }
}
