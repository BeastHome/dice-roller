package dice

import (
	"math"
	"testing"
)

func TestComputeStats_EmptyReturnsZero(t *testing.T) {
	s := ComputeStats(nil)
	if s != (Stats{}) {
		t.Fatalf("expected zero Stats for empty input, got %#v", s)
	}
}

func TestComputeStats_SingleValue(t *testing.T) {
	s := ComputeStats([]int{7})
	want := Stats{Count: 1, Sum: 7, Min: 7, Max: 7, Average: 7, Median: 7, StdDev: 0}
	if s != want {
		t.Fatalf("expected %#v, got %#v", want, s)
	}
}

func TestComputeStats_OddCountUsesMiddleValue(t *testing.T) {
	// sorted: [1, 3, 5] → median is 3.
	s := ComputeStats([]int{3, 5, 1})
	if s.Median != 3 {
		t.Fatalf("expected median=3 for odd count, got %v", s.Median)
	}
}

func TestComputeStats_EvenCountAveragesTwoMiddles(t *testing.T) {
	// sorted: [1, 2, 4, 8] → median is (2+4)/2 = 3.
	s := ComputeStats([]int{8, 2, 1, 4})
	if s.Median != 3 {
		t.Fatalf("expected median=3 for even count, got %v", s.Median)
	}
}

func TestComputeStats_KnownAverageAndMinMax(t *testing.T) {
	s := ComputeStats([]int{2, 4, 6, 8})
	if s.Count != 4 {
		t.Fatalf("expected count=4, got %d", s.Count)
	}
	if s.Sum != 20 {
		t.Fatalf("expected sum=20, got %d", s.Sum)
	}
	if s.Average != 5 {
		t.Fatalf("expected average=5, got %v", s.Average)
	}
	if s.Min != 2 || s.Max != 8 {
		t.Fatalf("expected min=2 max=8, got min=%d max=%d", s.Min, s.Max)
	}
}

func TestComputeStats_PopulationStdDev(t *testing.T) {
	// {2, 4, 4, 4, 5, 5, 7, 9} has population stddev = 2.0 exactly.
	s := ComputeStats([]int{2, 4, 4, 4, 5, 5, 7, 9})
	if math.Abs(s.StdDev-2.0) > 1e-9 {
		t.Fatalf("expected population stddev=2.0, got %v", s.StdDev)
	}
}

func TestComputeStats_AllIdenticalValuesHaveZeroStdDev(t *testing.T) {
	s := ComputeStats([]int{5, 5, 5, 5})
	if s.StdDev != 0 {
		t.Fatalf("expected stddev=0 for identical values, got %v", s.StdDev)
	}
}

func TestStatsFromResults_PullsTotals(t *testing.T) {
	rolls := []Result{
		{Total: 10},
		{Total: 20},
		{Total: 30},
	}
	s := StatsFromResults(rolls)
	if s.Sum != 60 || s.Average != 20 || s.Count != 3 {
		t.Fatalf("unexpected stats: %#v", s)
	}
}
