package stopcontext

import "sync"

func New() (*Stopper, *StopContext) {
	tmp := &stopContext{
		stopCh:    make(chan struct{}, 1),
		stoppedCh: make(chan struct{}, 1),
	}
	return &Stopper{
			tmp,
		}, &StopContext{
			tmp,
		}
}

type Stopper struct {
	*stopContext
}

func (this *Stopper) StopAndGetError() error {
	this.stopContext.stop()
	return this.stopContext.getError()
}

func (this *Stopper) Stop() {
	this.stopContext.stop()
}

func (this *Stopper) Wait() <-chan struct{} {
	return this.stopContext.wait()
}

func (this *Stopper) GetError() error {
	return this.stopContext.getError()
}

type StopContext struct {
	*stopContext
}

func (this *StopContext) Done() <-chan struct{} {
	return this.stopContext.done()
}

func (this *StopContext) SetError(err error) {
	this.stopContext.setError(err)
}

type stopContext struct {
	err   error
	mutex sync.Mutex

	stoping bool
	stopCh  chan struct{}

	stopped   bool
	stoppedCh chan struct{}
}

func (this *stopContext) stop() {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if !this.stoping {
		this.stoping = true
		close(this.stopCh)
	}
}

func (this *stopContext) done() <-chan struct{} {
	return this.stopCh
}

func (this *stopContext) setError(err error) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.err = err
	if !this.stopped {
		this.stopped = true
		close(this.stoppedCh)
	}
}

func (this *stopContext) wait() <-chan struct{} {
	return this.stoppedCh
}

func (this *stopContext) getError() error {
	<-this.wait()
	this.mutex.Lock()
	defer this.mutex.Unlock()
	return this.err
}
