package main

import (
	"reflect"
)

var visualizerIdCounter uint64 = 0

type Visualizer interface {
	OutputChan() <-chan []Led
	Start()

	// Methods for API
	Name() string
	ID() uint64
}

type VisualizerProperty struct {
	Value interface{}
	Type  string
}

func GetVisualizerProperties(v Visualizer) map[string]VisualizerProperty {
	value := reflect.Indirect(reflect.ValueOf(v))
	t := value.Type()

	// Scan fields
	properties := make(map[string]VisualizerProperty)
	for i := 0; i < t.NumField(); i++ {
		fieldType := t.Field(i)
		if property := fieldType.Tag.Get("property"); property != "" {
			fieldValue := value.Field(i)
			if fieldValue.CanInterface() {
				properties[fieldType.Name] = VisualizerProperty{
					Value: fieldValue.Interface(),
					Type:  fieldType.Type.Name(),
				}
			}
		}
	}
	return properties
}
