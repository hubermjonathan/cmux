package layout

import (
	"fmt"
	"reflect"
	"testing"
)

func TestComputeGrid(t *testing.T) {
	tests := []struct {
		n    int
		rows []int
	}{
		{1, []int{1}},
		{2, []int{1, 1}},
		{3, []int{1, 2}},
		{4, []int{2, 2}},
		{5, []int{2, 3}},
		{6, []int{3, 3}},
		{7, []int{3, 4}},
		{8, []int{4, 4}},
		{9, []int{3, 3, 3}},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("N=%d", tt.n), func(t *testing.T) {
			got := ComputeGrid(tt.n)
			if !reflect.DeepEqual(got.Rows, tt.rows) {
				t.Errorf("ComputeGrid(%d).Rows = %v, want %v", tt.n, got.Rows, tt.rows)
			}
			if got.Total() != tt.n {
				t.Errorf("ComputeGrid(%d).Total() = %d, want %d", tt.n, got.Total(), tt.n)
			}
		})
	}
}

func TestComputeGridZero(t *testing.T) {
	got := ComputeGrid(0)
	if len(got.Rows) != 0 {
		t.Errorf("ComputeGrid(0).Rows should be empty, got %v", got.Rows)
	}
}

func TestTiledFallback(t *testing.T) {
	got := ComputeGrid(12)
	total := got.Total()
	if total != 12 {
		t.Errorf("ComputeGrid(12).Total() = %d, want 12", total)
	}
}
