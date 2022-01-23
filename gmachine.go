// Package gmachine implements a simple virtual CPU, known as the G-machine.
package gmachine

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type Word uint64

type ElasticMemory map[Word]Word

type Option func(*Machine) error

func WithOutput(output io.Writer) Option {
	return func(m *Machine) error {
		m.output = output
		return nil
	}
}

func WithDebug() Option {
	return func(m *Machine) error {
		m.debug = true
		return nil
	}
}

type Machine struct {
	P        Word
	A        Word
	I        Word
	Memory   ElasticMemory
	FlagZero bool
	Stack    Stack

	output io.Writer
	debug  bool
}

func New(opts ...Option) *Machine {

	machine := &Machine{
		Memory: ElasticMemory{},
		output: os.Stdout,
	}

	for _, o := range opts {
		o(machine)
	}

	return machine
}

func (m *Machine) Run() {

	var err error

	for {

		if m.debug {
			fmt.Fprintln(m.output, m.DebugString())
		}

		opcode := m.Memory[m.P]
		m.P++
		switch opcode {
		case OpHALT:
			return
		case OpNOOP:
		case OpINCA:
			m.A++
		case OpDECA:
			m.A--
		case OpSETA:
			m.A = m.Next()
		case OpSETATOM:
			m.A = m.Memory[m.I]
		case OpSETI:
			m.I = m.Next()
		case OpINCI:
			m.I++
		case OpCMPI:
			iValue := m.Next()
			if iValue == m.I {
				m.FlagZero = true
			} else {
				m.FlagZero = false
			}
		case OpBIOS:
			io := m.Next()
			sendto := m.Next()

			if io == IOWrite {
				if sendto == SendToStdOut {
					fmt.Fprintf(m.output, "%c", m.A)
				}
			}
		case OpJUMP:
			m.P = m.Next()
		case OpJMPZ:
			if !m.FlagZero {
				m.P = m.Next()
			}
		case OpJSR:
			m.Stack.Push(m.P)
			m.P = m.Next()
		case OpRTS:
			m.P, err = m.Stack.Pop()
			if err != nil {
				fmt.Fprintf(m.output, "error with RTS, %s", err)
			}
		}
	}
}

func (m *Machine) Next() Word {

	location := m.P
	m.P++

	return m.Memory[location]
}

func (m *Machine) RunProgram(words []Word) {

	for k, v := range words {
		m.Memory[Word(k)] = v
	}

	m.Run()
}

func (m *Machine) DebugString() string {

	return fmt.Sprintf("Registers: P=%d, A=%d, I=%d\nMemory: %v\nStack: %v\n", m.P, m.A, m.I, m.Memory, m.Stack)
}

type Stack []Word

func (s *Stack) Push(word Word) {
	*s = append(*s, word)

}

func (s *Stack) Pop() (Word, error) {
	if len(*s) == 0 {
		return 0, fmt.Errorf("no values in stack.  cannot pop until a value is added to stack")
	}

	last := len(*s) - 1
	value := (*s)[last]
	*s = (*s)[:last]

	return value, nil
}

func AssembleFromString(text string) ([]Word, error) {

	reader := strings.NewReader(text)

	t := NewTokenizer()

	strs := t.Scanner(reader)

	words, err := Assemble(strs)
	if err != nil {
		return nil, err
	}

	return words, nil
}

func AssembleFromFile(path string) ([]Word, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %s", err)
	}
	defer file.Close()

	t := NewTokenizer()

	strs := t.Scanner(file)

	words, err := Assemble(strs)
	if err != nil {
		return nil, err
	}

	return words, nil

}

func Assemble(codes []string) ([]Word, error) {

	labels := map[string]Word{}
	refs := map[string][]Word{}

	var words []Word

	// first pass
	for index, code := range codes {

		// id labels
		if strings.HasSuffix(code, ":") {

			label := strings.ReplaceAll(code, ":", "")
			labels[label] = Word(index)
			continue
		}

		instruction, ok := TranslatorMap[code]
		if ok {

			words = append(words, instruction.Opcode)
			continue
		}

		inputOutput, ok := IOMap[code]
		if ok {

			words = append(words, Word(inputOutput))
			continue
		}

		// id label references
		if unicode.IsLetter(rune(code[0])) {
			refs[code] = append(refs[code], Word(index))
			words = append(words, 0)
			continue
		}

		// assemble data
		data, err := AssembleData(code)
		if err != nil {
			return nil, err
		}
		words = append(words, data...)

	}

	// second pass to populate assembly code with label address
	for label, addresses := range refs {

		for _, address := range addresses {
			words[address] = labels[label]
		}
	}

	return words, nil
}

func AssembleData(token string) ([]Word, error) {

	words := []Word{}

	token = strings.ReplaceAll(token, "#", "")

	if unicode.IsLetter(rune(token[0])) {
		for _, character := range token {
			words = append(words, Word(character))
		}
	} else if unicode.IsNumber(rune(token[0])) {
		num, err := strconv.Atoi(token)
		if err != nil {
			return nil, fmt.Errorf("unable to assemble data, %s", err)
		}
		words = append(words, Word(num))
	}

	return words, nil
}

func WriteWords(w io.Writer, words []Word) {

	for _, word := range words {
		raw := make([]byte, 8)
		binary.BigEndian.PutUint64(raw, uint64(word))
		w.Write(raw)
	}
}

func CreateBinary(sourcePath string, targetPath string) error {

	words, err := AssembleFromFile(sourcePath)
	if err != nil {
		return err
	}

	binaryFile, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer binaryFile.Close()

	WriteWords(binaryFile, words)

	return nil
}

func ReadWords(r io.Reader) []Word {
	raw := make([]byte, 8)
	words := []Word{}

	for {

		_, err := r.Read(raw)
		if err == io.EOF {
			break

		}
		if err != nil {
			return nil
		}
		bin := binary.BigEndian.Uint64(raw)
		words = append(words, Word(bin))
	}

	return words
}

// 6502, zeta
// assembler, create a binary, executor, create arm64 binary - string matching, bufio scanner, word by word
// opcodes with some operands vs none
// generate
// stack - store state, stack pointer
// tokens - const (e.g for)
// convert string to opcodes
// define const - pi
// label print
// assembler code
// subroutines print -
// go tool compile -S             compile but stop - no bin
// go tool compile -S gmachine.go | code -
// otool -vVt main | code -
// llvm
// not gate, register
// draw a triangle
// GPU draws triangle (hardware)
// mutex, concurrent
// locking
// disallow interrupt
// stack - push/pop
// mem fetch speed - get whole block
// delay in memory fetch
// cache - concurrency
// i/o - routine, memory address 9000 - print string, std lib - BIOS, ships with the machine
// memory mapped io - write to memory location that sends to
// DefaultMemSize is the number of 64-bit words of memory which will be
// allocated to a new G-machine by default.

// virtual memory
//
// dynamic memory sizing
// allocating more space as we go
// 9M

// layer to make app thinks it has access to machine - OS
// stdlib, runtime

// submit text instead of big array

// could i write tests as gmachine programs -- list of tests, testing framework
// opcode gmachine failtest

// const (
// 	DefaultMemSize = 1024
// )
