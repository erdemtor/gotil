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
	found, index := FindFirst([]string{"asd", "as", "da", "asdasd"}, finder)
	assert.Equal("as", found)
	assert.Equal(1, index)
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

func TestRandomStringOfLength(t *testing.T) {
	lenght := 5
	randStr := RandOfSize(lenght)
	assert.Len(t, randStr, lenght)
}
