package ringbuf_test

import (
	"reflect"
	"testing"

	"github.com/undernet/ringbuf"
)

func TestRingbuf_PushOutOfBounds(t *testing.T) {
	rbuf := ringbuf.New[int](10)
	if err := rbuf.Push([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}); err != ringbuf.ErrWriteTooLarge {
		t.Error("expected write too large")
	}
}

func TestRingbuf_PeekOutOfBounds(t *testing.T) {
	rbuf := ringbuf.New[int](10)
	buf := make([]int, 2)
	if err := rbuf.Peek(buf, 0); err != ringbuf.ErrReadTooLarge {
		t.Error("expected read to large")
	}
}
func TestRingbuf_Push(t *testing.T) {
	rbuf := ringbuf.New[int](10)
	if err := rbuf.Push([]int{1, 2, 3, 4, 5, 6}); err != nil {
		t.Error("expected not out of bounds")
	}
	if n := rbuf.Free(); n != 4 {
		t.Errorf("expected to have 4 free, got %d", n)
	}
}

func TestRingbuf_PushThenPush(t *testing.T) {
	rbuf := ringbuf.New[int](10)
	if err := rbuf.Push([]int{1, 2, 3, 4, 5, 6}); err != nil {
		t.Error("expected not out of bounds")
	}
	if err := rbuf.Push([]int{4, 5, 6, 7}); err != nil {
		t.Error("expected not out of bounds")
	}
	if n := rbuf.Free(); n != 0 {
		t.Errorf("expected to have 0 free, got %d", n)
	}
}

func TestRingbuf_PushThenAdvance(t *testing.T) {
	rbuf := ringbuf.New[int](10)
	if err := rbuf.Push([]int{1, 2, 3, 4, 5, 6}); err != nil {
		t.Error("expected not out of bounds")
	}
	if n := rbuf.Free(); n != 4 {
		t.Errorf("expected to have 4 free, got %d", n)
	}
	if err := rbuf.Advance(2); err != nil {
		t.Error("expected not out of bounds")
	}
	if n := rbuf.Free(); n != 6 {
		t.Errorf("expected to have 6 free, got %d", n)
	}
}

func TestRingbuf_PushThenAdvanceThenPeek(t *testing.T) {
	rbuf := ringbuf.New[int](10)
	if err := rbuf.Push([]int{1, 2, 3, 4, 5, 6}); err != nil {
		t.Error("expected not out of bounds")
	}
	buf := make([]int, 2)
	if err := rbuf.Peek(buf, 0); err != nil {
		t.Error("expected not out of bounds")
	}
	if !reflect.DeepEqual(buf, []int{1, 2}) {
		t.Errorf("expected {1, 2}, got %v", buf)
	}
	if n := rbuf.Free(); n != 4 {
		t.Errorf("expected to have 4 free, got %d", n)
	}

	if err := rbuf.Advance(2); err != nil {
		t.Error("expected not out of bounds")
	}
	if err := rbuf.Peek(buf, 0); err != nil {
		t.Error("expected not out of bounds")
	}
	if !reflect.DeepEqual(buf, []int{3, 4}) {
		t.Errorf("expected {3, 4}, got %v", buf)
	}
	if n := rbuf.Free(); n != 6 {
		t.Errorf("expected to have 6 free, got %d", n)
	}
}

func TestRingbuf_AdvanceBad(t *testing.T) {
	rbuf := ringbuf.New[int](10)
	if err := rbuf.Advance(2); err != ringbuf.ErrReadTooLarge {
		t.Error("expected out of bounds")
	}
}

func TestRingbuf_Wraparound(t *testing.T) {
	rbuf := ringbuf.New[int](10)
	if err := rbuf.Push([]int{1, 2, 3, 4, 5, 6}); err != nil {
		t.Error("expected not out of bounds")
	}
	if err := rbuf.Advance(6); err != nil {
		t.Error("expected not out of bounds")
	}

	if err := rbuf.Push([]int{1, 2, 3, 4, 5, 6}); err != nil {
		t.Error("expected not out of bounds")
	}
	buf := make([]int, 6)
	if err := rbuf.Peek(buf, 0); err != nil {
		t.Error("expected not out of bounds")
	}

	if !reflect.DeepEqual(buf, []int{1, 2, 3, 4, 5, 6}) {
		t.Errorf("expected {1, 2, 3, 4, 5, 6}, got %v", buf)
	}
}

func TestRingbuf_PeekOffset(t *testing.T) {
	rbuf := ringbuf.New[int](10)
	if err := rbuf.Push([]int{1, 2, 3, 4, 5, 6}); err != nil {
		t.Error("expected not out of bounds")
	}

	buf := make([]int, 3)
	if err := rbuf.Peek(buf, 3); err != nil {
		t.Error("expected not out of bounds")
	}

	if !reflect.DeepEqual(buf, []int{4, 5, 6}) {
		t.Errorf("expected {4, 5, 6}, got %v", buf)
	}
}

func TestRingbuf_PeekOffsetWraparound(t *testing.T) {
	rbuf := ringbuf.New[int](10)
	if err := rbuf.Push([]int{1, 2, 3, 4, 5, 6}); err != nil {
		t.Error("expected not out of bounds")
	}
	if err := rbuf.Advance(6); err != nil {
		t.Error("expected not out of bounds")
	}

	if err := rbuf.Push([]int{1, 2, 3, 4, 5, 6}); err != nil {
		t.Error("expected not out of bounds")
	}

	buf := make([]int, 3)
	if err := rbuf.Peek(buf, 3); err != nil {
		t.Error("expected not out of bounds")
	}

	if !reflect.DeepEqual(buf, []int{4, 5, 6}) {
		t.Errorf("expected {4, 5, 6}, got %v", buf)
	}
}

func TestRingbuf_PushEmptyDoesNothing(t *testing.T) {
	rbuf := ringbuf.New[int](10)
	if err := rbuf.Push([]int{}); err != nil {
		t.Error("expected not out of bounds")
	}
	if rbuf.Used() != 0 {
		t.Error("expected rbuf.Used() = 0")
	}
}

func TestRingbuf_AdvanceZeroDoesNothing(t *testing.T) {
	rbuf := ringbuf.New[int](10)
	if err := rbuf.Advance(0); err != nil {
		t.Error("expected not out of bounds")
	}
	if rbuf.Used() != 0 {
		t.Error("expected rbuf.Used() = 0")
	}
}

func TestRingbuf_PushWriteIndexBeforeReadIndex(t *testing.T) {
	rbuf := ringbuf.New[int](10)
	if err := rbuf.Push([]int{1, 2, 3, 4, 5, 6}); err != nil {
		t.Error("expected not out of bounds")
	}
	if err := rbuf.Advance(6); err != nil {
		t.Error("expected not out of bounds")
	}

	if err := rbuf.Push([]int{1, 2, 3, 4, 5, 6}); err != nil {
		t.Error("expected not out of bounds")
	}

	if err := rbuf.Push([]int{1, 2, 3, 4}); err != nil {
		t.Error("expected not out of bounds")
	}

	buf := make([]int, 10)
	if err := rbuf.Peek(buf, 10); err != nil {
		t.Error("expected not out of bounds")
	}

	if !reflect.DeepEqual(buf, []int{1, 2, 3, 4, 5, 6, 1, 2, 3, 4}) {
		t.Errorf("expected {1, 2, 3, 4, 5, 6, 1, 2, 3, 4}, got %v", buf)
	}
	if n := rbuf.Free(); n != 0 {
		t.Errorf("expected to have 0 free, got %d", n)
	}
}
