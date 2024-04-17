; Copyright Terence J. Boldt (c)2020-2024
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
 cpx #$20    ;ID bytes for ProDOS
 cpx #$00    ;
 cpx #$03    ;
 cpx #$3C    ;ID byte for Autostart ROM

 lda #$3f ;set all flags high and page 3 of EPROM for menu
 jmp PageJump

;The following bytes must exist in this location for Pascal communications
;as they are standard pointers for entry points
.byte     CommInit&$00FF ;low byte of rts for init of Pascal comms
.byte     $43 ; low byte of read for Pascal comms
.byte     $49 ; low byte of write for Pascal comms
.byte     $4F ; low byte of rts for status of Pascal comms

CommInit:
 lda #$0f ; set all flags high and page 0 for comm driver
 sta OutputFlags
 ldx #$00 ; set error code to 0 for success
 rts

PageJump:
 sta OutputFlags
 jmp Start ;this jump is only called if coming in from PageJump with A=$0f

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
 lda #$02
 sta vtab
 jsr BasCalc
 ldy #$00
 sty htab
PrintString:
 lda Text,y
 beq WaitForRPi
 ora #$80
 jsr PrintChar
 iny
 bne PrintString

.if HW_TYPE = 0

WaitForRPi:
 bit InputFlags
 bmi Reset
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

.else

WaitForRPi:
@1:
 bit InputFlags
 bmi @2 
 lda InputByte
 jmp @1
@2:
 bit InputFlags
 bpl @4
 bvs @3
 lda #ResetCommand
 sta OutputByte
@3:
 lda #$ff
 jsr Wait
 lda #'.'+$80
 jsr PrintChar
 jmp @2
@4:
 lda #$ff
 jsr Wait
@5:
 bit InputFlags
 bmi Ok 
 lda InputByte
 jmp @5

.endif

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
 bit InputFlags
 bvs SendByte
 sta OutputByte
.if HW_TYPE = 0
 lda #$3e ; set bit 0 low to indicate write started
 sta OutputFlags 
finishWrite:
 bit InputFlags
 bvc finishWrite
 lda #$3f
 sta OutputFlags
.endif
 rts

GetByte:
.if HW_TYPE = 0
 ldx #$3d ;set read flag low
 stx OutputFlags
.endif
waitRead:
 bit InputFlags
 bmi waitRead
 lda InputByte
.if HW_TYPE = 0
 ldx #$3f ;set all flags high
 stx OutputFlags
finishRead:
 bit InputFlags
 bpl finishRead
.endif
 rts

Text:
.byte	"Apple2-IO-RPi",$8d
.byte	"(c)2020-2024 Terence J. Boldt",$8d
.byte   $8d
.if HW_TYPE = 0
.byte	"Waiting for RPi FW:0011"
.else
.byte   "Waiting for RPi FW:8011"
.endif
end:
.byte	$00

.repeat	251-<end
.byte 0
.endrepeat

.byte      0,0     ;0000 blocks = check status
.byte      7       ;bit set(0=status 1=read 2=write) unset(3=format, 4/5=number of volumes, 6=interruptable, 7=removable)
.byte     DriverEntry&$00FF ;low byte of entry
