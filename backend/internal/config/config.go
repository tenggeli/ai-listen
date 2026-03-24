package config

import "github.com/spf13/viper"

type Config struct {
	AppName  string
	Env      string
	HTTPAddr string
}

func Load() Config {
	v := viper.New()
	v.AutomaticEnv()

	v.SetDefault("APP_NAME", "listen-backend")
	v.SetDefault("APP_ENV", "dev")
	v.SetDefault("HTTP_ADDR", ":8080")

	return Config{
		AppName:  v.GetString("APP_NAME"),
		Env:      v.GetString("APP_ENV"),
		HTTPAddr: v.GetString("HTTP_ADDR"),
	}
}
