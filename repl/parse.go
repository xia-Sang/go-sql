package repl

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
)

var (
	EXIT                  = errors.New("Exit")
	ERROR_COMMAND         = errors.New("error command")
	ERROR_BUFFER_OVERFLOW = errors.New("input exceeds buffer size")
)

const (
	MAX_BUFFER_SIZE = 512
)

// 展示go-sql信息
func showDb() {
	fd, err := os.Open("./repl/show.txt")
	if err != nil {
		panic(err)
	}
	defer fd.Close()
	data, err := io.ReadAll(fd)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", data)
}

// 读取新的一行
func (s *Scanner) readNewLine() error {
	reader := bufio.NewReader(os.Stdin)
	line, _, err := reader.ReadLine()
	if err != nil {
		os.Exit(1)
	}
	length := len(line)
	if length > MAX_BUFFER_SIZE {
		return ERROR_BUFFER_OVERFLOW
	}
	s.buffer = line
	s.length = length
	s.line += 1
	s.index = 0

	return nil
}

// 测试新的一行
func (s *Scanner) readScanner(reader *bufio.Reader) error {
	line, _, err := reader.ReadLine()
	if err != nil {
		os.Exit(1)
	}
	length := len(line)
	if length > MAX_BUFFER_SIZE {
		return ERROR_BUFFER_OVERFLOW
	}
	s.buffer = line
	s.length = length
	s.line += 1
	s.index = 0

	return nil
}

// prompt提示
func (s *Scanner) printPrompt() {
	c := color.New(color.FgHiGreen, color.Bold)
	c.Printf("db:%d> ", s.line)
}

// ParseCommand 解析命令行
func ParseCommand() {
	showDb()
	scan := NewScanner()
	for {
		scan.printPrompt()
		err := scan.CustomScanner()
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		if err := scan.DealCommand(); err != nil {
			if err == EXIT {
				os.Exit(0)
			}
			fmt.Println("Error:", err)
		}
		// fmt.Println(scan)
	}
}
