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
 lda #$1f ;set all flags high and page 1 of EPROM
PageJump:
 sta OutputFlags,x
 jmp Start ;this jump is only called if coming in from PageJump with A=$0f

;entry points for ProDOS
DriverEntry:
 lda #$0f ;set all flags high and page 0 of EPROM
 sta OutputFlags,x
 jmp Driver

;load first two blocks and execute to boot
Start:
 lda     #$01    ;set read command
 sta     Command
 
 lda     #$00    ;block 0
 sta     BlockLo
 sta     BlockHi
 sta     BufferLo   ;buffer at $800
 lda     #$08
 sta     BufferHi
 jsr     Driver  ;get the block

 jmp     $801    ;execute the block

; ProDOS Driver code
; First check that this is the right drive
Driver: 
 ldx Unit
 lda Command; Check which command is being requested
 beq GetStatus ;0 = Status command
 cmp #ReadBlockCommand
 beq ReadBlock
 cmp #WriteBlockCommand
 beq WriteBlock
 sec ;set carry as we don't support any other commands
 lda #$53 ;Invalid parameter error
 rts

; ProDOS Status Command Handler
GetStatus: 
 ldx #$ff ;low byte number of blocks 
 ldy #$ff ;high byte number of blocks
 lda #$0 ;zero accumulator and clear carry for success
 clc
 rts

; ProDOS Read Block Command
ReadBlock: 
 ldy #$00 ;Get the current time on each block read for now
 lda #GetTimeCommand
 jsr SendByte
getTimeByte:
 jsr GetByte
 sta $bf90,y
 iny
 cpy #$04
 bne getTimeByte
 lda #ReadBlockCommand ;read the block after setting the clock
 jsr SendByte
 lda BlockLo
 jsr SendByte
 lda BlockHi
 jsr SendByte
 ldy #$0
 jsr read256
 inc BufferHi
 jsr read256
 dec BufferHi
 lda #$0 ;zero accumulator and clear carry for success
 clc
 rts

read256: 
 jsr GetByte
 sta (BufferLo),y
 iny
 bne read256
 rts 
 
; ProDOS Write Block Command
WriteBlock:
 lda #WriteBlockCommand
 jsr SendByte
 lda BlockLo
 jsr SendByte
 lda BlockHi
 jsr SendByte
 ldy #$0
 jsr write256
 inc BufferHi
 jsr write256
 dec BufferHi
 lda #$0 ;zero accumulator and clear carry for success
 clc
 rts

write256:
 lda (BufferLo),y
 jsr SendByte
 iny
 bne write256
 rts

SendByte:
 pha 
waitWrite: 
 lda InputFlags,x
 rol
 rol 
 bcs waitWrite
 pla
 sta OutputByte,x
 lda #$0e ; set bit 0 low to indicate write started
 sta OutputFlags,x 
finishWrite:
 lda InputFlags,x
 rol
 rol
 bcc finishWrite
 lda #$0f
 sta OutputFlags,x
 rts

GetByte:
 lda #$0d ;set read flag low
 sta OutputFlags,x
waitRead:
 lda InputFlags,x
 rol
 bcs waitRead
 lda InputByte,x
 pha
 lda #$0f ;set all flags high
 sta OutputFlags,x
finishRead:
 lda InputFlags,x
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

