// Package gmachine implements a simple virtual CPU, known as the G-machine.
package gmachine

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// 6502, zeta
// assembler, create a binary, executor, create arm64 binary - string matching, bufio scanner, word by word
// opcodes with some operands vs none
// generate
// stack - store state, stack pointer
// tokens - const (e.g for)
// convert string to opcodes
// define const - pi
// label print
// assembler code
// subroutines print -
// go tool compile -S             compile but stop - no bin
// go tool compile -S gmachine.go | code -
// otool -vVt main | code -
// llvm
// not gate, register
// draw a triangle
// GPU draws triangle (hardware)
// mutex, concurrent
// locking
// disallow interrupt
// stack - push/pop
// mem fetch speed - get whole block
// delay in memory fetch
// cache - concurrency
// i/o - routine, memory address 9000 - print string, std lib - BIOS, ships with the machine
// memory mapped io - write to memory location that sends to
// DefaultMemSize is the number of 64-bit words of memory which will be
// allocated to a new G-machine by default.

// virtual memory
//
// dynamic memory sizing
// allocating more space as we go
// 9M

// layer to make app thinks it has access to machine - OS
// stdlib, runtime

// submit text instead of big array

// could i write tests as gmachine programs -- list of tests, testing framework
// opcode gmachine failtest

const (
	DefaultMemSize = 1024
)

type Word uint64

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

type Option func(*Machine) error

func WithOutput(output io.Writer) Option {
	return func(m *Machine) error {
		m.output = output
		return nil
	}
}

func WithInput(input io.Reader) Option {
	return func(m *Machine) error {
		m.input = input
		return nil
	}
}

// P is Program Counter
// A is Arithmatic
// I holds the index value of memory location
// FlagZero used for loop operations and any boolean state holder
type Machine struct {
	P        Word
	A        Word
	I        Word
	Memory   []Word
	FlagZero bool

	output io.Writer
	input  io.Reader
}

func New(opts ...Option) *Machine {

	machine := &Machine{
		Memory: make([]Word, DefaultMemSize),
		output: os.Stdout,
		input:  os.Stdin,
	}

	for _, o := range opts {
		o(machine)
	}

	return machine
}

func (m *Machine) Run() {

	for {

		opcode := m.Memory[m.P]
		m.P++
		switch opcode {
		case OpHALT:
			return
		case OpNOOP:
		case OpINCA:
			m.A++
		case OpDECA:
			m.A--
		case OpSETA:
			m.A = m.Next()
		case OpSETATOM:
			m.A = m.Memory[m.I]
		case OpSETI:
			m.I = m.Next()
		case OpINCI:
			m.I++
		case OpCMPI:
			iValue := m.Next()
			if iValue == m.I {
				m.FlagZero = true
			} else {
				m.FlagZero = false
			}
		case OpBIOS:
			io := m.Next()
			sendto := m.Next()

			if io == IOWrite {
				if sendto == SendToStdOut {
					fmt.Fprintf(m.output, "%c", m.A)
				}
			}
		case OpJUMP:
			m.P = m.Next()
		case OpJMPZ:
			if !m.FlagZero {
				m.P = m.Next()
			}
		}
	}

}

func (m *Machine) Next() Word {
	location := m.P
	m.P++

	return m.Memory[location]
}

func (m *Machine) RunProgram(opcodes []Word) {

	for k, v := range opcodes {
		m.Memory[k] = v
	}

	m.Run()
}

type Instruction struct {
	Opcode   Word
	Operands int
}

var TranslatorMap = map[string]Instruction{
	"HALT": {Opcode: OpHALT, Operands: 0},
	"NOOP": {Opcode: OpNOOP, Operands: 0},
	"INCA": {Opcode: OpINCA, Operands: 0},
	"DECA": {Opcode: OpDECA, Operands: 0},
	"SETA": {Opcode: OpSETA, Operands: 1},
}

func AssembleFromString(codeString string) []Word {

	scanner := bufio.NewScanner(strings.NewReader(codeString))
	scanner.Split(bufio.ScanWords)

	var codes []string
	for scanner.Scan() {

		codes = append(codes, scanner.Text())
	}

	words := Assemble(codes)

	return words
}

func Assemble(codes []string) []Word {

	var words []Word

	for _, code := range codes {
		instruction, ok := TranslatorMap[code]
		if ok {
			words = append(words, instruction.Opcode)
		} else {
			word, err := strconv.Atoi(code)
			if err != nil {
				fmt.Printf("%d of type %T", word, word)
			}
			words = append(words, Word(word))
		}

	}

	return words
}
