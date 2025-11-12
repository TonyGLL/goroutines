package main

import (
	"fmt"
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

	// Start worker goroutines, each with its own ID
	for i := 1; i <= pc.WorkerCount; i++ {
		go pc.worker(i)
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

	// Collect results as they come in
	results := []int{}
	for result := range pc.ResultsCh {
		results = append(results, result)
	}

	fmt.Println("\n=== SUMMARY ===")
	fmt.Printf("Total tasks processed: %d\n", len(results))

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
// id: unique identifier for the worker
func (pc *ProcessCenter) worker(id int) {
	for task := range pc.QueueCh {
		pc.processTask(id, task)
	}
}

// processTask simulates processing a task with a possible timeout
// workerID: the ID of the worker processing the task
// task: the task to process
func (pc *ProcessCenter) processTask(workerID int, task *Task) {
	defer pc.wg.Done() // Signal task completion

	resultCh := make(chan int, 1) // Buffer of 1 to avoid goroutine leaks

	// Internal goroutine: simulates slow work
	go func() {
		duration := 1000 // Default duration in milliseconds
		if task.ID%10 == 0 {
			duration = 3000 // Every 10th task takes longer
		}
		time.Sleep(time.Millisecond * time.Duration(duration)) // Simulate work
		resultCh <- task.ID                                    // Send result to buffer and exit without blocking
	}()

	// Select: wait for result or timeout
	select {
	case id := <-resultCh:
		fmt.Printf("Worker %d: Task %d processed\n", workerID, id)
		pc.ResultsCh <- id // Send result to results channel
	case <-time.After(2 * time.Second):
		fmt.Printf("Worker %d: Timeout processing task %d\n", workerID, task.ID)
	}
}
