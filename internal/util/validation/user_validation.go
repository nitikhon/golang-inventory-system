package validation

import (
	"errors"
	"regexp"
	"strings"

	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
)

func ValidateAndNormalizeUser(user *entity.User) (error) {
	// Validate required fields
	if err := validateRequiredFields(user); err != nil {
		return err
	}

	// Validate and normalize email
	if err := validateAndNormalizeEmail(user); err != nil {
		return err
	}

	// Validate and normalize phone number
	if err := validateAndNormalizePhone(user); err != nil {
		return err
	}

	// Validate and normalize username
	if err := validateAndNormalizeUsername(user); err != nil {
		return err
	}

	// Normalize names and password
	normalizeNames(user)
	normalizePassword(user)

	// Set default value for IsAdmin
	user.IsAdmin = false

	return nil
}

func validateRequiredFields(user *entity.User) error {
	// Check if all required fields are provided
	if user.Username == "" {
		return errors.New("username is required")
	}
	if user.Email == "" {
		return errors.New("email is required")
	}
	if user.Password == "" {
		return errors.New("password is required")
	}
	if user.Phone == "" {
		return errors.New("phone is required")
	}
	if user.FirstName == "" {
		return errors.New("first name is required")
	}
	if user.LastName == "" {
		return errors.New("last name is required")
	}
	return nil
}

func validateAndNormalizeEmail(user *entity.User) error {
	// Validate email format
	emailRegex := `^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`
	matched, err := regexp.MatchString(emailRegex, user.Email)
	if err != nil {
		return errors.New("failed to validate email format")
	}
	if !matched {
		return errors.New("invalid email format")
	}

	// Normalize email (lowercase entity part)
	atIndex := strings.LastIndex(user.Email, "@")
	if atIndex != -1 {
		localPart := user.Email[:atIndex]
		entityPart := strings.ToLower(user.Email[atIndex:])
		user.Email = localPart + entityPart
	}

	return nil
}

func validateAndNormalizePhone(user *entity.User) error {
	// Remove non-numeric characters from phone
	phoneRegex := `[^0-9]`
	phoneCleaned := regexp.MustCompile(phoneRegex).ReplaceAllString(user.Phone, "")
	if len(phoneCleaned) < 10 || len(phoneCleaned) > 10 {
		return errors.New("phone number must be at 10 digits")
	}
	user.Phone = phoneCleaned
	return nil
}

func validateAndNormalizeUsername(user *entity.User) error {
	// Validate username format
	usernameRegex := `^[a-zA-Z0-9._]+$`
	matched, err := regexp.MatchString(usernameRegex, user.Username)
	if err != nil {
		return errors.New("failed to validate username format")
	}
	if !matched {
		return errors.New("username can only contain alphabets, numbers, dots, and underscores")
	}

	// Trim spaces from username
	user.Username = strings.TrimSpace(user.Username)

	// Check username length
	if len(user.Username) < 3 || len(user.Username) > 20 {
		return errors.New("username must be 3-20 characters")
	}

	// Ensure username starts with a letter
	if !regexp.MustCompile(`^[a-zA-Z]`).MatchString(user.Username) {
		return errors.New("username must start with a letter")
	}

	// Check for invalid start/end characters (to prevent misinterpret as hiden directory (/home/.john) or command-line flag (ssh -help))
	if strings.HasPrefix(user.Username, ".") || strings.HasPrefix(user.Username, "_") ||
		strings.HasSuffix(user.Username, ".") || strings.HasSuffix(user.Username, "_") {
		return errors.New("username cannot start or end with '.' or '_'")
	}

	// Check for invalid sequences (prevent directory traversal attacks / ambiguity and parsing issues)
	if strings.Contains(user.Username, "..") || strings.Contains(user.Username, "__") ||
		strings.Contains(user.Username, "._") || strings.Contains(user.Username, "_.") {
		return errors.New("username cannot contain '..', '__', '._', or '_.'")
	}

	return nil
}

func normalizeNames(user *entity.User) {
	// Trim spaces from first and last names
	user.FirstName = strings.TrimSpace(user.FirstName)
	user.LastName = strings.TrimSpace(user.LastName)
}

func normalizePassword(user *entity.User) {
	// Trim spaces from password
	user.Password = strings.TrimSpace(user.Password)
}
