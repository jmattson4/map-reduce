package mr

import (
	"sync"
)

type HandlerType int

const (
	GetTask HandlerType = iota
	CompleteMap
	CompleteReduce
)

type TaskType int

const (
	Map TaskType = iota
	Reduce
)

type TaskStatus int

const (
	NotStarted TaskStatus = iota
	Started
	Finished
)

type ThreadSafeMap[K comparable, V any] struct {
	sync.RWMutex
	m map[K]V
}

func (tsM *ThreadSafeMap[K, V]) Put(key K, value V) {
	tsM.Lock()
	defer tsM.Unlock()
	tsM.m[key] = value
}

func (tsM *ThreadSafeMap[K, V]) Get(key K) (val V, ok bool) {
	tsM.RLock()
	defer tsM.RUnlock()
	val, ok = tsM.m[key]
	return
}

func NewThreadSafeMap[K comparable, V any]() *ThreadSafeMap[K, V] {
	return &ThreadSafeMap[K, V]{
		m: make(map[K]V),
	}
}

//
// Map functions return a slice of KeyValue.
//
type KeyValue struct {
	Key   string
	Value string
}

// for sorting by key.
type ByKey []KeyValue

// for sorting by key.
func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }
