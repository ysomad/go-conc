package main

import (
	"fmt"
	"sync"
)

func split[T any](in <-chan T, n int) []chan T {
	out := make([]chan T, n)

	for i := range n {
		out[i] = make(chan T)
	}

	go func() {
		for v := range in {
			for _, ch := range out {
				ch <- v
			}
		}

		for _, ch := range out {
			close(ch)
		}
	}()

	return out
}

func main() {
	in := make(chan int)

	go func() {
		for i := range 10000 {
			in <- i
		}

		close(in)
	}()

	chans := split(in, 10)

	wg := &sync.WaitGroup{}
	wg.Add(len(chans))

	for i, ch := range chans {
		go func() {
			defer wg.Done()

			for v := range ch {
				fmt.Printf("v = %v, ch = %d\n", v, i)
			}
		}()
	}

	wg.Wait()
}
