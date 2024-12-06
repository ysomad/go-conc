package main

import "fmt"

func join[T any](ch1 <-chan T, ch2 <-chan T) <-chan T {
	out := make(chan T)

	go func() {
		defer close(out)

		for ch1 != nil && ch2 != nil {
			select {
			case v, ok := <-ch1:
				if !ok {
					ch1 = nil
					continue
				}
				out <- v
			case v, ok := <-ch2:
				if !ok {
					ch2 = nil
					continue
				}
				out <- v
			}
		}
	}()

	return out
}

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		for i := 10000; i < 20000; i++ {
			ch1 <- i
		}
		close(ch1)
	}()

	go func() {
		for i := 30000; i < 40000; i++ {
			ch2 <- i
		}
		close(ch2)
	}()

	for v := range join(ch1, ch2) {
		fmt.Printf("GOT: %v\n", v)
	}
}
