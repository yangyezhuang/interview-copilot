package internal

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type LLMStreamCallbacks struct {
	OnThinking func(thought string)
	OnToken    func(token string)
	OnDone     func(fullText string)
	OnError    func(err error)
}

type LLMClient struct {
	httpClient *http.Client
}

func NewLLMClient() *LLMClient {
	return &LLMClient{
		httpClient: &http.Client{Timeout: 3 * time.Minute},
	}
}

func resolveBaseURL(settings Settings, provider string) string {
	if url, ok := settings.BaseUrls[provider]; ok && strings.TrimSpace(url) != "" {
		return strings.TrimRight(strings.TrimSpace(url), "/")
	}
	if settings.BaseUrl != "" && (provider == "custom" || provider == "ollama") {
		return strings.TrimRight(strings.TrimSpace(settings.BaseUrl), "/")
	}
	switch provider {
	case "deepseek":
		return "https://api.deepseek.com/v1"
	case "qwen":
		return "https://dashscope.aliyuncs.com/compatible-mode/v1"
	case "doubao":
		return "https://ark.cn-beijing.volces.com/api/v3"
	case "siliconflow":
		return "https://api.siliconflow.cn/v1"
	case "moonshot":
		return "https://api.moonshot.cn/v1"
	case "zhipu":
		return "https://open.bigmodel.cn/api/paas/v4"
	case "groq":
		return "https://api.groq.com/openai/v1"
	case "ollama":
		return "http://localhost:11434/v1"
	case "custom":
		return strings.TrimRight(settings.BaseUrl, "/")
	default: // "openai"
		return "https://api.openai.com/v1"
	}
}

func (c *LLMClient) Stream(
	ctx context.Context,
	settings Settings,
	systemPrompt string,
	userPrompt string,
	imageDataUrl string,
	cb LLMStreamCallbacks,
) {
	provider := settings.Provider
	if provider == "" {
		provider = "openai"
	}

	modelPair, ok := settings.Models[provider]
	var model string
	if ok {
		if settings.Smart {
			model = modelPair.Smart
		} else {
			model = modelPair.Fast
		}
	}
	if model == "" {
		switch provider {
		case "deepseek":
			if settings.Smart {
				model = "deepseek-reasoner"
			} else {
				model = "deepseek-chat"
			}
		case "qwen":
			if settings.Smart {
				model = "qwen-max"
			} else {
				model = "qwen-plus"
			}
		case "doubao":
			if settings.Smart {
				model = "doubao-1.5-pro-32k"
			} else {
				model = "doubao-1.5-lite-32k"
			}
		case "siliconflow":
			if settings.Smart {
				model = "deepseek-ai/DeepSeek-R1"
			} else {
				model = "deepseek-ai/DeepSeek-V3"
			}
		case "moonshot":
			if settings.Smart {
				model = "moonshot-v1-32k"
			} else {
				model = "moonshot-v1-8k"
			}
		case "zhipu":
			if settings.Smart {
				model = "glm-4-plus"
			} else {
				model = "glm-4-flash"
			}
		case "anthropic":
			if settings.Smart {
				model = "claude-3-5-sonnet-20241022"
			} else {
				model = "claude-3-5-haiku-20241022"
			}
		case "gemini":
			if settings.Smart {
				model = "gemini-1.5-pro"
			} else {
				model = "gemini-2.5-flash"
			}
		case "ollama":
			if settings.Smart {
				model = "deepseek-r1:7b"
			} else {
				model = "qwen2.5:7b"
			}
		case "groq":
			if settings.Smart {
				model = "llama-3.3-70b-versatile"
			} else {
				model = "llama-3.1-8b-instant"
			}
		default:
			if settings.Smart {
				model = "gpt-4o"
			} else {
				model = "gpt-4o-mini"
			}
		}
	}

	apiKey := settings.ApiKeys[provider]

	switch provider {
	case "anthropic":
		c.streamAnthropic(ctx, apiKey, model, systemPrompt, userPrompt, imageDataUrl, cb)
	case "gemini":
		c.streamGemini(ctx, apiKey, model, systemPrompt, userPrompt, imageDataUrl, cb)
	default:
		baseURL := resolveBaseURL(settings, provider)
		if baseURL == "" {
			cb.OnError(fmt.Errorf("provider %s Base URL is not set", provider))
			return
		}
		c.streamOpenAICompatible(ctx, baseURL, apiKey, model, systemPrompt, userPrompt, imageDataUrl, cb)
	}
}

// ----------------- OpenAI-compatible (DeepSeek, Qwen, Doubao, SiliconFlow, Ollama, etc.) -----------------

func (c *LLMClient) streamOpenAICompatible(
	ctx context.Context,
	baseURL string,
	apiKey string,
	model string,
	systemPrompt string,
	userPrompt string,
	imageDataUrl string,
	cb LLMStreamCallbacks,
) {
	endpoint := baseURL
	if !strings.HasSuffix(endpoint, "/chat/completions") {
		endpoint = strings.TrimRight(endpoint, "/") + "/chat/completions"
	}

	messages := make([]map[string]interface{}, 0)
	if systemPrompt != "" {
		messages = append(messages, map[string]interface{}{
			"role":    "system",
			"content": systemPrompt,
		})
	}

	if imageDataUrl != "" {
		messages = append(messages, map[string]interface{}{
			"role": "user",
			"content": []map[string]interface{}{
				{"type": "text", "text": userPrompt},
				{"type": "image_url", "image_url": map[string]string{"url": imageDataUrl}},
			},
		})
	} else {
		messages = append(messages, map[string]interface{}{
			"role":    "user",
			"content": userPrompt,
		})
	}

	reqBody := map[string]interface{}{
		"model":    model,
		"messages": messages,
		"stream":   true,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		cb.OnError(err)
		return
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(data))
	if err != nil {
		cb.OnError(err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		cb.OnError(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		cb.OnError(fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body)))
		return
	}

	reader := bufio.NewReader(resp.Body)
	var fullSb strings.Builder

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			cb.OnError(err)
			return
		}

		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "data:") {
			continue
		}

		payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if payload == "[DONE]" {
			break
		}

		var chunk struct {
			Choices []struct {
				Delta struct {
					Content          string `json:"content"`
					ReasoningContent string `json:"reasoning_content"`
				} `json:"delta"`
			} `json:"choices"`
		}

		if err := json.Unmarshal([]byte(payload), &chunk); err == nil && len(chunk.Choices) > 0 {
			delta := chunk.Choices[0].Delta
			if delta.ReasoningContent != "" && cb.OnThinking != nil {
				cb.OnThinking(delta.ReasoningContent)
			}
			if delta.Content != "" {
				fullSb.WriteString(delta.Content)
				if cb.OnToken != nil {
					cb.OnToken(delta.Content)
				}
			}
		}
	}

	if cb.OnDone != nil {
		cb.OnDone(fullSb.String())
	}
}

// ----------------- Anthropic Claude -----------------

func (c *LLMClient) streamAnthropic(
	ctx context.Context,
	apiKey string,
	model string,
	systemPrompt string,
	userPrompt string,
	imageDataUrl string,
	cb LLMStreamCallbacks,
) {
	endpoint := "https://api.anthropic.com/v1/messages"

	var userContents []map[string]interface{}
	if imageDataUrl != "" {
		parts := strings.SplitN(imageDataUrl, ",", 2)
		var b64Data string
		mediaType := "image/jpeg"
		if len(parts) == 2 {
			b64Data = parts[1]
			if strings.Contains(parts[0], "image/png") {
				mediaType = "image/png"
			}
		} else {
			b64Data = imageDataUrl
		}

		userContents = append(userContents, map[string]interface{}{
			"type": "image",
			"source": map[string]interface{}{
				"type":       "base64",
				"media_type": mediaType,
				"data":       b64Data,
			},
		})
	}

	userContents = append(userContents, map[string]interface{}{
		"type": "text",
		"text": userPrompt,
	})

	messages := []map[string]interface{}{
		{
			"role":    "user",
			"content": userContents,
		},
	}

	reqBody := map[string]interface{}{
		"model":      model,
		"max_tokens": 4096,
		"messages":   messages,
		"stream":     true,
	}
	if systemPrompt != "" {
		reqBody["system"] = systemPrompt
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		cb.OnError(err)
		return
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(data))
	if err != nil {
		cb.OnError(err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		cb.OnError(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		cb.OnError(fmt.Errorf("Anthropic API error (%d): %s", resp.StatusCode, string(body)))
		return
	}

	reader := bufio.NewReader(resp.Body)
	var fullSb strings.Builder

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			cb.OnError(err)
			return
		}

		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "data:") {
			continue
		}

		payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		var event struct {
			Type  string `json:"type"`
			Delta struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"delta"`
		}

		if err := json.Unmarshal([]byte(payload), &event); err == nil {
			if event.Type == "content_block_delta" && event.Delta.Text != "" {
				fullSb.WriteString(event.Delta.Text)
				if cb.OnToken != nil {
					cb.OnToken(event.Delta.Text)
				}
			}
		}
	}

	if cb.OnDone != nil {
		cb.OnDone(fullSb.String())
	}
}

// ----------------- Google Gemini -----------------

func (c *LLMClient) streamGemini(
	ctx context.Context,
	apiKey string,
	model string,
	systemPrompt string,
	userPrompt string,
	imageDataUrl string,
	cb LLMStreamCallbacks,
) {
	endpoint := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:streamGenerateContent?alt=sse&key=%s", model, apiKey)

	var parts []map[string]interface{}
	if imageDataUrl != "" {
		split := strings.SplitN(imageDataUrl, ",", 2)
		b64 := imageDataUrl
		mime := "image/jpeg"
		if len(split) == 2 {
			b64 = split[1]
			if strings.Contains(split[0], "image/png") {
				mime = "image/png"
			}
		}
		parts = append(parts, map[string]interface{}{
			"inline_data": map[string]interface{}{
				"mime_type": mime,
				"data":      b64,
			},
		})
	}
	parts = append(parts, map[string]interface{}{
		"text": userPrompt,
	})

	reqBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"role":  "user",
				"parts": parts,
			},
		},
	}
	if systemPrompt != "" {
		reqBody["systemInstruction"] = map[string]interface{}{
			"parts": []map[string]interface{}{
				{"text": systemPrompt},
			},
		}
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		cb.OnError(err)
		return
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(data))
	if err != nil {
		cb.OnError(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		cb.OnError(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		cb.OnError(fmt.Errorf("Gemini API error (%d): %s", resp.StatusCode, string(body)))
		return
	}

	reader := bufio.NewReader(resp.Body)
	var fullSb strings.Builder

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			cb.OnError(err)
			return
		}

		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "data:") {
			continue
		}

		payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		var event struct {
			Candidates []struct {
				Content struct {
					Parts []struct {
						Text string `json:"text"`
					} `json:"parts"`
				} `json:"content"`
			} `json:"candidates"`
		}

		if err := json.Unmarshal([]byte(payload), &event); err == nil && len(event.Candidates) > 0 {
			parts := event.Candidates[0].Content.Parts
			for _, p := range parts {
				if p.Text != "" {
					fullSb.WriteString(p.Text)
					if cb.OnToken != nil {
						cb.OnToken(p.Text)
					}
				}
			}
		}
	}

	if cb.OnDone != nil {
		cb.OnDone(fullSb.String())
	}
}
