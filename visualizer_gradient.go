package main

import ()

type GradientVisualizer struct {
	BaseVisualizer

	length     int
	StartColor ColorRGB `property:"rw"`
	EndColor   ColorRGB `property:"rw"`
}

func NewGradientVisualizer(length int, start, end ColorRGB) *GradientVisualizer {
	return &GradientVisualizer{
		BaseVisualizer: *NewBaseVisualizer("Gradient"),
		length:         length,
		StartColor:     start,
		EndColor:       end,
	}
}

func (v *GradientVisualizer) Start() {
	v.SendColors()
}

func (v *GradientVisualizer) OnPropertyChanged(string) {
	v.SendColors()
}

func (v *GradientVisualizer) SendColors() {
	d := make([]Led, v.length)
	length := float64(v.length)
	dr := (v.EndColor.Red - v.StartColor.Red) / (length - 1)
	dg := (v.EndColor.Green - v.StartColor.Green) / (length - 1)
	db := (v.EndColor.Blue - v.StartColor.Blue) / (length - 1)
	led := Led{v.StartColor.Red, v.StartColor.Green, v.StartColor.Blue, 0}
	for i := range d {
		d[i] = led
		led.Red += dr
		led.Green += dg
		led.Blue += db
	}
	v.SendData(d)
}
