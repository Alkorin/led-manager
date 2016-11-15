package main

import ()

func main() {
	l := NewLedManager([]Renderer{NewTerminalRenderer(128)})
	l.AttachVisualizer(NewRainbowVisualizer(128), 0, 127)
	l.Start()
	<-make(chan struct{})
}
