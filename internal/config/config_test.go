package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	t.Setenv("APP_ENV", "")
	t.Setenv("APP_PORT", "")
	t.Setenv("LOG_LEVEL", "")
	t.Setenv("DATABASE_PATH", "")
	t.Setenv("TEMPLATE_DIR", "")

	cfg := Load()

	if cfg.Environment != "development" {
		t.Errorf(
			"expected environment %q, got %q",
			"development",
			cfg.Environment,
		)
	}

	if cfg.Port != 8080 {
		t.Errorf(
			"expected port %d, got %d",
			8080,
			cfg.Port,
		)
	}

	if cfg.LogLevel != "info" {
		t.Errorf(
			"expected log level %q, got %q",
			"info",
			cfg.LogLevel,
		)
	}

	if cfg.DatabasePath != "./data/app.db" {
		t.Errorf(
			"expected database path %q, got %q",
			"./data/app.db",
			cfg.DatabasePath,
		)
	}

	if cfg.TemplateDir != "web/templates" {
		t.Errorf(
			"expected template directory %q, got %q",
			"web/templates",
			cfg.TemplateDir,
		)
	}
}

func TestLoadEnvironmentOverrides(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_PORT", "9090")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("DATABASE_PATH", "/tmp/test.db")
	t.Setenv("TEMPLATE_DIR", "/tmp/templates")

	cfg := Load()

	if cfg.Environment != "test" {
		t.Errorf(
			"expected environment %q, got %q",
			"test",
			cfg.Environment,
		)
	}

	if cfg.Port != 9090 {
		t.Errorf(
			"expected port %d, got %d",
			9090,
			cfg.Port,
		)
	}

	if cfg.LogLevel != "debug" {
		t.Errorf(
			"expected log level %q, got %q",
			"debug",
			cfg.LogLevel,
		)
	}

	if cfg.DatabasePath != "/tmp/test.db" {
		t.Errorf(
			"expected database path %q, got %q",
			"/tmp/test.db",
			cfg.DatabasePath,
		)
	}

	if cfg.TemplateDir != "/tmp/templates" {
		t.Errorf(
			"expected template directory %q, got %q",
			"/tmp/templates",
			cfg.TemplateDir,
		)
	}
}

func TestLoadInvalidPortUsesDefault(t *testing.T) {
	t.Setenv("APP_PORT", "not-a-number")

	cfg := Load()

	if cfg.Port != 8080 {
		t.Errorf(
			"expected invalid port to use default %d, got %d",
			8080,
			cfg.Port,
		)
	}
}
