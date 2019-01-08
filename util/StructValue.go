package util

import (
	"reflect"
	"strconv"
)

// StructValue as struct for utility
type StructValue struct{}

// StructValueInterface for interfacing function
type StructValueInterface interface {
	SetDefaultValueStruct(str interface{})
}

var sValue *StructValue

// NewServiceStructValue function for new structValueService
func NewServiceStructValue() *StructValue {
	if sValue == nil {
		sValue = &StructValue{}
	}
	return sValue
}

// SetDefaultValueStruct for set default value struct by tagging default
func (st *StructValue) SetDefaultValueStruct(str interface{}) {
	var value reflect.Value
	var prop reflect.Type

	typeStr := reflect.TypeOf(str).Kind()

	switch typeStr {
	case reflect.Struct:
		value = reflect.ValueOf(str)
		prop = reflect.TypeOf(str)

	case reflect.Ptr:
		value = reflect.ValueOf(str).Elem()
		prop = reflect.TypeOf(str).Elem()
	default:
		return
	}

	if value.Kind() == reflect.Slice {
		return
	}

	for i := 0; i < value.NumField(); i++ {
		fieldName := value.Type().Field(i).Name
		fieldValue := value.Field(i)
		v, _ := prop.FieldByName(fieldName)

		switch fieldValue.Kind() {
		case reflect.Int:
			tag, ok := v.Tag.Lookup("default")
			canS := fieldValue.CanSet()
			if fieldValue.Interface().(int) == 0 && ok && canS {
				setVal, _ := strconv.Atoi(tag)
				fieldValue.SetInt(int64(setVal))
			}
		case reflect.String:
			tag, ok := v.Tag.Lookup("default")
			canS := fieldValue.CanSet()
			if fieldValue.Interface().(string) == "" && ok && canS {
				fieldValue.SetString(tag)
			}
		}
	}

}

// GetNameOfStruct to get name of struct
func (st *StructValue) GetNameOfStruct(str interface{}) string {
	var structName string
	typ := reflect.TypeOf(str)

	switch typ.Kind() {
	case reflect.Struct:
		structName = typ.Name()
	case reflect.Ptr:
		structName = typ.Elem().Name() // if type is ptr you must get Elem first
	}

	return structName
}

// SetNilValue to set nil and return new interface
func (st *StructValue) SetNilValue(str interface{}) interface{} {
	var newValue interface{}

	typ := reflect.TypeOf(str).Kind()
	val := reflect.ValueOf(str)

	if typ == reflect.Ptr {
		val = reflect.Indirect(val)
	}

	newValue = reflect.New(val.Type()).Interface()

	return newValue
}
