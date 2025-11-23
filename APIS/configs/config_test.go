package configs

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func createTestEnvFile(t *testing.T, dir string, content string) string {
	t.Helper()
	envPath := filepath.Join(dir, ".env")
	err := os.WriteFile(envPath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to create test .env file: %v", err)
	}
	return envPath
}

func cleanupViper() {
	viper.Reset()
	os.Unsetenv("DB_DRIVER")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("WEB_PORT")
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("JWT_EXPIRATION")
}

func TestLoadDbConfig(t *testing.T) {
	tests := []struct {
		name           string
		envContent     string
		expectedConfig *confDB
	}{
		{
			name: "valid config with all fields",
			envContent: `DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=usuario_teste
DB_PASSWORD=senha_teste
DB_NAME=banco_teste`,
			expectedConfig: &confDB{
				DBDriver:   "postgres",
				DBHost:     "localhost",
				DBPort:     "5432",
				DBUser:     "usuario_teste",
				DBPassword: "senha_teste",
				DBName:     "banco_teste",
			},
		},
		{
			name: "config with special characters",
			envContent: `DB_DRIVER=mysql
DB_HOST=bd.exemplo.com
DB_PORT=3306
DB_USER=usuario@dominio
DB_PASSWORD="s3nh@!#123"
DB_NAME=meu-banco`,
			expectedConfig: &confDB{
				DBDriver:   "mysql",
				DBHost:     "bd.exemplo.com",
				DBPort:     "3306",
				DBUser:     "usuario@dominio",
				DBPassword: "s3nh@!#123",
				DBName:     "meu-banco",
			},
		},
		{
			name: "config with spaces in values",
			envContent: `DB_DRIVER=sqlite
DB_HOST=host local
DB_PORT=0
DB_USER=usuario teste
DB_PASSWORD=senha teste
DB_NAME=banco teste`,
			expectedConfig: &confDB{
				DBDriver:   "sqlite",
				DBHost:     "host local",
				DBPort:     "0",
				DBUser:     "usuario teste",
				DBPassword: "senha teste",
				DBName:     "banco teste",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanupViper()
			defer cleanupViper()

			tmpDir := t.TempDir()
			createTestEnvFile(t, tmpDir, tt.envContent)

			config, err := LoadDbConfig(tmpDir)
			if err != nil {
				t.Fatalf("LoadDbConfig() error = %v, want nil", err)
			}

			if config.DBDriver != tt.expectedConfig.DBDriver {
				t.Errorf("DBDriver = %v, want %v", config.DBDriver, tt.expectedConfig.DBDriver)
			}
			if config.DBHost != tt.expectedConfig.DBHost {
				t.Errorf("DBHost = %v, want %v", config.DBHost, tt.expectedConfig.DBHost)
			}
			if config.DBPort != tt.expectedConfig.DBPort {
				t.Errorf("DBPort = %v, want %v", config.DBPort, tt.expectedConfig.DBPort)
			}
			if config.DBUser != tt.expectedConfig.DBUser {
				t.Errorf("DBUser = %v, want %v", config.DBUser, tt.expectedConfig.DBUser)
			}
			if config.DBPassword != tt.expectedConfig.DBPassword {
				t.Errorf("DBPassword = %v, want %v", config.DBPassword, tt.expectedConfig.DBPassword)
			}
			if config.DBName != tt.expectedConfig.DBName {
				t.Errorf("DBName = %v, want %v", config.DBName, tt.expectedConfig.DBName)
			}
		})
	}
}

// TestLoadDbConfigErrors tests error handling for database configuration
func TestLoadDbConfigErrors(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func(t *testing.T) string
		expectError bool
	}{
		{
			name: "missing .env file",
			setupFunc: func(t *testing.T) string {
				return t.TempDir()
			},
			expectError: true,
		},
		{
			name: "malformed env file",
			setupFunc: func(t *testing.T) string {
				tmpDir := t.TempDir()
				createTestEnvFile(t, tmpDir, "INVALID===CONTENT\n{{{")
				return tmpDir
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanupViper()
			defer cleanupViper()

			path := tt.setupFunc(t)
			config, err := LoadDbConfig(path)

			if tt.expectError && err == nil {
				t.Errorf("LoadDbConfig() expected error, got nil")
			}
			if tt.expectError && config != nil {
				t.Errorf("LoadDbConfig() expected nil config on error, got %v", config)
			}
			if !tt.expectError && err != nil {
				t.Errorf("LoadDbConfig() unexpected error = %v", err)
			}
		})
	}
}

// TestLoadWebConfig tests loading web configuration with valid inputs
func TestLoadWebConfig(t *testing.T) {
	tests := []struct {
		name           string
		envContent     string
		expectedConfig *confWeb
	}{
		{
			name: "valid config with all fields",
			envContent: `WEB_PORT=8080
JWT_SECRET=mysecretkey
JWT_EXPIRATION=3600`,
			expectedConfig: &confWeb{
				WebServerPort: "8080",
				JWTSecret:     "mysecretkey",
				JWTExpiration: 3600,
			},
		},
		{
			name: "config with special characters in secret",
			envContent: `WEB_PORT=9000
JWT_SECRET=my$ecr3t!k3y@2024
JWT_EXPIRATION=7200`,
			expectedConfig: &confWeb{
				WebServerPort: "9000",
				JWTSecret:     "my$ecr3t!k3y@2024",
				JWTExpiration: 7200,
			},
		},
		{
			name: "config with large expiration",
			envContent: `WEB_PORT=3000
JWT_SECRET=longsecretkey123456789
JWT_EXPIRATION=86400`,
			expectedConfig: &confWeb{
				WebServerPort: "3000",
				JWTSecret:     "longsecretkey123456789",
				JWTExpiration: 86400,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanupViper()
			defer cleanupViper()

			tmpDir := t.TempDir()
			createTestEnvFile(t, tmpDir, tt.envContent)

			config, err := LoadWebConfig(tmpDir)
			if err != nil {
				t.Fatalf("LoadWebConfig() error = %v, want nil", err)
			}

			if config.WebServerPort != tt.expectedConfig.WebServerPort {
				t.Errorf("WebServerPort = %v, want %v", config.WebServerPort, tt.expectedConfig.WebServerPort)
			}
			if config.JWTSecret != tt.expectedConfig.JWTSecret {
				t.Errorf("JWTSecret = %v, want %v", config.JWTSecret, tt.expectedConfig.JWTSecret)
			}
			if config.JWTExpiration != tt.expectedConfig.JWTExpiration {
				t.Errorf("JWTExpiration = %v, want %v", config.JWTExpiration, tt.expectedConfig.JWTExpiration)
			}
		})
	}
}

// TestJWTInitialization tests that JWT authenticator is properly initialized
func TestJWTInitialization(t *testing.T) {
	tests := []struct {
		name       string
		envContent string
		jwtSecret  string
	}{
		{
			name: "simple secret",
			envContent: `WEB_PORT=8080
JWT_SECRET=simplesecret
JWT_EXPIRATION=3600`,
			jwtSecret: "simplesecret",
		},
		{
			name: "complex secret with special chars",
			envContent: `WEB_PORT=8080
JWT_SECRET="c0mpl3x!S3cr3t@2024#"
JWT_EXPIRATION=3600`,
			jwtSecret: "c0mpl3x!S3cr3t@2024#",
		},
		{
			name: "long secret",
			envContent: `WEB_PORT=8080
JWT_SECRET=verylongsecretkeyforjwtauthentication123456789
JWT_EXPIRATION=3600`,
			jwtSecret: "verylongsecretkeyforjwtauthentication123456789",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanupViper()
			defer cleanupViper()

			tmpDir := t.TempDir()
			createTestEnvFile(t, tmpDir, tt.envContent)

			config, err := LoadWebConfig(tmpDir)
			if err != nil {
				t.Fatalf("LoadWebConfig() error = %v, want nil", err)
			}

			if config.TokenAuth == nil {
				t.Fatal("TokenAuth is nil, expected non-nil JWT authenticator")
			}

			if config.JWTSecret != tt.jwtSecret {
				t.Errorf("JWTSecret = %v, want %v", config.JWTSecret, tt.jwtSecret)
			}
		})
	}
}

// TestLoadWebConfigErrors tests error handling for web configuration
func TestLoadWebConfigErrors(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func(t *testing.T) string
		expectError bool
	}{
		{
			name: "missing .env file",
			setupFunc: func(t *testing.T) string {
				return t.TempDir()
			},
			expectError: true,
		},
		{
			name: "invalid JWT_EXPIRATION (non-numeric)",
			setupFunc: func(t *testing.T) string {
				tmpDir := t.TempDir()
				createTestEnvFile(t, tmpDir, `WEB_PORT=8080
JWT_SECRET=secret
JWT_EXPIRATION=notanumber`)
				return tmpDir
			},
			expectError: true,
		},
		{
			name: "malformed env file",
			setupFunc: func(t *testing.T) string {
				tmpDir := t.TempDir()
				createTestEnvFile(t, tmpDir, "INVALID===CONTENT\n{{{")
				return tmpDir
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanupViper()
			defer cleanupViper()

			path := tt.setupFunc(t)
			config, err := LoadWebConfig(path)

			if tt.expectError && err == nil {
				t.Errorf("LoadWebConfig() expected error, got nil")
			}
			if tt.expectError && config != nil {
				t.Errorf("LoadWebConfig() expected nil config on error, got %v", config)
			}
			if !tt.expectError && err != nil {
				t.Errorf("LoadWebConfig() unexpected error = %v", err)
			}
		})
	}
}

// TestDbConfigEnvironmentOverrides tests environment variable overrides for DB config
func TestDbConfigEnvironmentOverrides(t *testing.T) {
	tests := []struct {
		name           string
		envContent     string
		envVars        map[string]string
		expectedConfig *confDB
	}{
		{
			name: "single field override",
			envContent: `DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=testuser
DB_PASSWORD=testpass
DB_NAME=testdb`,
			envVars: map[string]string{
				"DB_HOST": "overridden-host",
			},
			expectedConfig: &confDB{
				DBDriver:   "postgres",
				DBHost:     "overridden-host",
				DBPort:     "5432",
				DBUser:     "testuser",
				DBPassword: "testpass",
				DBName:     "testdb",
			},
		},
		{
			name: "multiple field overrides",
			envContent: `DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=testuser
DB_PASSWORD=testpass
DB_NAME=testdb`,
			envVars: map[string]string{
				"DB_HOST":     "prod-host",
				"DB_PORT":     "5433",
				"DB_PASSWORD": "prod-password",
			},
			expectedConfig: &confDB{
				DBDriver:   "postgres",
				DBHost:     "prod-host",
				DBPort:     "5433",
				DBUser:     "testuser",
				DBPassword: "prod-password",
				DBName:     "testdb",
			},
		},
		{
			name: "partial overrides",
			envContent: `DB_DRIVER=mysql
DB_HOST=localhost
DB_PORT=3306
DB_USER=user
DB_PASSWORD=pass
DB_NAME=db`,
			envVars: map[string]string{
				"DB_USER": "admin",
				"DB_NAME": "production",
			},
			expectedConfig: &confDB{
				DBDriver:   "mysql",
				DBHost:     "localhost",
				DBPort:     "3306",
				DBUser:     "admin",
				DBPassword: "pass",
				DBName:     "production",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanupViper()
			defer cleanupViper()

			tmpDir := t.TempDir()
			createTestEnvFile(t, tmpDir, tt.envContent)

			for key, value := range tt.envVars {
				t.Setenv(key, value)
			}

			config, err := LoadDbConfig(tmpDir)
			if err != nil {
				t.Fatalf("LoadDbConfig() error = %v, want nil", err)
			}

			if config.DBDriver != tt.expectedConfig.DBDriver {
				t.Errorf("DBDriver = %v, want %v", config.DBDriver, tt.expectedConfig.DBDriver)
			}
			if config.DBHost != tt.expectedConfig.DBHost {
				t.Errorf("DBHost = %v, want %v", config.DBHost, tt.expectedConfig.DBHost)
			}
			if config.DBPort != tt.expectedConfig.DBPort {
				t.Errorf("DBPort = %v, want %v", config.DBPort, tt.expectedConfig.DBPort)
			}
			if config.DBUser != tt.expectedConfig.DBUser {
				t.Errorf("DBUser = %v, want %v", config.DBUser, tt.expectedConfig.DBUser)
			}
			if config.DBPassword != tt.expectedConfig.DBPassword {
				t.Errorf("DBPassword = %v, want %v", config.DBPassword, tt.expectedConfig.DBPassword)
			}
			if config.DBName != tt.expectedConfig.DBName {
				t.Errorf("DBName = %v, want %v", config.DBName, tt.expectedConfig.DBName)
			}
		})
	}
}

// TestWebConfigEnvironmentOverrides tests environment variable overrides for Web config
func TestWebConfigEnvironmentOverrides(t *testing.T) {
	tests := []struct {
		name           string
		envContent     string
		envVars        map[string]string
		expectedConfig *confWeb
	}{
		{
			name: "single field override",
			envContent: `WEB_PORT=8080
JWT_SECRET=secret
JWT_EXPIRATION=3600`,
			envVars: map[string]string{
				"WEB_PORT": "9090",
			},
			expectedConfig: &confWeb{
				WebServerPort: "9090",
				JWTSecret:     "secret",
				JWTExpiration: 3600,
			},
		},
		{
			name: "multiple field overrides",
			envContent: `WEB_PORT=8080
JWT_SECRET=secret
JWT_EXPIRATION=3600`,
			envVars: map[string]string{
				"JWT_SECRET":     "prod-secret",
				"JWT_EXPIRATION": "7200",
			},
			expectedConfig: &confWeb{
				WebServerPort: "8080",
				JWTSecret:     "prod-secret",
				JWTExpiration: 7200,
			},
		},
		{
			name: "partial overrides",
			envContent: `WEB_PORT=3000
JWT_SECRET=dev-secret
JWT_EXPIRATION=1800`,
			envVars: map[string]string{
				"WEB_PORT": "4000",
			},
			expectedConfig: &confWeb{
				WebServerPort: "4000",
				JWTSecret:     "dev-secret",
				JWTExpiration: 1800,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanupViper()
			defer cleanupViper()

			tmpDir := t.TempDir()
			createTestEnvFile(t, tmpDir, tt.envContent)

			for key, value := range tt.envVars {
				t.Setenv(key, value)
			}

			config, err := LoadWebConfig(tmpDir)
			if err != nil {
				t.Fatalf("LoadWebConfig() error = %v, want nil", err)
			}

			if config.WebServerPort != tt.expectedConfig.WebServerPort {
				t.Errorf("WebServerPort = %v, want %v", config.WebServerPort, tt.expectedConfig.WebServerPort)
			}
			if config.JWTSecret != tt.expectedConfig.JWTSecret {
				t.Errorf("JWTSecret = %v, want %v", config.JWTSecret, tt.expectedConfig.JWTSecret)
			}
			if config.JWTExpiration != tt.expectedConfig.JWTExpiration {
				t.Errorf("JWTExpiration = %v, want %v", config.JWTExpiration, tt.expectedConfig.JWTExpiration)
			}
		})
	}
}

// TestGetDSN tests DSN generation for different database drivers
func TestGetDSN(t *testing.T) {
	tests := []struct {
		name        string
		envContent  string
		expectedDSN string
		expectError bool
	}{
		{
			name: "postgres DSN",
			envContent: `DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=testuser
DB_PASSWORD=testpass
DB_NAME=testdb`,
			expectedDSN: "host=localhost port=5432 user=testuser password=testpass dbname=testdb sslmode=disable",
			expectError: false,
		},
		{
			name: "mysql DSN",
			envContent: `DB_DRIVER=mysql
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=secret
DB_NAME=mydb`,
			expectedDSN: "root:secret@tcp(localhost:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local",
			expectError: false,
		},
		{
			name: "sqlite DSN",
			envContent: `DB_DRIVER=sqlite
DB_HOST=
DB_PORT=
DB_USER=
DB_PASSWORD=
DB_NAME=test.db`,
			expectedDSN: "test.db",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanupViper()
			defer cleanupViper()

			tmpDir := t.TempDir()
			createTestEnvFile(t, tmpDir, tt.envContent)

			_, err := LoadDbConfig(tmpDir)
			if err != nil {
				t.Fatalf("LoadDbConfig() error = %v", err)
			}

			dsn, err := GetDSN()
			if tt.expectError && err == nil {
				t.Errorf("GetDSN() expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("GetDSN() unexpected error = %v", err)
			}
			if !tt.expectError && dsn != tt.expectedDSN {
				t.Errorf("GetDSN() = %v, want %v", dsn, tt.expectedDSN)
			}
		})
	}
}

// TestGetDSN_WithoutLoadingConfig tests error when config not loaded
func TestGetDSN_WithoutLoadingConfig(t *testing.T) {
	cleanupViper()
	defer cleanupViper()

	// Reset dbConfig to nil
	dbConfig = nil

	dsn, err := GetDSN()
	if err == nil {
		t.Errorf("GetDSN() expected error when config not loaded, got nil")
	}
	if dsn != "" {
		t.Errorf("GetDSN() expected empty string, got %v", dsn)
	}
}

// TestNewDB_WithoutLoadingConfig tests error when config not loaded
func TestNewDB_WithoutLoadingConfig(t *testing.T) {
	cleanupViper()
	defer cleanupViper()

	// Reset dbConfig to nil
	dbConfig = nil

	db, err := NewDB()
	if err == nil {
		t.Errorf("NewDB() expected error when config not loaded, got nil")
	}
	if db != nil {
		t.Errorf("NewDB() expected nil db, got %v", db)
	}
}

// TestNewDB_UnsupportedDriver tests error for unsupported drivers
func TestNewDB_UnsupportedDriver(t *testing.T) {
	cleanupViper()
	defer cleanupViper()

	tmpDir := t.TempDir()
	createTestEnvFile(t, tmpDir, `DB_DRIVER=mongodb
DB_HOST=localhost
DB_PORT=27017
DB_USER=user
DB_PASSWORD=pass
DB_NAME=testdb`)

	_, err := LoadDbConfig(tmpDir)
	if err != nil {
		t.Fatalf("LoadDbConfig() error = %v", err)
	}

	db, err := NewDB()
	if err == nil {
		t.Errorf("NewDB() expected error for unsupported driver, got nil")
	}
	if db != nil {
		t.Errorf("NewDB() expected nil db for unsupported driver, got %v", db)
	}
}

// TestNewDB_DriversNotInstalled tests that appropriate errors are returned
// when GORM drivers are not installed (or database connection fails)
func TestNewDB_DriversNotInstalled(t *testing.T) {
	tests := []struct {
		driver      string
		expectError bool
	}{
		{"postgres", true}, // Will fail to connect (no test DB running)
		{"mysql", true},    // Driver not installed
		{"sqlite", true},   // Driver not installed
	}

	for _, tt := range tests {
		t.Run(tt.driver, func(t *testing.T) {
			cleanupViper()
			defer cleanupViper()

			tmpDir := t.TempDir()
			envContent := fmt.Sprintf(`DB_DRIVER=%s
DB_HOST=localhost
DB_PORT=5432
DB_USER=user
DB_PASSWORD=pass
DB_NAME=testdb`, tt.driver)
			createTestEnvFile(t, tmpDir, envContent)

			_, err := LoadDbConfig(tmpDir)
			if err != nil {
				t.Fatalf("LoadDbConfig() error = %v", err)
			}

			db, err := NewDB()
			if tt.expectError && err == nil {
				t.Errorf("NewDB() expected error for %s, got nil", tt.driver)
			}
			if tt.expectError && db != nil {
				t.Errorf("NewDB() expected nil db, got %v", db)
			}
		})
	}
}

// TestGetDbConfig_ReturnsLoadedConfig tests that GetDbConfig returns the loaded config
func TestGetDbConfig_ReturnsLoadedConfig(t *testing.T) {
	cleanupViper()
	defer cleanupViper()

	tmpDir := t.TempDir()
	createTestEnvFile(t, tmpDir, `DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=testuser
DB_PASSWORD=testpass
DB_NAME=testdb`)

	config, err := LoadDbConfig(tmpDir)
	if err != nil {
		t.Fatalf("LoadDbConfig() error = %v", err)
	}

	retrieved := GetDbConfig()
	if retrieved == nil {
		t.Fatal("GetDbConfig() returned nil")
	}
	if retrieved.DBDriver != config.DBDriver {
		t.Errorf("GetDbConfig().DBDriver = %v, want %v", retrieved.DBDriver, config.DBDriver)
	}
	if retrieved.DBHost != config.DBHost {
		t.Errorf("GetDbConfig().DBHost = %v, want %v", retrieved.DBHost, config.DBHost)
	}
}

// TestGetWebConfig_ReturnsLoadedConfig tests that GetWebConfig returns the loaded config
func TestGetWebConfig_ReturnsLoadedConfig(t *testing.T) {
	cleanupViper()
	defer cleanupViper()

	tmpDir := t.TempDir()
	createTestEnvFile(t, tmpDir, `WEB_PORT=8080
JWT_SECRET=secret
JWT_EXPIRATION=3600`)

	config, err := LoadWebConfig(tmpDir)
	if err != nil {
		t.Fatalf("LoadWebConfig() error = %v", err)
	}

	retrieved := GetWebConfig()
	if retrieved == nil {
		t.Fatal("GetWebConfig() returned nil")
	}
	if retrieved.WebServerPort != config.WebServerPort {
		t.Errorf("GetWebConfig().WebServerPort = %v, want %v", retrieved.WebServerPort, config.WebServerPort)
	}
	if retrieved.JWTSecret != config.JWTSecret {
		t.Errorf("GetWebConfig().JWTSecret = %v, want %v", retrieved.JWTSecret, config.JWTSecret)
	}
}

// TestGetDB_BeforeInitialization tests that GetDB returns nil before NewDB is called
func TestGetDB_BeforeInitialization(t *testing.T) {
	cleanupViper()
	defer cleanupViper()

	// Reset db to nil
	db = nil

	retrieved := GetDB()
	if retrieved != nil {
		t.Errorf("GetDB() expected nil before initialization, got %v", retrieved)
	}
}

// TestGetDB_AfterSuccessfulInitialization tests that GetDB returns the database instance
// after a successful NewDB call. This test requires a running PostgreSQL instance.
func TestGetDB_AfterSuccessfulInitialization(t *testing.T) {
	cleanupViper()
	defer cleanupViper()

	// Reset db to nil
	db = nil

	tmpDir := t.TempDir()
	// Use the actual working credentials from cmd/server/.env
	createTestEnvFile(t, tmpDir, `DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=myapp`)

	_, err := LoadDbConfig(tmpDir)
	if err != nil {
		t.Fatalf("LoadDbConfig() error = %v", err)
	}

	// Try to initialize DB
	dbInstance, err := NewDB()
	if err != nil {
		// If DB connection fails, skip this test (DB might not be running)
		t.Skipf("Skipping test: database not available: %v", err)
	}

	// GetDB should return the same instance
	retrieved := GetDB()
	if retrieved == nil {
		t.Error("GetDB() expected non-nil after successful NewDB call")
	}
	if retrieved != dbInstance {
		t.Error("GetDB() should return the same instance as NewDB()")
	}
}
