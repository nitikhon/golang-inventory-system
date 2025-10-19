package validation

import (
	"errors"
	"regexp"
	"strings"

	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
)

// ValidateUser validates all user fields without modifying them
func ValidateUser(user *entity.User) error {
	// Validate required fields
	if err := ValidateRequiredFields(user); err != nil {
		return err
	}

	// Validate email format
	if err := ValidateEmail(user.Email); err != nil {
		return err
	}

	// Validate phone number
	if err := ValidatePhone(user.Phone); err != nil {
		return err
	}

	// Validate username
	if err := ValidateUsername(user.Username); err != nil {
		return err
	}

	// Validate password
	if err := ValidatePassword(user.Password); err != nil {
		return err
	}

	// Validate names
	if err := ValidateNames(user.FirstName, user.LastName); err != nil {
		return err
	}

	return nil
}

// NormalizeUser normalizes all user fields
func NormalizeUser(user *entity.User) {
	// Normalize email
	user.Email = NormalizeEmail(user.Email)

	// Normalize phone
	user.Phone = NormalizePhone(user.Phone)

	// Normalize username
	user.Username = NormalizeUsername(user.Username)

	// Normalize names
	user.FirstName, user.LastName = NormalizeNames(user.FirstName, user.LastName)

	// Normalize password
	user.Password = NormalizePassword(user.Password)

	// Set default values
	user.IsAdmin = false
}

// ValidateAndNormalizeUser validates and normalizes user data (convenience function)
func ValidateAndNormalizeUser(user *entity.User) error {
	// First normalize the data
	NormalizeUser(user)

	// Then validate the normalized data
	if err := ValidateUser(user); err != nil {
		return err
	}

	return nil
}

func ValidateRequiredFields(user *entity.User) error {
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

func ValidateEmail(email string) error {
	if len(email) < 6 {
		return errors.New("email must be at least 6 characters long")
	}

	if len(email) > 128 {
		return errors.New("email must not exceed 128 characters")
	}

	// Validate email format
	emailRegex := `^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`
	matched, err := regexp.MatchString(emailRegex, email)
	if err != nil {
		return errors.New("failed to validate email format")
	}
	if !matched {
		return errors.New("invalid email format")
	}

	return nil
}

func NormalizeEmail(email string) string {
	// Normalize email (lowercase domain part)
	atIndex := strings.LastIndex(email, "@")
	if atIndex != -1 {
		localPart := email[:atIndex]
		domainPart := strings.ToLower(email[atIndex:])
		return localPart + domainPart
	}
	return email
}

func ValidateAndNormalizeEmail(email string) (string, error) {
	// First normalize
	normalizedEmail := NormalizeEmail(email)

	// Then validate
	err := ValidateEmail(normalizedEmail)
	return normalizedEmail, err
}

func ValidatePhone(phone string) error {
	// Check if phone contains only numeric characters after normalization
	phoneRegex := `[^0-9]`
	phoneCleaned := regexp.MustCompile(phoneRegex).ReplaceAllString(phone, "")
	if len(phoneCleaned) != 10 {
		return errors.New("phone number must be exactly 10 digits")
	}
	return nil
}

func NormalizePhone(phone string) string {
	// Remove non-numeric characters from phone
	phoneRegex := `[^0-9]`
	return regexp.MustCompile(phoneRegex).ReplaceAllString(phone, "")
}

func ValidateAndNormalizePhone(phone string) (string, error) {
	// First normalize
	normalizedPhone := NormalizePhone(phone)

	// Then validate
	err := ValidatePhone(normalizedPhone)
	return normalizedPhone, err
}

func ValidateUsername(username string) error {
	// Validate username format
	usernameRegex := `^[a-zA-Z0-9._]+$`
	matched, err := regexp.MatchString(usernameRegex, username)
	if err != nil {
		return errors.New("failed to validate username format")
	}
	if !matched {
		return errors.New("username can only contain alphabets, numbers, dots, and underscores")
	}

	// Check username length
	if len(username) < 3 || len(username) > 20 {
		return errors.New("username must be 3-20 characters")
	}

	// Check for invalid start/end characters (to prevent misinterpret as hiden directory (/home/.john) or command-line flag (ssh -help))
	if strings.HasPrefix(username, ".") || strings.HasPrefix(username, "_") ||
		strings.HasSuffix(username, ".") || strings.HasSuffix(username, "_") {
		return errors.New("username cannot start or end with '.' or '_'")
	}

	// Check for invalid sequences (prevent directory traversal attacks / ambiguity and parsing issues)
	if strings.Contains(username, "..") || strings.Contains(username, "__") ||
		strings.Contains(username, "._") || strings.Contains(username, "_.") {
		return errors.New("username cannot contain '..', '__', '._', or '_.'")
	}

	return nil
}

func NormalizeUsername(username string) string {
	// Trim spaces from username
	return strings.TrimSpace(username)
}

func ValidateAndNormalizeUsername(username string) (string, error) {
	// First normalize
	normalizedUsername := NormalizeUsername(username)

	// Then validate
	err := ValidateUsername(normalizedUsername)
	return normalizedUsername, err
}

func ValidatePassword(password string) error {
	// Check minimum length
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	// Check maximum length (prevent DoS attacks)
	if len(password) > 128 {
		return errors.New("password must not exceed 128 characters")
	}

	// Check for at least one lowercase letter
	hasLower, _ := regexp.MatchString(`[a-z]`, password)
	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}

	// Check for at least one uppercase letter
	hasUpper, _ := regexp.MatchString(`[A-Z]`, password)
	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}

	// Check for at least one digit
	hasDigit, _ := regexp.MatchString(`[0-9]`, password)
	if !hasDigit {
		return errors.New("password must contain at least one digit")
	}

	// Check for at least one special character
	hasSpecial, _ := regexp.MatchString(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`, password)
	if !hasSpecial {
		return errors.New("password must contain at least one special character")
	}

	return nil
}

func NormalizePassword(password string) string {
	// Trim spaces from password
	return strings.TrimSpace(password)
}

func ValidateNames(firstName, lastName string) error {
	// Validate first name
	if err := ValidateFirstName(firstName); err != nil {
		return err
	}

	// Validate last name
	if err := ValidateLastName(lastName); err != nil {
		return err
	}

	return nil
}

func NormalizeNames(firstName, lastName string) (string, string) {
	// Normalize both names
	normalizedFirstName := NormalizeName(firstName)
	normalizedLastName := NormalizeName(lastName)
	return normalizedFirstName, normalizedLastName
}

func ValidateAndNormalizeNames(firstName, lastName string) (string, string, error) {
	// First normalize
	normalizedFirstName, normalizedLastName := NormalizeNames(firstName, lastName)

	// Then validate
	err := ValidateNames(normalizedFirstName, normalizedLastName)
	return normalizedFirstName, normalizedLastName, err
}

func ValidateFirstName(firstName string) error {
	// Check length
	if len(firstName) < 1 {
		return errors.New("first name is required")
	}
	if len(firstName) > 200 {
		return errors.New("first name must not exceed 200 characters")
	}

	// Check for valid characters (letters, spaces, hyphens, apostrophes)
	nameRegex := `^[a-zA-ZÀ-ÿ\s\-']+$`
	matched, err := regexp.MatchString(nameRegex, firstName)
	if err != nil {
		return errors.New("failed to validate first name format")
	}
	if !matched {
		return errors.New("first name can only contain letters, spaces, hyphens, and apostrophes")
	}

	// Check for consecutive special characters
	if strings.Contains(firstName, "--") || strings.Contains(firstName, "''") || strings.Contains(firstName, "  ") {
		return errors.New("first name cannot contain consecutive special characters or spaces")
	}

	// Check start/end with valid characters (not space, hyphen, or apostrophe)
	firstChar := firstName[0]
	lastChar := firstName[len(firstName)-1]
	if firstChar == ' ' || firstChar == '-' || firstChar == '\'' ||
		lastChar == ' ' || lastChar == '-' || lastChar == '\'' {
		return errors.New("first name cannot start or end with space, hyphen, or apostrophe")
	}

	return nil
}

func ValidateLastName(lastName string) error {
	// Check length
	if len(lastName) < 1 {
		return errors.New("last name is required")
	}
	if len(lastName) > 200 {
		return errors.New("last name must not exceed 200 characters")
	}

	// Check for valid characters (letters, spaces, hyphens, apostrophes)
	nameRegex := `^[a-zA-ZÀ-ÿ\s\-']+$`
	matched, err := regexp.MatchString(nameRegex, lastName)
	if err != nil {
		return errors.New("failed to validate last name format")
	}
	if !matched {
		return errors.New("last name can only contain letters, spaces, hyphens, and apostrophes")
	}

	// Check for consecutive special characters
	if strings.Contains(lastName, "--") || strings.Contains(lastName, "''") || strings.Contains(lastName, "  ") {
		return errors.New("last name cannot contain consecutive special characters or spaces")
	}

	// Check start/end with valid characters (not space, hyphen, or apostrophe)
	firstChar := lastName[0]
	lastChar := lastName[len(lastName)-1]
	if firstChar == ' ' || firstChar == '-' || firstChar == '\'' ||
		lastChar == ' ' || lastChar == '-' || lastChar == '\'' {
		return errors.New("last name cannot start or end with space, hyphen, or apostrophe")
	}

	return nil
}

func NormalizeName(name string) string {
	// Trim spaces
	name = strings.TrimSpace(name)

	if len(name) == 0 {
		return name
	}

	// Convert to proper case (capitalize first letter and after separators)
	runes := []rune(strings.ToLower(name))
	capitalize := true

	for i, r := range runes {
		if capitalize && r >= 'a' && r <= 'z' {
			runes[i] = r - 'a' + 'A'
			capitalize = false
		} else if r == ' ' || r == '-' || r == '\'' {
			capitalize = true
		} else {
			capitalize = false
		}
	}

	return string(runes)
}
