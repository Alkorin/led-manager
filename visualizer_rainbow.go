package main

import (
	"time"
)

type RainbowVisualizer struct {
	BaseVisualizer

	length int

	Speed      float64 `property:"rw,min=0,max=0.02"`
	Luminosity float64 `property:"rw"`
	Width      float64 `property:"rw,min=0,max=0.05"`
}

func NewRainbowVisualizer(length int) *RainbowVisualizer {
	return &RainbowVisualizer{
		BaseVisualizer: *NewBaseVisualizer("Rainbow"),
		length:         length,
		Speed:          0.005,
		Luminosity:     1.0,
		Width:          0.01,
	}
}

func (v *RainbowVisualizer) Start() {
	go v.Run()
}

func (v *RainbowVisualizer) Run() {
	d := make([]Led, v.length)
	j := 0.0

	ticker := time.NewTicker(10 * time.Millisecond)

	for {
		select {
		case <-v.quit:
			ticker.Stop()
			return
		case <-ticker.C:
			j += v.Speed
			for i := 0; i < v.length; i++ {
				r, g, b := hueToRGB(j + float64(i)*v.Width)
				d[i].Red = r * v.Luminosity
				d[i].Green = g * v.Luminosity
				d[i].Blue = b * v.Luminosity
			}
			v.SendData(d)
		}
	}
}
