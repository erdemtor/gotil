package semaphore

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sync"
	"testing"
	"time"
)


func TestSemaphore_Lock(t *testing.T) {
	t.Parallel()
	for k := 0; k < 20; k++ {
		t.Run(fmt.Sprintf("%d", k), func(t *testing.T) {
			t.Parallel()
			rand.Seed(time.Now().Unix())
			semaphoreSize := rand.Intn(100)
			goRoutineCount := 200 + rand.Intn(100)
			sleepDuration := rand.Intn(100)
			s := New(semaphoreSize)
			start := time.Now()
			wg := sync.WaitGroup{}
			wg.Add(goRoutineCount)
			for i := 0; i < goRoutineCount; i++ {
				go func() {
					s.Lock()
					time.Sleep(time.Millisecond * time.Duration(sleepDuration))
					s.UnLock()
					wg.Done()

				}()
			}
			wg.Wait()
			end := time.Now()
			var elapsedTime = float64(end.Sub(start) / time.Millisecond)
			expectedTotalDuration := float64(goRoutineCount) / float64(semaphoreSize) * float64(sleepDuration)
			assert.True(t, elapsedTime > expectedTotalDuration && elapsedTime < expectedTotalDuration+1000)

		})

	}

}
