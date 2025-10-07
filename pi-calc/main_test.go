package main

import (
	"context"
	"errors"
	"math"
	"sync"
	"testing"
	"time"
)

// Test cho hàm SumLeibniz
func TestSumLeibniz(t *testing.T) {
	tests := []struct {
		name      string
		left      int
		right     int
		expected  float64
		wantError error
	}{
		{
			name:      "Single point [0, 0]",
			left:      0,
			right:     0,
			expected:  1.0,
			wantError: nil,
		},
		{
			name:      "Range [0, 2]",
			left:      0,
			right:     2,
			expected:  1.0 - 1.0/3.0 + 1.0/5.0,
			wantError: nil,
		},
		{
			name:      "Single point [5, 5]",
			left:      5,
			right:     5,
			expected:  -1.0 / 11.0,
			wantError: nil,
		},
		{
			name:      "Invalid range left > right",
			left:      5,
			right:     3,
			expected:  0.0,
			wantError: ErrInvalidRange,
		},
		{
			name:      "Range too large",
			left:      0,
			right:     2_000_000_000,
			expected:  0.0,
			wantError: ErrRangeTooLarge,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := SumLeibniz(tc.left, tc.right)

			// Check for error
			if tc.wantError != nil {
				if !errors.Is(err, tc.wantError) {
					t.Errorf("SumLeibniz(%d, %d) expected error %v, got %v", tc.left, tc.right, tc.wantError, err)
				}
				return
			}

			// Check for having an error when not expected
			if err != nil {
				t.Errorf("SumLeibniz(%d, %d) unexpected error: %v", tc.left, tc.right, err)
				return
			}

			// Check result
			if math.Abs(result-tc.expected) > 1e-10 {
				t.Errorf("SumLeibniz(%d, %d) = %.15f, want %.15f, diff = %e", tc.left, tc.right, result, tc.expected, math.Abs(result-tc.expected))
			}
		})
	}
}

// Test cho LeibnizPi (Parallel)
func TestLeibnizPi(t *testing.T) {
	tests := []struct {
		name      string
		k         int
		wantDelta float64
		wantError error
	}{
		{
			name:      "k=-1",
			k:         -1,
			wantDelta: 0,
			wantError: ErrNegativeK,
		},
		{
			name:      "k=0",
			k:         0,
			wantDelta: 1.0,
			wantError: nil,
		},
		{
			name:      "k=10",
			k:         10,
			wantDelta: 0.1,
			wantError: nil,
		},
		{
			name:      "k=1_000",
			k:         1_000,
			wantDelta: 0.001,
			wantError: nil,
		},
		{
			name:      "k=1_000_000",
			k:         1_000_000,
			wantDelta: 0.000001,
			wantError: nil,
		},
		{
			name:      "k too large causes error in parallel chunks",
			k:         10_000_000_000,
			wantDelta: 0,
			wantError: ErrRangeTooLarge,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lp := LeibnizPi{k: tc.k}
			result, err := lp.Do()

			// Check error
			if tc.wantError != nil {
				if !errors.Is(err, tc.wantError) {
					t.Errorf("LeibnizPi(k=%d) expected error %v, got %v", tc.k, tc.wantError, err)
				}
				return
			}

			// Check for having an error when not expected
			if err != nil {
				t.Errorf("LeibnizPi(k=%d) unexpected error: %v", tc.k, err)
				return
			}

			// Check delta
			delta := math.Abs(result - math.Pi)
			if delta > tc.wantDelta {
				t.Errorf("LeibnizPi(k=%d) = %.10f, delta = %.6f, want delta <= %.6f", tc.k, result, delta, tc.wantDelta)
			}

			t.Logf("LeibnizPi(k=%d) = %.10f | Pi = %.10f | delta = %.7e", tc.k, result, math.Pi, delta)
		})
	}
}

// Test cho LeibnizPiV2 (Sequential)
func TestLeibnizPiV2(t *testing.T) {
	tests := []struct {
		name      string
		k         int
		wantDelta float64
		wantError error
	}{
		{
			name:      "k=-1 (negative)",
			k:         -1,
			wantDelta: 0,
			wantError: ErrNegativeK,
		},
		{
			name:      "k=0",
			k:         0,
			wantDelta: 1.0,
			wantError: nil,
		},
		{
			name:      "k=10",
			k:         10,
			wantDelta: 0.1,
			wantError: nil,
		},
		{
			name:      "k=1_000",
			k:         1_000,
			wantDelta: 0.001,
			wantError: nil,
		},
		{
			name:      "k=1_000_000",
			k:         1_000_000,
			wantDelta: 0.000001,
			wantError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lp := LeibnizPiV2{k: tc.k}
			result, err := lp.Do()

			// Check error
			if tc.wantError != nil {
				if !errors.Is(err, tc.wantError) {
					t.Errorf("LeibnizPi(k=%d) expected error %v, got %v", tc.k, tc.wantError, err)
				}
				return
			}

			// Check no error when not expected
			if err != nil {
				t.Errorf("LeibnizPi(k=%d) unexpected error: %v", tc.k, err)
				return
			}

			// Check delta
			delta := math.Abs(result - math.Pi)
			if delta > tc.wantDelta {
				t.Errorf("LeibnizPi(k=%d) = %.10f, delta = %.6f, want delta <= %.6f", tc.k, result, delta, tc.wantDelta)
			}

			t.Logf("LeibnizPi(k=%d) = %.10f | Pi = %.10f | delta = %.7e", tc.k, result, math.Pi, delta)
		})
	}
}

// Test recordRuntime
func TestRecordRuntime(t *testing.T) {
	tests := []struct {
		name      string
		impl      myFunction
		wantError error
	}{
		{
			name:      "LeibnizPi with k=100",
			impl:      LeibnizPi{k: 100},
			wantError: nil,
		},
		{
			name:      "LeibnizPi with negative k",
			impl:      LeibnizPi{k: -1},
			wantError: ErrNegativeK,
		},
		{
			name:      "LeibnizPiV2 with k=100",
			impl:      LeibnizPiV2{k: 100},
			wantError: nil,
		},
		{
			name:      "LeibnizPiV2 with negative k",
			impl:      LeibnizPiV2{k: -1},
			wantError: ErrNegativeK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := recordRuntime(tc.impl)

			if tc.wantError != nil {
				if !errors.Is(tc.wantError, err) {
					t.Errorf("recordRuntime()  expected error %v, got %v", tc.wantError, err)
				}
				if result == nil {
					t.Error("recordRuntime() should return Result even on error")
				}
				return
			}

			if err != nil {
				t.Errorf("recordRuntime() unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Error("recordRuntime() returned nil result")
				return
			}

			t.Logf("Result: Pi=%.10f, Delta=%.10e, Time=%dms", result.Pi, result.Delta, result.ExecutionTimeInMillisecond)
		})
	}
}

// Test cho ComputeChunk với context cancellation
func TestComputeChunk(t *testing.T) {
	t.Run("Normal execution without cancellation", func(t *testing.T) {
		ctx := context.Background()
		resultChan := make(chan float64, 1)
		errChan := make(chan error, 1)
		var wg sync.WaitGroup

		wg.Add(1)
		go ComputeChunk(ctx, resultChan, errChan, 0, 2, &wg)
		wg.Wait()
		close(resultChan)
		close(errChan)

		// Check for errors
		select {
		case err := <-errChan:
			if err != nil {
				t.Errorf("ComputeChunk unexpected error: %v", err)
			}
		default:
		}

		// Check result
		result := <-resultChan
		expected := 1.0 - 1.0/3.0 + 1.0/5.0
		if math.Abs(result-expected) > 1e-10 {
			t.Errorf("ComputeChunk result = %.10f, want %.10f", result, expected)
		}
	})

	t.Run("Context cancelled before execution", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		resultChan := make(chan float64, 1)
		errChan := make(chan error, 1)
		var wg sync.WaitGroup

		wg.Add(1)
		go ComputeChunk(ctx, resultChan, errChan, 0, 2, &wg)
		wg.Wait()
		close(resultChan)
		close(errChan)

		// When context is cancelled before execution, the function returns early
		// and may not send a result. The select with default allows it to continue.
		// So we just verify no error occurred
		select {
		case err := <-errChan:
			if err != nil {
				t.Errorf("ComputeChunk unexpected error: %v", err)
			}
		default:
		}
	})

	t.Run("Invalid range error", func(t *testing.T) {
		ctx := context.Background()
		resultChan := make(chan float64, 1)
		errChan := make(chan error, 1)
		var wg sync.WaitGroup

		wg.Add(1)
		go ComputeChunk(ctx, resultChan, errChan, 5, 3, &wg) // Invalid range
		wg.Wait()
		close(resultChan)
		close(errChan)

		// Check for error
		err := <-errChan
		if err == nil {
			t.Error("ComputeChunk expected error for invalid range, got nil")
		}
		if !errors.Is(err, ErrInvalidRange) {
			t.Errorf("ComputeChunk expected ErrInvalidRange, got %v", err)
		}
	})

	t.Run("Context cancelled during result send", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		resultChan := make(chan float64)
		errChan := make(chan error, 1)
		var wg sync.WaitGroup

		wg.Add(1)
		go ComputeChunk(ctx, resultChan, errChan, 0, 100, &wg)

		// Give goroutine time to compute and reach the select for sending
		// Then cancel context before we read from resultChan
		time.Sleep(1 * time.Millisecond)
		cancel()

		wg.Wait()
		close(resultChan)
		close(errChan)
	})
}
