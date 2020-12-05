	ldx	#$50	;slot 5 for this test
start:	lda	#$80	;set read flag low (ready to read)
	sta	$c08d,x	;bit 1 low for writing values
waitwl:	lda	$c08e,x	;bit 0 low for reading values
	bmi	waitwl	;wait for write flag low	
	jsr	$fde3	;print nibble of data
	lda 	#$c0	;set read flag high (done reading)
	sta	$c08d,x
waitwh:	lda	$c08e,x
	bpl	waitwh	;wait for write flag high
	bmi	start	;go around againg for next nibble
