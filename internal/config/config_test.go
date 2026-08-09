package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	t.Setenv("APP_ENV", "")
	t.Setenv("APP_PORT", "")
	t.Setenv("DATABASE_PATH", "")

	cfg := Load()

	if cfg.Environment != "development" {
		t.Errorf("expected environment development, got %q", cfg.Environment)
	}

	if cfg.Port != 8080 {
		t.Errorf("expected port 8080, got %d", cfg.Port)
	}

	if cfg.DatabasePath != "./data/app.db" {
		t.Errorf("expected database path ./data/app.db, got %q", cfg.DatabasePath)
	}
}

func TestLoadEnvironmentOverrides(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_PORT", "9090")
	t.Setenv("DATABASE_PATH", "/tmp/test.db")

	cfg := Load()

	if cfg.Environment != "test" {
		t.Errorf("expected environment test, got %q", cfg.Environment)
	}

	if cfg.Port != 9090 {
		t.Errorf("expected port 9090, got %d", cfg.Port)
	}

	if cfg.DatabasePath != "/tmp/test.db" {
		t.Errorf("expected database path /tmp/test.db, got %q", cfg.DatabasePath)
	}
}

func TestLoadInvalidPortUsesDefault(t *testing.T) {
	t.Setenv("APP_PORT", "not-a-number")

	cfg := Load()

	if cfg.Port != 8080 {
		t.Errorf("expected invalid port to use default 8080, got %d", cfg.Port)
	}
}
