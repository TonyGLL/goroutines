package main

import (
	"fmt"
	"sync"
	"time"
)

type Product struct {
	ID   int
	Name string
}

func main() {
	// Creating a buffered channel with a capacity of 3
	ch := make(chan *Product, 3)

	// Start time (monotonic clock is used when paired with time.Since)
	start := time.Now()

	// WaitGroup to wait for the consumer to finish
	var wg sync.WaitGroup

	// Starting the consumer goroutine
	wg.Add(1)

	// Start the consumer goroutine
	go processProduct(ch, &wg)

	// Producing 10 products
	for i := 1; i <= 10; i++ {
		product := &Product{ID: i, Name: fmt.Sprintf("Product %d", i)}
		fmt.Printf("Generating and enqueing %s\n", product.Name)

		// Sending the product to the channel
		ch <- product
	}

	// Closing the channel to signal no more products will be sent
	close(ch)

	// Waiting for the consumer to finish processing
	wg.Wait()

	// Measure elapsed wall-clock time since start
	elapsed := time.Since(start)
	fmt.Printf("All products have been consumed and processed.\nElapsed: %s (ms=%d, ns=%d)\n", elapsed, elapsed.Milliseconds(), elapsed.Nanoseconds())
}

func processProduct(ch <-chan *Product, wg *sync.WaitGroup) {
	// Signal that the goroutine is done when this function exits
	defer wg.Done()

	// Consuming products from the channel
	for product := range ch {
		fmt.Printf("Processing %s\n", product.Name)
		time.Sleep(500 * time.Millisecond)
	}
}
