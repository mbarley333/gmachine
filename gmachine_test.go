package gmachine_test

import (
	"gmachine"
	"testing"
)

func TestNew(t *testing.T) {
	t.Parallel()
	g := gmachine.New()
	wantMemSize := gmachine.DefaultMemSize
	gotMemSize := len(g.Memory)
	if wantMemSize != gotMemSize {
		t.Errorf("want %d words of memory, got %d", wantMemSize, gotMemSize)
	}
	var wantP gmachine.Word = 0
	if wantP != g.P {
		t.Errorf("want initial P value %d, got %d", wantP, g.P)
	}
	var wantMemValue gmachine.Word = 0
	gotMemValue := g.Memory[gmachine.DefaultMemSize-1]
	if wantMemValue != gotMemValue {
		t.Errorf("want last memory location to contain %d, got %d", wantMemValue, gotMemValue)
	}
	var wantA gmachine.Word = 0
	if wantA != g.A {
		t.Errorf("want initial A value %d, got %d", wantA, g.A)
	}
}

func TestHalt(t *testing.T) {
	t.Parallel()

	g := gmachine.New()
	g.Run()

	want := gmachine.Word(1)
	got := g.P

	if want != got {
		t.Fatalf("want: %d, got: %d", want, got)
	}

}

func TestNOOP(t *testing.T) {
	t.Parallel()

	g := gmachine.New()

	g.Memory[0] = gmachine.NOOP
	g.Run()

	want := gmachine.Word(2)
	got := g.P

	if want != got {
		t.Fatalf("want: %d, got: %d", want, got)
	}

}

func TestINCA(t *testing.T) {
	t.Parallel()

	g := gmachine.New()

	g.Memory[0] = gmachine.INCA
	g.Run()

	want := gmachine.Word(1)
	got := g.A

	if want != got {
		t.Fatalf("want: %d, got: %d", want, got)
	}

}

func TestDECA(t *testing.T) {
	t.Parallel()

	g := gmachine.New()

	g.Memory[0] = gmachine.DECA
	g.A = 2
	g.Run()

	want := gmachine.Word(1)
	got := g.A

	if want != got {
		t.Fatalf("want: %d, got: %d", want, got)
	}

}

func TestSETA(t *testing.T) {
	t.Parallel()

	g := gmachine.New()

	g.Memory[0] = gmachine.SETA
	g.Memory[1] = 3

	g.Run()

	wantA := gmachine.Word(3)

	gotA := g.A

	if wantA != gotA {
		t.Fatalf("SETA want: %d, got: %d", wantA, gotA)
	}

	wantP := gmachine.Word(3)
	gotP := g.P

	if wantP != gotP {
		t.Fatalf("P want: %d, got: %d", wantP, gotP)
	}

}

func TestCalculate(t *testing.T) {
	t.Parallel()

	g := gmachine.New()
	g.A = 3

	opcodes := []gmachine.Word{gmachine.DECA, gmachine.DECA}

	g.Calculate(opcodes)

	want := gmachine.Word(1)
	got := g.A

	if want != got {
		t.Fatalf("want: %d, got: %d", want, got)
	}
}
