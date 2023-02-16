package gutil

import (
	"reflect"
)

// RIsEmpty 判断给定值是否为空
func RIsEmpty(value interface{}) bool {
	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		return val.Len() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return val.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return val.Float() == 0.0
	case reflect.Bool:
		return !val.Bool()
	case reflect.Ptr, reflect.Interface:
		return val.IsNil()
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			if !RIsEmpty(val.Field(i).Interface()) {
				return false
			}
		}
		return true
	}

	return reflect.DeepEqual(value, reflect.Zero(reflect.TypeOf(value)).Interface())
}

// RGetTypeName 获取给定值的类型名称
func RGetTypeName(value interface{}) string {
	return reflect.TypeOf(value).Name()
}

// RIsPointer 判断给定值是否为指针类型
func RIsPointer(value interface{}) bool {
	return reflect.ValueOf(value).Kind() == reflect.Ptr
}

// RGetFieldsValues 获取给定结构体的所有字段名和字段值
func RGetFieldsValues(obj interface{}) map[string]interface{} {
	fieldsValues := make(map[string]interface{})
	objVal := reflect.ValueOf(obj).Elem()
	objType := objVal.Type()
	for i := 0; i < objVal.NumField(); i++ {
		fieldVal := objVal.Field(i)
		fieldValue := fieldVal.Interface()
		fieldType := objType.Field(i)
		fieldName := fieldType.Name
		fieldsValues[fieldName] = fieldValue
	}
	return fieldsValues
}
