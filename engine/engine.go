package engine

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"code.int.thoseyears.com/golang/ppgo/helpers"
)

type Engine struct {
	components       map[string]Component
	componentConfigs map[string]*ComponentConfig
	instances        map[string]interface{}
	config           *Config
}

func New(cfg *Config) *Engine {
	return &Engine{
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
	if c, ok := this.components[config.Name]; ok {
		instance, err := c.Create(config.Config)
		return instance, err
	}
	return nil, fmt.Errorf("component `%v` Not Registered", config.Name)
}

func (this *Engine) Init() error {
	for k, v := range this.config.GetStringMap("components") {
		cc := new(ComponentConfig)
		err := helpers.StructDecode(v, cc, "json")
		if err != nil {
			return err
		}
		this.componentConfigs[k] = cc
		instance, err := this.CreateComponent(cc)
		if err != nil {
			return err
		}
		this.instances[k] = instance

	}
	return nil
}

func (this *Engine) GetInstance(name string) interface{} {
	return this.instances[name]
}

func (this *Engine) Start() error {
	go this.handleSysSignal()
	wg := sync.WaitGroup{}
	wg.Add(len(this.instances))
	for name, instance := range this.instances {
		go func(name string, instance interface{}) {
			defer wg.Done()
			log.Println(fmt.Sprintf("start component: %v", name))
			this.components[this.componentConfigs[name].Name].Start(this.componentConfigs[name].Config, instance)
		}(name, instance)
	}
	wg.Wait()
	return nil
}

func (this *Engine) Stop() error {
	for name, instance := range this.instances {
		this.components[this.componentConfigs[name].Name].Stop(this.componentConfigs[name].Config, instance)
		log.Println(fmt.Sprintf("stop component: %v", name))
	}
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
