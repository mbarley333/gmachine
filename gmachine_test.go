package gmachine_test

import (
	"bytes"
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

	opcodes := []gmachine.Word{
		gmachine.NOOP,
	}
	g.RunProgram(opcodes)

	want := gmachine.Word(2)
	got := g.P

	if want != got {
		t.Fatalf("want: %d, got: %d", want, got)
	}

}

func TestINCA(t *testing.T) {
	t.Parallel()

	g := gmachine.New()

	opcodes := []gmachine.Word{
		gmachine.INCA,
	}
	g.RunProgram(opcodes)

	want := gmachine.Word(1)
	got := g.A

	if want != got {
		t.Fatalf("want: %d, got: %d", want, got)
	}

}

func TestDECA(t *testing.T) {
	t.Parallel()

	g := gmachine.New()
	g.A = 2

	opcodes := []gmachine.Word{
		gmachine.DECA,
	}
	g.RunProgram(opcodes)

	want := gmachine.Word(1)
	got := g.A

	if want != got {
		t.Fatalf("want: %d, got: %d", want, got)
	}

}

func TestSETA(t *testing.T) {
	t.Parallel()

	g := gmachine.New()

	opcodes := []gmachine.Word{
		gmachine.SETA,
		3,
	}
	g.RunProgram(opcodes)

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

func TestRunProgram(t *testing.T) {
	t.Parallel()

	g := gmachine.New()

	type testCase struct {
		opcodes        []gmachine.Word
		expectedResult gmachine.Word
	}

	tcs := []testCase{
		{
			opcodes:        []gmachine.Word{gmachine.SETA, 5, gmachine.DECA, gmachine.DECA},
			expectedResult: gmachine.Word(3),
		},
		{
			opcodes:        []gmachine.Word{gmachine.SETA, 7, gmachine.DECA, gmachine.DECA},
			expectedResult: gmachine.Word(5),
		},
		{
			opcodes:        []gmachine.Word{gmachine.SETA, 2, gmachine.DECA, gmachine.DECA},
			expectedResult: gmachine.Word(0),
		},
	}

	for _, tc := range tcs {
		g.P = gmachine.Word(0)
		g.RunProgram(tc.opcodes)

		want := tc.expectedResult
		got := g.A

		if want != got {
			t.Fatalf("want: %d, got: %d", want, got)
		}

	}

}

func TestBIOSWrite(t *testing.T) {
	t.Parallel()

	output := &bytes.Buffer{}
	g := gmachine.New(
		gmachine.WithOutput(output),
	)

	opcodes := []gmachine.Word{
		gmachine.SETA,
		'J',
		gmachine.BIOS,
		gmachine.IOPWrite,
		gmachine.SendToStdOut,
	}
	g.RunProgram(opcodes)

	want := "J"
	got := output.String()

	if want != got {
		t.Fatalf("want: %q, got: %q", want, got)
	}

}
