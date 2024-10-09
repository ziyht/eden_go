package erb

import (
	"testing"

	"github.com/antlabs/gstl/rbtree"
)

func BenchmarkMulAdd(b *testing.B) {
	rb := New[int, int]()
	for i := 0; i < b.N; i ++ {
		rb.MulAdd(1, i + 1)
	}
}

func BenchmarkGSTLSet(b *testing.B) {
	rb := rbtree.New[int, int]()
	for i := 0; i < b.N; i ++ {
		rb.Set(i, i + 1)
	}
}

func BenchmarkSet(b *testing.B) {
	rb := New[int, int]()
	for i := 0; i < b.N; i ++ {
		rb.Set(i, i + 1)
	}
}

// go test -bench='BenchmarkRangeIter|BenchmarkRangeStack' -benchtime=10000000x -count=1 -benchmem
func BenchmarkRangeStack(b *testing.B) {
	t := New[int, int]()

	for i := 0; i < b.N; i ++ {
		t.Set(i, i + 1)
	}

	b.ResetTimer()

	t.Range(func(int, int)bool{return true})
}