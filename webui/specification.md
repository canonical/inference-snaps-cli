# Abstract
A web-based OpenAI chat client that connects to OpenAI-compliant REST API servers for interaction with a model.


# Specification

This must be a single-page application that can be hosted with a static web server. 
It should be responsive and work well on both desktop and mobile devices.

## Chat

### Input text
The chat interface must offer a text input field for user messages and a display area for the conversation history.

Messages should be displayed in a clear and organized manner, distinguishing between user messages and model responses.

The chat input should support basic text formatting (e.g., line breaks).

### Input images
If the model supports images in input, the chat interface should allow attaching images to messages. Other file types should not be allowed.

A small preview of the attached image should be displayed in the chat history.

The user should be able to view the attached full image by clicking on the preview.

Large images should be resized to less than 0.5 mega pixels before being sent to the model, while maintaining aspect ratio.

### Progress
The interface should show a spinner while waiting for the model's response.


### Responses
The model responses should be streamed and displayed incrementally as they are received from the API, rather than waiting for the entire response to be available.

For thinking output of reasoning models, the UI should display a distinct visual indicator to differentiate the thinking process and reasoning output from the final response.

The UI should have a button to cancel a request. Any partial response should stay in the UI. The user request that led to the cancelled response and the partial response should be excluded from the context of future interactions with the model, and marked as "Cancelled" in the UI.

### Error handling
Capture any error returned by the API but do not display the raw error message to the user. Instead, show a user-friendly error message in the chat history, prefixed with an error icon (Vanilla Framework's p-icon--error). Provide a link to view the raw error details in a modal dialog for users who want to see more information.

When an error occurs, the user should be able to retry sending the prompt. 

If the model name cannot be retrieved from the /models endpoint of the OpenAI API, display a warning message in the settings section indicating that the model name could not be retrieved.

## Configuration
The web application should consume static configurations of the following parameters:
- openAIBaseURL(string): OpenAI base URL
- capabilities(string array): Model capabilities (text, vision, audio, etc) - only text and vision should be supported for now.
- instanceName(string): used to construct the UI title, e.g. "gemma3" -> "Gemma3 Inference Snap"
- engineName(string): for display in the UI

These configurations are served by the web server at /config.

The model name should be queried from the /models endpoint of the OpenAI API.

### Config display
The UI should display the configurations as read-only in a settings section, indicated by an information icon (Vanilla Framework's p-icon--information) positioned in the right-most corner.

Do not display the UI name in the settings section as this is already shown in the header.

## UI frameworks
Use Vue.js as the JavaScript framework for building the user interface.

Use Canonical's Vanilla Framework (https://vanillaframework.io/) for styling and components.


## Serving the application
The application should be served using a server written in Go.

The Go server should serve the configurations at /config over HTTP. 
Use hardcoded values in Go for now.


## Further considerations (not in scope for this implementation)
- Hot reloading of configurations without refreshing the page.
- Conversation history
- Starting a new conversation
- Control the thinking for models that support it
- Make image attachment resize configurable
- Support documents in input
