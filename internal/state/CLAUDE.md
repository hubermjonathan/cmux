# internal/state — Pane State Machine

Reads and writes pane status via tmux session environment variables.

## Key types

- `Status` — one of `idle`, `working`, `waiting`
- `Transition(from, event)` — pure state machine function

## tmux env var format

- `CMUX_PANES` — comma-separated pane IDs
- `CMUX_PANE_{id}_STATUS` — current status
- `CMUX_PANE_{id}_LAST_EVENT` — unix timestamp

## Adding a new state

1. Add the constant to the Status type
2. Update Transition() with valid transitions
3. Update CountByStatus in tests
