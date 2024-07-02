package main

import (
	"awesomeProject/NewTime"
	"awesomeProject/User"
	"awesomeProject/calculate"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os/signal"
	"runtime"
	"sort"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

// 1. Исправить гонку
// 2. Бенчмарк для сравнения первого подхода (синхронный) и второго (асинхронный)
// 3. Реализовал преждевременное прерывание программы с помощью ctrl+c с выводом промежуточного результата. Здесь нужно использовать context (signal.NotifyContext())

func main() {
	firstPart()
	//secondPart()
	//thirdPart()
	//fourthPart()
}

func fourthPart() {
	fmt.Println(calculate.Calculate(1, "2", calculate.AmountStruct{Amount: 3}))
}

func thirdPart() {
	users := User.Users{User.User{
		ID:        "1",
		CreatedAt: time.Time{}.AddDate(2014, 7, 24),
	},
		User.User{
			ID:        "2",
			CreatedAt: time.Time{}.AddDate(2013, 7, 24),
		},
		User.User{
			ID:        "3",
			CreatedAt: time.Time{}.AddDate(2017, 7, 24),
		}}

	sort.Sort(users)
	fmt.Println(users)
}

func secondPart() {
	timeISO := NewTime.TimeISO8601{}
	timeUNIX := NewTime.TimeUnix{}

	unix := 1719928793
	iso := "2020-07-10 15:00:00.000"
	jsonUNIX, err := json.Marshal(unix)

	if err != nil {
		panic(err)
	}

	jsonISO, _ := json.Marshal(iso)
	err = timeISO.UnmarshalJSON(jsonISO)
	if err != nil {
		panic(err)
	}

	err = timeUNIX.UnmarshalJSON(jsonUNIX)
	if err != nil {
		panic(err)
	}

	marshalJsonIso, err := timeISO.MarshalJSON()
	marshalJsonUNIX, err := timeUNIX.MarshalJSON()

	fmt.Print(timeISO.Time, timeUNIX.Time, string(marshalJsonIso), string(marshalJsonUNIX))
}

func firstPart() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer stop()
	fmt.Println(sum(ctx, rand.Perm(10000000)))
	fmt.Println(sumAsync(ctx, rand.Perm(10000000)))
}

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
