package utils

import (
	"errors"
	"reflect"
	"regexp"
	"strings"
)

func ValidateStruct(s interface{}) error {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return errors.New("input is not a struct")
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		tag := field.Tag.Get("validate")

		if tag == "" {
			continue
		}

		rules := strings.Split(tag, ",")
		for _, rule := range rules {
			if err := validateField(field.Name, value, rule); err != nil {
				return err
			}
		}
	}

	return nil
}

func validateField(fieldName string, value reflect.Value, rule string) error {
	switch {
	case rule == "required":
		if value.Kind() == reflect.String && value.String() == "" {
			return errors.New(fieldName + " is required")
		}
	case rule == "email":
		if value.Kind() == reflect.String {
			email := value.String()
			emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
			matched, _ := regexp.MatchString(emailRegex, email)
			if !matched {
				return errors.New(fieldName + " must be a valid email")
			}
		}
	case strings.HasPrefix(rule, "min="):
		minStr := strings.TrimPrefix(rule, "min=")
		if value.Kind() == reflect.String {
			if len(value.String()) < parseInt(minStr) {
				return errors.New(fieldName + " must be at least " + minStr + " characters")
			}
		}
	}
	return nil
}

func parseInt(s string) int {
	result := 0
	for _, r := range s {
		if r >= '0' && r <= '9' {
			result = result*10 + int(r-'0')
		}
	}
	return result
}
