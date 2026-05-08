# cmux

tmux-based orchestrator for multiple Claude Code sessions.

## Build

```bash
cd ~/Code/cmux-cli
go build -o cmux .
```

Install to PATH:

```bash
cp cmux ~/.local/bin/
```

Cross-platform:

```bash
GOOS=linux GOARCH=amd64 go build -o cmux-linux .
GOOS=darwin GOARCH=arm64 go build -o cmux-darwin .
```

## Run tests

```bash
go test ./...
```

## Project structure

- `main.go` — entry point and command dispatch (switch on os.Args[1])
- `spawn.go`, `kill.go`, `attach.go`, `ls.go`, `send.go`, `focus.go` — user-facing commands
- `status.go` — pane status display
- `hook.go` — hook handlers called by Claude Code (not by users directly)
- `hooks.go` — `cmux hooks install/uninstall` command
- `config.go` — `cmux config show/init` command
- `docs.go` — `cmux docs` command (prints CLI reference for LLMs)
- `internal/tmux/` — thin wrapper around the `tmux` binary
- `internal/layout/` — grid layout computation for pane arrangement
- `internal/config/` — reads `~/.config/cmux/config.toml`
- `internal/state/` — pane state machine via tmux session env vars
- `internal/hooks/` — patches `~/.claude/settings.json` for lifecycle hooks
- `internal/notify/` — bell and OSC9 notification helpers
- `docs/cmux-docs.md` — CLI reference (kept in sync with commands)

## Adding a new command

1. Add a case to the switch in `main.go`
2. Create `<name>.go` in the root package with `func run<Name>(args []string) error`
3. Use `internal/tmux` for tmux operations, `internal/state` for pane state
4. Add to `docs/cmux-docs.md`
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
