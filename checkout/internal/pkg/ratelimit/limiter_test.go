// Rate limiter tests
package ratelimit

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

// Test if the number of requests per second limit is correctly enforced
func TestLimiter_RPSCondition(t * testing.T) {
	rps := 10
	maxConcurrency := 5
	longRequestCount := 10
	shortRequestCount := 15
	dateFormat := "15:04:05.9999"

	testChan := make(chan time.Time, longRequestCount + shortRequestCount)

	ctx := context.Background()

	limiter := New(uint16(rps), uint16(maxConcurrency))
	var wg sync.WaitGroup
	wg.Add(longRequestCount + shortRequestCount)

	for i := 1; i <= longRequestCount; i++ {
		go func(reqNum int) {
			defer wg.Done()
			
			err := limiter.Acquire(ctx)
			defer limiter.Release()
			if err != nil {
				t.Errorf("wait: %v", err)
				return
			}

			startTime := time.Now()
			testChan <- startTime
			// Simulate a long-running operation
			time.Sleep(2 * time.Second)
			fmt.Printf("Long request %d completed, start: %v, end: %v\n", reqNum, startTime.Format(dateFormat), time.Now().Format(dateFormat))

		}(i)
	}

	for i := 1; i <= shortRequestCount; i++ {
		go func(reqNum int) {
			defer wg.Done()
			
			err := limiter.Acquire(ctx)
			defer limiter.Release()
			if err != nil {
				t.Errorf("wait: %v", err)
				return
			}

			startTime := time.Now()
			testChan <- startTime

			// Simulate a short-running operation
			time.Sleep(150 * time.Millisecond)
			fmt.Printf("Short request %d completed, start: %v, end: %v\n", reqNum, startTime.Format(dateFormat), time.Now().Format(dateFormat))

		}(i)
	}

	wg.Wait()
	limiter.Close()
	close(testChan)
	timeArr := make([]time.Time, 0, longRequestCount + shortRequestCount)

	for d := range testChan {
		timeArr = append(timeArr, d)
	}

	for i := 0; i < longRequestCount + shortRequestCount; i++ {
		fmt.Println(timeArr[i].Format(dateFormat))
	}

	for i := 0; i < longRequestCount + shortRequestCount - rps; i += 1 {
		for j := i + 1; j < longRequestCount + shortRequestCount; j++ {
			if timeArr[j].Sub(timeArr[i]) > time.Second {
				if j - i > rps {
					t.Errorf("More than %v requests per second: first: %v, last: %v", rps, timeArr[i], timeArr[i + rps])
				}
				break
			}
		}
    }
}
