# cmux

A tmux-based orchestrator for running multiple Claude Code sessions in a grid layout.

## Install

```bash
go install github.com/jon-huber/cmux@latest
```

Or build from source:

```bash
cd ~/Code/cmux-cli
go build -o cmux .
cp cmux ~/.local/bin/   # or anywhere on your PATH
```

## Quick Start

```bash
# Install hooks into Claude Code
cmux hooks install

# Spawn 4 Claude Code sessions in a 2x2 grid
cmux spawn 4 --name myproject -- --model opus --effort high

# Check status
cmux status myproject

# Send text to pane 2
cmux send 2 "fix the failing tests"

# Kill the session
cmux kill myproject
```

## Configuration

Create `~/.config/cmux/config.toml`:

```toml
[claude]
args = "--model opus --effort high --permission-mode plan"
cwd = "."

[notifications]
bell = true
osc9 = true
auto-focus = true
```

## For LLMs

Run `cmux docs` to get the full CLI reference.
