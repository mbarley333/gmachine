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
		gmachine.OpNOOP,
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
		gmachine.OpINCA,
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
		gmachine.OpDECA,
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
		gmachine.OpSETA,
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
			opcodes:        []gmachine.Word{gmachine.OpSETA, 5, gmachine.OpDECA, gmachine.OpDECA},
			expectedResult: gmachine.Word(3),
		},
		{
			opcodes:        []gmachine.Word{gmachine.OpSETA, 7, gmachine.OpDECA, gmachine.OpDECA},
			expectedResult: gmachine.Word(5),
		},
		{
			opcodes:        []gmachine.Word{gmachine.OpSETA, 2, gmachine.OpDECA, gmachine.OpDECA},
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
		gmachine.OpSETA,
		'J',
		gmachine.OpBIOS,
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
		gmachine.OpSETI,
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
		gmachine.OpSETI,
		3,
		gmachine.OpINCI,
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
		gmachine.OpJUMP,
		3,
		'A',
		gmachine.OpSETI,
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
		gmachine.OpSETI,
		2,
		72,
		gmachine.OpSETATOM,
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
		{opcodes: []gmachine.Word{gmachine.OpSETI, 2, gmachine.OpCMPI, 2}, want: true, description: "2 == 2"},
		{opcodes: []gmachine.Word{gmachine.OpSETI, 2, gmachine.OpCMPI, 3}, want: false, description: "2 != 3"},
		{opcodes: []gmachine.Word{gmachine.OpSETI, 2, gmachine.OpCMPI, 3, gmachine.OpSETI, 3, gmachine.OpCMPI, 3}, want: true, description: "3 != 2, 3 == 3"},
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

		gmachine.OpINCI,
		gmachine.OpCMPI,
		10,
		gmachine.OpJMPZ,
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
		gmachine.OpJUMP,
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
		gmachine.OpSETI,
		2,
		gmachine.OpSETATOM,
		gmachine.OpBIOS,
		gmachine.IOPWrite,
		gmachine.SendToStdOut,
		gmachine.OpINCI,
		gmachine.OpCMPI,
		12,
		gmachine.OpJMPZ,
		14,
	}

	g.RunProgram(opcodes)

	want := "HelloWorld"
	got := output.String()

	if want != got {
		t.Fatalf("want: %q, got: %q", want, got)
	}

}
