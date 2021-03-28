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
 jsr $fc58
 ldy #$00
PrintString:
 lda Text,y
 beq GetChar
 ora #$80
 jsr $fded
 iny
 bne PrintString

GetChar: 
 jsr $fd1b
 sec	;subtract ascii "1" to get 0 - 3 from "1" to "4"
 sbc #$b1
 asl	;put in top nibble as EPROM page 
 asl
 asl
 asl
 ora #$0f ;set all flags high
 jmp PageJump

Text:

.byte	"Apple2-IO-RPi",$8d
.byte	"(c)2020-2021 Terence J. Boldt",$8d,$8d
.byte	"1. Boot",$8d
.byte	"2. File Access",$8d
.byte	"3. Command Line",$8d,$00

end:
 rts

.repeat	251-<end
.byte 0
.endrepeat
.byte      0,0     ;0000 blocks = check status
.byte      7       ;bit set(0=status 1=read 2=write) unset(3=format, 4/5=number of volumes, 6=interruptable, 7=removable)
.byte     DriverEntry&$00FF ;low byte of entry

