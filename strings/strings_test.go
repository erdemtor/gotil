package strings

import (
	"reflect"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArrayContains(t *testing.T) {
	assert := assert.New(t)
	contains := ArrayContains([]string{}, "s")
	assert.False(contains)
	contains = ArrayContains([]string{"a", "s", "d"}, "s")
	assert.True(contains)
	var arr []string
	contains = ArrayContains(arr, "s")
	assert.False(contains)
}

func TestArraysIntersect(t *testing.T) {
	assert := assert.New(t)
	res := ArraysIntersect([]string{"a", "s", "d"}, []string{"d", "a", "e"})
	assert.Contains(res, "a")
	assert.Contains(res, "d")
	assert.Len(res, 2)

	res = ArraysIntersect([]string{"a", "s", "d"}, []string{})
	assert.Len(res, 0)
}

func TestFindFirst(t *testing.T) {
	assert := assert.New(t)
	finder := func(a string) bool {
		return len(a) == 2
	}
	assert.Equal("as", FindFirst([]string{"asd", "as", "da", "asdasd"}, finder))
	assert.Empty(FindFirst(nil, finder))
	assert.Empty(FindFirst([]string{}, finder))
}

func TestGetKeys(t *testing.T) {
	assert := assert.New(t)
	stringMap := map[string]string{"a": "1", "b": "2", "c": "3"}

	keys := GetKeys(stringMap)
	assert.Len(keys, len(stringMap))
	sort.Strings(keys)
	assert.True(reflect.DeepEqual(keys, []string{"a", "b", "c"}))
}

func TestToJSONStringIgnoreError(t *testing.T) {
	data := map[string]int{"key": 1, "key2": 2}
	stringified := ToJSONStringIgnoreError(data)
	assert.NotEmpty(t, stringified)
	assert.Contains(t, stringified, "key2")
	stringified = ToJSONStringIgnoreError(func() {}) // functions are not marshall-able
	assert.Empty(t, stringified)
}

func TestRandomStringOfLength(t *testing.T) {
	lenght := 5
	randStr := RandomStringOfLength(lenght)
	assert.Len(t, randStr, lenght)
}
