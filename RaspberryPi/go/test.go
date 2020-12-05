package main

import (
    "time"
    "periph.io/x/periph/conn/gpio"
    "periph.io/x/periph/conn/gpio/gpioreg"
    "periph.io/x/periph/host"
    "fmt"
)

func main() {
    host.Init()

    startTime := time.Now()

    out_write := gpioreg.ByName("GPIO5")
    out_read := gpioreg.ByName("GPIO11")
    out_command1 := gpioreg.ByName("GPIO9")
    out_command2 := gpioreg.ByName("GPIO10")
    out_bit3 := gpioreg.ByName("GPIO22")
    out_bit2 := gpioreg.ByName("GPIO27")
    out_bit1 := gpioreg.ByName("GPIO17")
    out_bit0 := gpioreg.ByName("GPIO4")
    in_read := gpioreg.ByName("GPIO16")

    out_command1.Out(gpio.Low)
    out_command2.Out(gpio.Low)
    out_read.Out(gpio.High)
    out_write.Out(gpio.High)

    for i := 0; i < 4096; i++ {
      bit3 := gpio.Low
      bit2 := gpio.Low
      bit1 := gpio.Low
      bit0 := gpio.Low
      if ((i & 8) >> 3) == 1 {
        bit3 = gpio.High
      }
      out_bit3.Out(bit3)
      
      if ((i & 4) >> 2) == 1 {
        bit2 = gpio.High
      }
      out_bit2.Out(bit2)

      if ((i & 2) >> 1) == 1 {
        bit1 = gpio.High
      }
      out_bit1.Out(bit1)

      if (i & 1) == 1 {
        bit0 = gpio.High
      }
      out_bit0.Out(bit0)

      out_write.Out(gpio.Low)

      for in_read.Read() == gpio.Low {
        in_read.WaitForEdge(-1)
      }

      out_write.Out(gpio.High)

      for in_read.Read() == gpio.High {
        in_read.WaitForEdge(-1)
      }
    }

    elapsedTime := time.Since(startTime)

    fmt.Printf("Sent 2 KiB in %s", elapsedTime)
}
