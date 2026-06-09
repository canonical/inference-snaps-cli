package commands

import (
	"fmt"
	"testing"

	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/pkg/models"
	"github.com/canonical/inference-snaps-cli/pkg/storage"
)

func prepareModelsTestData() (*listModelsCommand, *outputModels, error) {
	cache := storage.NewMockCache()
	err := cache.SetActiveModel("4b-it-int4-fq-ov")
	if err != nil {
		return nil, nil, fmt.Errorf("error setting active model name: %v", err)
	}

	allModels, err := models.LoadManifests("../../../test_data/models")
	if err != nil {
		return nil, nil, fmt.Errorf("error loading models: %v", err)
	}

	ctx := &common.Context{
		ModelsDir: "../../../test_data/models",
		Cache:     cache,
		Config:    nil,
	}
	cmd := listModelsCommand{Context: ctx}

	activeModel, err := cmd.Cache.GetActiveModel()
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %v", common.LookingUpActiveModel, err)
	}

	modelsList := outputModels{
		ActiveModel: activeModel,
		Models:      allModels,
	}

	return &cmd, &modelsList, nil
}

func TestListModelsJson(t *testing.T) {
	cmd, modelsList, err := prepareModelsTestData()
	if err != nil {
		t.Fatalf("Error preparing test data: %v", err)
	}

	err = cmd.printModelsJson(*modelsList)
	if err != nil {
		t.Fatal(err)
	}
}

func TestListModelsTable(t *testing.T) {
	cmd, modelsList, err := prepareModelsTestData()
	if err != nil {
		t.Fatalf("Error preparing test data: %v", err)
	}

	err = cmd.printModelsTable(*modelsList)
	if err != nil {
		t.Fatal(err)
	}
}

func Example_printModelsJson() {
	cmd, modelsList, err := prepareModelsTestData()
	if err != nil {
		panic(fmt.Sprintf("Error preparing test data: %v", err))
	}

	// Use only the 4b-it-int4-fq-ov model to keep output concise
	var filtered []models.Manifest
	for _, m := range modelsList.Models {
		if m.ID == "4b-it-int4-fq-ov" {
			filtered = append(filtered, m)
		}
	}
	modelsList.Models = filtered

	err = cmd.printModelsJson(*modelsList)
	if err != nil {
		panic(fmt.Sprintf("Error printing models json: %v", err))
	}

	// Output:
	// {
	//   "active-model": "4b-it-int4-fq-ov",
	//   "models": [
	//     {
	//       "id": "4b-it-int4-fq-ov",
	//       "name": "4b",
	//       "description": "OpenVino 4b test model",
	//       "model-card-url": "https://example.com/model-card",
	//       "quantization": "int4-fq",
	//       "capabilities": [
	//         "text"
	//       ],
	//       "disk-size": "6G",
	//       "components": [
	//         "model-4b-it-int4-fq-ov"
	//       ],
	//       "environment": [
	//         "MODEL_NAME=gemma-4-4b-it-int4-fq-ov",
	//         "MODEL_PATH=$SNAP_COMPONENTS/model-4b-it-int4-fq-ov"
	//       ]
	//     }
	//   ]
	// }
}
