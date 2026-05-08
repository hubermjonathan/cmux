package main

import (
	"fmt"
	"strings"

	"github.com/hubermjonathan/cmux/internal/tmux"
)

func runSend(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: cmux send <pane> <text> [--session <name>]")
	}

	_, paneID, remaining, err := resolvePane(args)
	if err != nil {
		return err
	}

	text := strings.Join(remaining, " ")
	target := tmux.PaneTarget(paneID)
	if _, err := tmux.Run("send-keys", "-t", target, "-l", text); err != nil {
		return err
	}
	tmux.RunSilent("send-keys", "-t", target, "Enter")
	return nil
}
