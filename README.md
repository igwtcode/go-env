# `go-env`

[![Build Status](https://github.com/igwtcode/go-env/actions/workflows/build.yaml/badge.svg)](https://github.com/igwtcode/go-env/actions/workflows/build.yaml) [![Go Reference](https://pkg.go.dev/badge/github.com/igwtcode/go-env.svg)](https://pkg.go.dev/github.com/igwtcode/go-env)

A Simple yet Configurable Environment Variable Parser

The `go-env` package provides a simple, powerful, and flexible way to populate Go structs from environment variables. It is designed for developers who want a lightweight, dependency-free solution to read configuration from environment variables.

### Key Features:

- **Pure Go**: No third-party dependencies; only Go's built-in libraries.
- **Configurable Parsing**: Customize tag options, slice separators, and even add a prefix to all environment variable names.
- **Supports Structs**: Handles nested and embedded structs effortlessly.
- **Field Types**: Supports a wide range of Go types, including string, uint, int, float, bool and slices.
- **Error Handling**: Provides clear error messages for missing required fields or invalid values.

## Why Use This Package?

If you need to load configuration for your Go application from environment variables in a structured and maintainable way, `go-env` makes it easy. With support for common use cases like default values, required fields, basic validations and even custom parsing options, `go-env` simplifies the process.

## Installation

```bash
go get github.com/igwtcode/go-env
```

## Usage

### Basic Usage

You can configure the parser and use the `Unmarshal` method to populate your Go struct from environment variables.

```go
type Config struct {
	Hosts   []string `env:"name=hostslist|TARGET_HOSTS,default=localhost|127.0.0.1"`
	Port    int      `env:"required"`
	Timeout int      `env:"name=time_out,default=30"`
	Debug   bool     `env:"name=DEBUG,default=false"`
}

parser := env.NewParser()
var cfg Config
err := parser.Unmarshal(&cfg)
if err != nil {
    log.Fatalf("Failed to load config: %v", err)
}
```

### Configuring the Parser

You can configure the parser with options. By default:

- **Tag Option Separator**: `,`
- **Slice Value Separator**: `|`
- **Environment Variable Name Prefix**: ``

> [!WARNING]
> The tag option separator and slice value separator must not be the same. This will cause a panic.

#### 1. Default Configuration

```go
parser := env.NewParser()
```

#### 2. Custom Slice Separator

```go
parser := env.NewParser().WithSliceValueSeparator("/")
```

#### 3. Custom Tag and Slice Separators

```go
parser := env.NewParser().
    WithTagOptionSeparator(";").
    WithSliceValueSeparator(":")
```

If both separators are set to the same value, it will panic.

#### 4. Adding a Prefix to Environment Variables

```go
parser := env.NewParser().WithNamePrefix("MYAPP_")
```

## Example

```go
package main

import (
    "fmt"
    "log"

    "github.com/igwtcode/go-env"
)

type Config struct {
    Port    int    `env:"name=PORT,default=8080"`
    Timeout int    `env:"name=TIMEOUT,default=30"`
    Hosts   []string `env:"name=HOSTS,default=localhost|127.0.0.1"`
}

func main() {
    parser := env.NewParser()
    var cfg Config
    err := parser.Unmarshal(&cfg)
    if err != nil {
        log.Fatalf("Error loading config: %v", err)
    }
    fmt.Printf("%+v\n", cfg)
}
```

## Tag Options

The `go-env` tags control how the environment variables are mapped to struct fields and validated. These are the supported options:

- `name`: Defines the environment variable(s) to use. Multiple names can be provided, separated by the slice separator. First it will lookup all provided names in this list, if not provided with any, uses the field name, then upper and lower case of it respectively for lookup.
- `default`: Provides a default value if the environment variable is not set.
- `required`: Ensures that the field has a value; returns an error if missing.
- `lower`: Converts the value to lowercase.
- `upper`: Converts the value to uppercase. (if both `lower` and `upper` are used, final value would be **uppercase**)
- `notrim`: Disables trimming of leading/trailing whitespace. (used for the value as a whole and all its list items, in case of a slice)
- `min`/`max`: Enforces numeric ranges for integer and float fields. (`min` and `max` included: `[min, max]`)

### Examples:

```go
type Config struct {
    Port    int     `env:"name=PORT,min=1024,max=65534,default=8080"`
    Timeout int     `env:"name=TIMEOUT,required"`
    Hosts   []string `env:"name=HOSTS,default=localhost|127.0.0.1"`
    MyValue string   `env:"notrim,required,upper"`
}
```

## License

This project is licensed under the [MIT License](./LICENSE).
