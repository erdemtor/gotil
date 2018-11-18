package set_test

import (
	"fmt"
	"gotil/set"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOfStrings(t *testing.T) {
	t.Parallel()

	set := set.OfStrings("1", "2", "3", "3")
	assert.EqualValues(t, 3, set.Size())
}

func TestStrings_Put(t *testing.T) {
	t.Parallel()

	set := set.OfStrings()
	itemCount := 5000
	for i := 0; i < itemCount; i++ {
		set.Put(fmt.Sprintf("%d", i))
	}
	assert.EqualValues(t, itemCount, set.Size())
}

func TestStrings_Equals(t *testing.T) {
	t.Parallel()

	set1 := set.OfStrings()
	set2 := set.OfStrings()
	itemCount := 5000
	for i := 0; i < itemCount; i++ {
		set1.Put(fmt.Sprintf("%d", i))
		set2.Put(fmt.Sprintf("%d", i))
	}
	assert.True(t, set1.Equals(set2))
	assert.True(t, set2.Equals(set1))
	set1.Pop("1")
	assert.False(t, set1.Equals(set2))
	set1.Put("-1")
	assert.False(t, set1.Equals(set2))
}

func TestStrings_Delete(t *testing.T) {
	t.Parallel()

	set := set.OfStrings()
	itemCount := 5000
	for i := 0; i < itemCount; i++ {
		set.Put(fmt.Sprintf("%d", i))
	}
	assert.EqualValues(t, itemCount, set.Size())

	for i := itemCount - 1; i >= 0; i-- {
		assert.True(t, set.Pop(fmt.Sprintf("%d", i)))
	}
	assert.EqualValues(t, 0, set.Size())

	assert.False(t, set.Pop("nonexisting"))
}
