package main

import (
	"fmt"

	"github.com/jon-huber/cmux/internal/tmux"
)

func runAttach(args []string) error {
	session := ""
	if len(args) > 0 {
		session = "cmux-" + args[0]
	} else {
		session = mostRecentSession()
	}
	if session == "" {
		return fmt.Errorf("no cmux session found")
	}
	if !tmux.HasSession(session) {
		return fmt.Errorf("session %q not found", session)
	}
	return tmux.Exec("attach-session", "-t", session)
}
