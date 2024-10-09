package erb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestMulAdd(t *testing.T) {

	rb := New[int, int]()

	for i := 0; i < 10000; i++ {
		rb.MulAdd(i, i)
	}
	assert.Equal(t, int64(10000), rb.Size())

	for i := 0; i < 10000; i++ {
		rb.MulAdd(i, i + 1)
	}
	assert.Equal(t, int64(20000), rb.Size())

	for i := 0; i < 10000; i++ {
		v := rb.ValFirst(i)
		assert.Equal(t, i, v)
	}
}

func TestSetAll(t *testing.T) {
	rb := New[int, int]()

	for i := 0; i < 1000; i++ {
		rb.MulAdd(i/10, i)
	}

	cnt := rb.SetAll(-1, 1)
	assert.Equal(t, int64(0), cnt)

	cnt = rb.SetAll(0, 888)
	assert.Equal(t, int64(10), cnt)
	rb.RangeFrom(0, 0, func(k, v int) bool{
		assert.Equal(t, 0, k)
		assert.Equal(t, 888, v)
		return true
	})

	i := 10
	rb.RangeFrom(1, 1, func(k, v int) bool{
		assert.Equal(t, 1, k)
		assert.Equal(t, i, v)
		i++
		return true
	})
}	

func TestValFirstLast(t *testing.T) {
	rb := New[int, int]()

	for i := 0; i < 10000; i++ {
		rb.MulAdd(1, i)
	}
	assert.Equal(t, int64(10000), rb.Size())

	for i := 0; i < 10000; i++ {
		assert.Equal(t, i, rb.ValFirst(1))
		assert.Equal(t, 9999, rb.ValLast(1))

		rb.PopFirst()
		assert.True(t, rb.isRBTree())
	}
	assert.Equal(t, int64(0), rb.Size())

	for i := 0; i < 10000; i++ {
		rb.MulAdd(1, i)
	}
	assert.Equal(t, int64(10000), rb.Size())

	for i := 0; i < 10000; i++ {
		assert.Equal(t, 0, rb.ValFirst(1))
		assert.Equal(t, 9999 - i, rb.ValLast(1))

		rb.PopLast()
		assert.True(t, rb.isRBTree())
	}
	assert.Equal(t, int64(0), rb.Size())
}

func TestDelFirst(t *testing.T) {
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

func TestDelLast(t *testing.T) {
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

func TestDelFirstLast1(t *testing.T) {
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

func TestDelFirstLast2(t *testing.T) {
	rb := New[int, int]()

	//
	for i := 0; i < 10000; i++ {
		rb.MulAdd(1, i)
	}
	assert.Equal(t, int64(10000), rb.Size())

	for i := 0; i < 10000; i++ {
		assert.Equal(t, i, rb.ValFirst(1))
		assert.Equal(t, 9999, rb.ValLast(1))

		n := rb.DelFirst(1)
		assert.Equal(t, i, n.Val)
		assert.True(t, rb.isRBTree())
	}
	assert.Equal(t, int64(0), rb.Size())

	//
	for i := 0; i < 10000; i++ {
		rb.MulAdd(1, i)
	}
	assert.Equal(t, int64(10000), rb.Size())

	for i := 0; i < 10000; i++ {
		assert.Equal(t, 0, rb.ValFirst(1))
		assert.Equal(t, 9999 - i, rb.ValLast(1))

		rb.DelLast(1)
		assert.True(t, rb.isRBTree())
	}
	assert.Equal(t, int64(0), rb.Size())

	//
	for i := 0; i < 10000; i++ {
		rb.MulAdd(1, i)
		rb.MulAdd(2, i)
		rb.MulAdd(3, i)
	}

	for i := 0; i < 10000; i++ {
		assert.Equal(t, i, rb.ValFirst(2))
		assert.Equal(t, 9999, rb.ValLast(2))

		rb.DelFirst(2)
		assert.True(t, rb.isRBTree())
	}
	assert.Equal(t, int64(20000), rb.Size())

	for i := 0; i < 10000; i++ {
		rb.MulAdd(2, i)
	}
	assert.Equal(t, int64(30000), rb.Size())

	for i := 0; i < 10000; i++ {
		assert.Equal(t, 0, rb.ValFirst(2))
		assert.Equal(t, 9999 - i, rb.ValLast(2))

		rb.DelLast(2)
		assert.True(t, rb.isRBTree())
	}
	assert.Equal(t, int64(20000), rb.Size())
}

func TestDelAll(t *testing.T) {

	rb := New[int, int]()

	for i := 0; i < 10000; i++ {
		rb.MulAdd(i / 10, i)
	}

	for i := 0; i < 1000; i++ {
		cnt := rb.DelAll(i)
		assert.Equal(t, int64(10), cnt)
		assert.Equal(t, int64(10000 - 10 * (i + 1)), rb.Size())
	}
}

func TestRangFromMulti(t *testing.T) {

	rb := New[int, int]()

	for i := 0; i < 10000; i++ {
		rb.MulAdd(i / 10, i)
	}

	for i := 0; i < 1000; i++ {
		j := 0
		rb.RangeFrom(i, i, func(k int, v int) bool{
			assert.Equal(t, i, k)
			assert.Equal(t, i*10+j, v)
			j++
			return true
		})
	}

}
