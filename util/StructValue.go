package util

import (
	"fmt"
	"reflect"
	"strconv"
)

// StructValue as struct for utility
type StructValue struct{}

// StructValueInterface for interfacing function
type StructValueInterface interface {
	SetDefaultValueStruct(str interface{})
}

// RebuildProperty properties rebuild of struct
type RebuildProperty struct {
	IgnoreFieldString []string
	IgnoreFieldType   []reflect.Type
	MoveToMember      []string
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

// RebuilToNewStruct for create new struct and ignore some fields
func (st *StructValue) RebuilToNewStruct(str interface{}, props *RebuildProperty, withValue bool) (interface{}, error) {

	var val reflect.Value

	switch reflect.ValueOf(str).Kind() {
	case reflect.Ptr:
		val = reflect.ValueOf(str).Elem()
	case reflect.Struct:
		val = reflect.ValueOf(str)
	default:
		return nil, fmt.Errorf("Invalid Type")
	}

	var (
		fieldName, finalName string
		fieldValue           reflect.Value
		sf                   []reflect.StructField
		fs                   []string
	)

	fillStruct := make(map[string]reflect.Value)

	general := GeneralUtilService()

	for i := 0; i < val.NumField(); i++ {
		fieldName = val.Type().Field(i).Name
		fieldValue = val.Field(i)
		finalName = fieldName

		if ok, _ := general.InArray(fieldName, props.IgnoreFieldString); ok {
			continue
		}

		if ok, _ := general.InArray(reflect.Indirect(fieldValue).Type(), props.IgnoreFieldType); ok {
			continue
		}

		if ok, _ := general.InArray(finalName, props.MoveToMember); ok &&
			fieldValue.Kind() == reflect.Ptr || fieldValue.Kind() == reflect.Struct {

			for j := 0; j < fieldValue.NumField(); j++ {
				finalName = fieldValue.Type().Field(j).Name
				subValue := fieldValue.Field(j)

				if ok, _ := general.InArray(finalName, fs); ok {
					finalName = fieldValue.Type().Field(j).Name + "_" + fieldName
				}

				sf = append(sf, reflect.StructField{
					Name: finalName,
					Type: fieldValue.Type().Field(j).Type,
				})
				fs = append(fs, finalName)
				fillStruct[finalName] = subValue
			}
			continue
		}

		if ok, _ := general.InArray(finalName, fs); ok {
			finalName = finalName + "_Duplicate_" + strconv.Itoa(i)
		}

		sf = append(sf, reflect.StructField{
			Name: finalName,
			Type: val.Type().Field(i).Type,
		})

		fs = append(fs, finalName)
		fillStruct[finalName] = fieldValue
	}

	newStructField := reflect.StructOf(sf)
	newOfStruct := reflect.New(newStructField).Elem()

	if !withValue {
		return newOfStruct.Interface(), nil
	}

	for k, v := range fillStruct {
		if !reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface()) {
			newOfStruct.FieldByName(k).Set(v)
			continue
		}
	}

	return newOfStruct.Interface(), nil

}
