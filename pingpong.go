package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	count := make(chan int, 1000)
	done := make(chan bool)
	var wg sync.WaitGroup

	for i := 0; i < 12; i++ {
		pingChan := make(chan int, 1)
		pongChan := make(chan int, 1)

		wg.Add(1)
		go ping(done, pingChan, pongChan, count, &wg)
		wg.Add(1)
		go pong(done, pongChan, pingChan, count, &wg)

		pingChan <- 1
	}

	time.Sleep(10 * time.Second)
	close(done)
	wg.Wait()

	close(count)
	total := 0
	for c := range count {
		total += c
	}
	fmt.Printf("pings or pongs per second: %d \n", total/10)
}

func ping(done chan bool, pingChan <-chan int, pongChan, count chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	pingCount := 0
	for {
		select {
		case <-done:
			count <- pingCount
			return
		case ball := <-pingChan:
			pingCount++
			pongChan <- ball
		}
	}
}

func pong(done chan bool, pongChan <-chan int, pingChan, count chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	pongCount := 0
	for {
		select {
		case <-done:
			count <- pongCount
			return
		case ball := <-pongChan:
			pongCount++
			pingChan <- ball
		}
	}
}
