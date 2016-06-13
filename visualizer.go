package main

import (
	"reflect"
)

type Visualizer interface {
	OutputChan() <-chan []Led
	Start()

	// Methods for API
	Name() string
}

type VisualizerProperty struct {
	Name  string
	Value interface{}
}

func GetVisualizerProperties(v Visualizer) []VisualizerProperty {
	value := reflect.Indirect(reflect.ValueOf(v))
	t := value.Type()

	// Scan fields
	properties := []VisualizerProperty{}
	for i := 0; i < t.NumField(); i++ {
		fieldType := t.Field(i)
		if property := fieldType.Tag.Get("property"); property != "" {
			fieldValue := value.Field(i)
			if fieldValue.CanInterface() {
				properties = append(properties, VisualizerProperty{
					Name:  fieldType.Name,
					Value: fieldValue.Interface(),
				})
			}
		}
	}
	return properties
}
