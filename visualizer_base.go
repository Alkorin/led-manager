package main

import ()

type BaseVisualizer struct {
	outputChan chan []Led
}

func NewBaseVisualizer() *BaseVisualizer {
	return &BaseVisualizer{
		outputChan: make(chan []Led, 1),
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
