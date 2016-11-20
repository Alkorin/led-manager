package main

import (
	"math/rand"
	"time"
)

type StepColorVisualizer struct {
	BaseVisualizer

	length int

	update chan bool

	ticker       *time.Ticker
	tickerKicked chan bool

	Speed      float64 `property:"rw,min=0.1,max=3"`
	Luminosity float64 `property:"rw"`
	Width      int     `property:"rw,min=1,max=50,step=1"`
	Direction  string  `property:"rw" enum:"L2R,R2L,Double"`
}

func NewStepColorVisualizer(length int) *StepColorVisualizer {
	return &StepColorVisualizer{
		BaseVisualizer: *NewBaseVisualizer("StepColor"),
		length:         length,
		Speed:          0.5,
		Luminosity:     1.0,
		Width:          5,
		Direction:      "L2R",
		ticker:         time.NewTicker(500 * time.Millisecond),
		tickerKicked:   make(chan bool, 1),
		update:         make(chan bool, 1),
	}
}

func (v *StepColorVisualizer) Start() {
	go v.Run()
}

func (v *StepColorVisualizer) OnPropertyChanged(propertyName string) {
	switch propertyName {
	case "Luminosity":
		select {
		case v.update <- true:
		default:
			// Already some update pending
		}
	case "Speed":
		v.ticker.Stop()
		v.ticker = time.NewTicker(time.Duration(v.Speed * float64(time.Second)))
		v.tickerKicked <- true
	}
}

func (v *StepColorVisualizer) Run() {
	// Init
	buffer := make([]Led, v.length)
	out := make([]Led, v.length)
	curPos := -1
	curDir := -1
	hue := rand.Float64()
	// Loop
	for {
		select {
		case <-v.quit:
			v.ticker.Stop()
			return
		case <-v.update:
			for i := 0; i < len(buffer); i++ {
				out[i].Red = buffer[i].Red * v.Luminosity
				out[i].Green = buffer[i].Green * v.Luminosity
				out[i].Blue = buffer[i].Blue * v.Luminosity
			}
			v.SendData(out)
		case <-v.ticker.C:
			if curPos < 0 || curPos >= len(buffer) {
				// Invalid value, choose new one
				switch v.Direction {
				case "L2R":
					curPos = 0
					curDir = 1
				case "R2L":
					curPos = len(buffer) - 1
					curDir = -1
				case "Double":
					if curDir == 1 {
						curDir = -1
						curPos = len(buffer) - 1
					} else {
						curDir = 1
						curPos = 0
					}
				}
				// Compute newHue and avoid a too similar one
				hue += (1.0 + 4.0*rand.Float64()) / 6.0
			}
			newPos := curPos + v.Width*curDir
			for i := curPos; i != newPos && i < len(buffer) && i >= 0; i += curDir {
				r, g, b := hueToRGB(hue)
				buffer[i].Red = r
				buffer[i].Green = g
				buffer[i].Blue = b
			}
			curPos = newPos
			v.update <- true
		case <-v.tickerKicked:
		}
	}
}
