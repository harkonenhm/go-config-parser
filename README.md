# Config Parser Package

This Go package provides a simple way to parse configuration files and bind their values to a struct.

## Features

- Supports parsing configuration files in key-value format
- Automatically binds parsed values to a struct using reflection
- Handles different data types (string, int, float, bool)

## Example Usage

The configuration file is a simple key-value text file with the following format:

```
port: 8080
address: localhost
```

In this example, the `config.txt` file contains two lines: one for the port number and one for the server address.

```go
package main

import (
	"config"
	"fmt"
)

type Config struct {
	Port    int   `config:"port"`
	Address string `config:"address"`
}

func main() {
	var config Config
	err := config.ParseConfigFromExeFolder(&config, "config.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Port: %d\n", config.Port) // prints 8080
	fmt.Printf("Address: %s\n", config.Address) // prints localhost
}
```

