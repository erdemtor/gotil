package set

//OfStrings creates a strings with optional initial values
func OfStrings(initialVals ...string) *strings {
	set := &strings{
		threadUnsafe: threadUnsafe{m : map[interface{}]struct{}{}},
	}
	for _, initialVal := range initialVals {
		set.Put(initialVal)
	}
	return set
}

//strings is the struct for storing unique strings
type strings struct {
	threadUnsafe
}

//Put adds the given string to the set
func (s *strings) Put(new string) {
	s.threadUnsafe.Put(new)
}

//Contains checks if the given string exists
func (s *strings) Contains(existing string) bool {
	return s.threadUnsafe.Contains(existing)
}

//Pop removes the given string and returns true, if the string doesn't exist returns false
func (s *strings) Pop(input string) bool {
	return s.threadUnsafe.Pop(input)
}

//Equals checks the sets have one-to-one correspondence and have the same sizes
func (s *strings) Equals(sOther *strings) bool {
	if s.Size() != sOther.Size() {
		return false
	}
	for _, key := range s.Keys() {
		if !sOther.Contains(key) {
			return false
		}
	}
	return true
}

//Size returns the count of the keys
func (s *strings) Size() int {
	return s.threadUnsafe.Size()
}

//Keys returns an array of the keys of the set, order is unpredictable
func (s *strings) Keys() []string {
	var keys []string
	for k := range s.m {
		keys = append(keys, k.(string))
	}
	return keys
}
