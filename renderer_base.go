package main

import (
	"sync/atomic"
)

var rendererIdCounter uint64 = 0

type BaseRenderer struct {
	getters []getterFunc
	id      uint64
	name    string
}

func NewBaseRenderer(name string) *BaseRenderer {
	return &BaseRenderer{
		id:   atomic.AddUint64(&rendererIdCounter, 1),
		name: name,
	}
}

func (r *BaseRenderer) SetGetters(g []getterFunc) {
	r.getters = g
}

func (r *BaseRenderer) GetData(i int) []Led {
	return r.getters[i]()
}

func (r *BaseRenderer) ID() uint64 {
	return r.id
}

func (r *BaseRenderer) Name() string {
	return r.name
}
