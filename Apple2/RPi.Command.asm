; Copyright Terence J. Boldt (c)2021-2024
; Use of this source code is governed by an MIT
; license that can be found in the LICENSE file.

; This file contains the source for the RPI.COMMAND
; application that runs on the Apple II and extends
; ProDOS BASIC.SYSTEM to add the RPI command which
; allows commands to be executed on the Raspberry Pi

            .ORG  $300
INBUF      =  $200     ;GETLN input buffer.
WAIT       =  $FCA8    ;Monitor wait routine.
BELL       =  $FF3A    ;Monitor bell routine.
EXTRNCMD   =  $BE06    ;External cmd JMP vector.
XTRNADDR   =  $BE50    ;Ext cmd implementation addr.
XLEN       =  $BE52    ;length of command string-1.
XCNUM      =  $BE53    ;CI cmd no. (ext cmd - 0).
PBITS      =  $BE54    ;Command parameter bits.
XRETURN    =  $FF58    ;Known RTS instruction.
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
MenuCommand = $08

InputString = $fd67
PrintChar = $fded
Keyboard = $c000
ClearKeyboard = $c010
Wait = $fca8

LastChar = $06
SlotL = $fe
SlotH = $ff
ESC = $9b

;macro for string with high-bit set
.macro aschi str
.repeat .strlen (str), c
.byte .strat (str, c) | $80
.endrep
.endmacro

 .org $2000
 ldx #$07 ; start at slot 7
DetectSlot:
 ldy #$00
 lda #$fc
 sta SlotL
 txa
 ora #$c0
 sta SlotH 
 lda (SlotL),y
 bne nextSlot
 iny
 lda (SlotL),y
 bne nextSlot
 iny 
 lda (SlotL),y
 cmp #$17
 bne nextSlot
 iny
 lda (SlotL),y
 cmp #$14
 bne nextSlot
 txa
 asl
 asl
 asl
 asl
 tax
 clc
 bcc Start
nextSlot:
 dex
 bne DetectSlot
 rts 
Start:
 stx slotx + $1e01 ;set the slot for the driver
  ldy #$00
PrintString:
 lda Text,y
 beq copyDriver
 ora #$80
 jsr PrintChar
 iny
 bne PrintString
copyDriver: 
 ldy #$00
copyDriverByte:
 lda $2100,y
 sta $0300,y
 iny
 cpy #$e6
 bne copyDriverByte
 ;
 ; FIRST SAVE THE EXTERNAL COMMAND ADDRESS SO YOU WON'T
 ; DISCONNECT ANY PREVIOUSLY CONNECTED COMMAND.
 ;
            LDA  EXTRNCMD+1
            STA  NXTCMD
            LDA  EXTRNCMD+2
            STA  NXTCMD+1
 ;
            LDA  #<RPI      ;Install the address of our
            STA  EXTRNCMD+1  ; command handler in the
            LDA  #>RPI      ; external command JMP
            STA  EXTRNCMD+2  ; vector.

 lda #ExecCommand
 jsr SendByte
 ldy #$00
nextCommandByte:
 lda a2help, y
 beq finishCommand
 jsr SendByte
 iny
 jmp nextCommandByte
finishCommand:
 lda #$00
 jsr SendByte
showVersion:
 jsr GetByte
 cmp #$00
 beq FinishDriver
 jsr PrintChar
 jmp showVersion
FinishDriver:
 rts

a2help:
 .byte "a2help", $00

Text:
 aschi "RPI command version: 000E"
 .byte $8d
end:
 .byte $00 

.repeat	255-<end
.byte 0
.endrepeat

.org $0300
 RPI:       LDX  #0          ;Check for our command.
 NXTCHR:     LDA  INBUF,X     ;Get first character.
            ora  #$20        ;Make it lower case
            CMP  CMD,X       ;Does it match?
            BNE  NOTOURS     ;No, back to CI.
            INX              ;Next character
            CPX  #CMDLEN     ;All characters yet?
            BNE  NXTCHR      ;No, read next one.
 ;
            LDA  #CMDLEN-1   ;Our cmd! Put cmd length-1
            ;lda #$8d
            ;sta $02ff
            ;lda #$fe
            STA  XLEN        ; in CI global XLEN.
            LDA  #<XRETURN   ;Point XTRNADDR to a known
            STA  XTRNADDR    ; RTS since we'll handle
            LDA  #>XRETURN   ; at the time we intercept

            STA  XTRNADDR+1  ; our command.
            LDA  #0          ;Mark the cmd number as
            STA  XCNUM       ; zero (external).
            STA  PBITS       ;And indicate no parameters
            STA  PBITS+1     ; to be parsed.
            lda #$8d
            jsr $fded
slotx:      ldx #$70        ; set x to slot # in high nibble
            clc
            bcc SendCommand
  ;
 NOTOURS:    SEC              ; ALWAYS SET CARRY IF NOT YOUR
            JMP  (NXTCMD)    ; CMD AND LET NEXT COMMAND TRY
 ;                           ; TO CLAIM IT.

SendCommand:
 bit ClearKeyboard
 lda #$05 ;send command 5 = exec
 jsr SendByte
 ldy #$03 ;skip over "RPI"
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
 clc
 jmp (NXTCMD) 

HelpCommand:
 .byte "a2help",$00

SendByte:
 pha 
waitWrite: 
 lda InputFlags,x
 rol
 rol 
 bcs waitWrite
 pla
 sta OutputByte,x
 lda #$1e ; set bit 0 low to indicate write started
 sta OutputFlags,x 
finishWrite:
 lda InputFlags,x
 rol
 rol
 bcc finishWrite
 lda #$1f
 sta OutputFlags,x
 rts

GetByte:
 lda #$1d ;set read flag low
 sta OutputFlags,x
waitRead:
 lda InputFlags,x
 rol
 bcs waitRead
readByte:
 lda InputByte,x
 pha
 lda #$1f ;set all flags high
 sta OutputFlags,x
finishRead:
 lda InputFlags,x
 rol
 bcc finishRead
 pla
 clc ;success
 rts


CMD:   aschi   "rpi"
 CMDLEN     =  3       ;Our command length
 ;
 NXTCMD:    .byte    0,0           ; STORE THE NEXT EXT CMD'S
                             ; ADDRESS HERE.


