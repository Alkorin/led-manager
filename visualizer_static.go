package main

import ()

type StaticColorVisualizer struct {
	BaseVisualizer

	length int
	Color  Led `property:"rw"`
}

func NewStaticColorVisualizer(length int, color Led) *StaticColorVisualizer {
	return &StaticColorVisualizer{
		BaseVisualizer: *NewBaseVisualizer("Static"),
		length:         length,
		Color:          color,
	}
}

func (v *StaticColorVisualizer) Start() {
	d := make([]Led, v.length)
	for i := range d {
		d[i] = v.Color
	}
	v.SendData(d)
}
