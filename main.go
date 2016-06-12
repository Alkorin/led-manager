package main

import ()

func main() {
	l := NewLedManager()
	l.AttachRenderer(NewTerminalRenderer(128))
	l.AttachVisualizer(NewRainbowVisualizer(128), 0, 127)
	l.Start()
	<-make(chan struct{})
}
