package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func ParseConfigFromExeFolder(configStruct interface{}, configFile string) error {
	executable, err := os.Executable()
	if err != nil {
		return err
	}

	configPath := filepath.Dir(executable) + "/" + configFile
	return ParseConfigFromFile(configStruct, configPath)
}

func ParseConfigFromFile(configStruct interface{}, configPath string) error {
	configStr, err := readConfigFile(configPath)
	if err != nil {
		return err
	}

	return ParseConfig(configStruct, configStr)
}

func ParseConfig(configStruct interface{}, configStr string) error {
	keyValueConfig, err := extractValueMap(configStr)
	if err != nil {
		return err
	}

	for key, value := range keyValueConfig {
		err := setField(configStruct, key, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func readConfigFile(configPath string) (string, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func setField(s interface{}, tagValue, value string) error {
	// Get the value of the struct (must be a pointer to the struct)
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("s must be a pointer to a struct")
	}

	// Get the underlying struct value
	structValue := v.Elem()
	structType := structValue.Type()

	for i := 0; i < structType.NumField(); i++ {
		field := structValue.Field(i)
		fieldType := structType.Field(i)

		tagVal := fieldType.Tag.Get("config")
		if tagVal == tagValue {

			// Ensure the field is settable
			if !field.CanSet() {
				return fmt.Errorf("cannot set field %s", fieldType.Name)
			}

			// Convert the string value to the field's type and set it
			switch field.Kind() {
			case reflect.String:
				field.SetString(value)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				intVal, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return fmt.Errorf("failed to convert value to int: %v", err)
				}
				field.SetInt(intVal)
			case reflect.Float32, reflect.Float64:
				floatVal, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return fmt.Errorf("failed to convert value to float: %v", err)
				}
				field.SetFloat(floatVal)
			case reflect.Bool:
				boolVal, err := strconv.ParseBool(value)
				if err != nil {
					return fmt.Errorf("failed to convert value to bool: %v", err)
				}
				field.SetBool(boolVal)
			default:
				return fmt.Errorf("unsupported field type: %s", field.Kind())
			}
		}
	}

	return nil
}

func extractValueMap(str string) (map[string]string, error) {
	if len(str) == 0 {
		return nil, fmt.Errorf("empty config string")
	}
	//  remove whitespace characters but retain \n
	re := regexp.MustCompile("[ \t\r\f\v]+")
	str = re.ReplaceAllString(str, "")

	strSlice := strings.Split(str, "\n")

	m := make(map[string]string)
	for _, s := range strSlice {
		// ignore empty rows
		if len(s) == 0 {
			continue
		}
		if strings.Count(s, ":") != 1 {
			return nil, fmt.Errorf("invalid config key/value")
		}
		keyValue := strings.Split(s, ":")
		m[keyValue[0]] = keyValue[1]
	}
	return m, nil
}
