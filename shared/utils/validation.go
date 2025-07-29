package utils

import (
	"github.com/go-playground/validator/v10"
	"net/mail"
)

func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func ValidationErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fe.Field() + " is required"
	case "email":
		return fe.Field() + " must be a valid email"
	case "min":
		return fe.Field() + " must be at least " + fe.Param() + " characters"
	case "max":
		return fe.Field() + " must be at most " + fe.Param() + " characters"
	case "url":
		return fe.Field() + " must be a valid URL"
	case "position_enum":
		return fe.Field() + " must be a valid position enum: PENYERANG, GELANDANG, BERTAHAN or PENJAGA_GAWANG"
	case "datetime":
		return fe.Field() + " must be a valid datetime. Example: 2006-01-02 15:04:05"
	default:
		return fe.Field() + " is not valid"
	}
}
