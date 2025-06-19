package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is an HTTP client for interacting with Ollama
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new Ollama client
func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// GenerateRequest represents a request to the Ollama generate endpoint
type GenerateRequest struct {
	Model       string              `json:"model"`
	Prompt      string              `json:"prompt"`
	System      string              `json:"system,omitempty"`
	Template    string              `json:"template,omitempty"`
	Context     []int               `json:"context,omitempty"`
	Options     map[string]interface{} `json:"options,omitempty"`
	Format      string              `json:"format,omitempty"`
	Stream      bool                `json:"stream,omitempty"`
	Raw         bool                `json:"raw,omitempty"`
}

// GenerateResponse represents a response from the Ollama generate endpoint
type GenerateResponse struct {
	Model     string  `json:"model"`
	Created   int64   `json:"created_at"`
	Response  string  `json:"response"`
	Done      bool    `json:"done"`
	Context   []int   `json:"context,omitempty"`
	TotalDuration int64 `json:"total_duration,omitempty"`
	LoadDuration int64 `json:"load_duration,omitempty"`
	PromptEvalCount int `json:"prompt_eval_count,omitempty"`
	EvalCount int `json:"eval_count,omitempty"`
	EvalDuration int64 `json:"eval_duration,omitempty"`
}

// ChatRequest represents a request to the Ollama chat endpoint
type ChatRequest struct {
	Model   string     `json:"model"`
	Messages []Message `json:"messages"`
	Stream  bool       `json:"stream,omitempty"`
	Options map[string]interface{} `json:"options,omitempty"`
}

// Message represents a message in a chat request/response
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatResponse represents a response from the Ollama chat endpoint
type ChatResponse struct {
	Model    string  `json:"model"`
	Created  int64   `json:"created_at"`
	Message  Message `json:"message"`
	Done     bool    `json:"done"`
	TotalDuration int64 `json:"total_duration,omitempty"`
	LoadDuration int64 `json:"load_duration,omitempty"`
	PromptEvalCount int `json:"prompt_eval_count,omitempty"`
	EvalCount int `json:"eval_count,omitempty"`
	EvalDuration int64 `json:"eval_duration,omitempty"`
}

// EmbeddingRequest represents a request to the Ollama embedding endpoint
type EmbeddingRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// EmbeddingResponse represents a response from the Ollama embedding endpoint
type EmbeddingResponse struct {
	Embedding []float64 `json:"embedding"`
}

// ModelInfo represents information about an Ollama model
type ModelInfo struct {
	Name        string    `json:"name"`
	ModifiedAt  time.Time `json:"modified_at"`
	Size        int64     `json:"size"`
	Digest      string    `json:"digest"`
	Details     ModelDetails `json:"details,omitempty"`
}

// ModelDetails represents detailed information about an Ollama model
type ModelDetails struct {
	Format      string `json:"format"`
	Family      string `json:"family"`
	ParameterSize string `json:"parameter_size"`
	QuantizationLevel string `json:"quantization_level"`
}

// ListModelsResponse represents a response from the Ollama list models endpoint
type ListModelsResponse struct {
	Models []ModelInfo `json:"models"`
}

// Generate sends a generate request to Ollama
func (c *Client) Generate(ctx context.Context, req GenerateRequest) (*GenerateResponse, error) {
	url := fmt.Sprintf("%s/api/generate", c.BaseURL)
	
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	httpReq.Header.Set("Content-Type", "application/json")
	
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}
	
	var result GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return &result, nil
}

// Chat sends a chat request to Ollama
func (c *Client) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	url := fmt.Sprintf("%s/api/chat", c.BaseURL)
	
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	httpReq.Header.Set("Content-Type", "application/json")
	
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}
	
	var result ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return &result, nil
}

// Embeddings sends an embedding request to Ollama
func (c *Client) Embeddings(ctx context.Context, req EmbeddingRequest) (*EmbeddingResponse, error) {
	url := fmt.Sprintf("%s/api/embeddings", c.BaseURL)
	
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	httpReq.Header.Set("Content-Type", "application/json")
	
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}
	
	var result EmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return &result, nil
}

// ListModels lists available models from Ollama
func (c *Client) ListModels(ctx context.Context) (*ListModelsResponse, error) {
	url := fmt.Sprintf("%s/api/tags", c.BaseURL)
	
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}
	
	var result ListModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return &result, nil
}