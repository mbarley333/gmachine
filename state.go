package gmachine

import (
	"bufio"
	"fmt"
	"strings"
)

type StateFunc func(string) error

type Tokenizer struct {
	State        StateFunc
	OperandCount int
	LabelMap     map[string]bool
}

func NewTokenizer() *Tokenizer {
	s := &Tokenizer{
		LabelMap: map[string]bool{},
	}

	return s
}

func (t Tokenizer) Scanner(text string) []string {

	scanner := bufio.NewScanner(strings.NewReader(text))
	scanner.Split(bufio.ScanWords)

	var strs []string
	for scanner.Scan() {

		strs = append(strs, scanner.Text())
	}

	return strs
}

func (t *Tokenizer) StateOperand(str string) error {

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

	t.OperandCount -= 1

	return nil

}

func (t Tokenizer) StateOpCode(str string) error {
	return nil
}

func (t *Tokenizer) StateLabel(str string) error {
	_, ok := t.LabelMap[str]

	if ok {
		return fmt.Errorf("multiple label definition for: %s", str)
	} else {
		t.LabelMap[str] = true
	}
	return nil

}

func (t *Tokenizer) Tokenize(strs []string) error {

	if len(strs) == 0 {
		return fmt.Errorf("unable to Tokenize zero length slice")
	}

	for _, str := range strs {

		result, opCodeOk := TranslatorMap[str]

		switch {
		case t.OperandCount > 0:
			t.State = t.StateOperand
		case opCodeOk:
			t.State = t.StateOpCode
			t.OperandCount = result.Operands
		case isLabel(str):
			t.State = t.StateLabel
		}

		err := t.State(str)
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
