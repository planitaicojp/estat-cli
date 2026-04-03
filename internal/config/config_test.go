package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_noFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ESTAT_CONFIG_DIR", dir)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if cfg.Format != DefaultFormat {
		t.Errorf("Format = %q, want %q", cfg.Format, DefaultFormat)
	}
	if cfg.Lang != DefaultLang {
		t.Errorf("Lang = %q, want %q", cfg.Lang, DefaultLang)
	}
}

func TestLoad_withFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ESTAT_CONFIG_DIR", dir)

	content := []byte("app_id: test-id\nformat: json\nlang: E\n")
	if err := os.WriteFile(filepath.Join(dir, "config.yaml"), content, 0600); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if cfg.AppID != "test-id" {
		t.Errorf("AppID = %q, want %q", cfg.AppID, "test-id")
	}
	if cfg.Format != "json" {
		t.Errorf("Format = %q, want %q", cfg.Format, "json")
	}
	if cfg.Lang != "E" {
		t.Errorf("Lang = %q, want %q", cfg.Lang, "E")
	}
}

func TestSave(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ESTAT_CONFIG_DIR", dir)

	cfg := &Config{AppID: "save-test", Format: "csv", Lang: "J"}
	if err := cfg.Save(); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if loaded.AppID != "save-test" {
		t.Errorf("AppID = %q, want %q", loaded.AppID, "save-test")
	}
}

func TestEnvOr(t *testing.T) {
	t.Setenv("TEST_VAR", "hello")
	if got := EnvOr("TEST_VAR", "fallback"); got != "hello" {
		t.Errorf("EnvOr = %q, want %q", got, "hello")
	}
	if got := EnvOr("NONEXISTENT_VAR", "fallback"); got != "fallback" {
		t.Errorf("EnvOr = %q, want %q", got, "fallback")
	}
}
