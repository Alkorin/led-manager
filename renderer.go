package main

import ()

type getterFunc func() []Led

type Renderer interface {
	Size() []int
	SetGetters([]getterFunc)
	Start()
}
