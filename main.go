package main

import (
	"fmt"
	"sync"
)

// Barrier represents a synchronization point for a group of goroutines.
type Barrier struct {
	toggle         bool
	swapped        bool
	initialCount   int
	workCount      int
	wg             sync.WaitGroup
	mu             sync.Mutex
	workNotify     chan struct{}
	decisionNotify chan bool
}

// NewBarrier creates a new Barrier with the specified count.
func NewBarrier(count int) *Barrier {
	return &Barrier{
		toggle:         false,
		swapped:        true,
		initialCount:   count + 1,
		workCount:      count + 1,
		workNotify:     make(chan struct{}),
		decisionNotify: make(chan bool),
	}
}

// count: Tracks the number of goroutines that need to reach the barrier.
// workNotify: A channel used to broadcast a signal to waiting goroutines.
// mu: A mutex to ensure atomic updates to count.

func (b *Barrier) WaitForWorkloadComplete() {
	b.mu.Lock()
	b.workCount--
	if b.workCount == 0 {
		close(b.workNotify) // Notify all waiting goroutines
	}
	b.mu.Unlock()

	// Wait for the notification
	<-b.workNotify
}

func (b *Barrier) WaitForDecision() bool {
	decision := <-b.decisionNotify
	return decision
}

func (barrier *Barrier) barrieredCompareAndSwap(worker int, array *[10]int) {
startWork:
	fmt.Println("Worker", worker, "started")
	// Simulate work
	// time.Sleep(time.Second * time.Duration(worker+1))

	index := worker * 2

	if barrier.toggle {
		index = (worker * 2) + 1
		if index == (len(array) - 1) {
			fmt.Println("last val in array, return")
			barrier.WaitForWorkloadComplete()
			if barrier.WaitForDecision() {
				goto startWork
			} else {
				return
			}
		}
	}

	fmt.Printf("worker: %v, index: %v", worker, index)

	if array[index] > array[index+1] {
		barrier.mu.Lock()
		temp := array[index]
		array[index] = array[index+1]
		array[index+1] = temp
		// barrier.mu.Lock()
		barrier.swapped = true
		barrier.mu.Unlock()
	}

	fmt.Println("Worker", worker, "finished work, waiting")

	// Wait for other workers to complete
	barrier.WaitForWorkloadComplete()

	if barrier.WaitForDecision() {
		// Continue with the next phase of work
		goto startWork
	} else {
		fmt.Println("Worker", worker, "returning from goroutine")
		return
	}

}

func main() {
	array := [...]int{8, 6, 10, 9, 5, 2, 7, 1, 4, 3}

	fmt.Println("array:", array)

	const numWorkers = len(array) / 2

	// Create a barrier with the number of expected workers
	barrier := NewBarrier(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go barrier.barrieredCompareAndSwap(i, &array)
	}

	for {
		// Start multiple workers
		barrier.WaitForWorkloadComplete()
		fmt.Println("barrier swapped value: ", barrier.swapped)
		if barrier.swapped == false {
			// for all worker threads
			// send decision false
			for j := 0; j < numWorkers; j++ {
				barrier.decisionNotify <- false
			}
			break
		}
		barrier.mu.Lock()
		barrier.workNotify = make(chan struct{})
		barrier.workCount = barrier.initialCount
		barrier.swapped = false
		barrier.toggle = !barrier.toggle
		barrier.mu.Unlock()
		for j := 0; j < numWorkers; j++ {
			barrier.decisionNotify <- true
		}
	}

	// Wait for all workers to complete
	// barrier.WaitForWorkloadComplete()

	fmt.Println(array)

	fmt.Println("All workers completed")
}
