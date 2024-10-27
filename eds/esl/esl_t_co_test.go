package esl

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)


func prepare[K int, V int](splitIn int, scale int) (out [][]elem[K, V]) {

	for i := 0; i < splitIn; i++ {
		out = append(out, make([]elem[K, V], 0))
	}

	for i := 0; i < scale;  {
		for j := 0; j < splitIn; j++ {
			out[j] = append(out[j], elem[K, V]{key: K(i), Val: V(i)})
			i++
			if i >= scale {
				break
			}
		}
	}

	return
}

func addToList[K int, V int](sl *ESL[K, V], elems []elem[K, V]) (okCnt int) {
	for _, e := range elems {
		if sl.Add(e.key, e.Val) {okCnt++}
	}
	return
}

func delFromList[K int, V int](sl *ESL[K, V], elems []elem[K, V]) (okCnt int) {
	for _, e := range elems {
		if sl.Remove(e.key) {okCnt++}
	}
	return
}

func TestCoSaveTestAddDel(t *testing.T) {
	sl := New[int, int]()
	var wg sync.WaitGroup

	scale   := 10000000
	threads := 10

	inputs := prepare[int, int](threads, scale)

	// test add in multiThread
	wg.Add(threads)
	for i := 0; i < threads; i++ {
		a := i
		go func() {
			addToList(sl, inputs[a])
			wg.Done()
		}()
	}
	wg.Wait()

	assert.Equal(t, int64(scale), sl.Len())
	for i := 0; i < scale; i++ {
		if ! assert.Equal(t, i, sl.Val(i)) {
			t.FailNow()
		}
	}

	// test del in multiThread
	wg.Add(threads)
	for i := 0; i < threads; i++ {
		a := i
		go func() {
			delFromList(sl, inputs[a])
			wg.Done()
		}()
	}
	wg.Wait()
	assert.Equal(t, int64(0), sl.Len())
	assert.Nil(t, sl.First())
}