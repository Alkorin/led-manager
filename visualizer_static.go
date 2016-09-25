package main

import ()

type StaticColorVisualizer struct {
	BaseVisualizer

	length int
	Color  ColorRGB `property:"rw"`
	White  float64  `property:"rw"`
}

func NewStaticColorVisualizer(length int, color Led) *StaticColorVisualizer {
	return &StaticColorVisualizer{
		BaseVisualizer: *NewBaseVisualizer("StaticColor"),
		length:         length,
		Color:          ColorRGB{color.Red, color.Green, color.Blue},
		White:          color.White,
	}
}

func (v *StaticColorVisualizer) Start() {
	v.SendColor()
}

func (v *StaticColorVisualizer) OnPropertyChanged(string) {
	v.SendColor()
}

func (v *StaticColorVisualizer) SendColor() {
	d := make([]Led, v.length)
	color := Led{v.Color.Red, v.Color.Green, v.Color.Blue, v.White}
	for i := range d {
		d[i] = color
	}
	v.SendData(d)
}
