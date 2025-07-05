package utils

import "strings"

func SanitizePhoneNumber(phoneNumber string) string {
	// Change from 081234567890 to 6281234567890
	if strings.HasPrefix(phoneNumber, "0") {
		phoneNumber = "62" + phoneNumber[1:]
	}

	// Remove all non-numeric characters
	phoneNumber = strings.ReplaceAll(phoneNumber, " ", "")
	phoneNumber = strings.ReplaceAll(phoneNumber, "-", "")

	// Add + if there's no + in the phone number
	if !strings.HasPrefix(phoneNumber, "+") {
		phoneNumber = "+" + phoneNumber
	}

	return phoneNumber
}
