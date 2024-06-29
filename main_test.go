package main

import (
	"context"
	"math/rand"
	"os"
	"os/signal"
	"testing"
)

var num = 1000000

type sumTest struct {
	arg      []int
	expected int32
}

var sumTests = []sumTest{
	sumTest{[]int{1, 2, 5, 6, 4, 1, 3, 4, 5, 6, 7, 7}, 51},
	sumTest{[]int{1, 1, 2, 4}, 8},
	sumTest{[]int{7, 7, 4, 8, 7, 2, 3, 4, 5, 1, 2, 4, 4}, 58},
}

func BenchmarkSumAsync(b *testing.B) {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	for i := 0; i < b.N; i++ {
		sumAsync(ctx, rand.Perm(num))
	}
}

func BenchmarkSum(b *testing.B) {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	for i := 0; i < b.N; i++ {
		sum(ctx, rand.Perm(num))
	}
}
