; Copyright Terence J. Boldt (c)2024
; Use of this source code is governed by an MIT
; license that can be found in the LICENSE file.

; This file contains the source for the clock
; driver for the Apple2-IO-RPi

EnableWriteLang = $C08B
DisableWriteLang = $C082

ProdosJump = $BF06
ProdosClockCode = $BF07
ProdosMachineId = $BF98

PrintChar = $FDED
PrintHex = $FDE3
PrintByte = $FDDA

GetTimeCommand = $03

.org $2000
; Find Apple2-IO-RPi card
 ldx #$06
 ldy #$09
CheckIdBytes:
 lda $C700,y ; !! Self modifying code
 cmp IdBytes,y
 bne NextCard
 dey
 bne CheckIdBytes
 jmp FoundCard
NextCard:
 dec CheckIdBytes+2
 ldy #$09
 dex
 bne CheckIdBytes
CardNotFound:
 ldy #$00
PrintCardNotFound:
 lda TextCardNotFound,y
 beq Failed
 ora #$80
 jsr PrintChar
 iny
 bne PrintCardNotFound
Failed:
 rts
FoundCard:
 ldy #$00
PrintCardFound:
 lda TextCardFound,y
 beq PrintCardNumber
 ora #$80
 jsr PrintChar
 iny
 bne PrintCardFound
PrintCardNumber:
 inx
 txa
 jsr PrintHex
 lda #$8D
 jsr PrintChar


; Change driver code to point to the card
 txa
 ora #$C0
 sta GetByte+2
 sta SendByte+2

; Change destination to be ProDOS clock code location
 lda ProdosClockCode
 sta DriverDestination+1
 lda ProdosClockCode+1
 sta DriverDestination+2

; Changing RTS to JMP enables clock driver
 lda #$4C ; jump instruction
 sta ProdosJump

; Enable writing to language card RAM
; by triggering switch twice
 lda EnableWriteLang
 lda EnableWriteLang

; write driver code to language card RAM
 ldy #$00
WriteDriver:
 lda Driver,y
DriverDestination:
 sta $D742,y ; !! this address gets modified above
 iny
 cpy #EndDriver-Driver+2
 bne WriteDriver

; Disable writing to language card RAM
 lda DisableWriteLang

; Update ProDOS Machine ID to mark clock as enabled
 lda ProdosMachineId
 ora #$01
 sta ProdosMachineId

 ldy #$00
PrintFinished:
 lda TextDriverInstalled,y
 beq Finished
 ora #$80
 jsr PrintChar
 iny
 bne PrintFinished
Finished:
 rts

IdBytes:
.byte  $E0,$20,$E0,$00,$E0,$03,$E0,$3C,$A9,$3F

TextCardFound:
.byte "Found Apple2-IO-RPi in slot ",$00

TextCardNotFound:
.byte "Apple2-IO-RPi not found",$8D,$00

TextDriverInstalled:
.byte "Clock driver installed",$8D,$00

Driver:
 lda #GetTimeCommand
SendByte:
 jsr $C749 ; !! address gets modified on installation
 ldy #$00
getTimeByte:
GetByte:
 jsr $C743 ; !! address gets modified on installation
 sta $bf90,y
 iny
 cpy #$04
 bne getTimeByte
EndDriver:
 rts
