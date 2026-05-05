package commands

import (
	"sort"
	"strings"

	"github.com/canonical/inference-snaps-cli/pkg/storage"
	"github.com/spf13/cobra"
)

func completeConfigKeys(config storage.Config, toComplete string, appendEquals bool, excludedKeys map[string]struct{}) ([]string, cobra.ShellCompDirective) {
	values, err := config.GetAll()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	completions := make([]string, 0, len(values))
	for key := range values {
		if _, isExcluded := excludedKeys[key]; isExcluded {
			continue
		}

		candidate := key
		if appendEquals {
			candidate += "="
		}

		if strings.HasPrefix(candidate, toComplete) {
			completions = append(completions, candidate)
		}
	}

	sort.Strings(completions)
	return completions, cobra.ShellCompDirectiveNoFileComp
}
