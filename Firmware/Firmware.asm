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

SlotDrive = $50
InputByte = $c08e
OutputByte = $c08d
ReadBlockCommand = $01
WriteBlockCommand = $02
NibbleStorage = $1d

 .org STARTSLOT
;ID bytes for booting and drive detection
 cpx #$20    ;ID bytes for ProDOS and the
 cpx #$00    ; Apple Autostart ROM
 cpx #$03    ;
 cpx #$3C    ;this one for older II's

;load first two blocks and execute to boot
Boot:
 lda     #$01    ;set read command
 sta     Command
 
 jsr $ff58
 tsx
 lda $0100,x
 asl
 asl
 asl
 asl
 sta $2b
 sta Unit 
 tax

 lda     #$00    ;block 0
 sta     BlockLo
 sta     BlockHi
 sta     BufferLo   ;buffer at $800
 lda     #$08
 sta     BufferHi
 jsr     Driver  ;get the block

; lda     #$00    ;block 1
; sta     BlockLo
; sta     BlockHi
; sta     BufferLo   ;buffer at $A00
; lda     #$0A
; sta     BufferHi
; jsr     Driver   ;get the block

; ldx     #sdrive ;set up for slot n
 jmp     $801    ;execute the block

;;
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
 lda #ReadBlockCommand
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
 lsr
 lsr
 lsr
 lsr
 jsr SendNibble
 pla
 jsr SendNibble
 rts

SendNibble:
 and #$0F
 ora #$70 ;Write bit low
 pha 
waitWrite: 
 lda InputByte,x
 asl ;Second highest bit goes low when ready
 bmi waitWrite
 pla
 sta OutputByte,x
finishWrite:
 lda InputByte,x
 asl
 bpl finishWrite
 lda #$FF
 sta OutputByte,x
 rts

GetByte:
 jsr GetNibble
 asl
 asl
 asl
 asl
 sta NibbleStorage 
 jsr GetNibble
 and #$0f
 ora NibbleStorage
 rts

GetNibble:
 lda #$b0 ;set read flag low
 sta OutputByte,x
waitRead:
 lda InputByte,x
 bmi waitRead
 ora #$f0 ;set all flags high
 sta OutputByte,x
 pha
finishRead:
 lda InputByte,x
 bpl finishRead
 pla
end:
 rts

.repeat	251-<end
.byte 0
.endrepeat
.byte      0,0     ;0000 blocks = check status
.byte      3       ;bit 0=read 1=status
.byte     Driver&$00FF ;low byte of entry
