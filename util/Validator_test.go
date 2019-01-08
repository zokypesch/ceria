package util

import (
	"fmt"
	"testing"

	"github.com/zokypesch/ceria/util/mocks"

	"github.com/stretchr/testify/assert"
)

type UserTest struct {
	UserID string
	Name   string
	Age    string `validate:"required"`
	err    []error
}

func TestValidator(t *testing.T) {
	assert := assert.New(t)

	t.Run("Check positive case", func(t *testing.T) {

		mockValidation := new(mocks.ConfValidatorRepo)

		mockValidation.On("Validate").Return(nil)

		actValidator := NewUtilService(UserTest{
			UserID: "27",
			Name:   "FirstLast",
			Age:    "30",
		})

		assert.Equal(actValidator.Validate(), mockValidation.Validate())
		mockValidation.AssertExpectations(t)

		mockValidation.AssertNumberOfCalls(t, "Validate", 1)
	})

	t.Run("Check negative case", func(t *testing.T) {
		mockValidation := new(mocks.ConfValidatorRepo)

		mockValidation.On("Validate").Return(fmt.Errorf("Key: 'UserTest.Age' Error:Field validation for 'Age' failed on the 'required' tag"))

		actValidator := NewUtilService(UserTest{
			UserID: "27",
			Name:   "FirstLast",
			Age:    "",
		})

		errExpec := mockValidation.Validate()
		errAct := actValidator.Validate()

		assert.IsType(errExpec.Error(), errAct.Error())

		assert.Equal(errExpec.Error(), errAct.Error())

		mockValidation.AssertExpectations(t)

		mockValidation.AssertNumberOfCalls(t, "Validate", 1)

	})

}
