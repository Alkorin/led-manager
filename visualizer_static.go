package main

import ()

type StaticColorVisualizer struct {
	BaseVisualizer

	length int
	color  Led
}

func NewStaticColorVisualizer(length int, color Led) *StaticColorVisualizer {
	return &StaticColorVisualizer{
		BaseVisualizer: *NewBaseVisualizer(),
		length:         length,
		color:          color,
	}
}

func (v *StaticColorVisualizer) Start() {
	d := make([]Led, v.length)
	for i := range d {
		d[i] = v.color
	}
	v.SendData(d)
}
