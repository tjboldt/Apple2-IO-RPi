; Copyright Terence J. Boldt (c)2020-2024
; Use of this source code is governed by an MIT
; license that can be found in the LICENSE file.

; This file formerly contained the source for the
; firmware that had a pseudo command prompt
; but is now empty except required common bytes

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
 cpx #$20    ;ID bytes for ProDOS
 cpx #$00    ;
 cpx #$03    ;
 cpx #$3C    ;ID byte for Autostart ROM

 lda #$3f ;set all flags high and page 3 of EPROM for menu
 jmp PageJump

;The following bytes must exist in this location for Pascal communications
;as they are standard pointers for entry points
.byte     CommInit&$00FF ;low byte of rts for init of Pascal comms
.byte     $43 ; low byte of read for Pascal comms
.byte     $49 ; low byte of write for Pascal comms
.byte     $4F ; low byte of rts for status of Pascal comms

CommInit:
 lda #$0f ; set all flags high and page 0 for comm driver
 sta OutputFlags
 ldx #$00 ; set error code to 0 for success
 rts

PageJump:
 sta OutputFlags
 jmp Start ;this jump is only called if coming in from PageJump with A=$0f

;entry points for ProDOS
DriverEntry:
 lda #$0f ;set all flags high and page 0 of EPROM
 sta OutputFlags
 jmp Start ; this is never called as the EPROM page changes

;load first two blocks and execute to boot
Start:
end:
 rts

.repeat	251-<end
.byte 0
.endrepeat

.byte      0,0     ;0000 blocks = check status
.byte      23       ;bit set(0=status 1=read 2=write) unset(3=format, 4/5=number of volumes, 6=interruptable, 7=removable)
.byte     DriverEntry&$00FF ;low byte of entry

