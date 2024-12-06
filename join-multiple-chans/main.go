package main

import (
	"fmt"
	"sync"
)

func join[T any](chans ...chan T) <-chan T {
	out := make(chan T)
	wg := &sync.WaitGroup{}
	wg.Add(len(chans))

	go func() {
		wg.Wait()
		close(out)
	}()

	for _, ch := range chans {
		go func() {
			defer wg.Done()

			for v := range ch {
				out <- v
			}
		}()
	}

	return out
}

func main() {
	var chans []chan int

	for i := range 100 {
		ch := make(chan int)
		chans = append(chans, ch)

		go func() {
			for j := i; j < i*100; j++ {
				ch <- j
			}

			close(ch)
		}()
	}

	for v := range join(chans...) {
		fmt.Println(v)
	}
}
