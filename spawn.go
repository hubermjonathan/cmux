package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jon-huber/cmux/internal/config"
	"github.com/jon-huber/cmux/internal/hooks"
	"github.com/jon-huber/cmux/internal/layout"
	"github.com/jon-huber/cmux/internal/tmux"
)

func runSpawn(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: cmux spawn <N> [--name <name>] [--prompt <text>] [--cwd <path>] [-- <claude-args>...]")
	}

	n, err := strconv.Atoi(args[0])
	if err != nil || n < 1 || n > 20 {
		return fmt.Errorf("N must be an integer between 1 and 20")
	}

	name, prompt, cwd, claudeArgs := parseSpawnFlags(args[1:])

	cfg := config.Load("")
	if cwd == "" {
		cwd = cfg.Claude.Cwd
	}
	mergedArgs := config.MergeArgs(cfg.Claude.Args, claudeArgs)

	sessionName := "cmux-" + name

	// Install hooks
	settingsPath := hooks.ClaudeSettingsPath()
	hooks.Install(settingsPath)

	// Create session
	if tmux.HasSession(sessionName) {
		return fmt.Errorf("session %q already exists", sessionName)
	}

	_, err = tmux.Run("new-session", "-d", "-s", sessionName, "-x", "200", "-y", "50")
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	// Build grid
	grid := layout.ComputeGrid(n)
	paneIDs, err := buildLayout(sessionName, grid)
	if err != nil {
		tmux.RunSilent("kill-session", "-t", sessionName)
		return fmt.Errorf("failed to build layout: %w", err)
	}

	// Store pane mapping
	tmux.SetEnv(sessionName, "CMUX_PANES", joinInts(paneIDs))
	for i, id := range paneIDs {
		tmux.SetEnv(sessionName, fmt.Sprintf("CMUX_PANE_INDEX_%d", i+1), strconv.Itoa(id))
		tmux.SetEnv(sessionName, fmt.Sprintf("CMUX_PANE_%d_STATUS", id), "idle")
	}

	// Launch claude in each pane
	claudeCmd := "claude"
	if mergedArgs != "" {
		claudeCmd += " " + mergedArgs
	}
	for _, id := range paneIDs {
		paneTarget := fmt.Sprintf("%%%d", id)
		if cwd != "." && cwd != "" {
			tmux.Run("send-keys", "-t", paneTarget, fmt.Sprintf("cd %s", shellQuote(cwd)), "Enter")
			time.Sleep(100 * time.Millisecond)
		}
		exportCmd := fmt.Sprintf("export CMUX_PANE_ID=%d CMUX_SESSION=%s; %s", id, sessionName, claudeCmd)
		tmux.Run("send-keys", "-t", paneTarget, exportCmd, "Enter")
	}

	// Send prompt
	if prompt != "" {
		time.Sleep(3 * time.Second)
		for _, id := range paneIDs {
			paneTarget := fmt.Sprintf("%%%d", id)
			tmux.Run("send-keys", "-t", paneTarget, "-l", prompt)
			tmux.Run("send-keys", "-t", paneTarget, "Enter")
		}
	}

	// Configure status bar
	tmux.Run("set-option", "-t", sessionName, "status-right",
		fmt.Sprintf("#(cmux status --tmux --session %s)", sessionName))
	tmux.Run("set-option", "-t", sessionName, "status-interval", "2")

	// Attach
	return tmux.Exec("attach-session", "-t", sessionName)
}

func parseSpawnFlags(args []string) (name, prompt, cwd, claudeArgs string) {
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
		case "--":
			claudeArgs = strings.Join(args[i+1:], " ")
			return
		}
	}
	return
}

func buildLayout(session string, grid layout.Grid) ([]int, error) {
	// Get the first pane ID (created with new-session)
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

	// Split to create remaining panes
	for i := 1; i < total; i++ {
		direction := "-v"
		if len(grid.Rows) == 1 || (len(grid.Rows) == 2 && i <= grid.Rows[0]) {
			direction = "-h"
		}
		target := fmt.Sprintf("%%%d", paneIDs[0])
		out, err := tmux.Run("split-window", direction, "-t", target, "-P", "-F", "#{pane_id}")
		if err != nil {
			return nil, err
		}
		paneIDs = append(paneIDs, parsePaneID(strings.TrimSpace(out)))
	}

	// Equalize layout
	tmux.Run("select-layout", "-t", session, "tiled")

	return paneIDs, nil
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
