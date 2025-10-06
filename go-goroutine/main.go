package main

import (
	"context"
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

// Đọc từ nhiều channel
func selectExample() {
	chann1 := make(chan string, 5)
	chann2 := make(chan string, 5)
	// Gửi vào channel 1
	go func() {
		for i := 1; i <= 5; i++ {
			time.Sleep(100 * time.Millisecond)
			chann1 <- fmt.Sprintf("Channel1: %d", i)
		}
		close(chann1)
	}()
	// Gửi vào channel 2
	go func() {
		for i := 1; i <= 5; i++ {
			time.Sleep(150 * time.Millisecond)
			chann2 <- fmt.Sprintf("Channel2: %d", i)
		}
		close(chann2)
	}()
	// Đọc từ cả 2 channel
	openChannels := 2
	for openChannels > 0 {
		select {
		case msg, ok := <-chann1:
			if !ok {
				chann1 = nil
				openChannels--
			} else {
				fmt.Println("Nhận từ", msg)
			}
		case msg, ok := <-chann2:
			if !ok {
				chann2 = nil
				openChannels--
			} else {
				fmt.Println("Nhận từ", msg)
			}
		}
	}
}

// ============= CONTEXT EXAMPLES =============

// Ví dụ 1: Context với Timeout
func contextWithTimeout() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		select {
		case <-time.After(5 * time.Second):
			fmt.Println("Công việc hoàn thành")
		case <-ctx.Done():
			fmt.Println("Context timeout:", ctx.Err())
		}
	}()

	fmt.Println("Main: Context đã kết thúc")

	wg.Wait()
}

// Ví dụ 2: Context với Deadline
func contextWithDeadline() {

	// Deadline: 2 giây từ bây giờ
	deadline := time.Now().Add(2 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	// Kiểm tra deadline
	d, ok := ctx.Deadline()
	if ok {
		fmt.Printf("Deadline at: %s\n", d.Format(deadline.String()))
	}

	// Thực hiện API call giả lập
	go func() {
		select {
		case <-time.After(3 * time.Second):
			fmt.Println("API call hoàn thành")
		case <-ctx.Done():
			fmt.Println("API call bị hủy:", ctx.Err())
		}
	}()

	<-ctx.Done()
	fmt.Println("Context đã hết hạn")
}

func main() {
	dataRace()

	fmt.Println("\n=== Ví dụ: Nhiều channel ===")
	selectExample()

	// Chạy các ví dụ về Context
	fmt.Println("\n--- Context với Timeout (3 giây) ---")
	contextWithTimeout()
	time.Sleep(1 * time.Second)

		fmt.Println("\n--- Context với Deadline (2 giây) ---")
	contextWithDeadline()
	time.Sleep(1 * time.Second)
}
