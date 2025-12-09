package core

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestRetryWithCondition(t *testing.T) {
	attempts := 0
	err := RetryWithCondition(
		context.Background(),
		3,
		10*time.Millisecond,
		func(ctx context.Context) error {
			attempts++
			if attempts < 3 {
				return fmt.Errorf("attempt %d", attempts)
			}
			return nil
		},
		func(err error) bool {
			return true // Always retry
		},
	)

	if err != nil {
		t.Errorf("Expected success, got error: %v", err)
	}
	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

func TestParallelWithLimit(t *testing.T) {
	results := make([]int, 0)
	mu := sync.Mutex{}

	fns := []func() error{
		func() error {
			mu.Lock()
			results = append(results, 1)
			mu.Unlock()
			return nil
		},
		func() error {
			mu.Lock()
			results = append(results, 2)
			mu.Unlock()
			return nil
		},
		func() error {
			mu.Lock()
			results = append(results, 3)
			mu.Unlock()
			return nil
		},
	}

	err := ParallelWithLimit(2, fns...)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}
}

func TestDebounce(t *testing.T) {
	count := 0
	fn := func() {
		count++
	}

	debounced := Debounce(100*time.Millisecond, fn)

	// Call multiple times quickly
	debounced()
	debounced()
	debounced()

	time.Sleep(150 * time.Millisecond)

	if count != 1 {
		t.Errorf("Expected 1 call, got %d", count)
	}
}

func TestThrottle(t *testing.T) {
	count := 0
	fn := func() {
		count++
	}

	throttled := Throttle(100*time.Millisecond, fn)

	// Call multiple times
	throttled()
	throttled()
	throttled()

	time.Sleep(50 * time.Millisecond)
	throttled()

	if count > 2 {
		t.Errorf("Expected at most 2 calls, got %d", count)
	}
}

func TestMemoize(t *testing.T) {
	callCount := 0
	fn := func(x int) int {
		callCount++
		return x * 2
	}

	memoized := Memoize(fn)

	result1 := memoized(5)
	result2 := memoized(5)

	if result1 != 10 || result2 != 10 {
		t.Error("Results should be equal")
	}
	if callCount != 1 {
		t.Errorf("Expected 1 call, got %d", callCount)
	}
}

func TestTransform(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	doubled := Transform(slice, func(x int) int {
		return x * 2
	})

	if len(doubled) != 5 {
		t.Errorf("Expected length 5, got %d", len(doubled))
	}
	if doubled[0] != 2 {
		t.Errorf("Expected 2, got %d", doubled[0])
	}
}

func TestReduce(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	sum := Reduce(slice, 0, func(acc, val int) int {
		return acc + val
	})

	if sum != 15 {
		t.Errorf("Expected 15, got %d", sum)
	}
}

func TestPartition(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5, 6}
	evens, odds := Partition(slice, func(x int) bool {
		return x%2 == 0
	})

	if len(evens) != 3 || len(odds) != 3 {
		t.Errorf("Expected 3 evens and 3 odds, got %d and %d", len(evens), len(odds))
	}
}

func TestZip(t *testing.T) {
	slice1 := []int{1, 2, 3}
	slice2 := []string{"a", "b", "c"}
	zipped := Zip(slice1, slice2)

	if len(zipped) != 3 {
		t.Errorf("Expected length 3, got %d", len(zipped))
	}
	if zipped[0].First != 1 || zipped[0].Second != "a" {
		t.Error("Zip should match elements")
	}
}

func TestFlatten(t *testing.T) {
	slices := [][]int{{1, 2}, {3, 4}, {5, 6}}
	flattened := Flatten(slices)

	if len(flattened) != 6 {
		t.Errorf("Expected length 6, got %d", len(flattened))
	}
}

func TestIntersect(t *testing.T) {
	slice1 := []int{1, 2, 3, 4}
	slice2 := []int{3, 4, 5, 6}
	intersection := Intersect(slice1, slice2)

	if len(intersection) != 2 {
		t.Errorf("Expected 2 elements, got %d", len(intersection))
	}
}

func TestDifference(t *testing.T) {
	slice1 := []int{1, 2, 3, 4}
	slice2 := []int{3, 4, 5, 6}
	difference := Difference(slice1, slice2)

	if len(difference) != 2 {
		t.Errorf("Expected 2 elements, got %d", len(difference))
	}
}

func TestUnion(t *testing.T) {
	slice1 := []int{1, 2, 3}
	slice2 := []int{3, 4, 5}
	union := Union(slice1, slice2)

	if len(union) != 5 {
		t.Errorf("Expected 5 elements, got %d", len(union))
	}
}

func TestAllMatch(t *testing.T) {
	slice := []int{2, 4, 6, 8}
	allEven := AllMatch(slice, func(x int) bool {
		return x%2 == 0
	})

	if !allEven {
		t.Error("AllMatch should be even")
	}
}

func TestAny(t *testing.T) {
	slice := []int{1, 3, 4, 5}
	hasEven := Any(slice, func(x int) bool {
		return x%2 == 0
	})

	if !hasEven {
		t.Error("Should have even number")
	}
}

func TestCount(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	count := Count(slice, func(x int) bool {
		return x%2 == 0
	})

	if count != 2 {
		t.Errorf("Expected 2, got %d", count)
	}
}

func TestFirst(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	first, found := First(slice, func(x int) bool {
		return x > 3
	})

	if !found {
		t.Error("Should find element")
	}
	if first != 4 {
		t.Errorf("Expected 4, got %d", first)
	}
}

func TestTake(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	taken := Take(slice, 3)

	if len(taken) != 3 {
		t.Errorf("Expected length 3, got %d", len(taken))
	}
}

func TestDrop(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	dropped := Drop(slice, 2)

	if len(dropped) != 3 {
		t.Errorf("Expected length 3, got %d", len(dropped))
	}
	if dropped[0] != 3 {
		t.Errorf("Expected 3, got %d", dropped[0])
	}
}

func TestReverse(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	reversed := Reverse(slice)

	if reversed[0] != 5 {
		t.Errorf("Expected 5, got %d", reversed[0])
	}
}

func TestGroupBy(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5, 6}
	groups := GroupBy(slice, func(x int) string {
		if x%2 == 0 {
			return "even"
		}
		return "odd"
	})

	if len(groups["even"]) != 3 {
		t.Errorf("Expected 3 evens, got %d", len(groups["even"]))
	}
}

func TestSum(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	sum := Sum(slice)

	if sum != 15 {
		t.Errorf("Expected 15, got %d", sum)
	}
}

func TestAverage(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	avg := Average(slice)

	if avg != 3.0 {
		t.Errorf("Expected 3.0, got %f", avg)
	}
}

func TestMinValue(t *testing.T) {
	slice := []int{5, 2, 8, 1, 9}
	min := MinValue(slice)

	if min != 1 {
		t.Errorf("Expected 1, got %d", min)
	}
}

func TestMaxValue(t *testing.T) {
	slice := []int{5, 2, 8, 1, 9}
	max := MaxValue(slice)

	if max != 9 {
		t.Errorf("Expected 9, got %d", max)
	}
}

