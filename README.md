# gotil

Collection of often needed utility packages for go projects.

## Set

Exports the following interfaces & implementations.
```go
type Set interface {
	Put(elements ...interface{})
	Contains(element interface{}) bool
	Pop(input interface{}) bool
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
```

- Initialise a set:

```go 
  safeSet := set.ThreadSafe("1", "2", 3, 3)
  unsafeSet := set.ThreadUnSafe("1", "2", "3", "3")
```


## Loadbalancer

Provides a not-so-resource-greedy framework to submit work units to be executed concurrently. 

The balancer scales the number of goroutines up & down based on the both frequency of the work-units being submitted and the duration of execution.

The total number of goroutines are increased if the number of tasks submitted and waiting in the queue reaches a certain threshold, on the other hand the goroutines are reused if there are still work to do instead of being recreated.

Besides, the idle goroutines will kill themselves to release the resources.

- Simply initialise a new loadbalancer with a consumer function
```go
balancer := loadbalancer.New(func(data interface{}) {
   time.Sleep(time.Millisecond * time.Duration(float64(data.(int))
  }) 
```

- Then start submitting work-units to the balancer any time in the future within the lifecycle.

```go
for i := 0; i < 10000; i++ {
  balancer.Submit(i)
}
```


