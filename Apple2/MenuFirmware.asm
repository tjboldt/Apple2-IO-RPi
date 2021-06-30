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
 jsr $fc58	;clear screen and show menu options
 ldy #$00
PrintString:
 lda Text,y
 beq WaitForRPi
 ora #$80
 jsr $fded
 iny
 bne PrintString

WaitForRPi:
 lda InputFlags
 rol
 bcs OK
 lda #$ff
 jsr $fca8
 lda #$ae
 jsr $fded
 jmp WaitForRPi

OK:
 jsr $fc58 ;clear screen

 lda #MenuCommand ;request menu text from RPi
 jsr SendByte

DumpOutput:
 jsr GetByte
 cmp #$00
 beq GetChar
 jsr $fded
 clc
 bcc DumpOutput

GetChar: 
 jsr $fd0c
 sec	;subtract ascii "1" to get 0 - 3 from "1" to "4"
 sbc #$b1
 asl	;put in top nibble as EPROM page 
 asl
 asl
 asl
 ora #$0f ;set all flags high
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

.repeat	183-<end
.byte 0
.endrepeat

Text:
.byte	"Apple2-IO-RPi",$8d
.byte	"(c)2020-2021 Terence J. Boldt",$8d
.byte   $8d
.byte	"Waiting for RPi...",$00

.byte <GetByte  ;all firmware pages have pointer to GetByte routine here
.byte >GetByte
.byte <SendByte ;all firmware pages have pointer to SendByte routine here
.byte >SendByte

.byte      0,0     ;0000 blocks = check status
.byte      7       ;bit set(0=status 1=read 2=write) unset(3=format, 4/5=number of volumes, 6=interruptable, 7=removable)
.byte     DriverEntry&$00FF ;low byte of entry

