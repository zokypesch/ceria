package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInArray(t *testing.T) {
	util := GeneralUtilService()

	t.Run("Test with failed null zero array values", func(t *testing.T) {
		var arr interface{}
		value := "null"

		exist, _ := util.InArray(value, arr)

		assert.Equal(t, false, exist)

	})

	t.Run("Test with return expected", func(t *testing.T) {
		arr := []string{
			"name",
			"age",
			"addr",
		}
		value := "name"
		exist, index := util.InArray(value, arr)

		assert.Equal(t, true, exist)
		assert.Equal(t, 0, index)

		// test with value not found
		existF, indexF := util.InArray("", arr)

		assert.Equal(t, false, existF)
		assert.Equal(t, -1, indexF)
	})
}
