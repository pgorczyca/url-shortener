package counter

import (
	"errors"
	"sync/atomic"
)

type counter struct {
	start   uint64
	end     uint64
	current *uint64
}

func newCounter(start uint64, end uint64) *counter {
	return &counter{start: start, end: end, current: &start}
}

func (c *counter) increment() (uint64, error) {
	curr := atomic.LoadUint64(c.current)
	if curr > c.end {
		return 0, errors.New("Out of range")
	}
	return atomic.AddUint64(c.current, 1), nil
}

func (c *counter) isAlmostComplete() bool {
	current := atomic.LoadUint64(c.current)
	return float64(current/(c.end-c.start)) > 0.9
}
