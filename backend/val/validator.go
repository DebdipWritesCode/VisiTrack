package val

import "fmt"

func ValidateString(value string, minLength int, maxLength int) error {
	n := len(value)

	if n < minLength || n > maxLength {
		return fmt.Errorf("value must be between %d and %d characters long, got %d", minLength, maxLength, n)
	}
	return nil
}

func ValidatePhoneNumber(phoneNumber string) error {
	if len(phoneNumber) != 10 {
		return fmt.Errorf("phone number must be exactly 10 digits long, got %d", len(phoneNumber))
	}

	for _, char := range phoneNumber {
		if char < '0' || char > '9' {
			return fmt.Errorf("phone number must contain only digits, found '%c'", char)
		}
	}
	return nil
}

func ValidateFirstName(firstName string) error {
	if err := ValidateString(firstName, 1, 30); err != nil {
		return err
	}

	for _, char := range firstName {
		if char < 'A' || char > 'Z' {
			if char < 'a' || char > 'z' {
				return fmt.Errorf("first name must contain only alphabetic characters, found '%c'", char)
			}
		}
	}
	return nil
}

func ValidateLastName(lastName string) error {
	if err := ValidateString(lastName, 1, 30); err != nil {
		return err
	}

	for _, char := range lastName {
		if char < 'A' || char > 'Z' {
			if char < 'a' || char > 'z' {
				return fmt.Errorf("last name must contain only alphabetic characters, found '%c'", char)
			}
		}
	}
	return nil
}

func ValidateRole(role string) error {
	if role != "" && role != "admin" && role != "user" {
		return fmt.Errorf("role must be either 'admin' or 'user', got '%s'", role)
	}
	return nil
}

func ValidateID(id int64) error {
	if id <= 0 {
		return fmt.Errorf("ID must be a positive integer, got %d", id)
	}
	return nil
}
