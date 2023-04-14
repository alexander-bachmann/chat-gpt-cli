package main

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"fmt"
	"go-gpt/gpt"
	"log"
	"os"
	"strings"
)

//go:embed api.key
var apiKey string

func main() {
	model := "gpt-3.5-turbo"
	gpt := gpt.New(strings.TrimSpace(apiKey), model)

	for {
		prompt := userPrompt()
		answer, err := gpt.Chat(prompt)
		if err != nil {
			log.Fatal("Chat error:", err)
		}
		fmt.Println(answer)
		// prettyPrint(gpt.Messages)
	}
}

func userPrompt() string {
	fmt.Print("> ")
	prompt, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		log.Fatal("Error reading prompt:", err)
	}
	return strings.TrimSpace(prompt)
}

func prettyPrint(v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatal("Error marshalling JSON:", err)
	}
	fmt.Println(string(b))
}
