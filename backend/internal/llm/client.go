package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var ErrNotConfigured = errors.New("LLM is not configured")

type Usage struct {
	PromptTokens     int
	CompletionTokens int
	LatencyMS        int
}

type Client interface {
	CompleteJSON(context.Context, string, string) (json.RawMessage, Usage, error)
	Model() string
}

type OpenAICompatibleClient struct {
	baseURL    string
	model      string
	apiKey     string
	httpClient *http.Client
}

func NewOpenAICompatibleClient(baseURL, model, apiKey string) *OpenAICompatibleClient {
	return &OpenAICompatibleClient{
		baseURL: strings.TrimRight(baseURL, "/"), model: model, apiKey: apiKey,
		httpClient: &http.Client{Timeout: 45 * time.Second},
	}
}

func (c *OpenAICompatibleClient) Model() string { return c.model }

func (c *OpenAICompatibleClient) CompleteJSON(ctx context.Context, systemPrompt, userPrompt string) (json.RawMessage, Usage, error) {
	if c.baseURL == "" || c.model == "" || c.apiKey == "" {
		return nil, Usage{}, ErrNotConfigured
	}
	payload := map[string]any{
		"model":           c.model,
		"temperature":     0.1,
		"response_format": map[string]string{"type": "json_object"},
		"messages":        []map[string]string{{"role": "system", "content": systemPrompt}, {"role": "user", "content": userPrompt}},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, Usage{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, Usage{}, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	started := time.Now()
	resp, err := c.httpClient.Do(req)
	usage := Usage{LatencyMS: int(time.Since(started).Milliseconds())}
	if err != nil {
		return nil, usage, err
	}
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, usage, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		message := string(responseBody)
		if len(message) > 2048 {
			message = message[:2048]
		}
		return nil, usage, fmt.Errorf("LLM HTTP %d: %s", resp.StatusCode, message)
	}
	var decoded struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
		} `json:"usage"`
	}
	if err := json.Unmarshal(responseBody, &decoded); err != nil {
		return nil, usage, err
	}
	usage.PromptTokens = decoded.Usage.PromptTokens
	usage.CompletionTokens = decoded.Usage.CompletionTokens
	if len(decoded.Choices) == 0 {
		return nil, usage, errors.New("LLM response has no choices")
	}
	raw := json.RawMessage(decoded.Choices[0].Message.Content)
	var value any
	if !json.Valid(raw) || json.Unmarshal(raw, &value) != nil {
		return nil, usage, errors.New("LLM response content is not valid JSON")
	}
	return raw, usage, nil
}
