package repl

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

const MAX_BUFFER_SIZE = 512

type Scanner struct {
	buffer []byte
	index  int
	line   int
}

func newScanner() *Scanner {
	return &Scanner{
		buffer: make([]byte, MAX_BUFFER_SIZE),
		index:  0,
		line:   0,
	}
}

var (
	EXIT = errors.New("Exit")
)

func (s *Scanner) readInput() error {
	reader := bufio.NewReader(os.Stdin)
	reader.ReadByte('.')
	line, err := reader.ReadBytes(';')
	if err != nil {
		os.Exit(1)
	}
	if len(line) > MAX_BUFFER_SIZE {
		return errors.New("")
	}
	s.buffer = line
	s.line += 1
	s.index = 0
	return nil
}
func (s *Scanner) dealCommand() error {
	if s.index == 0 {
		if s.buffer[s.index] == '.' {
			if s.buffer[s.index+1] == 'e' && s.buffer[s.index+2] == 'x' && s.buffer[s.index+3] == 'i' && s.buffer[s.index+4] == 't' {
				return EXIT
			}
			return nil
		}
		return nil
	}
	fmt.Printf("Unrecognized command '%s'.\n", s.buffer)
	return nil
}
func printPrompt() {
	fmt.Print("db> ")
}
func dealCommand() {

}
func ParseCommand() {
	// args := os.Args

	// if len(args) < 2 {
	// 	return
	// }

	// for i, arg := range args[1:] {
	// 	fmt.Printf("Argument %d: %s\n", i+1, arg)
	// }
	scan := newScanner()
	for {
		printPrompt()
		scan.readInput()
		if err := scan.dealCommand(); err != nil {
			if err == EXIT {
				break
			}
		}
	}
}
