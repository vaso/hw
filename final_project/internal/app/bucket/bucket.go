package bucket

import (
	"log"
	"sync"
	"time"
)

type LeakyBucket interface {
	Add() bool
	Reset()
}

type bucket struct {
	Max       int
	Count     int
	leakSpeed float64
	lastLeak  time.Time

	lock  sync.Mutex
	timer Timer
}

func NewBucketWithTimer(capacity int, leakFrame time.Duration, timer Timer) LeakyBucket {
	return &bucket{
		Max:       capacity,
		Count:     0,
		leakSpeed: float64(capacity) / leakFrame.Seconds(),
		lastLeak:  time.Now(),
		timer:     timer,
	}
}

func NewBucket(capacity int, leakFrame time.Duration) LeakyBucket {
	return NewBucketWithTimer(capacity, leakFrame, NewTimer())
}

func (buck *bucket) Add() bool {
	buck.lock.Lock()
	defer buck.lock.Unlock()

	curTime := buck.timer.Now()
	log.Printf("timer: %s", curTime.String())
	newLeakDuration := curTime.Sub(buck.lastLeak)
	leakSeconds := newLeakDuration.Seconds()
	leakAmount := int(leakSeconds * buck.leakSpeed)
	if leakAmount > 0 {
		buck.lastLeak = curTime
	}
	newBucketCount := buck.Count - leakAmount + 1
	if newBucketCount > buck.Max {
		buck.Count = buck.Max
		return false
	}
	if newBucketCount <= 1 {
		buck.Count = 1
		return true
	}
	buck.Count = newBucketCount
	return true
}

func (buck *bucket) Reset() {
	buck.lock.Lock()
	defer buck.lock.Unlock()

	buck.Count = 0
	buck.lastLeak = buck.timer.Now()
}
