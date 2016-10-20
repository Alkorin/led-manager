package main

import (
	"sync/atomic"
)

var visualizerIdCounter uint64 = 0

type BaseVisualizer struct {
	outputChan chan []Led
	name       string
	id         uint64
}

func NewBaseVisualizer(name string) *BaseVisualizer {
	return &BaseVisualizer{
		outputChan: make(chan []Led, 1),
		name:       name,
		id:         atomic.AddUint64(&visualizerIdCounter, 1),
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

func (v *BaseVisualizer) ID() uint64 {
	return v.id
}

func (v *BaseVisualizer) OnPropertyChanged(string) {
}
