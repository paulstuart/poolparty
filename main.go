// Package main is an examination of the issues brought up by https://www.arp242.net/go-easy.html

// You can edit this code!
// Click here and start typing.
package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

const (
	jobs   = 20 // Run 20 jobs in total.
	limit  = 3
	millis = 500
)

var ()

func main() {
	// orig
	//redo()
	//sema()
	//less()
	okay()
}

// redo leverages channels for job propagation
func redo() {
	in := make(chan int)
	out := make(chan string)

	// feed the worker queue
	go func() {
		for i := 1; i <= jobs; i++ {
			in <- i
		}
		close(in)
	}()

	// create the worker pool
	var wg sync.WaitGroup
	wg.Add(limit)
	for i := 0; i < limit; i++ {
		go func(i int) {
			for x := range in {
				// "do work"
				time.Sleep(millis * time.Millisecond)
				out <- fmt.Sprintf("[%d]: %d", i, x)
			}
			wg.Done()
		}(i)
	}

	// this links the completion of the processing to the finalizing of the collection
	go func() {
		wg.Wait()
		close(out)
	}()

	// collect worker results
	for s := range out {
		fmt.Println("GOT:", s)
	}

	fmt.Println("done")
}

// sema uses semaphores for job distribution
func sema() {
	ctx := context.Background()
	sem := semaphore.NewWeighted(3)
	var wg sync.WaitGroup // Keep track of which jobs are finished.
	wg.Add(jobs)

	for i := 1; i <= jobs; i++ {
		if err := sem.Acquire(ctx, 1); err != nil {
			fmt.Println("sem err say:", err)
			continue
		}

		// Start a job.
		go func(i int) {
			// "do work"
			time.Sleep(millis * time.Millisecond)
			fmt.Println(i)

			sem.Release(1)
			wg.Done() // Signal that this job is done.
		}(i)
	}

	wg.Wait() // Wait until all jobs are done.
	fmt.Println("done")
}

func okay() {
	// first create a pool of workers that take their inputs from a channel
	var wg sync.WaitGroup
	wg.Add(jobs)
	ch := make(chan int)
	for i := 0; i < limit; i++ {

		// Start jobs.
		go func(i int) {
			for x := range ch {
				// "do work"
				time.Sleep(millis * time.Millisecond)
				fmt.Printf("[%d]: %d\n", i, x)
				wg.Done()
			}
		}(i)
	}

	// now feed the work into the workers channel
	for i := 1; i <= jobs; i++ {
		ch <- i
	}

	wg.Wait()
	fmt.Println("closing channel")
	close(ch)
	fmt.Println("done")
}

// orig is the sample shown by https://www.arp242.net/go-easy.html
func orig() {
	running := make(chan bool, 3) // Limit concurrent jobs to 3.
	var wg sync.WaitGroup         // Keep track of which jobs are finished.
	wg.Add(jobs)

	for i := 1; i <= jobs; i++ {
		running <- true // Fill running; this will block and wait if it's already full.
		// Start a job.
		go func(i int) {
			defer func() {
				<-running // Drain running so new jobs can be added.
				wg.Done() // Signal that this job is done.
			}()

			// "do work"
			time.Sleep(millis * time.Millisecond)
			fmt.Println(i)
		}(i)
	}

	wg.Wait() // Wait until all jobs are done.
	fmt.Println("done")
}

// less is orig() with the superfluous defer's removed
func less() {
	running := make(chan bool, 3) // Limit concurrent jobs to 3.
	var wg sync.WaitGroup         // Keep track of which jobs are finished.
	wg.Add(jobs)

	for i := 1; i <= jobs; i++ {
		running <- true // Fill running; this will block and wait if it's already full.
		// Start a job.
		go func(i int) {
			// "do work"
			time.Sleep(millis * time.Millisecond)
			fmt.Println(i)
			<-running // Drain running so new jobs can be added.
			wg.Done() // Signal that this job is done.
		}(i)
	}

	wg.Wait() // Wait until all jobs are done.
	fmt.Println("done")
}
