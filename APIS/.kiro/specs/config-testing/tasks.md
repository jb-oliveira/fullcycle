# Implementation Plan

- [x] 1. Refactor config.go to return errors instead of panicking
  - Remove the `init()` function
  - Change `LoadDbConfig` to return `(*confDB, error)`
  - Change `LoadWebConfig` to return `(*confWeb, error)`
  - Replace all `panic(err)` calls with proper error returns using `fmt.Errorf`
  - Add `GetDbConfig()` and `GetWebConfig()` getter functions
  - Ensure errors are wrapped with context
  - _Requirements: 5.1, 5.2, 5.4_

- [x] 2. Create test infrastructure and helper functions
  - [x] 2.1 Set up config_test.go with package declaration and imports
    - Import testing, os, path/filepath, and github.com/spf13/viper
    - _Requirements: 1.1, 1.2_
  
  - [x] 2.2 Implement test helper functions
    - Write `createTestEnvFile` to create temporary .env files
    - Write `cleanupViper` to reset Viper state between tests
    - _Requirements: 1.1, 1.2_

- [x] 3. Implement unit tests for DB configuration loading
  - [x] 3.1 Write table-driven test for LoadDbConfig with valid configurations
    - Test all DB fields are populated correctly
    - Test various valid values including special characters
    - _Requirements: 1.1, 1.4_
  
  - [x] 3.2 Write table-driven test for LoadDbConfig error cases
    - Test missing .env file returns error
    - Test malformed data returns error
    - Test missing required fields
    - Verify no panics occur
    - _Requirements: 2.1, 2.3, 2.4, 5.1_

- [x] 4. Implement unit tests for Web configuration loading
  - [x] 4.1 Write table-driven test for LoadWebConfig with valid configurations
    - Test all Web fields are populated correctly
    - Test JWT_EXPIRATION is parsed as integer
    - Test various valid values
    - _Requirements: 1.2, 4.1_
  
  - [x] 4.2 Write test for JWT authenticator initialization
    - Verify TokenAuth is non-nil after successful load
    - Verify HS256 algorithm is used
    - Test with various JWT_SECRET values
    - _Requirements: 1.3_
  
  - [x] 4.3 Write table-driven test for LoadWebConfig error cases
    - Test missing .env file returns error
    - Test invalid JWT_EXPIRATION (non-numeric) returns error
    - Verify no panics occur
    - _Requirements: 2.2, 4.3, 5.2_

- [x] 5. Implement tests for environment variable overrides
  - [x] 5.1 Write table-driven test for DB config environment overrides
    - Test single field override
    - Test multiple field overrides
    - Test partial overrides (some from env, some from file)
    - _Requirements: 3.1, 3.3_
  
  - [x] 5.2 Write table-driven test for Web config environment overrides
    - Test single field override
    - Test multiple field overrides
    - Test partial overrides
    - _Requirements: 3.2, 3.3_

- [ ] 6. Implement property-based tests using gopter
  - [ ] 6.1 Add gopter dependency and set up generators
    - Add gopter to go.mod
    - Create generators for valid DB config data
    - Create generators for valid Web config data
    - Configure tests to run 100+ iterations
    - _Requirements: 1.1, 1.2_
  
  - [ ]* 6.2 Write property test for DB Config Field Population
    - **Property 1: DB Config Field Population**
    - **Validates: Requirements 1.1**
    - Generate random valid DB configurations
    - Verify all fields are populated correctly
  
  - [ ]* 6.3 Write property test for Web Config Field Population
    - **Property 2: Web Config Field Population**
    - **Validates: Requirements 1.2**
    - Generate random valid Web configurations
    - Verify all fields are populated correctly
  
  - [ ]* 6.4 Write property test for JWT Authenticator Initialization
    - **Property 3: JWT Authenticator Initialization**
    - **Validates: Requirements 1.3**
    - Generate random JWT secrets
    - Verify TokenAuth is always initialized correctly
  
  - [ ]* 6.5 Write property test for Configuration Value Preservation
    - **Property 4: Configuration Value Preservation**
    - **Validates: Requirements 1.4, 4.2**
    - Generate strings with special characters, spaces, Unicode
    - Verify values are preserved exactly
  
  - [ ]* 6.6 Write property test for DB Environment Variable Override
    - **Property 5: DB Environment Variable Override**
    - **Validates: Requirements 3.1**
    - Generate random field names and values
    - Verify env vars always override file values
  
  - [ ]* 6.7 Write property test for Web Environment Variable Override
    - **Property 6: Web Environment Variable Override**
    - **Validates: Requirements 3.2**
    - Generate random field names and values
    - Verify env vars always override file values
  
  - [ ]* 6.8 Write property test for Partial Override Composition
    - **Property 7: Partial Override Composition**
    - **Validates: Requirements 3.3**
    - Generate random combinations of overridden and non-overridden fields
    - Verify correct composition of values
  
  - [ ]* 6.9 Write property test for Integer Type Parsing
    - **Property 8: Integer Type Parsing**
    - **Validates: Requirements 4.1**
    - Generate random numeric strings
    - Verify correct integer parsing
  
  - [ ]* 6.10 Write property test for Error Return Pattern (DB)
    - **Property 9: Error Return Pattern for DB Config**
    - **Validates: Requirements 5.1**
    - Generate various error conditions
    - Verify errors are returned, not panicked
  
  - [ ]* 6.11 Write property test for Error Return Pattern (Web)
    - **Property 10: Error Return Pattern for Web Config**
    - **Validates: Requirements 5.2**
    - Generate various error conditions
    - Verify errors are returned, not panicked
  
  - [ ]* 6.12 Write property test for Error Context Wrapping
    - **Property 11: Error Context Wrapping**
    - **Validates: Requirements 5.4**
    - Generate various error scenarios
    - Verify error messages contain context

- [x] 7. Update calling code to handle refactored API
  - [x] 7.1 Update cmd/server/main.go to call LoadDbConfig and LoadWebConfig
    - Add explicit calls to load functions
    - Add error handling for configuration loading
    - Use getter functions to access configuration
    - _Requirements: 5.1, 5.2_

- [x] 8. Checkpoint - Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.
