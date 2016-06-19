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
}

func NewEsp8266Renderer(size int, dest string) *Esp8266Renderer {
	// Check dest validity
	destAddr, err := net.ResolveUDPAddr("udp", dest)
	if err != nil {
		panic(err)
	}

	// Create socket
	saddr, _ := net.ResolveUDPAddr("udp", ":0")
	conn, _ := net.ListenUDP("udp", saddr)

	return &Esp8266Renderer{
		SingleRenderer: *NewSingleRenderer(size),
		buffer:         make([]byte, 4*size),
		conn:           conn,
		destAddr:       destAddr,
	}
}

func (r *Esp8266Renderer) Start() {
	for range time.Tick(10 * time.Millisecond) {
		data := r.GetData()
		for i, color := range data {
			r.buffer[4*i+0] = byte(255 * color.Green)
			r.buffer[4*i+1] = byte(255 * color.Red)
			r.buffer[4*i+2] = byte(255 * color.Blue)
			r.buffer[4*i+3] = byte(255 * color.White)
		}
		r.conn.WriteToUDP(r.buffer, r.destAddr)
	}
}
