package main

import (
	"time"
)

type RainbowVisualizer struct {
	BaseVisualizer

	length int

	Speed      float64 `property:"rw"`
	Luminosity float64 `property:"rw"`
}

func NewRainbowVisualizer(length int) *RainbowVisualizer {
	return &RainbowVisualizer{
		BaseVisualizer: *NewBaseVisualizer("Rainbow"),
		length:         length,
		Speed:          0.5,
		Luminosity:     1.0,
	}
}

func (v *RainbowVisualizer) Start() {
	d := make([]Led, v.length)
	j := 0.0
	for range time.Tick(10 * time.Millisecond) {
		j += v.Speed / 100
		for i := 0; i < v.length; i++ {
			r, g, b := hueToRGB(j + float64(i)/128.0)
			d[i] = Led{r * v.Luminosity, g * v.Luminosity, b * v.Luminosity, 0}
		}
		v.SendData(d)
	}
}
