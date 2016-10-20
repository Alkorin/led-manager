package main

import ()

type Visualizer interface {
	PropertyHandler

	OutputChan() <-chan []Led
	Start()

	// Methods for API
	Name() string
	ID() uint64
}
