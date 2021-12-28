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
	HALT = iota
	NOOP
	INCA
	DECA
	SETA
	BIOS
	SETI
	INCI
	CMPI
	JUMP
	SETATOM
)

const (
	IOPNone = iota
	IOPWrite
	IOPRead
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

type Machine struct {
	P      Word
	A      Word
	I      Word
	Memory []Word
	Zero   bool

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
		case HALT:
			return
		case NOOP:
		case INCA:
			m.A++
		case DECA:
			m.A--
		case SETA:
			m.A = m.Next()
		case SETATOM:
			m.A = m.Memory[m.I]
		case SETI:
			m.I = m.Next()
		case INCI:
			m.I++
		case CMPI:
			iValue := m.Next()
			if iValue == m.I {
				m.Zero = true
			} else {
				m.Zero = false
			}

		case BIOS:
			io := m.Next()
			sendto := m.Next()

			if io == IOPWrite {
				if sendto == SendToStdOut {
					fmt.Fprintf(m.output, "%c", m.A)
				}
			}
		case JUMP:
			m.P = m.Next()

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
