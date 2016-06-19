package main

import ()

type ApiVisualizer struct {
	Name       string                        `json:"name"`
	ID         uint64                        `json:"id"`
	Properties map[string]VisualizerProperty `json:"properties"`
}

func NewApiVisualizer(v Visualizer) *ApiVisualizer {
	return &ApiVisualizer{
		Name:       v.Name(),
		ID:         v.ID(),
		Properties: GetVisualizerProperties(v),
	}
}

type ApiEvent struct {
	EventType string
}

type ApiBufferEvent struct {
	ApiEvent

	Data []Led
}

func NewApiBufferEvent(data []Led) *ApiBufferEvent {
	return &ApiBufferEvent{
		ApiEvent{"bufferUpdate"},
		data,
	}
}

type ApiVisualizerPropertiesChangedEvent struct {
	ApiEvent

	VisualizerId uint64
}

func NewApiVisualizerPropertiesChangedEvent(visualizerId uint64) *ApiVisualizerPropertiesChangedEvent {
	return &ApiVisualizerPropertiesChangedEvent{
		ApiEvent{"visualizerPropertiesChanged"},
		visualizerId,
	}
}
