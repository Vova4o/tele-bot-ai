package summary

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type OpenAISummarizer struct {
	client  *genai.Client
	prompt  string
	model   string
	enabled bool
	mu      sync.Mutex
}

func NewOpenAISummarizer(apiKey, model, prompt string) *OpenAISummarizer {
	ctx := context.Background()
	clientNew, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("genai.NewClient: %v", err)
	}

	s := &OpenAISummarizer{
		client: clientNew,
		prompt: prompt,
		model:  model,
	}

	log.Printf("openai summarizer is enabled: %v", apiKey != "")

	if apiKey != "" {
		s.enabled = true
	}

	return s
}

func (s *OpenAISummarizer) Summarize(link string, title string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.enabled {
		return "", fmt.Errorf("openai summarizer is disabled")
	}

	model := s.client.GenerativeModel(s.model)

	var temperature float32 = 0.7

	model.SetTemperature(temperature)

	request := link + " " + title+ " " +s.prompt 

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	resp, err := model.GenerateContent(ctx, genai.Text(request))
	if err != nil {
		log.Fatal(err)
	}

	// resp.Candidates[0].Content.Parts[0]

	if len(resp.Candidates) == 0 {
		return "", errors.New("no choices in openai response")
	}

	rawSummary := strings.TrimSpace(fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0]))
	if strings.HasSuffix(rawSummary, ".") {
		return rawSummary, nil
	}

	// cut all after the last ".":
	sentences := strings.Split(rawSummary, ".")

	return strings.Join(sentences[:len(sentences)-1], ".") + ".", nil
}
