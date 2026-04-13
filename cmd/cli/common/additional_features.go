package common

import (
	"os"
	"strings"
)

type additionalFeatures struct {
	Chat bool
	Ui   bool
}

func AdditionalFeatures() additionalFeatures {
	const (
		additionalFeaturesEnv = "ADDITIONAL_FEATURES"
		featureChat           = "chat"
		featureUi             = "ui"
	)

	featuresCsv, found := os.LookupEnv(additionalFeaturesEnv)
	if !found {
		return additionalFeatures{}
	}

	var features additionalFeatures
	for feature := range strings.SplitSeq(featuresCsv, ",") {
		switch strings.TrimSpace(feature) {
		case featureChat:
			features.Chat = true
		case featureUi:
			features.Ui = true
		}
	}

	return features
}
