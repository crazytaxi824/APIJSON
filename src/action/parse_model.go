package action

import (
	"errors"
	"log"
	"reflect"
	"strings"
)

type StructFieldMap struct {
	FieldName string
	TagName   string
}

// 表名-json字段名-sql字段名&结构体属性名
var TableJsonStructMap map[string]map[string]StructFieldMap

// 结构体类型名-表名
var TypeToTablenameMap map[string]string

// 结构体类型名-结构体
var TypeToStruct map[string]interface{}

func init() {
	TableJsonStructMap = make(map[string]map[string]StructFieldMap)
	TypeToStruct = make(map[string]interface{})
	TypeToTablenameMap = make(map[string]string)
}

// 传入实体空结构体Slice, 缓存到 TableJsonStructMap
func CacheStructSlice(models []interface{}) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	for _, model := range models {
		tmp := make(map[string]StructFieldMap)

		// 获取结构体类型
		typ := reflect.TypeOf(model)
		count := typ.NumField()
		r, ok := typ.FieldByName("tableName")
		if !ok {
			continue
		}

		TypeToStruct[typ.String()] = model

		// 获取tableName
		tableName := r.Tag.Get("sql")

		TypeToTablenameMap[typ.String()] = tableName

		for i := 0; i < count; i++ {
			var t StructFieldMap
			t.FieldName = typ.Field(i).Name

			tagName := typ.Field(i).Tag.Get("sql")
			if strings.Contains(tagName, ",") {
				tmpTagName := strings.Split(tagName, ",")
				t.TagName = strings.TrimSpace(tmpTagName[0])
			} else {
				t.TagName = strings.TrimSpace(tagName)
			}
			if t.TagName == "-" || t.TagName == "" {
				continue
			}

			keyName := typ.Field(i).Tag.Get("json")
			if strings.Contains(keyName, ",") {
				tmpKeyName := strings.Split(keyName, ",")
				keyName = strings.TrimSpace(tmpKeyName[0])
			} else {
				keyName = strings.TrimSpace(keyName)
			}
			if keyName == "-" || keyName == "" {
				continue
			}

			tmp[keyName] = t
		}

		TableJsonStructMap[tableName] = tmp
	}
	return
}

// 传入实体空结构体，缓存到 TableJsonStructMap
func CacheStruct(model interface{}) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	tmp := make(map[string]StructFieldMap)

	// 获取结构体类型
	typ := reflect.TypeOf(model)

	// 获取tableName
	count := typ.NumField()
	r, ok := typ.FieldByName("tableName")
	if !ok {
		return
	}

	TypeToStruct[typ.String()] = model

	tableName := r.Tag.Get("sql")

	// 缓存tablename
	TypeToTablenameMap[typ.String()] = tableName

	for i := 0; i < count; i++ {
		var t StructFieldMap
		t.FieldName = typ.Field(i).Name

		tagName := typ.Field(i).Tag.Get("sql")
		if strings.Contains(tagName, ",") {
			tmpTagName := strings.Split(tagName, ",")
			t.TagName = strings.TrimSpace(tmpTagName[0])
		} else {
			t.TagName = strings.TrimSpace(tagName)
		}
		if t.TagName == "-" || t.TagName == "" {
			continue
		}

		keyName := typ.Field(i).Tag.Get("json")
		if strings.Contains(keyName, ",") {
			tmpKeyName := strings.Split(keyName, ",")
			keyName = strings.TrimSpace(tmpKeyName[0])
		} else {
			keyName = strings.TrimSpace(keyName)
		}
		if keyName == "-" || keyName == "" {
			continue
		}

		tmp[keyName] = t
	}
	TableJsonStructMap[tableName] = tmp
	return
}

// 检查传入字段是否合法,同时转换成sql字段
func CheckColumn(s []string, tableName string) (column []string, err error) {
	tmpStructMap, ok := TableJsonStructMap[tableName]
	if !ok {
		return nil, errors.New("table is not available")
	}

	for _, c := range s {
		mark := false
		for k := range tmpStructMap {
			if c == k {
				column = append(column, tmpStructMap[c].TagName)
				mark = true
			}
		}

		if !mark {
			return nil, errors.New("\"" + tableName + "\" field: \"" + c + "\" is not available")
		}
	}

	return column, nil
}

// 传指针，将pg获取到的结构体的值传入map, struct to map
func StructToMap(s interface{}, column []string) (map[string]interface{}, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	// 判断column是否为空
	if len(column) < 1 {
		return nil, nil
	}

	result := make(map[string]interface{})

	// 反射结构体，获取tableName
	value := reflect.ValueOf(s).Elem()
	r, ok := value.Type().FieldByName("tableName")
	if !ok {
		return nil, errors.New("no table")
	}
	tableName := r.Tag.Get("sql")

	// 遍历column赋值到map
	for _, v := range column {
		result[v] = value.FieldByName(TableJsonStructMap[tableName][v].FieldName).Interface()
	}

	return result, nil
}

// 传指针，将pg获取到的结构体Slice的值传入 []map, []struct to []map
func StructSliceToMap(s interface{}, column []string) ([]map[string]interface{}, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	// 判断column是否为空
	if len(column) < 1 {
		return nil, nil
	}

	// 反射slice
	valueSlice := reflect.ValueOf(s).Elem()

	// 获取slice长度
	valueLen := valueSlice.Len()

	// 判断slice是否为空
	if valueLen < 1 {
		return nil, nil
	}

	// 获取tableName
	r, ok := valueSlice.Index(0).Type().FieldByName("tableName")
	// 如果tableName不存在，返回error
	if !ok {
		return nil, errors.New("no table")
	}

	// 获取tableName
	tableName := r.Tag.Get("sql")

	// 遍历valueSlice,赋值给[]map
	finalResult := make([]map[string]interface{}, 0)
	for i := 0; i < valueLen; i++ {
		result := make(map[string]interface{})
		// 遍历column,赋值给map
		for _, v := range column {
			result[v] = valueSlice.Index(i).FieldByName(TableJsonStructMap[tableName][v].FieldName).Interface()
		}
		finalResult = append(finalResult, result)
	}

	return finalResult, nil
}
