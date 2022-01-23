package gmachine_test

import (
	"gmachine"
	"testing"
)

func TestState(t *testing.T) {

	type testCase struct {
		text        string
		wantError   bool
		description string
	}

	tcs := []testCase{
		{
			text:        "SETA 1",
			wantError:   false,
			description: "OpCode Operand",
		},
		{
			text:        "BIOS IOWrite SendToStdOut",
			wantError:   false,
			description: "OpCode Multi Operand",
		},
		{
			text:        "SETA INCA",
			wantError:   true,
			description: "Missing operand",
		},
		{
			text:        "LABEL: LABEL:",
			wantError:   true,
			description: "Same label definition",
		},
	}

	s := gmachine.NewStatemachine()

	gotError := false
	for _, tc := range tcs {

		strs := s.Scanner(tc.text)

		err := s.Tokenize(strs)
		if err != nil {
			gotError = true
		}

		if tc.wantError != gotError {
			t.Fatalf("%s wantError: %v, gotError:%v", tc.description, tc.wantError, gotError)
		}

	}

}
