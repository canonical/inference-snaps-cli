package commands

import (
	"fmt"
	"testing"

	"github.com/canonical/inference-snaps-cli/v2/cmd/modelctl/common"
	"github.com/canonical/inference-snaps-cli/v2/pkg/engines"
	"github.com/canonical/inference-snaps-cli/v2/pkg/models"
	"github.com/canonical/inference-snaps-cli/v2/pkg/storage"
)

func TestModelUnsupportedFormatResultsInError(t *testing.T) {
	manifest, err := models.LoadManifest("../../../test_data/models", "4b-it-int4-fq-ov")
	if err != nil {
		t.Fatalf("could not load model manifest: %v", err)
	}
	engine, err := engines.LoadManifest("../../../test_data/engines", "intel-gpu")
	if err != nil {
		t.Fatalf("could not load engine manifest: %v", err)
	}

	cmd := modelCommand{
		Context: &common.Context{
			EnginesDir:  "../../../test_data/engines",
			ModelsDir:   "../../../test_data/models",
			RuntimesDir: "../../../test_data/runtimes"},
		format: "invalid-format",
	}

	machine, err := machineFixture("dummy-machine")
	if err != nil {
		t.Fatalf("could not create machine fixture: %v", err)
	}

	scoredModel, err := common.GetScoredModel(cmd.Context, *engine, *manifest, machine)

	modelDetails, err := common.NewModelDetails(&scoredModel)
	if err != nil {
		t.Fatalf("could not create model details: %v", err)
	}
	err = cmd.printModelDetails(modelDetails)

	if err == nil {
		t.Fatalf("expected unsupported format to error out, got nil error")
	}
}

func Example_modelCommand_printModelDetailsYaml() {
	cache := storage.NewMockCache()
	if err := cache.SetActiveEngine("intel-gpu"); err != nil {
		panic(fmt.Sprintf("failed to set active engine: %v", err))
	}

	ctx := &common.Context{
		ModelsDir:   "../../../test_data/models",
		EnginesDir:  "../../../test_data/engines",
		RuntimesDir: "../../../test_data/runtimes",
		Cache:       cache,
	}

	manifest, err := models.LoadManifest("../../../test_data/models", "4b-it-int4-fq-ov")
	if err != nil {
		panic(fmt.Sprintf("could not load model manifest: %v", err))
	}

	engine, err := engines.LoadManifest("../../../test_data/engines", "intel-gpu")
	if err != nil {
		panic(fmt.Sprintf("could not load engine manifest: %v", err))
	}

	machine, err := machineFixture("dummy-machine")
	if err != nil {
		panic(fmt.Sprintf("could not create machine fixture: %v", err))
	}

	scoredModel, err := common.GetScoredModel(ctx, *engine, *manifest, machine)

	if err != nil {
		panic(fmt.Sprintf("failed to load model details: %v", err))
	}

	cmd := modelCommand{Context: ctx, format: "yaml"}
	modelDetails, err := common.NewModelDetails(&scoredModel)
	if err != nil {
		panic(fmt.Sprintf("could not create model details: %v", err))
	}
	if err := cmd.printModelDetails(modelDetails); err != nil {
		panic(fmt.Sprintf("failed to print model manifest: %v", err))
	}

	// Output:
	// name: 4b-it-int4-fq-ov
	// alias: 4b-it
	// description: OpenVino 4b test model
	// model-card-url: https://example.com/model-card
	// format: OpenVINO_IR
	// quantization: int4-fq
	// capabilities:
	//     - text
	// disk-size: 6G
	// required-memory: 8.7G
	// components:
	//     - model-4b-it-int4-fq-ov
	// compatible-engines:
	//     - intel-cpu
	//     - intel-gpu
	//     - intel-npu
	// compatible: true
}

func Example_modelCommand_printModelDetailsJson() {
	cache := storage.NewMockCache()
	if err := cache.SetActiveEngine("intel-gpu"); err != nil {
		panic(fmt.Sprintf("failed to set active engine: %v", err))
	}

	ctx := &common.Context{
		ModelsDir:   "../../../test_data/models",
		EnginesDir:  "../../../test_data/engines",
		RuntimesDir: "../../../test_data/runtimes",
		Cache:       cache,
	}

	manifest, err := models.LoadManifest("../../../test_data/models", "4b-it-int4-fq-ov")
	if err != nil {
		panic(fmt.Sprintf("could not load model manifest: %v", err))
	}

	engine, err := engines.LoadManifest("../../../test_data/engines", "intel-gpu")
	if err != nil {
		panic(fmt.Sprintf("could not load engine manifest: %v", err))
	}

	machine, err := machineFixture("dummy-machine")
	if err != nil {
		panic(fmt.Sprintf("could not create machine fixture: %v", err))
	}

	scoredModel, err := common.GetScoredModel(ctx, *engine, *manifest, machine)
	if err != nil {
		panic(fmt.Sprintf("failed to load model details: %v", err))
	}
	modelDetails, err := common.NewModelDetails(&scoredModel)
	if err != nil {
		panic(fmt.Sprintf("could not create model details: %v", err))
	}
	cmd := modelCommand{Context: ctx, format: "json"}
	if err := cmd.printModelDetails(modelDetails); err != nil {
		panic(fmt.Sprintf("failed to print model manifest: %v", err))
	}

	// Output:
	// {
	//   "name": "4b-it-int4-fq-ov",
	//   "alias": "4b-it",
	//   "description": "OpenVino 4b test model",
	//   "model-card-url": "https://example.com/model-card",
	//   "format": "OpenVINO_IR",
	//   "quantization": "int4-fq",
	//   "capabilities": [
	//     "text"
	//   ],
	//   "disk-size": "6G",
	//   "required-memory": "8.7G",
	//   "components": [
	//     "model-4b-it-int4-fq-ov"
	//   ],
	//   "compatible-engines": [
	//     "intel-cpu",
	//     "intel-gpu",
	//     "intel-npu"
	//   ],
	//   "compatible": true
	// }
}

func Example_modelCommand_printModelDetailsIncomaptible() {
	cache := storage.NewMockCache()
	if err := cache.SetActiveEngine("intel-gpu"); err != nil {
		panic(fmt.Sprintf("failed to set active engine: %v", err))
	}

	ctx := &common.Context{
		ModelsDir:   "../../../test_data/models",
		EnginesDir:  "../../../test_data/engines",
		RuntimesDir: "../../../test_data/runtimes",
		Cache:       cache,
	}

	manifest, err := models.LoadManifest("../../../test_data/models", "4b-it-int4-fq-ov")
	if err != nil {
		panic(fmt.Sprintf("could not load model manifest: %v", err))
	}

	engine, err := engines.LoadManifest("../../../test_data/engines", "intel-gpu")
	if err != nil {
		panic(fmt.Sprintf("could not load engine manifest: %v", err))
	}

	machine, err := machineFixture("no-disk-available-machine")
	if err != nil {
		panic(fmt.Sprintf("could not create machine fixture: %v", err))
	}

	scoredModel, err := common.GetScoredModel(ctx, *engine, *manifest, machine)

	if err != nil {
		panic(fmt.Sprintf("failed to load model details: %v", err))
	}

	cmd := modelCommand{Context: ctx, format: "yaml"}
	modelDetails, err := common.NewModelDetails(&scoredModel)
	if err != nil {
		panic(fmt.Sprintf("could not create model details: %v", err))
	}
	if err := cmd.printModelDetails(modelDetails); err != nil {
		panic(fmt.Sprintf("failed to print model manifest: %v", err))
	}

	// Output:
	// name: 4b-it-int4-fq-ov
	// alias: 4b-it
	// description: OpenVino 4b test model
	// model-card-url: https://example.com/model-card
	// format: OpenVINO_IR
	// quantization: int4-fq
	// capabilities:
	//     - text
	// disk-size: 6G
	// required-memory: 8.7G
	// components:
	//     - model-4b-it-int4-fq-ov
	// compatible-engines:
	//     - intel-cpu
	//     - intel-gpu
	//     - intel-npu
	// compatible: false
	// compatibility-issues:
	//     - insufficient memory
	//     - insufficient disk space
}
