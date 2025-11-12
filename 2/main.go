package main // Package declaration

import (
    "fmt"
    "sync"
)

// main sets up the ping-pong game using goroutines and channels
func main() {
	// Create channels for communication between goroutines
	pingCh := make(chan string)
	pongCh := make(chan string)
	var wg sync.WaitGroup // WaitGroup to synchronize goroutines

	wg.Add(2) // We will wait for two goroutines
	go pinger(pingCh, pongCh, &wg) // Start the pinger goroutine
	go ponger(pingCh, pongCh, &wg) // Start the ponger goroutine

	wg.Wait() // Wait for both goroutines to finish

	fmt.Println("Game completed!")
}

// pinger sends "ping" messages and waits for "pong" responses
// pingCh: channel to send "ping" messages
// pongCh: channel to receive "pong" messages
// wg: WaitGroup to signal completion
func pinger(pingCh chan<- string, pongCh <-chan string, wg *sync.WaitGroup) {
	defer wg.Done() // Signal that this goroutine is done

	for range 5 { // Repeat 5 times
		pingCh <- "ping" // Send "ping" to ponger
		msg := <-pongCh // Wait for "pong" response
		fmt.Println(msg, "->") // Print the received message
	}
}

// ponger waits for "ping" messages and replies with "pong"
// pingCh: channel to receive "ping" messages
// pongCh: channel to send "pong" messages
// wg: WaitGroup to signal completion
func ponger(pingCh <-chan string, pongCh chan<- string, wg *sync.WaitGroup) {
	defer wg.Done() // Signal that this goroutine is done

	for range 5 { // Repeat 5 times
		msg := <-pingCh // Wait for "ping" from pinger
		fmt.Println(msg, "<-") // Print the received message
		pongCh <- "pong" // Send back "pong"
	}
}
