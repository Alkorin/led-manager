package main

import ()

type Visualizer interface {
	GetOutputChan() <-chan []Led
	Start()
}
