// Package gmachine implements a simple virtual CPU, known as the G-machine.
package gmachine

import (
	"fmt"
	"io"
	"os"
)

// DefaultMemSize is the number of 64-bit words of memory which will be
// allocated to a new G-machine by default.
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

		instruction := m.Memory[m.P]
		m.P++
		switch instruction {
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
