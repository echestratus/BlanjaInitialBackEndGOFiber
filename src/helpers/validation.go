package helpers

import (
	"regexp"
	"unicode"

	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	FailedField string `json:"failedField"`
	Tag         string `json:"tag"`
	Value       string `json:"value"`
}

var validate = validator.New()

func ValidateStruct(param any) []*ErrorResponse {
	var errors []*ErrorResponse
	validate.RegisterValidation("passwordValidation", validatePassword)
	err := validate.Struct(param)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			element := ErrorResponse{
				FailedField: err.StructNamespace(),
				Tag:         err.Tag(),
				Value:       err.Param(),
			}
			errors = append(errors, &element)
		}
	}
	return errors
}

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	hasCapital := false
	hasDigit := false

	for _, v := range password {
		if unicode.IsUpper(v) {
			hasCapital = true
		}
		if unicode.IsDigit(v) {
			hasDigit = true
		}
	}
	charSpecial := regexp.MustCompile(`[_!@#\$%\^&\*]`)

	hasSpecialChar := charSpecial.MatchString(password)

	result := hasCapital && hasDigit && hasSpecialChar

	return result
}
