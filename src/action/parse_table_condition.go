package action

import (
	"errors"
	"strings"
)

type TableStruct struct {
	TableName string                 // 表名
	Column    []string               // 需要获取的字段
	Order     []OrderStruct          // 排序
	Search    map[string]interface{} // 搜索条件
}

func (t *TableStruct) parseTableCondition(condition map[string]interface{}) error {
	// column 条件
	if v, ok := condition["@column"]; ok {
		columnStr, ok := v.(string)
		if !ok {
			return errors.New("order error")
		}
		column := strings.Split(columnStr, ",")
		t.Column = column
	}

	// order 条件
	if v, ok := condition["@order"]; ok {
		orderStr, ok := v.(string)
		if !ok {
			return errors.New("order error")
		}
		t.Order = parseOrderCondition(orderStr)
	}

	return nil
}
