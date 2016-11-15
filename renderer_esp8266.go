package main

import (
	"math"
	"net"
	"time"
)

type Esp8266Renderer struct {
	SingleRenderer

	Gamma float64 `property:"rw,min=1,max=2.2"`

	buffer   []byte
	conn     *net.UDPConn
	destAddr *net.UDPAddr
	reverse  bool
	gammaMap [256]byte
}

func NewEsp8266Renderer(size int, dest string, reverse bool) *Esp8266Renderer {
	// Check dest validity
	destAddr, err := net.ResolveUDPAddr("udp", dest)
	if err != nil {
		panic(err)
	}

	// Create socket
	saddr, _ := net.ResolveUDPAddr("udp", ":0")
	conn, _ := net.ListenUDP("udp", saddr)

	r := &Esp8266Renderer{
		SingleRenderer: *NewSingleRenderer(size, "Esp8266"),
		buffer:         make([]byte, 4*size),
		conn:           conn,
		destAddr:       destAddr,
		reverse:        reverse,
		Gamma:          2.2,
	}

	r.generateGammaMap()

	return r
}

func (r *Esp8266Renderer) generateGammaMap() {
	for i := range r.gammaMap {
		r.gammaMap[i] = byte(math.Pow(float64(i)/255.0, r.Gamma) * 255)
	}
}

func (r *Esp8266Renderer) OnPropertyChanged(propertyName string) {
	switch propertyName {
	case "Gamma":
		r.generateGammaMap()
	}
}

func (r *Esp8266Renderer) Start() {
	for range time.Tick(10 * time.Millisecond) {
		data := r.GetData()
		if r.reverse {
			l := len(r.buffer)
			for i, color := range data {
				r.buffer[l-(4*(i+1))+0] = r.gammaMap[byte(255*color.Green)]
				r.buffer[l-(4*(i+1))+1] = r.gammaMap[byte(255*color.Red)]
				r.buffer[l-(4*(i+1))+2] = r.gammaMap[byte(255*color.Blue)]
				r.buffer[l-(4*(i+1))+3] = r.gammaMap[byte(255*color.White)]
			}
		} else {
			for i, color := range data {
				r.buffer[4*i+0] = r.gammaMap[byte(255*color.Green)]
				r.buffer[4*i+1] = r.gammaMap[byte(255*color.Red)]
				r.buffer[4*i+2] = r.gammaMap[byte(255*color.Blue)]
				r.buffer[4*i+3] = r.gammaMap[byte(255*color.White)]
			}
		}
		r.conn.WriteToUDP(r.buffer, r.destAddr)
	}
}
