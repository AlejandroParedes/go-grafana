package config

import (
	"bytes"
	"os"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNewConfig(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		cfg := NewConfig()
		if cfg.Server.Port != "8080" {
			t.Errorf("expected port 8080, got %s", cfg.Server.Port)
		}
		if cfg.Database.Host != "localhost" {
			t.Errorf("expected db host localhost, got %s", cfg.Database.Host)
		}
		if cfg.Logging.Level != "info" {
			t.Errorf("expected log level info, got %s", cfg.Logging.Level)
		}
	})

	t.Run("with env variables", func(t *testing.T) {
		os.Setenv("SERVER_PORT", "9090")
		os.Setenv("DB_HOST", "testdb")
		os.Setenv("LOG_LEVEL", "debug")
		os.Setenv("SERVER_READ_TIMEOUT", "10s")

		defer os.Unsetenv("SERVER_PORT")
		defer os.Unsetenv("DB_HOST")
		defer os.Unsetenv("LOG_LEVEL")
		defer os.Unsetenv("SERVER_READ_TIMEOUT")

		cfg := NewConfig()
		if cfg.Server.Port != "9090" {
			t.Errorf("expected port 9090, got %s", cfg.Server.Port)
		}
		if cfg.Database.Host != "testdb" {
			t.Errorf("expected db host testdb, got %s", cfg.Database.Host)
		}
		if cfg.Logging.Level != "debug" {
			t.Errorf("expected log level debug, got %s", cfg.Logging.Level)
		}
		if cfg.Server.ReadTimeout != 10*time.Second {
			t.Errorf("expected read timeout 10s, got %s", cfg.Server.ReadTimeout)
		}
	})
}

func Test_getEnv(t *testing.T) {
	t.Run("env not set", func(t *testing.T) {
		if val := getEnv("NON_EXISTENT_VAR", "default"); val != "default" {
			t.Errorf("expected default, got %s", val)
		}
	})

	t.Run("env is set", func(t *testing.T) {
		os.Setenv("EXISTENT_VAR", "value")
		defer os.Unsetenv("EXISTENT_VAR")
		if val := getEnv("EXISTENT_VAR", "default"); val != "value" {
			t.Errorf("expected value, got %s", val)
		}
	})
}

func Test_getDurationEnv(t *testing.T) {
	t.Run("env not set", func(t *testing.T) {
		if val := getDurationEnv("NON_EXISTENT_VAR", 5*time.Second); val != 5*time.Second {
			t.Errorf("expected 5s, got %s", val)
		}
	})

	t.Run("env set with duration string", func(t *testing.T) {
		os.Setenv("DURATION_VAR", "15s")
		defer os.Unsetenv("DURATION_VAR")
		if val := getDurationEnv("DURATION_VAR", 5*time.Second); val != 15*time.Second {
			t.Errorf("expected 15s, got %s", val)
		}
	})

	t.Run("env set with number string", func(t *testing.T) {
		os.Setenv("DURATION_VAR", "20")
		defer os.Unsetenv("DURATION_VAR")
		if val := getDurationEnv("DURATION_VAR", 5*time.Second); val != 20*time.Second {
			t.Errorf("expected 20s, got %s", val)
		}
	})

	t.Run("env set with invalid string", func(t *testing.T) {
		os.Setenv("DURATION_VAR", "invalid")
		defer os.Unsetenv("DURATION_VAR")
		if val := getDurationEnv("DURATION_VAR", 5*time.Second); val != 5*time.Second {
			t.Errorf("expected 5s, got %s", val)
		}
	})
}

func TestConfig_GetDSN(t *testing.T) {
	cfg := &Config{
		Database: DatabaseConfig{
			Host:     "host",
			Port:     "port",
			User:     "user",
			Password: "password",
			DBName:   "dbname",
			SSLMode:  "disable",
		},
	}
	expectedDSN := "host=host port=port user=user password=password dbname=dbname sslmode=disable"
	if dsn := cfg.GetDSN(); dsn != expectedDSN {
		t.Errorf("expected DSN '%s', got '%s'", expectedDSN, dsn)
	}
}

func TestConfig_Validate(t *testing.T) {
	cfg := &Config{}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestConfig_LogConfig(t *testing.T) {
	var buffer bytes.Buffer
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	core := zapcore.NewCore(encoder, zapcore.AddSync(&buffer), zap.InfoLevel)
	logger := zap.New(core)

	cfg := NewConfig()
	cfg.LogConfig(logger)

	logOutput := buffer.String()
	t.Log(logOutput) // to see what is being logged

	if !bytes.Contains(buffer.Bytes(), []byte(`"msg":"Configuration loaded"`)) {
		t.Error("log output should contain 'Configuration loaded'")
	}
	if !bytes.Contains(buffer.Bytes(), []byte(`"server_port":"8080"`)) {
		t.Error("log output should contain server port")
	}
	if !bytes.Contains(buffer.Bytes(), []byte(`"db_host":"localhost"`)) {
		t.Error("log output should contain db host")
	}
}
