package container

import (
	"sync"

	"github.com/mitchellh/mapstructure"
	"github.com/puper/errors"
)

type Creator func(cfg interface{}) (interface{}, error)

type Container struct {
	sync.RWMutex
	instances map[string]interface{}
	creators  map[string]Creator
}

func NewContainer() *Container {
	return &Container{
		instances: make(map[string]interface{}),
		creators:  make(map[string]Creator),
	}
}

func (this *Container) Set(name string, instance interface{}) {
	this.Lock()
	defer this.Unlock()
	this.instances[name] = instance
}

func (this *Container) Get(name string) (interface{}, error) {
	this.RLock()
	defer this.RUnlock()
	if instance, ok := this.instances[name]; ok {
		return instance, nil
	}
	return nil, errors.New("instance not found")
}

func (this *Container) Register(name string, creator Creator) {
	this.Lock()
	defer this.Unlock()
	this.creators[name] = creator
}

func (this *Container) Create(cfg interface{}) (interface{}, error) {
	this.RLock()
	defer this.RUnlock()
	return this.create(cfg)
}

func (this *Container) create(cfg interface{}) (interface{}, error) {
	var ccfg struct {
		Type   string
		Config interface{}
	}
	err := mapstructure.WeakDecode(cfg, &ccfg)
	if err != nil {
		return nil, errors.NotValidf("can not trans interface: %v", cfg)
	}
	if creator, ok := this.creators[ccfg.Type]; ok {
		return creator(ccfg.Config)
	}
	return nil, errors.NotFoundf("creator %s not found", ccfg.Type)
}

func (this *Container) Configure(name string, cfg interface{}) error {
	this.Lock()
	defer this.Unlock()

	instance, err := this.create(cfg)
	if err == nil {
		this.instances[name] = instance
	}
	return err
}

func (this *Container) ConfigureAll(cfg map[string]interface{}) error {
	this.Lock()
	defer this.Unlock()
	instances := make(map[string]interface{})
	for k, v := range cfg {
		instance, err := this.create(v)
		if err != nil {
			return err
		}
		instances[k] = instance
	}
	this.instances = instances
	return nil
}
