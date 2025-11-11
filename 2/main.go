package main

import (
	"fmt"
	"sync"
)

func main() {
	pingCh := make(chan string)
	pongCh := make(chan string)
	var wg sync.WaitGroup

	wg.Add(2)
	go pinger(pingCh, pongCh, &wg)
	go ponger(pingCh, pongCh, &wg)

	wg.Wait()

	fmt.Println("Game completed!")
}

func pinger(pingCh chan<- string, pongCh <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for range 5 {
		pingCh <- "ping"
		msg := <-pongCh
		fmt.Println(msg, "->")
	}
}

func ponger(pingCh <-chan string, pongCh chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	for range 5 {
		msg := <-pingCh
		fmt.Println(msg, "<-")
		pongCh <- "pong"
	}
}
