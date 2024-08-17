package spmcb

import (
	"strconv"
	"testing"
)

func TestNewBuffer(t *testing.T) {
	size := 10
	sb, err := New[string](size)
	if err != nil {
		t.Errorf("create new buffer failed: %v", err)
		return
	}
	if sb.size != size {
		t.Errorf("Create new buffer size not match, expect: %d, got: %d\n", size, sb.size)
		return
	}
	for i := 0; i < size+1; i++ {
		if err := sb.Produce(strconv.Itoa(i)); err != nil {
			t.Errorf("Produce new value failed: %v", err)
			return
		}
	}
	if len(sb.list) != size {
		t.Errorf("Buffer size not match, expect: %d, got: %d\n", size, len(sb.list))
		return
	}
	if sb.list[0] != strconv.Itoa(size) {
		t.Errorf("Buffer first element not match, expect: %d, got: %s\n", size, sb.list[0])
		return
	}
	for n, v := range sb.list {
		if n == 0 {
			continue
		}
		if v != strconv.Itoa(n) {
			t.Errorf("Value not match, expect: %d, got: %s\n", n, v)
			return
		}
	}
}

func TestNewConsumer(t *testing.T) {
	size := 10
	sb, err := New[int](size)
	if err != nil {
		t.Errorf("Create new buffer failed: %v\n", err)
		return
	}
	cs, err := NewConsumer(sb)
	if err != nil {
		t.Errorf("Create new consumer failed: %v\n", err)
		return
	}
	sb.Produce(1)
	v, err := cs.Consume()
	if err != nil {
		t.Errorf("Consume a buffer failed: %v\n", err)
		return
	}
	if v != 1 {
		t.Errorf("Consume buffer result not match, expect: %d, got: %d\n", 1, v)
		return
	}
}
