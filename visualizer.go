package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
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
	Min      *float64    `json:"min"`
	Max      *float64    `json:"max"`
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

			var min *float64
			var max *float64
			// Scan tag opts
			for _, opt := range strings.Split(property, ",") {
				if strings.HasPrefix(opt, "min=") {
					value, err := strconv.ParseFloat(strings.TrimPrefix(opt, "min="), 64)
					if err == nil {
						min = &value
					}
				} else if strings.HasPrefix(opt, "max=") {
					value, err := strconv.ParseFloat(strings.TrimPrefix(opt, "max="), 64)
					if err == nil {
						max = &value
					}
				}
			}

			if fieldValue.CanInterface() {
				properties[fieldType.Name] = VisualizerProperty{
					Value:    fieldValue.Interface(),
					TypeName: fieldType.Type.Name(),
					Min:      min,
					Max:      max,
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
