package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"regexp"
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

func (s *STTClient) Transcribe(wavBytes []byte, settings Settings, wm *WhisperManager) (string, error) {
	if len(wavBytes) < 100 {
		return "", fmt.Errorf("audio snippet too short")
	}

	// 1. Local transcription using whisper.cpp
	if settings.SttProvider == "local" {
		return s.TranscribeLocal(wavBytes, settings, wm)
	}

	// 2. Cloud transcription
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
	case "siliconflow":
		endpoint = "https://api.siliconflow.cn/v1/audio/transcriptions"
		apiKey = settings.SttKey
		if apiKey == "" {
			apiKey = settings.ApiKeys["siliconflow"]
		}
		model = settings.SttModel
		if model == "" {
			model = "FunAudioLLM/SenseVoiceSmall"
		}
	case "custom":
		base := strings.TrimRight(settings.SttBaseUrl, "/")
		if base == "" {
			base = strings.TrimRight(settings.BaseUrl, "/")
		}
		if base == "" {
			return "", fmt.Errorf("自定义 STT 端点未配置 Base URL")
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
		// Fallback to Groq, SiliconFlow, or OpenAI
		if settings.ApiKeys["groq"] != "" {
			endpoint = "https://api.groq.com/openai/v1/audio/transcriptions"
			apiKey = settings.ApiKeys["groq"]
			model = "whisper-large-v3-turbo"
		} else if settings.ApiKeys["siliconflow"] != "" {
			endpoint = "https://api.siliconflow.cn/v1/audio/transcriptions"
			apiKey = settings.ApiKeys["siliconflow"]
			model = "FunAudioLLM/SenseVoiceSmall"
		} else if settings.ApiKeys["openai"] != "" {
			endpoint = "https://api.openai.com/v1/audio/transcriptions"
			apiKey = settings.ApiKeys["openai"]
			model = "whisper-1"
		} else {
			return "", fmt.Errorf("未配置语音识别 (STT) 凭证，请在设置中配置本地模型或云端 API Key")
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

// TranscribeLocal transcribes WAV audio using local whisper.cpp executable
func (s *STTClient) TranscribeLocal(wavBytes []byte, settings Settings, wm *WhisperManager) (string, error) {
	if wm == nil {
		return "", fmt.Errorf("本地模型管理器未初始化")
	}

	modelName := settings.SttLocalModel
	if modelName == "" {
		modelName = "base"
	}

	if !wm.IsModelDownloaded(modelName) {
		spec, _ := LookupWhisperModel(modelName)
		displayName := modelName
		if spec.DisplayName != "" {
			displayName = spec.DisplayName
		}
		return "", fmt.Errorf("本地 Whisper 模型 '%s' 尚未下载。请前往【设置 -> 语音识别】点击下载模型", displayName)
	}

	modelPath := wm.GetModelPath(modelName)
	whisperBin := wm.FindWhisperBinary(settings.SttWhisperPath)
	if whisperBin == "" {
		return "", fmt.Errorf("未检测到本地 whisper-cli 执行文件。请在【设置 -> 语音识别】中配置 whisper-cli 路径，或将 whisper-cli.exe 放置于程序 bin 目录中")
	}

	// Create temporary WAV input file
	tmpWav, err := os.CreateTemp("", "whisper-in-*.wav")
	if err != nil {
		return "", fmt.Errorf("创建临时音频文件失败: %w", err)
	}
	tmpWavPath := tmpWav.Name()
	defer os.Remove(tmpWavPath)

	if _, err := tmpWav.Write(wavBytes); err != nil {
		tmpWav.Close()
		return "", fmt.Errorf("写入临时音频失败: %w", err)
	}
	tmpWav.Close()

	outPrefix := tmpWavPath + "-out"
	outTxtPath := outPrefix + ".txt"
	defer os.Remove(outTxtPath)

	args := []string{
		"-m", modelPath,
		"-f", tmpWavPath,
		"--no-timestamps",
		"-otxt",
		"-of", outPrefix,
	}

	lang := settings.SttLanguage
	if lang != "" && lang != "auto" {
		args = append(args, "-l", lang)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, whisperBin, args...)
	HideWindow(cmd)

	output, cmdErr := cmd.CombinedOutput()
	if cmdErr != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("whisper-cli 转录超时 (超过 60 秒)")
		}
		// If command failed, check if txt file was still written
		if txtBytes, err := os.ReadFile(outTxtPath); err == nil && len(txtBytes) > 0 {
			return strings.TrimSpace(string(txtBytes)), nil
		}
		return "", fmt.Errorf("whisper-cli 转录失败: %w, 详细信息: %s", cmdErr, string(output))
	}

	// Read generated .txt file
	txtBytes, err := os.ReadFile(outTxtPath)
	if err == nil && len(txtBytes) > 0 {
		return strings.TrimSpace(string(txtBytes)), nil
	}

	// Fallback to parsing stdout if .txt was not created
	cleaned := cleanWhisperOutput(string(output))
	if cleaned != "" {
		return cleaned, nil
	}

	return "", nil
}

// cleanWhisperOutput strips whisper logs and timestamps from stdout fallback
func cleanWhisperOutput(raw string) string {
	lines := strings.Split(raw, "\n")
	var result []string
	timestampRegex := regexp.MustCompile(`\[\d\d:\d\d:\d\d\.\d\d\d\s*-->\s*\d\d:\d\d:\d\d\.\d\d\d\]\s*`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "whisper_") || strings.HasPrefix(line, "system_info:") || strings.HasPrefix(line, "main:") {
			continue
		}
		line = timestampRegex.ReplaceAllString(line, "")
		line = strings.TrimSpace(line)
		if line != "" {
			result = append(result, line)
		}
	}
	return strings.Join(result, " ")
}
