package repl

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
)
var (
	EXIT                  = errors.New("Exit")
	ERROR_COMMAND         = errors.New("error command")
	ERROR_BUFFER_OVERFLOW = errors.New("input exceeds buffer size")
)

const (
	MAX_BUFFER_SIZE = 512
)

type Scanner struct {
	buffer []byte
	index  int
	length int
	line   int
}

func newScanner() *Scanner {
	return &Scanner{
		buffer: make([]byte, MAX_BUFFER_SIZE),
		index:  0,
		line:   1,
	}
}

func showDb(){
	fd,err:=os.Open("./repl/show.txt")
	if err!=nil{
		panic(err)
	}
	defer fd.Close()
	data,err:=io.ReadAll(fd)
	if err!=nil{
		panic(err)
	}
	fmt.Printf("%s\n",data)
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

func (s *Scanner) dotCommand() error {
	if bytes.Equal(s.buffer[s.index:s.index+5], []byte(".exit")) {
		return EXIT
	}
	return ERROR_COMMAND
}
// 多行命令 获取next
func (s *Scanner) next() (data []byte) {
    prev := s.index

	for prev < s.length && s.buffer[prev] == ' ' {
        prev++
    }

    s.index = prev
    for prev < s.length && s.buffer[prev] != ' ' {
        prev++
    }
    data = (s.buffer[s.index:prev])

    for prev < s.length && s.buffer[prev] == ' ' {
        prev++
    }
    s.index = prev
    return data
}

func (s *Scanner) normalCommand() error {
	line := ""
	index:=1
	flag:=false
	for s.index < s.length {
		currCommand:=string(s.next())
		if index==1{
			if currCommand=="select" || currCommand=="insert"{
				flag=true
			}else{
				//todo
			}
		}
		data := currCommand+" "
		line += data
	}
	// fmt.Println(line)
	if flag{
		fmt.Println("Executed.")
		return nil
	}else{
		fmt.Printf("Unrecognized command '%s'.\n", line)
		return ERROR_COMMAND
	}
}

func (s *Scanner) dealCommand() error {
	if s.index == 0 && s.buffer[s.index] == '.' {
		return s.dotCommand()
	} else {
		return s.normalCommand()
	}
}

func (s *Scanner) printPrompt() {
	fmt.Printf("db:%d> ", s.line)
}

func ParseCommand() {
	showDb()
	scan := newScanner()
	for {
		scan.printPrompt()
		err := scan.readNewLine()
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}
		if err := scan.dealCommand(); err != nil {
			if err == EXIT {
				break
			}
			fmt.Println("Error:", err)
		}
	}
}
