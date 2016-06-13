package main

import ()

type BaseRenderer struct {
	getters []getterFunc
}

func NewBaseRenderer() *BaseRenderer {
	return &BaseRenderer{}
}

func (r *BaseRenderer) SetGetters(g []getterFunc) {
	r.getters = g
}

func (r *BaseRenderer) GetData(i int) []Led {
	return r.getters[i]()
}
