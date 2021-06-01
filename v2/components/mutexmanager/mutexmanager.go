package mutexmanager

import (
	"sync"
)

var (
	defaultMutexManager = NewMutexManager()
)

func Default() *MutexManager {
	return defaultMutexManager
}

type Mutex struct {
	sync.RWMutex
	locks int64
}

func NewMutexManager() *MutexManager {
	return &MutexManager{
		mutexes: map[string]*Mutex{},
	}
}

type MutexManager struct {
	mutex   sync.Mutex
	mutexes map[string]*Mutex
}

func Lock(key string) {
	Default().Lock(key)
}
func Unlock(key string) {
	Default().Unlock(key)
}

func RLock(key string) {
	Default().RLock(key)
}
func RUnlock(key string) {
	Default().RUnlock(key)
}

func (this *MutexManager) Lock(key string) {
	this.mutex.Lock()
	if _, ok := this.mutexes[key]; !ok {
		this.mutexes[key] = &Mutex{}
	}
	this.mutexes[key].locks++
	this.mutex.Unlock()
	this.mutexes[key].Lock()

}

func (this *MutexManager) Unlock(key string) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if _, ok := this.mutexes[key]; ok {
		this.mutexes[key].Unlock()
		this.mutexes[key].locks--
		if this.mutexes[key].locks == 0 {
			delete(this.mutexes, key)
		}
	} else {
		panic("unlock of unlocked mutex")
	}
}

func (this *MutexManager) RLock(key string) {
	this.mutex.Lock()
	if _, ok := this.mutexes[key]; !ok {
		this.mutexes[key] = &Mutex{}
	}
	this.mutexes[key].locks++
	this.mutex.Unlock()
	this.mutexes[key].RLock()

}

func (this *MutexManager) RUnlock(key string) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if _, ok := this.mutexes[key]; ok {
		this.mutexes[key].RUnlock()
		this.mutexes[key].locks--
		if this.mutexes[key].locks == 0 {
			delete(this.mutexes, key)
		}
	} else {
		panic("unlock of unlocked mutex")
	}
}
