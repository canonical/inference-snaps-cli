package engines

import "fmt"

func Models(enginesDir string, engine string) ([]string, error) {
	engineManifest, err := LoadManifest(enginesDir, engine)
	if err != nil {
		return nil, fmt.Errorf("loading engine manifest: %w", err)
	}

	return engineManifest.Model.Options, nil
}

func DefaultModel(enginesDir string, engine string) (string, error) {
	engineManifest, err := LoadManifest(enginesDir, engine)
	if err != nil {
		return "", fmt.Errorf("loading engine manifest: %w", err)
	}
	return engineManifest.Model.Default, nil
}
