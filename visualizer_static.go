package main

import ()

type StaticColorVisualizer struct {
	BaseVisualizer

	length int
	Red    float64 `property:"rw"`
	Green  float64 `property:"rw"`
	Blue   float64 `property:"rw"`
	White  float64 `property:"rw"`
}

func NewStaticColorVisualizer(length int, color Led) *StaticColorVisualizer {
	return &StaticColorVisualizer{
		BaseVisualizer: *NewBaseVisualizer("StaticColor"),
		length:         length,
		Red:            color.Red,
		Green:          color.Green,
		Blue:           color.Blue,
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
	color := Led{v.Red, v.Green, v.Blue, v.White}
	for i := range d {
		d[i] = color
	}
	v.SendData(d)
}
