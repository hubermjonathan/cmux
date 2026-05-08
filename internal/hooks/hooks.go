package hooks

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

var cmuxHooks = map[string]string{
	"Stop":             "cmux hook stop",
	"UserPromptSubmit": "cmux hook prompt-submit",
	"SessionStart":     "cmux hook session-start",
}

func ClaudeSettingsPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".claude", "settings.json")
}

func Install(settingsPath string) error {
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		if os.IsNotExist(err) {
			data = []byte("{}")
		} else {
			return err
		}
	}

	var settings map[string]any
	if err := json.Unmarshal(data, &settings); err != nil {
		return err
	}

	hooks, _ := settings["hooks"].(map[string]any)
	if hooks == nil {
		hooks = map[string]any{}
	}

	for event, command := range cmuxHooks {
		entry := map[string]any{
			"hooks": []any{
				map[string]any{"type": "command", "command": command},
			},
		}

		existing, _ := hooks[event].([]any)
		if !containsCmuxHook(existing, command) {
			hooks[event] = append(existing, entry)
		}
	}

	settings["hooks"] = hooks
	out, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(settingsPath, out, 0644)
}

func Uninstall(settingsPath string) error {
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		return err
	}

	var settings map[string]any
	if err := json.Unmarshal(data, &settings); err != nil {
		return err
	}

	hooks, _ := settings["hooks"].(map[string]any)
	if hooks == nil {
		return nil
	}

	for event := range cmuxHooks {
		entries, _ := hooks[event].([]any)
		var filtered []any
		for _, entry := range entries {
			if !isCmuxEntry(entry) {
				filtered = append(filtered, entry)
			}
		}
		if filtered == nil {
			filtered = []any{}
		}
		hooks[event] = filtered
	}

	settings["hooks"] = hooks
	out, _ := json.MarshalIndent(settings, "", "  ")
	return os.WriteFile(settingsPath, out, 0644)
}

func IsInstalled(settingsPath string) bool {
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		return false
	}
	var settings map[string]any
	if err := json.Unmarshal(data, &settings); err != nil {
		return false
	}
	hooks, _ := settings["hooks"].(map[string]any)
	if hooks == nil {
		return false
	}
	for event, command := range cmuxHooks {
		entries, _ := hooks[event].([]any)
		if !containsCmuxHook(entries, command) {
			return false
		}
	}
	return true
}

func isCmuxEntry(entry any) bool {
	m, _ := entry.(map[string]any)
	if m == nil {
		return false
	}
	hooksList, _ := m["hooks"].([]any)
	for _, h := range hooksList {
		hm, _ := h.(map[string]any)
		if hm == nil {
			continue
		}
		cmd, _ := hm["command"].(string)
		if strings.HasPrefix(cmd, "cmux hook") {
			return true
		}
	}
	return false
}

func containsCmuxHook(entries []any, command string) bool {
	for _, entry := range entries {
		m, _ := entry.(map[string]any)
		if m == nil {
			continue
		}
		hooksList, _ := m["hooks"].([]any)
		for _, h := range hooksList {
			hm, _ := h.(map[string]any)
			if hm == nil {
				continue
			}
			if hm["command"] == command {
				return true
			}
		}
	}
	return false
}
