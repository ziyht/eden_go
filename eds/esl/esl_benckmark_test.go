package esl

import (
	// "math"
	"sync"
	"testing"

	"github.com/huandu/skiplist"
	"github.com/zhangyunhao116/skipset"
	"github.com/ziyht/eden_go/erand"
	cst "golang.org/x/exp/constraints"
)

const (
	// initsize = 1 << 10 // for `contains` `1Remove9Add90Contains` `1Range9Remove90Add900Contains`
	// randN    = math.MaxUint32
)

type anyskipset[T any] interface {
	Add(v T) bool
	Remove(v T) bool
	Contains(v T) bool
	Range(f func(v T) bool)
	Len() int
}

// go test -benchmem -run=^$ -bench ^BenchmarkInt64$ github.com/ziyht/eden_go/eds/esl -v -cpu=16
// go test -benchmem -run=^$ -bench ^BenchmarkInt64$ github.com/ziyht/eden_go/eds/esl -v -cpu=16 -benchtime=1000000x
func BenchmarkInt64(b *testing.B) {
	var all []benchTask[int64]

	all = append(all, benchTask[int64]{
		name: "skipset", New: func() anyskipset[int64] {
			return skipset.New[int64]()
		}})
	all = append(all, benchTask[int64]{
		name: "esl", New: func() anyskipset[int64] {
			return newBenckESL[int64]()
		}})
	all = append(all, benchTask[int64]{
		name: "skipset(func)", New: func() anyskipset[int64] {
			return skipset.NewFunc(func(a, b int64) bool {
				return a < b
			})
		}})
	all = append(all, benchTask[int64]{
		name: "sync.Map", New: func() anyskipset[int64] {
			return new(anySyncMap[int64])
		}})
	// all = append(all, benchTask[int64]{
	// 	name: "huandu/skiplist", New: func() anyskipset[int64] {
	// 		return newHuanduSkipListInt64[int64]()
	// 	}})

	b.SetParallelism(1)

	rng := erand.Int63

	benchAdd(b, rng, all)
	bench30Add70Contains(b, rng, all)
	bench1Remove9Add90Contains(b, rng, all)
	bench1Range9Remove90Add900Contains(b, rng, all)
}

func benchAdd[T any](b *testing.B, rng func() T, benchTasks []benchTask[T]) {
	for _, v := range benchTasks {
		b.Run("Add/"+v.name, func(b *testing.B) {
			s := v.New()
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					s.Add(rng())
				}
			})
		})
	}
}

func bench30Add70Contains[T any](b *testing.B, rng func() T, benchTasks []benchTask[T]) {
	for _, v := range benchTasks {
		b.Run("30Add70Contains/"+v.name, func(b *testing.B) {
			s := v.New()
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					u := erand.Uint32n(10)
					if u < 3 {
						s.Add(rng())
					} else {
						s.Contains(rng())
					}
				}
			})
		})
	}
}

func bench1Remove9Add90Contains[T any](b *testing.B, rng func() T, benchTasks []benchTask[T]) {
	for _, v := range benchTasks {
		b.Run("1Remove9Add90Contains/"+v.name, func(b *testing.B) {
			s := v.New()
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					u := erand.Uint32n(100)
					if u < 9 {
						s.Add(rng())
					} else if u == 10 {
						s.Remove(rng())
					} else {
						s.Contains(rng())
					}
				}
			})
		})
	}
}

func bench1Range9Remove90Add900Contains[T any](b *testing.B, rng func() T, benchTasks []benchTask[T]) {
	for _, v := range benchTasks {
		b.Run("1Range9Remove90Add900Contains/"+v.name, func(b *testing.B) {
			s := v.New()
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					u := erand.Uint32n(1000)
					if u == 0 {
						s.Range(func(score T) bool {
							return true
						})
					} else if u > 10 && u < 20 {
						s.Remove(rng())
					} else if u >= 100 && u < 190 {
						s.Add(rng())
					} else {
						s.Contains(rng())
					}
				}
			})
		})
	}
}

type benchTask[T any] struct {
	name string
	New  func() anyskipset[T]
}


type benckESL[T cst.Ordered] struct {
	*ESL[T, int64]
}

func newBenckESL[T cst.Ordered] () *benckESL[T] {
	return &benckESL[T]{New[T, int64]()}
}

func (m *benckESL[T]) Add(x T) bool {
	return m.ESL.Add(x, 0)
}

func (m *benckESL[T]) Contains(x T) bool {
	return m.ESL.Contains(x)
}

func (m *benckESL[T]) Remove(x T) bool {
	return m.ESL.Remove(x)
}

func (m *benckESL[T]) Range(f func(key T) bool) {
  m.ESL.Range(func(key T, val int64)bool { return f(key)})
}

func (m *benckESL[T]) RangeFrom(start T, f func(value T) bool) {
  m.ESL.RangeFrom(start, func(key T, val int64)bool { return f(key)})
}

func (m *benckESL[T]) Len() int {
	return int(m.ESL.Len())
}

//
// huanduSkipList
// 
type huanduSkipList [T any] struct {
	*skiplist.SkipList
}

func newHuanduSkipListInt64[T any] () *huanduSkipList[T] {
	return &huanduSkipList[T]{skiplist.New(skiplist.Int64)}
}

func (m *huanduSkipList[T]) Add(x T) bool {
	m.SkipList.Set(x, 0)
	return true
}

func (m *huanduSkipList[T]) Contains(x T) bool {
	ok := m.SkipList.Find(x)
	return ok != nil
}

func (m *huanduSkipList[T]) Remove(x T) bool {
	return m.SkipList.Remove(x) != nil
}

func (m *huanduSkipList[T]) Range(f func(value T) bool) {
	panic("TODO")
}

func (m *huanduSkipList[T]) RangeFrom(start T, f func(value T) bool) {
	panic("TODO")
}

func (m *huanduSkipList[T]) Len() int {
	return m.SkipList.Len()
}

//
// anySyncMap
// 
type anySyncMap[T any] struct {
	data sync.Map
}

func (m *anySyncMap[T]) Add(x T) bool {
	m.data.Store(x, struct{}{})
	return true
}

func (m *anySyncMap[T]) Contains(x T) bool {
	_, ok := m.data.Load(x)
	return ok
}

func (m *anySyncMap[T]) Remove(x T) bool {
	m.data.Delete(x)
	return true
}

func (m *anySyncMap[T]) Range(f func(value T) bool) {
	m.data.Range(func(key, _ any) bool {
		return !f(key.(T))
	})
}

func (m *anySyncMap[T]) RangeFrom(start T, f func(value T) bool) {
	panic("TODO")
}

func (m *anySyncMap[T]) Len() int {
	var i int
	m.data.Range(func(_, _ any) bool {
		i++
		return true
	})
	return i
}
