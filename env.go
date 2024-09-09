package env

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/igwtcode/go-env/internal/topt"
)

const (
	DefaultTagOptionSeparator  = "," // Default separator for tag options
	DefaultSliceValueSeparator = "|" // Default separator for slice values
)

// Parser represents a configurable environment variable parser.
type Parser struct {
	TagOptionSeparator  string // Separator for options in the tag (e.g., ',')
	SliceValueSeparator string // Separator for values in slices (e.g., '|')
	NamePrefix          string // Name prefix for environment variables
}

// NewParser creates a new Parser with default configuration.
func NewParser() *Parser {
	return &Parser{
		TagOptionSeparator:  DefaultTagOptionSeparator,
		SliceValueSeparator: DefaultSliceValueSeparator,
	}
}

// WithTagOptionSeparator configures the separator for tag options (default: ',').
func (p *Parser) WithTagOptionSeparator(separator string) *Parser {
	if separator == p.SliceValueSeparator {
		panic("tag option separator and slice value separator must not be the same")
	}
	p.TagOptionSeparator = separator
	return p
}

// WithSliceValueSeparator configures the separator for slice values (default: '|').
func (p *Parser) WithSliceValueSeparator(separator string) *Parser {
	if separator == p.TagOptionSeparator {
		panic("slice value separator and tag option separator must not be the same")
	}
	p.SliceValueSeparator = separator
	return p
}

// WithNamePrefix configures the prefix to add to environment variable names.
func (p *Parser) WithNamePrefix(prefix string) *Parser {
	p.NamePrefix = prefix
	return p
}

// parseTag parses the tag string into a map of options (e.g., "required", "default=foo").
func (p *Parser) parseTag(tag string) map[string]string {
	options := map[string]string{}
	parts := strings.Split(tag, p.TagOptionSeparator)
	for _, part := range parts {
		kv := strings.SplitN(part, "=", 2)
		key := strings.TrimSpace(strings.ToLower(kv[0]))
		if len(kv) == 2 {
			options[key] = kv[1]
		} else {
			options[key] = ""
		}
	}
	return options
}

// Unmarshal reads environment variables and populates the struct fields.
func (p *Parser) Unmarshal(envStruct interface{}) error {
	v := reflect.ValueOf(envStruct).Elem()
	t := reflect.TypeOf(envStruct).Elem()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		// Skip unexported fields
		if !fieldValue.CanSet() {
			continue
		}

		// Recursively handle embedded structs
		if fieldValue.Kind() == reflect.Struct {
			if err := p.Unmarshal(fieldValue.Addr().Interface()); err != nil {
				return err
			}
			continue
		}

		// Parse the `env` tagVal for environment variable options
		tagVal, tagOk := field.Tag.Lookup("env")
		if !tagOk {
			continue
		}
		tagOptions := p.parseTag(tagVal)

		// Get the lookup order for environment variables, ensuring unique names
		envNames := getEnvNames(field.Name, tagOptions, p)
		envVal := getEnvValue(envNames)

		// Apply trim by default, can be disabled with 'notrim' option
		if _, notrim := tagOptions[topt.NOTRIM]; !notrim {
			envVal = strings.TrimSpace(envVal)
		}

		// Handle default value
		if envVal == "" && tagOptions[topt.DEFAULT] != "" {
			envVal = tagOptions[topt.DEFAULT]
		}

		// Handle required fields
		if _, required := tagOptions[topt.REQUIRED]; required && envVal == "" {
			return fmt.Errorf("environment variable %s is required but not set", strings.Join(envNames, p.SliceValueSeparator))
		}

		// Handle lowercase
		if _, lower := tagOptions[topt.LOWER]; lower {
			envVal = strings.ToLower(envVal)
		}

		// Handle uppercase
		if _, upper := tagOptions[topt.UPPER]; upper {
			envVal = strings.ToUpper(envVal)
		}

		// Process slices using the configured slice value separator
		if fieldValue.Kind() == reflect.Slice {
			if err := handleSliceWithSeparator(fieldValue, envVal, tagOptions, p.SliceValueSeparator); err != nil {
				return err
			}
			continue
		}

		// Check if the field has an AWS-specific validation option and apply the validation
		if err := checkForAwsValidation(field.Name, envVal, tagOptions); err != nil {
			return err
		}

		// Set value to the appropriate field
		if err := setValue(fieldValue, envVal, tagOptions); err != nil {
			return err
		}
	}

	return nil
}

// awsValidationMap finds and applies the validation function for AWS-specific environment variables tag options.
func checkForAwsValidation(fieldName string, envVal string, tagOptions map[string]string) error {
	// if the field is not required and the env value is empty, return
	if _, ok := tagOptions[topt.REQUIRED]; !ok && envVal == "" {
		return nil
	}

	// Count how many v_aws_... validation options are provided
	vc := 0
	var vfn func(string) error

	// Check for v_aws_xxx options and validate exclusivity
	for tag, fn := range awsValidationMap {
		if _, ok := tagOptions[tag]; ok {
			vc++
			if vc > 1 {
				return fmt.Errorf("multiple v_aws validation options provided for field '%s': only one is allowed", fieldName)
			}
			vfn = fn
		}
	}

	// Apply the validation if v_aws validation option is found
	if vfn != nil {
		return vfn(envVal)
	}
	return nil
}

// getEnvNames returns a list of environment variable names to check, based on the 'name' tag option or the field name.
func getEnvNames(fieldName string, tagOptions map[string]string, p *Parser) []string {
	var envNames []string

	ap := func(sl []string) {
		for _, s := range sl {
			v := p.NamePrefix + s
			if !slices.Contains(envNames, v) {
				envNames = append(envNames, v)
			}
		}
	}

	// Check if `name` tag is provided, and split it into multiple names using the slice value separator.
	if name, ok := tagOptions[topt.NAME]; ok && name != "" {
		ap(strings.Split(name, p.SliceValueSeparator))
	}

	// Add the field name and the field name in upper and lower case
	ap([]string{fieldName, strings.ToUpper(fieldName), strings.ToLower(fieldName)})

	return envNames
}

// getEnvValue checks environment variables in order and returns the first non-empty value found.
func getEnvValue(envNames []string) string {
	for _, name := range envNames {
		if val := os.Getenv(name); val != "" {
			return val
		}
	}
	return ""
}

// setValue sets the value for a struct field based on its type.
func setValue(field reflect.Value, val string, tagOptions map[string]string) error {
	return setReflectValue(field, val, field.Kind(), tagOptions)
}

// setSliceValue sets the appropriate value for a slice element.
func setSliceValue(sliceElement reflect.Value, val string, kind reflect.Kind, tagOptions map[string]string) error {
	return setReflectValue(sliceElement, val, kind, tagOptions)
}

// setReflectValue sets the appropriate value based on the field's type.
func setReflectValue(field reflect.Value, val string, kind reflect.Kind, tagOptions map[string]string) error {
	switch kind {
	case reflect.String:
		field.SetString(val)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}
		if err := checkMinMax(intVal, tagOptions); err != nil {
			return err
		}
		field.SetInt(intVal)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return err
		}
		if err := checkMinMax(uintVal, tagOptions); err != nil {
			return err
		}
		field.SetUint(uintVal)
	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		}
		if err := checkMinMax(floatVal, tagOptions); err != nil {
			return err
		}
		field.SetFloat(floatVal)
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		field.SetBool(boolVal)
	default:
		return errors.New("unsupported field type")
	}
	return nil
}

// handleSliceWithSeparator processes slice types, splitting the input string using a specified separator.
func handleSliceWithSeparator(field reflect.Value, envVal string, tagOptions map[string]string, separator string) error {
	sliceType := field.Type().Elem().Kind()

	if envVal == "" {
		field.Set(reflect.MakeSlice(field.Type(), 0, 0))
		return nil
	}
	_, notrim := tagOptions[topt.NOTRIM]

	// Split the environment variable by the separator
	values := strings.Split(envVal, separator)
	// Filter out any empty elements (after trimming)
	filteredValues := []string{}
	for _, val := range values {
		if notrim {
			filteredValues = append(filteredValues, val)
		} else {
			trimmedVal := strings.TrimSpace(val)
			if trimmedVal != "" {
				filteredValues = append(filteredValues, trimmedVal)
			}
		}
	}

	// If all values are empty, set an empty slice
	if len(filteredValues) == 0 {
		field.Set(reflect.MakeSlice(field.Type(), 0, 0))
		return nil
	}

	// Create a new slice with the filtered values
	newSlice := reflect.MakeSlice(field.Type(), len(filteredValues), len(filteredValues))

	for i, val := range filteredValues {
		err := setSliceValue(newSlice.Index(i), val, sliceType, tagOptions)
		if err != nil {
			return err
		}
	}

	field.Set(newSlice)
	return nil
}

// checkMinMax validates if the value is within the range specified by the "min" and "max" tags.
func checkMinMax(val interface{}, tagOptions map[string]string) error {
	if minStr, ok := tagOptions[topt.MIN]; ok {
		min, err := strconv.ParseFloat(minStr, 64)
		if err != nil {
			return fmt.Errorf("invalid min value: %s", minStr)
		}
		if compareNumeric(val, min) < 0 {
			return fmt.Errorf("value %v is less than minimum allowed %v", val, min)
		}
	}

	if maxStr, ok := tagOptions[topt.MAX]; ok {
		max, err := strconv.ParseFloat(maxStr, 64)
		if err != nil {
			return fmt.Errorf("invalid max value: %s", maxStr)
		}
		if compareNumeric(val, max) > 0 {
			return fmt.Errorf("value %v is greater than maximum allowed %v", val, max)
		}
	}
	return nil
}

// compareNumeric compares two numeric values and returns -1 if val < threshold, 0 if equal, and 1 if val > threshold.
func compareNumeric(val interface{}, threshold float64) int {
	switch v := val.(type) {
	case int64:
		if float64(v) < threshold {
			return -1
		} else if float64(v) > threshold {
			return 1
		}
	case float64:
		if v < threshold {
			return -1
		} else if v > threshold {
			return 1
		}
	}
	return 0
}
