package main

import (
	"fmt"
	"reflect"
)

var visualizerIdCounter uint64 = 0

type Visualizer interface {
	OutputChan() <-chan []Led
	Start()

	// Methods for API
	Name() string
	ID() uint64

	OnPropertyChanged(string)
}

type VisualizerProperty struct {
	Value    interface{} `json:"value"`
	TypeName string      `json:"type"`
	object   reflect.Value
}

func GetVisualizerProperties(v Visualizer) map[string]VisualizerProperty {
	value := reflect.Indirect(reflect.ValueOf(v))

	// Scan fields
	properties := make(map[string]VisualizerProperty)
	for i := 0; i < value.Type().NumField(); i++ {
		fieldType := value.Type().Field(i)
		if property := fieldType.Tag.Get("property"); property != "" {
			fieldValue := value.Field(i)
			if fieldValue.CanInterface() {
				properties[fieldType.Name] = VisualizerProperty{
					Value:    fieldValue.Interface(),
					TypeName: fieldType.Type.Name(),
					object:   fieldValue,
				}
			}
		}
	}
	return properties
}

func SetVisualizerProperties(v Visualizer, data map[string]interface{}) error {
	properties := GetVisualizerProperties(v)
	for property, value := range data {
		visualizerProperty, ok := properties[property]
		if !ok {
			return fmt.Errorf("unknown property: %q", property)
		}

		switch visualizerProperty.object.Type().Kind() {
		case reflect.Float32, reflect.Float64:
			visualizerProperty.object.SetFloat(value.(float64))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			visualizerProperty.object.SetInt(int64(value.(float64)))
		default:
			return fmt.Errorf("unhandled property type: %q", visualizerProperty.object.Type().Kind().String())
		}
		v.OnPropertyChanged(property)
	}
	return nil
}
