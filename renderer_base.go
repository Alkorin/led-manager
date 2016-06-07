package main

import ()

type BaseRenderer struct {
	getData getterFunc
}

func NewBaseRenderer() *BaseRenderer {
	return &BaseRenderer{}
}

func (r *BaseRenderer) SetGetter(g getterFunc) {
	r.getData = g
}
