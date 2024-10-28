package erb

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestSet(t *testing.T) {

	rb := New[int32, int]()

	prev, ok := rb.Set(1, 1)
	assert.Equal(t, 0, prev)
	assert.Equal(t, false, ok)
	assert.Equal(t, int64(1), rb.Size())
	assert.Equal(t, 1, rb.Val(1))

	prev, ok = rb.Set(1, 2)
	assert.Equal(t, 1, prev)
	assert.Equal(t, true, ok)
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

func TestRange(t *testing.T) {
	rb := New[int, int]()

	rb.Range(func(k int, v int) bool {return true})
	rb.RangeRev(func(k int, v int) bool {return true})

	for i := 0; i < 10000; i++ {
		rb.MulAdd(i, i)
	}

	i := 0
	rb.Range(func(k int, v int) bool {
		assert.Equal(t, i, k)
		assert.Equal(t, k, v)
		i++
		return true
	})

	i = 9999
	rb.RangeRev(func(k int, v int) bool {
		assert.Equal(t, i, k)
		assert.Equal(t, k, v)
		i--
		return true
	})

	i = 0
	rb.Range(func(k int, v int) bool {
		assert.Equal(t, i, k)
		assert.Equal(t, k, v)
		i++
		return true
	}, 1000)
	assert.Equal(t, 1000, i)

	i = 9999
	rb.RangeRev(func(k int, v int) bool {
		assert.Equal(t, i, k)
		assert.Equal(t, k, v)
		i--
		return true
	}, 1000)
	assert.Equal(t, 8999, i)
}

func TestRangFrom(t *testing.T) {
	rb := New[int, int]()

	for i := 0; i < 100; i++ {
		rb.Add(i, i)
	}
	rb.Dels(10, 20)

	i := 21
	cnt := 0
	rb.RangeFrom(20, func(k int, v int) bool {
		assert.Equal(t, i, k)
		assert.Equal(t, i, v)
		i++
		cnt++
		return true
	})
	assert.Equal(t, 79, cnt)

	i = 21
	cnt = 0
	rb.RangeFrom(20, func(k int, v int) bool {
		assert.Equal(t, i, k)
		assert.Equal(t, i, v)
		i++
		cnt++
		return true
	}, 5)
	assert.Equal(t, 5, cnt)

	i = 30
	cnt = 0
	rb.RangeFrom(30, func(k int, v int) bool {
		assert.Equal(t, i, k)
		assert.Equal(t, i, v)
		i++
		cnt++
		return true
	})
	assert.Equal(t, 70, cnt)
}

func TestRangFromTo(t *testing.T) {

	rb := New[int, int]()

	for i := 0; i < 100; i++ {
		rb.Add(i, i)
	}
	rb.Dels(10, 20)

	i := 11
	cnt := 0
	rb.RangeFromTo(10, 20, func(k int, v int) bool {
		assert.Equal(t, i, k)
		assert.Equal(t, i, v)
		i++
		cnt++
		return true
	})
	assert.Equal(t, 9, cnt)

	i = 11
	cnt = 0
	rb.RangeFromTo(10, 19, func(k int, v int) bool {
		assert.Equal(t, i, k)
		assert.Equal(t, i, v)
		i++
		cnt++
		return true
	})
	assert.Equal(t, 9, cnt)

	i = 11
	cnt = 0
	rb.RangeFromTo(10, 20, func(k int, v int) bool {
		assert.Equal(t, i, k)
		assert.Equal(t, i, v)
		i++
		cnt++
		return true
	}, 5)
	assert.Equal(t, 5, cnt)

	i = 19
	cnt = 0
	rb.RangeFromTo(20, 10, func(k int, v int) bool {
		assert.Equal(t, i, k)
		assert.Equal(t, i, v)
		i--
		cnt++
		return true
	}, 5)
	assert.Equal(t, 5, cnt)

	i = 30
	cnt = 0
	rb.RangeFromTo(30, 40, func(k int, v int) bool {
		assert.Equal(t, i, k)
		assert.Equal(t, i, v)
		i++
		cnt++
		return true
	})
	assert.Equal(t, 11, cnt)

	i = 40
	cnt = 0
	rb.RangeFromTo(40, 30, func(k int, v int) bool {
		assert.Equal(t, i, k)
		assert.Equal(t, i, v)
		i--
		cnt++
		return true
	})
	assert.Equal(t, 11, cnt)
}

func TestRangIn(t *testing.T) {

	rb := New[int, int]()

	for i := 0; i < 100; i++ {
		rb.Add(i, i)
	}
	rb.Dels(10, 20)

	i := 11
	cnt := 0
	rb.RangeIn(10, 20, func(k int, v int) bool {
		assert.Equal(t, i, k)
		assert.Equal(t, i, v)
		i++
		cnt++
		return true
	})
	assert.Equal(t, 9, cnt)

	i = 11
	cnt = 0
	rb.RangeIn(10, 19, func(k int, v int) bool {
		assert.Equal(t, i, k)
		assert.Equal(t, i, v)
		i++
		cnt++
		return true
	})
	assert.Equal(t, 8, cnt)

	i = 11
	cnt = 0
	rb.RangeIn(10, 20, func(k int, v int) bool {
		assert.Equal(t, i, k)
		assert.Equal(t, i, v)
		i++
		cnt++
		return true
	}, 5)
	assert.Equal(t, 5, cnt)

	i = 19
	cnt = 0
	rb.RangeIn(20, 10, func(k int, v int) bool {
		assert.Equal(t, i, k)
		assert.Equal(t, i, v)
		i--
		cnt++
		return true
	}, 5)
	assert.Equal(t, 5, cnt)

	i = 30
	cnt = 0
	rb.RangeIn(30, 40, func(k int, v int) bool {
		assert.Equal(t, i, k)
		assert.Equal(t, i, v)
		i++
		cnt++
		return true
	})
	assert.Equal(t, 10, cnt)

	i = 40
	cnt = 0
	rb.RangeIn(40, 30, func(k int, v int) bool {
		assert.Equal(t, i, k)
		assert.Equal(t, i, v)
		i--
		cnt++
		return true
	})
	assert.Equal(t, 10, cnt)
}