package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestOpenAICompatibleClient(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer secret" {
			t.Errorf("missing bearer token")
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"choices":[{"message":{"content":"{\"valid\":true}"}}],"usage":{"prompt_tokens":4,"completion_tokens":2}}`)
	}))
	defer server.Close()
	client := NewOpenAICompatibleClient(server.URL, "test-model", "secret")
	result, usage, err := client.CompleteJSON(context.Background(), "system", "user")
	if err != nil {
		t.Fatalf("CompleteJSON() error = %v", err)
	}
	if string(result) != `{"valid":true}` {
		t.Fatalf("result = %s", result)
	}
	if usage.PromptTokens != 4 || usage.CompletionTokens != 2 {
		t.Fatalf("usage = %+v", usage)
	}
}

func TestOpenAICompatibleClientRequiresConfiguration(t *testing.T) {
	client := NewOpenAICompatibleClient("", "", "")
	if _, _, err := client.CompleteJSON(context.Background(), "", ""); err != ErrNotConfigured {
		t.Fatalf("error = %v", err)
	}
}

func TestRealOpenAICompatibleClient(t *testing.T) {
	if os.Getenv("LLM_INTEGRATION") != "1" {
		t.Skip("set LLM_INTEGRATION=1 to run the real provider check")
	}
	client := NewOpenAICompatibleClient(os.Getenv("LLM_BASE_URL"), os.Getenv("LLM_MODEL"), os.Getenv("LLM_API_KEY"))
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	raw, usage, err := client.CompleteJSON(ctx, "只返回 JSON 对象，不要输出额外文本。", `返回 {"ok":true}。`)
	if err != nil {
		t.Fatalf("real LLM check failed: %v", err)
	}
	var result struct {
		OK bool `json:"ok"`
	}
	if err = json.Unmarshal(raw, &result); err != nil || !result.OK {
		t.Fatalf("real LLM returned unexpected schema")
	}
	t.Logf("real LLM schema check passed: model=%s latency_ms=%d tokens=%d", client.Model(), usage.LatencyMS, usage.PromptTokens+usage.CompletionTokens)
}
