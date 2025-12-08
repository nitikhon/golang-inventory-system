package validation

import (
	"errors"
	"regexp"
	"strings"

	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
	"github.com/nitikhon/golang-inventory-system/internal/util/errormap"
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

// ValidateUserForUpdate validates user fields for update (excludes username)
func ValidateUserForUpdate(user *entity.User) error {
	// Validate required fields for update
	if err := ValidateRequiredFieldsForUpdate(user); err != nil {
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

// NormalizeUserForUpdate normalizes user fields for update (excludes username)
func NormalizeUserForUpdate(user *entity.User) {
	// Normalize email
	user.Email = NormalizeEmail(user.Email)

	// Normalize phone
	user.Phone = NormalizePhone(user.Phone)

	// Normalize names
	user.FirstName, user.LastName = NormalizeNames(user.FirstName, user.LastName)

	// Normalize password
	user.Password = NormalizePassword(user.Password)
}

// ValidateAndNormalizeUserForUpdate validates and normalizes user data for update (excludes username)
func ValidateAndNormalizeUserForUpdate(user *entity.User) error {
	// First normalize the data
	NormalizeUserForUpdate(user)

	// Then validate the normalized data
	if err := ValidateUserForUpdate(user); err != nil {
		return err
	}

	return nil
}

// ValidateRequiredFieldsForUpdate validates required fields for update (excludes username)
func ValidateRequiredFieldsForUpdate(user *entity.User) error {
	if user.Email == "" {
		return errors.New(errormap.ErrEmailRequired)
	}
	if user.Password == "" {
		return errors.New(errormap.ErrPasswordRequired)
	}
	if user.Phone == "" {
		return errors.New(errormap.ErrPhoneRequired)
	}
	if user.FirstName == "" {
		return errors.New(errormap.ErrFirstNameRequired)
	}
	if user.LastName == "" {
		return errors.New(errormap.ErrLastNameRequired)
	}
	return nil
}

func ValidateRequiredFields(user *entity.User) error {
	// Check if all required fields are provided
	if user.Username == "" {
		return errors.New(errormap.ErrUsernameRequired)
	}
	if user.Email == "" {
		return errors.New(errormap.ErrEmailRequired)
	}
	if user.Password == "" {
		return errors.New(errormap.ErrPasswordRequired)
	}
	if user.Phone == "" {
		return errors.New(errormap.ErrPhoneRequired)
	}
	if user.FirstName == "" {
		return errors.New(errormap.ErrFirstNameRequired)
	}
	if user.LastName == "" {
		return errors.New(errormap.ErrLastNameRequired)
	}
	return nil
}

func ValidateEmail(email string) error {
	if len(email) < 6 {
		return errors.New(errormap.ErrEmailMinLength)
	}

	if len(email) > 128 {
		return errors.New(errormap.ErrEmailMaxLength)
	}

	// Validate email format
	emailRegex := `^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`
	matched, err := regexp.MatchString(emailRegex, email)
	if err != nil {
		return errors.New(errormap.ErrEmailValidationFailed)
	}
	if !matched {
		return errors.New(errormap.ErrInvalidEmailFormat)
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
		return errors.New(errormap.ErrPhoneDigits)
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
		return errors.New(errormap.ErrUsernameValidationFailed)
	}
	if !matched {
		return errors.New(errormap.ErrUsernameInvalidChars)
	}

	// Check username length
	if len(username) < 3 || len(username) > 20 {
		return errors.New(errormap.ErrUsernameLength)
	}

	// Check for invalid start/end characters (to prevent misinterpret as hiden directory (/home/.john) or command-line flag (ssh -help))
	if strings.HasPrefix(username, ".") || strings.HasPrefix(username, "_") ||
		strings.HasSuffix(username, ".") || strings.HasSuffix(username, "_") {
		return errors.New(errormap.ErrUsernameInvalidStartEnd)
	}

	// Check for invalid sequences (prevent directory traversal attacks / ambiguity and parsing issues)
	if strings.Contains(username, "..") || strings.Contains(username, "__") ||
		strings.Contains(username, "._") || strings.Contains(username, "_.") {
		return errors.New(errormap.ErrUsernameInvalidSequence)
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
		return errors.New(errormap.ErrPasswordMinLength)
	}

	// Check maximum length (prevent DoS attacks)
	if len(password) > 128 {
		return errors.New(errormap.ErrPasswordMaxLength)
	}

	// Check for at least one lowercase letter
	hasLower, _ := regexp.MatchString(`[a-z]`, password)
	if !hasLower {
		return errors.New(errormap.ErrPasswordLowercase)
	}

	// Check for at least one uppercase letter
	hasUpper, _ := regexp.MatchString(`[A-Z]`, password)
	if !hasUpper {
		return errors.New(errormap.ErrPasswordUppercase)
	}

	// Check for at least one digit
	hasDigit, _ := regexp.MatchString(`[0-9]`, password)
	if !hasDigit {
		return errors.New(errormap.ErrPasswordDigit)
	}

	// Check for at least one special character
	hasSpecial, _ := regexp.MatchString(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`, password)
	if !hasSpecial {
		return errors.New(errormap.ErrPasswordSpecialChar)
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
		return errors.New(errormap.ErrFirstNameRequired)
	}
	if len(firstName) > 200 {
		return errors.New(errormap.ErrFirstNameMaxLength)
	}

	// Check for valid characters (letters, spaces, hyphens, apostrophes)
	nameRegex := `^[a-zA-ZÀ-ÿ\s\-']+$`
	matched, err := regexp.MatchString(nameRegex, firstName)
	if err != nil {
		return errors.New(errormap.ErrFirstNameValidationFailed)
	}
	if !matched {
		return errors.New(errormap.ErrFirstNameInvalidChars)
	}

	// Check for consecutive special characters
	if strings.Contains(firstName, "--") || strings.Contains(firstName, "''") || strings.Contains(firstName, "  ") {
		return errors.New(errormap.ErrFirstNameConsecutiveSpecial)
	}

	// Check start/end with valid characters (not space, hyphen, or apostrophe)
	firstChar := firstName[0]
	lastChar := firstName[len(firstName)-1]
	if firstChar == ' ' || firstChar == '-' || firstChar == '\'' ||
		lastChar == ' ' || lastChar == '-' || lastChar == '\'' {
		return errors.New(errormap.ErrFirstNameInvalidStartEnd)
	}

	return nil
}

func ValidateLastName(lastName string) error {
	// Check length
	if len(lastName) < 1 {
		return errors.New(errormap.ErrLastNameRequired)
	}
	if len(lastName) > 200 {
		return errors.New(errormap.ErrLastNameMaxLength)
	}

	// Check for valid characters (letters, spaces, hyphens, apostrophes)
	nameRegex := `^[a-zA-ZÀ-ÿ\s\-']+$`
	matched, err := regexp.MatchString(nameRegex, lastName)
	if err != nil {
		return errors.New(errormap.ErrLastNameValidationFailed)
	}
	if !matched {
		return errors.New(errormap.ErrLastNameInvalidChars)
	}

	// Check for consecutive special characters
	if strings.Contains(lastName, "--") || strings.Contains(lastName, "''") || strings.Contains(lastName, "  ") {
		return errors.New(errormap.ErrLastNameConsecutiveSpecial)
	}

	// Check start/end with valid characters (not space, hyphen, or apostrophe)
	firstChar := lastName[0]
	lastChar := lastName[len(lastName)-1]
	if firstChar == ' ' || firstChar == '-' || firstChar == '\'' ||
		lastChar == ' ' || lastChar == '-' || lastChar == '\'' {
		return errors.New(errormap.ErrLastNameInvalidStartEnd)
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
