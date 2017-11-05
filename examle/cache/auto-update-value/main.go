package main

import (
	"github.com/geruz/syncker/cache"
	"time"
	"fmt"
)



func main() {

	cache := cache.NewAutoValue(getProvider(), time.Second)
	start := time.Now()
	go func() {
		for {
			fmt.Printf("%v, value: %v\n", time.Since(start).Round(time.Millisecond), cache.Value())
			time.Sleep(time.Millisecond * 500)
		}
	}()
	time.Sleep(10 * time.Second)
	fmt.Println("stop update loop")
	cache.Stop()
	time.Sleep(10 * time.Second)
}

func getProvider()  cache.ValueProvider{
	i := 0
	return func(interface{}) (interface{}) {
		i++
		return i
	}
}

