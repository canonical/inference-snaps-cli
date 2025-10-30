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
		return // Do not add chat command if CHAT env var is not set
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

	// TODO
	//

	// var serverPort int
	// if v, err := config.Get("http.port"); err != nil {
	// 	return fmt.Errorf("could not get server port from config: %v", err)
	// } else {
	// 	var ok bool
	// 	serverPort, ok = v["http.port"].(int)
	// 	if !ok {
	// 		return fmt.Errorf("server port in config is not an integer")
	// 	}
	// }
	serverPortMap, err := config.Get("http.port")
	if err != nil {
		return fmt.Errorf("could not get server port from config: %v", err)
	}
	serverPort := serverPortMap["http.port"]

	apiBasePath, found := os.LookupEnv("OPENAI_BASE_PATH")
	if !found {
		return fmt.Errorf("OPENAI_BASE_PATH environment variable is not set")
	}

	// export OPENAI_BASE_URL="http://localhost:$port/$api_base_path"
	baseURL := fmt.Sprintf("http://localhost:%v/%s", serverPort, apiBasePath)

	cmd := exec.Command(chatPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = []string{"OPENAI_BASE_URL=" + baseURL}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running chat command: %v", err)
	}
	return nil
}
