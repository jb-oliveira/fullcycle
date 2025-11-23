# Design Document

## Overview

This design outlines the approach for creating comprehensive unit tests for the `configs` package and refactoring it to follow Go best practices. The current implementation uses `panic` for error handling and has an `init()` function that makes testing difficult. We will refactor the code to return errors properly and create table-driven tests that validate configuration loading under various scenarios.

## Architecture

The testing strategy involves:

1. **Code Refactoring**: Modify `LoadDbConfig` and `LoadWebConfig` to return errors instead of panicking
2. **Test Isolation**: Create temporary test directories and `.env` files for each test case
3. **Viper Reset**: Reset Viper's state between tests to ensure isolation
4. **Table-Driven Tests**: Use Go's standard table-driven test pattern for comprehensive coverage

### Component Changes

**Current Structure:**
```
configs/
  config.go (uses panic, has init())
```

**New Structure:**
```
configs/
  config.go (returns errors, no init())
  config_test.go (comprehensive table-driven tests)
```

## Components and Interfaces

### Refactored Configuration Functions

```go
// LoadDbConfig loads database configuration from the specified path
// Returns error if file is missing or cannot be parsed
func LoadDbConfig(path string) (*confDB, error)

// LoadWebConfig loads web server configuration from the specified path
// Returns error if file is missing or cannot be parsed
func LoadWebConfig(path string) (*confWeb, error)

// GetDbConfig returns the loaded database configuration
func GetDbConfig() *confDB

// GetWebConfig returns the loaded web configuration
func GetWebConfig() *confWeb
```

### Test Helpers

```go
// createTestEnvFile creates a temporary .env file with the given content
func createTestEnvFile(t *testing.T, dir string, content string) string

// cleanupViper resets Viper state between tests
func cleanupViper()

// setEnvVars sets environment variables for testing overrides
func setEnvVars(t *testing.T, vars map[string]string)

// unsetEnvVars cleans up environment variables after tests
func unsetEnvVars(t *testing.T, keys []string)
```

## Data Models

### Test Case Structure

```go
type dbConfigTestCase struct {
    name           string
    envContent     string
    envVars        map[string]string
    expectError    bool
    expectedConfig *confDB
}

type webConfigTestCase struct {
    name           string
    envContent     string
    envVars        map[string]string
    expectError    bool
    expectedConfig *confWeb
}
```

## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system-essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*


### Property 1: DB Config Field Population
*For any* valid .env file containing all required DB configuration fields, loading the configuration should populate all dbConfig struct fields with the exact values from the file.
**Validates: Requirements 1.1**

### Property 2: Web Config Field Population
*For any* valid .env file containing all required Web configuration fields, loading the configuration should populate all webConfig struct fields with the exact values from the file.
**Validates: Requirements 1.2**

### Property 3: JWT Authenticator Initialization
*For any* valid JWT_SECRET value in the configuration, LoadWebConfig should initialize the TokenAuth field with a non-nil JWT authenticator using the HS256 algorithm.
**Validates: Requirements 1.3**

### Property 4: Configuration Value Preservation
*For any* configuration value containing special characters, spaces, or non-ASCII characters, the loaded configuration should preserve the value exactly as specified in the file without modification.
**Validates: Requirements 1.4, 4.2**

### Property 5: DB Environment Variable Override
*For any* DB configuration field that has both a .env file value and an environment variable set, LoadDbConfig should use the environment variable value.
**Validates: Requirements 3.1**

### Property 6: Web Environment Variable Override
*For any* Web configuration field that has both a .env file value and an environment variable set, LoadWebConfig should use the environment variable value.
**Validates: Requirements 3.2**

### Property 7: Partial Override Composition
*For any* configuration with multiple fields where some have environment variable overrides and others don't, the loaded configuration should contain environment variable values for overridden fields and .env file values for non-overridden fields.
**Validates: Requirements 3.3**

### Property 8: Integer Type Parsing
*For any* valid numeric string value for JWT_EXPIRATION, LoadWebConfig should parse it as an integer type in the webConfig struct.
**Validates: Requirements 4.1**

### Property 9: Error Return Pattern for DB Config
*For any* error condition during DB configuration loading (missing file, parse error, unmarshal error), LoadDbConfig should return an error to the caller without panicking.
**Validates: Requirements 5.1**

### Property 10: Error Return Pattern for Web Config
*For any* error condition during Web configuration loading (missing file, parse error, unmarshal error), LoadWebConfig should return an error to the caller without panicking.
**Validates: Requirements 5.2**

### Property 11: Error Context Wrapping
*For any* error returned by configuration loading functions, the error message should contain contextual information about what operation failed using fmt.Errorf with %w.
**Validates: Requirements 5.4**

## Error Handling

The refactored configuration module will follow Go best practices for error handling:

1. **No Panics**: All functions return errors instead of calling `panic()`
2. **Error Wrapping**: Errors are wrapped with context using `fmt.Errorf("loading db config: %w", err)`
3. **Error Types**: Use `errors.Is` and `errors.As` for error checking in tests
4. **Graceful Degradation**: Missing optional fields should not cause errors

### Error Scenarios

- **File Not Found**: Return error with clear message about missing .env file
- **Parse Error**: Return error indicating malformed configuration data
- **Unmarshal Error**: Return error with details about which field failed to unmarshal
- **Type Mismatch**: Return error when numeric fields contain non-numeric data

## Testing Strategy

### Unit Testing Approach

We will use Go's standard `testing` package with table-driven tests. Each test function will:

1. Create a temporary directory for test files
2. Write a test `.env` file with specific content
3. Set environment variables if testing overrides
4. Call the configuration loading function
5. Assert the results match expectations
6. Clean up temporary files and environment variables

### Test Organization

```
configs/
  config_test.go
    - TestLoadDbConfig (table-driven)
    - TestLoadWebConfig (table-driven)
    - TestEnvironmentOverrides (table-driven)
    - TestErrorHandling (table-driven)
    - TestJWTInitialization (specific test)
```

### Property-Based Testing

We will use **gopter** (a Go property-based testing library) to implement property-based tests that validate the correctness properties defined above.

**Configuration:**
- Each property-based test will run a minimum of 100 iterations
- Tests will use custom generators for valid configuration data
- Edge cases (empty strings, special characters, large values) will be included in generators

**Test Tagging:**
Each property-based test will include a comment tag in this format:
```go
// Feature: config-testing, Property 1: DB Config Field Population
```

### Test Coverage Goals

- **Line Coverage**: Aim for 90%+ coverage of config.go
- **Branch Coverage**: Test all error paths and success paths
- **Edge Cases**: Empty values, special characters, missing files, malformed data
- **Integration**: Test interaction between Viper, environment variables, and file system

### Testing Tools

- **Standard Library**: `testing` package for unit tests
- **gopter**: Property-based testing library
- **testify/assert**: Optional, for cleaner assertions (standard library preferred)
- **os.TempDir**: For creating isolated test environments

### Test Isolation

Each test must be completely isolated:
- Use `t.TempDir()` for temporary directories (auto-cleanup)
- Reset Viper state with `viper.Reset()` between tests
- Use `t.Setenv()` for environment variables (auto-cleanup in Go 1.17+)
- Avoid global state dependencies

### Mock Strategy

We will NOT use mocks for this testing. Instead:
- Use real temporary files for .env files
- Use real environment variables (cleaned up after tests)
- Use real Viper library (integration testing approach)
- This ensures tests validate actual behavior, not mocked behavior

## Implementation Notes

### Refactoring Steps

1. Remove the `init()` function from config.go
2. Change `LoadDbConfig` signature to return `(*confDB, error)`
3. Change `LoadWebConfig` signature to return `(*confWeb, error)`
4. Replace all `panic(err)` calls with `return nil, fmt.Errorf("context: %w", err)`
5. Add getter functions `GetDbConfig()` and `GetWebConfig()` for accessing loaded configs
6. Update any calling code to handle returned errors

### Backward Compatibility

The refactoring will break existing code that relies on:
- The `init()` function auto-loading configuration
- The global `dbConfig` and `webConfig` variables being populated automatically

Calling code will need to be updated to:
- Explicitly call `LoadDbConfig()` and `LoadWebConfig()`
- Handle returned errors appropriately
- Use getter functions to access configuration

This is acceptable because it improves testability and follows Go best practices.
