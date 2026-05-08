package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Claude        Claude        `toml:"claude"`
	Layout        Layout        `toml:"layout"`
	Notifications Notifications `toml:"notifications"`
}

type Claude struct {
	Args string `toml:"args"`
	Cwd  string `toml:"cwd"`
}

type Layout struct {
	Strategy string `toml:"strategy"`
}

type Notifications struct {
	Bell      bool `toml:"bell"`
	OSC9      bool `toml:"osc9"`
	AutoFocus bool `toml:"auto-focus"`
}

func DefaultConfig() Config {
	return Config{
		Claude: Claude{Cwd: "."},
		Layout: Layout{Strategy: "grid"},
		Notifications: Notifications{
			Bell:      true,
			OSC9:      true,
			AutoFocus: true,
		},
	}
}

func DefaultPath() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "cmux", "config.toml")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "cmux", "config.toml")
}

func Load(path string) Config {
	cfg := DefaultConfig()

	if path == "" {
		path = DefaultPath()
	}
	if data, err := os.ReadFile(path); err == nil {
		toml.Unmarshal(data, &cfg)
	}

	if env := os.Getenv("CMUX_ARGS"); env != "" {
		cfg.Claude.Args = env
	}
	if env := os.Getenv("CMUX_CWD"); env != "" {
		cfg.Claude.Cwd = env
	}

	return cfg
}

func MergeArgs(configArgs, cliArgs string) string {
	if cliArgs == "" {
		return configArgs
	}
	if configArgs == "" {
		return cliArgs
	}
	return configArgs + " " + cliArgs
}
