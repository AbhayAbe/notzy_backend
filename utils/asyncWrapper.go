package utils

func AsyncWrapper(f func(interface{}) Result, data interface{}) <-chan Result {
	ch := make(chan Result)
	go func() {
		defer close(ch)
		res := f(data)
		ch <- res
	}()
	return ch
}
