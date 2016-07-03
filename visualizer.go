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
	Name     string      `json:"name"`
	Value    interface{} `json:"value"`
	TypeName string      `json:"type"`
	Min      *float64    `json:"min,omitempty"`
	Max      *float64    `json:"max,omitempty"`
	Enum     []string    `json:"enum,omitempty"`
	object   reflect.Value
}

func GetVisualizerProperties(v Visualizer) []VisualizerProperty {
	value := reflect.Indirect(reflect.ValueOf(v))

	// Scan fields
	properties := []VisualizerProperty{}
	for i := 0; i < value.Type().NumField(); i++ {
		fieldType := value.Type().Field(i)
		if property := fieldType.Tag.Get("property"); property != "" {
			fieldValue := value.Field(i)

			if fieldValue.CanInterface() {
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

				var enum []string
				if enumProperty := fieldType.Tag.Get("enum"); enumProperty != "" {
					enum = strings.Split(enumProperty, ",")
				}

				properties = append(properties, VisualizerProperty{
					Name:     fieldType.Name,
					Value:    fieldValue.Interface(),
					TypeName: fieldType.Type.Name(),
					Min:      min,
					Max:      max,
					Enum:     enum,
					object:   fieldValue,
				})
			}
		}
	}
	return properties
}

func SetVisualizerProperties(v Visualizer, data map[string]interface{}) error {
	visualizerProperties := GetVisualizerProperties(v)
PropertyLoop:
	for property, value := range data {
		for _, visualizerProperty := range visualizerProperties {
			if visualizerProperty.Name == property {
				switch visualizerProperty.object.Type().Kind() {
				case reflect.Float32, reflect.Float64:
					visualizerProperty.object.SetFloat(value.(float64))
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					visualizerProperty.object.SetInt(int64(value.(float64)))
				case reflect.String:
					visualizerProperty.object.SetString(value.(string))
				case reflect.Struct:
					structValue, ok := value.(map[string]interface{})
					if !ok {
						return fmt.Errorf("invalid data for %q, received %+v", property, value)
					}

					for k, v := range structValue {
						structField := visualizerProperty.object.FieldByName(k)
						if !structField.IsValid() || !structField.CanSet() {
							return fmt.Errorf("invalid data for %q: unknown field %q", property, k)
						}
						valueToSet := reflect.ValueOf(v)
						if !valueToSet.Type().AssignableTo(structField.Type()) {
							return fmt.Errorf("invalid data for %q: bad data for field %q: received %q, wanted %q", property, k, valueToSet.Type().Name(), structField.Type().Name())
						}
						structField.Set(reflect.ValueOf(v))
					}
				default:
					return fmt.Errorf("unhandled property type: %q", visualizerProperty.object.Type().Kind().String())
				}
				v.OnPropertyChanged(property)
				continue PropertyLoop
			}
		}
		return fmt.Errorf("unknown property: %q", property)
	}
	return nil
}
