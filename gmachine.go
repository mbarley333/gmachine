// Package gmachine implements a simple virtual CPU, known as the G-machine.
package gmachine

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
)

type Machine struct {
	P      Word
	A      Word
	Memory []Word
}

func New() *Machine {

	machine := &Machine{
		P:      0,
		A:      0,
		Memory: make([]Word, DefaultMemSize),
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
		}
	}
}

func (m *Machine) Calculate(opcodes []Word) {

	for k, v := range opcodes {
		m.Memory[k] = v
	}
	m.Run()

}
