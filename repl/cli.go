package repl

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"os"
)

var (
	EXIT                = errors.New("exit")
	ErrorCommand        = errors.New("error command")
	ErrorBufferOverflow = errors.New("input exceeds buffer size")
)

const (
	MaxBufferSize = 512
)

// 展示go-sql信息
func showDb() {
	// 输出字符串内容
	data, err := os.ReadFile("./repl/show.txt")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", data)
}

// prompt提示
func (s *Scanner) printPrompt() {
	c := color.New(color.FgHiGreen, color.Bold)
	// 避免golang错误提醒
	_, _ = c.Printf("db:%d> ", s.line)
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
			if errors.Is(err, EXIT) {
				os.Exit(0)
			}
			fmt.Println("Error:", err)
		}
		// fmt.Println(scan)
	}
}
