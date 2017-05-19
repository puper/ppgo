package ppgo

import (
	"sync"

	"github.com/puper/errors"
)

type Creator func(interface{}) (interface{}, error)

type Container struct {
	sync.RWMutex
	creators  map[string]Creator
	instances map[string]interface{}
}

func NewContainer() *Container {
	return &Container{
		creators:  make(map[string]Creator),
		instances: make(map[string]interface{}),
	}
}

func (this *Container) Register(name string, creator Creator) {
	this.Lock()
	defer this.Unlock()
	this.creators[name] = creator
}

func (this *Container) Configure(name string, value interface{}) error {
	this.Lock()
	defer this.Unlock()
	if creator, ok := this.creators[name]; ok {
		instance, err := creator(value)
		if err != nil {
			return errors.Annotatef(err, "configure %s failed use config: %s", name, value)
		}
		this.instances[name] = instance
		return nil
	}
	return errors.NotFoundf("creator %s not found", name)
}

func (this *Container) Set(name string, value interface{}) {
	this.Lock()
	defer this.Unlock()
	this.instances[name] = value
}

func (this *Container) Get(name string) (interface{}, error) {
	this.RLock()
	defer this.RUnlock()
	if instance, ok := this.instances[name]; ok {
		return instance, nil
	}
	return nil, errors.NotFoundf("instance %s not found", name)
}

func (this *Container) MustGet(name string) interface{} {
	if instance, err := this.Get(name); err == nil {
		return instance
	}
	return nil
}

func (this *Container) GetContainer(name string) (*Container, error) {
	instance, err := this.Get(name)
	if err != nil {
		return nil, err
	}
	if container, ok := instance.(*Container); ok {
		return container, nil
	}
	return nil, errors.NotValidf("instance %s is not a Container", name)
}
