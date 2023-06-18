package api

import (
	"github.com/KHarshit1203/simple-bank/util"
	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		return util.IsCurrencySupported(currency)
	}
	return false
}