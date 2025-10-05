package main

import (
	"errors"
	"math"
	"sync"
	"testing"
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
			name:      "Single middle point [5, 5]",
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
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := SumLeibniz(tc.left, tc.right)

			// Check error
			if tc.wantError != nil {
				if err == nil {
					t.Errorf("SumLeibniz(%d, %d) expected error %v, got nil", tc.left, tc.right, tc.wantError)
					return
				}
				if !errors.Is(err, tc.wantError) {
					t.Errorf("SumLeibniz(%d, %d) expected error %v, got %v", tc.left, tc.right, tc.wantError, err)
				}
				return
			}

			// Check no error when not expected
			if err != nil {
				t.Errorf("SumLeibniz(%d, %d) unexpected error: %v", tc.left, tc.right, err)
				return
			}

			// Check result
			if math.Abs(result-tc.expected) > 1e-10 {
				t.Errorf("SumLeibniz(%d, %d) = %.15f, want %.15f, diff = %e",
					tc.left, tc.right, result, tc.expected, math.Abs(result-tc.expected))
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
			name:      "k=1",
			k:         1,
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
			name:      "k=100",
			k:         100,
			wantDelta: 0.01,
			wantError: nil,
		},
		{
			name:      "k=1000",
			k:         1000,
			wantDelta: 0.002,
			wantError: nil,
		},
		{
			name:      "k=10000",
			k:         10000,
			wantDelta: 0.0002,
			wantError: nil,
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

			// Check no error when not expected
			if err != nil {
				t.Errorf("LeibnizPi(k=%d) unexpected error: %v", tc.k, err)
				return
			}

			// Check delta
			delta := math.Abs(result - math.Pi)
			if delta > tc.wantDelta {
				t.Errorf("LeibnizPi(k=%d) = %.10f, delta = %.6f, want delta <= %.6f",
					tc.k, result, delta, tc.wantDelta)
			}

			t.Logf("LeibnizPi(k=%d) = %.10f | Pi = %.10f | delta = %.7e",
				tc.k, result, math.Pi, delta)
		})
	}
}

// Test cho ComputeChunk function
func TestComputeChunk(t *testing.T) {
	t.Run("Valid range", func(t *testing.T) {
		resultChan := make(chan float64, 1)
		errChan := make(chan error, 1)
		var wg sync.WaitGroup

		wg.Add(1)
		go ComputeChunk(resultChan, errChan, 0, 5, &wg) // l=0, r=5
		wg.Wait()
		close(resultChan)
		close(errChan)

		// Check no error
		select {
		case err := <-errChan:
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		default:
		}

		// Check result received
		select {
		case result := <-resultChan:
			expected, _ := SumLeibniz(0, 5)
			if math.Abs(result-expected) > 1e-10 {
				t.Errorf("Expected %f, got %f", expected, result)
			}
		default:
			t.Error("No result received")
		}
	})

	t.Run("Invalid range - should send error", func(t *testing.T) {
		resultChan := make(chan float64, 1)
		errChan := make(chan error, 1)
		var wg sync.WaitGroup

		wg.Add(1)
		go ComputeChunk(resultChan, errChan, 10, 5, &wg) // left > right
		wg.Wait()
		close(resultChan)
		close(errChan)

		// Check error received
		gotError := false
		for err := range errChan {
			if err != nil {
				gotError = true
				if !errors.Is(err, ErrInvalidRange) {
					t.Errorf("Expected ErrInvalidRange, got %v", err)
				}
			}
		}
		if !gotError {
			t.Error("Expected error but got none")
		}

		// Check no result sent when error occurs
		resultCount := 0
		for range resultChan {
			resultCount++
		}
		if resultCount > 0 {
			t.Errorf("Should not receive result when error occurs, got %d results", resultCount)
		}
	})
}
