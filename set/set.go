package set

type Set interface {
	Put(elements ...interface{})
	Contains(element interface{}) bool
	Pop(input string) bool
	Size() int
	Keys() []interface{}
}

type ThreadSafeSet interface {
	Set
	Equals(sOther ThreadSafeSet) bool
	Unsafe() ThreadUnsafeSet
}

type ThreadUnsafeSet interface {
	Set
	Equals(sOther ThreadUnsafeSet) bool
}

func ThreadSafe(initialVals ...interface{}) ThreadSafeSet {
	s := &threadSafe{unsafe: &threadUnsafe{m: make(map[interface{}]struct{})}}
	s.Put(initialVals...)
	return s

}

func ThreadUnSafe(initialVals ...interface{}) ThreadUnsafeSet {
	s := &threadUnsafe{m: make(map[interface{}]struct{})}
	s.Put(initialVals...)
	return s

}
