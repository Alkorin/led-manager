package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

type LedManager struct {
	buffer []Led

	renderers   []Renderer
	visualizers []Visualizer
}

func NewLedManager() *LedManager {
	return &LedManager{}
}

func (l *LedManager) AttachRenderer(r Renderer) {
	l.renderers = append(l.renderers, r)
}

func (l *LedManager) AttachVisualizer(v Visualizer, start int, end int) {
	l.visualizers = append(l.visualizers, v)
	go func() {
		for {
			d := <-v.GetOutputChan()
			copy(l.buffer[start:end+1], d)
		}
	}()
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

	// Start Visualizers
	for _, v := range l.visualizers {
		go v.Start()
	}

	// Start HTTP server
	http.Handle("/buffer", websocket.Handler(func(ws *websocket.Conn) {
		for range time.Tick(100 * time.Millisecond) {
			j, _ := json.Marshal(l.buffer)
			ws.Write(j)
		}
	}))
	http.Handle("/", http.FileServer(http.Dir("./web")))
	go http.ListenAndServe(":8080", nil)
}
