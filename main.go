package main

import (
	"fmt"
	"sync"
)

func someNumberStrings() <-chan string {
	ch := make(chan string)
	numberStrings := []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten"}

	go func() {
		for _, numberString := range numberStrings {
			ch <- numberString
		}

		close(ch)
		return
	}()

	return ch
}

func someNumbers() <-chan string {
	ch := make(chan string)
	numbers := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}

	go func() {
		for _, number := range numbers {
			ch <- number
		}

		close(ch)
		return
	}()

	return ch
}

func fanIn(channels ...<-chan string) <-chan string {
	// create WaitGroup
	var wg sync.WaitGroup

	// create multiplexed stream
	multiplexedStream := make(chan string)

	// function that takes data from one of the two streams
	// and writes it into our multiplexed stream
	multiplex := func(c <-chan string) {
		// when we finish reading from the data stream, tell
		// our WaitGroup that we are no longer reading from that stream
		defer wg.Done()

		// write from the data stream to the multiplexed stream
		for i := range c {
			multiplexedStream <- i
		}
	}

	// tell our WaitGroup that we are waiting for two stream to finish
	wg.Add(len(channels))

	// read concurrently from the data streams into the multiplexed stream
	for _, c := range channels {
		go multiplex(c)
	}

	go func() {
		// wait for both streams to close and then close the multiplexed stream
		wg.Wait()
		close(multiplexedStream)
	}()

	return multiplexedStream
}

func main() {
	ch1 := someNumberStrings()
	ch2 := someNumbers()

	exit := make(chan struct{})
	mergedCh := fanIn(ch1, ch2)

	go func() {
		for val := range mergedCh {
			fmt.Println(val)
		}

		close(exit)
	}()

	<-exit

	fmt.Println("bye")
}
