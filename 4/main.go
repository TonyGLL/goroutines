package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Task struct {
	ID   int
	Name string
}

type ProcessCenter struct {
	QueueCh     chan *Task
	ResultsCh   chan int
	wg          sync.WaitGroup
	WorkerCount int
	TaskNumber  int
}

func main() {
	start := time.Now()

	pc := NewProcessCenter(4, 50)

	for i := 1; i <= pc.WorkerCount; i++ {
		go pc.worker()
	}

	for i := 1; i <= pc.TaskNumber; i++ {
		pc.enqueue(i)
	}

	close(pc.QueueCh)

	go func() {
		pc.wg.Wait()
		close(pc.ResultsCh)
	}()

	pc.printResults()

	elapsed := time.Since(start)
	fmt.Printf("Elapsed: %s\n", elapsed)
}

func NewProcessCenter(workerCount int, taskNumber int) *ProcessCenter {
	return &ProcessCenter{
		QueueCh:     make(chan *Task, taskNumber),
		ResultsCh:   make(chan int, taskNumber),
		WorkerCount: workerCount,
		TaskNumber:  taskNumber,
	}
}

func (pc *ProcessCenter) enqueue(id int) {
	pc.wg.Add(1)
	fmt.Println("Enqueuing task", id)
	pc.QueueCh <- &Task{ID: id, Name: fmt.Sprintf("Product %d", id)}
}

func (pc *ProcessCenter) worker() {
	for task := range pc.QueueCh {
		pc.processTask(task)
	}
}

func (pc *ProcessCenter) processTask(task *Task) {
	defer pc.wg.Done()
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)))
	fmt.Printf("Task %d processed\n", task.ID)
	pc.ResultsCh <- task.ID
}

/* func (pc *ProcessCenter) addResult(task *Task) {
	pc.ResultsMutex.Lock()
	pc.Results = append(pc.Results, task)
	pc.ResultsMutex.Unlock()
} */

func (pc *ProcessCenter) printResults() {
	results := []int{}
	for result := range pc.ResultsCh {
		results = append(results, result)
	}
	fmt.Println("\n=== SUMMARY ===")
	fmt.Printf("Total tasks processed: %d\n", len(results))
}
