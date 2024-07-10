package repl

import (
	"fmt"
	"strconv"
	"strings"
)

type StatementType string

const (
	UNSUPPORTED = "N/A"
	SELECT      = "SELECT"
	DELETE      = "DELETE"
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
	// todo 这个是存在问题的 我们是需要修复的
	//if columnCount == 0 {
	//	err = fmt.Errorf("%s expect VALUES or '(' here,error token:%s", s.buffer, s.curr)
	//	return
	//}
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
				var row []string
				if columnCount != 0 {
					row = make([]string, 0, columnCount)
				} else {
					row = make([]string, 0)
				}
				for {
					if s.next(); s.end {
						break rawLoop
					} else {
						currToken := string(s.curr)
						if currToken == "," {
							continue
						} else if currToken == ")" {
							if columnCount != 0 && len(row) != columnCount {
								err = fmt.Errorf(
									"%s expected column count is %d, got %d, %v",
									s.buffer, columnCount, len(row), row,
								)
								return
							}
							ast.Values = append(ast.Values, row)
							if columnCount == 0 {
								columnCount = len(row)
							}
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

/*
SelectTree 需要来实现select基础功能
* SELECT * FROM foo WHERE id < 3 LIMIT 1;
*/
type SelectTree struct {
	Projects []string
	Table    string   //table_name：要查询的表名称
	Where    []string //column1, column2, ...：要选择的字段名称，可以为多个字段。如果不指定字段名称，则会选择所有字段
	Limit    int64
}

func (s *Scanner) parseSelect() (ast *SelectTree, err error) {
	ast = &SelectTree{}
	//不需要对于select进行处理

	// 直接处理 */project
	ast.Projects = make([]string, 0)
	for {
		if s.next(); s.end {
			if len(ast.Projects) == 0 {
				err = fmt.Errorf("%s get select projects failed", s.buffer)
			}
			return
		} else {
			currToken := strings.ToUpper(string(s.curr))
			// *
			if currToken == ASTERISK {
				ast.Projects = append(ast.Projects, ASTERISK)
			} else {
				if currToken == "," {
					continue
				} else if strings.ToUpper(currToken) == FROM {
					break
				} else {
					ast.Projects = append(ast.Projects, currToken)
				}
			}
		}
	}
	// 获取到table
	if s.next(); s.end {
		return
	} else {
		ast.Table = string(s.curr)
	}
	// 获取到Where这个并不是必要的
	if s.next(); s.end {
		return
	}
	currToken := strings.ToUpper(string(s.curr))
	if currToken == WHERE {
		ast.Where = make([]string, 0)
		for {
			if s.next(); s.end {
				if len(ast.Where) == 0 {
					err = fmt.Errorf("missing WHERE clause")
				}
				return
			}
			currToken := string(s.curr)
			if strings.ToUpper(currToken) == LIMIT {
				break
			}
			ast.Where = append(ast.Where, currToken)
		}
	} else if currToken != LIMIT {
		err = fmt.Errorf("expect WHERE or LIMIT here")
		return
	}

	if s.next(); s.end {
		err = fmt.Errorf("expect LIMIT clause here")
		return
	}
	currToken = string(s.curr)
	ast.Limit, err = strconv.ParseInt(currToken, 10, 64)

	return
}

type DeleteTree struct {
	Table string   //table_name：需要插入新记录的表名
	Where []string //column1, column2, ...：要选择的字段名称，可以为多个字段。如果不指定字段名称，则会选择所有字段
}

func (s *Scanner) parseDelete() (ast *DeleteTree, err error) {
	ast = &DeleteTree{}
	// 不需要进行 解析 insert
	// 解析 into
	if s.next(); s.end || strings.ToUpper(string(s.curr)) != FROM {
		err = fmt.Errorf("%s is not DElETE statement,error token: %s", s.buffer, s.curr)
		return
	}
	// 解析table name
	if s.next(); s.end {
		err = fmt.Errorf("%s expect table after DELETE FROM", s.buffer)
		return
	} else {
		ast.Table = string(s.curr)
	}

	// 获取到Where这个并不是必要的
	if s.next(); s.end {
		return
	}
	currToken := strings.ToUpper(string(s.curr))
	if currToken == WHERE {
		ast.Where = make([]string, 0)
		for {
			if s.next(); s.end {
				if len(ast.Where) == 0 {
					err = fmt.Errorf("missing WHERE clause")
				}
				return
			}
			currToken := string(s.curr)
			if strings.ToUpper(currToken) == LIMIT {
				break
			}
			ast.Where = append(ast.Where, currToken)
		}
	}

	return
}
