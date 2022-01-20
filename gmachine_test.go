package gmachine_test

import (
	"bytes"
	"gmachine"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// func TestNew(t *testing.T) {

// 	t.Parallel()
// 	g := gmachine.New()

// 	want := &gmachine.Machine{
// 		P:        0,
// 		A:        0,
// 		I:        0,
// 		Memory:   gmachine.ElasticMemory{},
// 		FlagZero: false,
// 		Stack:    gmachine.Stack{},
// 	}

// 	got := g

// 	if !cmp.Equal(want, got) {
// 		t.Error(cmp.Diff(want, got))
// 	}

// }

// 	var wantP gmachine.Word = 0
// 	if wantP != g.P {
// 		t.Errorf("want initial P value %d, got %d", wantP, g.P)
// 	}
// 	var wantMemValue gmachine.Word = 0
// 	gotMemValue := g.Memory[gmachine.Word(len(g.Memory)-1)]
// 	if wantMemValue != gotMemValue {
// 		t.Errorf("want last memory location to contain %d, got %d", wantMemValue, gotMemValue)
// 	}
// 	var wantA gmachine.Word = 0
// 	if wantA != g.A {
// 		t.Errorf("want initial A value %d, got %d", wantA, g.A)
// 	}
// }

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
		gmachine.IOWrite,
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

func TestValidateInstructions(t *testing.T) {
	t.Parallel()

	type testCase struct {
		codes       []string
		ErrExpected bool
		description string
		operands    int
	}

	tcs := []testCase{
		{codes: []string{"SETA", "72"}, operands: 1, ErrExpected: false, description: "expect no error"},
		{codes: []string{"SETA"}, operands: 1, ErrExpected: true, description: "expect error"},
		{codes: []string{"SETA", "DECA"}, operands: 1, ErrExpected: true, description: "expect error"},
	}

	want := false
	got := false

	for _, tc := range tcs {

		want = tc.ErrExpected

		err := gmachine.ValidateInstructions(tc.codes, tc.operands)
		if err != nil {
			got = true
		}

		if want != got {
			t.Fatalf("%s want: %v, got: %v", tc.description, want, got)
		}

	}

}

func TestWriteWords(t *testing.T) {
	t.Parallel()

	output := &bytes.Buffer{}

	words := []gmachine.Word{gmachine.OpINCA, gmachine.OpDECA, gmachine.Word(72)}

	gmachine.WriteWords(output, words)

	want := []byte{0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 72}
	got := output.Bytes()

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}

}

func TestAssembleFromString(t *testing.T) {
	t.Parallel()

	str := "INCA DECA SETA 12"

	want := []gmachine.Word{
		gmachine.OpINCA,
		gmachine.OpDECA,
		gmachine.OpSETA,
		12,
	}

	got, err := gmachine.AssembleFromString(str)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}

}

func TestAssemble(t *testing.T) {
	t.Parallel()

	code := []string{"INCA", "DECA", "#72"}

	want := []gmachine.Word{gmachine.OpINCA, gmachine.OpDECA, gmachine.Word(72)}
	got, err := gmachine.Assemble(code)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}

}

func TestAssembleFromFile(t *testing.T) {
	t.Parallel()

	path := "testdata/testFile.gmachine"

	want := []gmachine.Word{gmachine.OpINCA, gmachine.OpDECA, gmachine.Word(72)}

	got, err := gmachine.AssembleFromFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}

}

func TestAssembleData(t *testing.T) {
	t.Parallel()

	text := "#HelloWorld"

	want := []gmachine.Word{
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
	}

	got, err := gmachine.AssembleData(text)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}

}

func TestReadWords(t *testing.T) {
	t.Parallel()

	sourcePath := "testdata/testFile.gmachine"

	targetPath := t.TempDir() + "/testFile.g"

	gmachine.CreateBinary(sourcePath, targetPath)

	file, err := os.Open(targetPath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	want := []gmachine.Word{
		gmachine.OpINCA,
		gmachine.OpDECA,
		gmachine.Word(72),
	}

	got := gmachine.ReadWords(file)

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}

}

func TestLabelAssemble(t *testing.T) {
	t.Parallel()

	str := "JUMP LABEL HALT LABEL: INCA JUMP 2"

	want := []gmachine.Word{
		gmachine.OpJUMP,
		gmachine.Word(3),
		gmachine.OpHALT,
		gmachine.OpINCA,
		gmachine.OpJUMP,
		gmachine.Word(2),
	}

	words, err := gmachine.AssembleFromString(str)
	if err != nil {
		t.Fatal(err)
	}

	got := words

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}

}

func TestDuplicateLabelAssemble(t *testing.T) {
	t.Parallel()

	str := "LABEL: INCA LABEL DECA"

	wantError := true
	_, err := gmachine.AssembleFromFile(str)

	gotError := false
	if err != nil {
		gotError = true
	}

	if wantError != gotError {
		t.Fatalf("want: %v, got %v", wantError, gotError)
	}

}

func TestStack(t *testing.T) {
	t.Parallel()

	s := gmachine.Stack{}

	// need a stack to track jump starting point
	want := 0

	got := len(s)

	if want != got {
		t.Fatalf("want: %d, got%d", want, got)
	}

	s.Push(gmachine.Word(1))

	wantPush := 1

	gotPush := len(s)

	if wantPush != gotPush {
		t.Fatalf("TestPush want: %d, got%d", wantPush, gotPush)
	}

	wantPop := gmachine.Word(1)

	gotPop, err := s.Pop()
	if err != nil {
		t.Fatal(err)
	}

	if wantPop != gotPop {
		t.Fatalf("TestPop want: %d, got:%d", wantPop, gotPop)
	}

}
