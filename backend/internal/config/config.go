package config

import "github.com/spf13/viper"

type Config struct {
	AppName               string
	Env                   string
	HTTPAddr              string
	MySQLDSN              string
	PaymentCallbackSecret string
}

func Load() Config {
	v := viper.New()
	v.AutomaticEnv()

	v.SetDefault("APP_NAME", "listen-backend")
	v.SetDefault("APP_ENV", "dev")
	v.SetDefault("HTTP_ADDR", ":18080")
	v.SetDefault("MYSQL_DSN", "hwd:hWd12300-@tcp(127.0.0.1:3306)/listen?charset=utf8mb4&parseTime=True&loc=Local")
	v.SetDefault("PAYMENT_CALLBACK_SECRET", "listen-dev-callback-secret")

	return Config{
		AppName:               v.GetString("APP_NAME"),
		Env:                   v.GetString("APP_ENV"),
		HTTPAddr:              v.GetString("HTTP_ADDR"),
		MySQLDSN:              v.GetString("MYSQL_DSN"),
		PaymentCallbackSecret: v.GetString("PAYMENT_CALLBACK_SECRET"),
	}
}
