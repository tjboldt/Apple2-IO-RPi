# Apple2-IO-RPi
Apple II expansion card using a Raspberry Pi for I/O

![Image of Board](/Hardware/Apple2IORPi.jpg)

## Project Status
This is a very early stage project. Currently one partially completed hand wired physical prototype exists that sends data unidirectional from the Raspberry Pi Zero W to the Apple II via 8 bit parallel data bus. The PCB layout has not been physically tested and no firmware or software exists yet.

## Purpose
The purpose of this project is to provide I/O for an Apple II series 8 bit computer via a Raspberry Pi Zero W which is powered by the Apple II expansion bus. Initially this would be storage via virtual ProDOS compatible drive. Next might be adding virtual serial card support over wifi. Future enhancements could use the RPi for more complex processing as per request from the Apple II. For example, the Apple II could request a web page or application and the RPi could calculate this in Apple II hi-res graphics mode and send the image data back to the II for display purposes.

## Similar Project
If you're looking for complete hardware design or prefer having Apple II peripherals control a control a Raspberry Pi rather than simply using the Raspberry Pi to provide storage and network access to the Apple II, have a look at David Schmenk's excellent [Apple2Pi](https://github.com/dschmenk/apple2pi) project.  
