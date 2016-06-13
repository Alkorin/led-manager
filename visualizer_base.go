package main

import ()

type BaseVisualizer struct {
	outputChan chan []Led
	name       string
}

func NewBaseVisualizer(name string) *BaseVisualizer {
	return &BaseVisualizer{
		outputChan: make(chan []Led, 1),
		name:       name,
	}
}

func (v *BaseVisualizer) SendData(d []Led) {
	// Send data if chan is free
	select {
	case v.outputChan <- d:
	default:
	}
}

func (v *BaseVisualizer) OutputChan() <-chan []Led {
	return v.outputChan
}

func (v *BaseVisualizer) Name() string {
	return v.name
}
