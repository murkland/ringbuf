package ringbuf

import "fmt"

type RingBuf[T any] struct {
	items      []T
	readIndex  int
	writeIndex int
	isFull     bool
}

func New[T any](capacity int) *RingBuf[T] {
	return &RingBuf[T]{items: make([]T, capacity)}
}

func (r *RingBuf[T]) Peek(items []T, offset int) {
	i := (r.readIndex + offset) % len(r.items)

	n := r.Used()
	if n < len(items) {
		panic(fmt.Sprintf("peek too large: %d < %d", n, len(items)))
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
}

func (r *RingBuf[T]) Push(items []T) {
	if len(items) == 0 {
		return
	}

	if r.Free() < len(items) {
		panic(fmt.Sprintf("push too large: %d < %d", r.Free(), len(items)))
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

	return
}

func (r *RingBuf[T]) Advance(n int) {
	if n == 0 {
		return
	}

	if n > r.Used() {
		panic(fmt.Sprintf("advance too large: %d < %d", r.Used(), n))
	}
	r.readIndex = (r.readIndex + n) % len(r.items)
	if r.readIndex == r.writeIndex {
		r.isFull = false
	}
}

func (r *RingBuf[T]) Pop(items []T, offset int) {
	r.Peek(items, offset)
	r.Advance(len(items))
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
