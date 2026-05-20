package commands

import (
	"reflect"
	"testing"

	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/pkg/storage"
	"github.com/spf13/cobra"
)

func TestGetCompletionSuggestsKnownKeys(t *testing.T) {
	cfg := storage.NewMockConfig()
	cfg.Set("api.port", "8080", storage.UserConfig)
	cfg.Set("api.endpoint", "https://example.com", storage.UserConfig)
	cfg.Set("model", "foo", storage.UserConfig)
	cmd := getCommand{Context: &common.Context{Config: cfg}}

	got, directive := cmd.completeKey(nil, nil, "api.")
	if directive != cobra.ShellCompDirectiveDefault {
		t.Fatalf("expected %v directive, got: %v", cobra.ShellCompDirectiveDefault, directive)
	}

	want := []string{"api.endpoint", "api.port"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestGetCompletionStopsAfterFirstArg(t *testing.T) {
	cfg := storage.NewMockConfig()
	cfg.Set("api.endpoint", "https://example.com", storage.UserConfig)
	cmd := getCommand{Context: &common.Context{Config: cfg}}

	got, directive := cmd.completeKey(nil, []string{"api.endpoint"}, "")
	if directive != cobra.ShellCompDirectiveError {
		t.Fatalf("expected %v directive, got: %v", cobra.ShellCompDirectiveError, directive)
	}
	if len(got) != 0 {
		t.Fatalf("expected no completion after first arg, got %v", got)
	}
}

func TestUnsetCompletionSuggestsKnownKeys(t *testing.T) {
	cfg := storage.NewMockConfig()
	cfg.Set("api.port", "8080", storage.UserConfig)
	cfg.Set("api.endpoint", "https://example.com", storage.UserConfig)
	cfg.Set("model", "foo", storage.UserConfig)
	cmd := unsetCommand{Context: &common.Context{Config: cfg}}

	got, directive := cmd.completeKey(nil, nil, "model")
	if directive != cobra.ShellCompDirectiveDefault {
		t.Fatalf("expected %v directive, got: %v", cobra.ShellCompDirectiveDefault, directive)
	}

	want := []string{"model"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}
