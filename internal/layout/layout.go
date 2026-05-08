package layout

type Grid struct {
	Rows []int // number of panes in each row
}

func ComputeGrid(n int) Grid {
	if n <= 0 {
		return Grid{Rows: []int{}}
	}

	switch n {
	case 1:
		return Grid{Rows: []int{1}}
	case 2:
		return Grid{Rows: []int{1, 1}}
	case 3:
		return Grid{Rows: []int{1, 2}}
	case 4:
		return Grid{Rows: []int{2, 2}}
	case 5:
		return Grid{Rows: []int{2, 3}}
	case 6:
		return Grid{Rows: []int{3, 3}}
	case 7:
		return Grid{Rows: []int{3, 4}}
	case 8:
		return Grid{Rows: []int{4, 4}}
	case 9:
		return Grid{Rows: []int{3, 3, 3}}
	default:
		return Grid{Rows: tiledFallback(n)}
	}
}

func tiledFallback(n int) []int {
	cols := intSqrt(n)
	rows := (n + cols - 1) / cols
	result := make([]int, rows)
	remaining := n
	for i := range result {
		perRow := remaining / (rows - i)
		result[i] = perRow
		remaining -= perRow
	}
	return result
}

func intSqrt(n int) int {
	for i := 1; ; i++ {
		if i*i >= n {
			return i
		}
	}
}

func (g Grid) Total() int {
	total := 0
	for _, r := range g.Rows {
		total += r
	}
	return total
}
