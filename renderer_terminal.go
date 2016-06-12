package main

import (
	"fmt"
	"time"
)

type TerminalRenderer struct {
	SingleRenderer
}

func NewTerminalRenderer(size int) *TerminalRenderer {
	return &TerminalRenderer{
		SingleRenderer: *NewSingleRenderer(size),
	}
}

func (r *TerminalRenderer) Start() {
	for range time.Tick(100 * time.Millisecond) {
		data := r.GetData()
		for _, color := range data {
			fmt.Printf("\x1b[48;2;%d;%d;%dm ", byte(color.Red*255), byte(color.Green*255), byte(color.Blue*255))
		}
		fmt.Print("\x1b[0m\r")
	}
}
