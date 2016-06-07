package main

import ()

type getterFunc func() []Led

type Renderer interface {
	Size() int
	SetGetter(getterFunc)
	Start()
}
