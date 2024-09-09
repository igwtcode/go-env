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
