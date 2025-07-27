package bucket

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockTimer struct {
	mock.Mock
}

func (m *MockTimer) Now() time.Time {
	args := m.Called()

	return args.Get(0).(time.Time)
}

func Test_bucket_Reset(t *testing.T) {
	lastLeakTime := time.Now()
	curCount := 3
	bucket := &bucket{
		Max:       3,
		Count:     curCount,
		leakSpeed: 20,
		lastLeak:  lastLeakTime,
		lock:      sync.Mutex{},
		timer:     NewTimer(),
	}
	t.Run("check Reset", func(t *testing.T) {
		require.Equal(t, 3, bucket.Count)
		require.Equal(t, lastLeakTime, bucket.lastLeak)

		bucket.Reset()

		require.Equal(t, 0, bucket.Count)
		require.Less(t, lastLeakTime, bucket.lastLeak)
	})
}

func Test_bucket_Add(t *testing.T) {
	timerMock := new(MockTimer)
	capacity := 60
	timeframe := time.Minute
	curTime := time.Now()
	leakSpeed := capacity / int(timeframe.Seconds())
	leakAmount := 15
	leakTime := curTime.Add(time.Duration(leakAmount*leakSpeed)*time.Second + 100*time.Millisecond) // skip 15.1 secs
	buck := NewBucketWithTimer(capacity, timeframe, timerMock)

	timerMock.On("Now").Return(curTime).Times(capacity + 1) // 1 is for false check
	// leak 15 after 15 sec
	timerMock.On("Now").Return(leakTime).Times(leakAmount + 1) // 1 is for false check
	wait := sync.WaitGroup{}
	t.Run("check Add", func(t *testing.T) {
		for i := 0; i < capacity; i++ {
			wait.Add(1)
			go func() {
				defer wait.Done()
				require.True(t, buck.Add())
			}()
		}
		wait.Wait()
		require.False(t, buck.Add())

		for i := 0; i < leakAmount; i++ {
			wait.Add(1)
			go func() {
				defer wait.Done()
				require.True(t, buck.Add())
			}()
		}
		wait.Wait()
		require.False(t, buck.Add())
	})
}
