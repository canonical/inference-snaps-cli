# Mock OpenAI API Server

A minimal Python HTTP server that mimics a subset of the [OpenAI REST API](https://platform.openai.com/docs/api-reference) for local testing. No third-party dependencies are required — it runs on the Python 3 standard library.

## Endpoints

Both `/v1` and `/v3` prefixes are served and return identical responses.

| Method | Path | Description |
|--------|------|-------------|
| `GET`  | `/{v1,v3}/models` | Returns a list of mock models |
| `POST` | `/{v1,v3}/chat/completions` | Returns a static chat completion response |

## Usage

```bash
python3 server.py [--host HOST] [--port PORT]
```

### Arguments

| Argument | Default | Description |
|----------|---------|-------------|
| `--host` | `127.0.0.1` | Network interface to bind to |
| `--port` | `8080` | TCP port to listen on |

### Examples

```bash
# Listen on localhost port 8080 (default)
python3 server.py

# Listen on all interfaces, port 11434
python3 server.py --host 0.0.0.0 --port 11434

# Listen on a specific interface
python3 server.py --host 192.168.1.10 --port 8080
```

## Sample Responses

### GET /v1/models

```json
{
  "object": "list",
  "data": [
    {
      "id": "mock-model",
      "object": "model",
      "created": 1712345678,
      "owned_by": "mock"
    },
    {
      "id": "mock-model-small",
      "object": "model",
      "created": 1712345678,
      "owned_by": "mock"
    }
  ]
}
```

### POST /v1/chat/completions

Request body is read and discarded; any valid JSON body is accepted.

```json
{
  "id": "chatcmpl-mock-0000000000000001",
  "object": "chat.completion",
  "created": 1712345678,
  "model": "mock-model",
  "choices": [
    {
      "index": 0,
      "message": {
        "role": "assistant",
        "content": "Hello! I am a mock assistant. How can I help you today?"
      },
      "finish_reason": "stop",
      "logprobs": null
    }
  ],
  "usage": {
    "prompt_tokens": 10,
    "completion_tokens": 12,
    "total_tokens": 22
  },
  "system_fingerprint": null
}
```

## Testing with curl

```bash
# List models
curl http://127.0.0.1:8080/v1/models | python3 -m json.tool

# Chat completion
curl http://127.0.0.1:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{"model":"mock-model","messages":[{"role":"user","content":"Hello"}]}' \
  | python3 -m json.tool
```

