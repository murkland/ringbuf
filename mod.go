package ringbuf

import (
	"errors"
)

var (
	ErrWriteTooLarge = errors.New("write too large")
	ErrReadTooLarge  = errors.New("read too large")
)

type RingBuf[T any] struct {
	items      []T
	readIndex  int
	writeIndex int
	isFull     bool
}

func New[T any](capacity int) *RingBuf[T] {
	return &RingBuf[T]{items: make([]T, capacity)}
}

func (r *RingBuf[T]) Peek(items []T, offset int) error {
	i := (r.readIndex + offset) % len(r.items)

	n := r.Used()
	if n < len(items) {
		return ErrReadTooLarge
	}

	if n > len(items) {
		n = len(items)
	}

	if r.writeIndex > i || i+n <= len(r.items) {
		copy(items, r.items[i:i+n])
	} else {
		copy(items, r.items[i:len(r.items)])
		copy(items[len(r.items)-i:], r.items[:n-len(r.items)+i])
	}

	return nil
}

func (r *RingBuf[T]) Push(items []T) error {
	if len(items) == 0 {
		return nil
	}

	if r.Free() < len(items) {
		return ErrWriteTooLarge
	}

	if r.writeIndex >= r.readIndex {
		i := len(r.items) - r.writeIndex
		if i >= len(items) {
			copy(r.items[r.writeIndex:], items)
			r.writeIndex += len(items)
		} else {
			copy(r.items[r.writeIndex:], items[:i])
			copy(r.items, items[i:])
			r.writeIndex = len(items) - i
		}
	} else {
		copy(r.items[r.writeIndex:], items)
		r.writeIndex += len(items)
	}

	r.writeIndex = r.writeIndex % len(r.items)
	if r.writeIndex == r.readIndex {
		r.isFull = true
	}

	return nil
}

func (r *RingBuf[T]) Advance(n int) error {
	if n == 0 {
		return nil
	}

	if n > r.Used() {
		return ErrReadTooLarge
	}
	r.readIndex = (r.readIndex + n) % len(r.items)
	if r.readIndex == r.writeIndex {
		r.isFull = false
	}
	return nil
}

func (r *RingBuf[T]) Free() int {
	return len(r.items) - r.Used()
}

func (r *RingBuf[T]) Used() int {
	if r.writeIndex == r.readIndex {
		if r.isFull {
			return len(r.items)
		}
		return 0
	}

	if r.writeIndex > r.readIndex {
		return r.writeIndex - r.readIndex
	}
	return len(r.items) - r.readIndex + r.writeIndex
}
