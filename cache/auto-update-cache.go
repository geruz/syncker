package cache

import (
	"sync"
	"time"
)

type Stringer interface {
	String() string
}
type CacheValueProvider func(query interface{}, prevValue interface{}) interface{}

type AutoUpdateCache interface {
	Get(q Stringer) interface{}
	ForceGet(q Stringer) interface{}
	Clear()
}

func NewAutoUpdateCache(pr CacheValueProvider, invalidateTime time.Duration) AutoUpdateCache {
	ch := autoUpdateCache{
		m:              make(map[string]*AutoValue),
		provider:       pr,
		lock:           &sync.Mutex{},
		invalidateTime: invalidateTime,
	}
	return &ch
}

type autoUpdateCache struct {
	m              map[string]*AutoValue
	provider       CacheValueProvider
	lock           *sync.Mutex
	invalidateTime time.Duration
}

func (uc *autoUpdateCache) Clear() {
	uc.lock.Lock()
	for k, c := range uc.m {
		delete(uc.m, k)
		c.Stop()
	}
	uc.m = make(map[string]*AutoValue)
	uc.lock.Unlock()
}

func (uc *autoUpdateCache) Get(q Stringer) interface{} {
	key := q.String()
	uc.lock.Lock()
	g, ok := uc.m[key]
	uc.lock.Unlock()
	if ok {
		return g.Value()
	}
	return uc.add(key, q)
}

func (uc *autoUpdateCache) ForceGet(q Stringer) interface{} {
	key := q.String()
	uc.lock.Lock()
	g, ok := uc.m[key]
	uc.lock.Unlock()
	if ok {
		return g.Force()
	}
	return uc.add(key, q)
}

func (uc *autoUpdateCache) add(key string, q Stringer) interface{} {
	provider := func(prev interface{}) interface{} {
		return uc.provider(q, prev)
	}
	g := NewAutoValue(provider, uc.invalidateTime)
	uc.lock.Lock()
	uc.m[key] = g
	uc.lock.Unlock()
	return g.Value()
}
