package main

import (
	"fmt"

	"github.com/hubermjonathan/cmux/internal/tmux"
)

func runFocus(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: cmux focus <pane> [--session <name>]")
	}

	_, paneID, _, err := resolvePane(args)
	if err != nil {
		return err
	}

	return tmux.RunSilent("select-pane", "-t", tmux.PaneTarget(paneID))
}
