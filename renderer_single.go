package main

import ()

type SingleRenderer struct {
	BaseRenderer

	size int
}

func NewSingleRenderer(size int, name string) *SingleRenderer {
	return &SingleRenderer{
		BaseRenderer: *NewBaseRenderer(name),
		size:         size,
	}
}

func (r *SingleRenderer) GetData() []Led {
	return r.getters[0]()
}

func (r *SingleRenderer) Size() []int {
	return []int{r.size}
}
