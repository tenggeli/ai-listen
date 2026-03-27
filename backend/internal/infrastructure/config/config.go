package config

import "os"

type ServerConfig struct {
	Port             string
	RepositoryDriver string
	MySQLDSN         string
}

func LoadServerConfig() ServerConfig {
	port := envOrDefault("LISTEN_SERVER_PORT", "8080")
	driver := envOrDefault("LISTEN_REPOSITORY_DRIVER", "memory")
	dsn := envOrDefault("LISTEN_MYSQL_DSN", "root:root@tcp(127.0.0.1:3306)/listen?parseTime=true&loc=Local")
	return ServerConfig{Port: port, RepositoryDriver: driver, MySQLDSN: dsn}
}

func envOrDefault(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
