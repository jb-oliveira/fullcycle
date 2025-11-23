package configs

import (
	"fmt"

	"github.com/go-chi/jwtauth"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbConfig  *confDB
	webConfig *confWeb
)

type confDB struct {
	DBDriver   string `mapstructure:"DB_DRIVER"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
}

type confWeb struct {
	WebServerPort string `mapstructure:"WEB_PORT"`
	JWTSecret     string `mapstructure:"JWT_SECRET"`
	JWTExpiration int    `mapstructure:"JWT_EXPIRATION"`
	TokenAuth     *jwtauth.JWTAuth
}

// LoadDbConfig loads database configuration from the specified path.
// Returns error if file is missing or cannot be parsed.
func LoadDbConfig(path string) (*confDB, error) {
	v := viper.New()
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(path)
	v.AutomaticEnv()

	err := v.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("reading db config file: %w", err)
	}

	var config confDB
	err = v.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling db config: %w", err)
	}

	dbConfig = &config
	return &config, nil
}

// LoadWebConfig loads web server configuration from the specified path.
// Returns error if file is missing or cannot be parsed.
func LoadWebConfig(path string) (*confWeb, error) {
	v := viper.New()
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(path)
	v.AutomaticEnv()

	err := v.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("reading web config file: %w", err)
	}

	var config confWeb
	err = v.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling web config: %w", err)
	}

	config.TokenAuth = jwtauth.New("HS256", []byte(config.JWTSecret), nil)
	webConfig = &config
	return &config, nil
}

// GetDbConfig returns the loaded database configuration.
func GetDbConfig() *confDB {
	return dbConfig
}

// GetWebConfig returns the loaded web configuration.
func GetWebConfig() *confWeb {
	return webConfig
}

// NewDB creates and returns a new GORM database connection using the loaded DB configuration.
// Returns error if dbConfig is not loaded or connection fails.
// Note: This function requires the appropriate GORM driver to be installed:
//   - postgres: gorm.io/driver/postgres (installed)
//   - mysql: gorm.io/driver/mysql
//   - sqlite: gorm.io/driver/sqlite
func NewDB() (*gorm.DB, error) {
	if dbConfig == nil {
		return nil, fmt.Errorf("database configuration not loaded: call LoadDbConfig first")
	}

	var dialector gorm.Dialector
	dsn := buildDSN(dbConfig)

	switch dbConfig.DBDriver {
	case "postgres", "postgresql":
		dialector = postgres.Open(dsn)
	case "mysql":
		return nil, fmt.Errorf("mysql driver not implemented: install gorm.io/driver/mysql")
	case "sqlite", "sqlite3":
		return nil, fmt.Errorf("sqlite driver not implemented: install gorm.io/driver/sqlite")
	default:
		return nil, fmt.Errorf("unsupported database driver: %s (supported: postgres, mysql, sqlite)", dbConfig.DBDriver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("opening database connection: %w", err)
	}

	return db, nil
}

// buildDSN constructs a database connection string from the configuration.
func buildDSN(config *confDB) string {
	switch config.DBDriver {
	case "postgres", "postgresql":
		return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
	case "sqlite", "sqlite3":
		return config.DBName // SQLite uses file path as DSN
	default:
		return ""
	}
}

// GetDSN returns the database connection string for the loaded configuration.
// Returns error if dbConfig is not loaded.
func GetDSN() (string, error) {
	if dbConfig == nil {
		return "", fmt.Errorf("database configuration not loaded: call LoadDbConfig first")
	}
	dsn := buildDSN(dbConfig)
	if dsn == "" {
		return "", fmt.Errorf("unsupported database driver: %s", dbConfig.DBDriver)
	}
	return dsn, nil
}
