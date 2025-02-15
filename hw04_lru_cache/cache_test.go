package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require" //nolint:all
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(2)

		_ = c.Set("aaa", 100)
		_ = c.Set("bbb", 200)
		_ = c.Set("ccc", 300)

		removed, ok := c.Get("aaa")
		require.False(t, ok)
		require.Nil(t, removed)
	})

	t.Run("purge less last active by read logic", func(t *testing.T) {
		c := NewCache(2)

		_ = c.Set("aaa", 100)
		_ = c.Set("bbb", 200)
		_, ok := c.Get("aaa") // move aaa to the top
		require.True(t, ok)
		_ = c.Set("ccc", 300)

		first, ok := c.Get("ccc")
		require.True(t, ok)
		require.Equal(t, 300, first)

		second, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, second)

		removed, ok := c.Get("bbb")
		require.False(t, ok)
		require.Nil(t, removed)
	})

	t.Run("purge less last active by update logic", func(t *testing.T) {
		c := NewCache(2)

		_ = c.Set("aaa", 100)
		_ = c.Set("bbb", 200)
		ok := c.Set("aaa", 110) // move aaa to the top and update its value to 110
		require.True(t, ok)
		_ = c.Set("ccc", 300)

		first, ok := c.Get("ccc")
		require.True(t, ok)
		require.Equal(t, 300, first)

		second, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 110, second)

		removed, ok := c.Get("bbb")
		require.False(t, ok)
		require.Nil(t, removed)
	})

	t.Run("clear cache", func(t *testing.T) {
		c := NewCache(10)

		_ = c.Set("aaa", 100)
		_ = c.Set("bbb", 200)

		c.Clear()

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})
}

func TestCacheMultithreading(_ *testing.T) {
	// t.Skip() // Remove me if task with asterisk completed.

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
