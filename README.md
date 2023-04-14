# Simple Chat GPT CLI written in Go

## Steps:
- create `api.key`
```bash
echo "put_your_OpenAI_API_key_here" > api.key
```

- build binary
```bash
go build -o bin/chat-gpt
```

- copy binary to directory in your $PATH (makes `chat-gpt` callable from shell)
```bash
cp bin/chat-gpt ~/.local/bin/chat-gpt
```

- if you don't care about the convenience of a binary, just run `go run main.go` at the root of the project directory
