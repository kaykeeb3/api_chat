
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	openAIEndpoint = "https://api.openai.com/v1/engines/davinci-codex/completions"
	openAIKey      = "YOUR_OPENAI_API_KEY"
)

type OpenAIRequest struct {
	Prompt     string  `json:"prompt"`
	MaxTokens  int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
}

type OpenAIResponse struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

func handleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	var requestData struct {
		Message string `json:"message"`
	}

	err = json.Unmarshal(body, &requestData)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	openAIRequest := OpenAIRequest{
		Prompt:     requestData.Message,
		MaxTokens:  50,
		Temperature: 0.8,
	}

	openAIRequestBody, err := json.Marshal(openAIRequest)
	if err != nil {
		http.Error(w, "Failed to prepare request to OpenAI", http.StatusInternalServerError)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, openAIEndpoint, bytes.NewBuffer(openAIRequestBody))
	if err != nil {
		http.Error(w, "Failed to create request to OpenAI", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", "Bearer "+openAIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to send request to OpenAI", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response from OpenAI", http.StatusInternalServerError)
		return
	}

	var openAIResponse OpenAIResponse
	err = json.Unmarshal(responseBody, &openAIResponse)
	if err != nil {
		http.Error(w, "Failed to parse response from OpenAI", http.StatusInternalServerError)
		return
	}

	response := openAIResponse.Choices[0].Text

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"response": "%s"}`, response)
}

func main() {
	http.HandleFunc("/chat", handleChat)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
