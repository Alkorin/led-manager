package main

import (
	"math"
)

type ColorRGB struct {
	Red   float64
	Green float64
	Blue  float64
}

func hueToRGB(h float64) (float64, float64, float64) {
	_, h = math.Modf(h)
	h *= 3.0
	if h < 1.0 {
		return 1 - h, h, 0
	} else if h < 2.0 {
		return 0, 1 - (h - 1), (h - 1)
	} else {
		return (h - 2), 0, 1 - (h - 2)
	}
}
