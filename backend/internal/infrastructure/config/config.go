package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type ServerConfig struct {
	Port             string
	RepositoryDriver string
	MySQLDSN         string
	AIMode           string
	MockPaymentOK    bool
}

func LoadServerConfig() ServerConfig {
	fileConfig := loadUserConfig()
	return ServerConfig{
		Port:             pickString(fileConfig, []string{"server.port", "LISTEN_SERVER_PORT"}, "LISTEN_SERVER_PORT", "8080"),
		RepositoryDriver: pickString(fileConfig, []string{"repository.driver", "LISTEN_REPOSITORY_DRIVER"}, "LISTEN_REPOSITORY_DRIVER", "memory"),
		MySQLDSN:         pickString(fileConfig, []string{"mysql.dsn", "LISTEN_MYSQL_DSN"}, "LISTEN_MYSQL_DSN", "hwd:hWd12300-@tcp(127.0.0.1:3306)/listen?parseTime=true&loc=Local"),
		AIMode:           pickString(fileConfig, []string{"ai.mode", "LISTEN_AI_MODE"}, "LISTEN_AI_MODE", "mock"),
		MockPaymentOK:    pickBool(fileConfig, []string{"mock.enable_payment_success", "LISTEN_MOCK_ENABLE_PAYMENT_SUCCESS"}, "LISTEN_MOCK_ENABLE_PAYMENT_SUCCESS", true),
	}
}

func loadUserConfig() map[string]string {
	path := resolveUserConfigPath()
	if path == "" {
		return map[string]string{}
	}
	file, err := os.Open(path)
	if err != nil {
		return map[string]string{}
	}
	defer file.Close()

	values := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") || strings.HasPrefix(line, ";") {
			continue
		}
		sep := strings.Index(line, "=")
		if sep <= 0 {
			continue
		}
		key := strings.TrimSpace(line[:sep])
		value := strings.TrimSpace(line[sep+1:])
		value = strings.Trim(value, `"'`)
		if key != "" && value != "" {
			values[key] = value
		}
	}
	return values
}

func resolveUserConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil || homeDir == "" {
		return ""
	}
	return filepath.Join(homeDir, "conf", "listenbase.cof")
}

func pickString(fileConfig map[string]string, fileKeys []string, envKey string, fallback string) string {
	for _, key := range fileKeys {
		if value := strings.TrimSpace(fileConfig[key]); value != "" {
			return value
		}
	}
	if value := strings.TrimSpace(os.Getenv(envKey)); value != "" {
		return value
	}
	return fallback
}

func pickBool(fileConfig map[string]string, fileKeys []string, envKey string, fallback bool) bool {
	for _, key := range fileKeys {
		if value, ok := fileConfig[key]; ok {
			parsed, err := strconv.ParseBool(strings.TrimSpace(value))
			if err == nil {
				return parsed
			}
		}
	}
	if value := strings.TrimSpace(os.Getenv(envKey)); value != "" {
		parsed, err := strconv.ParseBool(value)
		if err == nil {
			return parsed
		}
	}
	return fallback
}
