package set

import "sync"

type threadSafe struct {
	sync.RWMutex
	unsafe *threadUnsafe
}

func (t *threadSafe) Put(elements ...interface{}) {
	t.Lock()
	t.unsafe.Put(elements...)
	t.Unlock()
}
func (t *threadSafe) Contains(element interface{}) bool {
	t.RLock()
	defer t.RUnlock()
	return t.unsafe.Contains(element)

}
func (t *threadSafe) Pop(input string) bool {
	t.Lock()
	defer t.Unlock()
	return t.unsafe.Pop(input)
}

func (t *threadSafe) Equals(sOther ThreadSafeSet) bool {
	t.RLock()
	defer t.RUnlock()
	return t.unsafe.Equals(sOther.Unsafe())
}

func (t *threadSafe) Size() int {
	t.RLock()
	defer t.RUnlock()
	return t.unsafe.Size()

}

func (t *threadSafe) Unsafe() ThreadUnsafeSet {
	t.RLock()
	defer t.RUnlock()
	return t.unsafe

}

func (t *threadSafe) Keys() []interface{} {
	t.RLock()
	defer t.RUnlock()
	return t.unsafe.Keys()

}
