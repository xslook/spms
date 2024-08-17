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

	widx atomic.Int64
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
	for {
		n := rb.widx.Load()
		if !rb.widx.CompareAndSwap(n, n+1) {
			continue
		}
		rb.list[n%int64(rb.size)] = v
		break
	}
	return nil
}

type Consumer[T any] struct {
	rb   *Buffer[T]
	ridx atomic.Int64
}

func NewConsumer[T any](rb *Buffer[T]) (*Consumer[T], error) {
	if rb == nil {
		return nil, ErrNilBuffer
	}
	c := &Consumer[T]{
		rb: rb,
	}
	ridx := rb.widx.Load() - 1
	if ridx < 0 {
		ridx = 0
	}
	c.ridx.Store(ridx)
	return c, nil
}

// Consume the buffer
// It's a non-block method, e.g. it will return error while buffer has no more
// new items.
func (c *Consumer[T]) Consume() (T, error) {
	var t T
	if c == nil {
		return t, ErrNilBuffer
	}
	for {
		r := c.ridx.Load()
		w := c.rb.widx.Load()
		if r >= w {
			return t, ErrEmptyBuffer
		}
		if r < w-int64(c.rb.size) {
			if !c.ridx.CompareAndSwap(r, w-1) {
				continue
			}
			r = w - 1
		}
		if !c.ridx.CompareAndSwap(r, r+1) {
			continue
		}
		t = c.rb.list[r%int64(c.rb.size)]
		break
	}
	return t, nil
}
