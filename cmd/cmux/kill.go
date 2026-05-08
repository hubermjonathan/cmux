package main

import (
	"fmt"

	"github.com/hubermjonathan/cmux/internal/tmux"
)

func runKill(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: cmux kill <name> or cmux kill --all")
	}
	if args[0] == "--all" {
		sessions, _ := tmux.ListSessions("cmux-")
		for _, s := range sessions {
			tmux.RunSilent("kill-session", "-t", s)
		}
		fmt.Printf("killed %d session(s)\n", len(sessions))
		return nil
	}
	session := "cmux-" + args[0]
	if !tmux.HasSession(session) {
		return fmt.Errorf("session %q not found", args[0])
	}
	return tmux.RunSilent("kill-session", "-t", session)
}
