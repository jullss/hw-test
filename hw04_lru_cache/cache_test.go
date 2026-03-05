package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
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

	t.Run("purge logic", func(_ *testing.T) {
		c := NewCache(3)

		c.Set("aa", 10)
		c.Set("bb", 20)
		c.Set("cc", 30)

		c.Get("aa")
		c.Set("bb", 300)

		c.Set("dd", 400)

		_, ok := c.Get("cc")
		require.False(t, ok)

		_, ok = c.Get("aa")
		require.True(t, ok)

		val, _ := c.Get("bb")
		require.Equal(t, 300, val)
	})

	t.Run("displace by capacity", func(t *testing.T) {
		capacity := 3
		c := NewCache(capacity)

		c.Set("1", 10)
		c.Set("2", 20)
		c.Set("3", 30)
		require.Equal(t, 3, c.(*lruCache).queue.Len())

		c.Set("4", 40)
		require.Equal(t, 3, c.(*lruCache).queue.Len())

		_, ok := c.Get("1")
		require.False(t, ok)

		val, ok := c.Get("4")
		require.True(t, ok)
		require.Equal(t, 40, val)
	})

	t.Run("displace of least recently used", func(t *testing.T) {
		c := NewCache(3)

		c.Set("1", 10)
		c.Set("2", 20)
		c.Set("3", 30)

		c.Get("1")     // [1, 3, 2]
		c.Set("2", 22) // [2, 1, 3]

		c.Set("4", 4) // [4, 2, 1]

		_, ok := c.Get("3")
		require.False(t, ok)

		val, ok := c.Get("1")
		require.True(t, ok)
		require.Equal(t, 10, val)

		val, ok = c.Get("2")
		require.True(t, ok)
		require.Equal(t, 22, val)
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Skip() // Remove me if task with asterisk completed.

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
