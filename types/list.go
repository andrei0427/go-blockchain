package types

import (
	"fmt"
	"reflect"
)

type List[T any] struct {
	Data []T
}

func NewList[T any]() *List[T] {
	return &List[T]{
		Data: []T{},
	}
}

func (l *List[T]) Get(idx int) T {
	if idx < 0 || idx > l.Len()-1 {
		panic(fmt.Errorf("index (%d) out of bounds (%d)", idx, l.Len()-1))
	}

	return l.Data[idx]
}

func (l *List[T]) Insert(d T) {
	l.Data = append(l.Data, d)
}

func (l *List[T]) Clear() {
	l.Data = []T{}
}

func (l *List[T]) Len() int {
	return len(l.Data)
}

func (l *List[T]) IndexOf(v T) int {
	for i := 0; i < l.Len(); i++ {
		if reflect.DeepEqual(l.Data[i], v) {
			return i
		}
	}

	return -1
}

func (l *List[T]) Remove(v T) {
	idx := l.IndexOf(v)
	if idx == -1 {
		return
	}
	l.Pop(idx)
}

func (l *List[T]) Pop(idx int) {
	l.Data = append(l.Data[:idx], l.Data[idx+1:]...)
}

func (l *List[T]) Contains(v T) bool {
	return l.IndexOf(v) != -1
}

func (l *List[T]) Last() T {
	return l.Get(l.Len() - 1)
}
