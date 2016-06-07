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
		totalSize += r.Size()
	}

	log.Printf("Total renderer size: %d", totalSize)
	l.buffer = make([]Led, totalSize)

	// Attach getter for each renderers
	curPos := 0
	for _, r := range l.renderers {
		rendererSize := r.Size()
		start := curPos
		end := curPos + rendererSize
		r.SetGetter(func() []Led {
			return l.buffer[start:end]
		})
		go r.Start()
		curPos += rendererSize
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
