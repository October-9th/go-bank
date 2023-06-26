package api

import (
	"github.com/October-9th/simple-bank/util"
	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		// Check currency is supported
		return util.IsSupportedCurrency(currency)
	}
	return false

}
