package utils

import (
	"net/mail"
	"unicode"
)

func ValidatePassword(pass string) bool {
	var (
		upper, lower, number, symbol bool
		tot                          uint8
	)

	for _, char := range pass {
		switch {
		case unicode.IsUpper(char):
			upper = true
			tot++
		case unicode.IsLower(char):
			lower = true
			tot++
		case unicode.IsNumber(char):
			number = true
			tot++
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			symbol = true
			tot++
		default:
			return false
		}
	}

	if !upper || !lower || !number || !symbol || tot < 8 {
		return false
	}

	return true
}

func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
