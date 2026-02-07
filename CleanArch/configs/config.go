package config

import "github.com/spf13/viper"

type Config struct {
	DBDriver         string `mapstructure:"DB_DRIVER"`
	DBSource         string `mapstructure:"DB_SOURCE"`
	WebServerPort    string `mapstructure:"WEB_SERVER_PORT"`
	GrpcServerPort   string `mapstructure:"GRPC_SERVER_PORT"`
	GrapQLServerPort string `mapstructure:"GRAPQL_SERVER_PORT"`
}

func LoadConfig(path string) (*Config, error) {
	var config Config
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.SetConfigFile(".env")

	// Automatically override file values with Environment Variables if they exist
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
