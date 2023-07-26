package api

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/KHarshit1203/simple-bank/util"
	"github.com/go-playground/validator/v10"
)

type ValidatonError struct {
	Field   string `json:"error_field"`
	Tag     string `json:"tag"`
	Message string `json:"message"`
}

type ApiValidator struct {
	validator.Validate
}

func NewValidator() ApiValidator {
	validate := validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		queryName := strings.SplitN(fld.Tag.Get("query"), ",", 2)[0]
		if queryName != "" {
			if queryName == "-" {
				return ""
			}
			return queryName
		}

		jsonName := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if jsonName == "-" {
			return ""
		}
		return jsonName
	})

	return ApiValidator{*validate}
}

// validate Reuest Body validates given r
func (av *ApiValidator) validateRequest(request interface{}) []*ValidatonError {
	var errors []*ValidatonError
	if errs := av.Struct(request); errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			validationError := ValidatonError{
				Field:   err.Field(),
				Tag:     err.Tag(),
				Message: fmt.Sprintf("invalid value, %s validation failed on %s", err.ActualTag(), err.Field()),
			}
			errors = append(errors, &validationError)
		}
	}
	return errors
}

// validCurrency is custom validator function to check currency tag
var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		return util.IsCurrencySupported(currency)
	}
	return false
}
