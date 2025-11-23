package configs

import (
	"fmt"
	"log"

	"github.com/go-chi/jwtauth"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbConfig  *confDB
	webConfig *confWeb
	db        *gorm.DB
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

func LoadDbConfig(path string) (*confDB, error) {
	v := viper.New()
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(path)
	v.AutomaticEnv()

	err := v.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo de configuração do banco: %w", err)
	}

	var config confDB
	err = v.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("erro ao deserializar configuração do banco: %w", err)
	}

	dbConfig = &config
	return &config, nil
}

func LoadWebConfig(path string) (*confWeb, error) {
	v := viper.New()
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(path)
	v.AutomaticEnv()

	err := v.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo de configuração web: %w", err)
	}

	var config confWeb
	err = v.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("erro ao deserializar configuração web: %w", err)
	}

	config.TokenAuth = jwtauth.New("HS256", []byte(config.JWTSecret), nil)
	webConfig = &config
	return &config, nil
}

func GetDbConfig() *confDB {
	return dbConfig
}

func GetWebConfig() *confWeb {
	return webConfig
}

func InitGorm() error {
	if dbConfig == nil {
		return fmt.Errorf("configuração do banco não carregada: chame LoadDbConfig primeiro")
	}

	dsn := buildDSN(dbConfig)
	dialector := postgres.Open(dsn)

	var err error
	db, err = gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return fmt.Errorf("erro ao abrir conexão com banco: %w", err)
	}
	return nil
}

func GetDB() *gorm.DB {
	if db == nil {
		log.Fatalf("Banco de dados não foi inicializado!")
		return nil
	}
	return db
}

func buildDSN(config *confDB) string {
	switch config.DBDriver {
	case "postgres", "postgresql":
		return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
	case "sqlite", "sqlite3":
		return config.DBName
	default:
		return ""
	}
}

func GetDSN() (string, error) {
	if dbConfig == nil {
		return "", fmt.Errorf("configuração do banco não carregada: chame LoadDbConfig primeiro")
	}
	dsn := buildDSN(dbConfig)
	if dsn == "" {
		return "", fmt.Errorf("driver de banco não suportado: %s", dbConfig.DBDriver)
	}
	return dsn, nil
}
