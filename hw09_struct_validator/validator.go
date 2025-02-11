package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var sb strings.Builder
	for _, err := range v {
		sb.WriteString(fmt.Sprintf("Field %s: %v\n", err.Field, err.Err))
	}
	return sb.String()
}

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("input must be a struct or a pointer to a struct")
	}

	typ := val.Type()
	var validationErrors ValidationErrors

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		tag := fieldType.Tag.Get("validate")
		if tag == "" {
			continue
		}

		rules := strings.Split(tag, "|")
		for _, rule := range rules {
			if field.Kind() == reflect.Slice {
				for j := 0; j < field.Len(); j++ {
					if err := validateField(field.Index(j), rule); err != nil {
						validationErrors = append(validationErrors, ValidationError{
							Field: fieldType.Name,
							Err:   err,
						})
					}
				}
			} else {
				if err := validateField(field, rule); err != nil {
					validationErrors = append(validationErrors, ValidationError{
						Field: fieldType.Name,
						Err:   err,
					})
				}
			}
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}

func validateField(field reflect.Value, rule string) error {
	parts := strings.SplitN(rule, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid validation rule format")
	}
	ruleName, ruleValue := parts[0], parts[1]

	switch ruleName {
	case "len":
		return validateLen(field, ruleValue)
	case "regexp":
		return validateRegexp(field, ruleValue)
	case "in":
		return validateIn(field, ruleValue)
	case "min":
		return validateMin(field, ruleValue)
	case "max":
		return validateMax(field, ruleValue)
	default:
		return fmt.Errorf("unknown validation rule")
	}
}

func validateLen(field reflect.Value, value string) error {
	if field.Kind() == reflect.String {
		length, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		if len(field.String()) != length {
			return fmt.Errorf("length must be exactly %d", length)
		}
	}
	return nil
}

func validateRegexp(field reflect.Value, value string) error {
	if field.Kind() == reflect.String {
		matched, err := regexp.MatchString(value, field.String())
		if err != nil {
			return err
		}
		if !matched {
			return fmt.Errorf("must match regexp %s", value)
		}
	}
	return nil
}

func validateIn(field reflect.Value, value string) error {
	values := strings.Split(value, ",")
	switch field.Kind() { //nolint:exhaustive
	case reflect.String:
		for _, v := range values {
			if field.String() == v {
				return nil
			}
		}
	case reflect.Int:
		for _, v := range values {
			if intValue, err := strconv.Atoi(v); err == nil && field.Int() == int64(intValue) {
				return nil
			}
		}
	default:
		fmt.Println("unrecognized type")
	}
	return fmt.Errorf("must be one of %s", value)
}

func validateMin(field reflect.Value, value string) error {
	if field.Kind() == reflect.Int {
		min, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		if field.Int() < int64(min) {
			return fmt.Errorf("must be at least %d", min)
		}
	}
	return nil
}

func validateMax(field reflect.Value, value string) error {
	if field.Kind() == reflect.Int {
		max, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		if field.Int() > int64(max) {
			return fmt.Errorf("must be at most %d", max)
		}
	}
	return nil
}
