package task

import "github.com/spf13/viper"

func GetBaseUrl() string {
	domain := viper.GetString("config.domain")
	port := viper.GetString("config.port")
	return domain + ":" + port
}

func GetHeaders() {
	viper.GetStringSlice("config.headers")
}
