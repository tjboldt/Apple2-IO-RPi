; Copyright Terence J. Boldt (c)2020-2022
; Use of this source code is governed by an MIT
; license that can be found in the LICENSE file.

; This file contains the source for the firmware
; that displays the copyright message on boot
; and checks for the RPi status to be ready before
; attempting to boot

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

ResetCommand = $00
ReadBlockCommand = $01
WriteBlockCommand = $02
GetTimeCommand = $03
ChangeDriveCommand = $04
ExecCommand = $05
LoadFileCommand = $06
SaveFileCommand = $07
MenuCommand = $08

Wait = $fca8
PrintChar = $fded
Home = $fc58
ReadChar = $fd0c
BasCalc = $fbc1
htab = $24
vtab = $25
BasL = $28
htab80 = $057b

 .org SLOT*$100 + $C000
;ID bytes for booting and drive detection
 cpx #$20    ;ID bytes for ProDOS and the
 cpx #$00    ; Apple Autostart ROM
 cpx #$03    ;

 ldx #SLOT*$10
 stx $2b
 stx Unit 

;force EPROM to second page on boot
 lda #$3f ;set all flags high and page 1 of EPROMi
PageJump:
 sta OutputFlags
 jmp Start

;entry points for ProDOS
DriverEntry:
 lda #$0f ;set all flags high and page 0 of EPROM
 sta OutputFlags
;since the firmware page changes to 0, this falls through to the driver

Start:
 lda #$f0	;restore COUT after PR# called
 sta $36
 lda #$fd
 sta $37
 ;jsr Home	;clear screen and show menu options
 lda #$10
 sta vtab
 jsr BasCalc
 ldy #$00
PrintString:
 lda Text,y
 beq WaitForRPi
 ora #$80
 jsr PrintChar
 iny
 bne PrintString

WaitForRPi:
 lda InputFlags
 rol
 bcs Reset
 lda #$ff
 jsr Wait
 lda #'.'+$80
 jsr PrintChar
 jmp WaitForRPi

Reset:
 lda #'_'|$80
 jsr PrintChar
 lda #ResetCommand
 jsr SendByte
 lda #$88
 jsr PrintChar
 lda #'.'|$80
 jsr PrintChar
 jsr GetByte
 beq Ok
 jmp Reset

Ok:
 lda #$8D
 jsr PrintChar
 lda #'O'|$80
 jsr PrintChar
 lda #'K'|$80
 jsr PrintChar

Boot:
 lda #$0f 
 jmp PageJump

SendByte:
 pha 
waitWrite: 
 lda InputFlags
 rol
 rol 
 bcs waitWrite
 pla
 sta OutputByte
 lda #$3e ; set bit 0 low to indicate write started
 sta OutputFlags 
finishWrite:
 lda InputFlags
 rol
 rol
 bcc finishWrite
 lda #$3f
 sta OutputFlags
 rts

GetByte:
 lda #$3d ;set read flag low
 sta OutputFlags
waitRead:
 lda InputFlags
 rol
 bcs waitRead
 lda InputByte
 pha
 lda #$3f ;set all flags high
 sta OutputFlags
finishRead:
 lda InputFlags
 rol
 bcc finishRead
 pla
end:
 rts

; NOTE: The text below exactly fills the remaining 256 bytes of firmware
Text:
.byte	"Apple2-IO-RPi",$8d
.byte	"(c)2020-2022 Terence J. Boldt",$8d
.byte   $8d
.byte	"Waiting for RPi FW:000F...",$00

.byte      0,0     ;0000 blocks = check status
.byte      7       ;bit set(0=status 1=read 2=write) unset(3=format, 4/5=number of volumes, 6=interruptable, 7=removable)
.byte     DriverEntry&$00FF ;low byte of entry

