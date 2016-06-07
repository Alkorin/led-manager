package main

import ()

func main() {
	l := NewLedManager()
	l.AttachRenderer(NewTerminalRenderer(128))
	l.Start()
}
