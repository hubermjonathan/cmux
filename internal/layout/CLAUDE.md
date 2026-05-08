# internal/layout

Grid layout computation for tmux pane arrangement.

## How it works

ComputeGrid(n) returns a Grid with Rows — a slice of ints where each element is the number of panes in that row.

For N=1-9, there's a hardcoded lookup table optimized for visual balance:
- N=2: two columns (not two rows)
- N=3: 1 pane left, 2 stacked right
- N=4: 2x2 grid
- N=9: 3x3 grid

For N>9, tiledFallback computes the most square-like arrangement.

## Layout semantics

The "rows" concept maps to tmux as:
- Row 0 with N panes = N horizontal splits at the top
- Row 1 with M panes = M horizontal splits at the bottom

For layouts like [1, 2] (N=3): first split the window vertically into 2 halves, then split the right half horizontally.

## Adding new layout patterns

Add a case to the switch in ComputeGrid. Add a test case to TestComputeGrid.
