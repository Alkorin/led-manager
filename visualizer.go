package main

import ()

type Visualizer interface {
	PropertyHandler

	OutputChan() <-chan []Led
	Start()
	Close()

	// Methods for API
	Name() string
	ID() uint64
}
