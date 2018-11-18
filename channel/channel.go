package channel

import "sync"

func Merge(channels ...<-chan int) <-chan int {
	resChan := make(chan int)
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(len(channels))
		for i := range channels {
			c:= channels[i]
			go func() {
				for v := range c {
					resChan <- v
				}
				wg.Done()
			}()
		}
		wg.Wait()
		close(resChan)
	}()

	return resChan

}
