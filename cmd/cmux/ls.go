package main

import (
	"fmt"

	"github.com/hubermjonathan/cmux/internal/state"
	"github.com/hubermjonathan/cmux/internal/tmux"
)

func runLs(args []string) error {
	sessions, err := tmux.ListSessions("cmux-")
	if err != nil || len(sessions) == 0 {
		fmt.Println("no cmux sessions")
		return nil
	}
	for _, s := range sessions {
		panes := state.GetPaneList(s)
		fmt.Printf("%s (%d panes)\n", s, len(panes))
	}
	return nil
}
