package action

import (
	"strings"
)

type OrderStruct struct {
	ColumnName string
	OrderBy    string
}

// 解析 排序方式
func parseOrderCondition(orderStr string) []OrderStruct {
	orderSlice := strings.Split(orderStr, ",")

	var os []OrderStruct
	for k := range orderSlice {
		var o OrderStruct
		o.parseSingleOrderCondition(orderSlice[k])

		os = append(os, o)
	}
	return os
}

// 解析单个排序
func (order *OrderStruct) parseSingleOrderCondition(orderStr string) {
	if strings.Contains(orderStr, "+") {
		orderSlice := strings.Split(orderStr, "+")
		order.ColumnName = strings.TrimSpace(orderSlice[0])
		order.OrderBy = "ASC"
	} else if strings.Contains(orderStr, "-") {
		orderSlice := strings.Split(orderStr, "-")
		order.ColumnName = strings.TrimSpace(orderSlice[0])
		order.OrderBy = "DESC"
	} else {
		order.ColumnName = strings.TrimSpace(orderStr)
		order.OrderBy = "ASC"
	}
}
