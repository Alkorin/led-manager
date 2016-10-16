package main

import (
	"net"
	"time"
)

type Esp8266Renderer struct {
	SingleRenderer

	buffer   []byte
	conn     *net.UDPConn
	destAddr *net.UDPAddr
	reverse  bool
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

	return &Esp8266Renderer{
		SingleRenderer: *NewSingleRenderer(size, "Esp8266"),
		buffer:         make([]byte, 4*size),
		conn:           conn,
		destAddr:       destAddr,
		reverse:        reverse,
	}
}

func (r *Esp8266Renderer) Start() {
	for range time.Tick(10 * time.Millisecond) {
		data := r.GetData()
		if r.reverse {
			l := len(r.buffer)
			for i, color := range data {
				r.buffer[l-(4*(i+1))+0] = byte(255 * color.Green)
				r.buffer[l-(4*(i+1))+1] = byte(255 * color.Red)
				r.buffer[l-(4*(i+1))+2] = byte(255 * color.Blue)
				r.buffer[l-(4*(i+1))+3] = byte(255 * color.White)
			}
		} else {
			for i, color := range data {
				r.buffer[4*i+0] = byte(255 * color.Green)
				r.buffer[4*i+1] = byte(255 * color.Red)
				r.buffer[4*i+2] = byte(255 * color.Blue)
				r.buffer[4*i+3] = byte(255 * color.White)
			}
		}
		r.conn.WriteToUDP(r.buffer, r.destAddr)
	}
}
