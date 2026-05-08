package main

import (
	"fmt"

	"github.com/jon-huber/cmux/internal/hooks"
)

func runHooks(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: cmux hooks <install|uninstall>")
	}

	path := hooks.ClaudeSettingsPath()

	switch args[0] {
	case "install":
		if err := hooks.Install(path); err != nil {
			return fmt.Errorf("install failed: %w", err)
		}
		fmt.Println("hooks installed into", path)
	case "uninstall":
		if err := hooks.Uninstall(path); err != nil {
			return fmt.Errorf("uninstall failed: %w", err)
		}
		fmt.Println("hooks removed from", path)
	default:
		return fmt.Errorf("unknown subcommand: %s (use install or uninstall)", args[0])
	}
	return nil
}
