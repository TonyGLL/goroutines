package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Task represents a unit of work to be processed
type Task struct {
	ID   int    // Unique identifier for the task
	Name string // Name or description of the task
}

// ProcessCenter manages task distribution and processing using workers
type ProcessCenter struct {
	QueueCh     chan *Task     // Channel for enqueuing tasks
	ResultsCh   chan int       // Channel for collecting results
	wg          sync.WaitGroup // WaitGroup to synchronize task completion
	WorkerCount int            // Number of worker goroutines
	TaskNumber  int            // Total number of tasks to process
}

// main initializes the ProcessCenter, starts workers, enqueues tasks, and prints results
func main() {
	start := time.Now() // Record start time

	pc := NewProcessCenter(4, 50) // Create a ProcessCenter with 4 workers and 50 tasks

	// Start worker goroutines
	for i := 1; i <= pc.WorkerCount; i++ {
		go pc.worker()
	}

	// Enqueue all tasks
	for i := 1; i <= pc.TaskNumber; i++ {
		pc.enqueue(i)
	}

	close(pc.QueueCh) // No more tasks will be sent

	// Wait for all tasks to finish, then close results channel
	go func() {
		pc.wg.Wait()
		close(pc.ResultsCh)
	}()

	pc.printResults() // Print summary of processed tasks

	elapsed := time.Since(start) // Calculate elapsed time
	fmt.Printf("Elapsed: %s\n", elapsed)
}

// NewProcessCenter creates and initializes a new ProcessCenter
func NewProcessCenter(workerCount int, taskNumber int) *ProcessCenter {
	return &ProcessCenter{
		QueueCh:     make(chan *Task, taskNumber), // Buffered channel for tasks
		ResultsCh:   make(chan int, taskNumber),   // Buffered channel for results
		WorkerCount: workerCount,
		TaskNumber:  taskNumber,
	}
}

// enqueue adds a new task to the queue and increments the WaitGroup
func (pc *ProcessCenter) enqueue(id int) {
	pc.wg.Add(1) // Increment WaitGroup for each task
	fmt.Println("Enqueuing task", id)
	pc.QueueCh <- &Task{ID: id, Name: fmt.Sprintf("Product %d", id)} // Send task to queue
}

// worker processes tasks from the queue until the channel is closed
func (pc *ProcessCenter) worker() {
	for task := range pc.QueueCh {
		pc.processTask(task)
	}
}

// processTask simulates processing a task and sends the result
func (pc *ProcessCenter) processTask(task *Task) {
	defer pc.wg.Done()                                           // Signal task completion
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(500))) // Simulate work
	fmt.Printf("Task %d processed\n", task.ID)
	pc.ResultsCh <- task.ID // Send result to results channel
}

// printResults collects and prints all processed task results
func (pc *ProcessCenter) printResults() {
	results := []int{}
	for result := range pc.ResultsCh {
		results = append(results, result)
	}
	fmt.Println("\n=== SUMMARY ===")
	fmt.Printf("Total tasks processed: %d\n", len(results))
}
