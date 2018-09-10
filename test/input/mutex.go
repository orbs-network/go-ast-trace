package main

import "sync"

func main() {
	mutex := &sync.Mutex{}
	rwmutex := &sync.RWMutex{}

	normalLockExample(mutex)
	readLockExample(rwmutex)
	writeLockExample(rwmutex)
}

func normalLockExample(mutex *sync.Mutex) {
	mutex.Lock()
	defer mutex.Unlock()
}

func readLockExample(mutex *sync.RWMutex) {
	mutex.RLock()
	defer mutex.RUnlock()
}

func writeLockExample(mutex *sync.RWMutex) {
	mutex.Lock()
	defer mutex.Unlock()
}