package erb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/antlabs/gstl/rbtree"
	"golang.org/x/exp/rand"
)


func TestNew(t *testing.T) {

	root := New[int, int]()

	assert.Equal(t, int64(0), root.Size())
}

func TestAdd(t *testing.T) {

	root := New[int, int]()

	assert.True(t, root.Add(1, 1))
	assert.False(t, root.Add(1, 2))
	assert.False(t, root.Add(1, 3))

	assert.Equal(t, int64(1), root.Size())

	for i := 0; i < 10000; i++ {
		root.Add(i, i)
	}
	assert.Equal(t, int64(10000), root.Size())

	for i := 0; i < 10000; i++ {
		root.Add(i, i)
	}
	assert.Equal(t, int64(10000), root.Size())

	for i := 0; i < 10000; i++ {
		root.MulAdd(i, i)
	}
	assert.Equal(t, int64(20000), root.Size())

	for i := 0; i < 10000; i++ {
		root.MulAdd(i, i)
	}
	assert.Equal(t, int64(30000), root.Size())

	assert.Equal(t, true, root.isRBTree())
}

func TestIsRBTree(t *testing.T) {

	root := New[int, int]()

	assert.Equal(t, true, root.isRBTree())

	for i := 0; i < 10000; i++ {
		root.Add(i, i)
		assert.Equal(t, true, root.isRBTree())
	}

	for i := 0; i < 10000; i++ {
		root.MulAdd(i, i)
		assert.Equal(t, true, root.isRBTree())
	}

	for i := 0; i < 10000; i++ {
	  root.Del(i)
		assert.Equal(t, true, root.isRBTree())
	}
	assert.Equal(t, int64(10000), root.Size())

	for i := 0; i < 10000; i++ {
	  root.Del(i)
		assert.Equal(t, true, root.isRBTree())
	}
	assert.Equal(t, int64(0), root.Size())
}

func TestBasic(t *testing.T) {

	root := New[int32, int32]()

	var keys []int32

	for i := 0; i < 10000; i++ {
		keys = append(keys, rand.Int31())

		root.Add(keys[i], keys[i])
	}
	assert.Equal(t, int64(10000), root.Size())

	isRBTree, e := root.isRBTree2()
	assert.Equalf(t, true, isRBTree, "%s", e)

	iter := root.First()
	prev := iter.Val
	for i := 1; iter != nil; iter = iter.Next() {
		assert.True(t, prev <= iter.Val)
		//t.Logf("%d: <%d, %d>", i, iter.Key(), iter.Val)
		prev = iter.Val
		i++
	}

	iter = root.Last()
	next := iter.Val
	for i := 1; iter != nil; iter = iter.Prev() {
		assert.True(t, iter.Val <= next)
		// t.Logf("%d: %d", i, iter.Val)
		next = iter.Val
		i++
	}

	for i := 0 ; i < 10000; i++ {
		assert.Equal(t, keys[i], root.Val(keys[i]))
	}
}

func TestMultiInsert(t *testing.T) {

	rb := New[int, int]()

	assert.True (t, rb.Add(1, 1)) 
	assert.False(t, rb.Add(1, 2))
	assert.False(t, rb.Add(1, 3))

	assert.Equal(t, int64(1), rb.Size())

	rb = New[int, int]()
	cnt := 10000
	for i := 0; i < cnt; i ++ {
		rb.MulAdd((i / 100) + 1, i + 1)
	}
	assert.Equal(t, int64(cnt), rb.Size())

	iter := rb.First()
	for i := 0; iter != nil; iter = iter.Next() {
		assert.Equal(t, i + 1, iter.Val)
		i++
	}
}

func TestSet(t *testing.T) {

	rb := New[int32, int]()

	rb.Set(1, 1)
	assert.Equal(t, int64(1), rb.Size())
	assert.Equal(t, 1, rb.Val(1))

	rb.Set(1, 2)
	assert.Equal(t, int64(1), rb.Size())
	assert.Equal(t, 2, rb.Val(1))

	rb.Set(2, 1)
	assert.Equal(t, int64(2), rb.Size())
	assert.Equal(t, 2, rb.Val(1))
	assert.Equal(t, 1, rb.Val(2))

	rb.Set(2, 2)
	assert.Equal(t, int64(2), rb.Size())
	assert.Equal(t, 2, rb.Val(1))
	assert.Equal(t, 2, rb.Val(2))
}

func TestPop(t *testing.T) {

	rb := New[int, int]()
	cnt := 100
	for i := 0; i < cnt; i ++ {
		rb.MulAdd((i / 10) + 1, i + 1)
	}

	for i := 0; i < 50; i ++ {
		node := rb.PopFirst()
		assert.Equal(t, node.Val, i + 1)
		assert.True(t, rb.isRBTree())
	}

	for i := 0; i < 50; i ++ {
		node := rb.PopLast()
		assert.Equal(t, node.Val, 100 - i)
		assert.True(t, rb.isRBTree())
	}

	assert.Equal(t, int64(0), rb.Size())
	assert.Equal(t, (*Node[int,int])(nil), rb.PopFirst())
	assert.Equal(t, (*Node[int,int])(nil), rb.PopLast())
	assert.True(t, rb.isRBTree())
}

func TestPopKeyFirstAndLast(t *testing.T) {
	rb := New[int, int]()
	cnt := 100
	for i := 0; i < cnt; i ++ {
		rb.MulAdd((i / 10) + 1, i + 1) // 1:1, 1:2, 1:3... 10:99, 10:100
	}

	{
	  // 0 not exist
		node := rb.DelFirst(0)
		assert.Equal(t, (*Node[int,int])(nil), node)

		node = rb.DelLast(0)
		assert.Equal(t, (*Node[int,int])(nil), node)

		// 11 not exist
		node = rb.DelFirst(11)
		assert.Equal(t, (*Node[int,int])(nil), node)

		node = rb.DelLast(11)
		assert.Equal(t, (*Node[int,int])(nil), node)

		// 1 2 3 4 5 6 7 8 9 10
		// |
		node = rb.DelFirst(1)
		assert.Equal(t, 1, node.Val)
		assert.True(t, rb.isRBTree())

		// 2 3 4 5 6 7 8 9 10
		// |
		node = rb.DelFirst(1)
		assert.Equal(t, 2, node.Val)
		assert.True(t, rb.isRBTree())

		// 3 4 5 6 7 8 9 10
		//               |
		node = rb.DelLast(1)
		assert.Equal(t, 10, node.Val)
		assert.True(t, rb.isRBTree())

		// 3 4 5 6 7 8 9
		//             |
		node = rb.DelLast(1)
		assert.Equal(t, 9, node.Val)
		assert.True(t, rb.isRBTree())
	}
}

func TestPopKeyFirst(t *testing.T) {
	rb := New[int, int]()
	cnt := 100
	for i := 0; i < cnt; i ++ {
		rb.MulAdd((i / 10) + 1, i + 1)
	}

	{
		for i := 0; i < 10; i ++ {
			for j := 0; j < 10; j++ {
				node := rb.DelFirst(j + 1)
				assert.Equal(t, i + 1 + 10 * j, node.Val)
				assert.True(t, rb.isRBTree())
			}
		}

		assert.Equal(t, int64(0), rb.Size())
		assert.True(t, rb.isRBTree())
	}
}

func TestPopKeyLast(t *testing.T) {
	rb := New[int, int]()
	cnt := 100
	for i := 0; i < cnt; i ++ {
		rb.MulAdd((i / 10) + 1, i + 1)
	}

	{
		for i := 0; i < 10; i ++ {
			for j := 0; j < 10; j++ {
				node := rb.DelLast(j + 1)
				assert.Equal(t, 10 - i + 10 * j, node.Val)
				assert.True(t, rb.isRBTree())
			}
		}

		assert.Equal(t, int64(0), rb.Size())
		assert.Equal(t, (*Node[int,int])((*Node[int,int])(nil)), rb.PopLast())
		assert.True(t, rb.isRBTree())
	}
}

func BenchmarkAddMul(b *testing.B) {
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

func BenchmarkMySet1(b *testing.B) {
	rb := New[int, int]()
	for i := 0; i < b.N; i ++ {
		rb.Set(i, i + 1)
	}
}
