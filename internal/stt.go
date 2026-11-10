package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

type STTClient struct {
	httpClient *http.Client
}

func NewSTTClient() *STTClient {
	return &STTClient{
		httpClient: &http.Client{Timeout: 45 * time.Second},
	}
}

type STTResponse struct {
	Text  string `json:"text"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (s *STTClient) Transcribe(wavBytes []byte, settings Settings) (string, error) {
	if len(wavBytes) < 100 {
		return "", fmt.Errorf("audio snippet too short")
	}

	var endpoint string
	var apiKey string
	var model string

	switch settings.SttProvider {
	case "groq":
		endpoint = "https://api.groq.com/openai/v1/audio/transcriptions"
		apiKey = settings.SttKey
		if apiKey == "" {
			apiKey = settings.ApiKeys["groq"]
		}
		model = settings.SttModel
		if model == "" {
			model = "whisper-large-v3-turbo"
		}
	case "openai":
		endpoint = "https://api.openai.com/v1/audio/transcriptions"
		apiKey = settings.SttKey
		if apiKey == "" {
			apiKey = settings.ApiKeys["openai"]
		}
		model = settings.SttModel
		if model == "" {
			model = "whisper-1"
		}
	case "custom":
		base := strings.TrimRight(settings.SttBaseUrl, "/")
		if base == "" {
			base = strings.TrimRight(settings.BaseUrl, "/")
		}
		if base == "" {
			return "", fmt.Errorf("custom STT base URL is not configured")
		}
		if !strings.HasSuffix(base, "/audio/transcriptions") {
			endpoint = base + "/audio/transcriptions"
		} else {
			endpoint = base
		}
		apiKey = settings.SttKey
		if apiKey == "" {
			apiKey = settings.ApiKeys["custom"]
		}
		model = settings.SttModel
		if model == "" {
			model = "whisper-1"
		}
	default:
		// Fallback to Groq or OpenAI
		if settings.ApiKeys["groq"] != "" {
			endpoint = "https://api.groq.com/openai/v1/audio/transcriptions"
			apiKey = settings.ApiKeys["groq"]
			model = "whisper-large-v3-turbo"
		} else if settings.ApiKeys["openai"] != "" {
			endpoint = "https://api.openai.com/v1/audio/transcriptions"
			apiKey = settings.ApiKeys["openai"]
			model = "whisper-1"
		} else {
			return "", fmt.Errorf("no STT API key configured (please set Groq or OpenAI API key in Settings -> Audio)")
		}
	}

	if apiKey == "" && settings.SttProvider != "custom" {
		return "", fmt.Errorf("STT API key is empty")
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "speech.wav")
	if err != nil {
		return "", fmt.Errorf("failed to create multipart form file: %w", err)
	}
	if _, err := io.Copy(part, bytes.NewReader(wavBytes)); err != nil {
		return "", fmt.Errorf("failed to write audio to form file: %w", err)
	}

	if err := writer.WriteField("model", model); err != nil {
		return "", err
	}

	lang := settings.SttLanguage
	if lang != "" && lang != "auto" {
		_ = writer.WriteField("language", lang)
	}

	_ = writer.WriteField("response_format", "json")

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close multipart writer: %w", err)
	}

	req, err := http.NewRequest("POST", endpoint, body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("STT request error: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read STT response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("STT API error (status %d): %s", resp.StatusCode, string(respBytes))
	}

	var res STTResponse
	if err := json.Unmarshal(respBytes, &res); err != nil {
		return "", fmt.Errorf("failed to parse STT response: %w", err)
	}

	if res.Error != nil && res.Error.Message != "" {
		return "", fmt.Errorf("STT provider error: %s", res.Error.Message)
	}

	return strings.TrimSpace(res.Text), nil
}
