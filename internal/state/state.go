package state

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hubermjonathan/cmux/internal/tmux"
)

type Status string

const (
	Idle    Status = "idle"
	Working Status = "working"
	Waiting Status = "waiting"
)

func Transition(from Status, event string) Status {
	switch event {
	case "session-start", "prompt-submit":
		return Working
	case "stop":
		return Waiting
	default:
		return from
	}
}

func SetPaneStatus(session string, paneID int, status Status) error {
	key := fmt.Sprintf("CMUX_PANE_%d_STATUS", paneID)
	return tmux.SetEnv(session, key, string(status))
}

func SetPaneLastEvent(session string, paneID int) error {
	key := fmt.Sprintf("CMUX_PANE_%d_LAST_EVENT", paneID)
	return tmux.SetEnv(session, key, strconv.FormatInt(time.Now().Unix(), 10))
}

func GetPaneStatus(session string, paneID int) Status {
	key := fmt.Sprintf("CMUX_PANE_%d_STATUS", paneID)
	val, err := tmux.GetEnv(session, key)
	if err != nil {
		return Idle
	}
	return Status(val)
}

func GetPaneList(session string) []int {
	val, err := tmux.GetEnv(session, "CMUX_PANES")
	if err != nil {
		return nil
	}
	return parsePaneIDs(val)
}

func GetPaneLastEvent(session string, paneID int) int64 {
	key := fmt.Sprintf("CMUX_PANE_%d_LAST_EVENT", paneID)
	val, err := tmux.GetEnv(session, key)
	if err != nil {
		return 0
	}
	ts, _ := strconv.ParseInt(val, 10, 64)
	return ts
}

func CountByStatus(session string) map[Status]int {
	return countByStatusFromEnv(session)
}

func countByStatusFromEnv(session string) map[Status]int {
	counts := map[Status]int{}
	env, err := tmux.GetAllEnv(session)
	if err != nil {
		return counts
	}
	panes := parsePaneIDs(env["CMUX_PANES"])
	for _, id := range panes {
		key := fmt.Sprintf("CMUX_PANE_%d_STATUS", id)
		s := Status(env[key])
		if s == "" {
			s = Idle
		}
		counts[s]++
	}
	return counts
}

type PaneInfo struct {
	ID        int
	Status    Status
	LastEvent int64
}

func GetAllPanes(session string) []PaneInfo {
	env, err := tmux.GetAllEnv(session)
	if err != nil {
		return nil
	}
	panes := parsePaneIDs(env["CMUX_PANES"])
	infos := make([]PaneInfo, len(panes))
	for i, id := range panes {
		s := Status(env[fmt.Sprintf("CMUX_PANE_%d_STATUS", id)])
		if s == "" {
			s = Idle
		}
		ts, _ := strconv.ParseInt(env[fmt.Sprintf("CMUX_PANE_%d_LAST_EVENT", id)], 10, 64)
		infos[i] = PaneInfo{ID: id, Status: s, LastEvent: ts}
	}
	return infos
}

func parsePaneIDs(val string) []int {
	if val == "" {
		return nil
	}
	var ids []int
	for _, s := range strings.Split(val, ",") {
		if id, err := strconv.Atoi(strings.TrimSpace(s)); err == nil {
			ids = append(ids, id)
		}
	}
	return ids
}
