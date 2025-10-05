package main

import (
	"errors"
	"fmt"
	"math"
	"runtime"
	"sync"
	"time"
)

var threshold = 0.8
var numOfCoreCpu = runtime.GOMAXPROCS(0)
var workers = int(float64(numOfCoreCpu) * threshold)

type myFunction interface {
	Do() (float64, error)
}

type LeibnizPi struct{ k int }   // Parallel
type LeibnizPiV2 struct{ k int } // Sequential

var (
	ErrNegativeK    = errors.New("k must be non-negative")
	ErrInvalidRange = errors.New("left must be less than right")
)

// Sum cho đoạn [l, r]
func SumLeibniz(l, r int) (float64, error) {
	if l > r {
		return 0.0, ErrInvalidRange
	}
	sum := 0.0
	sign := 1.0
	if l&1 == 1 {
		sign = -1.0
	}
	for i := l; i <= r; i++ {
		sum += sign / float64(2*i+1)
		sign = -sign
	}
	return sum, nil
}

func ComputeChunk(resultChan chan<- float64, errChan chan<- error, left, right int, wg *sync.WaitGroup) {
	defer wg.Done()
	result, err := SumLeibniz(left, right)
	if err != nil {
		errChan <- err
		return
	}
	resultChan <- result
}

// ===== CÁCH 1: Parallel =====
func (lp LeibnizPi) Do() (float64, error) {
	if lp.k < 0 {
		return 0.0, ErrNegativeK
	}

	chunk := int(math.Ceil(float64(lp.k+1) / float64(workers)))

	actualWorkers := (lp.k + chunk) / chunk
	if actualWorkers > workers {
		actualWorkers = workers
	}

	resultChan := make(chan float64, actualWorkers)
	errChan := make(chan error, 1)

	var wg sync.WaitGroup

	for w := 0; w < workers; w++ {
		l := w * chunk
		if l > lp.k {
			break
		}
		r := l + chunk - 1
		if r > lp.k {
			r = lp.k
		}
		wg.Add(1)
		go ComputeChunk(resultChan, errChan, l, r, &wg)
	}

	go func() {
		wg.Wait()
		close(resultChan)
		close(errChan)
	}()

	select {
	case err := <-errChan:
		if err != nil {
			return 0.0, err
		}
	default:
	}

	total := 0.0
	for v := range resultChan {
		total += v
	}

	return 4 * total, nil
}

// ===== CÁCH 2: Sequential =====
func (lp LeibnizPiV2) Do() (float64, error) {
	if lp.k < 0 {
		return 0.0, ErrNegativeK
	}
	sign := 1.0
	sum := 0.0
	for i := 0; i <= lp.k; i++ {
		sum += sign / float64(2*i+1)
		sign = -sign
	}
	return 4 * sum, nil
}

type Result struct {
	Pi                         float64
	Delta                      float64
	ExecutionTimeInMillisecond int64
}

func recordRuntimeWithError(f myFunction) (*Result, error) {
	startTime := time.Now()

	res, err := f.Do()
	ms := time.Since(startTime).Milliseconds()

	if err != nil {
		return &Result{ExecutionTimeInMillisecond: ms}, err
	}

	result := &Result{
		Pi:                         res,
		Delta:                      math.Abs(res - math.Pi),
		ExecutionTimeInMillisecond: ms,
	}
	return result, nil
}

func main() {
	// iterations := []int{100_000_000_000}
	iterations := []int{1, 10, 100, 1_000_000}

	fmt.Println("=== Parallel Test ===")
	for _, K := range iterations {
		fmt.Printf("Parallel (k=%d)", K)
		resultPtr, err := recordRuntimeWithError(LeibnizPi{k: K})
		if err != nil {
			fmt.Println("Failed:", err)
			return
		}
		fmt.Printf("Pi=%.6f, Delta=%.6e, Time=%dms\n", resultPtr.Pi, resultPtr.Delta, resultPtr.ExecutionTimeInMillisecond)
	}

	fmt.Println("\n=== Sequential Test ===")
	for _, K := range iterations {
		fmt.Printf("Sequential (k=%d)", K)
		resultPtr, err := recordRuntimeWithError(LeibnizPiV2{k: K})
		if err != nil {
			fmt.Println("Failed:", err)
			return
		}
		fmt.Printf("Pi=%.6f, Delta=%.6e, Time=%dms\n", resultPtr.Pi, resultPtr.Delta, resultPtr.ExecutionTimeInMillisecond)
	}
}
