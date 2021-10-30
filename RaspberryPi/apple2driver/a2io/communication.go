package a2io

type A2Io interface {
	Init()
	WriteByte(data byte) error
	WriteString(outString string) error
	WriteBlock(buffer []byte) error
	WriteBuffer(buffer []byte) error
	ReadByte() (byte, error)
	ReadString() (string, error)
	ReadBlock(buffer []byte) error
}

type A2Comm struct {
	A2Io A2Io
}
