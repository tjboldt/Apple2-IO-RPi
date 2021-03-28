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
 lda $33
 pha
 lda #$a4
 sta $33
GetCommand:
 jsr $fd6a
 lda $0200
 cmp #$8d ;stop when return found
 beq ExitApp
 jsr DumpOutput
 clc
 bcc GetCommand
ExitApp:
 pla
 sta $33
 lda #$00
 sty $34
 rts

DumpOutput:
 ldx #$50
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
getOutput:
 jsr GetByte
 cmp #$00
 beq endOutput
 jsr $fded
 clc
 bcc getOutput
endOutput:
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
 lda #$0e ; set bit 0 low to indicate write started
 sta OutputFlags 
finishWrite:
 lda InputFlags
 rol
 rol
 bcc finishWrite
 lda #$0f
 sta OutputFlags
 rts

GetByte:
 lda #$0d ;set read flag low
 sta OutputFlags
waitRead:
 lda InputFlags
 rol
 bcs waitRead
 lda InputByte
 pha
 lda #$0f ;set all flags high
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

;Translation of code written in mini-assembler on Apple //e
;Currently only works if card is in slot 5

 .org $1000

