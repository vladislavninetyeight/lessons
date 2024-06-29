package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os/signal"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
)

// Дан слайс N чисел сгенерированных рандомно 1-100.
// 1. Подсчитать сумму всех чисел в слайсе синхронно
// 2. Подсчитать сумму всех чисел в слайсе разделив слайс на батчи и подсчитывая сумму каждого батча в отдельной горутине

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer stop()

	fmt.Println(sum(ctx, rand.Perm(10000000)))
	fmt.Println(sumAsync(ctx, rand.Perm(10000000)))
}

// 1. Исправить гонку
// 2. Бенчмарк для сравнения первого подхода (синхронный) и второго (асинхронный)
// 3. Реализовал преждевременное прерывание программы с помощью ctrl+c с выводом промежуточного результата. Здесь нужно использовать context (signal.NotifyContext())

func sumAsync(ctx context.Context, n []int) (res int32) {
	amountCPU := runtime.NumCPU()
	wg := sync.WaitGroup{}

	step := len(n) / amountCPU // 1
	if len(n)%amountCPU != 0 {
		step += 1 // 2
	}
	next := step // 2
	prev := 0    // 0

	if len(n) < amountCPU {
		next = len(n)
	}

	for _ = range amountCPU {
		wg.Add(1)
		part := n[prev:next]
		go func() {
			defer wg.Done()
			tempRes, _ := sumBatch(ctx, part)
			atomic.AddInt32(&res, tempRes)
		}()
		if next == len(n) { // len 13
			break
		}
		prev = next
		next = next + step
		if next > len(n) {
			next = len(n)
		}
	}
	wg.Wait()
	return
}

func sum(ctx context.Context, n []int) (res int32) {
	amountCPU := runtime.NumCPU()
	step := len(n) / amountCPU // 1
	if len(n)%amountCPU != 0 {
		step += 1 // 2
	}
	next := step // 2
	prev := 0    // 0

	if len(n) < amountCPU {
		next = len(n)
	}

	for _ = range amountCPU {
		tempRes, err := sumBatch(ctx, n[prev:next])
		res += tempRes

		if err != nil {
			break
		}

		if next == len(n) { // len 13
			break
		}
		prev = next
		next = next + step
		if next > len(n) {
			next = len(n)
		}
	}
	return
}

func sumBatch(ctx context.Context, n []int) (res int32, err error) {
	isCancel := false
	context.AfterFunc(ctx, func() {
		isCancel = true
	})
	for _, v := range n {
		if isCancel {
			return res, errors.New("cancel")
		}
		res += int32(v)
	}
	return
}
