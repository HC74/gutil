package gutil

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"reflect"
	"strings"
)

// DataSlice 数据分片
func DataSlice[T any](v []T, size int) [][]T {
	count := len(v) / size
	var results [][]T
	var temp = 0
	var tempSize = size
	for i := 0; i < count; i++ {
		results = append(results, v[temp:tempSize])
		temp += size
		tempSize += size
	}
	if len(v)%size != 0 {
		results = append(results, v[temp:])
	}
	return results
}

// Md5 MD5
func Md5(str string) string {
	hash := md5.Sum([]byte(str))
	md5str := hex.EncodeToString(hash[:])
	return md5str
}

// ToSqlCondition 转换为SQL条件
func ToSqlCondition(tableName string, v interface{}) map[string]interface{} {
	if v == nil {
		return nil
	}
	m := make(map[string]interface{})
	value := reflect.ValueOf(v)
	for value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if value.Kind() == reflect.Struct {
		tp := value.Type()
		num := value.NumField()

		for i := 0; i < num; i++ {
			fieldValue := value.Field(i)
			fieldType := tp.Field(i)

			if isBlank(fieldValue) {
				continue
			}

			if fieldValue.Kind() == reflect.Slice {
				m[tableName+"."+fieldType.Name+" in ?"] = fieldValue.Interface()
				continue
			}

			sqlTag := fieldType.Tag.Get("sql")
			if sqlTag == "-" {
				continue
			}

			field := fieldType.Tag.Get("json")
			if fieldType.Tag.Get("field") != "" {
				field = fieldType.Tag.Get("field")
			}

			if sqlTag == "" {
				m[tableName+"."+field+" = ?"] = fieldValue.Interface()
			} else {
				sqlTag = strings.ReplaceAll(sqlTag, "?", fmt.Sprint(fieldValue.Interface()))
				m[tableName+"."+field+" "+sqlTag] = nil
			}
		}
	}

	return m
}

// isBlank 条件为空
func isBlank(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}
