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

InputString = $fd67
PrintChar = $fded

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
 jsr $c300 ;enable 80 columns
 lda #$05 ;execute command
 jsr SendByte
 ldy #$00
sendHelp:
 lda HelpCommand,y
 beq endSendHelp
 jsr SendByte
 iny
 bne sendHelp
endSendHelp:
 lda #$00
 jsr SendByte
 jsr DumpOutput

 lda $33
 pha
 lda #$a4
 sta $33
GetCommand:
 jsr InputString
 lda $0200
 cmp #$8d ;skip when return found
 beq GetCommand
 jsr SendCommand
 clc
 bcc GetCommand

SendCommand:
 lda #$05 ;send command 5 = exec
 jsr SendByte
 ldy #$00
getInput:
 lda $0200,y
 cmp #$8d
 beq sendNullTerminator
 and #$7f
 jsr SendByte
 iny
 bne getInput
sendNullTerminator:
 lda #$00
 jsr SendByte
DumpOutput:
 jsr GetByte
 cmp #$00
 beq endOutput
 jsr PrintChar
 clc
 bcc DumpOutput
endOutput:
 rts

HelpCommand:
 .byte "a2help",$00

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
 bcs waitRead
 lda InputByte
 pha
 lda #$1f ;set all flags high
 sta OutputFlags
finishRead:
 lda InputFlags
 rol
 bcc finishRead
 pla
end:
 rts

.repeat	247-<end
.byte 0
.endrepeat

.byte <GetByte  ;all firmware pages have pointer to GetByte routine here
.byte >GetByte
.byte <SendByte ;all firmware pages have pointer to SendByte routine here
.byte >SendByte

.byte      0,0     ;0000 blocks = check status
.byte      7       ;bit set(0=status 1=read 2=write) unset(3=format, 4/5=number of volumes, 6=interruptable, 7=removable)
.byte     DriverEntry&$00FF ;low byte of entry
