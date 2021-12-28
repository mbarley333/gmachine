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

func TestSETI(t *testing.T) {
	t.Parallel()

	g := gmachine.New()

	wantNoI := gmachine.Word(0)

	gotNoI := g.I

	if wantNoI != gotNoI {
		t.Fatalf("want: %d, got: %d", wantNoI, gotNoI)
	}

	opcodes := []gmachine.Word{
		gmachine.SETI,
		3,
	}
	g.RunProgram(opcodes)

	want := gmachine.Word(3)
	got := g.I

	if want != got {
		t.Fatalf("want: %d, got: %d", want, got)
	}

}

func TestINCI(t *testing.T) {
	t.Parallel()

	g := gmachine.New()

	opcodes := []gmachine.Word{
		gmachine.SETI,
		3,
		gmachine.INCI,
	}
	g.RunProgram(opcodes)

	want := gmachine.Word(4)
	got := g.I

	if want != got {
		t.Fatalf("want: %d, got: %d", want, got)
	}

}

func TestJUMP(t *testing.T) {

	t.Parallel()

	g := gmachine.New()

	opcodes := []gmachine.Word{
		gmachine.JUMP,
		3,
		'A',
		gmachine.SETI,
		2,
	}

	g.RunProgram(opcodes)

	wantI := gmachine.Word(2)
	gotI := g.I

	if wantI != gotI {
		t.Fatalf("want: %d, got: %d", wantI, gotI)
	}

}

func TestSETATOM(t *testing.T) {

	t.Parallel()

	g := gmachine.New()

	opcodes := []gmachine.Word{
		gmachine.SETI,
		2,
		72,
		gmachine.SETATOM,
	}

	g.RunProgram(opcodes)

	want := gmachine.Word(72)
	got := g.A

	if want != got {
		t.Fatalf("want: %d, got: %d", want, got)
	}

}

func TestCMPI(t *testing.T) {
	t.Parallel()

	type testCase struct {
		opcodes     []gmachine.Word
		want        bool
		description string
	}

	tcs := []testCase{
		{opcodes: []gmachine.Word{gmachine.SETI, 2, gmachine.CMPI, 2}, want: true, description: "2 == 2"},
		{opcodes: []gmachine.Word{gmachine.SETI, 2, gmachine.CMPI, 3}, want: false, description: "2 != 3"},
		{opcodes: []gmachine.Word{gmachine.SETI, 2, gmachine.CMPI, 3, gmachine.SETI, 3, gmachine.CMPI, 3}, want: true, description: "3 != 2, 3 == 3"},
	}

	for _, tc := range tcs {

		g := gmachine.New()

		g.RunProgram(tc.opcodes)
		got := g.FlagZero

		if tc.want != got {
			t.Fatalf("%s want: %v, got: %v", tc.description, tc.want, got)
		}

	}

}

func TestLoopWithJMPZ(t *testing.T) {

	t.Parallel()

	g := gmachine.New()

	opcodes := []gmachine.Word{

		gmachine.INCI,
		gmachine.CMPI,
		10,
		gmachine.JMPZ,
		0,
	}

	g.RunProgram(opcodes)

	want := gmachine.Word(10)
	got := g.I

	if want != got {
		t.Fatalf("want: %d, got: %d", want, got)
	}

}

func TestHelloWorld(t *testing.T) {

	t.Parallel()

	output := &bytes.Buffer{}

	g := gmachine.New(
		gmachine.WithOutput(output),
	)

	opcodes := []gmachine.Word{
		gmachine.JUMP,
		12,
		72,
		101,
		108,
		108,
		111,
		87,
		111,
		114,
		108,
		100,
		gmachine.SETI,
		2,
		gmachine.SETATOM,
		gmachine.BIOS,
		gmachine.IOPWrite,
		gmachine.SendToStdOut,
		gmachine.INCI,
		gmachine.CMPI,
		12,
		gmachine.JMPZ,
		14,
	}

	g.RunProgram(opcodes)

	want := "HelloWorld"
	got := output.String()

	if want != got {
		t.Fatalf("want: %q, got: %q", want, got)
	}

}
