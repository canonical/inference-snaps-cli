package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func addChatCommand() {
	chatPath := os.Getenv("CHAT")
	if chatPath == "" {
		return // Add chat command only if CHAT is set
	}
	cmd := &cobra.Command{
		Use:               "chat",
		Short:             "Start the chat CLI",
		Long:              "Chat with the server via its OpenAI API.\nThis CLI supports text-based prompting only.",
		GroupID:           "basics",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              chat,
	}
	rootCmd.AddCommand(cmd)
}

func chat(_ *cobra.Command, args []string) error {
	// Run the app at path set in CHAT environment variable
	chatPath, found := os.LookupEnv("CHAT")
	if !found {
		return fmt.Errorf("CHAT environment variable is not set")
	}

	apiUrls, err := serverApiUrls()
	if err != nil {
		return fmt.Errorf("error getting server api urls: %v", err)
	}

	os.Setenv("OPENAI_BASE_URL", apiUrls["openai"])

	cmd := exec.Command(chatPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running chat command: %v", err)
	}
	return nil
}
