# cmux CLI Reference

cmux orchestrates multiple Claude Code sessions inside tmux with grid layouts, lifecycle hooks, and notifications.

## Quick Start

```bash
cmux spawn 4 --name backend --prompt "fix all failing tests" -- --model opus --effort high
```

## Commands

### cmux spawn <N> [options] [-- <claude-args>...]

Create a tmux session with N Claude Code panes in an optimal grid layout.

Options:
  --name <name>    Session name (default: cwd basename)
  --prompt <text>  Initial prompt sent to all panes after Claude starts
  --cwd <path>     Working directory for Claude instances
  --detach, -d     Don't attach after spawning; print attach command instead

Everything after `--` is passed verbatim to `claude`.

Examples:
  cmux spawn 4 --name backend -- --model opus --effort high
  cmux spawn 2 --prompt "review all open issues"

### cmux kill <name> | --all

Kill a named session or all cmux sessions.

Examples:
  cmux kill backend
  cmux kill --all

### cmux attach [<name>]

Attach to a cmux session. Defaults to most recent.

### cmux ls

List all cmux sessions with pane counts.

### cmux send <pane> <text> [--session <name>]

Send text to a specific pane (1-indexed). Uses tmux send-keys -l (literal mode).

Examples:
  cmux send 2 "now fix the integration tests"
  cmux send 1 "/review" --session backend

### cmux focus <pane> [--session <name>]

Focus a specific pane (1-indexed).

### cmux status [<name>]

Print pane status table (pane number, status, duration, last event).

### cmux status --tmux [--session <name>]

Print tmux-formatted status widget: [cmux: ●2 ◌1 ✓1]
- ● green = working
- ◌ yellow = waiting
- ✓ = idle

### cmux hook <event>

Hook handlers called by Claude Code (not for direct use):
  cmux hook session-start
  cmux hook prompt-submit
  cmux hook stop

### cmux hooks install|uninstall

Install/remove cmux hooks from ~/.claude/settings.json.

### cmux config [--path|--init]

Show current config, print config path, or initialize default config file.

### cmux docs [<command>]

Print this reference.

## Configuration

File: ~/.config/cmux/config.toml

```toml
[claude]
args = "--model opus --effort high"
cwd = "."

[layout]
strategy = "grid"

[notifications]
bell = true
osc9 = true
auto-focus = true
```

Override with environment variables: CMUX_ARGS, CMUX_CWD

## Environment Variables

Set in each pane before Claude launches:
- CMUX_PANE_ID — numeric pane ID
- CMUX_SESSION — session name (e.g., cmux-backend)

## Grid Layouts

| N | Layout | Description |
|---|--------|-------------|
| 1 | 1x1 | Full screen |
| 2 | 1x2 | Two columns |
| 3 | 1+2 | Left half, right split |
| 4 | 2x2 | Grid |
| 5 | 2+3 | Top 2, bottom 3 |
| 6 | 2x3 | Grid |
| 7 | 3+4 | Top 3, bottom 4 |
| 8 | 2x4 | Grid |
| 9 | 3x3 | Grid |

## State Machine

idle → working (on session-start or prompt-submit)
working → waiting (on stop)
waiting → working (on prompt-submit)
