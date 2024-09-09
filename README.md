# `go-env`

[![Build Status](https://github.com/igwtcode/go-env/actions/workflows/build.yaml/badge.svg)](https://github.com/igwtcode/go-env/actions/workflows/build.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/igwtcode/go-env)](https://goreportcard.com/report/github.com/igwtcode/go-env)
[![Go Reference](https://pkg.go.dev/badge/github.com/igwtcode/go-env.svg)](https://pkg.go.dev/github.com/igwtcode/go-env)

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

The `go-env` package allows you to control how environment variables are mapped to struct fields using tags. These tags provide powerful options to set defaults, enforce validation, and customize the behavior of how environment variables are parsed.

Here are the available options:

- **`name`**: Specifies the environment variable(s) to use for the field. Multiple names can be provided, separated by the slice separator (default `|`). The order of lookup is:

  1. All names listed in the `name` tag.
  2. The struct field name.
  3. The upper and lower case versions of the struct field name.

  Example: `name=AWS_DEFAULT_REGION|AWS_REGION`

- **`default`**: Defines a default value to use if the environment variable is not set.

  Example: `default=8080`

- **`required`**: Ensures the field must have a value. If no environment variable is set and no default is provided, an error is returned.

  Example: `required`

- **`lower`**: Converts the value to lowercase before setting the field.

  Example: `lower`

- **`upper`**: Converts the value to uppercase before setting the field. If both `lower` and `upper` are used, the final value will be **uppercase**.

  Example: `upper`

- **`notrim`**: Disables the default trimming of leading and trailing whitespace. Applies to both single values and list items in slices.

  Example: `notrim`

- **`min`/`max`**: Defines numeric range validation for integers or floats. If the environment variable value is outside the range, an error is returned.

  Example: `min=10,max=100`

- **`v_aws_region`**: Validates that the value is a valid AWS region name.

  Example: `v_aws_region`

- **`v_aws_account_id`**: Validates that the value is a valid 12-digit AWS account ID.

  Example: `v_aws_account_id`

- **`v_aws_role_arn`**: Validates that the value is a valid AWS Role ARN.

  Example: `v_aws_role_arn`

- **`v_aws_bucket_name`**: Validates that the value is a valid AWS S3 bucket name.

  Example: `v_aws_bucket_name`

> [!NOTE]
> AWS validators have no effect, when the field is not required and the env value is empty.

### [Examples](./_examples/)

Here is a comprehensive [example](./_examples/01/main.go) that demonstrates how to use the `go-env` package with options and features:

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/igwtcode/go-env"
)

type Config struct {
	// Example with no env tag, ignored by the parser
	NoEnvTag string

	// Example with trimming disabled and forced to uppercase
	MyValue string `env:"notrim,required,upper"`

	// AWS-specific validations
	AwsRegion    string `env:"name=AWS_DEFAULT_REGION|AWS_REGION,v_aws_region,required"`
	AccountID    string `env:"name=AWS_ACCOUNT_ID,v_aws_account_id,required"`
	RoleArn      string `env:"name=AWS_ROLE_ARN,v_aws_role_arn"`
	S3BucketName string `env:"name=AWS_BUCKET,v_aws_bucket_name"`

	// Additional fields with more combinations of options
	LogLevel string `env:"name=LOG_LEVEL,lower,default=info"`

	// Slice of strings with a default value (localhost and 127.0.0.1)
	Hosts []string `env:"default=localhost|127.0.0.1"`

	// Retry count with min and max validation, and a default value
	Retry uint `env:"name=RETRY_COUNT,min=0,max=10,default=3"`

	// Basic integer with min and max validation, and a default value
	Port int `env:"name=PORT,min=1024,max=65534,default=8080"`

	// Required value, environment variable must be set
	Timeout int `env:"required"`

	// Private field with no env tag, ignored by the parser
	privateField int
}

func main() {
	envVars := map[string]string{
		"DEMO_GO_ENV_MYVALUE":        "  MyValue  ",
		"DEMO_GO_ENV_AWS_REGION":     "us-west-2",
		"DEMO_GO_ENV_AWS_ACCOUNT_ID": "123456789012",
		"DEMO_GO_ENV_AWS_ROLE_ARN":   "arn:aws:iam::123456789012:role/MyRole",
		"DEMO_GO_ENV_AWS_BUCKET":     "my-s3-bucket",
		"DEMO_GO_ENV_LOG_LEVEL":      "DEBUG",
		"DEMO_GO_ENV_HOSTS":          "server1|server2|server3",
		"DEMO_GO_ENV_PORT":           "9090",
		"DEMO_GO_ENV_TIMEOUT":        "30",
	}

	setEnvVars(envVars)
	defer unsetEnvVars(envVars)

	var cfg Config

	parser := env.NewParser().WithNamePrefix("DEMO_GO_ENV_")
	err := parser.Unmarshal(&cfg)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Printf("Config: %+v\n", cfg)
}

// Helper function to set environment variables dynamically
func setEnvVars(envVars map[string]string) {
	for key, value := range envVars {
		os.Setenv(key, value)
	}
}

// Helper function to unset environment variables
func unsetEnvVars(envVars map[string]string) {
	for key := range envVars {
		os.Unsetenv(key)
	}
}
```

### Example Output:

```bash
Config: {NoEnvTag: MyValue:  MYVALUE   AwsRegion:us-west-2 AccountID:123456789012 RoleArn:arn:aws:iam::123456789012:role/MyRole S3BucketName:my-s3-bucket LogLevel:debug Hosts:[server1 server2 server3] Retry:3 Port:9090 Timeout:30 privateField:0}
```

## Related Projects

- [Netflix/go-env](https://github.com/Netflix/go-env): A similar Go package that reads environment variables into structs. While it offers very good functionality, `go-env` provides additional flexibility with configurable tag options, custom separators, and prefix support for environment variables.

## License

This project is licensed under the [MIT License](./LICENSE).
