package goman

import (
	"errors"
	"sync"
	"time"
)

type GoMan struct {
	sync.WaitGroup
}

func (this *GoMan) Go(f func()) {
	this.Add(1)
	go func() {
		defer this.Done()
		f()
	}()

}

func (this *GoMan) Wait(timeout time.Duration) error {
	if timeout > 0 {
		stopCh := make(chan bool)
		go func() {
			this.WaitGroup.Wait()
			select {
			case stopCh <- true:
			default:
			}
		}()
		select {
		case <-time.After(timeout):
			return errors.New("go wait timeout")
		case <-stopCh:
			return nil
		}
	}
	this.WaitGroup.Wait()
	return nil
}
