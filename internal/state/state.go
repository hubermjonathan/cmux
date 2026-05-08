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
	var ids []int
	for _, s := range strings.Split(val, ",") {
		if id, err := strconv.Atoi(strings.TrimSpace(s)); err == nil {
			ids = append(ids, id)
		}
	}
	return ids
}

func CountByStatus(session string) map[Status]int {
	counts := map[Status]int{}
	for _, id := range GetPaneList(session) {
		s := GetPaneStatus(session, id)
		counts[s]++
	}
	return counts
}
