package promise_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"gotil/concurrency/promise"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	p := promise.New(context.Background(), func() (interface{}, error) {
		time.Sleep(time.Millisecond)
		return 1, nil
	}, func(err error) {
		assert.Fail(t, "this should never be called")
	})
	res, err := p.Get()
	assert.NoError(t, err)
	assert.Equal(t, 1, res)
}

func TestCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var callerErr error
	p := promise.New(ctx, func() (interface{}, error) {
		time.Sleep(time.Second * 3)
		return 1, nil
	}, func(err error) {
		callerErr = err
	})
	cancel()
	res, err := p.Get()
	assert.Error(t, err)
	assert.Error(t, callerErr)
	assert.Equal(t, ctx.Err(), err)
	assert.Nil(t, res)
}

func TestFirst(t *testing.T) {
	t.Parallel()
	var promises []promise.Promise
	for i := 1; i < 10; i++ {
		iC := i
		promises = append(promises, promise.New(context.Background(), func() (interface{}, error) {
			time.Sleep(time.Second * time.Duration(iC))
			return iC, nil
		}, nil))

	}
	start := time.Now()
	res := promise.First(promises...)
	end := time.Now()
	assert.Equal(t, 1, res)
	elapsedTimeInSeconds := end.Sub(start).Seconds()
	assert.True(t, elapsedTimeInSeconds > 1 && elapsedTimeInSeconds < 2)
}

func TestAll(t *testing.T) {
	t.Parallel()

	var promises []promise.Promise
	for i := 1; i <= 4; i++ {
		iC := i
		promises = append(promises, promise.New(context.Background(), func() (interface{}, error) {
			time.Sleep(time.Second * time.Duration(iC))
			return iC, nil
		}, nil))

	}
	start := time.Now()
	res := promise.All(promises...)
	end := time.Now()
	elapsedTimeInSeconds := end.Sub(start).Seconds()
	assert.True(t, elapsedTimeInSeconds > 4 && elapsedTimeInSeconds < 5)

	assert.Len(t, res, 4)
}
