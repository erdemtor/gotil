package set_test

import (
	"github.com/stretchr/testify/assert"
	"gotil/set"
	"sync"
	"testing"
)

func TestThreadSafe(t *testing.T) {
	s := set.ThreadSafe("1", "2", 3, 3)
	assert.EqualValues(t, 3, s.Size())
}

func TestThreadSafe_Put(t *testing.T) {
	t.Parallel()
	wg := sync.WaitGroup{}
	s := set.ThreadSafe()
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(i int) {
			s.Put(i)
			wg.Done()
		}(i)
	}
	wg.Wait()
	assert.Equal(t, 20, s.Size())
	for i := 0; i < 20; i++ {
		assert.True(t, s.Contains(i))
	}
}

func TestThreadSafe_Equals(t *testing.T) {
	t.Parallel()

	s1 := set.ThreadSafe(1, "2", 3, "4")
	s2 := set.ThreadSafe("2", 3, "4", 1)
	s3 := set.ThreadSafe("a", "b", "c", "d")
	wg := sync.WaitGroup{}
	wg.Add(4)
	go func() {
		assert.True(t, s1.Equals(s2))
		wg.Done()
	}()
	go func() {
		assert.True(t, s2.Equals(s1))
		wg.Done()
	}()
	go func() {
		assert.False(t, s3.Equals(s1))
		wg.Done()
	}()
	go func() {
		assert.False(t, s1.Equals(s3))
		wg.Done()
	}()
	wg.Wait()
}
