package injection

import (
	"fmt"
	"sync"
)

var beforeNanos = make(map[uint64]int64)
var beforeNanosMutex = &sync.Mutex{}

func BeforeLock(id uint64) {
	ts, nano := getTimestamp()
	fn, cs := getCallStack()

	beforeNanosMutex.Lock()
	beforeNanos[id] = nano
	beforeNanosMutex.Unlock()

	fmt.Printf("ASTRACE %s BeforeLock %s [%s]\n", ts, fn, cs)
}

func AfterLock(id uint64) {
	ts, nano := getTimestamp()
	fn, cs := getCallStack()

	beforeNanosMutex.Lock()
	beforeNano := beforeNanos[id]
	delete(beforeNanos, id)
	beforeNanosMutex.Unlock()

	waitedMillis := (nano - beforeNano) / 1000000

	fmt.Printf("ASTRACE %s AfterLock %s %dms [%s]\n", ts, fn, waitedMillis, cs)
}