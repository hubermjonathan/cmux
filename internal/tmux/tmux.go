package tmux

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func Run(args ...string) (string, error) {
	cmd := exec.Command("tmux", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("tmux %s: %s (%w)", strings.Join(args, " "), strings.TrimSpace(string(out)), err)
	}
	return strings.TrimRight(string(out), "\n"), nil
}

func RunSilent(args ...string) error {
	_, err := Run(args...)
	return err
}

func Exec(args ...string) error {
	binary, err := exec.LookPath("tmux")
	if err != nil {
		return fmt.Errorf("tmux not found in PATH: %w", err)
	}
	fullArgs := append([]string{"tmux"}, args...)
	return syscall.Exec(binary, fullArgs, os.Environ())
}

func HasSession(name string) bool {
	err := RunSilent("has-session", "-t", name)
	return err == nil
}

func SetEnv(session, key, value string) error {
	return RunSilent("set-environment", "-t", session, key, value)
}

func GetEnv(session, key string) (string, error) {
	out, err := Run("show-environment", "-t", session, key)
	if err != nil {
		return "", err
	}
	parts := strings.SplitN(out, "=", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("unexpected format: %s", out)
	}
	return parts[1], nil
}

func ListSessions(prefix string) ([]string, error) {
	out, err := Run("list-sessions", "-F", "#{session_name}")
	if err != nil {
		return nil, err
	}
	if out == "" {
		return nil, nil
	}
	var sessions []string
	for _, line := range strings.Split(out, "\n") {
		if strings.HasPrefix(line, prefix) {
			sessions = append(sessions, line)
		}
	}
	return sessions, nil
}

func GetAllEnv(session string) (map[string]string, error) {
	out, err := Run("show-environment", "-t", session)
	if err != nil {
		return nil, err
	}
	env := map[string]string{}
	for _, line := range strings.Split(out, "\n") {
		if strings.HasPrefix(line, "-") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			env[parts[0]] = parts[1]
		}
	}
	return env, nil
}

func PaneTTY(paneID string) (string, error) {
	return Run("display-message", "-t", "%"+paneID, "-p", "#{pane_tty}")
}

func PaneTarget(id int) string {
	return fmt.Sprintf("%%%d", id)
}
