package main

import (
	"testing"
)

var num = 10000

type sumSquaresTest struct {
	arg      int
	expected int32
}

var sumSquaresTests = []sumSquaresTest{
	sumSquaresTest{1000, 332833500},
	sumSquaresTest{100, 328350},
	sumSquaresTest{700, 114088450},
}

func TestSumSquaresAsyncAtomic(t *testing.T) {
	for _, test := range sumSquaresTests {
		if output := sumSquaresAsyncAtomic(test.arg); output != test.expected {
			t.Errorf("Получившийся результат: %q не соответствует ожидаемому: %q", output, test.expected)
		}
	}
}

func TestSumSquaresAsyncMutex(t *testing.T) {
	for _, test := range sumSquaresTests {
		if output := sumSquaresAsyncMutex(test.arg); output != test.expected {
			t.Errorf("Получившийся результат: %q не соответствует ожидаемому: %q", output, test.expected)
		}
	}
}

func TestSumSquaresAsyncChannel(t *testing.T) {
	for _, test := range sumSquaresTests {
		if output := sumSquaresAsyncChannel(test.arg); output != test.expected {
			t.Errorf("Получившийся результат: %q не соответствует ожидаемому: %q", output, test.expected)
		}
	}
}

func TestSumSquares(t *testing.T) {
	for _, test := range sumSquaresTests {
		if output := sumSquares(test.arg); output != test.expected {
			t.Errorf("Получившийся результат: %q не соответствует ожидаемому: %q", output, test.expected)
		}
	}
}

func BenchmarkSumSquaresAsyncAtomic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sumSquaresAsyncAtomic(num)
	}
}

func BenchmarkSumSquaresAsyncMutex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sumSquaresAsyncMutex(num)
	}
}

func BenchmarkSumSquaresAsyncChannel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sumSquaresAsyncChannel(num)
	}
}

func BenchmarkSumSquares(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sumSquares(num)
	}
}
