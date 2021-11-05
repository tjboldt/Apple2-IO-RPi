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

ReadBlockCommand = $01
WriteBlockCommand = $02
GetTimeCommand = $03
ChangeDriveCommand = $04
ExecCommand = $05
LoadFileCommand = $06
SaveFileCommand = $07
MenuCommand = $08

InputString = $fd6a
StringBuffer = $0200
PrintChar = $fded
Keyboard = $c000
ClearKeyboard = $c010
Wait = $fca8
PromptChar = $33

ESC = $9b

 .org $2000
Start:
 lda PromptChar
 sta OldPromptChar
 lda #'$'|$80
 sta PromptChar
 lda #ExecCommand
 jsr SendByte
 ldx #$00
sendHelpCommand:
 lda HelpCommand,x
 cmp #$00
 beq sendHelpCommandEnd
 jsr SendByte
 inx
 bpl sendHelpCommand
sendHelpCommandEnd:
 lda #$00
 jsr SendByte 
 bit ClearKeyboard
 jsr DumpOutput

Prompt:
 lda #ExecCommand
 jsr SendByte
 ldx #$00
sendPromptCommand:
 lda PromptCommand,x
 cmp #$00
 beq sendPromptCommandEnd
 jsr SendByte
 inx
 bpl sendPromptCommand
sendPromptCommandEnd:
 lda #$00
 jsr SendByte 
 bit ClearKeyboard
 jsr DumpOutput

; get input
 jsr InputString
; check for "exit"
 lda StringBuffer
 cmp #'e'|$80
 bne Execute
 lda StringBuffer+1
 cmp #'x'|$80
 bne Execute
 lda StringBuffer+2
 cmp #'i'|$80
 bne Execute
 lda StringBuffer+3
 cmp #'t'|$80
 bne Execute
 lda OldPromptChar
 sta PromptChar
 rts
Execute: 
 bit ClearKeyboard
 lda #ExecCommand
 jsr SendByte
 ldy #$00
sendInput:
 lda $0200,y
 cmp #$8d
 beq sendNullTerminator
 and #$7f
 jsr SendByte
 iny
 bne sendInput
sendNullTerminator:
 lda #$00
 jsr SendByte
 jsr DumpOutput
 jmp Prompt

DumpOutput:
 jsr GetByte
 bcs skipOutput
 cmp #$00
 beq endOutput
 cmp #ESC
 beq escapeSequence
 jsr PrintChar
skipOutput:
 bit Keyboard ;check for keypress
 bpl DumpOutput ;keep dumping output if no keypress
 lda Keyboard ;send keypress to RPi
 ;jsr PrintChar
 and #$7f
 jsr SendByte
 bit ClearKeyboard
 clc
 bcc DumpOutput
endOutput:
 rts
escapeSequence:
 jsr ParseEscape
 clc
 bcc DumpOutput

ParseEscape:
 jsr GetByte ; expect first byte after ESC to be '['
 cmp #'['|$80
 beq endParse
checkLetter:
 jsr GetByte ; loop until there is a letter
 cmp #$C1
 bcc checkLetter
endParse:
 rts

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
 lda #$1d ;set read flag low
 sta OutputFlags
waitRead:
 lda InputFlags
 rol
 bcc readByte
 bit Keyboard ;keypress will abort waiting to read
 bpl waitRead
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

HelpCommand:
.byte	"a2help",$00
PromptCommand:
.byte   "a2prompt",$00
OldPromptChar:
.byte   "]"