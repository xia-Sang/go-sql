package repl

import (
	"bufio"
	"strings"
	"testing"
)

func Test_scanner(t *testing.T) {
	s := `INSERT INTO Websites (name, url, alexa, country)
VALUES ('百度','https://www.baidu.com/','4','CN');`
	s = `INSERT INTO Websites (name, url, country)
VALUES ('stackoverflow', 'http://stackoverflow.com/', 'IND');`
	scan := NewScanner()
	err := scan.scanData(bufio.NewReader(strings.NewReader(s)))
	t.Log(err)
	t.Log(scan)
	for !scan.end {
		scan.next()
		t.Log(string(scan.curr))
	}
}
