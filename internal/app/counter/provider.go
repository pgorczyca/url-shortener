package counter

type provider interface {
	getCounter() (*counter, error)
}

type etcdProvider struct {
}

func (e *etcdProvider) getCounter() (*counter, error) {
	return newCounter(0, 9), nil
}
