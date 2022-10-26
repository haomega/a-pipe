package task

import "github.com/spf13/viper"

type Config struct {
	Domain  string
	Port    string
	Headers []string
}

func LoadBaseConfig() *Config {
	domain := viper.GetString("config.domain")
	port := viper.GetString("config.port")
	headers := viper.GetStringSlice("config.headers")
	return &Config{
		Domain:  domain,
		Port:    port,
		Headers: headers,
	}
}
