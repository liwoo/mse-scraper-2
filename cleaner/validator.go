package cleaner

import (
	"fmt"
	"reflect"
	"strings"
)

const validateTagName = "validate"

type Validator interface {
	Validate(interface{}) (bool, error)
}

type DefaultValidator struct {
}

func (v DefaultValidator) Validate(val interface{}) (bool, error) {
	return true, nil
}

type StringValidator struct {
	isRequired bool
}

func (v StringValidator) Validate(val interface{}) (bool, error) {
	l := len(val.(string))

	if l == 0 && v.isRequired {
		return false, fmt.Errorf("cannot be blank")
	}

	return true, nil
}

type NumberValidator struct {
	isRequired bool
}

func (v NumberValidator) Validate(val interface{}) (bool, error) {
	num := val.(float64)

	if num == -1 && v.isRequired {
		return false, fmt.Errorf("cannot be blank")
	}

	return true, nil
}

func getValidatorFromTag(tag string) Validator {
	args := strings.Split(tag, ",")

	switch args[0] {
	case "number":
		validator := NumberValidator{}
		fmt.Sscanf(strings.Join(args[1:], ","), "required=%t", &validator.isRequired)
		return validator
	case "string":
		validator := StringValidator{}
		fmt.Sscanf(strings.Join(args[1:], ","), "required=%t", &validator.isRequired)
		return validator
	}
	return DefaultValidator{}
}

func validateStruct(s DailyCompanyRate) []error {
	errs := []error{}

	v := reflect.ValueOf(s)

	for i := 0; i < v.NumField(); i++ {
		tag := v.Type().Field(i).Tag.Get(validateTagName)

		if tag == "" || tag == "-" {
			continue
		}

		validator := getValidatorFromTag(tag)
		valid, err := validator.Validate(v.Field(i).Interface())

		if !valid && err != nil {
			errs = append(errs, fmt.Errorf("%s %s", v.Type().Field(i).Name, err.Error()))
		}
	}

	return errs
}
