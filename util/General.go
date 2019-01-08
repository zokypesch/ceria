package util

import (
	"reflect"
)

//GeneralUtilInterface for interfacing struct
type GeneralUtilInterface interface {
	InArray(val interface{}, array interface{}) (exists bool, index int)
}

// GeneralUtil struct defined
type GeneralUtil struct {
	GeneralUtilInterface
}

var utility *GeneralUtil

// GeneralUtilService for get general util
func GeneralUtilService() *GeneralUtil {
	if utility == nil {
		utility = &GeneralUtil{}
	}
	return utility
}

// InArray this func to check array value is exist
func (utils *GeneralUtil) InArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	//validate the interface
	if reflect.ValueOf(array).Kind() == reflect.Invalid {
		return
	}

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}

	return
}
