package engine

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/pkg/errors"
	"github.com/puper/ppgo/helpers"
)

type Engine struct {
	componentNames   []string
	components       map[string]Component
	componentConfigs map[string]*ComponentConfig
	instances        map[string]interface{}
	config           *Config

	stopMutex sync.Mutex
	stopped   bool
}

func New(cfg *Config) *Engine {
	return &Engine{
		componentNames:   make([]string, 0),
		components:       make(map[string]Component),
		componentConfigs: make(map[string]*ComponentConfig),
		instances:        make(map[string]interface{}),
		config:           cfg,
	}
}

func (this *Engine) RegisterComponent(name string, c Component) {
	this.components[name] = c
}

func (this *Engine) RegisterComponents(cs map[string]Component) {
	for k, v := range cs {
		this.RegisterComponent(k, v)
	}
}

func (this *Engine) CreateComponent(config *ComponentConfig) (interface{}, error) {
	if c, ok := this.components[config.BackendName]; ok {
		instance, err := c.Create(config.Config)
		return instance, err
	}
	return nil, fmt.Errorf("component `%v` Not Registered", config.Name)
}

func (this *Engine) Init() error {
	configs := make([]*ComponentConfig, 0)
	err := helpers.StructDecode(this.config.Get("components"), &configs, "json")
	if err != nil {
		return errors.WithMessage(err, "decode component config error")
	}
	for _, c := range configs {
		if _, ok := this.componentConfigs[c.Name]; ok {
			return fmt.Errorf("component `%v` already configured before", c.Name)
		}
		this.componentNames = append(this.componentNames, c.Name)
		this.componentConfigs[c.Name] = c
		instance, err := this.CreateComponent(c)
		if err != nil {
			return err
		}
		this.instances[c.Name] = instance

	}
	for _, k := range this.componentNames {
		if err := this.components[this.componentConfigs[k].BackendName].Init(this.componentConfigs[k].Config, this.instances[k]); err != nil {
			return errors.Wrapf(err, "init component `%v` error", k)
		}
	}
	return nil
}

func (this *Engine) GetInstance(name string) interface{} {
	return this.instances[name]
}

func (this *Engine) Start() error {
	go this.handleSysSignal()
	stopCh := make(chan struct{}, 1)
	for _, name := range this.componentNames {
		go func(name string, instance interface{}) {
			log.Println(fmt.Sprintf("start component: %v", name))
			err := this.components[this.componentConfigs[name].BackendName].Start(this.componentConfigs[name].Config, instance)
			if err != MethodNotImplemented {
				select {
				case stopCh <- struct{}{}:
				default:
				}
			}
		}(name, this.instances[name])
	}
	<-stopCh
	return this.Stop()
}

func (this *Engine) Stop() error {
	this.stopMutex.Lock()
	if this.stopped == false {
		this.stopped = true
	} else {
		this.stopMutex.Unlock()
		return nil
	}
	for i := len(this.componentNames) - 1; i >= 0; i-- {
		err := this.components[this.componentConfigs[this.componentNames[i]].BackendName].Stop(this.componentConfigs[this.componentNames[i]].Config, this.instances[this.componentNames[i]])
		log.Printf("stop component: %v, %v\n", this.componentNames[i], err)
	}
	this.stopMutex.Unlock()
	return nil
}

func (this *Engine) handleSysSignal() {
	sChan := make(chan os.Signal)
	for {
		signal.Notify(sChan, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		sig := <-sChan
		switch sig {
		case os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			this.Stop()
		}

	}
}

func (this *Engine) GetConfig() *Config {
	return this.config
}
