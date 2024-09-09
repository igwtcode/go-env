package env

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/igwtcode/go-env/internal/topt"
)

var (
	// AWS region name validation (e.g., us-east-1)
	awsRegionRgx = regexp.MustCompile(`^[a-z]{2}-[a-z-]+-[0-9]+$`)

	// AWS account ID validation (e.g., 123456789012)
	awsAccountIdRgx = regexp.MustCompile(`^\d{12}$`)

	// AWS S3 bucket name validation
	awsBucketNameRgx = regexp.MustCompile(`^[a-z0-9.-]{3,63}$`)

	// AWS IAM Role ARN validation (e.g., arn:aws:iam::123456789012:role/MyRole)
	awsRoleArnRgx = regexp.MustCompile(`^arn:aws:iam::\d{12}:role\/[a-zA-Z_+=,.@\-]{1,64}$`)
)

// Validation options map for v_aws_xxx exclusive options
var awsValidationMap = map[string]func(string) error{
	topt.V_AWS_REGION:      vAwsRegion,
	topt.V_AWS_ACCOUNT_ID:  vAwsAccountID,
	topt.V_AWS_BUCKET_NAME: vAwsBucketName,
	topt.V_AWS_ROLE_ARN:    vAwsRoleArn,
}

// vAwsRegion checks whether the provided AWS region name is valid based on the standard format.
// The valid format is "xx-xxxx-00" where 'x' represents lowercase letters and digits represent numbers.
//
// Returns an error if the validation fails.
func vAwsRegion(region string) error {
	if !awsRegionRgx.MatchString(region) {
		return fmt.Errorf("invalid AWS region name: %v. Expected format: xx-xxxx-00", region)
	}
	return nil
}

// vAwsAccountID checks whether the provided AWS account ID is valid.
// The account ID must be a 12-digit number.
//
// Returns an error if the validation fails.
func vAwsAccountID(id string) error {
	if !awsAccountIdRgx.MatchString(id) {
		return fmt.Errorf("invalid AWS account ID: %v. Must be a 12-digit number", id)
	}
	return nil
}

// vAwsBucketName validates the AWS bucket name.
//
// A valid AWS bucket name must be between 3 and 63 characters long and contain only lowercase letters, numbers, hyphens, and periods.
// It must not start or end with a period or hyphen, and it must not contain consecutive periods.
func vAwsBucketName(name string) error {
	// First, check if it matches the basic pattern
	if !awsBucketNameRgx.MatchString(name) {
		return fmt.Errorf("invalid AWS bucket name: %v. Must be between 3 and 63 characters long, containing only lowercase letters, numbers, hyphens, and periods", name)
	}

	// Ensure it doesn't start or end with a period or hyphen
	if strings.HasPrefix(name, ".") || strings.HasPrefix(name, "-") || strings.HasSuffix(name, ".") || strings.HasSuffix(name, "-") {
		return fmt.Errorf("invalid AWS bucket name: %v. Must not start or end with a period or hyphen", name)
	}

	// Ensure there are no consecutive periods
	if strings.Contains(name, "..") {
		return fmt.Errorf("invalid AWS bucket name: %v. Must not contain consecutive periods", name)
	}

	// If all checks pass, return nil (no error)
	return nil
}

// vAwsRoleArn checks whether the provided AWS Role ARN is valid.
//
// An AWS Role ARN should follow this pattern: arn:aws:iam::account-id:role/role-name
// where the account ID is a 12-digit number, and the role name is 1-64 characters long,
// consisting of letters, digits, and special characters.
//
// Returns an error if the validation fails.
func vAwsRoleArn(arn string) error {
	if !awsRoleArnRgx.MatchString(arn) {
		return fmt.Errorf("invalid AWS role ARN: %v. Must be in the format 'arn:aws:iam::account-id:role/role-name'", arn)
	}
	return nil
}
