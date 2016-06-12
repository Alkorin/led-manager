package main

import (
	"log"
	"time"
)

type LedManager struct {
	buffer []Led

	renderers []Renderer
}

func NewLedManager() *LedManager {
	return &LedManager{}
}

func (l *LedManager) AttachRenderer(r Renderer) {
	l.renderers = append(l.renderers, r)
}

func (l *LedManager) Start() {
	// Allocate needed memory
	totalSize := 0
	for _, r := range l.renderers {
		for _, s := range r.Size() {
			totalSize += s
		}
	}

	log.Printf("Total renderer size: %d", totalSize)
	l.buffer = make([]Led, totalSize)

	// Attach getter for each renderers
	curPos := 0
	for _, r := range l.renderers {
		getters := make([]getterFunc, len(r.Size()))
		for i, rendererSize := range r.Size() {
			start := curPos
			end := curPos + rendererSize
			getters[i] = func() []Led {
				return l.buffer[start:end]
			}
			curPos += rendererSize
		}
		r.SetGetters(getters)
		go r.Start()
	}

	// Do rainbow over buffer
	j := 0.0
	for range time.Tick(100 * time.Millisecond) {
		j += 0.03
		for i := 0; i < totalSize; i++ {
			r, g, b := hueToRGB(j + float64(i)/128.0)
			l.buffer[i] = Led{r, g, b, 0}
		}
	}
}
