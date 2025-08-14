package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	validate *validator.Validate
	trans    ut.Translator
)

// InitValidator initializes the validator with English translations
func InitValidator() {
	validate = validator.New()

	// Create English translator
	english := en.New()
	uni := ut.New(english, english)
	trans, _ = uni.GetTranslator("en")

	// Register English translations
	en_translations.RegisterDefaultTranslations(validate, trans)

	// Register custom validations
	registerCustomValidations()
}

// Validate validates a struct and returns validation errors
func Validate(s interface{}) []ValidationError {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	var errors []ValidationError
	for _, err := range err.(validator.ValidationErrors) {
		errors = append(errors, ValidationError{
			Field:   err.Field(),
			Tag:     err.Tag(),
			Value:   err.Value(),
			Message: err.Translate(trans),
		})
	}

	return errors
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string      `json:"field"`
	Tag     string      `json:"tag"`
	Value   interface{} `json:"value"`
	Message string      `json:"message"`
}

// Error returns the error message
func (e ValidationError) Error() string {
	return fmt.Sprintf("Field '%s' failed validation '%s': %s", e.Field, e.Tag, e.Message)
}

// registerCustomValidations registers custom validation functions
func registerCustomValidations() {
	// Custom validation for currency codes
	validate.RegisterValidation("currency", validateCurrency)

	// Custom validation for amount (positive decimal)
	validate.RegisterValidation("amount", validateAmount)

	// Custom validation for username (alphanumeric + underscore)
	validate.RegisterValidation("username", validateUsername)
}

// validateCurrency validates currency codes (3 letters)
func validateCurrency(fl validator.FieldLevel) bool {
	currency := fl.Field().String()
	return len(currency) == 3 && strings.ToUpper(currency) == currency
}

// validateAmount validates positive decimal amounts
func validateAmount(fl validator.FieldLevel) bool {
	amount := fl.Field().Float()
	return amount > 0
}

// validateUsername validates username format
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	// Username should be 3-50 characters, alphanumeric + underscore
	return len(username) >= 3 && len(username) <= 50
}

// GetValidator returns the global validator instance
func GetValidator() *validator.Validate {
	return validate
}

// GetTranslator returns the global translator instance
func GetTranslator() ut.Translator {
	return trans
}
