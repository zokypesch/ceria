package util

import (
	"gopkg.in/go-playground/validator.v9"
)

// ConfValidator variable
type ConfValidator struct {
	mystruct interface{}
	validate *validator.Validate
}

// ConfValidatorRepo interfaces
type ConfValidatorRepo interface {
	Validate() error
}

// var valid *ConfValidator

// Validate function
func (utilV *ConfValidator) Validate() error {

	err := utilV.validate.Struct(utilV.mystruct)
	return err

}

// NewUtilService service for util
func NewUtilService(props interface{}) *ConfValidator {
	return &ConfValidator{
		validate: validator.New(),
		mystruct: props,
	}
}
