package main

import (
	"github.com/hubermjonathan/cmux/internal/tmux"
)

func runAttach(args []string) error {
	name := ""
	if len(args) > 0 {
		name = args[0]
	}
	session, err := resolveSession(name)
	if err != nil {
		return err
	}
	return tmux.Exec("attach-session", "-t", session)
}
