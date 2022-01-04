package gmachine_test

import (
	"gmachine"
	"testing"
)

func TestElasticMemoryAdd(t *testing.T) {
	t.Parallel()

	m := gmachine.NewElasticMemory()

	wantInital := 0

	gotInitial := len(m)

	if wantInital != gotInitial {
		t.Fatalf("WantZero: want: %d, got: %d", wantInital, gotInitial)
	}

	wantAdd := 2

	opcode := []gmachine.Word{
		gmachine.OpINCA,
		gmachine.OpDECA,
	}

	m.Add(opcode)

	gotAdd := len(m)

	if wantAdd != gotAdd {
		t.Fatalf("WantOne: want: %d, got: %d", wantAdd, gotAdd)
	}

	wantAddToAddress := gmachine.OpINCA

	m.AddToAddress(1000, gmachine.OpINCA)

	gotAddToAddress := m[1000]

	if wantAdd != gotAdd {
		t.Fatalf("AddToAddress: want: %d, got: %d", wantAddToAddress, gotAddToAddress)
	}

}
