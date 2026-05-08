package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hubermjonathan/cmux/internal/state"
	"github.com/hubermjonathan/cmux/internal/tmux"
)

const sessionPrefix = "cmux-"

func mostRecentSession() string {
	sessions, _ := tmux.ListSessions(sessionPrefix)
	if len(sessions) > 0 {
		return sessions[len(sessions)-1]
	}
	return ""
}

func resolveSession(name string) (string, error) {
	session := ""
	if name != "" {
		if !strings.HasPrefix(name, sessionPrefix) {
			session = sessionPrefix + name
		} else {
			session = name
		}
	} else {
		session = mostRecentSession()
	}
	if session == "" {
		return "", fmt.Errorf("no cmux session found")
	}
	if !tmux.HasSession(session) {
		return "", fmt.Errorf("session %q not found", session)
	}
	return session, nil
}

func resolvePane(args []string) (session string, paneID int, remaining []string, err error) {
	if len(args) < 1 {
		err = fmt.Errorf("pane argument required")
		return
	}

	paneIdx, parseErr := strconv.Atoi(args[0])
	if parseErr != nil || paneIdx < 1 {
		err = fmt.Errorf("pane must be a positive integer (1-indexed)")
		return
	}

	var sessionName string
	for i := 1; i < len(args); i++ {
		if args[i] == "--session" && i+1 < len(args) {
			sessionName = args[i+1]
			i++
		} else {
			remaining = append(remaining, args[i])
		}
	}

	session, err = resolveSession(sessionName)
	if err != nil {
		return
	}

	panes := state.GetPaneList(session)
	if paneIdx > len(panes) {
		err = fmt.Errorf("pane %d does not exist (session has %d panes)", paneIdx, len(panes))
		return
	}

	paneID = panes[paneIdx-1]
	return
}

func parsePaneID(s string) int {
	s = strings.TrimPrefix(s, "%")
	id, _ := strconv.Atoi(s)
	return id
}

func shellQuote(s string) string {
	if strings.ContainsAny(s, " \t'\"\\$") {
		return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
	}
	return s
}

func cwdBasename() string {
	wd, err := os.Getwd()
	if err != nil {
		return "session"
	}
	return filepath.Base(wd)
}

func joinInts(ids []int) string {
	strs := make([]string, len(ids))
	for i, id := range ids {
		strs[i] = strconv.Itoa(id)
	}
	return strings.Join(strs, ",")
}
