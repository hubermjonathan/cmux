package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hubermjonathan/cmux/internal/config"
)

func runConfig(args []string) error {
	for _, arg := range args {
		switch arg {
		case "--path":
			fmt.Println(config.DefaultPath())
			return nil
		case "--init":
			return initConfig()
		}
	}
	cfg := config.Load("")
	fmt.Printf("[claude]\n  args = %q\n  cwd = %q\n\n[layout]\n  strategy = %q\n\n[notifications]\n  bell = %v\n  osc9 = %v\n  auto-focus = %v\n",
		cfg.Claude.Args, cfg.Claude.Cwd,
		cfg.Layout.Strategy,
		cfg.Notifications.Bell, cfg.Notifications.OSC9, cfg.Notifications.AutoFocus)
	return nil
}

func initConfig() error {
	path := config.DefaultPath()
	cfg := config.DefaultConfig()
	content := fmt.Sprintf(`[claude]
args = %q
cwd = %q

[layout]
strategy = %q

[notifications]
bell = %v
osc9 = %v
auto-focus = %v
`, cfg.Claude.Args, cfg.Claude.Cwd, cfg.Layout.Strategy,
		cfg.Notifications.Bell, cfg.Notifications.OSC9, cfg.Notifications.AutoFocus)

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return err
	}
	fmt.Println("config written to", path)
	return nil
}
