package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/hubermjonathan/cmux/internal/config"
	"github.com/hubermjonathan/cmux/internal/notify"
	"github.com/hubermjonathan/cmux/internal/state"
	"github.com/hubermjonathan/cmux/internal/tmux"
)

func runHook(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: cmux hook <session-start|prompt-submit|stop>")
	}

	event := args[0]
	paneIDStr := os.Getenv("CMUX_PANE_ID")
	session := os.Getenv("CMUX_SESSION")

	if paneIDStr == "" || session == "" {
		return nil
	}

	paneID, err := strconv.Atoi(paneIDStr)
	if err != nil {
		return nil
	}

	if !tmux.HasSession(session) {
		return nil
	}

	current := state.GetPaneStatus(session, paneID)
	newStatus := state.Transition(current, event)

	state.SetPaneStatus(session, paneID, newStatus)
	state.SetPaneLastEvent(session, paneID)

	if event == "stop" {
		cfg := config.Load("")
		if cfg.Notifications.Bell {
			notify.SendBell(paneID)
		}
		if cfg.Notifications.OSC9 {
			notify.SendOSC9(paneID, fmt.Sprintf("Claude waiting: pane %s [%s]", paneIDStr, session))
		}
		if cfg.Notifications.AutoFocus {
			autoFocusIfSingleWaiter(session, paneID)
		}
		tmux.RunSilent("refresh-client", "-S")
	}

	return nil
}

func autoFocusIfSingleWaiter(session string, paneID int) {
	counts := state.CountByStatus(session)
	if counts[state.Waiting] == 1 {
		tmux.RunSilent("select-pane", "-t", fmt.Sprintf("%%%d", paneID))
	}
}
