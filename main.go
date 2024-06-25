package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	fmt.Println(sumSquares(700))
	fmt.Println(sumSquaresAsyncAtomic(1000))
	fmt.Println(sumSquaresAsyncMutex(1000))
	fmt.Println(sumSquaresAsyncChannel(1000))
}

func sumSquaresAsyncAtomic(n int) (sum int32) {
	wg := sync.WaitGroup{}
	for i := range n {
		wg.Add(1)
		go func() {
			defer wg.Done()

			atomic.AddInt32(&sum, int32(i*i))
		}()
	}
	wg.Wait()
	return
}

func sumSquaresAsyncMutex(n int) (sum int32) {
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	for i := range n {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			sum += int32(i * i)
			mu.Unlock()
		}()
	}
	wg.Wait()
	return
}

func sumSquaresAsyncChannel(n int) (sum int32) {
	wg := sync.WaitGroup{}
	ch := make(chan struct{}, 1)
	for i := range n {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ch <- struct{}{}
			sum += int32(i * i)
			<-ch
		}()
	}
	wg.Wait()
	return
}

func sumSquares(n int) (sum int32) {
	for i := range n {
		sum += int32(i * i)
	}
	return
}
