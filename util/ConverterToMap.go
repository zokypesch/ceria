package util

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// ConverterToMap struct define for static
type ConverterToMap struct{}

var util *ConverterToMap

// NewUtilConvertToMap new service util
func NewUtilConvertToMap() *ConverterToMap {
	if util != nil {
		util = &ConverterToMap{}
	}

	return util
}

// ConvertToDynamicMap for convert to dynamic map
func (util *ConverterToMap) ConvertToDynamicMap(columns []string, values []interface{}) map[string]string {
	newMap := make(map[string]string)
	for i, col := range columns {

		val := values[i]

		valueToString := util.ConvertDataToString(val)
		newMap[col] = valueToString
	}
	return newMap
}

// ConvertMultiStructToMap for convert multi struct to map
func (util *ConverterToMap) ConvertMultiStructToMap(param interface{}) []map[string]interface{} {
	var listArr []map[string]interface{}

	switch reflect.TypeOf(param).Kind() {
	case reflect.Slice, reflect.Ptr:
		var s reflect.Value

		if reflect.TypeOf(param).Kind() == reflect.Ptr {
			s = reflect.ValueOf(param).Elem()
		} else if reflect.TypeOf(param).Kind() == reflect.Slice {
			s = reflect.ValueOf(param)
		}

		for i := 0; i < s.Len(); i++ {
			st := s.Index(i) // get value from index the result is reflect.Value
			// valueOfSt := reflect.ValueOf(st)
			// valueOfIndirect := reflect.Indirect(valueOfSt).Interface()

			// fmt.Println()
			// fmt.Println(valueOfIndirect)
			// fmt.Println(reflect.TypeOf(st))
			// list := make(map[string]interface{}, st.NumField())

			// if reflect.TypeOf(st).Kind() == reflect.Struct {
			// elem := st
			// newSt := reflect.Indirect(elem).Elem()
			// for k := 0; k < newSt.NumField(); k++ {
			// fieldName := newSt.Type().Field(k).Name
			// valueByField := newSt.Field(k)
			// fmt.Println(fieldName, valueByField)
			// list[fieldName] = util.ConvertDataToString(valueByField)
			// }
			// listArr = append(listArr, list)
			// }

			list := util.ConvertStructToSingeMap(st.Interface())
			listArr = append(listArr, list)

			// reflect.ValueOf(st).Interface()
			// list := util.ConvertStructToSingeMap(st) // bisa juga tapi ngaco karena dia tipenya reflect value interfacin dl donk

		}
	}

	return listArr
}

// ConvertDataToString anytype to string
func (util *ConverterToMap) ConvertDataToString(field interface{}) string {

	var data string

	var v interface{}

	b, ok := field.([]byte)

	if ok {
		v = string(b)
	} else {
		v = field
	}

	switch v.(type) {
	case string:
		data = v.(string)
	case int32, int64, int:
		data = strconv.Itoa(v.(int))
	case uint:
		conv := v.(uint)
		newconv := uint64(conv)
		data = strconv.FormatUint(newconv, 10)
	case time.Time:
		conv := v.(time.Time)
		data = conv.String()
	case *time.Time:
		conv := v.(*time.Time)
		if conv == nil {
			data = ""
		} else if conv != nil {
			data = conv.String()
		}
	default:
		if ok := v.(reflect.Value).CanInterface(); ok {
			mischef := v.(reflect.Value)
			return util.ConvertDataToString(mischef.Interface())
		}
		if v.(reflect.Value).Kind() == reflect.Invalid {
			data = "notset"
		}

		data = v.(reflect.Value).String() // get reflect value and go to string

	}
	return data
}

// ConvertInterfaceToKeyStr convert to string array
func (util *ConverterToMap) ConvertInterfaceToKeyStr(fieldList interface{}) []string {
	field := reflect.ValueOf(fieldList)

	list := make([]string, field.NumField())

	for i := 0; i < field.NumField(); i++ {
		list[i] = field.Type().Field(i).Name
	}

	return list

}

// ConvertStructToSingeMap this func for get result of single map
func (util *ConverterToMap) ConvertStructToSingeMap(fieldList interface{}) map[string]interface{} {

	var field reflect.Value
	var res map[string]interface{}

	switch reflect.ValueOf(fieldList).Kind() {
	case reflect.Ptr:
		field = reflect.ValueOf(fieldList).Elem()
	case reflect.Struct:
		field = reflect.ValueOf(fieldList)
	// case reflect.Invalid:
	// 	return res
	default:
		return res
	}

	// newfield := reflect.Indirect(field)
	// fmt.Println(newfield.Interface(), field.Interface(), field)
	// fmt.Println()
	res = make(map[string]interface{}, field.NumField())

	for i := 0; i < field.NumField(); i++ {
		fieldValue := field.Field(i)
		fieldName := util.ConvertDataToString(field.Type().Field(i).Name)

		// check the fieldValue kind only accept slice and struct
		switch fieldValue.Kind() {
		case reflect.Struct, reflect.Ptr:
			// if ok := fieldValue.CanInterface(); ok {
			// newValue = fieldValue.Interface()
			// }
			res[fieldName] = util.RefValueToInterface(fieldValue)
		case reflect.Slice:
			newSlice := make([]interface{}, fieldValue.Len())

			for i := 0; i < fieldValue.Len(); i++ {
				st := fieldValue.Index(i) // get value from index the result is reflect.Value

				if st.Kind() != reflect.Struct {
					newSlice[i] = st
					continue
				}

				newSlice[i] = util.RefValueToInterface(st)
			}

			res[fieldName] = newSlice
		case reflect.Map:
			mapRes := make(map[string]string, len(fieldValue.MapKeys()))
			for _, e := range fieldValue.MapKeys() {
				v := fieldValue.MapIndex(e)
				mapRes[e.String()] = util.ConvertDataToString(v)
			}
			res[fieldName] = mapRes
		default:
			if fieldValue.CanInterface() {
				res[fieldName] = util.ConvertDataToString(fieldValue.Interface())
			} else {
				res[fieldName] = fieldValue.String()
			}

		}

	}

	return res
}

// RefValueToInterface convert reflect.value for iterature to interface
func (util *ConverterToMap) RefValueToInterface(v reflect.Value) interface{} {
	var res interface{}

	if v.Kind() == reflect.Invalid || (v.Interface() == reflect.Zero(v.Type()).Interface() && v.Type().Name() == "") {
		return res
	}

	typ := v.Type()
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		v = v.Elem()
	}

	// Only structs are supported
	if typ.Kind() != reflect.Struct {
		return res
	}

	// is valid struct goes here
	newRes := make(map[string]interface{}, typ.NumField())

	for i := 0; i < typ.NumField(); i++ {
		p := typ.Field(i)
		fieldValue := v.Field(i)
		fieldName := util.ConvertDataToString(v.Type().Field(i).Name)

		if !p.Anonymous {
			// fmt.Println(p, fieldValue, fieldName, v.Type().Field(i).Type, )
			newRes[fieldName] = util.ConvertDataToString(fieldValue)
		} else { // Anonymus structues
			util.RefValueToInterface(v.Field(i).Addr())
		}
	}

	return newRes
}

// ConvertInterfaceMaptoMap for convert interfacemap to map
func (util *ConverterToMap) ConvertInterfaceMaptoMap(inter interface{}) (map[string]string, error) {
	var mp map[string]string
	var err error
	err = fmt.Errorf("FailedToConvert")

	val := reflect.ValueOf(inter)

	if val.Kind() != reflect.Map {
		return mp, err
	}

	mp = make(map[string]string, len(val.MapKeys()))
	for _, e := range val.MapKeys() {
		v := val.MapIndex(e)
		mp[e.String()] = util.ConvertDataToString(v)
	}

	return mp, nil
}

// RebuildToSlice for rebuild interface to struct
func (util *ConverterToMap) RebuildToSlice(inter interface{}) interface{} {
	var res interface{}

	typ := reflect.TypeOf(inter)
	vl := reflect.ValueOf(inter)
	// reflect.Indirect(vl) to use see a truth of type
	// fmt.Println(typ.Kind(), reflect.Indirect(vl).Kind(), vl.Kind(), vl.Interface()) // prints struct)
	// fmt.Println(reflect.Indirect(vl).Kind(), typ.Kind())

	if reflect.Indirect(vl).Kind() != reflect.Struct {
		return res
	}
	// 	if typ.Kind() != reflect.Struct || typ.Kind() != reflect.Ptr || reflect.Indirect(vl).Kind() != reflect.Slice {
	// 		return res
	// 	}

	slice := reflect.MakeSlice(reflect.SliceOf(typ), 0, 0)
	x := reflect.New(slice.Type()) // new reflect
	x.Elem().Set(slice)            // and set it to slice mode

	// conValues := reflect.ValueOf(slice).Interface() // get interface value of st
	// stAddr := conValues.(reflect.Value).Addr()      // st.Addr() // get address like &Model{}
	// realModel := stAddr.Interface()

	return x.Interface() // return as interface
}

// SetFieldNullByTag function for set struct field to PTR
func (util *ConverterToMap) SetFieldNullByTag(inter interface{}) interface{} {
	var res interface{}

	typ := reflect.TypeOf(inter).Kind()

	if typ != reflect.Struct && typ != reflect.Ptr {
		return res
	}

	var vl reflect.Value
	var prop reflect.Type

	switch typ {
	case reflect.Ptr:
		vl = reflect.ValueOf(inter).Elem()
		prop = reflect.TypeOf(inter).Elem()
	case reflect.Struct:
		vl = reflect.ValueOf(inter)
		prop = reflect.TypeOf(inter)
	}

	// initialize new struct

	for i := 0; i < vl.NumField(); i++ {
		field := vl.Field(i)
		vProp, _ := prop.FieldByName(vl.Type().Field(i).Name)

		tag, ok := vProp.Tag.Lookup("ceria")

		if field.CanSet() && ok && tag == "ignore_elastic" {
			// reflect.Indirect(field) result indirect is not pointer but existing value
			// reflect.New(reflect.TypeOf(field)) result is pointer
			// reflect.New(field.Type()).Elem() result is normal ->
			// because struct is pointer and need to elem() to get original value
			newStruct := reflect.New(field.Type()).Elem()
			field.Set(newStruct)
		}

	}

	res = vl.Interface()

	if vl.CanAddr() {
		res = vl.Addr().Interface()
	}

	return res
}
