package util

import (
	"reflect"
	"testing"

	"github.com/jinzhu/gorm"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Name string `default:"udin"`
	Age  int    `default:"100"`
}

var sInter StructValueInterface

func TestFillDefaultStruct(t *testing.T) {
	sInter := NewServiceStructValue()

	t.Run("Test with fill expectation", func(t *testing.T) {
		newStr := &testStruct{}

		sInter.SetDefaultValueStruct(newStr)

		assert.Equal(t, "udin", newStr.Name)
		assert.Equal(t, 100, newStr.Age)
	})

	t.Run("Test with no default expectation", func(t *testing.T) {
		newStr := &testStruct{Name: "Paijo", Age: 79}

		sInter.SetDefaultValueStruct(newStr)

		assert.NotEqual(t, "udin", newStr.Name)
		assert.NotEqual(t, 100, newStr.Age)
	})

	t.Run("Test with invalid struct", func(t *testing.T) {
		newParams := 5

		sInter.SetDefaultValueStruct(newParams)
		// do nothing and not panic error
	})

	t.Run("Tes Get Name failed params of struct", func(t *testing.T) {
		newParams := 5

		str := sInter.GetNameOfStruct(newParams)
		assert.Equal(t, str, "")

		str2 := sInter.GetNameOfStruct(&newParams)
		assert.Equal(t, str2, "int")

		// test of struct
		type TestMyStructName struct {
			Name string
		}

		assert.Equal(t, "TestMyStructName", sInter.GetNameOfStruct(TestMyStructName{}))
		assert.Equal(t, "TestMyStructName", sInter.GetNameOfStruct(&TestMyStructName{}))

	})

	t.Run("Tes Get Name failed params of struct", func(t *testing.T) {
		newParams := &testStruct{Name: "Paijo", Age: 79}

		str := sInter.GetNameOfStruct(newParams)
		assert.NotEmpty(t, str)
	})

	t.Run("Tes Run set nil value of interface", func(t *testing.T) {

		newParams := &testStruct{Name: "Paijo", Age: 79}

		str := sInter.SetNilValue(newParams)
		assert.Empty(t, str)

	})

	t.Run("Tes Run set nil value of address", func(t *testing.T) {

		newParams := testStruct{Name: "Paijo", Age: 79}

		str := sInter.SetNilValue(newParams)
		assert.Empty(t, str)

	})

	t.Run("Tes Run Rebuild struct failed with invalid params", func(t *testing.T) {
		str, err := sInter.RebuilToNewStruct(nil, &RebuildProperty{}, false)
		assert.Error(t, err)
		assert.Empty(t, str)

	})

	t.Run("Tes Run Rebuild struct failed with not struct params", func(t *testing.T) {
		str, err := sInter.RebuilToNewStruct("test", &RebuildProperty{}, false)
		assert.Error(t, err)
		assert.Empty(t, str)

	})

	t.Run("Tes Run Rebuild struct ecpected with struct params", func(t *testing.T) {

		type ParentStructParams struct {
			Email string
		}
		type NewParamStruct struct {
			gorm.Model
			ID      uint
			Name    string
			Age     string
			Melotot string
			Emails  []ParentStructParams
		}

		np := &NewParamStruct{
			gorm.Model{ID: 1000},
			9,
			"Udin", "joss", "iyes", []ParentStructParams{ParentStructParams{"udin@gmail.com"}}}

		var ignore []reflect.Type

		ignore = append(ignore, reflect.TypeOf(&ParentStructParams{}))

		// var newParamStruct
		str, err := sInter.RebuilToNewStruct(
			np,
			&RebuildProperty{IgnoreFieldString: []string{"Melotot"},
				IgnoreFieldType: ignore,
				MoveToMember:    []string{"Model"},
			},
			true,
		)
		assert.NoError(t, err)
		assert.NotEmpty(t, str)

	})
}
