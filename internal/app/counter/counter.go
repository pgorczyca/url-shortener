package counter

import (
	"errors"
	"sync/atomic"
)

var (
	ErrFullCounter = errors.New("full counter")
)

type counterRange struct {
	start uint64
	end   uint64
}

type counter struct {
	ctRange *counterRange
	current *uint64
}

func newCounter(cr *counterRange) *counter {
	return &counter{ctRange: cr, current: &cr.start}
}

func (c *counter) loadRange(counterRange) error {
	return nil
}

func (c *counter) increment() (uint64, error) {
	curr := atomic.LoadUint64(c.current)
	if curr > c.ctRange.end {
		return 0, ErrFullCounter
	}
	return atomic.AddUint64(c.current, 1), nil
}

func (c *counter) isAlmostComplete() bool {
	current := atomic.LoadUint64(c.current)
	return float64(current) > 0.8*float64(c.ctRange.end)
}

// func (c *counter) GetShort() string {
// 	defer c.increment()
// 	fmt.Println(*c.current)
// 	return Base62Encode(int(*c.current))
// }
