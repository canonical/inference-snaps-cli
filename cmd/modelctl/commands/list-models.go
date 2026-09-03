package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/canonical/inference-snaps-cli/v2/cmd/modelctl/common"
	"github.com/canonical/inference-snaps-cli/v2/pkg/engines"
	"github.com/canonical/inference-snaps-cli/v2/pkg/models"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/cobra"
)

type listModelsCommand struct {
	*common.Context

	// flags
	format string
	all    bool
}

type outputModels struct {
	ActiveModel string                                     `json:"active-model"`
	Models      []common.ModelDetailsWithCompatibleEngines `json:"models"`
}

func ListModels(ctx *common.Context) *cobra.Command {
	var cmd listModelsCommand
	cmd.Context = ctx

	cobraCmd := &cobra.Command{
		Use:               "list-models",
		Short:             "List available models",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              cmd.run,
	}

	// flags
	supportedFormats := []string{"table", "json"}
	cobraCmd.Flags().StringVar(
		&cmd.format,
		"format",
		"table",
		fmt.Sprintf("output format (%s)", strings.Join(supportedFormats, ", ")),
	)
	cobraCmd.Flags().BoolVar(
		&cmd.all,
		"all",
		false,
		"list all models, including those not supported by the active engine",
	)
	return cobraCmd
}

func (cmd *listModelsCommand) run(_ *cobra.Command, _ []string) error {
	activeEngine, err := cmd.Cache.GetActiveEngine()
	if err != nil {
		return fmt.Errorf("%s: %w", common.LookingUpActiveEngine, err)
	}

	engineManifest, err := engines.LoadManifest(cmd.EnginesDir, activeEngine)
	if err != nil {
		return fmt.Errorf("%s: %w", common.LoadingEngineManifest, err)
	}

	var modelsList outputModels
	allModelsWithEngines, err := common.GetAllModelsWithEngines(cmd.Context)
	if err != nil {
		return fmt.Errorf("%s: %w", common.LoadingModelManifests, err)
	}
	compatibleEnginesByModel := make(map[string][]string, len(allModelsWithEngines))
	for _, model := range allModelsWithEngines {
		compatibleEnginesByModel[model.Model.Name] = model.CompatibleEngines
	}

	if cmd.all {
		for _, modelManifest := range allModelsWithEngines {
			outputModel := modelManifest.Model
			modelsList.Models = append(modelsList.Models, common.ModelDetailsWithCompatibleEngines{
				Model:             outputModel,
				CompatibleEngines: modelManifest.CompatibleEngines,
			})
		}
	} else {
		for _, model := range engineManifest.Model.Options {
			modelManifest, err := models.LoadManifest(cmd.ModelsDir, model)
			if err != nil {
				return fmt.Errorf("loading model manifest for model %s: %v", model, err)
			}
			outputModel, err := common.NewModelDetails(modelManifest)
			if err != nil {
				return fmt.Errorf("creating model details for model %s: %v", model, err)
			}
			compatibleEngines, ok := compatibleEnginesByModel[model]
			if !ok {
				return fmt.Errorf("loading model manifest for model %s: no compatible engines metadata found", model)
			}
			if !slices.Contains(compatibleEngines, activeEngine) {
				return fmt.Errorf("loading model manifest for model %s: model is not compatible with active engine %s", model, activeEngine)
			}
			modelsList.Models = append(modelsList.Models, common.ModelDetailsWithCompatibleEngines{
				Model:             outputModel,
				CompatibleEngines: compatibleEngines,
			})
		}
	}

	activeModel, err := cmd.Cache.GetActiveModel()
	if err != nil {
		return fmt.Errorf("%s: %w", common.LookingUpActiveModel, err)
	}
	modelsList.ActiveModel = activeModel

	switch cmd.format {
	case "table", "":
		if err := cmd.printModelsTable(modelsList); err != nil {
			return fmt.Errorf("table: %w", err)
		}
	case "json":
		if err := cmd.printModelsJson(modelsList); err != nil {
			return fmt.Errorf("json: %w", err)
		}
	default:
		return fmt.Errorf("unknown format %q", cmd.format)
	}

	return nil
}

func (cmd *listModelsCommand) printModelsJson(modelsList outputModels) error {
	type jsonModel struct {
		common.ModelDetails
		CompatibleEngines []string `json:"compatible-engines,omitempty"`
	}
	type outputModelsJSON struct {
		ActiveModel string      `json:"active-model"`
		Models      []jsonModel `json:"models"`
	}

	jsonOutput := outputModelsJSON{
		ActiveModel: modelsList.ActiveModel,
		Models:      make([]jsonModel, 0, len(modelsList.Models)),
	}

	for _, modelWithEngines := range modelsList.Models {
		m := jsonModel{ModelDetails: modelWithEngines.Model}
		m.CompatibleEngines = modelWithEngines.CompatibleEngines
		jsonOutput.Models = append(jsonOutput.Models, m)
	}

	jsonString, err := json.MarshalIndent(jsonOutput, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling models: %v", err)
	}
	fmt.Printf("%s\n", jsonString)
	return nil
}

func (cmd *listModelsCommand) getModelsTable(modelsList outputModels) (string, error) {
	includeEnginesColumn := cmd.all
	headerRow := []string{"name", "capabilities", "disk"}
	if includeEnginesColumn {
		headerRow = append(headerRow, "engines")
	}
	tableRows := [][]string{headerRow}

	var modelNameMaxLen, modelCapabilitiesMaxLen, modelDiskMaxLen int

	for _, modelWithEngines := range modelsList.Models {
		model := modelWithEngines.Model
		name := model.Name
		// Mark active model with "*"
		if model.Name == modelsList.ActiveModel {
			name = name + "*"
		}

		capabilities := strings.Join(model.Capabilities, ", ")
		diskSize := model.DiskSize
		var engines string
		for _, engine := range modelWithEngines.CompatibleEngines {
			engines += engine + ", "
		}
		engines = strings.TrimSuffix(engines, ", ")
		// Find max name and capabilities lengths
		modelNameMaxLen = max(modelNameMaxLen, len(name), len(headerRow[0]))
		modelCapabilitiesMaxLen = max(modelCapabilitiesMaxLen, len(capabilities), len(headerRow[1]))
		modelDiskMaxLen = max(modelDiskMaxLen, len(diskSize), len(headerRow[2]))

		row := []string{name, capabilities, diskSize}
		if includeEnginesColumn {
			row = append(row, engines)
		}
		tableRows = append(tableRows, row)
	}

	tableMaxWidth := 80
	// Increase column widths to account for paddings
	modelNameMaxLen += 1
	modelCapabilitiesMaxLen += 2
	modelDiskMaxLen += 2
	modelEnginesMaxLen := 0
	if includeEnginesColumn {
		modelDiskMaxLen += 1
		// Engines column fills the remaining space
		modelEnginesMaxLen = tableMaxWidth - (modelNameMaxLen + modelCapabilitiesMaxLen + modelDiskMaxLen)
	}

	widths := tw.Mapper[int, int]{
		0: modelNameMaxLen,         // Model name
		1: modelCapabilitiesMaxLen, // Capabilities
		2: modelDiskMaxLen,         // Disk
	}
	headerPadding := []tw.Padding{
		{Overwrite: true, Right: " "},
		{Overwrite: true, Left: " ", Right: " "},
		{Overwrite: true, Left: " "},
	}
	rowPadding := []tw.Padding{
		{Overwrite: true, Right: " "},
		{Overwrite: true, Left: " ", Right: " "},
		{Overwrite: true, Left: " "},
	}
	if includeEnginesColumn {
		widths[3] = modelEnginesMaxLen // Engines
		headerPadding = append(headerPadding, tw.Padding{Overwrite: true, Left: " "})
		rowPadding = append(rowPadding, tw.Padding{Overwrite: true, Left: " "})
	}

	options := []tablewriter.Option{
		tablewriter.WithRenderer(renderer.NewColorized(renderer.ColorizedConfig{
			Header: renderer.Tint{
				FG: renderer.Colors{color.Bold},
			},
			Column: renderer.Tint{
				FG: renderer.Colors{color.Reset},
				BG: renderer.Colors{color.Reset},
			},
			Borders: tw.BorderNone,
			Settings: tw.Settings{
				Separators: tw.Separators{ShowHeader: tw.Off, ShowFooter: tw.Off, BetweenRows: tw.Off, BetweenColumns: tw.Off},
				Lines: tw.Lines{
					ShowTop:        tw.Off,
					ShowBottom:     tw.Off,
					ShowHeaderLine: tw.Off,
					ShowFooterLine: tw.Off,
				},
				CompactMode: tw.On,
			},
		})),
		tablewriter.WithConfig(tablewriter.Config{
			MaxWidth: tableMaxWidth,
			Widths: tw.CellWidth{
				PerColumn: widths,
			},
			Header: tw.CellConfig{
				Alignment: tw.CellAlignment{Global: tw.AlignLeft},
				Padding: tw.CellPadding{
					PerColumn: headerPadding,
				},
			},
			Row: tw.CellConfig{
				Formatting: tw.CellFormatting{AutoWrap: tw.WrapTruncate},
				Alignment:  tw.CellAlignment{Global: tw.AlignLeft},
				Padding: tw.CellPadding{
					PerColumn: rowPadding,
				},
			},
		}),
	}

	var tableOutput bytes.Buffer
	table := tablewriter.NewTable(&tableOutput, options...)
	table.Header(tableRows[0])
	err := table.Bulk(tableRows[1:])
	if err != nil {
		return "", fmt.Errorf("adding data: %v", err)
	}
	err = table.Render()
	if err != nil {
		return "", fmt.Errorf("rendering: %v", err)
	}

	allModels, err := common.GetAllModels(cmd.Context)
	if err != nil {
		return "", fmt.Errorf("getting all models: %w", err)
	}

	activeEngine, err := cmd.Cache.GetActiveEngine()
	if err != nil {
		return "", fmt.Errorf("%s: %w", common.LookingUpActiveEngine, err)
	}

	incompatibleModelsCount := len(allModels) - len(tableRows[1:])
	if incompatibleModelsCount != 0 {
		hint := common.SuggestListModels(incompatibleModelsCount, activeEngine)
		return tableOutput.String() + "\n" + hint + "\n", nil
	}
	return tableOutput.String(), nil
}

func (cmd *listModelsCommand) printModelsTable(modelsList outputModels) error {
	if len(modelsList.Models) == 0 {
		fmt.Fprintln(os.Stderr, "No models found.")
		return nil
	}

	tableOutput, err := cmd.getModelsTable(modelsList)
	if err != nil {
		return fmt.Errorf("generating table: %v", err)
	}

	fmt.Print(tableOutput)
	return nil
}
