package action

import (
	"errors"
	"strings"
)

type JoinTypeStruct struct {
	JoinType       string
	JoinTableName  string
	JoinColumnName string
}

func parseJoinCondition(join string) ([]JoinTypeStruct, error) {
	// 拆分多个join条件
	joinSlice := strings.Split(join, ",")

	var js []JoinTypeStruct
	for k := range joinSlice {
		// 分析每个join条件
		var j JoinTypeStruct
		err := j.parseSingleJoinCondition(joinSlice[k])
		if err != nil {
			return nil, err
		}

		js = append(js, j)
	}

	return js, nil
}

// 分析join条件
func (j *JoinTypeStruct) parseSingleJoinCondition(joinStr string) error {
	// 分解join的条件
	joinSlice := strings.Split(joinStr, "/")
	if len(joinSlice) != 3 {
		return errors.New("join condition error: " + "\"" + joinStr + "\"")
	}

	//获取需要联合查询的表名
	j.JoinTableName = strings.TrimSpace(joinSlice[1])

	//获取需要联合查询的表的字段
	j.JoinColumnName = strings.TrimSpace(joinSlice[2])

	// 获取 join 条件
	switch strings.TrimSpace(joinSlice[0]) {
	case "<":
		j.JoinType = "LEFT JOIN"
	case ">":
		j.JoinType = "RIGHT JOIN"
	case "&":
		j.JoinType = "INNER JOIN"
	case "|":
		j.JoinType = "FULL JOIN"
	case "":
		j.JoinType = "JOIN"
	case "!":
		j.JoinType = "OUTTER JOIN"
	}

	return nil
}
