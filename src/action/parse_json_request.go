package action

import (
	"encoding/json"
	"errors"
	"log"
	"strings"
)

type QueryConditionStruct struct {
	Page  int              // 第几页
	Count int              // 每页多少条数据
	Join  []JoinTypeStruct // 连表查询方式
	Query int
}

// 解析 APIJSON 请求
func ParseJsonRequest(req []byte) (map[string]interface{}, error) {
	reqMap := make(map[string]interface{})
	err := json.Unmarshal(req, &reqMap)
	if err != nil {
		return nil, err
	}
	return reqMap, err
}

// TODO 根据不同请求类型来实现不同业务
func parseReqType(reqMap map[string]interface{}, reqType string) error {
	switch reqType {
	case "get":
		//普通获取数据
	case "head":
		//普通获取数量
	case "gets":
		//安全/私密获取数据
	case "heads":
		//安全/私密获取数量
	case "post":
		//新增数据
	case "put":
		//修改数据，
		//只修改所传的字段
	case "delete":
		//删除数据
	}
	return nil
}

// 获取不同的 tableName
func ParseTableName(reqMap map[string]interface{}) error {
	for tableName, condition := range reqMap {
		//	判断 tableName 是否包含 []
		if strings.Contains(tableName, "[]") {
			// 判断slice条件
			c, ok := condition.(map[string]interface{})
			if !ok {
				return errors.New("condition error")
			}
			var qc QueryConditionStruct
			err := qc.parseGetSliceCondition(c)
			if err != nil {
				return err
			}
			log.Println(qc)
		} else {
			var table TableStruct
			table.TableName = tableName
			c, ok := condition.(map[string]interface{})
			if !ok {
				return errors.New("condition error")
			}
			log.Println(tableName)
			log.Println(c)
			//	TODO column, order
		}
	}
	return nil
}

// 分析 get，gets 中 Slice 的请求数据,
// page, count，join, query
func (q *QueryConditionStruct) parseGetSliceCondition(condition map[string]interface{}) error {
	// 判断是否有page条件
	if v, ok := condition["page"]; ok {
		page, ok := v.(float64)
		if !ok {
			return errors.New("page error: not a valid Int type")
		}
		q.Page = int(page)
	}

	// 判断是否有count条件
	if v, ok := condition["count"]; ok {
		count, ok := v.(float64)
		if !ok {
			return errors.New("count error: not a valid Int type")
		}
		q.Count = int(count)
	}

	// 判断是否有join条件
	if v, ok := condition["join"]; ok {
		join, ok := v.(string)
		if !ok {
			return errors.New("join error: not a valid String type")
		}

		// 分析join条件
		js, err := parseJoinCondition(join)
		if err != nil {
			return err
		}
		q.Join = js
	}

	// TODO query 方法

	return nil
}
