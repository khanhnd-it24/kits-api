package appvalidator

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"unicode"
)

func isStrongPassword(pass string) bool {
	var (
		upp, low, num, sym bool
		tot                uint8
	)

	for _, char := range pass {
		switch {
		case unicode.IsUpper(char):
			upp = true
			tot++
		case unicode.IsLower(char):
			low = true
			tot++
		case unicode.IsNumber(char):
			num = true
			tot++
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			sym = true
			tot++
		default:
			return false
		}
	}

	if !upp || !low || !num || !sym || tot < 6 {
		return false
	}

	return true
}

var strongPasswordValidator validator.Func = func(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	return isStrongPassword(password)
}

func registerStrongPasswordValidator(v *validator.Validate, trans ut.Translator) {
	registerValidator(v, trans, "strong_password", strongPasswordValidator, "{0} must have lower case, upper case, number and symbol")
}
