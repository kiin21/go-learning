package main

import (
	"fmt"
	"math"
	"runtime"
	"sync"
	"time"
)

type myFunction interface{ Do() float64 }

type LeibnizPi struct{ k int }   // Parallel
type LeibnizPiV2 struct{ k int } // Sequential

// Sum cho đoạn [l, r]
func SumLeibniz(l, r int) float64 {
	sum := 0.0
	sign := 1.0
	if l&1 == 1 {
		sign = -1.0
	}
	for i := l; i <= r; i++ {
		sum += sign / float64(2*i+1)
		sign = -sign
	}
	return sum
}

func ComputeChunk(part *float64, left, right int, wg *sync.WaitGroup) {
	defer wg.Done()
	*part = SumLeibniz(left, right)
}

func (lp LeibnizPi) Do() float64 {
	k := lp.k

	n := k + 1 // i = 0..k
	workers := runtime.GOMAXPROCS(0)

	chunk := (n-1)/workers + 1

	part := make([]float64, workers)
	var wg sync.WaitGroup

	for w := 0; w < workers; w++ {
		l := w * chunk
		if l >= n {
			break
		}
		r := l + chunk - 1
		if r >= n {
			r = n - 1
		}

		wg.Add(1)
		go ComputeChunk(&part[w], l, r, &wg)
	}
	wg.Wait()

	total := 0.0
	for _, v := range part {
		total += v
	}
	return 4 * total
}

// Linear calculation
func (lp LeibnizPiV2) Do() float64 {
	sign := 1.0
	sum := 0.0
	for i := 0; i <= lp.k; i++ {
		sum += sign / float64(2*i+1)
		sign = -sign
	}
	return 4 * sum
}

func recordRuntimeWithError(f myFunction, title string) {
	fmt.Println(title)
	startTime := time.Now()
	res := f.Do()
	ms := time.Since(startTime).Milliseconds()
	fmt.Printf("Pi = %.10f | delta=%.7e | %d ms\n\n", res, math.Abs(res-math.Pi), ms)
}

func main() {
	iterations := []int{1, 1_000, 100_000, 1_000_000_000}

	for _, K := range iterations {
		recordRuntimeWithError(LeibnizPi{k: K}, fmt.Sprintf("Parallel (k=%d)", K))
		recordRuntimeWithError(LeibnizPiV2{k: K}, fmt.Sprintf("Sequential (k=%d)", K))
	}
}
