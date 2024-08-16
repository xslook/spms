package spmcb

import (
	"fmt"
	"sync/atomic"
)

var (
	ErrEmptyBuffer = fmt.Errorf("Buffer is empty")
	ErrNilBuffer   = fmt.Errorf("Buffer is nil")
)

type Buffer[T any] struct {
	list []T

	idx  atomic.Int64
	size int
}

func New[T any](size int) (*Buffer[T], error) {
	if size <= 1 {
		return nil, ErrNilBuffer
	}
	b := &Buffer[T]{
		list: make([]T, size),
		size: size,
	}
	return b, nil
}

func (rb *Buffer[T]) Produce(v T) error {
	if rb == nil {
		return ErrNilBuffer
	}
	n := rb.idx.Add(1)
	rb.list[n%int64(rb.size)] = v
	return nil
}

type Consumer[T any] struct {
	rb  *Buffer[T]
	idx atomic.Int64
}

func NewConsumer[T any](rb *Buffer[T]) (*Consumer[T], error) {
	if rb == nil {
		return nil, ErrNilBuffer
	}
	c := &Consumer[T]{
		rb: rb,
	}
	c.idx.Store(rb.idx.Load())
	return c, nil
}

func (c *Consumer[T]) Consume() (T, error) {
	var t T
	if c == nil {
		return t, ErrNilBuffer
	}
	for {
		r := c.idx.Load()
		w := c.rb.idx.Load()
		if r >= w {
			break
		}
		if r < w-int64(c.rb.size) {
			if !c.idx.CompareAndSwap(r, w-1) {
				continue
			}
			r = w - 1
		}
		if !c.idx.CompareAndSwap(r, r+1) {
			continue
		}
		t = c.rb.list[(r+1)%int64(c.rb.size)]
		break
	}
	return t, nil
}
