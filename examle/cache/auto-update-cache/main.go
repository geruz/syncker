package main

import (
	"github.com/geruz/syncker/cache"
	"time"
	"sync"
)


func main() {
	c := cache.NewAutoUpdateCache(getValue, time.Second)
	wg := sync.WaitGroup{}
	wg.Add(3)
	f := func(s ss) {
		for i := 0; i < 50; i++ {
			println(i, " " + s.Key + ":", c.Get(s).(int))
			time.Sleep(time.Millisecond * 300)
		}
		wg.Done()
	}
	go f(ss{"a",1})
	go f(ss{"b",2})
	go f(ss{"с",3})

	go func(){
		for {
			time.Sleep(time.Second * 5)
			println("Clear")
			c.Clear()
		}
	}();

	go func(){
		for {
			time.Sleep(time.Millisecond * 300)
			println("force a:")
			c.ForceGet(ss{"a", 1})
		}
	}();

	wg.Wait()
}

type ss struct {
	Key string
	Delta int
}
func (s ss) String() string{
	return s.Key
}

func getValue (q interface{}, prev interface{}) interface{} {
	s := q.(ss)
	if prev == nil {
		return s.Delta
	}
	return s.Delta + prev.(int)
}
