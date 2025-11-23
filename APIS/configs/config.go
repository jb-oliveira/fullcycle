package configs

import (
	"fmt"

	"github.com/go-chi/jwtauth"
	"github.com/spf13/viper"
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
