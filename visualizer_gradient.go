package main

import (
	"math"

	"github.com/bthomson/go-color"
)

type GradientVisualizer struct {
	BaseVisualizer

	length        int
	Interpolation string   `property:"rw" enum:"RGB,HSV+,HSV-"`
	StartColor    ColorRGB `property:"rw"`
	EndColor      ColorRGB `property:"rw"`
}

func NewGradientVisualizer(length int, start, end ColorRGB) *GradientVisualizer {
	return &GradientVisualizer{
		BaseVisualizer: *NewBaseVisualizer("Gradient"),
		length:         length,
		StartColor:     start,
		EndColor:       end,
		Interpolation:  "RGB",
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
	switch v.Interpolation {
	case "RGB":
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
	case "HSV+", "HSV-":
		startHSL := color.RGB{v.StartColor.Red, v.StartColor.Green, v.StartColor.Blue}.ToHSL()
		endHSL := color.RGB{v.EndColor.Red, v.EndColor.Green, v.EndColor.Blue}.ToHSL()

		if v.Interpolation == "HSV+" && endHSL.H < startHSL.H {
			endHSL.H += 1
		} else if v.Interpolation == "HSV-" && endHSL.H > startHSL.H {
			endHSL.H -= 1
		}

		dh := (endHSL.H - startHSL.H) / (length - 1)
		ds := (endHSL.S - startHSL.S) / (length - 1)
		dl := (endHSL.L - startHSL.L) / (length - 1)
		for i := range d {
			rgb := startHSL.ToRGB()

			d[i].Red = rgb.R
			d[i].Green = rgb.G
			d[i].Blue = rgb.B

			_, startHSL.H = math.Modf(startHSL.H + dh)
			startHSL.S += ds
			startHSL.L += dl
		}
	}
	v.SendData(d)
}
