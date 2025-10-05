package main

import (
	"fmt"
	"sync"
	"time"
)

var mu sync.Mutex
var count int

func increment() {
	mu.Lock()
	count++
	mu.Unlock()
}

func dataRace() {
	var wg sync.WaitGroup

	startTime := time.Now()

	// Data race and fix with mutex
	for i := 0; i < 1000000; i++ {
		wg.Add(1)
		go func() {
			increment()
			wg.Done()
		}()
	}

	wg.Wait()
	duration := time.Since(startTime).Milliseconds()

	fmt.Printf("Count: %d, Time executed: %d ms\n", count, duration)
}




func main() {
	dataRace()
}
