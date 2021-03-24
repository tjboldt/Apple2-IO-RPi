;Translation of code written in mini-assembler on Apple //e
;Currently only works if card is in slot 5

 .org $1000
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
 
 brk
 brk
 brk
 brk
 brk
 brk
 brk
 brk
 brk
 brk
 brk
 brk

DumpOutput:
 ldx #$50
 lda #$05 ;send command 5 = exec
 jsr $c5aa
 ldy #$00
getInput:
 lda $0200,y
 cmp #$8d
 beq sendNullTerminator
 and #$7f
 jsr $c5aa
 iny
 bne getInput
sendNullTerminator:
 lda #$00
 jsr $c5aa ;send byte
getOutput:
 jsr $c5c8 ;get byte
 cmp #$00
 beq endOutput
 jsr $fded
 clc
 bcc getOutput
endOutput:
 rts


