package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewList(t *testing.T) {
	l := NewList[int]()

	assert.Equal(t, l.Data, []int{})
}

func TestListClear(t *testing.T) {
	l := NewList[int]()
	n := 100

	for i := 0; i < n; i++ {
		l.Insert(i)
	}

	assert.Equal(t, l.Len(), n)

	l.Clear()
	assert.Equal(t, l.Len(), 0)
}

func TestListContains(t *testing.T) {
	l := NewList[int]()
	n := 100

	for i := 0; i < n; i++ {
		l.Insert(i)
		assert.True(t, l.Contains(i))
	}
}

func TestListGetIndex(t *testing.T) {
	l := NewList[string]()
	n := 100

	for i := 0; i < n; i++ {
		data := (fmt.Sprintf("foo_%d", i))
		l.Insert(data)
		assert.Equal(t, l.IndexOf(data), i)
	}

	assert.Equal(t, l.IndexOf("bar"), -1)
}

func TestListRemove(t *testing.T) {
	l := NewList[string]()
	n := 100

	for i := 0; i < n; i++ {
		data := fmt.Sprintf("foo_%d", i)
		l.Insert(data)
		l.Remove(data)
		assert.False(t, l.Contains(data))
	}

	assert.Equal(t, l.Len(), 0)
}

func TestListGet(t *testing.T) {
	l := NewList[string]()
	n := 100

	for i := 0; i < n; i++ {
		data := fmt.Sprintf("foo_%d", i)
		l.Insert(data)

		assert.True(t, l.Contains(data))
		assert.Equal(t, l.Get(i), data)
	}
}

func TestListPop(t *testing.T) {
	l := NewList[string]()
	l.Insert("foo")
	l.Insert("bar")
	l.Insert("fizz")

	l.Pop(1)

	assert.Equal(t, "fizz", l.Get(1))
}

func TestListAdd(t *testing.T) {

	l := NewList[string]()
	n := 100

	for i := 0; i < n; i++ {
		data := fmt.Sprintf("foo_%d", i)
		l.Insert(data)
	}

	assert.Equal(t, l.Len(), n)
}

func TestListLast(t *testing.T) {

	l := NewList[string]()
	n := 100

	for i := 0; i <= n; i++ {
		data := fmt.Sprintf("foo_%d", i)
		l.Insert(data)
	}

	assert.Equal(t, l.Last(), fmt.Sprintf("foo_%d", n))
}
