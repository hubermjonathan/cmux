package hooks

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestInstallIntoEmptySettings(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.json")
	os.WriteFile(path, []byte(`{}`), 0644)

	err := Install(path)
	if err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(path)
	var result map[string]any
	json.Unmarshal(data, &result)

	hooks, _ := result["hooks"].(map[string]any)
	if hooks == nil {
		t.Fatal("hooks should exist")
	}

	stop, _ := hooks["Stop"].([]any)
	if len(stop) != 1 {
		t.Errorf("expected 1 Stop hook, got %d", len(stop))
	}

	submit, _ := hooks["UserPromptSubmit"].([]any)
	if len(submit) != 1 {
		t.Errorf("expected 1 UserPromptSubmit hook, got %d", len(submit))
	}

	start, _ := hooks["SessionStart"].([]any)
	if len(start) != 1 {
		t.Errorf("expected 1 SessionStart hook, got %d", len(start))
	}
}

func TestInstallPreservesExisting(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.json")
	os.WriteFile(path, []byte(`{"hooks":{"Stop":[{"hooks":[{"type":"command","command":"echo hi"}]}]}}`), 0644)

	Install(path)

	data, _ := os.ReadFile(path)
	var result map[string]any
	json.Unmarshal(data, &result)

	hooks := result["hooks"].(map[string]any)
	stop := hooks["Stop"].([]any)
	if len(stop) != 2 {
		t.Errorf("expected 2 Stop hooks (existing + cmux), got %d", len(stop))
	}
}

func TestInstallIdempotent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.json")
	os.WriteFile(path, []byte(`{}`), 0644)

	Install(path)
	Install(path)

	data, _ := os.ReadFile(path)
	var result map[string]any
	json.Unmarshal(data, &result)

	hooks := result["hooks"].(map[string]any)
	stop := hooks["Stop"].([]any)
	if len(stop) != 1 {
		t.Errorf("expected 1 Stop hook after double install, got %d", len(stop))
	}
}

func TestUninstall(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.json")
	os.WriteFile(path, []byte(`{"hooks":{"Stop":[{"hooks":[{"type":"command","command":"cmux hook stop"}]},{"hooks":[{"type":"command","command":"echo hi"}]}]}}`), 0644)

	Uninstall(path)

	data, _ := os.ReadFile(path)
	var result map[string]any
	json.Unmarshal(data, &result)

	hooks := result["hooks"].(map[string]any)
	stop := hooks["Stop"].([]any)
	if len(stop) != 1 {
		t.Errorf("expected 1 Stop hook after uninstall, got %d", len(stop))
	}
}

func TestInstallCreatesFileIfMissing(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.json")

	err := Install(path)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("settings.json should have been created")
	}
}
