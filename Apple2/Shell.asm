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

; hard code slot to 7 for now, will make it auto-detect later
SLOT = 7

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
ShellCommand = $09

InputString = $fd6a
StringBuffer = $0200
PrintChar = $fded
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

ESC = $9b

 .org $2000
Start:
 jsr $c300 ; force 80 columns
 lda LastChar
 pha
 bit ClearKeyboard
 lda #ResetCommand
 jsr SendByte
 lda #ShellCommand
 jsr SendByte
 jsr DumpOutput
 pla
 sta LastChar
 rts

DumpOutput:
 jsr GetByte
 bcs checkInput
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
checkInput:
 bit Keyboard ;check for keypress
 bpl DumpOutput ;keep dumping output if no keypress
 lda Keyboard ;send keypress to RPi
 and #$7f
 jsr SendByte
 bit ClearKeyboard
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
 lda InputFlags
 rol
 rol 
 bcs waitWrite
 pla
 sta OutputByte
 lda #$1e ; set bit 0 low to indicate write started
 sta OutputFlags 
finishWrite:
 lda InputFlags
 rol
 rol
 bcc finishWrite
 lda #$1f
 sta OutputFlags
 rts

GetByte:
 bit Keyboard ; skip byte read if key pressed
 bcc keyPressed
 lda #$1d ;set read flag low
 sta OutputFlags
waitRead:
 lda InputFlags
 rol
 bcc readByte
 bit Keyboard ;keypress will abort waiting to read
 bpl waitRead
keyPressed:
 lda #$1f ;set all flags high and exit
 sta OutputFlags
 sec ;failure
 rts 
readByte:
 lda InputByte
 pha
 lda #$1f ;set all flags high
 sta OutputFlags
finishRead:
 lda InputFlags
 rol
 bcc finishRead
 pla
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
