# internal/notify — Notification Helpers

Writes bell and OSC9 sequences to pane TTYs.

## Functions

- `BellBytes()` — returns `\a`
- `OSC9Bytes(msg)` — returns `\x1b]9;msg\x1b\\`
- `SendBell(paneID)` — writes bell to pane TTY (best-effort, no error)
- `SendOSC9(paneID, msg)` — writes OSC9 to pane TTY (best-effort)

## Design decisions

- All Send* functions are best-effort (silently ignore errors)
- Bell propagates over SSH reliably
- OSC9 is best-effort (works in iTerm2, kitty, foot)
