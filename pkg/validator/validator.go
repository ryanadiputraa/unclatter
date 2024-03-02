package validator

import (
	"errors"
	"fmt"
	"time"
	"unicode"

	"github.com/go-playground/validator/v10"
)

type Validator interface {
	Validate(val any) (error, map[string]string)
}

type validation struct {
	validator *validator.Validate
}

func NewValidator() Validator {
	v := validator.New()
	v.RegisterValidation("iso8601date", isISO8601Date)

	return &validation{
		validator: v,
	}
}

func (v *validation) Validate(val any) (error, map[string]string) {
	err := v.validator.Struct(val)
	errorMap := make(map[string]string)

	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, fieldErr := range validationErrors {
			field := fieldToSnakeCase(fieldErr.Field())
			errorMap[field] = FieldErrMsg(fieldErr)
		}
		return errors.New("invalid params"), errorMap
	}
	return err, errorMap
}

func FieldErrMsg(err validator.FieldError) string {
	field := fieldToSnakeCase(err.Field())
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "max":
		return fmt.Sprintf("%s should have a maximum length of %s", field, err.Param())
	case "min":
		return fmt.Sprintf("%s should have a minimum length of %s", field, err.Param())
	case "email":
		return fmt.Sprintf("%s should be a valid email address", field)
	case "http_url":
		return fmt.Sprintf("%s should be a valid http url", field)
	case "iso8601date":
		return fmt.Sprintf("%s should be a valid ISO8601 date", field)
	default:
		return err.Error()
	}
}

func fieldToSnakeCase(input string) string {
	prev := rune(0)
	var result []rune
	for i, char := range input {
		if i > 0 && unicode.IsUpper(char) && unicode.IsLower(prev) {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(char))
		prev = char
	}
	return string(result)
}

func isISO8601Date(fl validator.FieldLevel) bool {
	_, err := time.Parse(time.RFC3339Nano, fl.Field().String())
	return err == nil
}
