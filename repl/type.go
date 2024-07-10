package repl

import (
	"fmt"
	"strings"
)

type StatementType string

// SQL type tokens
const (
	UNSUPPORTED = "N/A"
	SELECT      = "SELECT"
	FROM        = "FROM"
	WHERE       = "WHERE"
	LIMIT       = "LIMIT"
	INSERT      = "INSERT"
	INTO        = "INTO"
	VALUES      = "VALUES"
	ASTERISK    = "*"
)

type InsertTree struct {
	Table   string     //table_name：需要插入新记录的表名
	Columns []string   //column1, column2, ...：需要插入的字段名
	Values  [][]string //value1, value2, ...：需要插入的字段值
}

// https://www.runoob.com/sql/sql-insert.html
/*
	INSERT INTO table_name
	VALUES (value1,value2,value3,...);

	INSERT INTO table_name (column1,column2,column3,...)
	VALUES (value1,value2,value3,...);
*/
// 实现对于insert的解析
func (s *Scanner) parseInsert1() (ast *InsertTree, err error) {
	for s.next(); !s.end; {
		fmt.Println("tok:", string(s.curr))
	}
	return nil, nil
}
func (s *Scanner) parseInsert() (ast *InsertTree, err error) {
	ast = &InsertTree{}
	// 不需要进行 解析 insert
	// 解析 into
	if s.next(); s.end || strings.ToUpper(string(s.curr)) != INTO {
		err = fmt.Errorf("%s is not INSERT statement,error token: %s", s.buffer, s.curr)
		return
	}
	// 解析table name
	if s.next(); s.end {
		err = fmt.Errorf("%s expect table after INSERT INTO", s.buffer)
		return
	} else {
		ast.Table = string(s.curr)
	}
	// 解析column或者values
	if s.next(); s.end {
		err = fmt.Errorf("%s expect VALUES or (colNames),error token:%s", s.buffer, s.curr)
		return
	} else {
		currToken := strings.ToUpper(string(s.curr))
		if currToken == "(" {
			ast.Columns = make([]string, 0)
			for {
				if s.next(); s.end {
					if len(ast.Columns) == 0 {
						err = fmt.Errorf("%s get Columns failed", s.buffer)
					}
					return
				} else {
					currToken := string(s.curr)
					if currToken == "," {
						continue
					} else if currToken == ")" {
						break
					} else if strings.ToUpper(currToken) == VALUES {
						break
					} else {
						ast.Columns = append(ast.Columns, currToken)
					}
				}
			}
		} else if currToken != VALUES {
			err = fmt.Errorf("%s expect VALUES or '(' here,error token:%s", s.buffer, s.curr)
			return
		}
	}
	columnCount := len(ast.Columns)
	if columnCount == 0 {
		err = fmt.Errorf("%s expect VALUES or '(' here", s.buffer)
		return
	}
	ast.Values = make([][]string, 0)

rawLoop:
	for {
		if s.next(); s.end {
			break rawLoop
		} else {
			currToken := string(s.curr)
			if currToken == "," {
				continue
			}
			if currToken == "(" {
				row := make([]string, 0, columnCount)
				for {
					if s.next(); s.end {
						break rawLoop
					} else {
						currToken := string(s.curr)
						if currToken == "," {
							continue
						} else if currToken == ")" {
							if len(row) != columnCount {
								err = fmt.Errorf(
									"%s expected column count is %d, got %d, %v",
									s.buffer, columnCount, len(row), row,
								)
								return
							}
							ast.Values = append(ast.Values, row)
							break
						} else {
							row = append(row, currToken)
						}
					}
				}
			}
		}
	}
	return
}
