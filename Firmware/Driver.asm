; ProDOS Global Page
Device5S1 = $bf1a ;POINTER FOR SLOT 5 DRIVE 1 DRIVER
DeviceCount = $bf31 ;DEVICE COUNT -1
DeviceList = $bf32 ;DEVICE LIST

; ProDOS Zero Page
Command = $42 ;ProDOS Command
Unit = $43 ;ProDOS unit (SDDD0000)
BufferLo = $44
BufferHi = $45
BlockLo = $46
BlockHi = $47

; ProDOS Error Codes
IOError = $27
NoDevice = $28
WriteProtect = $2B

SlotDrive = $50
InputByte = $c0de
OutputByte = $c0dd
ReadBlockCommand = $01
WriteBlockCommand = $02
NibbleStorage = $1d

 .org $1000

; Register the driver with ProDOS
 lda #<Driver
 sta Device5S1
 lda #>Driver
 sta Device5S1+1
; Add the drive to the device list
 inc DeviceCount
 lda DeviceCount
 ldy #SlotDrive
 sta DeviceList,y
 rts

; ProDOS Driver code
; First check that this is the right drive
Driver: 
 lda Unit
 cmp #SlotDrive
 beq DoCommand ;correct device, so proceed
 sec ;set carry as ProDOS treats this as an error
 lda #NoDevice ;put the error code in accumulator for ProDOS
 rts

; Check which command is being requested
DoCommand: 
 lda Command
 beq GetStatus ;0 = Status command
 cmp #ReadBlockCommand
 beq ReadBlock
 cmp #WriteBlockCommand
 beq WriteBlock
 sec ;set carry as we don't support any other commands
 lda #$53 ;Invalid parameter error
 rts

; ProDOS Status Command Handler
GetStatus: 
 ldx #$ff ;low byte number of blocks 
 ldy #$ff ;high byte number of blocks
 lda #$0 ;zero accumulator and clear carry for success
 clc
 rts

; ProDOS Read Block Command
ReadBlock: 
 lda #ReadBlockCommand
 jsr SendByte
 lda BlockLo
 jsr SendByte
 lda BlockHi
 jsr SendByte
 ldy #$0
 jsr read256
 inc BufferHi
 jsr read256
 dec BufferHi
 lda #$0 ;zero accumulator and clear carry for success
 clc
 rts

read256: 
 jsr GetByte
 sta (BufferLo),y
 iny
 bne read256
 rts 
 
; ProDOS Write Block Command
WriteBlock:
 lda #WriteBlockCommand
 jsr SendByte
 lda BlockLo
 jsr SendByte
 lda BlockHi
 jsr SendByte
 ldy #$0
 jsr write256
 inc BufferHi
 jsr write256
 dec BufferHi
 lda #$0 ;zero accumulator and clear carry for success
 clc
 rts

write256:
 lda (BufferLo),y
 jsr SendByte
 iny
 bne write256
 rts

SendByte:
 pha
 lsr
 lsr
 lsr
 lsr
 jsr SendNibble
 pla
 jsr SendNibble
 rts

SendNibble:
 and #$0F
 ora #$70 ;Write bit low
 pha 
waitWrite: 
 lda InputByte
 asl ;Second highest bit goes low when ready
 bmi waitWrite
 pla
 sta OutputByte
finishWrite:
 lda InputByte
 asl
 bpl finishWrite
 lda #$FF
 sta OutputByte
 rts

GetByte:
 jsr GetNibble
 asl
 asl
 asl
 asl
 sta NibbleStorage 
 jsr GetNibble
 and #$0f
 ora NibbleStorage
 lda NibbleStorage
 rts

GetNibble:
 lda #$b0 ;set read flag low
 sta OutputByte
waitRead:
 lda InputByte
 bmi waitRead
 and #$f0 ;set all flags high
 sta OutputByte
 rts
 
