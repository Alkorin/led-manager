package main

import (
	"sync/atomic"
)

var visualizerIdCounter uint64 = 0

type BaseVisualizer struct {
	Visualizer

	outputChan chan []Led
	name       string
	id         uint64

	isClosed bool
	quit     chan struct{}
}

func NewBaseVisualizer(name string) *BaseVisualizer {
	return &BaseVisualizer{
		outputChan: make(chan []Led, 1),
		name:       name,
		id:         atomic.AddUint64(&visualizerIdCounter, 1),
		quit:       make(chan struct{}),
	}
}

func (v *BaseVisualizer) Close() {
	v.isClosed = true
	close(v.quit)
	close(v.outputChan)
}

func (v *BaseVisualizer) SendData(d []Led) {
	// Don't try to send data on a closed channel, it panics
	if !v.isClosed {
		// Send data if chan is free
		select {
		case v.outputChan <- d:
		default:
		}
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
