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
	SendToStdOut
)

type Instruction struct {
	Opcode   Word
	Operands int
}

var TranslatorMap = map[string]Instruction{
	"HALT":    {Opcode: OpHALT, Operands: 0},
	"NOOP":    {Opcode: OpNOOP, Operands: 0},
	"INCA":    {Opcode: OpINCA, Operands: 0},
	"DECA":    {Opcode: OpDECA, Operands: 0},
	"SETA":    {Opcode: OpSETA, Operands: 1},
	"SETI":    {Opcode: OpSETI, Operands: 1},
	"BIOS":    {Opcode: OpBIOS, Operands: 2},
	"INCI":    {Opcode: OpINCI, Operands: 0},
	"CMPI":    {Opcode: OpCMPI, Operands: 1},
	"JUMP":    {Opcode: OpJUMP, Operands: 1},
	"JMPZ":    {Opcode: OpJMPZ, Operands: 1},
	"SETATOM": {Opcode: OpSETATOM, Operands: 0},
	"JSR":     {Opcode: OpJSR, Operands: 1},
	"RTS":     {Opcode: OpRTS, Operands: 0},
}

var IOMap = map[string]int{
	"IONone":       IONone,
	"IOWrite":      IOWrite,
	"IORead":       IORead,
	"SendToStdOut": SendToStdOut,
}
