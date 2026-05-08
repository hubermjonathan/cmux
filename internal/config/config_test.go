package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	cfg := Load("/nonexistent/path/config.toml")
	if cfg.Notifications.Bell != true {
		t.Error("bell should default to true")
	}
	if cfg.Notifications.OSC9 != true {
		t.Error("osc9 should default to true")
	}
	if cfg.Notifications.AutoFocus != true {
		t.Error("auto-focus should default to true")
	}
	if cfg.Claude.Cwd != "." {
		t.Errorf("cwd should default to '.', got %q", cfg.Claude.Cwd)
	}
	if cfg.Layout.Strategy != "grid" {
		t.Errorf("strategy should default to 'grid', got %q", cfg.Layout.Strategy)
	}
}

func TestLoadFromFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	os.WriteFile(path, []byte(`
[claude]
args = "--model opus --effort high"
cwd = "/tmp/test"

[notifications]
bell = false
`), 0644)

	cfg := Load(path)
	if cfg.Claude.Args != "--model opus --effort high" {
		t.Errorf("got args=%q", cfg.Claude.Args)
	}
	if cfg.Claude.Cwd != "/tmp/test" {
		t.Errorf("got cwd=%q", cfg.Claude.Cwd)
	}
	if cfg.Notifications.Bell != false {
		t.Error("bell should be false from file")
	}
	if cfg.Notifications.OSC9 != true {
		t.Error("osc9 should still be true (not overridden)")
	}
}

func TestEnvOverride(t *testing.T) {
	t.Setenv("CMUX_ARGS", "--effort xhigh")
	cfg := Load("/nonexistent/path/config.toml")
	if cfg.Claude.Args != "--effort xhigh" {
		t.Errorf("env override failed: got %q", cfg.Claude.Args)
	}
}

func TestMergeArgs(t *testing.T) {
	tests := []struct {
		config, cli, want string
	}{
		{"--model opus", "--effort high", "--model opus --effort high"},
		{"--model opus", "", "--model opus"},
		{"", "--effort high", "--effort high"},
		{"", "", ""},
	}
	for _, tt := range tests {
		got := MergeArgs(tt.config, tt.cli)
		if got != tt.want {
			t.Errorf("MergeArgs(%q, %q) = %q, want %q", tt.config, tt.cli, got, tt.want)
		}
	}
}
