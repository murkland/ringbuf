package ringbuf_test

import (
	"reflect"
	"testing"

	"github.com/murkland/ringbuf"
)

func TestRingbuf_PushOutOfBounds(t *testing.T) {
	rbuf := ringbuf.New[int](10)
	defer func() {
		if r := recover(); r == nil {
			t.Error("panic expected")
		}
	}()
	rbuf.Push([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11})
}

func TestRingbuf_PeekOutOfBounds(t *testing.T) {
	rbuf := ringbuf.New[int](10)
	buf := make([]int, 2)
	defer func() {
		if r := recover(); r == nil {
			t.Error("panic expected")
		}
	}()
	rbuf.Peek(buf, 0)
}
func TestRingbuf_Push(t *testing.T) {
	rbuf := ringbuf.New[int](10)
	rbuf.Push([]int{1, 2, 3, 4, 5, 6})
	if n := rbuf.Free(); n != 4 {
		t.Errorf("expected to have 4 free, got %d", n)
	}
}

func TestRingbuf_PushThenPush(t *testing.T) {
	rbuf := ringbuf.New[int](10)
	rbuf.Push([]int{1, 2, 3, 4, 5, 6})
	rbuf.Push([]int{4, 5, 6, 7})
	if n := rbuf.Free(); n != 0 {
		t.Errorf("expected to have 0 free, got %d", n)
	}
}

func TestRingbuf_PushThenAdvance(t *testing.T) {
	rbuf := ringbuf.New[int](10)
	rbuf.Push([]int{1, 2, 3, 4, 5, 6})
	if n := rbuf.Free(); n != 4 {
		t.Errorf("expected to have 4 free, got %d", n)
	}
	rbuf.Advance(2)
	if n := rbuf.Free(); n != 6 {
		t.Errorf("expected to have 6 free, got %d", n)
	}
}

func TestRingbuf_PushThenPop(t *testing.T) {
	rbuf := ringbuf.New[int](10)
	rbuf.Push([]int{1, 2, 3, 4, 5, 6})
	if n := rbuf.Free(); n != 4 {
		t.Errorf("expected to have 4 free, got %d", n)
	}
	buf := make([]int, 2)
	rbuf.Pop(buf, 0)
	if n := rbuf.Free(); n != 6 {
		t.Errorf("expected to have 6 free, got %d", n)
	}
	if !reflect.DeepEqual(buf, []int{1, 2}) {
		t.Errorf("expected {1, 2}, got %v", buf)
	}
}

func TestRingbuf_PushThenAdvanceThenPeek(t *testing.T) {
	rbuf := ringbuf.New[int](10)
	rbuf.Push([]int{1, 2, 3, 4, 5, 6})
	buf := make([]int, 2)
	rbuf.Peek(buf, 0)
	if !reflect.DeepEqual(buf, []int{1, 2}) {
		t.Errorf("expected {1, 2}, got %v", buf)
	}
	if n := rbuf.Free(); n != 4 {
		t.Errorf("expected to have 4 free, got %d", n)
	}

	rbuf.Advance(2)
	rbuf.Peek(buf, 0)
	if !reflect.DeepEqual(buf, []int{3, 4}) {
		t.Errorf("expected {3, 4}, got %v", buf)
	}
	if n := rbuf.Free(); n != 6 {
		t.Errorf("expected to have 6 free, got %d", n)
	}
}

func TestRingbuf_AdvanceBad(t *testing.T) {
	rbuf := ringbuf.New[int](10)
	defer func() {
		if r := recover(); r == nil {
			t.Error("panic expected")
		}
	}()
	rbuf.Advance(2)
}

func TestRingbuf_Wraparound(t *testing.T) {
	rbuf := ringbuf.New[int](10)
	rbuf.Push([]int{1, 2, 3, 4, 5, 6})
	rbuf.Advance(6)

	rbuf.Push([]int{1, 2, 3, 4, 5, 6})
	buf := make([]int, 6)
	rbuf.Peek(buf, 0)

	if !reflect.DeepEqual(buf, []int{1, 2, 3, 4, 5, 6}) {
		t.Errorf("expected {1, 2, 3, 4, 5, 6}, got %v", buf)
	}
}

func TestRingbuf_PeekOffset(t *testing.T) {
	rbuf := ringbuf.New[int](10)
	rbuf.Push([]int{1, 2, 3, 4, 5, 6})

	buf := make([]int, 3)
	rbuf.Peek(buf, 3)

	if !reflect.DeepEqual(buf, []int{4, 5, 6}) {
		t.Errorf("expected {4, 5, 6}, got %v", buf)
	}
}

func TestRingbuf_PeekOffsetWraparound(t *testing.T) {
	rbuf := ringbuf.New[int](10)
	rbuf.Push([]int{1, 2, 3, 4, 5, 6})
	rbuf.Advance(6)

	rbuf.Push([]int{1, 2, 3, 4, 5, 6})

	buf := make([]int, 3)
	rbuf.Peek(buf, 3)

	if !reflect.DeepEqual(buf, []int{4, 5, 6}) {
		t.Errorf("expected {4, 5, 6}, got %v", buf)
	}
}

func TestRingbuf_PushEmptyDoesNothing(t *testing.T) {
	rbuf := ringbuf.New[int](10)
	rbuf.Push([]int{})
	if rbuf.Used() != 0 {
		t.Error("expected rbuf.Used() = 0")
	}
}

func TestRingbuf_AdvanceZeroDoesNothing(t *testing.T) {
	rbuf := ringbuf.New[int](10)
	rbuf.Advance(0)
	if rbuf.Used() != 0 {
		t.Error("expected rbuf.Used() = 0")
	}
}

func TestRingbuf_PushWriteIndexBeforeReadIndex(t *testing.T) {
	rbuf := ringbuf.New[int](10)
	rbuf.Push([]int{1, 2, 3, 4, 5, 6})
	rbuf.Advance(6)

	rbuf.Push([]int{1, 2, 3, 4, 5, 6})

	rbuf.Push([]int{1, 2, 3, 4})

	buf := make([]int, 10)
	rbuf.Peek(buf, 10)

	if !reflect.DeepEqual(buf, []int{1, 2, 3, 4, 5, 6, 1, 2, 3, 4}) {
		t.Errorf("expected {1, 2, 3, 4, 5, 6, 1, 2, 3, 4}, got %v", buf)
	}
	if n := rbuf.Free(); n != 0 {
		t.Errorf("expected to have 0 free, got %d", n)
	}
}
