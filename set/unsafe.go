package set

type threadUnsafe struct {
	m map[interface{}]struct{}
}


func (t*threadUnsafe)Put(elements ...interface{})  {
	for _, elem := range elements {
		t.m[elem] = struct{}{}
	}
}

func (t *threadUnsafe) Contains(element interface{}) bool {
	_, exists := t.m[element]
	return exists
}

//Pop removes the given string and returns true, if the string doesn't exist returns false
func (t *threadUnsafe) Pop(input string) bool {
	_, exists := t.m[input]
	if exists {
		delete(t.m, input)
	}
	return exists
}

//Equals checks the sets have one-to-one correspondence and have the same sizes
func (t *threadUnsafe) Equals(sOther ThreadUnsafeSet) bool {
	if t.Size() != sOther.Size() {
		return false
	}
	for _, key := range t.Keys() {
		if !sOther.Contains(key) {
			return false
		}
	}
	return true
}

//Size returns the count of the keys
func (t *threadUnsafe) Size() int {
	return len(t.m)
}

//Keys returns an array of the keys of the set, order is unpredictable
func (t *threadUnsafe) Keys() []interface{} {
	var keys []interface{}
	for k := range t.m {
		keys = append(keys, k)
	}
	return keys
}

