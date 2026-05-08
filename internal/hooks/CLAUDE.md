# internal/hooks — Claude Code Hook Management

Patches ~/.claude/settings.json to install/uninstall cmux lifecycle hooks.

## Hook format in settings.json

```json
{
  "hooks": {
    "Stop": [{"hooks": [{"type": "command", "command": "cmux hook stop"}]}],
    "UserPromptSubmit": [{"hooks": [{"type": "command", "command": "cmux hook prompt-submit"}]}],
    "SessionStart": [{"hooks": [{"type": "command", "command": "cmux hook session-start"}]}]
  }
}
```

## Key behaviors

- Non-destructive merge: appends to existing hook arrays
- Identification: entries where command starts with "cmux hook"
- Idempotent: re-installing is safe (checks for existing entries)
