package gutil

type T any

var cache map[int]*Node[T]

type Node[T any] struct {
	Deadline int64
	Key      T
}
type HeapTimeout []Node[T]

func (ht HeapTimeout) Len() int {
	return len(ht)
}

func (ht HeapTimeout) Less(i, j int) bool {
	return ht[i].Deadline < ht[j].Deadline
}

func (ht HeapTimeout) Swap(i, j int) {
	ht[i], ht[j] = ht[j], ht[i]
}

func (ht HeapTimeout) Push(x any) {

}
