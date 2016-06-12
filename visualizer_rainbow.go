package main

import (
	"time"
)

type RainbowVisualizer struct {
	BaseVisualizer

	length int
}

func NewRainbowVisualizer(length int) *RainbowVisualizer {
	return &RainbowVisualizer{
		BaseVisualizer: *NewBaseVisualizer(),
		length:         length,
	}
}

func (v *RainbowVisualizer) Start() {
	d := make([]Led, v.length)
	j := 0.0
	for range time.Tick(100 * time.Millisecond) {
		j += 0.03
		for i := 0; i < v.length; i++ {
			r, g, b := hueToRGB(j + float64(i)/128.0)
			d[i] = Led{r, g, b, 0}
		}
		v.SendData(d)
	}
}
