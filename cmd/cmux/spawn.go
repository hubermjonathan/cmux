package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hubermjonathan/cmux/internal/config"
	"github.com/hubermjonathan/cmux/internal/hooks"
	"github.com/hubermjonathan/cmux/internal/layout"
	"github.com/hubermjonathan/cmux/internal/state"
	"github.com/hubermjonathan/cmux/internal/tmux"
)

func runSpawn(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: cmux spawn <N> [--name <name>] [--prompt <text>] [--cwd <path>] [--detach] [-- <claude-args>...]")
	}

	n, err := strconv.Atoi(args[0])
	if err != nil || n < 1 || n > 20 {
		return fmt.Errorf("N must be an integer between 1 and 20")
	}

	name, prompt, cwd, claudeArgs, detach := parseSpawnFlags(args[1:])

	cfg := config.Load("")
	if cwd == "" {
		cwd = cfg.Claude.Cwd
	}
	mergedArgs := config.MergeArgs(cfg.Claude.Args, claudeArgs)

	sessionName := sessionPrefix + name

	settingsPath := hooks.ClaudeSettingsPath()
	if !hooks.IsInstalled(settingsPath) {
		hooks.Install(settingsPath)
	}

	if tmux.HasSession(sessionName) {
		return fmt.Errorf("session %q already exists", sessionName)
	}

	_, err = tmux.Run("new-session", "-d", "-s", sessionName, "-x", "200", "-y", "50")
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	grid := layout.ComputeGrid(n)
	paneIDs, err := buildLayout(sessionName, grid)
	if err != nil {
		tmux.RunSilent("kill-session", "-t", sessionName)
		return fmt.Errorf("failed to build layout: %w", err)
	}

	tmux.SetEnv(sessionName, "CMUX_PANES", joinInts(paneIDs))
	for _, id := range paneIDs {
		state.SetPaneStatus(sessionName, id, state.Idle)
	}

	claudeCmd := "claude"
	if mergedArgs != "" {
		claudeCmd += " " + mergedArgs
	}
	for _, id := range paneIDs {
		target := tmux.PaneTarget(id)
		var cmd string
		if cwd != "." && cwd != "" {
			cmd = fmt.Sprintf("cd %s && export CMUX_PANE_ID=%d CMUX_SESSION=%s; %s", shellQuote(cwd), id, sessionName, claudeCmd)
		} else {
			cmd = fmt.Sprintf("export CMUX_PANE_ID=%d CMUX_SESSION=%s; %s", id, sessionName, claudeCmd)
		}
		tmux.Run("send-keys", "-t", target, cmd, "Enter")
	}

	if prompt != "" {
		time.Sleep(3 * time.Second)
		for _, id := range paneIDs {
			target := tmux.PaneTarget(id)
			tmux.Run("send-keys", "-t", target, "-l", prompt)
			tmux.Run("send-keys", "-t", target, "Enter")
		}
	}

	tmux.Run("set-option", "-t", sessionName, "status-right",
		fmt.Sprintf("#(cmux status --tmux --session %s)", sessionName))
	tmux.Run("set-option", "-t", sessionName, "status-interval", "2")

	if detach {
		fmt.Printf("session %q spawned (%d panes)\nattach with: cmux attach %s\n", sessionName, n, name)
		return nil
	}
	return tmux.Exec("attach-session", "-t", sessionName)
}

func parseSpawnFlags(args []string) (name, prompt, cwd, claudeArgs string, detach bool) {
	name = cwdBasename()
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--name":
			i++
			if i < len(args) {
				name = args[i]
			}
		case "--prompt":
			i++
			if i < len(args) {
				prompt = args[i]
			}
		case "--cwd":
			i++
			if i < len(args) {
				cwd = args[i]
			}
		case "--detach", "-d":
			detach = true
		case "--":
			claudeArgs = strings.Join(args[i+1:], " ")
			return
		}
	}
	return
}

func buildLayout(session string, grid layout.Grid) ([]int, error) {
	out, err := tmux.Run("list-panes", "-t", session, "-F", "#{pane_id}")
	if err != nil {
		return nil, err
	}
	firstID := parsePaneID(strings.TrimSpace(out))
	paneIDs := []int{firstID}

	total := grid.Total()
	if total <= 1 {
		return paneIDs, nil
	}

	for i := 1; i < total; i++ {
		direction := "-v"
		if len(grid.Rows) == 1 || (len(grid.Rows) == 2 && i <= grid.Rows[0]) {
			direction = "-h"
		}
		target := tmux.PaneTarget(paneIDs[0])
		out, err := tmux.Run("split-window", direction, "-t", target, "-P", "-F", "#{pane_id}")
		if err != nil {
			return nil, err
		}
		paneIDs = append(paneIDs, parsePaneID(strings.TrimSpace(out)))
	}

	tmux.Run("select-layout", "-t", session, "tiled")

	return paneIDs, nil
}
