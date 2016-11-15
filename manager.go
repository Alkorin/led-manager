package main

import (
	"log"

	"github.com/steeve/broadcaster"
)

type LedManager struct {
	buffer []Led

	renderers   map[uint64]Renderer
	visualizers map[uint64]Visualizer
	apiEvents   *broadcaster.Broadcaster
}

func NewLedManager(renderers []Renderer) *LedManager {
	// Allocate needed memory
	totalSize := 0
	for _, r := range renderers {
		for _, s := range r.Size() {
			totalSize += s
		}
	}

	log.Printf("Total renderer size: %d", totalSize)
	buffer := make([]Led, totalSize)

	// Attach getter for each renderers
	curPos := 0
	for _, r := range renderers {
		getters := make([]getterFunc, len(r.Size()))
		for i, rendererSize := range r.Size() {
			start := curPos
			end := curPos + rendererSize
			getters[i] = func() []Led {
				return buffer[start:end]
			}
			curPos += rendererSize
		}
		r.SetGetters(getters)
	}

	// Allocate renderers map
	rendererMap := make(map[uint64]Renderer)
	for _, r := range renderers {
		rendererMap[r.ID()] = r
	}

	return &LedManager{
		buffer:      buffer,
		apiEvents:   broadcaster.NewBroadcaster(),
		renderers:   rendererMap,
		visualizers: make(map[uint64]Visualizer),
	}
}

func (l *LedManager) AttachVisualizer(v Visualizer, start int, end int) {
	l.visualizers[v.ID()] = v
	go func() {
		for {
			d := <-v.OutputChan()
			copy(l.buffer[start:end+1], d)
		}
	}()
}

func (l *LedManager) Start() {
	// Start Renderers
	for _, r := range l.renderers {
		go r.Start()
	}

	// Start Visualizers
	for _, v := range l.visualizers {
		go v.Start()
	}

	// Start HTTP server
	go l.StartApi()
}
