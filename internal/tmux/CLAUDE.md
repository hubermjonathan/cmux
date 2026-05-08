# internal/tmux

Thin wrapper around the `tmux` binary. All tmux interaction in the project goes through this package.

## Key functions

- Run(args...): execute tmux command, return stdout
- RunSilent(args...): execute, discard output, return error
- Exec(args...): syscall.Exec into tmux (replaces process, used for attach)
- HasSession(name): check if session exists
- SetEnv/GetEnv: read/write tmux session environment variables
- ListSessions(prefix): list sessions matching prefix
- PaneTTY(paneID): get the TTY path for a pane

## Patterns

- Pane IDs are numeric strings (e.g. "42"). The % prefix is added by PaneTTY internally.
- Errors include the tmux stderr output for debugging.
- RunSilent is for commands where we don't need the output.
- Exec is only used for `cmux attach` (replaces the process with tmux attach).
