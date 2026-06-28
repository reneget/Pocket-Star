package log

import "sync"

type RingBuffer struct {
	mu    sync.Mutex
	lines []string
	head  int
	size  int
	cap   int
}

var Default = NewRingBuffer(500)

func NewRingBuffer(cap int) *RingBuffer {
	return &RingBuffer{
		lines: make([]string, cap),
		cap:   cap,
	}
}

func (rb *RingBuffer) Push(line string) {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	rb.lines[rb.head] = line
	rb.head = (rb.head + 1) % rb.cap
	if rb.size < rb.cap {
		rb.size++
	}
}

func (rb *RingBuffer) Slice() []string {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	if rb.size == 0 {
		return nil
	}
	result := make([]string, rb.size)
	start := rb.head - rb.size
	if start < 0 {
		start += rb.cap
	}
	for i := 0; i < rb.size; i++ {
		result[i] = rb.lines[(start+i)%rb.cap]
	}
	return result
}
