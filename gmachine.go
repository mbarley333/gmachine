// Package gmachine implements a simple virtual CPU, known as the G-machine.
package gmachine

// DefaultMemSize is the number of 64-bit words of memory which will be
// allocated to a new G-machine by default.
const (
	DefaultMemSize = 1024
)

type OpCode uint64

const (
	HALT = iota
	NOOP
)

type Machine struct {
	P      uint64
	Memory []uint64
}

func New() *Machine {

	machine := &Machine{
		P:      0,
		Memory: make([]uint64, DefaultMemSize),
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
		}
	}
}
