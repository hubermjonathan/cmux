# cmux

tmux-based orchestrator for multiple Claude Code sessions.

## Build

```bash
go build -o cmux ./cmd/cmux
cp cmux ~/.local/bin/   # or anywhere on your PATH
```

Cross-platform:

```bash
GOOS=linux GOARCH=amd64 go build -o cmux-linux ./cmd/cmux
GOOS=darwin GOARCH=arm64 go build -o cmux-darwin ./cmd/cmux
```

## Run tests

```bash
go test ./...
```

## Project structure

- `cmd/cmux/main.go` — entry point and command dispatch (switch on os.Args[1])
- `cmd/cmux/helpers.go` — shared utilities (mostRecentSession, parsePaneID, shellQuote, cwdBasename, joinInts)
- `cmd/cmux/spawn.go`, `kill.go`, `attach.go`, `ls.go`, `send.go`, `focus.go` — user-facing commands
- `cmd/cmux/status.go` — pane status display
- `cmd/cmux/hook.go` — hook handlers called by Claude Code (not by users directly)
- `cmd/cmux/hooks.go` — `cmux hooks install/uninstall` command
- `cmd/cmux/config.go` — `cmux config show/init` command
- `cmd/cmux/docs.go` — `cmux docs` command (prints CLI reference for LLMs)
- `cmd/cmux/docs/cmux-docs.md` — CLI reference (kept in sync with commands)
- `internal/tmux/` — thin wrapper around the `tmux` binary
- `internal/layout/` — grid layout computation for pane arrangement
- `internal/config/` — reads `~/.config/cmux/config.toml`
- `internal/state/` — pane state machine via tmux session env vars
- `internal/hooks/` — patches `~/.claude/settings.json` for lifecycle hooks
- `internal/notify/` — bell and OSC9 notification helpers

## Adding a new command

1. Add a case to the switch in `cmd/cmux/main.go`
2. Create `cmd/cmux/<name>.go` with `func run<Name>(args []string) error`
3. Use `internal/tmux` for tmux operations, `internal/state` for pane state
4. Add to `cmd/cmux/docs/cmux-docs.md`
5. Write tests for any new `internal/` logic

## How the hook system works

`cmux hook` subcommands (`stop`, `prompt-submit`, `session-start`) are called by
Claude Code's hook runner — not by users. They read `CMUX_SESSION` from the environment
to identify which session to update, then call `state.Transition` and write bell/OSC9
notifications via `internal/notify`.

`cmux hooks install/uninstall` patches `~/.claude/settings.json` to register these
handlers. See `internal/hooks/CLAUDE.md` for the settings.json format.

## Code style

- No third-party CLI framework (hand-rolled arg parsing)
- All tmux interaction goes through `internal/tmux`
- Functions that can fail return error (no panics)
- Tests are table-driven where possible
