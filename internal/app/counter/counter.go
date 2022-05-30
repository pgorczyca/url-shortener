package counter

import (
	"errors"
	"fmt"
	"sync/atomic"
)

var (
	ErrFullCounter = errors.New("full counter")
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
		return 0, ErrFullCounter
	}
	return atomic.AddUint64(c.current, 1), nil
}

func (c *counter) isAlmostComplete() bool {
	current := atomic.LoadUint64(c.current)
	return float64(current) > 0.9*float64(c.end)
}

func (c *counter) GetShort() string {
	c.increment()
	fmt.Println(*c.current)
	return Base62Encode(int(*c.current - 1))
}
