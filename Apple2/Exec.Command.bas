1000  PRINT "Slot number? "
1010  INPUT S
1015  REM  Set up InputByte, OutputByte, InputFlags and OutputFlags
1020 IB = 49294 + S * 16
1030 OB = 49293 + S * 16
1040 FI = 49291 + S * 16
1050 OF = 49287 + S * 16
1060  POKE OF,15: REM  set to bank zero with flags high
1070  GOSUB 2000
1999  END 
2000  REM  Execute Command
2010  PRINT "$";
2020  INPUT CM$
2030 A = 5: GOSUB 8000: REM  ExecCommand
2040  FOR X = 1 TO  LEN (CM$)
2050 A =  ASC ( MID$ (CM$,X,1))
2060  GOSUB 8000
2070  NEXT 
2080 A = 0: GOSUB 8000
2100  GOSUB 9000
2101  PRINT  CHR$ (A);
2110  IF A = 0 THEN  RETURN 
2130  GOTO 2100
7000  REM  Get Flags
7010 F =  PEEK (FI)
7020 RF = 0
7025 WF = 0
7030  IF F > 128 THEN WF = 1:F = F - 128
7040  IF F > 64 THEN RF = 1
7050  RETURN 
8000  REM  Send Byte
8010  GOSUB 7000: REM  Wait for read flag low
8020  IF RF = 1 THEN  GOTO 8010
8030  POKE OB,A
8040  POKE OF,14: REM  Set write flag low
8050  GOSUB 7000: REM  Wait for read flag high
8060  IF RF = 0 GOTO 8050
8070  POKE OF,15: REM  Set flags high
8080  RETURN 
9000  REM  Get Byte
9010  POKE OF,13: REM  Set read flag low
9020  GOSUB 7000: REM  Wait for write flag low
9030  IF WF = 1 THEN  GOTO 9020
9040 A =  PEEK (IB)
9050  POKE OF,15: REM  Set flags high
9060  GOSUB 7000: REM  Wait for write flag high
9070  IF WF = 0 THEN  GOTO 9060
9080  RETURN 
