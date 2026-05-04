package common

import (
	"fmt"
)

type ModelPropertiesStruct struct {
	Name                string `json:"name" yaml:"name"`
	Quantization        string `json:"quantization" yaml:"quantization"`
	MultimediaProjector string `json:"multimedia-projector" yaml:"multimedia-projector"`
}

func ModelProperties(ctx *Context) (ModelPropertiesStruct, error) {
	settings, err := EngineComponentSettings(ctx)
	if err != nil {
		return ModelPropertiesStruct{}, fmt.Errorf("loading engine environment: %v", err)
	}
	return modelProperties(settings)
}

func modelProperties(settingsCollection []ComponentSettings) (ModelPropertiesStruct, error) {
	properties := ModelPropertiesStruct{}
	for _, settings := range settingsCollection {
		if val, ok := settings.Properties["model-name"]; ok {
			properties.Name = val
		}

		if val, ok := settings.Properties["model-quantization"]; ok {
			properties.Quantization = val
		}

		if val, ok := settings.Properties["mmproj-name"]; ok {
			properties.MultimediaProjector = val
		}
	}

	return properties, nil
}
