package main

import ()

type Visualizer interface {
	OutputChan() <-chan []Led
	Start()
}
