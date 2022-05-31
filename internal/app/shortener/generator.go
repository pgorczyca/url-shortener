package shortener

import (
	"sync"
)

type ShortGenerator struct {
	mutex    sync.Mutex
	provider RangeProvider
	current  *counterRange
	backup   *counterRange
	counter  uint64
}

type counterRange struct {
	start    uint64
	end      uint64
	treshold uint64
}

func NewCounterManager(provider RangeProvider) (*ShortGenerator, error) {
	newRange, err := provider.GetRange()
	if err != nil {
		return nil, err
	}
	return &ShortGenerator{provider: provider, current: newRange, counter: newRange.start}, nil
}

func (sg *ShortGenerator) GetShort() (string, error) {
	next, err := sg.getNext()
	if err != nil {
		return "", err
	}
	return base62Encode(next), nil
}

func (sg *ShortGenerator) getNext() (uint64, error) {

	sg.mutex.Lock()
	defer sg.mutex.Unlock()

	sg.counter++

	if sg.counter > sg.current.treshold && sg.backup == nil {
		var err error
		sg.backup, err = sg.provider.GetRange()
		if err != nil {
			return 0, err
		}
	}

	if sg.counter == sg.current.end {
		sg.current = sg.backup
		sg.backup = nil
	}
	return sg.counter, nil
}
