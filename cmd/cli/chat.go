package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/openai/openai-go"
	openaiOption "github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/ssestream"
	"github.com/spf13/cobra"
)

var (
	verboseLogging bool
	chatModelName  string
	chatBaseUrl    string
)

func addChatCommand() {
	chatPath := os.Getenv(envChat)
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
	cmd.PersistentFlags().StringVar(&chatModelName, "model", "", "name of the model to use")
	cmd.PersistentFlags().StringVar(&chatBaseUrl, "base-url", "", "base URL of the OpenAI API server")
	rootCmd.AddCommand(cmd)
}

func chat(_ *cobra.Command, args []string) error {
	verboseLogging = true

	if chatBaseUrl == "" {
		apiUrls, err := serverApiUrls()
		if err != nil {
			return fmt.Errorf("error getting server api urls: %v", err)
		}
		chatBaseUrl = apiUrls[openAi]
	}

	reasoningModel := os.Getenv("REASONING_MODEL") == "true"

	if chatModelName == "" {
		modelService := openai.NewModelService(openaiOption.WithBaseURL(chatBaseUrl))
		modelPage, err := modelService.List(context.Background())
		if err != nil {
			log.Fatalf("Failed to list models: %v", err)
		}

		if len(modelPage.Data) == 0 {
			log.Fatalln("Server returned no models")
		} else if len(modelPage.Data) > 1 {
			names := make([]string, 0, len(modelPage.Data)) // Pre-allocate for efficiency
			for _, model := range modelPage.Data {
				names = append(names, model.ID)
			}
			return fmt.Errorf("server returned multiple models; please select one using:\n\t--model [%s]", strings.Join(names, "|"))
		}
		chatModelName = modelPage.Data[0].ID
	}

	if verboseLogging {
		fmt.Printf("Using model %v\n", chatModelName)
	}

	// OpenAI API Client
	client := openai.NewClient(openaiOption.WithBaseURL(chatBaseUrl))

	if err := checkServer(chatBaseUrl, client, chatModelName); err != nil {
		err = fmt.Errorf("%v\n\nUnable to chat. Make sure the server has started successfully.\n", err)
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("Connected to %s\n", chatBaseUrl)
	fmt.Println("Type your prompt, then ENTER to submit. CTRL-C to quit.")

	rl, err := readline.NewEx(&readline.Config{
		Prompt: color.RedString("» "),
		//HistoryFile:     "/tmp/readline.tmp",
		//AutoComplete:    completer,
		InterruptPrompt: "^C",
		//EOFPrompt:       "exit", // Does not work as expected

		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	})
	if err != nil {
		log.Fatalf("Can't init readline: %v", err)
	}
	defer rl.Close()
	//rl.CaptureExitSignal() // Should readline capture and handle the exit signal? - Can be used to interrupt the chat response stream.
	log.SetOutput(rl.Stderr())

	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are a helpful assistant."),
		},
		Model: chatModelName,
	}

	for {
		prompt, err := rl.Readline()
		if errors.Is(err, readline.ErrInterrupt) {
			if len(prompt) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}
		if prompt == "exit" {
			break
		}

		if len(prompt) > 0 {
			params = handlePrompt(client, params, reasoningModel, prompt)
		}
	}
	fmt.Println("Closing chat")

	return nil
}

func checkServer(baseURL string, client openai.Client, modelName string) error {

	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("Are you up?"),
		},
		Model:               modelName,
		MaxCompletionTokens: openai.Int(1),
		MaxTokens:           openai.Int(1), // for runtimes that don't yet support MaxCompletionTokens
	}

	stopProgress := startProgressSpinner("Connecting to " + baseURL)
	defer stopProgress()

	ctx := context.Background()
	_, err := client.Chat.Completions.New(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

func handlePrompt(client openai.Client, params openai.ChatCompletionNewParams, reasoningModel bool, prompt string) openai.ChatCompletionNewParams {
	params.Messages = append(params.Messages, openai.UserMessage(prompt))

	paramDebugString, _ := json.Marshal(params)

	if verboseLogging {
		fmt.Printf("Sending request: %s\n", paramDebugString)
	}

	stopProgress := startProgressSpinner("Waiting for a response")
	stream := client.Chat.Completions.NewStreaming(context.Background(), params)
	stopProgress()

	appendParam := processStream(stream, reasoningModel)

	// Store previous prompts for context
	if appendParam != nil {
		params.Messages = append(params.Messages, *appendParam)
	}
	fmt.Println()

	return params
}

func processStream(stream *ssestream.Stream[openai.ChatCompletionChunk], printThinking bool) *openai.ChatCompletionMessageParamUnion {

	// optionally, an accumulator helper can be used
	acc := openai.ChatCompletionAccumulator{}

	// For reasoning models we assume the first output is them thinking, because the opening <think> tag is not always present.
	thinking := printThinking

	for stream.Next() {
		chunk := stream.Current()
		acc.AddChunk(chunk)

		if _, ok := acc.JustFinishedContent(); ok {
			//fmt.Println("\nContent stream finished")
		}

		// if using tool calls
		if tool, ok := acc.JustFinishedToolCall(); ok {
			fmt.Printf("Tool call stream finished %d: %s %s", tool.Index, tool.Name, tool.Arguments)
		}

		if refusal, ok := acc.JustFinishedRefusal(); ok {
			fmt.Printf("Refusal stream finished: %s", refusal)
		}

		// Print chunks as they are received
		if len(chunk.Choices) > 0 {
			lastChunk := chunk.Choices[0].Delta.Content

			if strings.Contains(lastChunk, "<think>") {
				thinking = true
				fmt.Printf("%s", color.BlueString(lastChunk))
			} else if strings.Contains(lastChunk, "</think>") {
				thinking = false
				fmt.Printf("%s", color.BlueString(lastChunk))

			} else if thinking {
				fmt.Printf("%s", color.BlueString(lastChunk))

			} else {
				fmt.Printf("%s", lastChunk)
			}
		}
	}

	if stream.Err() != nil {
		err := fmt.Errorf("\n\nError reading response stream: %v\n", stream.Err())
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	// After the stream is finished, acc can be used like a ChatCompletion
	appendParam := acc.Choices[0].Message.ToParam()
	if acc.Choices[0].Message.Content == "" {
		return nil
	}
	return &appendParam
}

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}
