package esl

//
// rebuild from https://github.com/zhangyunhao116/skipset
// 

const (
	maxLevel            = 16
	p                   = 0.25
	defaultHighestLevel = 3
	VER                 = "0.0.1"
)

func New[K ordered, V any]() *ESL[K, V] {
	return &ESL[K, V]{
		header:       newHead[K, V](maxLevel),
		maxLevel: defaultHighestLevel,
	}
}