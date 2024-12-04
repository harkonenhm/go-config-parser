package config

import (
	"os"
	"reflect"
	"testing"
)

func TestParseConfigFromFile(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}

	configPath := wd + "/test_config.conf"
	result := TestStruct{}
	err = ParseConfigFromFile(&result, configPath)
	if err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}

	exp := TestStruct{
		"/path/to/folder/",
		123,
		1.23,
		true,
	}
	if result != exp {
		t.Errorf("Expected %v but got %v", exp, result)
	}
}

func TestParseConfig(t *testing.T) {
	var tests = []struct {
		name          string
		configStr     string
		expTestStruct TestStruct
		expError      bool
	}{
		{
			"no error",
			"config1: /path/to/folder/\nconfig2:123\nconfig3:1.23\nconfig4:  false",
			TestStruct{
				"/path/to/folder/",
				123,
				1.23,
				false,
			},
			false,
		},
		{
			"error empty config string",
			"",
			TestStruct{},
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := TestStruct{}
			err := ParseConfig(&result, test.configStr)
			if test.expError {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("Expected no error but got %v", err)
				}
				if result != test.expTestStruct {
					t.Errorf("Expected %v but got %v", test.expTestStruct, result)
				}
			}
		})
	}
}

func TestSetField(t *testing.T) {

	var tests = []struct {
		name         string
		tag          string
		tagValue     string
		expTestStruc TestStruct
		expError     bool
	}{
		{
			"string value",
			"config1",
			"/path/to/folder/",
			TestStruct{
				ConfigStr: "/path/to/folder/",
			},
			false,
		},
		{
			"int value",
			"config2",
			"10",
			TestStruct{
				ConfigInt: 10,
			},
			false,
		},
		{
			"float value",
			"config3",
			"10.12",
			TestStruct{
				ConfigFloat: 10.12,
			},
			false,
		},
		{
			"bool value",
			"config4",
			"true",
			TestStruct{
				ConfigBool: true,
			},
			false,
		},
		{
			"failed conversion",
			"config2",
			"not an int",
			TestStruct{
				ConfigInt: 10,
			},
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := TestStruct{}
			err := setField(&result, test.tag, test.tagValue)
			if test.expError {
				if err == nil {
					t.Errorf("Expected error but got nil.")
				}
			} else {
				if err != nil {
					t.Fatalf("Expected no error but got %v", err)
				}
				if result != test.expTestStruc {
					t.Errorf("Expected %v but got %v", test.expTestStruc, result)
				}

			}
		})
	}
}

func TestExtractValueMap(t *testing.T) {
	tests := []struct {
		name        string
		configStr   string
		expKeyValue map[string]string
		expError    bool
	}{
		{
			"valid config",
			`config1:/path/to/folder/
      config2:10`,
			map[string]string{
				"config1": "/path/to/folder/",
				"config2": "10",
			},
			false,
		},
		{
			"valid config: extra whitespace characters",
			"\nconfig1:  /path/to/folder/\t\nconfig2    :10  \n",
			map[string]string{
				"config1": "/path/to/folder/",
				"config2": "10",
			},
			false,
		},
		{
			"valid config: empty lines",
			"\n\nconfig1:  /path/to/folder/\t\n\nconfig2    :10  \n",
			map[string]string{
				"config1": "/path/to/folder/",
				"config2": "10",
			},
			false,
		},
		{
			"no valid values: empty string",
			"",
			map[string]string{},
			true,
		},
		{
			"no valid values: multiple :",
			"config1:/path/to/folder/ config2:10",
			map[string]string{},
			true,
		},
		{
			"no valid values: missing :",
			"some text \n more test \n 10",
			map[string]string{},
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := extractValueMap(test.configStr)
			if test.expError {
				if err == nil {
					t.Errorf("Expected error but got nil.")
				}
			} else {
				if err != nil {
					t.Fatalf("Expected no error but %v", err)
				}
				if !reflect.DeepEqual(result, test.expKeyValue) {
					t.Errorf("Expected %v but got %v", test.expKeyValue, result)
				}
			}
		})
	}
}

func TestReadFile(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}
	result, err := readConfigFile(wd + "/test_config.conf")
	if err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}
	exp := `config2: 123 
config1: /path/to/folder/
config3: 1.23
config4: true
`
	if result != exp {
		t.Errorf("Expect %v got %v", exp, result)
	}
}

type TestStruct struct {
	ConfigStr   string  `config:"config1"`
	ConfigInt   int     `config:"config2"`
	ConfigFloat float64 `config:"config3"`
	ConfigBool  bool    `config:"config4"`
}
