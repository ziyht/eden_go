package erb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/rand"
)

func TestBasic(t *testing.T) {

	root := New[int32, int32]()

	var keys []int32

	for i := 0; i < 10000; i++ {
		keys = append(keys, rand.Int31())

		root.Add(keys[i], keys[i])
	}
	assert.Equal(t, int64(10000), root.Size())

	iter := root.First()
	prev := iter.Val
	for i := 1; iter != nil; iter = iter.Next() {
		assert.True(t, prev <= iter.Val)
		//t.Logf("%d: %d", i, iter.Val)
		prev = iter.Val
		i++
	}

	iter = root.Last()
	next := iter.Val
	for i := 1; iter != nil; iter = iter.Prev() {
		assert.True(t, iter.Val <= next)
		t.Logf("%d: %d", i, iter.Val)
		next = iter.Val
		i++
	}

	for i := 0 ; i < 10000; i++ {
		assert.True(t, keys[i] == root.Val(keys[i]))
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
		rb.Add((i / 100) + 1, i + 1, true)
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
		rb.Add((i / 10) + 1, i + 1, true)
	}

	for i := 0; i < 10; i ++ {
		node := rb.PopFirst()
		assert.Equal(t, node.Val, i + 1)
	}

	for i := 0; i < 10; i ++ {
		node := rb.PopLast()
		assert.Equal(t, node.Val, 100 - i)
	}
}

func TestPopFirst(t *testing.T) {

	rb := New[int, int]()
	cnt := 100
	for i := 0; i < cnt; i ++ {
		rb.Add((i / 10) + 1, i + 1, true)
	}
	// iter := rb.First()
	// for ; iter != nil; iter = iter.Next() {
	// 	t.Logf("%d: %d", iter.Key(), iter.Val)
	// }

	{
		for i := 0; i < 10; i ++ {
			node := rb.PopKeyFirst(i + 1)
			t.Logf("%d: %d", node.Key(), node.Val)
		}
		iter := rb.First()
		for ; iter != nil; iter = iter.Next() {
			t.Logf("%d: %d", iter.Key(), iter.Val)
		}
	}

	{
		for i := 0; i < 10; i ++ {
			node := rb.PopKeyFirst(i + 1)
			t.Logf("%d: %d", node.Key(), node.Val)
		}
		iter := rb.First()
		for ; iter != nil; iter = iter.Next() {
			t.Logf("%d: %d", iter.Key(), iter.Val)
		}
	}
}