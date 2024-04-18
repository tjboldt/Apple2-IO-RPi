; Copyright Terence J. Boldt (c)2021-2024
; Use of this source code is governed by an MIT
; license that can be found in the LICENSE file.

; This file contains the source for the SHELL
; application that runs on the Apple II to talk
; to the Raspberry Pi

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

InputByte = $c08e
OutputByte = $c08d
InputFlags = $c08b
OutputFlags = $c087

ResetCommand = $00
ReadBlockCommand = $01
WriteBlockCommand = $02
GetTimeCommand = $03
ChangeDriveCommand = $04
ExecCommand = $05
LoadFileCommand = $06
SaveFileCommand = $07
MenuCommand = $08
ShellCommand = $09

InputString = $fd6a
StringBuffer = $0200
PrintChar = $fded
PrintHex = $FDE3
PrintByte = $FDDA

Keyboard = $c000
ClearKeyboard = $c010
Home = $fc58
Wait = $fca8
PromptChar = $33
Read80Col = $c01f
TextPage1 = $c054
TextPage2 = $c055

htab = $24
vtab = $25
BasL = $28
htab80 = $057b
BasCalc = $fbc1

LastChar = $06
SlotL = $fe
SlotH = $ff
ESC = $9b

 .org $2000
; Find Apple2-IO-RPi card
 ldx #$06
 ldy #$09
CheckIdBytes:
 lda $C700,y ; !! Self modifying code
 cmp IdBytes,y
 bne NextCard
 dey
 bne CheckIdBytes
 jmp FoundCard
NextCard:
 dec CheckIdBytes+2
 ldy #$09
 dex
 bne CheckIdBytes
CardNotFound:
 ldy #$00
PrintCardNotFound:
 lda TextCardNotFound,y
 beq Failed
 ora #$80
 jsr PrintChar
 iny
 bne PrintCardNotFound
Failed:
 rts
FoundCard:
 ldy #$00
PrintCardFound:
 lda TextCardFound,y
 beq PrintCardNumber
 ora #$80
 jsr PrintChar
 iny
 bne PrintCardFound
PrintCardNumber:
 inx
 txa
 jsr PrintHex
 lda #$8D
 jsr PrintChar 

SetOffsetForCard:
 txa
 asl
 asl
 asl
 asl
 tax

Init:
 lda #$8d
 jsr $c300 ; force 80 columns

 ldy #$00
PrintString:
 lda Text,y
 beq Start
 ora #$80
 jsr PrintChar
 iny
 bne PrintString

Start:
 lda LastChar
 pha
 bit ClearKeyboard
; lda #ResetCommand
; jsr SendByte
 lda #ShellCommand
 jsr SendByte
 jsr DumpOutput
 pla
 sta LastChar
 rts

DumpOutput:
 jsr GetByte
 cmp #$00
 beq endOutput
 pha
 jsr ClearCursor
 pla
 cmp #'H'
 beq setColumn
 cmp #'V'
 beq setRow
 cmp #'C'
 beq clearScreen
 cmp #'T'
 beq setTop
 cmp #'B'
 beq setBottom
 cmp #'U'
 beq moveUp
 jsr PrintChar
 jsr SetCursor
 jmp DumpOutput
endOutput:
 rts
clearScreen:
 jsr Home
 jsr SetCursor
 jmp DumpOutput
setColumn:
 jsr GetByte
 sta htab
 sta htab80
 jsr SetCursor
 jmp DumpOutput
setRow:
 jsr GetByte
 sta vtab
 jsr BasCalc
 jsr SetCursor
 jmp DumpOutput
setTop:
 jsr GetByte
 sta $22
 jmp DumpOutput
setBottom:
 jsr GetByte
 sta $23
 jmp DumpOutput
moveUp:
 dec vtab
 lda vtab
 jsr BasCalc
 jsr SetCursor
 jmp DumpOutput

SendByte:
 pha 
waitWrite: 
 lda InputFlags,x
 rol
 rol 
 bcs waitWrite
 pla
 sta OutputByte,x
.if HW_TYPE = 0
 lda #$1e ; set bit 0 low to indicate write started
 sta OutputFlags,x 
finishWrite:
 lda InputFlags,x
 rol
 rol
 bcc finishWrite
 lda #$1f
 sta OutputFlags,x
.endif
 rts

GetByte:
.if HW_TYPE = 0
 lda #$1d ;set read flag low
 sta OutputFlags,x
.endif
waitRead:
 lda InputFlags,x
 rol
 bcc readByte
 bit Keyboard ;keypress will abort waiting to read
 bpl waitRead
keyPressed:
 lda InputFlags,x
 rol
 rol
 bcs keyPressed
 lda Keyboard ;send keypress to RPi
 and #$7f
 sta OutputByte,x
 bit ClearKeyboard
.if HW_TYPE = 0
 lda #$1c ;set write flag low too
 sta OutputFlags,x
finishKeyPress:
 lda InputFlags,x
 rol
 rol
 bcc finishKeyPress
 lda #$1d ;set flags back for reading
 sta OutputFlags,x
.endif
 jmp waitRead
readByte:
 lda InputByte,x
.if HW_TYPE = 0
 pha
 lda #$1f ;set all flags high
 sta OutputFlags,x
finishRead:
 lda InputFlags,x
 rol
 bcc finishRead
 pla
.endif
 clc ;success
end:
 rts

SetCursor:
 lda htab80 ;get horizontal location / 2
 lsr
 tay
 lda TextPage2
 bcc setChar
 lda TextPage1
setChar:
 lda (BasL),y
 sta LastChar ; save so ClearCursor will pick it up
 cmp #$e0
 bpl lowerCase
 cmp #$c0
 bpl upperCase
 cmp #$a0
 bpl symbol
 cmp #$80
 bpl noop
symbol:
lowerCase:
invert:
 eor #$80
 jmp storeChar
upperCase:
 and #$1f
 jmp storeChar
noop:
storeChar:
 sta (BasL),y
 lda TextPage1
 rts

ClearCursor:
 lda htab80 ;get horizontal location / 2
 lsr
 tay
 lda TextPage2
 bcc restoreChar
 lda TextPage1
restoreChar:
 lda LastChar 
 sta (BasL),y
 lda TextPage1
 rts

Text:
.if HW_TYPE = 0
.byte	"Apple2-IO-RPi Shell Version 0010 (classic)",$8d
.else
.byte	"Apple2-IO-RPi Shell Version 8010 (pico)",$8d
.endif
.byte	"(c)2020-2024 Terence J. Boldt",$8d
.byte   $8d
.byte   $00

IdBytes:
.byte  $E0,$20,$E0,$00,$E0,$03,$E0,$3C,$A9,$3F

TextCardFound:
.byte "Found Apple2-IO-RPi in slot ",$00

TextCardNotFound:
.byte "Apple2-IO-RPi not found",$8D,$00
