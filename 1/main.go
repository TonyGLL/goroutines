package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func worker(id int, wg *sync.WaitGroup) {
	// Notify the WaitGroup that this goroutine is done when it exits
	defer wg.Done()

	fmt.Printf("Initializing task %d\n", id)
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	fmt.Printf("Task %d completed\n", id)
}

func main() {
	// Declare the WaitGroup
	var wg sync.WaitGroup

	// Launch 100 goroutines
	for i := range 100 {
		// Increment the WaitGroup counter
		wg.Add(1)

		// Start the goroutine
		go worker(i, &wg)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// All tasks are done
	fmt.Println("All tasks completed")
}
