package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadServerConfig_PreferUserConfigFile(t *testing.T) {
	tempHome := t.TempDir()
	confDir := filepath.Join(tempHome, "conf")
	if err := os.MkdirAll(confDir, 0o755); err != nil {
		t.Fatalf("mkdir conf dir: %v", err)
	}

	content := `
server.port=9099
repository.driver=mysql
mysql.dsn=root:pass@tcp(127.0.0.1:3306)/listen?parseTime=true
ai.mode=real
mock.enable_payment_success=false
`
	if err := os.WriteFile(filepath.Join(confDir, "listenbase.cof"), []byte(content), 0o644); err != nil {
		t.Fatalf("write config file: %v", err)
	}

	t.Setenv("HOME", tempHome)
	t.Setenv("LISTEN_SERVER_PORT", "8000")
	t.Setenv("LISTEN_REPOSITORY_DRIVER", "memory")
	t.Setenv("LISTEN_MYSQL_DSN", "env-dsn")
	t.Setenv("LISTEN_AI_MODE", "mock")
	t.Setenv("LISTEN_MOCK_ENABLE_PAYMENT_SUCCESS", "true")

	cfg := LoadServerConfig()
	if cfg.Port != "9099" {
		t.Fatalf("expected file port, got %s", cfg.Port)
	}
	if cfg.RepositoryDriver != "mysql" {
		t.Fatalf("expected file repository driver, got %s", cfg.RepositoryDriver)
	}
	if cfg.MySQLDSN != "root:pass@tcp(127.0.0.1:3306)/listen?parseTime=true" {
		t.Fatalf("expected file dsn, got %s", cfg.MySQLDSN)
	}
	if cfg.AIMode != "real" {
		t.Fatalf("expected file ai mode, got %s", cfg.AIMode)
	}
	if cfg.MockPaymentOK {
		t.Fatalf("expected file mock payment false")
	}
}

func TestLoadServerConfig_FallbackToEnvThenDefault(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)
	t.Setenv("LISTEN_SERVER_PORT", "7000")
	t.Setenv("LISTEN_REPOSITORY_DRIVER", "mysql")
	t.Setenv("LISTEN_MYSQL_DSN", "env-dsn")
	t.Setenv("LISTEN_AI_MODE", "real")
	t.Setenv("LISTEN_MOCK_ENABLE_PAYMENT_SUCCESS", "false")

	cfg := LoadServerConfig()
	if cfg.Port != "7000" {
		t.Fatalf("expected env port, got %s", cfg.Port)
	}
	if cfg.RepositoryDriver != "mysql" {
		t.Fatalf("expected env repository driver, got %s", cfg.RepositoryDriver)
	}
	if cfg.MySQLDSN != "env-dsn" {
		t.Fatalf("expected env dsn, got %s", cfg.MySQLDSN)
	}
	if cfg.AIMode != "real" {
		t.Fatalf("expected env ai mode, got %s", cfg.AIMode)
	}
	if cfg.MockPaymentOK {
		t.Fatalf("expected env mock payment false")
	}
}
