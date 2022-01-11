; Copyright Terence J. Boldt (c)2020-2022
; Use of this source code is governed by an MIT
; license that can be found in the LICENSE file.

; This file contains the source for the firmware
; that was formerly used to copy files from RPi
; to Apple II RAM

;ProDOS Zero Page
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

InputByte = $c08e+SLOT*$10
OutputByte = $c08d+SLOT*$10
InputFlags = $c08b+SLOT*$10
OutputFlags = $c087+SLOT*$10

ReadBlockCommand = $01
WriteBlockCommand = $02
GetTimeCommand = $03
ChangeDriveCommand = $04
ExecCommand = $05
LoadFileCommand = $06
SaveFileCommand = $07
MenuCommand = $08

InputString = $fd67
Monitor = $ff59

 .org SLOT*$100 + $C000
;ID bytes for booting and drive detection
 cpx #$20    ;ID bytes for ProDOS and the
 cpx #$00    ; Apple Autostart ROM
 cpx #$03    ;

 ldx #SLOT*$10
 stx $2b
 stx Unit 

;force EPROM to second page on boot
 lda #$3f ;set all flags high and page 3 of EPROM for menu
PageJump:
 sta OutputFlags
 jmp Start ;this jump is only called if coming in from PageJump with A=$2f

;entry points for ProDOS
DriverEntry:
 lda #$0f ;set all flags high and page 0 of EPROM
 sta OutputFlags

Start:
 lda #$a4
 sta $33

GetFilename:
 jsr InputString

LoadFile:
 lda #$00
 sta BufferLo
 lda #$20
 sta BufferHi
 lda #$06 ; send command 6 = load
 jsr SendByte
 ldy #$00
sendFilename:
 lda $0200,y
 cmp #$8d
 beq sendNullTerminator
 and #$7f
 jsr SendByte
 iny
 bne sendFilename 
sendNullTerminator:
 lda #$00
 jsr SendByte
  
 jsr GetByte
 sta BlockLo ; not really a block, just using the memory space
 jsr GetByte
 sta BlockHi
NextPage:
 lda BlockHi
 beq ReadFinalPage
 ldy #$00
NextByte:
 jsr GetByte
 sta (BufferLo),y
 iny
 bne NextByte
 inc BufferHi
 dec BlockHi
 bne NextPage
ReadFinalPage:
 lda BlockLo
 beq ExitToMonitor
 ldy #$00
NextByteFinal:
 jsr GetByte
 sta (BufferLo),y
 iny
 cpy BlockLo
 bne NextByteFinal
ExitToMonitor:
 jsr Monitor 

SendByte:
 pha 
waitWrite: 
 lda InputFlags
 rol
 rol 
 bcs waitWrite
 pla
 sta OutputByte
 lda #$2e ; set bit 0 low to indicate write started
 sta OutputFlags 
finishWrite:
 lda InputFlags
 rol
 rol
 bcc finishWrite
 lda #$2f
 sta OutputFlags
 rts

GetByte:
 lda #$2d ;set read flag low
 sta OutputFlags
waitRead:
 lda InputFlags
 rol
 bcs waitRead
 lda InputByte
 pha
 lda #$2f ;set all flags high
 sta OutputFlags
finishRead:
 lda InputFlags
 rol
 bcc finishRead
 pla
end:
 rts

.repeat	251-<end
.byte 0
.endrepeat

.byte      0,0     ;0000 blocks = check status
.byte      7       ;bit set(0=status 1=read 2=write) unset(3=format, 4/5=number of volumes, 6=interruptable, 7=removable)
.byte     DriverEntry&$00FF ;low byte of entry

