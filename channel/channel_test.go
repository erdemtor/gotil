package channel_test

import (
	"github.com/stretchr/testify/assert"
	"gotil/channel"
	"math/rand"
	"testing"
	"time"
)

func TestMerge(t *testing.T) {
	t.Parallel()
	var m = map[int]struct{}{}
	rand.Seed(time.Now().Unix())
	size := rand.Intn(100)
	var chans = make([]<-chan int, size, size)
	for i := 0; i < size; i++ {
		chans[i] = asChan(i)
		m[i] = struct{}{}
	}
	merged := channel.Merge(chans...)

	for v := range merged {
		_, exists := m[v]
		assert.True(t, exists)
		delete(m, v)
	}
	assert.Len(t, m, 0)

}

func asChan(nums ...int) <-chan int {
	resChan := make(chan int)
	go func() {
		for _, num := range nums {
			resChan <- num
		}
		close(resChan)
	}()
	return resChan

}
