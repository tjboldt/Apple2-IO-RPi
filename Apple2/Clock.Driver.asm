enableWriteLang = $C088
disableWriteLang = $C082

prodosJump = $BF06
prodosClockCode = $BF07
prodosMachineId = $BF98

; Find Apple2-IO-RPi card

; Change driver code to point to the card

; Change destination to be ProDOS clock code location
 lda prodosClockCode
 sta DriverDestination+1
 lda prodosClockCode+1
 sta DriverDestination+2

; Changing RTS to JMP enables clock driver
 lda #$4C ; jump instruction
 sta prodosJump

; Enable writing to language card RAM
; by triggering switch twice
 lda enableWriteLang
 lda enableWriteLang

; write driver code to language card RAM
WriteDriver:
 ldy #EndDriver-Driver+1
 lda Driver,y
DriverDestination:
 sta $D000,y ; this address gets modified above
 dey
 bne WriteDriver

; Disable writing to language card RAM
 lda disableWriteLang

; Update ProDOS Machine ID to mark clock as enabled
 lda prodosMachineId
 ora #$01
 sta prodosMachineId

 rts

Driver:
EndDriver:
 rts
