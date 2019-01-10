package util

import (
	"testing"

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
}
