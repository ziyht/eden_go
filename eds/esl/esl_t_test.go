package esl

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {

	sl := New[int, int]()

	assert.True(t, sl.Add(1, 1))
	assert.False(t, sl.Add(1, 2))
	assert.False(t, sl.Add(1, 3))

	assert.Equal(t, int64(1), sl.Size())

	for i := 0; i < 10000; i++ {
		sl.Add(i, i)
	}
	assert.Equal(t, int64(10000), sl.Size())

	for i := 0; i < 10000; i++ {
		sl.Add(i, i)
	}
	assert.Equal(t, int64(10000), sl.Size())
}

func TestBasic(t *testing.T) {
	sl := New[int, int]()
	sl.Add(1, 1)
	sl.Add(2, 2)
	sl.Add(3, 3)	
	assert.Equal(t, int64(3), sl.Len())
	assert.True(t, sl.Contains(1))
	assert.True(t, sl.Contains(2))
	assert.True(t, sl.Contains(3))
	assert.False(t, sl.Contains(4))
	sl.Del(2)
	assert.Equal(t, int64(2), sl.Len())
	assert.False(t, sl.Contains(2))

	assert.Equal(t, 1, sl.Val(1))
	assert.Equal(t, 3, sl.Val(3))
}

func TestBasic2(t *testing.T) {
	sl := New[int32, int32]()

	var keys []int32

	for i := 0; i < 10000; i++ {
		keys = append(keys, rand.Int31())

		sl.Add(keys[i], keys[i])
	}
	assert.Equal(t, int64(10000), sl.Size())


	iter := sl.First()
	prev := iter.Val
	for i := 1; iter != nil; iter = iter.Next() {
		assert.True(t, prev <= iter.Val)
		//t.Logf("%d: <%d, %d>", i, iter.Key(), iter.Val)
		prev = iter.Val
		i++
	}

	for i := 0 ; i < 10000; i++ {
		assert.Equal(t, keys[i], sl.Val(keys[i]))
	}

}

func TestSet(t *testing.T) {

	sl := New[int32, int]()

	prev, ok := sl.Set(1, 1)
	assert.Equal(t, 0, prev)
	assert.Equal(t, false, ok)
	assert.Equal(t, int64(1), sl.Size())
	assert.Equal(t, 1, sl.Val(1))

	prev, ok = sl.Set(1, 2)
	assert.Equal(t, 1, prev)
	assert.Equal(t, true, ok)
	assert.Equal(t, int64(1), sl.Size())
	assert.Equal(t, 2, sl.Val(1))

	sl.Set(2, 1)
	assert.Equal(t, int64(2), sl.Size())
	assert.Equal(t, 2, sl.Val(1))
	assert.Equal(t, 1, sl.Val(2))

	sl.Set(2, 2)
	assert.Equal(t, int64(2), sl.Size())
	assert.Equal(t, 2, sl.Val(1))
	assert.Equal(t, 2, sl.Val(2))
}

func TestPop(t *testing.T) {

	rb := New[int, int]()
	cnt := 100
	for i := 0; i < cnt; i ++ {
		rb.Add(i, i + 1)
	}

	for i := 0; i < 50; i ++ {
		node := rb.PopFirst()
		assert.Equal(t, node.Val, i + 1)
	}

	assert.Equal(t, int64(50), rb.Size())
}

func TestRange(t *testing.T) {
	sl := New[int, int]()

	sl.Range(func(k int, v int) bool {return true})

	for i := 0; i < 10000; i++ {
		sl.Add(i, i)
	}

	i := 0
	sl.Range(func(k int, v int) bool {
		assert.Equal(t, i, k)
		assert.Equal(t, k, v)
		i++
		return true
	})


	i = 0
	sl.Range(func(k int, v int) bool {
		assert.Equal(t, i, k)
		assert.Equal(t, k, v)
		i++
		return true
	}, 1000)
	assert.Equal(t, 1000, i)
}

func TestRangFrom(t *testing.T) {
	sl := New[int, int]()

	for i := 0; i < 100; i++ {
		sl.Add(i, i)
	}
	sl.Dels(10, 20)

	i := 21
	cnt := 0
	sl.RangeFrom(20, func(k int, v int) bool {
		assert.Equal(t, i, k)
		assert.Equal(t, i, v)
		i++
		cnt++
		return true
	})
	assert.Equal(t, 79, cnt)

	i = 21
	cnt = 0
	sl.RangeFrom(20, func(k int, v int) bool {
		assert.Equal(t, i, k)
		assert.Equal(t, i, v)
		i++
		cnt++
		return true
	}, 5)
	assert.Equal(t, 5, cnt)

	i = 30
	cnt = 0
	sl.RangeFrom(30, func(k int, v int) bool {
		assert.Equal(t, i, k)
		assert.Equal(t, i, v)
		i++
		cnt++
		return true
	})
	assert.Equal(t, 70, cnt)
}

func TestRangFromTo(t *testing.T) {

	sl := New[int, int]()

	for i := 0; i < 100; i++ {
		sl.Add(i, i)
	}
	sl.Dels(10, 20)

	i := 11
	cnt := 0
	sl.RangeFromTo(10, 20, func(k int, v int) bool {
		assert.Equal(t, i, k)
		assert.Equal(t, i, v)
		i++
		cnt++
		return true
	})
	assert.Equal(t, 9, cnt)

	i = 11
	cnt = 0
	sl.RangeFromTo(10, 20, func(k int, v int) bool {
		assert.Equal(t, i, k)
		assert.Equal(t, i, v)
		i++
		cnt++
		return true
	}, 5)
	assert.Equal(t, 5, cnt)

	// TODO
	i = 19
	cnt = 0
	sl.RangeFromTo(20, 10, func(k int, v int) bool {
		assert.Equal(t, i, k)
		assert.Equal(t, i, v)
		i--
		cnt++
		return true
	}, 5)
	assert.Equal(t, 0, cnt)

	i = 30
	cnt = 0
	sl.RangeFromTo(30, 40, func(k int, v int) bool {
		assert.Equal(t, i, k)
		assert.Equal(t, i, v)
		i++
		cnt++
		return true
	})
	assert.Equal(t, 11, cnt)

	// TODO
	i = 40
	cnt = 0
	sl.RangeFromTo(40, 30, func(k int, v int) bool {
		assert.Equal(t, i, k)
		assert.Equal(t, i, v)
		i--
		cnt++
		return true
	})
	assert.Equal(t, 0, cnt)
}

func TestRangIn(t *testing.T) {

	sl := New[int, int]()

	for i := 0; i < 100; i++ {
		sl.Add(i, i)
	}
	sl.Dels(10, 20)

	i := 11
	cnt := 0
	sl.RangeIn(10, 20, func(k int, v int) bool {
		assert.Equal(t, i, k)
		assert.Equal(t, i, v)
		i++
		cnt++
		return true
	})
	assert.Equal(t, 9, cnt)

	i = 11
	cnt = 0
	sl.RangeIn(10, 20, func(k int, v int) bool {
		assert.Equal(t, i, k)
		assert.Equal(t, i, v)
		i++
		cnt++
		return true
	}, 5)
	assert.Equal(t, 5, cnt)

	// TODO
	i = 19
	cnt = 0
	sl.RangeIn(20, 10, func(k int, v int) bool {
		assert.Equal(t, i, k)
		assert.Equal(t, i, v)
		i--
		cnt++
		return true
	}, 5)
	assert.Equal(t, 0, cnt)

	i = 30
	cnt = 0
	sl.RangeIn(30, 40, func(k int, v int) bool {
		assert.Equal(t, i, k)
		assert.Equal(t, i, v)
		i++
		cnt++
		return true
	})
	assert.Equal(t, 10, cnt)

	// TODO
	i = 40
	cnt = 0
	sl.RangeIn(40, 30, func(k int, v int) bool {
		assert.Equal(t, i, k)
		assert.Equal(t, i, v)
		i--
		cnt++
		return true
	})
	assert.Equal(t, 0, cnt)
}

func TestKeys(t *testing.T) {

	sl := New[int, int]()

	for i := 0; i < 1000; i++ {
		sl.Add(i, i)
	}
	
	keys := sl.Keys()
	for i := 0; i < 1000; i++ {
		assert.Equal(t, i, keys[i])
	}

	vals := sl.Vals()
	for i := 0; i < 1000; i++ {
		assert.Equal(t, i, vals[i])
	}
}
