package main

import ()

type ApiVisualizer struct {
	Name       string                        `json:"name"`
	ID         uint64                        `json:"id"`
	Properties map[string]VisualizerProperty `json:"properties"`
}
