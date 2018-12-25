package prodcons_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"gotil/concurrency/prodcons"
	"gotil/set"
	"gotil/strings"
	"testing"
)

func TestProdCons(t *testing.T) {
	t.Parallel()
	t.Run("nothing provided ends", func(t *testing.T) {
		t.Parallel()
		err := (&prodcons.Runner{}).Run()
		assert.NoError(t, err)
	})

	t.Run("the consumer receives exactly one time", func(t *testing.T) {
		t.Parallel()
		var producerCount = 1000

		produced := set.OfStrings()
		runner := prodcons.Runner{}
		for i := 0; i < producerCount; i++ {
			producedStr := strings.RandOfSize(5)
			produced.Put(producedStr)
			runner.Producer(func() (interface{}, error) { return producedStr, nil })

		}

		consumed := set.OfStrings()
		consumerCalledCount := 0
		err := runner.Consumer(func(data interface{}) error {
			consumerCalledCount++
			s, ok := data.(string)
			if !ok {
				return errors.New("unable to cast to string")
			}
			consumed.Put(s)
			return nil
		}).Run()
		assert.NoError(t, err)
		assert.Equal(t, producerCount, consumerCalledCount)
		assert.True(t, consumed.Equals(produced))
	})

	t.Run("error is returned", func(t *testing.T) {
		t.Parallel()
		runner := prodcons.Runner{}

		err := runner.
			Consumer(func(i interface{}) error { return nil }).
			Producer(func() (interface{}, error) { return 0, errors.New("err") }).
			Run()
		assert.Error(t, err)

		var producerCount = 1000
		runner = prodcons.Runner{}
		for i := 0; i < producerCount; i++ {
			runner.Producer(func() (interface{}, error) { return 0, nil })
		}
		runner.Producer(func() (interface{}, error) { return -1, errors.New("error") })

		err = runner.Consumer(func(interface{}) error { return nil }).Run()
		assert.Error(t, err)
		runner = prodcons.Runner{}
		err = runner.
			Consumer(func(i interface{}) error { return errors.New("consumer failed") }).
			Producer(func() (interface{}, error) { return 1, nil }).
			Run()
		assert.Error(t, err)

	})

	t.Run("no producer doesn't block", func(t *testing.T) {
		t.Parallel()
		runner := prodcons.Runner{}

		err := runner.Consumer(func(i interface{}) error { return nil }).Run()
		assert.NoError(t, err)
	})

}

func TestManyProducersManyConsumers(t *testing.T) {
	r := prodcons.Runner{}
	m := map[int]struct{}{}

	producerCount := 500
	for i := 0; i < producerCount; i++ {
		iC := i
		m[iC] = struct{}{}
		r.Producer(func() (interface{}, error) {
			return iC, nil
		})
	}

	producedVals := make(chan int, producerCount)
	for i := 0; i < 10; i++ {
		r.Consumer(func(i interface{}) error {
			producedVals <- i.(int)
			return nil
		})
	}
	err := r.Run()
	close(producedVals)
	assert.NoError(t, err)

	for val := range producedVals {
		_, exists := m[val]
		assert.True(t, exists)
		delete(m, val)
	}
	assert.Len(t, m, 0)
}
