10  HOME 
100  PRINT  CHR$ (4)"BLOAD AT28C64B,A$2000"
200  PRINT "Program the firmare in slot #"
300  INPUT SL
400 FW = 8192 + 256 * SL: REM  Firmware source 
500 PS = 49287 + SL * 16: REM  Firmware page selection
600 EP = 49152 + SL * 256: REM  EPROM location
900  HOME 
1000  FOR PG = 0 TO 3
1004  PRINT : PRINT 
1005  PRINT "Writing page "PG" to slot "SL
1006  PRINT "_______________________________________"
1007  PRINT "_______________________________________";
1008  HTAB 1
1010  POKE PS,PG * 16 + 15: REM  Set firmware page
1020  FOR X = 0 TO 255
1030 A =  PEEK (FW + PG * 2048 + X)
1035 B =  PEEK (EP + X): REM  Skip if unchanged
1037  IF (A = B) THEN  GOTO 1045
1040  POKE EP + X,A
1041 B =  PEEK (EP + X)
1042  IF (B <  > A) THEN  GOTO 1041: REM Wait for write
1045  HTAB (X / 256) * 39 + 1
1046  INVERSE : PRINT " ";: NORMAL 
1050  NEXT X
1060  NEXT PG
1900  PRINT 
2000  PRINT "Firmware Update Complete"
