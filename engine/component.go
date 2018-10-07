package engine

import "errors"

var (
	MethodNotImplemented = errors.New("Method Not Implemented")
)

type ComponentConfig struct {
	Name        string      `json:"name"`
	BackendName string      `json:"backend_name"`
	Config      interface{} `json:"config"`
}

type Component interface {
	Create(interface{}) (interface{}, error)
	Init(interface{}, interface{}) error
	Start(interface{}, interface{}) error
	Stop(interface{}, interface{}) error
	Update(interface{}, interface{}) error
}

type BaseComponent struct {
}

func (this *BaseComponent) Create(interface{}) (interface{}, error) {
	return nil, MethodNotImplemented
}

func (this *BaseComponent) Init(interface{}, interface{}) error {
	return nil
}

func (this *BaseComponent) Start(interface{}, interface{}) error {
	return MethodNotImplemented
}

func (this *BaseComponent) Stop(interface{}, interface{}) error {
	return nil
}

func (this *BaseComponent) Update(interface{}, interface{}) error {
	return nil
}
