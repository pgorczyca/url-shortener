package counter

import "fmt"

type ShortGenerator struct {
	provider RangeProvider
	counter  *counter
	backup   *counterRange
}

func NewCounterManager(provider RangeProvider) (*ShortGenerator, error) {
	newRange, err := provider.GetRange()
	if err != nil {
		return nil, err
	}
	return &ShortGenerator{provider: provider, counter: newCounter(newRange)}, nil
}

func (sg *ShortGenerator) GetShort() (string, error) {
	next, _ := sg.getNext()
	fmt.Println(next)
	return Base62Encode(next), nil
}

func (sg *ShortGenerator) getNext() (uint64, error) {

	// jezeli < 90% to return next
	if !sg.counter.isAlmostComplete() {
		next, err := sg.counter.increment()
		if err != nil {
			return 0, err
		}
		return next, nil
	}
	// jezeli powyzej 90% to

	//  czy backup nil, stowrzyc
	var err error
	if sg.backup == nil {
		fmt.Println("backup created")
		sg.backup, err = sg.provider.GetRange()
		if err != nil {
			return 0, err
		}
	}

	// swap jak next == end i ustawienie backup = nil
	if next == sg.counter.ctRange.end {
		fmt.Println("swap backup")
		sg.counter.ctRange = sg.backup
		sg.backup = nil
	}

	// return powyzej 90%
	return next, nil

}
