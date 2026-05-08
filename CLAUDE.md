# cmux

tmux-based orchestrator for multiple Claude Code sessions.

## Build

go build -o cmux .

## Run tests

go test ./...

## Project structure

- main.go: entry point and command dispatch
- cmd/: one file per CLI command
- internal/: shared logic (tmux wrapper, layout, config, state, hooks, notifications)

## Adding a new command

1. Add a case to the switch in main.go
2. Create cmd/<name>.go with func run<Name>(args []string) error
3. Use internal/tmux for tmux operations, internal/state for pane state
4. Add to docs/cmux-docs.md
5. Write tests for any new internal/ logic

## Code style

- No third-party CLI framework (hand-rolled arg parsing)
- All tmux interaction goes through internal/tmux
- Functions that can fail return error (no panics)
- Tests are table-driven where possible
