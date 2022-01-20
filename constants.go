package gmachine

const (
	OpHALT = iota
	OpNOOP
	OpINCA
	OpDECA
	OpSETA
	OpBIOS
	OpSETI
	OpINCI
	OpCMPI
	OpJUMP
	OpJMPZ
	OpSETATOM
	OpJSR
	OpRTS
)

const (
	IONone = iota
	IOWrite
	IORead
)

const (
	SendToNone = iota
	SendToStdOut
	ReadFromStdin
)
