package main

import (
	"time"
)

type ChaserVisualizer struct {
	BaseVisualizer

	length int
	Color  ColorRGB `property:"rw"`
	Speed  float64  `property:"rw,min=0,max=1"`
	Size   int      `property:"rw"`

	way float64
}

func NewChaserVisualizer(length int, color Led) *ChaserVisualizer {
	return &ChaserVisualizer{
		BaseVisualizer: *NewBaseVisualizer("Chaser"),
		length:         length,
		Color:          ColorRGB{color.Red, color.Green, color.Blue},
		Speed:          0.1,
		Size:           1,
		way:            1.0,
	}
}

func (v *ChaserVisualizer) Start() {
	go v.Run()
}

func (v *ChaserVisualizer) Run() {
	d := make([]Led, v.length)
	j := 0.0

	for range time.Tick(10 * time.Millisecond) {
		j += v.Speed * v.way
		if int(j) >= v.length {
			v.way = -1.0
		} else if int(j) == 0 {
			v.way = 1.0
		}
		for i := range d {
			if i > int(j)-v.Size && i < int(j)+v.Size {
				d[i].Red = v.Color.Red
				d[i].Green = v.Color.Green
				d[i].Blue = v.Color.Blue
			} else {
				d[i].Red = 0
				d[i].Green = 0
				d[i].Blue = 0
			}
		}
		v.SendData(d)
	}

}
