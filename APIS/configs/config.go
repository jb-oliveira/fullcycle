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

func init() {
	LoadDbConfig(".")
	LoadWebConfig(".")
	fmt.Printf("DBConfig: %v", dbConfig)
	fmt.Printf("WebConfig: %v", webConfig)
}

func LoadDbConfig(path string) {
	viper.SetConfigName("db_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	// viper.AutomaticEnv() = Isso faz com que ele sobscreva as variaveis de ambiente
	// ao inves do que tiver no .env  Caso  não tenha nas variaveis de ambiente ele carrega o que ta .env
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&dbConfig)
	if err != nil {
		panic(err)
	}

}

func LoadWebConfig(path string) {
	viper.SetConfigName("web_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	// viper.AutomaticEnv() = Isso faz com que ele sobscreva as variaveis de ambiente
	// ao inves do que tiver no .env  Caso  não tenha nas variaveis de ambiente ele carrega o que ta .env
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&webConfig)
	if err != nil {
		panic(err)
	}

	webConfig.TokenAuth = jwtauth.New("HS256", []byte(webConfig.JWTSecret), nil)

}
