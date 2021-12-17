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
	var wantP uint64 = 0
	if wantP != g.P {
		t.Errorf("want initial P value %d, got %d", wantP, g.P)
	}
	var wantMemValue uint64 = 0
	gotMemValue := g.Memory[gmachine.DefaultMemSize-1]
	if wantMemValue != gotMemValue {
		t.Errorf("want last memory location to contain %d, got %d", wantMemValue, gotMemValue)
	}
	var wantA uint64 = 0
	if wantA != g.A {
		t.Errorf("want initial A value %d, got %d", wantA, g.A)
	}
}

func TestHalt(t *testing.T) {
	t.Parallel()

	g := gmachine.New()
	g.Run()

	var want uint64
	want = 1
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

	var want uint64
	want = 2
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

	var want uint64
	want = 1
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

	var want uint64
	want = 1
	got := g.A

	if want != got {
		t.Fatalf("want: %d, got: %d", want, got)
	}

}
