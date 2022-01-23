package gmachine

import (
	"bufio"
	"fmt"
	"strings"
)

type StateFunc func(string) error

type Statemachine struct {
	State        StateFunc
	OperandCount int
	LabelMap     map[string]bool
}

func NewStatemachine() *Statemachine {
	s := &Statemachine{
		LabelMap: map[string]bool{},
	}

	return s
}

func (s Statemachine) Scanner(text string) []string {

	scanner := bufio.NewScanner(strings.NewReader(text))
	scanner.Split(bufio.ScanWords)

	var strs []string
	for scanner.Scan() {

		strs = append(strs, scanner.Text())
	}

	return strs
}

func (s *Statemachine) StateOperand(str string) error {

	_, okOpCode := TranslatorMap[str]
	if okOpCode {
		return fmt.Errorf("expecting operand, got OpCode: %s", str)
	}

	if strings.Index(":", str) > 0 {
		return fmt.Errorf("expecting operand, got Label: %s", str)
	}

	if strings.Index("#", str) > 0 {
		return fmt.Errorf("expecting operand, got #string: %s", str)
	}

	s.OperandCount -= 1

	return nil

}

func (s Statemachine) StateOpCode(str string) error {
	return nil
}

func (s *Statemachine) StateLabel(str string) error {
	_, ok := s.LabelMap[str]

	if ok {
		return fmt.Errorf("multiple label definition for: %s", str)
	} else {
		s.LabelMap[str] = true
	}
	return nil

}

func (s *Statemachine) Tokenize(strs []string) error {

	if len(strs) == 0 {
		return fmt.Errorf("unable to Tokenize zero length slice")
	}

	for _, str := range strs {

		result, opCodeOk := TranslatorMap[str]

		switch {
		case s.OperandCount > 0:
			s.State = s.StateOperand
		case opCodeOk:
			s.State = s.StateOpCode
			s.OperandCount = result.Operands
		case isLabel(str):
			s.State = s.StateLabel
		}

		err := s.State(str)
		if err != nil {
			return err
		}
	}

	return nil
}

func isLabel(str string) bool {

	ret := false
	if strings.Index(":", str) > 0 {
		ret = true
	}

	return ret
}
