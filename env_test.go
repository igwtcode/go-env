package env_test

import (
	"os"
	"testing"

	"github.com/igwtcode/go-env"
)

func TestBasicValues(t *testing.T) {
	type Config struct {
		AppUser string `env:""`
		Mode    string `env:"required"`
		Timeout int    `env:"name=time_out,required,notrim"`
	}

	os.Setenv("mode", "PROD")
	os.Setenv("AppUser", "john")
	os.Setenv("time_out", "30")
	defer os.Unsetenv("mode")
	defer os.Unsetenv("AppUser")
	defer os.Unsetenv("time_out")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.AppUser != "john" {
		t.Errorf("expected AppUser to be 'john', got %v", cfg.AppUser)
	}
	if cfg.Mode != "PROD" {
		t.Errorf("expected Mode to be 'PROD', got %v", cfg.Mode)
	}
	if cfg.Timeout != 30 {
		t.Errorf("expected Timeout to be '30', got %v", cfg.Timeout)
	}
}

func TestSliceValues(t *testing.T) {
	type Config struct {
		Hosts  []string  `env:"name=hostslist|TARGET_HOSTS,default=localhost|127.0.0.1"`
		Limits []float32 `env:"name=limits,default=1.2|2.1|3.3"`
		IDs    []uint16  `env:"name=user_id_list"`
	}

	os.Setenv("user_id_list", "23456|13332|45678|31234")
	os.Setenv("TARGET_HOSTS", "host1|host2|host3")
	defer os.Unsetenv("user_id_list")
	defer os.Unsetenv("TARGET_HOSTS")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedHosts := []string{"host1", "host2", "host3"}
	if len(cfg.Hosts) != len(expectedHosts) {
		t.Errorf("expected %d hosts, got %d", len(expectedHosts), len(cfg.Hosts))
	}
	for i, host := range expectedHosts {
		if cfg.Hosts[i] != host {
			t.Errorf("expected Hosts[%d] to be %v, got %v", i, host, cfg.Hosts[i])
		}
	}

	expectedLimits := []float32{1.2, 2.1, 3.3}
	if len(cfg.Limits) != len(expectedLimits) {
		t.Errorf("expected %d limits, got %d", len(expectedLimits), len(cfg.Limits))
	}
	for i, limit := range expectedLimits {
		if cfg.Limits[i] != limit {
			t.Errorf("expected Limits[%d] to be %v, got %v", i, limit, cfg.Limits[i])
		}
	}

	expectedIDs := []uint16{uint16(23456), uint16(13332), uint16(45678), uint16(31234)}
	if len(cfg.IDs) != len(expectedIDs) {
		t.Errorf("expected %d IDs, got %d", len(expectedIDs), len(cfg.IDs))
	}
	for i, id := range expectedIDs {
		if cfg.IDs[i] != id {
			t.Errorf("expected IDs[%d] to be %v, got %v", i, id, cfg.IDs[i])
		}
	}
}

func TestNumericAndBooleanValues(t *testing.T) {
	type Config struct {
		Port  uint `env:"min=1024,max=65535"`
		Debug bool `env:"name=DEBUG,default=false"`
	}

	os.Setenv("PORT", "8080")
	os.Setenv("DEBUG", "true")
	defer os.Unsetenv("PORT")
	defer os.Unsetenv("DEBUG")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Port != 8080 {
		t.Errorf("expected Port to be '8080', got %v", cfg.Port)
	}
	if cfg.Debug != true {
		t.Errorf("expected Debug to be 'true', got %v", cfg.Debug)
	}
}

func Test2(t *testing.T) {
	type Config struct {
		Mode string `env:"required"`
	}

	os.Clearenv()

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err == nil {
		t.Fatalf("expected an error for missing required field, got none")
	}
}

func Test3(t *testing.T) {
	type Config struct {
		Modes []string `env:"upper"`
	}

	os.Setenv("MYAPP_MODES", "dev/Staging/prOd")
	defer os.Unsetenv("MYAPP_MODES")

	parser := env.NewParser().WithSliceValueSeparator("/").WithNamePrefix("MYAPP_")
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []string{"DEV", "STAGING", "PROD"}
	if len(cfg.Modes) != len(expected) {
		t.Fatalf("expected %d modes, got %d", len(expected), len(cfg.Modes))
	}
	for i, mode := range expected {
		if cfg.Modes[i] != mode {
			t.Errorf("expected Modes[%d] to be %v, got %v", i, mode, cfg.Modes[i])
		}
	}
}

func Test4(t *testing.T) {
	type Config struct {
		Hosts []string `env:"name=hostlist,target_hosts#lower"`
	}

	os.Setenv("MYAPP_target_hosts", "example1.com,EXAMPLE2.COM,Example3.com")
	defer os.Unsetenv("MYAPP_target_hosts")

	parser := env.NewParser().
		WithTagOptionSeparator("#").
		WithSliceValueSeparator(",").
		WithNamePrefix("MYAPP_")
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []string{"example1.com", "example2.com", "example3.com"}
	if len(cfg.Hosts) != len(expected) {
		t.Fatalf("expected %d hosts, got %d", len(expected), len(cfg.Hosts))
	}
	for i, mode := range expected {
		if cfg.Hosts[i] != mode {
			t.Errorf("expected Hosts[%d] to be %v, got %v", i, mode, cfg.Hosts[i])
		}
	}
}

func Test5(t *testing.T) {
	type Config struct {
		WrongUint uint16 `env:"name=wrong_uint"`
	}

	os.Setenv("wrong_uint", "-ab536")
	defer os.Unsetenv("wrong_uint")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err == nil {
		t.Fatalf("expected an error for invalid uint value, got none")
	}
}

// 1. Test basic string field assignment
func TestBasicStringAssignment(t *testing.T) {
	type Config struct {
		Host string `env:"name=HOST"`
	}

	os.Setenv("HOST", "localhost")
	defer os.Unsetenv("HOST")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Host != "localhost" {
		t.Errorf("expected Host to be 'localhost', got %v", cfg.Host)
	}
}

// 2. Test basic int field assignment
func TestBasicIntAssignment(t *testing.T) {
	type Config struct {
		Port int `env:"name=PORT"`
	}

	os.Setenv("PORT", "8080")
	defer os.Unsetenv("PORT")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Port != 8080 {
		t.Errorf("expected Port to be '8080', got %v", cfg.Port)
	}
}

// 3. Test basic float field assignment
func TestBasicFloatAssignment(t *testing.T) {
	type Config struct {
		Rate float64 `env:"name=RATE"`
	}

	os.Setenv("RATE", "3.14")
	defer os.Unsetenv("RATE")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Rate != 3.14 {
		t.Errorf("expected Rate to be '3.14', got %v", cfg.Rate)
	}
}

// 4. Test basic bool field assignment
func TestBasicBoolAssignment(t *testing.T) {
	type Config struct {
		Debug bool `env:"name=DEBUG"`
	}

	os.Setenv("DEBUG", "true")
	defer os.Unsetenv("DEBUG")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !cfg.Debug {
		t.Errorf("expected Debug to be 'true', got %v", cfg.Debug)
	}
}

// 5. Test slice of strings with default separator
func TestSliceOfStringWithDefaultSeparator(t *testing.T) {
	type Config struct {
		Hosts []string `env:"name=HOSTS"`
	}

	os.Setenv("HOSTS", "host1|host2|host3")
	defer os.Unsetenv("HOSTS")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []string{"host1", "host2", "host3"}
	if len(cfg.Hosts) != len(expected) {
		t.Fatalf("expected %d hosts, got %d", len(expected), len(cfg.Hosts))
	}

	for i, host := range expected {
		if cfg.Hosts[i] != host {
			t.Errorf("expected Hosts[%d] to be %v, got %v", i, host, cfg.Hosts[i])
		}
	}
}

// 6. Test slice of integers with default separator
func TestSliceOfIntWithDefaultSeparator(t *testing.T) {
	type Config struct {
		Ports []int `env:"name=PORTS"`
	}

	os.Setenv("PORTS", "8080|8081|8082")
	defer os.Unsetenv("PORTS")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []int{8080, 8081, 8082}
	if len(cfg.Ports) != len(expected) {
		t.Fatalf("expected %d ports, got %d", len(expected), len(cfg.Ports))
	}

	for i, port := range expected {
		if cfg.Ports[i] != port {
			t.Errorf("expected Ports[%d] to be %v, got %v", i, port, cfg.Ports[i])
		}
	}
}

// 7. Test custom slice value separator
func TestCustomSliceValueSeparator(t *testing.T) {
	type Config struct {
		Hosts []string `env:"name=HOSTS"`
	}

	os.Setenv("HOSTS", "host1,host2,host3")
	defer os.Unsetenv("HOSTS")

	parser := env.NewParser().WithTagOptionSeparator(";").WithSliceValueSeparator(",")
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []string{"host1", "host2", "host3"}
	if len(cfg.Hosts) != len(expected) {
		t.Fatalf("expected %d hosts, got %d", len(expected), len(cfg.Hosts))
	}

	for i, host := range expected {
		if cfg.Hosts[i] != host {
			t.Errorf("expected Hosts[%d] to be %v, got %v", i, host, cfg.Hosts[i])
		}
	}
}

// 8. Test custom tag option separator
func TestCustomTagOptionSeparator(t *testing.T) {
	type Config struct {
		Host string `env:"name=HOST;default=localhost"`
	}

	os.Clearenv()

	parser := env.NewParser().WithTagOptionSeparator(";")
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Host != "localhost" {
		t.Errorf("expected Host to be 'localhost', got %v", cfg.Host)
	}
}

// 9. Test required field without environment variable set
func TestRequiredFieldWithoutEnv(t *testing.T) {
	type Config struct {
		Host string `env:"name=HOST,required"`
	}

	os.Clearenv()

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err == nil {
		t.Fatalf("expected an error for missing required field, got none")
	}
}

// 10. Test handling of nested structs
func TestNestedStruct(t *testing.T) {
	type Database struct {
		Host string `env:"name=DB_HOST"`
		Port int    `env:"name=DB_PORT"`
	}
	type Config struct {
		Database Database
	}

	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	defer os.Unsetenv("DB_HOST")
	defer os.Unsetenv("DB_PORT")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Database.Host != "localhost" {
		t.Errorf("expected Database.Host to be 'localhost', got %v", cfg.Database.Host)
	}
	if cfg.Database.Port != 5432 {
		t.Errorf("expected Database.Port to be '5432', got %v", cfg.Database.Port)
	}
}

// 11. Test default value for string field when environment variable is missing
func TestStringDefaultValue(t *testing.T) {
	type Config struct {
		Host string `env:"name=HOST,default=localhost"`
	}

	os.Clearenv()

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Host != "localhost" {
		t.Errorf("expected Host to default to 'localhost', got %v", cfg.Host)
	}
}

// 12. Test default value for int field when environment variable is missing
func TestIntDefaultValue(t *testing.T) {
	type Config struct {
		Port int `env:"name=PORT,default=8080"`
	}

	os.Clearenv()

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Port != 8080 {
		t.Errorf("expected Port to default to 8080, got %v", cfg.Port)
	}
}

// 13. Test default value for bool field when environment variable is missing
func TestBoolDefaultValue(t *testing.T) {
	type Config struct {
		Debug bool `env:"name=DEBUG,default=true"`
	}

	os.Clearenv()

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !cfg.Debug {
		t.Errorf("expected Debug to default to true, got %v", cfg.Debug)
	}
}

// 14. Test default value for float field when environment variable is missing
func TestFloatDefaultValue(t *testing.T) {
	type Config struct {
		Rate float64 `env:"name=RATE,default=3.14"`
	}

	os.Clearenv()

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Rate != 3.14 {
		t.Errorf("expected Rate to default to 3.14, got %v", cfg.Rate)
	}
}

// 15. Test slice of strings with default value
func TestSliceOfStringWithDefaultValue(t *testing.T) {
	type Config struct {
		Hosts []string `env:"name=HOSTS,default=host1|host2|host3"`
	}

	os.Clearenv()

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []string{"host1", "host2", "host3"}
	if len(cfg.Hosts) != len(expected) {
		t.Fatalf("expected %d hosts, got %d", len(expected), len(cfg.Hosts))
	}

	for i, host := range expected {
		if cfg.Hosts[i] != host {
			t.Errorf("expected Hosts[%d] to be %v, got %v", i, host, cfg.Hosts[i])
		}
	}
}

// 16. Test int slice with default value
func TestSliceOfIntWithDefaultValue(t *testing.T) {
	type Config struct {
		Ports []int `env:"name=PORTS,default=8080|8081|8082"`
	}

	os.Clearenv()

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []int{8080, 8081, 8082}
	if len(cfg.Ports) != len(expected) {
		t.Fatalf("expected %d ports, got %d", len(expected), len(cfg.Ports))
	}

	for i, port := range expected {
		if cfg.Ports[i] != port {
			t.Errorf("expected Ports[%d] to be %v, got %v", i, port, cfg.Ports[i])
		}
	}
}

// 17. Test required string field when environment variable is set
func TestRequiredStringFieldWithEnv(t *testing.T) {
	type Config struct {
		Host string `env:"name=HOST,required"`
	}

	os.Setenv("HOST", "localhost")
	defer os.Unsetenv("HOST")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Host != "localhost" {
		t.Errorf("expected Host to be 'localhost', got %v", cfg.Host)
	}
}

// 18. Test required int field when environment variable is set
func TestRequiredIntFieldWithEnv(t *testing.T) {
	type Config struct {
		Port int `env:"name=PORT,required"`
	}

	os.Setenv("PORT", "8080")
	defer os.Unsetenv("PORT")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Port != 8080 {
		t.Errorf("expected Port to be 8080, got %v", cfg.Port)
	}
}

// 19. Test required field with multiple names (fallback behavior)
func TestRequiredFieldWithMultipleNames(t *testing.T) {
	type Config struct {
		Region string `env:"name=AWS_DEFAULT_REGION|AWS_REGION,required"`
	}

	os.Setenv("AWS_REGION", "us-west-1")
	defer os.Unsetenv("AWS_REGION")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Region != "us-west-1" {
		t.Errorf("expected Region to be 'us-west-1', got %v", cfg.Region)
	}
}

// 20. Test handling of nested struct with default values
func TestNestedStructWithDefaultValues(t *testing.T) {
	type Database struct {
		Host string `env:"name=DB_HOST,default=localhost"`
		Port int    `env:"name=DB_PORT,default=5432"`
	}
	type Config struct {
		Database Database
	}

	os.Clearenv()

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Database.Host != "localhost" {
		t.Errorf("expected Database.Host to default to 'localhost', got %v", cfg.Database.Host)
	}
	if cfg.Database.Port != 5432 {
		t.Errorf("expected Database.Port to default to 5432, got %v", cfg.Database.Port)
	}
}

// 21. Test slice of booleans with default value
func TestSliceOfBoolWithDefaultValue(t *testing.T) {
	type Config struct {
		Flags []bool `env:"name=FLAGS,default=true|false|true"`
	}

	os.Clearenv()

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []bool{true, false, true}
	if len(cfg.Flags) != len(expected) {
		t.Fatalf("expected %d flags, got %d", len(expected), len(cfg.Flags))
	}

	for i, flag := range expected {
		if cfg.Flags[i] != flag {
			t.Errorf("expected Flags[%d] to be %v, got %v", i, flag, cfg.Flags[i])
		}
	}
}

// 22. Test invalid int value in environment variable
func TestInvalidIntValue(t *testing.T) {
	type Config struct {
		Port int `env:"name=PORT"`
	}

	os.Setenv("PORT", "invalid")
	defer os.Unsetenv("PORT")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err == nil {
		t.Fatalf("expected an error for invalid int value, got none")
	}
}

// 23. Test invalid bool value in environment variable
func TestInvalidBoolValue(t *testing.T) {
	type Config struct {
		Debug bool `env:"name=DEBUG"`
	}

	os.Setenv("DEBUG", "invalid_bool")
	defer os.Unsetenv("DEBUG")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err == nil {
		t.Fatalf("expected an error for invalid bool value, got none")
	}
}

// 24. Test invalid float value in environment variable
func TestInvalidFloatValue(t *testing.T) {
	type Config struct {
		Rate float64 `env:"name=RATE"`
	}

	os.Setenv("RATE", "invalid_float")
	defer os.Unsetenv("RATE")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err == nil {
		t.Fatalf("expected an error for invalid float value, got none")
	}
}

// 25. Test default value with invalid type in environment variable (int field with string default)
func TestDefaultWithInvalidType(t *testing.T) {
	type Config struct {
		Port int `env:"name=PORT,default=invalid"`
	}

	os.Clearenv()

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err == nil {
		t.Fatalf("expected an error for invalid default type, got none")
	}
}

// 26. Test handling of min/max for int field (valid value)
func TestMinMaxForIntFieldValid(t *testing.T) {
	type Config struct {
		Port int `env:"name=PORT,min=8000,max=9000"`
	}

	os.Setenv("PORT", "8080")
	defer os.Unsetenv("PORT")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Port != 8080 {
		t.Errorf("expected Port to be 8080, got %v", cfg.Port)
	}
}

// 27. Test handling of min/max for int field (invalid value below min)
func TestMinMaxForIntFieldInvalidBelowMin(t *testing.T) {
	type Config struct {
		Port int `env:"name=PORT,min=8000,max=9000"`
	}

	os.Setenv("PORT", "7000")
	defer os.Unsetenv("PORT")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err == nil {
		t.Fatalf("expected an error for value below min, got none")
	}
}

// 28. Test handling of min/max for int field (invalid value above max)
func TestMinMaxForIntFieldInvalidAboveMax(t *testing.T) {
	type Config struct {
		Port int `env:"name=PORT,min=8000,max=9000"`
	}

	os.Setenv("PORT", "9100")
	defer os.Unsetenv("PORT")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err == nil {
		t.Fatalf("expected an error for value above max, got none")
	}
}

// 29. Test lowercase option for string field
func TestLowercaseOptionForStringField(t *testing.T) {
	type Config struct {
		Host string `env:"name=HOST,lower"`
	}

	os.Setenv("HOST", "LOCALHOST")
	defer os.Unsetenv("HOST")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Host != "localhost" {
		t.Errorf("expected Host to be 'localhost' after applying 'lower', got %v", cfg.Host)
	}
}

// 30. Test uppercase option for string field
func TestUppercaseOptionForStringField(t *testing.T) {
	type Config struct {
		Host string `env:"name=HOST,upper"`
	}

	os.Setenv("HOST", "localhost")
	defer os.Unsetenv("HOST")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Host != "LOCALHOST" {
		t.Errorf("expected Host to be 'LOCALHOST' after applying 'upper', got %v", cfg.Host)
	}
}

// 31. Test trim option is applied by default for string field
func TestTrimAppliedByDefault(t *testing.T) {
	type Config struct {
		Host string `env:"name=HOST"`
	}

	os.Setenv("HOST", "  localhost  ")
	defer os.Unsetenv("HOST")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Host != "localhost" {
		t.Errorf("expected Host to be 'localhost' with default trimming, got '%v'", cfg.Host)
	}
}

// 32. Test notrim option prevents trimming of string field
func TestNoTrimOption(t *testing.T) {
	type Config struct {
		Host string `env:"name=HOST,notrim"`
	}

	os.Setenv("HOST", "  localhost  ")
	defer os.Unsetenv("HOST")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Host != "  localhost  " {
		t.Errorf("expected Host to retain spaces with 'notrim', got '%v'", cfg.Host)
	}
}

// 33. Test case-insensitive lookup for environment variable
func TestCaseInsensitiveEnvVarLookup(t *testing.T) {
	type Config struct {
		Region string `env:"name=REGION"`
	}

	os.Setenv("region", "us-west-1") // Lowercase variable name in the environment
	defer os.Unsetenv("region")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Region != "us-west-1" {
		t.Errorf("expected Region to be 'us-west-1', got %v", cfg.Region)
	}
}

// 34. Test empty slice when no environment variable is set
func TestEmptySliceWhenNoEnvVar(t *testing.T) {
	type Config struct {
		Hosts []string `env:"name=HOSTS"`
	}

	os.Clearenv()

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(cfg.Hosts) != 0 {
		t.Errorf("expected Hosts to be empty, got %v", cfg.Hosts)
	}
}

// 35. Test multiple environment variables with fallback
func TestMultipleEnvVarsFallback(t *testing.T) {
	type Config struct {
		Region string `env:"name=AWS_DEFAULT_REGION|AWS_REGION"`
	}

	os.Setenv("AWS_DEFAULT_REGION", "")
	os.Setenv("AWS_REGION", "us-east-1")
	defer os.Unsetenv("AWS_DEFAULT_REGION")
	defer os.Unsetenv("AWS_REGION")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Region != "us-east-1" {
		t.Errorf("expected Region to be 'us-east-1', got %v", cfg.Region)
	}
}

// 36. Test slice with invalid values in environment variable (e.g., non-integer in int slice)
func TestInvalidSliceValues(t *testing.T) {
	type Config struct {
		Ports []int `env:"name=PORTS"`
	}

	os.Setenv("PORTS", "8080|invalid|8082")
	defer os.Unsetenv("PORTS")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err == nil {
		t.Fatalf("expected an error for invalid int value in slice, got none")
	}
}

// 37. Test default slice separator with custom separator for options
func TestDefaultSliceSeparatorWithCustomOptionSeparator(t *testing.T) {
	type Config struct {
		Hosts []string `env:"name=HOSTS;default=host1|host2|host3"`
	}

	os.Clearenv()

	parser := env.NewParser().WithTagOptionSeparator(";")
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []string{"host1", "host2", "host3"}
	if len(cfg.Hosts) != len(expected) {
		t.Fatalf("expected %d hosts, got %d", len(expected), len(cfg.Hosts))
	}

	for i, host := range expected {
		if cfg.Hosts[i] != host {
			t.Errorf("expected Hosts[%d] to be %v, got %v", i, host, cfg.Hosts[i])
		}
	}
}

// 38. Test custom slice separator with custom option separator
func TestCustomSliceSeparatorWithCustomOptionSeparator(t *testing.T) {
	type Config struct {
		Hosts []string `env:"name=HOSTS;default=host1/host2/host3"`
	}

	os.Clearenv()

	parser := env.NewParser().WithTagOptionSeparator(";").WithSliceValueSeparator("/")
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []string{"host1", "host2", "host3"}
	if len(cfg.Hosts) != len(expected) {
		t.Fatalf("expected %d hosts, got %d", len(expected), len(cfg.Hosts))
	}

	for i, host := range expected {
		if cfg.Hosts[i] != host {
			t.Errorf("expected Hosts[%d] to be %v, got %v", i, host, cfg.Hosts[i])
		}
	}
}

// 39. Test handling of default slice when environment variable is not set
func TestDefaultSliceWhenNoEnvVar(t *testing.T) {
	type Config struct {
		Ports []int `env:"name=PORTS,default=8080|8081|8082"`
	}

	os.Clearenv()

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []int{8080, 8081, 8082}
	if len(cfg.Ports) != len(expected) {
		t.Fatalf("expected %d ports, got %d", len(expected), len(cfg.Ports))
	}

	for i, port := range expected {
		if cfg.Ports[i] != port {
			t.Errorf("expected Ports[%d] to be %v, got %v", i, port, cfg.Ports[i])
		}
	}
}

// 40. Test default value when required tag is missing
func TestDefaultWhenRequiredTagMissing(t *testing.T) {
	type Config struct {
		Host string `env:"name=HOST,default=localhost"`
	}

	os.Clearenv()

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Host != "localhost" {
		t.Errorf("expected Host to default to 'localhost', got %v", cfg.Host)
	}
}

// 41. Test default value for nested struct when environment variable is missing
func TestDefaultValueForNestedStruct(t *testing.T) {
	type Database struct {
		Host string `env:"name=DB_HOST,default=localhost"`
		Port int    `env:"name=DB_PORT,default=5432"`
	}
	type Config struct {
		Database Database
	}

	os.Clearenv()

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Database.Host != "localhost" {
		t.Errorf("expected Database.Host to be 'localhost', got %v", cfg.Database.Host)
	}
	if cfg.Database.Port != 5432 {
		t.Errorf("expected Database.Port to be 5432, got %v", cfg.Database.Port)
	}
}

// 42. Test nested struct with environment variable set for one field
func TestNestedStructWithEnvVarSetForOneField(t *testing.T) {
	type Database struct {
		Host string `env:"name=DB_HOST"`
		Port int    `env:"name=DB_PORT,default=5432"`
	}
	type Config struct {
		Database Database
	}

	os.Setenv("DB_HOST", "db.example.com")
	defer os.Unsetenv("DB_HOST")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Database.Host != "db.example.com" {
		t.Errorf("expected Database.Host to be 'db.example.com', got %v", cfg.Database.Host)
	}
	if cfg.Database.Port != 5432 {
		t.Errorf("expected Database.Port to be 5432, got %v", cfg.Database.Port)
	}
}

// 43. Test nested struct with environment variable set for all fields
func TestNestedStructWithEnvVarsSetForAllFields(t *testing.T) {
	type Database struct {
		Host string `env:"name=DB_HOST"`
		Port int    `env:"name=DB_PORT"`
	}
	type Config struct {
		Database Database
	}

	os.Setenv("DB_HOST", "db.example.com")
	os.Setenv("DB_PORT", "3306")
	defer os.Unsetenv("DB_HOST")
	defer os.Unsetenv("DB_PORT")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Database.Host != "db.example.com" {
		t.Errorf("expected Database.Host to be 'db.example.com', got %v", cfg.Database.Host)
	}
	if cfg.Database.Port != 3306 {
		t.Errorf("expected Database.Port to be 3306, got %v", cfg.Database.Port)
	}
}

// 44. Test invalid int value in nested struct
func TestInvalidIntInNestedStruct(t *testing.T) {
	type Database struct {
		Host string `env:"name=DB_HOST"`
		Port int    `env:"name=DB_PORT"`
	}
	type Config struct {
		Database Database
	}

	os.Setenv("DB_HOST", "db.example.com")
	os.Setenv("DB_PORT", "invalid_port")
	defer os.Unsetenv("DB_HOST")
	defer os.Unsetenv("DB_PORT")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err == nil {
		t.Fatalf("expected an error for invalid int value in nested struct, got none")
	}
}

// 45. Test setting float field with invalid float value in environment variable
func TestInvalidFloatValueInEnvVar(t *testing.T) {
	type Config struct {
		Rate float64 `env:"name=RATE"`
	}

	os.Setenv("RATE", "invalid_float")
	defer os.Unsetenv("RATE")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err == nil {
		t.Fatalf("expected an error for invalid float value, got none")
	}
}

// 46. Test handling of empty string for string field
func TestEmptyStringValue(t *testing.T) {
	type Config struct {
		Host string `env:"name=HOST"`
	}

	os.Setenv("HOST", "")
	defer os.Unsetenv("HOST")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Host != "" {
		t.Errorf("expected Host to be empty, got '%v'", cfg.Host)
	}
}

// 47. Test handling of empty int field with default
func TestEmptyIntFieldWithDefault(t *testing.T) {
	type Config struct {
		Port int `env:"name=PORT,default=8080"`
	}

	os.Setenv("PORT", "")
	defer os.Unsetenv("PORT")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Port != 8080 {
		t.Errorf("expected Port to default to 8080, got %v", cfg.Port)
	}
}

// 48. Test handling of empty float field with default
func TestEmptyFloatFieldWithDefault(t *testing.T) {
	type Config struct {
		Rate float64 `env:"name=RATE,default=3.14"`
	}

	os.Setenv("RATE", "")
	defer os.Unsetenv("RATE")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Rate != 3.14 {
		t.Errorf("expected Rate to default to 3.14, got %v", cfg.Rate)
	}
}

// 49. Test case with missing required field
func TestMissingRequiredField(t *testing.T) {
	type Config struct {
		Host string `env:"name=HOST,required"`
	}

	os.Clearenv()

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err == nil {
		t.Fatalf("expected an error for missing required field, got none")
	}
}

// 50. Test slice with no values (empty slice)
func TestEmptySlice(t *testing.T) {
	type Config struct {
		Hosts []string `env:"name=HOSTS"`
	}

	os.Setenv("HOSTS", "")
	defer os.Unsetenv("HOSTS")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(cfg.Hosts) != 0 {
		t.Errorf("expected Hosts to be an empty slice, got %v", cfg.Hosts)
	}
}

// 51. Test min and max validation for float field (valid value)
func TestMinMaxForFloatFieldValid(t *testing.T) {
	type Config struct {
		Rate float64 `env:"name=RATE,min=1.5,max=4.0"`
	}

	os.Setenv("RATE", "3.0")
	defer os.Unsetenv("RATE")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Rate != 3.0 {
		t.Errorf("expected Rate to be 3.0, got %v", cfg.Rate)
	}
}

// 52. Test min validation for float field (value below min)
func TestMinForFloatFieldInvalid(t *testing.T) {
	type Config struct {
		Rate float64 `env:"name=RATE,min=1.5"`
	}

	os.Setenv("RATE", "1.0")
	defer os.Unsetenv("RATE")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err == nil {
		t.Fatalf("expected an error for value below min, got none")
	}
}

// 53. Test max validation for float field (value above max)
func TestMaxForFloatFieldInvalid(t *testing.T) {
	type Config struct {
		Rate float64 `env:"name=RATE,max=4.0"`
	}

	os.Setenv("RATE", "5.0")
	defer os.Unsetenv("RATE")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err == nil {
		t.Fatalf("expected an error for value above max, got none")
	}
}

// 54. Test required int field with empty string in environment variable
func TestRequiredIntFieldWithEmptyString(t *testing.T) {
	type Config struct {
		Port int `env:"name=PORT,required"`
	}

	os.Setenv("PORT", "")
	defer os.Unsetenv("PORT")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err == nil {
		t.Fatalf("expected an error for missing required int field, got none")
	}
}

// 55. Test required string field with whitespace value
func TestRequiredStringFieldWithWhitespace(t *testing.T) {
	type Config struct {
		Host string `env:"name=HOST,required"`
	}

	os.Setenv("HOST", "   ")
	defer os.Unsetenv("HOST")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err == nil {
		t.Fatalf("expected an error for whitespace value in required field, got none")
	}
}

// 56. Test bool field with "1" and "0" as valid boolean values
func TestBoolFieldWithOneAndZero(t *testing.T) {
	type Config struct {
		Debug bool `env:"name=DEBUG"`
	}

	os.Setenv("DEBUG", "1")
	defer os.Unsetenv("DEBUG")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !cfg.Debug {
		t.Errorf("expected Debug to be true, got %v", cfg.Debug)
	}

	// Test with "0"
	os.Setenv("DEBUG", "0")
	err = parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Debug {
		t.Errorf("expected Debug to be false, got %v", cfg.Debug)
	}
}

// 57. Test bool field with invalid string ("yes" or "no")
func TestInvalidBoolField(t *testing.T) {
	type Config struct {
		Debug bool `env:"name=DEBUG"`
	}

	os.Setenv("DEBUG", "yes")
	defer os.Unsetenv("DEBUG")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err == nil {
		t.Fatalf("expected an error for invalid bool value, got none")
	}
}

// 58. Test default slice separator on empty slice
func TestEmptySliceWithDefaultSeparator(t *testing.T) {
	type Config struct {
		Hosts []string `env:"name=HOSTS"`
	}

	os.Setenv("HOSTS", "")
	defer os.Unsetenv("HOSTS")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(cfg.Hosts) != 0 {
		t.Errorf("expected Hosts to be an empty slice, got %v", cfg.Hosts)
	}
}

// 59. Test handling of default slice with custom separator
func TestDefaultSliceWithCustomSeparator(t *testing.T) {
	type Config struct {
		Ports []int `env:"name=PORTS,default=8080/8081/8082"`
	}

	os.Clearenv()

	parser := env.NewParser().WithSliceValueSeparator("/")
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []int{8080, 8081, 8082}
	if len(cfg.Ports) != len(expected) {
		t.Fatalf("expected %d ports, got %d", len(expected), len(cfg.Ports))
	}

	for i, port := range expected {
		if cfg.Ports[i] != port {
			t.Errorf("expected Ports[%d] to be %v, got %v", i, port, cfg.Ports[i])
		}
	}
}

// 60. Test setting multiple environment variables with default values
func TestMultipleEnvVarsWithDefaultValues(t *testing.T) {
	type Config struct {
		Region string `env:"name=AWS_DEFAULT_REGION/AWS_REGION,default=us-east-1"`
	}

	os.Clearenv()

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Region != "us-east-1" {
		t.Errorf("expected Region to default to 'us-east-1', got %v", cfg.Region)
	}
}

// 61. Test int field with min value equal to the value
func TestMinEqualToValueForIntField(t *testing.T) {
	type Config struct {
		Port int `env:"name=PORT,min=8080"`
	}

	os.Setenv("PORT", "8080")
	defer os.Unsetenv("PORT")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Port != 8080 {
		t.Errorf("expected Port to be 8080, got %v", cfg.Port)
	}
}

// 62. Test float field with max value equal to the value
func TestMaxEqualToValueForFloatField(t *testing.T) {
	type Config struct {
		Rate float64 `env:"name=RATE,max=3.14"`
	}

	os.Setenv("RATE", "3.14")
	defer os.Unsetenv("RATE")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Rate != 3.14 {
		t.Errorf("expected Rate to be 3.14, got %v", cfg.Rate)
	}
}

// 63. Test float field with value equal to min and max
func TestMinAndMaxEqualForFloatField(t *testing.T) {
	type Config struct {
		Rate float64 `env:"name=RATE,min=3.14,max=3.14"`
	}

	os.Setenv("RATE", "3.14")
	defer os.Unsetenv("RATE")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Rate != 3.14 {
		t.Errorf("expected Rate to be 3.14, got %v", cfg.Rate)
	}
}

// 64. Test default slice separator when there are multiple spaces between values
func TestSliceWithMultipleSpacesBetweenValues(t *testing.T) {
	type Config struct {
		Hosts []string `env:"name=HOSTS"`
	}

	os.Setenv("HOSTS", "host1  |  host2  |  host3")
	defer os.Unsetenv("HOSTS")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []string{"host1", "host2", "host3"}
	if len(cfg.Hosts) != len(expected) {
		t.Fatalf("expected %d hosts, got %d", len(expected), len(cfg.Hosts))
	}

	for i, host := range expected {
		if cfg.Hosts[i] != host {
			t.Errorf("expected Hosts[%d] to be %v, got %v", i, host, cfg.Hosts[i])
		}
	}
}

// 65. Test custom slice separator with multiple spaces between values
func TestCustomSliceWithMultipleSpacesBetweenValues(t *testing.T) {
	type Config struct {
		Ports []int `env:"name=PORTS"`
	}

	os.Setenv("PORTS", "8080  /  8081  /  8082")
	defer os.Unsetenv("PORTS")

	parser := env.NewParser().WithSliceValueSeparator("/")
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []int{8080, 8081, 8082}
	if len(cfg.Ports) != len(expected) {
		t.Fatalf("expected %d ports, got %d", len(expected), len(cfg.Ports))
	}

	for i, port := range expected {
		if cfg.Ports[i] != port {
			t.Errorf("expected Ports[%d] to be %v, got %v", i, port, cfg.Ports[i])
		}
	}
}

// 66. Test slice with leading and trailing spaces in environment variable values
func TestSliceWithLeadingAndTrailingSpaces(t *testing.T) {
	type Config struct {
		Hosts []string `env:"name=HOSTS"`
	}

	os.Setenv("HOSTS", "  host1|host2  |  host3  ")
	defer os.Unsetenv("HOSTS")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []string{"host1", "host2", "host3"}
	if len(cfg.Hosts) != len(expected) {
		t.Fatalf("expected %d hosts, got %d", len(expected), len(cfg.Hosts))
	}

	for i, host := range expected {
		if cfg.Hosts[i] != host {
			t.Errorf("expected Hosts[%d] to be %v, got %v", i, host, cfg.Hosts[i])
		}
	}
}

// 67. Test invalid float value in slice
func TestInvalidFloatInSlice(t *testing.T) {
	type Config struct {
		Rates []float64 `env:"name=RATES"`
	}

	os.Setenv("RATES", "3.14|invalid|1.618")
	defer os.Unsetenv("RATES")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err == nil {
		t.Fatalf("expected an error for invalid float value in slice, got none")
	}
}

// 68. Test required float field with value "0"
func TestRequiredFloatFieldWithZeroValue(t *testing.T) {
	type Config struct {
		Rate float64 `env:"name=RATE,required"`
	}

	os.Setenv("RATE", "0")
	defer os.Unsetenv("RATE")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Rate != 0 {
		t.Errorf("expected Rate to be 0, got %v", cfg.Rate)
	}
}

// 69. Test slice with all empty elements
func TestSliceWithEmptyElements(t *testing.T) {
	type Config struct {
		Hosts []string `env:"name=HOSTS"`
	}

	os.Setenv("HOSTS", "|||")
	defer os.Unsetenv("HOSTS")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(cfg.Hosts) != 0 {
		t.Errorf("expected Hosts to be an empty slice, got %v", cfg.Hosts)
	}
}

// 70. Test custom separator with default value in environment variable
func TestCustomSeparatorWithDefaultValueInEnvVar(t *testing.T) {
	type Config struct {
		Regions []string `env:"name=REGIONS,default=us-east-1|us-west-1"`
	}

	os.Setenv("REGIONS", "us-central-1/us-east-2")
	defer os.Unsetenv("REGIONS")

	parser := env.NewParser().WithSliceValueSeparator("/")
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []string{"us-central-1", "us-east-2"}
	if len(cfg.Regions) != len(expected) {
		t.Fatalf("expected %d regions, got %d", len(expected), len(cfg.Regions))
	}

	for i, region := range expected {
		if cfg.Regions[i] != region {
			t.Errorf("expected Regions[%d] to be %v, got %v", i, region, cfg.Regions[i])
		}
	}
}

// 71. Test int field with negative value
func TestNegativeIntValue(t *testing.T) {
	type Config struct {
		Port int `env:"name=PORT"`
	}

	os.Setenv("PORT", "-8080")
	defer os.Unsetenv("PORT")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Port != -8080 {
		t.Errorf("expected Port to be -8080, got %v", cfg.Port)
	}
}

// 72. Test float field with negative value
func TestNegativeFloatValue(t *testing.T) {
	type Config struct {
		Rate float64 `env:"name=RATE"`
	}

	os.Setenv("RATE", "-3.14")
	defer os.Unsetenv("RATE")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Rate != -3.14 {
		t.Errorf("expected Rate to be -3.14, got %v", cfg.Rate)
	}
}

// 73. Test slice with single value
func TestSingleValueInSlice(t *testing.T) {
	type Config struct {
		Hosts []string `env:"name=HOSTS"`
	}

	os.Setenv("HOSTS", "localhost")
	defer os.Unsetenv("HOSTS")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []string{"localhost"}
	if len(cfg.Hosts) != len(expected) {
		t.Fatalf("expected 1 host, got %d", len(cfg.Hosts))
	}

	if cfg.Hosts[0] != expected[0] {
		t.Errorf("expected Hosts[0] to be %v, got %v", expected[0], cfg.Hosts[0])
	}
}

// 74. Test slice with trailing separator
func TestSliceWithTrailingSeparator(t *testing.T) {
	type Config struct {
		Ports []int `env:"name=PORTS"`
	}

	os.Setenv("PORTS", "8080|8081|")
	defer os.Unsetenv("PORTS")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []int{8080, 8081}
	if len(cfg.Ports) != len(expected) {
		t.Fatalf("expected %d ports, got %d", len(expected), len(cfg.Ports))
	}

	for i, port := range expected {
		if cfg.Ports[i] != port {
			t.Errorf("expected Ports[%d] to be %v, got %v", i, port, cfg.Ports[i])
		}
	}
}

// 75. Test invalid bool field with numeric value greater than 1
func TestInvalidBoolWithNumericValue(t *testing.T) {
	type Config struct {
		Debug bool `env:"name=DEBUG"`
	}

	os.Setenv("DEBUG", "2")
	defer os.Unsetenv("DEBUG")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err == nil {
		t.Fatalf("expected an error for invalid bool value '2', got none")
	}
}

// 76. Test invalid bool field with negative numeric value
func TestInvalidBoolWithNegativeNumericValue(t *testing.T) {
	type Config struct {
		Debug bool `env:"name=DEBUG"`
	}

	os.Setenv("DEBUG", "-1")
	defer os.Unsetenv("DEBUG")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err == nil {
		t.Fatalf("expected an error for invalid bool value '-1', got none")
	}
}

// 77. Test slice with mix of valid and invalid int values
func TestSliceWithMixOfValidAndInvalidIntValues(t *testing.T) {
	type Config struct {
		Ports []int `env:"name=PORTS"`
	}

	os.Setenv("PORTS", "8080|invalid|8081")
	defer os.Unsetenv("PORTS")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err == nil {
		t.Fatalf("expected an error for invalid int value, got none")
	}
}

// 78. Test string field with embedded spaces
func TestStringFieldWithEmbeddedSpaces(t *testing.T) {
	type Config struct {
		Host string `env:"name=HOST"`
	}

	os.Setenv("HOST", "my host with spaces")
	defer os.Unsetenv("HOST")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Host != "my host with spaces" {
		t.Errorf("expected Host to be 'my host with spaces', got %v", cfg.Host)
	}
}

// 79. Test custom slice separator with empty slice in environment variable
func TestCustomSliceSeparatorWithEmptySlice(t *testing.T) {
	type Config struct {
		Hosts []string `env:"name=HOSTS"`
	}

	os.Setenv("HOSTS", "")
	defer os.Unsetenv("HOSTS")

	parser := env.NewParser().WithSliceValueSeparator("/")
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(cfg.Hosts) != 0 {
		t.Errorf("expected Hosts to be an empty slice, got %v", cfg.Hosts)
	}
}

// 80. Test int field with large value
func TestLargeIntValue(t *testing.T) {
	type Config struct {
		Port int `env:"name=PORT"`
	}

	os.Setenv("PORT", "2147483647") // Maximum value for int32
	defer os.Unsetenv("PORT")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Port != 2147483647 {
		t.Errorf("expected Port to be 2147483647, got %v", cfg.Port)
	}
}

// 81. Test nested struct with both fields set through environment variables
func TestNestedStructWithEnvVars(t *testing.T) {
	type Database struct {
		Host string `env:"name=DB_HOST"`
		Port int    `env:"name=DB_PORT"`
	}
	type Config struct {
		Database Database
	}

	os.Setenv("DB_HOST", "db.example.com")
	os.Setenv("DB_PORT", "5432")
	defer os.Unsetenv("DB_HOST")
	defer os.Unsetenv("DB_PORT")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Database.Host != "db.example.com" {
		t.Errorf("expected Database.Host to be 'db.example.com', got %v", cfg.Database.Host)
	}
	if cfg.Database.Port != 5432 {
		t.Errorf("expected Database.Port to be 5432, got %v", cfg.Database.Port)
	}
}

// 82. Test nested struct with default values
func TestNestedStructWithDefaults(t *testing.T) {
	type Database struct {
		Host string `env:"name=DB_HOST,default=localhost"`
		Port int    `env:"name=DB_PORT,default=3306"`
	}
	type Config struct {
		Database Database
	}

	os.Clearenv()

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Database.Host != "localhost" {
		t.Errorf("expected Database.Host to be 'localhost', got %v", cfg.Database.Host)
	}
	if cfg.Database.Port != 3306 {
		t.Errorf("expected Database.Port to be 3306, got %v", cfg.Database.Port)
	}
}

// 83. Test nested struct with one field set through environment variables and another using default
func TestNestedStructPartialEnvVars(t *testing.T) {
	type Database struct {
		Host string `env:"name=DB_HOST,default=localhost"`
		Port int    `env:"name=DB_PORT,default=3306"`
	}
	type Config struct {
		Database Database
	}

	os.Setenv("DB_HOST", "db.example.com")
	defer os.Unsetenv("DB_HOST")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Database.Host != "db.example.com" {
		t.Errorf("expected Database.Host to be 'db.example.com', got %v", cfg.Database.Host)
	}
	if cfg.Database.Port != 3306 {
		t.Errorf("expected Database.Port to be 3306, got %v", cfg.Database.Port)
	}
}

// 84. Test deeply nested struct with default values
func TestDeeplyNestedStructWithDefaults(t *testing.T) {
	type DB struct {
		Host string `env:"name=DB_HOST,default=localhost"`
		Port int    `env:"name=DB_PORT,default=5432"`
	}
	type Config struct {
		Service string `env:"name=SERVICE,default=my-service"`
		DB      DB
	}

	os.Clearenv()

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Service != "my-service" {
		t.Errorf("expected Service to be 'my-service', got %v", cfg.Service)
	}
	if cfg.DB.Host != "localhost" {
		t.Errorf("expected DB.Host to be 'localhost', got %v", cfg.DB.Host)
	}
	if cfg.DB.Port != 5432 {
		t.Errorf("expected DB.Port to be 5432, got %v", cfg.DB.Port)
	}
}

// 85. Test deeply nested struct with environment variables
func TestDeeplyNestedStructWithEnvVars(t *testing.T) {
	type DB struct {
		Host string `env:"name=DB_HOST"`
		Port int    `env:"name=DB_PORT"`
	}
	type Config struct {
		Service string `env:"name=SERVICE"`
		DB      DB
	}

	os.Setenv("SERVICE", "custom-service")
	os.Setenv("DB_HOST", "db.example.com")
	os.Setenv("DB_PORT", "3306")
	defer os.Unsetenv("SERVICE")
	defer os.Unsetenv("DB_HOST")
	defer os.Unsetenv("DB_PORT")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Service != "custom-service" {
		t.Errorf("expected Service to be 'custom-service', got %v", cfg.Service)
	}
	if cfg.DB.Host != "db.example.com" {
		t.Errorf("expected DB.Host to be 'db.example.com', got %v", cfg.DB.Host)
	}
	if cfg.DB.Port != 3306 {
		t.Errorf("expected DB.Port to be 3306, got %v", cfg.DB.Port)
	}
}

// 86. Test struct with private field (should ignore private field)
func TestStructWithPrivateField(t *testing.T) {
	type Config struct {
		PublicField  string `env:"name=PUBLIC_FIELD"`
		privateField string
	}

	os.Setenv("PUBLIC_FIELD", "public_value")
	defer os.Unsetenv("PUBLIC_FIELD")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.PublicField != "public_value" {
		t.Errorf("expected PublicField to be 'public_value', got %v", cfg.PublicField)
	}

	// Private field should remain its zero value
	if cfg.privateField != "" {
		t.Errorf("expected privateField to be '', got %v", cfg.privateField)
	}
}

// 87. Test deeply nested struct with a private field (should ignore private field)
func TestDeeplyNestedStructWithPrivateField(t *testing.T) {
	type DB struct {
		Host        string `env:"name=DB_HOST"`
		privatePort int
	}
	type Config struct {
		Service string `env:"name=SERVICE"`
		DB      DB
	}

	os.Setenv("SERVICE", "my-service")
	os.Setenv("DB_HOST", "db.example.com")
	defer os.Unsetenv("SERVICE")
	defer os.Unsetenv("DB_HOST")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Service != "my-service" {
		t.Errorf("expected Service to be 'my-service', got %v", cfg.Service)
	}
	if cfg.DB.Host != "db.example.com" {
		t.Errorf("expected DB.Host to be 'db.example.com', got %v", cfg.DB.Host)
	}

	// Private field should remain zero
	if cfg.DB.privatePort != 0 {
		t.Errorf("expected privatePort to be 0, got %v", cfg.DB.privatePort)
	}
}

// 88. Test struct with private slice field (should ignore private slice)
func TestStructWithPrivateSliceField(t *testing.T) {
	type Config struct {
		PublicField  string `env:"name=PUBLIC_FIELD"`
		privateSlice []string
	}

	os.Setenv("PUBLIC_FIELD", "public_value")
	defer os.Unsetenv("PUBLIC_FIELD")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.PublicField != "public_value" {
		t.Errorf("expected PublicField to be 'public_value', got %v", cfg.PublicField)
	}

	// Private slice should remain its zero value
	if cfg.privateSlice != nil {
		t.Errorf("expected privateSlice to be nil, got %v", cfg.privateSlice)
	}
}

// 89. Test deeply nested struct with both public and private fields
func TestDeeplyNestedStructWithPublicAndPrivateFields(t *testing.T) {
	type DB struct {
		PublicHost  string `env:"name=DB_HOST"`
		privatePort int
	}
	type Config struct {
		Service string `env:"name=SERVICE"`
		DB      DB
	}

	os.Setenv("SERVICE", "my-service")
	os.Setenv("DB_HOST", "db.example.com")
	defer os.Unsetenv("SERVICE")
	defer os.Unsetenv("DB_HOST")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Service != "my-service" {
		t.Errorf("expected Service to be 'my-service', got %v", cfg.Service)
	}
	if cfg.DB.PublicHost != "db.example.com" {
		t.Errorf("expected DB.PublicHost to be 'db.example.com', got %v", cfg.DB.PublicHost)
	}

	// Private field should remain zero
	if cfg.DB.privatePort != 0 {
		t.Errorf("expected privatePort to be 0, got %v", cfg.DB.privatePort)
	}
}

// 90. Test struct with multiple nested structs
func TestStructWithMultipleNestedStructs(t *testing.T) {
	type DB struct {
		Host string `env:"name=DB_HOST"`
		Port int    `env:"name=DB_PORT"`
	}
	type Cache struct {
		Host string `env:"name=CACHE_HOST"`
		Port int    `env:"name=CACHE_PORT"`
	}
	type Config struct {
		DB    DB
		Cache Cache
	}

	os.Setenv("DB_HOST", "db.example.com")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("CACHE_HOST", "cache.example.com")
	os.Setenv("CACHE_PORT", "6379")
	defer os.Unsetenv("DB_HOST")
	defer os.Unsetenv("DB_PORT")
	defer os.Unsetenv("CACHE_HOST")
	defer os.Unsetenv("CACHE_PORT")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.DB.Host != "db.example.com" {
		t.Errorf("expected DB.Host to be 'db.example.com', got %v", cfg.DB.Host)
	}
	if cfg.DB.Port != 5432 {
		t.Errorf("expected DB.Port to be 5432, got %v", cfg.DB.Port)
	}
	if cfg.Cache.Host != "cache.example.com" {
		t.Errorf("expected Cache.Host to be 'cache.example.com', got %v", cfg.Cache.Host)
	}
	if cfg.Cache.Port != 6379 {
		t.Errorf("expected Cache.Port to be 6379, got %v", cfg.Cache.Port)
	}
}

func TestValidAwsRegion(t *testing.T) {
	type Config struct {
		Region string `env:"name=AWS_REGION,v_aws_region"`
	}

	os.Setenv("AWS_REGION", "us-east-1")
	defer os.Unsetenv("AWS_REGION")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Region != "us-east-1" {
		t.Errorf("expected Region to be 'us-east-1', got %v", cfg.Region)
	}
}

func TestInvalidAwsRegion(t *testing.T) {
	type Config struct {
		Region string `env:"name=AWS_REGION,v_aws_region"`
	}

	os.Setenv("AWS_REGION", "invalid-region")
	defer os.Unsetenv("AWS_REGION")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err == nil {
		t.Fatalf("expected an error for invalid AWS region, got none")
	}
}

func TestMissingAwsRegionWithDefault(t *testing.T) {
	type Config struct {
		Region string `env:"name=AWS_REGION,v_aws_region,default=us-west-2"`
	}

	os.Clearenv()

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Region != "us-west-2" {
		t.Errorf("expected Region to be 'us-west-2', got %v", cfg.Region)
	}
}

func TestValidAwsAccountID(t *testing.T) {
	type Config struct {
		AccountID string `env:"name=AWS_ACCOUNT_ID,v_aws_account_id"`
	}

	os.Setenv("AWS_ACCOUNT_ID", "123456789012")
	defer os.Unsetenv("AWS_ACCOUNT_ID")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.AccountID != "123456789012" {
		t.Errorf("expected AccountID to be '123456789012', got %v", cfg.AccountID)
	}
}

func TestInvalidAwsAccountID(t *testing.T) {
	type Config struct {
		AccountID string `env:"name=AWS_ACCOUNT_ID,v_aws_account_id"`
	}

	os.Setenv("AWS_ACCOUNT_ID", "invalid-account-id")
	defer os.Unsetenv("AWS_ACCOUNT_ID")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err == nil {
		t.Fatalf("expected an error for invalid AWS account ID, got none")
	}
}

func TestMissingAwsAccountIDWithDefault(t *testing.T) {
	type Config struct {
		AccountID string `env:"name=AWS_ACCOUNT_ID,v_aws_account_id,default=123456789012"`
	}

	os.Clearenv()

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.AccountID != "123456789012" {
		t.Errorf("expected AccountID to be '123456789012', got %v", cfg.AccountID)
	}
}

func TestValidAwsRoleArn(t *testing.T) {
	type Config struct {
		RoleArn string `env:"name=AWS_ROLE_ARN,v_aws_role_arn"`
	}

	os.Setenv("AWS_ROLE_ARN", "arn:aws:iam::123456789012:role/MyRole")
	defer os.Unsetenv("AWS_ROLE_ARN")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.RoleArn != "arn:aws:iam::123456789012:role/MyRole" {
		t.Errorf("expected RoleArn to be 'arn:aws:iam::123456789012:role/MyRole', got %v", cfg.RoleArn)
	}
}

func TestInvalidAwsRoleArn(t *testing.T) {
	type Config struct {
		RoleArn string `env:"name=AWS_ROLE_ARN,v_aws_role_arn"`
	}

	os.Setenv("AWS_ROLE_ARN", "invalid-arn")
	defer os.Unsetenv("AWS_ROLE_ARN")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err == nil {
		t.Fatalf("expected an error for invalid AWS role ARN, got none")
	}
}

func TestMissingAwsRoleArnWithDefault(t *testing.T) {
	type Config struct {
		RoleArn string `env:"name=AWS_ROLE_ARN,v_aws_role_arn,default=arn:aws:iam::123456789012:role/DefaultRole"`
	}

	os.Clearenv()

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.RoleArn != "arn:aws:iam::123456789012:role/DefaultRole" {
		t.Errorf("expected RoleArn to be 'arn:aws:iam::123456789012:role/DefaultRole', got %v", cfg.RoleArn)
	}
}

func TestValidAwsBucketName(t *testing.T) {
	type Config struct {
		BucketName string `env:"name=AWS_BUCKET_NAME,v_aws_bucket_name"`
	}

	os.Setenv("AWS_BUCKET_NAME", "my-valid-bucket")
	defer os.Unsetenv("AWS_BUCKET_NAME")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.BucketName != "my-valid-bucket" {
		t.Errorf("expected BucketName to be 'my-valid-bucket', got %v", cfg.BucketName)
	}
}

func TestInvalidAwsBucketName(t *testing.T) {
	type Config struct {
		BucketName string `env:"name=AWS_BUCKET_NAME,v_aws_bucket_name"`
	}

	os.Setenv("AWS_BUCKET_NAME", "Invalid_Bucket_Name")
	defer os.Unsetenv("AWS_BUCKET_NAME")

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err == nil {
		t.Fatalf("expected an error for invalid AWS bucket name, got none")
	}
}

func TestMissingAwsBucketNameWithDefault(t *testing.T) {
	type Config struct {
		BucketName string `env:"name=AWS_BUCKET_NAME,v_aws_bucket_name,default=my-default-bucket"`
	}

	os.Clearenv()

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.BucketName != "my-default-bucket" {
		t.Errorf("expected BucketName to be 'my-default-bucket', got %v", cfg.BucketName)
	}
}

func TestMultipleAwsValidators(t *testing.T) {
	type Config struct {
		BucketName string `env:"v_aws_bucket_name,v_aws_region,default=my-default-bucket"`
	}

	os.Clearenv()

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err == nil {
		t.Fatalf("expected error for multiple AWS validators, got none")
	}
}

func TestRequiredFieldWithAwsValidator(t *testing.T) {
	type Config struct {
		BucketName string `env:"required,v_aws_bucket_name"`
	}

	os.Clearenv()

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err == nil {
		t.Fatalf("expected error for empty required field with AWS validators, got none")
	}
}

func TestOptionalFieldWithAwsValidator(t *testing.T) {
	type Config struct {
		BucketName string `env:"v_aws_bucket_name"`
	}

	os.Clearenv()

	parser := env.NewParser()
	var cfg Config
	err := parser.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
