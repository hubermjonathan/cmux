package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	var err error
	switch cmd {
	case "spawn":
		err = runSpawn(args)
	case "kill":
		err = runKill(args)
	case "attach":
		err = runAttach(args)
	case "ls":
		err = runLs(args)
	case "send":
		err = runSend(args)
	case "focus":
		err = runFocus(args)
	case "status":
		err = runStatus(args)
	case "hook":
		err = runHook(args)
	case "hooks":
		err = runHooks(args)
	case "config":
		err = runConfig(args)
	case "docs":
		err = runDocs(args)
	case "help", "--help", "-h":
		printUsage()
	case "version", "--version", "-v":
		fmt.Println("cmux v0.1.0")
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`Usage: cmux <command> [args]

Commands:
  spawn <N>     Create session with N Claude Code panes in a grid
  kill <name>   Kill a named session (or --all)
  attach <name> Attach to a session
  ls            List cmux sessions
  send          Send text to a pane
  focus         Focus a pane
  status        Show pane statuses
  hook          Hook handlers (called by Claude Code)
  hooks         Install/uninstall Claude Code hooks
  config        Show/init configuration
  docs          Print CLI reference for LLMs`)
}
