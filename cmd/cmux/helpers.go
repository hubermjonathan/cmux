package main

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hubermjonathan/cmux/internal/tmux"
)

func mostRecentSession() string {
	sessions, _ := tmux.ListSessions("cmux-")
	if len(sessions) > 0 {
		return sessions[len(sessions)-1]
	}
	return ""
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
