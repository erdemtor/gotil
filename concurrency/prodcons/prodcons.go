package prodcons


type Producer func() (interface{}, error)
type Consumer func(interface{}) error




type Runner struct {
	producers []Producer
	consumers []Consumer
}

func (r *Runner) Producer(p Producer) *Runner {
	r.producers = append(r.producers, p)
	return r
}

func (r *Runner) Consumer(c Consumer) *Runner {
	r.consumers = append(r.consumers, c)
	return r
}

func (r *Runner) Run() error {
	var buffer = make(chan interface{}, len(r.producers))
	var consumerErrs = make(chan error, len(r.producers))

	for i := range r.consumers {
		consumer := r.consumers[i]
		go func() {
			for data := range buffer {
				err := consumer(data)
				consumerErrs <- err
			}
		}()
	}
	var producerErrs = make(chan error, len(r.producers))

	for i := range r.producers {
		producer := r.producers[i]
		go func() {
			data, err := producer()
			buffer <- data
			producerErrs <- err
		}()
	}
	for i := 0; i < len(r.producers); i++ {
		err := <-producerErrs
		if err != nil {
			return err
		}
	}
	close(buffer)

	for i := 0; i < len(r.producers); i++ {
		err := <-consumerErrs
		if err != nil {
			return err
		}
	}
	return nil

}
