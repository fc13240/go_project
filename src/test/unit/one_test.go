package unit

import (
	"fmt"
	"testing"
)

type H int

func TestSomething(t *testing.T) {
	t.Logf("method TestSomething")
}

func TestSomething_suffix(t *testing.T) {
	t.Logf("method TestSomething_suffix")
}

func ExampleBuffer_reader() {
	fmt.Println("method ExampleBuffer_reader")
}

func BenchmarkOne(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fmt.Println(i)
	}
}
