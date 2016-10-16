package main

import ()

type getterFunc func() []Led

type Renderer interface {
	Size() []int
	SetGetters([]getterFunc)
	Start()

	// Methods for API
	ID() uint64
	Name() string
}
