package cache

import (
	"time"
	"sync"
)

type ValueProvider func(interface{}) (interface{})

type AutoValue struct {
	value 	 interface{}
	provider ValueProvider
	nextIter chan struct{}
	sleep 	 time.Duration
	lock 	 *sync.Mutex
}
func NewAutoValue (provider ValueProvider, sleep time.Duration) *AutoValue{
	av := &AutoValue{
		value: nil,
		provider: provider,
		sleep:sleep,
		nextIter: make(chan struct{}),
		lock: &sync.Mutex{},
	}
	av.start()
	return av
}
func(av *AutoValue) Force()  (interface{}) {
	value := av.update()
	av.nextIter <- struct{}{}
	return value
}
func(av *AutoValue) Stop() {
	close(av.nextIter)
}

func (av *AutoValue) Value() interface{} {
	av.lock.Lock()
	defer av.lock.Unlock()
	return av.value
}
func (av *AutoValue) SetValue (v interface{}) {
	av.lock.Lock()
	defer av.lock.Unlock()
	av.value = v
}


func(av *AutoValue) start() {
	av.update()
	go func(){
		for {
			select {
				case <-time.After(av.sleep):
					av.update()
				case _, isClose := <-av.nextIter:
					if isClose {
						return
					}
					break
			}
		}
	}()
}
func(av *AutoValue) update()  (interface{}) {
	value := av.provider(av.Value())
	av.SetValue(value)
	return value
}