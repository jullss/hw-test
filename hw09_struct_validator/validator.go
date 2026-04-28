package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

var (
	ErrMinLimit = errors.New("value is less than min threshold")
	ErrMaxLimit = errors.New("value is more than max threshold")
	ErrInSet    = errors.New("value must be in allowed set")
	ErrLen      = errors.New("value is not with allowed length")
	ErrRegexp   = errors.New("does not match regexp")
)

var (
	ErrNotStruct  = errors.New("should be a struct")
	ErrInvalidTag = errors.New("invalid tag syntax")
)

func (v ValidationErrors) Error() string {
	res := make([]string, 0, len(v))

	for _, err := range v {
		res = append(res, err.Err.Error())
	}

	return strings.Join(res, "; ")
}

func (v ValidationErrors) Unwrap() []error {
	res := make([]error, len(v))

	for i, v := range v {
		res[i] = v.Err
	}

	return res
}

func Validate(v interface{}) error {
	t := reflect.TypeOf(v)
	val := reflect.ValueOf(v)
	ve := make(ValidationErrors, 0)

	if t.Kind() == reflect.Ptr {
		if val.IsNil() {
			return fmt.Errorf("pass nil pointer")
		}

		t = t.Elem()
		val = val.Elem()
	}

	if t.Kind() != reflect.Struct {
		return fmt.Errorf("%s '%w'", t.Kind(), ErrNotStruct)
	}

	for i := 0; i < t.NumField(); i++ {
		fieldType := t.Field(i)
		fieldVal := val.Field(i)

		tagVal, ok := fieldType.Tag.Lookup("validate")
		if !ok {
			continue
		}

		rules := strings.Split(tagVal, "|")

		for _, rule := range rules {
			validErr, sysErr := processField(fieldVal, fieldType.Name, rule)

			if sysErr != nil {
				return sysErr
			}

			if len(validErr) > 0 {
				ve = append(ve, validErr...)
			}
		}
	}

	if len(ve) == 0 {
		return nil
	}

	return ve
}

func processField(fieldVal reflect.Value, fieldName string, rule string) ([]ValidationError, error) {
	var resErrors []ValidationError
	kind := fieldVal.Kind()

	//nolint:exhaustive
	switch kind {
	case reflect.Int:
		err := validateInt(fieldVal, rule)
		if err != nil {
			if errors.Is(err, ErrInvalidTag) {
				return nil, err
			}
			resErrors = append(resErrors, ValidationError{Field: fieldName, Err: err})
		}
	case reflect.String:
		err := validateString(fieldVal, rule)
		if err != nil {
			if errors.Is(err, ErrInvalidTag) {
				return nil, err
			}
			resErrors = append(resErrors, ValidationError{Field: fieldName, Err: err})
		}
	case reflect.Slice:
		validErrs, sysErr := validateSlice(fieldVal, fieldName, rule)
		if sysErr != nil {
			return nil, sysErr
		}

		resErrors = append(resErrors, validErrs...)
	}

	return resErrors, nil
}

func validateSlice(fieldVal reflect.Value, fieldName string, rule string) ([]ValidationError, error) {
	var ve []ValidationError
	for i := 0; i < fieldVal.Len(); i++ {
		elem := fieldVal.Index(i)
		var err error

		if elem.Kind() == reflect.Int {
			err = validateInt(elem, rule)
		} else if elem.Kind() == reflect.String {
			err = validateString(elem, rule)
		}

		if err != nil {
			if errors.Is(err, ErrInvalidTag) {
				return ve, err
			}

			ve = append(ve, ValidationError{
				Field: fmt.Sprintf("%s[%d]", fieldName, i),
				Err:   err,
			})
		}
	}
	return ve, nil
}

func validateInt(v reflect.Value, tagVal string) error {
	value := v.Int()
	ruleName, ruleVal, found := strings.Cut(tagVal, ":")

	if !found {
		return nil
	}

	switch ruleName {
	case "min":
		minV, err := strconv.Atoi(ruleVal)
		if err != nil {
			return fmt.Errorf("%w: '%s'", ErrInvalidTag, ruleVal)
		}
		if value < int64(minV) {
			return ErrMinLimit
		}
	case "max":
		maxV, err := strconv.Atoi(ruleVal)
		if err != nil {
			return fmt.Errorf("%w: '%s'", ErrInvalidTag, ruleVal)
		}
		if value > int64(maxV) {
			return ErrMaxLimit
		}
	case "in":
		match := false
		allowed := strings.Split(ruleVal, ",")

		for _, n := range allowed {
			n = strings.TrimSpace(n)
			allowedNum, err := strconv.Atoi(n)
			if err != nil {
				return fmt.Errorf("%w: '%s'", ErrInvalidTag, ruleVal)
			}

			if value == int64(allowedNum) {
				match = true
				break
			}
		}

		if !match {
			return ErrInSet
		}
	}

	return nil
}

func validateString(v reflect.Value, tagVal string) error {
	value := strings.TrimSpace(v.String())
	ruleName, ruleVal, found := strings.Cut(tagVal, ":")

	if !found {
		return nil
	}

	switch ruleName {
	case "len":
		allowedLen, err := strconv.Atoi(ruleVal)
		if err != nil {
			return fmt.Errorf("%w: '%s'", ErrInvalidTag, ruleVal)
		}
		if allowedLen != utf8.RuneCountInString(value) {
			return ErrLen
		}
	case "in":
		match := false
		allowedVals := strings.Split(ruleVal, ",")

		for _, w := range allowedVals {
			if value == strings.TrimSpace(w) {
				match = true
				break
			}
		}

		if !match {
			return ErrInSet
		}
	case "regexp":
		re, err := regexp.Compile(ruleVal)
		if err != nil {
			return fmt.Errorf("%w: '%s'", ErrInvalidTag, ruleVal)
		}
		if !re.MatchString(value) {
			return ErrRegexp
		}
	}

	return nil
}
