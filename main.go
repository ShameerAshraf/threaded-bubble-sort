package main

import (
	"fmt"
	"sync"
)

// Barrier represents a synchronization point for a group of goroutines.
type Barrier struct {
	toggle       bool
	swapped      bool
	initialCount int
	count        int
	wg           sync.WaitGroup
	mu           sync.Mutex
	notify       chan struct{}
}

// NewBarrier creates a new Barrier with the specified count.
func NewBarrier(count int) *Barrier {
	return &Barrier{
		toggle:       false,
		swapped:      true,
		initialCount: count + 1,
		count:        count + 1,
		notify:       make(chan struct{}),
	}
}

// count: Tracks the number of goroutines that need to reach the barrier.
// notify: A channel used to broadcast a signal to waiting goroutines.
// mu: A mutex to ensure atomic updates to count.

func (b *Barrier) Wait() {
	b.mu.Lock()
	b.count--
	if b.count == 0 {
		close(b.notify) // Notify all waiting goroutines
	}
	b.mu.Unlock()

	// Wait for the notification
	<-b.notify
}

func (barrier *Barrier) barrieredCompareAndSwap(worker int, array *[10]int) {
	fmt.Println("Worker", worker, "started")
	// Simulate work
	// time.Sleep(time.Second * time.Duration(worker+1))

	index := worker * 2

	if barrier.toggle {
		index = (worker * 2) + 1
		if index == (len(array) - 1) {
			fmt.Println("last val in array, return")
			barrier.Wait()
			return
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
	barrier.Wait()

	// Continue with the next phase of work
	fmt.Println("Worker", worker, "returning from goroutine")
}

func main() {
	array := [...]int{8, 6, 10, 9, 5, 2, 7, 1, 4, 3}

	fmt.Println("array:", array)

	const numWorkers = len(array) / 2

	// Create a barrier with the number of expected workers
	barrier := NewBarrier(numWorkers)

	for {
		// Start multiple workers
		for i := 0; i < numWorkers; i++ {
			go barrier.barrieredCompareAndSwap(i, &array)
		}
		barrier.Wait()
		fmt.Println("barrier swapped value: ", barrier.swapped)
		if barrier.swapped == false {
			break
		}
		barrier.mu.Lock()
		barrier.notify = make(chan struct{})
		barrier.count = barrier.initialCount
		barrier.swapped = false
		barrier.toggle = !barrier.toggle
		barrier.mu.Unlock()
	}

	// Wait for all workers to complete
	// barrier.Wait()

	fmt.Println(array)

	fmt.Println("All workers completed")
}
