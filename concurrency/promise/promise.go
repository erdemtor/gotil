package promise

import "context"

type Promise struct {
	onErr         func(error)
	resultingChan chan result
}

type result struct {
	res interface{}
	err error
}

func New(ctx context.Context, f func() (interface{}, error), onErr func(error)) Promise {
	var p = Promise{
		resultingChan: make(chan result),
		onErr:         onErr,
	}
	var innerResultChan = make(chan result)
	go func() {
		res, err := f()
		innerResultChan <- result{res: res, err: err}
	}()

	go func() {
		select {
		case <-ctx.Done():
			p.onErr(ctx.Err())
			p.resultingChan <- result{err: ctx.Err()}
		case res := <-innerResultChan:
			p.resultingChan <- res
		}

	}()
	return p
}

func All(promises ...Promise) []interface{} {
	var res []interface{}
	for ind := range promises {
		p := promises[ind]
		r, err := p.Get()
		if err != nil {
			p.onErr(err)
			continue
		}
		res = append(res, r)
	}
	return res
}

func OnCompletion(results chan<- interface{}, promises ...Promise) error {
	var errChan = make(chan error, len(promises))
	for ind := range promises {
		p := promises[ind]
		go func() {
			get, err := p.Get()
			if err != nil {
				errChan <- err
			} else {
				results <- get
			}
		}()
	}
	return <-errChan
}

func First(promises ...Promise) interface{} {
	res := make(chan interface{})
	ctx, cancel := context.WithCancel(context.Background())
	for ind := range promises {
		p := promises[ind]
		go func() {
			select {
			case <-ctx.Done():
			case returnedRes := <-p.resultingChan:
				if returnedRes.err != nil {
					p.onErr(returnedRes.err)
				} else {
					res <- returnedRes.res
				}
			}
		}()
	}
	first := <-res

	cancel()
	return first
}

func (p *Promise) Get() (interface{}, error) {

	select {
	case res := <-p.resultingChan:
		return res.res, res.err

	}
}
