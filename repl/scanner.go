package repl

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

// Scanner 定义一个扫描器
type Scanner struct {
	buffer []byte
	index  int
	length int
	line   int

	curr []byte //存储当前的token

	end  bool //标记是否结束
	flag bool //标记是否是行
}

// NewScanner 生成一个最基础的扫描器
func NewScanner() *Scanner {
	return &Scanner{
		buffer: nil,
		index:  0,
		length: 0,
		line:   1,
		curr:   nil,
		end:    false,
		flag:   false,
	}
}

// CustomScanner 实现一个简单的扫描器功能
// 如果第一个字符是.那么读取一行
// 否则读取到以；结尾
func (s *Scanner) CustomScanner() error {
	reader := bufio.NewReader(os.Stdin)
	return s.scanData(reader)

}

// 扫描数据 获取数据输入
// 注意：stdin输入的处理并不是安全的
func (s *Scanner) scanData(reader *bufio.Reader) error {

	var data []byte
	var flag bool
	firstChar, err := reader.ReadByte()
	if err != nil {
		return err // 如果到达输入末尾，结束循环
	}
	//检测是否为.命令 否则就是以;结尾的常规命令
	if firstChar == '.' {
		line, _, err := reader.ReadLine()
		if err != nil {
			return err // 如果读取失败，结束循环
		}
		data = line
	} else {
		var text bytes.Buffer
		text.WriteByte(firstChar)
		for {
			char, err := reader.ReadByte()
			if err != nil {
				os.Exit(1)
			}
			if char == ';' {
				break // 如果读取失败或者遇到分号，结束循环
			}
			text.WriteByte(char)
		}
		data = text.Bytes()
		flag = true
	}

	//超过最大长度 标记错误
	if len(data) > MaxBufferSize {
		return ErrorBufferOverflow
	}

	s.buffer = data
	s.line += 1
	s.length = len(data)
	s.flag = flag

	s.reset()
	return nil
}

// 一定是要进行更新的 不更新会出错的
func (s *Scanner) reset() {
	s.index = 0
	s.curr = []byte{}
	s.end = false
}

// 多行命令 进行不断地解析处理
func (s *Scanner) next() {
	for s.index < s.length && (s.buffer[s.index] == ' ' || s.buffer[s.index] == '\n' || s.buffer[s.index] == '\r') {
		s.index++
	}
	//只有每一次进来的时候 需要来更新一下 next
	if s.index >= s.length {
		s.end = true
		return
	}

	start := s.index

	if s.buffer[s.index] == ',' || s.buffer[s.index] == '(' || s.buffer[s.index] == ')' {
		s.curr = s.buffer[s.index : s.index+1]
		s.index++
		return
	}

	for s.index < s.length && s.buffer[s.index] != ' ' && s.buffer[s.index] != '\n' && s.buffer[s.index] != '\r' && s.buffer[s.index] != ',' && s.buffer[s.index] != '(' && s.buffer[s.index] != ')' {
		s.index++
	}

	s.curr = s.buffer[start:s.index]

	for s.index < s.length && s.buffer[s.index] == ' ' {
		s.index++
	}

}

// DealCommand 处理命令
func (s *Scanner) DealCommand() error {
	if !s.flag {
		return s.dotCommand()
	} else {
		return s.normalCommand()
	}
}

// 处理开头是.的命令
func (s *Scanner) dotCommand() error {
	if s.next(); bytes.Equal(s.curr, []byte("exit")) {
		return EXIT
	}
	return ErrorCommand
}

// 处理常规命令
func (s *Scanner) normalCommand() error {
	if s.next(); s.end {
		fmt.Printf("Unrecognized command '%s' %s.\n", s.buffer, s.curr)
		return ErrorCommand
	} else {
		currCommand := strings.ToUpper(string(s.curr))
		if currCommand == SELECT {
			//实现基础的select解析
			val, err := s.parseSelect()
			if err != nil {
				return err
			}
			fmt.Println(val)
		}
		if currCommand == INSERT {
			//实现基础的insert解析
			val, err := s.parseInsert()
			if err != nil {
				return err
			}
			fmt.Println(val)
		}
		if currCommand == DELETE {
			//实现基础的delete解析
			val, err := s.parseDelete()
			if err != nil {
				return err
			}
			fmt.Println(val)
		}
		fmt.Println("currCommand:", currCommand)
	}

	fmt.Println("Executed.")
	return nil
}
