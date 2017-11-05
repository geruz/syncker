package lockers

import (
	"sync"
)

type locks struct {
	m      map[string]*sync.Mutex
	global *sync.Mutex
}

func NewLocks() locks {
	return locks{map[string]*sync.Mutex{}, &sync.Mutex{}}
}

type Locks interface {
	Lock(name string)
	Unlock(name string)
}
func (l locks) Lock(name string) {
	l.global.Lock()
	if _, ok := l.m[name]; !ok {
		l.m[name] = &sync.Mutex{}
	}
	c := l.m[name]
	l.global.Unlock()
	c.Lock()

}
func (l locks) Unlock(name string) {
	l.global.Lock()
	c := l.m[name]
	l.global.Unlock()
	c.Unlock()
}
